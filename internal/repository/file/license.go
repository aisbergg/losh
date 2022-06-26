package file

import (
	"encoding/json"
	"losh/internal/models"

	"github.com/aisbergg/go-pathlib"
	"github.com/rotisserie/eris"
)

// GetLicense returns a license by its ID.
func (fr *FileRepository) GetLicense(id string) (*models.License, error) {
	licenses, err := fr.GetAllLicenses()
	if err != nil {
		return nil, err
	}
	for _, l := range licenses {
		if l.Xid == id {
			return l, nil
		}
	}
	return nil, eris.New("license not found")
}

// GetAllLicenses returns a list of all licenses
func (fr *FileRepository) GetAllLicenses() ([]*models.License, error) {
	content, err := readFile(fr.Path)
	if err != nil {
		return nil, err
	}
	var licenses []*models.License
	err = json.Unmarshal(content, &licenses)
	if err != nil {
		return nil, eris.Wrap(err, "failed to read license file")
	}

	// remove Dgraph ID and normalize license type
	for _, l := range licenses {
		l.ID = ""
		l.Type = models.AsLicenseType(string(l.Type))
	}

	return licenses, nil
}

// SaveLicenses writes the licenses to the file.
func (fr *FileRepository) SaveLicenses(licenses []*models.License) error {
	b, err := json.Marshal(licenses)
	if err != nil {
		return eris.Wrap(err, "failed to marshal licenses")
	}

	// remove Dgraph ID
	for _, l := range licenses {
		l.ID = ""
	}

	err = saveFile(fr.Path, b)
	if err != nil {
		return eris.Wrap(err, "failed to write licenses to file")
	}
	return nil
}

// DeleteLicense implements the Repository interface.
func (fr *FileRepository) DeleteLicense(id string) error {
	licenses, err := fr.GetAllLicenses()
	if err != nil {
		return err
	}
	for i, l := range licenses {
		if l.Xid == id {
			licenses = append(licenses[:i], licenses[i+1:]...)
			break
		}
	}
	return fr.SaveLicenses(licenses)
}

// DeleteAllLicenses implements the Repository interface.
func (fr *FileRepository) DeleteAllLicenses() error {
	return saveFile(fr.Path, []byte("{}"))
}

// readFile reads a file and returns its content.
func readFile(path *pathlib.Path) ([]byte, error) {
	if exists, err := path.Exists(); err != nil || !exists {
		return nil, eris.New("file does not exist")
	}
	if isFile, err := path.IsFile(); err != nil || !isFile {
		return nil, eris.New("path is not a file")
	}

	content, err := path.ReadFile()
	if err != nil {
		return nil, eris.Wrap(err, "failed to read license file")
	}

	return content, nil
}

// saveFile reads a file and returns its content.
func saveFile(path *pathlib.Path, content []byte) error {
	if exists, _ := path.Exists(); exists {
		if isFile, _ := path.IsFile(); !isFile {
			return eris.New("path is not a file")
		}
	} else if exists, _ := path.Parent().Exists(); !exists {
		return eris.New("parent directory does not exist")
	}
	return path.WriteFile(content)
}
