package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nutrition-platform/security"
	"nutrition-platform/middleware"
	"nutrition-platform/models"
	"nutrition-platform/services"
)

// TestAPIKeyGeneration tests the API key generation functionality
func TestAPIKeyGeneration(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		length   int
		wantErr  bool
		errMsg   string
	}{
		{
			name:   "valid_key_generation",
			prefix: "nk_",
			length: 64,
			wantErr: false,
		},
		{
			name:   "minimum_length_key",
			prefix: "np_",
			length: 32,
			wantErr: false,
		},
		{
			name:    "too_short_key",
			prefix:  "nt_",
			length:  20,
			wantErr: true,
			errMsg:  "API key length must be at least 32 characters",
		},
		{
			name:    "prefix_too_long",
			prefix:  "very_long_prefix_that_exceeds_length_",
			length:  32,
			wantErr: true,
			errMsg:  "prefix is too long for the specified key length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey, err := security.GenerateSecureAPIKey(tt.prefix, tt.length)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.length, len(apiKey))
			assert.True(t, strings.HasPrefix(apiKey, tt.prefix))

			// Test uniqueness by generating multiple keys
			apiKey2, err := security.GenerateSecureAPIKey(tt.prefix, tt.length)
			require.NoError(t, err)
			assert.NotEqual(t, apiKey, apiKey2, "Generated keys should be unique")
		})
	}
}

// TestAPIKeyValidation tests the comprehensive API key validation
func TestAPIKeyValidation(t *testing.T) {
	validator := security.NewAPIKeyValidator()

	tests := []struct {
		name           string
		apiKey         string
		expectedValid  bool
		expectedLevel  security.SecurityLevel
		minScore       int
		maxScore       int
		expectedIssues []string
	}{
		{
			name:           "secure_key",
			apiKey:         "nk_8Kf9mN2pQ7rS3tU6vW9xY1zA4bC7dE0fG",
			expectedValid:  true,
			expectedLevel:  security.SecurityLevelHigh,
			minScore:       70,
			maxScore:       100,
			expectedIssues: []string{},
		},
		{
			name:           "too_short_key",
			apiKey:         "nk_short",
			expectedValid:  false,
			expectedLevel:  security.SecurityLevelLow,
			minScore:       0,
			maxScore:       30,
			expectedIssues: []string{"length"},
		},
		{
			name:           "no_prefix_key",
			apiKey:         "8Kf9mN2pQ7rS3tU6vW9xY1zA4bC7dE0fG3hI5j",
			expectedValid:  true,
			expectedLevel:  security.SecurityLevelMedium,
			minScore:       50,
			maxScore:       90,
			expectedIssues: []string{"prefix"},
		},
		{
			name:           "pattern_issues_key",
			apiKey:         "nk_password123456789012345678901234567890",
			expectedValid:  false,
			expectedLevel:  security.SecurityLevelLow,
			minScore:       0,
			maxScore:       50,
			expectedIssues: []string{"pattern"},
		},
		{
			name:           "repetitive_chars_key",
			apiKey:         "nk_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expectedValid:  true,
			expectedLevel:  security.SecurityLevelLow,
			minScore:       0,
			maxScore:       60,
			expectedIssues: []string{"pattern", "distribution", "entropy"},
		},
		{
			name:           "keyboard_pattern_key",
			apiKey:         "nk_qwerty1234567890abcdefghijklmnopqr",
			expectedValid:  true,
			expectedLevel:  security.SecurityLevelMedium,
			minScore:       30,
			maxScore:       80,
			expectedIssues: []string{"pattern"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateAPIKey(tt.apiKey)

			assert.Equal(t, tt.expectedValid, result.Valid, "Validation result mismatch")
			assert.Equal(t, tt.expectedLevel, result.SecurityLevel, "Security level mismatch")
			assert.GreaterOrEqual(t, result.Score, tt.minScore, "Score too low")
			assert.LessOrEqual(t, result.Score, tt.maxScore, "Score too high")

			// Check for expected issue categories
			for _, expectedIssue := range tt.expectedIssues {
				found := false
				for _, issue := range result.Issues {
					if issue.Category == expectedIssue {
						found = true
						break
					}
				}
				assert.True(t, found, fmt.Sprintf("Expected issue category '%s' not found", expectedIssue))
			}

			// Ensure recommendations are provided for low-scoring keys
			if result.Score < 90 {
				assert.NotEmpty(t, result.Recommendations, "Recommendations should be provided for low-scoring keys")
			}
		})
	}
}

// TestAPIKeyHashing tests the hashing and validation functionality
func TestAPIKeyHashing(t *testing.T) {
	apiKey := "nk_8Kf9mN2pQ7rS3tU6vW9xY1zA4bC7dE0fG"

	// Test hashing
	hash1 := security.HashAPIKey(apiKey)
	hash2 := security.HashAPIKey(apiKey)

	assert.Equal(t, hash1, hash2, "Same input should produce same hash")
	assert.NotEqual(t, apiKey, hash1, "Hash should be different from original")
	assert.Equal(t, 64, len(hash1), "SHA-256 hash should be 64 characters")

	// Test validation
	assert.True(t, security.ValidateAPIKeyHash(apiKey, hash1), "Valid key should pass validation")
	assert.False(t, security.ValidateAPIKeyHash("wrong_key", hash1), "Invalid key should fail validation")
	assert.False(t, security.ValidateAPIKeyHash(apiKey, "wrong_hash"), "Wrong hash should fail validation")
}

// TestAPIKeyMiddleware tests the Echo middleware functionality
func TestAPIKeyMiddleware(t *testing.T) {
	e := echo.New()

	// Mock API key service
	mockService := &MockAPIKeyService{
		keys: map[string]*models.APIKey{
			"valid_hash": {
				ID:          "test-id",
				Name:        "Test Key",
				KeyHash:     "valid_hash",
				Scopes:      []string{"nutrition:read", "nutrition:write"},
				Status:      "active",
				RateLimit:   1000,
				ExpiresAt:   time.Now().Add(24 * time.Hour),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	tests := []struct {
		name           string
		header         string
		value          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid_api_key",
			header:         "X-API-Key",
			value:          "nk_valid_key_for_testing_purposes_123",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "missing_api_key",
			header:         "",
			value:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "API key required",
		},
		{
			name:           "invalid_api_key",
			header:         "X-API-Key",
			value:          "invalid_key",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid API key",
		},
		{
			name:           "authorization_header",
			header:         "Authorization",
			value:          "Bearer nk_valid_key_for_testing_purposes_123",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.header != "" {
				req.Header.Set(tt.header, tt.value)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create middleware with mock service
			mw := middleware.APIKeyMiddleware(mockService)
			handler := mw(func(c echo.Context) error {
				return c.String(http.StatusOK, "success")
			})

			err := handler(c)
			if tt.expectedStatus == http.StatusOK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			}
		})
	}
}

// TestRateLimiting tests the rate limiting functionality
func TestRateLimiting(t *testing.T) {
	e := echo.New()
	mockService := &MockAPIKeyService{
		keys: map[string]*models.APIKey{
			"rate_limited_hash": {
				ID:        "rate-test-id",
				Name:      "Rate Limited Key",
				KeyHash:   "rate_limited_hash",
				Scopes:    []string{"nutrition:read"},
				Status:    "active",
				RateLimit: 2, // Very low limit for testing
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		usage: make(map[string]int),
	}

	mw := middleware.APIKeyMiddleware(mockService)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// First request should succeed
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.Header.Set("X-API-Key", "nk_rate_limited_key_for_testing")
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)
	err1 := handler(c1)
	assert.NoError(t, err1)

	// Second request should succeed
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set("X-API-Key", "nk_rate_limited_key_for_testing")
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	err2 := handler(c2)
	assert.NoError(t, err2)

	// Third request should be rate limited
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.Header.Set("X-API-Key", "nk_rate_limited_key_for_testing")
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	err3 := handler(c3)
	assert.Error(t, err3)
	if he, ok := err3.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusTooManyRequests, he.Code)
	}
}

// TestScopeValidation tests API key scope validation
func TestScopeValidation(t *testing.T) {
	e := echo.New()
	mockService := &MockAPIKeyService{
		keys: map[string]*models.APIKey{
			"read_only_hash": {
				ID:        "read-only-id",
				Name:      "Read Only Key",
				KeyHash:   "read_only_hash",
				Scopes:    []string{"nutrition:read"},
				Status:    "active",
				RateLimit: 1000,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			"write_access_hash": {
				ID:        "write-access-id",
				Name:      "Write Access Key",
				KeyHash:   "write_access_hash",
				Scopes:    []string{"nutrition:read", "nutrition:write"},
				Status:    "active",
				RateLimit: 1000,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	tests := []struct {
		name           string
		apiKey         string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "read_only_get_allowed",
			apiKey:         "nk_read_only_key_for_testing",
			method:         http.MethodGet,
			path:           "/api/v1/nutrition/meals",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "read_only_post_denied",
			apiKey:         "nk_read_only_key_for_testing",
			method:         http.MethodPost,
			path:           "/api/v1/nutrition/meals",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "write_access_post_allowed",
			apiKey:         "nk_write_access_key_for_testing",
			method:         http.MethodPost,
			path:           "/api/v1/nutrition/meals",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "write_access_get_allowed",
			apiKey:         "nk_write_access_key_for_testing",
			method:         http.MethodGet,
			path:           "/api/v1/nutrition/meals",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("X-API-Key", tt.apiKey)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			mw := middleware.APIKeyScopeMiddleware([]string{"nutrition:read", "nutrition:write"}, mockService)
			handler := mw(func(c echo.Context) error {
				return c.String(http.StatusOK, "success")
			})

			err := handler(c)
			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			}
		})
	}
}

// TestSecurityAudit tests the security audit functionality
func TestSecurityAudit(t *testing.T) {
	audit := security.PerformSecurityAudit()

	assert.NotNil(t, audit)
	assert.NotZero(t, audit.Timestamp)
	assert.GreaterOrEqual(t, audit.OverallScore, 70, "Security audit score should be high")
	assert.NotEmpty(t, audit.Findings, "Audit should have findings")
	assert.NotEmpty(t, audit.Summary, "Audit should have summary")

	// Check for required audit categories
	requiredCategories := []string{"key_generation", "storage", "comparison", "rate_limiting", "monitoring"}
	for _, category := range requiredCategories {
		found := false
		for _, finding := range audit.Findings {
			if finding.Category == category {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("Required audit category '%s' not found", category))
	}

	// Verify summary statistics
	summary := audit.Summary
	assert.Contains(t, summary, "total_findings")
	assert.Contains(t, summary, "critical_issues")
	assert.Contains(t, summary, "high_issues")
	assert.Contains(t, summary, "medium_issues")
	assert.Contains(t, summary, "low_issues")
	assert.Contains(t, summary, "info_items")
}

// TestAPIKeyEndpoints tests the API key management endpoints
func TestAPIKeyEndpoints(t *testing.T) {
	e := echo.New()
	mockService := &MockAPIKeyService{
		keys: make(map[string]*models.APIKey),
	}

	// Test create API key endpoint
	t.Run("create_api_key", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "Test API Key",
			"scopes":      []string{"nutrition:read", "nutrition:write"},
			"rate_limit":  1000,
			"expires_in":  "30d",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/admin/api-keys", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := handlers.CreateAPIKey(mockService)
		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "api_key")
		assert.Contains(t, response, "key_id")
	})

	// Test get API key endpoint
	t.Run("get_api_key", func(t *testing.T) {
		// First create a key
		mockService.keys["test-hash"] = &models.APIKey{
			ID:        "test-id",
			Name:      "Test Key",
			KeyHash:   "test-hash",
			Scopes:    []string{"nutrition:read"},
			Status:    "active",
			RateLimit: 1000,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := httptest.NewRequest(http.MethodGet, "/admin/api-keys/test-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		handler := handlers.GetAPIKey(mockService)
		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-id", response.ID)
		assert.Equal(t, "Test Key", response.Name)
	})
}

// MockAPIKeyService implements the APIKeyService interface for testing
type MockAPIKeyService struct {
	keys  map[string]*models.APIKey
	usage map[string]int
}

func (m *MockAPIKeyService) CreateAPIKey(name string, scopes []string, rateLimit int, expiresIn time.Duration, userID string) (*models.APIKey, string, error) {
	apiKey, err := security.GenerateSecureAPIKey("nk_", 64)
	if err != nil {
		return nil, "", err
	}

	hash := security.HashAPIKey(apiKey)
	key := &models.APIKey{
		ID:        fmt.Sprintf("key-%d", len(m.keys)+1),
		Name:      name,
		KeyHash:   hash,
		Scopes:    scopes,
		Status:    "active",
		RateLimit: rateLimit,
		ExpiresAt: time.Now().Add(expiresIn),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.keys[hash] = key
	return key, apiKey, nil
}

func (m *MockAPIKeyService) ValidateAPIKey(keyHash string) (*models.APIKey, error) {
	if key, exists := m.keys[keyHash]; exists {
		if key.Status == "active" && key.ExpiresAt.After(time.Now()) {
			return key, nil
		}
		return nil, fmt.Errorf("API key is inactive or expired")
	}
	return nil, fmt.Errorf("API key not found")
}

func (m *MockAPIKeyService) UpdateUsage(keyID, endpoint, method, userAgent, ipAddress string, responseTime time.Duration, statusCode int) error {
	if m.usage == nil {
		m.usage = make(map[string]int)
	}
	m.usage[keyID]++
	return nil
}

func (m *MockAPIKeyService) GetAPIKey(keyID string) (*models.APIKey, error) {
	for _, key := range m.keys {
		if key.ID == keyID {
			return key, nil
		}
	}
	return nil, fmt.Errorf("API key not found")
}

func (m *MockAPIKeyService) RevokeAPIKey(keyID string) error {
	for _, key := range m.keys {
		if key.ID == keyID {
			key.Status = "revoked"
			key.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("API key not found")
}

func (m *MockAPIKeyService) GetAPIKeyStats(keyID string) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"total_requests": m.usage[keyID],
		"last_used":      time.Now(),
		"status":         "active",
	}
	return stats, nil
}

func (m *MockAPIKeyService) CheckRateLimit(keyID string) (bool, error) {
	if key, exists := m.keys[keyID]; exists {
		return m.usage[keyID] < key.RateLimit, nil
	}
	return false, fmt.Errorf("API key not found")
}

// BenchmarkAPIKeyValidation benchmarks the API key validation performance
func BenchmarkAPIKeyValidation(b *testing.B) {
	validator := security.NewAPIKeyValidator()
	apiKey := "nk_8Kf9mN2pQ7rS3tU6vW9xY1zA4bC7dE0fG3hI5j"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateAPIKey(apiKey)
	}
}

// BenchmarkAPIKeyHashing benchmarks the API key hashing performance
func BenchmarkAPIKeyHashing(b *testing.B) {
	apiKey := "nk_8Kf9mN2pQ7rS3tU6vW9xY1zA4bC7dE0fG3hI5j"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		security.HashAPIKey(apiKey)
	}
}

// BenchmarkConstantTimeCompare benchmarks the constant-time comparison
func BenchmarkConstantTimeCompare(b *testing.B) {
	hash1 := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	hash2 := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		security.ValidateAPIKeyHash("test", hash1)
		security.ValidateAPIKeyHash("test", hash2)
	}
}