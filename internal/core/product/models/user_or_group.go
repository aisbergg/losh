package models

import (
	"encoding/json"
	"strings"
)

type UserOrGroup interface {
	IsUserOrGroup()
	GetName() *string
	GetID() *string
	GetAltID() *string
	// GetXid() *string
	// GetHost() *Host
	// GetFullName() *string
	// GetEmail() *string
	// GetAvatar() *File
	// GetUrl() *string
	// GetMemberOf() []*Group
	// GetProducts() []*Product
}

func unmarshalUserOrGroup(msg *json.RawMessage) (res UserOrGroup, err error) {
	if msg == nil {
		return res, nil
	}
	var o map[string]*json.RawMessage
	err = json.Unmarshal(*msg, &o)
	if err != nil {
		return nil, err
	}
	tn, ok := o["__typename"]
	typename := "User"
	if ok {
		typename = string(*tn)
	}
	switch strings.ToLower(typename) {
	case `"group"`:
		res = &Group{}
	default:
		// default to User
		res = &User{}
	}

	return res, nil
}

// -----------------------------------------------------------------------------

var _ Node = (*User)(nil)
var _ UserOrGroup = (*Group)(nil)

type User struct {
	ID       *string    `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Xid      *string    `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Host     *Host      `mandatory:"true" json:"host,omitempty" graphql:"host"`
	Name     *string    `mandatory:"true" json:"name,omitempty" graphql:"name"`
	FullName *string    `json:"fullName,omitempty" graphql:"fullName"`
	Email    *string    `json:"email,omitempty" graphql:"email"`
	Avatar   *File      `json:"avatar,omitempty" graphql:"avatar"`
	URL      *string    `json:"url,omitempty" graphql:"url"`
	MemberOf []*Group   `json:"memberOf,omitempty" graphql:"memberOf"`
	Products []*Product `json:"products,omitempty" graphql:"products"`
	Locale   *string    `json:"locale,omitempty" graphql:"locale"`
}

// GetID returns the ID of the node.
func (u *User) GetID() *string { return u.ID }

// GetAltID returns the alternative IDs of the node.
func (u *User) GetAltID() *string { return u.Xid }

// GetName returns the name of the user.
func (u *User) GetName() *string { return u.Name }

// MarshalJSON implements the json.Marshaler interface.
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Typename string `json:"__typename"`
		Alias
	}{
		Typename: "User",
		Alias:    (Alias)(*u),
	})
}

func (*User) IsNode()        {}
func (*User) IsUserOrGroup() {}

// -----------------------------------------------------------------------------

var _ Node = (*Group)(nil)
var _ UserOrGroup = (*Group)(nil)

type Group struct {
	ID       *string       `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Xid      *string       `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Host     *Host         `mandatory:"true" json:"host,omitempty" graphql:"host"`
	Name     *string       `mandatory:"true" json:"name,omitempty" graphql:"name"`
	FullName *string       `json:"fullName,omitempty" graphql:"fullName"`
	Email    *string       `json:"email,omitempty" graphql:"email"`
	Avatar   *File         `json:"avatar,omitempty" graphql:"avatar"`
	URL      *string       `json:"url,omitempty" graphql:"url"`
	MemberOf []*Group      `json:"memberOf,omitempty" graphql:"memberOf"`
	Products []*Product    `json:"products,omitempty" graphql:"products"`
	Members  []UserOrGroup `json:"members,omitempty" graphql:"members"`
}

// GetID returns the ID of the node.
func (g *Group) GetID() *string { return g.ID }

// GetAltID returns the alternative IDs of the node.
func (g *Group) GetAltID() *string { return g.Xid }

// GetName returns the name of the group.
func (g *Group) GetName() *string { return g.Name }

// MarshalJSON implements the json.Marshaler interface.
func (g *Group) MarshalJSON() ([]byte, error) {
	type Alias Group
	return json.Marshal(&struct {
		Typename string `json:"__typename"`
		Alias
	}{
		Typename: "Group",
		Alias:    (Alias)(*g),
	})
}

func (*Group) IsNode()        {}
func (*Group) IsUserOrGroup() {}
