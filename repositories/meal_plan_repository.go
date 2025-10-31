package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"nutrition-platform/models"
)

type MealPlanRepository struct {
	db *sql.DB
}

func NewMealPlanRepository(db *sql.DB) *MealPlanRepository {
	return &MealPlanRepository{db: db}
}

// CreateMealPlan creates a new meal plan with its meals
func (r *MealPlanRepository) CreateMealPlan(mealPlan *models.MealPlan) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert meal plan
	mealPlanQuery := `
		INSERT INTO meal_plans (user_id, name, description, start_date, end_date, 
			goal_type, target_calories, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err = tx.QueryRow(mealPlanQuery,
		mealPlan.UserID,
		mealPlan.Name,
		mealPlan.Description,
		mealPlan.StartDate,
		mealPlan.EndDate,
		mealPlan.GoalType,
		mealPlan.TargetCalories,
		mealPlan.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&mealPlan.ID)

	if err != nil {
		log.Printf("Error creating meal plan: %v", err)
		return fmt.Errorf("failed to create meal plan: %w", err)
	}

	// Insert planned meals
	for _, plannedMeal := range mealPlan.PlannedMeals {
		plannedMealQuery := `
			INSERT INTO planned_meals (meal_plan_id, day_of_week, meal_type, food_id, recipe_id, quantity, unit)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := tx.Exec(plannedMealQuery,
			mealPlan.ID,
			plannedMeal.DayOfWeek,
			plannedMeal.MealType,
			plannedMeal.FoodID,
			plannedMeal.RecipeID,
			plannedMeal.Quantity,
			plannedMeal.Unit,
		)

		if err != nil {
			log.Printf("Error creating planned meal: %v", err)
			return fmt.Errorf("failed to create planned meal: %w", err)
		}
	}

	return tx.Commit()
}

// GetMealPlanByID retrieves a meal plan by its ID with all planned meals
func (r *MealPlanRepository) GetMealPlanByID(id, userID string) (*models.MealPlan, error) {
	// Get meal plan details
	mealPlanQuery := `
		SELECT id, user_id, name, description, start_date, end_date, 
			goal_type, target_calories, is_active, created_at, updated_at
		FROM meal_plans 
		WHERE id = $1 AND user_id = $2`

	var mealPlan models.MealPlan

	err := r.db.QueryRow(mealPlanQuery, id, userID).Scan(
		&mealPlan.ID,
		&mealPlan.UserID,
		&mealPlan.Name,
		&mealPlan.Description,
		&mealPlan.StartDate,
		&mealPlan.EndDate,
		&mealPlan.GoalType,
		&mealPlan.TargetCalories,
		&mealPlan.IsActive,
		&mealPlan.CreatedAt,
		&mealPlan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("meal plan not found")
		}
		return nil, fmt.Errorf("failed to get meal plan: %w", err)
	}

	// Get planned meals
	plannedMealsQuery := `
		SELECT day_of_week, meal_type, food_id, recipe_id, quantity, unit
		FROM planned_meals
		WHERE meal_plan_id = $1
		ORDER BY day_of_week, meal_type`

	rows, err := r.db.Query(plannedMealsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get planned meals: %w", err)
	}
	defer rows.Close()

	var plannedMeals []*models.PlannedMeal
	for rows.Next() {
		var plannedMeal models.PlannedMeal
		var foodID, recipeID sql.NullString

		err := rows.Scan(
			&plannedMeal.DayOfWeek,
			&plannedMeal.MealType,
			&foodID,
			&recipeID,
			&plannedMeal.Quantity,
			&plannedMeal.Unit,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan planned meal row: %w", err)
		}

		if foodID.Valid {
			plannedMeal.FoodID = foodID.String
		}
		if recipeID.Valid {
			plannedMeal.RecipeID = recipeID.String
		}

		plannedMeals = append(plannedMeals, &plannedMeal)
	}
	mealPlan.PlannedMeals = plannedMeals

	return &mealPlan, nil
}

// GetUserMealPlans retrieves all meal plans for a user
func (r *MealPlanRepository) GetUserMealPlans(userID string, limit, offset int) ([]*models.MealPlan, error) {
	query := `
		SELECT id, user_id, name, description, start_date, end_date, 
			goal_type, target_calories, is_active, created_at, updated_at
		FROM meal_plans 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meal plans: %w", err)
	}
	defer rows.Close()

	var mealPlans []*models.MealPlan
	for rows.Next() {
		var mealPlan models.MealPlan

		err := rows.Scan(
			&mealPlan.ID,
			&mealPlan.UserID,
			&mealPlan.Name,
			&mealPlan.Description,
			&mealPlan.StartDate,
			&mealPlan.EndDate,
			&mealPlan.GoalType,
			&mealPlan.TargetCalories,
			&mealPlan.IsActive,
			&mealPlan.CreatedAt,
			&mealPlan.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan meal plan row: %w", err)
		}

		mealPlans = append(mealPlans, &mealPlan)
	}

	return mealPlans, nil
}

// GetActiveMealPlan retrieves the active meal plan for a user
func (r *MealPlanRepository) GetActiveMealPlan(userID string) (*models.MealPlan, error) {
	query := `
		SELECT id, user_id, name, description, start_date, end_date, 
			goal_type, target_calories, is_active, created_at, updated_at
		FROM meal_plans 
		WHERE user_id = $1 AND is_active = true
		LIMIT 1`

	var mealPlan models.MealPlan

	err := r.db.QueryRow(query, userID).Scan(
		&mealPlan.ID,
		&mealPlan.UserID,
		&mealPlan.Name,
		&mealPlan.Description,
		&mealPlan.StartDate,
		&mealPlan.EndDate,
		&mealPlan.GoalType,
		&mealPlan.TargetCalories,
		&mealPlan.IsActive,
		&mealPlan.CreatedAt,
		&mealPlan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active meal plan found")
		}
		return nil, fmt.Errorf("failed to get active meal plan: %w", err)
	}

	return &mealPlan, nil
}

// UpdateMealPlan updates an existing meal plan and its planned meals
func (r *MealPlanRepository) UpdateMealPlan(mealPlan *models.MealPlan) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update meal plan
	mealPlanQuery := `
		UPDATE meal_plans 
		SET name = $2, description = $3, start_date = $4, end_date = $5,
			goal_type = $6, target_calories = $7, is_active = $8, updated_at = $9
		WHERE id = $1 AND user_id = $10`

	_, err = tx.Exec(mealPlanQuery,
		mealPlan.ID,
		mealPlan.Name,
		mealPlan.Description,
		mealPlan.StartDate,
		mealPlan.EndDate,
		mealPlan.GoalType,
		mealPlan.TargetCalories,
		mealPlan.IsActive,
		time.Now(),
		mealPlan.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update meal plan: %w", err)
	}

	// Delete existing planned meals
	_, err = tx.Exec("DELETE FROM planned_meals WHERE meal_plan_id = $1", mealPlan.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing planned meals: %w", err)
	}

	// Insert new planned meals
	for _, plannedMeal := range mealPlan.PlannedMeals {
		plannedMealQuery := `
			INSERT INTO planned_meals (meal_plan_id, day_of_week, meal_type, food_id, recipe_id, quantity, unit)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err := tx.Exec(plannedMealQuery,
			mealPlan.ID,
			plannedMeal.DayOfWeek,
			plannedMeal.MealType,
			plannedMeal.FoodID,
			plannedMeal.RecipeID,
			plannedMeal.Quantity,
			plannedMeal.Unit,
		)

		if err != nil {
			return fmt.Errorf("failed to create planned meal: %w", err)
		}
	}

	return tx.Commit()
}

// DeleteMealPlan deletes a meal plan
func (r *MealPlanRepository) DeleteMealPlan(id, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete planned meals first
	_, err = tx.Exec("DELETE FROM planned_meals WHERE meal_plan_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete planned meals: %w", err)
	}

	// Delete meal plan
	_, err = tx.Exec("DELETE FROM meal_plans WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete meal plan: %w", err)
	}

	return tx.Commit()
}

// DeactivateAllMealPlans deactivates all meal plans for a user
func (r *MealPlanRepository) DeactivateAllMealPlans(userID string) error {
	query := `UPDATE meal_plans SET is_active = false, updated_at = $1 WHERE user_id = $2`
	
	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate meal plans: %w", err)
	}

	return nil
}

// SetActiveMealPlan sets a meal plan as active and deactivates others
func (r *MealPlanRepository) SetActiveMealPlan(id, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Deactivate all meal plans for user
	_, err = tx.Exec("UPDATE meal_plans SET is_active = false, updated_at = $1 WHERE user_id = $2", time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate meal plans: %w", err)
	}

	// Activate specified meal plan
	_, err = tx.Exec("UPDATE meal_plans SET is_active = true, updated_at = $1 WHERE id = $2 AND user_id = $3", time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to activate meal plan: %w", err)
	}

	return tx.Commit()
}

// GetMealPlanByDateRange retrieves meal plans for a user within a date range
func (r *MealPlanRepository) GetMealPlanByDateRange(userID string, startDate, endDate time.Time) ([]*models.MealPlan, error) {
	query := `
		SELECT id, user_id, name, description, start_date, end_date, 
			goal_type, target_calories, is_active, created_at, updated_at
		FROM meal_plans 
		WHERE user_id = $1 AND (
			(start_date <= $2 AND end_date >= $2) OR
			(start_date <= $3 AND end_date >= $3) OR
			(start_date >= $2 AND end_date <= $3)
		)
		ORDER BY start_date ASC`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal plans by date range: %w", err)
	}
	defer rows.Close()

	var mealPlans []*models.MealPlan
	for rows.Next() {
		var mealPlan models.MealPlan

		err := rows.Scan(
			&mealPlan.ID,
			&mealPlan.UserID,
			&mealPlan.Name,
			&mealPlan.Description,
			&mealPlan.StartDate,
			&mealPlan.EndDate,
			&mealPlan.GoalType,
			&mealPlan.TargetCalories,
			&mealPlan.IsActive,
			&mealPlan.CreatedAt,
			&mealPlan.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan meal plan row: %w", err)
		}

		mealPlans = append(mealPlans, &mealPlan)
	}

	return mealPlans, nil
}

// GetPlannedMealsForDate retrieves planned meals for a specific date
func (r *MealPlanRepository) GetPlannedMealsForDate(userID string, date time.Time) ([]*models.PlannedMeal, error) {
	// Get day of week as string (0-6, Sunday=0)
	dayOfWeek := date.Weekday().String()
	
	query := `
		SELECT pm.day_of_week, pm.meal_type, pm.food_id, pm.recipe_id, pm.quantity, pm.unit
		FROM planned_meals pm
		INNER JOIN meal_plans mp ON pm.meal_plan_id = mp.id
		WHERE mp.user_id = $1 AND mp.is_active = true 
			AND pm.day_of_week = $2
			AND mp.start_date <= $3 AND mp.end_date >= $3
		ORDER BY pm.meal_type`

	rows, err := r.db.Query(query, userID, dayOfWeek, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get planned meals for date: %w", err)
	}
	defer rows.Close()

	var plannedMeals []*models.PlannedMeal
	for rows.Next() {
		var plannedMeal models.PlannedMeal
		var foodID, recipeID sql.NullString

		err := rows.Scan(
			&plannedMeal.DayOfWeek,
			&plannedMeal.MealType,
			&foodID,
			&recipeID,
			&plannedMeal.Quantity,
			&plannedMeal.Unit,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan planned meal row: %w", err)
		}

		if foodID.Valid {
			plannedMeal.FoodID = foodID.String
		}
		if recipeID.Valid {
			plannedMeal.RecipeID = recipeID.String
		}

		plannedMeals = append(plannedMeals, &plannedMeal)
	}

	return plannedMeals, nil
}
