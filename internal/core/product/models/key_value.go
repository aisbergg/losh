package models

import (
	"encoding/json"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
)

type StringOrFloat interface {
	IsStringOrFloat()
}

// TODO: use for unmarshalling KeyValue
func unmarshalStringOrFloat(msg *json.RawMessage) (res StringOrFloat, err error) {
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
	case `"stringv"`:
		res = &StringV{}
	case `"floatv"`:
		res = &FloatV{}
	default:
		return nil, errors.New("unknown type: " + typename)
	}

	return res, nil
}

// ----------------------------------------------------------------------------

var _ Node = (*KeyValue)(nil)

type KeyValue struct {
	ID    *string       `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Key   *string       `mandatory:"true" json:"key,omitempty" graphql:"key"`
	Value StringOrFloat `mandatory:"true" json:"value,omitempty" graphql:"value"`
}

// GetID returns the ID of the node.
func (kv *KeyValue) GetID() *string {
	return kv.ID
}

// GetAltID returns the alternative IDs of the node.
func (kv *KeyValue) GetAltID() *string {
	return nil
}

func (*KeyValue) IsNode() {}

// -----------------------------------------------------------------------------

var _ Node = (*FloatV)(nil)

type FloatV struct {
	ID    *string  `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Value *float64 `mandatory:"true" json:"value,omitempty" graphql:"value"`
}

// GetID returns the ID of the node.
func (fv *FloatV) GetID() *string {
	return fv.ID
}

// GetAltID returns the alternative IDs of the node.
func (fv *FloatV) GetAltID() *string {
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v *FloatV) MarshalJSON() ([]byte, error) {
	type Alias FloatV
	return json.Marshal(&struct {
		Typename string `json:"__typename"`
		Alias
	}{
		Typename: "FloatV",
		Alias:    (Alias)(*v),
	})
}

func (*FloatV) IsStringOrFloat() {}
func (*FloatV) IsNode()          {}

// -----------------------------------------------------------------------------

var _ Node = (*StringV)(nil)

type StringV struct {
	ID    *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Value *string `mandatory:"true" json:"value,omitempty" graphql:"value"`
}

// GetID returns the ID of the node.
func (sv *StringV) GetID() *string {
	return sv.ID
}

// GetAltID returns the alternative IDs of the node.
func (sv *StringV) GetAltID() *string {
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v *StringV) MarshalJSON() ([]byte, error) {
	type Alias StringV
	return json.Marshal(&struct {
		Typename string `json:"__typename"`
		Alias
	}{
		Typename: "StringV",
		Alias:    (Alias)(*v),
	})
}

func (*StringV) IsStringOrFloat() {}
func (*StringV) IsNode()          {}
