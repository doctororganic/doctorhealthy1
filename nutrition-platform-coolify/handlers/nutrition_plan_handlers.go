package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"nutrition-platform/models"
	"nutrition-platform/services"

	"github.com/labstack/echo/v4"
)

// NutritionPlanHandler handles nutrition plan endpoints
type NutritionPlanHandler struct {
	nutritionPlanService *services.NutritionPlanService
	healthService        *services.HealthService
}

// NewNutritionPlanHandler creates a new nutrition plan handler
func NewNutritionPlanHandler(nutritionPlanService *services.NutritionPlanService, healthService *services.HealthService) *NutritionPlanHandler {
	return &NutritionPlanHandler{
		nutritionPlanService: nutritionPlanService,
		healthService:        healthService,
	}
}

// GetNutritionPlanRecommendations provides AI-powered nutrition plan recommendations
func (h *NutritionPlanHandler) GetNutritionPlanRecommendations(c echo.Context) error {
	var req models.HealthAssessmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "NUTRITION_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "NUTRITION_002",
		})
	}

	recommendations, err := h.nutritionPlanService.RecommendNutritionPlan(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "recommendation_failed",
			"message": err.Error(),
			"code":    "NUTRITION_003",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"recommendations": recommendations,
		"total_plans":     len(recommendations),
		"disclaimer":      "These recommendations are based on scientific evidence and general guidelines. Always consult with a healthcare provider or registered dietitian before making significant dietary changes, especially if you have medical conditions or take medications.",
		"generated_at":    "now",
	})
}

// GetQuickNutritionAssessment provides a quick nutrition assessment with plan suggestions
func (h *NutritionPlanHandler) GetQuickNutritionAssessment(c echo.Context) error {
	// Parse query parameters for quick assessment
	age, _ := strconv.Atoi(c.QueryParam("age"))
	height, _ := strconv.ParseFloat(c.QueryParam("height"), 64)
	weight, _ := strconv.ParseFloat(c.QueryParam("weight"), 64)
	gender := c.QueryParam("gender")
	activityLevel := c.QueryParam("activity_level")
	goal := c.QueryParam("goal")

	// Validate required parameters
	if age <= 0 || height <= 0 || weight <= 0 || gender == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_parameters",
			"message": "Age, height, weight, and gender are required",
			"code":    "NUTRITION_004",
		})
	}

	// Create simplified assessment request
	req := models.HealthAssessmentRequest{
		Age:           age,
		Height:        height,
		Weight:        weight,
		Gender:        gender,
		ActivityLevel: activityLevel,
		HealthGoals:   []string{goal},
	}

	// Set defaults if not provided
	if req.ActivityLevel == "" {
		req.ActivityLevel = "moderate"
	}

	recommendations, err := h.nutritionPlanService.RecommendNutritionPlan(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "assessment_failed",
			"message": err.Error(),
			"code":    "NUTRITION_005",
		})
	}

	// Calculate basic metrics for response
	heightM := height / 100
	bmi := weight / (heightM * heightM)

	var bmiCategory string
	if bmi < 18.5 {
		bmiCategory = "Underweight"
	} else if bmi < 25 {
		bmiCategory = "Normal weight"
	} else if bmi < 30 {
		bmiCategory = "Overweight"
	} else {
		bmiCategory = "Obese"
	}

	// Return top 3 recommendations for quick assessment
	topRecommendations := recommendations
	if len(topRecommendations) > 3 {
		topRecommendations = topRecommendations[:3]
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"assessment": map[string]interface{}{
			"bmi":            bmi,
			"bmi_category":   bmiCategory,
			"age_group":      getAgeGroup(age),
			"activity_level": activityLevel,
		},
		"top_recommendations": topRecommendations,
		"quick_tips":          getQuickNutritionTips(bmiCategory, goal),
		"disclaimer":          "This is a simplified assessment. For comprehensive nutrition planning, use the detailed assessment endpoint.",
	})
}

// GetNutritionPlanTypes returns all available nutrition plan types with descriptions
func (h *NutritionPlanHandler) GetNutritionPlanTypes(c echo.Context) error {
	planTypes := make(map[string]interface{})

	for planType, planInfo := range services.NutritionPlanTypes {
		planTypes[planType] = map[string]interface{}{
			"name":                planInfo.Name,
			"description":         planInfo.Description,
			"macro_ratio":         planInfo.MacroRatio,
			"benefits":            planInfo.Benefits,
			"best_for":            planInfo.BestFor,
			"restrictions":        planInfo.Restrictions,
			"duration":            planInfo.Duration,
			"evidence_level":      planInfo.EvidenceLevel,
			"monitoring_required": planInfo.MonitoringRequired,
			"medical_approval":    planInfo.MedicalApproval,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"plan_types": planTypes,
		"total":      len(planTypes),
	})
}

// GetNutritionPlanDetails returns detailed information about a specific plan type
func (h *NutritionPlanHandler) GetNutritionPlanDetails(c echo.Context) error {
	planType := c.Param("plan_type")
	if planType == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "missing_plan_type",
			"message": "Plan type parameter is required",
			"code":    "NUTRITION_006",
		})
	}

	planInfo, exists := services.NutritionPlanTypes[planType]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "plan_not_found",
			"message": "Nutrition plan type not found",
			"code":    "NUTRITION_007",
		})
	}

	// Get scientific evidence for this plan
	nutritionService := services.NewNutritionPlanService(nil) // We only need the method, not DB
	evidence := nutritionService.GetScientificEvidence(planType)

	// Get sample meal ideas for this plan type
	sampleMeals := getSampleMealsForPlan(planType)

	// Get implementation tips
	implementationTips := getImplementationTips(planType)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"plan": map[string]interface{}{
			"type":                planType,
			"name":                planInfo.Name,
			"description":         planInfo.Description,
			"macro_ratio":         planInfo.MacroRatio,
			"benefits":            planInfo.Benefits,
			"best_for":            planInfo.BestFor,
			"restrictions":        planInfo.Restrictions,
			"duration":            planInfo.Duration,
			"evidence_level":      planInfo.EvidenceLevel,
			"monitoring_required": planInfo.MonitoringRequired,
			"medical_approval":    planInfo.MedicalApproval,
			"scientific_evidence": evidence,
			"sample_meals":        sampleMeals,
			"implementation_tips": implementationTips,
		},
	})
}

// CreatePersonalizedNutritionPlan creates a personalized nutrition plan for a user
func (h *NutritionPlanHandler) CreatePersonalizedNutritionPlan(c echo.Context) error {
	var req struct {
		PlanType          string                         `json:"plan_type" validate:"required"`
		HealthAssessment  models.HealthAssessmentRequest `json:"health_assessment" validate:"required"`
		CustomPreferences map[string]interface{}         `json:"custom_preferences,omitempty"`
		DurationWeeks     int                            `json:"duration_weeks" validate:"min=1,max=52"`
		StartDate         string                         `json:"start_date,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "NUTRITION_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "NUTRITION_002",
		})
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "authentication_required",
			"message": "User authentication required",
			"code":    "NUTRITION_008",
		})
	}

	// Validate plan type exists
	planInfo, exists := services.NutritionPlanTypes[req.PlanType]
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_plan_type",
			"message": "Invalid nutrition plan type",
			"code":    "NUTRITION_009",
		})
	}

	// Get recommendations to validate suitability
	recommendations, err := h.nutritionPlanService.RecommendNutritionPlan(&req.HealthAssessment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "recommendation_failed",
			"message": err.Error(),
			"code":    "NUTRITION_003",
		})
	}

	// Check if selected plan is in recommendations
	var selectedRecommendation *services.PlanRecommendation
	for _, rec := range recommendations {
		if rec.PlanType == req.PlanType {
			selectedRecommendation = &rec
			break
		}
	}

	if selectedRecommendation == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":        "plan_not_suitable",
			"message":      "Selected plan type is not recommended for your profile",
			"code":         "NUTRITION_010",
			"alternatives": recommendations[:min(3, len(recommendations))],
		})
	}

	// Create the personalized plan
	plan := createPersonalizedPlan(userID, req.PlanType, planInfo, selectedRecommendation, &req.HealthAssessment, req.DurationWeeks)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success":              true,
		"plan":                 plan,
		"message":              "Personalized nutrition plan created successfully",
		"recommendation_score": selectedRecommendation.Score,
		"confidence":           selectedRecommendation.Confidence,
	})
}

// GetNutritionPlanComparison compares multiple nutrition plans for a user
func (h *NutritionPlanHandler) GetNutritionPlanComparison(c echo.Context) error {
	var req struct {
		PlanTypes        []string                       `json:"plan_types" validate:"required,min=2,max=5"`
		HealthAssessment models.HealthAssessmentRequest `json:"health_assessment" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "invalid_request",
			"message": "Invalid request format",
			"code":    "NUTRITION_001",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
			"code":    "NUTRITION_002",
		})
	}

	// Validate all plan types exist
	for _, planType := range req.PlanTypes {
		if _, exists := services.NutritionPlanTypes[planType]; !exists {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "invalid_plan_type",
				"message": fmt.Sprintf("Invalid nutrition plan type: %s", planType),
				"code":    "NUTRITION_009",
			})
		}
	}

	// Get all recommendations
	allRecommendations, err := h.nutritionPlanService.RecommendNutritionPlan(&req.HealthAssessment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "recommendation_failed",
			"message": err.Error(),
			"code":    "NUTRITION_003",
		})
	}

	// Filter recommendations for requested plan types
	var comparison []map[string]interface{}
	for _, planType := range req.PlanTypes {
		var recommendation *services.PlanRecommendation
		for _, rec := range allRecommendations {
			if rec.PlanType == planType {
				recommendation = &rec
				break
			}
		}

		planInfo := services.NutritionPlanTypes[planType]
		comparisonItem := map[string]interface{}{
			"plan_type":   planType,
			"name":        planInfo.Name,
			"score":       0.0,
			"confidence":  "low",
			"suitability": "not_recommended",
		}

		if recommendation != nil {
			comparisonItem["score"] = recommendation.Score
			comparisonItem["confidence"] = recommendation.Confidence
			comparisonItem["rationale"] = recommendation.Rationale
			comparisonItem["benefits"] = recommendation.Benefits
			comparisonItem["considerations"] = recommendation.Considerations
			comparisonItem["macro_distribution"] = recommendation.MacroDistribution
			comparisonItem["monitoring_required"] = recommendation.MonitoringRequired
			comparisonItem["medical_approval"] = recommendation.MedicalApproval

			if recommendation.Score >= 70 {
				comparisonItem["suitability"] = "highly_suitable"
			} else if recommendation.Score >= 50 {
				comparisonItem["suitability"] = "suitable"
			} else {
				comparisonItem["suitability"] = "not_recommended"
			}
		}

		comparison = append(comparison, comparisonItem)
	}

	// Sort by score
	sort.Slice(comparison, func(i, j int) bool {
		return comparison[i]["score"].(float64) > comparison[j]["score"].(float64)
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"comparison": comparison,
		"best_match": comparison[0]["plan_type"],
		"disclaimer": "Comparison based on your health profile and scientific evidence. Consult healthcare providers before making dietary changes.",
	})
}

// RegisterRoutes registers nutrition plan routes
func (h *NutritionPlanHandler) RegisterRoutes(e *echo.Group) {
	nutrition := e.Group("/nutrition-plans")

	// Plan recommendations and analysis
	nutrition.POST("/recommendations", h.GetNutritionPlanRecommendations)
	nutrition.GET("/quick-assessment", h.GetQuickNutritionAssessment)
	nutrition.POST("/comparison", h.GetNutritionPlanComparison)

	// Plan information
	nutrition.GET("/types", h.GetNutritionPlanTypes)
	nutrition.GET("/types/:plan_type", h.GetNutritionPlanDetails)

	// Personalized plans
	nutrition.POST("/personalized", h.CreatePersonalizedNutritionPlan)
}

// Helper functions

func getAgeGroup(age int) string {
	if age < 18 {
		return "adolescent"
	} else if age < 30 {
		return "young_adult"
	} else if age < 50 {
		return "middle_aged"
	} else if age < 65 {
		return "older_adult"
	} else {
		return "senior"
	}
}

func getQuickNutritionTips(bmiCategory, goal string) []string {
	tips := []string{
		"Focus on whole, unprocessed foods",
		"Stay hydrated with plenty of water",
		"Include a variety of colorful fruits and vegetables",
	}

	switch bmiCategory {
	case "Underweight":
		tips = append(tips, "Increase calorie-dense, nutritious foods", "Consider healthy fats like nuts and avocados")
	case "Overweight", "Obese":
		tips = append(tips, "Focus on portion control", "Increase fiber intake to feel fuller")
	}

	switch goal {
	case "weight_loss":
		tips = append(tips, "Create a moderate calorie deficit", "Combine diet with regular exercise")
	case "muscle_building":
		tips = append(tips, "Ensure adequate protein intake", "Time protein around workouts")
	case "heart_health":
		tips = append(tips, "Limit sodium and saturated fats", "Include omega-3 rich foods")
	}

	return tips
}

func getSampleMealsForPlan(planType string) map[string][]string {
	mealPlans := map[string]map[string][]string{
		"mediterranean": {
			"breakfast": {"Greek yogurt with berries and nuts", "Whole grain toast with olive oil and tomato", "Vegetable omelet with herbs"},
			"lunch":     {"Mediterranean salad with chickpeas", "Grilled fish with vegetables", "Lentil soup with whole grain bread"},
			"dinner":    {"Baked salmon with quinoa", "Chicken with roasted vegetables", "Vegetable pasta with olive oil"},
		},
		"ketogenic": {
			"breakfast": {"Eggs with avocado and bacon", "Keto smoothie with MCT oil", "Cheese and vegetable omelet"},
			"lunch":     {"Chicken salad with olive oil", "Zucchini noodles with meat sauce", "Tuna salad lettuce wraps"},
			"dinner":    {"Steak with buttered broccoli", "Salmon with asparagus", "Pork chops with cauliflower mash"},
		},
		"plant_based": {
			"breakfast": {"Oatmeal with fruit and nuts", "Smoothie bowl with plant protein", "Chia pudding with berries"},
			"lunch":     {"Buddha bowl with quinoa", "Lentil and vegetable curry", "Black bean and sweet potato salad"},
			"dinner":    {"Tofu stir-fry with brown rice", "Chickpea and vegetable stew", "Stuffed bell peppers with quinoa"},
		},
	}

	if meals, exists := mealPlans[planType]; exists {
		return meals
	}

	return map[string][]string{
		"breakfast": {"Balanced breakfast with protein, carbs, and healthy fats"},
		"lunch":     {"Nutrient-dense lunch with vegetables and lean protein"},
		"dinner":    {"Well-balanced dinner following plan guidelines"},
	}
}

func getImplementationTips(planType string) []string {
	tips := map[string][]string{
		"mediterranean": {
			"Start by replacing butter with olive oil",
			"Add fish to your diet 2-3 times per week",
			"Snack on nuts instead of processed foods",
			"Use herbs and spices instead of salt for flavor",
		},
		"ketogenic": {
			"Gradually reduce carbs over 1-2 weeks",
			"Track your macros carefully, especially in the beginning",
			"Stay hydrated and supplement electrolytes",
			"Expect an adaptation period of 2-4 weeks",
		},
		"dash": {
			"Read nutrition labels to monitor sodium content",
			"Gradually increase fruits and vegetables",
			"Choose low-fat dairy products",
			"Limit processed and restaurant foods",
		},
		"plant_based": {
			"Start with one plant-based meal per day",
			"Learn to combine proteins for complete amino acids",
			"Supplement with vitamin B12",
			"Experiment with new plant-based recipes",
		},
	}

	if planTips, exists := tips[planType]; exists {
		return planTips
	}

	return []string{
		"Start gradually and make sustainable changes",
		"Plan your meals in advance",
		"Stay consistent with your eating patterns",
		"Monitor how you feel and adjust as needed",
	}
}

func createPersonalizedPlan(userID, planType string, planInfo services.PlanTypeInfo, recommendation *services.PlanRecommendation, assessment *models.HealthAssessmentRequest, durationWeeks int) map[string]interface{} {
	// Calculate user metrics for personalization
	heightM := assessment.Height / 100
	bmi := assessment.Weight / (heightM * heightM)

	var bmr float64
	if assessment.Gender == "male" {
		bmr = 10*assessment.Weight + 6.25*assessment.Height - 5*float64(assessment.Age) + 5
	} else {
		bmr = 10*assessment.Weight + 6.25*assessment.Height - 5*float64(assessment.Age) - 161
	}

	activityMultipliers := map[string]float64{
		"sedentary":   1.2,
		"light":       1.375,
		"moderate":    1.55,
		"active":      1.725,
		"very_active": 1.9,
	}
	tdee := bmr * activityMultipliers[assessment.ActivityLevel]

	return map[string]interface{}{
		"user_id":               userID,
		"plan_type":             planType,
		"plan_name":             planInfo.Name,
		"duration_weeks":        durationWeeks,
		"daily_calorie_target":  int(tdee),
		"bmi":                   bmi,
		"macro_targets":         recommendation.MacroDistribution,
		"personalized_benefits": recommendation.Benefits,
		"considerations":        recommendation.Considerations,
		"monitoring_required":   recommendation.MonitoringRequired,
		"medical_approval":      recommendation.MedicalApproval,
		"confidence_score":      recommendation.Score,
		"evidence_level":        planInfo.EvidenceLevel,
		"sample_meals":          getSampleMealsForPlan(planType),
		"implementation_tips":   getImplementationTips(planType),
		"created_at":            time.Now(),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
