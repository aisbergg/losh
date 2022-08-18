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
