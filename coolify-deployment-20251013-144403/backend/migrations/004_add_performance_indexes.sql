-- Migration: Add Performance Indexes
-- Version: 004
-- Description: Adds missing indexes to improve query performance and prevent N+1 queries

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_activity_level ON users(activity_level);
CREATE INDEX IF NOT EXISTS idx_users_preferred_language ON users(preferred_language);

-- Foods table indexes
CREATE INDEX IF NOT EXISTS idx_foods_name ON foods(name);
CREATE INDEX IF NOT EXISTS idx_foods_name_ar ON foods(name_ar);
CREATE INDEX IF NOT EXISTS idx_foods_category ON foods(category);
CREATE INDEX IF NOT EXISTS idx_foods_subcategory ON foods(subcategory);
CREATE INDEX IF NOT EXISTS idx_foods_brand ON foods(brand);
CREATE INDEX IF NOT EXISTS idx_foods_barcode ON foods(barcode);
CREATE INDEX IF NOT EXISTS idx_foods_verified ON foods(verified);
CREATE INDEX IF NOT EXISTS idx_foods_is_halal ON foods(is_halal);
CREATE INDEX IF NOT EXISTS idx_foods_is_vegetarian ON foods(is_vegetarian);
CREATE INDEX IF NOT EXISTS idx_foods_is_vegan ON foods(is_vegan);
CREATE INDEX IF NOT EXISTS idx_foods_contains_alcohol ON foods(contains_alcohol);
CREATE INDEX IF NOT EXISTS idx_foods_contains_pork ON foods(contains_pork);
CREATE INDEX IF NOT EXISTS idx_foods_category_subcategory ON foods(category, subcategory);
CREATE INDEX IF NOT EXISTS idx_foods_name_category ON foods(name, category);

-- Exercises table indexes
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises(name);
CREATE INDEX IF NOT EXISTS idx_exercises_name_ar ON exercises(name_ar);
CREATE INDEX IF NOT EXISTS idx_exercises_category ON exercises(category);
CREATE INDEX IF NOT EXISTS idx_exercises_equipment ON exercises(equipment);
CREATE INDEX IF NOT EXISTS idx_exercises_difficulty_level ON exercises(difficulty_level);
CREATE INDEX IF NOT EXISTS idx_exercises_verified ON exercises(verified);
CREATE INDEX IF NOT EXISTS idx_exercises_category_difficulty ON exercises(category, difficulty_level);

-- Meal plans table indexes
CREATE INDEX IF NOT EXISTS idx_meal_plans_user_id ON meal_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_meal_plans_is_active ON meal_plans(is_active);
CREATE INDEX IF NOT EXISTS idx_meal_plans_start_date ON meal_plans(start_date);
CREATE INDEX IF NOT EXISTS idx_meal_plans_end_date ON meal_plans(end_date);
CREATE INDEX IF NOT EXISTS idx_meal_plans_created_at ON meal_plans(created_at);
CREATE INDEX IF NOT EXISTS idx_meal_plans_user_active ON meal_plans(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_meal_plans_user_dates ON meal_plans(user_id, start_date, end_date);

-- Workout plans table indexes
CREATE INDEX IF NOT EXISTS idx_workout_plans_user_id ON workout_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_plans_is_active ON workout_plans(is_active);
CREATE INDEX IF NOT EXISTS idx_workout_plans_difficulty_level ON workout_plans(difficulty_level);
CREATE INDEX IF NOT EXISTS idx_workout_plans_goal ON workout_plans(goal);
CREATE INDEX IF NOT EXISTS idx_workout_plans_start_date ON workout_plans(start_date);
CREATE INDEX IF NOT EXISTS idx_workout_plans_end_date ON workout_plans(end_date);
CREATE INDEX IF NOT EXISTS idx_workout_plans_created_at ON workout_plans(created_at);
CREATE INDEX IF NOT EXISTS idx_workout_plans_user_active ON workout_plans(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_workout_plans_user_dates ON workout_plans(user_id, start_date, end_date);

-- User food logs table indexes
CREATE INDEX IF NOT EXISTS idx_user_food_logs_user_id ON user_food_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_food_id ON user_food_logs(food_id);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_meal_type ON user_food_logs(meal_type);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_consumed_at ON user_food_logs(consumed_at);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_created_at ON user_food_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_user_consumed ON user_food_logs(user_id, consumed_at);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_user_meal_type ON user_food_logs(user_id, meal_type);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_user_date_meal ON user_food_logs(user_id, DATE(consumed_at), meal_type);

-- User exercise logs table indexes
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_user_id ON user_exercise_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_exercise_id ON user_exercise_logs(exercise_id);
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_performed_at ON user_exercise_logs(performed_at);
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_created_at ON user_exercise_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_user_performed ON user_exercise_logs(user_id, performed_at);
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_user_date ON user_exercise_logs(user_id, DATE(performed_at));

-- System metrics table indexes
CREATE INDEX IF NOT EXISTS idx_system_metrics_metric_name ON system_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_system_metrics_metric_type ON system_metrics(metric_type);
CREATE INDEX IF NOT EXISTS idx_system_metrics_timestamp ON system_metrics(timestamp);
CREATE INDEX IF NOT EXISTS idx_system_metrics_name_timestamp ON system_metrics(metric_name, timestamp);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_foods_search_composite ON foods(name, category, is_halal, is_vegetarian, is_vegan);
CREATE INDEX IF NOT EXISTS idx_exercises_search_composite ON exercises(name, category, difficulty_level, equipment);
CREATE INDEX IF NOT EXISTS idx_user_logs_daily_summary ON user_food_logs(user_id, DATE(consumed_at), meal_type);
CREATE INDEX IF NOT EXISTS idx_user_exercise_daily_summary ON user_exercise_logs(user_id, DATE(performed_at));

-- Text search indexes (for PostgreSQL full-text search)
-- These will be ignored in SQLite but useful for PostgreSQL
CREATE INDEX IF NOT EXISTS idx_foods_name_gin ON foods USING gin(to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_foods_description_gin ON foods USING gin(to_tsvector('english', description));
CREATE INDEX IF NOT EXISTS idx_exercises_name_gin ON exercises USING gin(to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_exercises_description_gin ON exercises USING gin(to_tsvector('english', description));

-- Partial indexes for active records only (PostgreSQL specific)
CREATE INDEX IF NOT EXISTS idx_users_active_email ON users(email) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_meal_plans_active_user ON meal_plans(user_id, created_at) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_workout_plans_active_user ON workout_plans(user_id, created_at) WHERE is_active = true;

-- Add comments for documentation
COMMENT ON INDEX idx_foods_search_composite IS 'Composite index for food search queries with dietary filters';
COMMENT ON INDEX idx_exercises_search_composite IS 'Composite index for exercise search queries with filters';
COMMENT ON INDEX idx_user_logs_daily_summary IS 'Optimizes daily nutrition summary queries';
COMMENT ON INDEX idx_user_exercise_daily_summary IS 'Optimizes daily exercise summary queries';