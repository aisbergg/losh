package dgraph

import (
	"context"
	"losh/internal/core/product/models"
	"losh/internal/lib/errors"
	"losh/internal/repository"
)

var (
	errGetTagStr    = "failed to get tag(s)"
	errSaveTagStr   = "failed to save tag(s)"
	errDeleteTagStr = "failed to delete tag(s)"
)

// GetTag returns a `Tag` object by its ID.
func (dr *DgraphRepository) GetTag(id, xid *string) (*models.Tag, error) {
	ctx := context.Background()
	getTag, err := dr.client.GetTag(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetTagStr).
			Add("tagId", id).Add("tagXid", xid)
	}
	if getTag.GetTag == nil { // not found
		return nil, nil
	}
	tag := &models.Tag{ID: *id}
	if err = dr.dataCopier.CopyTo(getTag.GetTag, tag); err != nil {
		panic(err)
	}
	return tag, nil
}

// GetTags returns a list of `Tag` objects matching the filter criteria.
func (dr *DgraphRepository) GetTags(filter *models.TagFilter, order *models.TagOrder, first *int64, offset *int64) ([]*models.Tag, error) {
	ctx := context.Background()
	getTags, err := dr.client.GetTags(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetTagStr)
	}
	tags := make([]*models.Tag, 0, len(getTags.QueryTag))
	for _, x := range getTags.QueryTag {
		tag := &models.Tag{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, tag); err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// GetAllTags returns a list of all `Tag` objects.
func (dr *DgraphRepository) GetAllTags() ([]*models.Tag, error) {
	return dr.GetTags(nil, nil, nil, nil)
}

// SaveTag saves a `Tag` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveTag(tag *models.Tag) (err error) {
	err = dr.SaveTags([]*models.Tag{tag})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("tagId", tag.ID).Add("tagXid", tag.Xid)
	}
	return
}

// SaveTags saves `Tag` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveTags(tags []*models.Tag) error {
	reqData := make([]*models.AddTagInput, 0, len(tags))
	for _, x := range tags {
		if x.ID != "" {
			continue
		}
		tag := &models.AddTagInput{}
		if err := dr.dataCopier.CopyTo(x, tag); err != nil {
			return repository.WrapRepoError(err, errSaveTagStr).
				Add("tagId", x.ID).Add("tagXid", x.Xid)
		}
		reqData = append(reqData, tag)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveTags(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveTagStr)
	}
	// save ID from response
	for i, x := range tags {
		x.ID = respData.AddTag.Tag[i].ID
	}
	return nil
}

// DeleteTag deletes a `Tag` object.
func (dr *DgraphRepository) DeleteTag(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.TagFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteTag(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteTagStr).
			Add("tagId", id).Add("tagXid", xid)
	}
	return nil
}

// DeleteAllTags deletes all `Tags` objects.
func (dr *DgraphRepository) DeleteAllTags() error {
	return dr.DeleteTag(nil, nil)
}

// saveTagIfNecessary saves a `Tag` object if it is not already saved.
func (dr *DgraphRepository) saveTagIfNecessary(tag *models.Tag) (*models.TagRef, error) {
	if tag == nil {
		return nil, nil
	}
	if tag.ID == "" {
		if err := dr.SaveTag(tag); err != nil {
			return nil, err
		}
	}
	return &models.TagRef{ID: &tag.ID}, nil
}
