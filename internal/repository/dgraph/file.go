package dgraph

import (
	"context"
	"losh/crawler/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetFileStr    = "failed to get file(s)"
	errSaveFileStr   = "failed to save file(s)"
	errDeleteFileStr = "failed to delete file(s)"
)

// GetFile returns a `File` object by its ID.
func (dr *DgraphRepository) GetFile(id string) (*models.File, error) {
	ctx := context.Background()
	getFile, err := dr.client.GetFile(ctx, id)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetFileStr).
			AddIfNotNil("fileId", id)
	}
	file := &models.File{ID: id}
	if err = copier.CopyWithOption(file, getFile.GetFile, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return file, nil
}

// GetFiles returns a list of `File` objects matching the filter criteria.
func (dr *DgraphRepository) GetFiles(filter *models.FileFilter, order *models.FileOrder, first *int64, offset *int64) ([]*models.File, error) {
	ctx := context.Background()
	getFiles, err := dr.client.GetFiles(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.NewRepoErrorWrap(err, errGetFileStr)
	}
	files := make([]*models.File, 0, len(getFiles.QueryFile))
	for _, x := range getFiles.QueryFile {
		file := &models.File{ID: x.ID}
		if err = copier.CopyWithOption(file, x, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		files = append(files, file)
	}
	return files, nil
}

// GetAllFiles returns a list of all `File` objects.
func (dr *DgraphRepository) GetAllFiles() ([]*models.File, error) {
	return dr.GetFiles(nil, nil, nil, nil)
}

// SaveFile saves a `File` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveFile(file *models.File) (err error) {
	err = dr.SaveFiles([]*models.File{file})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.AddIfNotNil("fileId", file.ID)
	}
	return
}

// SaveFiles saves `File` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveFiles(files []*models.File) error {
	reqData := make([]*models.AddFileInput, 0, len(files))
	for _, x := range files {
		if x.ID != "" {
			continue
		}
		file := &models.AddFileInput{}
		if err := copier.CopyWithOption(file, x,
			copier.Option{Converters: dr.convertersForSave, DeepCopy: true, IgnoreEmpty: true}); err != nil {
			return repository.NewRepoErrorWrap(err, errSaveFileStr).
				AddIfNotNil("fileId", x.ID)
		}
		reqData = append(reqData, file)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveFiles(ctx, reqData, []string{})
	if err != nil {
		return repository.NewRepoErrorWrap(err, errSaveFileStr)
	}
	// save ID from response
	for i, x := range files {
		x.ID = respData.AddFile.File[i].ID
	}
	return nil
}

// DeleteFile deletes a `File` object.
func (dr *DgraphRepository) DeleteFile(id *string) error {
	ctx := context.Background()
	delFilter := models.FileFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	_, err := dr.client.DeleteFile(ctx, delFilter)
	if err != nil {
		return repository.NewRepoErrorWrap(err, errDeleteFileStr).
			AddIfNotNil("fileId", id)
	}
	return nil
}

// DeleteAllFiles deletes all `Files` objects.
func (dr *DgraphRepository) DeleteAllFiles() error {
	return dr.DeleteFile(nil)
}

// saveFileIfNecessary saves a `File` object if it is not already saved.
func (dr *DgraphRepository) saveFileIfNecessary(file *models.File) (*models.FileRef, error) {
	if file == nil {
		return nil, nil
	}
	if file.ID == "" {
		if err := dr.SaveFile(file); err != nil {
			return nil, err
		}
	}
	return &models.FileRef{ID: &file.ID}, nil
}
