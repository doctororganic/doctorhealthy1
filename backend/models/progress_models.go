package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// BodyMeasurement represents a user's body measurement entry
type BodyMeasurement struct {
	ID               uint      `json:"id" db:"id"`
	UserID           uint      `json:"user_id" db:"user_id"`
	Date             time.Time `json:"date" db:"date"`
	Weight           *float64  `json:"weight,omitempty" db:"weight"`
	BodyFat          *float64  `json:"body_fat,omitempty" db:"body_fat"`
	MuscleMass       *float64  `json:"muscle_mass,omitempty" db:"muscle_mass"`
	Waist            *float64  `json:"waist,omitempty" db:"waist"`
	Chest            *float64  `json:"chest,omitempty" db:"chest"`
	Arms             *float64  `json:"arms,omitempty" db:"arms"`
	Thighs           *float64  `json:"thighs,omitempty" db:"thighs"`
	Hips             *float64  `json:"hips,omitempty" db:"hips"`
	Calories         *float64  `json:"calories,omitempty" db:"calories"`
	ActivityLevel    *string   `json:"activity_level,omitempty" db:"activity_level"`
	Notes            *string   `json:"notes,omitempty" db:"notes"`
	Photos           PhotoList `json:"photos,omitempty" db:"photos"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// PhotoList is a custom type for handling JSON arrays of photo URLs
type PhotoList []string

// Value implements the driver.Valuer interface for PhotoList
func (p PhotoList) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for PhotoList
func (p *PhotoList) Scan(value interface{}) error {
	if value == nil {
		*p = PhotoList{}
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return fmt.Errorf("cannot scan %T into PhotoList", value)
	}
}

// WeightTrendPoint represents a data point for weight trends
type WeightTrendPoint struct {
	Date   time.Time `json:"date"`
	Weight float64   `json:"weight"`
}

// BodyFatTrendPoint represents a data point for body fat trends
type BodyFatTrendPoint struct {
	Date    time.Time `json:"date"`
	BodyFat float64   `json:"body_fat"`
}

// MeasurementStats represents statistical analysis of measurements
type MeasurementStats struct {
	Current       float64   `json:"current"`
	Previous      float64   `json:"previous"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
	Min           float64   `json:"min"`
	Max           float64   `json:"max"`
	Average       float64   `json:"average"`
	Trend         string    `json:"trend"` // "increasing", "decreasing", "stable"
	Period        string    `json:"period"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
}

// MeasurementComparison represents comparison between two measurement sets
type MeasurementComparison struct {
	CurrentMeasurement BodyMeasurement     `json:"current_measurement"`
	PreviousMeasurement *BodyMeasurement   `json:"previous_measurement,omitempty"`
	Differences        map[string]float64  `json:"differences"`
	PercentChanges     map[string]float64  `json:"percent_changes"`
	TimeDifference     int                 `json:"time_difference_days"`
	Analysis           map[string]string   `json:"analysis"`
}

// ProgressPhoto represents a progress photo entry

// []string is a custom type for handling JSON arrays of tags

// Scan implements the sql.Scanner interface for []string

// Milestone represents a user's milestone achievements
type Milestone struct {
	ID          uint      `json:"id" db:"id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	Type        string    `json:"type" db:"type"` // "weight", "measurement", "workout", "custom"
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	TargetValue *float64  `json:"target_value,omitempty" db:"target_value"`
	AchievedAt  time.Time `json:"achieved_at" db:"achieved_at"`
	PhotoURL    *string   `json:"photo_url,omitempty" db:"photo_url"`
	IsAchieved  bool      `json:"is_achieved" db:"is_achieved"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WeightGoal represents a user's weight goal
type WeightGoal struct {
	ID             uint       `json:"id" db:"id"`
	UserID         uint       `json:"user_id" db:"user_id"`
	TargetWeight   float64    `json:"target_weight" db:"target_weight"`
	CurrentWeight  float64    `json:"current_weight" db:"current_weight"`
	StartDate      time.Time  `json:"start_date" db:"start_date"`
	TargetDate     *time.Time `json:"target_date,omitempty" db:"target_date"`
	WeeklyGoal     float64    `json:"weekly_goal" db:"weekly_goal"` // kg per week, can be negative for weight loss
	ActivityLevel  string     `json:"activity_level" db:"activity_level"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	AchievedAt     *time.Time `json:"achieved_at,omitempty" db:"achieved_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// ProgressSummary represents a summary of user's overall progress
type ProgressSummary struct {
	UserID              uint                `json:"user_id"`
	CurrentWeight       *float64            `json:"current_weight,omitempty"`
	WeightChange        float64             `json:"weight_change"`
	WeightChangePercent float64             `json:"weight_change_percent"`
	BodyFatChange       *float64            `json:"body_fat_change,omitempty"`
	BodyFatChangePercent *float64           `json:"body_fat_change_percent,omitempty"`
	PhotosCount         int                 `json:"photos_count"`
	MilestonesCount     int                 `json:"milestones_count"`
	ActiveGoals         int                 `json:"active_goals"`
	RecentMeasurement   *BodyMeasurement    `json:"recent_measurement,omitempty"`
	RecentPhotos        []ProgressPhoto     `json:"recent_photos,omitempty"`
	UpcomingMilestones  []Milestone         `json:"upcoming_milestones,omitempty"`
	LastUpdated         time.Time           `json:"last_updated"`
}

// PersonalRecord represents a user's personal record for an exercise
type PersonalRecord struct {
	ID           uint      `json:"id" db:"id"`
	UserID       uint      `json:"user_id" db:"user_id"`
	ExerciseID   uint      `json:"exercise_id" db:"exercise_id"`
	RecordType   string    `json:"record_type" db:"record_type"` // "weight", "reps", "time", "distance"
	Value        float64   `json:"value" db:"value"`
	Unit         string    `json:"unit" db:"unit"`
	Date         time.Time `json:"date" db:"date"`
	Notes        *string   `json:"notes,omitempty" db:"notes"`
	WorkoutLogID *uint     `json:"workout_log_id,omitempty" db:"workout_log_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// PersonalRecordType represents the type of personal record
type PersonalRecordType string

const (
	RecordWeight    PersonalRecordType = "weight"
	RecordReps      PersonalRecordType = "reps"
	RecordTime      PersonalRecordType = "time"
	RecordDistance  PersonalRecordType = "distance"
)

// ListPersonalRecordsRequest represents a request to list personal records
type ListPersonalRecordsRequest struct {
	ExerciseID *uint    `json:"exercise_id,omitempty"`
	RecordType *string  `json:"record_type,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Page       int       `json:"page,omitempty" validate:"omitempty,min=1"`
	Limit      int       `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
}

// PersonalRecordListResponse represents a response for personal record listing
type PersonalRecordListResponse struct {
	Records []PersonalRecord `json:"records"`
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
	HasNext bool             `json:"has_next"`
}

// PersonalRecordWithExercise represents a personal record with exercise details
type PersonalRecordWithExercise struct {
	PersonalRecord
	Exercise Exercise `json:"exercise"`
}

// MilestoneProgress represents progress towards a milestone
type MilestoneProgress struct {
	MilestoneID      uint    `json:"milestone_id"`
	MilestoneTitle   string  `json:"milestone_title"`
	CurrentValue     float64 `json:"current_value"`
	TargetValue      float64 `json:"target_value"`
	ProgressPercent  float64 `json:"progress_percent"`
	IsCompleted      bool    `json:"is_completed"`
	EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
}

// MilestoneStats represents statistics for milestones
type MilestoneStats struct {
	TotalMilestones     int                     `json:"total_milestones"`
	CompletedMilestones int                     `json:"completed_milestones"`
	InProgressMilestones int                    `json:"in_progress_milestones"`
	CompletionRate      float64                 `json:"completion_rate"`
	RecentAchievements  []Milestone             `json:"recent_achievements"`
	UpcomingMilestones  []MilestoneProgress     `json:"upcoming_milestones"`
	ProgressByType      map[string]MilestoneStats `json:"progress_by_type"`
}

// WeightProgress represents weight progress over time
type WeightProgress struct {
	CurrentWeight     float64             `json:"current_weight"`
	TargetWeight      float64             `json:"target_weight"`
	WeightChange      float64             `json:"weight_change"`
	WeightChangeRate  float64             `json:"weight_change_rate"` // kg per week
	ProgressPercent   float64             `json:"progress_percent"`
	EstimatedCompletion *time.Time        `json:"estimated_completion,omitempty"`
	TrendPoints       []WeightTrendPoint  `json:"trend_points"`
	IsOnTrack         bool                `json:"is_on_track"`
}

// MeasurementTrend represents measurement trends over time
type MeasurementTrend struct {
	MeasurementType string              `json:"measurement_type"` // "weight", "body_fat", "muscle_mass", etc.
	Current         float64             `json:"current"`
	Previous        float64             `json:"previous"`
	Change          float64             `json:"change"`
	ChangePercent   float64             `json:"change_percent"`
	Trend           string              `json:"trend"` // "increasing", "decreasing", "stable"
	TrendPoints     []WeightTrendPoint  `json:"trend_points"` // Using WeightTrendPoint struct for any measurement
	Unit            string              `json:"unit"`
}
