package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nutrition-platform/models"
)

type WorkoutPlanRepository struct {
	db *sql.DB
}

func NewWorkoutPlanRepository(db *sql.DB) *WorkoutPlanRepository {
	return &WorkoutPlanRepository{db: db}
}

// CreateWorkoutPlan creates a new workout plan
func (r *WorkoutPlanRepository) CreateWorkoutPlan(ctx context.Context, plan *models.WorkoutPlan) error {
	query := `
		INSERT INTO workout_plans (
			name, description, exercises, user_id, is_public, is_template,
			duration_weeks, difficulty, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		plan.Name,
		plan.Description,
		models.WorkoutPlanExercises(plan.Exercises).Value(),
		plan.UserID,
		plan.IsPublic,
		plan.IsTemplate,
		plan.DurationWeeks,
		plan.Difficulty,
		now,
		now,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create workout plan: %w", err)
	}

	return nil
}

// GetWorkoutPlanByID retrieves a workout plan by ID
func (r *WorkoutPlanRepository) GetWorkoutPlanByID(ctx context.Context, id int64) (*models.WorkoutPlan, error) {
	query := `
		SELECT id, name, description, exercises, user_id, is_public, is_template,
			   duration_weeks, difficulty, created_at, updated_at
		FROM workout_plans 
		WHERE id = $1`

	plan := &models.WorkoutPlan{}
	var exercisesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&exercisesJSON,
		&plan.UserID,
		&plan.IsPublic,
		&plan.IsTemplate,
		&plan.DurationWeeks,
		&plan.Difficulty,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workout plan not found")
		}
		return nil, fmt.Errorf("failed to get workout plan: %w", err)
	}

	// Scan exercises JSON
	exercises := models.WorkoutPlanExercises{}
	if err := exercises.Scan(exercisesJSON); err != nil {
		return nil, fmt.Errorf("failed to scan exercises: %w", err)
	}
	plan.Exercises = exercises

	return plan, nil
}

// UpdateWorkoutPlan updates an existing workout plan
func (r *WorkoutPlanRepository) UpdateWorkoutPlan(ctx context.Context, plan *models.WorkoutPlan) error {
	query := `
		UPDATE workout_plans 
		SET name = $1, description = $2, exercises = $3, is_public = $4,
			is_template = $5, duration_weeks = $6, difficulty = $7, updated_at = $8
		WHERE id = $9
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		plan.Name,
		plan.Description,
		models.WorkoutPlanExercises(plan.Exercises).Value(),
		plan.IsPublic,
		plan.IsTemplate,
		plan.DurationWeeks,
		plan.Difficulty,
		now,
		plan.ID,
	).Scan(&plan.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("workout plan not found")
		}
		return fmt.Errorf("failed to update workout plan: %w", err)
	}

	return nil
}

// DeleteWorkoutPlan soft deletes a workout plan
func (r *WorkoutPlanRepository) DeleteWorkoutPlan(ctx context.Context, id int64) error {
	query := `UPDATE workout_plans SET deleted_at = $1 WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete workout plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workout plan not found")
	}

	return nil
}

// ListWorkoutPlans retrieves workout plans with filtering and pagination
func (r *WorkoutPlanRepository) ListWorkoutPlans(ctx context.Context, req *models.ListWorkoutPlansRequest) (*models.WorkoutPlanListResponse, error) {
	whereConditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if req.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.IsPublic != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_public = $%d", argIndex))
		args = append(args, *req.IsPublic)
		argIndex++
	}

	if req.IsTemplate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_template = $%d", argIndex))
		args = append(args, *req.IsTemplate)
		argIndex++
	}

	if req.Difficulty != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("difficulty = $%d", argIndex))
		args = append(args, req.Difficulty)
		argIndex++
	}

	// Add search
	if req.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
		argIndex += 2
	}

	whereClause := "WHERE " + strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM workout_plans %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count workout plans: %w", err)
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, name, description, exercises, user_id, is_public, is_template,
			   duration_weeks, difficulty, created_at, updated_at
		FROM workout_plans 
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause,
		req.SortBy,
		argIndex,
		argIndex+1,
	)

	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query workout plans: %w", err)
	}
	defer rows.Close()

	plans := []*models.WorkoutPlan{}
	for rows.Next() {
		plan := &models.WorkoutPlan{}
		var exercisesJSON []byte

		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&exercisesJSON,
			&plan.UserID,
			&plan.IsPublic,
			&plan.IsTemplate,
			&plan.DurationWeeks,
			&plan.Difficulty,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workout plan: %w", err)
		}

		// Scan exercises JSON
		exercises := models.WorkoutPlanExercises{}
		if err := exercises.Scan(exercisesJSON); err != nil {
			return nil, fmt.Errorf("failed to scan exercises: %w", err)
		}
		plan.Exercises = exercises

		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workout plan rows: %w", err)
	}

	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)

	return &models.WorkoutPlanListResponse{
		WorkoutPlans: plans,
		Total:        total,
		Page:         int(req.Offset/req.Limit) + 1,
		Limit:        req.Limit,
		TotalPages:   int(totalPages),
	}, nil
}

// DuplicateWorkoutPlan creates a copy of an existing workout plan
func (r *WorkoutPlanRepository) DuplicateWorkoutPlan(ctx context.Context, originalID, newUserID int64, newName string) (*models.WorkoutPlan, error) {
	// Get original plan
	original, err := r.GetWorkoutPlanByID(ctx, originalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original workout plan: %w", err)
	}

	// Create duplicate
	duplicate := &models.WorkoutPlan{
		Name:          newName,
		Description:   original.Description,
		Exercises:     original.Exercises,
		UserID:        newUserID,
		IsPublic:      false,
		IsTemplate:    false,
		DurationWeeks: original.DurationWeeks,
		Difficulty:    original.Difficulty,
	}

	err = r.CreateWorkoutPlan(ctx, duplicate)
	if err != nil {
		return nil, fmt.Errorf("failed to create duplicate workout plan: %w", err)
	}

	return duplicate, nil
}

// GetWorkoutPlanStats returns statistics for a workout plan
func (r *WorkoutPlanRepository) GetWorkoutPlanStats(ctx context.Context, planID int64) (*models.WorkoutPlanStats, error) {
	stats := &models.WorkoutPlanStats{}

	// Get total workout logs for this plan
	query := `
		SELECT 
			COUNT(*) as total_workouts,
			COUNT(DISTINCT DATE(workout_date)) as unique_days,
			COALESCE(AVG(duration_minutes), 0) as avg_duration,
			COALESCE(SUM(duration_minutes), 0) as total_duration
		FROM workout_logs 
		WHERE workout_plan_id = $1`

	err := r.db.QueryRowContext(ctx, query, planID).Scan(
		&stats.TotalWorkouts,
		&stats.UniqueDays,
		&stats.AvgDuration,
		&stats.TotalDuration,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get workout plan stats: %w", err)
	}

	// Get last workout date
	lastWorkoutQuery := `
		SELECT MAX(workout_date) 
		FROM workout_logs 
		WHERE workout_plan_id = $1`

	err = r.db.QueryRowContext(ctx, lastWorkoutQuery, planID).Scan(&stats.LastWorkout)
	if err != nil {
		return nil, fmt.Errorf("failed to get last workout date: %w", err)
	}

	return stats, nil
}

// GetUserWorkoutPlanHistory retrieves user's workout plan history
func (r *WorkoutPlanRepository) GetUserWorkoutPlanHistory(ctx context.Context, userID int64, limit, offset int) ([]*models.WorkoutPlanHistory, error) {
	query := `
		SELECT 
			wp.id, wp.name, wl.workout_date, wl.duration_minutes,
			COUNT(wl.id) OVER (PARTITION BY wp.id) as workout_count,
			MAX(wl.workout_date) OVER (PARTITION BY wp.id) as last_workout
		FROM workout_plans wp
		LEFT JOIN workout_logs wl ON wp.id = wl.workout_plan_id
		WHERE wp.user_id = $1 AND wp.deleted_at IS NULL
		ORDER BY last_workout DESC NULLS LAST, wp.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query workout plan history: %w", err)
	}
	defer rows.Close()

	history := []*models.WorkoutPlanHistory{}
	for rows.Next() {
		item := &models.WorkoutPlanHistory{}
		
		err := rows.Scan(
			&item.PlanID,
			&item.PlanName,
			&item.LastWorkoutDate,
			&item.LastWorkoutDuration,
			&item.WorkoutCount,
			&item.LastWorkoutDate,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workout plan history: %w", err)
		}

		history = append(history, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workout plan history rows: %w", err)
	}

	return history, nil
}

// GetPopularWorkoutPlans retrieves popular workout plans
func (r *WorkoutPlanRepository) GetPopularWorkoutPlans(ctx context.Context, limit int) ([]*models.WorkoutPlan, error) {
	query := `
		SELECT wp.id, wp.name, wp.description, wp.exercises, wp.user_id, 
			   wp.is_public, wp.is_template, wp.duration_weeks, wp.difficulty,
			   wp.created_at, wp.updated_at,
			   COUNT(wl.id) as usage_count
		FROM workout_plans wp
		LEFT JOIN workout_logs wl ON wp.id = wl.workout_plan_id
		WHERE wp.is_public = true AND wp.deleted_at IS NULL
		GROUP BY wp.id, wp.name, wp.description, wp.exercises, wp.user_id,
				 wp.is_public, wp.is_template, wp.duration_weeks, wp.difficulty,
				 wp.created_at, wp.updated_at
		ORDER BY usage_count DESC, wp.created_at DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query popular workout plans: %w", err)
	}
	defer rows.Close()

	plans := []*models.WorkoutPlan{}
	for rows.Next() {
		plan := &models.WorkoutPlan{}
		var exercisesJSON []byte
		var usageCount int64

		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&exercisesJSON,
			&plan.UserID,
			&plan.IsPublic,
			&plan.IsTemplate,
			&plan.DurationWeeks,
			&plan.Difficulty,
			&plan.CreatedAt,
			&plan.UpdatedAt,
			&usageCount,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan popular workout plan: %w", err)
		}

		// Scan exercises JSON
		exercises := models.WorkoutPlanExercises{}
		if err := exercises.Scan(exercisesJSON); err != nil {
			return nil, fmt.Errorf("failed to scan exercises: %w", err)
		}
		plan.Exercises = exercises

		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating popular workout plan rows: %w", err)
	}

	return plans, nil
}
