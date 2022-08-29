package dgraph

import "github.com/fenos/dqlx"

var RepositoryDQLFragment = dqlx.Query(nil).Raw(`
	uid
	Repository.id
	Repository.xid
	Repository.url
	Repository.permaUrl
	Repository.host {
		uid
		Host.name
	}
	Repository.name
	Repository.reference
	Repository.path
`).EdgeFromQuery(UserOrGroupBasicDQLFragment.Name("Repository.owner"))
