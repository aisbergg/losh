package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetMaterialStr    = "failed to get material(s)"
	errSaveMaterialStr   = "failed to save material(s)"
	errDeleteMaterialStr = "failed to delete material(s)"
)

// GetMaterial returns a `Material` object by its ID.
func (dr *DgraphRepository) GetMaterial(id string) (*models.Material, error) {
	ctx := context.Background()
	getMaterial, err := dr.client.GetMaterial(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetMaterialStr).
			Add("materialId", id)
	}
	if getMaterial.GetMaterial == nil { // not found
		return nil, nil
	}
	material := &models.Material{ID: id}
	if err = dr.dataCopier.CopyTo(getMaterial.GetMaterial, material); err != nil {
		panic(err)
	}
	return material, nil
}

// GetMaterials returns a list of `Material` objects matching the filter criteria.
func (dr *DgraphRepository) GetMaterials(filter *models.MaterialFilter, order *models.MaterialOrder, first *int64, offset *int64) ([]*models.Material, error) {
	ctx := context.Background()
	getMaterials, err := dr.client.GetMaterials(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetMaterialStr)
	}
	materials := make([]*models.Material, 0, len(getMaterials.QueryMaterial))
	for _, x := range getMaterials.QueryMaterial {
		material := &models.Material{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, material); err != nil {
			panic(err)
		}
		materials = append(materials, material)
	}
	return materials, nil
}

// GetAllMaterials returns a list of all `Material` objects.
func (dr *DgraphRepository) GetAllMaterials() ([]*models.Material, error) {
	return dr.GetMaterials(nil, nil, nil, nil)
}

// SaveMaterial saves a `Material` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveMaterial(material *models.Material) (err error) {
	err = dr.SaveMaterials([]*models.Material{material})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("materialId", material.ID)
	}
	return
}

// SaveMaterials saves `Material` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveMaterials(materials []*models.Material) error {
	reqData := make([]*models.AddMaterialInput, 0, len(materials))
	for _, x := range materials {
		if x.ID != "" {
			continue
		}
		material := &models.AddMaterialInput{}
		if err := dr.dataCopier.CopyTo(x, material); err != nil {
			return repository.WrapRepoError(err, errSaveMaterialStr).
				Add("materialId", x.ID)
		}
		reqData = append(reqData, material)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveMaterials(ctx, reqData, []string{})
	if err != nil {
		return repository.WrapRepoError(err, errSaveMaterialStr)
	}
	// save ID from response
	for i, x := range materials {
		x.ID = respData.AddMaterial.Material[i].ID
	}
	return nil
}

// DeleteMaterial deletes a `Material` object.
func (dr *DgraphRepository) DeleteMaterial(id *string) error {
	ctx := context.Background()
	delFilter := models.MaterialFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteMaterials(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteMaterialStr).
			Add("materialId", id)
	}
	return nil
}

// DeleteAllMaterials deletes all `Materials` objects.
func (dr *DgraphRepository) DeleteAllMaterials() error {
	return dr.DeleteMaterial(nil)
}

// saveMaterialIfNecessary saves a `Material` object if it is not already saved.
func (dr *DgraphRepository) saveMaterialIfNecessary(material *models.Material) (*models.MaterialRef, error) {
	if material == nil {
		return nil, nil
	}
	if material.ID == "" {
		if err := dr.SaveMaterial(material); err != nil {
			return nil, err
		}
	}
	return &models.MaterialRef{ID: &material.ID}, nil
}
