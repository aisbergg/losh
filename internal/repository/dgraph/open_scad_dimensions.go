package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetOpenSCADDimensionsStr    = "failed to get osd(s)"
	errSaveOpenSCADDimensionsStr   = "failed to save osd(s)"
	errDeleteOpenSCADDimensionsStr = "failed to delete osd(s)"
)

// GetOpenSCADDimensions returns a `OpenSCADDimensions` object by its ID.
func (dr *DgraphRepository) GetOpenSCADDimensions(id string) (*models.OpenSCADDimensions, error) {
	ctx := context.Background()
	getOpenSCADDimensions, err := dr.client.GetOpenSCADDimensions(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetOpenSCADDimensionsStr).
			Add("osdId", id)
	}
	}
	osd := &models.OpenSCADDimensions{ID: id}
	if err = copier.CopyWithOption(osd, getOpenSCADDimensions.GetOpenSCADDimensions, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return osd, nil
}

// GetOpenSCADDimensionss returns a list of `OpenSCADDimensions` objects matching the filter criteria.
func (dr *DgraphRepository) GetOpenSCADDimensionss(filter *models.OpenSCADDimensionsFilter, order *models.OpenSCADDimensionsOrder, first *int64, offset *int64) ([]*models.OpenSCADDimensions, error) {
	ctx := context.Background()
	getOpenSCADDimensionss, err := dr.client.GetOpenSCADDimensionss(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetOpenSCADDimensionsStr)
	}
	osds := make([]*models.OpenSCADDimensions, 0, len(getOpenSCADDimensionss.QueryOpenSCADDimensions))
	for _, x := range getOpenSCADDimensionss.QueryOpenSCADDimensions {
		osd := &models.OpenSCADDimensions{ID: x.ID}
		if err = copier.CopyWithOption(osd, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		osds = append(osds, osd)
	}
	return osds, nil
}

// GetAllOpenSCADDimensionss returns a list of all `OpenSCADDimensions` objects.
func (dr *DgraphRepository) GetAllOpenSCADDimensionss() ([]*models.OpenSCADDimensions, error) {
	return dr.GetOpenSCADDimensionss(nil, nil, nil, nil)
}

// SaveOpenSCADDimensions saves a `OpenSCADDimensions` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveOpenSCADDimensions(osd *models.OpenSCADDimensions) (err error) {
	err = dr.SaveOpenSCADDimensionss([]*models.OpenSCADDimensions{osd})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("osdId", osd.ID)
	}
	return
}

// SaveOpenSCADDimensionss saves `OpenSCADDimensions` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveOpenSCADDimensionss(osds []*models.OpenSCADDimensions) error {
	reqData := make([]*models.AddOpenSCADDimensionsInput, 0, len(osds))
	for _, x := range osds {
		if x.ID != "" {
			continue
		}
		osd := &models.AddOpenSCADDimensionsInput{}
			return repository.WrapRepoError(err, errSaveOpenSCADDimensionsStr).
				Add("osdId", x.ID)
		}
		reqData = append(reqData, osd)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveOpenSCADDimensionss(ctx, reqData, []string{})
	if err != nil {
		return repository.WrapRepoError(err, errSaveOpenSCADDimensionsStr)
	}
	// save ID from response
	for i, x := range osds {
		x.ID = respData.AddOpenSCADDimensions.OpenSCADDimensions[i].ID
	}
	return nil
}

// DeleteOpenSCADDimensions deletes a `OpenSCADDimensions` object.
func (dr *DgraphRepository) DeleteOpenSCADDimensions(id *string) error {
	ctx := context.Background()
	delFilter := models.OpenSCADDimensionsFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteOpenSCADDimensionss(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteOpenSCADDimensionsStr).
			Add("osdId", id)
	}
	return nil
}

// DeleteAllOpenSCADDimensionss deletes all `OpenSCADDimensionss` objects.
func (dr *DgraphRepository) DeleteAllOpenSCADDimensionss() error {
	return dr.DeleteOpenSCADDimensions(nil)
}

// saveOpenSCADDimensionsIfNecessary saves a `OpenSCADDimensions` object if it is not already saved.
func (dr *DgraphRepository) saveOpenSCADDimensionsIfNecessary(osd *models.OpenSCADDimensions) (*models.OpenSCADDimensionsRef, error) {
	if osd == nil {
		return nil, nil
	}
	if osd.ID == "" {
		if err := dr.SaveOpenSCADDimensions(osd); err != nil {
			return nil, err
		}
	}
	return &models.OpenSCADDimensionsRef{ID: &osd.ID}, nil
}
