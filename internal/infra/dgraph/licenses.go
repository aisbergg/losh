package dgraph

import (
	"context"

	"losh/internal/core/product/models"
)

// GetAllLicensesBasic returns a list of all `License` objects with basic
// information.
func (dr *DgraphRepository) GetAllLicensesBasic(ctx context.Context) ([]*models.License, error) {
	rsp := struct {
		Licenses []*models.License "json:\"queryLicense\" graphql:\"queryLicense\""
	}{}
	if err := dr.client.GetAllLicensesBasicWithResponse(ctx, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetLicenseStr)
	}
	return rsp.Licenses, nil
}
