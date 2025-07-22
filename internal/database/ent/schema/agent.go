package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Agent holds the schema definition for the Agent entity.
type Agent struct {
	ent.Schema
}

// Mixin returns the mixins for the Agent entity.
func (Agent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseEntityMixin{},
	}
}

// Fields of the Agent.
func (Agent) Fields() []ent.Field {
	return []ent.Field{
		// Agent-specific fields
		field.String("config_path").
			Optional().
			Nillable(),
		field.JSON("agent_config", map[string]interface{}{}).
			Optional(),
		field.String("llm_provider").
			Optional().
			Nillable(),
		field.String("system_prompt_id").
			Optional().
			Nillable(),
		field.JSON("tool_dependencies", []string{}).
			Optional(),
		field.JSON("prompt_dependencies", []string{}).
			Optional(),
		field.JSON("agent_dependencies", []string{}).
			Optional(),
		field.Enum("agent_type").
			Values("CONVERSATIONAL", "TASK_ORIENTED", "SPECIALIZED", "COMPOSITE").
			Default("CONVERSATIONAL"),
		field.JSON("capabilities", []string{}).
			Optional(),
		field.JSON("supported_languages", []string{}).
			Optional(),
		field.Bool("supports_memory").
			Default(false),
		field.Bool("supports_tools").
			Default(false),
		field.Bool("supports_multimodal").
			Default(false),
		field.JSON("model_preferences", []string{}).
			Optional(),
		field.Float("default_temperature").
			Optional().
			Nillable(),
		field.Int("default_max_tokens").
			Optional().
			Nillable(),
		field.Int("session_timeout_minutes").
			Default(30),
	}
}

// Edges of the Agent.
func (Agent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Ref("agents").
			Unique().
			Required(),
		edge.To("dependencies", AgentDependency.Type),
	}
}

// Indexes of the Agent.
func (Agent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("agent_type"),
		index.Fields("llm_provider"),
		index.Fields("is_installed"),
		index.Fields("stability"),
		index.Fields("supports_tools"),
		index.Fields("name", "version").
			Edges("repository").
			Unique(),
	}
}