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

var _ Node = (*Host)(nil)

type Host struct {
	ID     *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Domain *string `altID:"true" mandatory:"true" json:"domain,omitempty" graphql:"domain" dql:"Host.domain"`
	Name   *string `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Host.name"`
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
