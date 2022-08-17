package models

import "time"

var _ Node = (*File)(nil)

type File struct {
	ID            *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	DiscoveredAt  *time.Time  `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt"`
	LastIndexedAt *time.Time  `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt"`
	DataSource    *Repository `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource"`
	Xid           *string     `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid"`
	Name          *string     `mandatory:"true" json:"name,omitempty" graphql:"name"`
	Path          *string     `mandatory:"true" json:"path,omitempty" graphql:"path"`
	MimeType      *string     `json:"mimeType,omitempty" graphql:"mimeType"`
	URL           *string     `mandatory:"true" json:"url,omitempty" graphql:"url"`
	CreatedAt     *time.Time  `json:"createdAt,omitempty" graphql:"createdAt"`
}

// GetID returns the ID of the node.
func (f *File) GetID() *string {
	return f.ID
}

// GetAltID returns the alternative IDs of the node.
func (f *File) GetAltID() *string {
	return f.Xid
}

func (*File) IsNode()        {}
func (*File) IsCrawlerMeta() {}
