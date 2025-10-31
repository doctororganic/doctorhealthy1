package handlers

import (
	"net/http"
	"strconv"

	"nutrition-platform/middleware"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
)

// CreateSupplement creates a new supplement
func CreateSupplement(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var supplement services.Supplement
	if err := c.Bind(&supplement); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if supplement.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Supplement name is required"})
	}

	supplement.UserID = userID

	err := services.CreateSupplement(&supplement)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create supplement"})
	}

	return c.JSON(http.StatusCreated, supplement)
}

// GetSupplements retrieves supplements for the authenticated user
func GetSupplements(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Check if only active supplements are requested
	activeOnly := c.QueryParam("active") == "true"

	var supplements []services.Supplement
	var err error

	if activeOnly {
		supplements, err = services.GetActiveSupplementsByUserID(userID)
	} else {
		supplements, err = services.GetSupplementsByUserID(userID)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve supplements"})
	}

	return c.JSON(http.StatusOK, supplements)
}

// GetSupplementByID retrieves a specific supplement by ID
func GetSupplementByID(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	supplementID := c.Param("id")
	if supplementID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Supplement ID is required"})
	}

	supplement, err := services.GetSupplementByID(supplementID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplement not found"})
	}

	// Check if user owns this supplement
	if supplement.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	return c.JSON(http.StatusOK, supplement)
}

// UpdateSupplement updates an existing supplement
func UpdateSupplement(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	supplementID := c.Param("id")
	if supplementID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Supplement ID is required"})
	}

	// Check if supplement exists and user owns it
	existingSupplement, err := services.GetSupplementByID(supplementID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplement not found"})
	}

	if existingSupplement.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	var updatedSupplement services.Supplement
	if err := c.Bind(&updatedSupplement); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if updatedSupplement.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Supplement name is required"})
	}

	updatedSupplement.UserID = userID

	err = services.UpdateSupplement(supplementID, &updatedSupplement)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update supplement"})
	}

	return c.JSON(http.StatusOK, updatedSupplement)
}

// DeleteSupplement deletes a supplement
func DeleteSupplement(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	supplementID := c.Param("id")
	if supplementID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Supplement ID is required"})
	}

	err := services.DeleteSupplement(supplementID, userID)
	if err != nil {
		if err.Error() == "supplement not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplement not found"})
		}
		if err.Error() == "unauthorized: supplement belongs to another user" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete supplement"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Supplement deleted successfully"})
}

// ToggleSupplementStatus toggles the active status of a supplement
func ToggleSupplementStatus(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	supplementID := c.Param("id")
	if supplementID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Supplement ID is required"})
	}

	err := services.ToggleSupplementStatus(supplementID, userID)
	if err != nil {
		if err.Error() == "supplement not found or unauthorized" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplement not found or access denied"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to toggle supplement status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Supplement status toggled successfully"})
}

// SearchSupplements searches supplements with filters
func SearchSupplements(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	query := c.QueryParam("q")
	filters := make(map[string]interface{})

	// Parse filters from query parameters
	if category := c.QueryParam("category"); category != "" {
		filters["category"] = category
	}
	if form := c.QueryParam("form"); form != "" {
		filters["form"] = form
	}
	if brand := c.QueryParam("brand"); brand != "" {
		filters["brand"] = brand
	}
	if isHalal := c.QueryParam("is_halal"); isHalal != "" {
		if halal, err := strconv.ParseBool(isHalal); err == nil {
			filters["is_halal"] = halal
		}
	}
	if isVegetarian := c.QueryParam("is_vegetarian"); isVegetarian != "" {
		if vegetarian, err := strconv.ParseBool(isVegetarian); err == nil {
			filters["is_vegetarian"] = vegetarian
		}
	}
	if isVegan := c.QueryParam("is_vegan"); isVegan != "" {
		if vegan, err := strconv.ParseBool(isVegan); err == nil {
			filters["is_vegan"] = vegan
		}
	}
	if isOrganic := c.QueryParam("is_organic"); isOrganic != "" {
		if organic, err := strconv.ParseBool(isOrganic); err == nil {
			filters["is_organic"] = organic
		}
	}
	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filters["is_active"] = active
		}
	}

	supplements, err := services.SearchSupplements(userID, query, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search supplements"})
	}

	return c.JSON(http.StatusOK, supplements)
}

// GetSupplementsByCategory retrieves supplements by category
func GetSupplementsByCategory(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category is required"})
	}

	supplements, err := services.GetSupplementsByCategory(userID, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve supplements"})
	}

	return c.JSON(http.StatusOK, supplements)
}

// GetSupplementStats retrieves supplement statistics for the user
func GetSupplementStats(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	stats, err := services.GetSupplementStats(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve supplement statistics"})
	}

	return c.JSON(http.StatusOK, stats)
}
