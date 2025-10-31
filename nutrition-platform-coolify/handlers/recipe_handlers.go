package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"nutrition-platform/models"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
)

// RecipeHandler handles recipe-related endpoints
type RecipeHandler struct {
	recipeService *services.RecipeService
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

// CreateRecipe creates a new recipe
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	var req models.CreateRecipeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "RECIPE_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "RECIPE_002",
		})
	}

	// Get user ID from context
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "RECIPE_003",
		})
	}

	recipe, err := h.recipeService.CreateRecipe(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "creation_failed",
			"message": err.Error(),
			"code":    "RECIPE_004",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"recipe":  recipe,
		"message": "Recipe created successfully",
	})
}

// GetRecipe retrieves a recipe by ID
func (h *RecipeHandler) GetRecipe(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_recipe_id",
			"message": "Recipe ID is required",
			"code":    "RECIPE_005",
		})
	}

	recipe, err := h.recipeService.GetRecipe(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "recipe_not_found",
			"message": "Recipe not found",
			"code":    "RECIPE_006",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"recipe":  recipe,
	})
}

// SearchRecipes searches for recipes
func (h *RecipeHandler) SearchRecipes(c echo.Context) error {
	req := &models.RecipeSearchRequest{
		Query:           c.QueryParam("q"),
		Cuisine:         c.QueryParam("cuisine"),
		Country:         c.QueryParam("country"),
		DifficultyLevel: c.QueryParam("difficulty"),
		Page:            1,
		Limit:           20,
	}

	// Parse optional parameters
	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil && page > 0 {
		req.Page = page
	}

	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil && limit > 0 && limit <= 100 {
		req.Limit = limit
	}

	if maxPrepTime, err := strconv.Atoi(c.QueryParam("max_prep_time")); err == nil && maxPrepTime > 0 {
		req.MaxPrepTime = maxPrepTime
	}

	if maxCookTime, err := strconv.Atoi(c.QueryParam("max_cook_time")); err == nil && maxCookTime > 0 {
		req.MaxCookTime = maxCookTime
	}

	if minRating, err := strconv.ParseFloat(c.QueryParam("min_rating"), 64); err == nil && minRating > 0 {
		req.MinRating = minRating
	}

	// Parse boolean parameters
	if isHalal := c.QueryParam("is_halal"); isHalal != "" {
		if halal, err := strconv.ParseBool(isHalal); err == nil {
			req.IsHalal = &halal
		}
	}

	if isKosher := c.QueryParam("is_kosher"); isKosher != "" {
		if kosher, err := strconv.ParseBool(isKosher); err == nil {
			req.IsKosher = &kosher
		}
	}

	// Parse array parameters
	if dietaryTags := c.QueryParam("dietary_tags"); dietaryTags != "" {
		req.DietaryTags = parseCommaSeparated(dietaryTags)
	}

	if allergens := c.QueryParam("exclude_allergens"); allergens != "" {
		req.Allergens = parseCommaSeparated(allergens)
	}

	response, err := h.recipeService.SearchRecipes(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "search_failed",
			"message": err.Error(),
			"code":    "RECIPE_007",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// UpdateRecipe updates a recipe
func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_recipe_id",
			"message": "Recipe ID is required",
			"code":    "RECIPE_005",
		})
	}

	var req models.UpdateRecipeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "RECIPE_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "RECIPE_002",
		})
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "RECIPE_003",
		})
	}

	recipe, err := h.recipeService.UpdateRecipe(id, userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "update_failed",
			"message": err.Error(),
			"code":    "RECIPE_008",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"recipe":  recipe,
		"message": "Recipe updated successfully",
	})
}

// DeleteRecipe deletes a recipe
func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_recipe_id",
			"message": "Recipe ID is required",
			"code":    "RECIPE_005",
		})
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "RECIPE_003",
		})
	}

	err := h.recipeService.DeleteRecipe(id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "deletion_failed",
			"message": err.Error(),
			"code":    "RECIPE_009",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Recipe deleted successfully",
	})
}

// RateRecipe rates a recipe
func (h *RecipeHandler) RateRecipe(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_recipe_id",
			"message": "Recipe ID is required",
			"code":    "RECIPE_005",
		})
	}

	var req models.RecipeRatingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "RECIPE_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "RECIPE_002",
		})
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "RECIPE_003",
		})
	}

	err := h.recipeService.RateRecipe(id, userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "rating_failed",
			"message": err.Error(),
			"code":    "RECIPE_010",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Recipe rated successfully",
	})
}

// GetRecipesByCountry retrieves recipes by country
func (h *RecipeHandler) GetRecipesByCountry(c echo.Context) error {
	country := c.Param("country")
	if country == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_country",
			"message": "Country parameter is required",
			"code":    "RECIPE_011",
		})
	}

	page := 1
	limit := 20

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	response, err := h.recipeService.GetRecipesByCountry(country, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "retrieval_failed",
			"message": err.Error(),
			"code":    "RECIPE_012",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// GetRecipesByCuisine retrieves recipes by cuisine
func (h *RecipeHandler) GetRecipesByCuisine(c echo.Context) error {
	cuisine := c.Param("cuisine")
	if cuisine == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_cuisine",
			"message": "Cuisine parameter is required",
			"code":    "RECIPE_013",
		})
	}

	page := 1
	limit := 20

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	response, err := h.recipeService.GetRecipesByCuisine(cuisine, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "retrieval_failed",
			"message": err.Error(),
			"code":    "RECIPE_012",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// RegisterRoutes registers recipe routes
func (h *RecipeHandler) RegisterRoutes(e *echo.Group) {
	recipes := e.Group("/recipes")

	// CRUD operations
	recipes.POST("", h.CreateRecipe)
	recipes.GET("", h.SearchRecipes)
	recipes.GET("/:id", h.GetRecipe)
	recipes.PUT("/:id", h.UpdateRecipe)
	recipes.DELETE("/:id", h.DeleteRecipe)

	// Additional operations
	recipes.POST("/:id/rate", h.RateRecipe)
	recipes.GET("/country/:country", h.GetRecipesByCountry)
	recipes.GET("/cuisine/:cuisine", h.GetRecipesByCuisine)
}

// Helper functions
func getUserIDFromContext(c echo.Context) string {
	if userID, ok := c.Get("user_id").(string); ok {
		return userID
	}
	return ""
}

func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	// Simple comma separation - in production, you might want more sophisticated parsing
	var result []string
	for _, item := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
