package models

import (
	"encoding/json"
	"time"

	"losh/internal/infra/dgraph/dgclient"

	"github.com/aisbergg/go-errors/pkg/errors"
)

var _ Node = (*Product)(nil)

type Product struct {
	ID                    *string                `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	DiscoveredAt          *time.Time             `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt" dql:"CrawlerMeta.discoveredAt"`
	LastIndexedAt         *time.Time             `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt" dql:"CrawlerMeta.lastIndexedAt"`
	DataSource            *Repository            `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource" dql:"Product.dataSource"`
	Xid                   *string                `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"Product.xid"`
	Name                  *string                `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Product.name"`
	Description           *string                `mandatory:"true" json:"description,omitempty" graphql:"description" dql:"Product.description"`
	DocumentationLanguage *string                `mandatory:"true" json:"documentationLanguage,omitempty" graphql:"documentationLanguage" dql:"Product.documentationLanguage"`
	Version               *string                `mandatory:"true" json:"version,omitempty" graphql:"version" dql:"Product.version"`
	License               *License               `json:"license,omitempty" graphql:"license" dql:"Product.license"`
	Licensor              UserOrGroup            `mandatory:"true" json:"licensor,omitempty" graphql:"licensor" dql:"Product.licensor"`
	Website               *string                `json:"website,omitempty" graphql:"website" dql:"Product.website"`
	State                 *dgclient.ProductState `mandatory:"true" json:"state,omitempty" graphql:"state" dql:"Product.state"`
	Release               *Component             `mandatory:"true" json:"release,omitempty" graphql:"release" dql:"Product.release"`
	Releases              []*Component           `mandatory:"true" json:"releases,omitempty" graphql:"releases" dql:"Product.releases"`
	RenamedTo             *Product               `json:"renamedTo,omitempty" graphql:"renamedTo" dql:"Product.renamedTo"`
	RenamedFrom           *Product               `json:"renamedFrom,omitempty" graphql:"renamedFrom" dql:"Product.renamedFrom"`
	ForkOf                *Product               `json:"forkOf,omitempty" graphql:"forkOf" dql:"Product.forkOf"`
	Forks                 []*Product             `json:"forks,omitempty" graphql:"forks" dql:"Product.forks"`
	ForkCount             *int64                 `json:"forkCount" graphql:"forkCount" dql:"Product.forkCount"`
	StarCount             *int64                 `json:"starCount" graphql:"starCount" dql:"Product.starCount"`
	Tags                  []*Tag                 `json:"tags,omitempty" graphql:"tags" dql:"Product.tags"`
	Category              *Category              `json:"category,omitempty" graphql:"category" dql:"Product.category"`
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
