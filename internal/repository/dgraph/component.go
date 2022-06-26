package dgraph

import (
	"context"
	"losh/crawler/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetComponentStr    = "failed to get component(s)"
	errSaveComponentStr   = "failed to save component(s)"
	errDeleteComponentStr = "failed to delete component(s)"
)

// GetComponent returns a `Component` object by its ID.
func (dr *DgraphRepository) GetComponent(id, xid *string) (*models.Component, error) {
	ctx := context.Background()
	getComponent, err := dr.client.GetComponent(ctx, id, xid)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetComponentStr).
			AddIfNotNil("componentId", id).AddIfNotNil("componentXid", xid)
	}
	component := &models.Component{ID: *id}
	if err = copier.CopyWithOption(component, getComponent.GetComponent, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return component, nil
}

// GetComponents returns a list of `Component` objects matching the filter criteria.
func (dr *DgraphRepository) GetComponents(filter *models.ComponentFilter, order *models.ComponentOrder, first *int64, offset *int64) ([]*models.Component, error) {
	ctx := context.Background()
	getComponents, err := dr.client.GetComponents(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetComponentStr)
	}
	components := make([]*models.Component, 0, len(getComponents.QueryComponent))
	for _, x := range getComponents.QueryComponent {
		component := &models.Component{ID: x.ID}
		if err = copier.CopyWithOption(component, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		components = append(components, component)
	}
	return components, nil
}

// GetAllComponents returns a list of all `Component` objects.
func (dr *DgraphRepository) GetAllComponents() ([]*models.Component, error) {
	return dr.GetComponents(nil, nil, nil, nil)
}

// SaveComponent saves a `Component` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveComponent(component *models.Component) (err error) {
	err = dr.SaveComponents([]*models.Component{component})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.AddIfNotNil("componentId", component.ID).AddIfNotNil("componentXid", component.Xid)
	}
	return
}

// SaveComponents saves `Component` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveComponents(components []*models.Component) error {
	reqData := make([]*models.AddComponentInput, 0, len(components))
	for _, x := range components {
		if x.ID != "" {
			continue
		}
		component := &models.AddComponentInput{}
		if err := copier.CopyWithOption(component, x,
			copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
			return repository.NewRepoErrorWrap(err, errSaveComponentStr).
				AddIfNotNil("componentId", x.ID).AddIfNotNil("componentXid", x.Xid)
		}
		reqData = append(reqData, component)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveComponents(ctx, reqData)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errSaveComponentStr)
	}
	// save ID from response
	for i, x := range components {
		x.ID = respData.AddComponent.Component[i].ID
	}
	return nil
}

// DeleteComponent deletes a `Component` object.
func (dr *DgraphRepository) DeleteComponent(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.ComponentFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteComponent(ctx, delFilter)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errDeleteComponentStr).
			AddIfNotNil("componentId", id).AddIfNotNil("componentXid", xid)
	}
	return nil
}

// DeleteAllComponents deletes all `Components` objects.
func (dr *DgraphRepository) DeleteAllComponents() error {
	return dr.DeleteComponent(nil, nil)
}

// saveComponentIfNecessary saves a `Component` object if it is not already saved.
func (dr *DgraphRepository) saveComponentIfNecessary(component *models.Component) (*models.ComponentRef, error) {
	if component == nil {
		return nil, nil
	}
	if component.ID == "" {
		if err := dr.SaveComponent(component); err != nil {
			return nil, err
		}
	}
	return &models.ComponentRef{ID: &component.ID}, nil
}
