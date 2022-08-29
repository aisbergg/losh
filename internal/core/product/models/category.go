package models

var _ Node = (*Category)(nil)

type Category struct {
	ID          *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid         *string     `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"Category.xid"`
	FullName    *string     `mandatory:"true" json:"fullName,omitempty" graphql:"fullName" dql:"Category.fullName"`
	Name        *string     `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Category.name"`
	Description *string     `json:"description,omitempty" graphql:"description" dql:"Category.description"`
	Parent      *Category   `json:"parent,omitempty" graphql:"parent" dql:"Category.parent"`
	Children    []*Category `json:"children,omitempty" graphql:"children" dql:"Category.children"`
	Products    []*Product  `json:"products,omitempty" graphql:"products" dql:"Category.products"`
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
