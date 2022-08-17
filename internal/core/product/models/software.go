package models

import "time"

var _ Node = (*Software)(nil)

type Software struct {
	ID                    *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	DiscoveredAt          *time.Time  `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt"`
	LastIndexedAt         *time.Time  `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt"`
	DataSource            *Repository `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource"`
	Release               *string     `json:"release,omitempty" graphql:"release"`
	InstallationGuide     *File       `json:"installationGuide,omitempty" graphql:"installationGuide"`
	DocumentationLanguage *string     `json:"documentationLanguage,omitempty" graphql:"documentationLanguage"`
	License               *License    `json:"license,omitempty" graphql:"license"`
	Licensor              *string     `json:"licensor,omitempty" graphql:"licensor"`
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
