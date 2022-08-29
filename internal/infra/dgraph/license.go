package dgraph

import "github.com/fenos/dqlx"

var LicenseBasicDQLFragment = dqlx.QueryBuilder{}.Raw(`
	uid
	License.xid
	License.name
	License.isSpdx
	License.isDeprecated
	License.isOsiApproved
	License.isFsfLibre
	License.isBlocked
`)
