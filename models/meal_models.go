package models

import (
	"time"
)

// Food represents a food item in the system
type Food struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Brand       *string   `json:"brand" db:"brand"`
	Barcode     *string   `json:"barcode" db:"barcode"`
	Category    *string   `json:"category" db:"category"`
	Calories    float64   `json:"calories" db:"calories"`
	Protein     float64   `json:"protein" db:"protein"`
	Carbs       float64   `json:"carbs" db:"carbs"`
	Fat         float64   `json:"fat" db:"fat"`
	Fiber       float64   `json:"fiber" db:"fiber"`
	Sugar       float64   `json:"sugar" db:"sugar"`
	Sodium      int       `json:"sodium" db:"sodium"`
	ServingSize string    `json:"serving_size" db:"serving_size"`
	ServingUnit string    `json:"serving_unit" db:"serving_unit"`
	UserID      *uint     `json:"user_id" db:"user_id"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Recipe represents a recipe that combines multiple foods
type Recipe struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Instructions *string   `json:"instructions" db:"instructions"`
	PrepTime    *int      `json:"prep_time" db:"prep_time"`        // minutes
	CookTime    *int      `json:"cook_time" db:"cook_time"`        // minutes
	Servings    int       `json:"servings" db:"servings"`
	Category    *string   `json:"category" db:"category"`
	Cuisine     *string   `json:"cuisine" db:"cuisine"`
	Difficulty  *string   `json:"difficulty" db:"difficulty"`      // easy, medium, hard
	UserID      uint      `json:"user_id" db:"user_id"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RecipeIngredient represents ingredients in a recipe
type RecipeIngredient struct {
	ID        uint      `json:"id" db:"id"`
	RecipeID  uint      `json:"recipe_id" db:"recipe_id"`
	FoodID    uint      `json:"food_id" db:"food_id"`
	Quantity  float64   `json:"quantity" db:"quantity"`
	Unit      string    `json:"unit" db:"unit"`
	Notes     *string   `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	
	// Join for API responses
	Food      *Food    `json:"food,omitempty" db:"food"`
}

// FoodLog represents a food log entry for tracking meals
type FoodLog struct {
	ID             uint           `json:"id" db:"id"`
	UserID         uint           `json:"user_id" db:"user_id"`
	FoodID         *uint          `json:"food_id" db:"food_id"`
	RecipeID       *uint          `json:"recipe_id" db:"recipe_id"`
	MealType       string         `json:"meal_type" db:"meal_type"` // breakfast, lunch, dinner, snack
	Quantity       float64        `json:"quantity" db:"quantity"`
	Unit           string         `json:"unit" db:"unit"`
	LogDate        time.Time      `json:"log_date" db:"log_date"`
	Notes          *string        `json:"notes" db:"notes"`
	Calories       float64        `json:"calories" db:"calories"`
	Protein        float64        `json:"protein" db:"protein"`
	Carbs          float64        `json:"carbs" db:"carbs"`
	Fat            float64        `json:"fat" db:"fat"`
	Fiber          float64        `json:"fiber" db:"fiber"`
	Sugar          float64        `json:"sugar" db:"sugar"`
	Sodium         int            `json:"sodium" db:"sodium"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	
	// Joins for API responses
	Food   *Food   `json:"food,omitempty" db:"food"`
	Recipe *Recipe `json:"recipe,omitempty" db:"recipe"`
}

// MealPlan represents a meal plan for a specific date
type MealPlan struct {
	ID             uint      `json:"id" db:"id"`
	UserID         uint      `json:"user_id" db:"user_id"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description" db:"description"`
	StartDate      time.Time `json:"start_date" db:"start_date"`
	EndDate        time.Time `json:"end_date" db:"end_date"`
	TargetCalories *int      `json:"target_calories" db:"target_calories"`
	TargetProtein  *float64  `json:"target_protein" db:"target_protein"`
	TargetCarbs    *float64  `json:"target_carbs" db:"target_carbs"`
	TargetFat      *float64  `json:"target_fat" db:"target_fat"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// MealPlanItem represents a meal within a meal plan
type MealPlanItem struct {
	ID           uint      `json:"id" db:"id"`
	MealPlanID   uint      `json:"meal_plan_id" db:"meal_plan_id"`
	FoodID       *uint     `json:"food_id" db:"food_id"`
	RecipeID     *uint     `json:"recipe_id" db:"recipe_id"`
	MealType     string    `json:"meal_type" db:"meal_type"`
	DayOfWeek    int       `json:"day_of_week" db:"day_of_week"` // 0-6, Sunday=0
	Quantity     float64   `json:"quantity" db:"quantity"`
	Unit         string    `json:"unit" db:"unit"`
	Notes        *string   `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	
	// Joins for API responses
	Food   *Food   `json:"food,omitempty" db:"food"`
	Recipe *Recipe `json:"recipe,omitempty" db:"recipe"`
}

// DailyNutritionSummary represents daily nutrition totals
type DailyNutritionSummary struct {
	UserID    uint           `json:"user_id"`
	Date      time.Time      `json:"date"`
	Calories  float64        `json:"calories"`
	Protein   float64        `json:"protein"`
	Carbs     float64        `json:"carbs"`
	Fat       float64        `json:"fat"`
	Fiber     float64        `json:"fiber"`
	Sugar     float64        `json:"sugar"`
	Sodium    int            `json:"sodium"`
	Meals     []FoodLog      `json:"meals,omitempty"`
	Goals     *NutritionGoal `json:"goals,omitempty"`
}

// Request/Response DTOs

// CreateFoodRequest represents a request to create a food
type CreateFoodRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	Barcode     *string `json:"barcode,omitempty"`
	Category    *string `json:"category,omitempty"`
	Calories    float64 `json:"calories" validate:"required,min=0"`
	Protein     float64 `json:"protein" validate:"required,min=0"`
	Carbs       float64 `json:"carbs" validate:"required,min=0"`
	Fat         float64 `json:"fat" validate:"required,min=0"`
	Fiber       float64 `json:"fiber" validate:"min=0"`
	Sugar       float64 `json:"sugar" validate:"min=0"`
	Sodium      int     `json:"sodium" validate:"min=0"`
	ServingSize string  `json:"serving_size" validate:"required,min=1,max=50"`
	ServingUnit string  `json:"serving_unit" validate:"required,min=1,max=20"`
}

// UpdateFoodRequest represents a request to update a food
type UpdateFoodRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string `json:"description,omitempty"`
	Brand       *string `json:"brand,omitempty"`
	Barcode     *string `json:"barcode,omitempty"`
	Category    *string `json:"category,omitempty"`
	Calories    *float64 `json:"calories,omitempty" validate:"omitempty,min=0"`
	Protein     *float64 `json:"protein,omitempty" validate:"omitempty,min=0"`
	Carbs       *float64 `json:"carbs,omitempty" validate:"omitempty,min=0"`
	Fat         *float64 `json:"fat,omitempty" validate:"omitempty,min=0"`
	Fiber       *float64 `json:"fiber,omitempty" validate:"omitempty,min=0"`
	Sugar       *float64 `json:"sugar,omitempty" validate:"omitempty,min=0"`
	Sodium      *int     `json:"sodium,omitempty" validate:"omitempty,min=0"`
	ServingSize *string `json:"serving_size,omitempty" validate:"omitempty,min=1,max=50"`
	ServingUnit *string `json:"serving_unit,omitempty" validate:"omitempty,min=1,max=20"`
}

// CreateRecipeRequest represents a request to create a recipe
type CreateRecipeRequest struct {
	Name         string                     `json:"name" validate:"required,min=2,max=200"`
	Description  *string                    `json:"description,omitempty"`
	Instructions *string                    `json:"instructions,omitempty"`
	PrepTime     *int                       `json:"prep_time,omitempty" validate:"omitempty,min=0"`
	CookTime     *int                       `json:"cook_time,omitempty" validate:"omitempty,min=0"`
	Servings     int                        `json:"servings" validate:"required,min=1,max=100"`
	Category     *string                    `json:"category,omitempty"`
	Cuisine      *string                    `json:"cuisine,omitempty"`
	Difficulty   *string                    `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard"`
	IsPublic     bool                       `json:"is_public"`
	Ingredients  []CreateRecipeIngredient   `json:"ingredients" validate:"required,min=1,dive"`
}

// CreateRecipeIngredient represents a recipe ingredient
type CreateRecipeIngredient struct {
	FoodID   uint    `json:"food_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
	Unit     string  `json:"unit" validate:"required,min=1,max=20"`
	Notes    *string `json:"notes,omitempty"`
}

// UpdateRecipeRequest represents a request to update a recipe
type UpdateRecipeRequest struct {
	Name         *string                         `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description  *string                        `json:"description,omitempty"`
	Instructions *string                        `json:"instructions,omitempty"`
	PrepTime     *int                            `json:"prep_time,omitempty" validate:"omitempty,min=0"`
	CookTime     *int                            `json:"cook_time,omitempty" validate:"omitempty,min=0"`
	Servings     *int                            `json:"servings,omitempty" validate:"omitempty,min=1,max=100"`
	Category     *string                        `json:"category,omitempty"`
	Cuisine      *string                        `json:"cuisine,omitempty"`
	Difficulty   *string                        `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard"`
	IsPublic     *bool                           `json:"is_public,omitempty"`
	Ingredients  *[]CreateRecipeIngredient       `json:"ingredients,omitempty" validate:"omitempty,dive"`
}

// CreateFoodLogRequest represents a request to create a food log entry
type CreateFoodLogRequest struct {
	FoodID   *uint   `json:"food_id,omitempty"`
	RecipeID *uint   `json:"recipe_id,omitempty"`
	MealType string  `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
	Unit     string  `json:"unit" validate:"required,min=1,max=20"`
	LogDate  string  `json:"log_date" validate:"required"` // YYYY-MM-DD format
	Notes    *string `json:"notes,omitempty"`
}

// UpdateFoodLogRequest represents a request to update a food log entry
type UpdateFoodLogRequest struct {
	FoodID   *uint   `json:"food_id,omitempty"`
	RecipeID *uint   `json:"recipe_id,omitempty"`
	MealType *string `json:"meal_type,omitempty" validate:"omitempty,oneof=breakfast lunch dinner snack"`
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit     *string `json:"unit,omitempty" validate:"omitempty,min=1,max=20"`
	LogDate  *string `json:"log_date,omitempty"`
	Notes    *string `json:"notes,omitempty"`
}

// CreateMealPlanRequest represents a request to create a meal plan
type CreateMealPlanRequest struct {
	Name           string                `json:"name" validate:"required,min=2,max=200"`
	Description    *string               `json:"description,omitempty"`
	StartDate      string                `json:"start_date" validate:"required"` // YYYY-MM-DD format
	EndDate        string                `json:"end_date" validate:"required"`   // YYYY-MM-DD format
	TargetCalories *int                  `json:"target_calories,omitempty" validate:"omitempty,min=0"`
	TargetProtein  *float64              `json:"target_protein,omitempty" validate:"omitempty,min=0"`
	TargetCarbs    *float64              `json:"target_carbs,omitempty" validate:"omitempty,min=0"`
	TargetFat      *float64              `json:"target_fat,omitempty" validate:"omitempty,min=0"`
	Items          []CreateMealPlanItem  `json:"items,omitempty" validate:"omitempty,dive"`
}

// CreateMealPlanItem represents a meal plan item
type CreateMealPlanItem struct {
	FoodID   *uint   `json:"food_id,omitempty"`
	RecipeID *uint   `json:"recipe_id,omitempty"`
	MealType string  `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	DayOfWeek int    `json:"day_of_week" validate:"required,min=0,max=6"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
	Unit     string  `json:"unit" validate:"required,min=1,max=20"`
	Notes    *string `json:"notes,omitempty"`
}

// UpdateMealPlanRequest represents a request to update a meal plan
type UpdateMealPlanRequest struct {
	Name           *string               `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description    *string               `json:"description,omitempty"`
	StartDate      *string               `json:"start_date,omitempty"`
	EndDate        *string               `json:"end_date,omitempty"`
	TargetCalories *int                  `json:"target_calories,omitempty" validate:"omitempty,min=0"`
	TargetProtein  *float64              `json:"target_protein,omitempty" validate:"omitempty,min=0"`
	TargetCarbs    *float64              `json:"target_carbs,omitempty" validate:"omitempty,min=0"`
	TargetFat      *float64              `json:"target_fat,omitempty" validate:"omitempty,min=0"`
	IsActive       *bool                 `json:"is_active,omitempty"`
	Items          *[]CreateMealPlanItem  `json:"items,omitempty" validate:"omitempty,dive"`
}

// SearchParams represents search parameters
type SearchParams struct {
	Query    string `query:"q"`
	Category string `query:"category"`
	UserOnly bool   `query:"user_only"`
	Verified bool   `query:"verified"`
}

// FoodLogQueryParams represents food log query parameters
type FoodLogQueryParams struct {
	StartDate string `query:"start_date" validate:"omitempty"`
	EndDate   string `query:"end_date" validate:"omitempty"`
	MealType  string `query:"meal_type" validate:"omitempty,oneof=breakfast lunch dinner snack"`
}

// Helper methods

// IsValidMealType checks if meal type is valid
func IsValidMealType(mealType string) bool {
	validTypes := []string{"breakfast", "lunch", "dinner", "snack"}
	for _, valid := range validTypes {
		if mealType == valid {
			return true
		}
	}
	return false
}

// IsValidDifficulty checks if difficulty is valid
func IsValidDifficulty(difficulty string) bool {
	validDifficulties := []string{"easy", "medium", "hard"}
	for _, valid := range validDifficulties {
		if difficulty == valid {
			return true
		}
	}
	return false
}

// ParseDate parses a YYYY-MM-DD string to time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ToFoodLogWithJoins converts FoodLog to include food/recipe data
func (fl *FoodLog) ToFoodLogWithJoins(food *Food, recipe *Recipe) *FoodLog {
	fl.Food = food
	fl.Recipe = recipe
	return fl
}

// ToMealPlanItemWithJoins converts MealPlanItem to include food/recipe data
func (mpi *MealPlanItem) ToMealPlanItemWithJoins(food *Food, recipe *Recipe) *MealPlanItem {
	mpi.Food = food
	mpi.Recipe = recipe
	return mpi
}

// CalculateNutrition calculates total nutrition for a recipe
func (r *Recipe) CalculateNutrition(ingredients []RecipeIngredient) (calories, protein, carbs, fat float64) {
	for _, ing := range ingredients {
		if ing.Food != nil {
			food := ing.Food
			// Calculate nutrition based on quantity ratio
			ratio := ing.Quantity / food.CalculateServingWeight()
			
			calories += food.Calories * ratio
			protein += food.Protein * ratio
			carbs += food.Carbs * ratio
			fat += food.Fat * ratio
		}
	}
	return calories, protein, carbs, fat
}

// CalculateServingWeight estimates serving weight based on serving size
func (f *Food) CalculateServingWeight() float64 {
	// This is a simplified calculation
	// In a real implementation, you'd have a database of food densities
	return 100.0 // Default to 100g
}

// CalculateDailyTotal calculates daily nutrition from food logs
func CalculateDailyTotal(logs []FoodLog) (calories, protein, carbs, fat, fiber, sugar float64, sodium int) {
	for _, log := range logs {
		calories += log.Calories
		protein += log.Protein
		carbs += log.Carbs
		fat += log.Fat
		fiber += log.Fiber
		sugar += log.Sugar
		sodium += log.Sodium
	}
	return calories, protein, carbs, fat, fiber, sugar, sodium
}
