package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/database"
	"nutrition-platform/models"

	"github.com/lib/pq"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *database.Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_verified, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsVerified,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, 
		       is_verified, is_active, last_login, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role,
		       is_verified, is_active, last_login, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.Exec(query, user.ID, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// UpdateLastLogin updates the last login time for a user
func (r *UserRepository) UpdateLastLogin(userID uint) error {
	query := `
		UPDATE users
		SET last_login = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	return nil
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(userID uint, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.Exec(query, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// DeactivateUser deactivates a user account
func (r *UserRepository) DeactivateUser(userID uint) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// VerifyEmail marks user email as verified
func (r *UserRepository) VerifyEmail(userID uint) error {
	query := `
		UPDATE users
		SET is_verified = true, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// GetUserProfile retrieves user profile
func (r *UserRepository) GetUserProfile(userID uint) (*models.UserProfile, error) {
	query := `
		SELECT id, user_id, date_of_birth, gender, height, weight,
		       activity_level, goal, dietary_restrictions, allergies,
		       preferred_units, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`
	
	profile := &models.UserProfile{}
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.DateOfBirth,
		&profile.Gender,
		&profile.Height,
		&profile.Weight,
		&profile.ActivityLevel,
		&profile.Goal,
		&profile.DietaryRestrictions,
		&profile.Allergies,
		&profile.PreferredUnits,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user profile not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	
	return profile, nil
}

// CreateUserProfile creates a user profile
func (r *UserRepository) CreateUserProfile(profile *models.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, date_of_birth, gender, height, weight,
		                          activity_level, goal, dietary_restrictions, allergies,
		                          preferred_units)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query,
		profile.UserID,
		profile.DateOfBirth,
		profile.Gender,
		profile.Height,
		profile.Weight,
		profile.ActivityLevel,
		profile.Goal,
		profile.DietaryRestrictions,
		profile.Allergies,
		profile.PreferredUnits,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	
	return nil
}

// UpdateUserProfile updates user profile
func (r *UserRepository) UpdateUserProfile(profile *models.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET date_of_birth = $2, gender = $3, height = $4, weight = $5,
		    activity_level = $6, goal = $7, dietary_restrictions = $8,
		    allergies = $9, preferred_units = $10, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	
	result, err := r.db.Exec(query,
		profile.UserID,
		profile.DateOfBirth,
		profile.Gender,
		profile.Height,
		profile.Weight,
		profile.ActivityLevel,
		profile.Goal,
		profile.DietaryRestrictions,
		profile.Allergies,
		profile.PreferredUnits,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user profile not found")
	}
	
	return nil
}

// GetUserPreferences retrieves user preferences
func (r *UserRepository) GetUserPreferences(userID uint) (*models.UserPreferences, error) {
	query := `
		SELECT id, user_id, language, timezone, notifications_enabled,
		       email_notifications, push_notifications, units, dark_mode,
		       created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`
	
	prefs := &models.UserPreferences{}
	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID,
		&prefs.UserID,
		&prefs.Language,
		&prefs.Timezone,
		&prefs.NotificationsEnabled,
		&prefs.EmailNotifications,
		&prefs.PushNotifications,
		&prefs.Units,
		&prefs.DarkMode,
		&prefs.CreatedAt,
		&prefs.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default preferences if none found
			return &models.UserPreferences{
				UserID:               userID,
				Language:             "en",
				Timezone:             "UTC",
				NotificationsEnabled: true,
				EmailNotifications:   true,
				PushNotifications:    true,
				Units:                "metric",
				DarkMode:             false,
			}, nil
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}
	
	return prefs, nil
}

// CreateOrUpdateUserPreferences creates or updates user preferences
func (r *UserRepository) CreateOrUpdateUserPreferences(prefs *models.UserPreferences) error {
	query := `
		INSERT INTO user_preferences (user_id, language, timezone, notifications_enabled,
		                             email_notifications, push_notifications, units, dark_mode)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id)
		DO UPDATE SET
			language = EXCLUDED.language,
			timezone = EXCLUDED.timezone,
			notifications_enabled = EXCLUDED.notifications_enabled,
			email_notifications = EXCLUDED.email_notifications,
			push_notifications = EXCLUDED.push_notifications,
			units = EXCLUDED.units,
			dark_mode = EXCLUDED.dark_mode,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query,
		prefs.UserID,
		prefs.Language,
		prefs.Timezone,
		prefs.NotificationsEnabled,
		prefs.EmailNotifications,
		prefs.PushNotifications,
		prefs.Units,
		prefs.DarkMode,
	).Scan(&prefs.ID, &prefs.CreatedAt, &prefs.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create/update user preferences: %w", err)
	}
	
	return nil
}

// GetActiveNutritionGoals retrieves active nutrition goals for a user
func (r *UserRepository) GetActiveNutritionGoals(userID uint) ([]*models.NutritionGoal, error) {
	query := `
		SELECT id, user_id, daily_calories, protein_grams, carbs_grams, fat_grams,
		       fiber_grams, sugar_grams, sodium_mg, water_ml, is_active,
		       start_date, end_date, created_at, updated_at
		FROM nutrition_goals
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nutrition goals: %w", err)
	}
	defer rows.Close()
	
	var goals []*models.NutritionGoal
	for rows.Next() {
		goal := &models.NutritionGoal{}
		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.DailyCalories,
			&goal.ProteinGrams,
			&goal.CarbsGrams,
			&goal.FatGrams,
			&goal.FiberGrams,
			&goal.SugarGrams,
			&goal.SodiumMg,
			&goal.WaterMl,
			&goal.IsActive,
			&goal.StartDate,
			&goal.EndDate,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan nutrition goal: %w", err)
		}
		goals = append(goals, goal)
	}
	
	return goals, nil
}

// CreateNutritionGoal creates a nutrition goal
func (r *UserRepository) CreateNutritionGoal(goal *models.NutritionGoal) error {
	query := `
		INSERT INTO nutrition_goals (user_id, daily_calories, protein_grams, carbs_grams,
		                            fat_grams, fiber_grams, sugar_grams, sodium_mg, water_ml,
		                            is_active, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query,
		goal.UserID,
		goal.DailyCalories,
		goal.ProteinGrams,
		goal.CarbsGrams,
		goal.FatGrams,
		goal.FiberGrams,
		goal.SugarGrams,
		goal.SodiumMg,
		goal.WaterMl,
		goal.IsActive,
		goal.StartDate,
		goal.EndDate,
	).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create nutrition goal: %w", err)
	}
	
	return nil
}

// GetUsers retrieves users with pagination
func (r *UserRepository) GetUsers(page, perPage int) ([]*models.User, int64, error) {
	offset := (page - 1) * perPage
	
	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM users WHERE is_active = true"
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users count: %w", err)
	}
	
	// Get users
	query := `
		SELECT id, email, password_hash, first_name, last_name, role,
		       is_verified, is_active, last_login, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()
	
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.IsVerified,
			&user.IsActive,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, total, nil
}

// SearchUsers searches users by email or name
func (r *UserRepository) SearchUsers(search string, page, perPage int) ([]*models.User, int64, error) {
	offset := (page - 1) * perPage
	searchPattern := "%" + search + "%"
	
	// Get total count
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM users 
		WHERE is_active = true 
		AND (email ILIKE $1 OR first_name ILIKE $1 OR last_name ILIKE $1)
	`
	err := r.db.QueryRow(countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users count: %w", err)
	}
	
	// Get users
	query := `
		SELECT id, email, password_hash, first_name, last_name, role,
		       is_verified, is_active, last_login, created_at, updated_at
		FROM users
		WHERE is_active = true 
		AND (email ILIKE $1 OR first_name ILIKE $1 OR last_name ILIKE $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, searchPattern, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()
	
	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.IsVerified,
			&user.IsActive,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, total, nil
}
