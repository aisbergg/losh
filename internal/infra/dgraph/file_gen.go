// Code generated by codegen, DO NOT EDIT.

package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net/request"
)

// make sure the struct implements the interface
var _ FileRepository = (*DgraphRepository)(nil)

// FileRepository is an interface for getting and saving `File` objects to a repository.
type FileRepository interface {
	GetFile(ctx context.Context, id, xid *string) (*models.File, error)
	GetFiles(ctx context.Context, filter *dgclient.FileFilter, order *dgclient.FileOrder, first *int64, offset *int64) ([]*models.File, int64, error)
	GetAllFiles(ctx context.Context) ([]*models.File, int64, error)
	CreateFile(ctx context.Context, input *models.File) error
	CreateFiles(ctx context.Context, input []*models.File) error
	UpdateFile(ctx context.Context, input *models.File) error
	DeleteFile(ctx context.Context, id, xid *string) error
	DeleteAllFiles(ctx context.Context) error
}

var (
	errGetFileStr    = "failed to get file(s)"
	errSaveFileStr   = "failed to save file(s)"
	errDeleteFileStr = "failed to delete file(s)"
)

// GetFile returns a `File` object by its ID.
func (dr *DgraphRepository) GetFile(ctx context.Context, id, xid *string) (*models.File, error) {
	var rspData interface{}
	if id != nil {
		dr.log.Debugw("get File", "id", *id)
		rsp, err := dr.client.GetFileByID(ctx, *id)
		if err != nil {
			return nil, WrapRepoError(err, errGetFileStr).Add("fileId", id)
		}
		rspData = rsp.GetFile
	} else if xid != nil {
		dr.log.Debugw("get File", "xid", *xid)
		rsp, err := dr.client.GetFileByXid(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetFileStr).Add("fileXid", xid)
		}
		rspData = rsp.GetFile
	} else {
		panic("must specify id or xid")
	}

	if rspData == nil {
		return nil, nil
	}
	ret := &models.File{}
	if err := dr.copier.CopyTo(rspData, ret); err != nil {
		panic(err)
	}
	return ret, nil
}

// GetFileID returns the ID of an existing `File` object.
func (dr *DgraphRepository) GetFileID(ctx context.Context, xid *string) (*string, error) {
	if xid != nil {
		dr.log.Debugw("get File", "xid", *xid)
		rsp, err := dr.client.GetFileID(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetFileStr).Add("fileXid", xid)
		}
		if rsp.GetFile == nil {
			return nil, nil
		}
		return &rsp.GetFile.ID, nil
	}

	panic("must specify xid")
}

// GetFiles returns a list of `File` objects matching the filter criteria.
func (dr *DgraphRepository) GetFiles(ctx context.Context, filter *dgclient.FileFilter, order *dgclient.FileOrder, first *int64, offset *int64) ([]*models.File, int64, error) {
	dr.log.Debugw("get Files")
	rsp, err := dr.client.GetFiles(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetFileStr)
	}
	ret := make([]*models.File, 0, len(rsp.QueryFile))
	if err = dr.copier.CopyTo(rsp.QueryFile, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateFile.Count, nil
}

// GetAllFiles returns a list of all `File` objects.
func (dr *DgraphRepository) GetAllFiles(ctx context.Context) ([]*models.File, int64, error) {
	return dr.GetFiles(ctx, nil, nil, nil, nil)
}

// GetFileWithCustomQuery returns a `File` object by its ID.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetFileWithCustomQuery(ctx context.Context, operationName, query string, id, xid *string) (*models.File, error) {
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
		File *models.File "json:\"getFile\" graphql:\"getFile\""
	}{}
	dr.log.Debugw("get File with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetFileStr)
	}
	return rsp.File, nil
}

// GetFilesWithCustomQuery returns a list of `File` objects matching the filter criteria.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetFilesWithCustomQuery(ctx context.Context, operationName, query string, filter *dgclient.FileFilter, order *dgclient.FileOrder, first *int64, offset *int64) ([]*models.File, error) {
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
		Files []*models.File "json:\"queryFile\" graphql:\"queryFile\""
	}{}
	dr.log.Debugw("get Files with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetFileStr)
	}
	return rsp.Files, nil
}

// GetAllFilesWithCustomQuery returns a list of all `File` objects.
func (dr *DgraphRepository) GetAllFilesWithCustomQuery(ctx context.Context, operationName, query string) ([]*models.File, error) {
	return dr.GetFilesWithCustomQuery(ctx, operationName, query, nil, nil, nil, nil)
}

// CreateFile creates a new `File` object.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateFile(ctx context.Context, input *models.File) error {
	dr.log.Debugw("create File", []interface{}{"xid", *input.Xid}...)
	inputData := dgclient.AddFileInput{}
	dr.copyORMStruct(input, &inputData)
	rsp, err := dr.client.CreateFiles(ctx, []*dgclient.AddFileInput{&inputData})
	if err != nil {
		return WrapRepoError(err, "failed to create file").
			Add("fileId", input.ID).Add("fileXid", input.Xid)
	}
	// save ID from response
	input.ID = &rsp.AddFile.File[0].ID
	return nil
}

// CreateFiles creates new `File` objects.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateFiles(ctx context.Context, input []*models.File) error {
	inputData := make([]*dgclient.AddFileInput, 0, len(input))
	for _, v := range input {
		iv := &dgclient.AddFileInput{}
		dr.copyORMStruct(v, iv)
		inputData = append(inputData, iv)
	}

	dr.log.Debugw("create Files")
	rsp, err := dr.client.CreateFiles(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to create files")
	}

	// save ID from response
	for i, v := range input {
		v.ID = &rsp.AddFile.File[i].ID
	}

	return nil
}

// UpdateFile updates an existing `File` object.
func (dr *DgraphRepository) UpdateFile(ctx context.Context, input *models.File) error {
	dr.log.Debugw("update File", []interface{}{"id", *input.ID, "xid", *input.Xid}...)
	if *input.ID == "" {
		return WrapRepoError(nil, "missing ID").Add("fileXid", input.Xid)
	}
	patch := &dgclient.FilePatch{}
	dr.copyORMStruct(input, patch)
	patch.Xid = nil
	inputData := dgclient.UpdateFileInput{
		Filter: dgclient.FileFilter{
			ID: []string{*input.ID},
		},
		Set: patch,
	}
	_, err := dr.client.UpdateFiles(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to update file").
			Add("fileId", *input.ID).Add("fileXid", input.Xid)
	}
	return nil
}

// DeleteFile deletes a `File` object.
func (dr *DgraphRepository) DeleteFile(ctx context.Context, id, xid *string) error {
	delFilter := dgclient.FileFilter{}
	if id != nil && xid != nil {
		return NewRepoError("must specify either id or xid")
	}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &dgclient.StringHashFilter{Eq: xid}
	}

	dr.log.Debugw("delete File")
	if _, err := dr.client.DeleteFiles(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteFileStr).
			Add("fileId", id).Add("fileXid", xid)
	}
	return nil
}

// DeleteAllFiles deletes all `File` objects.
func (dr *DgraphRepository) DeleteAllFiles(ctx context.Context) error {
	delFilter := dgclient.FileFilter{}
	dr.log.Debugw("delete all File")
	if _, err := dr.client.DeleteFiles(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteFileStr)
	}
	return nil
}
