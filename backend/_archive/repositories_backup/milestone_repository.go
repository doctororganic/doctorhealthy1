package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

type MilestoneRepository struct {
	db *sql.DB
}

func NewMilestoneRepository(db *sql.DB) *MilestoneRepository {
	return &MilestoneRepository{db: db}
}

// CreateMilestone creates a new milestone
func (r *MilestoneRepository) CreateMilestone(ctx context.Context, milestone *models.Milestone) error {
	query := `
		INSERT INTO milestones (
			user_id, title, description, milestone_type, target_value,
			current_value, target_date, is_achieved, achieved_date,
			category, priority, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		milestone.UserID,
		milestone.Title,
		milestone.Description,
		milestone.MilestoneType,
		milestone.TargetValue,
		milestone.CurrentValue,
		milestone.TargetDate,
		milestone.IsAchieved,
		milestone.AchievedDate,
		milestone.Category,
		milestone.Priority,
		milestone.IsActive,
		now,
		now,
	).Scan(&milestone.ID, &milestone.CreatedAt, &milestone.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	return nil
}

// GetMilestoneByID retrieves a milestone by ID
func (r *MilestoneRepository) GetMilestoneByID(ctx context.Context, id int64) (*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE id = $1`

	var milestone models.Milestone
	var achievedDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&milestone.ID,
		&milestone.UserID,
		&milestone.Title,
		&milestone.Description,
		&milestone.MilestoneType,
		&milestone.TargetValue,
		&milestone.CurrentValue,
		&milestone.TargetDate,
		&milestone.IsAchieved,
		&achievedDate,
		&milestone.Category,
		&milestone.Priority,
		&milestone.IsActive,
		&milestone.CreatedAt,
		&milestone.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("milestone not found")
		}
		return nil, fmt.Errorf("failed to get milestone: %w", err)
	}

	if achievedDate.Valid {
		milestone.AchievedDate = &achievedDate.Time
	}

	return &milestone, nil
}

// GetMilestonesByUserID retrieves milestones for a user with pagination
func (r *MilestoneRepository) GetMilestonesByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// UpdateMilestone updates an existing milestone
func (r *MilestoneRepository) UpdateMilestone(ctx context.Context, milestone *models.Milestone) error {
	query := `
		UPDATE milestones
		SET title = $1, description = $2, milestone_type = $3, target_value = $4,
			current_value = $5, target_date = $6, is_achieved = $7, achieved_date = $8,
			category = $9, priority = $10, is_active = $11, updated_at = $12
		WHERE id = $13 AND user_id = $14
		RETURNING updated_at`

	now := time.Now()
	var achievedDate interface{} = nil
	if milestone.AchievedDate != nil {
		achievedDate = *milestone.AchievedDate
	}

	err := r.db.QueryRowContext(ctx, query,
		milestone.Title,
		milestone.Description,
		milestone.MilestoneType,
		milestone.TargetValue,
		milestone.CurrentValue,
		milestone.TargetDate,
		milestone.IsAchieved,
		achievedDate,
		milestone.Category,
		milestone.Priority,
		milestone.IsActive,
		now,
		milestone.ID,
		milestone.UserID,
	).Scan(&milestone.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("milestone not found or access denied")
		}
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	return nil
}

// DeleteMilestone deletes a milestone
func (r *MilestoneRepository) DeleteMilestone(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM milestones WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone not found or access denied")
	}

	return nil
}

// GetActiveMilestones retrieves active milestones for a user
func (r *MilestoneRepository) GetActiveMilestones(ctx context.Context, userID int64, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1 AND is_active = true AND is_achieved = false
		ORDER BY priority DESC, target_date ASC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get active milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// GetAchievedMilestones retrieves achieved milestones for a user
func (r *MilestoneRepository) GetAchievedMilestones(ctx context.Context, userID int64, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1 AND is_achieved = true
		ORDER BY achieved_date DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get achieved milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// GetMilestonesByType retrieves milestones filtered by type
func (r *MilestoneRepository) GetMilestonesByType(ctx context.Context, userID int64, milestoneType string, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1 AND milestone_type = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, userID, milestoneType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones by type: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// GetMilestonesByCategory retrieves milestones filtered by category
func (r *MilestoneRepository) GetMilestonesByCategory(ctx context.Context, userID int64, category string, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1 AND category = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, userID, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones by category: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// GetUpcomingMilestones retrieves milestones approaching their target date
func (r *MilestoneRepository) GetUpcomingMilestones(ctx context.Context, userID int64, days int, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1 
			AND is_active = true 
			AND is_achieved = false
			AND target_date BETWEEN NOW() AND NOW() + INTERVAL '%d days'
		ORDER BY target_date ASC, priority DESC
		LIMIT $2 OFFSET $3`

	query = fmt.Sprintf(query, days)
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// GetMilestoneProgress calculates progress for active milestones
func (r *MilestoneRepository) GetMilestoneProgress(ctx context.Context, userID int64) ([]*models.MilestoneProgress, error) {
	query := `
		SELECT 
			id,
			title,
			milestone_type,
			target_value,
			current_value,
			CASE 
				WHEN target_value > 0 THEN (current_value / target_value) * 100
				ELSE 0
			END as progress_percentage,
			target_date,
			CASE 
				WHEN target_date > NOW() THEN EXTRACT(DAYS FROM target_date - NOW())
				ELSE 0
			END as days_remaining
		FROM milestones
		WHERE user_id = $1 AND is_active = true AND is_achieved = false
		ORDER BY priority DESC, target_date ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone progress: %w", err)
	}
	defer rows.Close()

	var progressList []*models.MilestoneProgress
	for rows.Next() {
		var progress models.MilestoneProgress
		err := rows.Scan(
			&progress.ID,
			&progress.Title,
			&progress.MilestoneType,
			&progress.TargetValue,
			&progress.CurrentValue,
			&progress.ProgressPercentage,
			&progress.TargetDate,
			&progress.DaysRemaining,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone progress: %w", err)
		}

		progressList = append(progressList, &progress)
	}

	return progressList, nil
}

// AchieveMilestone marks a milestone as achieved
func (r *MilestoneRepository) AchieveMilestone(ctx context.Context, id, userID int64) error {
	query := `
		UPDATE milestones
		SET is_achieved = true, achieved_date = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4 AND is_active = true
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, now, now, id, userID).Scan(new(time.Time))
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("milestone not found, access denied, or already achieved")
		}
		return fmt.Errorf("failed to achieve milestone: %w", err)
	}

	return nil
}

// UpdateMilestoneProgress updates the current value of a milestone
func (r *MilestoneRepository) UpdateMilestoneProgress(ctx context.Context, id, userID int64, currentValue float64) error {
	query := `
		UPDATE milestones
		SET current_value = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4 AND is_active = true
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, currentValue, now, id, userID).Scan(new(time.Time))
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("milestone not found or access denied")
		}
		return fmt.Errorf("failed to update milestone progress: %w", err)
	}

	return nil
}

// GetMilestoneStats calculates statistics for a user's milestones
func (r *MilestoneRepository) GetMilestoneStats(ctx context.Context, userID int64) (*models.MilestoneStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_milestones,
			COUNT(CASE WHEN is_achieved = true THEN 1 END) as achieved_milestones,
			COUNT(CASE WHEN is_active = true AND is_achieved = false THEN 1 END) as active_milestones,
			COUNT(CASE WHEN target_date <= NOW() AND is_achieved = false THEN 1 END) as overdue_milestones,
			COUNT(CASE WHEN target_date BETWEEN NOW() AND NOW() + INTERVAL '7 days' 
				AND is_achieved = false AND is_active = true THEN 1 END) as upcoming_milestones
		FROM milestones
		WHERE user_id = $1`

	var stats models.MilestoneStats
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.TotalMilestones,
		&stats.AchievedMilestones,
		&stats.ActiveMilestones,
		&stats.OverdueMilestones,
		&stats.UpcomingMilestones,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get milestone stats: %w", err)
	}

	// Calculate achievement rate
	if stats.TotalMilestones > 0 {
		stats.AchievementRate = float64(stats.AchievedMilestones) / float64(stats.TotalMilestones) * 100
	}

	return &stats, nil
}

// SearchMilestones searches milestones by title or description
func (r *MilestoneRepository) SearchMilestones(ctx context.Context, userID int64, searchTerm string, limit, offset int) ([]*models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, target_value,
			   current_value, target_date, is_achieved, achieved_date,
			   category, priority, is_active, created_at, updated_at
		FROM milestones
		WHERE user_id = $1 AND (title ILIKE $2 OR description ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, userID, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search milestones: %w", err)
	}
	defer rows.Close()

	var milestones []*models.Milestone
	for rows.Next() {
		var milestone models.Milestone
		var achievedDate sql.NullTime

		err := rows.Scan(
			&milestone.ID,
			&milestone.UserID,
			&milestone.Title,
			&milestone.Description,
			&milestone.MilestoneType,
			&milestone.TargetValue,
			&milestone.CurrentValue,
			&milestone.TargetDate,
			&milestone.IsAchieved,
			&achievedDate,
			&milestone.Category,
			&milestone.Priority,
			&milestone.IsActive,
			&milestone.CreatedAt,
			&milestone.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}

		if achievedDate.Valid {
			milestone.AchievedDate = &achievedDate.Time
		}

		milestones = append(milestones, &milestone)
	}

	return milestones, nil
}

// GetMilestoneCountByUserID gets the total count of milestones for a user
func (r *MilestoneRepository) GetMilestoneCountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM milestones WHERE user_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get milestone count: %w", err)
	}

	return count, nil
}
