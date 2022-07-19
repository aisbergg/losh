package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetTechnologySpecificDocumentationCriteriaStr    = "failed to get tsdc(s)"
	errSaveTechnologySpecificDocumentationCriteriaStr   = "failed to save tsdc(s)"
	errDeleteTechnologySpecificDocumentationCriteriaStr = "failed to delete tsdc(s)"
)

// GetTechnologySpecificDocumentationCriteria returns a `TechnologySpecificDocumentationCriteria` object by its ID.
func (dr *DgraphRepository) GetTechnologySpecificDocumentationCriteria(id, xid *string) (*models.TechnologySpecificDocumentationCriteria, error) {
	ctx := context.Background()
	getTechnologySpecificDocumentationCriteria, err := dr.client.GetTechnologySpecificDocumentationCriteria(ctx, id, xid)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetTechnologySpecificDocumentationCriteriaStr).
			AddIfNotNil("tsdcId", id).AddIfNotNil("tsdcXid", xid)
	}
	tsdc := &models.TechnologySpecificDocumentationCriteria{ID: *id}
	if err = copier.CopyWithOption(tsdc, getTechnologySpecificDocumentationCriteria.GetTechnologySpecificDocumentationCriteria, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return tsdc, nil
}

// GetTechnologySpecificDocumentationCriterias returns a list of `TechnologySpecificDocumentationCriteria` objects matching the filter criteria.
func (dr *DgraphRepository) GetTechnologySpecificDocumentationCriterias(filter *models.TechnologySpecificDocumentationCriteriaFilter, order *models.TechnologySpecificDocumentationCriteriaOrder, first *int64, offset *int64) ([]*models.TechnologySpecificDocumentationCriteria, error) {
	ctx := context.Background()
	getTechnologySpecificDocumentationCriterias, err := dr.client.GetTechnologySpecificDocumentationCriterias(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetTechnologySpecificDocumentationCriteriaStr)
	}
	tsdcs := make([]*models.TechnologySpecificDocumentationCriteria, 0, len(getTechnologySpecificDocumentationCriterias.QueryTechnologySpecificDocumentationCriteria))
	for _, x := range getTechnologySpecificDocumentationCriterias.QueryTechnologySpecificDocumentationCriteria {
		tsdc := &models.TechnologySpecificDocumentationCriteria{ID: x.ID}
		if err = copier.CopyWithOption(tsdc, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		tsdcs = append(tsdcs, tsdc)
	}
	return tsdcs, nil
}

// GetAllTechnologySpecificDocumentationCriterias returns a list of all `TechnologySpecificDocumentationCriteria` objects.
func (dr *DgraphRepository) GetAllTechnologySpecificDocumentationCriterias() ([]*models.TechnologySpecificDocumentationCriteria, error) {
	return dr.GetTechnologySpecificDocumentationCriterias(nil, nil, nil, nil)
}

// SaveTechnologySpecificDocumentationCriteria saves a `TechnologySpecificDocumentationCriteria` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveTechnologySpecificDocumentationCriteria(tsdc *models.TechnologySpecificDocumentationCriteria) (err error) {
	err = dr.SaveTechnologySpecificDocumentationCriterias([]*models.TechnologySpecificDocumentationCriteria{tsdc})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.AddIfNotNil("tsdcId", tsdc.ID).AddIfNotNil("tsdcXid", tsdc.Xid)
	}
	return
}

// SaveTechnologySpecificDocumentationCriterias saves `TechnologySpecificDocumentationCriteria` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveTechnologySpecificDocumentationCriterias(tsdcs []*models.TechnologySpecificDocumentationCriteria) error {
	reqData := make([]*models.AddTechnologySpecificDocumentationCriteriaInput, 0, len(tsdcs))
	for _, x := range tsdcs {
		if x.ID != "" {
			continue
		}
		tsdc := &models.AddTechnologySpecificDocumentationCriteriaInput{}
		if err := copier.CopyWithOption(tsdc, x,
			copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
			return repository.NewRepoErrorWrap(err, errSaveTechnologySpecificDocumentationCriteriaStr).
				AddIfNotNil("tsdcId", x.ID).AddIfNotNil("tsdcXid", x.Xid)
		}
		reqData = append(reqData, tsdc)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveTechnologySpecificDocumentationCriterias(ctx, reqData)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errSaveTechnologySpecificDocumentationCriteriaStr)
	}
	// save ID from response
	for i, x := range tsdcs {
		x.ID = respData.AddTechnologySpecificDocumentationCriteria.TechnologySpecificDocumentationCriteria[i].ID
	}
	return nil
}

// DeleteTechnologySpecificDocumentationCriteria deletes a `TechnologySpecificDocumentationCriteria` object.
func (dr *DgraphRepository) DeleteTechnologySpecificDocumentationCriteria(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.TechnologySpecificDocumentationCriteriaFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteTechnologySpecificDocumentationCriteria(ctx, delFilter)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errDeleteTechnologySpecificDocumentationCriteriaStr).
			AddIfNotNil("tsdcId", id).AddIfNotNil("tsdcXid", xid)
	}
	return nil
}

// DeleteAllTechnologySpecificDocumentationCriterias deletes all `TechnologySpecificDocumentationCriterias` objects.
func (dr *DgraphRepository) DeleteAllTechnologySpecificDocumentationCriterias() error {
	return dr.DeleteTechnologySpecificDocumentationCriteria(nil, nil)
}

// saveTechnologySpecificDocumentationCriteriaIfNecessary saves a `TechnologySpecificDocumentationCriteria` object if it is not already saved.
func (dr *DgraphRepository) saveTechnologySpecificDocumentationCriteriaIfNecessary(tsdc *models.TechnologySpecificDocumentationCriteria) (*models.TechnologySpecificDocumentationCriteriaRef, error) {
	if tsdc == nil {
		return nil, nil
	}
	if tsdc.ID == "" {
		if err := dr.SaveTechnologySpecificDocumentationCriteria(tsdc); err != nil {
			return nil, err
		}
	}
	return &models.TechnologySpecificDocumentationCriteriaRef{ID: &tsdc.ID}, nil
}
