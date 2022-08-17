// Code generated by codegen, DO NOT EDIT.

package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net/request"
)

// make sure the struct implements the interface
var _ TagRepository = (*DgraphRepository)(nil)

// TagRepository is an interface for getting and saving `Tag` objects to a repository.
type TagRepository interface {
	GetTag(ctx context.Context, id, xid *string) (*models.Tag, error)
	GetTags(ctx context.Context, filter *dgclient.TagFilter, order *dgclient.TagOrder, first *int64, offset *int64) ([]*models.Tag, int64, error)
	GetAllTags(ctx context.Context) ([]*models.Tag, int64, error)
	CreateTag(ctx context.Context, input *models.Tag) error
	CreateTags(ctx context.Context, input []*models.Tag) error
	UpdateTag(ctx context.Context, input *models.Tag) error
	DeleteTag(ctx context.Context, id, xid *string) error
	DeleteAllTags(ctx context.Context) error
}

var (
	errGetTagStr    = "failed to get tag(s)"
	errSaveTagStr   = "failed to save tag(s)"
	errDeleteTagStr = "failed to delete tag(s)"
)

// GetTag returns a `Tag` object by its ID.
func (dr *DgraphRepository) GetTag(ctx context.Context, id, xid *string) (*models.Tag, error) {
	var rspData interface{}
	if id != nil {
		dr.log.Debugw("get Tag", "id", *id)
		rsp, err := dr.client.GetTagByID(ctx, *id)
		if err != nil {
			return nil, WrapRepoError(err, errGetTagStr).Add("tagId", id)
		}
		rspData = rsp.GetTag
	} else if xid != nil {
		dr.log.Debugw("get Tag", "xid", *xid)
		rsp, err := dr.client.GetTagByXid(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetTagStr).Add("tagXid", xid)
		}
		rspData = rsp.GetTag
	} else {
		panic("must specify id or xid")
	}

	if rspData == nil {
		return nil, nil
	}
	ret := &models.Tag{}
	if err := dr.copier.CopyTo(rspData, ret); err != nil {
		panic(err)
	}
	return ret, nil
}

// GetTagID returns the ID of an existing `Tag` object.
func (dr *DgraphRepository) GetTagID(ctx context.Context, xid *string) (*string, error) {
	if xid != nil {
		dr.log.Debugw("get Tag", "xid", *xid)
		rsp, err := dr.client.GetTagID(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetTagStr).Add("tagXid", xid)
		}
		if rsp.GetTag == nil {
			return nil, nil
		}
		return &rsp.GetTag.ID, nil
	}

	panic("must specify xid")
}

// GetTags returns a list of `Tag` objects matching the filter criteria.
func (dr *DgraphRepository) GetTags(ctx context.Context, filter *dgclient.TagFilter, order *dgclient.TagOrder, first *int64, offset *int64) ([]*models.Tag, int64, error) {
	dr.log.Debugw("get Tags")
	rsp, err := dr.client.GetTags(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetTagStr)
	}
	ret := make([]*models.Tag, 0, len(rsp.QueryTag))
	if err = dr.copier.CopyTo(rsp.QueryTag, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateTag.Count, nil
}

// GetAllTags returns a list of all `Tag` objects.
func (dr *DgraphRepository) GetAllTags(ctx context.Context) ([]*models.Tag, int64, error) {
	return dr.GetTags(ctx, nil, nil, nil, nil)
}

// GetTagWithCustomQuery returns a `Tag` object by its ID.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetTagWithCustomQuery(ctx context.Context, operationName, query string, id, xid *string) (*models.Tag, error) {
	req := request.GraphQLRequest{
		Ctx:           ctx,
		OperationName: operationName,
		Query:         query,
		Variables: map[string]interface{}{
			"id":  id,
			"xid": xid,
		},
	}
	rsp := struct {
		Tag *models.Tag "json:\"getTag\" graphql:\"getTag\""
	}{}
	dr.log.Debugw("get Tag with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetTagStr)
	}
	return rsp.Tag, nil
}

// GetTagsWithCustomQuery returns a list of `Tag` objects matching the filter criteria.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetTagsWithCustomQuery(ctx context.Context, operationName, query string, filter *dgclient.TagFilter, order *dgclient.TagOrder, first *int64, offset *int64) ([]*models.Tag, error) {
	req := request.GraphQLRequest{
		Ctx:           ctx,
		OperationName: operationName,
		Query:         query,
		Variables: map[string]interface{}{
			"filter": filter,
			"order":  order,
			"first":  first,
			"offset": offset,
		},
	}
	rsp := struct {
		Tags []*models.Tag "json:\"queryTag\" graphql:\"queryTag\""
	}{}
	dr.log.Debugw("get Tags with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetTagStr)
	}
	return rsp.Tags, nil
}

// GetAllTagsWithCustomQuery returns a list of all `Tag` objects.
func (dr *DgraphRepository) GetAllTagsWithCustomQuery(ctx context.Context, operationName, query string) ([]*models.Tag, error) {
	return dr.GetTagsWithCustomQuery(ctx, operationName, query, nil, nil, nil, nil)
}

// CreateTag creates a new `Tag` object.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateTag(ctx context.Context, input *models.Tag) error {
	dr.log.Debugw("create Tag", []interface{}{"xid", *input.Xid}...)
	inputData := dgclient.AddTagInput{}
	dr.copyORMStruct(input, &inputData)
	rsp, err := dr.client.CreateTags(ctx, []*dgclient.AddTagInput{&inputData})
	if err != nil {
		return WrapRepoError(err, "failed to create tag").
			Add("tagId", input.ID).Add("tagXid", input.Xid)
	}
	// save ID from response
	input.ID = &rsp.AddTag.Tag[0].ID
	return nil
}

// CreateTags creates new `Tag` objects.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateTags(ctx context.Context, input []*models.Tag) error {
	inputData := make([]*dgclient.AddTagInput, 0, len(input))
	for _, v := range input {
		iv := &dgclient.AddTagInput{}
		dr.copyORMStruct(v, iv)
		inputData = append(inputData, iv)
	}

	dr.log.Debugw("create Tags")
	rsp, err := dr.client.CreateTags(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to create tags")
	}

	// save ID from response
	for i, v := range input {
		v.ID = &rsp.AddTag.Tag[i].ID
	}

	return nil
}

// UpdateTag updates an existing `Tag` object.
func (dr *DgraphRepository) UpdateTag(ctx context.Context, input *models.Tag) error {
	dr.log.Debugw("update Tag", []interface{}{"id", *input.ID, "xid", *input.Xid}...)
	if *input.ID == "" {
		return WrapRepoError(nil, "missing ID").Add("tagXid", input.Xid)
	}
	patch := &dgclient.TagPatch{}
	dr.copyORMStruct(input, patch)
	patch.Xid = nil
	inputData := dgclient.UpdateTagInput{
		Filter: dgclient.TagFilter{
			ID: []string{*input.ID},
		},
		Set: patch,
	}
	_, err := dr.client.UpdateTags(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to update tag").
			Add("tagId", *input.ID).Add("tagXid", input.Xid)
	}
	return nil
}

// DeleteTag deletes a `Tag` object.
func (dr *DgraphRepository) DeleteTag(ctx context.Context, id, xid *string) error {
	delFilter := dgclient.TagFilter{}
	if id != nil && xid != nil {
		return NewRepoError("must specify either id or xid")
	}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &dgclient.StringHashFilter{Eq: xid}
	}

	dr.log.Debugw("delete Tag")
	if _, err := dr.client.DeleteTags(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteTagStr).
			Add("tagId", id).Add("tagXid", xid)
	}
	return nil
}

// DeleteAllTags deletes all `Tag` objects.
func (dr *DgraphRepository) DeleteAllTags(ctx context.Context) error {
	delFilter := dgclient.TagFilter{}
	dr.log.Debugw("delete all Tag")
	if _, err := dr.client.DeleteTags(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteTagStr)
	}
	return nil
}