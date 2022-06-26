package license

import (
	"losh/internal/models"
	"losh/internal/util/stringutil"

	"github.com/rotisserie/eris"
)

// Provider is the interface for providing SPDX licenses.
type Provider interface {
	// GetLicense returns a license by its SPDX ID.
	GetLicense(id, xid *string) (*models.License, error)

	// GetAllLicenses returns a list of all licenses.
	GetAllLicenses() ([]*models.License, error)
}

// Repository is the interface for getting and saving licenses to a repository.
type Repository interface {
	// GetLicense returns a license by its SPDX ID.
	GetLicense(id, xid *string) (*models.License, error)

	// GetAllLicenses returns a list of all licenses.
	GetAllLicenses() ([]*models.License, error)

	// SaveLicenses saves the licenses.
	SaveLicenses(licenses []*models.License) error

	// DeleteLicense deletes a license.
	DeleteLicense(id, xid *string) error

	// DeleteAllLicenses deletes all licenses.
	DeleteAllLicenses() error
}

type LicenseCache struct {
	licenses map[string]*models.License
	nameToID map[string]string
	provider Provider
}

func NewLicenseCache(provider Provider) *LicenseCache {
	return &LicenseCache{
		provider: provider,
	}
}

// Reload reloads the license cache.
func (lc *LicenseCache) Reload() error {
	return lc.loadLicenses()
}

// Get returns a license by its SPDX identifier.
func (lc *LicenseCache) Get(id string) *models.License {
	id = stringutil.NormalizeName(id)
	l, ok := lc.licenses[id]
	if ok {
		return l
	}
	return nil
}

// GetByIDOrName returns a license by its SPDX identifier or name.
func (lc *LicenseCache) GetByIDOrName(idOrName string) *models.License {
	idOrName = stringutil.NormalizeName(idOrName)
	// is name
	if id, ok := lc.nameToID[idOrName]; ok {
		l := lc.licenses[id]
		return l
	}
	// is ID
	l, ok := lc.licenses[idOrName]
	if ok {
		return l
	}
	return nil
}

// loadLicenses loads licenses from the provider.
func (lc *LicenseCache) loadLicenses() error {
	// get licenses from provider
	licenses, err := lc.provider.GetAllLicenses()
	if err != nil {
		return eris.Wrap(err, "failed to get a list of licenses")
	}

	lc.licenses = make(map[string]*models.License, len(licenses))
	lc.nameToID = make(map[string]string, len(licenses))
	for _, l := range licenses {
		normalizedID := stringutil.NormalizeName(l.ID)
		lc.licenses[normalizedID] = l
		lc.nameToID[stringutil.NormalizeName(l.Name)] = normalizedID
	}
	return nil
}
