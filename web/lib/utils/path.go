package utils

import (
	"os"

	"github.com/chigopher/pathlib"
	"github.com/rotisserie/eris"
)

func ResolveExecRelPath(p string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return p, eris.New("failed to resolve path")
	}
	path := pathlib.NewPath(execPath).Parent().Join(p)
	if exists, err := path.Exists(); err != nil || !exists {
		return path.String(), eris.New("path does not exists")
	}
	path, err = path.ResolveAll()
	if err != nil {
		return path.String(), eris.New("failed to resolve path")
	}
	path = path.Clean()
	return path.String(), nil
}
