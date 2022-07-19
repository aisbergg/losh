package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetCategoryStr    = "failed to get category(s)"
	errSaveCategoryStr   = "failed to save category(s)"
	errDeleteCategoryStr = "failed to delete category(s)"
)

// GetCategory returns a `Category` object by its ID.
func (dr *DgraphRepository) GetCategory(id, xid *string) (*models.Category, error) {
	ctx := context.Background()
	getCategory, err := dr.client.GetCategory(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetCategoryStr).
			Add("categoryId", id).Add("categoryXid", xid)
	}
	if getCategory.GetCategory == nil { // not found
		return nil, nil
	}
	category := &models.Category{ID: *id}
	if err = dr.dataCopier.CopyTo(getCategory.GetCategory, category); err != nil {
		panic(err)
	}
	return category, nil
}

// GetCategories returns a list of `Category` objects matching the filter criteria.
func (dr *DgraphRepository) GetCategories(filter *models.CategoryFilter, order *models.CategoryOrder, first *int64, offset *int64) ([]*models.Category, error) {
	ctx := context.Background()
	getCategories, err := dr.client.GetCategories(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetCategoryStr)
	}
	categories := make([]*models.Category, 0, len(getCategories.QueryCategory))
	for _, x := range getCategories.QueryCategory {
		category := &models.Category{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, category); err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetAllCategories returns a list of all `Category` objects.
func (dr *DgraphRepository) GetAllCategories() ([]*models.Category, error) {
	return dr.GetCategories(nil, nil, nil, nil)
}

// SaveCategory saves a `Category` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveCategory(category *models.Category) (err error) {
	err = dr.SaveCategories([]*models.Category{category})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("categoryId", category.ID).Add("categoryXid", category.Xid)
	}
	return
}

// SaveCategories saves `Category` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveCategories(categories []*models.Category) error {
	reqData := make([]*models.AddCategoryInput, 0, len(categories))
	for _, x := range categories {
		if x.ID != "" {
			continue
		}
		category := &models.AddCategoryInput{}
		if err := dr.dataCopier.CopyTo(x, category); err != nil {
			return repository.WrapRepoError(err, errSaveCategoryStr).
				Add("categoryId", x.ID).Add("categoryXid", x.Xid)
		}
		reqData = append(reqData, category)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveCategories(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveCategoryStr)
	}
	// save ID from response
	for i, x := range categories {
		x.ID = respData.AddCategory.Category[i].ID
	}
	return nil
}

// DeleteCategory deletes a `Category` object.
func (dr *DgraphRepository) DeleteCategory(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.CategoryFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteCategory(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteCategoryStr).
			Add("categoryId", id).Add("categoryXid", xid)
	}
	return nil
}

// DeleteAllCategories deletes all `Category` objects.
func (dr *DgraphRepository) DeleteAllCategories() error {
	return dr.DeleteCategory(nil, nil)
}

// saveCategoryIfNecessary saves a `Category` object if it is not already saved.
func (dr *DgraphRepository) saveCategoryIfNecessary(category *models.Category) (*models.CategoryRef, error) {
	if category == nil {
		return nil, nil
	}
	if category.ID == "" {
		if err := dr.SaveCategory(category); err != nil {
			return nil, err
		}
	}
	return &models.CategoryRef{ID: &category.ID}, nil
}
