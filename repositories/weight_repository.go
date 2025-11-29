package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/database"
	"nutrition-platform/models"
)

type WeightRepository struct {
	db *database.Database
}

func NewWeightRepository(db *database.Database) *WeightRepository {
	return &WeightRepository{db: db}
}

// CreateWeightLog creates a new weight log entry
func (r *WeightRepository) CreateWeightLog(log *models.WeightLog) error {
	query := `
		INSERT INTO weight_logs (user_id, weight, date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, log.UserID, log.Weight, log.Date, log.Notes).
		Scan(&log.ID, &log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create weight log: %w", err)
	}

	return nil
}

// GetWeightLogs retrieves weight logs for a user
func (r *WeightRepository) GetWeightLogs(userID uint, limit, offset int, startDate, endDate *time.Time) ([]*models.WeightLog, error) {
	query := `
		SELECT id, user_id, weight, date, notes, created_at, updated_at
		FROM weight_logs
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIndex := 2

	if startDate != nil {
		query += fmt.Sprintf(" AND date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY date DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.WeightLog
	for rows.Next() {
		log := &models.WeightLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Weight, &log.Date,
			&log.Notes, &log.CreatedAt, &log.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weight log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetWeightLogByID retrieves a specific weight log
func (r *WeightRepository) GetWeightLogByID(id, userID uint) (*models.WeightLog, error) {
	query := `
		SELECT id, user_id, weight, date, notes, created_at, updated_at
		FROM weight_logs
		WHERE id = $1 AND user_id = $2
	`

	log := &models.WeightLog{}
	err := r.db.QueryRow(query, id, userID).Scan(
		&log.ID, &log.UserID, &log.Weight, &log.Date,
		&log.Notes, &log.CreatedAt, &log.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("weight log not found")
		}
		return nil, fmt.Errorf("failed to get weight log: %w", err)
	}

	return log, nil
}

// UpdateWeightLog updates a weight log entry
func (r *WeightRepository) UpdateWeightLog(log *models.WeightLog) error {
	query := `
		UPDATE weight_logs
		SET weight = $1, date = $2, notes = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, log.Weight, log.Date, log.Notes, log.ID, log.UserID).
		Scan(&log.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("weight log not found")
		}
		return fmt.Errorf("failed to update weight log: %w", err)
	}

	return nil
}

// DeleteWeightLog deletes a weight log entry
func (r *WeightRepository) DeleteWeightLog(id, userID uint) error {
	query := `DELETE FROM weight_logs WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete weight log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("weight log not found")
	}

	return nil
}
