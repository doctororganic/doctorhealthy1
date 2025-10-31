package models

import (
	"nutrition-platform/errors"
	"time"
)

// Food represents a food item in the database
type Food struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;index" validate:"required,min=1,max=200"`
	Category    string    `json:"category" gorm:"not null" validate:"required,oneof=fruits vegetables grains proteins dairy fats beverages other"`
	Calories    float64   `json:"calories" gorm:"not null" validate:"min=0"`            // per 100g
	Protein     float64   `json:"protein" gorm:"not null;default:0" validate:"min=0"`   // grams per 100g
	Carbs       float64   `json:"carbs" gorm:"not null;default:0" validate:"min=0"`     // grams per 100g
	Fat         float64   `json:"fat" gorm:"not null;default:0" validate:"min=0"`       // grams per 100g
	Fiber       float64   `json:"fiber" gorm:"not null;default:0" validate:"min=0"`     // grams per 100g
	Sugar       float64   `json:"sugar" gorm:"not null;default:0" validate:"min=0"`     // grams per 100g
	Sodium      float64   `json:"sodium" gorm:"not null;default:0" validate:"min=0"`    // mg per 100g
	Potassium   float64   `json:"potassium" gorm:"not null;default:0" validate:"min=0"` // mg per 100g
	Vitamin_C   float64   `json:"vitamin_c" gorm:"not null;default:0" validate:"min=0"` // mg per 100g
	Calcium     float64   `json:"calcium" gorm:"not null;default:0" validate:"min=0"`   // mg per 100g
	Iron        float64   `json:"iron" gorm:"not null;default:0" validate:"min=0"`      // mg per 100g
	ServingSize string    `json:"serving_size" gorm:"not null;default:'100g'" validate:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for the Food model
func (Food) TableName() string {
	return "foods"
}

// Validate validates the food model
func (f *Food) Validate() error {
	if f.Name == "" {
		return errors.ErrInvalidInputError("name is required")
	}
	if len(f.Name) > 200 {
		return errors.ErrInvalidInputError("name must be less than 200 characters")
	}
	if f.Category == "" {
		return errors.ErrInvalidInputError("category is required")
	}
	validCategories := []string{"fruits", "vegetables", "grains", "proteins", "dairy", "fats", "beverages", "other"}
	if !contains(validCategories, f.Category) {
		return errors.ErrInvalidInputError("invalid category")
	}
	if f.Calories < 0 {
		return errors.ErrInvalidInputError("calories cannot be negative")
	}
	if f.Protein < 0 || f.Carbs < 0 || f.Fat < 0 || f.Fiber < 0 || f.Sugar < 0 {
		return errors.ErrInvalidInputError("nutrient values cannot be negative")
	}
	return nil
}

// CalculateNutrition calculates nutrition info for a specific quantity
func (f *Food) CalculateNutrition(quantity float64) *NutritionInfo {
	scale := quantity / 100.0 // nutrients are per 100g

	return &NutritionInfo{
		Calories:      f.Calories * scale,
		Protein:       f.Protein * scale,
		Carbohydrates: f.Carbs * scale,
		Fat:           f.Fat * scale,
		Fiber:         f.Fiber * scale,
		Sugar:         f.Sugar * scale,
		Sodium:        f.Sodium * scale,
		VitaminC:      f.Vitamin_C * scale,
		Calcium:       f.Calcium * scale,
		Iron:          f.Iron * scale,
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
