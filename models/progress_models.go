package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ProgressPhoto represents a user's progress photo
type ProgressPhoto struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	Date           time.Time `json:"date" db:"date"`
	PhotoURL       string    `json:"photo_url" db:"photo_url"`
	ThumbnailURL   string    `json:"thumbnail_url" db:"thumbnail_url"`
	Weight         float64   `json:"weight" db:"weight"`
	Notes          string    `json:"notes" db:"notes"`
	Visibility     string    `json:"visibility" db:"visibility"` // private, coach, public
	Tags           string    `json:"tags" db:"tags"`             // JSON array of tags
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// BodyMeasurement represents a user's body measurements
type BodyMeasurement struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	Date            time.Time `json:"date" db:"date"`
	Weight          float64   `json:"weight" db:"weight"`
	BodyFat         float64   `json:"body_fat" db:"body_fat"`
	MuscleMass      float64   `json:"muscle_mass" db:"muscle_mass"`
	Waist           float64   `json:"waist" db:"waist"`
	Chest           float64   `json:"chest" db:"chest"`
	Arms            float64   `json:"arms" db:"arms"`
	Thighs          float64   `json:"thighs" db:"thighs"`
	Hips            float64   `json:"hips" db:"hips"`
	Shoulders       float64   `json:"shoulders" db:"shoulders"`
	Calves          float64   `json:"calves" db:"calves"`
	Forearms        float64   `json:"forearms" db:"forearms"`
	Neck            float64   `json:"neck" db:"neck"`
	BMI             float64   `json:"bmi" db:"bmi"`
	Notes           string    `json:"notes" db:"notes"`
	MeasurementUnit string    `json:"measurement_unit" db:"measurement_unit"` // metric, imperial
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Milestone represents a user's achievement milestone
type Milestone struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`              // weight_loss, weight_gain, muscle_gain, strength, consistency, custom
	TargetValue float64   `json:"target_value" db:"target_value"`
	CurrentValue float64  `json:"current_value" db:"current_value"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	TargetDate  time.Time `json:"target_date" db:"target_date"`
	AchievedAt  *time.Time `json:"achieved_at" db:"achieved_at"`
	Status      string    `json:"status" db:"status"`           // active, achieved, paused, cancelled
	BadgeURL    string    `json:"badge_url" db:"badge_url"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WeightGoal represents a user's weight goal
type WeightGoal struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	StartingWeight float64   `json:"starting_weight" db:"starting_weight"`
	TargetWeight   float64   `json:"target_weight" db:"target_weight"`
	CurrentWeight  float64   `json:"current_weight" db:"current_weight"`
	WeeklyGoal     float64   `json:"weekly_goal" db:"weekly_goal"` // kg per week
	ActivityLevel  string    `json:"activity_level" db:"activity_level"`
	TargetDate     time.Time `json:"target_date" db:"target_date"`
	Status         string    `json:"status" db:"status"`          // active, achieved, paused
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ProgressSummary represents a summary of user's progress
type ProgressSummary struct {
	Period                    string  `json:"period"`
	StartDate                 time.Time `json:"start_date"`
	EndDate                   time.Time `json:"end_date"`
	WeightChange              float64 `json:"weight_change"`
	BodyFatChange             float64 `json:"body_fat_change"`
	MuscleMassChange          float64 `json:"muscle_mass_change"`
	WaistChange               float64 `json:"waist_change"`
	WeightTrend               string  `json:"weight_trend"` // losing, gaining, stable
	AverageWeeklyWeightChange float64 `json:"average_weekly_weight_change"`
	MeasurementsCount         int     `json:"measurements_count"`
	PhotosCount               int     `json:"photos_count"`
	MilestonesAchieved        int     `json:"milestones_achieved"`
	ProgressPercentage        float64 `json:"progress_percentage"`
}

// WeightTrendData represents weight trend over time
type WeightTrendData struct {
	Date   time.Time `json:"date"`
	Weight float64   `json:"weight"`
	BMI    float64   `json:"bmi"`
	BodyFat float64  `json:"body_fat,omitempty"`
}

// Request DTOs
type CreateProgressPhotoRequest struct {
	Date         time.Time `json:"date" validate:"required"`
	PhotoURL     string    `json:"photo_url" validate:"required,url"`
	Weight       float64   `json:"weight" validate:"min=0,max=500"`
	Notes        string    `json:"notes" validate:"max=1000"`
	Visibility   string    `json:"visibility" validate:"oneof=private coach public"`
	Tags         []string  `json:"tags"`
}

type UpdateProgressPhotoRequest struct {
	Date       *time.Time `json:"date"`
	PhotoURL   *string    `json:"photo_url" validate:"omitempty,url"`
	Weight     *float64   `json:"weight" validate:"omitempty,min=0,max=500"`
	Notes      *string    `json:"notes" validate:"omitempty,max=1000"`
	Visibility *string    `json:"visibility" validate:"omitempty,oneof=private coach public"`
	Tags       []string   `json:"tags"`
}

type CreateBodyMeasurementRequest struct {
	Date            time.Time `json:"date" validate:"required"`
	Weight          float64   `json:"weight" validate:"min=0,max=500"`
	BodyFat         float64   `json:"body_fat" validate:"min=0,max=100"`
	MuscleMass      float64   `json:"muscle_mass" validate:"min=0,max=500"`
	Waist           float64   `json:"waist" validate:"min=0,max=500"`
	Chest           float64   `json:"chest" validate:"min=0,max=500"`
	Arms            float64   `json:"arms" validate:"min=0,max=500"`
	Thighs          float64   `json:"thighs" validate:"min=0,max=500"`
	Hips            float64   `json:"hips" validate:"min=0,max=500"`
	Shoulders       float64   `json:"shoulders" validate:"min=0,max=500"`
	Calves          float64   `json:"calves" validate:"min=0,max=500"`
	Forearms        float64   `json:"forearms" validate:"min=0,max=500"`
	Neck            float64   `json:"neck" validate:"min=0,max=500"`
	Notes           string    `json:"notes" validate:"max=1000"`
	MeasurementUnit string    `json:"measurement_unit" validate:"oneof=metric imperial"`
}

type UpdateBodyMeasurementRequest struct {
	Date            *time.Time `json:"date"`
	Weight          *float64   `json:"weight" validate:"omitempty,min=0,max=500"`
	BodyFat         *float64   `json:"body_fat" validate:"omitempty,min=0,max=100"`
	MuscleMass      *float64   `json:"muscle_mass" validate:"omitempty,min=0,max=500"`
	Waist           *float64   `json:"waist" validate:"omitempty,min=0,max=500"`
	Chest           *float64   `json:"chest" validate:"omitempty,min=0,max=500"`
	Arms            *float64   `json:"arms" validate:"omitempty,min=0,max=500"`
	Thighs          *float64   `json:"thighs" validate:"omitempty,min=0,max=500"`
	Hips            *float64   `json:"hips" validate:"omitempty,min=0,max=500"`
	Shoulders       *float64   `json:"shoulders" validate:"omitempty,min=0,max=500"`
	Calves          *float64   `json:"calves" validate:"omitempty,min=0,max=500"`
	Forearms        *float64   `json:"forearms" validate:"omitempty,min=0,max=500"`
	Neck            *float64   `json:"neck" validate:"omitempty,min=0,max=500"`
	Notes           *string    `json:"notes" validate:"omitempty,max=1000"`
	MeasurementUnit *string    `json:"measurement_unit" validate:"omitempty,oneof=metric imperial"`
}

type CreateMilestoneRequest struct {
	Title       string    `json:"title" validate:"required,max=200"`
	Description string    `json:"description" validate:"max=1000"`
	Type        string    `json:"type" validate:"required,oneof=weight_loss weight_gain muscle_gain strength consistency custom"`
	TargetValue float64   `json:"target_value" validate:"required,min=0"`
	TargetDate  time.Time `json:"target_date" validate:"required"`
	IsPublic    bool      `json:"is_public"`
}

type UpdateMilestoneRequest struct {
	Title       *string    `json:"title" validate:"omitempty,max=200"`
	Description *string    `json:"description" validate:"omitempty,max=1000"`
	Type        *string    `json:"type" validate:"omitempty,oneof=weight_loss weight_gain muscle_gain strength consistency custom"`
	TargetValue *float64   `json:"target_value" validate:"omitempty,min=0"`
	TargetDate  *time.Time `json:"target_date"`
	Status      *string    `json:"status" validate:"omitempty,oneof=active achieved paused cancelled"`
	BadgeURL    *string    `json:"badge_url" validate:"omitempty,url"`
	IsPublic    *bool      `json:"is_public"`
}

type CreateWeightGoalRequest struct {
	StartingWeight float64   `json:"starting_weight" validate:"required,min=0,max=500"`
	TargetWeight   float64   `json:"target_weight" validate:"required,min=0,max=500"`
	WeeklyGoal     float64   `json:"weekly_goal" validate:"required,min=-2,max=2"`
	ActivityLevel  string    `json:"activity_level" validate:"required,oneof=sedentary light moderate active very_active"`
	TargetDate     time.Time `json:"target_date" validate:"required"`
}

type UpdateWeightGoalRequest struct {
	CurrentWeight *float64   `json:"current_weight" validate:"omitempty,min=0,max=500"`
	WeeklyGoal    *float64   `json:"weekly_goal" validate:"omitempty,min=-2,max=2"`
	ActivityLevel *string    `json:"activity_level" validate:"omitempty,oneof=sedentary light moderate active very_active"`
	TargetDate    *time.Time `json:"target_date"`
	Status        *string    `json:"status" validate:"omitempty,oneof=active achieved paused"`
}

// Response DTOs
type ProgressPhotoResponse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Date         time.Time `json:"date"`
	PhotoURL     string    `json:"photo_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Weight       float64   `json:"weight"`
	Notes        string    `json:"notes"`
	Visibility   string    `json:"visibility"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BodyMeasurementResponse struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Date            time.Time `json:"date"`
	Weight          float64   `json:"weight"`
	BodyFat         float64   `json:"body_fat"`
	MuscleMass      float64   `json:"muscle_mass"`
	Waist           float64   `json:"waist"`
	Chest           float64   `json:"chest"`
	Arms            float64   `json:"arms"`
	Thighs          float64   `json:"thighs"`
	Hips            float64   `json:"hips"`
	Shoulders       float64   `json:"shoulders"`
	Calves          float64   `json:"calves"`
	Forearms        float64   `json:"forearms"`
	Neck            float64   `json:"neck"`
	BMI             float64   `json:"bmi"`
	Notes           string    `json:"notes"`
	MeasurementUnit string    `json:"measurement_unit"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MilestoneResponse struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Type         string     `json:"type"`
	TargetValue  float64    `json:"target_value"`
	CurrentValue float64    `json:"current_value"`
	StartDate    time.Time  `json:"start_date"`
	TargetDate   time.Time  `json:"target_date"`
	AchievedAt   *time.Time `json:"achieved_at"`
	Status       string     `json:"status"`
	BadgeURL     string     `json:"badge_url"`
	IsPublic     bool       `json:"is_public"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type WeightGoalResponse struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	StartingWeight float64   `json:"starting_weight"`
	TargetWeight   float64   `json:"target_weight"`
	CurrentWeight  float64   `json:"current_weight"`
	WeeklyGoal     float64   `json:"weekly_goal"`
	ActivityLevel  string    `json:"activity_level"`
	TargetDate     time.Time `json:"target_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Query filters
type ProgressPhotoFilter struct {
	UserID     int64     `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Visibility string    `json:"visibility"`
	Tags       []string  `json:"tags"`
}

type BodyMeasurementFilter struct {
	UserID          int64     `json:"user_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	MeasurementUnit string    `json:"measurement_unit"`
}

type MilestoneFilter struct {
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// Database helper types
type TagsArray []string

func (t TagsArray) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

func (t *TagsArray) Scan(value interface{}) error {
	if value == nil {
		*t = TagsArray{}
		return nil
	}
	
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), t)
	case []byte:
		return json.Unmarshal(v, t)
	default:
		return fmt.Errorf("cannot scan %T into TagsArray", value)
	}
}

// ProgressAnalytics represents detailed progress analytics
type ProgressAnalytics struct {
	WeightAnalysis    WeightAnalysis    `json:"weight_analysis"`
	MeasurementAnalysis MeasurementAnalysis `json:"measurement_analysis"`
	PhotoAnalysis     PhotoAnalysis     `json:"photo_analysis"`
	MilestoneAnalysis MilestoneAnalysis `json:"milestone_analysis"`
	Trends            []TrendData       `json:"trends"`
}

type WeightAnalysis struct {
	CurrentWeight     float64   `json:"current_weight"`
	StartingWeight    float64   `json:"starting_weight"`
	TargetWeight      float64   `json:"target_weight"`
	TotalChange       float64   `json:"total_change"`
	WeeklyAverage     float64   `json:"weekly_average"`
	PredictedDate     time.Time `json:"predicted_date"`
	OnTrack           bool      `json:"on_track"`
	BMI               float64   `json:"bmi"`
	BMICategory       string    `json:"bmi_category"`
}

type MeasurementAnalysis struct {
	MostRecentMeasurement BodyMeasurement `json:"most_recent_measurement"`
	PreviousMeasurement    BodyMeasurement `json:"previous_measurement"`
	Changes               map[string]float64 `json:"changes"`
	ProgressRate          map[string]float64 `json:"progress_rate"`
}

type PhotoAnalysis struct {
	TotalPhotos      int       `json:"total_photos"`
	FirstPhotoDate   time.Time `json:"first_photo_date"`
	MostRecentDate   time.Time `json:"most_recent_date"`
	PhotoFrequency   float64   `json:"photo_frequency"`
	WeightDifference float64   `json:"weight_difference"`
	VisualProgress   string    `json:"visual_progress"`
}

type MilestoneAnalysis struct {
	TotalMilestones     int     `json:"total_milestones"`
	AchievedMilestones  int     `json:"achieved_milestones"`
	ActiveMilestones    int     `json:"active_milestones"`
	CompletionRate      float64 `json:"completion_rate"`
	RecentAchievements  []Milestone `json:"recent_achievements"`
	UpcomingDeadlines   []Milestone `json:"upcoming_deadlines"`
}

type TrendData struct {
	Date         time.Time `json:"date"`
	Weight       float64   `json:"weight"`
	BodyFat      float64   `json:"body_fat"`
	MuscleMass   float64   `json:"muscle_mass"`
	Waist        float64   `json:"waist"`
	Chest        float64   `json:"chest"`
	Arms         float64   `json:"arms"`
	Trend        string    `json:"trend"`
}
// Additional model types needed for analytics
type WeightTrendPoint struct {
	Date   time.Time `json:"date"`
	Weight float64   `json:"weight"`
	BMI    float64   `json:"bmi"`
}

type BodyFatTrendPoint struct {
	Date    time.Time `json:"date"`
	BodyFat float64   `json:"body_fat"`
}

type MeasurementStats struct {
	TotalMeasurements int     `json:"total_measurements"`
	AverageWeight     float64 `json:"average_weight"`
	MinWeight         float64 `json:"min_weight"`
	MaxWeight         float64 `json:"max_weight"`
	WeightChange      float64 `json:"weight_change"`
	AverageBodyFat    float64 `json:"average_body_fat"`
	BodyFatChange     float64 `json:"body_fat_change"`
}

type MeasurementComparison struct {
	CurrentWeight     float64 `json:"current_weight"`
	PreviousWeight    float64 `json:"previous_weight"`
	WeightDifference  float64 `json:"weight_difference"`
	WeightChangePercent float64 `json:"weight_change_percent"`
	CurrentBodyFat    float64 `json:"current_body_fat"`
	PreviousBodyFat   float64 `json:"previous_body_fat"`
	BodyFatDifference float64 `json:"body_fat_difference"`
	BodyFatChangePercent float64 `json:"body_fat_change_percent"`
	Period            string  `json:"period"`
	ComparisonDate    time.Time `json:"comparison_date"`
}

// PersonalRecordStats represents personal record statistics
type PersonalRecordStats struct {
	UserID              int64     `json:"user_id"`
	TotalRecords        int       `json:"total_records"`
	RecentRecords       int       `json:"recent_records"`
	MostRecentRecord    time.Time `json:"most_recent_record"`
	ImprovementRate     float64   `json:"improvement_rate"`
	RecordsByType       map[string]int `json:"records_by_type"`
	ProgressThisMonth   int       `json:"progress_this_month"`
	ProgressThisYear    int       `json:"progress_this_year"`
}

// WeightProgressAnalytics represents weight progress analytics
type WeightProgressAnalytics struct {
	UserID              int64     `json:"user_id"`
	Period              string    `json:"period"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	StartingWeight      float64   `json:"starting_weight"`
	CurrentWeight       float64   `json:"current_weight"`
	TargetWeight        float64   `json:"target_weight"`
	TotalChange         float64   `json:"total_change"`
	WeeklyAverage       float64   `json:"weekly_average"`
	PredictedDate       time.Time `json:"predicted_date"`
	OnTrack             bool      `json:"on_track"`
	Trend               string    `json:"trend"`
	BMI                 float64   `json:"bmi"`
	BMICategory         string    `json:"bmi_category"`
	ProgressPercentage  float64   `json:"progress_percentage"`
}

// BodyCompositionAnalytics represents body composition analytics
type BodyCompositionAnalytics struct {
	UserID              int64     `json:"user_id"`
	Period              string    `json:"period"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	WeightChange        float64   `json:"weight_change"`
	BodyFatChange       float64   `json:"body_fat_change"`
	MuscleMassChange    float64   `json:"muscle_mass_change"`
	WeightTrend         string    `json:"weight_trend"`
	BodyFatTrend        string    `json:"body_fat_trend"`
	MuscleMassTrend     string    `json:"muscle_mass_trend"`
	CompositionRatio    float64   `json:"composition_ratio"`
	HealthScore         float64   `json:"health_score"`
}

// MilestoneAnalytics represents milestone analytics
type MilestoneAnalytics struct {
	UserID              int64     `json:"user_id"`
	TotalMilestones     int       `json:"total_milestones"`
	AchievedMilestones  int       `json:"achieved_milestones"`
	ActiveMilestones    int       `json:"active_milestones"`
	CompletionRate      float64   `json:"completion_rate"`
	RecentAchievements  []Milestone `json:"recent_achievements"`
	UpcomingDeadlines   []Milestone `json:"upcoming_deadlines"`
	AverageTimeToAchieve int       `json:"average_time_to_achieve"`
	MostAchievedType    string    `json:"most_achieved_type"`
	ProgressThisMonth   int       `json:"progress_this_month"`
}

// PhotoAnalytics represents photo analytics
type PhotoAnalytics struct {
	UserID              int64     `json:"user_id"`
	TotalPhotos         int       `json:"total_photos"`
	FirstPhotoDate      time.Time `json:"first_photo_date"`
	MostRecentDate      time.Time `json:"most_recent_date"`
	PhotoFrequency      float64   `json:"photo_frequency"`
	WeightDifference    float64   `json:"weight_difference"`
	VisualProgress      string    `json:"visual_progress"`
	PhotosByMonth       map[string]int `json:"photos_by_month"`
	ConsistencyScore    float64   `json:"consistency_score"`
	TagsDistribution    map[string]int `json:"tags_distribution"`
}

// ProgressPredictions represents progress predictions
type ProgressPredictions struct {
	UserID              int64     `json:"user_id"`
	PredictionType      string    `json:"prediction_type"`
	TargetDate          time.Time `json:"target_date"`
	PredictedWeight     float64   `json:"predicted_weight"`
	PredictedBodyFat    float64   `json:"predicted_body_fat"`
	PredictedMuscleMass float64   `json:"predicted_muscle_mass"`
	Confidence          float64   `json:"confidence"`
	TimeToGoal          int       `json:"time_to_goal"`
	Recommendations     []string  `json:"recommendations"`
	RiskFactors         []string  `json:"risk_factors"`
}

// ConsistencyAnalytics represents consistency analytics
type ConsistencyAnalytics struct {
	UserID              int64     `json:"user_id"`
	OverallScore        float64   `json:"overall_score"`
	WorkoutConsistency  float64   `json:"workout_consistency"`
	MealConsistency     float64   `json:"meal_consistency"`
	MeasurementConsistency float64 `json:"measurement_consistency"`
	PhotoConsistency    float64   `json:"photo_consistency"`
	StreakDays          int       `json:"streak_days"`
	LongestStreak       int       `json:"longest_streak"`
	MissedDays          int       `json:"missed_days"`
	PerfectWeeks        int       `json:"perfect_weeks"`
	MonthlyTrends       map[string]float64 `json:"monthly_trends"`
}

// AchievementAnalytics represents achievement analytics
type AchievementAnalytics struct {
	UserID              int64     `json:"user_id"`
	TotalAchievements   int       `json:"total_achievements"`
	AchievementsByType  map[string]int `json:"achievements_by_type"`
	RecentAchievements  []string  `json:"recent_achievements"`
	ProgressToNext      float64   `json:"progress_to_next"`
	NextAchievement     string    `json:"next_achievement"`
	AchievementRate     float64   `json:"achievement_rate"`
	MilestoneProgress   map[string]float64 `json:"milestone_progress"`
	BadgeCount          int       `json:"badge_count"`
}

// WeightGoalProgress represents weight goal progress tracking
type WeightGoalProgress struct {
	UserID              int64     `json:"user_id"`
	GoalID              int64     `json:"goal_id"`
	CurrentWeight       float64   `json:"current_weight"`
	TargetWeight        float64   `json:"target_weight"`
	StartingWeight      float64   `json:"starting_weight"`
	ProgressPercentage  float64   `json:"progress_percentage"`
	RemainingWeight     float64   `json:"remaining_weight"`
	WeeklyProgress      float64   `json:"weekly_progress"`
	EstimatedCompletion time.Time `json:"estimated_completion"`
	OnTrack             bool      `json:"on_track"`
	DaysToGoal          int       `json:"days_to_goal"`
	ProgressRate        float64   `json:"progress_rate"`
}

// WeightGoalStats represents weight goal statistics
type WeightGoalStats struct {
	UserID              int64     `json:"user_id"`
	ActiveGoals         int       `json:"active_goals"`
	CompletedGoals      int       `json:"completed_goals"`
	AbandonedGoals      int       `json:"abandoned_goals"`
	TotalWeightLost     float64   `json:"total_weight_lost"`
	TotalWeightGained   float64   `json:"total_weight_gained"`
	AverageTimeToGoal   int       `json:"average_time_to_goal"`
	SuccessRate         float64   `json:"success_rate"`
	CurrentStreak       int       `json:"current_streak"`
	LongestStreak       int       `json:"longest_streak"`
	GoalCompletionHistory []GoalCompletion `json:"goal_completion_history"`
}

// GoalCompletion represents individual goal completion data
type GoalCompletion struct {
	GoalID              int64     `json:"goal_id"`
	GoalType            string    `json:"goal_type"`
	StartDate           time.Time `json:"start_date"`
	CompletionDate      time.Time `json:"completion_date"`
	Duration            int       `json:"duration"`
	Success             bool      `json:"success"`
	FinalWeight         float64   `json:"final_weight"`
	WeightDifference    float64   `json:"weight_difference"`
}
