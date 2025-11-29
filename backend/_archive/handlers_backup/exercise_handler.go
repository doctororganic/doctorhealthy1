package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"nutrition-platform/models"
	"nutrition-platform/repositories"

	"github.com/labstack/echo/v4"
)

// ExerciseHandler handles exercise CRUD operations
type ExerciseHandler struct {
	exerciseRepo *repositories.ExerciseRepository
}

// NewExerciseHandler creates a new exercise handler
func NewExerciseHandler(db *sql.DB) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseRepo: repositories.NewExerciseRepository(db),
	}
}

// GetExercises returns paginated list of exercises
func (h *ExerciseHandler) GetExercises(c echo.Context) error {
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

	// Parse filters
	req := &models.ListExercisesRequest{
		Limit:  limit,
		Offset: (page - 1) * limit,
		SortBy: "created_at DESC",
	}

	if search := c.QueryParam("search"); search != "" {
		req.Search = search
	}

	if muscleGroup := c.QueryParam("muscle_group"); muscleGroup != "" {
		req.MuscleGroup = muscleGroup
	}

	if equipment := c.QueryParam("equipment"); equipment != "" {
		req.Equipment = equipment
	}

	if difficulty := c.QueryParam("difficulty"); difficulty != "" {
		req.Difficulty = difficulty
	}

	// Convert userID to int64 for repository
	var userIDInt64 int64
	switch v := userID.(type) {
	case int64:
		userIDInt64 = v
	case uint:
		userIDInt64 = int64(v)
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			userIDInt64 = id
		}
	default:
		userIDInt64 = 0
	}

	// Only show public exercises or user's own exercises
	isPublic := true
	req.IsPublic = &isPublic
	req.CreatedBy = &userIDInt64

	ctx := context.Background()
	response, err := h.exerciseRepo.ListExercises(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch exercises: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   response.Exercises,
		"pagination": map[string]interface{}{
			"page":        response.Page,
			"limit":       response.Limit,
			"total":       response.Total,
			"total_pages": response.TotalPages,
			"has_next":    response.Page < response.TotalPages,
			"has_prev":    response.Page > 1,
		},
	})
}

// GetExercise returns a specific exercise by ID
func (h *ExerciseHandler) GetExercise(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	exerciseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid exercise ID",
		})
	}

	ctx := context.Background()
	exercise, err := h.exerciseRepo.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		if err.Error() == "exercise not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Exercise not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch exercise: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   exercise,
	})
}

// CreateExercise creates a new exercise
func (h *ExerciseHandler) CreateExercise(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Convert userID to int64
	var userIDInt64 int64
	switch v := userID.(type) {
	case int64:
		userIDInt64 = v
	case uint:
		userIDInt64 = int64(v)
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			userIDInt64 = id
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid user ID",
			})
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID type",
		})
	}

	var req struct {
		Name         string   `json:"name" validate:"required"`
		Description  string   `json:"description"`
		MuscleGroups []string `json:"muscle_groups" validate:"required"`
		Equipment    []string `json:"equipment"`
		Difficulty   string   `json:"difficulty" validate:"required,oneof=beginner intermediate advanced"`
		Instructions string   `json:"instructions"`
		Tips         string   `json:"tips"`
		IsPublic     bool     `json:"is_public"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	exercise := &models.Exercise{
		Name:         req.Name,
		Description:  req.Description,
		MuscleGroups: models.MuscleGroups(req.MuscleGroups),
		Equipment:    models.Equipment(req.Equipment),
		Difficulty:   req.Difficulty,
		Instructions: req.Instructions,
		Tips:         req.Tips,
		CreatedBy:    userIDInt64,
		IsPublic:     req.IsPublic,
	}

	ctx := context.Background()
	err := h.exerciseRepo.CreateExercise(ctx, exercise)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create exercise: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Exercise created successfully",
		"data":    exercise,
	})
}

// UpdateExercise updates an existing exercise
func (h *ExerciseHandler) UpdateExercise(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	exerciseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid exercise ID",
		})
	}

	// Get existing exercise to verify ownership
	ctx := context.Background()
	existingExercise, err := h.exerciseRepo.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		if err.Error() == "exercise not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Exercise not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch exercise",
		})
	}

	// Check ownership (only creator can update)
	var userIDInt64 int64
	switch v := userID.(type) {
	case int64:
		userIDInt64 = v
	case uint:
		userIDInt64 = int64(v)
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			userIDInt64 = id
		}
	}

	if existingExercise.CreatedBy != userIDInt64 {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You can only update your own exercises",
		})
	}

	var req struct {
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		MuscleGroups []string `json:"muscle_groups"`
		Equipment    []string `json:"equipment"`
		Difficulty   string   `json:"difficulty"`
		Instructions string   `json:"instructions"`
		Tips         string   `json:"tips"`
		IsPublic     *bool    `json:"is_public"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Update fields
	if req.Name != "" {
		existingExercise.Name = req.Name
	}
	if req.Description != "" {
		existingExercise.Description = req.Description
	}
	if len(req.MuscleGroups) > 0 {
		existingExercise.MuscleGroups = models.MuscleGroups(req.MuscleGroups)
	}
	if len(req.Equipment) > 0 {
		existingExercise.Equipment = models.Equipment(req.Equipment)
	}
	if req.Difficulty != "" {
		existingExercise.Difficulty = req.Difficulty
	}
	if req.Instructions != "" {
		existingExercise.Instructions = req.Instructions
	}
	if req.Tips != "" {
		existingExercise.Tips = req.Tips
	}
	if req.IsPublic != nil {
		existingExercise.IsPublic = *req.IsPublic
	}

	err = h.exerciseRepo.UpdateExercise(ctx, existingExercise)
	if err != nil {
		if err.Error() == "exercise not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Exercise not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update exercise: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Exercise updated successfully",
		"data":    existingExercise,
	})
}

// DeleteExercise deletes an exercise
func (h *ExerciseHandler) DeleteExercise(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	exerciseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid exercise ID",
		})
	}

	// Verify ownership before deleting
	ctx := context.Background()
	existingExercise, err := h.exerciseRepo.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		if err.Error() == "exercise not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Exercise not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch exercise",
		})
	}

	// Check ownership
	var userIDInt64 int64
	switch v := userID.(type) {
	case int64:
		userIDInt64 = v
	case uint:
		userIDInt64 = int64(v)
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			userIDInt64 = id
		}
	}

	if existingExercise.CreatedBy != userIDInt64 {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You can only delete your own exercises",
		})
	}

	err = h.exerciseRepo.DeleteExercise(ctx, exerciseID)
	if err != nil {
		if err.Error() == "exercise not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Exercise not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete exercise: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Exercise deleted successfully",
	})
}

// SearchExercises searches for exercises
func (h *ExerciseHandler) SearchExercises(c echo.Context) error {
	return h.GetExercises(c) // Reuse GetExercises which already supports search
}

