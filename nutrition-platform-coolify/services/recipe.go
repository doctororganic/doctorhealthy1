package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Recipe represents a recipe
type Recipe struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PrepTime     int       `json:"prep_time"` // in minutes
	CookTime     int       `json:"cook_time"` // in minutes
	Servings     int       `json:"servings"`
	Difficulty   string    `json:"difficulty"` // easy, medium, hard
	Calories     int       `json:"calories"`
	Protein      float64   `json:"protein"`
	Carbs        float64   `json:"carbs"`
	Fat          float64   `json:"fat"`
	Fiber        float64   `json:"fiber"`
	Sugar        float64   `json:"sugar"`
	Sodium       float64   `json:"sodium"`
	Category     string    `json:"category"` // appetizer, main, dessert, beverage
	Cuisine      string    `json:"cuisine"`  // italian, chinese, arabic, etc.
	IsHalal      bool      `json:"is_halal"`
	IsVegetarian bool      `json:"is_vegetarian"`
	IsVegan      bool      `json:"is_vegan"`
	Tags         []string  `json:"tags"`
	ImageURL     string    `json:"image_url,omitempty"`
	Rating       float64   `json:"rating"`
	RatingCount  int       `json:"rating_count"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RecipeData represents the structure of recipes.json
type RecipeData struct {
	Recipes  []Recipe `json:"recipes"`
	Metadata Metadata `json:"metadata"`
}

const recipesFile = "backend/data/recipes.json"

// CreateRecipe creates a new recipe
func CreateRecipe(recipe *Recipe) error {
	// Generate ID and timestamps
	recipe.ID = uuid.New().String()
	recipe.CreatedAt = time.Now()
	recipe.UpdatedAt = time.Now()

	// Validate difficulty
	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
	}
	if !validDifficulties[recipe.Difficulty] {
		recipe.Difficulty = "medium" // default
	}

	// Validate category
	validCategories := map[string]bool{
		"appetizer": true,
		"main":      true,
		"dessert":   true,
		"beverage":  true,
		"snack":     true,
	}
	if !validCategories[recipe.Category] {
		recipe.Category = "main" // default
	}

	// Auto-detect dietary restrictions
	recipe.IsHalal = isHalalRecipe(recipe.Ingredients)
	recipe.IsVegetarian = isVegetarianRecipe(recipe.Ingredients)
	recipe.IsVegan = isVeganRecipe(recipe.Ingredients)

	// Initialize rating
	recipe.Rating = 0
	recipe.RatingCount = 0

	return AppendJSON(recipesFile, recipe)
}

// GetRecipesByUserID retrieves all recipes for a specific user
func GetRecipesByUserID(userID string) ([]Recipe, error) {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return nil, err
	}

	var userRecipes []Recipe
	for _, recipe := range data.Recipes {
		if recipe.UserID == userID {
			userRecipes = append(userRecipes, recipe)
		}
	}

	return userRecipes, nil
}

// GetPublicRecipes retrieves all public recipes
func GetPublicRecipes() ([]Recipe, error) {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return nil, err
	}

	var publicRecipes []Recipe
	for _, recipe := range data.Recipes {
		if recipe.IsPublic {
			publicRecipes = append(publicRecipes, recipe)
		}
	}

	return publicRecipes, nil
}

// GetRecipeByID retrieves a specific recipe by ID
func GetRecipeByID(recipeID string) (*Recipe, error) {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return nil, err
	}

	for _, recipe := range data.Recipes {
		if recipe.ID == recipeID {
			return &recipe, nil
		}
	}

	return nil, fmt.Errorf("recipe not found")
}

// UpdateRecipe updates an existing recipe
func UpdateRecipe(recipeID string, updatedRecipe *Recipe) error {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return err
	}

	for i, recipe := range data.Recipes {
		if recipe.ID == recipeID {
			// Preserve original ID, created time, and rating data
			updatedRecipe.ID = recipe.ID
			updatedRecipe.CreatedAt = recipe.CreatedAt
			updatedRecipe.UpdatedAt = time.Now()
			updatedRecipe.Rating = recipe.Rating
			updatedRecipe.RatingCount = recipe.RatingCount

			// Auto-detect dietary restrictions
			updatedRecipe.IsHalal = isHalalRecipe(updatedRecipe.Ingredients)
			updatedRecipe.IsVegetarian = isVegetarianRecipe(updatedRecipe.Ingredients)
			updatedRecipe.IsVegan = isVeganRecipe(updatedRecipe.Ingredients)

			data.Recipes[i] = *updatedRecipe
			data.Metadata.UpdatedAt = time.Now()

			return WriteJSON(recipesFile, data)
		}
	}

	return fmt.Errorf("recipe not found")
}

// DeleteRecipe deletes a recipe
func DeleteRecipe(recipeID string, userID string) error {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return err
	}

	for i, recipe := range data.Recipes {
		if recipe.ID == recipeID {
			// Check if user owns this recipe
			if recipe.UserID != userID {
				return fmt.Errorf("unauthorized: recipe belongs to another user")
			}

			// Remove recipe from slice
			data.Recipes = append(data.Recipes[:i], data.Recipes[i+1:]...)
			data.Metadata.UpdatedAt = time.Now()

			return WriteJSON(recipesFile, data)
		}
	}

	return fmt.Errorf("recipe not found")
}

// SearchRecipes searches recipes by various criteria
func SearchRecipes(query string, filters map[string]interface{}) ([]Recipe, error) {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return nil, err
	}

	var results []Recipe
	queryLower := strings.ToLower(query)

	for _, recipe := range data.Recipes {
		// Only search public recipes or user's own recipes
		userID, hasUserID := filters["user_id"].(string)
		if !recipe.IsPublic && (!hasUserID || recipe.UserID != userID) {
			continue
		}

		// Text search
		if query != "" {
			matchesQuery := false

			// Search in name
			if strings.Contains(strings.ToLower(recipe.Name), queryLower) {
				matchesQuery = true
			}

			// Search in description
			if !matchesQuery && strings.Contains(strings.ToLower(recipe.Description), queryLower) {
				matchesQuery = true
			}

			// Search in ingredients
			if !matchesQuery {
				for _, ingredient := range recipe.Ingredients {
					if strings.Contains(strings.ToLower(ingredient), queryLower) {
						matchesQuery = true
						break
					}
				}
			}

			// Search in tags
			if !matchesQuery {
				for _, tag := range recipe.Tags {
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
		if !matchesFilters(recipe, filters) {
			continue
		}

		results = append(results, recipe)
	}

	return results, nil
}

// matchesFilters checks if a recipe matches the given filters
func matchesFilters(recipe Recipe, filters map[string]interface{}) bool {
	if category, ok := filters["category"].(string); ok && category != "" {
		if recipe.Category != category {
			return false
		}
	}

	if cuisine, ok := filters["cuisine"].(string); ok && cuisine != "" {
		if recipe.Cuisine != cuisine {
			return false
		}
	}

	if difficulty, ok := filters["difficulty"].(string); ok && difficulty != "" {
		if recipe.Difficulty != difficulty {
			return false
		}
	}

	if isHalal, ok := filters["is_halal"].(bool); ok {
		if recipe.IsHalal != isHalal {
			return false
		}
	}

	if isVegetarian, ok := filters["is_vegetarian"].(bool); ok {
		if recipe.IsVegetarian != isVegetarian {
			return false
		}
	}

	if isVegan, ok := filters["is_vegan"].(bool); ok {
		if recipe.IsVegan != isVegan {
			return false
		}
	}

	if maxPrepTime, ok := filters["max_prep_time"].(float64); ok {
		if float64(recipe.PrepTime) > maxPrepTime {
			return false
		}
	}

	if maxCookTime, ok := filters["max_cook_time"].(float64); ok {
		if float64(recipe.CookTime) > maxCookTime {
			return false
		}
	}

	if minRating, ok := filters["min_rating"].(float64); ok {
		if recipe.Rating < minRating {
			return false
		}
	}

	return true
}

// RateRecipe adds or updates a rating for a recipe
func RateRecipe(recipeID string, rating float64) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return err
	}

	for i, recipe := range data.Recipes {
		if recipe.ID == recipeID {
			// Calculate new average rating
			totalRating := recipe.Rating * float64(recipe.RatingCount)
			totalRating += rating
			data.Recipes[i].RatingCount++
			data.Recipes[i].Rating = totalRating / float64(data.Recipes[i].RatingCount)
			data.Recipes[i].UpdatedAt = time.Now()
			data.Metadata.UpdatedAt = time.Now()

			return WriteJSON(recipesFile, data)
		}
	}

	return fmt.Errorf("recipe not found")
}

// GetRecipesByCategory retrieves recipes by category
func GetRecipesByCategory(category string) ([]Recipe, error) {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return nil, err
	}

	var categoryRecipes []Recipe
	for _, recipe := range data.Recipes {
		if recipe.IsPublic && recipe.Category == category {
			categoryRecipes = append(categoryRecipes, recipe)
		}
	}

	return categoryRecipes, nil
}

// GetTopRatedRecipes retrieves top-rated public recipes
func GetTopRatedRecipes(limit int) ([]Recipe, error) {
	var data RecipeData
	err := ReadJSON(recipesFile, &data)
	if err != nil {
		return nil, err
	}

	// Filter public recipes with ratings
	var ratedRecipes []Recipe
	for _, recipe := range data.Recipes {
		if recipe.IsPublic && recipe.RatingCount > 0 {
			ratedRecipes = append(ratedRecipes, recipe)
		}
	}

	// Sort by rating (simple bubble sort for small datasets)
	for i := 0; i < len(ratedRecipes)-1; i++ {
		for j := 0; j < len(ratedRecipes)-i-1; j++ {
			if ratedRecipes[j].Rating < ratedRecipes[j+1].Rating {
				ratedRecipes[j], ratedRecipes[j+1] = ratedRecipes[j+1], ratedRecipes[j]
			}
		}
	}

	// Return top recipes up to limit
	if limit > 0 && limit < len(ratedRecipes) {
		return ratedRecipes[:limit], nil
	}

	return ratedRecipes, nil
}

// Dietary restriction detection functions
func isHalalRecipe(ingredients []string) bool {
	nonHalalIngredients := []string{
		"pork", "ham", "bacon", "sausage", "pepperoni", "prosciutto",
		"alcohol", "wine", "beer", "rum", "vodka", "whiskey",
		"gelatin", "lard", "pancetta", "chorizo",
	}

	for _, ingredient := range ingredients {
		ingredientLower := strings.ToLower(ingredient)
		for _, nonHalal := range nonHalalIngredients {
			if strings.Contains(ingredientLower, nonHalal) {
				return false
			}
		}
	}

	return true
}

func isVegetarianRecipe(ingredients []string) bool {
	meatIngredients := []string{
		"beef", "chicken", "pork", "lamb", "turkey", "duck", "fish",
		"salmon", "tuna", "shrimp", "crab", "lobster", "meat", "ham",
		"bacon", "sausage", "pepperoni", "anchovy", "prosciutto",
	}

	for _, ingredient := range ingredients {
		ingredientLower := strings.ToLower(ingredient)
		for _, meat := range meatIngredients {
			if strings.Contains(ingredientLower, meat) {
				return false
			}
		}
	}

	return true
}

func isVeganRecipe(ingredients []string) bool {
	if !isVegetarianRecipe(ingredients) {
		return false
	}

	animalProducts := []string{
		"milk", "cheese", "butter", "cream", "yogurt", "egg", "honey",
		"gelatin", "whey", "casein", "lactose", "mayonnaise",
	}

	for _, ingredient := range ingredients {
		ingredientLower := strings.ToLower(ingredient)
		for _, animal := range animalProducts {
			if strings.Contains(ingredientLower, animal) {
				return false
			}
		}
	}

	return true
}
