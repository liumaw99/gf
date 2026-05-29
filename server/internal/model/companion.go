package model

import "time"

// CompanionConfig 伙伴配置模型
type CompanionConfig struct {
	ID                 string    `json:"id" db:"id"`
	UserID             string    `json:"user_id" db:"user_id"`
	Name               string    `json:"name" db:"name"`
	AvatarType         string    `json:"avatar_type" db:"avatar_type"`
	PersonalityType    string    `json:"personality_type" db:"personality_type"`
	EncouragementStyle string    `json:"encouragement_style" db:"encouragement_style"`
	VoiceType          *string   `json:"voice_type,omitempty" db:"voice_type"`
	GrowthLevel        int       `json:"growth_level" db:"growth_level"`
	RelationshipScore  float64   `json:"relationship_score" db:"relationship_score"`
	CharacterID        *string   `json:"character_id,omitempty" db:"character_id"`
	CustomName         *string   `json:"custom_name,omitempty" db:"custom_name"`
	CustomAvatarURL    *string   `json:"custom_avatar_url,omitempty" db:"custom_avatar_url"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// UserCharacterConfig 用户角色配置模型
type UserCharacterConfig struct {
	ID                     string    `json:"id" db:"id"`
	UserID                 string    `json:"user_id" db:"user_id"`
	CharacterID            string    `json:"character_id" db:"character_id"`
	CustomName             *string   `json:"custom_name,omitempty" db:"custom_name"`
	CustomAvatarURL        *string   `json:"custom_avatar_url,omitempty" db:"custom_avatar_url"`
	IntimacyLevel          int       `json:"intimacy_level" db:"intimacy_level"`
	IntimacyScore          int       `json:"intimacy_score" db:"intimacy_score"`
	UnlockedStories        []string  `json:"unlocked_stories" db:"unlocked_stories"`
	UnlockedExpressions    []string  `json:"unlocked_expressions" db:"unlocked_expressions"`
	UnlockedDeepNightMode  bool      `json:"unlocked_deep_night_mode" db:"unlocked_deep_night_mode"`
	TotalMessages          int       `json:"total_messages" db:"total_messages"`
	TotalGoalsCompleted    int       `json:"total_goals_completed" db:"total_goals_completed"`
	IsActive               bool      `json:"is_active" db:"is_active"`
	IsFavorite             bool      `json:"is_favorite" db:"is_favorite"`
	SwitchCountThisMonth   int       `json:"switch_count_this_month" db:"switch_count_this_month"`
	LastSwitchDate         *time.Time `json:"last_switch_date,omitempty" db:"last_switch_date"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}
