package models

import "time"

type CrawlerMetaImpl struct {
	DiscoveredAt  time.Time        `json:"discoveredAt"`
	LastIndexedAt time.Time        `json:"lastIndexedAt"`
	DataSource    *ComponentSource `json:"dataSource,omitempty"`
}

// IsCrawlerMeta implements the CrawlerMeta interface.
func (c *CrawlerMetaImpl) IsCrawlerMeta() {}
