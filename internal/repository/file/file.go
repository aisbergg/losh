package file

import "github.com/aisbergg/go-pathlib"

// FileRepository is the license repository that uses a file as the backend.
type FileRepository struct {
	Path *pathlib.Path
}

// NewFileRepository creates a new FileProvider.
func NewFileRepository(path *pathlib.Path) *FileRepository {
	return &FileRepository{
		Path: path,
	}
}
