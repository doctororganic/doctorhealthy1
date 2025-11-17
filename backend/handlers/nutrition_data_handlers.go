package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

// MetabolismResponse represents metabolism guide response
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
func (h *NutritionDataHandler) GetWorkoutTechniques(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	technique := c.QueryParam("technique")

	filePath := filepath.Join(h.DataPath, "techniques.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read workout techniques file",
			"code":  "FILE_READ_ERROR",
		})
	}
	
	var workoutData map[string]interface{}
	if err := json.Unmarshal(data, &workoutData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse workout techniques data",
			"code": "PARSE_ERROR",
		})
	}
	
	workoutTechniques := workoutData["workout_techniques"].([]interface{})
	
	// Filter by specific technique if requested
	if technique != "" {
		for _, tech := range workoutTechniques {
			techMap := tech.(map[string]interface{})
			if name, ok := techMap["name"].(map[string]interface{}); ok && name[lang] == technique {
				workoutTechniques = []interface{}{tech}
			}
		}
	}
	
	response := WorkoutTechniquesResponse{
		WorkoutTechniques: workoutTechniques,
		Status:            "success",
		Message:         fmt.Sprintf("Workout techniques retrieved successfully in %s", lang),
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetMealPlans returns meal plans data
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

	// Extract JSON from file
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
func (h *NutritionDataHandler) GetVitaminsAndMinerals(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}
	nutrient := c.QueryParam("nutrient")
	population := c.QueryParam("population")

	filePath := filepath.Join(h.DataPath, "vitamins.json")
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

	filePath := filepath.Join(h.DataPath, "diseases.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Disease data not found",
			"code":  "DISEASE_DATA_NOT_FOUND",
		})
	}

	var diseaseData map[string]interface{}
	if err := json.Unmarshal(data, &diseaseData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse disease data",
			"code":  "PARSE_ERROR",
		})
	}
	
	var healthConditions []interface{}
	if conditions, ok := diseaseData["health_conditions"].([]interface{}); ok {
		healthConditions = conditions
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid disease data format",
			"code":  "INVALID_FORMAT",
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

// MealPlanGenerateRequest represents the request for generating a meal plan
type MealPlanGenerateRequest struct {
	Calories            int      `json:"calories" validate:"required,min=1000,max=5000"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
	Goals               string   `json:"goals"`
	MealsPerDay         int      `json:"meals_per_day"`
	PlanDuration        int      `json:"plan_duration"`
}

// MealPlanGenerateResponse represents the response for generated meal plan
type MealPlanGenerateResponse struct {
	PlanID      string                 `json:"plan_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Calories    int                    `json:"calories"`
	Goals       string                 `json:"goals"`
	Duration    int                    `json:"duration"`
	Meals       []map[string]interface{} `json:"meals"`
	GeneratedAt time.Time              `json:"generated_at"`
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
}

// GenerateMealPlan generates a personalized meal plan
func (h *NutritionDataHandler) GenerateMealPlan(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	var req MealPlanGenerateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
	}

	// Set defaults if not provided
	if req.MealsPerDay == 0 {
		req.MealsPerDay = 3
	}
	if req.PlanDuration == 0 {
		req.PlanDuration = 7
	}

	// Generate unique plan ID
	planID := fmt.Sprintf("meal_plan_%d", time.Now().Unix())
	
	// Create seed for randomness
	rand.Seed(time.Now().UnixNano())

	// Generate meal plan based on goals and calories
	mealPlan := h.generateMealPlanContent(req, lang)

	response := MealPlanGenerateResponse{
		PlanID:      planID,
		Name:        fmt.Sprintf("%d-Day Personalized Meal Plan", req.PlanDuration),
		Description: fmt.Sprintf("Personalized %d-calorie meal plan for %s", req.Calories, req.Goals),
		Calories:    req.Calories,
		Goals:       req.Goals,
		Duration:    req.PlanDuration,
		Meals:       mealPlan,
		GeneratedAt: time.Now(),
		Status:      "success",
		Message:     fmt.Sprintf("Meal plan generated successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// generateMealPlanContent creates the actual meal plan content
func (h *NutritionDataHandler) generateMealPlanContent(req MealPlanGenerateRequest, lang string) []map[string]interface{} {
	meals := []map[string]interface{}{}
	
	// Calculate calories per meal
	caloriesPerMeal := req.Calories / req.MealsPerDay
	
	// Meal templates based on goals
	mealTemplates := h.getMealTemplates(req.Goals, caloriesPerMeal, lang)
	
	// Generate meals for each day
	for day := 1; day <= req.PlanDuration; day++ {
		dayMeals := []map[string]interface{}{}
		
		for mealNum := 1; mealNum <= req.MealsPerDay; mealNum++ {
			// Select random meal template
			template := mealTemplates[rand.Intn(len(mealTemplates))]
			
			meal := map[string]interface{}{
				"day":         day,
				"meal_number": mealNum,
				"meal_type":   h.getMealType(mealNum, lang),
				"name":        template["name"],
				"description": template["description"],
				"calories":    template["calories"],
				"protein":     template["protein"],
				"carbs":       template["carbs"],
				"fat":         template["fat"],
				"ingredients": template["ingredients"],
				"instructions": template["instructions"],
				"prep_time":   template["prep_time"],
				"cook_time":   template["cook_time"],
			}
			
			dayMeals = append(dayMeals, meal)
		}
		
		meals = append(meals, map[string]interface{}{
			"day":   day,
			"date":  time.Now().AddDate(0, 0, day-1).Format("2006-01-02"),
			"meals": dayMeals,
		})
	}
	
	return meals
}

// getMealTemplates returns meal templates based on goals
func (h *NutritionDataHandler) getMealTemplates(goals string, caloriesPerMeal int, lang string) []map[string]interface{} {
	
	// Base templates
	baseTemplates := []map[string]interface{}{
		{
			"name":        h.getMealName("protein_bowl", lang),
			"description": h.getMealDescription("protein_bowl", lang),
			"calories":    caloriesPerMeal,
			"protein":     float64(caloriesPerMeal) * 0.3 / 4.0,
			"carbs":       float64(caloriesPerMeal) * 0.4 / 4.0,
			"fat":         float64(caloriesPerMeal) * 0.3 / 9.0,
			"ingredients": []string{"Chicken breast", "Quinoa", "Mixed vegetables", "Olive oil"},
			"instructions": h.getInstructions("grill_chicken", lang),
			"prep_time":   "15 min",
			"cook_time":   "20 min",
		},
		{
			"name":        h.getMealName("salad_bowl", lang),
			"description": h.getMealDescription("salad_bowl", lang),
			"calories":    caloriesPerMeal,
			"protein":     float64(caloriesPerMeal) * 0.25 / 4.0,
			"carbs":       float64(caloriesPerMeal) * 0.5 / 4.0,
			"fat":         float64(caloriesPerMeal) * 0.25 / 9.0,
			"ingredients": []string{"Mixed greens", "Chickpeas", "Vegetables", "Lemon dressing"},
			"instructions": h.getInstructions("assemble_salad", lang),
			"prep_time":   "10 min",
			"cook_time":   "0 min",
		},
		{
			"name":        h.getMealName("smoothie", lang),
			"description": h.getMealDescription("smoothie", lang),
			"calories":    caloriesPerMeal,
			"protein":     float64(caloriesPerMeal) * 0.2 / 4.0,
			"carbs":       float64(caloriesPerMeal) * 0.6 / 4.0,
			"fat":         float64(caloriesPerMeal) * 0.2 / 9.0,
			"ingredients": []string{"Protein powder", "Banana", "Berries", "Almond milk"},
			"instructions": h.getInstructions("blend_smoothie", lang),
			"prep_time":   "5 min",
			"cook_time":   "0 min",
		},
		{
			"name":        h.getMealName("stir_fry", lang),
			"description": h.getMealDescription("stir_fry", lang),
			"calories":    caloriesPerMeal,
			"protein":     float64(caloriesPerMeal) * 0.3 / 4.0,
			"carbs":       float64(caloriesPerMeal) * 0.4 / 4.0,
			"fat":         float64(caloriesPerMeal) * 0.3 / 9.0,
			"ingredients": []string{"Tofu", "Brown rice", "Mixed vegetables", "Soy sauce"},
			"instructions": h.getInstructions("stir_fry", lang),
			"prep_time":   "15 min",
			"cook_time":   "15 min",
		},
	}
	
	// Adjust templates based on goals
	if goals == "weight_loss" {
		// Lower carb, higher protein for weight loss
		for _, template := range baseTemplates {
			template["protein"] = template["protein"].(float64) * 1.2
			template["carbs"] = template["carbs"].(float64) * 0.8
		}
	} else if goals == "muscle_gain" {
		// Higher protein for muscle gain
		for _, template := range baseTemplates {
			template["protein"] = template["protein"].(float64) * 1.4
		}
	}
	
	return baseTemplates
}

// Helper functions for meal names and descriptions
func (h *NutritionDataHandler) getMealName(mealType, lang string) string {
	names := map[string]map[string]string{
		"protein_bowl": {
			"en": "Protein Power Bowl",
			"ar": "طبق البروتين القوي",
		},
		"salad_bowl": {
			"en": "Garden Fresh Salad Bowl",
			"ar": "سلطة الحديقة الطازجة",
		},
		"smoothie": {
			"en": "Nutrient-Dense Smoothie",
			"ar": "سموذي غني بالمغذيات",
		},
		"stir_fry": {
			"en": "Vegetable Stir Fry",
			"ar": "خضروات مقلية",
		},
	}
	
	if name, ok := names[mealType][lang]; ok {
		return name
	}
	return names[mealType]["en"]
}

func (h *NutritionDataHandler) getMealDescription(mealType, lang string) string {
	descriptions := map[string]map[string]string{
		"protein_bowl": {
			"en": "High-protein meal with balanced nutrients",
			"ar": "وجبة عالية البروتين مع مغذيات متوازنة",
		},
		"salad_bowl": {
			"en": "Fresh and light salad perfect for any time",
			"ar": "سلطة طازجة وخفيفة مناسبة لأي وقت",
		},
		"smoothie": {
			"en": "Quick and nutritious blended meal",
			"ar": "وجبة سريعة ومغذية مختلطة",
		},
		"stir_fry": {
			"en": "Flavorful vegetable dish with complete protein",
			"ar": "طبق خضروات لذيذ مع بروتين كامل",
		},
	}
	
	if desc, ok := descriptions[mealType][lang]; ok {
		return desc
	}
	return descriptions[mealType]["en"]
}

func (h *NutritionDataHandler) getInstructions(cookingType, lang string) string {
	instructions := map[string]map[string]string{
		"grill_chicken": {
			"en": "Season chicken and grill until cooked through. Serve with quinoa and steamed vegetables.",
			"ar": "تبل الدجاج واشويه حتى ينضج تمامًا. قدمه مع الكينوا والخضار المطهو على البخار.",
		},
		"assemble_salad": {
			"en": "Mix all ingredients in a bowl and toss with dressing. Serve immediately.",
			"ar": "اخلط جميع المكونات في وعاء ورشها بالصلصة. قدمها فورًا.",
		},
		"blend_smoothie": {
			"en": "Blend all ingredients until smooth. Add ice if desired and blend again.",
			"ar": "اخلط جميع المكونات حتى تصبح ناعمة. أضف الثلج إذا رغبت واخلط مرة أخرى.",
		},
		"stir_fry": {
			"en": "Heat oil in pan, stir-fry vegetables and tofu. Add cooked rice and sauce, toss to combine.",
			"ar": "سخن الزيت في مقلاة، قلّب الخضار والتوفو. أضف الأرز المطبوخ والصلصة، قلب للمزج.",
		},
	}
	
	if inst, ok := instructions[cookingType][lang]; ok {
		return inst
	}
	return instructions[cookingType]["en"]
}

func (h *NutritionDataHandler) getMealType(mealNum int, lang string) string {
	mealTypes := map[int]map[string]string{
		1: {"en": "Breakfast", "ar": "الإفطار"},
		2: {"en": "Lunch", "ar": "الغداء"},
		3: {"en": "Dinner", "ar": "العشاء"},
		4: {"en": "Snack", "ar": "وجبة خفيفة"},
		5: {"en": "Post-Workout", "ar": "بعد التمرين"},
		6: {"en": "Pre-Workout", "ar": "قبل التمرين"},
	}
	
	if mealType, ok := mealTypes[mealNum][lang]; ok {
		return mealType
	}
	return mealTypes[mealNum]["en"]
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
	api.POST("/meal-plans/generate", h.GenerateMealPlan)

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

// GetCalories returns calories data
func GetCalories(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock calories data
	caloriesData := map[string]interface{}{
		"categories": []map[string]interface{}{
			{
				"name": map[string]string{
					"en": "Low Calorie",
					"ar": "منخفض السعرات الحرارية",
				},
				"description": map[string]string{
					"en": "Meals under 400 calories",
					"ar": "وجبات تحت 400 سعرة حرارية",
				},
				"range": "0-400",
			},
			{
				"name": map[string]string{
					"en": "Medium Calorie",
					"ar": "متوسط السعرات الحرارية",
				},
				"description": map[string]string{
					"en": "Meals between 400-800 calories",
					"ar": "وجبات بين 400-800 سعرة حرارية",
				},
				"range": "400-800",
			},
		},
		"status":  "success",
		"message": fmt.Sprintf("Calories data retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, caloriesData)
}

// GetCaloriesByCategory returns calories data by category
func GetCaloriesByCategory(c echo.Context) error {
	category := c.Param("category")
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock category-specific data
	categoryData := map[string]interface{}{
		"category": category,
		"meals": []map[string]interface{}{
			{
				"id":       "meal_001",
				"name":     "Healthy Salad",
				"calories": 350,
				"category": category,
			},
		},
		"status":  "success",
		"message": fmt.Sprintf("Calories data for category '%s' retrieved successfully in %s", category, lang),
	}

	return c.JSON(http.StatusOK, categoryData)
}

// GetSkills returns skills data
func GetSkills(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock skills data
	skillsData := map[string]interface{}{
		"skills": []map[string]interface{}{
			{
				"id":          "skill_001",
				"name":        "Meal Planning",
				"difficulty":  "beginner",
				"description": "Learn to plan healthy meals",
			},
		},
		"status":  "success",
		"message": fmt.Sprintf("Skills data retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, skillsData)
}

// GetSkillsByDifficulty returns skills data by difficulty
func GetSkillsByDifficulty(c echo.Context) error {
	difficulty := c.Param("difficulty")
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock difficulty-specific data
	skillsData := map[string]interface{}{
		"difficulty": difficulty,
		"skills": []map[string]interface{}{
			{
				"id":          "skill_001",
				"name":        "Meal Planning",
				"difficulty":  difficulty,
				"description": "Learn to plan healthy meals",
			},
		},
		"status":  "success",
		"message": fmt.Sprintf("Skills data for difficulty '%s' retrieved successfully in %s", difficulty, lang),
	}

	return c.JSON(http.StatusOK, skillsData)
}

// GetTypePlans returns type plans data
func GetTypePlans(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock type plans data
	typePlansData := map[string]interface{}{
		"plans": []map[string]interface{}{
			{
				"id":          "plan_001",
				"name":        "Weight Loss Plan",
				"type":        "weight_loss",
				"description": "Comprehensive weight loss program",
			},
		},
		"status":  "success",
		"message": fmt.Sprintf("Type plans data retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, typePlansData)
}

// GetTypePlansByType returns type plans data by type
func GetTypePlansByType(c echo.Context) error {
	planType := c.Param("type")
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock type-specific data
	typePlansData := map[string]interface{}{
		"type": planType,
		"plans": []map[string]interface{}{
			{
				"id":          "plan_001",
				"name":        "Weight Loss Plan",
				"type":        planType,
				"description": "Comprehensive weight loss program",
			},
		},
		"status":  "success",
		"message": fmt.Sprintf("Type plans data for type '%s' retrieved successfully in %s", planType, lang),
	}

	return c.JSON(http.StatusOK, typePlansData)
}
// GenerateWorkoutPlan generates a personalized workout plan
func (h *NutritionDataHandler) GenerateWorkoutPlan(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	duration := c.QueryParam("duration")
	level := c.QueryParam("level")

	// Set defaults
	if duration == "" {
		duration = "4"
	}
	if level == "" {
		level = "intermediate"
	}

	durationInt, _ := strconv.Atoi(duration)
	if durationInt == 0 {
		durationInt = 4
	}

	// Generate unique plan ID
	planID := fmt.Sprintf("workout_plan_%d", time.Now().Unix())

	// Create seed for randomness
	rand.Seed(time.Now().UnixNano())

	// Generate workout plan
	workoutPlan := h.generateWorkoutPlanContent(durationInt, level, lang)

	response := map[string]interface{}{
		"plan_id":     planID,
		"name":        fmt.Sprintf("%d-Day Personalized Workout Plan", durationInt),
		"description": fmt.Sprintf("Personalized %s workout plan for %d days", level, durationInt),
		"duration":    durationInt,
		"level":       level,
		"workouts":    workoutPlan,
		"generated_at": time.Now(),
		"status":      "success",
		"message":     fmt.Sprintf("Workout plan generated successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// generateWorkoutPlanContent creates the actual workout plan content
func (h *NutritionDataHandler) generateWorkoutPlanContent(duration int, level, lang string) []map[string]interface{} {
	workouts := []map[string]interface{}{}

	// Workout templates based on level
	templates := h.getWorkoutTemplates(level, lang)

	// Generate workouts for each day
	for day := 1; day <= duration; day++ {
		// Select random template
		template := templates[rand.Intn(len(templates))]

		workout := map[string]interface{}{
			"day":          day,
			"name":         template["name"],
			"description":  template["description"],
			"focus":        template["focus"],
			"exercises":    template["exercises"],
			"duration":     template["duration"],
			"difficulty":   level,
			"equipment":    template["equipment"],
		}

		workouts = append(workouts, workout)
	}

	return workouts
}

// getWorkoutTemplates returns workout templates based on level
func (h *NutritionDataHandler) getWorkoutTemplates(level, lang string) []map[string]interface{} {
	templates := []map[string]interface{}{
		{
			"name":        h.getWorkoutName("upper_body", lang),
			"description": h.getWorkoutDescription("upper_body", lang),
			"focus":       "Upper Body Strength",
			"exercises": []map[string]interface{}{
				{"name": "Push-ups", "sets": 3, "reps": 12, "rest": "60s"},
				{"name": "Pull-ups", "sets": 3, "reps": 8, "rest": "90s"},
				{"name": "Shoulder Press", "sets": 3, "reps": 10, "rest": "60s"},
			},
			"duration":  "45 minutes",
			"equipment": []string{"dumbbells", "pull-up bar"},
		},
		{
			"name":        h.getWorkoutName("lower_body", lang),
			"description": h.getWorkoutDescription("lower_body", lang),
			"focus":       "Lower Body Strength",
			"exercises": []map[string]interface{}{
				{"name": "Squats", "sets": 4, "reps": 12, "rest": "90s"},
				{"name": "Lunges", "sets": 3, "reps": 10, "rest": "60s"},
				{"name": "Calf Raises", "sets": 3, "reps": 15, "rest": "45s"},
			},
			"duration":  "50 minutes",
			"equipment": []string{"dumbbells", "squat rack"},
		},
		{
			"name":        h.getWorkoutName("cardio", lang),
			"description": h.getWorkoutDescription("cardio", lang),
			"focus":       "Cardiovascular Fitness",
			"exercises": []map[string]interface{}{
				{"name": "Running", "duration": "20 min", "intensity": "moderate"},
				{"name": "Jumping Jacks", "duration": "5 min", "intensity": "high"},
				{"name": "Burpees", "sets": 3, "reps": 10, "rest": "60s"},
			},
			"duration":  "30 minutes",
			"equipment": []string{"treadmill", "none"},
		},
	}

	// Adjust based on level
	if level == "beginner" {
		for _, template := range templates {
			if exercises, ok := template["exercises"].([]map[string]interface{}); ok {
				for _, exercise := range exercises {
					if reps, ok := exercise["reps"].(int); ok {
						exercise["reps"] = int(float64(reps) * 0.7)
					}
				}
			}
		}
	} else if level == "advanced" {
		for _, template := range templates {
			if exercises, ok := template["exercises"].([]map[string]interface{}); ok {
				for _, exercise := range exercises {
					if reps, ok := exercise["reps"].(int); ok {
						exercise["reps"] = int(float64(reps) * 1.3)
					}
				}
			}
		}
	}

	return templates
}

// Helper functions for workout names and descriptions
func (h *NutritionDataHandler) getWorkoutName(workoutType, lang string) string {
	names := map[string]map[string]string{
		"upper_body": {
			"en": "Upper Body Power Workout",
			"ar": "تمرين القوة الجزء العلوي من الجسم",
		},
		"lower_body": {
			"en": "Lower Body Strength Session",
			"ar": "جلسة قوة الجزء السفلي من الجسم",
		},
		"cardio": {
			"en": "Cardio Blast Session",
			"ar": "جلسة الكارديو المكثفة",
		},
	}

	if name, ok := names[workoutType][lang]; ok {
		return name
	}
	return names[workoutType]["en"]
}

func (h *NutritionDataHandler) getWorkoutDescription(workoutType, lang string) string {
	descriptions := map[string]map[string]string{
		"upper_body": {
			"en": "Focus on building upper body strength and muscle definition",
			"ar": "التركيز على بناء قوة الجزء العلوي من الجسم وتعريف العضلات",
		},
		"lower_body": {
			"en": "Strengthen and tone lower body muscles with compound movements",
			"ar": "تقوية وتنغيم عضلات الجزء السفلي من الجسم بالحركات المركبة",
		},
		"cardio": {
			"en": "Improve cardiovascular health and burn calories effectively",
			"ar": "تحسين صحة القلب والأوعية الدموية وحرق السعرات الحرارية بفعالية",
		},
	}

	if desc, ok := descriptions[workoutType][lang]; ok {
		return desc
	}
	return descriptions[workoutType]["en"]
}

// GetSupplementsRecommendations returns supplement recommendations
func (h *NutritionDataHandler) GetSupplementsRecommendations(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	filePath := filepath.Join(h.DataPath, "supplements.json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Supplements data not found",
			"code":  "SUPPLEMENTS_NOT_FOUND",
		})
	}

	var supplementsData map[string]interface{}
	if err := json.Unmarshal(data, &supplementsData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse supplements data",
			"code":  "PARSE_ERROR",
		})
	}

	var supplements []interface{}
	if supp, ok := supplementsData["supplements"].([]interface{}); ok {
		supplements = supp
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid supplements data format",
			"code":  "INVALID_FORMAT",
		})
	}

	response := map[string]interface{}{
		"supplements": supplements,
		"status":      "success",
		"message":     fmt.Sprintf("Supplement recommendations retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}

// GetDrugInteractions returns drug interaction information
func (h *NutritionDataHandler) GetDrugInteractions(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	// Mock drug interactions data
	interactions := []map[string]interface{}{
		{
			"drug_name": map[string]string{
				"en": "Warfarin",
				"ar": "الوارفارين",
			},
			"interactions": []string{
				"Vitamin K can reduce effectiveness",
				"Avoid large amounts of green leafy vegetables",
			},
			"severity": "high",
		},
		{
			"drug_name": map[string]string{
				"en": "Metformin",
				"ar": "الميتفورمين",
			},
			"interactions": []string{
				"Vitamin B12 deficiency possible with long-term use",
				"Monitor B12 levels regularly",
			},
			"severity": "medium",
		},
	}

	response := map[string]interface{}{
		"drug_interactions": interactions,
		"status":            "success",
		"message":           fmt.Sprintf("Drug interactions retrieved successfully in %s", lang),
	}

	return c.JSON(http.StatusOK, response)
}
