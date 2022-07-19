package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetComponentSourceStr    = "failed to get componentSource(s)"
	errSaveComponentSourceStr   = "failed to save componentSource(s)"
	errDeleteComponentSourceStr = "failed to delete componentSource(s)"
)

// GetComponentSource returns a `ComponentSource` object by its ID.
func (dr *DgraphRepository) GetComponentSource(id, xid *string) (*models.ComponentSource, error) {
	ctx := context.Background()
	getComponentSource, err := dr.client.GetComponentSource(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetComponentSourceStr).
			Add("componentSourceId", id).Add("componentSourceXid", xid)
	}
	}
	componentSource := &models.ComponentSource{ID: *id}
	if err = copier.CopyWithOption(componentSource, getComponentSource.GetComponentSource, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return componentSource, nil
}

// GetComponentSources returns a list of `ComponentSource` objects matching the filter criteria.
func (dr *DgraphRepository) GetComponentSources(filter *models.ComponentSourceFilter, order *models.ComponentSourceOrder, first *int64, offset *int64) ([]*models.ComponentSource, error) {
	ctx := context.Background()
	getComponentSources, err := dr.client.GetComponentSources(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetComponentSourceStr)
	}
	componentSources := make([]*models.ComponentSource, 0, len(getComponentSources.QueryComponentSource))
	for _, x := range getComponentSources.QueryComponentSource {
		componentSource := &models.ComponentSource{ID: x.ID}
		if err = copier.CopyWithOption(componentSource, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		componentSources = append(componentSources, componentSource)
	}
	return componentSources, nil
}

// GetAllComponentSources returns a list of all `ComponentSource` objects.
func (dr *DgraphRepository) GetAllComponentSources() ([]*models.ComponentSource, error) {
	return dr.GetComponentSources(nil, nil, nil, nil)
}

// SaveComponentSource saves a `ComponentSource` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveComponentSource(componentSource *models.ComponentSource) (err error) {
	err = dr.SaveComponentSources([]*models.ComponentSource{componentSource})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("componentSourceId", componentSource.ID).Add("componentSourceXid", componentSource.Xid)
	}
	return
}

// SaveComponentSources saves `ComponentSource` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveComponentSources(componentSources []*models.ComponentSource) error {
	reqData := make([]*models.AddComponentSourceInput, 0, len(componentSources))
	for _, x := range componentSources {
		if x.ID != "" {
			continue
		}
		componentSource := &models.AddComponentSourceInput{}
			return repository.WrapRepoError(err, errSaveComponentSourceStr).
				Add("componentSourceId", x.ID).Add("componentSourceXid", x.Xid)
		}
		reqData = append(reqData, componentSource)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveComponentSources(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveComponentSourceStr)
	}
	// save ID from response
	for i, x := range componentSources {
		x.ID = respData.AddComponentSource.ComponentSource[i].ID
	}
	return nil
}

// DeleteComponentSource deletes a `ComponentSource` object.
func (dr *DgraphRepository) DeleteComponentSource(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.ComponentSourceFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteComponentSource(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteComponentSourceStr).
			Add("componentSourceId", id).Add("componentSourceXid", xid)
	}
	return nil
}

// DeleteAllComponentSources deletes all `ComponentSources` objects.
func (dr *DgraphRepository) DeleteAllComponentSources() error {
	return dr.DeleteComponentSource(nil, nil)
}

// saveComponentSourceIfNecessary saves a `ComponentSource` object if it is not already saved.
func (dr *DgraphRepository) saveComponentSourceIfNecessary(componentSource *models.ComponentSource) (*models.ComponentSourceRef, error) {
	if componentSource == nil {
		return nil, nil
	}
	if componentSource.ID == "" {
		if err := dr.SaveComponentSource(componentSource); err != nil {
			return nil, err
		}
	}
	return &models.ComponentSourceRef{ID: &componentSource.ID}, nil
}
