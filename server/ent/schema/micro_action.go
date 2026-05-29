package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MicroAction holds the schema definition for the MicroAction entity.
type MicroAction struct {
	ent.Schema
}

// Fields of the MicroAction.
func (MicroAction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("habit_id", uuid.UUID{}),
		field.Int("phase").
			Positive(),
		field.Int("day_in_phase").
			Positive(),
		field.String("action_text").
			NotEmpty(),
		field.Int("estimated_duration").
			Optional().
			Nillable(),
		field.Bool("is_default").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the MicroAction.
func (MicroAction) Edges() []ent.Edge {
	return nil
}

// Indexes of the MicroAction.
func (MicroAction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("habit_id"),
	}
}
