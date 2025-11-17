package services

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// RecipeService handles recipe-related business logic
type RecipeService struct {
	db *sql.DB
}

// NewRecipeService creates a new recipe service
func NewRecipeService(db *sql.DB) *RecipeService {
	return &RecipeService{
		db: db,
	}
}

// GetRecipes retrieves recipes (public access)
func (s *RecipeService) GetRecipes(limit, offset int, filters map[string]interface{}) ([]models.Recipe, error) {
	query := `
		SELECT id, name, description, instructions, prep_time, cook_time,
			   servings, calories, protein, carbs, fat, is_public,
			   created_by, created_at, updated_at
		FROM recipes 
		WHERE is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var createdBy sql.NullInt64
		
		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Instructions,
			&recipe.PrepTime, &recipe.CookTime, &recipe.Servings,
			&recipe.Calories, &recipe.Protein, &recipe.Carbs, &recipe.Fat,
			&recipe.IsPublic, &createdBy, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if createdBy.Valid {
			recipe.CreatedBy = uint(createdBy.Int64)
		}
		
		recipes = append(recipes, recipe)
	}
	
	return recipes, nil
}

// CreateRecipe creates a new recipe
func (s *RecipeService) CreateRecipe(recipe *models.Recipe) error {
	query := `
		INSERT INTO recipes (name, description, instructions, prep_time, cook_time,
							servings, calories, protein, carbs, fat, is_public,
							created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		recipe.Name, recipe.Description, recipe.Instructions,
		recipe.PrepTime, recipe.CookTime, recipe.Servings,
		recipe.Calories, recipe.Protein, recipe.Carbs, recipe.Fat,
		recipe.IsPublic, recipe.CreatedBy, now, now,
	).Scan(&recipe.ID)
	
	if err != nil {
		return err
	}
	
	recipe.CreatedAt = now
	recipe.UpdatedAt = now
	
	return nil
}

// GetRecipe retrieves a recipe by ID
func (s *RecipeService) GetRecipe(id uint) (*models.Recipe, error) {
	query := `
		SELECT id, name, description, instructions, prep_time, cook_time,
			   servings, calories, protein, carbs, fat, is_public,
			   created_by, created_at, updated_at
		FROM recipes 
		WHERE id = $1
	`
	
	var recipe models.Recipe
	var createdBy sql.NullInt64
	
	err := s.db.QueryRow(query, id).Scan(
		&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Instructions,
		&recipe.PrepTime, &recipe.CookTime, &recipe.Servings,
		&recipe.Calories, &recipe.Protein, &recipe.Carbs, &recipe.Fat,
		&recipe.IsPublic, &createdBy, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("recipe not found")
		}
		return nil, err
	}
	
	if createdBy.Valid {
		recipe.CreatedBy = uint(createdBy.Int64)
	}
	
	return &recipe, nil
}

// UpdateRecipe updates an existing recipe
func (s *RecipeService) UpdateRecipe(recipe *models.Recipe) error {
	query := `
		UPDATE recipes 
		SET name = $2, description = $3, instructions = $4, prep_time = $5,
			cook_time = $6, servings = $7, calories = $8, protein = $9,
			carbs = $10, fat = $11, is_public = $12, updated_at = $13
		WHERE id = $1
	`
	
	recipe.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(query,
		recipe.ID, recipe.Name, recipe.Description, recipe.Instructions,
		recipe.PrepTime, recipe.CookTime, recipe.Servings,
		recipe.Calories, recipe.Protein, recipe.Carbs, recipe.Fat,
		recipe.IsPublic, recipe.UpdatedAt,
	)
	
	return err
}

// DeleteRecipe deletes a recipe
func (s *RecipeService) DeleteRecipe(id uint) error {
	query := `DELETE FROM recipes WHERE id = $1`
	
	_, err := s.db.Exec(query, id)
	return err
}

// SearchRecipes searches for recipes by name or description
func (s *RecipeService) SearchRecipes(query string, limit, offset int) ([]models.Recipe, error) {
	searchQuery := `
		SELECT id, name, description, instructions, prep_time, cook_time,
			   servings, calories, protein, carbs, fat, is_public,
			   created_by, created_at, updated_at
		FROM recipes 
		WHERE is_public = true AND (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.Query(searchQuery, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var createdBy sql.NullInt64
		
		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Instructions,
			&recipe.PrepTime, &recipe.CookTime, &recipe.Servings,
			&recipe.Calories, &recipe.Protein, &recipe.Carbs, &recipe.Fat,
			&recipe.IsPublic, &createdBy, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if createdBy.Valid {
			recipe.CreatedBy = uint(createdBy.Int64)
		}
		
		recipes = append(recipes, recipe)
	}
	
	return recipes, nil
}
