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

var _ Node = (*Tag)(nil)

type Tag struct {
	ID      *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Name    *string `altID:"true" mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Tag.name"`
	Aliases []*Tag  `json:"aliases,omitempty" graphql:"aliases" dql:"Tag.aliases"`
	Related []*Tag  `json:"related,omitempty" graphql:"related" dql:"Tag.related"`
}

// GetID returns the ID of the node.
func (t *Tag) GetID() *string {
	return t.ID
}

// GetAltID returns the alternative IDs of the node.
func (t *Tag) GetAltID() *string {
	return t.Name
}

func (*Tag) IsNode() {}
