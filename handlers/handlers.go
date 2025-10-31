package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"nutrition-platform/cache"
	"nutrition-platform/config"
	"nutrition-platform/database"
	"nutrition-platform/logger"
	"nutrition-platform/models"
	"nutrition-platform/repositories"

	"github.com/labstack/echo/v4"
)

// Handlers struct holds all dependencies
type Handlers struct {
	Config  *config.Config
	DB      *database.Database
	Cache   cache.Cache
	Logger  *logger.Logger
	Audit   *logger.AuditLogger
	PerfLog *logger.PerformanceLogger
}

// NewHandlers creates a new handlers instance
func NewHandlers(cfg *config.Config, db *database.Database, cache cache.Cache, logger *logger.Logger) *Handlers {
	return &Handlers{
		Config:  cfg,
		DB:      db,
		Cache:   cache,
		Logger:  logger,
		Audit:   logger.NewAuditLogger(cfg.Logging),
		PerfLog: logger.NewPerformanceLogger(cfg.Logging),
	}
}

// Response represents a standard API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination and metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
	Count      int   `json:"count,omitempty"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse returns a successful response
func (h *Handlers) SuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithMeta returns a successful response with metadata
func (h *Handlers) SuccessResponseWithMeta(c echo.Context, data interface{}, meta Meta) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

// ErrorResponse returns an error response
func (h *Handlers) ErrorResponse(c echo.Context, status int, code, message string, details ...string) error {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
	
	if len(details) > 0 {
		response.Details = details[0]
	}
	
	return c.JSON(status, response)
}

// ValidationErrorResponse returns a validation error response
func (h *Handlers) ValidationErrorResponse(c echo.Context, errors map[string]string) error {
	return c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   "Validation failed",
		Code:    "VALIDATION_ERROR",
		Data:    errors,
	})
}

// GetUserIDFromContext extracts user ID from echo context
func (h *Handlers) GetUserIDFromContext(c echo.Context) (uint, error) {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return 0, h.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}
	return userID, nil
}

// GetCurrentUser retrieves current user from database
func (h *Handlers) GetCurrentUser(c echo.Context) (*models.User, error) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	repo := repositories.NewUserRepository(h.DB)
	user, err := repo.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, h.ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		h.Logger.Error("Failed to get user", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, h.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get user")
	}

	return user, nil
}

// CheckAdminRole checks if user has admin role
func (h *Handlers) CheckAdminRole(c echo.Context) error {
	userRole, ok := c.Get("user_role").(string)
	if !ok {
		return h.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User role not found")
	}

	if userRole != "admin" {
		h.Audit.LogDataAccess(0, "admin_access", "denied", c.RealIP())
		return h.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Admin access required")
	}

	return nil
}

// ParsePagination parses pagination parameters
func (h *Handlers) ParsePagination(c echo.Context) (page, perPage int) {
	page = 1
	perPage = 20

	if p := c.QueryParam("page"); p != "" {
		if parsed, err := parseInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := c.QueryParam("per_page"); pp != "" {
		if parsed, err := parseInt(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	return page, perPage
}

// Helper function to parse integer
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
