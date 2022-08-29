package models

var _ Node = (*Material)(nil)

type Material struct {
	ID          *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Name        *string `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Material.name"`
	Description *string `json:"description,omitempty" graphql:"description" dql:"Material.description"`
}

// GetID returns the ID of the node.
func (m *Material) GetID() *string {
	return m.ID
}

// GetAltID returns the alternative IDs of the node.
func (m *Material) GetAltID() *string {
	return nil
}

func (*Material) IsNode() {}
