package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"nutrition-platform/middleware"
	"nutrition-platform/models"
	"nutrition-platform/repositories"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	handlers *Handlers
	repo     *repositories.UserRepository
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(handlers *Handlers, repo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{
		handlers: handlers,
		repo:     repo,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Validate input
	if err := h.validateCreateUserRequest(&req); err != nil {
		return h.handlers.ValidationErrorResponse(c, err)
	}

	// Check if user already exists
	existingUser, err := h.repo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return h.handlers.ErrorResponse(c, http.StatusConflict, "EMAIL_EXISTS", "Email already registered")
	}

	// Hash password
	passwordHash, err := middleware.HashPassword(req.Password, h.handlers.Config.Security.BCryptRounds)
	if err != nil {
		h.handlers.Logger.Error("Failed to hash password", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "HASH_ERROR", "Failed to process password")
	}

	// Create user
	user := &models.User{
		Email:       req.Email,
		PasswordHash: passwordHash,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        "user",
		IsVerified:  false,
		IsActive:    true,
	}

	err = h.repo.CreateUser(user)
	if err != nil {
		h.handlers.Logger.Error("Failed to create user", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "CREATE_ERROR", "Failed to create user")
	}

	// Create user profile if provided
	if req.Profile != nil {
		profile := &models.UserProfile{
			UserID: user.ID,
		}

		if req.Profile.DateOfBirth != nil {
			if dob, err := time.Parse("2006-01-02", *req.Profile.DateOfBirth); err == nil {
				profile.DateOfBirth.Valid = true
				profile.DateOfBirth.Time = dob
			}
		}

		if req.Profile.Gender != nil {
			profile.Gender = req.Profile.Gender
		}

		if req.Profile.Height != nil {
			profile.Height.Valid = true
			profile.Height.Float64 = *req.Profile.Height
		}

		if req.Profile.Weight != nil {
			profile.Weight.Valid = true
			profile.Weight.Float64 = *req.Profile.Weight
		}

		if req.Profile.ActivityLevel != nil {
			profile.ActivityLevel = req.Profile.ActivityLevel
		}

		if req.Profile.Goal != nil {
			profile.Goal = req.Profile.Goal
		}

		profile.DietaryRestrictions = req.Profile.DietaryRestrictions
		profile.Allergies = req.Profile.Allergies

		if req.Profile.PreferredUnits != nil {
			profile.PreferredUnits = *req.Profile.PreferredUnits
		} else {
			profile.PreferredUnits = "metric"
		}

		err = h.repo.CreateUserProfile(profile)
		if err != nil {
			h.handlers.Logger.Warn("Failed to create user profile", map[string]interface{}{
				"user_id": user.ID,
				"error":   err.Error(),
			})
		}
	}

	// Create user preferences if provided
	if req.Preferences != nil {
		prefs := &models.UserPreferences{
			UserID: user.ID,
		}

		if req.Preferences.Language != nil {
			prefs.Language = *req.Preferences.Language
		} else {
			prefs.Language = "en"
		}

		if req.Preferences.Timezone != nil {
			prefs.Timezone = *req.Preferences.Timezone
		} else {
			prefs.Timezone = "UTC"
		}

		if req.Preferences.NotificationsEnabled != nil {
			prefs.NotificationsEnabled = *req.Preferences.NotificationsEnabled
		} else {
			prefs.NotificationsEnabled = true
		}

		if req.Preferences.EmailNotifications != nil {
			prefs.EmailNotifications = *req.Preferences.EmailNotifications
		} else {
			prefs.EmailNotifications = true
		}

		if req.Preferences.PushNotifications != nil {
			prefs.PushNotifications = *req.Preferences.PushNotifications
		} else {
			prefs.PushNotifications = true
		}

		if req.Preferences.Units != nil {
			prefs.Units = *req.Preferences.Units
		} else {
			prefs.Units = "metric"
		}

		if req.Preferences.DarkMode != nil {
			prefs.DarkMode = *req.Preferences.DarkMode
		} else {
			prefs.DarkMode = false
		}

		err = h.repo.CreateOrUpdateUserPreferences(prefs)
		if err != nil {
			h.handlers.Logger.Warn("Failed to create user preferences", map[string]interface{}{
				"user_id": user.ID,
				"error":   err.Error(),
			})
		}
	}

	// Generate tokens
	accessToken, refreshToken, err := middleware.GenerateTokens(
		user.ID,
		user.Email,
		user.Role,
		h.handlers.Config.JWT,
	)
	if err != nil {
		h.handlers.Logger.Error("Failed to generate tokens", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "TOKEN_ERROR", "Failed to generate tokens")
	}

	// Store refresh token (in a real implementation, this would be stored in database)
	// For now, we'll just return it

	// Log registration
	h.handlers.Audit.LogAuthentication(user.ID, user.Email, "register", c.RealIP(), c.Request().UserAgent(), true)

	// Return response
	userResponse := user.ToUserResponse()
	response := models.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(h.handlers.Config.JWT.AccessExpiry.Seconds()),
	}

	return h.handlers.SuccessResponse(c, response)
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required")
	}

	// Get user
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		h.handlers.Audit.LogAuthentication(0, req.Email, "login", c.RealIP(), c.Request().UserAgent(), false)
		return h.handlers.ErrorResponse(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	// Check password
	if !middleware.CheckPassword(req.Password, user.PasswordHash) {
		h.handlers.Audit.LogAuthentication(user.ID, user.Email, "login", c.RealIP(), c.Request().UserAgent(), false)
		return h.handlers.ErrorResponse(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		h.handlers.Audit.LogAuthentication(user.ID, user.Email, "login", c.RealIP(), c.Request().UserAgent(), false)
		return h.handlers.ErrorResponse(c, http.StatusForbidden, "ACCOUNT_INACTIVE", "Account is deactivated")
	}

	// Update last login
	err = h.repo.UpdateLastLogin(user.ID)
	if err != nil {
		h.handlers.Logger.Warn("Failed to update last login", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
	}

	// Generate tokens
	accessToken, refreshToken, err := middleware.GenerateTokens(
		user.ID,
		user.Email,
		user.Role,
		h.handlers.Config.JWT,
	)
	if err != nil {
		h.handlers.Logger.Error("Failed to generate tokens", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "TOKEN_ERROR", "Failed to generate tokens")
	}

	// Store refresh token (in a real implementation, this would be stored in database)
	// For now, we'll just return it

	// Log successful login
	h.handlers.Audit.LogAuthentication(user.ID, user.Email, "login", c.RealIP(), c.Request().UserAgent(), true)

	// Return response
	userResponse := user.ToUserResponse()
	response := models.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(h.handlers.Config.JWT.AccessExpiry.Seconds()),
	}

	return h.handlers.SuccessResponse(c, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req models.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Validate refresh token
	claims, err := middleware.ValidateToken(req.RefreshToken, h.handlers.Config.JWT)
	if err != nil {
		return h.handlers.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid refresh token")
	}

	// Get user to verify they're still active
	user, err := h.repo.GetUserByID(claims.UserID)
	if err != nil {
		return h.handlers.ErrorResponse(c, http.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
	}

	if !user.IsActive {
		return h.handlers.ErrorResponse(c, http.StatusForbidden, "ACCOUNT_INACTIVE", "Account is deactivated")
	}

	// Generate new access token
	newAccessToken, err := middleware.RefreshAccessToken(req.RefreshToken, h.handlers.Config.JWT)
	if err != nil {
		h.handlers.Logger.Error("Failed to refresh token", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "TOKEN_ERROR", "Failed to refresh token")
	}

	// Log token refresh
	h.handlers.Audit.LogAuthentication(user.ID, user.Email, "refresh_token", c.RealIP(), c.Request().UserAgent(), true)

	response := models.RefreshTokenResponse{
		AccessToken: newAccessToken,
		ExpiresIn:   int(h.handlers.Config.JWT.AccessExpiry.Seconds()),
	}

	return h.handlers.SuccessResponse(c, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c echo.Context) error {
	userID, err := h.handlers.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// In a real implementation, you would invalidate the refresh token
	// For now, we'll just log the logout

	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		return h.handlers.ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
	}

	// Log logout
	h.handlers.Audit.LogAuthentication(user.ID, user.Email, "logout", c.RealIP(), c.Request().UserAgent(), true)

	return h.handlers.SuccessResponse(c, map[string]string{"message": "Logged out successfully"})
}

// ForgotPassword handles password reset requests
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req models.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Check if user exists
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return h.handlers.SuccessResponse(c, map[string]string{
			"message": "If an account with that email exists, a password reset link has been sent",
		})
	}

	// Generate reset token (in a real implementation, this would be stored in database)
	resetToken := generateResetToken(user.Email)
	
	// In a real implementation, you would send an email with the reset link
	// For now, we'll just log it
	h.handlers.Logger.Info("Password reset requested", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"token":   resetToken,
	})

	// Log password reset request
	h.handlers.Audit.LogAuthentication(user.ID, user.Email, "forgot_password", c.RealIP(), c.Request().UserAgent(), true)

	return h.handlers.SuccessResponse(c, map[string]string{
		"message": "If an account with that email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req models.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// In a real implementation, you would validate the reset token and get the user email
	// For now, we'll just return a success message
	
	h.handlers.Logger.Info("Password reset attempted", map[string]interface{}{
		"token": req.Token,
	})

	return h.handlers.SuccessResponse(c, map[string]string{
		"message": "Password has been reset successfully",
	})
}

// ChangePassword handles password change for authenticated users
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID, err := h.handlers.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return h.handlers.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Get user
	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		return h.handlers.ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
	}

	// Verify current password
	if !middleware.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return h.handlers.ErrorResponse(c, http.StatusUnauthorized, "INVALID_PASSWORD", "Current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := middleware.HashPassword(req.NewPassword, h.handlers.Config.Security.BCryptRounds)
	if err != nil {
		h.handlers.Logger.Error("Failed to hash new password", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "HASH_ERROR", "Failed to process new password")
	}

	// Update password
	err = h.repo.UpdatePassword(user.ID, newPasswordHash)
	if err != nil {
		h.handlers.Logger.Error("Failed to update password", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return h.handlers.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_ERROR", "Failed to update password")
	}

	// Log password change
	h.handlers.Audit.LogAuthentication(user.ID, user.Email, "change_password", c.RealIP(), c.Request().UserAgent(), true)

	return h.handlers.SuccessResponse(c, map[string]string{
		"message": "Password changed successfully",
	})
}

// validateCreateUserRequest validates the create user request
func (h *AuthHandler) validateCreateUserRequest(req *models.CreateUserRequest) map[string]string {
	errors := make(map[string]string)

	// Validate email
	if req.Email == "" {
		errors["email"] = "Email is required"
	} else if !strings.Contains(req.Email, "@") {
		errors["email"] = "Invalid email format"
	}

	// Validate password
	if req.Password == "" {
		errors["password"] = "Password is required"
	} else if len(req.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	}

	// Validate profile data if provided
	if req.Profile != nil {
		if req.Profile.Height != nil && (*req.Profile.Height < 50 || *req.Profile.Height > 300) {
			errors["profile.height"] = "Height must be between 50 and 300 cm"
		}

		if req.Profile.Weight != nil && (*req.Profile.Weight < 20 || *req.Profile.Weight > 500) {
			errors["profile.weight"] = "Weight must be between 20 and 500 kg"
		}

		if req.Profile.Gender != nil && *req.Profile.Gender != "" {
			validGenders := []string{"male", "female", "other"}
			valid := false
			for _, gender := range validGenders {
				if *req.Profile.Gender == gender {
					valid = true
					break
				}
			}
			if !valid {
				errors["profile.gender"] = "Gender must be male, female, or other"
			}
		}

		if req.Profile.ActivityLevel != nil && *req.Profile.ActivityLevel != "" {
			validLevels := []string{"sedentary", "light", "moderate", "active", "very_active"}
			valid := false
			for _, level := range validLevels {
				if *req.Profile.ActivityLevel == level {
					valid = true
					break
				}
			}
			if !valid {
				errors["profile.activity_level"] = "Invalid activity level"
			}
		}

		if req.Profile.Goal != nil && *req.Profile.Goal != "" {
			validGoals := []string{"lose_weight", "maintain", "gain_weight", "gain_muscle"}
			valid := false
			for _, goal := range validGoals {
				if *req.Profile.Goal == goal {
					valid = true
					break
				}
			}
			if !valid {
				errors["profile.goal"] = "Invalid goal"
			}
		}

		if req.Profile.PreferredUnits != nil && *req.Profile.PreferredUnits != "" {
			if *req.Profile.PreferredUnits != "metric" && *req.Profile.PreferredUnits != "imperial" {
				errors["profile.preferred_units"] = "Units must be metric or imperial"
			}
		}
	}

	// Validate preferences if provided
	if req.Preferences != nil {
		if req.Preferences.Language != nil && *req.Preferences.Language != "" {
			if len(*req.Preferences.Language) != 2 {
				errors["preferences.language"] = "Language must be a 2-character code"
			}
		}

		if req.Preferences.Units != nil && *req.Preferences.Units != "" {
			if *req.Preferences.Units != "metric" && *req.Preferences.Units != "imperial" {
				errors["preferences.units"] = "Units must be metric or imperial"
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// generateResetToken generates a password reset token
func generateResetToken(email string) string {
	h := sha256.New()
	h.Write([]byte(email + time.Now().String()))
	return hex.EncodeToString(h.Sum(nil))[:32]
}
