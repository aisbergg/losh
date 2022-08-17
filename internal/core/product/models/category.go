package models

var _ Node = (*Category)(nil)

type Category struct {
	ID          *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Xid         *string     `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	FullName    *string     `mandatory:"true" json:"fullName,omitempty" graphql:"fullName"`
	Name        *string     `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Description *string     `json:"description,omitempty" graphql:"description"`
	Parent      *Category   `json:"parent,omitempty" graphql:"parent"`
	Children    []*Category `json:"children,omitempty" graphql:"children"`
	Products    []*Product  `json:"products,omitempty" graphql:"products"`
}

// GetID returns the ID of the node.
func (c *Category) GetID() *string {
	return c.ID
}

// GetAltID returns the alternative IDs of the node.
func (c *Category) GetAltID() *string {
	return c.Xid
}

func (*Category) IsNode() {}
