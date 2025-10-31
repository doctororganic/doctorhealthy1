package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"nutrition-platform/models"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles health-related endpoints
type HealthHandler struct {
	healthService *services.HealthService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// CreateHealthComplaint creates a new health complaint
func (h *HealthHandler) CreateHealthComplaint(c echo.Context) error {
	var req models.CreateHealthComplaintRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "HEALTH_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "HEALTH_002",
		})
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "HEALTH_003",
		})
	}

	complaint, err := h.healthService.CreateHealthComplaint(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "creation_failed",
			"message": err.Error(),
			"code":    "HEALTH_004",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success":   true,
		"complaint": complaint,
		"message":   "Health complaint recorded successfully",
	})
}

// CreateUserInjury creates a new user injury record
func (h *HealthHandler) CreateUserInjury(c echo.Context) error {
	var req models.CreateUserInjuryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "HEALTH_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "HEALTH_002",
		})
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "HEALTH_003",
		})
	}

	injury, err := h.healthService.CreateUserInjury(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "creation_failed",
			"message": err.Error(),
			"code":    "HEALTH_005",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"injury":  injury,
		"message": "Injury record created successfully",
	})
}

// GetHealthConditions retrieves health conditions
func (h *HealthHandler) GetHealthConditions(c echo.Context) error {
	category := c.QueryParam("category")
	page := 1
	limit := 20

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	conditions, err := h.healthService.GetHealthConditions(category, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "retrieval_failed",
			"message": err.Error(),
			"code":    "HEALTH_006",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"conditions": conditions,
		"page":       page,
		"limit":      limit,
	})
}

// GetUserHealthComplaints retrieves user's health complaints
func (h *HealthHandler) GetUserHealthComplaints(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "HEALTH_003",
		})
	}

	status := c.QueryParam("status")

	complaints, err := h.healthService.GetUserHealthComplaints(userID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "retrieval_failed",
			"message": err.Error(),
			"code":    "HEALTH_007",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"complaints": complaints,
	})
}

// GetUserInjuries retrieves user's injuries
func (h *HealthHandler) GetUserInjuries(c echo.Context) error {
	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "HEALTH_003",
		})
	}

	status := c.QueryParam("status")

	injuries, err := h.healthService.GetUserInjuries(userID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "retrieval_failed",
			"message": err.Error(),
			"code":    "HEALTH_008",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"injuries": injuries,
	})
}

// PerformHealthAssessment performs a comprehensive health assessment
func (h *HealthHandler) PerformHealthAssessment(c echo.Context) error {
	var req models.HealthAssessmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "HEALTH_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "HEALTH_002",
		})
	}

	assessment, err := h.healthService.PerformHealthAssessment(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "assessment_failed",
			"message": err.Error(),
			"code":    "HEALTH_009",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"assessment": assessment,
		"message":    "Health assessment completed successfully",
	})
}

// GetHealthRiskAssessment provides a quick health risk assessment
func (h *HealthHandler) GetHealthRiskAssessment(c echo.Context) error {
	var req models.HealthAssessmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "HEALTH_001",
		})
	}

	// Simplified validation for risk assessment
	if req.Age <= 0 || req.Height <= 0 || req.Weight <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_data",
			"message": "Age, height, and weight are required",
			"code":    "HEALTH_010",
		})
	}

	assessment, err := h.healthService.PerformHealthAssessment(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "assessment_failed",
			"message": err.Error(),
			"code":    "HEALTH_009",
		})
	}

	// Return simplified risk assessment
	riskAssessment := map[string]interface{}{
		"bmi":                 assessment.BMI,
		"bmi_category":        assessment.BMICategory,
		"risk_score":          assessment.RiskScore,
		"health_risk_factors": assessment.HealthRiskFactors,
		"red_flags":           assessment.RedFlags,
		"key_recommendations": getKeyRecommendations(assessment.Recommendations),
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"risk_assessment": riskAssessment,
	})
}

// GetSymptomChecker provides symptom checking functionality
func (h *HealthHandler) GetSymptomChecker(c echo.Context) error {
	symptoms := c.QueryParam("symptoms")
	if symptoms == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_symptoms",
			"message": "Symptoms parameter is required",
			"code":    "HEALTH_011",
		})
	}

	symptomList := parseCommaSeparated(symptoms)
	if len(symptomList) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_symptoms",
			"message": "At least one symptom must be provided",
			"code":    "HEALTH_012",
		})
	}

	// Simple symptom checking logic
	result := h.performSymptomCheck(symptomList)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"result":     result,
		"disclaimer": "This is for informational purposes only and should not replace professional medical advice. Consult with a healthcare provider for proper diagnosis and treatment.",
	})
}

// GetHealthTips provides general health tips
func (h *HealthHandler) GetHealthTips(c echo.Context) error {
	category := c.QueryParam("category")
	age := c.QueryParam("age")
	gender := c.QueryParam("gender")

	tips := h.generateHealthTips(category, age, gender)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"tips":    tips,
	})
}

// RegisterRoutes registers health routes
func (h *HealthHandler) RegisterRoutes(e *echo.Group) {
	health := e.Group("/health")

	// Health complaints
	health.POST("/complaints", h.CreateHealthComplaint)
	health.GET("/complaints", h.GetUserHealthComplaints)

	// Injuries
	health.POST("/injuries", h.CreateUserInjury)
	health.GET("/injuries", h.GetUserInjuries)

	// Health conditions
	health.GET("/conditions", h.GetHealthConditions)

	// Assessments
	health.POST("/assessment", h.PerformHealthAssessment)
	health.POST("/risk-assessment", h.GetHealthRiskAssessment)

	// Utilities
	health.GET("/symptom-checker", h.GetSymptomChecker)
	health.GET("/tips", h.GetHealthTips)
}

// Helper functions

func getKeyRecommendations(recommendations []models.HealthRecommendation) []string {
	var key []string
	for _, rec := range recommendations {
		if rec.Priority == "high" {
			key = append(key, rec.Title)
		}
	}
	if len(key) == 0 && len(recommendations) > 0 {
		// If no high priority, take first few recommendations
		for i, rec := range recommendations {
			if i >= 3 {
				break
			}
			key = append(key, rec.Title)
		}
	}
	return key
}

func (h *HealthHandler) performSymptomCheck(symptoms []string) map[string]interface{} {
	// Simple symptom checking logic - in production, this would be more sophisticated
	urgentSymptoms := []string{
		"chest pain", "difficulty breathing", "severe headache", "loss of consciousness",
		"severe abdominal pain", "high fever", "severe bleeding", "stroke symptoms",
	}

	moderateSymptoms := []string{
		"persistent cough", "fever", "nausea", "dizziness", "fatigue",
		"muscle pain", "joint pain", "headache",
	}

	var urgency string
	var possibleConditions []string
	var recommendations []string

	// Check for urgent symptoms
	hasUrgent := false
	for _, symptom := range symptoms {
		for _, urgent := range urgentSymptoms {
			if strings.Contains(strings.ToLower(symptom), urgent) {
				hasUrgent = true
				break
			}
		}
		if hasUrgent {
			break
		}
	}

	if hasUrgent {
		urgency = "urgent"
		recommendations = []string{
			"Seek immediate medical attention",
			"Call emergency services if symptoms are severe",
			"Do not delay medical care",
		}
	} else {
		// Check for moderate symptoms
		hasModerate := false
		for _, symptom := range symptoms {
			for _, moderate := range moderateSymptoms {
				if strings.Contains(strings.ToLower(symptom), moderate) {
					hasModerate = true
					break
				}
			}
			if hasModerate {
				break
			}
		}

		if hasModerate {
			urgency = "moderate"
			recommendations = []string{
				"Consider consulting with a healthcare provider",
				"Monitor symptoms and seek care if they worsen",
				"Rest and stay hydrated",
			}
		} else {
			urgency = "low"
			recommendations = []string{
				"Monitor symptoms",
				"Consider home remedies if appropriate",
				"Consult healthcare provider if symptoms persist",
			}
		}
	}

	// Simple condition mapping (very basic)
	if containsInsensitive(symptoms, "fever") && containsInsensitive(symptoms, "cough") {
		possibleConditions = append(possibleConditions, "Upper respiratory infection", "Flu")
	}
	if containsInsensitive(symptoms, "headache") && containsInsensitive(symptoms, "nausea") {
		possibleConditions = append(possibleConditions, "Migraine", "Tension headache")
	}
	if containsInsensitive(symptoms, "fatigue") && containsInsensitive(symptoms, "muscle pain") {
		possibleConditions = append(possibleConditions, "Viral infection", "Overexertion")
	}

	return map[string]interface{}{
		"urgency":             urgency,
		"possible_conditions": possibleConditions,
		"recommendations":     recommendations,
		"symptoms_analyzed":   symptoms,
	}
}

func (h *HealthHandler) generateHealthTips(category, age, gender string) []string {
	var tips []string

	switch category {
	case "nutrition":
		tips = []string{
			"Eat a variety of colorful fruits and vegetables daily",
			"Choose whole grains over refined grains",
			"Include lean proteins in your meals",
			"Stay hydrated by drinking plenty of water",
			"Limit processed foods and added sugars",
		}
	case "exercise":
		tips = []string{
			"Aim for at least 150 minutes of moderate exercise per week",
			"Include both cardio and strength training",
			"Start slowly and gradually increase intensity",
			"Find activities you enjoy to stay motivated",
			"Don't forget to warm up and cool down",
		}
	case "sleep":
		tips = []string{
			"Aim for 7-9 hours of sleep per night",
			"Maintain a consistent sleep schedule",
			"Create a relaxing bedtime routine",
			"Keep your bedroom cool, dark, and quiet",
			"Avoid screens before bedtime",
		}
	case "mental_health":
		tips = []string{
			"Practice stress management techniques",
			"Stay connected with friends and family",
			"Take time for activities you enjoy",
			"Consider meditation or mindfulness practices",
			"Don't hesitate to seek professional help when needed",
		}
	default:
		tips = []string{
			"Maintain a balanced diet with variety",
			"Exercise regularly and stay active",
			"Get adequate sleep and rest",
			"Manage stress effectively",
			"Stay hydrated throughout the day",
			"Don't skip regular health checkups",
			"Practice good hygiene",
			"Limit alcohol and avoid smoking",
		}
	}

	// Add age-specific tips
	if age != "" {
		if ageInt, err := strconv.Atoi(age); err == nil {
			if ageInt > 65 {
				tips = append(tips, "Focus on balance exercises to prevent falls")
				tips = append(tips, "Ensure adequate calcium and vitamin D intake")
			} else if ageInt < 18 {
				tips = append(tips, "Focus on growth and development nutrition")
				tips = append(tips, "Limit screen time and encourage outdoor activities")
			}
		}
	}

	// Add gender-specific tips
	if gender == "female" {
		tips = append(tips, "Ensure adequate iron intake, especially during menstruation")
		tips = append(tips, "Consider calcium supplements for bone health")
	}

	return tips
}

func containsInsensitive(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(strings.ToLower(s), strings.ToLower(item)) {
			return true
		}
	}
	return false
}
