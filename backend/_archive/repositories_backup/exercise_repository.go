package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nutrition-platform/models"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

// CreateExercise creates a new exercise
func (r *ExerciseRepository) CreateExercise(ctx context.Context, exercise *models.Exercise) error {
	query := `
		INSERT INTO exercises (
			name, description, muscle_groups, equipment, difficulty, 
			instructions, tips, created_by, is_public, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		exercise.Name,
		exercise.Description,
		models.MuscleGroups(exercise.MuscleGroups).Value(),
		models.Equipment(exercise.Equipment).Value(),
		exercise.Difficulty,
		exercise.Instructions,
		exercise.Tips,
		exercise.CreatedBy,
		exercise.IsPublic,
		now,
		now,
	).Scan(&exercise.ID, &exercise.CreatedAt, &exercise.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create exercise: %w", err)
	}

	return nil
}

// GetExerciseByID retrieves an exercise by ID
func (r *ExerciseRepository) GetExerciseByID(ctx context.Context, id int64) (*models.Exercise, error) {
	query := `
		SELECT id, name, description, muscle_groups, equipment, difficulty,
			   instructions, tips, created_by, is_public, created_at, updated_at
		FROM exercises 
		WHERE id = $1`

	exercise := &models.Exercise{}
	var muscleGroupsJSON, equipmentJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&muscleGroupsJSON,
		&equipmentJSON,
		&exercise.Difficulty,
		&exercise.Instructions,
		&exercise.Tips,
		&exercise.CreatedBy,
		&exercise.IsPublic,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	// Scan JSON fields
	mg := models.MuscleGroups{}
	if err := mg.Scan(muscleGroupsJSON); err != nil {
		return nil, fmt.Errorf("failed to scan muscle groups: %w", err)
	}
	exercise.MuscleGroups = mg

	eq := models.Equipment{}
	if err := eq.Scan(equipmentJSON); err != nil {
		return nil, fmt.Errorf("failed to scan equipment: %w", err)
	}
	exercise.Equipment = eq

	return exercise, nil
}

// UpdateExercise updates an existing exercise
func (r *ExerciseRepository) UpdateExercise(ctx context.Context, exercise *models.Exercise) error {
	query := `
		UPDATE exercises 
		SET name = $1, description = $2, muscle_groups = $3, equipment = $4,
			difficulty = $5, instructions = $6, tips = $7, is_public = $8, updated_at = $9
		WHERE id = $10
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		exercise.Name,
		exercise.Description,
		models.MuscleGroups(exercise.MuscleGroups).Value(),
		models.Equipment(exercise.Equipment).Value(),
		exercise.Difficulty,
		exercise.Instructions,
		exercise.Tips,
		exercise.IsPublic,
		now,
		exercise.ID,
	).Scan(&exercise.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("exercise not found")
		}
		return fmt.Errorf("failed to update exercise: %w", err)
	}

	return nil
}

// DeleteExercise soft deletes an exercise
func (r *ExerciseRepository) DeleteExercise(ctx context.Context, id int64) error {
	query := `UPDATE exercises SET deleted_at = $1 WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("exercise not found")
	}

	return nil
}

// ListExercises retrieves exercises with filtering and pagination
func (r *ExerciseRepository) ListExercises(ctx context.Context, req *models.ListExercisesRequest) (*models.ExerciseListResponse, error) {
	whereConditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if req.MuscleGroup != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("muscle_groups ? $%d", argIndex))
		args = append(args, req.MuscleGroup)
		argIndex++
	}

	if req.Equipment != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("equipment ? $%d", argIndex))
		args = append(args, req.Equipment)
		argIndex++
	}

	if req.Difficulty != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("difficulty = $%d", argIndex))
		args = append(args, req.Difficulty)
		argIndex++
	}

	if req.CreatedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("created_by = $%d", argIndex))
		args = append(args, *req.CreatedBy)
		argIndex++
	}

	if req.IsPublic != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_public = $%d", argIndex))
		args = append(args, *req.IsPublic)
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM exercises %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count exercises: %w", err)
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, name, description, muscle_groups, equipment, difficulty,
			   instructions, tips, created_by, is_public, created_at, updated_at
		FROM exercises 
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
		return nil, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer rows.Close()

	exercises := []*models.Exercise{}
	for rows.Next() {
		exercise := &models.Exercise{}
		var muscleGroupsJSON, equipmentJSON []byte

		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&muscleGroupsJSON,
			&equipmentJSON,
			&exercise.Difficulty,
			&exercise.Instructions,
			&exercise.Tips,
			&exercise.CreatedBy,
			&exercise.IsPublic,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise: %w", err)
		}

		// Scan JSON fields
		mg := models.MuscleGroups{}
		if err := mg.Scan(muscleGroupsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan muscle groups: %w", err)
		}
		exercise.MuscleGroups = mg

		eq := models.Equipment{}
		if err := eq.Scan(equipmentJSON); err != nil {
			return nil, fmt.Errorf("failed to scan equipment: %w", err)
		}
		exercise.Equipment = eq

		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercise rows: %w", err)
	}

	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)

	return &models.ExerciseListResponse{
		Exercises: exercises,
		Total:     total,
		Page:      int(req.Offset/req.Limit) + 1,
		Limit:     req.Limit,
		TotalPages: int(totalPages),
	}, nil
}

// GetFavoriteExercises retrieves user's favorite exercises
func (r *ExerciseRepository) GetFavoriteExercises(ctx context.Context, userID int64, limit, offset int) (*models.ExerciseListResponse, error) {
	query := `
		SELECT e.id, e.name, e.description, e.muscle_groups, e.equipment, e.difficulty,
			   e.instructions, e.tips, e.created_by, e.is_public, e.created_at, e.updated_at
		FROM exercises e
		INNER JOIN favorites f ON e.id = f.item_id
		WHERE f.user_id = $1 AND f.item_type = 'exercise' AND e.deleted_at IS NULL
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorite exercises: %w", err)
	}
	defer rows.Close()

	exercises := []*models.Exercise{}
	for rows.Next() {
		exercise := &models.Exercise{}
		var muscleGroupsJSON, equipmentJSON []byte

		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&muscleGroupsJSON,
			&equipmentJSON,
			&exercise.Difficulty,
			&exercise.Instructions,
			&exercise.Tips,
			&exercise.CreatedBy,
			&exercise.IsPublic,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan favorite exercise: %w", err)
		}

		// Scan JSON fields
		mg := models.MuscleGroups{}
		if err := mg.Scan(muscleGroupsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan muscle groups: %w", err)
		}
		exercise.MuscleGroups = mg

		eq := models.Equipment{}
		if err := eq.Scan(equipmentJSON); err != nil {
			return nil, fmt.Errorf("failed to scan equipment: %w", err)
		}
		exercise.Equipment = eq

		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating favorite exercise rows: %w", err)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM exercises e
		INNER JOIN favorites f ON e.id = f.item_id
		WHERE f.user_id = $1 AND f.item_type = 'exercise' AND e.deleted_at IS NULL`
	
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count favorite exercises: %w", err)
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &models.ExerciseListResponse{
		Exercises: exercises,
		Total:     total,
		Page:      int(offset/limit) + 1,
		Limit:     limit,
		TotalPages: int(totalPages),
	}, nil
}

// AddToFavorites adds an exercise to user's favorites
func (r *ExerciseRepository) AddToFavorites(ctx context.Context, userID, exerciseID int64) error {
	query := `
		INSERT INTO favorites (user_id, item_id, item_type, created_at)
		VALUES ($1, $2, 'exercise', $3)
		ON CONFLICT (user_id, item_id, item_type) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, exerciseID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add exercise to favorites: %w", err)
	}

	return nil
}

// RemoveFromFavorites removes an exercise from user's favorites
func (r *ExerciseRepository) RemoveFromFavorites(ctx context.Context, userID, exerciseID int64) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND item_id = $2 AND item_type = 'exercise'`
	
	result, err := r.db.ExecContext(ctx, query, userID, exerciseID)
	if err != nil {
		return fmt.Errorf("failed to remove exercise from favorites: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("exercise not found in favorites")
	}

	return nil
}

// GetMuscleGroups returns all available muscle groups
func (r *ExerciseRepository) GetMuscleGroups(ctx context.Context) ([]string, error) {
	// Return static list since we defined these as constants
	return models.AllMuscleGroups, nil
}

// GetEquipmentTypes returns all available equipment types
func (r *ExerciseRepository) GetEquipmentTypes(ctx context.Context) ([]string, error) {
	// Return static list since we defined these as constants
	return models.AllEquipment, nil
}
