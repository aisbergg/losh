// Package licensefile provides a license repository that loads/stores licenses
// in a file. It is mostly used for testing.
package licensefile

import (
	"context"
	"encoding/json"

	"losh/internal/core/product/models"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
)

// FileRepository is the license repository that uses a file as the backend.
type FileRepository struct {
	Path pathlib.Path
}

// NewFileRepository creates a new FileProvider.
func NewFileRepository(path pathlib.Path) *FileRepository {
	return &FileRepository{
		Path: path,
	}
}

// GetLicense returns a license by its ID.
func (fr *FileRepository) GetLicense(ctx context.Context, id string) (*models.License, error) {
	licenses, err := fr.GetAllLicenses(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range licenses {
		if l.Xid == &id {
			return l, nil
		}
	}
	return nil, errors.New("license not found")
}

// GetAllLicenses returns a list of all licenses
func (fr *FileRepository) GetAllLicenses(ctx context.Context) ([]*models.License, error) {
	content, err := readFile(fr.Path)
	if err != nil {
		return nil, err
	}
	var licenses []*models.License
	err = json.Unmarshal(content, &licenses)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read license file")
	}

	// remove Dgraph ID and normalize license type
	for _, l := range licenses {
		l.ID = nil
		lt := models.AsLicenseType(string(*l.Type))
		l.Type = &lt
	}

	return licenses, nil
}

// SaveLicenses writes the licenses to the file.
func (fr *FileRepository) SaveLicenses(ctx context.Context, licenses []*models.License) error {
	b, err := json.Marshal(licenses)
	if err != nil {
		return errors.Wrap(err, "failed to marshal licenses")
	}

	// remove Dgraph ID
	for _, l := range licenses {
		l.ID = nil
	}

	err = saveFile(fr.Path, b)
	if err != nil {
		return errors.Wrap(err, "failed to write licenses to file")
	}
	return nil
}

// DeleteLicense implements the Repository interface.
func (fr *FileRepository) DeleteLicense(ctx context.Context, id string) error {
	licenses, err := fr.GetAllLicenses(ctx)
	if err != nil {
		return err
	}
	for i, l := range licenses {
		if l.Xid == &id {
			licenses = append(licenses[:i], licenses[i+1:]...)
			break
		}
	}
	return fr.SaveLicenses(ctx, licenses)
}

// DeleteAllLicenses implements the Repository interface.
func (fr *FileRepository) DeleteAllLicenses(ctx context.Context) error {
	return saveFile(fr.Path, []byte("{}"))
}

// readFile reads a file and returns its content.
func readFile(path pathlib.Path) ([]byte, error) {
	if exists, err := path.Exists(); err != nil || !exists {
		return nil, errors.New("file does not exist")
	}
	if isFile, err := path.IsFile(); err != nil || !isFile {
		return nil, errors.New("path is not a file")
	}

	content, err := path.ReadFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read license file")
	}

	return content, nil
}

// saveFile reads a file and returns its content.
func saveFile(path pathlib.Path, content []byte) error {
	if exists, _ := path.Exists(); exists {
		if isFile, _ := path.IsFile(); !isFile {
			return errors.New("path is not a file")
		}
	} else if exists, _ := path.Parent().Exists(); !exists {
		return errors.New("parent directory does not exist")
	}
	return path.WriteFile(content)
}
