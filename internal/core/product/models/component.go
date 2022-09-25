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

import (
	"encoding/json"
	"time"

	"losh/internal/infra/dgraph/dgclient"
)

var _ Node = (*Component)(nil)

type Component struct {
	ID                          *string                                  `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	DiscoveredAt                *time.Time                               `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt" dql:"CrawlerMeta.discoveredAt"`
	LastIndexedAt               *time.Time                               `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt" dql:"CrawlerMeta.lastIndexedAt"`
	DataSource                  *Repository                              `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource" dql:"Component.dataSource"`
	Xid                         *string                                  `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"Component.xid"`
	Name                        *string                                  `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"Component.name"`
	Description                 *string                                  `mandatory:"true" json:"description,omitempty" graphql:"description" dql:"Component.description"`
	Version                     *string                                  `mandatory:"true" json:"version,omitempty" graphql:"version" dql:"Component.version"`
	CreatedAt                   *time.Time                               `mandatory:"true" json:"createdAt,omitempty" graphql:"createdAt" dql:"Component.createdAt"`
	Releases                    []*Component                             `json:"releases,omitempty" graphql:"releases" dql:"Component.releases"`
	IsLatest                    *bool                                    `mandatory:"true" json:"isLatest,omitempty" graphql:"isLatest" dql:"Component.isLatest"`
	Repository                  *Repository                              `mandatory:"true" json:"repository,omitempty" graphql:"repository" dql:"Component.repository"`
	License                     *License                                 `json:"license,omitempty" graphql:"license" dql:"Component.license"`
	AdditionalLicenses          []*License                               `json:"additionalLicenses,omitempty" graphql:"additionalLicenses" dql:"Component.additionalLicenses"`
	Licensor                    UserOrGroup                              `mandatory:"true" json:"licensor,omitempty" graphql:"licensor" dql:"Component.licensor"`
	DocumentationLanguage       *string                                  `mandatory:"true" json:"documentationLanguage,omitempty" graphql:"documentationLanguage" dql:"Component.documentationLanguage"`
	TechnologyReadinessLevel    *dgclient.TechnologyReadinessLevel       `mandatory:"true" json:"technologyReadinessLevel,omitempty" graphql:"technologyReadinessLevel" dql:"Component.technologyReadinessLevel"`
	DocumentationReadinessLevel *dgclient.DocumentationReadinessLevel    `mandatory:"true" json:"documentationReadinessLevel,omitempty" graphql:"documentationReadinessLevel" dql:"Component.documentationReadinessLevel"`
	Attestation                 *string                                  `json:"attestation,omitempty" graphql:"attestation" dql:"Component.attestation"`
	Publication                 *string                                  `json:"publication,omitempty" graphql:"publication" dql:"Component.publication"`
	Issues                      *string                                  `json:"issues,omitempty" graphql:"issues" dql:"Component.issues"`
	CompliesWith                *TechnicalStandard                       `json:"compliesWith,omitempty" graphql:"compliesWith" dql:"Component.compliesWith"`
	CpcPatentClass              *string                                  `json:"cpcPatentClass,omitempty" graphql:"cpcPatentClass" dql:"Component.cpcPatentClass"`
	Tsdc                        *TechnologySpecificDocumentationCriteria `json:"tsdc,omitempty" graphql:"tsdc" dql:"Component.tsdc"`
	Components                  []*Component                             `json:"components,omitempty" graphql:"components" dql:"Component.components"`
	Software                    []*Software                              `json:"software,omitempty" graphql:"software" dql:"Component.software"`
	Image                       *File                                    `json:"image,omitempty" graphql:"image" dql:"Component.image"`
	Readme                      *File                                    `json:"readme,omitempty" graphql:"readme" dql:"Component.readme"`
	ContributionGuide           *File                                    `json:"contributionGuide,omitempty" graphql:"contributionGuide" dql:"Component.contributionGuide"`
	Bom                         *File                                    `json:"bom,omitempty" graphql:"bom" dql:"Component.bom"`
	ManufacturingInstructions   *File                                    `json:"manufacturingInstructions,omitempty" graphql:"manufacturingInstructions" dql:"Component.manufacturingInstructions"`
	UserManual                  *File                                    `json:"userManual,omitempty" graphql:"userManual" dql:"Component.userManual"`
	Product                     *Product                                 `json:"product,omitempty" graphql:"product" dql:"Component.product"`
	UsedIn                      []*Component                             `json:"usedIn,omitempty" graphql:"usedIn" dql:"Component.usedIn"`
	Source                      *File                                    `json:"source,omitempty" graphql:"source" dql:"Component.source"`
	Export                      []*File                                  `json:"export,omitempty" graphql:"export" dql:"Component.export"`
	Auxiliary                   []*File                                  `json:"auxiliary,omitempty" graphql:"auxiliary" dql:"Component.auxiliary"`
	Organization                *Group                                   `json:"organization,omitempty" graphql:"organization" dql:"Component.organization"`
	Mass                        *float64                                 `json:"mass,omitempty" graphql:"mass" dql:"Component.mass"`
	OuterDimensions             OuterDimensions                          `json:"outerDimensions,omitempty" graphql:"outerDimensions" dql:"Component.outerDimensions"`
	Material                    *Material                                `json:"material,omitempty" graphql:"material" dql:"Component.material"`
	ManufacturingProcess        *ManufacturingProcess                    `json:"manufacturingProcess,omitempty" graphql:"manufacturingProcess" dql:"Component.manufacturingProcess"`
	ProductionMetadata          []*KeyValue                              `json:"productionMetadata,omitempty" graphql:"productionMetadata" dql:"Component.productionMetadata"`
}

type componentAlias Component

func (c *Component) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	// unmarshal UserOrGroup field
	raw, _ := objMap["licensor"]
	c.Licensor, err = unmarshalUserOrGroup(raw)
	if err != nil {
		return err
	}
	// unmarshal OuterDimensions field
	raw, _ = objMap["outerDimensions"]
	c.OuterDimensions, err = unmarshalOuterDimensions(raw)
	if err != nil {
		return err
	}
	// unmarshal the component itself
	err = json.Unmarshal(b, (*componentAlias)(c))
	if err != nil {
		return err
	}

	return nil
}

// GetID returns the ID of the node.
func (c *Component) GetID() *string {
	return c.ID
}

// GetAltID returns the alternative IDs of the node.
func (c *Component) GetAltID() *string {
	return c.Xid
}

func (*Component) IsNode()        {}
func (*Component) IsCrawlerMeta() {}
