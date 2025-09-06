package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"nutrition-platform/middleware"
	"nutrition-platform/services"
)

// CreateMedicalPlan creates a new medical plan
func CreateMedicalPlan(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var plan services.MedicalPlan
	if err := c.Bind(&plan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if plan.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan name is required"})
	}
	if plan.Type == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan type is required"})
	}

	plan.UserID = userID

	err := services.CreateMedicalPlan(&plan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create medical plan"})
	}

	return c.JSON(http.StatusCreated, plan)
}

// GetMedicalPlans retrieves medical plans for the authenticated user
func GetMedicalPlans(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Check if only active plans are requested
	activeOnly := c.QueryParam("active") == "true"
	// Check if public plans should be included
	includePublic := c.QueryParam("include_public") == "true"

	var plans []services.MedicalPlan
	var err error

	if includePublic {
		// Get public plans
		publicPlans, err := services.GetPublicMedicalPlans()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve public plans"})
		}
		plans = append(plans, publicPlans...)
	}

	// Get user's own plans
	var userPlans []services.MedicalPlan
	if activeOnly {
		userPlans, err = services.GetActiveMedicalPlansByUserID(userID)
	} else {
		userPlans, err = services.GetMedicalPlansByUserID(userID)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve medical plans"})
	}

	plans = append(plans, userPlans...)

	return c.JSON(http.StatusOK, plans)
}

// GetPublicMedicalPlans retrieves public medical plans
func GetPublicMedicalPlans(c echo.Context) error {
	plans, err := services.GetPublicMedicalPlans()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve public medical plans"})
	}

	return c.JSON(http.StatusOK, plans)
}

// GetMedicalPlanByID retrieves a specific medical plan by ID
func GetMedicalPlanByID(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	planID := c.Param("id")
	if planID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan ID is required"})
	}

	plan, err := services.GetMedicalPlanByID(planID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Medical plan not found"})
	}

	// Check if user can access this plan (owns it or it's public)
	if plan.UserID != userID && !plan.IsPublic {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	return c.JSON(http.StatusOK, plan)
}

// UpdateMedicalPlan updates an existing medical plan
func UpdateMedicalPlan(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	planID := c.Param("id")
	if planID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan ID is required"})
	}

	// Check if plan exists and user owns it
	existingPlan, err := services.GetMedicalPlanByID(planID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Medical plan not found"})
	}

	if existingPlan.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	var updatedPlan services.MedicalPlan
	if err := c.Bind(&updatedPlan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if updatedPlan.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan name is required"})
	}
	if updatedPlan.Type == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan type is required"})
	}

	updatedPlan.UserID = userID

	err = services.UpdateMedicalPlan(planID, &updatedPlan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update medical plan"})
	}

	return c.JSON(http.StatusOK, updatedPlan)
}

// DeleteMedicalPlan deletes a medical plan
func DeleteMedicalPlan(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	planID := c.Param("id")
	if planID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan ID is required"})
	}

	err := services.DeleteMedicalPlan(planID, userID)
	if err != nil {
		if err.Error() == "medical plan not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Medical plan not found"})
		}
		if err.Error() == "unauthorized: plan belongs to another user" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete medical plan"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Medical plan deleted successfully"})
}

// SearchMedicalPlans searches medical plans with filters
func SearchMedicalPlans(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	query := c.QueryParam("q")
	includePublic := c.QueryParam("include_public") == "true"
	filters := make(map[string]interface{})

	// Parse filters from query parameters
	if planType := c.QueryParam("type"); planType != "" {
		filters["type"] = planType
	}
	if category := c.QueryParam("category"); category != "" {
		filters["category"] = category
	}
	if createdBy := c.QueryParam("created_by"); createdBy != "" {
		filters["created_by"] = createdBy
	}
	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filters["is_active"] = active
		}
	}
	if isPublic := c.QueryParam("is_public"); isPublic != "" {
		if public, err := strconv.ParseBool(isPublic); err == nil {
			filters["is_public"] = public
		}
	}
	if minRating := c.QueryParam("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			filters["min_rating"] = rating
		}
	}
	if maxDuration := c.QueryParam("max_duration"); maxDuration != "" {
		if duration, err := strconv.Atoi(maxDuration); err == nil {
			filters["max_duration"] = duration
		}
	}
	if minDuration := c.QueryParam("min_duration"); minDuration != "" {
		if duration, err := strconv.Atoi(minDuration); err == nil {
			filters["min_duration"] = duration
		}
	}

	plans, err := services.SearchMedicalPlans(userID, query, filters, includePublic)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search medical plans"})
	}

	return c.JSON(http.StatusOK, plans)
}

// GetMedicalPlansByCategory retrieves medical plans by category
func GetMedicalPlansByCategory(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category is required"})
	}

	includePublic := c.QueryParam("include_public") == "true"

	plans, err := services.GetMedicalPlansByCategory(category, includePublic, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve medical plans"})
	}

	return c.JSON(http.StatusOK, plans)
}

// RateMedicalPlan adds a rating to a medical plan
func RateMedicalPlan(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	planID := c.Param("id")
	if planID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan ID is required"})
	}

	var ratingRequest struct {
		Rating float64 `json:"rating"`
	}

	if err := c.Bind(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if ratingRequest.Rating < 1.0 || ratingRequest.Rating > 5.0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be between 1.0 and 5.0"})
	}

	// Check if plan exists and is accessible
	plan, err := services.GetMedicalPlanByID(planID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Medical plan not found"})
	}

	// Users can only rate public plans or their own plans
	if plan.UserID != userID && !plan.IsPublic {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	err = services.RateMedicalPlan(planID, ratingRequest.Rating)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to rate medical plan"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Rating added successfully"})
}

// UpdateCheckpoint updates a checkpoint in a medical plan
func UpdateCheckpoint(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	planID := c.Param("id")
	checkpointID := c.Param("checkpoint_id")

	if planID == "" || checkpointID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plan ID and Checkpoint ID are required"})
	}

	var checkpointUpdate struct {
		Completed bool                   `json:"completed"`
		Notes     string                 `json:"notes"`
		Metrics   map[string]interface{} `json:"metrics"`
	}

	if err := c.Bind(&checkpointUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := services.UpdateCheckpoint(planID, checkpointID, userID, checkpointUpdate.Completed, checkpointUpdate.Notes, checkpointUpdate.Metrics)
	if err != nil {
		if err.Error() == "medical plan not found or unauthorized" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Medical plan not found or access denied"})
		}
		if err.Error() == "checkpoint not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Checkpoint not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update checkpoint"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Checkpoint updated successfully"})
}