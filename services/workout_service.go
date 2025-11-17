package services

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// WorkoutService handles workout-related business logic
type WorkoutService struct {
	db *sql.DB
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(db *sql.DB) *WorkoutService {
	return &WorkoutService{
		db: db,
	}
}

// GetWorkouts retrieves workouts for a user
func (s *WorkoutService) GetWorkouts(userID uint, limit, offset int) ([]models.Workout, error) {
	query := `
		SELECT id, user_id, name, description, duration, intensity,
			   calories_burned, workout_type, completed_at, created_at, updated_at
		FROM workouts 
		WHERE user_id = $1 
		ORDER BY completed_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var workouts []models.Workout
	for rows.Next() {
		var workout models.Workout
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&workout.ID, &workout.UserID, &workout.Name, &workout.Description,
			&workout.Duration, &workout.Intensity, &workout.CaloriesBurned,
			&workout.WorkoutType, &completedAt, &workout.CreatedAt, &workout.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if completedAt.Valid {
			workout.CompletedAt = &completedAt.Time
		}
		
		workouts = append(workouts, workout)
	}
	
	return workouts, nil
}

// CreateWorkout creates a new workout
func (s *WorkoutService) CreateWorkout(workout *models.Workout) error {
	query := `
		INSERT INTO workouts (user_id, name, description, duration, intensity,
							calories_burned, workout_type, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		workout.UserID, workout.Name, workout.Description, workout.Duration,
		workout.Intensity, workout.CaloriesBurned, workout.WorkoutType,
		workout.CompletedAt, now, now,
	).Scan(&workout.ID)
	
	if err != nil {
		return err
	}
	
	workout.CreatedAt = now
	workout.UpdatedAt = now
	
	return nil
}

// GetWorkout retrieves a workout by ID
func (s *WorkoutService) GetWorkout(id, userID uint) (*models.Workout, error) {
	query := `
		SELECT id, user_id, name, description, duration, intensity,
			   calories_burned, workout_type, completed_at, created_at, updated_at
		FROM workouts 
		WHERE id = $1 AND user_id = $2
	`
	
	var workout models.Workout
	var completedAt sql.NullTime
	
	err := s.db.QueryRow(query, id, userID).Scan(
		&workout.ID, &workout.UserID, &workout.Name, &workout.Description,
		&workout.Duration, &workout.Intensity, &workout.CaloriesBurned,
		&workout.WorkoutType, &completedAt, &workout.CreatedAt, &workout.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workout not found")
		}
		return nil, err
	}
	
	if completedAt.Valid {
		workout.CompletedAt = &completedAt.Time
	}
	
	return &workout, nil
}

// GetWorkoutPlans retrieves workout plans for a user
func (s *WorkoutService) GetWorkoutPlans(userID uint) ([]models.WorkoutPlan, error) {
	query := `
		SELECT id, user_id, name, description, duration_weeks, difficulty_level,
			   is_active, start_date, end_date, created_at, updated_at
		FROM workout_plans 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var plans []models.WorkoutPlan
	for rows.Next() {
		var plan models.WorkoutPlan
		var startDate, endDate sql.NullTime
		
		err := rows.Scan(
			&plan.ID, &plan.UserID, &plan.Name, &plan.Description,
			&plan.DurationWeeks, &plan.DifficultyLevel, &plan.IsActive,
			&startDate, &endDate, &plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if startDate.Valid {
			plan.StartDate = &startDate.Time
		}
		if endDate.Valid {
			plan.EndDate = &endDate.Time
		}
		
		plans = append(plans, plan)
	}
	
	return plans, nil
}

// CreateWorkoutPlan creates a new workout plan
func (s *WorkoutService) CreateWorkoutPlan(plan *models.WorkoutPlan) error {
	query := `
		INSERT INTO workout_plans (user_id, name, description, duration_weeks,
								 difficulty_level, is_active, start_date, end_date,
								 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		plan.UserID, plan.Name, plan.Description, plan.DurationWeeks,
		plan.DifficultyLevel, plan.IsActive, plan.StartDate,
		plan.EndDate, now, now,
	).Scan(&plan.ID)
	
	if err != nil {
		return err
	}
	
	plan.CreatedAt = now
	plan.UpdatedAt = now
	
	return nil
}
