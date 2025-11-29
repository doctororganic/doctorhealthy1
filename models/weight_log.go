package models

import "time"

// WeightLog represents a weight tracking entry
type WeightLog struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Weight    float64   `json:"weight" db:"weight"` // in kg
	Date      time.Time `json:"date" db:"date"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
