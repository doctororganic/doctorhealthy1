package services

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// HealthService handles health-related business logic
type HealthService struct {
	db *sql.DB
}

// NewHealthService creates a new health service
func NewHealthService(db *sql.DB) *HealthService {
	return &HealthService{
		db: db,
	}
}

// GetHealthProfile retrieves health profile for a user
func (s *HealthService) GetHealthProfile(userID uint) (*models.HealthProfile, error) {
	query := `
		SELECT id, user_id, date_of_birth, gender, height, weight, blood_type,
			   allergies, medical_conditions, medications, emergency_contact,
			   created_at, updated_at
		FROM health_profiles 
		WHERE user_id = $1
	`
	
	var profile models.HealthProfile
	var dob sql.NullTime
	var height, weight sql.NullFloat64
	
	err := s.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &dob, &profile.Gender,
		&height, &weight, &profile.BloodType, &profile.Allergies,
		&profile.MedicalConditions, &profile.Medications,
		&profile.EmergencyContact, &profile.CreatedAt, &profile.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("health profile not found")
		}
		return nil, err
	}
	
	if dob.Valid {
		profile.DateOfBirth = &dob.Time
	}
	if height.Valid {
		profile.Height = &height.Float64
	}
	if weight.Valid {
		profile.Weight = &weight.Float64
	}
	
	return &profile, nil
}

// CreateHealthProfile creates a new health profile
func (s *HealthService) CreateHealthProfile(profile *models.HealthProfile) error {
	query := `
		INSERT INTO health_profiles (user_id, date_of_birth, gender, height, weight,
									blood_type, allergies, medical_conditions, medications,
									emergency_contact, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		profile.UserID, profile.DateOfBirth, profile.Gender,
		profile.Height, profile.Weight, profile.BloodType,
		profile.Allergies, profile.MedicalConditions,
		profile.Medications, profile.EmergencyContact, now, now,
	).Scan(&profile.ID)
	
	if err != nil {
		return err
	}
	
	profile.CreatedAt = now
	profile.UpdatedAt = now
	
	return nil
}

// GetHealthConditions retrieves health conditions for a user
func (s *HealthService) GetHealthConditions(userID uint) ([]models.HealthCondition, error) {
	query := `
		SELECT id, user_id, name, description, severity, diagnosed_date,
			   is_active, created_at, updated_at
		FROM health_conditions 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var conditions []models.HealthCondition
	for rows.Next() {
		var condition models.HealthCondition
		var diagnosedDate sql.NullTime
		
		err := rows.Scan(
			&condition.ID, &condition.UserID, &condition.Name,
			&condition.Description, &condition.Severity,
			&diagnosedDate, &condition.IsActive,
			&condition.CreatedAt, &condition.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if diagnosedDate.Valid {
			condition.DiagnosedDate = &diagnosedDate.Time
		}
		
		conditions = append(conditions, condition)
	}
	
	return conditions, nil
}

// CreateHealthCondition creates a new health condition
func (s *HealthService) CreateHealthCondition(condition *models.HealthCondition) error {
	query := `
		INSERT INTO health_conditions (user_id, name, description, severity,
									  diagnosed_date, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		condition.UserID, condition.Name, condition.Description,
		condition.Severity, condition.DiagnosedDate, condition.IsActive,
		now, now,
	).Scan(&condition.ID)
	
	if err != nil {
		return err
	}
	
	condition.CreatedAt = now
	condition.UpdatedAt = now
	
	return nil
}

// GetProgressData retrieves progress data for a user
func (s *HealthService) GetProgressData(userID uint) ([]models.ProgressEntry, error) {
	query := `
		SELECT id, user_id, weight, body_fat_percentage, measurements,
			   notes, recorded_at, created_at, updated_at
		FROM progress_entries 
		WHERE user_id = $1 
		ORDER BY recorded_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []models.ProgressEntry
	for rows.Next() {
		var entry models.ProgressEntry
		var weight, bodyFat sql.NullFloat64
		
		err := rows.Scan(
			&entry.ID, &entry.UserID, &weight, &bodyFat,
			&entry.Measurements, &entry.Notes, &entry.RecordedAt,
			&entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if weight.Valid {
			entry.Weight = &weight.Float64
		}
		if bodyFat.Valid {
			entry.BodyFatPercentage = &bodyFat.Float64
		}
		
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// CreateProgressEntry creates a new progress entry
func (s *HealthService) CreateProgressEntry(entry *models.ProgressEntry) error {
	query := `
		INSERT INTO progress_entries (user_id, weight, body_fat_percentage,
									 measurements, notes, recorded_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	
	now := time.Now()
	if entry.RecordedAt.IsZero() {
		entry.RecordedAt = now
	}
	
	err := s.db.QueryRow(query,
		entry.UserID, entry.Weight, entry.BodyFatPercentage,
		entry.Measurements, entry.Notes, entry.RecordedAt, now, now,
	).Scan(&entry.ID)
	
	if err != nil {
		return err
	}
	
	entry.CreatedAt = now
	entry.UpdatedAt = now
	
	return nil
}
