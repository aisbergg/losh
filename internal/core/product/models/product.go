package models

import (
	"encoding/json"
	"time"

	"github.com/aisbergg/go-errors/pkg/errors"
)

var _ Node = (*Product)(nil)

type Product struct {
	License               *License               `json:"license,omitempty" graphql:"license" dql:"Product.license"`
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
