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
	"strings"

	"losh/internal/infra/dgraph/dgclient"
)

var _ Node = (*License)(nil)

type License struct {
	ID            *string               `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid           *string               `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"License.xid"`
	Name          *string               `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"License.name"`
	Text          *string               `json:"text,omitempty" graphql:"text" dql:"License.text"`
	TextHTML      *string               `json:"textHTML,omitempty" graphql:"textHTML" dql:"License.textHTML"`
	ReferenceURL  *string               `json:"referenceURL,omitempty" graphql:"referenceURL" dql:"License.referenceURL"`
	DetailsURL    *string               `json:"detailsURL,omitempty" graphql:"detailsURL" dql:"License.detailsURL"`
	Type          *dgclient.LicenseType `mandatory:"true" json:"type,omitempty" graphql:"type" dql:"License.type"`
	IsSpdx        *bool                 `mandatory:"true" json:"isSpdx,omitempty" graphql:"isSpdx" dql:"License.isSpdx"`
	IsDeprecated  *bool                 `mandatory:"true" json:"isDeprecated,omitempty" graphql:"isDeprecated" dql:"License.isDeprecated"`
	IsOsiApproved *bool                 `mandatory:"true" json:"isOsiApproved,omitempty" graphql:"isOsiApproved" dql:"License.isOsiApproved"`
	IsFsfLibre    *bool                 `mandatory:"true" json:"isFsfLibre,omitempty" graphql:"isFsfLibre" dql:"License.isFsfLibre"`
	IsBlocked     *bool                 `mandatory:"true" json:"isBlocked,omitempty" graphql:"isBlocked" dql:"License.isBlocked"`
}

// GetID returns the ID of the node.
func (l *License) GetID() *string {
	return l.ID
}

// GetAltID returns the alternative IDs of the node.
func (l *License) GetAltID() *string {
	return l.Xid
}

func (*License) IsNode() {}

// AsLicenseType returns a license type from a string.
func AsLicenseType(s string) dgclient.LicenseType {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "STRONG":
		return dgclient.LicenseTypeStrong
	case "WEAK":
		return dgclient.LicenseTypeWeak
	case "PERMISSIVE":
		return dgclient.LicenseTypePermissive
	default:
		return dgclient.LicenseTypeUnknown
	}
}
