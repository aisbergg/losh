package models

import "encoding/json"

var _ Node = (*Repository)(nil)

type Repository struct {
	ID        *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid       *string     `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"Repository.xid"`
	URL       *string     `mandatory:"true" json:"url,omitempty" graphql:"url" dql:"Repository.url"`
	PermaURL  *string     `mandatory:"true" json:"permaUrl,omitempty" graphql:"permaUrl" dql:"Repository.permaUrl"`
	Host      *Host       `mandatory:"true" json:"host,omitempty" graphql:"host" dql:"Repository.host"`
	Owner     UserOrGroup `json:"owner,omitempty" graphql:"owner" dql:"Repository.owner"`
	Name      *string     `json:"name,omitempty" graphql:"name" dql:"Repository.name"`
	Reference *string     `json:"reference,omitempty" graphql:"reference" dql:"Repository.reference"`
	Path      *string     `json:"path,omitempty" graphql:"path" dql:"Repository.path"`
}

// GetID returns the ID of the node.
func (r *Repository) GetID() *string {
	return r.ID
}

// GetAltID returns the alternative IDs of the node.
func (r *Repository) GetAltID() *string {
	return r.Xid
}

type repositoryAlias Repository

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Repository) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	// unmarshal UserOrGroup field
	rawOwner, _ := objMap["owner"]
	r.Owner, err = unmarshalUserOrGroup(rawOwner)
	if err != nil {
		return err
	}
	// unmarshal the repository itself
	err = json.Unmarshal(b, (*repositoryAlias)(r))
	if err != nil {
		return err
	}

	return nil
}

func (*Repository) IsNode() {}
