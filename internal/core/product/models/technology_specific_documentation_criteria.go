// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
