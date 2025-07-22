package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Prompt holds the schema definition for the Prompt entity.
type Prompt struct {
	ent.Schema
}

// Mixin returns the mixins for the Prompt entity.
func (Prompt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseEntityMixin{},
	}
}

// Fields of the Prompt.
func (Prompt) Fields() []ent.Field {
	return []ent.Field{
		// Prompt-specific fields
		field.String("template_path").
			Optional().
			Nillable(),
		field.Text("template_content").
			Optional().
			Nillable(),
		field.JSON("variables_schema", map[string]interface{}{}).
			Optional(),
		field.Enum("prompt_type").
			Values("SYSTEM", "USER", "ASSISTANT", "FUNCTION", "TEMPLATE").
			Default("TEMPLATE"),
		field.Int("context_window").
			Optional().
			Nillable(),
		field.JSON("default_variables", map[string]interface{}{}).
			Optional(),
		field.JSON("required_variables", []string{}).
			Optional(),
		field.String("language").
			Default("en"),
		field.Bool("supports_streaming").
			Default(false),
		field.JSON("model_preferences", []string{}).
			Optional(),
		field.Float("temperature").
			Optional().
			Nillable(),
		field.Int("max_tokens").
			Optional().
			Nillable(),
		field.JSON("stop_sequences", []string{}).
			Optional(),
	}
}

// Edges of the Prompt.
func (Prompt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Ref("prompts").
			Unique().
			Required(),
		edge.To("dependencies", PromptDependency.Type),
	}
}

// Indexes of the Prompt.
func (Prompt) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("prompt_type"),
		index.Fields("is_installed"),
		index.Fields("stability"),
		index.Fields("language"),
		index.Fields("name", "version").
			Edges("repository").
			Unique(),
	}
}