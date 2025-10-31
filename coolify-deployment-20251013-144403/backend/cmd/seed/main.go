package main

import (
	"fmt"
	"log"
	"time"

	"nutrition-platform/config"
	"nutrition-platform/models"
	"nutrition-platform/services"
)

// SeedData contains initial data for the application
type SeedData struct {
	APIKeys []APIKeySeed `json:"api_keys"`
	Users   []UserSeed   `json:"users"`
	Foods   []FoodSeed   `json:"foods"`
}

type APIKeySeed struct {
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	RateLimit int      `json:"rate_limit"`
}

type UserSeed struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FoodSeed struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	CaloriesPer100g float64 `json:"calories_per_100g"`
	ProteinPer100g  float64 `json:"protein_per_100g"`
	CarbsPer100g    float64 `json:"carbs_per_100g"`
	FatPer100g      float64 `json:"fat_per_100g"`
	IsHalal         bool    `json:"is_halal"`
	IsVegetarian    bool    `json:"is_vegetarian"`
	IsVegan         bool    `json:"is_vegan"`
}

func main() {
	log.Println("Starting database seeding...")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := models.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDatabase()

	// Initialize services
	apiKeyService := services.NewAPIKeyService(models.DB)

	// Seed data
	if err := seedAPIKeys(apiKeyService); err != nil {
		log.Fatalf("Failed to seed API keys: %v", err)
	}

	if err := seedFoods(); err != nil {
		log.Fatalf("Failed to seed foods: %v", err)
	}

	if err := seedRecipes(); err != nil {
		log.Fatalf("Failed to seed recipes: %v", err)
	}

	if err := seedHealthConditions(); err != nil {
		log.Fatalf("Failed to seed health conditions: %v", err)
	}

	if err := seedInjuries(); err != nil {
		log.Fatalf("Failed to seed injuries: %v", err)
	}

	if err := seedMedications(); err != nil {
		log.Fatalf("Failed to seed medications: %v", err)
	}

	if err := seedVitaminsMinerals(); err != nil {
		log.Fatalf("Failed to seed vitamins and minerals: %v", err)
	}

	if err := seedWorkoutPrograms(); err != nil {
		log.Fatalf("Failed to seed workout programs: %v", err)
	}

	if err := seedNutritionPlans(); err != nil {
		log.Fatalf("Failed to seed nutrition plans: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

func seedAPIKeys(apiKeyService *services.APIKeyService) error {
	log.Println("Seeding API keys...")

	// Create admin API key
	adminReq := &models.CreateAPIKeyRequest{
		Name:      "Admin Key",
		Scopes:    []models.APIKeyScope{models.ScopeAdmin},
		RateLimit: 1000,
	}

	adminResponse, err := apiKeyService.CreateAPIKey("admin-user", adminReq)
	if err != nil {
		return err
	}

	log.Printf("Created admin API key: %s", adminResponse.Key)

	// Create read-only API key
	readOnlyReq := &models.CreateAPIKeyRequest{
		Name:      "Read Only Key",
		Scopes:    []models.APIKeyScope{models.ScopeReadOnly, models.ScopeNutrition},
		RateLimit: 100,
	}

	readOnlyResponse, err := apiKeyService.CreateAPIKey("readonly-user", readOnlyReq)
	if err != nil {
		return err
	}

	log.Printf("Created read-only API key: %s", readOnlyResponse.Key)

	return nil
}

func seedFoods() error {
	log.Println("Seeding foods...")

	foods := []FoodSeed{
		{
			Name:            "Chicken Breast",
			Category:        "Meat",
			CaloriesPer100g: 165,
			ProteinPer100g:  31,
			CarbsPer100g:    0,
			FatPer100g:      3.6,
			IsHalal:         true,
			IsVegetarian:    false,
			IsVegan:         false,
		},
		{
			Name:            "Brown Rice",
			Category:        "Grains",
			CaloriesPer100g: 111,
			ProteinPer100g:  2.6,
			CarbsPer100g:    23,
			FatPer100g:      0.9,
			IsHalal:         true,
			IsVegetarian:    true,
			IsVegan:         true,
		},
		{
			Name:            "Broccoli",
			Category:        "Vegetables",
			CaloriesPer100g: 34,
			ProteinPer100g:  2.8,
			CarbsPer100g:    7,
			FatPer100g:      0.4,
			IsHalal:         true,
			IsVegetarian:    true,
			IsVegan:         true,
		},
		{
			Name:            "Salmon",
			Category:        "Fish",
			CaloriesPer100g: 208,
			ProteinPer100g:  25,
			CarbsPer100g:    0,
			FatPer100g:      13,
			IsHalal:         true,
			IsVegetarian:    false,
			IsVegan:         false,
		},
		{
			Name:            "Quinoa",
			Category:        "Grains",
			CaloriesPer100g: 120,
			ProteinPer100g:  4.4,
			CarbsPer100g:    22,
			FatPer100g:      1.9,
			IsHalal:         true,
			IsVegetarian:    true,
			IsVegan:         true,
		},
		{
			Name:            "Greek Yogurt",
			Category:        "Dairy",
			CaloriesPer100g: 59,
			ProteinPer100g:  10,
			CarbsPer100g:    3.6,
			FatPer100g:      0.4,
			IsHalal:         true,
			IsVegetarian:    true,
			IsVegan:         false,
		},
		{
			Name:            "Almonds",
			Category:        "Nuts",
			CaloriesPer100g: 579,
			ProteinPer100g:  21,
			CarbsPer100g:    22,
			FatPer100g:      50,
			IsHalal:         true,
			IsVegetarian:    true,
			IsVegan:         true,
		},
		{
			Name:            "Sweet Potato",
			Category:        "Vegetables",
			CaloriesPer100g: 86,
			ProteinPer100g:  1.6,
			CarbsPer100g:    20,
			FatPer100g:      0.1,
			IsHalal:         true,
			IsVegetarian:    true,
			IsVegan:         true,
		},
	}

	for _, food := range foods {
		query := `
			INSERT INTO foods (name, category, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g, is_halal, is_vegetarian, is_vegan, verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, $10, $10)
			ON CONFLICT (name) DO NOTHING
		`

		_, err := models.DB.Exec(query,
			food.Name,
			food.Category,
			food.CaloriesPer100g,
			food.ProteinPer100g,
			food.CarbsPer100g,
			food.FatPer100g,
			food.IsHalal,
			food.IsVegetarian,
			food.IsVegan,
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded food: %s", food.Name)
	}

	return nil
}
func seedRecipes() error {
	log.Println("Seeding recipes...")

	recipes := []map[string]interface{}{
		{
			"name":              "Mediterranean Grilled Chicken",
			"cuisine":           "Mediterranean",
			"country":           "Greece",
			"difficulty_level":  "medium",
			"prep_time_minutes": 20,
			"cook_time_minutes": 25,
			"servings":          4,
			"ingredients":       `[{"name":"Chicken breast","amount":4,"unit":"pieces"},{"name":"Olive oil","amount":3,"unit":"tbsp"},{"name":"Lemon juice","amount":2,"unit":"tbsp"},{"name":"Garlic","amount":3,"unit":"cloves"},{"name":"Oregano","amount":1,"unit":"tsp"}]`,
			"instructions":      `[{"step_number":1,"instruction":"Marinate chicken in olive oil, lemon juice, garlic, and oregano for 30 minutes"},{"step_number":2,"instruction":"Preheat grill to medium-high heat"},{"step_number":3,"instruction":"Grill chicken for 6-7 minutes per side until cooked through"}]`,
			"dietary_tags":      `["gluten-free","dairy-free","high-protein"]`,
			"is_halal":          true,
			"verified":          true,
		},
		{
			"name":              "Quinoa Buddha Bowl",
			"cuisine":           "International",
			"difficulty_level":  "easy",
			"prep_time_minutes": 15,
			"cook_time_minutes": 20,
			"servings":          2,
			"ingredients":       `[{"name":"Quinoa","amount":1,"unit":"cup"},{"name":"Sweet potato","amount":1,"unit":"large"},{"name":"Chickpeas","amount":1,"unit":"can"},{"name":"Spinach","amount":2,"unit":"cups"},{"name":"Avocado","amount":1,"unit":"piece"}]`,
			"instructions":      `[{"step_number":1,"instruction":"Cook quinoa according to package instructions"},{"step_number":2,"instruction":"Roast sweet potato cubes at 400°F for 20 minutes"},{"step_number":3,"instruction":"Assemble bowl with quinoa, roasted sweet potato, chickpeas, spinach, and avocado"}]`,
			"dietary_tags":      `["vegan","vegetarian","gluten-free","high-fiber"]`,
			"is_halal":          true,
			"is_vegetarian":     true,
			"is_vegan":          true,
			"verified":          true,
		},
		{
			"name":              "Salmon Teriyaki",
			"cuisine":           "Japanese",
			"country":           "Japan",
			"difficulty_level":  "medium",
			"prep_time_minutes": 10,
			"cook_time_minutes": 15,
			"servings":          4,
			"ingredients":       `[{"name":"Salmon fillets","amount":4,"unit":"pieces"},{"name":"Soy sauce","amount":3,"unit":"tbsp"},{"name":"Mirin","amount":2,"unit":"tbsp"},{"name":"Sugar","amount":1,"unit":"tbsp"},{"name":"Ginger","amount":1,"unit":"tsp"}]`,
			"instructions":      `[{"step_number":1,"instruction":"Mix soy sauce, mirin, sugar, and ginger for teriyaki sauce"},{"step_number":2,"instruction":"Pan-fry salmon fillets for 4-5 minutes per side"},{"step_number":3,"instruction":"Brush with teriyaki sauce and cook for 2 more minutes"}]`,
			"dietary_tags":      `["gluten-free","high-protein","omega-3"]`,
			"is_halal":          true,
			"verified":          true,
		},
	}

	for _, recipe := range recipes {
		query := `
			INSERT INTO recipes (
				id, name, cuisine, country, difficulty_level, prep_time_minutes, 
				cook_time_minutes, total_time_minutes, servings, ingredients, 
				instructions, dietary_tags, allergens, is_halal, is_kosher, 
				is_vegetarian, is_vegan, rating, rating_count, verified, 
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
			) ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("recipe_%d", time.Now().UnixNano())
		totalTime := recipe["prep_time_minutes"].(int) + recipe["cook_time_minutes"].(int)

		_, err := models.DB.Exec(query,
			id,
			recipe["name"],
			recipe["cuisine"],
			recipe["country"],
			recipe["difficulty_level"],
			recipe["prep_time_minutes"],
			recipe["cook_time_minutes"],
			totalTime,
			recipe["servings"],
			recipe["ingredients"],
			recipe["instructions"],
			recipe["dietary_tags"],
			`[]`, // allergens
			recipe["is_halal"],
			false, // is_kosher
			recipe["is_vegetarian"] != nil && recipe["is_vegetarian"].(bool),
			recipe["is_vegan"] != nil && recipe["is_vegan"].(bool),
			4.5, // rating
			10,  // rating_count
			recipe["verified"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded recipe: %s", recipe["name"])
	}

	return nil
}

func seedHealthConditions() error {
	log.Println("Seeding health conditions...")

	conditions := []map[string]interface{}{
		{
			"name":                         "Type 2 Diabetes",
			"category":                     "metabolic",
			"icd_10_code":                  "E11",
			"description":                  "A chronic condition that affects the way the body processes blood sugar (glucose)",
			"symptoms":                     `["increased thirst","frequent urination","increased hunger","fatigue","blurred vision"]`,
			"dietary_recommendations":      `["limit refined carbohydrates","choose whole grains","eat regular meals","monitor portion sizes","include fiber-rich foods"]`,
			"exercise_recommendations":     `["150 minutes moderate aerobic activity per week","strength training 2+ times per week","monitor blood sugar before and after exercise"]`,
			"is_chronic":                   true,
			"requires_medical_supervision": true,
		},
		{
			"name":                         "Hypertension",
			"category":                     "cardiovascular",
			"icd_10_code":                  "I10",
			"description":                  "High blood pressure that can lead to serious health complications",
			"symptoms":                     `["headaches","shortness of breath","nosebleeds","chest pain"]`,
			"dietary_recommendations":      `["reduce sodium intake","increase potassium-rich foods","limit alcohol","maintain healthy weight","eat plenty of fruits and vegetables"]`,
			"exercise_recommendations":     `["regular aerobic exercise","aim for 30 minutes most days","include resistance training","start slowly if sedentary"]`,
			"is_chronic":                   true,
			"requires_medical_supervision": true,
		},
		{
			"name":                         "Iron Deficiency Anemia",
			"category":                     "nutritional",
			"icd_10_code":                  "D50",
			"description":                  "A condition where blood lacks adequate healthy red blood cells due to iron deficiency",
			"symptoms":                     `["fatigue","weakness","pale skin","shortness of breath","cold hands and feet"]`,
			"dietary_recommendations":      `["increase iron-rich foods","combine iron with vitamin C","avoid tea and coffee with meals","include lean meats and leafy greens"]`,
			"exercise_recommendations":     `["light to moderate exercise","avoid intense exercise until iron levels improve","focus on gentle activities like walking"]`,
			"is_chronic":                   false,
			"requires_medical_supervision": true,
		},
	}

	for _, condition := range conditions {
		query := `
			INSERT INTO health_conditions (
				id, name, category, icd_10_code, description, symptoms,
				dietary_recommendations, exercise_recommendations, 
				is_chronic, requires_medical_supervision, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("condition_%d", time.Now().UnixNano())

		_, err := models.DB.Exec(query,
			id,
			condition["name"],
			condition["category"],
			condition["icd_10_code"],
			condition["description"],
			condition["symptoms"],
			condition["dietary_recommendations"],
			condition["exercise_recommendations"],
			condition["is_chronic"],
			condition["requires_medical_supervision"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded health condition: %s", condition["name"])
	}

	return nil
}

func seedInjuries() error {
	log.Println("Seeding injuries...")

	injuries := []map[string]interface{}{
		{
			"name":                  "Lower Back Strain",
			"category":              "musculoskeletal",
			"body_part":             "lower back",
			"severity_level":        "moderate",
			"description":           "Injury to muscles, ligaments, or tendons in the lower back",
			"symptoms":              `["lower back pain","muscle spasms","stiffness","limited range of motion"]`,
			"causes":                `["heavy lifting","sudden movements","poor posture","muscle imbalance"]`,
			"treatment_options":     `["rest","ice/heat therapy","gentle stretching","physical therapy","pain medication"]`,
			"recovery_time_days":    14,
			"exercise_restrictions": `["avoid heavy lifting","no twisting movements","limit bending"]`,
			"recommended_exercises": `["gentle walking","pelvic tilts","knee-to-chest stretches","cat-cow stretches"]`,
		},
		{
			"name":                  "Ankle Sprain",
			"category":              "sports",
			"body_part":             "ankle",
			"severity_level":        "mild",
			"description":           "Stretching or tearing of ligaments in the ankle",
			"symptoms":              `["ankle pain","swelling","bruising","difficulty walking"]`,
			"causes":                `["rolling ankle","uneven surfaces","sports activities","weak ankle muscles"]`,
			"treatment_options":     `["RICE protocol","elevation","compression","gradual return to activity"]`,
			"recovery_time_days":    7,
			"exercise_restrictions": `["avoid running","no jumping","limit weight bearing initially"]`,
			"recommended_exercises": `["ankle circles","calf raises","balance exercises","swimming"]`,
		},
		{
			"name":                  "Tennis Elbow",
			"category":              "overuse",
			"body_part":             "elbow",
			"severity_level":        "moderate",
			"description":           "Inflammation of tendons on the outside of the elbow",
			"symptoms":              `["elbow pain","weakness in grip","pain when lifting","tenderness on outer elbow"]`,
			"causes":                `["repetitive arm motions","tennis","poor technique","muscle imbalance"]`,
			"treatment_options":     `["rest","ice","anti-inflammatory medication","physical therapy","elbow strap"]`,
			"recovery_time_days":    21,
			"exercise_restrictions": `["avoid gripping activities","no tennis/racquet sports","limit lifting"]`,
			"recommended_exercises": `["gentle stretching","eccentric strengthening","wrist flexor stretches"]`,
		},
	}

	for _, injury := range injuries {
		query := `
			INSERT INTO injuries (
				id, name, category, body_part, severity_level, description,
				symptoms, causes, treatment_options, recovery_time_days,
				exercise_restrictions, recommended_exercises, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("injury_%d", time.Now().UnixNano())

		_, err := models.DB.Exec(query,
			id,
			injury["name"],
			injury["category"],
			injury["body_part"],
			injury["severity_level"],
			injury["description"],
			injury["symptoms"],
			injury["causes"],
			injury["treatment_options"],
			injury["recovery_time_days"],
			injury["exercise_restrictions"],
			injury["recommended_exercises"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded injury: %s", injury["name"])
	}

	return nil
}

func seedMedications() error {
	log.Println("Seeding medications...")

	medications := []map[string]interface{}{
		{
			"name":                  "Metformin",
			"generic_name":          "Metformin Hydrochloride",
			"drug_class":            "Biguanide",
			"category":              "prescription",
			"description":           "Medication used to treat type 2 diabetes",
			"indications":           `["type 2 diabetes","prediabetes","PCOS"]`,
			"side_effects":          `["nausea","diarrhea","stomach upset","metallic taste"]`,
			"food_interactions":     `["take with meals to reduce stomach upset","avoid excessive alcohol"]`,
			"affects_nutrition":     true,
			"nutritional_effects":   `{"vitamin_b12": "may decrease absorption", "folate": "may decrease levels"}`,
			"requires_prescription": true,
		},
		{
			"name":                  "Ibuprofen",
			"generic_name":          "Ibuprofen",
			"drug_class":            "NSAID",
			"category":              "otc",
			"description":           "Nonsteroidal anti-inflammatory drug for pain and inflammation",
			"indications":           `["pain relief","inflammation","fever reduction","headaches"]`,
			"side_effects":          `["stomach upset","heartburn","dizziness","drowsiness"]`,
			"food_interactions":     `["take with food to reduce stomach irritation","avoid alcohol"]`,
			"affects_nutrition":     false,
			"requires_prescription": false,
		},
		{
			"name":                  "Lisinopril",
			"generic_name":          "Lisinopril",
			"drug_class":            "ACE Inhibitor",
			"category":              "prescription",
			"description":           "Medication used to treat high blood pressure and heart failure",
			"indications":           `["hypertension","heart failure","post-heart attack"]`,
			"side_effects":          `["dry cough","dizziness","fatigue","hyperkalemia"]`,
			"food_interactions":     `["avoid salt substitutes with potassium","limit high-potassium foods"]`,
			"affects_nutrition":     true,
			"nutritional_effects":   `{"potassium": "may increase levels", "sodium": "monitor intake"}`,
			"requires_prescription": true,
		},
	}

	for _, medication := range medications {
		query := `
			INSERT INTO medications (
				id, name, generic_name, drug_class, category, description,
				indications, side_effects, food_interactions, affects_nutrition,
				nutritional_effects, requires_prescription, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("med_%d", time.Now().UnixNano())

		_, err := models.DB.Exec(query,
			id,
			medication["name"],
			medication["generic_name"],
			medication["drug_class"],
			medication["category"],
			medication["description"],
			medication["indications"],
			medication["side_effects"],
			medication["food_interactions"],
			medication["affects_nutrition"],
			medication["nutritional_effects"],
			medication["requires_prescription"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded medication: %s", medication["name"])
	}

	return nil
}

func seedVitaminsMinerals() error {
	log.Println("Seeding vitamins and minerals...")

	nutrients := []map[string]interface{}{
		{
			"name":                "Vitamin D",
			"type":                "vitamin",
			"category":            "fat_soluble",
			"description":         "Essential vitamin for bone health and immune function",
			"functions":           `["bone health","immune function","calcium absorption","muscle function"]`,
			"deficiency_symptoms": `["bone pain","muscle weakness","fatigue","increased infections"]`,
			"food_sources":        `["fatty fish","egg yolks","fortified milk","mushrooms","sunlight exposure"]`,
			"daily_requirements":  `{"adults": "600-800 IU", "elderly": "800-1000 IU"}`,
			"upper_limit":         `{"adults": "4000 IU"}`,
			"best_taken_with":     `["fat","calcium"]`,
		},
		{
			"name":                "Iron",
			"type":                "mineral",
			"category":            "trace_mineral",
			"description":         "Essential mineral for oxygen transport and energy production",
			"functions":           `["oxygen transport","energy production","immune function","cognitive function"]`,
			"deficiency_symptoms": `["fatigue","weakness","pale skin","shortness of breath","cold hands and feet"]`,
			"food_sources":        `["red meat","poultry","fish","beans","spinach","fortified cereals"]`,
			"daily_requirements":  `{"men": "8 mg", "women": "18 mg", "pregnant": "27 mg"}`,
			"upper_limit":         `{"adults": "45 mg"}`,
			"best_taken_with":     `["vitamin C","citrus fruits"]`,
			"avoid_taking_with":   `["calcium","tea","coffee"]`,
		},
		{
			"name":                "Omega-3 Fatty Acids",
			"type":                "fatty_acid",
			"category":            "essential_fatty_acid",
			"description":         "Essential fatty acids important for heart and brain health",
			"functions":           `["heart health","brain function","inflammation reduction","eye health"]`,
			"deficiency_symptoms": `["dry skin","fatigue","poor memory","mood swings","heart problems"]`,
			"food_sources":        `["fatty fish","walnuts","flaxseeds","chia seeds","algae oil"]`,
			"daily_requirements":  `{"adults": "1.1-1.6 g ALA", "EPA+DHA": "250-500 mg"}`,
			"best_taken_with":     `["meals","fat"]`,
		},
	}

	for _, nutrient := range nutrients {
		query := `
			INSERT INTO vitamins_minerals (
				id, name, type, category, description, functions,
				deficiency_symptoms, food_sources, daily_requirements,
				upper_limit, best_taken_with, avoid_taking_with,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("nutrient_%d", time.Now().UnixNano())

		_, err := models.DB.Exec(query,
			id,
			nutrient["name"],
			nutrient["type"],
			nutrient["category"],
			nutrient["description"],
			nutrient["functions"],
			nutrient["deficiency_symptoms"],
			nutrient["food_sources"],
			nutrient["daily_requirements"],
			nutrient["upper_limit"],
			nutrient["best_taken_with"],
			nutrient["avoid_taking_with"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded nutrient: %s", nutrient["name"])
	}

	return nil
}

func seedWorkoutPrograms() error {
	log.Println("Seeding workout programs...")

	programs := []map[string]interface{}{
		{
			"name":                     "Beginner Full Body Workout",
			"program_type":             "strength",
			"fitness_level":            "beginner",
			"duration_weeks":           8,
			"days_per_week":            3,
			"session_duration_minutes": 45,
			"equipment_required":       `["dumbbells","resistance bands","mat"]`,
			"target_goals":             `["muscle building","strength","general fitness"]`,
			"muscle_groups_targeted":   `["full body","core","legs","arms","back","chest"]`,
			"difficulty_rating":        2,
			"calorie_burn_estimate":    300,
		},
		{
			"name":                     "HIIT Cardio Blast",
			"program_type":             "cardio",
			"fitness_level":            "intermediate",
			"duration_weeks":           6,
			"days_per_week":            4,
			"session_duration_minutes": 30,
			"equipment_required":       `["none","bodyweight"]`,
			"target_goals":             `["weight loss","cardiovascular fitness","endurance"]`,
			"muscle_groups_targeted":   `["full body","core","legs"]`,
			"difficulty_rating":        4,
			"calorie_burn_estimate":    400,
		},
		{
			"name":                     "Yoga Flow for Flexibility",
			"program_type":             "flexibility",
			"fitness_level":            "beginner",
			"duration_weeks":           12,
			"days_per_week":            5,
			"session_duration_minutes": 60,
			"equipment_required":       `["yoga mat","blocks","strap"]`,
			"target_goals":             `["flexibility","stress relief","balance","mindfulness"]`,
			"muscle_groups_targeted":   `["full body","core","back","hips"]`,
			"difficulty_rating":        2,
			"calorie_burn_estimate":    200,
		},
	}

	for _, program := range programs {
		query := `
			INSERT INTO workout_programs (
				id, name, program_type, fitness_level, duration_weeks,
				days_per_week, session_duration_minutes, equipment_required,
				target_goals, muscle_groups_targeted, difficulty_rating,
				calorie_burn_estimate, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("program_%d", time.Now().UnixNano())

		_, err := models.DB.Exec(query,
			id,
			program["name"],
			program["program_type"],
			program["fitness_level"],
			program["duration_weeks"],
			program["days_per_week"],
			program["session_duration_minutes"],
			program["equipment_required"],
			program["target_goals"],
			program["muscle_groups_targeted"],
			program["difficulty_rating"],
			program["calorie_burn_estimate"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded workout program: %s", program["name"])
	}

	return nil
}
func seedNutritionPlans() error {
	log.Println("Seeding nutrition plans...")

	plans := []map[string]interface{}{
		{
			"name":                      "Mediterranean Weight Loss Plan",
			"plan_type":                 "mediterranean",
			"description":               "Mediterranean diet adapted for healthy weight loss",
			"duration_weeks":            12,
			"daily_calorie_target":      1800,
			"macro_targets":             `{"protein_grams": 135, "protein_percent": 15, "carbs_grams": 248, "carbs_percent": 55, "fat_grams": 60, "fat_percent": 30, "fiber_grams": 35, "sodium_mg": 2300}`,
			"meal_timing":               `[{"meal_type": "breakfast", "recommended_time": "7:00-9:00 AM", "calorie_percent": 25}, {"meal_type": "lunch", "recommended_time": "12:00-2:00 PM", "calorie_percent": 35}, {"meal_type": "dinner", "recommended_time": "6:00-8:00 PM", "calorie_percent": 30}, {"meal_type": "snack", "recommended_time": "3:00-4:00 PM", "calorie_percent": 10}]`,
			"recommended_foods":         `["olive oil", "fish", "vegetables", "fruits", "whole grains", "nuts", "legumes"]`,
			"foods_to_avoid":            `["processed foods", "refined sugars", "trans fats", "excessive red meat"]`,
			"hydration_target_ml":       2500,
			"medical_approval_required": false,
			"is_active":                 true,
		},
		{
			"name":                      "DASH Hypertension Management Plan",
			"plan_type":                 "dash",
			"description":               "DASH diet specifically designed for blood pressure management",
			"duration_weeks":            16,
			"daily_calorie_target":      2000,
			"macro_targets":             `{"protein_grams": 90, "protein_percent": 18, "carbs_grams": 275, "carbs_percent": 55, "fat_grams": 60, "fat_percent": 27, "fiber_grams": 30, "sodium_mg": 1500}`,
			"meal_timing":               `[{"meal_type": "breakfast", "recommended_time": "7:00-8:00 AM", "calorie_percent": 25}, {"meal_type": "lunch", "recommended_time": "12:00-1:00 PM", "calorie_percent": 30}, {"meal_type": "dinner", "recommended_time": "6:00-7:00 PM", "calorie_percent": 35}, {"meal_type": "snack", "recommended_time": "3:00 PM", "calorie_percent": 10}]`,
			"recommended_foods":         `["low-fat dairy", "fruits", "vegetables", "whole grains", "lean proteins", "nuts", "seeds"]`,
			"foods_to_avoid":            `["high sodium foods", "processed meats", "sugary drinks", "excessive alcohol"]`,
			"hydration_target_ml":       2000,
			"medical_approval_required": false,
			"is_active":                 true,
		},
		{
			"name":                      "Ketogenic Metabolic Reset Plan",
			"plan_type":                 "ketogenic",
			"description":               "Medically supervised ketogenic diet for metabolic health",
			"duration_weeks":            8,
			"daily_calorie_target":      1600,
			"macro_targets":             `{"protein_grams": 80, "protein_percent": 20, "carbs_grams": 20, "carbs_percent": 5, "fat_grams": 133, "fat_percent": 75, "fiber_grams": 25, "sodium_mg": 3000}`,
			"meal_timing":               `[{"meal_type": "breakfast", "recommended_time": "8:00-10:00 AM", "calorie_percent": 30}, {"meal_type": "lunch", "recommended_time": "1:00-3:00 PM", "calorie_percent": 35}, {"meal_type": "dinner", "recommended_time": "6:00-8:00 PM", "calorie_percent": 35}]`,
			"recommended_foods":         `["fatty fish", "meat", "eggs", "cheese", "nuts", "seeds", "low-carb vegetables", "healthy oils"]`,
			"foods_to_avoid":            `["grains", "sugar", "fruits", "starchy vegetables", "legumes", "most dairy"]`,
			"hydration_target_ml":       3000,
			"special_instructions":      "Monitor ketone levels daily. Supplement with electrolytes. Expect adaptation period of 2-4 weeks.",
			"medical_approval_required": true,
			"is_active":                 true,
		},
	}

	for _, plan := range plans {
		query := `
			INSERT INTO nutritional_plans (
				id, name, plan_type, description, duration_weeks,
				daily_calorie_target, macro_targets, meal_timing,
				recommended_foods, foods_to_avoid, hydration_target_ml,
				special_instructions, medical_approval_required, is_active,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			ON CONFLICT (name) DO NOTHING
		`

		id := fmt.Sprintf("plan_%d", time.Now().UnixNano())

		_, err := models.DB.Exec(query,
			id,
			plan["name"],
			plan["plan_type"],
			plan["description"],
			plan["duration_weeks"],
			plan["daily_calorie_target"],
			plan["macro_targets"],
			plan["meal_timing"],
			plan["recommended_foods"],
			plan["foods_to_avoid"],
			plan["hydration_target_ml"],
			plan["special_instructions"],
			plan["medical_approval_required"],
			plan["is_active"],
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		}

		log.Printf("Seeded nutrition plan: %s", plan["name"])
	}

	return nil
}
