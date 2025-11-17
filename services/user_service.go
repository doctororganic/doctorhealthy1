package services

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// UserService handles user-related business logic
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, 
			   is_verified, is_active, last_login, created_at, updated_at
		FROM users 
		WHERE id = $1
	`
	
	var user models.User
	var lastLogin sql.NullTime
	
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, 
		&user.LastName, &user.Role, &user.IsVerified, &user.IsActive,
		&lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, 
			   is_verified, is_active, last_login, created_at, updated_at
		FROM users 
		WHERE email = $1
	`
	
	var user models.User
	var lastLogin sql.NullTime
	
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, 
		&user.LastName, &user.Role, &user.IsVerified, &user.IsActive,
		&lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	
	return &user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role, 
						   is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	
	now := time.Now()
	err := s.db.QueryRow(query,
		user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.IsVerified, user.IsActive, now, now,
	).Scan(&user.ID)
	
	if err != nil {
		return err
	}
	
	user.CreatedAt = now
	user.UpdatedAt = now
	
	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET email = $2, password_hash = $3, first_name = $4, last_name = $5,
			role = $6, is_verified = $7, is_active = $8, last_login = $9, updated_at = $10
		WHERE id = $1
	`
	
	user.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.IsVerified, user.IsActive, user.LastLogin, user.UpdatedAt,
	)
	
	return err
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	
	_, err := s.db.Exec(query, id)
	return err
}

// UpdateLastLogin updates the last login time for a user
func (s *UserService) UpdateLastLogin(id uint) error {
	query := `UPDATE users SET last_login = $1 WHERE id = $2`
	
	_, err := s.db.Exec(query, time.Now(), id)
	return err
}
