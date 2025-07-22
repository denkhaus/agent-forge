package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ComponentDependency holds the schema definition for the ComponentDependency entity.
type ComponentDependency struct {
	ent.Schema
}

// Fields of the ComponentDependency.
func (ComponentDependency) Fields() []ent.Field {
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

// Edges of the ComponentDependency.
func (ComponentDependency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("component", Component.Type).
			Ref("dependencies").
			Unique().
			Required(),
	}
}

// Indexes of the ComponentDependency.
func (ComponentDependency) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("is_required"),
		index.Fields("dependency_name", "dependency_version").
			Edges("component").
			Unique(),
	}
}
