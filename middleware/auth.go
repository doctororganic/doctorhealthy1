package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nutrition-platform/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IssuedAt int64  `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}

// JWTAuth middleware for JWT authentication
func JWTAuth(cfg config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip authentication for public routes
			if isPublicRoute(c.Request().URL.Path) {
				return next(c)
			}

			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Authorization header required",
					"code":  "AUTH_REQUIRED",
				})
			}

			// Check Bearer token format
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid authorization header format",
					"code":  "INVALID_AUTH_FORMAT",
				})
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenParts[1], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.Secret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid token",
					"code":  "INVALID_TOKEN",
					"details": err.Error(),
				})
			}

			// Validate claims
			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				// Check if token is expired
				if time.Now().Unix() > claims.ExpiresAt {
					return c.JSON(http.StatusUnauthorized, map[string]interface{}{
						"error": "Token expired",
						"code":  "TOKEN_EXPIRED",
					})
				}

				// Store user information in context
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
				c.Set("jwt_claims", claims)

				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token claims",
				"code":  "INVALID_CLAIMS",
			})
		}
	}
}

// AdminAuth middleware for admin-only routes
func AdminAuth(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if user has admin role
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "User role not found",
					"code":  "ROLE_NOT_FOUND",
				})
			}

			if userRole != "admin" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Admin access required",
					"code":  "ADMIN_REQUIRED",
				})
			}

			return next(c)
		}
	}
}

// OptionalAuth middleware that doesn't fail if no token is provided
func OptionalAuth(cfg config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c) // No token, continue without authentication
			}

			// Try to authenticate if token is provided
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				token, err := jwt.ParseWithClaims(tokenParts[1], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(cfg.Secret), nil
				})

				if err == nil {
					if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
						if time.Now().Unix() <= claims.ExpiresAt {
							c.Set("user_id", claims.UserID)
							c.Set("user_email", claims.Email)
							c.Set("user_role", claims.Role)
							c.Set("jwt_claims", claims)
						}
					}
				}
			}

			return next(c)
		}
	}
}

// GenerateTokens generates access and refresh tokens
func GenerateTokens(userID uint, email, role string, cfg config.JWTConfig) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessClaims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(cfg.AccessExpiry).Unix(),
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(cfg.RefreshExpiry).Unix(),
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(cfg.RefreshSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns claims
func ValidateToken(tokenString string, cfg config.JWTConfig) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// RefreshAccessToken generates a new access token from a refresh token
func RefreshAccessToken(refreshTokenString string, cfg config.JWTConfig) (string, error) {
	// Validate refresh token using refresh secret
	claims, err := func() (*JWTClaims, error) {
		token, err := jwt.ParseWithClaims(refreshTokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.RefreshSecret), nil
		})

		if err != nil {
			return nil, err
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			return claims, nil
		}

		return nil, errors.New("invalid token claims")
	}()

	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token is expired
	if time.Now().Unix() > claims.ExpiresAt {
		return "", errors.New("refresh token expired")
	}

	// Generate new access token
	accessClaims := &JWTClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(cfg.AccessExpiry).Unix(),
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return accessTokenObj.SignedString([]byte(cfg.Secret))
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string, rounds int) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), rounds)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// isPublicRoute checks if a route should be accessible without authentication
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/health",
		"/api/v1/info",
		"/api/v1/public/",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/forgot-password",
		"/api/v1/auth/reset-password",
		"/api/v1/auth/verify-email",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}

	return false
}
