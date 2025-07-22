package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SyncOperation holds the schema definition for the SyncOperation entity.
type SyncOperation struct {
	ent.Schema
}

// Fields of the SyncOperation.
func (SyncOperation) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.Enum("type").
			Values("PULL", "PUSH", "CLONE", "FORK", "MERGE"),
		field.Enum("status").
			Values("PENDING", "RUNNING", "COMPLETED", "FAILED", "CANCELLED").
			Default("PENDING"),
		field.Enum("direction").
			Values("UPSTREAM_TO_LOCAL", "LOCAL_TO_UPSTREAM", "BIDIRECTIONAL"),
		field.String("repository_id").
			Optional().
			Nillable(),
		field.String("component_id").
			Optional().
			Nillable(),
		field.String("source_commit").
			Optional().
			Nillable(),
		field.String("target_commit").
			Optional().
			Nillable(),
		field.String("branch").
			Optional().
			Nillable(),
		field.Time("started_at").
			Default(time.Now).
			Immutable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
		field.Text("error_message").
			Optional().
			Nillable(),
		field.Int("total_steps").
			Default(1),
		field.Int("completed_steps").
			Default(0),
	}
}

// Edges of the SyncOperation.
func (SyncOperation) Edges() []ent.Edge {
	return nil
}

// Indexes of the SyncOperation.
func (SyncOperation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("type"),
		index.Fields("started_at"),
	}
}
