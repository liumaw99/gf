package model

import "time"

// ChatMessage 对话消息模型
type ChatMessage struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	CharacterID *string   `json:"character_id,omitempty" db:"character_id"`
	Role        string    `json:"role" db:"role"` // user / assistant / system
	Content     string    `json:"content" db:"content"`
	IsComplete  bool      `json:"is_complete" db:"is_complete"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
