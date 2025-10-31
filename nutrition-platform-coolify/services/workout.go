package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Exercise represents a single exercise in a workout
type Exercise struct {
	Name        string  `json:"name"`
	Sets        int     `json:"sets"`
	Reps        int     `json:"reps"`
	Weight      float64 `json:"weight,omitempty"`   // in kg
	Duration    int     `json:"duration,omitempty"` // in seconds
	Distance    float64 `json:"distance,omitempty"` // in km
	Calories    int     `json:"calories,omitempty"`
	RestTime    int     `json:"rest_time,omitempty"` // in seconds
	Notes       string  `json:"notes,omitempty"`
	MuscleGroup string  `json:"muscle_group"`
	Equipment   string  `json:"equipment,omitempty"`
}

// Workout represents a workout session
type Workout struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Exercises    []Exercise `json:"exercises"`
	Duration     int        `json:"duration"`      // total duration in minutes
	Calories     int        `json:"calories"`      // total calories burned
	Difficulty   string     `json:"difficulty"`    // beginner, intermediate, advanced
	Type         string     `json:"type"`          // strength, cardio, flexibility, sports
	TargetMuscle []string   `json:"target_muscle"` // chest, back, legs, etc.
	Equipment    []string   `json:"equipment"`     // dumbbells, barbell, bodyweight, etc.
	Tags         []string   `json:"tags"`
	IsPublic     bool       `json:"is_public"`
	Rating       float64    `json:"rating"`
	RatingCount  int        `json:"rating_count"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// WorkoutData represents the structure of workouts.json
type WorkoutData struct {
	Workouts []Workout `json:"workouts"`
	Metadata Metadata  `json:"metadata"`
}

const workoutsFile = "backend/data/workouts.json"

// CreateWorkout creates a new workout
func CreateWorkout(workout *Workout) error {
	// Generate ID and timestamps
	workout.ID = uuid.New().String()
	workout.CreatedAt = time.Now()
	workout.UpdatedAt = time.Now()

	// Validate difficulty
	validDifficulties := map[string]bool{
		"beginner":     true,
		"intermediate": true,
		"advanced":     true,
	}
	if !validDifficulties[workout.Difficulty] {
		workout.Difficulty = "beginner" // default
	}

	// Validate type
	validTypes := map[string]bool{
		"strength":    true,
		"cardio":      true,
		"flexibility": true,
		"sports":      true,
		"mixed":       true,
	}
	if !validTypes[workout.Type] {
		workout.Type = "mixed" // default
	}

	// Calculate total duration and calories if not provided
	if workout.Duration == 0 {
		workout.Duration = calculateWorkoutDuration(workout.Exercises)
	}
	if workout.Calories == 0 {
		workout.Calories = calculateWorkoutCalories(workout.Exercises, workout.Duration)
	}

	// Extract target muscles and equipment from exercises
	if len(workout.TargetMuscle) == 0 {
		workout.TargetMuscle = extractTargetMuscles(workout.Exercises)
	}
	if len(workout.Equipment) == 0 {
		workout.Equipment = extractEquipment(workout.Exercises)
	}

	// Initialize rating
	workout.Rating = 0
	workout.RatingCount = 0

	return AppendJSON(workoutsFile, workout)
}

// GetWorkoutsByUserID retrieves all workouts for a specific user
func GetWorkoutsByUserID(userID string) ([]Workout, error) {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return nil, err
	}

	var userWorkouts []Workout
	for _, workout := range data.Workouts {
		if workout.UserID == userID {
			userWorkouts = append(userWorkouts, workout)
		}
	}

	return userWorkouts, nil
}

// GetPublicWorkouts retrieves all public workouts
func GetPublicWorkouts() ([]Workout, error) {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return nil, err
	}

	var publicWorkouts []Workout
	for _, workout := range data.Workouts {
		if workout.IsPublic {
			publicWorkouts = append(publicWorkouts, workout)
		}
	}

	return publicWorkouts, nil
}

// GetWorkoutByID retrieves a specific workout by ID
func GetWorkoutByID(workoutID string) (*Workout, error) {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return nil, err
	}

	for _, workout := range data.Workouts {
		if workout.ID == workoutID {
			return &workout, nil
		}
	}

	return nil, fmt.Errorf("workout not found")
}

// UpdateWorkout updates an existing workout
func UpdateWorkout(workoutID string, updatedWorkout *Workout) error {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return err
	}

	for i, workout := range data.Workouts {
		if workout.ID == workoutID {
			// Preserve original ID, created time, and rating data
			updatedWorkout.ID = workout.ID
			updatedWorkout.CreatedAt = workout.CreatedAt
			updatedWorkout.UpdatedAt = time.Now()
			updatedWorkout.Rating = workout.Rating
			updatedWorkout.RatingCount = workout.RatingCount

			// Recalculate duration and calories
			if updatedWorkout.Duration == 0 {
				updatedWorkout.Duration = calculateWorkoutDuration(updatedWorkout.Exercises)
			}
			if updatedWorkout.Calories == 0 {
				updatedWorkout.Calories = calculateWorkoutCalories(updatedWorkout.Exercises, updatedWorkout.Duration)
			}

			// Update target muscles and equipment
			if len(updatedWorkout.TargetMuscle) == 0 {
				updatedWorkout.TargetMuscle = extractTargetMuscles(updatedWorkout.Exercises)
			}
			if len(updatedWorkout.Equipment) == 0 {
				updatedWorkout.Equipment = extractEquipment(updatedWorkout.Exercises)
			}

			data.Workouts[i] = *updatedWorkout
			data.Metadata.UpdatedAt = time.Now()

			return WriteJSON(workoutsFile, data)
		}
	}

	return fmt.Errorf("workout not found")
}

// DeleteWorkout deletes a workout
func DeleteWorkout(workoutID string, userID string) error {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return err
	}

	for i, workout := range data.Workouts {
		if workout.ID == workoutID {
			// Check if user owns this workout
			if workout.UserID != userID {
				return fmt.Errorf("unauthorized: workout belongs to another user")
			}

			// Remove workout from slice
			data.Workouts = append(data.Workouts[:i], data.Workouts[i+1:]...)
			data.Metadata.UpdatedAt = time.Now()

			return WriteJSON(workoutsFile, data)
		}
	}

	return fmt.Errorf("workout not found")
}

// CompleteWorkout marks a workout as completed
func CompleteWorkout(workoutID string, userID string) error {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return err
	}

	for i, workout := range data.Workouts {
		if workout.ID == workoutID && workout.UserID == userID {
			now := time.Now()
			data.Workouts[i].CompletedAt = &now
			data.Workouts[i].UpdatedAt = now
			data.Metadata.UpdatedAt = now

			return WriteJSON(workoutsFile, data)
		}
	}

	return fmt.Errorf("workout not found or unauthorized")
}

// SearchWorkouts searches workouts by various criteria
func SearchWorkouts(query string, filters map[string]interface{}) ([]Workout, error) {
	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return nil, err
	}

	var results []Workout
	queryLower := strings.ToLower(query)

	for _, workout := range data.Workouts {
		// Only search public workouts or user's own workouts
		userID, hasUserID := filters["user_id"].(string)
		if !workout.IsPublic && (!hasUserID || workout.UserID != userID) {
			continue
		}

		// Text search
		if query != "" {
			matchesQuery := false

			// Search in name
			if strings.Contains(strings.ToLower(workout.Name), queryLower) {
				matchesQuery = true
			}

			// Search in description
			if !matchesQuery && strings.Contains(strings.ToLower(workout.Description), queryLower) {
				matchesQuery = true
			}

			// Search in exercise names
			if !matchesQuery {
				for _, exercise := range workout.Exercises {
					if strings.Contains(strings.ToLower(exercise.Name), queryLower) {
						matchesQuery = true
						break
					}
				}
			}

			// Search in tags
			if !matchesQuery {
				for _, tag := range workout.Tags {
					if strings.Contains(strings.ToLower(tag), queryLower) {
						matchesQuery = true
						break
					}
				}
			}

			if !matchesQuery {
				continue
			}
		}

		// Apply filters
		if !matchesWorkoutFilters(workout, filters) {
			continue
		}

		results = append(results, workout)
	}

	return results, nil
}

// matchesWorkoutFilters checks if a workout matches the given filters
func matchesWorkoutFilters(workout Workout, filters map[string]interface{}) bool {
	if workoutType, ok := filters["type"].(string); ok && workoutType != "" {
		if workout.Type != workoutType {
			return false
		}
	}

	if difficulty, ok := filters["difficulty"].(string); ok && difficulty != "" {
		if workout.Difficulty != difficulty {
			return false
		}
	}

	if targetMuscle, ok := filters["target_muscle"].(string); ok && targetMuscle != "" {
		found := false
		for _, muscle := range workout.TargetMuscle {
			if muscle == targetMuscle {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if equipment, ok := filters["equipment"].(string); ok && equipment != "" {
		found := false
		for _, eq := range workout.Equipment {
			if eq == equipment {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if maxDuration, ok := filters["max_duration"].(float64); ok {
		if float64(workout.Duration) > maxDuration {
			return false
		}
	}

	if minRating, ok := filters["min_rating"].(float64); ok {
		if workout.Rating < minRating {
			return false
		}
	}

	return true
}

// RateWorkout adds or updates a rating for a workout
func RateWorkout(workoutID string, rating float64) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	var data WorkoutData
	err := ReadJSON(workoutsFile, &data)
	if err != nil {
		return err
	}

	for i, workout := range data.Workouts {
		if workout.ID == workoutID {
			// Calculate new average rating
			totalRating := workout.Rating * float64(workout.RatingCount)
			totalRating += rating
			data.Workouts[i].RatingCount++
			data.Workouts[i].Rating = totalRating / float64(data.Workouts[i].RatingCount)
			data.Workouts[i].UpdatedAt = time.Now()
			data.Metadata.UpdatedAt = time.Now()

			return WriteJSON(workoutsFile, data)
		}
	}

	return fmt.Errorf("workout not found")
}

// GetWorkoutStats returns statistics about user's workouts
func GetWorkoutStats(userID string) (map[string]interface{}, error) {
	workouts, err := GetWorkoutsByUserID(userID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_workouts":     len(workouts),
		"completed_workouts": 0,
		"total_duration":     0,
		"total_calories":     0,
		"workout_types":      map[string]int{},
		"difficulty_levels":  map[string]int{},
		"favorite_muscles":   map[string]int{},
	}

	completedCount := 0
	totalDuration := 0
	totalCalories := 0
	workoutTypes := map[string]int{}
	difficultyLevels := map[string]int{}
	muscleCount := map[string]int{}

	for _, workout := range workouts {
		if workout.CompletedAt != nil {
			completedCount++
			totalDuration += workout.Duration
			totalCalories += workout.Calories
		}
		workoutTypes[workout.Type]++
		difficultyLevels[workout.Difficulty]++

		// Count target muscles
		for _, muscle := range workout.TargetMuscle {
			muscleCount[muscle]++
		}
	}

	stats["completed_workouts"] = completedCount
	stats["total_duration"] = totalDuration
	stats["total_calories"] = totalCalories
	stats["workout_types"] = workoutTypes
	stats["difficulty_levels"] = difficultyLevels
	stats["favorite_muscles"] = muscleCount

	return stats, nil
}

// Helper functions
func calculateWorkoutDuration(exercises []Exercise) int {
	totalDuration := 0
	for _, exercise := range exercises {
		if exercise.Duration > 0 {
			totalDuration += exercise.Duration
		} else {
			// Estimate duration based on sets and reps (assuming 2 seconds per rep + rest time)
			estimatedTime := exercise.Sets * (exercise.Reps*2 + exercise.RestTime)
			totalDuration += estimatedTime
		}
	}
	return totalDuration / 60 // convert to minutes
}

func calculateWorkoutCalories(exercises []Exercise, duration int) int {
	totalCalories := 0
	for _, exercise := range exercises {
		if exercise.Calories > 0 {
			totalCalories += exercise.Calories
		}
	}

	// If no calories specified, estimate based on duration and intensity
	if totalCalories == 0 {
		// Rough estimation: 5-10 calories per minute depending on intensity
		totalCalories = duration * 7 // average estimation
	}

	return totalCalories
}

func extractTargetMuscles(exercises []Exercise) []string {
	muscleMap := make(map[string]bool)
	for _, exercise := range exercises {
		if exercise.MuscleGroup != "" {
			muscleMap[exercise.MuscleGroup] = true
		}
	}

	var muscles []string
	for muscle := range muscleMap {
		muscles = append(muscles, muscle)
	}
	return muscles
}

func extractEquipment(exercises []Exercise) []string {
	equipmentMap := make(map[string]bool)
	for _, exercise := range exercises {
		if exercise.Equipment != "" {
			equipmentMap[exercise.Equipment] = true
		} else {
			equipmentMap["bodyweight"] = true
		}
	}

	var equipment []string
	for eq := range equipmentMap {
		equipment = append(equipment, eq)
	}
	return equipment
}
