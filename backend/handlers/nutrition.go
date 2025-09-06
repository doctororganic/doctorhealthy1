package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// NutritionResponse represents a standardized API response for nutrition data
type NutritionResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// MetaInfo provides pagination and filtering information
type MetaInfo struct {
	Total       int    `json:"total"`
	Page        int    `json:"page"`
	PerPage     int    `json:"per_page"`
	TotalPages  int    `json:"total_pages"`
	Language    string `json:"language"`
	Country     string `json:"country"`
	APIVersion  string `json:"api_version"`
	CacheExpiry int    `json:"cache_expiry_seconds"`
}

// GetNutritionMeals returns meal and recipe data from the unified nutrition schema
func GetNutritionMeals(c echo.Context) error {
	// Extract query parameters
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en" // Default to English
	}
	
	country := c.QueryParam("country")
	if country == "" {
		country = "US" // Default country
	}
	
	category := c.QueryParam("category")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Mock data - In production, this would come from database
	mealsData := map[string]interface{}{
		"meals": []map[string]interface{}{
			{
				"id":          "meal_001",
				"name":        map[string]string{"en": "Mediterranean Bowl", "ar": "وعاء البحر الأبيض المتوسط"},
				"description": map[string]string{"en": "Healthy Mediterranean-style bowl", "ar": "وعاء صحي على الطراز المتوسطي"},
				"category":    "healthy",
				"cuisine":     "mediterranean",
				"prep_time":   15,
				"cook_time":   20,
				"servings":    2,
				"difficulty":  "easy",
				"nutrition": map[string]interface{}{
					"calories":      450,
					"protein":       25.5,
					"carbohydrates": 35.2,
					"fat":           18.7,
					"fiber":         8.3,
					"sodium":        680,
				},
				"ingredients": []map[string]interface{}{
					{
						"name":        map[string]string{"en": "Quinoa", "ar": "الكينوا"},
						"amount":      "1 cup",
						"alternatives": []string{"brown rice", "bulgur"},
					},
					{
						"name":        map[string]string{"en": "Chicken Breast", "ar": "صدر الدجاج"},
						"amount":      "200g",
						"alternatives": []string{"tofu", "chickpeas"},
					},
				},
				"instructions": []map[string]string{
					{"en": "Cook quinoa according to package instructions", "ar": "اطبخ الكينوا حسب تعليمات العبوة"},
					{"en": "Grill chicken breast until cooked through", "ar": "اشوي صدر الدجاج حتى ينضج تماماً"},
					{"en": "Combine all ingredients in a bowl", "ar": "اخلط جميع المكونات في وعاء"},
				},
				"tags":        []string{"healthy", "protein-rich", "gluten-free"},
				"allergens":   []string{"none"},
				"dietary":     []string{"halal", "gluten-free"},
			},
		},
		"recipes": []map[string]interface{}{
			{
				"id":          "recipe_001",
				"name":        map[string]string{"en": "Koshari", "ar": "كشري"},
				"description": map[string]string{"en": "Traditional Egyptian comfort food", "ar": "طعام مصري تقليدي مريح"},
				"category":    "traditional",
				"cuisine":     "egyptian",
				"region":      "middle_east",
				"prep_time":   30,
				"cook_time":   45,
				"servings":    4,
				"difficulty":  "medium",
				"nutrition": map[string]interface{}{
					"calories":      520,
					"protein":       18.2,
					"carbohydrates": 78.5,
					"fat":           12.3,
					"fiber":         12.1,
					"sodium":        890,
				},
				"cultural_significance": map[string]string{
					"en": "National dish of Egypt, popular street food",
					"ar": "الطبق الوطني لمصر، طعام شارع شعبي",
				},
			},
		},
	}

	// Filter by category if specified
	if category != "" {
		// Filter logic would go here
	}

	// Calculate pagination
	totalMeals := len(mealsData["meals"].([]map[string]interface{})) + len(mealsData["recipes"].([]map[string]interface{}))
	totalPages := (totalMeals + perPage - 1) / perPage

	meta := &MetaInfo{
		Total:       totalMeals,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		Language:    lang,
		Country:     country,
		APIVersion:  "1.0",
		CacheExpiry: 3600, // 1 hour
	}

	response := NutritionResponse{
		Success: true,
		Message: "Meals and recipes retrieved successfully",
		Data:    mealsData,
		Meta:    meta,
	}

	// Set cache headers
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	c.Response().Header().Set("ETag", `"meals-v1"`)

	return c.JSON(http.StatusOK, response)
}

// GetNutritionRecipes returns recipe data specifically
func GetNutritionRecipes(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	cuisine := c.QueryParam("cuisine")
	region := c.QueryParam("region")
	dietary := c.QueryParam("dietary")

	// Mock recipe data
	recipesData := map[string]interface{}{
		"recipes": []map[string]interface{}{
			{
				"id":          "recipe_koshari",
				"name":        map[string]string{"en": "Egyptian Koshari", "ar": "كشري مصري"},
				"description": map[string]string{"en": "Traditional Egyptian mixed rice dish", "ar": "طبق الأرز المصري التقليدي المختلط"},
				"cuisine":     "egyptian",
				"region":      "middle_east",
				"dietary":     []string{"vegetarian", "halal"},
				"prep_time":   30,
				"cook_time":   45,
				"total_time":  75,
				"servings":    6,
				"difficulty":  "medium",
				"cost_level":  "budget",
				"nutrition": map[string]interface{}{
					"per_serving": map[string]interface{}{
						"calories":      485,
						"protein":       16.8,
						"carbohydrates": 82.3,
						"fat":           8.9,
						"fiber":         11.2,
						"sugar":         12.5,
						"sodium":        920,
						"potassium":     680,
						"iron":          4.2,
						"calcium":       85,
					},
				},
				"ingredients": []map[string]interface{}{
					{
						"name":        map[string]string{"en": "Basmati Rice", "ar": "أرز بسمتي"},
						"amount":      "2 cups",
						"unit":        "cups",
						"weight":      "400g",
						"alternatives": []string{"jasmine rice", "long grain rice"},
						"optional":     false,
					},
					{
						"name":        map[string]string{"en": "Brown Lentils", "ar": "عدس بني"},
						"amount":      "1 cup",
						"unit":        "cups",
						"weight":      "200g",
						"alternatives": []string{"green lentils"},
						"optional":     false,
					},
				},
				"instructions": []map[string]interface{}{
					{
						"step":        1,
						"description": map[string]string{"en": "Rinse and cook lentils until tender", "ar": "اغسل واطبخ العدس حتى ينضج"},
						"time":        15,
						"temperature": "medium heat",
						"tips":        map[string]string{"en": "Don't overcook to avoid mushy texture", "ar": "لا تفرط في الطبخ لتجنب القوام الطري"},
					},
					{
						"step":        2,
						"description": map[string]string{"en": "Cook rice separately until fluffy", "ar": "اطبخ الأرز منفصلاً حتى يصبح رقيقاً"},
						"time":        18,
						"temperature": "medium-low heat",
					},
				},
				"equipment": []string{"large pot", "strainer", "wooden spoon"},
				"storage": map[string]interface{}{
					"refrigerator": "3-4 days",
					"freezer":      "2-3 months",
					"instructions": map[string]string{"en": "Store in airtight container", "ar": "احفظ في وعاء محكم الإغلاق"},
				},
				"cultural_info": map[string]interface{}{
					"origin":      "Egypt",
					"significance": map[string]string{"en": "National dish and popular street food", "ar": "الطبق الوطني وطعام الشارع الشعبي"},
					"occasions":   []string{"daily meals", "family gatherings", "street food"},
				},
			},
		},
		"filters": map[string]interface{}{
			"cuisines":     []string{"egyptian", "lebanese", "moroccan", "turkish", "persian"},
			"regions":      []string{"middle_east", "north_africa", "mediterranean"},
			"dietary":      []string{"vegetarian", "vegan", "halal", "gluten-free", "dairy-free"},
			"difficulties": []string{"easy", "medium", "hard"},
			"cost_levels":  []string{"budget", "moderate", "premium"},
		},
	}

	// Apply filters
	if cuisine != "" || region != "" || dietary != "" {
		// Filter logic would be implemented here
	}

	response := NutritionResponse{
		Success: true,
		Message: "Recipes retrieved successfully",
		Data:    recipesData,
		Meta: &MetaInfo{
			Language:    lang,
			APIVersion:  "1.0",
			CacheExpiry: 3600,
		},
	}

	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return c.JSON(http.StatusOK, response)
}

// GetNutritionWorkouts returns workout and exercise data
func GetNutritionWorkouts(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	workoutType := c.QueryParam("type")
	level := c.QueryParam("level")
	duration := c.QueryParam("duration")

	// Mock workout data
	workoutsData := map[string]interface{}{
		"workouts": []map[string]interface{}{
			{
				"id":          "workout_shred",
				"name":        map[string]string{"en": "Shred Workout Plan", "ar": "خطة تمرين التقطيع"},
				"description": map[string]string{"en": "High-intensity fat burning workout", "ar": "تمرين عالي الكثافة لحرق الدهون"},
				"type":        "strength",
				"level":       "intermediate",
				"duration":    45,
				"equipment":   []string{"dumbbells", "resistance bands", "yoga mat"},
				"target_areas": []string{"full body", "core", "cardio"},
				"exercises": []map[string]interface{}{
					{
						"name":        map[string]string{"en": "Burpees", "ar": "بيربيز"},
						"description": map[string]string{"en": "Full body explosive movement", "ar": "حركة انفجارية للجسم كاملاً"},
						"sets":        3,
						"reps":        "12-15",
						"rest":        "60 seconds",
						"calories_burned": 12,
						"instructions": []map[string]string{
							{"en": "Start in standing position", "ar": "ابدأ في وضعية الوقوف"},
							{"en": "Drop into squat position", "ar": "انزل إلى وضعية القرفصاء"},
							{"en": "Jump back to plank", "ar": "اقفز للخلف إلى وضعية البلانك"},
							{"en": "Do a push-up", "ar": "قم بضغطة واحدة"},
							{"en": "Jump feet back to squat", "ar": "اقفز بالقدمين للخلف للقرفصاء"},
							{"en": "Explode up with arms overhead", "ar": "انفجر للأعلى مع رفع الذراعين"},
						},
						"common_mistakes": []map[string]string{
							{"en": "Landing hard on feet", "ar": "الهبوط بقوة على القدمين"},
							{"en": "Not maintaining plank position", "ar": "عدم الحفاظ على وضعية البلانك"},
							{"en": "Rushing through movements", "ar": "التسرع في الحركات"},
						},
						"injury_risks": []map[string]string{
							{"en": "Wrist strain from improper landing", "ar": "إجهاد المعصم من الهبوط غير الصحيح"},
							{"en": "Lower back stress", "ar": "ضغط أسفل الظهر"},
						},
						"alternatives": []map[string]interface{}{
							{
								"name": map[string]string{"en": "Modified Burpee", "ar": "بيربي معدل"},
								"description": map[string]string{"en": "Step back instead of jumping", "ar": "خطوة للخلف بدلاً من القفز"},
							},
							{
								"name": map[string]string{"en": "Half Burpee", "ar": "نصف بيربي"},
								"description": map[string]string{"en": "Without push-up component", "ar": "بدون مكون الضغط"},
							},
						},
					},
				},
				"warm_up": map[string]interface{}{
					"duration": 10,
					"exercises": []string{"arm circles", "leg swings", "light cardio"},
					"tips": map[string]string{"en": "Gradually increase intensity", "ar": "زد الكثافة تدريجياً"},
				},
				"cool_down": map[string]interface{}{
					"duration": 10,
					"exercises": []string{"static stretching", "deep breathing"},
					"tips": map[string]string{"en": "Hold stretches for 30 seconds", "ar": "احتفظ بالتمدد لمدة 30 ثانية"},
				},
			},
		},
		"exercise_library": []map[string]interface{}{
			{
				"id":          "exercise_pushup",
				"name":        map[string]string{"en": "Push-up", "ar": "ضغط"},
				"category":    "strength",
				"muscle_groups": []string{"chest", "shoulders", "triceps", "core"},
				"equipment":   "bodyweight",
				"difficulty":  "beginner",
				"calories_per_rep": 0.5,
			},
		},
	}

	// Apply filters
	if workoutType != "" || level != "" || duration != "" {
		// Filter logic would be implemented here
	}

	response := NutritionResponse{
		Success: true,
		Message: "Workouts retrieved successfully",
		Data:    workoutsData,
		Meta: &MetaInfo{
			Language:    lang,
			APIVersion:  "1.0",
			CacheExpiry: 3600,
		},
	}

	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return c.JSON(http.StatusOK, response)
}

// GetPregnancyNutrition returns pregnancy-specific nutrition data
func GetPregnancyNutrition(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	trimester := c.QueryParam("trimester")
	condition := c.QueryParam("condition")

	// Mock pregnancy nutrition data
	pregnancyData := map[string]interface{}{
		"nutrition_guidelines": map[string]interface{}{
			"normal_pregnancy": map[string]interface{}{
				"calories": map[string]interface{}{
					"first_trimester":  "No additional calories needed",
					"second_trimester": "+340 calories/day",
					"third_trimester":  "+450 calories/day",
					"description": map[string]string{
						"en": "Caloric needs increase gradually during pregnancy",
						"ar": "تزداد الاحتياجات من السعرات الحرارية تدريجياً أثناء الحمل",
					},
				},
				"protein": map[string]interface{}{
					"daily_requirement": "71g/day",
					"sources": []string{"lean meats", "fish", "eggs", "dairy", "legumes", "nuts"},
					"description": map[string]string{
						"en": "Essential for fetal growth and maternal tissue development",
						"ar": "ضروري لنمو الجنين وتطوير أنسجة الأم",
					},
				},
				"folate": map[string]interface{}{
					"daily_requirement": "600mcg/day",
					"sources": []string{"leafy greens", "citrus fruits", "fortified grains", "legumes"},
					"importance": map[string]string{
						"en": "Prevents neural tube defects",
						"ar": "يمنع عيوب الأنبوب العصبي",
					},
				},
				"iron": map[string]interface{}{
					"daily_requirement": "27mg/day",
					"sources": []string{"red meat", "poultry", "fish", "dried beans", "fortified cereals"},
					"importance": map[string]string{
						"en": "Prevents anemia and supports increased blood volume",
						"ar": "يمنع فقر الدم ويدعم زيادة حجم الدم",
					},
				},
				"calcium": map[string]interface{}{
					"daily_requirement": "1000mg/day",
					"sources": []string{"dairy products", "fortified plant milks", "leafy greens", "sardines"},
					"importance": map[string]string{
						"en": "Essential for fetal bone development",
						"ar": "ضروري لتطوير عظام الجنين",
					},
				},
			},
		},
		"foods_to_avoid": []map[string]interface{}{
			{
				"category": "high_mercury_fish",
				"items": []string{"shark", "swordfish", "king mackerel", "tilefish"},
				"reason": map[string]string{
					"en": "High mercury levels can harm fetal nervous system",
					"ar": "مستويات الزئبق العالية يمكن أن تضر بالجهاز العصبي للجنين",
				},
			},
			{
				"category": "raw_undercooked",
				"items": []string{"raw fish", "raw eggs", "undercooked meat", "unpasteurized dairy"},
				"reason": map[string]string{
					"en": "Risk of foodborne illness",
					"ar": "خطر الأمراض المنقولة بالغذاء",
				},
			},
		},
		"supplements": []map[string]interface{}{
			{
				"name": "Prenatal Vitamin",
				"description": map[string]string{
					"en": "Comprehensive vitamin and mineral supplement",
					"ar": "مكمل شامل للفيتامينات والمعادن",
				},
				"key_nutrients": []string{"folic acid", "iron", "calcium", "vitamin D", "DHA"},
				"timing": map[string]string{
					"en": "Take with food to reduce nausea",
					"ar": "تناول مع الطعام لتقليل الغثيان",
				},
			},
		},
		"meal_planning": map[string]interface{}{
			"sample_day": []map[string]interface{}{
				{
					"meal": "breakfast",
					"foods": []string{"fortified cereal", "milk", "banana", "orange juice"},
					"nutrients_focus": []string{"folate", "calcium", "vitamin C"},
				},
				{
					"meal": "lunch",
					"foods": []string{"spinach salad", "grilled chicken", "quinoa", "avocado"},
					"nutrients_focus": []string{"iron", "protein", "folate", "healthy fats"},
				},
				{
					"meal": "dinner",
					"foods": []string{"salmon", "sweet potato", "broccoli", "brown rice"},
					"nutrients_focus": []string{"DHA", "vitamin A", "fiber", "complex carbs"},
				},
			},
		},
	}

	// Filter by trimester or condition if specified
	if trimester != "" || condition != "" {
		// Filter logic would be implemented here
	}

	response := NutritionResponse{
		Success: true,
		Message: "Pregnancy nutrition information retrieved successfully",
		Data:    pregnancyData,
		Meta: &MetaInfo{
			Language:    lang,
			APIVersion:  "1.0",
			CacheExpiry: 7200, // 2 hours
		},
	}

	c.Response().Header().Set("Cache-Control", "public, max-age=7200")
	return c.JSON(http.StatusOK, response)
}

// GetPediatricNutrition returns pediatric nutrition data
func GetPediatricNutrition(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	ageGroup := c.QueryParam("age_group")
	condition := c.QueryParam("condition")

	// Mock pediatric nutrition data
	pediatricData := map[string]interface{}{
		"age_groups": map[string]interface{}{
			"infants_0_6_months": map[string]interface{}{
				"feeding": "exclusive_breastfeeding",
				"description": map[string]string{
					"en": "Breast milk provides all necessary nutrients",
					"ar": "حليب الأم يوفر جميع العناصر الغذائية الضرورية",
				},
				"supplements": []string{"vitamin D"},
			},
			"infants_6_12_months": map[string]interface{}{
				"feeding": "breastfeeding_plus_solids",
				"first_foods": []string{"iron-fortified cereal", "pureed vegetables", "pureed fruits"},
				"foods_to_avoid": []string{"honey", "cow's milk", "choking hazards"},
				"description": map[string]string{
					"en": "Gradual introduction of solid foods",
					"ar": "إدخال تدريجي للأطعمة الصلبة",
				},
			},
			"toddlers_1_3_years": map[string]interface{}{
				"calories_per_day": "1000-1400",
				"key_nutrients": []string{"iron", "calcium", "vitamin D", "healthy fats"},
				"portion_sizes": map[string]string{
					"grains": "1/4 to 1/2 cup",
					"vegetables": "1/4 to 1/2 cup",
					"fruits": "1/4 to 1/2 cup",
					"protein": "1-2 tablespoons",
				},
				"feeding_tips": []map[string]string{
					{"en": "Offer variety and let child explore", "ar": "قدم التنوع ودع الطفل يستكشف"},
					{"en": "Be patient with picky eating", "ar": "كن صبوراً مع الأكل الانتقائي"},
					{"en": "Establish regular meal times", "ar": "حدد أوقات وجبات منتظمة"},
				},
			},
			"preschoolers_3_5_years": map[string]interface{}{
				"calories_per_day": "1200-1600",
				"focus_areas": []string{"balanced meals", "healthy snacks", "hydration"},
				"sample_meals": []map[string]interface{}{
					{
						"meal": "breakfast",
						"example": "whole grain cereal with milk and berries",
						"nutrients": []string{"fiber", "calcium", "antioxidants"},
					},
					{
						"meal": "lunch",
						"example": "turkey and cheese sandwich with vegetables",
						"nutrients": []string{"protein", "calcium", "vitamins"},
					},
				},
			},
		},
		"common_concerns": []map[string]interface{}{
			{
				"concern": "picky_eating",
				"description": map[string]string{
					"en": "Child refuses many foods or has limited food preferences",
					"ar": "الطفل يرفض العديد من الأطعمة أو لديه تفضيلات غذائية محدودة",
				},
				"strategies": []map[string]string{
					{"en": "Offer new foods multiple times", "ar": "قدم الأطعمة الجديدة عدة مرات"},
					{"en": "Make mealtimes pleasant", "ar": "اجعل أوقات الوجبات ممتعة"},
					{"en": "Be a good role model", "ar": "كن قدوة جيدة"},
				},
			},
			{
				"concern": "food_allergies",
				"description": map[string]string{
					"en": "Adverse reactions to specific foods",
					"ar": "ردود فعل سلبية لأطعمة معينة",
				},
				"common_allergens": []string{"milk", "eggs", "peanuts", "tree nuts", "soy", "wheat", "fish", "shellfish"},
				"management": []map[string]string{
					{"en": "Read food labels carefully", "ar": "اقرأ ملصقات الطعام بعناية"},
					{"en": "Have emergency action plan", "ar": "احتفظ بخطة عمل طوارئ"},
					{"en": "Work with healthcare provider", "ar": "اعمل مع مقدم الرعاية الصحية"},
				},
			},
		},
		"nutritional_supplements": map[string]interface{}{
			"vitamin_d": map[string]interface{}{
				"importance": map[string]string{
					"en": "Essential for bone development and immune function",
					"ar": "ضروري لتطوير العظام ووظيفة المناعة",
				},
				"recommended_dose": "400 IU/day for infants, 600 IU/day for children",
			},
			"iron": map[string]interface{}{
				"importance": map[string]string{
					"en": "Prevents iron deficiency anemia",
					"ar": "يمنع فقر الدم الناتج عن نقص الحديد",
				},
				"food_sources": []string{"fortified cereals", "lean meats", "beans", "spinach"},
			},
		},
	}

	// Filter by age group or condition if specified
	if ageGroup != "" || condition != "" {
		// Filter logic would be implemented here
	}

	response := NutritionResponse{
		Success: true,
		Message: "Pediatric nutrition information retrieved successfully",
		Data:    pediatricData,
		Meta: &MetaInfo{
			Language:    lang,
			APIVersion:  "1.0",
			CacheExpiry: 7200, // 2 hours
		},
	}

	c.Response().Header().Set("Cache-Control", "public, max-age=7200")
	return c.JSON(http.StatusOK, response)
}

// GetSportsNutrition returns sports and performance nutrition data
func GetSportsNutrition(c echo.Context) error {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = "en"
	}

	sportType := c.QueryParam("sport_type")
	goal := c.QueryParam("goal")
	intensity := c.QueryParam("intensity")

	// Mock sports nutrition data
	sportsData := map[string]interface{}{
		"energy_systems": map[string]interface{}{
			"phosphagen_system": map[string]interface{}{
				"duration": "0-10 seconds",
				"primary_fuel": "creatine phosphate",
				"sports_examples": []string{"weightlifting", "sprinting", "jumping"},
				"nutrition_focus": []string{"creatine", "adequate protein"},
				"description": map[string]string{
					"en": "Immediate energy for explosive movements",
					"ar": "طاقة فورية للحركات الانفجارية",
				},
			},
			"glycolytic_system": map[string]interface{}{
				"duration": "10 seconds - 2 minutes",
				"primary_fuel": "glucose/glycogen",
				"sports_examples": []string{"400m sprint", "swimming laps", "basketball"},
				"nutrition_focus": []string{"carbohydrates", "glycogen loading"},
				"description": map[string]string{
					"en": "High-intensity energy from carbohydrates",
					"ar": "طاقة عالية الكثافة من الكربوهيدرات",
				},
			},
			"oxidative_system": map[string]interface{}{
				"duration": "2+ minutes",
				"primary_fuel": "fats and carbohydrates",
				"sports_examples": []string{"marathon", "cycling", "triathlon"},
				"nutrition_focus": []string{"endurance nutrition", "fat adaptation", "hydration"},
				"description": map[string]string{
					"en": "Sustained energy for endurance activities",
					"ar": "طاقة مستدامة لأنشطة التحمل",
				},
			},
		},
		"macronutrient_needs": map[string]interface{}{
			"endurance_sports": map[string]interface{}{
				"carbohydrates": "6-10g/kg body weight",
				"protein": "1.2-1.4g/kg body weight",
				"fat": "20-35% of total calories",
				"hydration": "150-250ml every 15-20 minutes",
				"timing": map[string]interface{}{
					"pre_exercise": "3-4 hours before: high carb, low fat meal",
					"during_exercise": "30-60g carbs per hour for >1 hour exercise",
					"post_exercise": "1.5g carbs/kg + 0.25g protein/kg within 30 minutes",
				},
			},
			"strength_sports": map[string]interface{}{
				"carbohydrates": "3-5g/kg body weight",
				"protein": "1.6-2.2g/kg body weight",
				"fat": "20-35% of total calories",
				"timing": map[string]interface{}{
					"pre_workout": "Balanced meal 2-3 hours before",
					"post_workout": "20-40g protein within 2 hours",
					"daily": "Distribute protein throughout the day",
				},
			},
			"team_sports": map[string]interface{}{
				"carbohydrates": "5-7g/kg body weight",
				"protein": "1.2-1.7g/kg body weight",
				"fat": "20-35% of total calories",
				"considerations": []map[string]string{
					{"en": "Variable intensity requires flexible fueling", "ar": "الكثافة المتغيرة تتطلب تزويد مرن بالوقود"},
					{"en": "Focus on recovery between games", "ar": "ركز على التعافي بين المباريات"},
				},
			},
		},
		"supplements": []map[string]interface{}{
			{
				"name": "Creatine Monohydrate",
				"benefits": map[string]string{
					"en": "Improves power output and muscle mass",
					"ar": "يحسن القوة الناتجة وكتلة العضلات",
				},
				"dosage": "3-5g daily",
				"best_for": []string{"strength training", "high-intensity intervals"},
				"evidence_level": "strong",
			},
			{
				"name": "Caffeine",
				"benefits": map[string]string{
					"en": "Enhances endurance and reduces perceived exertion",
					"ar": "يعزز التحمل ويقلل الجهد المدرك",
				},
				"dosage": "3-6mg/kg body weight",
				"timing": "30-60 minutes before exercise",
				"best_for": []string{"endurance sports", "team sports"},
				"evidence_level": "strong",
			},
			{
				"name": "Beta-Alanine",
				"benefits": map[string]string{
					"en": "Reduces muscle fatigue in high-intensity exercise",
					"ar": "يقلل إجهاد العضلات في التمارين عالية الكثافة",
				},
				"dosage": "3-5g daily (divided doses)",
				"best_for": []string{"1-4 minute high-intensity efforts"},
				"evidence_level": "moderate",
			},
		},
		"hydration_guidelines": map[string]interface{}{
			"daily_needs": "35-40ml/kg body weight",
			"exercise_needs": map[string]interface{}{
				"before": "400-600ml 2-3 hours before",
				"during": "150-250ml every 15-20 minutes",
				"after": "150% of fluid lost through sweat",
			},
			"electrolyte_replacement": map[string]interface{}{
				"when_needed": "Exercise >1 hour or heavy sweating",
				"sodium": "200-700mg per hour",
				"potassium": "150-300mg per hour",
			},
			"signs_of_dehydration": []map[string]string{
				{"en": "Dark yellow urine", "ar": "بول أصفر داكن"},
				{"en": "Decreased performance", "ar": "انخفاض الأداء"},
				{"en": "Increased heart rate", "ar": "زيادة معدل ضربات القلب"},
				{"en": "Fatigue and dizziness", "ar": "التعب والدوخة"},
			},
		},
	}

	// Apply filters
	if sportType != "" || goal != "" || intensity != "" {
		// Filter logic would be implemented here
	}

	response := NutritionResponse{
		Success: true,
		Message: "Sports nutrition information retrieved successfully",
		Data:    sportsData,
		Meta: &MetaInfo{
			Language:    lang,
			APIVersion:  "1.0",
			CacheExpiry: 3600,
		},
	}

	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return c.JSON(http.StatusOK, response)
}

// Helper function to validate and sanitize query parameters
func validateQueryParams(c echo.Context) map[string]string {
	params := make(map[string]string)
	
	// Language validation
	lang := c.QueryParam("lang")
	validLangs := []string{"en", "ar", "fr", "es"}
	if contains(validLangs, lang) {
		params["lang"] = lang
	} else {
		params["lang"] = "en"
	}
	
	// Country validation
	country := c.QueryParam("country")
	validCountries := []string{"US", "SA", "AE", "EG", "MA", "TN", "JO", "LB"}
	if contains(validCountries, country) {
		params["country"] = country
	} else {
		params["country"] = "US"
	}
	
	return params
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to paginate results
func paginateResults(data interface{}, page, perPage int) (interface{}, *MetaInfo) {
	// This would implement actual pagination logic
	// For now, return mock pagination info
	meta := &MetaInfo{
		Page:       page,
		PerPage:    perPage,
		Total:      100, // Mock total
		TotalPages: 5,   // Mock total pages
		APIVersion: "1.0",
	}
	return data, meta
}

// Helper function to filter data based on criteria
func filterNutritionData(data map[string]interface{}, filters map[string]string) map[string]interface{} {
	// This would implement actual filtering logic
	// For now, return data as-is
	return data
}

// Helper function to get localized content
func getLocalizedContent(content map[string]string, lang string) string {
	if val, exists := content[lang]; exists {
		return val
	}
	// Fallback to English
	if val, exists := content["en"]; exists {
		return val
	}
	// Return first available language
	for _, val := range content {
		return val
	}
	return ""
}

// GetCalories retrieves calories data with pagination and filtering
func GetCalories(c echo.Context) error {
	if err := validateNutritionAccess(c, "nutrition_read"); err != nil {
		return err
	}

	// Get query parameters
	category := c.QueryParam("category")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Mock calories data structure
	caloriesData := map[string]interface{}{
		"starches": []map[string]interface{}{
			{
				"name_en": "White Rice",
				"name_ar": "الأرز الأبيض",
				"serving_size": "1/2 cup cooked",
				"calories": 85,
				"carbs_g": 15,
				"protein_g": 3,
				"fat_g": 0,
			},
			{
				"name_en": "Brown Rice",
				"name_ar": "الأرز البني",
				"serving_size": "1/2 cup cooked",
				"calories": 110,
				"carbs_g": 23,
				"protein_g": 3,
				"fat_g": 1,
			},
		},
		"proteins": []map[string]interface{}{
			{
				"name_en": "Chicken Breast",
				"name_ar": "صدر الدجاج",
				"serving_size": "3 oz cooked",
				"calories": 140,
				"carbs_g": 0,
				"protein_g": 26,
				"fat_g": 3,
			},
		},
		"vegetables": []map[string]interface{}{
			{
				"name_en": "Broccoli",
				"name_ar": "البروكلي",
				"serving_size": "1 cup cooked",
				"calories": 25,
				"carbs_g": 5,
				"protein_g": 3,
				"fat_g": 0,
			},
		},
		"fruits": []map[string]interface{}{
			{
				"name_en": "Apple",
				"name_ar": "التفاح",
				"serving_size": "1 medium",
				"calories": 80,
				"carbs_g": 22,
				"protein_g": 0,
				"fat_g": 0,
			},
		},
		"fats": []map[string]interface{}{
			{
				"name_en": "Olive Oil",
				"name_ar": "زيت الزيتون",
				"serving_size": "1 tbsp",
				"calories": 120,
				"carbs_g": 0,
				"protein_g": 0,
				"fat_g": 14,
			},
		},
	}

	// Filter by category if specified
	var responseData interface{}
	if category != "" {
		if categoryData, exists := caloriesData[category]; exists {
			responseData = categoryData
		} else {
			return c.JSON(http.StatusNotFound, NutritionResponse{
				Success: false,
				Message: "Category not found",
			})
		}
	} else {
		responseData = caloriesData
	}

	return c.JSON(http.StatusOK, NutritionResponse{
		Success: true,
		Message: "Calories data retrieved successfully",
		Data:    responseData,
		Meta: &MetaInfo{
			Total:      len(caloriesData),
			Page:       page,
			PerPage:    perPage,
			TotalPages: 1,
			Language:   "en",
			Country:    "US",
			APIVersion: "1.0",
			CacheExpiry: 3600,
		},
	})
}

// GetCaloriesByCategory retrieves calories data for a specific category
func GetCaloriesByCategory(c echo.Context) error {
	if err := validateNutritionAccess(c, "nutrition_read"); err != nil {
		return err
	}

	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, NutritionResponse{
			Success: false,
			Message: "Category parameter is required",
		})
	}

	// Set category as query parameter and call GetCalories
	c.QueryParams().Set("category", category)
	return GetCalories(c)
}

// GetSkills retrieves cooking skills data with pagination and filtering
func GetSkills(c echo.Context) error {
	if err := validateNutritionAccess(c, "nutrition_read"); err != nil {
		return err
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Mock skills data structure
	skillsData := map[string]interface{}{
		"cooking_skills": []map[string]interface{}{
			{
				"name": map[string]string{
					"en": "Knife Skills",
					"ar": "مهارات السكين",
				},
				"description": map[string]string{
					"en": "Master the art of chopping, dicing, and slicing with proper knife techniques for efficient and safe cooking.",
					"ar": "إتقان فن التقطيع والتقطيع إلى مكعبات والتقطيع بتقنيات السكين المناسبة للطبخ الفعال والآمن.",
				},
				"difficulty": "beginner",
				"category": "basic_techniques",
			},
			{
				"name": map[string]string{
					"en": "Sautéing",
					"ar": "القلي السريع",
				},
				"description": map[string]string{
					"en": "Learn to cook food quickly in a small amount of oil or fat over high heat while stirring frequently.",
					"ar": "تعلم طبخ الطعام بسرعة في كمية قليلة من الزيت أو الدهن على نار عالية مع التحريك المتكرر.",
				},
				"difficulty": "intermediate",
				"category": "cooking_methods",
			},
			{
				"name": map[string]string{
					"en": "Braising",
					"ar": "الطبخ البطيء",
				},
				"description": map[string]string{
					"en": "A combination cooking method that uses both wet and dry heat to break down tough cuts of meat.",
					"ar": "طريقة طبخ مختلطة تستخدم الحرارة الرطبة والجافة لتفكيك قطع اللحم القاسية.",
				},
				"difficulty": "advanced",
				"category": "cooking_methods",
			},
			{
				"name": map[string]string{
					"en": "Seasoning and Flavoring",
					"ar": "التتبيل والنكهة",
				},
				"description": map[string]string{
					"en": "Understanding how to balance flavors using herbs, spices, acids, and aromatics to enhance dishes.",
					"ar": "فهم كيفية توازن النكهات باستخدام الأعشاب والتوابل والأحماض والعطريات لتحسين الأطباق.",
				},
				"difficulty": "intermediate",
				"category": "flavor_development",
			},
			{
				"name": map[string]string{
					"en": "Food Safety and Hygiene",
					"ar": "سلامة الغذاء والنظافة",
				},
				"description": map[string]string{
					"en": "Essential knowledge of proper food handling, storage, and preparation to prevent foodborne illnesses.",
					"ar": "المعرفة الأساسية للتعامل السليم مع الطعام وتخزينه وتحضيره لمنع الأمراض المنقولة بالغذاء.",
				},
				"difficulty": "beginner",
				"category": "safety",
			},
		},
	}

	return c.JSON(http.StatusOK, NutritionResponse{
		Success: true,
		Message: "Skills data retrieved successfully",
		Data:    skillsData,
		Meta: &MetaInfo{
			Total:      len(skillsData["cooking_skills"].([]map[string]interface{})),
			Page:       page,
			PerPage:    perPage,
			TotalPages: 1,
			Language:   "en",
			Country:    "US",
			APIVersion: "1.0",
			CacheExpiry: 3600,
		},
	})
}

// GetSkillsByDifficulty retrieves skills data filtered by difficulty level
func GetSkillsByDifficulty(c echo.Context) error {
	if err := validateNutritionAccess(c, "nutrition_read"); err != nil {
		return err
	}

	difficulty := c.Param("difficulty")
	if difficulty == "" {
		return c.JSON(http.StatusBadRequest, NutritionResponse{
			Success: false,
			Message: "Difficulty parameter is required",
		})
	}

	// Mock filtered skills data
	var filteredSkills []map[string]interface{}
	allSkills := []map[string]interface{}{
		{
			"name": map[string]string{
				"en": "Knife Skills",
				"ar": "مهارات السكين",
			},
			"difficulty": "beginner",
		},
		{
			"name": map[string]string{
				"en": "Sautéing",
				"ar": "القلي السريع",
			},
			"difficulty": "intermediate",
		},
		{
			"name": map[string]string{
				"en": "Braising",
				"ar": "الطبخ البطيء",
			},
			"difficulty": "advanced",
		},
	}

	for _, skill := range allSkills {
		if skill["difficulty"] == difficulty {
			filteredSkills = append(filteredSkills, skill)
		}
	}

	if len(filteredSkills) == 0 {
		return c.JSON(http.StatusNotFound, NutritionResponse{
			Success: false,
			Message: "No skills found for the specified difficulty level",
		})
	}

	return c.JSON(http.StatusOK, NutritionResponse{
		Success: true,
		Message: "Skills data retrieved successfully",
		Data:    map[string]interface{}{"cooking_skills": filteredSkills},
		Meta: &MetaInfo{
			Total:      len(filteredSkills),
			Page:       1,
			PerPage:    20,
			TotalPages: 1,
			Language:   "en",
			Country:    "US",
			APIVersion: "1.0",
			CacheExpiry: 3600,
		},
	})
}

// GetTypePlans retrieves diet type plans data with pagination and filtering
func GetTypePlans(c echo.Context) error {
	if err := validateNutritionAccess(c, "nutrition_read"); err != nil {
		return err
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Mock type plans data structure
	typePlansData := map[string]interface{}{
		"diets": []map[string]interface{}{
			{
				"name": map[string]string{
					"en": "Mediterranean Diet",
					"ar": "النظام الغذائي المتوسطي",
				},
				"allowed_foods": map[string]string{
					"en": "Fruits, vegetables, whole grains, legumes, nuts, olive oil, fish, and moderate amounts of dairy and poultry.",
					"ar": "الفواكه والخضروات والحبوب الكاملة والبقوليات والمكسرات وزيت الزيتون والأسماك وكميات معتدلة من منتجات الألبان والدواجن.",
				},
				"minimal_foods": map[string]string{
					"en": "Red meat, processed foods, refined sugars, and highly processed oils.",
					"ar": "اللحوم الحمراء والأطعمة المصنعة والسكريات المكررة والزيوت المعالجة بشدة.",
				},
				"type": "balanced",
				"difficulty": "moderate",
			},
			{
				"name": map[string]string{
					"en": "Ketogenic Diet",
					"ar": "النظام الغذائي الكيتوني",
				},
				"allowed_foods": map[string]string{
					"en": "High-fat foods like avocados, nuts, seeds, oils, fatty fish, meat, and low-carb vegetables.",
					"ar": "الأطعمة عالية الدهون مثل الأفوكادو والمكسرات والبذور والزيوت والأسماك الدهنية واللحوم والخضروات منخفضة الكربوهيدرات.",
				},
				"minimal_foods": map[string]string{
					"en": "Grains, sugar, fruits (except berries), starchy vegetables, and most dairy products.",
					"ar": "الحبوب والسكر والفواكه (باستثناء التوت) والخضروات النشوية ومعظم منتجات الألبان.",
				},
				"type": "low_carb",
				"difficulty": "challenging",
			},
			{
				"name": map[string]string{
					"en": "Vegetarian Diet",
					"ar": "النظام الغذائي النباتي",
				},
				"allowed_foods": map[string]string{
					"en": "Fruits, vegetables, grains, legumes, nuts, seeds, dairy products, and eggs.",
					"ar": "الفواكه والخضروات والحبوب والبقوليات والمكسرات والبذور ومنتجات الألبان والبيض.",
				},
				"minimal_foods": map[string]string{
					"en": "All meat, poultry, fish, and seafood.",
					"ar": "جميع اللحوم والدواجن والأسماك والمأكولات البحرية.",
				},
				"type": "plant_based",
				"difficulty": "easy",
			},
			{
				"name": map[string]string{
					"en": "Paleo Diet",
					"ar": "النظام الغذائي الباليو",
				},
				"allowed_foods": map[string]string{
					"en": "Meat, fish, eggs, vegetables, fruits, nuts, and seeds.",
					"ar": "اللحوم والأسماك والبيض والخضروات والفواكه والمكسرات والبذور.",
				},
				"minimal_foods": map[string]string{
					"en": "Grains, legumes, dairy, refined sugar, and processed foods.",
					"ar": "الحبوب والبقوليات ومنتجات الألبان والسكر المكرر والأطعمة المصنعة.",
				},
				"type": "whole_foods",
				"difficulty": "moderate",
			},
			{
				"name": map[string]string{
					"en": "DASH Diet",
					"ar": "نظام داش الغذائي",
				},
				"allowed_foods": map[string]string{
					"en": "Fruits, vegetables, whole grains, lean proteins, low-fat dairy, nuts, and seeds.",
					"ar": "الفواكه والخضروات والحبوب الكاملة والبروتينات الخالية من الدهون ومنتجات الألبان قليلة الدسم والمكسرات والبذور.",
				},
				"minimal_foods": map[string]string{
					"en": "High-sodium foods, red meat, sweets, and sugary beverages.",
					"ar": "الأطعمة عالية الصوديوم واللحوم الحمراء والحلويات والمشروبات السكرية.",
				},
				"type": "heart_healthy",
				"difficulty": "easy",
			},
		},
	}

	return c.JSON(http.StatusOK, NutritionResponse{
		Success: true,
		Message: "Type plans data retrieved successfully",
		Data:    typePlansData,
		Meta: &MetaInfo{
			Total:      len(typePlansData["diets"].([]map[string]interface{})),
			Page:       page,
			PerPage:    perPage,
			TotalPages: 1,
			Language:   "en",
			Country:    "US",
			APIVersion: "1.0",
			CacheExpiry: 3600,
		},
	})
}

// GetTypePlansByType retrieves diet plans filtered by type
func GetTypePlansByType(c echo.Context) error {
	if err := validateNutritionAccess(c, "nutrition_read"); err != nil {
		return err
	}

	dietType := c.Param("type")
	if dietType == "" {
		return c.JSON(http.StatusBadRequest, NutritionResponse{
			Success: false,
			Message: "Diet type parameter is required",
		})
	}

	// Mock filtered diet plans data
	var filteredPlans []map[string]interface{}
	allPlans := []map[string]interface{}{
		{
			"name": map[string]string{
				"en": "Mediterranean Diet",
				"ar": "النظام الغذائي المتوسطي",
			},
			"type": "balanced",
		},
		{
			"name": map[string]string{
				"en": "Ketogenic Diet",
				"ar": "النظام الغذائي الكيتوني",
			},
			"type": "low_carb",
		},
		{
			"name": map[string]string{
				"en": "Vegetarian Diet",
				"ar": "النظام الغذائي النباتي",
			},
			"type": "plant_based",
		},
	}

	for _, plan := range allPlans {
		if plan["type"] == dietType {
			filteredPlans = append(filteredPlans, plan)
		}
	}

	if len(filteredPlans) == 0 {
		return c.JSON(http.StatusNotFound, NutritionResponse{
			Success: false,
			Message: "No diet plans found for the specified type",
		})
	}

	return c.JSON(http.StatusOK, NutritionResponse{
		Success: true,
		Message: "Type plans data retrieved successfully",
		Data:    map[string]interface{}{"diets": filteredPlans},
		Meta: &MetaInfo{
			Total:      len(filteredPlans),
			Page:       1,
			PerPage:    20,
			TotalPages: 1,
			Language:   "en",
			Country:    "US",
			APIVersion: "1.0",
			CacheExpiry: 3600,
		},
	})
}

// Helper function to validate API key scopes for nutrition endpoints
func validateNutritionAccess(c echo.Context, requiredScope string) error {
	// Check if user has required scope for nutrition data access
	apiKeyID := c.Get("api_key_id")
	if apiKeyID == nil {
		// No API key, check if user is authenticated
		userID := c.Get("user_id")
		if userID == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
		}
		return nil
	}
	
	// API key is present, check scopes
	scopes := c.Get("api_key_scopes")
	if scopes == nil {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid API key scopes")
	}
	
	scopeList, ok := scopes.([]string)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid scope format")
	}
	
	if !contains(scopeList, requiredScope) && !contains(scopeList, "admin") {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
	}
	
	return nil
}