package model

import "time"

// Habit 习惯模型
type Habit struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	Category     string    `json:"category" db:"category"`
	Description  *string   `json:"description,omitempty" db:"description"`
	Priority     int       `json:"priority" db:"priority"`
	FrequencyType string   `json:"frequency_type" db:"frequency_type"`
	TargetDays   []int     `json:"target_days" db:"target_days"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CurrentPhase int       `json:"current_phase" db:"current_phase"`
	HSI          float64   `json:"hsi" db:"hsi"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// MicroAction 微行动模型
type MicroAction struct {
	ID                string    `json:"id" db:"id"`
	HabitID           string    `json:"habit_id" db:"habit_id"`
	Phase             int       `json:"phase" db:"phase"`
	DayInPhase        int       `json:"day_in_phase" db:"day_in_phase"`
	ActionText        string    `json:"action_text" db:"action_text"`
	EstimatedDuration *int      `json:"estimated_duration,omitempty" db:"estimated_duration"`
	IsDefault         bool      `json:"is_default" db:"is_default"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
