package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Product represents a nutrition product
type Product struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Brand        string    `json:"brand"`
	Description  string    `json:"description"`
	Category     string    `json:"category"` // supplement, food, beverage, equipment
	Subcategory  string    `json:"subcategory"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	ImageURL     string    `json:"image_url,omitempty"`
	Images       []string  `json:"images,omitempty"`
	Nutrition    Nutrition `json:"nutrition,omitempty"`
	Ingredients  []string  `json:"ingredients"`
	Allergens    []string  `json:"allergens"`
	IsHalal      bool      `json:"is_halal"`
	IsVegetarian bool      `json:"is_vegetarian"`
	IsVegan      bool      `json:"is_vegan"`
	IsOrganic    bool      `json:"is_organic"`
	Tags         []string  `json:"tags"`
	Rating       float64   `json:"rating"`
	RatingCount  int       `json:"rating_count"`
	Stock        int       `json:"stock"`
	SKU          string    `json:"sku"`
	Barcode      string    `json:"barcode,omitempty"`
	Manufacturer string    `json:"manufacturer"`
	Country      string    `json:"country"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	IsActive     bool      `json:"is_active"`
	SubmittedBy  string    `json:"submitted_by,omitempty"`
	ApprovedBy   string    `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PendingProduct represents a product awaiting approval
type PendingProduct struct {
	ID          string    `json:"id"`
	Product     Product   `json:"product"`
	SubmittedBy string    `json:"submitted_by"`
	Reason      string    `json:"reason,omitempty"`
	Status      string    `json:"status"` // pending, approved, rejected
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Nutrition represents nutritional information
type Nutrition struct {
	ServingSize   string  `json:"serving_size"`
	Calories      int     `json:"calories"`
	Protein       float64 `json:"protein"`
	Carbs         float64 `json:"carbs"`
	Fat           float64 `json:"fat"`
	Fiber         float64 `json:"fiber"`
	Sugar         float64 `json:"sugar"`
	Sodium        float64 `json:"sodium"`
	Cholesterol   float64 `json:"cholesterol"`
	VitaminA      float64 `json:"vitamin_a,omitempty"`
	VitaminC      float64 `json:"vitamin_c,omitempty"`
	Calcium       float64 `json:"calcium,omitempty"`
	Iron          float64 `json:"iron,omitempty"`
}

// ProductData represents the structure of products.json
type ProductData struct {
	Products []Product `json:"products"`
	Metadata Metadata  `json:"metadata"`
}

// PendingProductData represents the structure of pending-products.json
type PendingProductData struct {
	PendingProducts []PendingProduct `json:"pending_products"`
	Metadata        Metadata         `json:"metadata"`
}

const (
	productsFile        = "backend/data/products.json"
	pendingProductsFile = "backend/data/pending-products.json"
)

// SubmitProduct submits a product for approval
func SubmitProduct(product *Product, submittedBy string) error {
	// Generate ID and timestamps
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.SubmittedBy = submittedBy
	product.IsActive = false // Not active until approved

	// Auto-detect dietary restrictions
	product.IsHalal = isHalalProduct(product.Ingredients)
	product.IsVegetarian = isVegetarianProduct(product.Ingredients)
	product.IsVegan = isVeganProduct(product.Ingredients)

	// Generate SKU if not provided
	if product.SKU == "" {
		product.SKU = generateSKU(product.Name, product.Brand)
	}

	// Create pending product entry
	pendingProduct := PendingProduct{
		ID:          uuid.New().String(),
		Product:     *product,
		SubmittedBy: submittedBy,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return AppendJSON(pendingProductsFile, pendingProduct)
}

// GetPendingProducts retrieves all pending products (admin only)
func GetPendingProducts() ([]PendingProduct, error) {
	var data PendingProductData
	err := ReadJSON(pendingProductsFile, &data)
	if err != nil {
		return nil, err
	}

	var pendingProducts []PendingProduct
	for _, product := range data.PendingProducts {
		if product.Status == "pending" {
			pendingProducts = append(pendingProducts, product)
		}
	}

	return pendingProducts, nil
}

// ApproveProduct approves a pending product
func ApproveProduct(pendingID string, approvedBy string, notes string) error {
	// Get pending product
	var pendingData PendingProductData
	err := ReadJSON(pendingProductsFile, &pendingData)
	if err != nil {
		return err
	}

	var pendingProduct *PendingProduct
	var pendingIndex int
	for i, p := range pendingData.PendingProducts {
		if p.ID == pendingID && p.Status == "pending" {
			pendingProduct = &p
			pendingIndex = i
			break
		}
	}

	if pendingProduct == nil {
		return fmt.Errorf("pending product not found")
	}

	// Update product for approval
	product := pendingProduct.Product
	product.IsActive = true
	product.ApprovedBy = approvedBy
	now := time.Now()
	product.ApprovedAt = &now
	product.UpdatedAt = now

	// Add to products
	err = AppendJSON(productsFile, product)
	if err != nil {
		return err
	}

	// Update pending product status
	pendingData.PendingProducts[pendingIndex].Status = "approved"
	pendingData.PendingProducts[pendingIndex].Notes = notes
	pendingData.PendingProducts[pendingIndex].UpdatedAt = now
	pendingData.Metadata.UpdatedAt = now

	return WriteJSON(pendingProductsFile, pendingData)
}

// RejectProduct rejects a pending product
func RejectProduct(pendingID string, rejectedBy string, reason string) error {
	var data PendingProductData
	err := ReadJSON(pendingProductsFile, &data)
	if err != nil {
		return err
	}

	for i, product := range data.PendingProducts {
		if product.ID == pendingID && product.Status == "pending" {
			data.PendingProducts[i].Status = "rejected"
			data.PendingProducts[i].Reason = reason
			data.PendingProducts[i].Notes = fmt.Sprintf("Rejected by %s", rejectedBy)
			data.PendingProducts[i].UpdatedAt = time.Now()
			data.Metadata.UpdatedAt = time.Now()
			
			return WriteJSON(pendingProductsFile, data)
		}
	}

	return fmt.Errorf("pending product not found")
}

// GetProducts retrieves all active products
func GetProducts() ([]Product, error) {
	var data ProductData
	err := ReadJSON(productsFile, &data)
	if err != nil {
		return nil, err
	}

	var activeProducts []Product
	for _, product := range data.Products {
		if product.IsActive {
			activeProducts = append(activeProducts, product)
		}
	}

	return activeProducts, nil
}

// GetProductByID retrieves a specific product by ID
func GetProductByID(productID string) (*Product, error) {
	var data ProductData
	err := ReadJSON(productsFile, &data)
	if err != nil {
		return nil, err
	}

	for _, product := range data.Products {
		if product.ID == productID && product.IsActive {
			return &product, nil
		}
	}

	return nil, fmt.Errorf("product not found")
}

// SearchProducts searches products by various criteria
func SearchProducts(query string, filters map[string]interface{}) ([]Product, error) {
	var data ProductData
	err := ReadJSON(productsFile, &data)
	if err != nil {
		return nil, err
	}

	var results []Product
	queryLower := strings.ToLower(query)

	for _, product := range data.Products {
		if !product.IsActive {
			continue
		}

		// Text search
		if query != "" {
			matchesQuery := false
			
			// Search in name
			if strings.Contains(strings.ToLower(product.Name), queryLower) {
				matchesQuery = true
			}
			
			// Search in brand
			if !matchesQuery && strings.Contains(strings.ToLower(product.Brand), queryLower) {
				matchesQuery = true
			}
			
			// Search in description
			if !matchesQuery && strings.Contains(strings.ToLower(product.Description), queryLower) {
				matchesQuery = true
			}
			
			// Search in ingredients
			if !matchesQuery {
				for _, ingredient := range product.Ingredients {
					if strings.Contains(strings.ToLower(ingredient), queryLower) {
						matchesQuery = true
						break
					}
				}
			}
			
			// Search in tags
			if !matchesQuery {
				for _, tag := range product.Tags {
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
		if !matchesProductFilters(product, filters) {
			continue
		}

		results = append(results, product)
	}

	return results, nil
}

// matchesProductFilters checks if a product matches the given filters
func matchesProductFilters(product Product, filters map[string]interface{}) bool {
	if category, ok := filters["category"].(string); ok && category != "" {
		if product.Category != category {
			return false
		}
	}

	if subcategory, ok := filters["subcategory"].(string); ok && subcategory != "" {
		if product.Subcategory != subcategory {
			return false
		}
	}

	if brand, ok := filters["brand"].(string); ok && brand != "" {
		if product.Brand != brand {
			return false
		}
	}

	if isHalal, ok := filters["is_halal"].(bool); ok {
		if product.IsHalal != isHalal {
			return false
		}
	}

	if isVegetarian, ok := filters["is_vegetarian"].(bool); ok {
		if product.IsVegetarian != isVegetarian {
			return false
		}
	}

	if isVegan, ok := filters["is_vegan"].(bool); ok {
		if product.IsVegan != isVegan {
			return false
		}
	}

	if isOrganic, ok := filters["is_organic"].(bool); ok {
		if product.IsOrganic != isOrganic {
			return false
		}
	}

	if minPrice, ok := filters["min_price"].(float64); ok {
		if product.Price < minPrice {
			return false
		}
	}

	if maxPrice, ok := filters["max_price"].(float64); ok {
		if product.Price > maxPrice {
			return false
		}
	}

	if minRating, ok := filters["min_rating"].(float64); ok {
		if product.Rating < minRating {
			return false
		}
	}

	if inStock, ok := filters["in_stock"].(bool); ok && inStock {
		if product.Stock <= 0 {
			return false
		}
	}

	return true
}

// RateProduct adds or updates a rating for a product
func RateProduct(productID string, rating float64) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	var data ProductData
	err := ReadJSON(productsFile, &data)
	if err != nil {
		return err
	}

	for i, product := range data.Products {
		if product.ID == productID {
			// Calculate new average rating
			totalRating := product.Rating * float64(product.RatingCount)
			totalRating += rating
			data.Products[i].RatingCount++
			data.Products[i].Rating = totalRating / float64(data.Products[i].RatingCount)
			data.Products[i].UpdatedAt = time.Now()
			data.Metadata.UpdatedAt = time.Now()
			
			return WriteJSON(productsFile, data)
		}
	}

	return fmt.Errorf("product not found")
}

// UpdateProductStock updates the stock quantity of a product
func UpdateProductStock(productID string, newStock int) error {
	var data ProductData
	err := ReadJSON(productsFile, &data)
	if err != nil {
		return err
	}

	for i, product := range data.Products {
		if product.ID == productID {
			data.Products[i].Stock = newStock
			data.Products[i].UpdatedAt = time.Now()
			data.Metadata.UpdatedAt = time.Now()
			
			return WriteJSON(productsFile, data)
		}
	}

	return fmt.Errorf("product not found")
}

// GetProductsByCategory retrieves products by category
func GetProductsByCategory(category string) ([]Product, error) {
	var data ProductData
	err := ReadJSON(productsFile, &data)
	if err != nil {
		return nil, err
	}

	var categoryProducts []Product
	for _, product := range data.Products {
		if product.IsActive && product.Category == category {
			categoryProducts = append(categoryProducts, product)
		}
	}

	return categoryProducts, nil
}

// Helper functions
func isHalalProduct(ingredients []string) bool {
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

func isVegetarianProduct(ingredients []string) bool {
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

func isVeganProduct(ingredients []string) bool {
	if !isVegetarianProduct(ingredients) {
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

func generateSKU(name, brand string) string {
	// Simple SKU generation: first 3 chars of brand + first 3 chars of name + timestamp
	brandCode := strings.ToUpper(brand)
	if len(brandCode) > 3 {
		brandCode = brandCode[:3]
	}
	nameCode := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(nameCode) > 3 {
		nameCode = nameCode[:3]
	}
	timestamp := time.Now().Unix() % 10000
	return fmt.Sprintf("%s%s%04d", brandCode, nameCode, timestamp)
}