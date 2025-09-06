package services

import "time"

// Metadata represents common metadata for all data structures
type Metadata struct {
	TotalCount  int       `json:"total_count"`
	LastUpdated string    `json:"last_updated"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewMetadata creates a new metadata instance
func NewMetadata() Metadata {
	return Metadata{
		TotalCount:  0,
		LastUpdated: time.Now().Format(time.RFC3339),
		Version:     "1.0",
		CreatedAt:   time.Now(),
	}
}