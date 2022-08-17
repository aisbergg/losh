// Code generated by codegen, DO NOT EDIT.

package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net/request"
)

// make sure the struct implements the interface
var _ MaterialRepository = (*DgraphRepository)(nil)

// MaterialRepository is an interface for getting and saving `Material` objects to a repository.
type MaterialRepository interface {
	GetMaterial(ctx context.Context, id *string) (*models.Material, error)
	GetMaterials(ctx context.Context, filter *dgclient.MaterialFilter, order *dgclient.MaterialOrder, first *int64, offset *int64) ([]*models.Material, int64, error)
	GetAllMaterials(ctx context.Context) ([]*models.Material, int64, error)
	CreateMaterial(ctx context.Context, input *models.Material) error
	CreateMaterials(ctx context.Context, input []*models.Material) error
	UpdateMaterial(ctx context.Context, input *models.Material) error
	DeleteMaterial(ctx context.Context, id *string) error
	DeleteAllMaterials(ctx context.Context) error
}

var (
	errGetMaterialStr    = "failed to get material(s)"
	errSaveMaterialStr   = "failed to save material(s)"
	errDeleteMaterialStr = "failed to delete material(s)"
)

// GetMaterial returns a `Material` object by its ID.
func (dr *DgraphRepository) GetMaterial(ctx context.Context, id *string) (*models.Material, error) {
	var rspData interface{}
	if id != nil {
		dr.log.Debugw("get Material", "id", *id)
		rsp, err := dr.client.GetMaterialByID(ctx, *id)
		if err != nil {
			return nil, WrapRepoError(err, errGetMaterialStr).Add("materialId", id)
		}
		rspData = rsp.GetMaterial
	} else {
		panic("must specify id")
	}

	if rspData == nil {
		return nil, nil
	}
	ret := &models.Material{}
	if err := dr.copier.CopyTo(rspData, ret); err != nil {
		panic(err)
	}
	return ret, nil
}

// GetMaterials returns a list of `Material` objects matching the filter criteria.
func (dr *DgraphRepository) GetMaterials(ctx context.Context, filter *dgclient.MaterialFilter, order *dgclient.MaterialOrder, first *int64, offset *int64) ([]*models.Material, int64, error) {
	dr.log.Debugw("get Materials")
	rsp, err := dr.client.GetMaterials(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetMaterialStr)
	}
	ret := make([]*models.Material, 0, len(rsp.QueryMaterial))
	if err = dr.copier.CopyTo(rsp.QueryMaterial, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateMaterial.Count, nil
}

// GetAllMaterials returns a list of all `Material` objects.
func (dr *DgraphRepository) GetAllMaterials(ctx context.Context) ([]*models.Material, int64, error) {
	return dr.GetMaterials(ctx, nil, nil, nil, nil)
}

// GetMaterialWithCustomQuery returns a `Material` object by its ID.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetMaterialWithCustomQuery(ctx context.Context, operationName, query string, id *string) (*models.Material, error) {
	req := request.GraphQLRequest{
		Ctx:           ctx,
		OperationName: operationName,
		Query:         query,
		Variables: map[string]interface{}{
			"id": id,
		},
	}
	rsp := struct {
		Material *models.Material "json:\"getMaterial\" graphql:\"getMaterial\""
	}{}
	dr.log.Debugw("get Material with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetMaterialStr)
	}
	return rsp.Material, nil
}

// GetMaterialsWithCustomQuery returns a list of `Material` objects matching the filter criteria.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetMaterialsWithCustomQuery(ctx context.Context, operationName, query string, filter *dgclient.MaterialFilter, order *dgclient.MaterialOrder, first *int64, offset *int64) ([]*models.Material, error) {
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
		Materials []*models.Material "json:\"queryMaterial\" graphql:\"queryMaterial\""
	}{}
	dr.log.Debugw("get Materials with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetMaterialStr)
	}
	return rsp.Materials, nil
}

// GetAllMaterialsWithCustomQuery returns a list of all `Material` objects.
func (dr *DgraphRepository) GetAllMaterialsWithCustomQuery(ctx context.Context, operationName, query string) ([]*models.Material, error) {
	return dr.GetMaterialsWithCustomQuery(ctx, operationName, query, nil, nil, nil, nil)
}

// CreateMaterial creates a new `Material` object.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateMaterial(ctx context.Context, input *models.Material) error {
	dr.log.Debugw("create Material", []interface{}{}...)
	inputData := dgclient.AddMaterialInput{}
	dr.copyORMStruct(input, &inputData)
	rsp, err := dr.client.CreateMaterials(ctx, []*dgclient.AddMaterialInput{&inputData})
	if err != nil {
		return WrapRepoError(err, "failed to create material").
			Add("materialId", input.ID)
	}
	// save ID from response
	input.ID = &rsp.AddMaterial.Material[0].ID
	return nil
}

// CreateMaterials creates new `Material` objects.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateMaterials(ctx context.Context, input []*models.Material) error {
	inputData := make([]*dgclient.AddMaterialInput, 0, len(input))
	for _, v := range input {
		iv := &dgclient.AddMaterialInput{}
		dr.copyORMStruct(v, iv)
		inputData = append(inputData, iv)
	}

	dr.log.Debugw("create Materials")
	rsp, err := dr.client.CreateMaterials(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to create materials")
	}

	// save ID from response
	for i, v := range input {
		v.ID = &rsp.AddMaterial.Material[i].ID
	}

	return nil
}

// UpdateMaterial updates an existing `Material` object.
func (dr *DgraphRepository) UpdateMaterial(ctx context.Context, input *models.Material) error {
	dr.log.Debugw("update Material", []interface{}{"id", *input.ID}...)
	if *input.ID == "" {
		return WrapRepoError(nil, "missing ID")
	}
	patch := &dgclient.MaterialPatch{}
	dr.copyORMStruct(input, patch)
	inputData := dgclient.UpdateMaterialInput{
		Filter: dgclient.MaterialFilter{
			ID: []string{*input.ID},
		},
		Set: patch,
	}
	_, err := dr.client.UpdateMaterials(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to update material").
			Add("materialId", *input.ID)
	}
	return nil
}

// DeleteMaterial deletes a `Material` object.
func (dr *DgraphRepository) DeleteMaterial(ctx context.Context, id *string) error {
	delFilter := dgclient.MaterialFilter{}
	if id == nil {
		return NewRepoError("must specify id")
	}
	delFilter.ID = []string{*id}

	dr.log.Debugw("delete Material")
	if _, err := dr.client.DeleteMaterials(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteMaterialStr).
			Add("materialId", id)
	}
	return nil
}

// DeleteAllMaterials deletes all `Material` objects.
func (dr *DgraphRepository) DeleteAllMaterials(ctx context.Context) error {
	delFilter := dgclient.MaterialFilter{}
	dr.log.Debugw("delete all Material")
	if _, err := dr.client.DeleteMaterials(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteMaterialStr)
	}
	return nil
}
