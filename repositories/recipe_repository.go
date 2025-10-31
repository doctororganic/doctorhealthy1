package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"nutrition-platform/models"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

// CreateRecipe creates a new recipe with its ingredients
func (r *RecipeRepository) CreateRecipe(recipe *models.Recipe) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert recipe
	recipeQuery := `
		INSERT INTO recipes (user_id, name, description, instructions, prep_time_minutes, 
			cook_time_minutes, servings, difficulty, cuisine_type, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	err = tx.QueryRow(recipeQuery,
		recipe.UserID,
		recipe.Name,
		recipe.Description,
		recipe.Instructions,
		recipe.PrepTimeMinutes,
		recipe.CookTimeMinutes,
		recipe.Servings,
		recipe.Difficulty,
		recipe.CuisineType,
		recipe.IsPublic,
		time.Now(),
		time.Now(),
	).Scan(&recipe.ID)

	if err != nil {
		log.Printf("Error creating recipe: %v", err)
		return fmt.Errorf("failed to create recipe: %w", err)
	}

	// Insert recipe ingredients
	for _, ingredient := range recipe.Ingredients {
		ingredientQuery := `
			INSERT INTO recipe_ingredients (recipe_id, food_id, quantity, unit, notes)
			VALUES ($1, $2, $3, $4, $5)`

		_, err := tx.Exec(ingredientQuery,
			recipe.ID,
			ingredient.FoodID,
			ingredient.Quantity,
			ingredient.Unit,
			ingredient.Notes,
		)

		if err != nil {
			log.Printf("Error creating recipe ingredient: %v", err)
			return fmt.Errorf("failed to create recipe ingredient: %w", err)
		}
	}

	// Insert nutrition summary
	nutritionQuery := `
		INSERT INTO recipe_nutrition (recipe_id, calories, protein, carbs, fat, 
			saturated_fat, fiber, sugar, sodium, cholesterol, potassium)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(nutritionQuery,
		recipe.ID,
		recipe.NutritionSummary.Calories,
		recipe.NutritionSummary.Protein,
		recipe.NutritionSummary.Carbs,
		recipe.NutritionSummary.Fat,
		recipe.NutritionSummary.SaturatedFat,
		recipe.NutritionSummary.Fiber,
		recipe.NutritionSummary.Sugar,
		recipe.NutritionSummary.Sodium,
		recipe.NutritionSummary.Cholesterol,
		recipe.NutritionSummary.Potassium,
	)

	if err != nil {
		log.Printf("Error creating recipe nutrition: %v", err)
		return fmt.Errorf("failed to create recipe nutrition: %w", err)
	}

	return tx.Commit()
}

// GetRecipeByID retrieves a recipe by its ID with all details
func (r *RecipeRepository) GetRecipeByID(id, userID string) (*models.Recipe, error) {
	// Get recipe details
	recipeQuery := `
		SELECT id, user_id, name, description, instructions, prep_time_minutes,
			cook_time_minutes, servings, difficulty, cuisine_type, is_public, created_at, updated_at
		FROM recipes 
		WHERE id = $1 AND (user_id = $2 OR is_public = true)`

	var recipe models.Recipe
	var cuisineType sql.NullString

	err := r.db.QueryRow(recipeQuery, id, userID).Scan(
		&recipe.ID,
		&recipe.UserID,
		&recipe.Name,
		&recipe.Description,
		&recipe.Instructions,
		&recipe.PrepTimeMinutes,
		&recipe.CookTimeMinutes,
		&recipe.Servings,
		&recipe.Difficulty,
		&cuisineType,
		&recipe.IsPublic,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("recipe not found")
		}
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if cuisineType.Valid {
		recipe.CuisineType = cuisineType.String
	}

	// Get recipe ingredients
	ingredientsQuery := `
		SELECT food_id, quantity, unit, notes
		FROM recipe_ingredients
		WHERE recipe_id = $1`

	rows, err := r.db.Query(ingredientsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []*models.RecipeIngredient
	for rows.Next() {
		var ingredient models.RecipeIngredient
		var notes sql.NullString

		err := rows.Scan(
			&ingredient.FoodID,
			&ingredient.Quantity,
			&ingredient.Unit,
			&notes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan ingredient row: %w", err)
		}

		if notes.Valid {
			ingredient.Notes = notes.String
		}

		ingredients = append(ingredients, &ingredient)
	}
	recipe.Ingredients = ingredients

	// Get nutrition summary
	nutritionQuery := `
		SELECT calories, protein, carbs, fat, saturated_fat, fiber, sugar, sodium, cholesterol, potassium
		FROM recipe_nutrition
		WHERE recipe_id = $1`

	var nutrition models.NutritionInfo
	err = r.db.QueryRow(nutritionQuery, id).Scan(
		&nutrition.Calories,
		&nutrition.Protein,
		&nutrition.Carbs,
		&nutrition.Fat,
		&nutrition.SaturatedFat,
		&nutrition.Fiber,
		&nutrition.Sugar,
		&nutrition.Sodium,
		&nutrition.Cholesterol,
		&nutrition.Potassium,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get recipe nutrition: %w", err)
		}
		// If no nutrition data, set defaults
	}
	recipe.NutritionSummary = nutrition

	return &recipe, nil
}

// SearchRecipes searches for recipes based on query and filters
func (r *RecipeRepository) SearchRecipes(userID, query string, filters models.RecipeSearchFilters, limit, offset int) ([]*models.Recipe, error) {
	whereClauses := []string{"(user_id = $1 OR is_public = true)"}
	args := []interface{}{userID}
	argIndex := 2

	if query != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, "%"+query+"%", "%"+query+"%")
		argIndex += 2
	}

	if filters.CuisineType != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("cuisine_type = $%d", argIndex))
		args = append(args, filters.CuisineType)
		argIndex++
	}

	if filters.Difficulty != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("difficulty = $%d", argIndex))
		args = append(args, filters.Difficulty)
		argIndex++
	}

	if filters.MaxPrepTime != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("prep_time_minutes <= $%d", argIndex))
		args = append(args, *filters.MaxPrepTime)
		argIndex++
	}

	if filters.IsPublic != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_public = $%d", argIndex))
		args = append(args, *filters.IsPublic)
		argIndex++
	}

	baseQuery := `
		SELECT id, user_id, name, description, instructions, prep_time_minutes,
			cook_time_minutes, servings, difficulty, cuisine_type, is_public, created_at, updated_at
		FROM recipes`

	whereClause := " WHERE " + strings.Join(whereClauses, " AND ")
	
	// Add ordering
	orderBy := " ORDER BY created_at DESC"
	if filters.SortBy != "" {
		direction := "DESC"
		if filters.SortDirection == "asc" {
			direction = "ASC"
		}
		orderBy = fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, direction)
	}

	// Add pagination
	pagination := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	fullQuery := baseQuery + whereClause + orderBy + pagination

	rows, err := r.db.Query(fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var cuisineType sql.NullString

		err := rows.Scan(
			&recipe.ID,
			&recipe.UserID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Instructions,
			&recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes,
			&recipe.Servings,
			&recipe.Difficulty,
			&cuisineType,
			&recipe.IsPublic,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe row: %w", err)
		}

		if cuisineType.Valid {
			recipe.CuisineType = cuisineType.String
		}

		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

// UpdateRecipe updates an existing recipe and its ingredients
func (r *RecipeRepository) UpdateRecipe(recipe *models.Recipe) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update recipe
	recipeQuery := `
		UPDATE recipes 
		SET name = $2, description = $3, instructions = $4, prep_time_minutes = $5,
			cook_time_minutes = $6, servings = $7, difficulty = $8, cuisine_type = $9,
			is_public = $10, updated_at = $11
		WHERE id = $1 AND user_id = $12`

	_, err = tx.Exec(recipeQuery,
		recipe.ID,
		recipe.Name,
		recipe.Description,
		recipe.Instructions,
		recipe.PrepTimeMinutes,
		recipe.CookTimeMinutes,
		recipe.Servings,
		recipe.Difficulty,
		recipe.CuisineType,
		recipe.IsPublic,
		time.Now(),
		recipe.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}

	// Delete existing ingredients
	_, err = tx.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = $1", recipe.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing ingredients: %w", err)
	}

	// Insert new ingredients
	for _, ingredient := range recipe.Ingredients {
		ingredientQuery := `
			INSERT INTO recipe_ingredients (recipe_id, food_id, quantity, unit, notes)
			VALUES ($1, $2, $3, $4, $5)`

		_, err := tx.Exec(ingredientQuery,
			recipe.ID,
			ingredient.FoodID,
			ingredient.Quantity,
			ingredient.Unit,
			ingredient.Notes,
		)

		if err != nil {
			return fmt.Errorf("failed to create recipe ingredient: %w", err)
		}
	}

	// Update nutrition summary
	nutritionQuery := `
		UPDATE recipe_nutrition 
		SET calories = $2, protein = $3, carbs = $4, fat = $5, saturated_fat = $6,
			fiber = $7, sugar = $8, sodium = $9, cholesterol = $10, potassium = $11
		WHERE recipe_id = $1`

	_, err = tx.Exec(nutritionQuery,
		recipe.ID,
		recipe.NutritionSummary.Calories,
		recipe.NutritionSummary.Protein,
		recipe.NutritionSummary.Carbs,
		recipe.NutritionSummary.Fat,
		recipe.NutritionSummary.SaturatedFat,
		recipe.NutritionSummary.Fiber,
		recipe.NutritionSummary.Sugar,
		recipe.NutritionSummary.Sodium,
		recipe.NutritionSummary.Cholesterol,
		recipe.NutritionSummary.Potassium,
	)

	if err != nil {
		return fmt.Errorf("failed to update recipe nutrition: %w", err)
	}

	return tx.Commit()
}

// DeleteRecipe soft deletes a recipe by marking it as deleted
func (r *RecipeRepository) DeleteRecipe(id, userID string) error {
	query := `UPDATE recipes SET updated_at = $1 WHERE id = $2 AND user_id = $3`
	
	_, err := r.db.Exec(query, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	return nil
}

// GetUserRecipes retrieves all recipes created by a user
func (r *RecipeRepository) GetUserRecipes(userID string, limit, offset int) ([]*models.Recipe, error) {
	query := `
		SELECT id, user_id, name, description, instructions, prep_time_minutes,
			cook_time_minutes, servings, difficulty, cuisine_type, is_public, created_at, updated_at
		FROM recipes 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var cuisineType sql.NullString

		err := rows.Scan(
			&recipe.ID,
			&recipe.UserID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Instructions,
			&recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes,
			&recipe.Servings,
			&recipe.Difficulty,
			&cuisineType,
			&recipe.IsPublic,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe row: %w", err)
		}

		if cuisineType.Valid {
			recipe.CuisineType = cuisineType.String
		}

		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

// GetPublicRecipes retrieves all public recipes
func (r *RecipeRepository) GetPublicRecipes(limit, offset int) ([]*models.Recipe, error) {
	query := `
		SELECT id, user_id, name, description, instructions, prep_time_minutes,
			cook_time_minutes, servings, difficulty, cuisine_type, is_public, created_at, updated_at
		FROM recipes 
		WHERE is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var cuisineType sql.NullString

		err := rows.Scan(
			&recipe.ID,
			&recipe.UserID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Instructions,
			&recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes,
			&recipe.Servings,
			&recipe.Difficulty,
			&cuisineType,
			&recipe.IsPublic,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe row: %w", err)
		}

		if cuisineType.Valid {
			recipe.CuisineType = cuisineType.String
		}

		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}
