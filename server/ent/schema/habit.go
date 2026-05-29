package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Habit holds the schema definition for the Habit entity.
type Habit struct {
	ent.Schema
}

// Fields of the Habit.
func (Habit) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("title").
			NotEmpty(),
		field.String("category").
			NotEmpty(),
		field.String("description").
			Optional().
			Nillable(),
		field.Int("priority").
			Default(3),
		field.String("frequency_type").
			Default("daily"),
		field.Ints("target_days").
			Default([]int{}),
		field.Bool("is_active").
			Default(true),
		field.Int("current_phase").
			Default(1),
		field.Float("hsi").
			Default(0.0),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(func() time.Time { return time.Now() }),
	}
}

// Edges of the Habit.
func (Habit) Edges() []ent.Edge {
	return nil
}

// Indexes of the Habit.
func (Habit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("user_id", "is_active"),
	}
}
