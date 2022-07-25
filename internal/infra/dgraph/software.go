package dgraph

import (
	"context"
	"losh/internal/core/product/models"
	"losh/internal/lib/errors"
	"losh/internal/repository"
)

var (
	errGetSoftwareStr    = "failed to get software(s)"
	errSaveSoftwareStr   = "failed to save software(s)"
	errDeleteSoftwareStr = "failed to delete software(s)"
)

// GetSoftware returns a `Software` object by its ID.
func (dr *DgraphRepository) GetSoftware(id string) (*models.Software, error) {
	ctx := context.Background()
	getSoftware, err := dr.client.GetSoftware(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetSoftwareStr).
			Add("softwareId", id)
	}
	if getSoftware.GetSoftware == nil { // not found
		return nil, nil
	}
	software := &models.Software{ID: id}
	if err = dr.dataCopier.CopyTo(getSoftware.GetSoftware, software); err != nil {
		panic(err)
	}
	return software, nil
}

// GetSoftwares returns a list of `Software` objects matching the filter criteria.
func (dr *DgraphRepository) GetSoftwares(filter *models.SoftwareFilter, order *models.SoftwareOrder, first *int64, offset *int64) ([]*models.Software, error) {
	ctx := context.Background()
	getSoftwares, err := dr.client.GetSoftwares(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetSoftwareStr)
	}
	softwares := make([]*models.Software, 0, len(getSoftwares.QuerySoftware))
	for _, x := range getSoftwares.QuerySoftware {
		software := &models.Software{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, software); err != nil {
			panic(err)
		}
		softwares = append(softwares, software)
	}
	return softwares, nil
}

// GetAllSoftwares returns a list of all `Software` objects.
func (dr *DgraphRepository) GetAllSoftwares() ([]*models.Software, error) {
	return dr.GetSoftwares(nil, nil, nil, nil)
}

// SaveSoftware saves a `Software` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveSoftware(software *models.Software) (err error) {
	err = dr.SaveSoftwares([]*models.Software{software})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("softwareId", software.ID)
	}
	return
}

// SaveSoftwares saves `Software` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveSoftwares(softwares []*models.Software) error {
	reqData := make([]*models.AddSoftwareInput, 0, len(softwares))
	for _, x := range softwares {
		if x.ID != "" {
			continue
		}
		software := &models.AddSoftwareInput{}
		if err := dr.dataCopier.CopyTo(x, software); err != nil {
			return repository.WrapRepoError(err, errSaveSoftwareStr).
				Add("softwareId", x.ID)
		}
		reqData = append(reqData, software)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveSoftwares(ctx, reqData, []string{})
	if err != nil {
		return repository.WrapRepoError(err, errSaveSoftwareStr)
	}
	// save ID from response
	for i, x := range softwares {
		x.ID = respData.AddSoftware.Software[i].ID
	}
	return nil
}

// DeleteSoftware deletes a `Software` object.
func (dr *DgraphRepository) DeleteSoftware(id *string) error {
	ctx := context.Background()
	delFilter := models.SoftwareFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteSoftware(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteSoftwareStr).
			Add("softwareId", id)
	}
	return nil
}

// DeleteAllSoftwares deletes all `Softwares` objects.
func (dr *DgraphRepository) DeleteAllSoftwares() error {
	return dr.DeleteSoftware(nil)
}

// saveSoftwareIfNecessary saves a `Software` object if it is not already saved.
func (dr *DgraphRepository) saveSoftwareIfNecessary(software *models.Software) (*models.SoftwareRef, error) {
	if software == nil {
		return nil, nil
	}
	if software.ID == "" {
		if err := dr.SaveSoftware(software); err != nil {
			return nil, err
		}
	}
	return &models.SoftwareRef{ID: &software.ID}, nil
}
