package models

var _ Node = (*ManufacturingProcess)(nil)

type ManufacturingProcess struct {
	ID          *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Name        *string `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"ManufacturingProcess.name"`
	Description *string `json:"description,omitempty" graphql:"description" dql:"ManufacturingProcess.description"`
}

// GetID returns the ID of the node.
func (mp *ManufacturingProcess) GetID() *string {
	return mp.ID
}

// GetAltID returns the alternative IDs of the node.
func (mp *ManufacturingProcess) GetAltID() *string {
	return nil
}

func (*ManufacturingProcess) IsNode() {}
