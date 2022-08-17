package models

import "time"

type CrawlerMetaImpl struct {
	DiscoveredAt  *time.Time  `json:"discoveredAt"`
	LastIndexedAt *time.Time  `json:"lastIndexedAt"`
	DataSource    *Repository `json:"dataSource"`
}

// IsCrawlerMeta implements the CrawlerMeta interface.
func (*CrawlerMetaImpl) IsCrawlerMeta() {}
