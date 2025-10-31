package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

type WeightGoalRepository struct {
	db *sql.DB
}

func NewWeightGoalRepository(db *sql.DB) *WeightGoalRepository {
	return &WeightGoalRepository{db: db}
}

// CreateWeightGoal creates a new weight goal
func (r *WeightGoalRepository) CreateWeightGoal(ctx context.Context, goal *models.WeightGoal) error {
	query := `
		INSERT INTO weight_goals (
			user_id, start_weight, target_weight, current_weight,
			goal_type, target_date, weekly_target, activity_level,
			is_active, achieved_date, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	var achievedDate interface{} = nil
	if goal.AchievedDate != nil {
		achievedDate = *goal.AchievedDate
	}

	err := r.db.QueryRowContext(ctx, query,
		goal.UserID,
		goal.StartWeight,
		goal.TargetWeight,
		goal.CurrentWeight,
		goal.GoalType,
		goal.TargetDate,
		goal.WeeklyTarget,
		goal.ActivityLevel,
		goal.IsActive,
		achievedDate,
		goal.Notes,
		now,
		now,
	).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create weight goal: %w", err)
	}

	return nil
}

// GetWeightGoalByID retrieves a weight goal by ID
func (r *WeightGoalRepository) GetWeightGoalByID(ctx context.Context, id int64) (*models.WeightGoal, error) {
	query := `
		SELECT id, user_id, start_weight, target_weight, current_weight,
			   goal_type, target_date, weekly_target, activity_level,
			   is_active, achieved_date, notes, created_at, updated_at
		FROM weight_goals
		WHERE id = $1`

	var goal models.WeightGoal
	var achievedDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.StartWeight,
		&goal.TargetWeight,
		&goal.CurrentWeight,
		&goal.GoalType,
		&goal.TargetDate,
		&goal.WeeklyTarget,
		&goal.ActivityLevel,
		&goal.IsActive,
		&achievedDate,
		&goal.Notes,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("weight goal not found")
		}
		return nil, fmt.Errorf("failed to get weight goal: %w", err)
	}

	if achievedDate.Valid {
		goal.AchievedDate = &achievedDate.Time
	}

	return &goal, nil
}

// GetWeightGoalsByUserID retrieves weight goals for a user with pagination
func (r *WeightGoalRepository) GetWeightGoalsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.WeightGoal, error) {
	query := `
		SELECT id, user_id, start_weight, target_weight, current_weight,
			   goal_type, target_date, weekly_target, activity_level,
			   is_active, achieved_date, notes, created_at, updated_at
		FROM weight_goals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight goals: %w", err)
	}
	defer rows.Close()

	var goals []*models.WeightGoal
	for rows.Next() {
		var goal models.WeightGoal
		var achievedDate sql.NullTime

		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.StartWeight,
			&goal.TargetWeight,
			&goal.CurrentWeight,
			&goal.GoalType,
			&goal.TargetDate,
			&goal.WeeklyTarget,
			&goal.ActivityLevel,
			&goal.IsActive,
			&achievedDate,
			&goal.Notes,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan weight goal: %w", err)
		}

		if achievedDate.Valid {
			goal.AchievedDate = &achievedDate.Time
		}

		goals = append(goals, &goal)
	}

	return goals, nil
}

// UpdateWeightGoal updates an existing weight goal
func (r *WeightGoalRepository) UpdateWeightGoal(ctx context.Context, goal *models.WeightGoal) error {
	query := `
		UPDATE weight_goals
		SET start_weight = $1, target_weight = $2, current_weight = $3,
			goal_type = $4, target_date = $5, weekly_target = $6, 
			activity_level = $7, is_active = $8, achieved_date = $9,
			notes = $10, updated_at = $11
		WHERE id = $12 AND user_id = $13
		RETURNING updated_at`

	now := time.Now()
	var achievedDate interface{} = nil
	if goal.AchievedDate != nil {
		achievedDate = *goal.AchievedDate
	}

	err := r.db.QueryRowContext(ctx, query,
		goal.StartWeight,
		goal.TargetWeight,
		goal.CurrentWeight,
		goal.GoalType,
		goal.TargetDate,
		goal.WeeklyTarget,
		goal.ActivityLevel,
		goal.IsActive,
		achievedDate,
		goal.Notes,
		now,
		goal.ID,
		goal.UserID,
	).Scan(&goal.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("weight goal not found or access denied")
		}
		return fmt.Errorf("failed to update weight goal: %w", err)
	}

	return nil
}

// DeleteWeightGoal deletes a weight goal
func (r *WeightGoalRepository) DeleteWeightGoal(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM weight_goals WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete weight goal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("weight goal not found or access denied")
	}

	return nil
}

// GetActiveWeightGoal gets the active weight goal for a user
func (r *WeightGoalRepository) GetActiveWeightGoal(ctx context.Context, userID int64) (*models.WeightGoal, error) {
	query := `
		SELECT id, user_id, start_weight, target_weight, current_weight,
			   goal_type, target_date, weekly_target, activity_level,
			   is_active, achieved_date, notes, created_at, updated_at
		FROM weight_goals
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT 1`

	var goal models.WeightGoal
	var achievedDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.StartWeight,
		&goal.TargetWeight,
		&goal.CurrentWeight,
		&goal.GoalType,
		&goal.TargetDate,
		&goal.WeeklyTarget,
		&goal.ActivityLevel,
		&goal.IsActive,
		&achievedDate,
		&goal.Notes,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active weight goal found")
		}
		return nil, fmt.Errorf("failed to get active weight goal: %w", err)
	}

	if achievedDate.Valid {
		goal.AchievedDate = &achievedDate.Time
	}

	return &goal, nil
}

// UpdateCurrentWeight updates the current weight for a goal
func (r *WeightGoalRepository) UpdateCurrentWeight(ctx context.Context, id, userID int64, currentWeight float64) error {
	query := `
		UPDATE weight_goals
		SET current_weight = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4 AND is_active = true
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, currentWeight, now, id, userID).Scan(new(time.Time))
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("weight goal not found, access denied, or not active")
		}
		return fmt.Errorf("failed to update current weight: %w", err)
	}

	return nil
}

// AchieveWeightGoal marks a weight goal as achieved
func (r *WeightGoalRepository) AchieveWeightGoal(ctx context.Context, id, userID int64) error {
	query := `
		UPDATE weight_goals
		SET is_active = false, achieved_date = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4 AND is_active = true
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, now, now, id, userID).Scan(new(time.Time))
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("weight goal not found, access denied, or already achieved")
		}
		return fmt.Errorf("failed to achieve weight goal: %w", err)
	}

	return nil
}

// GetWeightGoalProgress calculates progress for active weight goals
func (r *WeightGoalRepository) GetWeightGoalProgress(ctx context.Context, userID int64) (*models.WeightGoalProgress, error) {
	query := `
		SELECT 
			id,
			start_weight,
			target_weight,
			current_weight,
			goal_type,
			weekly_target,
			target_date,
			CASE 
				WHEN goal_type = 'lose' THEN 
					CASE WHEN start_weight > target_weight THEN 
						((start_weight - current_weight) / (start_weight - target_weight)) * 100
					ELSE 0
					END
				WHEN goal_type = 'gain' THEN 
					CASE WHEN target_weight > start_weight THEN 
						((current_weight - start_weight) / (target_weight - start_weight)) * 100
					ELSE 0
					END
				ELSE 0
			END as progress_percentage,
			CASE 
				WHEN target_date > NOW() THEN EXTRACT(DAYS FROM target_date - NOW())
				ELSE 0
			END as days_remaining,
			EXTRACT(DAYS FROM NOW() - created_at) as days_elapsed
		FROM weight_goals
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT 1`

	var progress models.WeightGoalProgress
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&progress.ID,
		&progress.StartWeight,
		&progress.TargetWeight,
		&progress.CurrentWeight,
		&progress.GoalType,
		&progress.WeeklyTarget,
		&progress.TargetDate,
		&progress.ProgressPercentage,
		&progress.DaysRemaining,
		&progress.DaysElapsed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active weight goal found")
		}
		return nil, fmt.Errorf("failed to get weight goal progress: %w", err)
	}

	// Calculate weight remaining and weekly progress
	progress.WeightRemaining = progress.TargetWeight - progress.CurrentWeight
	if progress.DaysElapsed > 0 {
		weeklyChange := (progress.CurrentWeight - progress.StartWeight) / (progress.DaysElapsed / 7.0)
		progress.WeeklyProgress = &weeklyChange
	}

	return &progress, nil
}

// GetWeightGoalHistory retrieves historical weight goal data
func (r *WeightGoalRepository) GetWeightGoalHistory(ctx context.Context, userID int64, limit, offset int) ([]*models.WeightGoal, error) {
	query := `
		SELECT id, user_id, start_weight, target_weight, current_weight,
			   goal_type, target_date, weekly_target, activity_level,
			   is_active, achieved_date, notes, created_at, updated_at
		FROM weight_goals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	return r.GetWeightGoalsByUserID(ctx, userID, limit, offset)
}

// GetWeightGoalStats calculates statistics for a user's weight goals
func (r *WeightGoalRepository) GetWeightGoalStats(ctx context.Context, userID int64) (*models.WeightGoalStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_goals,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_goals,
			COUNT(CASE WHEN is_active = false AND achieved_date IS NOT NULL THEN 1 END) as achieved_goals,
			COUNT(CASE WHEN is_active = false AND achieved_date IS NULL THEN 1 END) as abandoned_goals,
			COUNT(CASE WHEN goal_type = 'lose' THEN 1 END) as weight_loss_goals,
			COUNT(CASE WHEN goal_type = 'gain' THEN 1 END) as weight_gain_goals,
			COUNT(CASE WHEN goal_type = 'maintain' THEN 1 END) as maintenance_goals
		FROM weight_goals
		WHERE user_id = $1`

	var stats models.WeightGoalStats
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.TotalGoals,
		&stats.ActiveGoals,
		&stats.AchievedGoals,
		&stats.AbandonedGoals,
		&stats.WeightLossGoals,
		&stats.WeightGainGoals,
		&stats.MaintenanceGoals,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get weight goal stats: %w", err)
	}

	// Calculate success rate
	if stats.TotalGoals > 0 {
		stats.SuccessRate = float64(stats.AchievedGoals) / float64(stats.TotalGoals) * 100
	}

	return &stats, nil
}

// GetWeightGoalsByType retrieves weight goals filtered by type
func (r *WeightGoalRepository) GetWeightGoalsByType(ctx context.Context, userID int64, goalType string, limit, offset int) ([]*models.WeightGoal, error) {
	query := `
		SELECT id, user_id, start_weight, target_weight, current_weight,
			   goal_type, target_date, weekly_target, activity_level,
			   is_active, achieved_date, notes, created_at, updated_at
		FROM weight_goals
		WHERE user_id = $1 AND goal_type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, userID, goalType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight goals by type: %w", err)
	}
	defer rows.Close()

	var goals []*models.WeightGoal
	for rows.Next() {
		var goal models.WeightGoal
		var achievedDate sql.NullTime

		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.StartWeight,
			&goal.TargetWeight,
			&goal.CurrentWeight,
			&goal.GoalType,
			&goal.TargetDate,
			&goal.WeeklyTarget,
			&goal.ActivityLevel,
			&goal.IsActive,
			&achievedDate,
			&goal.Notes,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan weight goal: %w", err)
		}

		if achievedDate.Valid {
			goal.AchievedDate = &achievedDate.Time
		}

		goals = append(goals, &goal)
	}

	return goals, nil
}

// GetWeightGoalPredictions predicts goal completion based on current progress
func (r *WeightGoalRepository) GetWeightGoalPredictions(ctx context.Context, userID int64) (*models.WeightGoalPrediction, error) {
	query := `
		SELECT 
			wg.id,
			wg.start_weight,
			wg.target_weight,
			wg.current_weight,
			wg.goal_type,
			wg.weekly_target,
			wg.target_date,
			EXTRACT(DAYS FROM NOW() - wg.created_at) as days_elapsed,
			-- Get recent weight changes from body measurements
			(
				SELECT AVG(weight - LAG(weight, 1, weight) OVER (ORDER BY measurement_date))
				FROM body_measurements
				WHERE user_id = $1 
					AND measurement_date >= NOW() - INTERVAL '30 days'
					AND measurement_date <= NOW()
				ORDER BY measurement_date DESC
				LIMIT 5
			) as recent_weekly_change
		FROM weight_goals wg
		WHERE wg.user_id = $1 AND wg.is_active = true
		ORDER BY wg.created_at DESC
		LIMIT 1`

	var prediction models.WeightGoalPrediction
	var recentWeeklyChange sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&prediction.GoalID,
		&prediction.StartWeight,
		&prediction.TargetWeight,
		&prediction.CurrentWeight,
		&prediction.GoalType,
		&prediction.WeeklyTarget,
		&prediction.TargetDate,
		&prediction.DaysElapsed,
		&recentWeeklyChange,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active weight goal found")
		}
		return nil, fmt.Errorf("failed to get weight goal predictions: %w", err)
	}

	// Calculate predictions
	if recentWeeklyChange.Valid && recentWeeklyChange.Float64 != 0 {
		weightDiff := prediction.TargetWeight - prediction.CurrentWeight
		if recentWeeklyChange.Float64 != 0 {
			weeksToGoal := weightDiff / recentWeeklyChange.Float64
			if weeksToGoal > 0 {
				prediction.PredictedCompletionDate = time.Now().AddDate(0, 0, int(weeksToGoal*7))
				prediction.IsOnTrack = weeksToGoal <= (prediction.TargetDate.Sub(time.Now()).Hours() / 24 / 7)
			}
		}
	}

	return &prediction, nil
}

// SearchWeightGoals searches weight goals by notes or goal type
func (r *WeightGoalRepository) SearchWeightGoals(ctx context.Context, userID int64, searchTerm string, limit, offset int) ([]*models.WeightGoal, error) {
	query := `
		SELECT id, user_id, start_weight, target_weight, current_weight,
			   goal_type, target_date, weekly_target, activity_level,
			   is_active, achieved_date, notes, created_at, updated_at
		FROM weight_goals
		WHERE user_id = $1 AND (notes ILIKE $2 OR goal_type ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, userID, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search weight goals: %w", err)
	}
	defer rows.Close()

	var goals []*models.WeightGoal
	for rows.Next() {
		var goal models.WeightGoal
		var achievedDate sql.NullTime

		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.StartWeight,
			&goal.TargetWeight,
			&goal.CurrentWeight,
			&goal.GoalType,
			&goal.TargetDate,
			&goal.WeeklyTarget,
			&goal.ActivityLevel,
			&goal.IsActive,
			&achievedDate,
			&goal.Notes,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan weight goal: %w", err)
		}

		if achievedDate.Valid {
			goal.AchievedDate = &achievedDate.Time
		}

		goals = append(goals, &goal)
	}

	return goals, nil
}

// GetWeightGoalCountByUserID gets the total count of weight goals for a user
func (r *WeightGoalRepository) GetWeightGoalCountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM weight_goals WHERE user_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get weight goal count: %w", err)
	}

	return count, nil
}
