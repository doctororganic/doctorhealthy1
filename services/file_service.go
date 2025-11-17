package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nutrition-platform/config"
	"nutrition-platform/models"
)

// FileService handles file upload and management business logic
type FileService struct {
	db         *sql.DB
	uploadPath string
	maxSize    int64
}

// NewFileService creates a new file service
func NewFileService(cfg config.UploadConfig) *FileService {
	return &FileService{
		uploadPath: cfg.Path,
		maxSize:    cfg.MaxSize,
	}
}

// SetDatabase sets the database connection (for dependency injection)
func (s *FileService) SetDatabase(db *sql.DB) {
	s.db = db
}

// UploadFile uploads a file and saves its metadata
func (s *FileService) UploadFile(file *multipart.FileHeader, userID uint) (*models.FileMetadata, error) {
	// Validate file size
	if file.Size > s.maxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size")
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Calculate file hash
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// Reset file pointer
	src.Close()
	src, err = file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to reopen uploaded file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s_%d%s", fileHash[:8], time.Now().Unix(), ext)
	filePath := filepath.Join(s.uploadPath, uniqueFilename)

	// Ensure upload directory exists
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath) // Clean up on failure
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Save file metadata to database
	fileMetadata := &models.FileMetadata{
		UserID:      userID,
		Filename:    file.Filename,
		StoragePath: filePath,
		FileSize:    file.Size,
		ContentType: file.Header.Get("Content-Type"),
		FileHash:    fileHash,
		IsPublic:    false,
	}

	if err := s.CreateFileMetadata(fileMetadata); err != nil {
		os.Remove(filePath) // Clean up on failure
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return fileMetadata, nil
}

// GetFile retrieves file metadata by ID
func (s *FileService) GetFile(id, userID uint) (*models.FileMetadata, error) {
	query := `
		SELECT id, user_id, filename, storage_path, file_size, content_type,
			   file_hash, is_public, download_count, created_at, updated_at
		FROM file_metadata 
		WHERE id = $1 AND user_id = $2
	`
	
	var file models.FileMetadata
	
	err := s.db.QueryRow(query, id, userID).Scan(
		&file.ID, &file.UserID, &file.Filename, &file.StoragePath,
		&file.FileSize, &file.ContentType, &file.FileHash, &file.IsPublic,
		&file.DownloadCount, &file.CreatedAt, &file.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, err
	}
	
	return &file, nil
}

// GetPublicFile retrieves public file metadata
func (s *FileService) GetPublicFile(id uint) (*models.FileMetadata, error) {
	query := `
		SELECT id, user_id, filename, storage_path, file_size, content_type,
			   file_hash, is_public, download_count, created_at, updated_at
		FROM file_metadata 
		WHERE id = $1 AND is_public = true
	`
	
	var file models.FileMetadata
	
	err := s.db.QueryRow(query, id).Scan(
		&file.ID, &file.UserID, &file.Filename, &file.StoragePath,
		&file.FileSize, &file.ContentType, &file.FileHash, &file.IsPublic,
		&file.DownloadCount, &file.CreatedAt, &file.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, err
	}
	
	return &file, nil
}

// CreateFileMetadata creates file metadata record
func (s *FileService) CreateFileMetadata(file *models.FileMetadata) error {
	query := `
		INSERT INTO file_metadata (user_id, filename, storage_path, file_size,
								 content_type, file_hash, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	
	now := time.Now()
	
	err := s.db.QueryRow(query,
		file.UserID, file.Filename, file.StoragePath, file.FileSize,
		file.ContentType, file.FileHash, file.IsPublic, now, now,
	).Scan(&file.ID)
	
	if err != nil {
		return err
	}
	
	file.CreatedAt = now
	file.UpdatedAt = now
	
	return nil
}

// DeleteFile deletes a file and its metadata
func (s *FileService) DeleteFile(id, userID uint) error {
	// Get file metadata first
	file, err := s.GetFile(id, userID)
	if err != nil {
		return err
	}

	// Delete file from disk
	if err := os.Remove(file.StoragePath); err != nil {
		// Log error but continue with database deletion
		fmt.Printf("Warning: failed to delete file from disk: %v\n", err)
	}

	// Delete metadata from database
	query := `DELETE FROM file_metadata WHERE id = $1 AND user_id = $2`
	
	_, err = s.db.Exec(query, id, userID)
	return err
}

// IncrementDownloadCount increments the download count for a file
func (s *FileService) IncrementDownloadCount(id uint) error {
	query := `UPDATE file_metadata SET download_count = download_count + 1 WHERE id = $1`
	
	_, err := s.db.Exec(query, id)
	return err
}

// ValidateFileType validates if the file type is allowed
func (s *FileService) ValidateFileType(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Define allowed extensions
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".txt":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".csv":  true,
	}
	
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type %s is not allowed", ext)
	}
	
	return nil
}
