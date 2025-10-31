package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"

	"nutrition-platform/models"
	"nutrition-platform/repositories"
	"nutrition-platform/middleware"
)

type MealHandler struct {
	foodRepo       *repositories.FoodRepository
	recipeRepo     *repositories.RecipeRepository
	foodLogRepo    *repositories.FoodLogRepository
	mealPlanRepo   *repositories.MealPlanRepository
	validator      *validator.Validate
}

func NewMealHandler(
	foodRepo *repositories.FoodRepository,
	recipeRepo *repositories.RecipeRepository,
	foodLogRepo *repositories.FoodLogRepository,
	mealPlanRepo *repositories.MealPlanRepository,
) *MealHandler {
	return &MealHandler{
		foodRepo:     foodRepo,
		recipeRepo:   recipeRepo,
		foodLogRepo:  foodLogRepo,
		mealPlanRepo: mealPlanRepo,
		validator:    validator.New(),
	}
}

// ========== FOOD ENDPOINTS ==========

// CreateFood creates a new food entry
func (h *MealHandler) CreateFood(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req models.FoodRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, middleware.FormatValidationError(err))
	}

	food := models.FoodRequestToFood(&req)
	food.UserID = userID
	food.SourceType = "custom"
	food.Verified = false

	if err := h.foodRepo.CreateFood(food); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create food")
	}

	response := models.FoodToResponse(food)
	return c.JSON(http.StatusCreated, response)
}

// GetFood retrieves a food by ID
func (h *MealHandler) GetFood(c echo.Context) error {
	userID := c.Get("user_id").(string)
	foodID := c.Param("id")

	food, err := h.foodRepo.GetFoodByID(foodID, userID)
	if err != nil {
		if err.Error() == "food not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Food not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food")
	}

	response := models.FoodToResponse(food)
	return c.JSON(http.StatusOK, response)
}

// SearchFoods searches for foods
func (h *MealHandler) SearchFoods(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse query parameters
	query := c.QueryParam("q")
	
	// Parse filters
	filters := models.FoodSearchFilters{
		Brand:        c.QueryParam("brand"),
		SourceType:   c.QueryParam("source_type"),
		SortBy:       c.QueryParam("sort_by"),
		SortDirection: c.QueryParam("sort_direction"),
	}

	// Parse boolean filter
	if verifiedStr := c.QueryParam("verified"); verifiedStr != "" {
		if verified, err := strconv.ParseBool(verifiedStr); err == nil {
			filters.Verified = &verified
		}
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	foods, err := h.foodRepo.SearchFoods(userID, query, filters, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to search foods")
	}

	responses := make([]*models.FoodResponse, len(foods))
	for i, food := range foods {
		responses[i] = models.FoodToResponse(food)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"foods":      responses,
		"page":       page,
		"limit":      limit,
		"total":      len(responses), // TODO: Implement total count query
	})
}

// UpdateFood updates a food entry
func (h *MealHandler) UpdateFood(c echo.Context) error {
	userID := c.Get("user_id").(string)
	foodID := c.Param("id")

	var req models.FoodRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, middleware.FormatValidationError(err))
	}

	// Check if food exists and belongs to user
	existingFood, err := h.foodRepo.GetFoodByID(foodID, userID)
	if err != nil {
		if err.Error() == "food not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Food not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food")
	}

	if existingFood.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Can only update own foods")
	}

	food := models.FoodRequestToFood(&req)
	food.ID = foodID
	food.UserID = userID
	food.SourceType = existingFood.SourceType
	food.Verified = existingFood.Verified

	if err := h.foodRepo.UpdateFood(food); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update food")
	}

	response := models.FoodToResponse(food)
	return c.JSON(http.StatusOK, response)
}

// DeleteFood deletes a food entry
func (h *MealHandler) DeleteFood(c echo.Context) error {
	userID := c.Get("user_id").(string)
	foodID := c.Param("id")

	// Check if food exists and belongs to user
	existingFood, err := h.foodRepo.GetFoodByID(foodID, userID)
	if err != nil {
		if err.Error() == "food not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Food not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food")
	}

	if existingFood.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Can only delete own foods")
	}

	if err := h.foodRepo.DeleteFood(foodID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete food")
	}

	return c.NoContent(http.StatusNoContent)
}

// GetFoodByBarcode retrieves food by barcode
func (h *MealHandler) GetFoodByBarcode(c echo.Context) error {
	userID := c.Get("user_id").(string)
	barcode := c.Param("barcode")

	food, err := h.foodRepo.GetFoodByBarcode(barcode, userID)
	if err != nil {
		if err.Error() == "food not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Food not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food by barcode")
	}

	response := models.FoodToResponse(food)
	return c.JSON(http.StatusOK, response)
}

// ========== RECIPE ENDPOINTS ==========

// CreateRecipe creates a new recipe
func (h *MealHandler) CreateRecipe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req models.RecipeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, middleware.FormatValidationError(err))
	}

	recipe := models.RecipeRequestToRecipe(&req)
	recipe.UserID = userID

	if err := h.recipeRepo.CreateRecipe(recipe); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create recipe")
	}

	response := models.RecipeToResponse(recipe)
	return c.JSON(http.StatusCreated, response)
}

// GetRecipe retrieves a recipe by ID
func (h *MealHandler) GetRecipe(c echo.Context) error {
	userID := c.Get("user_id").(string)
	recipeID := c.Param("id")

	recipe, err := h.recipeRepo.GetRecipeByID(recipeID, userID)
	if err != nil {
		if err.Error() == "recipe not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Recipe not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get recipe")
	}

	response := models.RecipeToResponse(recipe)
	return c.JSON(http.StatusOK, response)
}

// SearchRecipes searches for recipes
func (h *MealHandler) SearchRecipes(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse query parameters
	query := c.QueryParam("q")
	
	// Parse filters
	filters := models.RecipeSearchFilters{
		CuisineType:    c.QueryParam("cuisine_type"),
		Difficulty:     c.QueryParam("difficulty"),
		SortBy:         c.QueryParam("sort_by"),
		SortDirection:  c.QueryParam("sort_direction"),
	}

	// Parse numeric filters
	if maxPrepTimeStr := c.QueryParam("max_prep_time"); maxPrepTimeStr != "" {
		if maxPrepTime, err := strconv.Atoi(maxPrepTimeStr); err == nil {
			filters.MaxPrepTime = &maxPrepTime
		}
	}

	// Parse boolean filter
	if isPublicStr := c.QueryParam("is_public"); isPublicStr != "" {
		if isPublic, err := strconv.ParseBool(isPublicStr); err == nil {
			filters.IsPublic = &isPublic
		}
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	recipes, err := h.recipeRepo.SearchRecipes(userID, query, filters, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to search recipes")
	}

	responses := make([]*models.RecipeResponse, len(recipes))
	for i, recipe := range recipes {
		responses[i] = models.RecipeToResponse(recipe)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recipes":    responses,
		"page":       page,
		"limit":      limit,
		"total":      len(responses), // TODO: Implement total count query
	})
}

// ========== FOOD LOG ENDPOINTS ==========

// CreateFoodLog creates a new food log entry
func (h *MealHandler) CreateFoodLog(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req models.FoodLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, middleware.FormatValidationError(err))
	}

	logEntry := models.FoodLogRequestToFoodLog(&req)
	logEntry.UserID = userID

	if err := h.foodLogRepo.CreateFoodLog(logEntry); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create food log")
	}

	response := models.FoodLogToResponse(logEntry)
	return c.JSON(http.StatusCreated, response)
}

// GetFoodLogs retrieves food logs for a date
func (h *MealHandler) GetFoodLogs(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse date parameter
	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Date parameter is required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
	}

	logs, err := h.foodLogRepo.GetFoodLogsByDate(userID, date)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food logs")
	}

	responses := make([]*models.FoodLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = models.FoodLogToResponse(log)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetNutritionSummary retrieves daily nutrition summary
func (h *MealHandler) GetNutritionSummary(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse date parameter
	dateStr := c.QueryParam("date")
	if dateStr == "" {
		// Default to today
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
	}

	summary, err := h.foodLogRepo.GetDailyNutritionSummary(userID, date)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get nutrition summary")
	}

	return c.JSON(http.StatusOK, summary)
}

// UpdateFoodLog updates a food log entry
func (h *MealHandler) UpdateFoodLog(c echo.Context) error {
	userID := c.Get("user_id").(string)
	logID := c.Param("id")

	var req models.FoodLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, middleware.FormatValidationError(err))
	}

	// Check if log exists and belongs to user
	existingLog, err := h.foodLogRepo.GetFoodLogByID(logID, userID)
	if err != nil {
		if err.Error() == "food log not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Food log not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food log")
	}

	logEntry := models.FoodLogRequestToFoodLog(&req)
	logEntry.ID = logID
	logEntry.UserID = userID

	if err := h.foodLogRepo.UpdateFoodLog(logEntry); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update food log")
	}

	response := models.FoodLogToResponse(logEntry)
	return c.JSON(http.StatusOK, response)
}

// DeleteFoodLog deletes a food log entry
func (h *MealHandler) DeleteFoodLog(c echo.Context) error {
	userID := c.Get("user_id").(string)
	logID := c.Param("id")

	// Check if log exists and belongs to user
	existingLog, err := h.foodLogRepo.GetFoodLogByID(logID, userID)
	if err != nil {
		if err.Error() == "food log not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Food log not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get food log")
	}

	if existingLog.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Can only delete own food logs")
	}

	if err := h.foodLogRepo.DeleteFoodLog(logID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete food log")
	}

	return c.NoContent(http.StatusNoContent)
}

// ========== MEAL PLAN ENDPOINTS ==========

// CreateMealPlan creates a new meal plan
func (h *MealHandler) CreateMealPlan(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req models.MealPlanRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, middleware.FormatValidationError(err))
	}

	mealPlan := models.MealPlanRequestToMealPlan(&req)
	mealPlan.UserID = userID

	if err := h.mealPlanRepo.CreateMealPlan(mealPlan); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create meal plan")
	}

	response := models.MealPlanToResponse(mealPlan)
	return c.JSON(http.StatusCreated, response)
}

// GetMealPlan retrieves a meal plan by ID
func (h *MealHandler) GetMealPlan(c echo.Context) error {
	userID := c.Get("user_id").(string)
	planID := c.Param("id")

	mealPlan, err := h.mealPlanRepo.GetMealPlanByID(planID, userID)
	if err != nil {
		if err.Error() == "meal plan not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Meal plan not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get meal plan")
	}

	response := models.MealPlanToResponse(mealPlan)
	return c.JSON(http.StatusOK, response)
}

// GetActiveMealPlan retrieves the active meal plan
func (h *MealHandler) GetActiveMealPlan(c echo.Context) error {
	userID := c.Get("user_id").(string)

	mealPlan, err := h.mealPlanRepo.GetActiveMealPlan(userID)
	if err != nil {
		if err.Error() == "no active meal plan found" {
			return echo.NewHTTPError(http.StatusNotFound, "No active meal plan found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get active meal plan")
	}

	response := models.MealPlanToResponse(mealPlan)
	return c.JSON(http.StatusOK, response)
}

// SetActiveMealPlan sets a meal plan as active
func (h *MealHandler) SetActiveMealPlan(c echo.Context) error {
	userID := c.Get("user_id").(string)
	planID := c.Param("id")

	if err := h.mealPlanRepo.SetActiveMealPlan(planID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set active meal plan")
	}

	return c.NoContent(http.StatusOK)
}
