package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Component holds the schema definition for the Component entity.
type Component struct {
	ent.Schema
}

// Fields of the Component.
func (Component) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name"),
		field.String("namespace").
			Default("default"),
		field.String("version"),
		field.Enum("kind").
			Values("TOOL", "PROMPT", "AGENT"),
		field.Text("description"),
		field.String("author"),
		field.String("license"),
		field.String("homepage").
			Optional().
			Nillable(),
		field.String("documentation").
			Optional().
			Nillable(),
		field.JSON("tags", []string{}).
			Optional(),
		field.JSON("categories", []string{}).
			Optional(),
		field.JSON("keywords", []string{}).
			Optional(),
		field.Enum("stability").
			Values("EXPERIMENTAL", "BETA", "STABLE", "DEPRECATED").
			Default("EXPERIMENTAL"),
		field.Enum("maturity").
			Values("ALPHA", "BETA", "STABLE", "MATURE").
			Default("ALPHA"),
		field.String("forge_version"),
		field.JSON("platforms", []string{}).
			Optional(),
		field.Text("spec"),
		field.String("spec_hash"),
		field.Bool("is_installed").
			Default(false),
		field.String("install_path").
			Optional().
			Nillable(),
		field.Time("installed_at").
			Optional().
			Nillable(),
		field.String("commit_hash"),
		field.String("branch").
			Default("main"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Component.
func (Component) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Ref("components").
			Unique().
			Required(),
		edge.To("dependencies", ComponentDependency.Type),
	}
}

// Indexes of the Component.
func (Component) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind"),
		index.Fields("stability"),
		index.Fields("is_installed"),
		index.Fields("commit_hash"),
		index.Fields("name", "version").
			Edges("repository").
			Unique(),
	}
}
