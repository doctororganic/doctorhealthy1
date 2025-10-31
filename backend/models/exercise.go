package models

import "time"

type Exercise struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListExercisesRequest represents a request to list exercises with filters
type ListExercisesRequest struct {
	Name           string `json:"name,omitempty"`
	Category       string `json:"category,omitempty"`
	MuscleGroup    string `json:"muscle_group,omitempty"`
	Equipment      string `json:"equipment,omitempty"`
	Difficulty     string `json:"difficulty,omitempty"`
	Page           int    `json:"page,omitempty" validate:"omitempty,min=1"`
	Limit          int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
}

// ExerciseListResponse represents a response for exercise listing
type ExerciseListResponse struct {
	Exercises []Exercise `json:"exercises"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
	HasNext   bool       `json:"has_next"`
}
