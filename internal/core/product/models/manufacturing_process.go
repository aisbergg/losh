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
