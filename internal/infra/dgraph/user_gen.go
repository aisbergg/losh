// Code generated by codegen, DO NOT EDIT.

package dgraph

import (
	"context"

	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net/request"
)

// make sure the struct implements the interface
var _ UserRepository = (*DgraphRepository)(nil)

// UserRepository is an interface for getting and saving `User` objects to a repository.
type UserRepository interface {
	GetUser(ctx context.Context, id, xid *string) (*models.User, error)
	GetUsers(ctx context.Context, filter *dgclient.UserFilter, order *dgclient.UserOrder, first *int64, offset *int64) ([]*models.User, int64, error)
	GetAllUsers(ctx context.Context) ([]*models.User, int64, error)
	CreateUser(ctx context.Context, input *models.User) error
	CreateUsers(ctx context.Context, input []*models.User) error
	UpdateUser(ctx context.Context, input *models.User) error
	DeleteUser(ctx context.Context, id, xid *string) error
	DeleteAllUsers(ctx context.Context) error
}

var (
	errGetUserStr    = "failed to get user(s)"
	errSaveUserStr   = "failed to save user(s)"
	errDeleteUserStr = "failed to delete user(s)"
)

// GetUser returns a `User` object by its ID.
func (dr *DgraphRepository) GetUser(ctx context.Context, id, xid *string) (*models.User, error) {
	var rspData interface{}
	if id != nil {
		dr.log.Debugw("get User", "id", *id)
		rsp, err := dr.client.GetUserByID(ctx, *id)
		if err != nil {
			return nil, WrapRepoError(err, errGetUserStr).Add("userId", id)
		}
		rspData = rsp.GetUser
	} else if xid != nil {
		dr.log.Debugw("get User", "xid", *xid)
		rsp, err := dr.client.GetUserByXid(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetUserStr).Add("userXid", xid)
		}
		rspData = rsp.GetUser
	} else {
		panic("must specify id or xid")
	}

	if rspData == nil {
		return nil, nil
	}
	ret := &models.User{}
	if err := dr.copier.CopyTo(rspData, ret); err != nil {
		panic(err)
	}
	return ret, nil
}

// GetUserID returns the ID of an existing `User` object.
func (dr *DgraphRepository) GetUserID(ctx context.Context, xid *string) (*string, error) {
	if xid != nil {
		dr.log.Debugw("get User", "xid", *xid)
		rsp, err := dr.client.GetUserID(ctx, *xid)
		if err != nil {
			return nil, WrapRepoError(err, errGetUserStr).Add("userXid", xid)
		}
		if rsp.GetUser == nil {
			return nil, nil
		}
		return &rsp.GetUser.ID, nil
	}

	panic("must specify xid")
}

// GetUsers returns a list of `User` objects matching the filter criteria.
func (dr *DgraphRepository) GetUsers(ctx context.Context, filter *dgclient.UserFilter, order *dgclient.UserOrder, first *int64, offset *int64) ([]*models.User, int64, error) {
	dr.log.Debugw("get Users")
	rsp, err := dr.client.GetUsers(ctx, filter, order, first, offset)
	if err != nil {
		return nil, 0, WrapRepoError(err, errGetUserStr)
	}
	ret := make([]*models.User, 0, len(rsp.QueryUser))
	if err = dr.copier.CopyTo(rsp.QueryUser, &ret); err != nil {
		panic(err)
	}
	return ret, *rsp.AggregateUser.Count, nil
}

// GetAllUsers returns a list of all `User` objects.
func (dr *DgraphRepository) GetAllUsers(ctx context.Context) ([]*models.User, int64, error) {
	return dr.GetUsers(ctx, nil, nil, nil, nil)
}

// GetUserWithCustomQuery returns a `User` object by its ID.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetUserWithCustomQuery(ctx context.Context, operationName, query string, id, xid *string) (*models.User, error) {
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
		User *models.User "json:\"getUser\" graphql:\"getUser\""
	}{}
	dr.log.Debugw("get User with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetUserStr)
	}
	return rsp.User, nil
}

// GetUsersWithCustomQuery returns a list of `User` objects matching the filter criteria.
// The given query controls the amount of information to be returned.
func (dr *DgraphRepository) GetUsersWithCustomQuery(ctx context.Context, operationName, query string, filter *dgclient.UserFilter, order *dgclient.UserOrder, first *int64, offset *int64) ([]*models.User, error) {
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
		Users []*models.User "json:\"queryUser\" graphql:\"queryUser\""
	}{}
	dr.log.Debugw("get Users with custom query")
	if err := dr.requester.Do(req, &rsp); err != nil {
		return nil, WrapRepoError(err, errGetUserStr)
	}
	return rsp.Users, nil
}

// GetAllUsersWithCustomQuery returns a list of all `User` objects.
func (dr *DgraphRepository) GetAllUsersWithCustomQuery(ctx context.Context, operationName, query string) ([]*models.User, error) {
	return dr.GetUsersWithCustomQuery(ctx, operationName, query, nil, nil, nil, nil)
}

// CreateUser creates a new `User` object.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateUser(ctx context.Context, input *models.User) error {
	dr.log.Debugw("create User", []interface{}{"xid", *input.Xid}...)
	inputData := dgclient.AddUserInput{}
	dr.copyORMStruct(input, &inputData)
	rsp, err := dr.client.CreateUsers(ctx, []*dgclient.AddUserInput{&inputData})
	if err != nil {
		return WrapRepoError(err, "failed to create user").
			Add("userId", input.ID).Add("userXid", input.Xid)
	}
	// save ID from response
	input.ID = &rsp.AddUser.User[0].ID
	return nil
}

// CreateUsers creates new `User` objects.
// After successful creation the ID field of the input will be populated with
// the ID assigned by the DB.
func (dr *DgraphRepository) CreateUsers(ctx context.Context, input []*models.User) error {
	inputData := make([]*dgclient.AddUserInput, 0, len(input))
	for _, v := range input {
		iv := &dgclient.AddUserInput{}
		dr.copyORMStruct(v, iv)
		inputData = append(inputData, iv)
	}

	dr.log.Debugw("create Users")
	rsp, err := dr.client.CreateUsers(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to create users")
	}

	// save ID from response
	for i, v := range input {
		v.ID = &rsp.AddUser.User[i].ID
	}

	return nil
}

// UpdateUser updates an existing `User` object.
func (dr *DgraphRepository) UpdateUser(ctx context.Context, input *models.User) error {
	dr.log.Debugw("update User", []interface{}{"id", *input.ID, "xid", *input.Xid}...)
	if *input.ID == "" {
		return WrapRepoError(nil, "missing ID").Add("userXid", input.Xid)
	}
	patch := &dgclient.UserPatch{}
	dr.copyORMStruct(input, patch)
	patch.Xid = nil
	inputData := dgclient.UpdateUserInput{
		Filter: dgclient.UserFilter{
			ID: []string{*input.ID},
		},
		Set: patch,
	}
	_, err := dr.client.UpdateUsers(ctx, inputData)
	if err != nil {
		return WrapRepoError(err, "failed to update user").
			Add("userId", *input.ID).Add("userXid", input.Xid)
	}
	return nil
}

// DeleteUser deletes a `User` object.
func (dr *DgraphRepository) DeleteUser(ctx context.Context, id, xid *string) error {
	delFilter := dgclient.UserFilter{}
	if id != nil && xid != nil {
		return NewRepoError("must specify either id or xid")
	}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &dgclient.StringHashFilter{Eq: xid}
	}

	dr.log.Debugw("delete User")
	if _, err := dr.client.DeleteUsers(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteUserStr).
			Add("userId", id).Add("userXid", xid)
	}
	return nil
}

// DeleteAllUsers deletes all `User` objects.
func (dr *DgraphRepository) DeleteAllUsers(ctx context.Context) error {
	delFilter := dgclient.UserFilter{}
	dr.log.Debugw("delete all User")
	if _, err := dr.client.DeleteUsers(ctx, delFilter); err != nil {
		return WrapRepoError(err, errDeleteUserStr)
	}
	return nil
}