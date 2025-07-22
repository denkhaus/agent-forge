package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tool holds the schema definition for the Tool entity.
type Tool struct {
	ent.Schema
}

// Mixin returns the mixins for the Tool entity.
func (Tool) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseEntityMixin{},
	}
}

// Fields of the Tool.
func (Tool) Fields() []ent.Field {
	return []ent.Field{
		// Tool-specific fields
		field.Enum("execution_type").
			Values("MCP", "HTTP", "BINARY", "FUNCTION").
			Default("MCP"),
		field.String("schema_path").
			Optional().
			Nillable(),
		field.JSON("server_config", map[string]interface{}{}).
			Optional(),
		field.JSON("capabilities", []string{}).
			Optional(),
		field.String("entry_point").
			Optional().
			Nillable(),
		field.JSON("environment_variables", map[string]string{}).
			Optional(),
		field.JSON("required_permissions", []string{}).
			Optional(),
		field.Int("timeout_seconds").
			Default(30),
		field.Bool("supports_streaming").
			Default(false),
		field.JSON("input_schema", map[string]interface{}{}).
			Optional(),
		field.JSON("output_schema", map[string]interface{}{}).
			Optional(),
	}
}

// Edges of the Tool.
func (Tool) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Ref("tools").
			Unique().
			Required(),
		edge.To("dependencies", ToolDependency.Type),
	}
}

// Indexes of the Tool.
func (Tool) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("execution_type"),
		index.Fields("is_installed"),
		index.Fields("stability"),
		index.Fields("name", "version").
			Edges("repository").
			Unique(),
	}
}