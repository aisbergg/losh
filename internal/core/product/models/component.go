package models

import (
	"encoding/json"
	"time"

	"losh/internal/infra/dgraph/dgclient"
)

var _ Node = (*Component)(nil)

type Component struct {
	ID                          *string                                  `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	DiscoveredAt                *time.Time                               `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt"`
	LastIndexedAt               *time.Time                               `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt"`
	DataSource                  *Repository                              `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource"`
	Xid                         *string                                  `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Name                        *string                                  `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Description                 *string                                  `mandatory:"true" json:"description,omitempty" graphql:"description"`
	Version                     *string                                  `mandatory:"true" json:"version,omitempty" graphql:"version"`
	CreatedAt                   *time.Time                               `mandatory:"true" json:"createdAt,omitempty" graphql:"createdAt"`
	Releases                    []*Component                             `json:"releases,omitempty" graphql:"releases"`
	IsLatest                    *bool                                    `mandatory:"true" json:"isLatest,omitempty" graphql:"isLatest"`
	Repository                  *Repository                              `mandatory:"true" json:"repository,omitempty" graphql:"repository"`
	License                     *License                                 `mandatory:"true" json:"license,omitempty" graphql:"license"`
	AdditionalLicenses          []*License                               `mandatory:"true" json:"additionalLicenses,omitempty" graphql:"additionalLicenses"`
	Licensor                    UserOrGroup                              `mandatory:"true" json:"licensor,omitempty" graphql:"licensor"`
	DocumentationLanguage       *string                                  `mandatory:"true" json:"documentationLanguage,omitempty" graphql:"documentationLanguage"`
	TechnologyReadinessLevel    *dgclient.TechnologyReadinessLevel       `mandatory:"true" json:"technologyReadinessLevel,omitempty" graphql:"technologyReadinessLevel"`
	DocumentationReadinessLevel *dgclient.DocumentationReadinessLevel    `mandatory:"true" json:"documentationReadinessLevel,omitempty" graphql:"documentationReadinessLevel"`
	Attestation                 *string                                  `json:"attestation,omitempty" graphql:"attestation"`
	Publication                 *string                                  `json:"publication,omitempty" graphql:"publication"`
	Issues                      *string                                  `json:"issues,omitempty" graphql:"issues"`
	CompliesWith                *TechnicalStandard                       `json:"compliesWith,omitempty" graphql:"compliesWith"`
	CpcPatentClass              *string                                  `json:"cpcPatentClass,omitempty" graphql:"cpcPatentClass"`
	Tsdc                        *TechnologySpecificDocumentationCriteria `json:"tsdc,omitempty" graphql:"tsdc"`
	Components                  []*Component                             `json:"components,omitempty" graphql:"components"`
	Software                    []*Software                              `json:"software,omitempty" graphql:"software"`
	Image                       *File                                    `json:"image,omitempty" graphql:"image"`
	Readme                      *File                                    `json:"readme,omitempty" graphql:"readme"`
	ContributionGuide           *File                                    `json:"contributionGuide,omitempty" graphql:"contributionGuide"`
	Bom                         *File                                    `json:"bom,omitempty" graphql:"bom"`
	ManufacturingInstructions   *File                                    `json:"manufacturingInstructions,omitempty" graphql:"manufacturingInstructions"`
	UserManual                  *File                                    `json:"userManual,omitempty" graphql:"userManual"`
	Product                     *Product                                 `json:"product,omitempty" graphql:"product"`
	UsedIn                      []*Component                             `json:"usedIn,omitempty" graphql:"usedIn"`
	Source                      *File                                    `json:"source,omitempty" graphql:"source"`
	Export                      []*File                                  `json:"export,omitempty" graphql:"export"`
	Auxiliary                   []*File                                  `json:"auxiliary,omitempty" graphql:"auxiliary"`
	Organization                *Group                                   `json:"organization,omitempty" graphql:"organization"`
	Mass                        *float64                                 `json:"mass,omitempty" graphql:"mass"`
	OuterDimensions             OuterDimensions                          `json:"outerDimensions,omitempty" graphql:"outerDimensions"`
	Material                    *Material                                `json:"material,omitempty" graphql:"material"`
	ManufacturingProcess        *ManufacturingProcess                    `json:"manufacturingProcess,omitempty" graphql:"manufacturingProcess"`
	ProductionMetadata          []*KeyValue                              `json:"productionMetadata,omitempty" graphql:"productionMetadata"`
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
