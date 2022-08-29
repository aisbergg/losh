package dgraph

import "github.com/fenos/dqlx"

var CategoryDQLFragment = dqlx.QueryBuilder{}.Raw(`
	uid
	Category.xid
	Category.fullName
	Category.name
	Category.description
	Category.parent {
		uid
		Category.fullName
	}
	Category.children {
		uid
		Category.fullName
	}
	Category.products {
		uid
		Product.name
	}
`)
