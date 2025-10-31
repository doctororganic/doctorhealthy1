package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nutrition-platform/models"
)

type WorkoutLogRepository struct {
	db *sql.DB
}

func NewWorkoutLogRepository(db *sql.DB) *WorkoutLogRepository {
	return &WorkoutLogRepository{db: db}
}

// CreateWorkoutLog creates a new workout log
func (r *WorkoutLogRepository) CreateWorkoutLog(ctx context.Context, log *models.WorkoutLog) error {
	query := `
		INSERT INTO workout_logs (
			user_id, workout_plan_id, workout_date, duration_minutes,
			completed_exercises, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		log.UserID,
		log.WorkoutPlanID,
		log.WorkoutDate,
		log.DurationMinutes,
		models.CompletedExercises(log.CompletedExercises).Value(),
		log.Notes,
		now,
		now,
	).Scan(&log.ID, &log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create workout log: %w", err)
	}

	return nil
}

// GetWorkoutLogByID retrieves a workout log by ID
func (r *WorkoutLogRepository) GetWorkoutLogByID(ctx context.Context, id int64) (*models.WorkoutLog, error) {
	query := `
		SELECT id, user_id, workout_plan_id, workout_date, duration_minutes,
			   completed_exercises, notes, created_at, updated_at
		FROM workout_logs 
		WHERE id = $1`

	log := &models.WorkoutLog{}
	var completedExercisesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.UserID,
		&log.WorkoutPlanID,
		&log.WorkoutDate,
		&log.DurationMinutes,
		&completedExercisesJSON,
		&log.Notes,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workout log not found")
		}
		return nil, fmt.Errorf("failed to get workout log: %w", err)
	}

	// Scan completed exercises JSON
	completedExercises := models.CompletedExercises{}
	if err := completedExercises.Scan(completedExercisesJSON); err != nil {
		return nil, fmt.Errorf("failed to scan completed exercises: %w", err)
	}
	log.CompletedExercises = completedExercises

	return log, nil
}

// UpdateWorkoutLog updates an existing workout log
func (r *WorkoutLogRepository) UpdateWorkoutLog(ctx context.Context, log *models.WorkoutLog) error {
	query := `
		UPDATE workout_logs 
		SET workout_plan_id = $1, workout_date = $2, duration_minutes = $3,
			completed_exercises = $4, notes = $5, updated_at = $6
		WHERE id = $7
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		log.WorkoutPlanID,
		log.WorkoutDate,
		log.DurationMinutes,
		models.CompletedExercises(log.CompletedExercises).Value(),
		log.Notes,
		now,
		log.ID,
	).Scan(&log.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("workout log not found")
		}
		return fmt.Errorf("failed to update workout log: %w", err)
	}

	return nil
}

// DeleteWorkoutLog deletes a workout log
func (r *WorkoutLogRepository) DeleteWorkoutLog(ctx context.Context, id int64) error {
	query := `DELETE FROM workout_logs WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workout log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workout log not found")
	}

	return nil
}

// ListWorkoutLogs retrieves workout logs with filtering and pagination
func (r *WorkoutLogRepository) ListWorkoutLogs(ctx context.Context, req *models.ListWorkoutLogsRequest) (*models.WorkoutLogListResponse, error) {
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if req.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("wl.user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.WorkoutPlanID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("wl.workout_plan_id = $%d", argIndex))
		args = append(args, *req.WorkoutPlanID)
		argIndex++
	}

	if req.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("wl.workout_date >= $%d", argIndex))
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("wl.workout_date <= $%d", argIndex))
		args = append(args, *req.EndDate)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM workout_logs wl %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count workout logs: %w", err)
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT wl.id, wl.user_id, wl.workout_plan_id, wl.workout_date, wl.duration_minutes,
			   wl.completed_exercises, wl.notes, wl.created_at, wl.updated_at,
			   wp.name as workout_plan_name
		FROM workout_logs wl
		LEFT JOIN workout_plans wp ON wl.workout_plan_id = wp.id
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
		return nil, fmt.Errorf("failed to query workout logs: %w", err)
	}
	defer rows.Close()

	logs := []*models.WorkoutLogWithPlan{}
	for rows.Next() {
		log := &models.WorkoutLogWithPlan{}
		var completedExercisesJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.WorkoutPlanID,
			&log.WorkoutDate,
			&log.DurationMinutes,
			&completedExercisesJSON,
			&log.Notes,
			&log.CreatedAt,
			&log.UpdatedAt,
			&log.WorkoutPlanName,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workout log: %w", err)
		}

		// Scan completed exercises JSON
		completedExercises := models.CompletedExercises{}
		if err := completedExercises.Scan(completedExercisesJSON); err != nil {
			return nil, fmt.Errorf("failed to scan completed exercises: %w", err)
		}
		log.CompletedExercises = completedExercises

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workout log rows: %w", err)
	}

	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)

	return &models.WorkoutLogListResponse{
		WorkoutLogs: logs,
		Total:       total,
		Page:        int(req.Offset/req.Limit) + 1,
		Limit:       req.Limit,
		TotalPages:  int(totalPages),
	}, nil
}

// GetUserWorkoutStats retrieves workout statistics for a user
func (r *WorkoutLogRepository) GetUserWorkoutStats(ctx context.Context, userID int64, period string) (*models.UserWorkoutStats, error) {
	stats := &models.UserWorkoutStats{}

	// Set date range based on period
	var startDate time.Time
	now := time.Now()
	
	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0) // Default to month
	}

	// Get basic stats
	query := `
		SELECT 
			COUNT(*) as total_workouts,
			COUNT(DISTINCT DATE(workout_date)) as unique_days,
			COALESCE(SUM(duration_minutes), 0) as total_duration,
			COALESCE(AVG(duration_minutes), 0) as avg_duration
		FROM workout_logs 
		WHERE user_id = $1 AND workout_date >= $2`

	err := r.db.QueryRowContext(ctx, query, userID, startDate).Scan(
		&stats.TotalWorkouts,
		&stats.UniqueDays,
		&stats.TotalDuration,
		&stats.AvgDuration,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get workout stats: %w", err)
	}

	// Get current streak
	streakQuery := `
		WITH dated_workouts AS (
			SELECT DISTINCT DATE(workout_date) as workout_date
			FROM workout_logs 
			WHERE user_id = $1 AND workout_date >= $2
			ORDER BY workout_date DESC
		),
		streak_groups AS (
			SELECT 
				workout_date,
				ROW_NUMBER() OVER (ORDER BY workout_date DESC) as rn,
				workout_date - INTERVAL '1 day' * ROW_NUMBER() OVER (ORDER BY workout_date DESC) as grp
			FROM dated_workouts
		)
		SELECT COUNT(*) as current_streak
		FROM streak_groups
		WHERE grp = (SELECT grp FROM streak_groups ORDER BY workout_date DESC LIMIT 1)`

	err = r.db.QueryRowContext(ctx, streakQuery, userID, startDate).Scan(&stats.CurrentStreak)
	if err != nil {
		return nil, fmt.Errorf("failed to get current streak: %w", err)
	}

	// Get personal records count
	prQuery := `
		SELECT COUNT(*) 
		FROM personal_records 
		WHERE user_id = $1 AND achieved_at >= $2`

	err = r.db.QueryRowContext(ctx, prQuery, userID, startDate).Scan(&stats.PersonalRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal records count: %w", err)
	}

	return stats, nil
}

// GetWorkoutCalendar retrieves workout calendar data for a user
func (r *WorkoutLogRepository) GetWorkoutCalendar(ctx context.Context, userID int64, year, month int) ([]*models.WorkoutCalendarDay, error) {
	// Set date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	query := `
		SELECT 
			DATE(workout_date) as workout_date,
			COUNT(*) as workout_count,
			COALESCE(SUM(duration_minutes), 0) as total_duration,
			COALESCE(AVG(duration_minutes), 0) as avg_duration
		FROM workout_logs 
		WHERE user_id = $1 AND workout_date >= $2 AND workout_date < $3
		GROUP BY DATE(workout_date)
		ORDER BY workout_date`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query workout calendar: %w", err)
	}
	defer rows.Close()

	calendar := []*models.WorkoutCalendarDay{}
	for rows.Next() {
		day := &models.WorkoutCalendarDay{}
		
		err := rows.Scan(
			&day.Date,
			&day.WorkoutCount,
			&day.TotalDuration,
			&day.AvgDuration,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workout calendar day: %w", err)
		}

		calendar = append(calendar, day)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workout calendar rows: %w", err)
	}

	return calendar, nil
}

// GetRecentWorkoutLogs retrieves recent workout logs for a user
func (r *WorkoutLogRepository) GetRecentWorkoutLogs(ctx context.Context, userID int64, limit int) ([]*models.WorkoutLogWithPlan, error) {
	query := `
		SELECT wl.id, wl.user_id, wl.workout_plan_id, wl.workout_date, wl.duration_minutes,
			   wl.completed_exercises, wl.notes, wl.created_at, wl.updated_at,
			   wp.name as workout_plan_name
		FROM workout_logs wl
		LEFT JOIN workout_plans wp ON wl.workout_plan_id = wp.id
		WHERE wl.user_id = $1
		ORDER BY wl.workout_date DESC, wl.created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent workout logs: %w", err)
	}
	defer rows.Close()

	logs := []*models.WorkoutLogWithPlan{}
	for rows.Next() {
		log := &models.WorkoutLogWithPlan{}
		var completedExercisesJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.WorkoutPlanID,
			&log.WorkoutDate,
			&log.DurationMinutes,
			&completedExercisesJSON,
			&log.Notes,
			&log.CreatedAt,
			&log.UpdatedAt,
			&log.WorkoutPlanName,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan recent workout log: %w", err)
		}

		// Scan completed exercises JSON
		completedExercises := models.CompletedExercises{}
		if err := completedExercises.Scan(completedExercisesJSON); err != nil {
			return nil, fmt.Errorf("failed to scan completed exercises: %w", err)
		}
		log.CompletedExercises = completedExercises

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recent workout log rows: %w", err)
	}

	return logs, nil
}

// GetWorkoutLogsByDateRange retrieves workout logs within a date range
func (r *WorkoutLogRepository) GetWorkoutLogsByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.WorkoutLogWithPlan, error) {
	query := `
		SELECT wl.id, wl.user_id, wl.workout_plan_id, wl.workout_date, wl.duration_minutes,
			   wl.completed_exercises, wl.notes, wl.created_at, wl.updated_at,
			   wp.name as workout_plan_name
		FROM workout_logs wl
		LEFT JOIN workout_plans wp ON wl.workout_plan_id = wp.id
		WHERE wl.user_id = $1 AND wl.workout_date >= $2 AND wl.workout_date <= $3
		ORDER BY wl.workout_date DESC, wl.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query workout logs by date range: %w", err)
	}
	defer rows.Close()

	logs := []*models.WorkoutLogWithPlan{}
	for rows.Next() {
		log := &models.WorkoutLogWithPlan{}
		var completedExercisesJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.WorkoutPlanID,
			&log.WorkoutDate,
			&log.DurationMinutes,
			&completedExercisesJSON,
			&log.Notes,
			&log.CreatedAt,
			&log.UpdatedAt,
			&log.WorkoutPlanName,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workout log by date range: %w", err)
		}

		// Scan completed exercises JSON
		completedExercises := models.CompletedExercises{}
		if err := completedExercises.Scan(completedExercisesJSON); err != nil {
			return nil, fmt.Errorf("failed to scan completed exercises: %w", err)
		}
		log.CompletedExercises = completedExercises

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workout log by date range rows: %w", err)
	}

	return logs, nil
}
