package models

type Database struct {
	ID      *string `id:"true" mandatory:"true" json:"id,omitempty" graphql:"id"`
	Version *string `mandatory:"true" json:"version,omitempty" graphql:"version"`
}

func (*Database) IsNode() {}
