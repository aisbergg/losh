package models

import "time"

var _ Node = (*File)(nil)

type File struct {
	ID            *string     `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	DiscoveredAt  *time.Time  `mandatory:"true" json:"discoveredAt,omitempty" graphql:"discoveredAt" dql:"CrawlerMeta.discoveredAt"`
	LastIndexedAt *time.Time  `mandatory:"true" json:"lastIndexedAt,omitempty" graphql:"lastIndexedAt" dql:"CrawlerMeta.lastIndexedAt"`
	DataSource    *Repository `mandatory:"true" json:"dataSource,omitempty" graphql:"dataSource" dql:"File.dataSource"`
	Xid           *string     `altID:"true" mandatory:"true" json:"xid,omitempty" graphql:"xid" dql:"File.xid"`
	Name          *string     `mandatory:"true" json:"name,omitempty" graphql:"name" dql:"File.name"`
	Path          *string     `mandatory:"true" json:"path,omitempty" graphql:"path" dql:"File.path"`
	MimeType      *string     `json:"mimeType,omitempty" graphql:"mimeType" dql:"File.mimeType"`
	URL           *string     `mandatory:"true" json:"url,omitempty" graphql:"url" dql:"File.url"`
	CreatedAt     *time.Time  `json:"createdAt,omitempty" graphql:"createdAt" dql:"File.createdAt"`
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
