package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	ID          uint      `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	PasswordHash string   `json:"-" db:"password_hash"`
	FirstName   *string   `json:"first_name" db:"first_name"`
	LastName    *string   `json:"last_name" db:"last_name"`
	Role        string    `json:"role" db:"role"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	LastLogin   *time.Time `json:"last_login" db:"last_login"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	ID                   uint            `json:"id" db:"id"`
	UserID               uint            `json:"user_id" db:"user_id"`
	DateOfBirth          *sql.NullTime   `json:"date_of_birth" db:"date_of_birth"`
	Gender               *string         `json:"gender" db:"gender"`
	Height               *sql.NullFloat64 `json:"height" db:"height"`           // in cm
	Weight               *sql.NullFloat64 `json:"weight" db:"weight"`           // in kg
	ActivityLevel        *string         `json:"activity_level" db:"activity_level"`
	Goal                 *string         `json:"goal" db:"goal"`
	DietaryRestrictions  pq.StringArray  `json:"dietary_restrictions" db:"dietary_restrictions"`
	Allergies            pq.StringArray  `json:"allergies" db:"allergies"`
	PreferredUnits       string          `json:"preferred_units" db:"preferred_units"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// UserPreferences represents user preferences
type UserPreferences struct {
	ID                    uint      `json:"id" db:"id"`
	UserID                uint      `json:"user_id" db:"user_id"`
	Language              string    `json:"language" db:"language"`
	Timezone              string    `json:"timezone" db:"timezone"`
	NotificationsEnabled  bool      `json:"notifications_enabled" db:"notifications_enabled"`
	EmailNotifications    bool      `json:"email_notifications" db:"email_notifications"`
	PushNotifications     bool      `json:"push_notifications" db:"push_notifications"`
	Units                 string    `json:"units" db:"units"`
	DarkMode              bool      `json:"dark_mode" db:"dark_mode"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// NutritionGoal represents nutrition goals for a user
type NutritionGoal struct {
	ID              uint       `json:"id" db:"id"`
	UserID          uint       `json:"user_id" db:"user_id"`
	DailyCalories   *int       `json:"daily_calories" db:"daily_calories"`
	ProteinGrams    *float64   `json:"protein_grams" db:"protein_grams"`
	CarbsGrams      *float64   `json:"carbs_grams" db:"carbs_grams"`
	FatGrams        *float64   `json:"fat_grams" db:"fat_grams"`
	FiberGrams      *float64   `json:"fiber_grams" db:"fiber_grams"`
	SugarGrams      *float64   `json:"sugar_grams" db:"sugar_grams"`
	SodiumMg        *int       `json:"sodium_mg" db:"sodium_mg"`
	WaterMl         *int       `json:"water_ml" db:"water_ml"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	StartDate       *time.Time `json:"start_date" db:"start_date"`
	EndDate         *time.Time `json:"end_date" db:"end_date"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// APIKey represents an API key for a user
type APIKey struct {
	ID           uint       `json:"id" db:"id"`
	UserID       uint       `json:"user_id" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	APIKey       string     `json:"api_key" db:"api_key"`
	Permissions  []string   `json:"permissions" db:"permissions"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"`
	LastUsedAt   *time.Time `json:"last_used_at" db:"last_used_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	ID         uint       `json:"id" db:"id"`
	UserID     uint       `json:"user_id" db:"user_id"`
	TokenHash  string     `json:"-" db:"token_hash"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	IsRevoked  bool       `json:"is_revoked" db:"is_revoked"`
	DeviceInfo *string    `json:"device_info" db:"device_info"`
	IPAddress  *string    `json:"ip_address" db:"ip_address"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// User represents a user response with profile information
type UserResponse struct {
	ID         uint                   `json:"id"`
	Email      string                 `json:"email"`
	FirstName  *string                `json:"first_name"`
	LastName   *string                `json:"last_name"`
	Role       string                 `json:"role"`
	IsVerified bool                   `json:"is_verified"`
	LastLogin  *time.Time             `json:"last_login"`
	CreatedAt  time.Time              `json:"created_at"`
	Profile    *UserProfile           `json:"profile,omitempty"`
	Preferences *UserPreferences      `json:"preferences,omitempty"`
	ActiveGoals []*NutritionGoal      `json:"active_goals,omitempty"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email           string                 `json:"email" validate:"required,email"`
	Password        string                 `json:"password" validate:"required,min=8"`
	FirstName       *string                `json:"first_name,omitempty"`
	LastName        *string                `json:"last_name,omitempty"`
	Profile         *CreateProfileRequest  `json:"profile,omitempty"`
	Preferences     *CreatePreferencesRequest `json:"preferences,omitempty"`
}

// CreateProfileRequest represents a request to create a user profile
type CreateProfileRequest struct {
	DateOfBirth         *string    `json:"date_of_birth,omitempty"`
	Gender              *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Height              *float64   `json:"height,omitempty" validate:"omitempty,min=50,max=300"`
	Weight              *float64   `json:"weight,omitempty" validate:"omitempty,min=20,max=500"`
	ActivityLevel       *string    `json:"activity_level,omitempty" validate:"omitempty,oneof=sedentary light moderate active very_active"`
	Goal                *string    `json:"goal,omitempty" validate:"omitempty,oneof=lose_weight maintain gain_weight gain_muscle"`
	DietaryRestrictions []string   `json:"dietary_restrictions,omitempty"`
	Allergies           []string   `json:"allergies,omitempty"`
	PreferredUnits      *string    `json:"preferred_units,omitempty" validate:"omitempty,oneof=metric imperial"`
}

// CreatePreferencesRequest represents a request to create user preferences
type CreatePreferencesRequest struct {
	Language              *string `json:"language,omitempty" validate:"omitempty,len=2"`
	Timezone              *string `json:"timezone,omitempty"`
	NotificationsEnabled  *bool   `json:"notifications_enabled,omitempty"`
	EmailNotifications    *bool   `json:"email_notifications,omitempty"`
	PushNotifications     *bool   `json:"push_notifications,omitempty"`
	Units                 *string `json:"units,omitempty" validate:"omitempty,oneof=metric imperial"`
	DarkMode              *bool   `json:"dark_mode,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

// UpdateProfileRequest represents a request to update a user profile
type UpdateProfileRequest struct {
	DateOfBirth         *string  `json:"date_of_birth,omitempty" validate:"omitempty"`
	Gender              *string  `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Height              *float64 `json:"height,omitempty" validate:"omitempty,min=50,max=300"`
	Weight              *float64 `json:"weight,omitempty" validate:"omitempty,min=20,max=500"`
	ActivityLevel       *string  `json:"activity_level,omitempty" validate:"omitempty,oneof=sedentary light moderate active very_active"`
	Goal                *string  `json:"goal,omitempty" validate:"omitempty,oneof=lose_weight maintain gain_weight gain_muscle"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
	Allergies           []string `json:"allergies,omitempty"`
	PreferredUnits      *string  `json:"preferred_units,omitempty" validate:"omitempty,oneof=metric imperial"`
}

// UpdatePreferencesRequest represents a request to update user preferences
type UpdatePreferencesRequest struct {
	Language              *string `json:"language,omitempty" validate:"omitempty,len=2"`
	Timezone              *string `json:"timezone,omitempty"`
	NotificationsEnabled  *bool   `json:"notifications_enabled,omitempty"`
	EmailNotifications    *bool   `json:"email_notifications,omitempty"`
	PushNotifications     *bool   `json:"push_notifications,omitempty"`
	Units                 *string `json:"units,omitempty" validate:"omitempty,oneof=metric imperial"`
	DarkMode              *bool   `json:"dark_mode,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents a refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// CreateAPIKeyRequest represents a request to create an API key
type CreateAPIKeyRequest struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
	ExpiresAt   *string  `json:"expires_at,omitempty"`
}

// BMRRequest represents a BMR calculation request
type BMRRequest struct {
	Age          int     `json:"age" validate:"required,min=1,max=120"`
	Gender       string  `json:"gender" validate:"required,oneof=male female"`
	Height       float64 `json:"height" validate:"required,min=50,max=300"`
	Weight       float64 `json:"weight" validate:"required,min=20,max=500"`
	ActivityLevel string  `json:"activity_level" validate:"required,oneof=sedentary light moderate active very_active"`
}

// BMRResponse represents a BMR calculation response
type BMRResponse struct {
	BMR      float64 `json:"bmr"`
	TDEE     float64 `json:"tdee"`
	Formula  string  `json:"formula"`
	Activity string  `json:"activity_level"`
}

// BMIRequest represents a BMI calculation request
type BMIRequest struct {
	Height float64 `json:"height" validate:"required,min=50,max=300"`
	Weight float64 `json:"weight" validate:"required,min=20,max=500"`
	Units  string  `json:"units" validate:"required,oneof=metric imperial"`
}

// BMIResponse represents a BMI calculation response
type BMIResponse struct {
	BMI        float64 `json:"bmi"`
	Category   string  `json:"category"`
	Status     string  `json:"status"`
	Height     float64 `json:"height"`
	Weight     float64 `json:"weight"`
	Units      string  `json:"units"`
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
	}
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName != nil && u.LastName != nil {
		return *u.FirstName + " " + *u.LastName
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	if u.LastName != nil {
		return *u.LastName
	}
	return u.Email
}

// GetDisplayName returns the best display name for the user
func (u *User) GetDisplayName() string {
	if name := u.GetFullName(); name != u.Email {
		return name
	}
	return u.Email
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission string) bool {
	// Admin has all permissions
	if u.Role == "admin" {
		return true
	}
	
	// Define user permissions based on role
	permissions := map[string][]string{
		"user": {"read_own_profile", "update_own_profile", "create_own_data", "read_own_data"},
		"premium": {"read_own_profile", "update_own_profile", "create_own_data", "read_own_data", "advanced_analytics"},
		"admin": nil, // Admin has all permissions
	}
	
	userPerms, exists := permissions[u.Role]
	if !exists {
		return false
	}
	
	for _, perm := range userPerms {
		if perm == permission {
			return true
		}
	}
	
	return false
}

// IsPremium checks if user has premium access
func (u *User) IsPremium() bool {
	return u.Role == "premium" || u.Role == "admin"
}

// ParseActivityLevel converts activity level string to multiplier
func ParseActivityLevel(level string) float64 {
	multipliers := map[string]float64{
		"sedentary":     1.2,
		"light":         1.375,
		"moderate":      1.55,
		"active":        1.725,
		"very_active":   1.9,
	}
	
	if multiplier, exists := multipliers[level]; exists {
		return multiplier
	}
	return 1.2 // Default to sedentary
}

// CalculateBMR calculates Basal Metabolic Rate using Mifflin-St Jeor equation
func CalculateBMR(req BMRRequest) BMRResponse {
	var bmr float64
	
	if req.Gender == "male" {
		bmr = 10*req.Weight + 6.25*req.Height - 5*float64(req.Age) + 5
	} else {
		bmr = 10*req.Weight + 6.25*req.Height - 5*float64(req.Age) - 161
	}
	
	activityMultiplier := ParseActivityLevel(req.ActivityLevel)
	tdee := bmr * activityMultiplier
	
	return BMRResponse{
		BMR:      bmr,
		TDEE:     tdee,
		Formula:  "Mifflin-St Jeor",
		Activity: req.ActivityLevel,
	}
}

// CalculateBMI calculates Body Mass Index
func CalculateBMI(req BMIRequest) BMIResponse {
	var bmi float64
	var height float64
	var weight float64
	
	if req.Units == "imperial" {
		// Convert to metric
		height = req.Height * 2.54 // inches to cm
		weight = req.Weight * 0.453592 // lbs to kg
	} else {
		height = req.Height
		weight = req.Weight
	}
	
	// Calculate BMI (weight in kg / height in meters squared)
	heightInMeters := height / 100
	bmi = weight / (heightInMeters * heightInMeters)
	
	// Determine category
	var category, status string
	switch {
	case bmi < 18.5:
		category = "underweight"
		status = "Underweight"
	case bmi < 25:
		category = "normal"
		status = "Normal weight"
	case bmi < 30:
		category = "overweight"
		status = "Overweight"
	default:
		category = "obese"
		status = "Obese"
	}
	
	return BMIResponse{
		BMI:      bmi,
		Category: category,
		Status:   status,
		Height:   req.Height,
		Weight:   req.Weight,
		Units:    req.Units,
	}
}
