package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nutrition-platform/models"
)

type PersonalRecordRepository struct {
	db *sql.DB
}

func NewPersonalRecordRepository(db *sql.DB) *PersonalRecordRepository {
	return &PersonalRecordRepository{db: db}
}

// CreatePersonalRecord creates a new personal record
func (r *PersonalRecordRepository) CreatePersonalRecord(ctx context.Context, pr *models.PersonalRecord) error {
	// Check if this is actually a new PR
	existing, err := r.GetUserBestRecord(ctx, pr.UserID, pr.ExerciseID, pr.RecordType)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("failed to check existing record: %w", err)
	}

	// If existing record exists, only create if this is better
	if existing != nil {
		if !r.isNewRecordBetter(pr, existing, pr.RecordType) {
			return fmt.Errorf("new record (%.2f) is not better than existing record (%.2f)", pr.RecordValue, existing.RecordValue)
		}
	}

	query := `
		INSERT INTO personal_records (
			user_id, exercise_id, record_type, record_value, 
			workout_log_id, achieved_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err = r.db.QueryRowContext(ctx, query,
		pr.UserID,
		pr.ExerciseID,
		pr.RecordType,
		pr.RecordValue,
		pr.WorkoutLogID,
		pr.AchievedAt,
		time.Now(),
	).Scan(&pr.ID, &pr.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create personal record: %w", err)
	}

	return nil
}

// GetPersonalRecordByID retrieves a personal record by ID
func (r *PersonalRecordRepository) GetPersonalRecordByID(ctx context.Context, id int64) (*models.PersonalRecord, error) {
	query := `
		SELECT id, user_id, exercise_id, record_type, record_value,
			   workout_log_id, achieved_at, created_at
		FROM personal_records 
		WHERE id = $1`

	pr := &models.PersonalRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pr.ID,
		&pr.UserID,
		&pr.ExerciseID,
		&pr.RecordType,
		&pr.RecordValue,
		&pr.WorkoutLogID,
		&pr.AchievedAt,
		&pr.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("personal record not found")
		}
		return nil, fmt.Errorf("failed to get personal record: %w", err)
	}

	return pr, nil
}

// UpdatePersonalRecord updates an existing personal record
func (r *PersonalRecordRepository) UpdatePersonalRecord(ctx context.Context, pr *models.PersonalRecord) error {
	query := `
		UPDATE personal_records 
		SET record_type = $1, record_value = $2, workout_log_id = $3, achieved_at = $4
		WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query,
		pr.RecordType,
		pr.RecordValue,
		pr.WorkoutLogID,
		pr.AchievedAt,
		pr.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update personal record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("personal record not found")
	}

	return nil
}

// DeletePersonalRecord deletes a personal record
func (r *PersonalRecordRepository) DeletePersonalRecord(ctx context.Context, id int64) error {
	query := `DELETE FROM personal_records WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete personal record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("personal record not found")
	}

	return nil
}

// ListPersonalRecords retrieves personal records with filtering and pagination
func (r *PersonalRecordRepository) ListPersonalRecords(ctx context.Context, req *models.ListPersonalRecordsRequest) (*models.PersonalRecordListResponse, error) {
	whereConditions := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if req.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pr.user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.ExerciseID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pr.exercise_id = $%d", argIndex))
		args = append(args, *req.ExerciseID)
		argIndex++
	}

	if req.RecordType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("pr.record_type = $%d", argIndex))
		args = append(args, req.RecordType)
		argIndex++
	}

	if req.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pr.achieved_at >= $%d", argIndex))
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pr.achieved_at <= $%d", argIndex))
		args = append(args, *req.EndDate)
		argIndex++
	}

	whereClause := "WHERE " + strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM personal_records pr %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count personal records: %w", err)
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT pr.id, pr.user_id, pr.exercise_id, pr.record_type, pr.record_value,
			   pr.workout_log_id, pr.achieved_at, pr.created_at,
			   e.name as exercise_name, e.muscle_groups
		FROM personal_records pr
		INNER JOIN exercises e ON pr.exercise_id = e.id
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
		return nil, fmt.Errorf("failed to query personal records: %w", err)
	}
	defer rows.Close()

	records := []*models.PersonalRecordWithExercise{}
	for rows.Next() {
		record := &models.PersonalRecordWithExercise{}
		var muscleGroupsJSON []byte

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.ExerciseID,
			&record.RecordType,
			&record.RecordValue,
			&record.WorkoutLogID,
			&record.AchievedAt,
			&record.CreatedAt,
			&record.ExerciseName,
			&muscleGroupsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan personal record: %w", err)
		}

		// Scan muscle groups JSON
		mg := models.MuscleGroups{}
		if err := mg.Scan(muscleGroupsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan muscle groups: %w", err)
		}
		record.MuscleGroups = mg

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating personal record rows: %w", err)
	}

	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)

	return &models.PersonalRecordListResponse{
		Records:    records,
		Total:      total,
		Page:       int(req.Offset/req.Limit) + 1,
		Limit:      req.Limit,
		TotalPages: int(totalPages),
	}, nil
}

// GetUserBestRecord retrieves the best record for a user and exercise
func (r *PersonalRecordRepository) GetUserBestRecord(ctx context.Context, userID, exerciseID int64, recordType models.PersonalRecordType) (*models.PersonalRecord, error) {
	var orderBy string
	switch recordType {
	case models.PRTypeWeight:
		orderBy = "record_value DESC"
	case models.PRTypeTime:
		orderBy = "record_value ASC"
	case models.PRTypeReps:
		orderBy = "record_value DESC"
	default:
		orderBy = "record_value DESC"
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, exercise_id, record_type, record_value,
			   workout_log_id, achieved_at, created_at
		FROM personal_records 
		WHERE user_id = $1 AND exercise_id = $2 AND record_type = $3
		ORDER BY %s, achieved_at DESC
		LIMIT 1`, orderBy)

	pr := &models.PersonalRecord{}
	err := r.db.QueryRowContext(ctx, query, userID, exerciseID, recordType).Scan(
		&pr.ID,
		&pr.UserID,
		&pr.ExerciseID,
		&pr.RecordType,
		&pr.RecordValue,
		&pr.WorkoutLogID,
		&pr.AchievedAt,
		&pr.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("personal record not found")
		}
		return nil, fmt.Errorf("failed to get user best record: %w", err)
	}

	return pr, nil
}

// GetUserPersonalRecords retrieves all personal records for a user
func (r *PersonalRecordRepository) GetUserPersonalRecords(ctx context.Context, userID int64) ([]*models.PersonalRecordWithExercise, error) {
	query := `
		SELECT DISTINCT ON (pr.exercise_id, pr.record_type)
			pr.id, pr.user_id, pr.exercise_id, pr.record_type, pr.record_value,
			pr.workout_log_id, pr.achieved_at, pr.created_at,
			e.name as exercise_name, e.muscle_groups
		FROM personal_records pr
		INNER JOIN exercises e ON pr.exercise_id = e.id
		WHERE pr.user_id = $1
		ORDER BY pr.exercise_id, pr.record_type, 
			CASE 
				WHEN pr.record_type = 'weight' THEN pr.record_value DESC
				WHEN pr.record_type = 'time' THEN pr.record_value ASC
				ELSE pr.record_value DESC
			END,
			pr.achieved_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user personal records: %w", err)
	}
	defer rows.Close()

	records := []*models.PersonalRecordWithExercise{}
	for rows.Next() {
		record := &models.PersonalRecordWithExercise{}
		var muscleGroupsJSON []byte

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.ExerciseID,
			&record.RecordType,
			&record.RecordValue,
			&record.WorkoutLogID,
			&record.AchievedAt,
			&record.CreatedAt,
			&record.ExerciseName,
			&muscleGroupsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user personal record: %w", err)
		}

		// Scan muscle groups JSON
		mg := models.MuscleGroups{}
		if err := mg.Scan(muscleGroupsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan muscle groups: %w", err)
		}
		record.MuscleGroups = mg

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user personal record rows: %w", err)
	}

	return records, nil
}

// GetPersonalRecordHistory retrieves the history of personal records for a specific exercise
func (r *PersonalRecordRepository) GetPersonalRecordHistory(ctx context.Context, userID, exerciseID int64, recordType models.PersonalRecordType) ([]*models.PersonalRecord, error) {
	query := `
		SELECT id, user_id, exercise_id, record_type, record_value,
			   workout_log_id, achieved_at, created_at
		FROM personal_records 
		WHERE user_id = $1 AND exercise_id = $2 AND record_type = $3
		ORDER BY achieved_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, exerciseID, recordType)
	if err != nil {
		return nil, fmt.Errorf("failed to query personal record history: %w", err)
	}
	defer rows.Close()

	records := []*models.PersonalRecord{}
	for rows.Next() {
		record := &models.PersonalRecord{}

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.ExerciseID,
			&record.RecordType,
			&record.RecordValue,
			&record.WorkoutLogID,
			&record.AchievedAt,
			&record.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan personal record history: %w", err)
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating personal record history rows: %w", err)
	}

	return records, nil
}

// GetRecentPersonalRecords retrieves recent personal records for a user
func (r *PersonalRecordRepository) GetRecentPersonalRecords(ctx context.Context, userID int64, limit int) ([]*models.PersonalRecordWithExercise, error) {
	query := `
		SELECT pr.id, pr.user_id, pr.exercise_id, pr.record_type, pr.record_value,
			   pr.workout_log_id, pr.achieved_at, pr.created_at,
			   e.name as exercise_name, e.muscle_groups
		FROM personal_records pr
		INNER JOIN exercises e ON pr.exercise_id = e.id
		WHERE pr.user_id = $1
		ORDER BY pr.achieved_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent personal records: %w", err)
	}
	defer rows.Close()

	records := []*models.PersonalRecordWithExercise{}
	for rows.Next() {
		record := &models.PersonalRecordWithExercise{}
		var muscleGroupsJSON []byte

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.ExerciseID,
			&record.RecordType,
			&record.RecordValue,
			&record.WorkoutLogID,
			&record.AchievedAt,
			&record.CreatedAt,
			&record.ExerciseName,
			&muscleGroupsJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan recent personal record: %w", err)
		}

		// Scan muscle groups JSON
		mg := models.MuscleGroups{}
		if err := mg.Scan(muscleGroupsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan muscle groups: %w", err)
		}
		record.MuscleGroups = mg

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recent personal record rows: %w", err)
	}

	return records, nil
}

// isNewRecordBetter checks if the new record is better than the existing one
func (r *PersonalRecordRepository) isNewRecordBetter(new, existing *models.PersonalRecord, recordType models.PersonalRecordType) bool {
	switch recordType {
	case models.PRTypeWeight:
		return new.RecordValue > existing.RecordValue
	case models.PRTypeTime:
		// For time records, lower is better
		return new.RecordValue < existing.RecordValue
	case models.PRTypeReps:
		return new.RecordValue > existing.RecordValue
	default:
		return new.RecordValue > existing.RecordValue
	}
}

// GetPersonalRecordStats retrieves statistics about personal records
func (r *PersonalRecordRepository) GetPersonalRecordStats(ctx context.Context, userID int64) (*models.PersonalRecordStats, error) {
	stats := &models.PersonalRecordStats{}

	// Get total PRs
	query := `SELECT COUNT(*) FROM personal_records WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&stats.TotalRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to get total records: %w", err)
	}

	// Get PRs by type
	typeQuery := `
		SELECT record_type, COUNT(*) 
		FROM personal_records 
		WHERE user_id = $1 
		GROUP BY record_type`
	
	rows, err := r.db.QueryContext(ctx, typeQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get records by type: %w", err)
	}
	defer rows.Close()

	stats.ByType = make(map[models.PersonalRecordType]int64)
	for rows.Next() {
		var recordType models.PersonalRecordType
		var count int64
		err := rows.Scan(&recordType, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan record type stats: %w", err)
		}
		stats.ByType[recordType] = count
	}

	// Get recent PR count (last 30 days)
	recentQuery := `
		SELECT COUNT(*) 
		FROM personal_records 
		WHERE user_id = $1 AND achieved_at >= $2`
	
	err = r.db.QueryRowContext(ctx, recentQuery, userID, time.Now().AddDate(0, 0, -30)).Scan(&stats.RecentRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent records: %w", err)
	}

	return stats, nil
}
