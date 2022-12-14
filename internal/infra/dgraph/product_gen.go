// Code generated by codegen, DO NOT EDIT.

package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net/request"
)

// make sure the struct implements the interface
var _ ProductRepository = (*DgraphRepository)(nil)

// ProductRepository is an interface for getting and saving `Product` objects to a repository.
type ProductRepository interface {
	GetProduct(ctx context.Context, id, xid *string) (*models.Product, error)
	GetProducts(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, int64, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, int64, error)
	CreateProduct(ctx context.Context, input *models.Product) error
	CreateProducts(ctx context.Context, input []*models.Product) error
	UpdateProduct(ctx context.Context, input *models.Product) error
	DeleteProduct(ctx context.Context, id, xid *string) error
	DeleteAllProducts(ctx context.Context) error
}

var (
	errGetProductStr    = "failed to get product(s)"
	errSaveProductStr   = "failed to save product(s)"
	errDeleteProductStr = "failed to delete product(s)"
)

// GetProduct returns a `Product` object by its ID.
func (dr *DgraphRepository) GetProduct(ctx context.Context, id, xid *string) (*models.Product, error) {
	var rspData interface{}
	if id != nil {
		dr.log.Debugw("get Product", "id", *id)
		rsp, err := dr.client.GetProductByID(ctx, *id)
		if err != nil {
			return nil, WrapRepoError(err, errGetProductStr).Add("productId", id)
		}
		rspData = rsp.GetProduct
	} else if xid != nil {
		dr.log.Debugw("get Product", "xid", *xid)
		rsp, err := dr.client.GetProductByXid(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetProductStr).Add("productXid", xid)
		}
		rspData = rsp.GetProduct
	} else {
		panic("must specify id or xid")
	}

	if rspData == nil {
		return nil, nil
	}
	ret := &models.Product{}
	if err := dr.copier.CopyTo(rspData, ret); err != nil {
		panic(err)
	}
	return ret, nil
}

// GetProductID returns the ID of an existing `Product` object.
func (dr *DgraphRepository) GetProductID(ctx context.Context, xid *string) (*string, error) {
	if xid != nil {
		dr.log.Debugw("get Product", "xid", *xid)
		rsp, err := dr.client.GetProductID(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetProductStr).Add("productXid", xid)
		}
		if rsp.GetProduct == nil {
			return nil, nil
		}
		return &rsp.GetProduct.ID, nil
	}

	panic("must specify xid")
}

// GetProducts returns a list of `Product` objects matching the filter criteria.
func (dr *DgraphRepository) GetProducts(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, int64, error) {
	dr.log.Debugw("get Products")
	rsp, err := dr.client.GetProducts(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetProductStr)
	}
	ret := make([]*models.Product, 0, len(rsp.QueryProduct))
	if err = dr.copier.CopyTo(rsp.QueryProduct, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateProduct.Count, nil
}

// GetAllProducts returns a list of all `Product` objects.
func (dr *DgraphRepository) GetAllProducts(ctx context.Context) ([]*models.Product, int64, error) {
	return dr.GetProducts(ctx, nil, nil, nil, nil)
}

// GetProductWithCustomQuery returns a `Product` object by its ID.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetProductWithCustomQuery(ctx context.Context, operationName, query string, id, xid *string) (*models.Product, error) {
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
		Product *models.Product "json:\"getProduct\" graphql:\"getProduct\""
	}{}
	dr.log.Debugw("get Product with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetProductStr)
	}
	return rsp.Product, nil
}

// GetProductsWithCustomQuery returns a list of `Product` objects matching the filter criteria.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetProductsWithCustomQuery(ctx context.Context, operationName, query string, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, error) {
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
		Products []*models.Product "json:\"queryProduct\" graphql:\"queryProduct\""
	}{}
	dr.log.Debugw("get Products with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetProductStr)
	}
	return rsp.Products, nil
}

// GetAllProductsWithCustomQuery returns a list of all `Product` objects.
func (dr *DgraphRepository) GetAllProductsWithCustomQuery(ctx context.Context, operationName, query string) ([]*models.Product, error) {
	return dr.GetProductsWithCustomQuery(ctx, operationName, query, nil, nil, nil, nil)
}

// CreateProduct creates a new `Product` object.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateProduct(ctx context.Context, input *models.Product) error {
	dr.log.Debugw("create Product", []interface{}{"xid", *input.Xid}...)
	inputData := dgclient.AddProductInput{}
	dr.copyORMStruct(input, &inputData)
	rsp, err := dr.client.CreateProducts(ctx, []*dgclient.AddProductInput{&inputData})
	if err != nil {
		return WrapRepoError(err, "failed to create product").
			Add("productId", input.ID).Add("productXid", input.Xid)
	}
	// save ID from response
	input.ID = &rsp.AddProduct.Product[0].ID
	return nil
}

// CreateProducts creates new `Product` objects.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateProducts(ctx context.Context, input []*models.Product) error {
	inputData := make([]*dgclient.AddProductInput, 0, len(input))
	for _, v := range input {
		iv := &dgclient.AddProductInput{}
		dr.copyORMStruct(v, iv)
		inputData = append(inputData, iv)
	}

	dr.log.Debugw("create Products")
	rsp, err := dr.client.CreateProducts(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to create products")
	}

	// save ID from response
	for i, v := range input {
		v.ID = &rsp.AddProduct.Product[i].ID
	}

	return nil
}

// UpdateProduct updates an existing `Product` object.
func (dr *DgraphRepository) UpdateProduct(ctx context.Context, input *models.Product) error {
	dr.log.Debugw("update Product", []interface{}{"id", *input.ID, "xid", *input.Xid}...)
	if *input.ID == "" {
		return WrapRepoError(nil, "missing ID").Add("productXid", input.Xid)
	}
	patch := &dgclient.ProductPatch{}
	dr.copyORMStruct(input, patch)
	patch.Xid = nil
	inputData := dgclient.UpdateProductInput{
		Filter: dgclient.ProductFilter{
			ID: []string{*input.ID},
		},
		Set: patch,
	}
	_, err := dr.client.UpdateProducts(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to update product").
			Add("productId", *input.ID).Add("productXid", input.Xid)
	}
	return nil
}

// DeleteProduct deletes a `Product` object.
func (dr *DgraphRepository) DeleteProduct(ctx context.Context, id, xid *string) error {
	delFilter := dgclient.ProductFilter{}
	if id != nil && xid != nil {
		return NewRepoError("must specify either id or xid")
	}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &dgclient.StringHashFilter{Eq: xid}
	}

	dr.log.Debugw("delete Product")
	if _, err := dr.client.DeleteProducts(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteProductStr).
			Add("productId", id).Add("productXid", xid)
	}
	return nil
}

// DeleteAllProducts deletes all `Product` objects.
func (dr *DgraphRepository) DeleteAllProducts(ctx context.Context) error {
	delFilter := dgclient.ProductFilter{}
	dr.log.Debugw("delete all Product")
	if _, err := dr.client.DeleteProducts(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteProductStr)
	}
	return nil
}
