package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CompanionConfig holds the schema definition for the CompanionConfig entity.
type CompanionConfig struct {
	ent.Schema
}

// Fields of the CompanionConfig.
func (CompanionConfig) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Unique(),
		field.String("name").
			Default("Lagom"),
		field.String("avatar_type").
			Default("fox"),
		field.String("personality_type").
			Default("gentle"),
		field.String("encouragement_style").
			Default("soft"),
		field.String("voice_type").
			Optional().
			Nillable(),
		field.Int("growth_level").
			Default(1),
		field.Float("relationship_score").
			Default(0.0),
		field.String("character_id").
			Optional().
			Nillable(),
		field.String("custom_name").
			Optional().
			Nillable(),
		field.String("custom_avatar_url").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(func() time.Time { return time.Now() }),
	}
}

// Edges of the CompanionConfig.
func (CompanionConfig) Edges() []ent.Edge {
	return nil
}

// Indexes of the CompanionConfig.
func (CompanionConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id").
			Unique(),
	}
}
