package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserCharacterConfig holds the schema definition for the UserCharacterConfig entity.
type UserCharacterConfig struct {
	ent.Schema
}

// Fields of the UserCharacterConfig.
func (UserCharacterConfig) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("character_id").
			NotEmpty(),
		field.String("custom_name").
			Optional().
			Nillable(),
		field.String("custom_avatar_url").
			Optional().
			Nillable(),
		field.Int("intimacy_level").
			Default(1),
		field.Int("intimacy_score").
			Default(0),
		field.Strings("unlocked_stories").
			Default([]string{}),
		field.Strings("unlocked_expressions").
			Default([]string{}),
		field.Bool("unlocked_deep_night_mode").
			Default(false),
		field.Int("total_messages").
			Default(0),
		field.Int("total_goals_completed").
			Default(0),
		field.Bool("is_active").
			Default(true),
		field.Bool("is_favorite").
			Default(false),
		field.Int("switch_count_this_month").
			Default(0),
		field.Time("last_switch_date").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(func() time.Time { return time.Now() }),
	}
}

// Edges of the UserCharacterConfig.
func (UserCharacterConfig) Edges() []ent.Edge {
	return nil
}

// Indexes of the UserCharacterConfig.
func (UserCharacterConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("user_id", "is_active"),
		index.Fields("user_id", "character_id").
			Unique(),
	}
}
