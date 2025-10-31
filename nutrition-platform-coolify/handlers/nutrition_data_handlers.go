package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// NutritionDataHandler handles all nutrition-related data endpoints
type NutritionDataHandler struct {
	DataPath string
}

// NewNutritionDataHandler creates a new nutrition data handler
func NewNutritionDataHandler(dataPath string) *NutritionDataHandler {
	return &NutritionDataHandler{
		DataPath: dataPath,
	}
}

// MetabolismResponse represents the metabolism guide response
type MetabolismResponse struct {
	MetabolismGuide interface{} `json:"metabolism_guide"`
	Status          string      `json:"status"`
	Message         string      `json:"message"`
}

// WorkoutTechniquesResponse represents workout techniques response
type WorkoutTechniquesResponse struct {
	WorkoutTechniques []interface{} `json:"workout_techniques"`
	Status            string        `json:"status"`
	Message           string        `json:"message"`
}

// MealPlansResponse represents meal plans response
type MealPlansResponse struct {
	MealPlans interface{} `json:"meal_plans"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
}

// VitaminsResponse represents vitamins and minerals response
type VitaminsResponse struct {
	Nutrients []interface{} `json:"nutrients"`
	Status    string        `json:"status"`
	Message   string        `json:"message"`
}

// DrugsNutritionResponse represents drugs affecting nutrition response
type DrugsNutritionResponse struct {
	NutritionalRecommendations interface{} `json:"nutritional_recommendations"`
	Status                     string      `json:"status"`
	Message                    string      `json:"message"`
}

// DiseaseResponse represents disease data response
type DiseaseResponse struct {
	HealthConditions []interface{} `json:"health_conditions"`
	Status           string        `json:"status"`
	Message          string        `json:"message"`
}

// GetMetabolismGuide returns metabolism guide data
// @Summary Get metabolism guide
// @Description Retrieve comprehensive metabolism guide with nutrition and exercise information
// @Tags Metabolism
// @Accept json
// @Produce json
// @Param lang query string false "Language (en/ar)" default(en)
// @Success 200 {object} MetabolismResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/metabolism [get]
func (h *NutritionDataHandler) GetMetabolismGuide(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	filePath := filepath.Join(h.DataPath, "metabolism.js")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Metabolism data not found",
			"code":  "METABOLISM_NOT_FOUND",
		})
	}

	var metabolismData map[string]interface{}
	if err := json.Unmarshal(data, &metabolismData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse metabolism data",
			"code":  "PARSE_ERROR",
		})
	}

	response := MetabolismResponse{
		MetabolismGuide: metabolismData["metabolism_guide"],
		Status:          "success",
		Message:         fmt.Sprintf("Metabolism guide retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// GetWorkoutTechniques returns workout techniques data
// @Summary Get workout techniques
// @Description Retrieve workout techniques with detailed instructions and benefits
// @Tags Workouts
// @Accept json
// @Produce json
// @Param lang query string false "Language (en/ar)" default(en)
// @Param technique query string false "Specific technique name"
// @Success 200 {object} WorkoutTechniquesResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/workout-techniques [get]
func (h *NutritionDataHandler) GetWorkoutTechniques(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	technique := c.QueryParam("technique")

	filePath := filepath.Join(h.DataPath, "workouts teq.js")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Workout techniques data not found",
			"code":  "WORKOUT_TECHNIQUES_NOT_FOUND",
		})
	}

	// Extract JSON from the file (skip the comment lines)
	jsonStart := strings.Index(string(data), "{")
	if jsonStart == -1 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid workout techniques data format",
			"code":  "INVALID_FORMAT",
		})
	}

	jsonData := data[jsonStart:]
	var workoutData map[string]interface{}
	if err := json.Unmarshal(jsonData, &workoutData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse workout techniques data",
			"code":  "PARSE_ERROR",
		})
	}

	workoutTechniques := workoutData["workout_techniques"].([]interface{})

	// Filter by specific technique if requested
	if technique != "" {
		for _, tech := range workoutTechniques {
			techMap := tech.(map[string]interface{})
			if name, ok := techMap["name"].(map[string]interface{}); ok {
				if name[lang] == technique {
					workoutTechniques = []interface{}{tech}
					break
				}
			}
		}
	}

	response := WorkoutTechniquesResponse{
		WorkoutTechniques: workoutTechniques,
		Status:            "success",
		Message:           fmt.Sprintf("Workout techniques retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// GetMealPlans returns meal plans data
// @Summary Get meal plans
// @Description Retrieve comprehensive meal plans for different fitness goals
// @Tags Meals
// @Accept json
// @Produce json
// @Param lang query string false "Language (en/ar)" default(en)
// @Param type query string false "Workout type (shred, bulk, etc.)"
// @Param week query int false "Week number (1-4)"
// @Success 200 {object} MealPlansResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/meal-plans [get]
func (h *NutritionDataHandler) GetMealPlans(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	workoutType := c.QueryParam("type")
	weekStr := c.QueryParam("week")

	filePath := filepath.Join(h.DataPath, "meals plans.js")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Meal plans data not found",
			"code":  "MEAL_PLANS_NOT_FOUND",
		})
	}

	// Extract JSON from the file
	jsonStart := strings.Index(string(data), "{")
	if jsonStart == -1 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid meal plans data format",
			"code":  "INVALID_FORMAT",
		})
	}

	jsonData := data[jsonStart:]
	var mealData map[string]interface{}
	if err := json.Unmarshal(jsonData, &mealData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse meal plans data",
			"code":  "PARSE_ERROR",
		})
	}

	// Filter by workout type and week if specified
	filteredData := mealData
	if workoutType != "" {
		if typeData, ok := mealData["workoutType"].(map[string]interface{}); ok {
			if typeData[lang] != workoutType {
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": "Workout type not found",
					"code":  "WORKOUT_TYPE_NOT_FOUND",
				})
			}
		}
	}

	if weekStr != "" {
		week, err := strconv.Atoi(weekStr)
		if err != nil || week < 1 || week > 4 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid week number (1-4)",
				"code":  "INVALID_WEEK",
			})
		}

		if weeks, ok := mealData["weeks"].([]interface{}); ok {
			for _, w := range weeks {
				weekMap := w.(map[string]interface{})
				if int(weekMap["week"].(float64)) == week {
					filteredData = map[string]interface{}{
						"workoutType": mealData["workoutType"],
						"week":        w,
					}
					break
				}
			}
		}
	}

	response := MealPlansResponse{
		MealPlans: filteredData,
		Status:    "success",
		Message:   fmt.Sprintf("Meal plans retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// GetVitaminsAndMinerals returns vitamins and minerals data
// @Summary Get vitamins and minerals
// @Description Retrieve comprehensive vitamins and minerals information with dosages
// @Tags Nutrition
// @Accept json
// @Produce json
// @Param lang query string false "Language (en/ar)" default(en)
// @Param nutrient query string false "Specific nutrient name"
// @Param population query string false "Target population (normal adult, bodybuilder, etc.)"
// @Success 200 {object} VitaminsResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/vitamins-minerals [get]
func (h *NutritionDataHandler) GetVitaminsAndMinerals(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	nutrient := c.QueryParam("nutrient")
	population := c.QueryParam("population")

	filePath := filepath.Join(h.DataPath, "vitamins and minerals.js")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Vitamins and minerals data not found",
			"code":  "VITAMINS_NOT_FOUND",
		})
	}

	var vitaminData map[string]interface{}
	if err := json.Unmarshal(data, &vitaminData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse vitamins data",
			"code":  "PARSE_ERROR",
		})
	}

	nutrients := vitaminData["nutrients"].([]interface{})

	// Filter by specific nutrient if requested
	if nutrient != "" {
		for _, nut := range nutrients {
			nutMap := nut.(map[string]interface{})
			if name, ok := nutMap["name"].(map[string]interface{}); ok {
				if strings.Contains(strings.ToLower(name[lang].(string)), strings.ToLower(nutrient)) {
					nutrients = []interface{}{nut}
					break
				}
			}
		}
	}

	// Filter by population if requested
	if population != "" {
		filteredNutrients := []interface{}{}
		for _, nut := range nutrients {
			nutMap := nut.(map[string]interface{})
			if dosage, ok := nutMap["dosage"].([]interface{}); ok {
				for _, dose := range dosage {
					doseMap := dose.(map[string]interface{})
					if pop, ok := doseMap["population"].(map[string]interface{}); ok {
						if strings.Contains(strings.ToLower(pop[lang].(string)), strings.ToLower(population)) {
							filteredNutrients = append(filteredNutrients, nut)
							break
						}
					}
				}
			}
		}
		nutrients = filteredNutrients
	}

	response := VitaminsResponse{
		Nutrients: nutrients,
		Status:    "success",
		Message:   fmt.Sprintf("Vitamins and minerals retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// GetDrugsNutritionInteractions returns drugs affecting nutrition data
// @Summary Get drugs nutrition interactions
// @Description Retrieve information about how drugs affect nutrition and recommendations
// @Tags Nutrition
// @Accept json
// @Produce json
// @Param lang query string false "Language (en/ar)" default(en)
// @Success 200 {object} DrugsNutritionResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/drugs-nutrition [get]
func (h *NutritionDataHandler) GetDrugsNutritionInteractions(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	filePath := filepath.Join(h.DataPath, "drugs affect nutrition.js")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Drugs nutrition data not found",
			"code":  "DRUGS_NUTRITION_NOT_FOUND",
		})
	}

	var drugsData map[string]interface{}
	if err := json.Unmarshal(data, &drugsData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse drugs nutrition data",
			"code":  "PARSE_ERROR",
		})
	}

	response := DrugsNutritionResponse{
		NutritionalRecommendations: drugsData["NutritionalRecommendations"],
		Status:                     "success",
		Message:                    fmt.Sprintf("Drugs nutrition interactions retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// GetDiseaseData returns disease and health conditions data
// @Summary Get disease data
// @Description Retrieve comprehensive disease and health conditions information
// @Tags Health
// @Accept json
// @Produce json
// @Param lang query string false "Language (en/ar)" default(en)
// @Param condition query string false "Specific health condition"
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} DiseaseResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/diseases [get]
func (h *NutritionDataHandler) GetDiseaseData(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	condition := c.QueryParam("condition")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	filePath := filepath.Join(h.DataPath, "old disease.js")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Disease data not found",
			"code":  "DISEASE_DATA_NOT_FOUND",
		})
	}

	// Extract JSON from the file (skip the initial text)
	jsonStart := strings.Index(string(data), "[")
	if jsonStart == -1 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid disease data format",
			"code":  "INVALID_FORMAT",
		})
	}

	jsonData := data[jsonStart:]
	var healthConditions []interface{}
	if err := json.Unmarshal(jsonData, &healthConditions); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse disease data",
			"code":  "PARSE_ERROR",
		})
	}

	// Filter by specific condition if requested
	if condition != "" {
		filteredConditions := []interface{}{}
		for _, cond := range healthConditions {
			condMap := cond.(map[string]interface{})
			if name, ok := condMap["name"].(map[string]interface{}); ok {
				if strings.Contains(strings.ToLower(name[lang].(string)), strings.ToLower(condition)) {
					filteredConditions = append(filteredConditions, cond)
				}
			}
		}
		healthConditions = filteredConditions
	}

	// Apply pagination
	total := len(healthConditions)
	start := offset
	end := offset + limit

	if start >= total {
		healthConditions = []interface{}{}
	} else {
		if end > total {
			end = total
		}
		healthConditions = healthConditions[start:end]
	}

	response := DiseaseResponse{
		HealthConditions: healthConditions,
		Status:           "success",
		Message:          fmt.Sprintf("Disease data retrieved successfully in %s (showing %d-%d of %d)", lang, start+1, end, total),
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers all nutrition data routes
func (h *NutritionDataHandler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	// Metabolism routes
	api.GET("/metabolism", h.GetMetabolismGuide)

	// Workout techniques routes
	api.GET("/workout-techniques", h.GetWorkoutTechniques)

	// Meal plans routes
	api.GET("/meal-plans", h.GetMealPlans)

	// Vitamins and minerals routes
	api.GET("/vitamins-minerals", h.GetVitaminsAndMinerals)

	// Drugs nutrition interactions routes
	api.GET("/drugs-nutrition", h.GetDrugsNutritionInteractions)

	// Disease data routes
	api.GET("/diseases", h.GetDiseaseData)

	// Calories data routes
	api.GET("/calories", GetCalories)
	api.GET("/calories/:category", GetCaloriesByCategory)

	// Skills data routes
	api.GET("/skills", GetSkills)
	api.GET("/skills/:difficulty", GetSkillsByDifficulty)

	// Type plans data routes
	api.GET("/type-plans", GetTypePlans)
	api.GET("/type-plans/:type", GetTypePlansByType)
}
