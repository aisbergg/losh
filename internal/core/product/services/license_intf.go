package services

import "losh/internal/core/product/models"

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
