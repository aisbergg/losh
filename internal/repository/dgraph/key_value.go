package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetKeyValueStr    = "failed to get keyValue(s)"
	errSaveKeyValueStr   = "failed to save keyValue(s)"
	errDeleteKeyValueStr = "failed to delete keyValue(s)"
)

// GetKeyValue returns a `KeyValue` object by its ID.
func (dr *DgraphRepository) GetKeyValue(id string) (*models.KeyValue, error) {
	ctx := context.Background()
	getKeyValue, err := dr.client.GetKeyValue(ctx, id)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetKeyValueStr).
			AddIfNotNil("keyValueId", id)
	}
	keyValue := &models.KeyValue{ID: id}
	if err = copier.CopyWithOption(keyValue, getKeyValue.GetKeyValue, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return keyValue, nil
}

// GetKeyValues returns a list of `KeyValue` objects matching the filter criteria.
func (dr *DgraphRepository) GetKeyValues(filter *models.KeyValueFilter, order *models.KeyValueOrder, first *int64, offset *int64) ([]*models.KeyValue, error) {
	ctx := context.Background()
	getKeyValues, err := dr.client.GetKeyValues(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetKeyValueStr)
	}
	keyValues := make([]*models.KeyValue, 0, len(getKeyValues.QueryKeyValue))
	for _, x := range getKeyValues.QueryKeyValue {
		keyValue := &models.KeyValue{ID: x.ID}
		if err = copier.CopyWithOption(keyValue, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		keyValues = append(keyValues, keyValue)
	}
	return keyValues, nil
}

// GetAllKeyValues returns a list of all `KeyValue` objects.
func (dr *DgraphRepository) GetAllKeyValues() ([]*models.KeyValue, error) {
	return dr.GetKeyValues(nil, nil, nil, nil)
}

// SaveKeyValue saves a `KeyValue` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveKeyValue(keyValue *models.KeyValue) (err error) {
	err = dr.SaveKeyValues([]*models.KeyValue{keyValue})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.AddIfNotNil("keyValueId", keyValue.ID)
	}
	return
}

// SaveKeyValues saves `KeyValue` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveKeyValues(keyValues []*models.KeyValue) error {
	reqData := make([]*models.AddKeyValueInput, 0, len(keyValues))
	for _, x := range keyValues {
		if x.ID != "" {
			continue
		}
		keyValue := &models.AddKeyValueInput{}
		if err := copier.CopyWithOption(keyValue, x,
			copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
			return repository.NewRepoErrorWrap(err, errSaveKeyValueStr).
				AddIfNotNil("keyValueId", x.ID)
		}
		reqData = append(reqData, keyValue)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveKeyValues(ctx, reqData)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errSaveKeyValueStr)
	}
	// save ID from response
	for i, x := range keyValues {
		x.ID = respData.AddKeyValue.KeyValue[i].ID
	}
	return nil
}

// DeleteKeyValue deletes a `KeyValue` object.
func (dr *DgraphRepository) DeleteKeyValue(id *string) error {
	ctx := context.Background()
	delFilter := models.KeyValueFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteKeyValue(ctx, delFilter)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errDeleteKeyValueStr).
			AddIfNotNil("keyValueId", id)
	}
	return nil
}

// DeleteAllKeyValues deletes all `KeyValues` objects.
func (dr *DgraphRepository) DeleteAllKeyValues() error {
	return dr.DeleteKeyValue(nil)
}

// saveKeyValueIfNecessary saves a `KeyValue` object if it is not already saved.
func (dr *DgraphRepository) saveKeyValueIfNecessary(keyValue *models.KeyValue) (*models.KeyValueRef, error) {
	if keyValue == nil {
		return nil, nil
	}
	if keyValue.ID == "" {
		if err := dr.SaveKeyValue(keyValue); err != nil {
			return nil, err
		}
	}
	return &models.KeyValueRef{ID: &keyValue.ID}, nil
}
