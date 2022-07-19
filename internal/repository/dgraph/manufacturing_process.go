package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetManufacturingProcessStr    = "failed to get manufacturingProcess(s)"
	errSaveManufacturingProcessStr   = "failed to save manufacturingProcess(s)"
	errDeleteManufacturingProcessStr = "failed to delete manufacturingProcess(s)"
)

// GetManufacturingProcess returns a `ManufacturingProcess` object by its ID.
func (dr *DgraphRepository) GetManufacturingProcess(id string) (*models.ManufacturingProcess, error) {
	ctx := context.Background()
	getManufacturingProcess, err := dr.client.GetManufacturingProcess(ctx, id)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetManufacturingProcessStr).
			Add("manufacturingProcessId", id)
	}
	if getManufacturingProcess.GetManufacturingProcess == nil { // not found
		return nil, nil
	}
	manufacturingProcess := &models.ManufacturingProcess{ID: id}
	if err = dr.dataCopier.CopyTo(getManufacturingProcess.GetManufacturingProcess, manufacturingProcess); err != nil {
		panic(err)
	}
	return manufacturingProcess, nil
}

// GetManufacturingProcesss returns a list of `ManufacturingProcess` objects matching the filter criteria.
func (dr *DgraphRepository) GetManufacturingProcesss(filter *models.ManufacturingProcessFilter, order *models.ManufacturingProcessOrder, first *int64, offset *int64) ([]*models.ManufacturingProcess, error) {
	ctx := context.Background()
	getManufacturingProcesss, err := dr.client.GetManufacturingProcesses(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetManufacturingProcessStr)
	}
	manufacturingProcesss := make([]*models.ManufacturingProcess, 0, len(getManufacturingProcesss.QueryManufacturingProcess))
	for _, x := range getManufacturingProcesss.QueryManufacturingProcess {
		manufacturingProcess := &models.ManufacturingProcess{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, manufacturingProcess); err != nil {
			panic(err)
		}
		manufacturingProcesss = append(manufacturingProcesss, manufacturingProcess)
	}
	return manufacturingProcesss, nil
}

// GetAllManufacturingProcesss returns a list of all `ManufacturingProcess` objects.
func (dr *DgraphRepository) GetAllManufacturingProcesss() ([]*models.ManufacturingProcess, error) {
	return dr.GetManufacturingProcesss(nil, nil, nil, nil)
}

// SaveManufacturingProcess saves a `ManufacturingProcess` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveManufacturingProcess(manufacturingProcess *models.ManufacturingProcess) (err error) {
	err = dr.SaveManufacturingProcesss([]*models.ManufacturingProcess{manufacturingProcess})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("manufacturingProcessId", manufacturingProcess.ID)
	}
	return
}

// SaveManufacturingProcesss saves `ManufacturingProcess` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveManufacturingProcesss(manufacturingProcesss []*models.ManufacturingProcess) error {
	reqData := make([]*models.AddManufacturingProcessInput, 0, len(manufacturingProcesss))
	for _, x := range manufacturingProcesss {
		if x.ID != "" {
			continue
		}
		manufacturingProcess := &models.AddManufacturingProcessInput{}
		if err := dr.dataCopier.CopyTo(x, manufacturingProcess); err != nil {
			return repository.WrapRepoError(err, errSaveManufacturingProcessStr).
				Add("manufacturingProcessId", x.ID)
		}
		reqData = append(reqData, manufacturingProcess)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveManufacturingProcesses(ctx, reqData, []string{})
	if err != nil {
		return repository.WrapRepoError(err, errSaveManufacturingProcessStr)
	}
	// save ID from response
	for i, x := range manufacturingProcesss {
		x.ID = respData.AddManufacturingProcess.ManufacturingProcess[i].ID
	}
	return nil
}

// DeleteManufacturingProcess deletes a `ManufacturingProcess` object.
func (dr *DgraphRepository) DeleteManufacturingProcess(id *string) error {
	ctx := context.Background()
	delFilter := models.ManufacturingProcessFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteManufacturingProcesses(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteManufacturingProcessStr).
			Add("manufacturingProcessId", id)
	}
	return nil
}

// DeleteAllManufacturingProcesss deletes all `ManufacturingProcesss` objects.
func (dr *DgraphRepository) DeleteAllManufacturingProcesss() error {
	return dr.DeleteManufacturingProcess(nil)
}

// saveManufacturingProcessIfNecessary saves a `ManufacturingProcess` object if it is not already saved.
func (dr *DgraphRepository) saveManufacturingProcessIfNecessary(manufacturingProcess *models.ManufacturingProcess) (*models.ManufacturingProcessRef, error) {
	if manufacturingProcess == nil {
		return nil, nil
	}
	if manufacturingProcess.ID == "" {
		if err := dr.SaveManufacturingProcess(manufacturingProcess); err != nil {
			return nil, err
		}
	}
	return &models.ManufacturingProcessRef{ID: &manufacturingProcess.ID}, nil
}
