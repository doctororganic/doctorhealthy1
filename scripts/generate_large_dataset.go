package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Large dataset generator for 10,000 users load testing
type DatasetGenerator struct {
	OutputDir    string
	UserCount    int
	WorkoutsPerUser int
	RecipesPerUser  int
	MealsPerUser    int
}

// Enhanced workout structure for large dataset
type LargeWorkout struct {
	ID          string                 `json:"id"`
	Name        map[string]string      `json:"name"`
	Type        string                 `json:"type"`
	Goal        string                 `json:"goal"`
	Level       string                 `json:"experience_level"`
	Duration    int                    `json:"duration"`
	Exercises   []Exercise             `json:"exercises"`
	Equipment   []string               `json:"equipment"`
	Calories    int                    `json:"calories_per_hour"`
	TargetMuscles []string             `json:"target_muscles"`
	Health      HealthInfo             `json:"health_conditions"`
	Created     time.Time              `json:"created_at"`
	UserIDs     []int                  `json:"suitable_for_users"`
}

type Exercise struct {
	Name        map[string]string `json:"name"`
	Sets        int              `json:"sets"`
	Reps        string           `json:"reps"`
	Duration    int              `json:"duration_seconds"`
	Rest        int              `json:"rest_seconds"`
	Difficulty  int              `json:"difficulty"`
	Instructions map[string]string `json:"instructions"`
}

type HealthInfo struct {
	SuitableFor []string `json:"suitable_for"`
	AvoidIf     []string `json:"avoid_if"`
	Modifications []string `json:"modifications"`
}

// Enhanced recipe structure for large dataset
type LargeRecipe struct {
	ID           string                 `json:"id"`
	Name         map[string]string      `json:"name"`
	Cuisine      string                 `json:"cuisine"`
	DietType     string                 `json:"diet_type"`
	Ingredients  []Ingredient           `json:"ingredients"`
	Instructions []map[string]string    `json:"instructions"`
	Nutrition    NutritionInfo          `json:"nutrition"`
	PrepTime     int                    `json:"prep_time"`
	CookTime     int                    `json:"cook_time"`
	Servings     int                    `json:"servings"`
	IsHalal      bool                   `json:"is_halal"`
	Allergens    []string               `json:"allergens"`
	Tags         []string               `json:"tags"`
	Created      time.Time              `json:"created_at"`
	UserRatings  []UserRating           `json:"user_ratings"`
}

type Ingredient struct {
	Name         map[string]string `json:"name"`
	Amount       float64           `json:"amount"`
	Unit         string            `json:"unit"`
	IsOptional   bool              `json:"optional"`
	HalalSub     string            `json:"halal_substitute,omitempty"`
}

type NutritionInfo struct {
	Calories     int     `json:"calories"`
	Protein      float64 `json:"protein"`
	Carbs        float64 `json:"carbs"`
	Fat          float64 `json:"fat"`
	Fiber        float64 `json:"fiber"`
	Sugar        float64 `json:"sugar"`
	Sodium       float64 `json:"sodium"`
	Vitamins     map[string]float64 `json:"vitamins"`
	Minerals     map[string]float64 `json:"minerals"`
}

type UserRating struct {
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Date      time.Time `json:"date"`
}

// User profile for personalization
type UserProfile struct {
	ID               int                    `json:"id"`
	Name             string                 `json:"name"`
	Age              int                    `json:"age"`
	Gender           string                 `json:"gender"`
	Weight           float64                `json:"weight_kg"`
	Height           float64                `json:"height_cm"`
	ActivityLevel    string                 `json:"activity_level"`
	Goals            []string               `json:"goals"`
	DietType         string                 `json:"diet_type"`
	HealthConditions []string               `json:"health_conditions"`
	Allergies        []string               `json:"allergies"`
	Language         string                 `json:"language"`
	Country          string                 `json:"country"`
	BMR              float64                `json:"bmr"`
	TDEE             float64                `json:"tdee"`
	MacroTargets     map[string]float64     `json:"macro_targets"`
	CreatedAt        time.Time              `json:"created_at"`
	LastActive       time.Time              `json:"last_active"`
}

func NewDatasetGenerator(outputDir string) *DatasetGenerator {
	return &DatasetGenerator{
		OutputDir:       outputDir,
		UserCount:       10000,
		WorkoutsPerUser: 50,   // 500,000 total workouts
		RecipesPerUser:  30,   // 300,000 total recipes
		MealsPerUser:    100,  // 1,000,000 total meals
	}
}

func (dg *DatasetGenerator) GenerateAllData() error {
	fmt.Printf("🎯 Generating large dataset for %d users...\n", dg.UserCount)
	
	// Create output directory
	if err := os.MkdirAll(dg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate user profiles first
	users, err := dg.generateUsers()
	if err != nil {
		return fmt.Errorf("failed to generate users: %w", err)
	}

	// Generate workouts with user associations
	if err := dg.generateWorkouts(users); err != nil {
		return fmt.Errorf("failed to generate workouts: %w", err)
	}

	// Generate recipes with user ratings
	if err := dg.generateRecipes(users); err != nil {
		return fmt.Errorf("failed to generate recipes: %w", err)
	}

	// Generate meal plans
	if err := dg.generateMealPlans(users); err != nil {
		return fmt.Errorf("failed to generate meal plans: %w", err)
	}

	// Generate health data
	if err := dg.generateHealthData(users); err != nil {
		return fmt.Errorf("failed to generate health data: %w", err)
	}

	fmt.Printf("✅ Successfully generated dataset for %d users\n", len(users))
	return nil
}

func (dg *DatasetGenerator) generateUsers() ([]UserProfile, error) {
	fmt.Printf("👥 Generating %d user profiles...\n", dg.UserCount)
	
	users := make([]UserProfile, dg.UserCount)
	
	// Demographics for realistic distribution
	names := []string{"Ahmed", "Fatima", "Mohammad", "Aisha", "Ali", "Sarah", "Omar", "Layla", "Hassan", "Zara"}
	countries := []string{"Egypt", "UAE", "Saudi Arabia", "Jordan", "Lebanon", "USA", "UK", "Canada"}
	goals := [][]string{
		{"weight_loss", "cardiovascular_health"},
		{"muscle_gain", "strength_building"},
		{"endurance", "general_fitness"},
		{"flexibility", "stress_relief"},
		{"conditioning", "athletic_performance"},
	}
	
	for i := 0; i < dg.UserCount; i++ {
		user := UserProfile{
			ID:            i + 1,
			Name:          fmt.Sprintf("%s_%d", names[rand.Intn(len(names))], i+1),
			Age:           20 + rand.Intn(50), // 20-70 years
			Gender:        []string{"male", "female"}[rand.Intn(2)],
			Weight:        50 + rand.Float64()*50,  // 50-100 kg
			Height:        150 + rand.Float64()*30, // 150-180 cm
			ActivityLevel: []string{"sedentary", "light", "moderate", "active", "very_active"}[rand.Intn(5)],
			Goals:         goals[rand.Intn(len(goals))],
			DietType:      []string{"balanced", "halal", "vegetarian", "keto", "mediterranean"}[rand.Intn(5)],
			Language:      []string{"en", "ar"}[rand.Intn(2)],
			Country:       countries[rand.Intn(len(countries))],
			CreatedAt:     time.Now().AddDate(0, 0, -rand.Intn(365)),
			LastActive:    time.Now().AddDate(0, 0, -rand.Intn(7)),
		}
		
		// Calculate BMR using Mifflin-St Jeor equation
		if user.Gender == "male" {
			user.BMR = 10*user.Weight + 6.25*user.Height - 5*float64(user.Age) + 5
		} else {
			user.BMR = 10*user.Weight + 6.25*user.Height - 5*float64(user.Age) - 161
		}
		
		// Calculate TDEE based on activity level
		multipliers := map[string]float64{
			"sedentary":    1.2,
			"light":        1.375,
			"moderate":     1.55,
			"active":       1.725,
			"very_active":  1.9,
		}
		user.TDEE = user.BMR * multipliers[user.ActivityLevel]
		
		// Set macro targets based on goals
		user.MacroTargets = dg.calculateMacroTargets(user.TDEE, user.Goals)
		
		// Add health conditions and allergies realistically
		if rand.Float32() < 0.3 { // 30% have health conditions
			user.HealthConditions = []string{
				[]string{"diabetes", "hypertension", "arthritis", "heart_disease", "asthma"}[rand.Intn(5)],
			}
		}
		
		if rand.Float32() < 0.2 { // 20% have allergies
			user.Allergies = []string{
				[]string{"nuts", "dairy", "gluten", "eggs", "shellfish"}[rand.Intn(5)],
			}
		}
		
		users[i] = user
	}
	
	// Save users to file
	if err := dg.saveToFile("users.json", users); err != nil {
		return nil, err
	}
	
	fmt.Printf("✅ Generated %d user profiles\n", len(users))
	return users, nil
}

func (dg *DatasetGenerator) generateWorkouts(users []UserProfile) error {
	fmt.Printf("💪 Generating workouts for %d users...\n", len(users))
	
	workoutTypes := []string{"hiit", "cardio", "strength", "yoga", "pilates", "crossfit", "dance", "martial_arts"}
	goals := []string{"weight_loss", "muscle_gain", "endurance", "flexibility", "strength", "conditioning"}
	levels := []string{"beginner", "intermediate", "advanced"}
	equipment := [][]string{
		{"none"},
		{"dumbbells", "resistance_bands"},
		{"barbell", "bench", "squat_rack"},
		{"treadmill", "elliptical"},
		{"yoga_mat", "blocks"},
	}
	
	exercises := []map[string]string{
		{"en": "Push-ups", "ar": "تمرين الضغط"},
		{"en": "Squats", "ar": "القرفصاء"},
		{"en": "Lunges", "ar": "الطعن"},
		{"en": "Plank", "ar": "البلانك"},
		{"en": "Burpees", "ar": "البيربي"},
		{"en": "Mountain Climbers", "ar": "متسلق الجبال"},
		{"en": "Jumping Jacks", "ar": "القفز المتباعد"},
	}
	
	totalWorkouts := len(users) * dg.WorkoutsPerUser
	workouts := make([]LargeWorkout, totalWorkouts)
	
	for i := 0; i < totalWorkouts; i++ {
		workoutType := workoutTypes[rand.Intn(len(workoutTypes))]
		goal := goals[rand.Intn(len(goals))]
		level := levels[rand.Intn(len(levels))]
		equip := equipment[rand.Intn(len(equipment))]
		
		workout := LargeWorkout{
			ID: fmt.Sprintf("workout_%d", i+1),
			Name: map[string]string{
				"en": fmt.Sprintf("%s %s Workout", strings.Title(level), strings.Title(workoutType)),
				"ar": fmt.Sprintf("تمرين %s %s", workoutType, level),
			},
			Type:          workoutType,
			Goal:          goal,
			Level:         level,
			Duration:      20 + rand.Intn(60), // 20-80 minutes
			Equipment:     equip,
			Calories:      200 + rand.Intn(400), // 200-600 calories/hour
			TargetMuscles: dg.getTargetMuscles(workoutType),
			Health: HealthInfo{
				SuitableFor:   dg.getSuitableConditions(workoutType, level),
				AvoidIf:       dg.getAvoidConditions(workoutType),
				Modifications: dg.getModifications(level),
			},
			Created: time.Now().AddDate(0, 0, -rand.Intn(365)),
			UserIDs: dg.getRandomUserIDs(users, 50+rand.Intn(100)), // 50-150 users per workout
		}
		
		// Generate exercises for this workout
		numExercises := 5 + rand.Intn(10) // 5-15 exercises
		workout.Exercises = make([]Exercise, numExercises)
		
		for j := 0; j < numExercises; j++ {
			exerciseNames := exercises[rand.Intn(len(exercises))]
			workout.Exercises[j] = Exercise{
				Name:         exerciseNames,
				Sets:         2 + rand.Intn(4), // 2-6 sets
				Reps:         fmt.Sprintf("%d-%d", 8+rand.Intn(5), 15+rand.Intn(10)),
				Duration:     30 + rand.Intn(60), // 30-90 seconds
				Rest:         15 + rand.Intn(45), // 15-60 seconds rest
				Difficulty:   1 + rand.Intn(5),   // 1-5 difficulty
				Instructions: map[string]string{
					"en": fmt.Sprintf("Perform %s with proper form", exerciseNames["en"]),
					"ar": fmt.Sprintf("أدِ %s بالشكل الصحيح", exerciseNames["ar"]),
				},
			}
		}
		
		workouts[i] = workout
		
		if (i+1)%10000 == 0 {
			fmt.Printf("Generated %d/%d workouts...\n", i+1, totalWorkouts)
		}
	}
	
	// Save workouts in chunks for better performance
	chunkSize := 10000
	for i := 0; i < len(workouts); i += chunkSize {
		end := i + chunkSize
		if end > len(workouts) {
			end = len(workouts)
		}
		
		filename := fmt.Sprintf("workouts_chunk_%d.json", i/chunkSize+1)
		if err := dg.saveToFile(filename, workouts[i:end]); err != nil {
			return err
		}
	}
	
	fmt.Printf("✅ Generated %d workouts in chunks\n", len(workouts))
	return nil
}

// Helper functions
func (dg *DatasetGenerator) calculateMacroTargets(tdee float64, goals []string) map[string]float64 {
	// Default balanced macros
	proteinPercent := 0.25
	carbsPercent := 0.45
	fatPercent := 0.30
	
	// Adjust based on goals
	for _, goal := range goals {
		switch goal {
		case "weight_loss":
			proteinPercent = 0.35
			carbsPercent = 0.35
			fatPercent = 0.30
		case "muscle_gain":
			proteinPercent = 0.30
			carbsPercent = 0.50
			fatPercent = 0.20
		}
	}
	
	return map[string]float64{
		"calories": tdee,
		"protein":  (tdee * proteinPercent) / 4, // 4 cal/g
		"carbs":    (tdee * carbsPercent) / 4,   // 4 cal/g
		"fat":      (tdee * fatPercent) / 9,     // 9 cal/g
	}
}

func (dg *DatasetGenerator) getTargetMuscles(workoutType string) []string {
	muscleMap := map[string][]string{
		"hiit":          {"full_body", "cardiovascular"},
		"cardio":        {"cardiovascular", "legs"},
		"strength":      {"chest", "back", "shoulders", "arms", "legs"},
		"yoga":          {"full_body", "core", "flexibility"},
		"pilates":       {"core", "flexibility", "stability"},
		"crossfit":      {"full_body", "functional"},
		"dance":         {"cardiovascular", "legs", "coordination"},
		"martial_arts":  {"full_body", "core", "flexibility", "power"},
	}
	return muscleMap[workoutType]
}

func (dg *DatasetGenerator) getSuitableConditions(workoutType, level string) []string {
	if level == "beginner" {
		return []string{"general_health", "beginner_friendly"}
	}
	
	suitableMap := map[string][]string{
		"yoga":     {"arthritis", "stress", "flexibility_issues"},
		"cardio":   {"weight_management", "cardiovascular_health"},
		"strength": {"bone_density", "muscle_building"},
		"hiit":     {"time_constrained", "advanced_fitness"},
	}
	return suitableMap[workoutType]
}

func (dg *DatasetGenerator) getAvoidConditions(workoutType string) []string {
	avoidMap := map[string][]string{
		"hiit":     {"heart_disease", "recent_injury"},
		"strength": {"acute_injury", "severe_arthritis"},
		"cardio":   {"severe_asthma", "heart_condition"},
	}
	return avoidMap[workoutType]
}

func (dg *DatasetGenerator) getModifications(level string) []string {
	modMap := map[string][]string{
		"beginner":     {"reduce_intensity", "longer_rest", "assisted_movements"},
		"intermediate": {"moderate_weights", "standard_rest"},
		"advanced":     {"increase_weight", "reduce_rest", "add_complexity"},
	}
	return modMap[level]
}

func (dg *DatasetGenerator) getRandomUserIDs(users []UserProfile, count int) []int {
	if count > len(users) {
		count = len(users)
	}
	
	selected := make(map[int]bool)
	userIDs := make([]int, 0, count)
	
	for len(userIDs) < count {
		userID := rand.Intn(len(users)) + 1
		if !selected[userID] {
			selected[userID] = true
			userIDs = append(userIDs, userID)
		}
	}
	
	return userIDs
}

func (dg *DatasetGenerator) saveToFile(filename string, data interface{}) error {
	filePath := filepath.Join(dg.OutputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Placeholder methods for recipes, meal plans, and health data
func (dg *DatasetGenerator) generateRecipes(users []UserProfile) error {
	fmt.Printf("🍽️ Generating recipes for %d users...\n", len(users))
	// Implementation similar to workouts...
	return nil
}

func (dg *DatasetGenerator) generateMealPlans(users []UserProfile) error {
	fmt.Printf("📋 Generating meal plans for %d users...\n", len(users))
	// Implementation for meal plans...
	return nil
}

func (dg *DatasetGenerator) generateHealthData(users []UserProfile) error {
	fmt.Printf("🏥 Generating health data for %d users...\n", len(users))
	// Implementation for health conditions, vitamins, etc...
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	generator := NewDatasetGenerator("../data/large_dataset")
	if err := generator.GenerateAllData(); err != nil {
		fmt.Printf("❌ Error generating dataset: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("🎉 Large dataset generation completed successfully!")
}