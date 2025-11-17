// Package main is the entry point for the nutrition platform backend
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nutrition-platform/config"
	"nutrition-platform/handlers"
	backendmodels "nutrition-platform/models"
	"nutrition-platform/services"
	"nutrition-platform/security"
	"nutrition-platform/validation"

	customMiddleware "nutrition-platform/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := backendmodels.InitDB(cfg.GetDatabaseURL())
	defer func() {
		if err := backendmodels.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize services
	healthService := services.NewHealthService(db)
	nutritionPlanService := services.NewNutritionPlanService(db)

	// Initialize Echo instance
	e := echo.New()

	// Set validator
	e.Validator = validation.NewInputValidator()

	// Middleware
	e.Use(customMiddleware.CustomLogger())
	e.Use(customMiddleware.CustomRecover())
	e.Use(customMiddleware.CORS())
	e.Use(customMiddleware.RateLimiter())

	// Custom middleware
	e.Use(customMiddleware.AuthMiddleware())
	e.Use(customMiddleware.SecurityHeaders())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(healthService)
	nutritionPlanHandler := handlers.NewNutritionPlanHandler(nutritionPlanService, healthService)
	nutritionDataHandler := handlers.NewNutritionDataHandler("./data")

	// Initialize JWT manager and auth handler
	jwtManager := security.NewJWTManager()
	authHandler := handlers.NewAuthHandler(nil, jwtManager) // UserService is nil for stub implementation

	// Routes
	api := e.Group("/api/v1")

	// Authentication routes (no auth middleware required)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout-all", handlers.LogoutAll)

	// Protected auth routes (require JWT authentication)
	protectedAuth := api.Group("/auth")
	protectedAuth.Use(customMiddleware.JWTAuth())
	protectedAuth.POST("/logout", authHandler.Logout)
	protectedAuth.GET("/profile", authHandler.GetProfile)
	protectedAuth.GET("/sessions", authHandler.GetSessions)
	protectedAuth.DELETE("/sessions/:id", authHandler.DeleteSession)
	protectedAuth.PUT("/profile", authHandler.UpdateProfile)
	protectedAuth.DELETE("/profile", authHandler.DeleteProfile)
	protectedAuth.POST("/change-password", authHandler.ChangePassword)

	// Admin auth routes (require JWT authentication)
	adminAuth := api.Group("/auth/admin")
	adminAuth.Use(customMiddleware.JWTAuth())
	adminAuth.Use(customMiddleware.AdminAuth())
	adminAuth.GET("/users", authHandler.GetAllUsers)
	adminAuth.DELETE("/users/:id", authHandler.DeleteUser)
	adminAuth.GET("/audit-logs", authHandler.GetAuditLogs)

	// Protected routes (require JWT authentication)
	protected := api.Group("")
	protected.Use(customMiddleware.JWTAuth())
	protected.GET("/dashboard", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to protected dashboard",
			"user_id": c.Get("user_id"),
		})
	})

	// Health routes
	health := api.Group("/health")
	health.POST("/complaints", healthHandler.CreateHealthComplaint)
	health.GET("/complaints", healthHandler.GetUserHealthComplaints)
	health.POST("/injuries", healthHandler.CreateUserInjury)
	health.GET("/injuries", healthHandler.GetUserInjuries)
	health.GET("/conditions", healthHandler.GetHealthConditions)
	health.POST("/assessment", healthHandler.PerformHealthAssessment)
	health.POST("/risk-assessment", healthHandler.GetHealthRiskAssessment)
	health.GET("/symptom-checker", healthHandler.GetSymptomChecker)
	health.GET("/tips", healthHandler.GetHealthTips)

	// Nutrition plan routes
	nutrition := api.Group("/nutrition-plans")
	nutrition.POST("/recommendations", nutritionPlanHandler.GetNutritionPlanRecommendations)
	nutrition.GET("/quick-assessment", nutritionPlanHandler.GetQuickNutritionAssessment)
	nutrition.POST("/comparison", nutritionPlanHandler.GetNutritionPlanComparison)
	nutrition.GET("/types", nutritionPlanHandler.GetNutritionPlanTypes)
	nutrition.GET("/types/:plan_type", nutritionPlanHandler.GetNutritionPlanDetails)
	nutrition.POST("/personalized", nutritionPlanHandler.CreatePersonalizedNutritionPlan)

	// Nutrition data routes
	api.GET("/metabolism", nutritionDataHandler.GetMetabolismGuide)
	api.GET("/workout-techniques", nutritionDataHandler.GetWorkoutTechniques)
	api.GET("/meal-plans", nutritionDataHandler.GetMealPlans)
	api.POST("/meal-plans/generate", nutritionDataHandler.GenerateMealPlan)
	api.GET("/vitamins-minerals", nutritionDataHandler.GetVitaminsAndMinerals)
	api.GET("/drugs-nutrition", nutritionDataHandler.GetDrugsNutritionInteractions)
	api.GET("/diseases", nutritionDataHandler.GetDiseaseData)
	api.GET("/calories", handlers.GetCalories)
	api.GET("/calories/:category", handlers.GetCaloriesByCategory)
	api.GET("/skills", handlers.GetSkills)
	api.GET("/skills/:difficulty", handlers.GetSkillsByDifficulty)
	api.GET("/type-plans", handlers.GetTypePlans)
	api.GET("/type-plans/:type", handlers.GetTypePlansByType)

	// Protected generation routes (require JWT authentication)
	protectedGen := api.Group("/generate")
	protectedGen.Use(customMiddleware.JWTAuth())
	protectedGen.GET("/meal-plans", nutritionDataHandler.GenerateMealPlan)
	protectedGen.GET("/workout-plans", nutritionDataHandler.GenerateWorkoutPlan)
	protectedGen.GET("/supplements", nutritionDataHandler.GetSupplementsRecommendations)
	protectedGen.GET("/drugs/interactions", nutritionDataHandler.GetDrugInteractions)

	// API routes
	api.GET("/meals", handlers.GetMealsAPI)
	api.POST("/meals", handlers.CreateMealAPI)
	api.GET("/meals/:id", handlers.GetMealAPI)
	api.PUT("/meals/:id", handlers.UpdateMealAPI)
	api.DELETE("/meals/:id", handlers.DeleteMealAPI)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "nutrition-platform-backend",
			"version":   "1.0.0",
		})
	})

	// Info endpoint
	e.GET("/api/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"service":     "Nutrition Platform API",
			"version":     "1.0.0",
			"environment": cfg.Environment,
			"endpoints": map[string]interface{}{
				"health":          "/health",
				"nutrition_data":  "/api/v1/metabolism, /api/v1/meal-plans, etc.",
				"nutrition_plans": "/api/v1/nutrition-plans/*",
				"health_services": "/api/v1/health/*",
				"api_endpoints":   "/api/v1/meals/*",
			},
		})
	})

	// Start server
	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on port %s", cfg.Port)

	// Start server in a goroutine
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
