package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"nutrition-platform/config"
	"nutrition-platform/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	log.Printf("Starting Nutrition Platform Backend on port %s", cfg.Port)

	// Create Echo instance
	e := echo.New()

	// Global middleware
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())

	// Correlation ID middleware for request tracking
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for existing correlation ID in header, otherwise generate one
			correlationID := c.Request().Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = generateCorrelationID()
			}

			// Set correlation ID in response header
			c.Response().Header().Set("X-Correlation-ID", correlationID)

			// Store in context for use in logging
			c.Set("correlation_id", correlationID)

			return next(c)
		}
	})

	// CORS configuration for frontend integration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost",
		},
		AllowMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization, "X-API-Key", "X-Requested-With",
			"X-Correlation-ID",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// API routes
	api := e.Group("/api/v1")
	{
		// Nutrition analysis - working endpoint
		api.POST("/nutrition/analyze", handlers.AnalyzeNutrition)

		// Placeholder endpoints for other features
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)

		api.GET("/foods", getFoods)
		api.POST("/foods", createFood)
		api.GET("/foods/:id", getFood)
		api.PUT("/foods/:id", updateFood)
		api.DELETE("/foods/:id", deleteFood)
		api.GET("/foods/search", searchFoods)

		api.GET("/exercises", getExercises)
		api.POST("/exercises", createExercise)
		api.GET("/exercises/:id", getExercise)
		api.PUT("/exercises/:id", updateExercise)
		api.DELETE("/exercises/:id", deleteExercise)
		api.GET("/exercises/search", searchExercises)

		api.GET("/meal-plans", getMealPlans)
		api.POST("/meal-plans", generateMealPlan)
		api.GET("/meal-plans/:id", getMealPlan)
		api.PUT("/meal-plans/:id", updateMealPlan)
		api.DELETE("/meal-plans/:id", deleteMealPlan)

		api.GET("/workout-plans", getWorkoutPlans)
		api.POST("/workout-plans", generateWorkouts)
		api.GET("/workout-plans/:id", getWorkoutPlan)
		api.PUT("/workout-plans/:id", updateWorkoutPlan)
		api.DELETE("/workout-plans/:id", deleteWorkoutPlan)

		api.GET("/food-logs", getFoodLogs)
		api.POST("/food-logs", createFoodLog)
		api.GET("/food-logs/:id", getFoodLog)
		api.PUT("/food-logs/:id", updateFoodLog)
		api.DELETE("/food-logs/:id", deleteFoodLog)

		api.GET("/exercise-logs", getExerciseLogs)
		api.POST("/exercise-logs", createExerciseLog)
		api.GET("/exercise-logs/:id", getExerciseLog)
		api.PUT("/exercise-logs/:id", updateExerciseLog)
		api.DELETE("/exercise-logs/:id", deleteExerciseLog)

		api.GET("/info", getAPIInfo)
	}

	// Health endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now(),
			"version":   "1.0.0",
		})
	})

	// Start server
	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: e,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.KeepAliveTimeout) * time.Second,
	}

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}


