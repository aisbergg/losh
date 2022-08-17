// Code generated by codegen, DO NOT EDIT.

package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net/request"
)

// make sure the struct implements the interface
var _ DatabaseRepository = (*DgraphRepository)(nil)

// DatabaseRepository is an interface for getting and saving `Database` objects to a repository.
type DatabaseRepository interface {
	GetDatabase(ctx context.Context, id *string) (*models.Database, error)
	GetDatabases(ctx context.Context, filter *dgclient.DatabaseFilter, order *dgclient.DatabaseOrder, first *int64, offset *int64) ([]*models.Database, int64, error)
	GetAllDatabases(ctx context.Context) ([]*models.Database, int64, error)
	CreateDatabase(ctx context.Context, input *models.Database) error
	CreateDatabases(ctx context.Context, input []*models.Database) error
	UpdateDatabase(ctx context.Context, input *models.Database) error
	DeleteDatabase(ctx context.Context, id *string) error
	DeleteAllDatabases(ctx context.Context) error
}

var (
	errGetDatabaseStr    = "failed to get database(s)"
	errSaveDatabaseStr   = "failed to save database(s)"
	errDeleteDatabaseStr = "failed to delete database(s)"
)

// GetDatabase returns a `Database` object by its ID.
func (dr *DgraphRepository) GetDatabase(ctx context.Context, id *string) (*models.Database, error) {
	var rspData interface{}
	if id != nil {
		dr.log.Debugw("get Database", "id", *id)
		rsp, err := dr.client.GetDatabaseByID(ctx, *id)
		if err != nil {
			return nil, WrapRepoError(err, errGetDatabaseStr).Add("databaseId", id)
		}
		rspData = rsp.GetDatabase
	} else {
		panic("must specify id")
	}

	if rspData == nil {
		return nil, nil
	}
	ret := &models.Database{}
	if err := dr.copier.CopyTo(rspData, ret); err != nil {
		panic(err)
	}
	return ret, nil
}

// GetDatabases returns a list of `Database` objects matching the filter criteria.
func (dr *DgraphRepository) GetDatabases(ctx context.Context, filter *dgclient.DatabaseFilter, order *dgclient.DatabaseOrder, first *int64, offset *int64) ([]*models.Database, int64, error) {
	dr.log.Debugw("get Databases")
	rsp, err := dr.client.GetDatabases(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetDatabaseStr)
	}
	ret := make([]*models.Database, 0, len(rsp.QueryDatabase))
	if err = dr.copier.CopyTo(rsp.QueryDatabase, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateDatabase.Count, nil
}

// GetAllDatabases returns a list of all `Database` objects.
func (dr *DgraphRepository) GetAllDatabases(ctx context.Context) ([]*models.Database, int64, error) {
	return dr.GetDatabases(ctx, nil, nil, nil, nil)
}

// GetDatabaseWithCustomQuery returns a `Database` object by its ID.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetDatabaseWithCustomQuery(ctx context.Context, operationName, query string, id *string) (*models.Database, error) {
	req := request.GraphQLRequest{
		Ctx:           ctx,
		OperationName: operationName,
		Query:         query,
		Variables: map[string]interface{}{
			"id": id,
		},
	}
	rsp := struct {
		Database *models.Database "json:\"getDatabase\" graphql:\"getDatabase\""
	}{}
	dr.log.Debugw("get Database with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetDatabaseStr)
	}
	return rsp.Database, nil
}

// GetDatabasesWithCustomQuery returns a list of `Database` objects matching the filter criteria.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetDatabasesWithCustomQuery(ctx context.Context, operationName, query string, filter *dgclient.DatabaseFilter, order *dgclient.DatabaseOrder, first *int64, offset *int64) ([]*models.Database, error) {
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
		Databases []*models.Database "json:\"queryDatabase\" graphql:\"queryDatabase\""
	}{}
	dr.log.Debugw("get Databases with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetDatabaseStr)
	}
	return rsp.Databases, nil
}

// GetAllDatabasesWithCustomQuery returns a list of all `Database` objects.
func (dr *DgraphRepository) GetAllDatabasesWithCustomQuery(ctx context.Context, operationName, query string) ([]*models.Database, error) {
	return dr.GetDatabasesWithCustomQuery(ctx, operationName, query, nil, nil, nil, nil)
}

// CreateDatabase creates a new `Database` object.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateDatabase(ctx context.Context, input *models.Database) error {
	dr.log.Debugw("create Database", []interface{}{}...)
	inputData := dgclient.AddDatabaseInput{}
	dr.copyORMStruct(input, &inputData)
	rsp, err := dr.client.CreateDatabases(ctx, []*dgclient.AddDatabaseInput{&inputData})
	if err != nil {
		return WrapRepoError(err, "failed to create database").
			Add("databaseId", input.ID)
	}
	// save ID from response
	input.ID = &rsp.AddDatabase.Database[0].ID
	return nil
}

// CreateDatabases creates new `Database` objects.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateDatabases(ctx context.Context, input []*models.Database) error {
	inputData := make([]*dgclient.AddDatabaseInput, 0, len(input))
	for _, v := range input {
		iv := &dgclient.AddDatabaseInput{}
		dr.copyORMStruct(v, iv)
		inputData = append(inputData, iv)
	}

	dr.log.Debugw("create Databases")
	rsp, err := dr.client.CreateDatabases(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to create databases")
	}

	// save ID from response
	for i, v := range input {
		v.ID = &rsp.AddDatabase.Database[i].ID
	}

	return nil
}

// UpdateDatabase updates an existing `Database` object.
func (dr *DgraphRepository) UpdateDatabase(ctx context.Context, input *models.Database) error {
	dr.log.Debugw("update Database", []interface{}{"id", *input.ID}...)
	if *input.ID == "" {
		return WrapRepoError(nil, "missing ID")
	}
	patch := &dgclient.DatabasePatch{}
	dr.copyORMStruct(input, patch)
	inputData := dgclient.UpdateDatabaseInput{
		Filter: dgclient.DatabaseFilter{
			ID: []string{*input.ID},
		},
		Set: patch,
	}
	_, err := dr.client.UpdateDatabases(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to update database").
			Add("databaseId", *input.ID)
	}
	return nil
}

// DeleteDatabase deletes a `Database` object.
func (dr *DgraphRepository) DeleteDatabase(ctx context.Context, id *string) error {
	delFilter := dgclient.DatabaseFilter{}
	if id == nil {
		return NewRepoError("must specify id")
	}
	delFilter.ID = []string{*id}

	dr.log.Debugw("delete Database")
	if _, err := dr.client.DeleteDatabases(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteDatabaseStr).
			Add("databaseId", id)
	}
	return nil
}

// DeleteAllDatabases deletes all `Database` objects.
func (dr *DgraphRepository) DeleteAllDatabases(ctx context.Context) error {
	delFilter := dgclient.DatabaseFilter{}
	dr.log.Debugw("delete all Database")
	if _, err := dr.client.DeleteDatabases(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteDatabaseStr)
	}
	return nil
}