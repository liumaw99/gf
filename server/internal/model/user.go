package model

import "time"

// User 用户模型
type User struct {
	ID                string    `json:"id" db:"id"`
	AppleID           *string   `json:"apple_id,omitempty" db:"apple_id"`
	GoogleID          *string   `json:"google_id,omitempty" db:"google_id"`
	Email             *string   `json:"email,omitempty" db:"email"`
	Name              *string   `json:"name,omitempty" db:"name"`
	AvatarURL         *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	SubscriptionTier  string    `json:"subscription_tier" db:"subscription_tier"`
	SubscriptionUntil *time.Time `json:"subscription_until,omitempty" db:"subscription_until"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
