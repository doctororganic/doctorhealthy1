package middleware

import (
	"net/http"
	"strings"
	"time"

	"nutrition-platform/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key-change-in-production")

// JWTAuth middleware for JWT authentication
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString == auth {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Bearer token required",
				})
			}

			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			// Check if user still exists
			user, err := services.GetUserByID(claims.UserID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "User not found",
				})
			}

			// Set user context
			c.Set("user", user)
			c.Set("user_id", claims.UserID)
			c.Set("is_admin", claims.IsAdmin)

			return next(c)
		}
	}
}

// AdminAuth middleware for admin-only routes
func AdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAdmin, ok := c.Get("is_admin").(bool)
			if !ok || !isAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Admin access required",
				})
			}
			return next(c)
		}
	}
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID, email, role string, isAdmin bool) (string, error) {
	claims := &Claims{
		UserID:  userID,
		Email:   email,
		Role:    role,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
