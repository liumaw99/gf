package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// FocusSession holds the schema definition for the FocusSession entity.
type FocusSession struct {
	ent.Schema
}

// Fields of the FocusSession.
func (FocusSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.Time("start_time"),
		field.Int("duration").
			Positive(),
		field.Float("focus_quality_score").
			Optional().
			Nillable(),
		field.Strings("distractions").
			Default([]string{}),
		field.String("context").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the FocusSession.
func (FocusSession) Edges() []ent.Edge {
	return nil
}

// Indexes of the FocusSession.
func (FocusSession) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "start_time"),
	}
}
