package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DailyCompletion holds the schema definition for the DailyCompletion entity.
type DailyCompletion struct {
	ent.Schema
}

// Fields of the DailyCompletion.
func (DailyCompletion) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("habit_id", uuid.UUID{}),
		field.Time("date"),
		field.Time("completed_at").
			Default(time.Now),
		field.Int("quality_score").
			Optional().
			Nillable(),
		field.String("note").
			Optional().
			Nillable(),
	}
}

// Edges of the DailyCompletion.
func (DailyCompletion) Edges() []ent.Edge {
	return nil
}

// Indexes of the DailyCompletion.
func (DailyCompletion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "habit_id", "date").
			Unique(),
		index.Fields("user_id", "date"),
	}
}
