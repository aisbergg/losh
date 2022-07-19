package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetFloatVStr    = "failed to get floatV(s)"
	errSaveFloatVStr   = "failed to save floatV(s)"
	errDeleteFloatVStr = "failed to delete floatV(s)"
)

// GetFloatV returns a `FloatV` object by its ID.
func (dr *DgraphRepository) GetFloatV(id string) (*models.FloatV, error) {
	ctx := context.Background()
	getFloatV, err := dr.client.GetFloatV(ctx, id)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetFloatVStr).
			AddIfNotNil("floatVId", id)
	}
	floatV := &models.FloatV{ID: id}
	if err = copier.CopyWithOption(floatV, getFloatV.GetFloatV, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return floatV, nil
}

// GetFloatVs returns a list of `FloatV` objects matching the filter criteria.
func (dr *DgraphRepository) GetFloatVs(filter *models.FloatVFilter, order *models.FloatVOrder, first *int64, offset *int64) ([]*models.FloatV, error) {
	ctx := context.Background()
	getFloatVs, err := dr.client.GetFloatVs(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetFloatVStr)
	}
	floatVs := make([]*models.FloatV, 0, len(getFloatVs.QueryFloatV))
	for _, x := range getFloatVs.QueryFloatV {
		floatV := &models.FloatV{ID: x.ID}
		if err = copier.CopyWithOption(floatV, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		floatVs = append(floatVs, floatV)
	}
	return floatVs, nil
}

// GetAllFloatVs returns a list of all `FloatV` objects.
func (dr *DgraphRepository) GetAllFloatVs() ([]*models.FloatV, error) {
	return dr.GetFloatVs(nil, nil, nil, nil)
}

// SaveFloatV saves a `FloatV` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveFloatV(floatV *models.FloatV) (err error) {
	err = dr.SaveFloatVs([]*models.FloatV{floatV})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.AddIfNotNil("floatVId", floatV.ID)
	}
	return
}

// SaveFloatVs saves `FloatV` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveFloatVs(floatVs []*models.FloatV) error {
	reqData := make([]*models.AddFloatVInput, 0, len(floatVs))
	for _, x := range floatVs {
		if x.ID != "" {
			continue
		}
		floatV := &models.AddFloatVInput{}
		if err := copier.CopyWithOption(floatV, x,
			copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
			return repository.NewRepoErrorWrap(err, errSaveFloatVStr).
				AddIfNotNil("floatVId", x.ID)
		}
		reqData = append(reqData, floatV)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveFloatVs(ctx, reqData, []string{})
	if err != nil {
		return repository.NewRepoErrorWrap(err, errSaveFloatVStr)
	}
	// save ID from response
	for i, x := range floatVs {
		x.ID = respData.AddFloatV.FloatV[i].ID
	}
	return nil
}

// DeleteFloatV deletes a `FloatV` object.
func (dr *DgraphRepository) DeleteFloatV(id *string) error {
	ctx := context.Background()
	delFilter := models.FloatVFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteFloatVs(ctx, delFilter)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errDeleteFloatVStr).
			AddIfNotNil("floatVId", id)
	}
	return nil
}

// DeleteAllFloatVs deletes all `FloatVs` objects.
func (dr *DgraphRepository) DeleteAllFloatVs() error {
	return dr.DeleteFloatV(nil)
}

// saveFloatVIfNecessary saves a `FloatV` object if it is not already saved.
func (dr *DgraphRepository) saveFloatVIfNecessary(floatV *models.FloatV) (*models.FloatVRef, error) {
	if floatV == nil {
		return nil, nil
	}
	if floatV.ID == "" {
		if err := dr.SaveFloatV(floatV); err != nil {
			return nil, err
		}
	}
	return &models.FloatVRef{ID: &floatV.ID}, nil
}
