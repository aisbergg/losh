package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetBoundingBoxDimensionsStr    = "failed to get bbd(s)"
	errSaveBoundingBoxDimensionsStr   = "failed to save bbd(s)"
	errDeleteBoundingBoxDimensionsStr = "failed to delete bbd(s)"
)

// GetBoundingBoxDimensions returns a `BoundingBoxDimensions` object by its ID.
func (dr *DgraphRepository) GetBoundingBoxDimensions(id string) (*models.BoundingBoxDimensions, error) {
	ctx := context.Background()
	getBoundingBoxDimensions, err := dr.client.GetBoundingBoxDimensions(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetBoundingBoxDimensionsStr).
			Add("bbdId", id)
	}
	if getBoundingBoxDimensions.GetBoundingBoxDimensions == nil { // not found
		return nil, nil
	}
	bbd := &models.BoundingBoxDimensions{ID: id}
	if err = dr.dataCopier.CopyTo(getBoundingBoxDimensions.GetBoundingBoxDimensions, bbd); err != nil {
		panic(err)
	}
	return bbd, nil
}

// GetBoundingBoxDimensionss returns a list of `BoundingBoxDimensions` objects matching the filter criteria.
func (dr *DgraphRepository) GetBoundingBoxDimensionss(filter *models.BoundingBoxDimensionsFilter, order *models.BoundingBoxDimensionsOrder, first *int64, offset *int64) ([]*models.BoundingBoxDimensions, error) {
	ctx := context.Background()
	getBoundingBoxDimensionss, err := dr.client.GetBoundingBoxDimensionss(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetBoundingBoxDimensionsStr)
	}
	bbds := make([]*models.BoundingBoxDimensions, 0, len(getBoundingBoxDimensionss.QueryBoundingBoxDimensions))
	for _, x := range getBoundingBoxDimensionss.QueryBoundingBoxDimensions {
		bbd := &models.BoundingBoxDimensions{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, bbd); err != nil {
			panic(err)
		}
		bbds = append(bbds, bbd)
	}
	return bbds, nil
}

// GetAllBoundingBoxDimensionss returns a list of all `BoundingBoxDimensions` objects.
func (dr *DgraphRepository) GetAllBoundingBoxDimensionss() ([]*models.BoundingBoxDimensions, error) {
	return dr.GetBoundingBoxDimensionss(nil, nil, nil, nil)
}

// SaveBoundingBoxDimensions saves a `BoundingBoxDimensions` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveBoundingBoxDimensions(bbd *models.BoundingBoxDimensions) (err error) {
	err = dr.SaveBoundingBoxDimensionss([]*models.BoundingBoxDimensions{bbd})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("bbdId", bbd.ID)
	}
	return
}

// SaveBoundingBoxDimensionss saves `BoundingBoxDimensions` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveBoundingBoxDimensionss(bbds []*models.BoundingBoxDimensions) error {
	reqData := make([]*models.AddBoundingBoxDimensionsInput, 0, len(bbds))
	for _, x := range bbds {
		if x.ID != "" {
			continue
		}
		bbd := &models.AddBoundingBoxDimensionsInput{}
		if err := dr.dataCopier.CopyTo(x, bbd); err != nil {
			return repository.WrapRepoError(err, errSaveBoundingBoxDimensionsStr).
				Add("bbdId", x.ID)
		}
		reqData = append(reqData, bbd)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveBoundingBoxDimensionss(ctx, reqData, []string{})
	if err != nil {
		return repository.WrapRepoError(err, errSaveBoundingBoxDimensionsStr)
	}
	// save ID from response
	for i, x := range bbds {
		x.ID = respData.AddBoundingBoxDimensions.BoundingBoxDimensions[i].ID
	}
	return nil
}

// DeleteBoundingBoxDimensions deletes a `BoundingBoxDimensions` object.
func (dr *DgraphRepository) DeleteBoundingBoxDimensions(id *string) error {
	ctx := context.Background()
	delFilter := models.BoundingBoxDimensionsFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteBoundingBoxDimensions(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteBoundingBoxDimensionsStr).
			Add("bbdId", id)
	}
	return nil
}

// DeleteAllBoundingBoxDimensionss deletes all `BoundingBoxDimensionss` objects.
func (dr *DgraphRepository) DeleteAllBoundingBoxDimensionss() error {
	return dr.DeleteBoundingBoxDimensions(nil)
}

// saveBoundingBoxDimensionsIfNecessary saves a `BoundingBoxDimensions` object if it is not already saved.
func (dr *DgraphRepository) saveBoundingBoxDimensionsIfNecessary(bbd *models.BoundingBoxDimensions) (*models.BoundingBoxDimensionsRef, error) {
	if bbd == nil {
		return nil, nil
	}
	if bbd.ID == "" {
		if err := dr.SaveBoundingBoxDimensions(bbd); err != nil {
			return nil, err
		}
	}
	return &models.BoundingBoxDimensionsRef{ID: &bbd.ID}, nil
}
