package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetStringVStr    = "failed to get stringV(s)"
	errSaveStringVStr   = "failed to save stringV(s)"
	errDeleteStringVStr = "failed to delete stringV(s)"
)

// GetStringV returns a `StringV` object by its ID.
func (dr *DgraphRepository) GetStringV(id string) (*models.StringV, error) {
	ctx := context.Background()
	getStringV, err := dr.client.GetStringV(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetStringVStr).
			Add("stringVId", id)
	}
	}
	stringV := &models.StringV{ID: id}
	if err = copier.CopyWithOption(stringV, getStringV.GetStringV, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return stringV, nil
}

// GetStringVs returns a list of `StringV` objects matching the filter criteria.
func (dr *DgraphRepository) GetStringVs(filter *models.StringVFilter, order *models.StringVOrder, first *int64, offset *int64) ([]*models.StringV, error) {
	ctx := context.Background()
	getStringVs, err := dr.client.GetStringVs(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetStringVStr)
	}
	stringVs := make([]*models.StringV, 0, len(getStringVs.QueryStringV))
	for _, x := range getStringVs.QueryStringV {
		stringV := &models.StringV{ID: x.ID}
		if err = copier.CopyWithOption(stringV, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		stringVs = append(stringVs, stringV)
	}
	return stringVs, nil
}

// GetAllStringVs returns a list of all `StringV` objects.
func (dr *DgraphRepository) GetAllStringVs() ([]*models.StringV, error) {
	return dr.GetStringVs(nil, nil, nil, nil)
}

// SaveStringV saves a `StringV` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveStringV(stringV *models.StringV) (err error) {
	err = dr.SaveStringVs([]*models.StringV{stringV})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("stringVId", stringV.ID)
	}
	return
}

// SaveStringVs saves `StringV` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveStringVs(stringVs []*models.StringV) error {
	reqData := make([]*models.AddStringVInput, 0, len(stringVs))
	for _, x := range stringVs {
		if x.ID != "" {
			continue
		}
		stringV := &models.AddStringVInput{}
			return repository.WrapRepoError(err, errSaveStringVStr).
				Add("stringVId", x.ID)
		}
		reqData = append(reqData, stringV)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveStringVs(ctx, reqData, []string{})
	if err != nil {
		return repository.WrapRepoError(err, errSaveStringVStr)
	}
	// save ID from response
	for i, x := range stringVs {
		x.ID = respData.AddStringV.StringV[i].ID
	}
	return nil
}

// DeleteStringV deletes a `StringV` object.
func (dr *DgraphRepository) DeleteStringV(id *string) error {
	ctx := context.Background()
	delFilter := models.StringVFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteStringVs(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteStringVStr).
			Add("stringVId", id)
	}
	return nil
}

// DeleteAllStringVs deletes all `StringVs` objects.
func (dr *DgraphRepository) DeleteAllStringVs() error {
	return dr.DeleteStringV(nil)
}

// saveStringVIfNecessary saves a `StringV` object if it is not already saved.
func (dr *DgraphRepository) saveStringVIfNecessary(stringV *models.StringV) (*models.StringVRef, error) {
	if stringV == nil {
		return nil, nil
	}
	if stringV.ID == "" {
		if err := dr.SaveStringV(stringV); err != nil {
			return nil, err
		}
	}
	return &models.StringVRef{ID: &stringV.ID}, nil
}
