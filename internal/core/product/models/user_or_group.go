// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	ID          *string    `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid         *string    `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"UserOrGroup.xid"`
	Host        *Host      `mandatory:"true" json:"host,omitempty" graphql:"host" dql:"UserOrGroup.host"`
	Name        *string    `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"UserOrGroup.name"`
	FullName    *string    `json:"fullName,omitempty" graphql:"fullName" dql:"UserOrGroup.fullName"`
	Email       *string    `json:"email,omitempty" graphql:"email" dql:"UserOrGroup.email"`
	Description *string    `json:"description,omitempty" graphql:"description" dql:"UserOrGroup.description"`
	Avatar      *File      `json:"avatar,omitempty" graphql:"avatar" dql:"UserOrGroup.avatar"`
	URL         *string    `json:"url,omitempty" graphql:"url" dql:"UserOrGroup.url"`
	MemberOf    []*Group   `json:"memberOf,omitempty" graphql:"memberOf" dql:"UserOrGroup.memberOf"`
	Products    []*Product `json:"products,omitempty" graphql:"products" dql:"UserOrGroup.products"`
	Locale      *string    `json:"locale,omitempty" graphql:"locale" dql:"User.locale"`
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
	ID          *string       `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid         *string       `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"UserOrGroup.xid"`
	Host        *Host         `mandatory:"true" json:"host,omitempty" graphql:"host" dql:"UserOrGroup.host"`
	Name        *string       `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"UserOrGroup.name"`
	FullName    *string       `json:"fullName,omitempty" graphql:"fullName" dql:"UserOrGroup.fullName"`
	Email       *string       `json:"email,omitempty" graphql:"email" dql:"UserOrGroup.email"`
	Description *string       `json:"description,omitempty" graphql:"description" dql:"UserOrGroup.description"`
	Avatar      *File         `json:"avatar,omitempty" graphql:"avatar" dql:"UserOrGroup.avatar"`
	URL         *string       `json:"url,omitempty" graphql:"url" dql:"UserOrGroup.url"`
	MemberOf    []*Group      `json:"memberOf,omitempty" graphql:"memberOf" dql:"UserOrGroup.memberOf"`
	Products    []*Product    `json:"products,omitempty" graphql:"products" dql:"UserOrGroup.products"`
	Members     []UserOrGroup `json:"members,omitempty" graphql:"members" dql:"Group.members"`
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
