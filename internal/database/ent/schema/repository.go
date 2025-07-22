package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Repository holds the schema definition for the Repository entity.
type Repository struct {
	ent.Schema
}

// Fields of the Repository.
func (Repository) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name").
			Unique(),
		field.String("url"),
		field.Enum("type").
			Values("GITHUB", "GITLAB", "BITBUCKET", "LOCAL", "OTHER").
			Default("GITHUB"),
		field.Bool("is_active").
			Default(true),
		field.String("default_branch").
			Default("main"),
		field.Time("last_sync").
			Optional().
			Nillable(),
		field.Enum("sync_status").
			Values("NEVER_SYNCED", "SYNCING", "UP_TO_DATE", "BEHIND", "AHEAD", "DIVERGED", "ERROR").
			Default("NEVER_SYNCED"),
		field.Text("manifest").
			Optional().
			Nillable(),
		field.String("manifest_hash").
			Optional().
			Nillable(),
		field.Bool("has_write_access").
			Default(false),
		field.String("access_token").
			Optional().
			Nillable().
			Sensitive(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Repository.
func (Repository) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("components", Component.Type),
		edge.To("tools", Tool.Type),
		edge.To("prompts", Prompt.Type),
		edge.To("agents", Agent.Type),
		edge.To("forks", Fork.Type),
	}
}

// Indexes of the Repository.
func (Repository) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("is_active"),
		index.Fields("sync_status"),
	}
}
