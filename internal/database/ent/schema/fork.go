package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Fork holds the schema definition for the Fork entity.
type Fork struct {
	ent.Schema
}

// Fields of the Fork.
func (Fork) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("fork_url"),
		field.String("fork_owner"),
		field.String("fork_name"),
		field.Bool("is_active").
			Default(true),
		field.Time("last_sync").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Fork.
func (Fork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("original_repo", Repository.Type).
			Ref("forks").
			Unique().
			Required(),
	}
}

// Indexes of the Fork.
func (Fork) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("fork_owner", "fork_name").
			Edges("original_repo").
			Unique(),
		index.Fields("is_active"),
	}
}
