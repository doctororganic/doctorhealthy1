package models

import (
	"time"
)

// FileMetadata represents file metadata information
type FileMetadata struct {
	ID            uint      `json:"id" db:"id"`
	UserID        uint      `json:"user_id" db:"user_id"`
	Filename      string    `json:"filename" db:"filename"`
	StoragePath   string    `json:"storage_path" db:"storage_path"`
	FileSize      int64     `json:"file_size" db:"file_size"`
	ContentType   string    `json:"content_type" db:"content_type"`
	FileHash      string    `json:"file_hash" db:"file_hash"`
	IsPublic      bool      `json:"is_public" db:"is_public"`
	DownloadCount int       `json:"download_count" db:"download_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Product represents a product in the system
type Product struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Category    string    `json:"category" db:"category"`
	Brand       string    `json:"brand" db:"brand"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	IsApproved  bool      `json:"is_approved" db:"is_approved"`
	CreatedBy   uint      `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Workout represents a workout session
type Workout struct {
	ID             uint       `json:"id" db:"id"`
	UserID         uint       `json:"user_id" db:"user_id"`
	Name           string     `json:"name" db:"name"`
	Description    string     `json:"description" db:"description"`
	Duration       int        `json:"duration" db:"duration"`        // in minutes
	Intensity      string     `json:"intensity" db:"intensity"`      // low, medium, high
	CaloriesBurned int        `json:"calories_burned" db:"calories_burned"`
	WorkoutType    string     `json:"workout_type" db:"workout_type"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// WorkoutPlan represents a workout plan
type WorkoutPlan struct {
	ID              uint       `json:"id" db:"id"`
	UserID          uint       `json:"user_id" db:"user_id"`
	Name            string     `json:"name" db:"name"`
	Description     string     `json:"description" db:"description"`
	DurationWeeks   int        `json:"duration_weeks" db:"duration_weeks"`
	DifficultyLevel string     `json:"difficulty_level" db:"difficulty_level"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	StartDate       *time.Time `json:"start_date" db:"start_date"`
	EndDate         *time.Time `json:"end_date" db:"end_date"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// Recipe represents a recipe
type Recipe struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Instructions string   `json:"instructions" db:"instructions"`
	PrepTime    int       `json:"prep_time" db:"prep_time"`    // in minutes
	CookTime    int       `json:"cook_time" db:"cook_time"`    // in minutes
	Servings    int       `json:"servings" db:"servings"`
	Calories    int       `json:"calories" db:"calories"`
	Protein     float64   `json:"protein" db:"protein"`
	Carbs       float64   `json:"carbs" db:"carbs"`
	Fat         float64   `json:"fat" db:"fat"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedBy   uint      `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// HealthProfile represents a user's health profile
type HealthProfile struct {
	ID               uint        `json:"id" db:"id"`
	UserID           uint        `json:"user_id" db:"user_id"`
	DateOfBirth      *time.Time  `json:"date_of_birth" db:"date_of_birth"`
	Gender           *string     `json:"gender" db:"gender"`
	Height           *float64    `json:"height" db:"height"`           // in cm
	Weight           *float64    `json:"weight" db:"weight"`           // in kg
	BloodType        string      `json:"blood_type" db:"blood_type"`
	Allergies        []string    `json:"allergies" db:"allergies"`
	MedicalConditions []string   `json:"medical_conditions" db:"medical_conditions"`
	Medications      []string    `json:"medications" db:"medications"`
	EmergencyContact string      `json:"emergency_contact" db:"emergency_contact"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// HealthCondition represents a health condition
type HealthCondition struct {
	ID            uint       `json:"id" db:"id"`
	UserID        uint       `json:"user_id" db:"user_id"`
	Name          string     `json:"name" db:"name"`
	Description   string     `json:"description" db:"description"`
	Severity      string     `json:"severity" db:"severity"`      // mild, moderate, severe
	DiagnosedDate *time.Time `json:"diagnosed_date" db:"diagnosed_date"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ProgressEntry represents a progress tracking entry
type ProgressEntry struct {
	ID                uint        `json:"id" db:"id"`
	UserID            uint        `json:"user_id" db:"user_id"`
	Weight            *float64    `json:"weight" db:"weight"`                     // in kg
	BodyFatPercentage *float64    `json:"body_fat_percentage" db:"body_fat_percentage"`
	Measurements      []string    `json:"measurements" db:"measurements"`        // JSON array of measurements
	Notes             string      `json:"notes" db:"notes"`
	RecordedAt        time.Time   `json:"recorded_at" db:"recorded_at"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

// Meal represents a meal entry
type Meal struct {
	ID          uint      `json:"id" db:"id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Calories    int       `json:"calories" db:"calories"`
	Protein     float64   `json:"protein" db:"protein"`
	Carbs       float64   `json:"carbs" db:"carbs"`
	Fat         float64   `json:"fat" db:"fat"`
	Fiber       float64   `json:"fiber" db:"fiber"`
	Sugar       float64   `json:"sugar" db:"sugar"`
	Sodium      int       `json:"sodium" db:"sodium"`
	MealType    string    `json:"meal_type" db:"meal_type"`    // breakfast, lunch, dinner, snack
	ConsumedAt  time.Time `json:"consumed_at" db:"consumed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NutritionPlan represents a nutrition plan
type NutritionPlan struct {
	ID            uint        `json:"id" db:"id"`
	UserID        uint        `json:"user_id" db:"user_id"`
	Name          string      `json:"name" db:"name"`
	Description   string      `json:"description" db:"description"`
	DailyCalories *int        `json:"daily_calories" db:"daily_calories"`
	ProteinGrams  *float64    `json:"protein_grams" db:"protein_grams"`
	CarbsGrams    *float64    `json:"carbs_grams" db:"carbs_grams"`
	FatGrams      *float64    `json:"fat_grams" db:"fat_grams"`
	FiberGrams    *float64    `json:"fiber_grams" db:"fiber_grams"`
	SugarGrams    *float64    `json:"sugar_grams" db:"sugar_grams"`
	SodiumMg      *int        `json:"sodium_mg" db:"sodium_mg"`
	WaterMl       *int        `json:"water_ml" db:"water_ml"`
	IsActive      bool        `json:"is_active" db:"is_active"`
	StartDate     *time.Time  `json:"start_date" db:"start_date"`
	EndDate       *time.Time  `json:"end_date" db:"end_date"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}
