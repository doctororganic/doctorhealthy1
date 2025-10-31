-- Migration: Create files management tables
-- Created: 2025-01-25
-- Description: Tables for file upload, storage, and management

-- Files table for storing file metadata
CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY,
    file_name TEXT NOT NULL,
    original_name TEXT NOT NULL,
    file_url TEXT NOT NULL,
    thumbnail_url TEXT,
    file_path TEXT NOT NULL,
    size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    uploader_id TEXT NOT NULL,
    purpose TEXT NOT NULL DEFAULT 'profile',
    status TEXT NOT NULL DEFAULT 'active',
    metadata TEXT, -- JSON field for additional metadata
    hash TEXT NOT NULL, -- SHA-256 hash for deduplication
    visibility TEXT NOT NULL DEFAULT 'private',
    uploaded_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    FOREIGN KEY (uploader_id) REFERENCES users(id) ON DELETE CASCADE
);

-- User profile images table
CREATE TABLE IF NOT EXISTS user_profile_images (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    file_id TEXT NOT NULL,
    file_url TEXT NOT NULL,
    thumb_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    position INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Meal photos table for meal-related images
CREATE TABLE IF NOT EXISTS meal_photos (
    id TEXT PRIMARY KEY,
    meal_id TEXT,
    user_id TEXT NOT NULL,
    file_id TEXT NOT NULL,
    file_url TEXT NOT NULL,
    thumb_url TEXT,
    caption TEXT,
    position INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (meal_id) REFERENCES meals(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Progress photos table for before/after photos
CREATE TABLE IF NOT EXISTS progress_photos (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    file_id TEXT NOT NULL,
    file_url TEXT NOT NULL,
    thumb_url TEXT,
    photo_type TEXT NOT NULL DEFAULT 'front', -- front, side, back, before, after
    weight REAL,
    body_fat REAL,
    measurements TEXT, -- JSON string for measurements
    notes TEXT,
    taken_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Recipe images table for recipe photos
CREATE TABLE IF NOT EXISTS recipe_images (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    file_id TEXT NOT NULL,
    file_url TEXT NOT NULL,
    thumb_url TEXT,
    caption TEXT,
    position INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- File upload sessions for chunked uploads
CREATE TABLE IF NOT EXISTS file_upload_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    chunk_size INTEGER NOT NULL,
    total_chunks INTEGER NOT NULL,
    uploaded_chunks INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'initiated',
    metadata TEXT, -- JSON field for additional metadata
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- File processing jobs for background processing
CREATE TABLE IF NOT EXISTS file_processing_jobs (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    job_type TEXT NOT NULL, -- thumbnail, resize, optimize, analyze, convert
    status TEXT NOT NULL DEFAULT 'pending',
    progress INTEGER NOT NULL DEFAULT 0, -- 0-100
    result TEXT, -- JSON field for job results
    error TEXT,
    created_at DATETIME NOT NULL,
    started_at DATETIME,
    completed_at DATETIME,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- File shares for public/private file sharing
CREATE TABLE IF NOT EXISTS file_shares (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    shared_by TEXT NOT NULL,
    share_token TEXT NOT NULL UNIQUE,
    share_type TEXT NOT NULL DEFAULT 'private', -- public, private, password
    password TEXT, -- hashed password for password-protected shares
    expires_at DATETIME,
    view_count INTEGER NOT NULL DEFAULT 0,
    max_views INTEGER,
    permissions TEXT NOT NULL DEFAULT 'view', -- JSON array of permissions
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_by) REFERENCES users(id) ON DELETE CASCADE
);

-- File validation results
CREATE TABLE IF NOT EXISTS file_validation (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    rules TEXT NOT NULL, -- JSON field for validation rules
    passed BOOLEAN NOT NULL,
    validation_log TEXT, -- JSON array of validation messages
    processed_at DATETIME NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Indexes for better performance
-- Files table indexes
CREATE INDEX IF NOT EXISTS idx_files_uploader_id ON files(uploader_id);
CREATE INDEX IF NOT EXISTS idx_files_purpose ON files(purpose);
CREATE INDEX IF NOT EXISTS idx_files_status ON files(status);
CREATE INDEX IF NOT EXISTS idx_files_content_type ON files(content_type);
CREATE INDEX IF NOT EXISTS idx_files_hash ON files(hash);
CREATE INDEX IF NOT EXISTS idx_files_uploaded_at ON files(uploaded_at);
CREATE INDEX IF NOT EXISTS idx_files_visibility ON files(visibility);
CREATE INDEX IF NOT EXISTS idx_files_uploader_purpose ON files(uploader_id, purpose);
CREATE INDEX IF NOT EXISTS idx_files_uploader_status ON files(uploader_id, status);

-- User profile images indexes
CREATE INDEX IF NOT EXISTS idx_user_profile_images_user_id ON user_profile_images(user_id);
CREATE INDEX IF NOT EXISTS idx_user_profile_images_is_active ON user_profile_images(is_active);
CREATE INDEX IF NOT EXISTS idx_user_profile_images_position ON user_profile_images(position);

-- Meal photos indexes
CREATE INDEX IF NOT EXISTS idx_meal_photos_meal_id ON meal_photos(meal_id);
CREATE INDEX IF NOT EXISTS idx_meal_photos_user_id ON meal_photos(user_id);
CREATE INDEX IF NOT EXISTS idx_meal_photos_position ON meal_photos(position);

-- Progress photos indexes
CREATE INDEX IF NOT EXISTS idx_progress_photos_user_id ON progress_photos(user_id);
CREATE INDEX IF NOT EXISTS idx_progress_photos_photo_type ON progress_photos(photo_type);
CREATE INDEX IF NOT EXISTS idx_progress_photos_taken_at ON progress_photos(taken_at);
CREATE INDEX IF NOT EXISTS idx_progress_photos_user_type ON progress_photos(user_id, photo_type);

-- Recipe images indexes
CREATE INDEX IF NOT EXISTS idx_recipe_images_recipe_id ON recipe_images(recipe_id);
CREATE INDEX IF NOT EXISTS idx_recipe_images_is_primary ON recipe_images(is_primary);
CREATE INDEX IF NOT EXISTS idx_recipe_images_position ON recipe_images(position);

-- Upload sessions indexes
CREATE INDEX IF NOT EXISTS idx_file_upload_sessions_user_id ON file_upload_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_file_upload_sessions_status ON file_upload_sessions(status);
CREATE INDEX IF NOT EXISTS idx_file_upload_sessions_expires_at ON file_upload_sessions(expires_at);

-- Processing jobs indexes
CREATE INDEX IF NOT EXISTS idx_file_processing_jobs_file_id ON file_processing_jobs(file_id);
CREATE INDEX IF NOT EXISTS idx_file_processing_jobs_status ON file_processing_jobs(status);
CREATE INDEX IF NOT EXISTS idx_file_processing_jobs_job_type ON file_processing_jobs(job_type);
CREATE INDEX IF NOT EXISTS idx_file_processing_jobs_created_at ON file_processing_jobs(created_at);

-- File shares indexes
CREATE INDEX IF NOT EXISTS idx_file_shares_file_id ON file_shares(file_id);
CREATE INDEX IF NOT EXISTS idx_file_shares_shared_by ON file_shares(shared_by);
CREATE INDEX IF NOT EXISTS idx_file_shares_share_token ON file_shares(share_token);
CREATE INDEX IF NOT EXISTS idx_file_shares_expires_at ON file_shares(expires_at);
CREATE INDEX IF NOT EXISTS idx_file_shares_share_type ON file_shares(share_type);

-- Validation indexes
CREATE INDEX IF NOT EXISTS idx_file_validation_file_id ON file_validation(file_id);
CREATE INDEX IF NOT EXISTS idx_file_validation_processed_at ON file_validation(processed_at);

-- Triggers for automatic timestamp updates
CREATE TRIGGER IF NOT EXISTS update_files_updated_at
    AFTER UPDATE ON files
    FOR EACH ROW
BEGIN
    UPDATE files SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_user_profile_images_updated_at
    AFTER UPDATE ON user_profile_images
    FOR EACH ROW
BEGIN
    UPDATE user_profile_images SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_meal_photos_updated_at
    AFTER UPDATE ON meal_photos
    FOR EACH ROW
BEGIN
    UPDATE meal_photos SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_progress_photos_updated_at
    AFTER UPDATE ON progress_photos
    FOR EACH ROW
BEGIN
    UPDATE progress_photos SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_recipe_images_updated_at
    AFTER UPDATE ON recipe_images
    FOR EACH ROW
BEGIN
    UPDATE recipe_images SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_file_upload_sessions_updated_at
    AFTER UPDATE ON file_upload_sessions
    FOR EACH ROW
BEGIN
    UPDATE file_upload_sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_file_shares_updated_at
    AFTER UPDATE ON file_shares
    FOR EACH ROW
BEGIN
    UPDATE file_shares SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Views for common queries
CREATE VIEW IF NOT EXISTS active_files AS
SELECT 
    f.*,
    u.username as uploader_name,
    u.email as uploader_email
FROM files f
JOIN users u ON f.uploader_id = u.id
WHERE f.status = 'active' AND f.deleted_at IS NULL;

CREATE VIEW IF NOT EXISTS user_file_stats AS
SELECT 
    uploader_id,
    COUNT(*) as total_files,
    SUM(size) as total_size,
    COUNT(CASE WHEN content_type LIKE 'image/%' THEN 1 END) as image_count,
    COUNT(CASE WHEN content_type LIKE 'video/%' THEN 1 END) as video_count,
    MAX(uploaded_at) as last_upload_date
FROM files
WHERE status = 'active' AND deleted_at IS NULL
GROUP BY uploader_id;

-- Storage limits and policies (could be moved to a separate policies table)
-- For now, these are default constraints that should be enforced at the application level

-- Example storage limits per user tier (these should be adjusted based on your business model)
-- Free tier: 100MB, Pro tier: 1GB, Enterprise tier: 10GB

-- File size limits by purpose
-- Profile images: 5MB max
-- Meal photos: 10MB max  
-- Progress photos: 10MB max
-- Recipe images: 5MB max
-- Documents: 20MB max
