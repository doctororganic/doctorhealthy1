package models

import "time"

type Exercise struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
