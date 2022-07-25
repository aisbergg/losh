package services

import (
	"losh/internal/core/product/models"
	"losh/internal/lib/util/stringutil"

	"github.com/aisbergg/go-errors/pkg/errors"
)

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
		return errors.Wrap(err, "failed to get a list of licenses")
	}

	lc.licenses = make(map[string]*models.License, len(licenses))
	lc.nameToID = make(map[string]string, len(licenses))
	for _, l := range licenses {
		// remove text and textHTML, because we don't need them here
		l.Text = nil
		l.TextHTML = nil

		// save license under their normalized ID/name
		normalizedID := stringutil.NormalizeName(l.Xid)
		lc.licenses[normalizedID] = l
		lc.nameToID[stringutil.NormalizeName(l.Name)] = normalizedID
	}
	return nil
}
