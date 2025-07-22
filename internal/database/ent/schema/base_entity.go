package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// BaseEntityMixin provides common fields for all AgentForge entities.
type BaseEntityMixin struct {
	mixin.Schema
}

// Fields returns the common fields for all entities.
func (BaseEntityMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("name"),
		field.String("namespace").
			Default("default"),
		field.String("version"),
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