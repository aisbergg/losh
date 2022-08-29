package dgraph

import "github.com/fenos/dqlx"

var UserOrGroupBasicDQLFragment = dqlx.QueryBuilder{}.
	Raw(`
		dgraph.type
		uid
		UserOrGroup.name
		UserOrGroup.fullName
	`)

var UserOrGroupFullDQLFragment = dqlx.QueryBuilder{}.Raw(`
	dgraph.type
	uid
	UserOrGroup.host {
		uid
		Host.name
	}
	UserOrGroup.name
	UserOrGroup.fullName
	UserOrGroup.email
	UserOrGroup.avatar {
		uid
		File.path
	}
	UserOrGroup.url
	UserOrGroup.memberOf {
		uid
		Group.fullName
	}
	UserOrGroup.products {
		uid
		Product.name
	}

	User.locale
	Group.members {
		dgraph.type
		uid
	}
`)
