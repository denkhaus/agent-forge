package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// LocalConfig holds the schema definition for the LocalConfig entity.
type LocalConfig struct {
	ent.Schema
}

// Fields of the LocalConfig.
func (LocalConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("key").
			Unique(),
		field.Text("value"),
		field.Enum("type").
			Values("STRING", "INTEGER", "BOOLEAN", "JSON").
			Default("STRING"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the LocalConfig.
func (LocalConfig) Edges() []ent.Edge {
	return nil
}

// Indexes of the LocalConfig.
func (LocalConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
	}
}
