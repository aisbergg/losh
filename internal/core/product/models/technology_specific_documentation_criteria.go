package models

var _ Node = (*TechnologySpecificDocumentationCriteria)(nil)

type TechnologySpecificDocumentationCriteria struct {
	ID              *string      `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Xid             *string      `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"TechnologySpecificDocumentationCriteria.xid"`
	Name            *string      `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"TechnologySpecificDocumentationCriteria.name"`
	Comment         *string      `json:"comment,omitempty" graphql:"comment" dql:"TechnologySpecificDocumentationCriteria.comment"`
	RequirementsURI *string      `json:"requirementsUri,omitempty" graphql:"requirementsUri" dql:"TechnologySpecificDocumentationCriteria.requirementsUri"`
	Components      []*Component `json:"components,omitempty" graphql:"components" dql:"TechnologySpecificDocumentationCriteria.components"`
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
