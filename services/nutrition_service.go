package services

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// NutritionService handles nutrition-related business logic
type NutritionService struct {
	db *sql.DB
}

// NewNutritionService creates a new nutrition service
func NewNutritionService(db *sql.DB) *NutritionService {
	return &NutritionService{
		db: db,
	}
}

// GetMeals retrieves meals for a user
func (s *NutritionService) GetMeals(userID uint, limit, offset int) ([]models.Meal, error) {
	query := `
		SELECT id, user_id, name, description, calories, protein, carbs, fat,
			   fiber, sugar, sodium, meal_type, consumed_at, created_at, updated_at
		FROM meals 
		WHERE user_id = $1 
		ORDER BY consumed_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var meals []models.Meal
	for rows.Next() {
		var meal models.Meal
		err := rows.Scan(
			&meal.ID, &meal.UserID, &meal.Name, &meal.Description,
			&meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat,
			&meal.Fiber, &meal.Sugar, &meal.Sodium, &meal.MealType,
			&meal.ConsumedAt, &meal.CreatedAt, &meal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	
	return meals, nil
}

// CreateMeal creates a new meal
func (s *NutritionService) CreateMeal(meal *models.Meal) error {
	query := `
		INSERT INTO meals (user_id, name, description, calories, protein, carbs, fat,
						  fiber, sugar, sodium, meal_type, consumed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	
	now := time.Now()
	if meal.ConsumedAt.IsZero() {
		meal.ConsumedAt = now
	}
	
	err := s.db.QueryRow(query,
		meal.UserID, meal.Name, meal.Description, meal.Calories, meal.Protein,
		meal.Carbs, meal.Fat, meal.Fiber, meal.Sugar, meal.Sodium,
		meal.MealType, meal.ConsumedAt, now, now,
	).Scan(&meal.ID)
	
	if err != nil {
		return err
	}
	
	meal.CreatedAt = now
	meal.UpdatedAt = now
	
	return nil
}

// GetMeal retrieves a meal by ID
func (s *NutritionService) GetMeal(id, userID uint) (*models.Meal, error) {
	query := `
		SELECT id, user_id, name, description, calories, protein, carbs, fat,
			   fiber, sugar, sodium, meal_type, consumed_at, created_at, updated_at
		FROM meals 
		WHERE id = $1 AND user_id = $2
	`
	
	var meal models.Meal
	err := s.db.QueryRow(query, id, userID).Scan(
		&meal.ID, &meal.UserID, &meal.Name, &meal.Description,
		&meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat,
		&meal.Fiber, &meal.Sugar, &meal.Sodium, &meal.MealType,
		&meal.ConsumedAt, &meal.CreatedAt, &meal.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("meal not found")
		}
		return nil, err
	}
	
	return &meal, nil
}

// UpdateMeal updates an existing meal
func (s *NutritionService) UpdateMeal(meal *models.Meal) error {
	query := `
		UPDATE meals 
		SET name = $2, description = $3, calories = $4, protein = $5, carbs = $6,
			fat = $7, fiber = $8, sugar = $9, sodium = $10, meal_type = $11,
			consumed_at = $12, updated_at = $13
		WHERE id = $1 AND user_id = $14
	`
	
	meal.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(query,
		meal.ID, meal.Name, meal.Description, meal.Calories, meal.Protein,
		meal.Carbs, meal.Fat, meal.Fiber, meal.Sugar, meal.Sodium,
		meal.MealType, meal.ConsumedAt, meal.UpdatedAt, meal.UserID,
	)
	
	return err
}

// DeleteMeal deletes a meal
func (s *NutritionService) DeleteMeal(id, userID uint) error {
	query := `DELETE FROM meals WHERE id = $1 AND user_id = $2`
	
	_, err := s.db.Exec(query, id, userID)
	return err
}

// GetNutritionPlans retrieves nutrition plans for a user
func (s *NutritionService) GetNutritionPlans(userID uint) ([]models.NutritionPlan, error) {
	query := `
		SELECT id, user_id, name, description, daily_calories, protein_grams,
			   carbs_grams, fat_grams, fiber_grams, sugar_grams, sodium_mg,
			   is_active, start_date, end_date, created_at, updated_at
		FROM nutrition_plans 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var plans []models.NutritionPlan
	for rows.Next() {
		var plan models.NutritionPlan
		var startDate, endDate sql.NullTime
		
		err := rows.Scan(
			&plan.ID, &plan.UserID, &plan.Name, &plan.Description,
			&plan.DailyCalories, &plan.ProteinGrams, &plan.CarbsGrams,
			&plan.FatGrams, &plan.FiberGrams, &plan.SugarGrams,
			&plan.SodiumMg, &plan.IsActive, &startDate, &endDate,
			&plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if startDate.Valid {
			plan.StartDate = &startDate.Time
		}
		if endDate.Valid {
			plan.EndDate = &endDate.Time
		}
		
		plans = append(plans, plan)
	}
	
	return plans, nil
}

// CreateNutritionPlan creates a new nutrition plan
func (s *NutritionService) CreateNutritionPlan(plan *models.NutritionPlan) error {
	query := `
		INSERT INTO nutrition_plans (user_id, name, description, daily_calories,
									protein_grams, carbs_grams, fat_grams, fiber_grams,
									sugar_grams, sodium_mg, is_active, start_date, 
									end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		plan.UserID, plan.Name, plan.Description, plan.DailyCalories,
		plan.ProteinGrams, plan.CarbsGrams, plan.FatGrams, plan.FiberGrams,
		plan.SugarGrams, plan.SodiumMg, plan.IsActive, plan.StartDate,
		plan.EndDate, now, now,
	).Scan(&plan.ID)
	
	if err != nil {
		return err
	}
	
	plan.CreatedAt = now
	plan.UpdatedAt = now
	
	return nil
}
