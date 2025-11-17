package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Response structure for API responses
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Simple data structures
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Updated  time.Time `json:"updated,omitempty"`
}

type Food struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Calories    int       `json:"calories"`
	Protein     int       `json:"protein"`
	Carbs       int       `json:"carbs"`
	Fat         int       `json:"fat"`
	Updated     time.Time `json:"updated,omitempty"`
}

type Workout struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	Difficulty  string    `json:"difficulty"`
	Description string    `json:"description"`
	Updated     time.Time `json:"updated,omitempty"`
}

type Recipe struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CookingTime int       `json:"cooking_time"`
	Updated     time.Time `json:"updated,omitempty"`
}

type Exercise struct {
	Name        string `json:"name"`
	Sets        int    `json:"sets"`
	Reps        int    `json:"reps"`
	Duration    int    `json:"duration"` // in seconds, for cardio exercises
	Rest        int    `json:"rest"`     // rest time in seconds
	Description string `json:"description"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

var foods = []Food{
	{ID: 1, Name: "Apple", Calories: 95, Protein: 0, Carbs: 25, Fat: 0},
	{ID: 2, Name: "Chicken Breast", Calories: 165, Protein: 31, Carbs: 0, Fat: 4},
	{ID: 3, Name: "Rice", Calories: 130, Protein: 3, Carbs: 28, Fat: 0},
}

var workouts = []Workout{
	{ID: 1, Name: "Push-ups", Duration: 15, Difficulty: "Beginner", Description: "Classic upper body exercise"},
	{ID: 2, Name: "Squats", Duration: 20, Difficulty: "Beginner", Description: "Lower body strength exercise"},
	{ID: 3, Name: "Pull-ups", Duration: 15, Difficulty: "Intermediate", Description: "Upper body pulling exercise"},
}

var recipes = []Recipe{
	{ID: 1, Name: "Grilled Chicken Salad", Description: "Healthy salad with grilled chicken", CookingTime: 25},
	{ID: 2, Name: "Vegetable Stir Fry", Description: "Mixed vegetables with rice", CookingTime: 20},
	{ID: 3, Name: "Protein Smoothie", Description: "Banana and protein powder smoothie", CookingTime: 5},
}

func main() {
	r := mux.NewRouter()
	
	// Enable CORS
	r.Use(corsMiddleware)
	
	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status:  "healthy",
			Message: "Server is running",
			Data: map[string]interface{}{
				"timestamp": time.Now(),
				"version":   "1.0.0",
			},
		})
	}).Methods("GET")
	
	// API info endpoint
	r.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status: "success",
			Data: map[string]interface{}{
				"service":     "Nutrition Platform API",
				"version":     "1.0.0",
				"endpoints":   []string{"/health", "/api/v1/users", "/api/v1/foods", "/api/v1/workouts", "/api/v1/recipes"},
				"methods":     []string{"GET", "POST", "PUT", "DELETE"},
			},
		})
	}).Methods("GET")
	
	// User endpoints
	r.HandleFunc("/api/v1/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", getUserHandler).Methods("GET")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/users", createUserHandler).Methods("POST")
	
	// Food endpoints
	r.HandleFunc("/api/v1/foods", getFoodsHandler).Methods("GET")
	r.HandleFunc("/api/v1/foods/{id:[0-9]+}", getFoodHandler).Methods("GET")
	r.HandleFunc("/api/v1/foods/{id:[0-9]+}", updateFoodHandler).Methods("PUT")
	r.HandleFunc("/api/v1/foods/{id:[0-9]+}", deleteFoodHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/foods", createFoodHandler).Methods("POST")
	
	// Workout endpoints
	r.HandleFunc("/api/v1/workouts", getWorkoutsHandler).Methods("GET")
	r.HandleFunc("/api/v1/workouts/{id:[0-9]+}", getWorkoutHandler).Methods("GET")
	r.HandleFunc("/api/v1/workouts/{id:[0-9]+}", updateWorkoutHandler).Methods("PUT")
	r.HandleFunc("/api/v1/workouts/{id:[0-9]+}", deleteWorkoutHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/workouts", createWorkoutHandler).Methods("POST")
	r.HandleFunc("/api/v1/workouts/generate", generateWorkoutHandler).Methods("POST")
	
	// Recipe endpoints
	r.HandleFunc("/api/v1/recipes", getRecipesHandler).Methods("GET")
	r.HandleFunc("/api/v1/recipes/{id:[0-9]+}", getRecipeHandler).Methods("GET")
	r.HandleFunc("/api/v1/recipes/{id:[0-9]+}", updateRecipeHandler).Methods("PUT")
	r.HandleFunc("/api/v1/recipes/{id:[0-9]+}", deleteRecipeHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/recipes", createRecipeHandler).Methods("POST")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("🚀 Nutrition Platform API Server starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API info: http://localhost:%s/api/info", port)
	log.Printf("👥 Users: http://localhost:%s/api/v1/users", port)
	log.Printf("🍎 Foods: http://localhost:%s/api/v1/foods", port)
	log.Printf("💪 Workouts: http://localhost:%s/api/v1/workouts", port)
	log.Printf("🍳 Recipes: http://localhost:%s/api/v1/recipes", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Handler functions
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data:   users,
	})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(Response{
				Status: "success",
				Data:   user,
			})
			return
		}
	}
	
	http.Error(w, "User not found", http.StatusNotFound)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	
	var updateUser User
	if err := json.Unmarshal(body, &updateUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Find and update user
	for i, user := range users {
		if user.ID == id {
			updateUser.ID = id
			updateUser.Updated = time.Now()
			users[i] = updateUser
			
			log.Printf("✅ Updated user %d: %s", id, updateUser.Name)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "User updated successfully",
				Data:    updateUser,
			})
			return
		}
	}
	
	http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			log.Printf("🗑️ Deleted user %d", id)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "User deleted successfully",
			})
			return
		}
	}
	
	http.Error(w, "User not found", http.StatusNotFound)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	
	log.Printf("➕ Created user %d: %s", newUser.ID, newUser.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    newUser,
	})
}

// Food handlers
func getFoodsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data:   foods,
	})
}

func getFoodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for _, food := range foods {
		if food.ID == id {
			json.NewEncoder(w).Encode(Response{
				Status: "success",
				Data:   food,
			})
			return
		}
	}
	
	http.Error(w, "Food not found", http.StatusNotFound)
}

func updateFoodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	var updateFood Food
	if err := json.NewDecoder(r.Body).Decode(&updateFood); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	for i, food := range foods {
		if food.ID == id {
			updateFood.ID = id
			updateFood.Updated = time.Now()
			foods[i] = updateFood
			
			log.Printf("✅ Updated food %d: %s", id, updateFood.Name)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Food updated successfully",
				Data:    updateFood,
			})
			return
		}
	}
	
	http.Error(w, "Food not found", http.StatusNotFound)
}

func deleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for i, food := range foods {
		if food.ID == id {
			foods = append(foods[:i], foods[i+1:]...)
			log.Printf("🗑️ Deleted food %d", id)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Food deleted successfully",
			})
			return
		}
	}
	
	http.Error(w, "Food not found", http.StatusNotFound)
}

func createFoodHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var newFood Food
	if err := json.NewDecoder(r.Body).Decode(&newFood); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	newFood.ID = len(foods) + 1
	foods = append(foods, newFood)
	
	log.Printf("➕ Created food %d: %s", newFood.ID, newFood.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Food created successfully",
		Data:    newFood,
	})
}

// Workout handlers
func getWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data:   workouts,
	})
}

func getWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for _, workout := range workouts {
		if workout.ID == id {
			json.NewEncoder(w).Encode(Response{
				Status: "success",
				Data:   workout,
			})
			return
		}
	}
	
	http.Error(w, "Workout not found", http.StatusNotFound)
}

func updateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	var updateWorkout Workout
	if err := json.NewDecoder(r.Body).Decode(&updateWorkout); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	for i, workout := range workouts {
		if workout.ID == id {
			updateWorkout.ID = id
			updateWorkout.Updated = time.Now()
			workouts[i] = updateWorkout
			
			log.Printf("✅ Updated workout %d: %s", id, updateWorkout.Name)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Workout updated successfully",
				Data:    updateWorkout,
			})
			return
		}
	}
	
	http.Error(w, "Workout not found", http.StatusNotFound)
}

func deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for i, workout := range workouts {
		if workout.ID == id {
			workouts = append(workouts[:i], workouts[i+1:]...)
			log.Printf("🗑️ Deleted workout %d", id)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Workout deleted successfully",
			})
			return
		}
	}
	
	http.Error(w, "Workout not found", http.StatusNotFound)
}

func createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var newWorkout Workout
	if err := json.NewDecoder(r.Body).Decode(&newWorkout); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	newWorkout.ID = len(workouts) + 1
	workouts = append(workouts, newWorkout)
	
	log.Printf("➕ Created workout %d: %s", newWorkout.ID, newWorkout.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Workout created successfully",
		Data:    newWorkout,
	})
}

// Generate workout handler
func generateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Define request structure
	type WorkoutRequest struct {
		FitnessLevel string `json:"fitness_level"` // beginner, intermediate, advanced
		Goals        string `json:"goals"`         // strength, cardio, flexibility, weight_loss
		Duration     int    `json:"duration"`      // workout duration in minutes
	}
	
	// Define generated workout structure
	type GeneratedWorkout struct {
		Name         string            `json:"name"`
		Duration     int               `json:"duration"`
		Difficulty   string            `json:"difficulty"`
		Description  string            `json:"description"`
		Exercises    []Exercise        `json:"exercises"`
		Goals        []string          `json:"goals"`
		FitnessLevel string            `json:"fitness_level"`
	}
	
	var req WorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Duration <= 0 || req.Duration > 180 {
		http.Error(w, "Duration must be between 1 and 180 minutes", http.StatusBadRequest)
		return
	}
	
	// Generate workout based on parameters
	generatedWorkout := GeneratedWorkout{
		Name:         generateWorkoutName(req.FitnessLevel, req.Goals),
		Duration:     req.Duration,
		Difficulty:   req.FitnessLevel,
		Description:  generateWorkoutDescription(req.FitnessLevel, req.Goals, req.Duration),
		Goals:        []string{req.Goals},
		FitnessLevel: req.FitnessLevel,
		Exercises:    generateExercises(req.FitnessLevel, req.Goals, req.Duration),
	}
	
	log.Printf("🎯 Generated workout: %s (%d min, %s, %s)", 
		generatedWorkout.Name, req.Duration, req.FitnessLevel, req.Goals)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Message: "Workout generated successfully",
		Data:    generatedWorkout,
	})
}

// Helper functions for workout generation
func generateWorkoutName(fitnessLevel, goals string) string {
	level := map[string]string{
		"beginner":     "Beginner",
		"intermediate": "Intermediate",
		"advanced":     "Advanced",
	}[fitnessLevel]
	
	goal := map[string]string{
		"strength":      "Strength",
		"cardio":        "Cardio",
		"flexibility":   "Flexibility",
		"weight_loss":   "Fat Burning",
	}[goals]
	
	return fmt.Sprintf("%s %s Workout", level, goal)
}

func generateWorkoutDescription(fitnessLevel, goals string, duration int) string {
	return fmt.Sprintf("A %d-minute %s workout designed for %s level to help you achieve your %s goals.", 
		duration, goals, fitnessLevel, goals)
}

func generateExercises(fitnessLevel, goals string, duration int) []Exercise {
	var exercises []Exercise
	
	// Exercise database based on goals
	exerciseDB := map[string][]map[string]interface{}{
		"strength": {
			{"name": "Push-ups", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 8, 12, 15), "rest": 60, "description": "Classic upper body exercise"},
			{"name": "Squats", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 10, 15, 20), "rest": 60, "description": "Lower body strength exercise"},
			{"name": "Plank", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 45, "description": "Core stability exercise"},
			{"name": "Lunges", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 8, 12, 15), "rest": 60, "description": "Leg and glute exercise"},
			{"name": "Dips", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 5, 10, 15), "rest": 60, "description": "Triceps and chest exercise"},
		},
		"cardio": {
			{"name": "Jumping Jacks", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 30, "description": "Full body cardio exercise"},
			{"name": "High Knees", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 30, "description": "Cardio and leg exercise"},
			{"name": "Burpees", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 5, 8, 12), "rest": 60, "description": "Full body conditioning"},
			{"name": "Mountain Climbers", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 30, "description": "Core and cardio exercise"},
			{"name": "Jump Rope", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 60, 90, 120), "rest": 45, "description": "Cardio exercise"},
		},
		"weight_loss": {
			{"name": "Jumping Jacks", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 45, 60, 90), "rest": 30, "description": "Full body cardio exercise"},
			{"name": "Bodyweight Squats", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 15, 20, 25), "rest": 45, "description": "Lower body strength and cardio"},
			{"name": "Push-ups", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 8, 12, 15), "rest": 45, "description": "Upper body strength"},
			{"name": "Plank", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 30, "description": "Core stability"},
			{"name": "Burpees", "sets": 3, "reps": getRepsByLevel(fitnessLevel, 5, 8, 10), "rest": 60, "description": "Full body conditioning"},
		},
		"flexibility": {
			{"name": "Yoga Sun Salutation", "sets": 3, "duration": 60, "rest": 30, "description": "Full body stretching sequence"},
			{"name": "Hamstring Stretch", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 15, "description": "Lower body flexibility"},
			{"name": "Shoulder Stretch", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 20, 30, 45), "rest": 15, "description": "Upper body flexibility"},
			{"name": "Hip Flexor Stretch", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 15, "description": "Hip and lower back flexibility"},
			{"name": "Child's Pose", "sets": 3, "duration": getDurationByLevel(fitnessLevel, 30, 45, 60), "rest": 15, "description": "Relaxation and back stretch"},
		},
	}
	
	// Get exercises for the specified goal
	if goalExercises, ok := exerciseDB[goals]; ok {
		// Select appropriate number of exercises based on duration
		numExercises := calculateNumExercises(duration)
		
		for i := 0; i < numExercises && i < len(goalExercises); i++ {
			exerciseData := goalExercises[i]
			exercise := Exercise{
				Name:        exerciseData["name"].(string),
				Description: exerciseData["description"].(string),
			}
			
			if sets, ok := exerciseData["sets"].(int); ok {
				exercise.Sets = sets
			}
			if reps, ok := exerciseData["reps"].(int); ok {
				exercise.Reps = reps
			}
			if duration, ok := exerciseData["duration"].(int); ok {
				exercise.Duration = duration
			}
			if rest, ok := exerciseData["rest"].(int); ok {
				exercise.Rest = rest
			}
			
			exercises = append(exercises, exercise)
		}
	}
	
	return exercises
}

func getRepsByLevel(level string, beginner, intermediate, advanced int) int {
	switch level {
	case "beginner":
		return beginner
	case "intermediate":
		return intermediate
	case "advanced":
		return advanced
	default:
		return beginner
	}
}

func getDurationByLevel(level string, beginner, intermediate, advanced int) int {
	switch level {
	case "beginner":
		return beginner
	case "intermediate":
		return intermediate
	case "advanced":
		return advanced
	default:
		return beginner
	}
}

func calculateNumExercises(duration int) int {
	if duration <= 20 {
		return 3
	} else if duration <= 40 {
		return 4
	} else {
		return 5
	}
}

// Recipe handlers
func getRecipesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data:   recipes,
	})
}

func getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for _, recipe := range recipes {
		if recipe.ID == id {
			json.NewEncoder(w).Encode(Response{
				Status: "success",
				Data:   recipe,
			})
			return
		}
	}
	
	http.Error(w, "Recipe not found", http.StatusNotFound)
}

func updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	var updateRecipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&updateRecipe); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	for i, recipe := range recipes {
		if recipe.ID == id {
			updateRecipe.ID = id
			updateRecipe.Updated = time.Now()
			recipes[i] = updateRecipe
			
			log.Printf("✅ Updated recipe %d: %s", id, updateRecipe.Name)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Recipe updated successfully",
				Data:    updateRecipe,
			})
			return
		}
	}
	
	http.Error(w, "Recipe not found", http.StatusNotFound)
}

func deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	for i, recipe := range recipes {
		if recipe.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			log.Printf("🗑️ Deleted recipe %d", id)
			json.NewEncoder(w).Encode(Response{
				Status:  "success",
				Message: "Recipe deleted successfully",
			})
			return
		}
	}
	
	http.Error(w, "Recipe not found", http.StatusNotFound)
}

func createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var newRecipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&newRecipe); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	newRecipe.ID = len(recipes) + 1
	recipes = append(recipes, newRecipe)
	
	log.Printf("➕ Created recipe %d: %s", newRecipe.ID, newRecipe.Name)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Recipe created successfully",
		Data:    newRecipe,
	})
}
