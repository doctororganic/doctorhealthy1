package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"nutrition-platform/models"
)

// APIKeyService handles API key operations
type APIKeyService struct {
	db    *sql.DB
	cache map[string]*models.APIKey
	mutex sync.RWMutex
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(db *sql.DB) *APIKeyService {
	return &APIKeyService{
		db:    db,
		cache: make(map[string]*models.APIKey),
	}
}

// CreateAPIKey creates a new API key
func (s *APIKeyService) CreateAPIKey(userID string, req *models.CreateAPIKeyRequest) (*models.CreateAPIKeyResponse, error) {
	// Validate scopes
	if err := models.ValidateScopes(req.Scopes); err != nil {
		return nil, fmt.Errorf("invalid scopes: %w", err)
	}

	// Generate API key
	key, keyHash, err := models.GenerateAPIKey("nk")
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Calculate expiration
	var expiresAt *time.Time
	if req.ExpiresIn != nil {
		expiry := time.Now().AddDate(0, 0, *req.ExpiresIn)
		expiresAt = &expiry
	}

	// Create API key model
	apiKey := &models.APIKey{
		ID:        generateID(),
		Name:      req.Name,
		KeyHash:   keyHash,
		Prefix:    models.GetPrefix(key),
		UserID:    userID,
		Status:    models.APIKeyStatusActive,
		Scopes:    req.Scopes,
		RateLimit: req.RateLimit,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  req.Metadata,
	}

	// Store in database
	if err := s.storeAPIKey(apiKey); err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}

	// Cache the API key
	s.mutex.Lock()
	s.cache[keyHash] = apiKey
	s.mutex.Unlock()

	// Prepare response
	response := &models.CreateAPIKeyResponse{
		APIKey: apiKey,
		Key:    key,
	}

	if expiresAt != nil {
		response.Warning = fmt.Sprintf("This API key will expire on %s", expiresAt.Format("2006-01-02"))
	}

	return response, nil
}

// ValidateAPIKey validates an API key and returns the associated API key model
func (s *APIKeyService) ValidateAPIKey(key string) (*models.APIKey, error) {
	// Validate format
	if !models.ValidateAPIKeyFormat(key) {
		return nil, fmt.Errorf("invalid API key format")
	}

	// Hash the key
	keyHash := models.HashAPIKey(key)

	// Check cache first
	s.mutex.RLock()
	apiKey, exists := s.cache[keyHash]
	s.mutex.RUnlock()

	if exists {
		if !apiKey.IsActive() {
			return nil, fmt.Errorf("API key is inactive or expired")
		}
		return apiKey, nil
	}

	// Load from database
	apiKey, err := s.loadAPIKeyByHash(keyHash)
	if err != nil {
		return nil, fmt.Errorf("API key not found or invalid: %w", err)
	}

	if !apiKey.IsActive() {
		return nil, fmt.Errorf("API key is inactive or expired")
	}

	// Cache the result
	s.mutex.Lock()
	s.cache[keyHash] = apiKey
	s.mutex.Unlock()

	return apiKey, nil
}

// UpdateAPIKeyUsage records API key usage
func (s *APIKeyService) UpdateAPIKeyUsage(apiKey *models.APIKey, endpoint, method string, statusCode int, responseTime int64, ipAddress, userAgent string) error {
	// Update last used timestamp
	apiKey.UpdateLastUsed()

	// Update in database
	if err := s.updateLastUsed(apiKey.ID); err != nil {
		log.Printf("Failed to update API key last used: %v", err)
	}

	// Record usage statistics
	usage := &models.APIKeyUsage{
		APIKeyID:     apiKey.ID,
		Endpoint:     endpoint,
		Method:       method,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Timestamp:    time.Now(),
	}

	return s.recordUsage(usage)
}

// GetAPIKeys retrieves API keys for a user
func (s *APIKeyService) GetAPIKeys(userID string, page, limit int) (*models.APIKeyListResponse, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, name, prefix, status, scopes, rate_limit, expires_at, last_used_at, created_at, updated_at, metadata,
		       COUNT(*) OVER() as total_count
		FROM api_keys 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	var totalCount int

	for rows.Next() {
		var apiKey models.APIKey
		var scopesJSON []byte
		var metadataJSON []byte

		err := rows.Scan(
			&apiKey.ID,
			&apiKey.Name,
			&apiKey.Prefix,
			&apiKey.Status,
			&scopesJSON,
			&apiKey.RateLimit,
			&apiKey.ExpiresAt,
			&apiKey.LastUsedAt,
			&apiKey.CreatedAt,
			&apiKey.UpdatedAt,
			&metadataJSON,
			&totalCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}

		// Parse JSON fields
		if err := json.Unmarshal(scopesJSON, &apiKey.Scopes); err != nil {
			return nil, fmt.Errorf("failed to parse scopes: %w", err)
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &apiKey.Metadata); err != nil {
				return nil, fmt.Errorf("failed to parse metadata: %w", err)
			}
		}

		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API keys: %w", err)
	}

	return &models.APIKeyListResponse{
		APIKeys: apiKeys,
		Total:   totalCount,
		Page:    page,
		Limit:   limit,
	}, nil
}

// RevokeAPIKey revokes an API key
func (s *APIKeyService) RevokeAPIKey(userID, apiKeyID string) error {
	query := `UPDATE api_keys SET status = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`
	_, err := s.db.Exec(query, models.APIKeyStatusRevoked, time.Now(), apiKeyID, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	// Remove from cache
	s.mutex.Lock()
	for hash, apiKey := range s.cache {
		if apiKey.ID == apiKeyID {
			delete(s.cache, hash)
			break
		}
	}
	s.mutex.Unlock()

	return nil
}

// GetAPIKeyStats retrieves usage statistics for an API key
func (s *APIKeyService) GetAPIKeyStats(userID, apiKeyID string, days int) (*models.APIKeyStatsResponse, error) {
	// Verify ownership
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM api_keys WHERE id = $1 AND user_id = $2)", apiKeyID, userID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to verify API key ownership: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("API key not found")
	}

	// Get basic stats
	stats := &models.APIKeyStatsResponse{
		APIKeyID:    apiKeyID,
		Endpoints:   make(map[string]int64),
		StatusCodes: make(map[string]int64),
	}

	// Get total requests and last used
	err = s.db.QueryRow(`
		SELECT COUNT(*), MAX(timestamp) 
		FROM api_key_usage 
		WHERE api_key_id = $1 AND timestamp >= $2
	`, apiKeyID, time.Now().AddDate(0, 0, -days)).Scan(&stats.TotalRequests, &stats.LastUsed)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	// Get endpoint stats
	endpointRows, err := s.db.Query(`
		SELECT endpoint, COUNT(*) 
		FROM api_key_usage 
		WHERE api_key_id = $1 AND timestamp >= $2 
		GROUP BY endpoint
	`, apiKeyID, time.Now().AddDate(0, 0, -days))
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint stats: %w", err)
	}
	defer endpointRows.Close()

	for endpointRows.Next() {
		var endpoint string
		var count int64
		if err := endpointRows.Scan(&endpoint, &count); err != nil {
			return nil, fmt.Errorf("failed to scan endpoint stats: %w", err)
		}
		stats.Endpoints[endpoint] = count
	}

	// Get status code stats
	statusRows, err := s.db.Query(`
		SELECT status_code, COUNT(*) 
		FROM api_key_usage 
		WHERE api_key_id = $1 AND timestamp >= $2 
		GROUP BY status_code
	`, apiKeyID, time.Now().AddDate(0, 0, -days))
	if err != nil {
		return nil, fmt.Errorf("failed to get status code stats: %w", err)
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var statusCode int
		var count int64
		if err := statusRows.Scan(&statusCode, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status code stats: %w", err)
		}
		stats.StatusCodes[fmt.Sprintf("%d", statusCode)] = count
	}

	// Get daily usage
	dailyRows, err := s.db.Query(`
		SELECT DATE(timestamp) as date, COUNT(*) 
		FROM api_key_usage 
		WHERE api_key_id = $1 AND timestamp >= $2 
		GROUP BY DATE(timestamp) 
		ORDER BY date
	`, apiKeyID, time.Now().AddDate(0, 0, -days))
	if err != nil {
		return nil, fmt.Errorf("failed to get daily usage: %w", err)
	}
	defer dailyRows.Close()

	for dailyRows.Next() {
		var date string
		var count int64
		if err := dailyRows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("failed to scan daily usage: %w", err)
		}
		stats.DailyUsage = append(stats.DailyUsage, models.DailyUsage{
			Date:     date,
			Requests: count,
		})
	}

	return stats, nil
}

// Helper methods

func (s *APIKeyService) storeAPIKey(apiKey *models.APIKey) error {
	scopesJSON, err := json.Marshal(apiKey.Scopes)
	if err != nil {
		return fmt.Errorf("failed to marshal scopes: %w", err)
	}

	metadataJSON, err := json.Marshal(apiKey.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO api_keys (id, name, key_hash, prefix, user_id, status, scopes, rate_limit, expires_at, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = s.db.Exec(query,
		apiKey.ID,
		apiKey.Name,
		apiKey.KeyHash,
		apiKey.Prefix,
		apiKey.UserID,
		apiKey.Status,
		scopesJSON,
		apiKey.RateLimit,
		apiKey.ExpiresAt,
		apiKey.CreatedAt,
		apiKey.UpdatedAt,
		metadataJSON,
	)

	return err
}

func (s *APIKeyService) loadAPIKeyByHash(keyHash string) (*models.APIKey, error) {
	query := `
		SELECT id, name, prefix, user_id, status, scopes, rate_limit, expires_at, last_used_at, created_at, updated_at, metadata
		FROM api_keys 
		WHERE key_hash = $1
	`

	var apiKey models.APIKey
	var scopesJSON []byte
	var metadataJSON []byte

	err := s.db.QueryRow(query, keyHash).Scan(
		&apiKey.ID,
		&apiKey.Name,
		&apiKey.Prefix,
		&apiKey.UserID,
		&apiKey.Status,
		&scopesJSON,
		&apiKey.RateLimit,
		&apiKey.ExpiresAt,
		&apiKey.LastUsedAt,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
		&metadataJSON,
	)

	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if err := json.Unmarshal(scopesJSON, &apiKey.Scopes); err != nil {
		return nil, fmt.Errorf("failed to parse scopes: %w", err)
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &apiKey.Metadata); err != nil {
			return nil, fmt.Errorf("failed to parse metadata: %w", err)
		}
	}

	apiKey.KeyHash = keyHash
	return &apiKey, nil
}

func (s *APIKeyService) updateLastUsed(apiKeyID string) error {
	query := `UPDATE api_keys SET last_used_at = $1, updated_at = $1 WHERE id = $2`
	_, err := s.db.Exec(query, time.Now(), apiKeyID)
	return err
}

func (s *APIKeyService) recordUsage(usage *models.APIKeyUsage) error {
	query := `
		INSERT INTO api_key_usage (api_key_id, endpoint, method, status_code, response_time, ip_address, user_agent, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.Exec(query,
		usage.APIKeyID,
		usage.Endpoint,
		usage.Method,
		usage.StatusCode,
		usage.ResponseTime,
		usage.IPAddress,
		usage.UserAgent,
		usage.Timestamp,
	)

	return err
}

// generateID generates a unique ID for API keys
func generateID() string {
	return fmt.Sprintf("ak_%d_%s", time.Now().Unix(), generateRandomString(8))
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
