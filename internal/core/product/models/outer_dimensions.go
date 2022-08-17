package models

import (
	"encoding/json"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
)

type OuterDimensions interface {
	IsOuterDimensions()
}

func unmarshalOuterDimensions(msg *json.RawMessage) (res OuterDimensions, err error) {
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
	case `"boundingboxdimensions"`:
		res = &BoundingBoxDimensions{}
	case `"openscaddimensions"`:
		res = &OpenSCADDimensions{}
	default:
		return nil, errors.New("unknown type: " + typename)
	}

	return res, nil
}

// -----------------------------------------------------------------------------

var _ Node = (*BoundingBoxDimensions)(nil)

type BoundingBoxDimensions struct {
	ID     *string  `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Height *float64 `mandatory:"true" json:"height,omitempty" graphql:"height"`
	Width  *float64 `mandatory:"true" json:"width,omitempty" graphql:"width"`
	Depth  *float64 `mandatory:"true" json:"depth,omitempty" graphql:"depth"`
}

// GetID returns the ID of the node.
func (bbd *BoundingBoxDimensions) GetID() *string {
	return bbd.ID
}

// GetAltID returns the alternative IDs of the node.
func (bbd *BoundingBoxDimensions) GetAltID() *string {
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v *BoundingBoxDimensions) MarshalJSON() ([]byte, error) {
	type Alias BoundingBoxDimensions
	return json.Marshal(&struct {
		Typename string `json:"__typename"`
		Alias
	}{
		Typename: "BoundingBoxDimensions",
		Alias:    (Alias)(*v),
	})
}

func (*BoundingBoxDimensions) IsOuterDimensions() {}
func (*BoundingBoxDimensions) IsNode()            {}

// -----------------------------------------------------------------------------

var _ Node = (*OpenSCADDimensions)(nil)

type OpenSCADDimensions struct {
	ID       *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Openscad *string `mandatory:"true" json:"openscad,omitempty" graphql:"openscad"`
	Unit     *string `mandatory:"true" json:"unit,omitempty" graphql:"unit"`
}

// GetID returns the ID of the node.
func (osd *OpenSCADDimensions) GetID() *string {
	return osd.ID
}

// GetAltID returns the alternative IDs of the node.
func (osd *OpenSCADDimensions) GetAltID() *string {
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v *OpenSCADDimensions) MarshalJSON() ([]byte, error) {
	type Alias OpenSCADDimensions
	return json.Marshal(&struct {
		Typename string `json:"__typename"`
		Alias
	}{
		Typename: "OpenSCADDimensions",
		Alias:    (Alias)(*v),
	})
}

func (*OpenSCADDimensions) IsOuterDimensions() {}
func (*OpenSCADDimensions) IsNode()            {}
