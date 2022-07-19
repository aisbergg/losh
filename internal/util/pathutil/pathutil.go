package pathutil

import (
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
)

// GetValidFilePath returns true if path is a valid file path.
func GetValidFilePath(pathStr string) (pathlib.Path, error) {
	path := pathlib.NewPath(pathStr)
	if exists, err := path.Exists(); err != nil || !exists {
		return path, errors.New("file does not exist")
	}
	path, err := path.ResolveAll()
	if err != nil {
		return path, errors.New("failed to resolve path")
	}
	if isFile, err := path.IsFile(); err != nil || !isFile {
		return path, errors.New("given path is not a file")
	}

	return path, nil
}
