-- Initial database schema for Nutrition Platform
-- This migration creates all the necessary tables for the application

-- Enable UUID extension for PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE,
    gender VARCHAR(10),
    height DECIMAL(5,2),
    weight DECIMAL(5,2),
    activity_level VARCHAR(20),
    goals TEXT,
    dietary_restrictions JSONB DEFAULT '[]',
    religious_filter_enabled BOOLEAN DEFAULT true,
    filter_alcohol BOOLEAN DEFAULT true,
    filter_pork BOOLEAN DEFAULT true,
    preferred_language VARCHAR(5) DEFAULT 'en',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Foods table
CREATE TABLE IF NOT EXISTS foods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    name_ar VARCHAR(100),
    description TEXT,
    description_ar TEXT,
    category VARCHAR(50),
    subcategory VARCHAR(50),
    barcode VARCHAR(50),
    brand VARCHAR(100),
    serving_size DECIMAL(8,2),
    serving_unit VARCHAR(20),
    calories_per_100g DECIMAL(8,2),
    protein_per_100g DECIMAL(8,2),
    carbs_per_100g DECIMAL(8,2),
    fat_per_100g DECIMAL(8,2),
    fiber_per_100g DECIMAL(8,2),
    sugar_per_100g DECIMAL(8,2),
    sodium_per_100g DECIMAL(8,2),
    ingredients JSONB DEFAULT '[]',
    allergens JSONB DEFAULT '[]',
    contains_alcohol BOOLEAN DEFAULT false,
    contains_pork BOOLEAN DEFAULT false,
    is_halal BOOLEAN DEFAULT true,
    is_kosher BOOLEAN DEFAULT false,
    is_vegetarian BOOLEAN DEFAULT false,
    is_vegan BOOLEAN DEFAULT false,
    image_url VARCHAR(255),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Exercises table
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    name_ar VARCHAR(100),
    description TEXT,
    description_ar TEXT,
    category VARCHAR(50),
    muscle_groups JSONB DEFAULT '[]',
    equipment VARCHAR(50),
    difficulty_level VARCHAR(20),
    instructions TEXT,
    instructions_ar TEXT,
    calories_per_minute DECIMAL(5,2),
    met_value DECIMAL(4,2),
    image_url VARCHAR(255),
    video_url VARCHAR(255),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- API keys table
CREATE TABLE IF NOT EXISTS api_keys (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    prefix VARCHAR(10) NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'active',
    scopes JSONB NOT NULL,
    rate_limit INTEGER DEFAULT 100,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}'
);

-- API key usage table
CREATE TABLE IF NOT EXISTS api_key_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_key_id VARCHAR(50) REFERENCES api_keys(id) ON DELETE CASCADE,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INTEGER NOT NULL,
    response_time BIGINT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- API metrics snapshots table
CREATE TABLE IF NOT EXISTS api_metrics_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_key_id VARCHAR(50) REFERENCES api_keys(id) ON DELETE CASCADE,
    metrics_data JSONB NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(api_key_id, timestamp)
);

-- Usage alerts table
CREATE TABLE IF NOT EXISTS usage_alerts (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    condition VARCHAR(100) NOT NULL,
    threshold DECIMAL(10,2) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    last_triggered TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Meal plans table
CREATE TABLE IF NOT EXISTS meal_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    start_date DATE,
    end_date DATE,
    total_calories INTEGER,
    total_protein DECIMAL(8,2),
    total_carbs DECIMAL(8,2),
    total_fat DECIMAL(8,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workout plans table
CREATE TABLE IF NOT EXISTS workout_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    start_date DATE,
    end_date DATE,
    difficulty_level VARCHAR(20),
    goal VARCHAR(50),
    days_per_week INTEGER,
    duration_weeks INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User food logs table
CREATE TABLE IF NOT EXISTS user_food_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    food_id UUID REFERENCES foods(id),
    quantity DECIMAL(8,2) NOT NULL,
    unit VARCHAR(20),
    meal_type VARCHAR(20),
    consumed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    calories DECIMAL(8,2),
    protein DECIMAL(8,2),
    carbs DECIMAL(8,2),
    fat DECIMAL(8,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User exercise logs table
CREATE TABLE IF NOT EXISTS user_exercise_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id),
    duration_minutes INTEGER,
    sets INTEGER,
    reps INTEGER,
    weight DECIMAL(6,2),
    calories_burned DECIMAL(8,2),
    notes TEXT,
    performed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Recipes table
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    name_ar VARCHAR(200),
    description TEXT,
    description_ar TEXT,
    cuisine VARCHAR(50),
    country VARCHAR(50),
    difficulty_level VARCHAR(20),
    prep_time_minutes INTEGER,
    cook_time_minutes INTEGER,
    total_time_minutes INTEGER,
    servings INTEGER,
    ingredients JSONB NOT NULL DEFAULT '[]',
    instructions JSONB NOT NULL DEFAULT '[]',
    nutrition_per_serving JSONB DEFAULT '{}',
    dietary_tags JSONB DEFAULT '[]', -- vegetarian, vegan, gluten-free, etc.
    allergens JSONB DEFAULT '[]',
    is_halal BOOLEAN DEFAULT true,
    is_kosher BOOLEAN DEFAULT false,
    image_url VARCHAR(255),
    video_url VARCHAR(255),
    rating DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    created_by UUID REFERENCES users(id),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Diseases/Health Conditions table
CREATE TABLE IF NOT EXISTS health_conditions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    name_ar VARCHAR(200),
    category VARCHAR(100), -- chronic, acute, metabolic, etc.
    icd_10_code VARCHAR(20),
    description TEXT,
    description_ar TEXT,
    symptoms JSONB DEFAULT '[]',
    risk_factors JSONB DEFAULT '[]',
    complications JSONB DEFAULT '[]',
    dietary_recommendations JSONB DEFAULT '[]',
    exercise_recommendations JSONB DEFAULT '[]',
    lifestyle_modifications JSONB DEFAULT '[]',
    severity_levels JSONB DEFAULT '[]', -- mild, moderate, severe
    is_chronic BOOLEAN DEFAULT false,
    requires_medical_supervision BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Health Complaints/Symptoms
CREATE TABLE IF NOT EXISTS user_health_complaints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    complaint_type VARCHAR(100) NOT NULL,
    severity INTEGER CHECK (severity >= 1 AND severity <= 10),
    description TEXT,
    symptoms JSONB DEFAULT '[]',
    duration_days INTEGER,
    frequency VARCHAR(50), -- daily, weekly, occasional
    triggers JSONB DEFAULT '[]',
    current_medications JSONB DEFAULT '[]',
    reported_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active', -- active, resolved, monitoring
    medical_attention_required BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Injuries table
CREATE TABLE IF NOT EXISTS injuries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    name_ar VARCHAR(200),
    category VARCHAR(100), -- acute, chronic, sports, occupational
    body_part VARCHAR(100),
    severity_level VARCHAR(50), -- minor, moderate, severe, critical
    description TEXT,
    description_ar TEXT,
    symptoms JSONB DEFAULT '[]',
    causes JSONB DEFAULT '[]',
    treatment_options JSONB DEFAULT '[]',
    recovery_time_days INTEGER,
    exercise_restrictions JSONB DEFAULT '[]',
    recommended_exercises JSONB DEFAULT '[]',
    nutrition_recommendations JSONB DEFAULT '[]',
    prevention_tips JSONB DEFAULT '[]',
    when_to_seek_help TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Injuries
CREATE TABLE IF NOT EXISTS user_injuries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    injury_id UUID REFERENCES injuries(id),
    custom_injury_name VARCHAR(200), -- if not in predefined list
    severity INTEGER CHECK (severity >= 1 AND severity <= 10),
    injury_date DATE,
    description TEXT,
    treatment_received TEXT,
    current_status VARCHAR(50) DEFAULT 'healing', -- healing, recovered, chronic
    affects_exercise BOOLEAN DEFAULT true,
    exercise_limitations JSONB DEFAULT '[]',
    medical_clearance_required BOOLEAN DEFAULT false,
    expected_recovery_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Nutritional Plans (more detailed than meal plans)
CREATE TABLE IF NOT EXISTS nutritional_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    plan_type VARCHAR(100), -- weight_loss, muscle_gain, medical, therapeutic
    health_condition_id UUID REFERENCES health_conditions(id),
    description TEXT,
    duration_weeks INTEGER,
    daily_calorie_target INTEGER,
    macro_targets JSONB NOT NULL, -- protein, carbs, fat percentages
    micro_targets JSONB DEFAULT '{}', -- vitamins, minerals
    meal_timing JSONB DEFAULT '[]',
    food_restrictions JSONB DEFAULT '[]',
    recommended_foods JSONB DEFAULT '[]',
    foods_to_avoid JSONB DEFAULT '[]',
    supplement_recommendations JSONB DEFAULT '[]',
    hydration_target_ml INTEGER,
    special_instructions TEXT,
    created_by VARCHAR(100), -- nutritionist, system, user
    medical_approval_required BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Drugs/Medications
CREATE TABLE IF NOT EXISTS medications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    name_ar VARCHAR(200),
    generic_name VARCHAR(200),
    brand_names JSONB DEFAULT '[]',
    drug_class VARCHAR(100),
    category VARCHAR(100), -- prescription, otc, supplement
    description TEXT,
    description_ar TEXT,
    indications JSONB DEFAULT '[]', -- what it treats
    contraindications JSONB DEFAULT '[]',
    side_effects JSONB DEFAULT '[]',
    drug_interactions JSONB DEFAULT '[]',
    food_interactions JSONB DEFAULT '[]',
    dosage_forms JSONB DEFAULT '[]', -- tablet, capsule, liquid
    typical_dosages JSONB DEFAULT '[]',
    administration_route VARCHAR(50), -- oral, topical, injection
    pregnancy_category VARCHAR(10),
    requires_prescription BOOLEAN DEFAULT true,
    affects_nutrition BOOLEAN DEFAULT false,
    nutritional_effects JSONB DEFAULT '{}',
    monitoring_requirements TEXT,
    storage_requirements TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Medications
CREATE TABLE IF NOT EXISTS user_medications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    medication_id UUID REFERENCES medications(id),
    custom_medication_name VARCHAR(200), -- if not in predefined list
    dosage VARCHAR(100),
    frequency VARCHAR(100),
    administration_time JSONB DEFAULT '[]', -- morning, evening, with meals
    start_date DATE,
    end_date DATE,
    prescribed_by VARCHAR(200),
    reason_for_taking TEXT,
    side_effects_experienced JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    adherence_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Vitamins and Minerals
CREATE TABLE IF NOT EXISTS vitamins_minerals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    name_ar VARCHAR(100),
    type VARCHAR(50), -- vitamin, mineral, trace_element
    category VARCHAR(50), -- fat_soluble, water_soluble, macro_mineral, trace_mineral
    chemical_name VARCHAR(200),
    description TEXT,
    description_ar TEXT,
    functions JSONB DEFAULT '[]', -- what it does in the body
    deficiency_symptoms JSONB DEFAULT '[]',
    toxicity_symptoms JSONB DEFAULT '[]',
    food_sources JSONB DEFAULT '[]',
    daily_requirements JSONB DEFAULT '{}', -- by age/gender groups
    upper_limit JSONB DEFAULT '{}',
    absorption_factors JSONB DEFAULT '[]',
    interactions JSONB DEFAULT '[]', -- with other nutrients/drugs
    best_taken_with JSONB DEFAULT '[]',
    avoid_taking_with JSONB DEFAULT '[]',
    supplement_forms JSONB DEFAULT '[]',
    stability_factors TEXT, -- heat, light, air sensitivity
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Vitamin/Mineral Tracking
CREATE TABLE IF NOT EXISTS user_supplements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    vitamin_mineral_id UUID REFERENCES vitamins_minerals(id),
    supplement_name VARCHAR(200),
    brand VARCHAR(100),
    dosage VARCHAR(100),
    form VARCHAR(50), -- tablet, capsule, liquid, powder
    frequency VARCHAR(100),
    taken_with_meals BOOLEAN DEFAULT true,
    start_date DATE,
    end_date DATE,
    reason_for_taking TEXT,
    prescribed_by VARCHAR(200),
    cost_per_month DECIMAL(10,2),
    effectiveness_rating INTEGER CHECK (effectiveness_rating >= 1 AND effectiveness_rating <= 5),
    side_effects TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced Workout Plans
CREATE TABLE IF NOT EXISTS workout_programs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    name_ar VARCHAR(200),
    description TEXT,
    description_ar TEXT,
    program_type VARCHAR(100), -- strength, cardio, flexibility, rehabilitation
    fitness_level VARCHAR(50), -- beginner, intermediate, advanced
    duration_weeks INTEGER,
    days_per_week INTEGER,
    session_duration_minutes INTEGER,
    equipment_required JSONB DEFAULT '[]',
    target_goals JSONB DEFAULT '[]', -- weight_loss, muscle_gain, endurance
    muscle_groups_targeted JSONB DEFAULT '[]',
    contraindications JSONB DEFAULT '[]',
    modifications_available JSONB DEFAULT '[]',
    progression_plan JSONB DEFAULT '[]',
    created_by VARCHAR(100),
    difficulty_rating INTEGER CHECK (difficulty_rating >= 1 AND difficulty_rating <= 5),
    calorie_burn_estimate INTEGER, -- per session
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workout Sessions (individual workouts within a program)
CREATE TABLE IF NOT EXISTS workout_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_program_id UUID REFERENCES workout_programs(id) ON DELETE CASCADE,
    session_number INTEGER,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    warm_up_exercises JSONB DEFAULT '[]',
    main_exercises JSONB DEFAULT '[]',
    cool_down_exercises JSONB DEFAULT '[]',
    estimated_duration_minutes INTEGER,
    estimated_calories_burned INTEGER,
    difficulty_level INTEGER CHECK (difficulty_level >= 1 AND difficulty_level <= 5),
    equipment_needed JSONB DEFAULT '[]',
    instructions TEXT,
    safety_notes TEXT,
    modifications JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Workout Tracking
CREATE TABLE IF NOT EXISTS user_workout_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workout_session_id UUID REFERENCES workout_sessions(id),
    workout_program_id UUID REFERENCES workout_programs(id),
    scheduled_date DATE,
    completed_date DATE,
    duration_minutes INTEGER,
    calories_burned INTEGER,
    perceived_exertion INTEGER CHECK (perceived_exertion >= 1 AND perceived_exertion <= 10),
    mood_before INTEGER CHECK (mood_before >= 1 AND mood_before <= 5),
    mood_after INTEGER CHECK (mood_after >= 1 AND mood_after <= 5),
    exercises_completed JSONB DEFAULT '[]',
    exercises_skipped JSONB DEFAULT '[]',
    modifications_used JSONB DEFAULT '[]',
    notes TEXT,
    injuries_reported JSONB DEFAULT '[]',
    status VARCHAR(50) DEFAULT 'scheduled', -- scheduled, completed, skipped, partial
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- System metrics table
CREATE TABLE IF NOT EXISTS system_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(15,6),
    metric_type VARCHAR(50),
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_key_usage_api_key_id ON api_key_usage(api_key_id);
CREATE INDEX IF NOT EXISTS idx_api_key_usage_timestamp ON api_key_usage(timestamp);
CREATE INDEX IF NOT EXISTS idx_foods_name ON foods(name);
CREATE INDEX IF NOT EXISTS idx_foods_category ON foods(category);
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises(name);
CREATE INDEX IF NOT EXISTS idx_exercises_category ON exercises(category);
CREATE INDEX IF NOT EXISTS idx_meal_plans_user_id ON meal_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_plans_user_id ON workout_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_user_food_logs_user_id ON user_food_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_exercise_logs_user_id ON user_exercise_logs(user_id);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_foods_updated_at BEFORE UPDATE ON foods FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_exercises_updated_at BEFORE UPDATE ON exercises FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_api_keys_updated_at BEFORE UPDATE ON api_keys FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_meal_plans_updated_at BEFORE UPDATE ON meal_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workout_plans_updated_at BEFORE UPDATE ON workout_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_usage_alerts_updated_at BEFORE UPDATE ON usage_alerts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();-- Additi
onal indexes for the new tables
CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);
CREATE INDEX IF NOT EXISTS idx_recipes_country ON recipes(country);
CREATE INDEX IF NOT EXISTS idx_recipes_difficulty ON recipes(difficulty_level);
CREATE INDEX IF NOT EXISTS idx_recipes_prep_time ON recipes(prep_time_minutes);
CREATE INDEX IF NOT EXISTS idx_recipes_rating ON recipes(rating);

CREATE INDEX IF NOT EXISTS idx_health_conditions_category ON health_conditions(category);
CREATE INDEX IF NOT EXISTS idx_health_conditions_icd10 ON health_conditions(icd_10_code);

CREATE INDEX IF NOT EXISTS idx_user_health_complaints_user_id ON user_health_complaints(user_id);
CREATE INDEX IF NOT EXISTS idx_user_health_complaints_type ON user_health_complaints(complaint_type);
CREATE INDEX IF NOT EXISTS idx_user_health_complaints_status ON user_health_complaints(status);

CREATE INDEX IF NOT EXISTS idx_injuries_category ON injuries(category);
CREATE INDEX IF NOT EXISTS idx_injuries_body_part ON injuries(body_part);
CREATE INDEX IF NOT EXISTS idx_injuries_severity ON injuries(severity_level);

CREATE INDEX IF NOT EXISTS idx_user_injuries_user_id ON user_injuries(user_id);
CREATE INDEX IF NOT EXISTS idx_user_injuries_status ON user_injuries(current_status);

CREATE INDEX IF NOT EXISTS idx_nutritional_plans_user_id ON nutritional_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_nutritional_plans_type ON nutritional_plans(plan_type);
CREATE INDEX IF NOT EXISTS idx_nutritional_plans_active ON nutritional_plans(is_active);

CREATE INDEX IF NOT EXISTS idx_medications_class ON medications(drug_class);
CREATE INDEX IF NOT EXISTS idx_medications_category ON medications(category);
CREATE INDEX IF NOT EXISTS idx_medications_generic ON medications(generic_name);

CREATE INDEX IF NOT EXISTS idx_user_medications_user_id ON user_medications(user_id);
CREATE INDEX IF NOT EXISTS idx_user_medications_active ON user_medications(is_active);

CREATE INDEX IF NOT EXISTS idx_vitamins_minerals_type ON vitamins_minerals(type);
CREATE INDEX IF NOT EXISTS idx_vitamins_minerals_category ON vitamins_minerals(category);

CREATE INDEX IF NOT EXISTS idx_user_supplements_user_id ON user_supplements(user_id);
CREATE INDEX IF NOT EXISTS idx_user_supplements_active ON user_supplements(is_active);

CREATE INDEX IF NOT EXISTS idx_workout_programs_type ON workout_programs(program_type);
CREATE INDEX IF NOT EXISTS idx_workout_programs_level ON workout_programs(fitness_level);

CREATE INDEX IF NOT EXISTS idx_workout_sessions_program_id ON workout_sessions(workout_program_id);

CREATE INDEX IF NOT EXISTS idx_user_workout_sessions_user_id ON user_workout_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_workout_sessions_date ON user_workout_sessions(completed_date);
CREATE INDEX IF NOT EXISTS idx_user_workout_sessions_status ON user_workout_sessions(status);

-- Create triggers for updated_at columns on new tables
CREATE TRIGGER update_recipes_updated_at BEFORE UPDATE ON recipes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_health_conditions_updated_at BEFORE UPDATE ON health_conditions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_health_complaints_updated_at BEFORE UPDATE ON user_health_complaints FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_injuries_updated_at BEFORE UPDATE ON injuries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_injuries_updated_at BEFORE UPDATE ON user_injuries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_nutritional_plans_updated_at BEFORE UPDATE ON nutritional_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_medications_updated_at BEFORE UPDATE ON medications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_medications_updated_at BEFORE UPDATE ON user_medications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_vitamins_minerals_updated_at BEFORE UPDATE ON vitamins_minerals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_supplements_updated_at BEFORE UPDATE ON user_supplements FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workout_programs_updated_at BEFORE UPDATE ON workout_programs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workout_sessions_updated_at BEFORE UPDATE ON workout_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_workout_sessions_updated_at BEFORE UPDATE ON user_workout_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();