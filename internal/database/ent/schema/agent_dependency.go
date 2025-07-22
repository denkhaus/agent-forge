package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AgentDependency holds the schema definition for the AgentDependency entity.
type AgentDependency struct {
	ent.Schema
}

// Fields of the AgentDependency.
func (AgentDependency) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.Enum("type").
			Values("RUNTIME", "BUILD", "OPTIONAL", "PEER").
			Default("RUNTIME"),
		field.String("dependency_name"),
		field.String("dependency_version"),
		field.String("version_range"),
		field.Bool("is_required").
			Default(true),
		field.String("condition").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the AgentDependency.
func (AgentDependency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("agent", Agent.Type).
			Ref("dependencies").
			Unique().
			Required(),
	}
}

// Indexes of the AgentDependency.
func (AgentDependency) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("is_required"),
		index.Fields("dependency_name", "dependency_version").
			Edges("agent").
			Unique(),
	}
}