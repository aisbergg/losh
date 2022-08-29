package dgraph

import "github.com/fenos/dqlx"

var FileDQLFragment = dqlx.QueryBuilder{}.Raw(`
	uid
	CrawlerMeta.discoveredAt
	CrawlerMeta.lastIndexedAt
	File.name
	File.path
	File.mimeType
	File.url
	File.createdAt
`)
