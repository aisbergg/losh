package liquid

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/utils"
	"gopkg.in/osteele/liquid.v1"
)

// Engine struct
type Engine struct {
	// views folder
	directory string
	// http.FileSystem supports embedded files
	fileSystem http.FileSystem
	// views extension
	extension string
	// layout variable name that incapsulates the template
	layout string
	// determines if the engine parsed all templates
	loaded bool
	// reload on each render
	reload bool
	// debug prints the parsed templates
	debug bool
	// lock for funcmap and templates
	mutex sync.RWMutex
	// templates
	Templates map[string]*liquid.Template
}

// New returns a Handlebar render engine for Fiber
func New(directory, extension string) *Engine {
	engine := &Engine{
		directory: directory,
		extension: extension,
		layout:    "embed",
	}
	return engine
}

func NewFileSystem(fs http.FileSystem, extension string) *Engine {
	engine := &Engine{
		directory:  "/",
		fileSystem: fs,
		extension:  extension,
		layout:     "embed",
	}
	return engine
}

// Layout defines the variable name that will incapsulate the template
func (e *Engine) Layout(key string) *Engine {
	e.layout = key
	return e
}

// Reload if set to true the templates are reloading on each render,
// use it when you're in development and you don't want to restart
// the application when you edit a template file.
func (e *Engine) Reload(enabled bool) *Engine {
	e.reload = enabled
	return e
}

// Debug will print the parsed templates when Load is triggered.
func (e *Engine) Debug(enabled bool) *Engine {
	e.debug = enabled
	return e
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	// race safe
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.Templates = make(map[string]*liquid.Template)
	liquidEngine := liquid.NewEngine()

	// Loop trough each directory and register template files
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, e.extension) {
			return nil
		}
		// Get the relative file path
		// ./views/html/index.tmpl -> index.tmpl
		rel, err := filepath.Rel(e.directory, path)
		if err != nil {
			return err
		}
		// Reverse slashes '\' -> '/' and
		// partials\footer.tmpl -> partials/footer.tmpl
		name := filepath.ToSlash(rel)
		// Remove ext from name 'index.tmpl' -> 'index'
		name = strings.TrimSuffix(name, e.extension)
		// name = strings.Replace(name, e.extension, "", -1)
		// Read the file
		// #gosec G304
		buf, err := utils.ReadFile(path, e.fileSystem)
		if err != nil {
			return err
		}
		// Create new template associated with the current one
		tmpl, err := liquidEngine.ParseTemplate(buf)
		if err != nil {
			return err
		}
		e.Templates[name] = tmpl
		// Debugging
		if e.debug {
			fmt.Printf("views: parsed template: %s\n", name)
		}
		return err
	}
	// notify engine that we parsed all templates
	e.loaded = true
	if e.fileSystem != nil {
		return utils.Walk(e.fileSystem, e.directory, walkFn)
	}
	return filepath.Walk(e.directory, walkFn)
}

func getLiquidBinding(binding interface{}) liquid.Bindings {
	if binding == nil {
		return nil
	}
	if binds, ok := binding.(liquid.Bindings); ok {
		return binds
	}
	if binds, ok := binding.(map[string]interface{}); ok {
		return binds
	}
	if binds, ok := binding.(fiber.Map); ok {
		bind := make(liquid.Bindings)
		for key, value := range binds {
			bind[key] = value
		}
		return bind
	}

	return nil
}

// Render will render the template by name
func (e *Engine) Render(out io.Writer, template string, binding interface{}, layout ...string) error {
	if !e.loaded || e.reload {
		if e.reload {
			e.loaded = false
		}
		if err := e.Load(); err != nil {
			return err
		}
	}
	tmpl, ok := e.Templates[template]
	if !ok {
		return fmt.Errorf("template %s does not exist", template)
	}

	var err error
	liquidBinding := getLiquidBinding(binding)
	rendered, err := tmpl.Render(liquidBinding)
	if err != nil {
		return err
	}
	if len(layout) > 0 && layout[0] != "" {
		if liquidBinding == nil {
			liquidBinding = make(map[string]interface{}, 1)
		}
		liquidBinding[e.layout] = rendered
		template := e.Templates[layout[0]]
		if template == nil {
			return fmt.Errorf("layout %s does not exist", layout[0])
		}
		rendered, err = tmpl.Render(liquidBinding)
	}
	if _, err = out.Write(rendered); err != nil {
		return err
	}
	return nil
}
