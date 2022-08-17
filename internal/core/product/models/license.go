package models

import (
	"strings"

	"losh/internal/infra/dgraph/dgclient"
)

var _ Node = (*License)(nil)

type License struct {
	ID            *string               `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Xid           *string               `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Name          *string               `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Text          *string               `json:"text,omitempty" graphql:"text"`
	TextHTML      *string               `json:"textHTML,omitempty" graphql:"textHTML"`
	ReferenceURL  *string               `json:"referenceURL,omitempty" graphql:"referenceURL"`
	DetailsURL    *string               `json:"detailsURL,omitempty" graphql:"detailsURL"`
	Type          *dgclient.LicenseType `mandatory:"true" json:"type,omitempty" graphql:"type"`
	IsSpdx        *bool                 `mandatory:"true" json:"isSpdx,omitempty" graphql:"isSpdx"`
	IsDeprecated  *bool                 `mandatory:"true" json:"isDeprecated,omitempty" graphql:"isDeprecated"`
	IsOsiApproved *bool                 `mandatory:"true" json:"isOsiApproved,omitempty" graphql:"isOsiApproved"`
	IsFsfLibre    *bool                 `mandatory:"true" json:"isFsfLibre,omitempty" graphql:"isFsfLibre"`
	IsBlocked     *bool                 `mandatory:"true" json:"isBlocked,omitempty" graphql:"isBlocked"`
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
