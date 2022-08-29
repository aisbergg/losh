package models

var _ Node = (*Tag)(nil)

type Tag struct {
	ID      *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid     *string `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"Tag.xid"`
	Name    *string `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Tag.name"`
	Aliases []*Tag  `json:"aliases,omitempty" graphql:"aliases" dql:"Tag.aliases"`
	Related []*Tag  `json:"related,omitempty" graphql:"related" dql:"Tag.related"`
}

// GetID returns the ID of the node.
func (t *Tag) GetID() *string {
	return t.ID
}

// GetAltID returns the alternative IDs of the node.
func (t *Tag) GetAltID() *string {
	return t.Xid
}

func (*Tag) IsNode() {}
