package models

type Database struct {
	ID      *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id" dql:"uid"`
	Version *string `mandatory:"true" json:"version,omitempty" graphql:"version" dql:"Database.version"`
}

func (*Database) IsNode() {}
