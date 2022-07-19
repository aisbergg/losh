package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetLicenseStr    = "failed to get license(s)"
	errSaveLicenseStr   = "failed to save license(s)"
	errDeleteLicenseStr = "failed to delete license(s)"
)

// GetLicense returns a `License` object by its ID.
func (dr *DgraphRepository) GetLicense(id, xid *string) (*models.License, error) {
	ctx := context.Background()
	getLicense, err := dr.client.GetLicense(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetLicenseStr).
			Add("licenseId", id).Add("licenseXid", xid)
	}
	if getLicense.GetLicense == nil { // not found
		return nil, nil
	}
	license := &models.License{ID: *id}
	if err = dr.dataCopier.CopyTo(getLicense.GetLicense, license); err != nil {
		panic(err)
	}
	return license, nil
}

// GetLicenses returns a list of `License` objects matching the filter criteria.
func (dr *DgraphRepository) GetLicenses(filter *models.LicenseFilter, order *models.LicenseOrder, first *int64, offset *int64) ([]*models.License, error) {
	ctx := context.Background()
	getLicenses, err := dr.client.GetLicenses(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetLicenseStr)
	}
	licenses := make([]*models.License, 0, len(getLicenses.QueryLicense))
	for _, x := range getLicenses.QueryLicense {
		license := &models.License{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, license); err != nil {
			panic(err)
		}
		licenses = append(licenses, license)
	}
	return licenses, nil
}

// GetAllLicenses returns a list of all `License` objects.
func (dr *DgraphRepository) GetAllLicenses() ([]*models.License, error) {
	return dr.GetLicenses(nil, nil, nil, nil)
}

// SaveLicense saves a `License` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveLicense(license *models.License) (err error) {
	err = dr.SaveLicenses([]*models.License{license})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("licenseId", license.ID).Add("licenseXid", license.Xid)
	}
	return
}

// SaveLicenses saves `License` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveLicenses(licenses []*models.License) error {
	reqData := make([]*models.AddLicenseInput, 0, len(licenses))
	for _, x := range licenses {
		if x.ID != "" {
			continue
		}
		license := &models.AddLicenseInput{}
		if err := dr.dataCopier.CopyTo(x, license); err != nil {
			return repository.WrapRepoError(err, errSaveLicenseStr).
				Add("licenseId", x.ID).Add("licenseXid", x.Xid)
		}
		reqData = append(reqData, license)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveLicenses(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveLicenseStr)
	}
	// save ID from response
	for i, x := range licenses {
		x.ID = respData.AddLicense.License[i].ID
	}
	return nil
}

// DeleteLicense deletes a `License` object.
func (dr *DgraphRepository) DeleteLicense(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.LicenseFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteLicense(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteLicenseStr).
			Add("licenseId", id).Add("licenseXid", xid)
	}
	return nil
}

// DeleteAllLicenses deletes all `Licenses` objects.
func (dr *DgraphRepository) DeleteAllLicenses() error {
	return dr.DeleteLicense(nil, nil)
}

// saveLicenseIfNecessary saves a `License` object if it is not already saved.
func (dr *DgraphRepository) saveLicenseIfNecessary(license *models.License) (*models.LicenseRef, error) {
	if license == nil {
		return nil, nil
	}
	if license.ID == "" {
		if err := dr.SaveLicense(license); err != nil {
			return nil, err
		}
	}
	return &models.LicenseRef{ID: &license.ID}, nil
}
