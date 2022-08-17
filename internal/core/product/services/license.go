package services

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/lib/util/stringutil"

	"github.com/aisbergg/go-errors/pkg/errors"
)

// // Get returns a license by its SPDX identifier.
// func (s *Service) Get(id string) *models.License {
// 	id = stringutil.NormalizeName(id)
// 	l, ok := s.licenses[id]
// 	if ok {
// 		return l
// 	}
// 	return nil
// }

// // GetLicensesFromSPDXOrg returns a list of licenses from the SPDX.org license
// // list.
// func (s *Service) GetLicensesFromSPDXOrg(ctx context.Context) ([]*models.License, error) {
// 	return s.spdxLicenseProvider.GetAllLicenses(ctx)
// }

// GetCachedLicenseByIDOrName returns a license by its SPDX identifier or name.
func (s *Service) GetCachedLicenseByIDOrName(idOrName string) *models.License {
	idOrName = stringutil.NormalizeName(idOrName)
	// is name
	if id, ok := s.nameToID[idOrName]; ok {
		l := s.licenses[id]
		return l
	}
	// is ID
	l, ok := s.licenses[idOrName]
	if ok {
		return l
	}
	return nil
}

// ReloadLicenseCache reloads the license cache.
func (s *Service) ReloadLicenseCache() error {
	return s.loadLicenses()
}

// loadLicenses loads licenses from the repo.
func (s *Service) loadLicenses() error {
	// get licenses from repo
	licenses, err := s.repo.GetAllLicensesBasic(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get a list of licenses")
	}

	s.licenses = make(map[string]*models.License, len(licenses))
	s.nameToID = make(map[string]string, len(licenses))
	for _, l := range licenses {
		// save license under their normalized ID/name
		normalizedID := stringutil.NormalizeName(*l.Xid)
		s.licenses[normalizedID] = l
		s.nameToID[stringutil.NormalizeName(*l.Name)] = normalizedID
	}
	return nil
}
