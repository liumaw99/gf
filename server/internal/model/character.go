package model

import (
	"encoding/json"
	"time"
)

// CelebrityCharacter 名人角色模型
type CelebrityCharacter struct {
	ID              string          `json:"id" db:"id"`
	Name            string          `json:"name" db:"name"`
	NameEn          *string         `json:"name_en,omitempty" db:"name_en"`
	Category        string          `json:"category" db:"category"`
	Era             *string         `json:"era,omitempty" db:"era"`
	Nationality     *string         `json:"nationality,omitempty" db:"nationality"`
	FamousFor       []string        `json:"famous_for" db:"famous_for"`
	AvatarURL       *string         `json:"avatar_url,omitempty" db:"avatar_url"`
	AvatarLottieURL *string         `json:"avatar_lottie_url,omitempty" db:"avatar_lottie_url"`
	StyleTags       []string        `json:"style_tags" db:"style_tags"`
	Distillate      json.RawMessage `json:"distillate" db:"distillate"`
	QualityScore    float64         `json:"quality_score" db:"quality_score"`
	Status          string          `json:"status" db:"status"`
	IsPremium       bool            `json:"is_premium" db:"is_premium"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}
