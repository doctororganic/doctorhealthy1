package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"nutrition-platform/models"
)

type FoodLogRepository struct {
	db *sql.DB
}

func NewFoodLogRepository(db *sql.DB) *FoodLogRepository {
	return &FoodLogRepository{db: db}
}

// CreateFoodLog creates a new food log entry
func (r *FoodLogRepository) CreateFoodLog(logEntry *models.FoodLog) error {
	query := `
		INSERT INTO food_logs (user_id, food_id, recipe_id, meal_type, quantity, unit, 
			log_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := r.db.QueryRow(query,
		logEntry.UserID,
		logEntry.FoodID,
		logEntry.RecipeID,
		logEntry.MealType,
		logEntry.Quantity,
		logEntry.Unit,
		logEntry.LogDate,
		time.Now(),
		time.Now(),
	).Scan(&logEntry.ID)

	if err != nil {
		log.Printf("Error creating food log: %v", err)
		return fmt.Errorf("failed to create food log: %w", err)
	}

	return nil
}

// GetFoodLogByID retrieves a food log entry by its ID
func (r *FoodLogRepository) GetFoodLogByID(id, userID string) (*models.FoodLog, error) {
	query := `
		SELECT id, user_id, food_id, recipe_id, meal_type, quantity, unit, 
			log_date, created_at, updated_at
		FROM food_logs 
		WHERE id = $1 AND user_id = $2`

	var logEntry models.FoodLog
	var foodID, recipeID sql.NullString

	err := r.db.QueryRow(query, id, userID).Scan(
		&logEntry.ID,
		&logEntry.UserID,
		&foodID,
		&recipeID,
		&logEntry.MealType,
		&logEntry.Quantity,
		&logEntry.Unit,
		&logEntry.LogDate,
		&logEntry.CreatedAt,
		&logEntry.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("food log not found")
		}
		return nil, fmt.Errorf("failed to get food log: %w", err)
	}

	if foodID.Valid {
		logEntry.FoodID = foodID.String
	}
	if recipeID.Valid {
		logEntry.RecipeID = recipeID.String
	}

	return &logEntry, nil
}

// GetFoodLogsByDate retrieves all food log entries for a user on a specific date
func (r *FoodLogRepository) GetFoodLogsByDate(userID string, date time.Time) ([]*models.FoodLog, error) {
	query := `
		SELECT id, user_id, food_id, recipe_id, meal_type, quantity, unit, 
			log_date, created_at, updated_at
		FROM food_logs 
		WHERE user_id = $1 AND DATE(log_date) = DATE($2)
		ORDER BY log_date ASC`

	rows, err := r.db.Query(query, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get food logs by date: %w", err)
	}
	defer rows.Close()

	var logs []*models.FoodLog
	for rows.Next() {
		var logEntry models.FoodLog
		var foodID, recipeID sql.NullString

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.UserID,
			&foodID,
			&recipeID,
			&logEntry.MealType,
			&logEntry.Quantity,
			&logEntry.Unit,
			&logEntry.LogDate,
			&logEntry.CreatedAt,
			&logEntry.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan food log row: %w", err)
		}

		if foodID.Valid {
			logEntry.FoodID = foodID.String
		}
		if recipeID.Valid {
			logEntry.RecipeID = recipeID.String
		}

		logs = append(logs, &logEntry)
	}

	return logs, nil
}

// GetFoodLogsByDateRange retrieves food log entries for a user within a date range
func (r *FoodLogRepository) GetFoodLogsByDateRange(userID string, startDate, endDate time.Time) ([]*models.FoodLog, error) {
	query := `
		SELECT id, user_id, food_id, recipe_id, meal_type, quantity, unit, 
			log_date, created_at, updated_at
		FROM food_logs 
		WHERE user_id = $1 AND log_date >= $2 AND log_date <= $3
		ORDER BY log_date ASC`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get food logs by date range: %w", err)
	}
	defer rows.Close()

	var logs []*models.FoodLog
	for rows.Next() {
		var logEntry models.FoodLog
		var foodID, recipeID sql.NullString

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.UserID,
			&foodID,
			&recipeID,
			&logEntry.MealType,
			&logEntry.Quantity,
			&logEntry.Unit,
			&logEntry.LogDate,
			&logEntry.CreatedAt,
			&logEntry.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan food log row: %w", err)
		}

		if foodID.Valid {
			logEntry.FoodID = foodID.String
		}
		if recipeID.Valid {
			logEntry.RecipeID = recipeID.String
		}

		logs = append(logs, &logEntry)
	}

	return logs, nil
}

// GetFoodLogsByMealType retrieves food log entries for a user by meal type on a specific date
func (r *FoodLogRepository) GetFoodLogsByMealType(userID string, mealType string, date time.Time) ([]*models.FoodLog, error) {
	query := `
		SELECT id, user_id, food_id, recipe_id, meal_type, quantity, unit, 
			log_date, created_at, updated_at
		FROM food_logs 
		WHERE user_id = $1 AND meal_type = $2 AND DATE(log_date) = DATE($3)
		ORDER BY log_date ASC`

	rows, err := r.db.Query(query, userID, mealType, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get food logs by meal type: %w", err)
	}
	defer rows.Close()

	var logs []*models.FoodLog
	for rows.Next() {
		var logEntry models.FoodLog
		var foodID, recipeID sql.NullString

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.UserID,
			&foodID,
			&recipeID,
			&logEntry.MealType,
			&logEntry.Quantity,
			&logEntry.Unit,
			&logEntry.LogDate,
			&logEntry.CreatedAt,
			&logEntry.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan food log row: %w", err)
		}

		if foodID.Valid {
			logEntry.FoodID = foodID.String
		}
		if recipeID.Valid {
			logEntry.RecipeID = recipeID.String
		}

		logs = append(logs, &logEntry)
	}

	return logs, nil
}

// UpdateFoodLog updates an existing food log entry
func (r *FoodLogRepository) UpdateFoodLog(logEntry *models.FoodLog) error {
	query := `
		UPDATE food_logs 
		SET food_id = $2, recipe_id = $3, meal_type = $4, quantity = $5, 
			unit = $6, log_date = $7, updated_at = $8
		WHERE id = $1 AND user_id = $9`

	_, err := r.db.Exec(query,
		logEntry.ID,
		logEntry.FoodID,
		logEntry.RecipeID,
		logEntry.MealType,
		logEntry.Quantity,
		logEntry.Unit,
		logEntry.LogDate,
		time.Now(),
		logEntry.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update food log: %w", err)
	}

	return nil
}

// DeleteFoodLog deletes a food log entry
func (r *FoodLogRepository) DeleteFoodLog(id, userID string) error {
	query := `DELETE FROM food_logs WHERE id = $1 AND user_id = $2`
	
	_, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete food log: %w", err)
	}

	return nil
}

// GetDailyNutritionSummary calculates daily nutrition totals for a user on a specific date
func (r *FoodLogRepository) GetDailyNutritionSummary(userID string, date time.Time) (*models.DailyNutritionSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(fl.quantity * f.calories / f.serving_size_numeric), 0) as total_calories,
			COALESCE(SUM(fl.quantity * f.protein / f.serving_size_numeric), 0) as total_protein,
			COALESCE(SUM(fl.quantity * f.carbs / f.serving_size_numeric), 0) as total_carbs,
			COALESCE(SUM(fl.quantity * f.fat / f.serving_size_numeric), 0) as total_fat,
			COALESCE(SUM(fl.quantity * f.saturated_fat / f.serving_size_numeric), 0) as total_saturated_fat,
			COALESCE(SUM(fl.quantity * f.fiber / f.serving_size_numeric), 0) as total_fiber,
			COALESCE(SUM(fl.quantity * f.sugar / f.serving_size_numeric), 0) as total_sugar,
			COALESCE(SUM(fl.quantity * f.sodium / f.serving_size_numeric), 0) as total_sodium,
			COUNT(fl.id) as total_items
		FROM food_logs fl
		LEFT JOIN foods f ON fl.food_id = f.id
		WHERE fl.user_id = $1 AND DATE(fl.log_date) = DATE($2)`

	var summary models.DailyNutritionSummary
	err := r.db.QueryRow(query, userID, date).Scan(
		&summary.TotalCalories,
		&summary.TotalProtein,
		&summary.TotalCarbs,
		&summary.TotalFat,
		&summary.TotalSaturatedFat,
		&summary.TotalFiber,
		&summary.TotalSugar,
		&summary.TotalSodium,
		&summary.TotalItems,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get daily nutrition summary: %w", err)
	}

	summary.Date = date
	summary.UserID = userID

	return &summary, nil
}

// GetWeeklyNutritionSummary calculates weekly nutrition totals for a user
func (r *FoodLogRepository) GetWeeklyNutritionSummary(userID string, startDate time.Time) ([]*models.DailyNutritionSummary, error) {
	query := `
		SELECT 
			DATE(fl.log_date) as date,
			COALESCE(SUM(fl.quantity * f.calories / f.serving_size_numeric), 0) as total_calories,
			COALESCE(SUM(fl.quantity * f.protein / f.serving_size_numeric), 0) as total_protein,
			COALESCE(SUM(fl.quantity * f.carbs / f.serving_size_numeric), 0) as total_carbs,
			COALESCE(SUM(fl.quantity * f.fat / f.serving_size_numeric), 0) as total_fat,
			COALESCE(SUM(fl.quantity * f.saturated_fat / f.serving_size_numeric), 0) as total_saturated_fat,
			COALESCE(SUM(fl.quantity * f.fiber / f.serving_size_numeric), 0) as total_fiber,
			COALESCE(SUM(fl.quantity * f.sugar / f.serving_size_numeric), 0) as total_sugar,
			COALESCE(SUM(fl.quantity * f.sodium / f.serving_size_numeric), 0) as total_sodium,
			COUNT(fl.id) as total_items
		FROM food_logs fl
		LEFT JOIN foods f ON fl.food_id = f.id
		WHERE fl.user_id = $1 AND fl.log_date >= $2 AND fl.log_date < $2 + INTERVAL '7 days'
		GROUP BY DATE(fl.log_date)
		ORDER BY date ASC`

	rows, err := r.db.Query(query, userID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly nutrition summary: %w", err)
	}
	defer rows.Close()

	var summaries []*models.DailyNutritionSummary
	for rows.Next() {
		var summary models.DailyNutritionSummary
		err := rows.Scan(
			&summary.Date,
			&summary.TotalCalories,
			&summary.TotalProtein,
			&summary.TotalCarbs,
			&summary.TotalFat,
			&summary.TotalSaturatedFat,
			&summary.TotalFiber,
			&summary.TotalSugar,
			&summary.TotalSodium,
			&summary.TotalItems,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan daily summary row: %w", err)
		}

		summary.UserID = userID
		summaries = append(summaries, &summary)
	}

	return summaries, nil
}

// SearchFoodLogs searches for food logs with filters
func (r *FoodLogRepository) SearchFoodLogs(userID string, filters models.FoodLogSearchFilters, limit, offset int) ([]*models.FoodLog, error) {
	whereClauses := []string{"user_id = $1"}
	args := []interface{}{userID}
	argIndex := 2

	if !filters.StartDate.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprintf("log_date >= $%d", argIndex))
		args = append(args, filters.StartDate)
		argIndex++
	}

	if !filters.EndDate.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprintf("log_date <= $%d", argIndex))
		args = append(args, filters.EndDate)
		argIndex++
	}

	if filters.MealType != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("meal_type = $%d", argIndex))
		args = append(args, filters.MealType)
		argIndex++
	}

	baseQuery := `
		SELECT id, user_id, food_id, recipe_id, meal_type, quantity, unit, 
			log_date, created_at, updated_at
		FROM food_logs`

	whereClause := " WHERE " + strings.Join(whereClauses, " AND ")
	orderBy := " ORDER BY log_date DESC"
	pagination := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	fullQuery := baseQuery + whereClause + orderBy + pagination

	rows, err := r.db.Query(fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search food logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.FoodLog
	for rows.Next() {
		var logEntry models.FoodLog
		var foodID, recipeID sql.NullString

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.UserID,
			&foodID,
			&recipeID,
			&logEntry.MealType,
			&logEntry.Quantity,
			&logEntry.Unit,
			&logEntry.LogDate,
			&logEntry.CreatedAt,
			&logEntry.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan food log row: %w", err)
		}

		if foodID.Valid {
			logEntry.FoodID = foodID.String
		}
		if recipeID.Valid {
			logEntry.RecipeID = recipeID.String
		}

		logs = append(logs, &logEntry)
	}

	return logs, nil
}
