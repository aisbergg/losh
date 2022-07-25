package dgraph

import (
	"context"
	"losh/internal/core/product/models"
	"losh/internal/lib/errors"
	"losh/internal/repository"
)

var (
	errGetProductStr    = "failed to get product(s)"
	errSaveProductStr   = "failed to save product(s)"
	errDeleteProductStr = "failed to delete product(s)"
)

// GetProduct returns a `Product` object by its ID.
func (dr *DgraphRepository) GetProduct(id, xid *string) (*models.Product, error) {
	ctx := context.Background()
	getProduct, err := dr.client.GetProduct(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetProductStr).
			Add("productId", id).Add("productXid", xid)
	}
	if getProduct.GetProduct == nil { // not found
		return nil, nil
	}
	product := &models.Product{ID: *id}
	if err = dr.dataCopier.CopyTo(getProduct.GetProduct, product); err != nil {
		panic(err)
	}
	return product, nil
}

// GetProducts returns a list of `Product` objects matching the filter criteria.
func (dr *DgraphRepository) GetProducts(filter *models.ProductFilter, order *models.ProductOrder, first *int64, offset *int64) ([]*models.Product, error) {
	ctx := context.Background()
	getProducts, err := dr.client.GetProducts(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetProductStr)
	}
	products := make([]*models.Product, 0, len(getProducts.QueryProduct))
	for _, x := range getProducts.QueryProduct {
		product := &models.Product{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, product); err != nil {
			panic(err)
		}
		products = append(products, product)
	}
	return products, nil
}

// GetAllProducts returns a list of all `Product` objects.
func (dr *DgraphRepository) GetAllProducts() ([]*models.Product, error) {
	return dr.GetProducts(nil, nil, nil, nil)
}

// SaveProduct saves a `Product` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveProduct(product *models.Product) (err error) {
	err = dr.SaveProducts([]*models.Product{product})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("productId", product.ID).Add("productXid", product.Xid)
	}
	return
}

// SaveProducts saves `Product` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveProducts(products []*models.Product) error {
	reqData := make([]*models.AddProductInput, 0, len(products))
	for _, x := range products {
		if x.ID != "" {
			continue
		}
		product := &models.AddProductInput{}
		if err := dr.dataCopier.CopyTo(x, product); err != nil {
			return repository.WrapRepoError(err, errSaveProductStr).
				Add("productId", x.ID).Add("productXid", x.Xid)
		}
		reqData = append(reqData, product)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveProducts(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveProductStr)
	}
	// save ID from response
	for i, x := range products {
		x.ID = respData.AddProduct.Product[i].ID
	}
	return nil
}

// DeleteProduct deletes a `Product` object.
func (dr *DgraphRepository) DeleteProduct(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.ProductFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteProduct(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteProductStr).
			Add("productId", id).Add("productXid", xid)
	}
	return nil
}

// DeleteAllProducts deletes all `Products` objects.
func (dr *DgraphRepository) DeleteAllProducts() error {
	return dr.DeleteProduct(nil, nil)
}

// saveProductIfNecessary saves a `Product` object if it is not already saved.
func (dr *DgraphRepository) saveProductIfNecessary(product *models.Product) (*models.ProductRef, error) {
	if product == nil {
		return nil, nil
	}
	if product.ID == "" {
		if err := dr.SaveProduct(product); err != nil {
			return nil, err
		}
	}
	return &models.ProductRef{ID: &product.ID}, nil
}
