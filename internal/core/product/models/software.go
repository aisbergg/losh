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
