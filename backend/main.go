package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"nutrition-platform/api"
	"nutrition-platform/config"
	"nutrition-platform/handlers"
	middlewareCustom "nutrition-platform/middleware/custom"
	"nutrition-platform/models"
	"nutrition-platform/validation"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := models.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDatabase()

	// Create Echo instance
	e := echo.New()

	// Set custom validator
	e.Validator = &validation.InputValidator{}

	// Initialize metrics
	metrics := middlewareCustom.NewMetrics()

	// Set custom error handler
	e.HTTPErrorHandler = middlewareCustom.ErrorHandler(cfg)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middlewareCustom.CorrelationID(cfg))
	e.Use(middlewareCustom.MetricsMiddleware(metrics))
	e.Use(middlewareCustom.SecurityHeaders(cfg))
	e.Use(middlewareCustom.RateLimiter(cfg))
	// Health check endpoint (before circuit breaker to avoid being blocked)
	e.GET("/health", middlewareCustom.HealthCheck())

	// e.Use(middlewareCustom.CircuitBreakerMiddleware(cfg, metrics)) // Temporarily disabled
	e.Use(middlewareCustom.RequestSizeLimit(cfg.MaxRequestSize))

	// Metrics endpoint
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// API routes
	api := e.Group("/api/v1")

	// User routes
	api.GET("/users", getUsers)
	api.POST("/users", createUser)
	api.GET("/users/:id", getUser)
	api.PUT("/users/:id", updateUser)
	api.DELETE("/users/:id", deleteUser)

	// Food routes
	api.GET("/foods", getFoods)
	api.POST("/foods", createFood)
	api.GET("/foods/:id", getFood)
	api.PUT("/foods/:id", updateFood)
	api.DELETE("/foods/:id", deleteFood)
	api.GET("/foods/search", searchFoods)

	// Exercise routes
	api.GET("/exercises", getExercises)
	api.POST("/exercises", createExercise)
	api.GET("/exercises/:id", getExercise)
	api.PUT("/exercises/:id", updateExercise)
	api.DELETE("/exercises/:id", deleteExercise)
	api.GET("/exercises/search", searchExercises)

	// Meal plan routes
	api.GET("/meal-plans", getMealPlans)
	api.POST("/meal-plans", createMealPlan)
	api.GET("/meal-plans/:id", getMealPlan)
	api.PUT("/meal-plans/:id", updateMealPlan)
	api.DELETE("/meal-plans/:id", deleteMealPlan)

	// Workout plan routes
	api.GET("/workout-plans", getWorkoutPlans)
	api.POST("/workout-plans", createWorkoutPlan)
	api.GET("/workout-plans/:id", getWorkoutPlan)
	api.PUT("/workout-plans/:id", updateWorkoutPlan)
	api.DELETE("/workout-plans/:id", deleteWorkoutPlan)

	// Food log routes
	api.GET("/food-logs", getFoodLogs)
	api.POST("/food-logs", createFoodLog)
	api.GET("/food-logs/:id", getFoodLog)
	api.PUT("/food-logs/:id", updateFoodLog)
	api.DELETE("/food-logs/:id", deleteFoodLog)

	// Exercise log routes
	api.GET("/exercise-logs", getExerciseLogs)
	api.POST("/exercise-logs", createExerciseLog)
	api.GET("/exercise-logs/:id", getExerciseLog)
	api.PUT("/exercise-logs/:id", updateExerciseLog)
	api.DELETE("/exercise-logs/:id", deleteExerciseLog)

	// Recipe routes
	recipeHandler := api.NewRecipeHandler("/Users/khaledahmedmohamed/Desktop/trae new healthy1/nutrition-platform/backend/data")
	api.GET("/recipes", recipeHandler.GetRecipes)
	api.GET("/recipes/:id", recipeHandler.GetRecipeByID)
	api.GET("/recipes/country/:country", recipeHandler.GetRecipesByCountry)

	// Initialize and register nutrition data handler
	nutritionHandler := handlers.NewNutritionDataHandler("/Users/khaledahmedmohamed/Desktop/trae new healthy1")
	nutritionHandler.RegisterRoutes(e)

	// Start server
	serverAddr := cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.KeepAliveTimeout) * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Placeholder handler functions

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

func createMealPlan(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Meal plan created successfully",
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

func createWorkoutPlan(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Workout plan created successfully",
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

// Food log handlers
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

// Exercise log handlers
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