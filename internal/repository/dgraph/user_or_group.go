package dgraph

import (
	"context"
	"fmt"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"

	"github.com/jinzhu/copier"
)

var (
	errGetUserStr     = "failed to get user(s)"
	errSaveUserStr    = "failed to save user(s)"
	errDeleteUserStr  = "failed to delete user(s)"
	errGetGroupStr    = "failed to get group(s)"
	errSaveGroupStr   = "failed to save group(s)"
	errDeleteGroupStr = "failed to delete group(s)"
)

// -----------------------------------------------------------------------------
//
// User
//
// -----------------------------------------------------------------------------

// GetUser returns a `User` object by its ID.
func (dr *DgraphRepository) GetUser(id, xid *string) (*models.User, error) {
	ctx := context.Background()
	getUser, err := dr.client.GetUser(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetUserStr).
			Add("userId", id).Add("userXid", xid)
	}
	}
	user := &models.User{ID: *id}
	if err = copier.CopyWithOption(user, getUser.GetUser, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return user, nil
}

// GetUsers returns a list of `User` objects matching the filter criteria.
func (dr *DgraphRepository) GetUsers(filter *models.UserFilter, order *models.UserOrder, first *int64, offset *int64) ([]*models.User, error) {
	ctx := context.Background()
	getUsers, err := dr.client.GetUsers(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetUserStr)
	}
	users := make([]*models.User, 0, len(getUsers.QueryUser))
	for _, x := range getUsers.QueryUser {
		user := &models.User{ID: x.ID}
		if err = copier.CopyWithOption(user, x.User, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return users, nil
}

// GetAllUsers returns a list of all `User` objects.
func (dr *DgraphRepository) GetAllUsers() ([]*models.User, error) {
	return dr.GetUsers(nil, nil, nil, nil)
}

// SaveUser saves a `User` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveUser(user *models.User) (err error) {
	err = dr.SaveUsers([]*models.User{user})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("userId", user.ID).Add("userXid", user.Xid)
	}
	return
}

// SaveUsers saves `User` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveUsers(users []*models.User) error {
	reqData := make([]*models.AddUserInput, 0, len(users))
	for _, x := range users {
		if x.ID != "" {
			continue
		}
		user := &models.AddUserInput{}
			return repository.WrapRepoError(err, errSaveUserStr).
				Add("userId", x.ID).Add("userXid", x.Xid)
		}
		reqData = append(reqData, user)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveUsers(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveUserStr)
	}
	// save ID from response
	for i, x := range users {
		x.ID = respData.AddUser.User[i].ID
	}
	return nil
}

// DeleteUser deletes a `User` object.
func (dr *DgraphRepository) DeleteUser(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.UserFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteUser(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteUserStr).
			Add("userId", id).Add("userXid", xid)
	}
	return nil
}

// DeleteAllUsers deletes all `Users` objects.
func (dr *DgraphRepository) DeleteAllUsers() error {
	return dr.DeleteUser(nil, nil)
}

// -----------------------------------------------------------------------------
//
// Group
//
// -----------------------------------------------------------------------------

// GetGroup returns a `Group` object by its ID.
func (dr *DgraphRepository) GetGroup(id, xid *string) (*models.Group, error) {
	ctx := context.Background()
	getGroup, err := dr.client.GetGroup(ctx, id, xid)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetGroupStr).
			Add("groupId", id).Add("groupXid", xid)
	}
	}
	group := &models.Group{ID: *id}
	if err = copier.CopyWithOption(group, getGroup.GetGroup, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
		panic(err)
	}
	return group, nil
}

// GetGroups returns a list of `Group` objects matching the filter criteria.
func (dr *DgraphRepository) GetGroups(filter *models.GroupFilter, order *models.GroupOrder, first *int64, offset *int64) ([]*models.Group, error) {
	ctx := context.Background()
	getGroups, err := dr.client.GetGroups(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetGroupStr)
	}
	groups := make([]*models.Group, 0, len(getGroups.QueryGroup))
	for _, x := range getGroups.QueryGroup {
		group := &models.Group{ID: x.ID}
		if err = copier.CopyWithOption(group, x.Group, copier.Option{DeepCopy: true, IgnoreEmpty: true}); err != nil {
			panic(err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// GetAllGroups returns a list of all `Group` objects.
func (dr *DgraphRepository) GetAllGroups() ([]*models.Group, error) {
	return dr.GetGroups(nil, nil, nil, nil)
}

// SaveGroup saves a `Group` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveGroup(group *models.Group) (err error) {
	err = dr.SaveGroups([]*models.Group{group})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("groupId", group.ID).Add("groupXid", group.Xid)
	}
	return
}

// SaveGroups saves `Group` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveGroups(groups []*models.Group) error {
	reqData := make([]*models.AddGroupInput, 0, len(groups))
	for _, x := range groups {
		if x.ID != "" {
			continue
		}
		group := &models.AddGroupInput{}
			return repository.WrapRepoError(err, errSaveGroupStr).
				Add("groupId", x.ID).Add("groupXid", x.Xid)
		}
		reqData = append(reqData, group)
	}
	ctx := context.Background()
	respData, err := dr.client.SaveGroups(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveGroupStr)
	}
	// save ID from response
	for i, x := range groups {
		x.ID = respData.AddGroup.Group[i].ID
	}
	return nil
}

// DeleteGroup deletes a `Group` object.
func (dr *DgraphRepository) DeleteGroup(id, xid *string) error {
	ctx := context.Background()
	delFilter := models.GroupFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if xid != nil {
		delFilter.Xid = &models.StringHashFilter{Eq: xid}
	}
	_, err := dr.client.DeleteGroup(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteGroupStr).
			Add("groupId", id).Add("groupXid", xid)
	}
	return nil
}

// DeleteAllGroups deletes all `Groups` objects.
func (dr *DgraphRepository) DeleteAllGroups() error {
	return dr.DeleteGroup(nil, nil)
}

// -----------------------------------------------------------------------------
//
// User Or Group
//
// -----------------------------------------------------------------------------

// saveUserOrGroupIfNecessary saves a `UserOrGroup` object if it is not already saved.
func (dr *DgraphRepository) saveUserOrGroupIfNecessary(userOrGroup models.UserOrGroup) (*models.UserOrGroupRef, error) {
	if userOrGroup == nil {
		return nil, nil
	}
	switch o := userOrGroup.(type) {
	case *models.User:
		userRef, err := dr.saveUserIfNecessary(o)
		if err != nil {
			return nil, err
		}
		return &models.UserOrGroupRef{ID: userRef.ID}, nil
	case *models.Group:
		groupRef, err := dr.saveGroupIfNecessary(o)
		if err != nil {
			return nil, err
		}
		return &models.UserOrGroupRef{ID: groupRef.ID}, nil
	}
	panic(fmt.Errorf("implementation error: unknown UserOrGroup type: %T", userOrGroup))
}

// saveUserIfNecessary saves a `User` object if it is not already saved.
func (dr *DgraphRepository) saveUserIfNecessary(user *models.User) (*models.UserRef, error) {
	if user == nil {
		return nil, nil
	}
	if user.ID == "" {
		if err := dr.SaveUser(user); err != nil {
			return nil, err
		}
	}
	return &models.UserRef{ID: &user.ID}, nil
}

// saveGroupIfNecessary saves a `Group` object if it is not already saved.
func (dr *DgraphRepository) saveGroupIfNecessary(group *models.Group) (*models.GroupRef, error) {
	if group == nil {
		return nil, nil
	}
	if group.ID == "" {
		if err := dr.SaveGroup(group); err != nil {
			return nil, err
		}
	}
	return &models.GroupRef{ID: &group.ID}, nil
}
