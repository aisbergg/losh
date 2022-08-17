package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
)

// SearchProducts returns a list of `Product` objects matching the filter criteria.
func (dr *DgraphRepository) SearchProducts(ctx context.Context, filter *dgclient.ProductFilter, order *dgclient.ProductOrder, first *int64, offset *int64) ([]*models.Product, int64, error) {
	dr.log.Debugw("search Products")
	rsp, err := dr.client.SearchProducts(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetProductStr)
	}
	ret := make([]*models.Product, 0, len(rsp.QueryProduct))
	if err = dr.copier.CopyTo(rsp.QueryProduct, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateProduct.Count, nil
}
