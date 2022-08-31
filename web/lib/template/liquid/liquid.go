package liquid

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	losherrors "losh/internal/lib/errors"
	"losh/internal/lib/util/reflectutil"

	"github.com/aisbergg/go-frontmatter/pkg/frontmatter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/utils"
	"github.com/osteele/liquid"
	"gopkg.in/yaml.v3"
)

func newLiquidEngine(templates map[string]*template) *liquid.Engine {
	engine := liquid.NewEngine()
	addTags(engine, templates)
	addFilters(engine)
	return engine
}

// Engine implements fiber.Views interface to load and render Liquid templates.
type Engine struct {
	// views folder
	directory string
	// http.FileSystem supports embedded files
	fileSystem http.FileSystem
	// views extension
	extension string
	// layoutVar variable name that encapsulates the template
	layoutVar string
	// determines if the engine parsed all templates
	loaded bool
	// reloadEnabled on each render
	reloadEnabled bool
	// lock for loading up the engine
	mutex sync.RWMutex
	// templates
	Templates map[string]*template
	// determines whether the frontmatterEnabled parsing shall be used
	frontmatterEnabled bool
}

// New creates a Liquid template render engine for Fiber.
func New(directory, extension string) *Engine {
	engine := &Engine{
		directory: directory,
		extension: extension,
		layoutVar: "embed",
	}
	return engine
}

// NewFileSystem creates a Liquid template render engine for Fiber.
func NewFileSystem(fs http.FileSystem, directory, extension string) *Engine {
	engine := &Engine{
		directory:  directory,
		fileSystem: fs,
		extension:  extension,
		layoutVar:  "embed",
	}
	return engine
}

// Layout defines the variable name that will incapsulate the template
func (e *Engine) Layout(key string) *Engine {
	e.layoutVar = key
	return e
}

// EnableReload if set to true the templates are reloading on each render,
// use it when you're in development and you don't want to restart
// the application when you edit a template file.
func (e *Engine) EnableReload(enabled bool) *Engine {
	e.reloadEnabled = enabled
	return e
}

// EnableFrontmatter enables the frontmatter parsing.
func (e *Engine) EnableFrontmatter(enabled bool) *Engine {
	e.frontmatterEnabled = enabled
	return e
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	// race safe
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// frontmatter
	frtmFmts := []*frontmatter.Format{
		frontmatter.NewFormat("---", "---", yaml.Unmarshal),
	}

	e.Templates = make(map[string]*template)
	liquidEngine := newLiquidEngine(e.Templates)

	// Loop trough each directory and register template files
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return NewLoadError(path, err, nil)
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, e.extension) {
			return nil
		}
		// get the relative file path, e.g: ./views/html/index.tmpl -> index.tmpl
		rel, err := filepath.Rel(e.directory, path)
		if err != nil {
			return NewLoadError(path, losherrors.NewAppError("failed make relative file path"), nil)
		}
		// reverse slashes '\' -> '/'
		name := filepath.ToSlash(rel)

		// read the file
		file, err := e.fileSystem.Open(path)
		if err != nil {
			return NewLoadError(path, err, nil)
		}
		defer file.Close()
		var body []byte
		var frtm map[string]interface{}
		if e.frontmatterEnabled {
			body, err = frontmatter.Parse(file, &frtm, frtmFmts...)
		} else {
			body, err = ioutil.ReadAll(file)
		}
		if err != nil {
			return NewLoadError(path, err, nil)
		}

		// create new template associated with the current one
		tmpl, err := liquidEngine.ParseTemplate(body)
		if err != nil {
			return NewLoadError(path, err, nil)
		}
		e.Templates[name] = newTemplate(tmpl, frtm)
		return nil
	}
	// notify engine that we parsed all templates
	e.loaded = true
	if e.fileSystem != nil {
		return utils.Walk(e.fileSystem, e.directory, walkFn)
	}
	return filepath.Walk(e.directory, walkFn)
}

func getLiquidBinding(binding interface{}) (liquid.Bindings, error) {
	var liqBinding liquid.Bindings
	if binding != nil {
		switch t := binding.(type) {
		case liquid.Bindings:
			liqBinding = t
		case map[string]interface{}:
			liqBinding = t
		case fiber.Map:
			liqBinding = make(liquid.Bindings)
			for key, value := range t {
				liqBinding[key] = value
			}
		default:
			return liqBinding, losherrors.NewAppError("invalid binding")
		}
	}
	// for capture_global tag
	liqBinding["captured_global"] = make(map[string][]string)

	return liqBinding, nil
}

// getTemplate returns the template by name or name+extension.
func (e *Engine) getTemplate(name string) (*template, error) {
	tmpl, ok := e.Templates[name]
	if ok {
		return tmpl, nil
	}
	tmpl, ok = e.Templates[name+e.extension]
	if ok {
		return tmpl, nil
	}
	return nil, losherrors.NewAppError("template %s does not exist", name)
}

// Render will render the template by name
func (e *Engine) Render(out io.Writer, name string, binding interface{}, layout ...string) error {
	if !e.loaded || e.reloadEnabled {
		if e.reloadEnabled {
			e.loaded = false
		}
		if err := e.Load(); err != nil {
			return NewRenderError(name, err, binding)
		}
	}

	var err error
	var tpl *template
	if tpl, err = e.getTemplate(name); err != nil {
		return NewRenderError(name, err, binding)
	}

	liquidBinding, err := getLiquidBinding(binding)
	if err != nil {
		return NewRenderError(name, err, binding)
	}

	// get stack of templates (layouts + template)
	tpls := make([]*template, 0, 5)
	var cbdFrtm map[string]interface{}
	i := 0
	for {
		i++
		tpls = append(tpls, tpl)
		frtm := tpl.frtm
		if frtm != nil {
			lyti, ok := frtm["layout"]
			if cbdFrtm == nil {
				cbdFrtm = frtm
			} else {
				mergeMap(cbdFrtm, frtm)
			}
			if ok {
				if lyt, ok := lyti.(string); ok {
					if tpl, err = e.getTemplate(lyt); err != nil {
						return NewRenderError(lyt, err, binding)
					}
					continue
				}
			}
		} else if i == 1 && len(layout) > 0 && layout[0] != "" {
			if tpl, err = e.getTemplate(layout[0]); err != nil {
				return NewRenderError(layout[0], err, binding)
			}
			continue
		}
		break
	}

	// render stack of templates
	if liquidBinding == nil {
		liquidBinding = make(map[string]interface{})
	}
	mergeMap(cbdFrtm, liquidBinding)
	liquidBinding[e.layoutVar] = ""
	var rendered []byte
	for _, tpl := range tpls {
		rendered, err = tpl.Render(liquidBinding, false)
		if err != nil {
			return NewRenderError(name, err, binding)
		}
		liquidBinding[e.layoutVar] = rendered
	}

	// write out the rendered template
	if _, err = out.Write(rendered); err != nil {
		return NewRenderError(name, err, binding)
	}

	return nil
}

// template is a Liquid template with frontmatter.
type template struct {
	*liquid.Template
	// frtm is the frontmatter data
	frtm map[string]interface{}
}

// newTemplate creates a new template from the given reader.
func newTemplate(lqdTpl *liquid.Template, frtm map[string]interface{}) *template {
	return &template{
		Template: lqdTpl,
		frtm:     frtm,
	}
}

func (t *template) Render(vars liquid.Bindings, inclFrtm bool) ([]byte, liquid.SourceError) {
	if inclFrtm && t.frtm != nil {
		mergeMap(t.frtm, &vars)
	}
	return t.Template.Render(vars)
}

func mergeMap(src, dst interface{}) {
	srcVal := reflectutil.Indirect(reflect.ValueOf(src))
	dstVal := reflectutil.Indirect(reflect.ValueOf(dst))
	if srcVal.Kind() != reflect.Map || dstVal.Kind() != reflect.Map {
		return
	}
	if !srcVal.IsValid() || !dstVal.IsValid() {
		return
	}

	for _, key := range srcVal.MapKeys() {
		if dstVal.MapIndex(key).IsValid() {
			mergeMap(srcVal.MapIndex(key).Interface(), dstVal.MapIndex(key).Interface())
		} else {
			dstVal.SetMapIndex(key, srcVal.MapIndex(key))
		}
	}
}
