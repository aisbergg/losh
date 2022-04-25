package liquid

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/utils"
	"github.com/osteele/liquid"
)

func newLiquidEngine(templates map[string]*liquid.Template) *liquid.Engine {
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
	// reload on each render
	reload bool
	// lock for loading up the engine
	mutex sync.RWMutex
	// templates
	Templates map[string]*liquid.Template
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
func NewFileSystem(fs http.FileSystem, extension string) *Engine {
	engine := &Engine{
		directory:  "/",
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

// Reload if set to true the templates are reloading on each render,
// use it when you're in development and you don't want to restart
// the application when you edit a template file.
func (e *Engine) Reload(enabled bool) *Engine {
	e.reload = enabled
	return e
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	// race safe
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.Templates = make(map[string]*liquid.Template)
	liquidEngine := newLiquidEngine(e.Templates)

	// Loop trough each directory and register template files
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return NewLoadError(path, err)
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
			return NewLoadError(path, errors.New("failed make relative file path"))
		}
		// reverse slashes '\' -> '/'
		name := filepath.ToSlash(rel)
		// read the file
		buf, err := utils.ReadFile(path, e.fileSystem)
		if err != nil {
			return NewLoadError(path, err)
		}
		// create new template associated with the current one
		tmpl, err := liquidEngine.ParseTemplate(buf)
		if err != nil {
			return NewLoadError(path, err)
		}
		e.Templates[name] = tmpl
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
			return liqBinding, errors.New("invalid binding")
		}
	}
	// for capture_global tag
	liqBinding["capture_global"] = make(map[string][]string)

	return liqBinding, nil
}

// getTemplate returns the template by name or name+extension.
func (e *Engine) getTemplate(name string) (*liquid.Template, error) {
	tmpl, ok := e.Templates[name]
	if ok {
		return tmpl, nil
	}
	tmpl, ok = e.Templates[name+e.extension]
	if ok {
		return tmpl, nil
	}
	return nil, fmt.Errorf("template %s does not exist", name)
}

// Render will render the template by name
func (e *Engine) Render(out io.Writer, template string, binding interface{}, layout ...string) error {
	if !e.loaded || e.reload {
		if e.reload {
			e.loaded = false
		}
		if err := e.Load(); err != nil {
			return NewRenderError(template, err)
		}
	}

	var err error
	var tmpl *liquid.Template
	if tmpl, err = e.getTemplate(template); err != nil {
		return NewRenderError(template, err)
	}

	liquidBinding, err := getLiquidBinding(binding)
	if err != nil {
		return NewRenderError(template, err)
	}
	rendered, err := tmpl.Render(liquidBinding)
	if err != nil {
		return NewRenderError(template, err)
	}
	if len(layout) > 0 && layout[0] != "" {
		if liquidBinding == nil {
			liquidBinding = make(map[string]interface{}, 1)
		}
		liquidBinding[e.layoutVar] = rendered
		if tmpl, err = e.getTemplate(layout[0]); err != nil {
			return NewRenderError(template, err)
		}
		rendered, err = tmpl.Render(liquidBinding)
		if err != nil {
			return NewRenderError(template, err)
		}
	}
	if _, err = out.Write(rendered); err != nil {
		return NewRenderError(template, err)
	}
	return nil
}
