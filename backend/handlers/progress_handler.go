package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"nutrition-platform/models"
	"github.com/labstack/echo/v4"
)

type ProgressHandler struct {
	progressPhotoRepo      ProgressPhotoRepository
	bodyMeasurementRepo    BodyMeasurementRepository
	milestoneRepo          MilestoneRepository
	weightGoalRepo         WeightGoalRepository
	progressAnalyticsRepo  ProgressAnalyticsRepository
}

func NewProgressHandler(
	progressPhotoRepo ProgressPhotoRepository,
	bodyMeasurementRepo BodyMeasurementRepository,
	milestoneRepo MilestoneRepository,
	weightGoalRepo WeightGoalRepository,
	progressAnalyticsRepo ProgressAnalyticsRepository,
) *ProgressHandler {
	return &ProgressHandler{
		progressPhotoRepo:      progressPhotoRepo,
		bodyMeasurementRepo:    bodyMeasurementRepo,
		milestoneRepo:          milestoneRepo,
		weightGoalRepo:         weightGoalRepo,
		progressAnalyticsRepo:  progressAnalyticsRepo,
	}
}

// Progress Photos

func (h *ProgressHandler) UploadProgressPhoto(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	description := c.FormValue("description")
	dateStr := c.FormValue("date")
	
	// Parse date
	var date time.Time
	var err error
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD"})
		}
	} else {
		date = time.Now()
	}
	
	// Get file
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}
	
	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file type. Only JPEG, PNG, and WebP are allowed"})
	}
	
	// Validate file size (10MB max)
	if file.Size > 10*1024*1024 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File too large. Maximum size is 10MB"})
	}
	
	// TODO: Implement file upload to cloud storage
	// For now, use a placeholder URL
	photoURL := "https://example.com/photos/" + file.Filename
	
	photo := &models.ProgressPhoto{
		UserID:      userID,
		Date:        date,
		PhotoURL:    photoURL,
		Description: description,
		Tags:        make(models.StringArray, 0),
	}
	
	err = h.progressPhotoRepo.CreatePhoto(c.Request().Context(), photo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload photo"})
	}
	
	return c.JSON(http.StatusCreated, photo)
}

func (h *ProgressHandler) GetProgressPhotos(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	photos, err := h.progressPhotoRepo.GetUserPhotos(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get photos"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"photos": photos,
		"count":  len(photos),
	})
}

func (h *ProgressHandler) DeleteProgressPhoto(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid photo ID"})
	}
	
	err = h.progressPhotoRepo.DeletePhoto(c.Request().Context(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete photo"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Photo deleted successfully"})
}

// Body Measurements

func (h *ProgressHandler) LogBodyMeasurement(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	var measurement models.BodyMeasurement
	if err := c.Bind(&measurement); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	measurement.UserID = userID
	
	// Validate
	if measurement.Weight <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Weight must be greater than 0"})
	}
	
	if err := measurement.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	err := h.bodyMeasurementRepo.CreateMeasurement(c.Request().Context(), &measurement)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to log measurement"})
	}
	
	return c.JSON(http.StatusCreated, measurement)
}

func (h *ProgressHandler) GetBodyMeasurementHistory(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	measurements, err := h.bodyMeasurementRepo.GetUserMeasurements(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get measurements"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"measurements": measurements,
		"count":        len(measurements),
	})
}

func (h *ProgressHandler) GetLatestBodyMeasurements(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	measurements, err := h.bodyMeasurementRepo.GetUserMeasurements(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get measurements"})
	}
	
	// Group by type and get latest for each type
	latestMeasurements := make(map[string]*models.BodyMeasurement)
	for _, measurement := range measurements {
		if existing, ok := latestMeasurements[measurement.Type]; !ok || measurement.Date.After(existing.Date) {
			latestMeasurements[measurement.Type] = measurement
		}
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"latest_measurements": latestMeasurements,
		"count":              len(latestMeasurements),
	})
}

// Milestones

func (h *ProgressHandler) CreateMilestone(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	var milestone models.Milestone
	if err := c.Bind(&milestone); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	milestone.UserID = userID
	
	if err := milestone.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	err := h.milestoneRepo.CreateMilestone(c.Request().Context(), &milestone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create milestone"})
	}
	
	return c.JSON(http.StatusCreated, milestone)
}

func (h *ProgressHandler) GetMilestones(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	milestones, err := h.milestoneRepo.GetUserMilestones(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get milestones"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"milestones": milestones,
		"count":      len(milestones),
	})
}

func (h *ProgressHandler) UpdateMilestone(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid milestone ID"})
	}
	
	var milestone models.Milestone
	if err := c.Bind(&milestone); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	milestone.ID = id
	milestone.UserID = userID
	
	if err := milestone.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	err = h.milestoneRepo.UpdateMilestone(c.Request().Context(), &milestone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update milestone"})
	}
	
	return c.JSON(http.StatusOK, milestone)
}

func (h *ProgressHandler) DeleteMilestone(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid milestone ID"})
	}
	
	err = h.milestoneRepo.DeleteMilestone(c.Request().Context(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete milestone"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Milestone deleted successfully"})
}

// Weight Goals

func (h *ProgressHandler) SetWeightGoal(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	var goal models.WeightGoal
	if err := c.Bind(&goal); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	goal.UserID = userID
	
	if err := goal.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	err := h.weightGoalRepo.CreateWeightGoal(c.Request().Context(), &goal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set weight goal"})
	}
	
	return c.JSON(http.StatusCreated, goal)
}

func (h *ProgressHandler) GetWeightGoals(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	goals, err := h.weightGoalRepo.GetUserWeightGoals(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get weight goals"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"goals": goals,
		"count": len(goals),
	})
}

func (h *ProgressHandler) GetActiveWeightGoal(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	goal, err := h.weightGoalRepo.GetActiveWeightGoal(c.Request().Context(), userID)
	if err != nil {
		if err.Error() == "no active weight goal found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No active weight goal found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get active weight goal"})
	}
	
	return c.JSON(http.StatusOK, goal)
}

// Progress Analytics

func (h *ProgressHandler) GetProgressSummary(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	summary, err := h.progressAnalyticsRepo.GetProgressSummary(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get progress summary"})
	}
	
	return c.JSON(http.StatusOK, summary)
}

func (h *ProgressHandler) GetWeightProgress(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	period := c.QueryParam("period")
	if period == "" {
		period = "1m" // Default to 1 month
	}
	
	progress, err := h.progressAnalyticsRepo.GetWeightProgress(c.Request().Context(), userID, period)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get weight progress"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"weight_progress": progress,
		"period":          period,
		"count":           len(progress),
	})
}

func (h *ProgressHandler) GetMeasurementTrends(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	
	measurementType := c.QueryParam("type")
	if measurementType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Measurement type is required"})
	}
	
	period := c.QueryParam("period")
	if period == "" {
		period = "1m" // Default to 1 month
	}
	
	trends, err := h.progressAnalyticsRepo.GetMeasurementTrends(c.Request().Context(), userID, measurementType, period)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get measurement trends"})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"trends": trends,
		"type":   measurementType,
		"period": period,
		"count":  len(trends),
	})
}

// Utility functions


// ProgressHandler interfaces (these would be defined in their respective repository files)

type ProgressPhotoRepository interface {
	CreatePhoto(ctx context.Context, photo *models.ProgressPhoto) error
	GetPhotosByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.ProgressPhoto, error)
	GetUserPhotos(ctx context.Context, userID int64) ([]*models.ProgressPhoto, error)
	DeletePhoto(ctx context.Context, id, userID int64) error
}

type BodyMeasurementRepository interface {
	CreateMeasurement(ctx context.Context, measurement *models.BodyMeasurement) error
	GetMeasurementsByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.BodyMeasurement, error)
	GetMeasurementsByType(ctx context.Context, userID int64, measurementType string) ([]*models.BodyMeasurement, error)
	GetUserMeasurements(ctx context.Context, userID int64) ([]*models.BodyMeasurement, error)
	GetLatestMeasurementByType(ctx context.Context, userID int64, measurementType string) (*models.BodyMeasurement, error)
	DeleteMeasurement(ctx context.Context, id, userID int64) error
}

type MilestoneRepository interface {
	CreateMilestone(ctx context.Context, milestone *models.Milestone) error
	GetMilestone(ctx context.Context, id, userID int64) (*models.Milestone, error)
	UpdateMilestone(ctx context.Context, milestone *models.Milestone) error
	DeleteMilestone(ctx context.Context, id, userID int64) error
	GetUserMilestones(ctx context.Context, userID int64) ([]*models.Milestone, error)
	GetActiveMilestones(ctx context.Context, userID int64) ([]*models.Milestone, error)
}

type WeightGoalRepository interface {
	CreateWeightGoal(ctx context.Context, goal *models.WeightGoal) error
	GetWeightGoal(ctx context.Context, id, userID int64) (*models.WeightGoal, error)
	UpdateWeightGoal(ctx context.Context, goal *models.WeightGoal) error
	DeleteWeightGoal(ctx context.Context, id, userID int64) error
	GetUserWeightGoals(ctx context.Context, userID int64) ([]*models.WeightGoal, error)
	GetActiveWeightGoal(ctx context.Context, userID int64) (*models.WeightGoal, error)
	GetWeightGoalHistory(ctx context.Context, userID int64) ([]*models.WeightGoal, error)
}

type ProgressAnalyticsRepository interface {
	GetProgressSummary(ctx context.Context, userID int64) (*models.ProgressSummary, error)
	GetWeightProgress(ctx context.Context, userID int64, period string) ([]*models.WeightProgress, error)
	GetMeasurementTrends(ctx context.Context, userID int64, measurementType, period string) ([]*models.MeasurementTrend, error)
}
