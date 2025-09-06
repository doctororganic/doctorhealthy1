package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"nutrition-platform/middleware"
	"nutrition-platform/models"
	"nutrition-platform/services"
)

// APIKeyHandler handles API key related requests
type APIKeyHandler struct {
	apiKeyService *services.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(apiKeyService *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// CreateAPIKey creates a new API key
// POST /api/v1/api-keys
func (h *APIKeyHandler) CreateAPIKey(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"error": "authentication_required",
			"message": "User authentication is required to create API keys.",
			"code": "AUTH_006",
		})
	}

	// Parse request
	var req models.CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid_request",
			"message": "Invalid request format. Please check your JSON payload.",
			"code": "REQ_001",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := validateCreateAPIKeyRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "validation_failed",
			"message": "Request validation failed.",
			"code": "VAL_001",
			"details": err.Error(),
		})
	}

	// Set default rate limit if not provided
	if req.RateLimit == 0 {
		req.RateLimit = 1000 // Default 1000 requests per minute
	}

	// Create API key
	response, err := h.apiKeyService.CreateAPIKey(userID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]interface{}{
			"error": "creation_failed",
			"message": "Failed to create API key. Please try again.",
			"code": "SRV_001",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "API key created successfully",
		"data": response,
		"security_notice": "Please store this API key securely. It will not be shown again.",
	})
}

// GetAPIKeys retrieves API keys for the authenticated user
// GET /api/v1/api-keys
func (h *APIKeyHandler) GetAPIKeys(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"error": "authentication_required",
			"message": "User authentication is required to view API keys.",
			"code": "AUTH_006",
		})
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// Get API keys
	response, err := h.apiKeyService.GetAPIKeys(userID, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]interface{}{
			"error": "retrieval_failed",
			"message": "Failed to retrieve API keys. Please try again.",
			"code": "SRV_002",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": response,
	})
}

// GetAPIKey retrieves a specific API key
// GET /api/v1/api-keys/:id
func (h *APIKeyHandler) GetAPIKey(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"error": "authentication_required",
			"message": "User authentication is required to view API key details.",
			"code": "AUTH_006",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "missing_parameter",
			"message": "API key ID is required.",
			"code": "REQ_002",
		})
	}

	// Get API key stats (default to 30 days)
	stats, err := h.apiKeyService.GetAPIKeyStats(userID, apiKeyID, 30)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]interface{}{
			"error": "not_found",
			"message": "API key not found or you don't have permission to view it.",
			"code": "RES_001",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": stats,
	})
}

// GetAPIKeyStats retrieves usage statistics for an API key
// GET /api/v1/api-keys/:id/stats
func (h *APIKeyHandler) GetAPIKeyStats(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"error": "authentication_required",
			"message": "User authentication is required to view API key statistics.",
			"code": "AUTH_006",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "missing_parameter",
			"message": "API key ID is required.",
			"code": "REQ_002",
		})
	}

	// Parse days parameter
	days, err := strconv.Atoi(c.QueryParam("days"))
	if err != nil || days < 1 || days > 365 {
		days = 30 // Default to 30 days
	}

	// Get statistics
	stats, err := h.apiKeyService.GetAPIKeyStats(userID, apiKeyID, days)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]interface{}{
			"error": "not_found",
			"message": "API key not found or you don't have permission to view it.",
			"code": "RES_001",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": stats,
		"period": map[string]interface{}{
			"days": days,
			"description": "Statistics for the last " + strconv.Itoa(days) + " days",
		},
	})
}

// RevokeAPIKey revokes an API key
// DELETE /api/v1/api-keys/:id
func (h *APIKeyHandler) RevokeAPIKey(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"error": "authentication_required",
			"message": "User authentication is required to revoke API keys.",
			"code": "AUTH_006",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "missing_parameter",
			"message": "API key ID is required.",
			"code": "REQ_002",
		})
	}

	// Revoke API key
	err := h.apiKeyService.RevokeAPIKey(userID, apiKeyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]interface{}{
			"error": "revocation_failed",
			"message": "Failed to revoke API key. It may not exist or you don't have permission.",
			"code": "SRV_003",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "API key revoked successfully",
		"api_key_id": apiKeyID,
	})
}

// ValidateAPIKey validates an API key (for testing purposes)
// POST /api/v1/api-keys/validate
func (h *APIKeyHandler) ValidateAPIKey(c echo.Context) error {
	var req struct {
		APIKey string `json:"api_key" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid_request",
			"message": "Invalid request format. Please provide an API key.",
			"code": "REQ_001",
		})
	}

	if req.APIKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": "missing_api_key",
			"message": "API key is required for validation.",
			"code": "REQ_003",
		})
	}

	// Validate API key
	apiKey, err := h.apiKeyService.ValidateAPIKey(req.APIKey)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"valid": false,
			"error": err.Error(),
			"message": "The provided API key is invalid, expired, or revoked.",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"valid": true,
		"api_key_info": map[string]interface{}{
			"id": apiKey.ID,
			"name": apiKey.Name,
			"prefix": apiKey.Prefix,
			"status": apiKey.Status,
			"scopes": apiKey.Scopes,
			"rate_limit": apiKey.RateLimit,
			"expires_at": apiKey.ExpiresAt,
			"last_used_at": apiKey.LastUsedAt,
			"created_at": apiKey.CreatedAt,
		},
		"message": "API key is valid and active.",
	})
}

// GetAPIKeyScopes returns available API key scopes
// GET /api/v1/api-keys/scopes
func (h *APIKeyHandler) GetAPIKeyScopes(c echo.Context) error {
	scopes := map[string]interface{}{
		"available_scopes": []map[string]interface{}{
			{
				"scope": models.ScopeReadOnly,
				"description": "Read-only access to all endpoints",
				"permissions": []string{"GET", "HEAD"},
			},
			{
				"scope": models.ScopeReadWrite,
				"description": "Read and write access to all endpoints",
				"permissions": []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			},
			{
				"scope": models.ScopeAdmin,
				"description": "Full administrative access",
				"permissions": []string{"ALL"},
			},
			{
				"scope": models.ScopeNutrition,
				"description": "Access to nutrition-related endpoints",
				"endpoints": []string{"/api/v1/nutrition/*"},
			},
			{
				"scope": models.ScopeWorkouts,
				"description": "Access to workout-related endpoints",
				"endpoints": []string{"/api/v1/workouts/*"},
			},
			{
				"scope": models.ScopeMeals,
				"description": "Access to meal-related endpoints",
				"endpoints": []string{"/api/v1/meals/*"},
			},
			{
				"scope": models.ScopeHealth,
				"description": "Access to health-related endpoints",
				"endpoints": []string{"/api/v1/health/*"},
			},
			{
				"scope": models.ScopeSupplements,
				"description": "Access to supplement-related endpoints",
				"endpoints": []string{"/api/v1/supplements/*"},
			},
		},
		"scope_combinations": map[string]interface{}{
			"basic_read": []models.APIKeyScope{models.ScopeReadOnly},
			"basic_write": []models.APIKeyScope{models.ScopeReadWrite},
			"nutrition_specialist": []models.APIKeyScope{models.ScopeNutrition, models.ScopeHealth, models.ScopeSupplements, models.ScopeReadOnly},
			"fitness_specialist": []models.APIKeyScope{models.ScopeWorkouts, models.ScopeMeals, models.ScopeReadOnly},
			"full_access": []models.APIKeyScope{models.ScopeAdmin},
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": scopes,
	})
}

// Helper functions

func validateCreateAPIKeyRequest(req *models.CreateAPIKeyRequest) error {
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "API key name is required")
	}

	if len(req.Name) < 3 || len(req.Name) > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "API key name must be between 3 and 100 characters")
	}

	if len(req.Scopes) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one scope is required")
	}

	if err := models.ValidateScopes(req.Scopes); err != nil {
		return err
	}

	if req.RateLimit < 0 || req.RateLimit > 10000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Rate limit must be between 1 and 10000 requests per minute")
	}

	if req.ExpiresIn != nil && (*req.ExpiresIn < 1 || *req.ExpiresIn > 3650) {
		return echo.NewHTTPError(http.StatusBadRequest, "Expiration must be between 1 and 3650 days")
	}

	// Validate metadata
	if req.Metadata != nil {
		for key, value := range req.Metadata {
			if len(key) > 50 {
				return echo.NewHTTPError(http.StatusBadRequest, "Metadata key length cannot exceed 50 characters")
			}
			if str, ok := value.(string); ok && len(str) > 500 {
				return echo.NewHTTPError(http.StatusBadRequest, "Metadata value length cannot exceed 500 characters")
			}
		}
	}

	return nil
}

// RegisterAPIKeyRoutes registers all API key routes
func RegisterAPIKeyRoutes(e *echo.Echo, apiKeyService *services.APIKeyService) {
	handler := NewAPIKeyHandler(apiKeyService)

	// API key management routes (require authentication)
	apiKeys := e.Group("/api/v1/api-keys")
	
	// Apply authentication middleware to all routes except validation and scopes
	apiKeys.Use(middleware.APIKeyMiddleware(apiKeyService))

	// CRUD operations
	apiKeys.POST("", handler.CreateAPIKey, middleware.ReadWriteMiddleware())
	apiKeys.GET("", handler.GetAPIKeys)
	apiKeys.GET("/:id", handler.GetAPIKey)
	apiKeys.GET("/:id/stats", handler.GetAPIKeyStats)
	apiKeys.DELETE("/:id", handler.RevokeAPIKey, middleware.ReadWriteMiddleware())

	// Public routes (no authentication required)
	public := e.Group("/api/v1/public")
	public.POST("/api-keys/validate", handler.ValidateAPIKey)
	public.GET("/api-keys/scopes", handler.GetAPIKeyScopes)
}