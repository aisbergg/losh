package dgraph

import (
	"context"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

// GetDatabaseInfo returns the database information. If non exists, it returns
// nil.
func (dr *DgraphRepository) GetDatabaseInfo(id string) (*models.Database, error) {
	ctx := context.Background()
	getDatabaseInfo, err := dr.client.GetDatabaseInfo(ctx)
	if err != nil {
		return nil, repository.WrapRepoError(err, "failed to get database information")
	}
	if len(getDatabaseInfo.QueryDatabase) == 0 {
		return nil, nil
	}
	databaseInfo := &models.Database{}
	if err = copier.CopyWithOption(databaseInfo, getDatabaseInfo.QueryDatabase[0], copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return databaseInfo, nil
}

// SaveDatabaseInfo save the database information. It will make sure to save
// only one object by deleting old ones and creating a new one with the new
// data.
func (dr *DgraphRepository) SaveDatabaseInfo(database *models.Database) (err error) {
	databaseInfo := &models.AddDatabaseInput{}
	if err = copier.CopyWithOption(databaseInfo, database,
		copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	reqData := []*models.AddDatabaseInput{databaseInfo}
	ctx := context.Background()
	respData, err := dr.client.SaveDatabaseInfo(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, "failed to save database information")
	}
	// save ID from response
	database.ID = respData.AddDatabase.Database[0].ID
	return nil
}
