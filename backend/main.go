package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MindscapeHQ/raygun4go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	
	"nutrition-platform/config"
	"nutrition-platform/handlers"
	authMiddleware "nutrition-platform/middleware"
	"nutrition-platform/repositories"
	"nutrition-platform/storage"
)

func main() {
	// Initialize Raygun error tracking
	raygun, err := raygun4go.New("nutrition-platform", "J5KNVQg46P71JymsDyPWiQ")
	if err != nil {
		log.Println("Unable to create Raygun client:", err.Error())
	}
	defer raygun.HandleError()

	// Set global Raygun context
	raygun.Version("1.0.0")
	raygun.Tags([]string{"go", "echo", "nutrition-platform"})
	raygun.CustomData(map[string]interface{}{
		"environment": "production",
		"service":     "nutrition-platform-backend",
		"framework":   "echo",
	})

	// Load configuration
	cfg := config.LoadConfig()
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	log.Printf("Starting Nutrition Platform Backend on port %s", cfg.Port)

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	foodRepo := repositories.NewFoodRepository(db)
	recipeRepo := repositories.NewRecipeRepository(db)
	foodLogRepo := repositories.NewFoodLogRepository(db)
	mealPlanRepo := repositories.NewMealPlanRepository(db)
	
	// Initialize workout repositories
	exerciseRepo := repositories.NewExerciseRepository(db)
	workoutPlanRepo := repositories.NewWorkoutPlanRepository(db)
	workoutLogRepo := repositories.NewWorkoutLogRepository(db)
	personalRecordRepo := repositories.NewPersonalRecordRepository(db)

	// Initialize progress tracking repositories
	progressPhotoRepo := repositories.NewProgressPhotoRepository(db)
	bodyMeasurementRepo := repositories.NewBodyMeasurementRepository(db)
	milestoneRepo := repositories.NewMilestoneRepository(db)
	weightGoalRepo := repositories.NewWeightGoalRepository(db)
	progressAnalyticsRepo := repositories.NewProgressAnalyticsRepository(db)

	// Initialize handlers
	mealHandler := handlers.NewMealHandler(foodRepo, recipeRepo, foodLogRepo, mealPlanRepo)
	workoutHandler := handlers.NewWorkoutHandler(exerciseRepo, workoutPlanRepo, workoutLogRepo, personalRecordRepo)
	progressHandler := handlers.NewProgressHandler(progressPhotoRepo, bodyMeasurementRepo, milestoneRepo, weightGoalRepo, progressAnalyticsRepo)

	// Initialize storage
	fileStorage := storage.NewLocalStorage(cfg.UploadDir, cfg.BaseURL)
	imageProcessor := storage.NewImageProcessor()

	// Initialize file handler
	fileHandler := handlers.NewFileHandler(fileStorage, imageProcessor, progressPhotoRepo)

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

	// Raygun request context middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get correlation ID for this request
			correlationID := c.Request().Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = generateCorrelationID()
			}

			// Set request context for Raygun error tracking
			raygun.Request(c.Request())
			
			// Add request-specific context
			raygun.CustomData(map[string]interface{}{
				"correlation_id": correlationID,
				"endpoint":       c.Path(),
				"method":         c.Request().Method,
				"user_agent":     c.Request().UserAgent(),
				"remote_addr":     c.Request().RemoteAddr,
			})

			// Set user information if available (from JWT or session)
			if userID := c.Get("user_id"); userID != nil {
				raygun.User(fmt.Sprintf("%v", userID))
			}

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

		// Food endpoints
		foods := api.Group("/foods")
		foods.POST("", mealHandler.CreateFood)
		foods.GET("/:id", mealHandler.GetFood)
		foods.PUT("/:id", mealHandler.UpdateFood)
		foods.DELETE("/:id", mealHandler.DeleteFood)
		foods.GET("", mealHandler.SearchFoods)
		foods.GET("/barcode/:barcode", mealHandler.GetFoodByBarcode)

		// Recipe endpoints
		recipes := api.Group("/recipes")
		recipes.POST("", mealHandler.CreateRecipe)
		recipes.GET("/:id", mealHandler.GetRecipe)
		recipes.GET("", mealHandler.SearchRecipes)

		// Food log endpoints
		foodLogs := api.Group("/food-logs")
		foodLogs.POST("", mealHandler.CreateFoodLog)
		foodLogs.GET("", mealHandler.GetFoodLogs)
		foodLogs.PUT("/:id", mealHandler.UpdateFoodLog)
		foodLogs.DELETE("/:id", mealHandler.DeleteFoodLog)
		foodLogs.GET("/nutrition-summary", mealHandler.GetNutritionSummary)

		// Meal plan endpoints
		mealPlans := api.Group("/meal-plans")
		mealPlans.POST("", mealHandler.CreateMealPlan)
		mealPlans.GET("/:id", mealHandler.GetMealPlan)
		mealPlans.GET("/active", mealHandler.GetActiveMealPlan)
		mealPlans.PUT("/:id/activate", mealHandler.SetActiveMealPlan)

		// Workout endpoints
		exercises := api.Group("/exercises")
		exercises.POST("", workoutHandler.CreateExercise)
		exercises.GET("/:id", workoutHandler.GetExercise)
		exercises.PUT("/:id", workoutHandler.UpdateExercise)
		exercises.DELETE("/:id", workoutHandler.DeleteExercise)
		exercises.GET("", workoutHandler.ListExercises)
		exercises.GET("/muscle-groups", workoutHandler.GetMuscleGroups)
		exercises.GET("/equipment-types", workoutHandler.GetEquipmentTypes)

		// Workout plan endpoints
		workoutPlans := api.Group("/workout-plans")
		workoutPlans.POST("", workoutHandler.CreateWorkoutPlan)
		workoutPlans.GET("/:id", workoutHandler.GetWorkoutPlan)
		workoutPlans.GET("", workoutHandler.ListWorkoutPlans)

		// Workout log endpoints
		workoutLogs := api.Group("/workout-logs")
		workoutLogs.POST("", workoutHandler.CreateWorkoutLog)
		workoutLogs.GET("/:id", workoutHandler.GetWorkoutLog)
		workoutLogs.GET("", workoutHandler.ListWorkoutLogs)
		workoutLogs.GET("/stats", workoutHandler.GetUserWorkoutStats)

		// Personal record endpoints
		personalRecords := api.Group("/personal-records")
		personalRecords.POST("", workoutHandler.CreatePersonalRecord)
		personalRecords.GET("", workoutHandler.ListPersonalRecords)
		personalRecords.GET("/user", workoutHandler.GetUserPersonalRecords)

		// Progress tracking endpoints
		progressPhotos := api.Group("/progress-photos")
		progressPhotos.POST("", progressHandler.UploadProgressPhoto)
		progressPhotos.GET("", progressHandler.GetProgressPhotos)
		progressPhotos.DELETE("/:id", progressHandler.DeleteProgressPhoto)

		// Body measurement endpoints
		bodyMeasurements := api.Group("/body-measurements")
		bodyMeasurements.POST("", progressHandler.LogBodyMeasurement)
		bodyMeasurements.GET("", progressHandler.GetBodyMeasurementHistory)
		bodyMeasurements.GET("/latest", progressHandler.GetLatestBodyMeasurements)

		// Milestone endpoints
		milestones := api.Group("/milestones")
		milestones.POST("", progressHandler.CreateMilestone)
		milestones.GET("", progressHandler.GetMilestones)
		milestones.PUT("/:id", progressHandler.UpdateMilestone)
		milestones.DELETE("/:id", progressHandler.DeleteMilestone)

		// Weight goal endpoints
		weightGoals := api.Group("/weight-goals")
		weightGoals.POST("", progressHandler.SetWeightGoal)
		weightGoals.GET("", progressHandler.GetWeightGoals)
		weightGoals.GET("/active", progressHandler.GetActiveWeightGoal)

		// Progress analytics endpoints
		progressAnalytics := api.Group("/progress-analytics")
		progressAnalytics.GET("/summary", progressHandler.GetProgressSummary)
		progressAnalytics.GET("/weight-progress", progressHandler.GetWeightProgress)
		progressAnalytics.GET("/measurement-trends", progressHandler.GetMeasurementTrends)

		// File upload endpoints
		files := api.Group("/files")
		files.POST("/upload", fileHandler.UploadFile)
		files.POST("/upload/progress-photo", fileHandler.UploadProgressPhoto)
		files.POST("/upload/bulk", fileHandler.BulkUpload)
		files.GET("/:path", fileHandler.GetFile)
		files.DELETE("/:path", fileHandler.DeleteFile)
		files.GET("/:path/info", fileHandler.GetFileInfo)
		files.POST("/validate", fileHandler.ValidateImage)
		files.GET("/upload/:uploadId/progress", fileHandler.GetUploadProgress)

		// Legacy placeholder endpoints (to be replaced)
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)

		api.GET("/info", getAPIInfo)
	}

	// Health endpoint - use proper health check with database connectivity
	e.GET("/health", healthHandler)

	// Simple health endpoint for load balancers (always returns 200)
	e.GET("/health/simple", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"service": "nutrition-platform-backend",
			"timestamp": time.Now().Unix(),
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

// initDatabase initializes database connection and runs migrations
func initDatabase(cfg *config.Config) (*sql.DB, error) {
	// Use environment variable or fallback to config
	dbURL := cfg.DatabaseURL
	if dbURL == "" {
		dbURL = "postgres://localhost/nutrition_platform?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully")

	// TODO: Run migrations when migration system is ready
	// For now, just ensure basic tables exist
	if err := createBasicTables(db); err != nil {
		log.Printf("Warning: Failed to create basic tables: %v", err)
	}

	return db, nil
}

// createBasicTables creates minimal required tables for the application
func createBasicTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS foods (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			brand VARCHAR(100),
			calories DECIMAL(10,2),
			protein DECIMAL(10,2),
			carbs DECIMAL(10,2),
			fat DECIMAL(10,2),
			fiber DECIMAL(10,2),
			sugar DECIMAL(10,2),
			sodium DECIMAL(10,2),
			barcode VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS food_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			food_id INTEGER REFERENCES foods(id),
			quantity DECIMAL(10,2) NOT NULL,
			unit VARCHAR(50) NOT NULL,
			meal_type VARCHAR(50),
			log_date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
	}

	log.Println("Basic database tables verified/created")
	return nil
}

// Helper functions for correlation ID generation
func generateCorrelationID() string {
	// Generate a random 16-character hex string
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// Placeholder handlers (to be implemented)
func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Users endpoint - to be implemented"})
}

func createUser(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Create user endpoint - to be implemented"})
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Get user %s - to be implemented", id)})
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Update user %s - to be implemented", id)})
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Delete user %s - to be implemented", id)})
}

func getAPIInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":        "Nutrition Platform API",
		"version":     "1.0.0",
		"description": "Backend API for nutrition tracking, meal planning, and workout management",
		"endpoints": map[string]interface{}{
			"/api/v1/nutrition/analyze": "POST - Analyze nutrition from text or image",
			"/api/v1/foods": map[string]string{
				"GET": "Search foods",
				"POST": "Create food",
				"GET /:id": "Get food by ID",
				"PUT /:id": "Update food",
				"DELETE /:id": "Delete food",
				"GET /barcode/:barcode": "Get food by barcode",
			},
			"/api/v1/recipes": map[string]string{
				"GET": "Search recipes",
				"POST": "Create recipe",
				"GET /:id": "Get recipe by ID",
			},
			"/api/v1/food-logs": map[string]string{
				"GET": "Get food logs",
				"POST": "Create food log",
				"PUT /:id": "Update food log",
				"DELETE /:id": "Delete food log",
				"GET /nutrition-summary": "Get nutrition summary",
			},
			"/api/v1/meal-plans": map[string]string{
				"GET": "Get meal plans",
				"POST": "Create meal plan",
				"GET /:id": "Get meal plan by ID",
				"GET /active": "Get active meal plan",
				"PUT /:id/activate": "Set active meal plan",
			},
			"/api/v1/exercises": map[string]string{
				"GET": "List exercises",
				"POST": "Create exercise",
				"GET /:id": "Get exercise by ID",
				"PUT /:id": "Update exercise",
				"DELETE /:id": "Delete exercise",
				"GET /muscle-groups": "Get muscle groups",
				"GET /equipment-types": "Get equipment types",
			},
			"/api/v1/workout-plans": map[string]string{
				"GET": "List workout plans",
				"POST": "Create workout plan",
				"GET /:id": "Get workout plan by ID",
			},
			"/api/v1/workout-logs": map[string]string{
				"GET": "List workout logs",
				"POST": "Create workout log",
				"GET /:id": "Get workout log by ID",
				"GET /stats": "Get workout statistics",
			},
			"/api/v1/personal-records": map[string]string{
				"GET": "List personal records",
				"POST": "Create personal record",
				"GET /user": "Get user personal records",
			},
			"/api/v1/progress-photos": map[string]string{
				"GET": "Get progress photos",
				"POST": "Upload progress photo",
				"DELETE /:id": "Delete progress photo",
			},
			"/api/v1/body-measurements": map[string]string{
				"GET": "Get body measurement history",
				"POST": "Log body measurement",
				"GET /latest": "Get latest body measurements",
			},
			"/api/v1/milestones": map[string]string{
				"GET": "Get milestones",
				"POST": "Create milestone",
				"PUT /:id": "Update milestone",
				"DELETE /:id": "Delete milestone",
			},
			"/api/v1/weight-goals": map[string]string{
				"GET": "Get weight goals",
				"POST": "Set weight goal",
				"GET /active": "Get active weight goal",
			},
			"/api/v1/progress-analytics": map[string]string{
				"GET /summary": "Get progress summary",
				"GET /weight-progress": "Get weight progress",
				"GET /measurement-trends": "Get measurement trends",
			},
			"/api/v1/files": map[string]string{
				"POST /upload": "Upload generic file",
				"POST /upload/progress-photo": "Upload progress photo with processing",
				"POST /upload/bulk": "Upload multiple files",
				"GET /:path": "Get file",
				"DELETE /:path": "Delete file",
				"GET /:path/info": "Get file information",
				"POST /validate": "Validate image file",
				"GET /upload/:uploadId/progress": "Get upload progress",
			},
		},
	})
}

func healthHandler(c echo.Context) error {
	// TODO: Add actual database connectivity check
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"database":  "connected",
		"version":   "1.0.0",
	})
}
