package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ToolDependency holds the schema definition for the ToolDependency entity.
type ToolDependency struct {
	ent.Schema
}

// Fields of the ToolDependency.
func (ToolDependency) Fields() []ent.Field {
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

// Edges of the ToolDependency.
func (ToolDependency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tool", Tool.Type).
			Ref("dependencies").
			Unique().
			Required(),
	}
}

// Indexes of the ToolDependency.
func (ToolDependency) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("is_required"),
		index.Fields("dependency_name", "dependency_version").
			Edges("tool").
			Unique(),
	}
}