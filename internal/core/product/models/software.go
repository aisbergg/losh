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

import "time"

var _ Node = (*Software)(nil)

type Software struct {
	ID                    *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	DiscoveredAt          *time.Time  `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt" dql:"CrawlerMeta.discoveredAt"`
	LastIndexedAt         *time.Time  `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt" dql:"CrawlerMeta.lastIndexedAt"`
	DataSource            *Repository `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource" dql:"Software.dataSource"`
	Release               *string     `json:"release,omitempty" graphql:"release" dql:"Software.release"`
	InstallationGuide     *File       `json:"installationGuide,omitempty" graphql:"installationGuide" dql:"Software.installationGuide"`
	DocumentationLanguage *string     `json:"documentationLanguage,omitempty" graphql:"documentationLanguage" dql:"Software.documentationLanguage"`
	License               *License    `json:"license,omitempty" graphql:"license" dql:"Software.license"`
	Licensor              *string     `json:"licensor,omitempty" graphql:"licensor" dql:"Software.licensor"`
}

// GetID returns the ID of the node.
func (s *Software) GetID() *string {
	return s.ID
}

// GetAltID returns the alternative IDs of the node.
func (s *Software) GetAltID() *string {
	return nil
}

func (*Software) IsNode()        {}
func (*Software) IsCrawlerMeta() {}
