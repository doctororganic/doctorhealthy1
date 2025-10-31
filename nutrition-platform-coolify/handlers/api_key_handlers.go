package handlers

import (
	"net/http"
	"strconv"

	"nutrition-platform/models"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
)

// APIKeyHandler handles API key management endpoints
type APIKeyHandler struct {
	apiKeyService    *services.APIKeyService
	analyticsService *services.AnalyticsService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(apiKeyService *services.APIKeyService, analyticsService *services.AnalyticsService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService:    apiKeyService,
		analyticsService: analyticsService,
	}
}

// CreateAPIKey creates a new API key
func (h *APIKeyHandler) CreateAPIKey(c echo.Context) error {
	var req models.CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "API_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "API_002",
		})
	}

	// Get user ID from context (set by authentication middleware)
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "API_003",
		})
	}

	response, err := h.apiKeyService.CreateAPIKey(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "creation_failed",
			"message": err.Error(),
			"code":    "API_004",
		})
	}

	return c.JSON(http.StatusCreated, response)
}

// GetAPIKeys retrieves API keys for the authenticated user
func (h *APIKeyHandler) GetAPIKeys(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "API_003",
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	response, err := h.apiKeyService.GetAPIKeys(userID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "retrieval_failed",
			"message": err.Error(),
			"code":    "API_005",
		})
	}

	return c.JSON(http.StatusOK, response)
}

// RevokeAPIKey revokes an API key
func (h *APIKeyHandler) RevokeAPIKey(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "API_003",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_api_key_id",
			"message": "API key ID is required",
			"code":    "API_006",
		})
	}

	err := h.apiKeyService.RevokeAPIKey(userID, apiKeyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "revocation_failed",
			"message": err.Error(),
			"code":    "API_007",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "API key revoked successfully",
		"api_key_id": apiKeyID,
	})
}

// GetAPIKeyStats retrieves usage statistics for an API key
func (h *APIKeyHandler) GetAPIKeyStats(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "API_003",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_api_key_id",
			"message": "API key ID is required",
			"code":    "API_006",
		})
	}

	days, _ := strconv.Atoi(c.QueryParam("days"))
	if days < 1 || days > 365 {
		days = 30 // Default to 30 days
	}

	stats, err := h.apiKeyService.GetAPIKeyStats(userID, apiKeyID, days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "stats_retrieval_failed",
			"message": err.Error(),
			"code":    "API_008",
		})
	}

	return c.JSON(http.StatusOK, stats)
}

// GetAPIKeyMetrics retrieves real-time metrics for an API key
func (h *APIKeyHandler) GetAPIKeyMetrics(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "API_003",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_api_key_id",
			"message": "API key ID is required",
			"code":    "API_006",
		})
	}

	metrics, err := h.analyticsService.GetAPIMetrics(apiKeyID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "metrics_not_found",
			"message": err.Error(),
			"code":    "API_009",
		})
	}

	return c.JSON(http.StatusOK, metrics)
}

// GetUsageReport generates a comprehensive usage report
func (h *APIKeyHandler) GetUsageReport(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "API_003",
		})
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_api_key_id",
			"message": "API key ID is required",
			"code":    "API_006",
		})
	}

	days, _ := strconv.Atoi(c.QueryParam("days"))
	if days < 1 || days > 365 {
		days = 30 // Default to 30 days
	}

	report, err := h.analyticsService.GetUsageReport(apiKeyID, days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "report_generation_failed",
			"message": err.Error(),
			"code":    "API_010",
		})
	}

	return c.JSON(http.StatusOK, report)
}

// GetGlobalMetrics retrieves global API metrics (admin only)
func (h *APIKeyHandler) GetGlobalMetrics(c echo.Context) error {
	// Check if user has admin privileges
	apiKey := c.Get("api_key").(*models.APIKey)
	if !apiKey.HasScope(models.ScopeAdmin) {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error":   "insufficient_permissions",
			"message": "Admin privileges required",
			"code":    "API_011",
		})
	}

	metrics, err := h.analyticsService.GetGlobalMetrics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "global_metrics_failed",
			"message": err.Error(),
			"code":    "API_012",
		})
	}

	return c.JSON(http.StatusOK, metrics)
}

// RegisterRoutes registers API key management routes
func (h *APIKeyHandler) RegisterRoutes(e *echo.Group) {
	apiKeys := e.Group("/api-keys")

	// API key CRUD operations
	apiKeys.POST("", h.CreateAPIKey)
	apiKeys.GET("", h.GetAPIKeys)
	apiKeys.DELETE("/:id", h.RevokeAPIKey)

	// API key analytics and monitoring
	apiKeys.GET("/:id/stats", h.GetAPIKeyStats)
	apiKeys.GET("/:id/metrics", h.GetAPIKeyMetrics)
	apiKeys.GET("/:id/report", h.GetUsageReport)

	// Global metrics (admin only)
	apiKeys.GET("/global/metrics", h.GetGlobalMetrics)
}
