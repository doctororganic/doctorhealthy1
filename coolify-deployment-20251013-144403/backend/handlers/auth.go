package handlers

import (
	"net/http"
	"strings"

	"nutrition-platform/middleware"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	FirstName       string `json:"first_name" validate:"required,min=2"`
	LastName        string `json:"last_name" validate:"required,min=2"`
	DateOfBirth     string `json:"date_of_birth" validate:"required"`
	Gender          string `json:"gender" validate:"required,oneof=male female"`
	Language        string `json:"language" validate:"required,oneof=en ar"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
	ExpiresIn    int64       `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Register handles user registration
func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	// Check if passwords match
	if req.Password != req.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Passwords do not match",
		})
	}

	// Check if user already exists
	existingUser, _ := services.GetUserByEmail(req.Email)
	if existingUser != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "User with this email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process password",
		})
	}

	// Create user
	user, err := services.CreateUser(
		req.Email,
		string(hashedPassword),
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	// Generate tokens
	accessToken, err := middleware.GenerateToken(user.ID, user.Email, user.Role, user.IsAdmin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate access token",
		})
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate refresh token",
		})
	}

	// Create session
	deviceInfo := c.Request().Header.Get("X-Device-Info")
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	_, err = services.CreateSession(user.ID, accessToken, refreshToken, deviceInfo, ipAddress, userAgent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	// Remove password from user object before sending response
	user.Password = ""

	return c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	})
}

// Login handles user authentication
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	// Get user by email
	user, err := services.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid email or password",
		})
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid email or password",
		})
	}

	// Generate tokens
	accessToken, err := middleware.GenerateToken(user.ID, user.Email, user.Role, user.IsAdmin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate access token",
		})
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate refresh token",
		})
	}

	// Create session
	deviceInfo := c.Request().Header.Get("X-Device-Info")
	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	_, err = services.CreateSession(user.ID, accessToken, refreshToken, deviceInfo, ipAddress, userAgent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	// Remove password from user object before sending response
	user.Password = ""

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	})
}

// Logout handles user logout
func Logout(c echo.Context) error {
	// Get token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization header required",
		})
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid authorization header format",
		})
	}

	accessToken := tokenParts[1]

	// Invalidate session
	err := services.InvalidateSessionByToken(accessToken)
	if err != nil {
		// Log error but still return success for security
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// RefreshToken handles token refresh
func RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Get session by refresh token
	session, err := services.GetSessionByRefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Session not found or expired",
		})
	}

	// Get user
	user, err := services.GetUserByID(session.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not found",
		})
	}

	// Generate new tokens
	newAccessToken, err := middleware.GenerateToken(user.ID, user.Email, user.Role, user.IsAdmin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate access token",
		})
	}

	newRefreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate refresh token",
		})
	}

	// Update session with new tokens
	err = services.UpdateSessionTokens(session.ID, newAccessToken, newRefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update session",
		})
	}

	// Remove password from user object before sending response
	user.Password = ""

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	})
}

// LogoutAll handles logout from all devices
func LogoutAll(c echo.Context) error {
	// Get user from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Invalidate all user sessions
	err := services.InvalidateAllUserSessions(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to logout from all devices",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out from all devices successfully",
	})
}

// GetProfile returns the current user's profile
func GetProfile(c echo.Context) error {
	// Get user from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Get user
	user, err := services.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// Remove password from response
	user.Password = ""

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// GetSessions returns the current user's active sessions
func GetSessions(c echo.Context) error {
	// Get user from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Get user sessions
	sessions, err := services.GetUserSessions(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve sessions",
		})
	}

	// Remove sensitive information from sessions
	for i := range sessions {
		sessions[i].AccessToken = ""
		sessions[i].RefreshToken = ""
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sessions": sessions,
	})
}

// InvalidateSessionHandler handles invalidating a specific session
func InvalidateSessionHandler(c echo.Context) error {
	sessionID := c.Param("session_id")

	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Session ID required",
		})
	}

	// Invalidate session
	err := services.InvalidateSession(sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to invalidate session",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Session invalidated successfully",
	})
}
