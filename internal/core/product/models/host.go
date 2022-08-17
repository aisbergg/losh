package models

var _ Node = (*Host)(nil)

type Host struct {
	ID     *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Domain *string `altID:"true" mandatory:"true" json:"domain,omitempty" graphql:"domain"`
	Name   *string `mandatory:"true" json:"name,omitempty" graphql:"name"`
}

// GetID returns the ID of the node.
func (h *Host) GetID() *string {
	return h.ID
}

// GetAltID returns the alternative IDs of the node.
func (h *Host) GetAltID() *string {
	return h.Domain
}

func (*Host) IsNode() {}
