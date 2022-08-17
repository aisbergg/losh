package models

import (
	"encoding/json"
	"time"

	"github.com/aisbergg/go-errors/pkg/errors"
)

var _ Node = (*Product)(nil)

type Product struct {
	ID                    *string      `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	DiscoveredAt          *time.Time   `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt"`
	LastIndexedAt         *time.Time   `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt"`
	DataSource            *Repository  `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource"`
	Xid                   *string      `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Name                  *string      `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Description           *string      `mandatory:"true" json:"description,omitempty" graphql:"description"`
	DocumentationLanguage *string      `mandatory:"true" json:"documentationLanguage,omitempty" graphql:"documentationLanguage"`
	Version               *string      `mandatory:"true" json:"version,omitempty" graphql:"version"`
	License               *License     `mandatory:"true" json:"license,omitempty" graphql:"license"`
	Licensor              UserOrGroup  `mandatory:"true" json:"licensor,omitempty" graphql:"licensor"`
	Website               *string      `json:"website,omitempty" graphql:"website"`
	Release               *Component   `mandatory:"true" json:"release,omitempty" graphql:"release"`
	Releases              []*Component `mandatory:"true" json:"releases,omitempty" graphql:"releases"`
	RenamedTo             *Product     `json:"renamedTo,omitempty" graphql:"renamedTo"`
	RenamedFrom           *Product     `json:"renamedFrom,omitempty" graphql:"renamedFrom"`
	ForkOf                *Product     `json:"forkOf,omitempty" graphql:"forkOf"`
	Forks                 []*Product   `json:"forks,omitempty" graphql:"forks"`
	Tags                  []*Tag       `json:"tags,omitempty" graphql:"tags"`
	Category              *Category    `json:"category,omitempty" graphql:"category"`
}

// GetID returns the ID of the node.
func (p *Product) GetID() *string {
	return p.ID
}

// GetAltID returns the alternative IDs of the node.
func (p *Product) GetAltID() *string {
	return p.Xid
}

type productAlias Product

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *Product) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal product")
	}
	// unmarshal UserOrGroup field
	rawLicensor := objMap["licensor"]
	p.Licensor, err = unmarshalUserOrGroup(rawLicensor)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal product licensor")
	}
	// unmarshal the product itself
	err = json.Unmarshal(b, (*productAlias)(p))
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal product")
	}

	return nil
}

func (*Product) IsNode()        {}
func (*Product) IsCrawlerMeta() {}
