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
