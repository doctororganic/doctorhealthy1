package services

import "time"

// Metadata represents metadata for data collections
type Metadata struct {
	Version   string    `json:"version"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMetadata creates a new metadata instance
func NewMetadata() Metadata {
	now := time.Now()
	return Metadata{
		Version:   "1.0",
		Count:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type UserService struct{}
type FoodService struct{}
type LogService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func NewFoodService() *FoodService {
	return &FoodService{}
}

func NewLogService() *LogService {
	return &LogService{}
}
