package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"nutrition-platform/models"
)

// RecipeService handles recipe operations
type RecipeService struct {
	db *sql.DB
}

// NewRecipeService creates a new recipe service
func NewRecipeService(db *sql.DB) *RecipeService {
	return &RecipeService{db: db}
}

// CreateRecipe creates a new recipe
func (s *RecipeService) CreateRecipe(userID string, req *models.CreateRecipeRequest) (*models.Recipe, error) {
	recipe := &models.Recipe{
		ID:                  generateUUID(),
		Name:                req.Name,
		NameAr:              req.NameAr,
		Description:         req.Description,
		DescriptionAr:       req.DescriptionAr,
		Cuisine:             req.Cuisine,
		Country:             req.Country,
		DifficultyLevel:     req.DifficultyLevel,
		PrepTimeMinutes:     req.PrepTimeMinutes,
		CookTimeMinutes:     req.CookTimeMinutes,
		Servings:            req.Servings,
		Ingredients:         req.Ingredients,
		Instructions:        req.Instructions,
		NutritionPerServing: req.NutritionPerServing,
		DietaryTags:         req.DietaryTags,
		Allergens:           req.Allergens,
		IsHalal:             req.IsHalal,
		IsKosher:            req.IsKosher,
		ImageURL:            req.ImageURL,
		VideoURL:            req.VideoURL,
		CreatedBy:           &userID,
		Rating:              0,
		RatingCount:         0,
		Verified:            false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Calculate total time
	if req.PrepTimeMinutes != nil && req.CookTimeMinutes != nil {
		totalTime := *req.PrepTimeMinutes + *req.CookTimeMinutes
		recipe.TotalTimeMinutes = &totalTime
	}

	// Store in database
	if err := s.storeRecipe(recipe); err != nil {
		return nil, fmt.Errorf("failed to store recipe: %w", err)
	}

	return recipe, nil
}

// GetRecipe retrieves a recipe by ID
func (s *RecipeService) GetRecipe(id string) (*models.Recipe, error) {
	query := `
		SELECT id, name, name_ar, description, description_ar, cuisine, country, 
		       difficulty_level, prep_time_minutes, cook_time_minutes, total_time_minutes,
		       servings, ingredients, instructions, nutrition_per_serving, dietary_tags,
		       allergens, is_halal, is_kosher, image_url, video_url, rating, rating_count,
		       created_by, verified, created_at, updated_at
		FROM recipes WHERE id = $1
	`

	var recipe models.Recipe
	var ingredientsJSON, instructionsJSON, nutritionJSON, dietaryTagsJSON, allergensJSON []byte

	err := s.db.QueryRow(query, id).Scan(
		&recipe.ID, &recipe.Name, &recipe.NameAr, &recipe.Description, &recipe.DescriptionAr,
		&recipe.Cuisine, &recipe.Country, &recipe.DifficultyLevel, &recipe.PrepTimeMinutes,
		&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings,
		&ingredientsJSON, &instructionsJSON, &nutritionJSON, &dietaryTagsJSON,
		&allergensJSON, &recipe.IsHalal, &recipe.IsKosher, &recipe.ImageURL,
		&recipe.VideoURL, &recipe.Rating, &recipe.RatingCount, &recipe.CreatedBy,
		&recipe.Verified, &recipe.CreatedAt, &recipe.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if err := json.Unmarshal(ingredientsJSON, &recipe.Ingredients); err != nil {
		return nil, fmt.Errorf("failed to parse ingredients: %w", err)
	}
	if err := json.Unmarshal(instructionsJSON, &recipe.Instructions); err != nil {
		return nil, fmt.Errorf("failed to parse instructions: %w", err)
	}
	if len(nutritionJSON) > 0 {
		if err := json.Unmarshal(nutritionJSON, &recipe.NutritionPerServing); err != nil {
			return nil, fmt.Errorf("failed to parse nutrition: %w", err)
		}
	}
	if err := json.Unmarshal(dietaryTagsJSON, &recipe.DietaryTags); err != nil {
		return nil, fmt.Errorf("failed to parse dietary tags: %w", err)
	}
	if err := json.Unmarshal(allergensJSON, &recipe.Allergens); err != nil {
		return nil, fmt.Errorf("failed to parse allergens: %w", err)
	}

	return &recipe, nil
}

// SearchRecipes searches for recipes based on criteria
func (s *RecipeService) SearchRecipes(req *models.RecipeSearchRequest) (*models.RecipeListResponse, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if req.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+req.Query+"%")
		argIndex++
	}

	if req.Cuisine != "" {
		conditions = append(conditions, fmt.Sprintf("cuisine = $%d", argIndex))
		args = append(args, req.Cuisine)
		argIndex++
	}

	if req.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, req.Country)
		argIndex++
	}

	if req.DifficultyLevel != "" {
		conditions = append(conditions, fmt.Sprintf("difficulty_level = $%d", argIndex))
		args = append(args, req.DifficultyLevel)
		argIndex++
	}

	if req.MaxPrepTime > 0 {
		conditions = append(conditions, fmt.Sprintf("prep_time_minutes <= $%d", argIndex))
		args = append(args, req.MaxPrepTime)
		argIndex++
	}

	if req.MaxCookTime > 0 {
		conditions = append(conditions, fmt.Sprintf("cook_time_minutes <= $%d", argIndex))
		args = append(args, req.MaxCookTime)
		argIndex++
	}

	if req.IsHalal != nil {
		conditions = append(conditions, fmt.Sprintf("is_halal = $%d", argIndex))
		args = append(args, *req.IsHalal)
		argIndex++
	}

	if req.IsKosher != nil {
		conditions = append(conditions, fmt.Sprintf("is_kosher = $%d", argIndex))
		args = append(args, *req.IsKosher)
		argIndex++
	}

	if req.MinRating > 0 {
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", argIndex))
		args = append(args, req.MinRating)
		argIndex++
	}

	// Handle dietary tags
	if len(req.DietaryTags) > 0 {
		for _, tag := range req.DietaryTags {
			conditions = append(conditions, fmt.Sprintf("dietary_tags ? $%d", argIndex))
			args = append(args, tag)
			argIndex++
		}
	}

	// Handle allergen exclusions
	if len(req.Allergens) > 0 {
		for _, allergen := range req.Allergens {
			conditions = append(conditions, fmt.Sprintf("NOT (allergens ? $%d)", argIndex))
			args = append(args, allergen)
			argIndex++
		}
	}

	// Build query
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total results
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM recipes %s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count recipes: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, name, name_ar, description, description_ar, cuisine, country,
		       difficulty_level, prep_time_minutes, cook_time_minutes, total_time_minutes,
		       servings, ingredients, instructions, nutrition_per_serving, dietary_tags,
		       allergens, is_halal, is_kosher, image_url, video_url, rating, rating_count,
		       created_by, verified, created_at, updated_at
		FROM recipes %s
		ORDER BY rating DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var ingredientsJSON, instructionsJSON, nutritionJSON, dietaryTagsJSON, allergensJSON []byte

		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.NameAr, &recipe.Description, &recipe.DescriptionAr,
			&recipe.Cuisine, &recipe.Country, &recipe.DifficultyLevel, &recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings,
			&ingredientsJSON, &instructionsJSON, &nutritionJSON, &dietaryTagsJSON,
			&allergensJSON, &recipe.IsHalal, &recipe.IsKosher, &recipe.ImageURL,
			&recipe.VideoURL, &recipe.Rating, &recipe.RatingCount, &recipe.CreatedBy,
			&recipe.Verified, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(ingredientsJSON, &recipe.Ingredients)
		json.Unmarshal(instructionsJSON, &recipe.Instructions)
		if len(nutritionJSON) > 0 {
			json.Unmarshal(nutritionJSON, &recipe.NutritionPerServing)
		}
		json.Unmarshal(dietaryTagsJSON, &recipe.DietaryTags)
		json.Unmarshal(allergensJSON, &recipe.Allergens)

		recipes = append(recipes, recipe)
	}

	return &models.RecipeListResponse{
		Recipes: recipes,
		Total:   total,
		Page:    req.Page,
		Limit:   req.Limit,
		HasNext: offset+req.Limit < total,
	}, nil
}

// UpdateRecipe updates a recipe
func (s *RecipeService) UpdateRecipe(id, userID string, req *models.UpdateRecipeRequest) (*models.Recipe, error) {
	// First check if recipe exists and user has permission
	recipe, err := s.GetRecipe(id)
	if err != nil {
		return nil, err
	}

	if recipe.CreatedBy == nil || *recipe.CreatedBy != userID {
		return nil, fmt.Errorf("permission denied")
	}

	// Build update query dynamically
	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	// Add other fields as needed...

	if len(setParts) == 0 {
		return recipe, nil // No updates
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf("UPDATE recipes SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return s.GetRecipe(id)
}

// DeleteRecipe deletes a recipe
func (s *RecipeService) DeleteRecipe(id, userID string) error {
	// Check permission
	recipe, err := s.GetRecipe(id)
	if err != nil {
		return err
	}

	if recipe.CreatedBy == nil || *recipe.CreatedBy != userID {
		return fmt.Errorf("permission denied")
	}

	query := "DELETE FROM recipes WHERE id = $1"
	_, err = s.db.Exec(query, id)
	return err
}

// RateRecipe adds or updates a recipe rating
func (s *RecipeService) RateRecipe(recipeID, userID string, req *models.RecipeRatingRequest) error {
	// Implementation would include rating logic
	// For now, simplified version
	query := `
		UPDATE recipes 
		SET rating = (rating * rating_count + $1) / (rating_count + 1),
		    rating_count = rating_count + 1,
		    updated_at = $2
		WHERE id = $3
	`

	_, err := s.db.Exec(query, req.Rating, time.Now(), recipeID)
	return err
}

// GetRecipesByCountry retrieves recipes by country
func (s *RecipeService) GetRecipesByCountry(country string, page, limit int) (*models.RecipeListResponse, error) {
	req := &models.RecipeSearchRequest{
		Country: country,
		Page:    page,
		Limit:   limit,
	}
	return s.SearchRecipes(req)
}

// GetRecipesByCuisine retrieves recipes by cuisine
func (s *RecipeService) GetRecipesByCuisine(cuisine string, page, limit int) (*models.RecipeListResponse, error) {
	req := &models.RecipeSearchRequest{
		Cuisine: cuisine,
		Page:    page,
		Limit:   limit,
	}
	return s.SearchRecipes(req)
}

// Helper methods

func (s *RecipeService) storeRecipe(recipe *models.Recipe) error {
	ingredientsJSON, _ := json.Marshal(recipe.Ingredients)
	instructionsJSON, _ := json.Marshal(recipe.Instructions)
	nutritionJSON, _ := json.Marshal(recipe.NutritionPerServing)
	dietaryTagsJSON, _ := json.Marshal(recipe.DietaryTags)
	allergensJSON, _ := json.Marshal(recipe.Allergens)

	query := `
		INSERT INTO recipes (
			id, name, name_ar, description, description_ar, cuisine, country,
			difficulty_level, prep_time_minutes, cook_time_minutes, total_time_minutes,
			servings, ingredients, instructions, nutrition_per_serving, dietary_tags,
			allergens, is_halal, is_kosher, image_url, video_url, rating, rating_count,
			created_by, verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err := s.db.Exec(query,
		recipe.ID, recipe.Name, recipe.NameAr, recipe.Description, recipe.DescriptionAr,
		recipe.Cuisine, recipe.Country, recipe.DifficultyLevel, recipe.PrepTimeMinutes,
		recipe.CookTimeMinutes, recipe.TotalTimeMinutes, recipe.Servings,
		ingredientsJSON, instructionsJSON, nutritionJSON, dietaryTagsJSON,
		allergensJSON, recipe.IsHalal, recipe.IsKosher, recipe.ImageURL,
		recipe.VideoURL, recipe.Rating, recipe.RatingCount, recipe.CreatedBy,
		recipe.Verified, recipe.CreatedAt, recipe.UpdatedAt,
	)

	return err
}

func generateUUID() string {
	// Simple UUID generation - in production use proper UUID library
	return fmt.Sprintf("recipe_%d", time.Now().UnixNano())
}
