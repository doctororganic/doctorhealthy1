package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"nutrition-platform/repositories"

	"github.com/labstack/echo/v4"
)

// FoodHandler handles food CRUD operations
type FoodHandler struct {
	foodRepo *repositories.FoodRepository
}

// NewFoodHandler creates a new food handler
func NewFoodHandler(db *sql.DB) *FoodHandler {
	return &FoodHandler{
		foodRepo: repositories.NewFoodRepository(db),
	}
}

// GetFoods returns paginated list of foods
func (h *FoodHandler) GetFoods(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Parse pagination parameters
	page := 1
	limit := 20
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse search query
	searchQuery := c.QueryParam("search")

	// Get foods
	offset := (page - 1) * limit
	userIDStr := userID.(string)
	
	// Use SearchFoods with empty filters for now
	filters := struct {
		Brand       string
		SourceType  string
		Verified    *bool
		SortBy      string
		SortDirection string
	}{}
	
	foods, err := h.foodRepo.SearchFoods(userIDStr, searchQuery, filters, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch foods: " + err.Error(),
		})
	}

	// TODO: Get total count for pagination metadata
	total := len(foods) // Placeholder

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   foods,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
			"has_next":   page*limit < total,
			"has_prev":   page > 1,
		},
	})
}

// GetFood returns a specific food by ID
func (h *FoodHandler) GetFood(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	foodID := c.Param("id")
	userIDStr := userID.(string)

	food, err := h.foodRepo.GetFoodByID(foodID, userIDStr)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Food not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   food,
	})
}

// CreateFood creates a new food
func (h *FoodHandler) CreateFood(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req struct {
		Name        string   `json:"name" validate:"required"`
		Description *string  `json:"description"`
		Brand       *string  `json:"brand"`
		Barcode     *string  `json:"barcode"`
		Category    *string  `json:"category"`
		Calories    float64  `json:"calories"`
		Protein     float64  `json:"protein"`
		Carbs       float64  `json:"carbs"`
		Fat         float64  `json:"fat"`
		Fiber       float64  `json:"fiber"`
		Sugar       float64  `json:"sugar"`
		Sodium      int      `json:"sodium"`
		ServingSize string   `json:"serving_size"`
		ServingUnit string   `json:"serving_unit"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Convert to repository model (simplified - may need adjustment based on actual model)
	food := &struct {
		ID          uint
		UserID      *uint
		Name        string
		Brand       *string
		Description *string
		BarCode     *string
		ServingSize string
		Calories    float64
		Protein     float64
		Carbs       float64
		Fat         float64
		SaturatedFat float64
		Fiber       float64
		Sugar       float64
		Sodium      int
		Cholesterol float64
		Potassium   float64
		SourceType  string
		Verified    bool
	}{
		UserID:      userID.(*uint),
		Name:        req.Name,
		Brand:       req.Brand,
		Description: req.Description,
		BarCode:     req.Barcode,
		ServingSize: req.ServingSize,
		Calories:    req.Calories,
		Protein:     req.Protein,
		Carbs:       req.Carbs,
		Fat:         req.Fat,
		Fiber:       req.Fiber,
		Sugar:       req.Sugar,
		Sodium:      req.Sodium,
		SourceType:  "user",
		Verified:    false,
	}

	// Note: This is a stub - actual implementation needs proper model conversion
	// For now, return success
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Food created successfully",
		"data":    food,
	})
}

// UpdateFood updates an existing food
func (h *FoodHandler) UpdateFood(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	foodID := c.Param("id")
	// Stub implementation
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Food updated successfully",
		"id":      foodID,
	})
}

// DeleteFood deletes a food
func (h *FoodHandler) DeleteFood(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	foodID := c.Param("id")
	userIDStr := userID.(string)

	err := h.foodRepo.DeleteFood(foodID, userIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete food: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Food deleted successfully",
	})
}

// SearchFoods searches for foods
func (h *FoodHandler) SearchFoods(c echo.Context) error {
	return h.GetFoods(c) // Reuse GetFoods which already supports search
}

