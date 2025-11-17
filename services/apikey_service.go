package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// APIKeyService handles API key management business logic
type APIKeyService struct {
	db *sql.DB
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(db *sql.DB) *APIKeyService {
	return &APIKeyService{
		db: db,
	}
}

// GetAPIKeys retrieves API keys for a user
func (s *APIKeyService) GetAPIKeys(userID uint) ([]models.APIKey, error) {
	query := `
		SELECT id, user_id, name, api_key, permissions, is_active,
			   expires_at, last_used_at, created_at
		FROM api_keys 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var expiresAt, lastUsedAt sql.NullTime
		
		err := rows.Scan(
			&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.APIKey,
			&apiKey.Permissions, &apiKey.IsActive, &expiresAt,
			&lastUsedAt, &apiKey.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if expiresAt.Valid {
			apiKey.ExpiresAt = &expiresAt.Time
		}
		if lastUsedAt.Valid {
			apiKey.LastUsedAt = &lastUsedAt.Time
		}
		
		apiKeys = append(apiKeys, apiKey)
	}
	
	return apiKeys, nil
}

// CreateAPIKey creates a new API key
func (s *APIKeyService) CreateAPIKey(apiKey *models.APIKey) error {
	// Generate a new API key
	generatedKey, err := s.generateAPIKey()
	if err != nil {
		return err
	}
	
	apiKey.APIKey = generatedKey
	
	query := `
		INSERT INTO api_keys (user_id, name, api_key, permissions, is_active,
							expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	
	now := time.Now()
	
	err = s.db.QueryRow(query,
		apiKey.UserID, apiKey.Name, apiKey.APIKey, apiKey.Permissions,
		apiKey.IsActive, apiKey.ExpiresAt, now,
	).Scan(&apiKey.ID)
	
	if err != nil {
		return err
	}
	
	apiKey.CreatedAt = now
	
	return nil
}

// GetAPIKey retrieves an API key by ID
func (s *APIKeyService) GetAPIKey(id, userID uint) (*models.APIKey, error) {
	query := `
		SELECT id, user_id, name, api_key, permissions, is_active,
			   expires_at, last_used_at, created_at
		FROM api_keys 
		WHERE id = $1 AND user_id = $2
	`
	
	var apiKey models.APIKey
	var expiresAt, lastUsedAt sql.NullTime
	
	err := s.db.QueryRow(query, id, userID).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.APIKey,
		&apiKey.Permissions, &apiKey.IsActive, &expiresAt,
		&lastUsedAt, &apiKey.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, err
	}
	
	if expiresAt.Valid {
		apiKey.ExpiresAt = &expiresAt.Time
	}
	if lastUsedAt.Valid {
		apiKey.LastUsedAt = &lastUsedAt.Time
	}
	
	return &apiKey, nil
}

// ValidateAPIKey validates an API key and returns user information
func (s *APIKeyService) ValidateAPIKey(keyString string) (*models.APIKey, error) {
	query := `
		SELECT id, user_id, name, api_key, permissions, is_active,
			   expires_at, last_used_at, created_at
		FROM api_keys 
		WHERE api_key = $1 AND is_active = true
	`
	
	var apiKey models.APIKey
	var expiresAt, lastUsedAt sql.NullTime
	
	err := s.db.QueryRow(query, keyString).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.APIKey,
		&apiKey.Permissions, &apiKey.IsActive, &expiresAt,
		&lastUsedAt, &apiKey.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid API key")
		}
		return nil, err
	}
	
	// Check if key has expired
	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("API key has expired")
	}
	
	if expiresAt.Valid {
		apiKey.ExpiresAt = &expiresAt.Time
	}
	if lastUsedAt.Valid {
		apiKey.LastUsedAt = &lastUsedAt.Time
	}
	
	// Update last used time
	s.updateLastUsed(apiKey.ID)
	
	return &apiKey, nil
}

// UpdateAPIKey updates an existing API key
func (s *APIKeyService) UpdateAPIKey(apiKey *models.APIKey) error {
	query := `
		UPDATE api_keys 
		SET name = $2, permissions = $3, is_active = $4, expires_at = $5
		WHERE id = $1 AND user_id = $6
	`
	
	_, err := s.db.Exec(query,
		apiKey.ID, apiKey.Name, apiKey.Permissions,
		apiKey.IsActive, apiKey.ExpiresAt, apiKey.UserID,
	)
	
	return err
}

// DeleteAPIKey deletes an API key
func (s *APIKeyService) DeleteAPIKey(id, userID uint) error {
	query := `DELETE FROM api_keys WHERE id = $1 AND user_id = $2`
	
	_, err := s.db.Exec(query, id, userID)
	return err
}

// RegenerateAPIKey regenerates an API key
func (s *APIKeyService) RegenerateAPIKey(id, userID uint) (string, error) {
	// Generate a new API key
	newKey, err := s.generateAPIKey()
	if err != nil {
		return "", err
	}
	
	query := `UPDATE api_keys SET api_key = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`
	
	_, err = s.db.Exec(query, newKey, time.Now(), id, userID)
	if err != nil {
		return "", err
	}
	
	return newKey, nil
}

// generateAPIKey generates a new random API key
func (s *APIKeyService) generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "nk_" + hex.EncodeToString(bytes), nil
}

// updateLastUsed updates the last used timestamp for an API key
func (s *APIKeyService) updateLastUsed(id uint) error {
	query := `UPDATE api_keys SET last_used_at = $1 WHERE id = $2`
	
	_, err := s.db.Exec(query, time.Now(), id)
	return err
}
