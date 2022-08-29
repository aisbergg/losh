package dgraph

import "github.com/fenos/dqlx"

var TagDQLFragment = dqlx.QueryBuilder{}.Raw(`
	uid
	Tag.xid
	Tag.name
	Tag.aliases {
		uid
		Tag.name
	}
	Tag.related {
		uid
		Tag.name
	}
`)
