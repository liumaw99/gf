package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ChatMessage holds the schema definition for the ChatMessage entity.
type ChatMessage struct {
	ent.Schema
}

// Fields of the ChatMessage.
func (ChatMessage) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("character_id").
			Optional().
			Nillable(),
		field.String("role").
			NotEmpty(),
		field.String("content").
			NotEmpty(),
		field.Bool("is_complete").
			Default(true),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the ChatMessage.
func (ChatMessage) Edges() []ent.Edge {
	return nil
}

// Indexes of the ChatMessage.
func (ChatMessage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "created_at"),
		index.Fields("user_id", "character_id", "created_at"),
	}
}
