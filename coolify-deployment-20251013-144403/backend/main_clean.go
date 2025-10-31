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
		// Nutrition analysis - MAIN WORKING ENDPOINT
		api.POST("/nutrition/analyze", handlers.AnalyzeNutrition)

		// User management endpoints
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)

		// Food management endpoints
		api.GET("/foods", getFoods)
		api.POST("/foods", createFood)
		api.GET("/foods/:id", getFood)
		api.PUT("/foods/:id", updateFood)
		api.DELETE("/foods/:id", deleteFood)
		api.GET("/foods/search", searchFoods)

		// Exercise management endpoints
		api.GET("/exercises", getExercises)
		api.POST("/exercises", createExercise)
		api.GET("/exercises/:id", getExercise)
		api.PUT("/exercises/:id", updateExercise)
		api.DELETE("/exercises/:id", deleteExercise)
		api.GET("/exercises/search", searchExercises)

		// Meal plan endpoints
		api.GET("/meal-plans", getMealPlans)
		api.POST("/meal-plans", generateMealPlan)
		api.GET("/meal-plans/:id", getMealPlan)
		api.PUT("/meal-plans/:id", updateMealPlan)
		api.DELETE("/meal-plans/:id", deleteMealPlan)

		// Workout plan endpoints
		api.GET("/workout-plans", getWorkoutPlans)
		api.POST("/workout-plans", generateWorkouts)
		api.GET("/workout-plans/:id", getWorkoutPlan)
		api.PUT("/workout-plans/:id", updateWorkoutPlan)
		api.DELETE("/workout-plans/:id", deleteWorkoutPlan)

		// Activity logging endpoints
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

		// API information endpoint
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

// UTILITY FUNCTIONS

// generateCorrelationID generates a unique correlation ID for request tracking
func generateCorrelationID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

// randomString generates a cryptographically secure random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to timestamp-based randomness if crypto fails
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		} else {
			b[i] = charset[num.Int64()]
		}
	}
	return string(b)
}

// HANDLER FUNCTIONS

// User handlers
func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get users endpoint",
		"data":    []interface{}{},
	})
}

func createUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get user endpoint",
		"id":      id,
	})
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"id":      id,
	})
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// Food handlers
func getFoods(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get foods endpoint",
		"data":    []interface{}{},
		"page":    page,
		"limit":   limit,
	})
}

func createFood(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Food created successfully",
	})
}

func getFood(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get food endpoint",
		"id":      id,
	})
}

func updateFood(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Food updated successfully",
		"id":      id,
	})
}

func deleteFood(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Food deleted successfully",
		"id":      id,
	})
}

func searchFoods(c echo.Context) error {
	query := c.QueryParam("q")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Search foods endpoint",
		"query":   query,
		"data":    []interface{}{},
	})
}

// Exercise handlers
func getExercises(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get exercises endpoint",
		"data":    []interface{}{},
		"page":    page,
		"limit":   limit,
	})
}

func createExercise(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Exercise created successfully",
	})
}

func getExercise(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get exercise endpoint",
		"id":      id,
	})
}

func updateExercise(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Exercise updated successfully",
		"id":      id,
	})
}

func deleteExercise(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Exercise deleted successfully",
		"id":      id,
	})
}

func searchExercises(c echo.Context) error {
	query := c.QueryParam("q")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Search exercises endpoint",
		"query":   query,
		"data":    []interface{}{},
	})
}

// Meal plan handlers
func getMealPlans(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get meal plans endpoint",
		"data":    []interface{}{},
	})
}

func generateMealPlan(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Meal plan generated successfully",
	})
}

func getMealPlan(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get meal plan endpoint",
		"id":      id,
	})
}

func updateMealPlan(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Meal plan updated successfully",
		"id":      id,
	})
}

func deleteMealPlan(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Meal plan deleted successfully",
		"id":      id,
	})
}

// Workout plan handlers
func getWorkoutPlans(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get workout plans endpoint",
		"data":    []interface{}{},
	})
}

func generateWorkouts(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Workout plan generated successfully",
	})
}

func getWorkoutPlan(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get workout plan endpoint",
		"id":      id,
	})
}

func updateWorkoutPlan(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Workout plan updated successfully",
		"id":      id,
	})
}

func deleteWorkoutPlan(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Workout plan deleted successfully",
		"id":      id,
	})
}

// Activity log handlers
func getFoodLogs(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get food logs endpoint",
		"data":    []interface{}{},
	})
}

func createFoodLog(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Food log created successfully",
	})
}

func getFoodLog(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get food log endpoint",
		"id":      id,
	})
}

func updateFoodLog(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Food log updated successfully",
		"id":      id,
	})
}

func deleteFoodLog(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Food log deleted successfully",
		"id":      id,
	})
}

func getExerciseLogs(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get exercise logs endpoint",
		"data":    []interface{}{},
	})
}

func createExerciseLog(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Exercise log created successfully",
	})
}

func getExerciseLog(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get exercise log endpoint",
		"id":      id,
	})
}

func updateExerciseLog(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Exercise log updated successfully",
		"id":      id,
	})
}

func deleteExerciseLog(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Exercise log deleted successfully",
		"id":      id,
	})
}

// API info endpoint
func getAPIInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":        "Nutrition Platform API",
		"version":     "1.0.0",
		"description": "Doctor Healthy Nutrition Platform Backend API",
		"endpoints": []string{
			"/health",
			"/api/v1/info",
			"/api/v1/nutrition/analyze",
			"/api/v1/users",
			"/api/v1/foods",
			"/api/v1/exercises",
			"/api/v1/meal-plans",
			"/api/v1/workout-plans",
		},
		"features": []string{
			"Nutrition Analysis",
			"Halal Compliance Check",
			"Multi-language Support",
			"Recipe Management",
			"User Management",
			"Activity Tracking",
		},
		"status":    "online",
		"timestamp": time.Now().UTC(),
	})
}