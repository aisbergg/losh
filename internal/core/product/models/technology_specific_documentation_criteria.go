package models

var _ Node = (*TechnologySpecificDocumentationCriteria)(nil)

type TechnologySpecificDocumentationCriteria struct {
	ID              *string      `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Xid             *string      `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Name            *string      `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Comment         *string      `json:"comment,omitempty" graphql:"comment"`
	RequirementsURI *string      `json:"requirementsUri,omitempty" graphql:"requirementsUri"`
	Components      []*Component `json:"components,omitempty" graphql:"components"`
}

// GetID returns the ID of the node.
func (tsdc *TechnologySpecificDocumentationCriteria) GetID() *string {
	return tsdc.ID
}

// GetAltID returns the alternative IDs of the node.
func (tsdc *TechnologySpecificDocumentationCriteria) GetAltID() *string {
	return tsdc.Xid
}

func (*TechnologySpecificDocumentationCriteria) IsNode() {}
