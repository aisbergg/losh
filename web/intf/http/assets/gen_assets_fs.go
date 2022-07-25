// go:build ignore

// This is a generator for the assets.go file. The assets.go file is used to
// include static assets (templates, css, js, etc.) into the application. In
// prod mode the assets will be embedded into the binary. In dev mode the assets
// will be loaded from the compiled assets dir.

package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/aisbergg/go-pathlib/pkg/pathlib"
)

// includes is a list of elements to include in the prod embed.FS.
var includes = []string{
	"icons",
	"static",
	"templates",
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		die(fmt.Errorf("Usage: %s <output-file> <output-dir>", os.Args[0]))
	}

	var tpl *template.Template
	mode := strings.ToLower(strings.TrimSpace(args[0]))
	switch mode {
	case "dev":
		tpl = devTemplate
	case "prod":
		tpl = prodTemplate
	default:
		die(fmt.Errorf("Unknown mode: %s", mode))
	}

	dstPth := pathlib.NewPath(args[1])
	dstFle, err := dstPth.OpenFile(os.O_WRONLY | os.O_CREATE | os.O_TRUNC)
	die(err)
	defer dstFle.Close()
	err = tpl.Execute(dstFle, struct {
		Includes []string
	}{
		Includes: includes,
	})
	die(err)

	fmt.Printf("Generated %s\n", dstPth)
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var prodTemplate = template.Must(template.New("").Parse(`// DO NOT EDIT

package assets

import (
	"embed"
	"net/http"

	"github.com/spf13/afero"
)

{{ range $inc := .Includes -}}
//go:embed {{ $inc }}
{{ end -}}
var embeddedAssets embed.FS

// AssetsHTTP holds the assets of the application for the http.Filesystem.
var AssetsHTTP = http.FS(embeddedAssets)

// AssetsAfero holds the assets of the application in form of afero.Fs.
var AssetsAfero = afero.FromIOFS{FS: embeddedAssets}
`))

var devTemplate = template.Must(template.New("").Parse(`// DO NOT EDIT

package assets

import (
	"os"
	"errors"

	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"github.com/spf13/afero"
)

// AssetsHTTP holds the assets of the application for the http.Filesystem.
var AssetsHTTP *afero.HttpFs

// AssetsAfero holds the assets of the application in form of afero.Fs.
var AssetsAfero afero.Fs

func init() {
	path, err := resolvePath("assets")
	if err != nil {
		panic(err)
	}
	AssetsAfero = afero.NewBasePathFs(afero.NewOsFs(), path)
	AssetsHTTP = afero.NewHttpFs(AssetsAfero)
}

func resolvePath(p string) (string, error) {
	// first try to resolve path relative to executable path
	execPath, err := os.Executable()
	if err != nil {
		return p, errors.New("failed to resolve path")
	}
	path := pathlib.NewPosixPath(execPath).Parent().Join(p)
	if exists, err := path.Exists(); err == nil || exists {
		path, err = path.ResolveAll()
		if err == nil {
			return path.Clean().String(), nil
		}
	}

	// then try to resolve in current working directory
	path = pathlib.NewPosixPath(p)
	if exists, err := path.Exists(); err == nil || exists {
		path, err = path.ResolveAll()
		if err == nil {
			return path.Clean().String(), nil
		}
	}

	return "", errors.New("failed to resolve path")
}
`))
