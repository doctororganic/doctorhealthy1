package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"nutrition-platform/middleware"
	"nutrition-platform/models"
	"nutrition-platform/repositories"
)

type WorkoutHandler struct {
	exerciseRepo        *repositories.ExerciseRepository
	workoutPlanRepo     *repositories.WorkoutPlanRepository
	workoutLogRepo      *repositories.WorkoutLogRepository
	personalRecordRepo  *repositories.PersonalRecordRepository
}

func NewWorkoutHandler(
	exerciseRepo *repositories.ExerciseRepository,
	workoutPlanRepo *repositories.WorkoutPlanRepository,
	workoutLogRepo *repositories.WorkoutLogRepository,
	personalRecordRepo *repositories.PersonalRecordRepository,
) *WorkoutHandler {
	return &WorkoutHandler{
		exerciseRepo:       exerciseRepo,
		workoutPlanRepo:    workoutPlanRepo,
		workoutLogRepo:     workoutLogRepo,
		personalRecordRepo: personalRecordRepo,
	}
}

// Exercise endpoints

// CreateExercise creates a new exercise
func (h *WorkoutHandler) CreateExercise(c echo.Context) error {
	var req models.CreateExerciseRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	exercise := &models.Exercise{
		Name:        req.Name,
		Description: req.Description,
		MuscleGroups: models.MuscleGroups(req.MuscleGroups),
		Equipment:   models.Equipment(req.Equipment),
		Difficulty:  req.Difficulty,
		Instructions: req.Instructions,
		Tips:        req.Tips,
		CreatedBy:   userID,
		IsPublic:    req.IsPublic,
	}

	if err := h.exerciseRepo.CreateExercise(context.Background(), exercise); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create exercise")
	}

	resp := &models.ExerciseResponse{
		ID:           exercise.ID,
		Name:         exercise.Name,
		Description:  exercise.Description,
		MuscleGroups: exercise.MuscleGroups,
		Equipment:    exercise.Equipment,
		Difficulty:   exercise.Difficulty,
		Instructions: exercise.Instructions,
		Tips:         exercise.Tips,
		CreatedBy:    exercise.CreatedBy,
		IsPublic:     exercise.IsPublic,
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Exercise created successfully",
		"exercise": resp,
	})
}

// GetExercise retrieves an exercise by ID
func (h *WorkoutHandler) GetExercise(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid exercise ID")
	}

	exercise, err := h.exerciseRepo.GetExerciseByID(context.Background(), id)
	if err != nil {
		if err.Error() == "exercise not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Exercise not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get exercise")
	}

	resp := &models.ExerciseResponse{
		ID:           exercise.ID,
		Name:         exercise.Name,
		Description:  exercise.Description,
		MuscleGroups: exercise.MuscleGroups,
		Equipment:    exercise.Equipment,
		Difficulty:   exercise.Difficulty,
		Instructions: exercise.Instructions,
		Tips:         exercise.Tips,
		CreatedBy:    exercise.CreatedBy,
		IsPublic:     exercise.IsPublic,
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateExercise updates an existing exercise
func (h *WorkoutHandler) UpdateExercise(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid exercise ID")
	}

	var req models.UpdateExerciseRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Get existing exercise
	exercise, err := h.exerciseRepo.GetExerciseByID(context.Background(), id)
	if err != nil {
		if err.Error() == "exercise not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Exercise not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get exercise")
	}

	// Check permissions
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	if exercise.CreatedBy != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You can only update your own exercises")
	}

	// Update fields
	if req.Name != nil {
		exercise.Name = *req.Name
	}
	if req.Description != nil {
		exercise.Description = *req.Description
	}
	if req.MuscleGroups != nil {
		exercise.MuscleGroups = models.MuscleGroups(req.MuscleGroups)
	}
	if req.Equipment != nil {
		exercise.Equipment = models.Equipment(req.Equipment)
	}
	if req.Difficulty != nil {
		exercise.Difficulty = *req.Difficulty
	}
	if req.Instructions != nil {
		exercise.Instructions = *req.Instructions
	}
	if req.Tips != nil {
		exercise.Tips = *req.Tips
	}
	if req.IsPublic != nil {
		exercise.IsPublic = *req.IsPublic
	}

	if err := h.exerciseRepo.UpdateExercise(context.Background(), exercise); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update exercise")
	}

	resp := &models.ExerciseResponse{
		ID:           exercise.ID,
		Name:         exercise.Name,
		Description:  exercise.Description,
		MuscleGroups: exercise.MuscleGroups,
		Equipment:    exercise.Equipment,
		Difficulty:   exercise.Difficulty,
		Instructions: exercise.Instructions,
		Tips:         exercise.Tips,
		CreatedBy:    exercise.CreatedBy,
		IsPublic:     exercise.IsPublic,
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Exercise updated successfully",
		"exercise": resp,
	})
}

// DeleteExercise deletes an exercise
func (h *WorkoutHandler) DeleteExercise(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid exercise ID")
	}

	// Get existing exercise
	exercise, err := h.exerciseRepo.GetExerciseByID(context.Background(), id)
	if err != nil {
		if err.Error() == "exercise not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Exercise not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get exercise")
	}

	// Check permissions
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	if exercise.CreatedBy != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You can only delete your own exercises")
	}

	if err := h.exerciseRepo.DeleteExercise(context.Background(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete exercise")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Exercise deleted successfully",
	})
}

// ListExercises retrieves exercises with filtering and pagination
func (h *WorkoutHandler) ListExercises(c echo.Context) error {
	req := &models.ListExercisesRequest{
		Limit: 20,
		Offset: 0,
		SortBy: "created_at DESC",
	}

	// Parse query parameters
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
	if createdBy := c.QueryParam("created_by"); createdBy != "" {
		if id, err := strconv.ParseInt(createdBy, 10, 64); err == nil {
			req.CreatedBy = &id
		}
	}
	if isPublic := c.QueryParam("is_public"); isPublic != "" {
		if val, err := strconv.ParseBool(isPublic); err == nil {
			req.IsPublic = &val
		}
	}
	if limit := c.QueryParam("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			req.Limit = val
		}
	}
	if offset := c.QueryParam("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			req.Offset = val
		}
	}
	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	response, err := h.exerciseRepo.ListExercises(context.Background(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list exercises")
	}

	return c.JSON(http.StatusOK, response)
}

// Workout Plan endpoints

// CreateWorkoutPlan creates a new workout plan
func (h *WorkoutHandler) CreateWorkoutPlan(c echo.Context) error {
	var req models.CreateWorkoutPlanRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	plan := &models.WorkoutPlan{
		Name:          req.Name,
		Description:   req.Description,
		Exercises:     models.WorkoutPlanExercises(req.Exercises),
		UserID:        userID,
		IsPublic:      req.IsPublic,
		IsTemplate:    req.IsTemplate,
		DurationWeeks: req.DurationWeeks,
		Difficulty:    req.Difficulty,
	}

	if err := h.workoutPlanRepo.CreateWorkoutPlan(context.Background(), plan); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create workout plan")
	}

	resp := &models.WorkoutPlanResponse{
		ID:            plan.ID,
		Name:          plan.Name,
		Description:   plan.Description,
		Exercises:     plan.Exercises,
		UserID:        plan.UserID,
		IsPublic:      plan.IsPublic,
		IsTemplate:    plan.IsTemplate,
		DurationWeeks: plan.DurationWeeks,
		Difficulty:    plan.Difficulty,
		CreatedAt:     plan.CreatedAt,
		UpdatedAt:     plan.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Workout plan created successfully",
		"plan":    resp,
	})
}

// GetWorkoutPlan retrieves a workout plan by ID
func (h *WorkoutHandler) GetWorkoutPlan(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workout plan ID")
	}

	plan, err := h.workoutPlanRepo.GetWorkoutPlanByID(context.Background(), id)
	if err != nil {
		if err.Error() == "workout plan not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Workout plan not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get workout plan")
	}

	resp := &models.WorkoutPlanResponse{
		ID:            plan.ID,
		Name:          plan.Name,
		Description:   plan.Description,
		Exercises:     plan.Exercises,
		UserID:        plan.UserID,
		IsPublic:      plan.IsPublic,
		IsTemplate:    plan.IsTemplate,
		DurationWeeks: plan.DurationWeeks,
		Difficulty:    plan.Difficulty,
		CreatedAt:     plan.CreatedAt,
		UpdatedAt:     plan.UpdatedAt,
	}

	return c.JSON(http.StatusOK, resp)
}

// ListWorkoutPlans retrieves workout plans with filtering and pagination
func (h *WorkoutHandler) ListWorkoutPlans(c echo.Context) error {
	req := &models.ListWorkoutPlansRequest{
		Limit: 20,
		Offset: 0,
		SortBy: "created_at DESC",
	}

	// Parse query parameters
	if search := c.QueryParam("search"); search != "" {
		req.Search = search
	}
	if userID := c.QueryParam("user_id"); userID != "" {
		if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
			req.UserID = &id
		}
	}
	if isPublic := c.QueryParam("is_public"); isPublic != "" {
		if val, err := strconv.ParseBool(isPublic); err == nil {
			req.IsPublic = &val
		}
	}
	if isTemplate := c.QueryParam("is_template"); isTemplate != "" {
		if val, err := strconv.ParseBool(isTemplate); err == nil {
			req.IsTemplate = &val
		}
	}
	if difficulty := c.QueryParam("difficulty"); difficulty != "" {
		req.Difficulty = difficulty
	}
	if limit := c.QueryParam("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			req.Limit = val
		}
	}
	if offset := c.QueryParam("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			req.Offset = val
		}
	}
	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	response, err := h.workoutPlanRepo.ListWorkoutPlans(context.Background(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list workout plans")
	}

	return c.JSON(http.StatusOK, response)
}

// Workout Log endpoints

// CreateWorkoutLog creates a new workout log
func (h *WorkoutHandler) CreateWorkoutLog(c echo.Context) error {
	var req models.CreateWorkoutLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	log := &models.WorkoutLog{
		UserID:           userID,
		WorkoutPlanID:    req.WorkoutPlanID,
		WorkoutDate:      req.WorkoutDate,
		DurationMinutes:  req.DurationMinutes,
		CompletedExercises: models.CompletedExercises(req.CompletedExercises),
		Notes:            req.Notes,
	}

	if err := h.workoutLogRepo.CreateWorkoutLog(context.Background(), log); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create workout log")
	}

	resp := &models.WorkoutLogResponse{
		ID:                log.ID,
		UserID:            log.UserID,
		WorkoutPlanID:     log.WorkoutPlanID,
		WorkoutDate:       log.WorkoutDate,
		DurationMinutes:   log.DurationMinutes,
		CompletedExercises: log.CompletedExercises,
		Notes:             log.Notes,
		CreatedAt:         log.CreatedAt,
		UpdatedAt:         log.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Workout log created successfully",
		"log":     resp,
	})
}

// GetWorkoutLog retrieves a workout log by ID
func (h *WorkoutHandler) GetWorkoutLog(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid workout log ID")
	}

	log, err := h.workoutLogRepo.GetWorkoutLogByID(context.Background(), id)
	if err != nil {
		if err.Error() == "workout log not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Workout log not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get workout log")
	}

	// Check permissions
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	if log.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You can only view your own workout logs")
	}

	resp := &models.WorkoutLogResponse{
		ID:                log.ID,
		UserID:            log.UserID,
		WorkoutPlanID:     log.WorkoutPlanID,
		WorkoutDate:       log.WorkoutDate,
		DurationMinutes:   log.DurationMinutes,
		CompletedExercises: log.CompletedExercises,
		Notes:             log.Notes,
		CreatedAt:         log.CreatedAt,
		UpdatedAt:         log.UpdatedAt,
	}

	return c.JSON(http.StatusOK, resp)
}

// ListWorkoutLogs retrieves workout logs with filtering and pagination
func (h *WorkoutHandler) ListWorkoutLogs(c echo.Context) error {
	req := &models.ListWorkoutLogsRequest{
		Limit: 20,
		Offset: 0,
		SortBy: "workout_date DESC",
	}

	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}
	req.UserID = &userID

	// Parse query parameters
	if workoutPlanID := c.QueryParam("workout_plan_id"); workoutPlanID != "" {
		if id, err := strconv.ParseInt(workoutPlanID, 10, 64); err == nil {
			req.WorkoutPlanID = &id
		}
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = &date
		}
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = &date
		}
	}
	if limit := c.QueryParam("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			req.Limit = val
		}
	}
	if offset := c.QueryParam("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			req.Offset = val
		}
	}
	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	response, err := h.workoutLogRepo.ListWorkoutLogs(context.Background(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list workout logs")
	}

	return c.JSON(http.StatusOK, response)
}

// GetUserWorkoutStats retrieves workout statistics for the authenticated user
func (h *WorkoutHandler) GetUserWorkoutStats(c echo.Context) error {
	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	period := c.QueryParam("period")
	if period == "" {
		period = "month" // Default to month
	}

	stats, err := h.workoutLogRepo.GetUserWorkoutStats(context.Background(), userID, period)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get workout stats")
	}

	return c.JSON(http.StatusOK, stats)
}

// Personal Record endpoints

// CreatePersonalRecord creates a new personal record
func (h *WorkoutHandler) CreatePersonalRecord(c echo.Context) error {
	var req models.CreatePersonalRecordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	pr := &models.PersonalRecord{
		UserID:       userID,
		ExerciseID:   req.ExerciseID,
		RecordType:   req.RecordType,
		RecordValue:  req.RecordValue,
		WorkoutLogID: req.WorkoutLogID,
		AchievedAt:   req.AchievedAt,
	}

	if err := h.personalRecordRepo.CreatePersonalRecord(context.Background(), pr); err != nil {
		if err.Error() == "new record" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create personal record")
	}

	resp := &models.PersonalRecordResponse{
		ID:          pr.ID,
		UserID:      pr.UserID,
		ExerciseID:  pr.ExerciseID,
		RecordType:  pr.RecordType,
		RecordValue: pr.RecordValue,
		WorkoutLogID: pr.WorkoutLogID,
		AchievedAt:  pr.AchievedAt,
		CreatedAt:   pr.CreatedAt,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Personal record created successfully",
		"record":  resp,
	})
}

// ListPersonalRecords retrieves personal records with filtering and pagination
func (h *WorkoutHandler) ListPersonalRecords(c echo.Context) error {
	req := &models.ListPersonalRecordsRequest{
		Limit: 20,
		Offset: 0,
		SortBy: "achieved_at DESC",
	}

	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}
	req.UserID = &userID

	// Parse query parameters
	if exerciseID := c.QueryParam("exercise_id"); exerciseID != "" {
		if id, err := strconv.ParseInt(exerciseID, 10, 64); err == nil {
			req.ExerciseID = &id
		}
	}
	if recordType := c.QueryParam("record_type"); recordType != "" {
		req.RecordType = recordType
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			req.StartDate = &date
		}
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			req.EndDate = &date
		}
	}
	if limit := c.QueryParam("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			req.Limit = val
		}
	}
	if offset := c.QueryParam("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			req.Offset = val
		}
	}
	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	response, err := h.personalRecordRepo.ListPersonalRecords(context.Background(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list personal records")
	}

	return c.JSON(http.StatusOK, response)
}

// GetUserPersonalRecords retrieves all personal records for the authenticated user
func (h *WorkoutHandler) GetUserPersonalRecords(c echo.Context) error {
	// Get user ID from context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	records, err := h.personalRecordRepo.GetUserPersonalRecords(context.Background(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user personal records")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"records": records,
	})
}

// GetMuscleGroups returns all available muscle groups
func (h *WorkoutHandler) GetMuscleGroups(c echo.Context) error {
	muscleGroups, err := h.exerciseRepo.GetMuscleGroups(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get muscle groups")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"muscle_groups": muscleGroups,
	})
}

// GetEquipmentTypes returns all available equipment types
func (h *WorkoutHandler) GetEquipmentTypes(c echo.Context) error {
	equipment, err := h.exerciseRepo.GetEquipmentTypes(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get equipment types")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"equipment": equipment,
	})
}

// Helper function to validate request and get correlation ID
func validateRequestAndGetCorrelationID(c echo.Context, req interface{}) (string, error) {
	if err := c.Bind(req); err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	correlationID := middleware.GetCorrelationID(c)
	if correlationID == "" {
		correlationID = "unknown"
	}

	return correlationID, nil
}
