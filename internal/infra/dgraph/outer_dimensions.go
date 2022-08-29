package dgraph

import "github.com/fenos/dqlx"

var OuterDimensionsDQLFragment = dqlx.QueryBuilder{}.Raw(`
	dgraph.type
	uid

	BoundingBoxDimensions.height
	BoundingBoxDimensions.width
	BoundingBoxDimensions.depth

	OpenSCADDimensions.openscad
	OpenSCADDimensions.unit
`)
