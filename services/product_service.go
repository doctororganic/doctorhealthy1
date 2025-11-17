package services

import (
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// ProductService handles product-related business logic
type ProductService struct {
	db *sql.DB
}

// NewProductService creates a new product service
func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

// GetProducts retrieves products (public access)
func (s *ProductService) GetProducts(limit, offset int, filters map[string]interface{}) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, category, brand, image_url,
			   is_approved, created_by, created_at, updated_at
		FROM products 
		WHERE is_approved = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdBy sql.NullInt64
		
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Category, &product.Brand, &product.ImageURL,
			&product.IsApproved, &createdBy, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if createdBy.Valid {
			product.CreatedBy = uint(createdBy.Int64)
		}
		
		products = append(products, product)
	}
	
	return products, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, category, brand,
							image_url, is_approved, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		product.Name, product.Description, product.Price, product.Category,
		product.Brand, product.ImageURL, product.IsApproved,
		product.CreatedBy, now, now,
	).Scan(&product.ID)
	
	if err != nil {
		return err
	}
	
	product.CreatedAt = now
	product.UpdatedAt = now
	
	return nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, category, brand, image_url,
			   is_approved, created_by, created_at, updated_at
		FROM products 
		WHERE id = $1
	`
	
	var product models.Product
	var createdBy sql.NullInt64
	
	err := s.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.Category, &product.Brand, &product.ImageURL,
		&product.IsApproved, &createdBy, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	
	if createdBy.Valid {
		product.CreatedBy = uint(createdBy.Int64)
	}
	
	return &product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $2, description = $3, price = $4, category = $5,
			brand = $6, image_url = $7, updated_at = $8
		WHERE id = $1
	`
	
	product.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(query,
		product.ID, product.Name, product.Description, product.Price,
		product.Category, product.Brand, product.ImageURL, product.UpdatedAt,
	)
	
	return err
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(id uint) error {
	query := `DELETE FROM products WHERE id = $1`
	
	_, err := s.db.Exec(query, id)
	return err
}

// GetPendingProducts retrieves pending products for admin review
func (s *ProductService) GetPendingProducts(limit, offset int) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, category, brand, image_url,
			   is_approved, created_by, created_at, updated_at
		FROM products 
		WHERE is_approved = false
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdBy sql.NullInt64
		
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Category, &product.Brand, &product.ImageURL,
			&product.IsApproved, &createdBy, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if createdBy.Valid {
			product.CreatedBy = uint(createdBy.Int64)
		}
		
		products = append(products, product)
	}
	
	return products, nil
}

// ApproveProduct approves a product
func (s *ProductService) ApproveProduct(id uint) error {
	query := `UPDATE products SET is_approved = true, updated_at = $1 WHERE id = $2`
	
	_, err := s.db.Exec(query, time.Now(), id)
	return err
}

// RejectProduct rejects a product
func (s *ProductService) RejectProduct(id uint) error {
	query := `DELETE FROM products WHERE id = $1`
	
	_, err := s.db.Exec(query, id)
	return err
}
