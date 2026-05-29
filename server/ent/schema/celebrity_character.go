package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CelebrityCharacter holds the schema definition for the CelebrityCharacter entity.
type CelebrityCharacter struct {
	ent.Schema
}

// Fields of the CelebrityCharacter.
func (CelebrityCharacter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(50).
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.String("name_en").
			Optional().
			Nillable(),
		field.String("category").
			NotEmpty(),
		field.String("era").
			Optional().
			Nillable(),
		field.String("nationality").
			Optional().
			Nillable(),
		field.Strings("famous_for").
			Default([]string{}),
		field.String("avatar_url").
			Optional().
			Nillable(),
		field.String("avatar_lottie_url").
			Optional().
			Nillable(),
		field.Strings("style_tags").
			Default([]string{}),
		field.JSON("distillate", map[string]any{}).
			Default(map[string]any{}),
		field.Float("quality_score").
			Default(0.0),
		field.String("status").
			Default("active"),
		field.Bool("is_premium").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(func() time.Time { return time.Now() }),
	}
}

// Edges of the CelebrityCharacter.
func (CelebrityCharacter) Edges() []ent.Edge {
	return nil
}

// Indexes of the CelebrityCharacter.
func (CelebrityCharacter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("category"),
		index.Fields("status"),
	}
}
