package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("apple_id").
			Unique().
			Optional().
			Nillable(),
		field.String("google_id").
			Unique().
			Optional().
			Nillable(),
		field.String("email").
			Optional().
			Nillable(),
		field.String("password_hash").
			Optional().
			Nillable().
			Sensitive(),
		field.String("name").
			Optional().
			Nillable(),
		field.String("avatar_url").
			Optional().
			Nillable(),
		field.String("subscription_tier").
			Default("free"),
		field.Time("subscription_until").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(func() time.Time { return time.Now() }),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// ent.Edge("companion_config", CompanionConfig.Type).
		// 	Unique().
		// 	Ref("user"),
		// ent.Edge("habits", Habit.Type).
		// 	Ref("user"),
		// ent.Edge("daily_logs", DailyLog.Type).
		// 	Ref("user"),
		// ent.Edge("character_configs", UserCharacterConfig.Type).
		// 	Ref("user"),
		// ent.Edge("chat_messages", ChatMessage.Type).
		// 	Ref("user"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("apple_id").
			Unique(),
		index.Fields("google_id").
			Unique(),
		index.Fields("email").
			Unique(),
	}
}
