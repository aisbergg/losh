package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetTechnicalStandardStr    = "failed to get technicalStandard(s)"
	errSaveTechnicalStandardStr   = "failed to save technicalStandard(s)"
	errDeleteTechnicalStandardStr = "failed to delete technicalStandard(s)"
)

// GetTechnicalStandard returns a `TechnicalStandard` object by its ID.
func (dr *DgraphRepository) GetTechnicalStandard(id, xid *string) (*models.TechnicalStandard, error) {
	ctx := context.Background()
	getTechnicalStandard, err := dr.client.GetTechnicalStandard(ctx, id, xid)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetTechnicalStandardStr).
			AddIfNotNil("technicalStandardId", id).AddIfNotNil("technicalStandardXid", xid)
	}
	technicalStandard := &models.TechnicalStandard{ID: *id}
	if err = copier.CopyWithOption(technicalStandard, getTechnicalStandard.GetTechnicalStandard, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return technicalStandard, nil
}

// GetTechnicalStandards returns a list of `TechnicalStandard` objects matching the filter criteria.
func (dr *DgraphRepository) GetTechnicalStandards(filter *models.TechnicalStandardFilter, order *models.TechnicalStandardOrder, first *int64, offset *int64) ([]*models.TechnicalStandard, error) {
	ctx := context.Background()
	getTechnicalStandards, err := dr.client.GetTechnicalStandards(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetTechnicalStandardStr)
	}
	technicalStandards := make([]*models.TechnicalStandard, 0, len(getTechnicalStandards.QueryTechnicalStandard))
	for _, x := range getTechnicalStandards.QueryTechnicalStandard {
		technicalStandard := &models.TechnicalStandard{ID: x.ID}
		if err = copier.CopyWithOption(technicalStandard, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		technicalStandards = append(technicalStandards, technicalStandard)
	}
	return technicalStandards, nil
}

// GetAllTechnicalStandards returns a list of all `TechnicalStandard` objects.
func (dr *DgraphRepository) GetAllTechnicalStandards() ([]*models.TechnicalStandard, error) {
	return dr.GetTechnicalStandards(nil, nil, nil, nil)
}

// SaveTechnicalStandard saves a `TechnicalStandard` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveTechnicalStandard(technicalStandard *models.TechnicalStandard) (err error) {
	err = dr.SaveTechnicalStandards([]*models.TechnicalStandard{technicalStandard})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.AddIfNotNil("technicalStandardId", technicalStandard.ID).AddIfNotNil("technicalStandardXid", technicalStandard.Xid)
	}
	return
}

// SaveTechnicalStandards saves `TechnicalStandard` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveTechnicalStandards(technicalStandards []*models.TechnicalStandard) error {
	reqData := make([]*models.AddTechnicalStandardInput, 0, len(technicalStandards))
	for _, x := range technicalStandards {
		if x.ID != "" {
			continue
		}
		technicalStandard := &models.AddTechnicalStandardInput{}
		if err := copier.CopyWithOption(technicalStandard, x,
			copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
			return repository.NewRepoErrorWrap(err, errSaveTechnicalStandardStr).
				AddIfNotNil("technicalStandardId", x.ID).AddIfNotNil("technicalStandardXid", x.Xid)
		}
		reqData = append(reqData, technicalStandard)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveTechnicalStandards(ctx, reqData)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errSaveTechnicalStandardStr)
	}
	// save ID from response
	for i, x := range technicalStandards {
		x.ID = respData.AddTechnicalStandard.TechnicalStandard[i].ID
	}
	return nil
}

// DeleteTechnicalStandard deletes a `TechnicalStandard` object.
func (dr *DgraphRepository) DeleteTechnicalStandard(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.TechnicalStandardFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteTechnicalStandard(ctx, delFilter)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errDeleteTechnicalStandardStr).
			AddIfNotNil("technicalStandardId", id).AddIfNotNil("technicalStandardXid", xid)
	}
	return nil
}

// DeleteAllTechnicalStandards deletes all `TechnicalStandards` objects.
func (dr *DgraphRepository) DeleteAllTechnicalStandards() error {
	return dr.DeleteTechnicalStandard(nil, nil)
}

// saveTechnicalStandardIfNecessary saves a `TechnicalStandard` object if it is not already saved.
func (dr *DgraphRepository) saveTechnicalStandardIfNecessary(technicalStandard *models.TechnicalStandard) (*models.TechnicalStandardRef, error) {
	if technicalStandard == nil {
		return nil, nil
	}
	if technicalStandard.ID == "" {
		if err := dr.SaveTechnicalStandard(technicalStandard); err != nil {
			return nil, err
		}
	}
	return &models.TechnicalStandardRef{ID: &technicalStandard.ID}, nil
}