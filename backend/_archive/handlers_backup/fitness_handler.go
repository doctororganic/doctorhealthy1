package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	backendmodels "nutrition-platform/models"
	"nutrition-platform/repositories"

	"github.com/labstack/echo/v4"
)

// FitnessHandler handles fitness-related endpoints (workout logging)
type FitnessHandler struct {
	workoutLogRepo *repositories.WorkoutLogRepository
}

// NewFitnessHandler creates a new fitness handler
func NewFitnessHandler(db *sql.DB) *FitnessHandler {
	return &FitnessHandler{
		workoutLogRepo: repositories.NewWorkoutLogRepository(db),
	}
}

// LogWorkout logs a workout session
func (h *FitnessHandler) LogWorkout(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Convert userID to uint
	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int64:
		userIDUint = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			userIDUint = uint(id)
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
		WorkoutPlanID    *uint                    `json:"workout_plan_id"`
		WorkoutDate      time.Time                `json:"workout_date"`
		DurationMinutes  int                       `json:"duration_minutes" validate:"required,min=1"`
		Notes            *string                  `json:"notes"`
		CaloriesBurned   *int                     `json:"calories_burned"`
		Completed        bool                     `json:"completed"`
		CompletedExercises []backendmodels.CompletedExerciseItem `json:"completed_exercises"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Use current date if not provided
	if req.WorkoutDate.IsZero() {
		req.WorkoutDate = time.Now()
	}

	log := &backendmodels.WorkoutLog{
		UserID:            userIDUint,
		WorkoutPlanID:     req.WorkoutPlanID,
		WorkoutDate:       req.WorkoutDate,
		DurationMinutes:   req.DurationMinutes,
		CompletedExercises: backendmodels.CompletedExercises(req.CompletedExercises),
		Notes:             req.Notes,
		CaloriesBurned:    req.CaloriesBurned,
		Completed:         req.Completed,
	}

	ctx := context.Background()
	err := h.workoutLogRepo.CreateWorkoutLog(ctx, log)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to log workout: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Workout logged successfully",
		"data":    log,
	})
}

// GetWorkouts returns workout history for the current user
func (h *FitnessHandler) GetWorkouts(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// Convert userID to int64 for repository
	var userIDInt64 int64
	switch v := userID.(type) {
	case uint:
		userIDInt64 = int64(v)
	case int64:
		userIDInt64 = v
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

	// Build request
	req := &backendmodels.ListWorkoutLogsRequest{
		Page:  page,
		Limit: limit,
	}

	userIDUint := uint(userIDInt64)
	req.UserID = &userIDUint

	// Parse date filters
	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = &startDate
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			req.EndDate = &endDate
		}
	}

	// Parse completed filter
	if completedStr := c.QueryParam("completed"); completedStr != "" {
		completed := completedStr == "true"
		req.Completed = &completed
	}

	ctx := context.Background()
	response, err := h.workoutLogRepo.ListWorkoutLogs(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch workouts: " + err.Error(),
		})
	}

	// Convert WorkoutLogWithPlan to WorkoutLog for response
	workoutLogs := make([]backendmodels.WorkoutLog, len(response.WorkoutLogs))
	for i, wlp := range response.WorkoutLogs {
		// WorkoutLogWithPlan embeds WorkoutLog, so we can assign directly
		workoutLogs[i] = backendmodels.WorkoutLog{
			ID:              wlp.ID,
			UserID:          wlp.UserID,
			WorkoutPlanID:   wlp.WorkoutPlanID,
			WorkoutDate:     wlp.WorkoutDate,
			DurationMinutes: wlp.DurationMinutes,
			CompletedExercises: wlp.CompletedExercises,
			Notes:           wlp.Notes,
			CaloriesBurned:  wlp.CaloriesBurned,
			Completed:       wlp.Completed,
			CreatedAt:       wlp.CreatedAt,
			UpdatedAt:       wlp.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   workoutLogs,
		"pagination": map[string]interface{}{
			"page":     response.Page,
			"limit":    response.Limit,
			"total":    response.Total,
			"has_next": response.HasNext,
		},
	})
}

// GetWorkout returns a specific workout log by ID
func (h *FitnessHandler) GetWorkout(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid workout ID",
		})
	}

	ctx := context.Background()
	log, err := h.workoutLogRepo.GetWorkoutLogByID(ctx, workoutID)
	if err != nil {
		if err.Error() == "workout log not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Workout not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch workout: " + err.Error(),
		})
	}

	// Verify ownership
	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int64:
		userIDUint = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			userIDUint = uint(id)
		}
	}

	if log.UserID != userIDUint {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You can only view your own workouts",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   log,
	})
}

// UpdateWorkout updates an existing workout log
func (h *FitnessHandler) UpdateWorkout(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid workout ID",
		})
	}

	// Get existing workout
	ctx := context.Background()
	existingLog, err := h.workoutLogRepo.GetWorkoutLogByID(ctx, workoutID)
	if err != nil {
		if err.Error() == "workout log not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Workout not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch workout",
		})
	}

	// Verify ownership
	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int64:
		userIDUint = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			userIDUint = uint(id)
		}
	}

	if existingLog.UserID != userIDUint {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You can only update your own workouts",
		})
	}

	var req struct {
		WorkoutPlanID     *uint                     `json:"workout_plan_id"`
		WorkoutDate       *time.Time                `json:"workout_date"`
		DurationMinutes   *int                      `json:"duration_minutes"`
		Notes             *string                   `json:"notes"`
		CaloriesBurned    *int                      `json:"calories_burned"`
		Completed         *bool                     `json:"completed"`
		CompletedExercises []backendmodels.CompletedExerciseItem `json:"completed_exercises"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Update fields
	if req.WorkoutPlanID != nil {
		existingLog.WorkoutPlanID = req.WorkoutPlanID
	}
	if req.WorkoutDate != nil {
		existingLog.WorkoutDate = *req.WorkoutDate
	}
	if req.DurationMinutes != nil {
		existingLog.DurationMinutes = *req.DurationMinutes
	}
	if req.Notes != nil {
		existingLog.Notes = req.Notes
	}
	if req.CaloriesBurned != nil {
		existingLog.CaloriesBurned = req.CaloriesBurned
	}
	if req.Completed != nil {
		existingLog.Completed = *req.Completed
	}
	if req.CompletedExercises != nil {
		existingLog.CompletedExercises = backendmodels.CompletedExercises(req.CompletedExercises)
	}

	err = h.workoutLogRepo.UpdateWorkoutLog(ctx, existingLog)
	if err != nil {
		if err.Error() == "workout log not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Workout not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update workout: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Workout updated successfully",
		"data":    existingLog,
	})
}

// DeleteWorkout deletes a workout log
func (h *FitnessHandler) DeleteWorkout(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid workout ID",
		})
	}

	// Verify ownership
	ctx := context.Background()
	existingLog, err := h.workoutLogRepo.GetWorkoutLogByID(ctx, workoutID)
	if err != nil {
		if err.Error() == "workout log not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Workout not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch workout",
		})
	}

	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int64:
		userIDUint = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			userIDUint = uint(id)
		}
	}

	if existingLog.UserID != userIDUint {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You can only delete your own workouts",
		})
	}

	err = h.workoutLogRepo.DeleteWorkoutLog(ctx, workoutID)
	if err != nil {
		if err.Error() == "workout log not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Workout not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete workout: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Workout deleted successfully",
	})
}

