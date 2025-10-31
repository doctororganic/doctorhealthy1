package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// NutritionRequest represents the request structure for nutrition analysis
type NutritionRequest struct {
	Food       string  `json:"food" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required,min=0"`
	Unit       string  `json:"unit" validate:"required"`
	CheckHalal bool    `json:"checkHalal"`
	Language   string  `json:"language"`
}

// NutritionResponse represents the response structure for nutrition analysis
type NutritionResponse struct {
	Food              string   `json:"food"`
	Quantity          float64  `json:"quantity"`
	Unit              string   `json:"unit"`
	Calories          int      `json:"calories"`
	Protein           float64  `json:"protein"`
	Carbohydrates     float64  `json:"carbohydrates"`
	Fat               float64  `json:"fat"`
	Fiber             float64  `json:"fiber"`
	HalalStatus       *bool    `json:"halalStatus,omitempty"`
	Recommendations   []string `json:"recommendations"`
	MedicalDisclaimer string   `json:"medicalDisclaimer"`
}

// AnalyzeNutrition handles nutrition analysis requests (exported function)
func AnalyzeNutrition(c echo.Context) error {
	var req NutritionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Additional validation
	if req.Quantity <= 0 || req.Quantity > 10000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Quantity must be between 0 and 10000"})
	}

	if len(req.Food) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Food name must be less than 100 characters"})
	}

	// Generate nutrition analysis
	response := generateNutritionAnalysis(req)

	return c.JSON(http.StatusOK, response)
}

// generateNutritionAnalysis generates nutrition data for a food item
func generateNutritionAnalysis(req NutritionRequest) NutritionResponse {
	// Mock nutrition data based on common foods
	baseNutrition := map[string]map[string]float64{
		"chicken": {"calories": 165, "protein": 31, "carbs": 0, "fat": 3.6, "fiber": 0},
		"rice":    {"calories": 130, "protein": 2.7, "carbs": 28, "fat": 0.3, "fiber": 0.4},
		"apple":   {"calories": 52, "protein": 0.3, "carbs": 14, "fat": 0.2, "fiber": 2.4},
		"bread":   {"calories": 265, "protein": 9, "carbs": 49, "fat": 3.2, "fiber": 2.7},
		"egg":     {"calories": 155, "protein": 13, "carbs": 1.1, "fat": 11, "fiber": 0},
		"milk":    {"calories": 42, "protein": 3.4, "carbs": 5, "fat": 1, "fiber": 0},
		"fish":    {"calories": 206, "protein": 22, "carbs": 0, "fat": 12, "fiber": 0},
	}

	// Default values
	nutrition := map[string]float64{"calories": 100, "protein": 5, "carbs": 15, "fat": 3, "fiber": 2}

	// Find matching food
	foodLower := strings.ToLower(req.Food)
	for food, values := range baseNutrition {
		if strings.Contains(foodLower, food) {
			nutrition = values
			break
		}
	}

	// Scale by quantity (assuming per 100g base)
	scale := req.Quantity / 100

	// Check halal status if requested
	var halalStatus *bool
	if req.CheckHalal {
		isHalal := checkHalalStatus(req.Food)
		halalStatus = &isHalal
	}

	// Generate recommendations
	recommendations := generateRecommendations(req.Food, nutrition)

	return NutritionResponse{
		Food:              req.Food,
		Quantity:          req.Quantity,
		Unit:              req.Unit,
		Calories:          int(nutrition["calories"] * scale),
		Protein:           nutrition["protein"] * scale,
		Carbohydrates:     nutrition["carbs"] * scale,
		Fat:               nutrition["fat"] * scale,
		Fiber:             nutrition["fiber"] * scale,
		HalalStatus:       halalStatus,
		Recommendations:   recommendations,
		MedicalDisclaimer: "This nutritional information is for educational purposes only and should not replace professional medical advice. Consult with a healthcare provider for personalized dietary recommendations.",
	}
}

// checkHalalStatus checks if a food is halal
func checkHalalStatus(food string) bool {
	nonHalalKeywords := []string{"pork", "ham", "bacon", "wine", "beer", "alcohol", "gelatin"}
	foodLower := strings.ToLower(food)

	for _, keyword := range nonHalalKeywords {
		if strings.Contains(foodLower, keyword) {
			return false
		}
	}
	return true
}

// generateRecommendations generates food recommendations
func generateRecommendations(food string, nutrition map[string]float64) []string {
	recommendations := []string{
		"Maintain a balanced diet with variety of foods",
		"Consider portion sizes appropriate for your daily needs",
	}

	// Add specific recommendations based on nutrition profile
	if nutrition["protein"] > 20 {
		recommendations = append(recommendations, "Excellent source of protein for muscle health")
	}
	if nutrition["fiber"] > 3 {
		recommendations = append(recommendations, "Good source of fiber for digestive health")
	}
	if nutrition["calories"] > 200 {
		recommendations = append(recommendations, "High calorie food - consider portion control")
	}

	return recommendations
}