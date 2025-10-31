#!/bin/bash

# ============================================
# REFACTOR MONOLITHIC MAIN.GO
# Break into modular structure
# ============================================

echo "🔧 Refactoring monolithic main.go..."

# Create modular structure
mkdir -p backend/cmd/server
mkdir -p backend/internal/router
mkdir -p backend/internal/handlers/{users,foods,exercises,meals,workouts,nutrition}
mkdir -p backend/internal/services
mkdir -p backend/internal/middleware
mkdir -p backend/config

echo "✅ Created modular directory structure"

# Create new main.go (minimal)
cat > backend/cmd/server/main.go << 'EOF'
package main

import (
    "log"
    "nutrition-platform/config"
    "nutrition-platform/internal/router"
    "nutrition-platform/pkg/database"
    "nutrition-platform/pkg/logger"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize logger
    if err := logger.Init(cfg.LogLevel); err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Close()
    
    // Initialize database
    if err := database.Init(cfg.Database); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.Close()
    
    // Setup router
    e := router.Setup(cfg)
    
    // Start server
    log.Printf("Starting server on port %s", cfg.Port)
    if err := e.Start(":" + cfg.Port); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
EOF

echo "✅ Created minimal main.go"

# Create router
cat > backend/internal/router/router.go << 'EOF'
package router

import (
    "nutrition-platform/config"
    "nutrition-platform/internal/handlers/users"
    "nutrition-platform/internal/handlers/foods"
    "nutrition-platform/internal/handlers/nutrition"
    "nutrition-platform/internal/middleware"
    
    "github.com/labstack/echo/v4"
    echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Setup(cfg *config.Config) *echo.Echo {
    e := echo.New()
    
    // Global middleware
    e.Use(echomiddleware.Logger())
    e.Use(echomiddleware.Recover())
    e.Use(echomiddleware.CORS())
    e.Use(middleware.Security())
    e.Use(middleware.RateLimit())
    
    // Health check
    e.GET("/health", healthHandler)
    
    // API routes
    api := e.Group("/api/v1")
    
    // User routes
    users.RegisterRoutes(api.Group("/users"))
    
    // Food routes
    foods.RegisterRoutes(api.Group("/foods"))
    
    // Nutrition routes
    nutrition.RegisterRoutes(api.Group("/nutrition"))
    
    return e
}

func healthHandler(c echo.Context) error {
    return c.JSON(200, map[string]string{"status": "healthy"})
}
EOF

echo "✅ Created modular router"

# Create user handlers
cat > backend/internal/handlers/users/users.go << 'EOF'
package users

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group) {
    g.GET("", GetUsers)
    g.POST("", CreateUser)
    g.GET("/:id", GetUser)
    g.PUT("/:id", UpdateUser)
    g.DELETE("/:id", DeleteUser)
}

func GetUsers(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]interface{}{
        "users": []interface{}{},
    })
}

func CreateUser(c echo.Context) error {
    return c.JSON(http.StatusCreated, map[string]string{
        "message": "User created",
    })
}

func GetUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{
        "id": id,
    })
}

func UpdateUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{
        "message": "User updated",
        "id": id,
    })
}

func DeleteUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{
        "message": "User deleted",
        "id": id,
    })
}
EOF

echo "✅ Created user handlers"

echo ""
echo "============================"
echo "✅ Refactoring complete!"
echo "============================"
echo ""
echo "New structure:"
echo "  backend/cmd/server/main.go (minimal entry point)"
echo "  backend/internal/router/ (routing logic)"
echo "  backend/internal/handlers/ (handler functions)"
echo "  backend/internal/services/ (business logic)"
echo "  backend/internal/middleware/ (middleware)"
echo "  backend/config/ (configuration)"
echo ""
echo "Next steps:"
echo "1. Move handlers from old main.go to new structure"
echo "2. Move business logic to services/"
echo "3. Test the refactored code"
echo "4. Remove old main.go"
echo ""
