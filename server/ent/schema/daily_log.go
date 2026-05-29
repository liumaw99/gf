package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DailyLog holds the schema definition for the DailyLog entity.
type DailyLog struct {
	ent.Schema
}

// Fields of the DailyLog.
func (DailyLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.Time("date"),
		field.Int("mood_score").
			Optional().
			Nillable(),
		field.String("reflection_note").
			Optional().
			Nillable(),
		field.String("summary").
			Optional().
			Nillable(),
	}
}

// Edges of the DailyLog.
func (DailyLog) Edges() []ent.Edge {
	return nil
}

// Indexes of the DailyLog.
func (DailyLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "date").
			Unique(),
	}
}
