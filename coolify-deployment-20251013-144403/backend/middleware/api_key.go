package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"nutrition-platform/models"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RateLimiter manages rate limiting for API keys
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
	}

	// Start cleanup routine
	go rl.cleanup()

	return rl
}

// IsAllowed checks if a request is allowed based on rate limit
func (rl *RateLimiter) IsAllowed(apiKeyID string, limit int) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	window := now.Add(-time.Minute) // 1-minute window

	// Get existing requests for this API key
	requests, exists := rl.requests[apiKeyID]
	if !exists {
		requests = []time.Time{}
	}

	// Filter requests within the time window
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(window) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if limit is exceeded
	if len(validRequests) >= limit {
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[apiKeyID] = validRequests

	return true
}

// cleanup removes old entries from the rate limiter
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		window := now.Add(-time.Minute)

		for apiKeyID, requests := range rl.requests {
			var validRequests []time.Time
			for _, req := range requests {
				if req.After(window) {
					validRequests = append(validRequests, req)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, apiKeyID)
			} else {
				rl.requests[apiKeyID] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// APIKeyMiddleware creates middleware for API key authentication
func APIKeyMiddleware(apiKeyService *services.APIKeyService) echo.MiddlewareFunc {
	rateLimiter := NewRateLimiter()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Extract API key from header
			apiKey := extractAPIKey(c)
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"error":   "missing_api_key",
					"message": "API key is required. Please provide it in the Authorization header as 'Bearer <api_key>' or in the X-API-Key header.",
					"code":    "AUTH_001",
				})
			}

			// Validate API key
			apiKeyModel, err := apiKeyService.ValidateAPIKey(apiKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"error":   "invalid_api_key",
					"message": "The provided API key is invalid, expired, or revoked.",
					"code":    "AUTH_002",
				})
			}

			// Check rate limit
			if !rateLimiter.IsAllowed(apiKeyModel.ID, apiKeyModel.RateLimit) {
				return echo.NewHTTPError(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "rate_limit_exceeded",
					"message":     fmt.Sprintf("Rate limit of %d requests per minute exceeded.", apiKeyModel.RateLimit),
					"code":        "RATE_001",
					"retry_after": 60,
				})
			}

			// Check endpoint access
			if !apiKeyModel.CanAccess(c.Request().URL.Path, c.Request().Method) {
				return echo.NewHTTPError(http.StatusForbidden, map[string]interface{}{
					"error":           "insufficient_permissions",
					"message":         "Your API key does not have permission to access this endpoint.",
					"code":            "AUTH_003",
					"required_scopes": getRequiredScopes(c.Request().URL.Path, c.Request().Method),
					"your_scopes":     apiKeyModel.Scopes,
				})
			}

			// Store API key info in context
			c.Set("api_key", apiKeyModel)
			c.Set("user_id", apiKeyModel.UserID)
			c.Set("api_key_id", apiKeyModel.ID)

			// Process request
			err = next(c)

			// Record usage statistics (async)
			go func() {
				responseTime := time.Since(start).Milliseconds()
				statusCode := c.Response().Status
				ipAddress := c.RealIP()
				userAgent := c.Request().UserAgent()

				if recordErr := apiKeyService.UpdateAPIKeyUsage(
					apiKeyModel,
					c.Request().URL.Path,
					c.Request().Method,
					statusCode,
					responseTime,
					ipAddress,
					userAgent,
				); recordErr != nil {
					// Log error but don't fail the request
					echo.New().Logger.Errorf("Failed to record API key usage: %v", recordErr)
				}
			}()

			return err
		}
	}
}

// OptionalAPIKeyMiddleware creates middleware for optional API key authentication
// If API key is provided, it validates it; if not, the request continues without authentication
func OptionalAPIKeyMiddleware(apiKeyService *services.APIKeyService) echo.MiddlewareFunc {
	rateLimiter := NewRateLimiter()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Extract API key from header
			apiKey := extractAPIKey(c)
			if apiKey == "" {
				// No API key provided, continue without authentication
				return next(c)
			}

			// Validate API key if provided
			apiKeyModel, err := apiKeyService.ValidateAPIKey(apiKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"error":   "invalid_api_key",
					"message": "The provided API key is invalid, expired, or revoked.",
					"code":    "AUTH_002",
				})
			}

			// Check rate limit
			if !rateLimiter.IsAllowed(apiKeyModel.ID, apiKeyModel.RateLimit) {
				return echo.NewHTTPError(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "rate_limit_exceeded",
					"message":     fmt.Sprintf("Rate limit of %d requests per minute exceeded.", apiKeyModel.RateLimit),
					"code":        "RATE_001",
					"retry_after": 60,
				})
			}

			// Store API key info in context
			c.Set("api_key", apiKeyModel)
			c.Set("user_id", apiKeyModel.UserID)
			c.Set("api_key_id", apiKeyModel.ID)
			c.Set("authenticated", true)

			// Process request
			err = next(c)

			// Record usage statistics (async)
			go func() {
				responseTime := time.Since(start).Milliseconds()
				statusCode := c.Response().Status
				ipAddress := c.RealIP()
				userAgent := c.Request().UserAgent()

				if recordErr := apiKeyService.UpdateAPIKeyUsage(
					apiKeyModel,
					c.Request().URL.Path,
					c.Request().Method,
					statusCode,
					responseTime,
					ipAddress,
					userAgent,
				); recordErr != nil {
					echo.New().Logger.Errorf("Failed to record API key usage: %v", recordErr)
				}
			}()

			return err
		}
	}
}

// ScopeMiddleware creates middleware that requires specific scopes
func ScopeMiddleware(requiredScopes ...models.APIKeyScope) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey, ok := c.Get("api_key").(*models.APIKey)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"error":   "authentication_required",
					"message": "This endpoint requires authentication.",
					"code":    "AUTH_004",
				})
			}

			// Check if API key has any of the required scopes
			hasRequiredScope := false
			for _, requiredScope := range requiredScopes {
				if apiKey.HasScope(requiredScope) {
					hasRequiredScope = true
					break
				}
			}

			if !hasRequiredScope {
				return echo.NewHTTPError(http.StatusForbidden, map[string]interface{}{
					"error":           "insufficient_scope",
					"message":         "Your API key does not have the required scope for this operation.",
					"code":            "AUTH_005",
					"required_scopes": requiredScopes,
					"your_scopes":     apiKey.Scopes,
				})
			}

			return next(c)
		}
	}
}

// AdminOnlyMiddleware creates middleware that requires admin scope
func AdminOnlyMiddleware() echo.MiddlewareFunc {
	return ScopeMiddleware(models.ScopeAdmin)
}

// ReadWriteMiddleware creates middleware that requires read-write or admin scope
func ReadWriteMiddleware() echo.MiddlewareFunc {
	return ScopeMiddleware(models.ScopeReadWrite, models.ScopeAdmin)
}

// Helper functions

// extractAPIKey extracts the API key from the request headers
func extractAPIKey(c echo.Context) string {
	// Try Authorization header first (Bearer token)
	auth := c.Request().Header.Get("Authorization")
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Try X-API-Key header
	apiKey := c.Request().Header.Get("X-API-Key")
	if apiKey != "" {
		return apiKey
	}

	// DO NOT extract from query parameters for security reasons
	// Query parameters are less secure and can be logged in server logs
	return ""
}

// getRequiredScopes returns the required scopes for an endpoint
func getRequiredScopes(endpoint, method string) []models.APIKeyScope {
	// Define endpoint to scope mapping
	endpointScopes := map[string][]models.APIKeyScope{
		"/api/v1/nutrition":   {models.ScopeNutrition, models.ScopeReadOnly},
		"/api/v1/meals":       {models.ScopeMeals, models.ScopeReadOnly},
		"/api/v1/workouts":    {models.ScopeWorkouts, models.ScopeReadOnly},
		"/api/v1/health":      {models.ScopeHealth, models.ScopeReadOnly},
		"/api/v1/supplements": {models.ScopeSupplements, models.ScopeReadOnly},
	}

	// Check for specific endpoint matches
	for endpointPrefix, scopes := range endpointScopes {
		if strings.HasPrefix(endpoint, endpointPrefix) {
			if method == "GET" || method == "HEAD" {
				return scopes
			} else {
				// For write operations, require read-write or admin
				return append(scopes, models.ScopeReadWrite, models.ScopeAdmin)
			}
		}
	}

	// Default scopes
	if method == "GET" || method == "HEAD" {
		return []models.APIKeyScope{models.ScopeReadOnly}
	}
	return []models.APIKeyScope{models.ScopeReadWrite}
}

// GetAPIKeyFromContext retrieves the API key from the Echo context
func GetAPIKeyFromContext(c echo.Context) *models.APIKey {
	apiKey, ok := c.Get("api_key").(*models.APIKey)
	if !ok {
		return nil
	}
	return apiKey
}

// GetUserIDFromContext retrieves the user ID from the Echo context
func GetUserIDFromContext(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}

// IsAuthenticated checks if the request is authenticated via API key
func IsAuthenticated(c echo.Context) bool {
	authenticated, ok := c.Get("authenticated").(bool)
	return ok && authenticated
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			return next(c)
		}
	}
}

// RequestLogging logs API requests with API key information
func RequestLogging() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `{"time":"${time_rfc3339}","method":"${method}","uri":"${uri}","status":${status},"latency":"${latency_human}","api_key_id":"${header:x-api-key-id}","user_id":"${header:x-user-id}","ip":"${remote_ip}","user_agent":"${user_agent}"}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	})
}
