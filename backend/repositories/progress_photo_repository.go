package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nutrition-platform/models"
)

type ProgressPhotoRepository struct {
	db *sql.DB
}

func NewProgressPhotoRepository(db *sql.DB) *ProgressPhotoRepository {
	return &ProgressPhotoRepository{db: db}
}

// CreateProgressPhoto creates a new progress photo
func (r *ProgressPhotoRepository) CreateProgressPhoto(ctx context.Context, photo *models.ProgressPhoto) error {
	query := `
		INSERT INTO progress_photos (
			user_id, photo_url, thumbnail_url, file_size, file_type,
			capture_date, weight, body_fat_percentage, notes, tags,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		photo.UserID,
		photo.PhotoURL,
		photo.ThumbnailURL,
		photo.FileSize,
		photo.FileType,
		photo.CaptureDate,
		photo.Weight,
		photo.BodyFatPercentage,
		photo.Notes,
		models.StringArray(photo.Tags).Value(),
		now,
		now,
	).Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create progress photo: %w", err)
	}

	return nil
}

// GetProgressPhotoByID retrieves a progress photo by ID
func (r *ProgressPhotoRepository) GetProgressPhotoByID(ctx context.Context, id int64) (*models.ProgressPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, file_size, file_type,
			   capture_date, weight, body_fat_percentage, notes, tags,
			   created_at, updated_at
		FROM progress_photos
		WHERE id = $1`

	var photo models.ProgressPhoto
	var tags sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&photo.ID,
		&photo.UserID,
		&photo.PhotoURL,
		&photo.ThumbnailURL,
		&photo.FileSize,
		&photo.FileType,
		&photo.CaptureDate,
		&photo.Weight,
		&photo.BodyFatPercentage,
		&photo.Notes,
		&tags,
		&photo.CreatedAt,
		&photo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("progress photo not found")
		}
		return nil, fmt.Errorf("failed to get progress photo: %w", err)
	}

	if tags.Valid {
		photo.Tags = models.StringArray{}.Scan(tags.String)
	}

	return &photo, nil
}

// GetProgressPhotosByUserID retrieves progress photos for a user with pagination
func (r *ProgressPhotoRepository) GetProgressPhotosByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.ProgressPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, file_size, file_type,
			   capture_date, weight, body_fat_percentage, notes, tags,
			   created_at, updated_at
		FROM progress_photos
		WHERE user_id = $1
		ORDER BY capture_date DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress photos: %w", err)
	}
	defer rows.Close()

	var photos []*models.ProgressPhoto
	for rows.Next() {
		var photo models.ProgressPhoto
		var tags sql.NullString

		err := rows.Scan(
			&photo.ID,
			&photo.UserID,
			&photo.PhotoURL,
			&photo.ThumbnailURL,
			&photo.FileSize,
			&photo.FileType,
			&photo.CaptureDate,
			&photo.Weight,
			&photo.BodyFatPercentage,
			&photo.Notes,
			&tags,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan progress photo: %w", err)
		}

		if tags.Valid {
			photo.Tags = models.StringArray{}.Scan(tags.String)
		}

		photos = append(photos, &photo)
	}

	return photos, nil
}

// UpdateProgressPhoto updates an existing progress photo
func (r *ProgressPhotoRepository) UpdateProgressPhoto(ctx context.Context, photo *models.ProgressPhoto) error {
	query := `
		UPDATE progress_photos
		SET photo_url = $1, thumbnail_url = $2, capture_date = $3,
			weight = $4, body_fat_percentage = $5, notes = $6, tags = $7,
			updated_at = $8
		WHERE id = $9 AND user_id = $10
		RETURNING updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		photo.PhotoURL,
		photo.ThumbnailURL,
		photo.CaptureDate,
		photo.Weight,
		photo.BodyFatPercentage,
		photo.Notes,
		models.StringArray(photo.Tags).Value(),
		now,
		photo.ID,
		photo.UserID,
	).Scan(&photo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("progress photo not found or access denied")
		}
		return fmt.Errorf("failed to update progress photo: %w", err)
	}

	return nil
}

// DeleteProgressPhoto deletes a progress photo
func (r *ProgressPhotoRepository) DeleteProgressPhoto(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM progress_photos WHERE id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete progress photo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("progress photo not found or access denied")
	}

	return nil
}

// GetProgressPhotosByDateRange retrieves photos within a date range
func (r *ProgressPhotoRepository) GetProgressPhotosByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time, limit, offset int) ([]*models.ProgressPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, file_size, file_type,
			   capture_date, weight, body_fat_percentage, notes, tags,
			   created_at, updated_at
		FROM progress_photos
		WHERE user_id = $1 AND capture_date BETWEEN $2 AND $3
		ORDER BY capture_date DESC, created_at DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress photos by date range: %w", err)
	}
	defer rows.Close()

	var photos []*models.ProgressPhoto
	for rows.Next() {
		var photo models.ProgressPhoto
		var tags sql.NullString

		err := rows.Scan(
			&photo.ID,
			&photo.UserID,
			&photo.PhotoURL,
			&photo.ThumbnailURL,
			&photo.FileSize,
			&photo.FileType,
			&photo.CaptureDate,
			&photo.Weight,
			&photo.BodyFatPercentage,
			&photo.Notes,
			&tags,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan progress photo: %w", err)
		}

		if tags.Valid {
			photo.Tags = models.StringArray{}.Scan(tags.String)
		}

		photos = append(photos, &photo)
	}

	return photos, nil
}

// GetProgressPhotosByTags retrieves photos filtered by tags
func (r *ProgressPhotoRepository) GetProgressPhotosByTags(ctx context.Context, userID int64, tags []string, limit, offset int) ([]*models.ProgressPhoto, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag must be provided")
	}

	placeholders := make([]string, len(tags))
	args := make([]interface{}, len(tags)+2)
	args[0] = userID

	for i, tag := range tags {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = tag
	}
	args[len(args)-1] = limit
	args[len(args)-2] = offset

	query := fmt.Sprintf(`
		SELECT id, user_id, photo_url, thumbnail_url, file_size, file_type,
			   capture_date, weight, body_fat_percentage, notes, tags,
			   created_at, updated_at
		FROM progress_photos
		WHERE user_id = $1 AND tags && ARRAY[%s]
		ORDER BY capture_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d`,
		strings.Join(placeholders, ","),
		len(tags)+2,
		len(tags)+3,
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress photos by tags: %w", err)
	}
	defer rows.Close()

	var photos []*models.ProgressPhoto
	for rows.Next() {
		var photo models.ProgressPhoto
		var tags sql.NullString

		err := rows.Scan(
			&photo.ID,
			&photo.UserID,
			&photo.PhotoURL,
			&photo.ThumbnailURL,
			&photo.FileSize,
			&photo.FileType,
			&photo.CaptureDate,
			&photo.Weight,
			&photo.BodyFatPercentage,
			&photo.Notes,
			&tags,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan progress photo: %w", err)
		}

		if tags.Valid {
			photo.Tags = models.StringArray{}.Scan(tags.String)
		}

		photos = append(photos, &photo)
	}

	return photos, nil
}

// GetPhotoCountByUserID gets the total count of photos for a user
func (r *ProgressPhotoRepository) GetPhotoCountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM progress_photos WHERE user_id = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get photo count: %w", err)
	}

	return count, nil
}

// GetLatestProgressPhoto gets the most recent progress photo for a user
func (r *ProgressPhotoRepository) GetLatestProgressPhoto(ctx context.Context, userID int64) (*models.ProgressPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, file_size, file_type,
			   capture_date, weight, body_fat_percentage, notes, tags,
			   created_at, updated_at
		FROM progress_photos
		WHERE user_id = $1
		ORDER BY capture_date DESC, created_at DESC
		LIMIT 1`

	var photo models.ProgressPhoto
	var tags sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&photo.ID,
		&photo.UserID,
		&photo.PhotoURL,
		&photo.ThumbnailURL,
		&photo.FileSize,
		&photo.FileType,
		&photo.CaptureDate,
		&photo.Weight,
		&photo.BodyFatPercentage,
		&photo.Notes,
		&tags,
		&photo.CreatedAt,
		&photo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no progress photos found")
		}
		return nil, fmt.Errorf("failed to get latest progress photo: %w", err)
	}

	if tags.Valid {
		photo.Tags = models.StringArray{}.Scan(tags.String)
	}

	return &photo, nil
}

// SearchProgressPhotos searches photos by notes text
func (r *ProgressPhotoRepository) SearchProgressPhotos(ctx context.Context, userID int64, searchTerm string, limit, offset int) ([]*models.ProgressPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, file_size, file_type,
			   capture_date, weight, body_fat_percentage, notes, tags,
			   created_at, updated_at
		FROM progress_photos
		WHERE user_id = $1 AND (notes ILIKE $2 OR tags && ARRAY[$3])
		ORDER BY capture_date DESC, created_at DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.db.QueryContext(ctx, query, userID, "%"+searchTerm+"%", searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search progress photos: %w", err)
	}
	defer rows.Close()

	var photos []*models.ProgressPhoto
	for rows.Next() {
		var photo models.ProgressPhoto
		var tags sql.NullString

		err := rows.Scan(
			&photo.ID,
			&photo.UserID,
			&photo.PhotoURL,
			&photo.ThumbnailURL,
			&photo.FileSize,
			&photo.FileType,
			&photo.CaptureDate,
			&photo.Weight,
			&photo.BodyFatPercentage,
			&photo.Notes,
			&tags,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan progress photo: %w", err)
		}

		if tags.Valid {
			photo.Tags = models.StringArray{}.Scan(tags.String)
		}

		photos = append(photos, &photo)
	}

	return photos, nil
}

// BatchCreateProgressPhotos creates multiple progress photos in a transaction
func (r *ProgressPhotoRepository) BatchCreateProgressPhotos(ctx context.Context, photos []*models.ProgressPhoto) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO progress_photos (
			user_id, photo_url, thumbnail_url, file_size, file_type,
			capture_date, weight, body_fat_percentage, notes, tags,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, photo := range photos {
		err := stmt.QueryRowContext(ctx,
			photo.UserID,
			photo.PhotoURL,
			photo.ThumbnailURL,
			photo.FileSize,
			photo.FileType,
			photo.CaptureDate,
			photo.Weight,
			photo.BodyFatPercentage,
			photo.Notes,
			models.StringArray(photo.Tags).Value(),
			now,
			now,
		).Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt)

		if err != nil {
			return fmt.Errorf("failed to create progress photo: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
