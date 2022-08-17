package models

var _ Node = (*TechnicalStandard)(nil)

type TechnicalStandard struct {
	ID          *string      `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Xid         *string      `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Name        *string      `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Description *string      `json:"description,omitempty" graphql:"description"`
	Components  []*Component `json:"components,omitempty" graphql:"components"`
}

// GetID returns the ID of the node.
func (ts *TechnicalStandard) GetID() *string {
	return ts.ID
}

// GetAltID returns the alternative IDs of the node.
func (ts *TechnicalStandard) GetAltID() *string {
	return ts.Xid
}

func (*TechnicalStandard) IsNode() {}
