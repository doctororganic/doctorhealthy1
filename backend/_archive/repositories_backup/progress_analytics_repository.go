package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nutrition-platform/models"
)

type ProgressAnalyticsRepository struct {
	db *sql.DB
}

func NewProgressAnalyticsRepository(db *sql.DB) *ProgressAnalyticsRepository {
	return &ProgressAnalyticsRepository{db: db}
}

// GetOverallProgressSummary provides a comprehensive overview of user progress
func (r *ProgressAnalyticsRepository) GetOverallProgressSummary(ctx context.Context, userID int64) (*models.ProgressSummary, error) {
	query := `
		WITH latest_measurements AS (
			SELECT * FROM body_measurements 
			WHERE user_id = $1 
			ORDER BY measurement_date DESC 
			LIMIT 1
		),
		first_measurements AS (
			SELECT * FROM body_measurements 
			WHERE user_id = $1 
			ORDER BY measurement_date ASC 
			LIMIT 1
		),
		active_goals AS (
			SELECT COUNT(*) as count FROM weight_goals 
			WHERE user_id = $1 AND is_active = true
		),
		achieved_milestones AS (
			SELECT COUNT(*) as count FROM milestones 
			WHERE user_id = $1 AND is_achieved = true
		),
		total_photos AS (
			SELECT COUNT(*) as count FROM progress_photos 
			WHERE user_id = $1
		)
		SELECT 
			lm.weight as current_weight,
			fm.weight as starting_weight,
			lm.body_fat_percentage as current_body_fat,
			fm.body_fat_percentage as starting_body_fat,
			lm.measurement_date as last_measurement_date,
			EXTRACT(DAYS FROM lm.measurement_date - fm.measurement_date) as days_tracked,
			ag.count as active_goals,
			am.count as achieved_milestones,
			tp.count as total_photos,
			CASE 
				WHEN fm.weight > 0 AND lm.weight > 0 THEN 
					((fm.weight - lm.weight) / fm.weight) * 100
				ELSE 0
			END as weight_change_percentage
		FROM latest_measurements lm
		CROSS JOIN first_measurements fm
		CROSS JOIN active_goals ag
		CROSS JOIN achieved_milestones am
		CROSS JOIN total_photos tp`

	var summary models.ProgressSummary
	var startingBodyFat, currentBodyFat sql.NullFloat64
	var lastMeasurementDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&summary.CurrentWeight,
		&summary.StartingWeight,
		&currentBodyFat,
		&startingBodyFat,
		&lastMeasurementDate,
		&summary.DaysTracked,
		&summary.ActiveGoals,
		&summary.AchievedMilestones,
		&summary.TotalPhotos,
		&summary.WeightChangePercentage,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty summary for new users
			return &models.ProgressSummary{}, nil
		}
		return nil, fmt.Errorf("failed to get progress summary: %w", err)
	}

	if currentBodyFat.Valid {
		summary.CurrentBodyFat = &currentBodyFat.Float64
	}
	if startingBodyFat.Valid {
		summary.StartingBodyFat = &startingBodyFat.Float64
	}
	if lastMeasurementDate.Valid {
		summary.LastMeasurementDate = &lastMeasurementDate.Time
	}

	// Calculate weight change
	summary.WeightChange = summary.StartingWeight - summary.CurrentWeight

	return &summary, nil
}

// GetWeightProgressAnalytics analyzes weight trends and patterns
func (r *ProgressAnalyticsRepository) GetWeightProgressAnalytics(ctx context.Context, userID int64, days int) (*models.WeightProgressAnalytics, error) {
	query := `
		WITH weight_data AS (
			SELECT 
				measurement_date,
				weight,
				LAG(weight, 1) OVER (ORDER BY measurement_date) as previous_weight,
				EXTRACT(DAYS FROM measurement_date - LAG(measurement_date, 1) OVER (ORDER BY measurement_date)) as days_diff
			FROM body_measurements
			WHERE user_id = $1 
				AND measurement_date >= NOW() - INTERVAL '%d days'
				AND weight IS NOT NULL
			ORDER BY measurement_date ASC
		),
		weekly_changes AS (
			SELECT 
				DATE_TRUNC('week', measurement_date) as week,
				AVG(weight - previous_weight) as avg_weekly_change
			FROM weight_data
			WHERE previous_weight IS NOT NULL AND days_diff <= 7
			GROUP BY DATE_TRUNC('week', measurement_date)
		)
		SELECT 
			COUNT(*) as total_measurements,
			MIN(weight) as min_weight,
			MAX(weight) as max_weight,
			AVG(weight) as avg_weight,
			AVG(weight - previous_weight) as avg_daily_change,
			COUNT(CASE WHEN weight - previous_weight < 0 THEN 1 END) as weight_loss_days,
			COUNT(CASE WHEN weight - previous_weight > 0 THEN 1 END) as weight_gain_days,
			COUNT(CASE WHEN weight - previous_weight = 0 THEN 1 END) as stable_days,
			AVG(w.avg_weekly_change) as avg_weekly_change,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY weight) as median_weight
		FROM weight_data wd
		LEFT JOIN weekly_changes wc ON true
		WHERE wd.previous_weight IS NOT NULL`

	query = fmt.Sprintf(query, days)
	var analytics models.WeightProgressAnalytics
	var avgWeeklyChange, medianWeight sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.TotalMeasurements,
		&analytics.MinWeight,
		&analytics.MaxWeight,
		&analytics.AvgWeight,
		&analytics.AvgDailyChange,
		&analytics.WeightLossDays,
		&analytics.WeightGainDays,
		&analytics.StableDays,
		&avgWeeklyChange,
		&medianWeight,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get weight progress analytics: %w", err)
	}

	if avgWeeklyChange.Valid {
		analytics.AvgWeeklyChange = &avgWeeklyChange.Float64
	}
	if medianWeight.Valid {
		analytics.MedianWeight = &medianWeight.Float64
	}

	// Calculate consistency score
	if analytics.TotalMeasurements > 0 {
		stableRatio := float64(analytics.StableDays) / float64(analytics.TotalMeasurements)
		analytics.ConsistencyScore = stableRatio * 100
	}

	// Determine trend
	if analytics.AvgDailyChange < -0.1 {
		analytics.Trend = "decreasing"
	} else if analytics.AvgDailyChange > 0.1 {
		analytics.Trend = "increasing"
	} else {
		analytics.Trend = "stable"
	}

	return &analytics, nil
}

// GetBodyCompositionAnalytics analyzes body composition changes
func (r *ProgressAnalyticsRepository) GetBodyCompositionAnalytics(ctx context.Context, userID int64, days int) (*models.BodyCompositionAnalytics, error) {
	query := `
		WITH composition_data AS (
			SELECT 
				measurement_date,
				weight,
				body_fat_percentage,
				waist,
				chest,
				hips,
				LAG(body_fat_percentage, 1) OVER (ORDER BY measurement_date) as prev_body_fat,
				LAG(waist, 1) OVER (ORDER BY measurement_date) as prev_waist,
				LAG(chest, 1) OVER (ORDER BY measurement_date) as prev_chest,
				LAG(hips, 1) OVER (ORDER BY measurement_date) as prev_hips
			FROM body_measurements
			WHERE user_id = $1 
				AND measurement_date >= NOW() - INTERVAL '%d days'
				AND (body_fat_percentage IS NOT NULL OR waist IS NOT NULL)
			ORDER BY measurement_date ASC
		)
		SELECT 
			COUNT(*) as total_measurements,
			AVG(body_fat_percentage) as avg_body_fat,
			MIN(body_fat_percentage) as min_body_fat,
			MAX(body_fat_percentage) as max_body_fat,
			AVG(waist) as avg_waist,
			MIN(waist) as min_waist,
			MAX(waist) as max_waist,
			AVG(body_fat_percentage - prev_body_fat) as avg_body_fat_change,
			AVG(waist - prev_waist) as avg_waist_change,
			COUNT(CASE WHEN body_fat_percentage - prev_body_fat < 0 THEN 1 END) as body_fat_decrease_days,
			COUNT(CASE WHEN waist - prev_waist < 0 THEN 1 END) as waist_decrease_days
		FROM composition_data
		WHERE prev_body_fat IS NOT NULL OR prev_waist IS NOT NULL`

	query = fmt.Sprintf(query, days)
	var analytics models.BodyCompositionAnalytics
	var avgBodyFat, minBodyFat, maxBodyFat, avgWaist, minWaist, maxWaist sql.NullFloat64
	var avgBodyFatChange, avgWaistChange sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.TotalMeasurements,
		&avgBodyFat,
		&minBodyFat,
		&maxBodyFat,
		&avgWaist,
		&minWaist,
		&maxWaist,
		&avgBodyFatChange,
		&avgWaistChange,
		&analytics.BodyFatDecreaseDays,
		&analytics.WaistDecreaseDays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get body composition analytics: %w", err)
	}

	if avgBodyFat.Valid {
		analytics.AvgBodyFat = &avgBodyFat.Float64
	}
	if minBodyFat.Valid {
		analytics.MinBodyFat = &minBodyFat.Float64
	}
	if maxBodyFat.Valid {
		analytics.MaxBodyFat = &maxBodyFat.Float64
	}
	if avgWaist.Valid {
		analytics.AvgWaist = &avgWaist.Float64
	}
	if minWaist.Valid {
		analytics.MinWaist = &minWaist.Float64
	}
	if maxWaist.Valid {
		analytics.MaxWaist = &maxWaist.Float64
	}
	if avgBodyFatChange.Valid {
		analytics.AvgBodyFatChange = &avgBodyFatChange.Float64
	}
	if avgWaistChange.Valid {
		analytics.AvgWaistChange = &avgWaistChange.Float64
	}

	return &analytics, nil
}

// GetMilestoneAnalytics analyzes milestone achievement patterns
func (r *ProgressAnalyticsRepository) GetMilestoneAnalytics(ctx context.Context, userID int64) (*models.MilestoneAnalytics, error) {
	query := `
		WITH milestone_data AS (
			SELECT 
				milestone_type,
				category,
				priority,
				is_achieved,
				EXTRACT(DAYS FROM achieved_date - created_at) as days_to_achieve,
				EXTRACT(DAYS FROM target_date - achieved_date) as days_ahead_or_behind
			FROM milestones
			WHERE user_id = $1
		),
		type_stats AS (
			SELECT 
				milestone_type,
				COUNT(*) as total,
				COUNT(CASE WHEN is_achieved = true THEN 1 END) as achieved,
				AVG(days_to_achieve) as avg_days_to_achieve
			FROM milestone_data
			GROUP BY milestone_type
		),
		category_stats AS (
			SELECT 
				category,
				COUNT(*) as total,
				COUNT(CASE WHEN is_achieved = true THEN 1 END) as achieved
			FROM milestone_data
			GROUP BY category
		)
		SELECT 
			COUNT(*) as total_milestones,
			COUNT(CASE WHEN is_achieved = true THEN 1 END) as achieved_milestones,
			COUNT(CASE WHEN priority = 'high' AND is_achieved = true THEN 1 END) as high_priority_achieved,
			COUNT(CASE WHEN priority = 'high' THEN 1 END) as total_high_priority,
			AVG(days_to_achieve) as avg_days_to_achieve,
			MIN(days_to_achieve) as fastest_achievement,
			MAX(days_to_achieve) as slowest_achievement,
			COUNT(CASE WHEN days_ahead_or_behind < 0 THEN 1 END) as ahead_of_schedule,
			COUNT(CASE WHEN days_ahead_or_behind > 0 THEN 1 END) as behind_schedule
		FROM milestone_data
		WHERE is_achieved = true`

	var analytics models.MilestoneAnalytics
	var avgDaysAchieve, fastestAchieve, slowestAchieve sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.TotalMilestones,
		&analytics.AchievedMilestones,
		&analytics.HighPriorityAchieved,
		&analytics.TotalHighPriority,
		&avgDaysAchieve,
		&fastestAchieve,
		&slowestAchieve,
		&analytics.AheadOfSchedule,
		&analytics.BehindSchedule,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get milestone analytics: %w", err)
	}

	if avgDaysAchieve.Valid {
		analytics.AvgDaysToAchieve = &avgDaysAchieve.Float64
	}
	if fastestAchieve.Valid {
		analytics.FastestAchievement = &fastestAchieve.Float64
	}
	if slowestAchieve.Valid {
		analytics.SlowestAchievement = &slowestAchieve.Float64
	}

	// Calculate achievement rate
	if analytics.TotalMilestones > 0 {
		analytics.AchievementRate = float64(analytics.AchievedMilestones) / float64(analytics.TotalMilestones) * 100
	}

	// Calculate high priority success rate
	if analytics.TotalHighPriority > 0 {
		analytics.HighPrioritySuccessRate = float64(analytics.HighPriorityAchieved) / float64(analytics.TotalHighPriority) * 100
	}

	return &analytics, nil
}

// GetPhotoAnalytics analyzes progress photo patterns and usage
func (r *ProgressAnalyticsRepository) GetPhotoAnalytics(ctx context.Context, userID int64) (*models.PhotoAnalytics, error) {
	query := `
		WITH photo_data AS (
			SELECT 
				capture_date,
				created_at,
				file_size,
				weight,
				body_fat_percentage,
				tags
			FROM progress_photos
			WHERE user_id = $1
		),
		monthly_photos AS (
			SELECT 
				DATE_TRUNC('month', capture_date) as month,
				COUNT(*) as photo_count
			FROM photo_data
			GROUP BY DATE_TRUNC('month', capture_date)
		),
		tag_analytics AS (
			SELECT 
				UNNEST(string_to_array(replace(replace(tags::text, '{', ''), '}', ''), ',')) as tag,
				COUNT(*) as tag_count
			FROM photo_data
			WHERE tags IS NOT NULL
			GROUP BY UNNEST(string_to_array(replace(replace(tags::text, '{', ''), '}', ''), ','))
			ORDER BY tag_count DESC
			LIMIT 10
		)
		SELECT 
			COUNT(*) as total_photos,
			COUNT(CASE WHEN weight IS NOT NULL THEN 1 END) as photos_with_weight,
			COUNT(CASE WHEN body_fat_percentage IS NOT NULL THEN 1 END) as photos_with_body_fat,
			AVG(file_size) as avg_file_size,
			MIN(capture_date) as first_photo_date,
			MAX(capture_date) as last_photo_date,
			COUNT(CASE WHEN tags IS NOT NULL THEN 1 END) as tagged_photos,
			AVG(photo_count) as avg_monthly_photos,
			MAX(photo_count) as max_monthly_photos
		FROM photo_data pd
		LEFT JOIN monthly_photos mp ON true`

	var analytics models.PhotoAnalytics
	var firstPhoto, lastPhoto sql.NullTime
	var avgFileSize, avgMonthly, maxMonthly sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.TotalPhotos,
		&analytics.PhotosWithWeight,
		&analytics.PhotosWithBodyFat,
		&avgFileSize,
		&firstPhoto,
		&lastPhoto,
		&analytics.TaggedPhotos,
		&avgMonthly,
		&maxMonthly,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get photo analytics: %w", err)
	}

	if firstPhoto.Valid {
		analytics.FirstPhotoDate = &firstPhoto.Time
	}
	if lastPhoto.Valid {
		analytics.LastPhotoDate = &lastPhoto.Time
	}
	if avgFileSize.Valid {
		analytics.AvgFileSize = &avgFileSize.Float64
	}
	if avgMonthly.Valid {
		analytics.AvgMonthlyPhotos = &avgMonthly.Float64
	}
	if maxMonthly.Valid {
		analytics.MaxMonthlyPhotos = &maxMonthly.Float64
	}

	// Calculate tagging rate
	if analytics.TotalPhotos > 0 {
		analytics.TaggingRate = float64(analytics.TaggedPhotos) / float64(analytics.TotalPhotos) * 100
	}

	return &analytics, nil
}

// GetProgressPredictions predicts future progress based on historical data
func (r *ProgressAnalyticsRepository) GetProgressPredictions(ctx context.Context, userID int64, daysToPredict int) (*models.ProgressPredictions, error) {
	query := `
		WITH recent_measurements AS (
			SELECT 
				weight,
				body_fat_percentage,
				measurement_date,
				EXTRACT(DAYS FROM NOW() - measurement_date) as days_ago
			FROM body_measurements
			WHERE user_id = $1 
				AND measurement_date >= NOW() - INTERVAL '60 days'
				AND weight IS NOT NULL
			ORDER BY measurement_date DESC
			LIMIT 10
		),
		weight_trend AS (
			SELECT 
				CORR(weight, days_ago) as weight_correlation,
				REGR_SLOPE(weight, days_ago) as weight_slope,
				REGR_INTERCEPT(weight, days_ago) as weight_intercept
			FROM recent_measurements
		),
		goal_data AS (
			SELECT 
				wg.target_weight,
				wg.goal_type,
				wg.weekly_target,
				wg.target_date
			FROM weight_goals wg
			WHERE wg.user_id = $1 AND wg.is_active = true
			LIMIT 1
		)
		SELECT 
			wt.weight_slope as daily_weight_change,
			wt.weight_slope * 7 as weekly_weight_change,
			gd.target_weight,
			gd.goal_type,
			gd.weekly_target,
			gd.target_date,
			CASE 
				WHEN wt.weight_slope < 0 AND gd.goal_type = 'lose' THEN 'on_track'
				WHEN wt.weight_slope > 0 AND gd.goal_type = 'gain' THEN 'on_track'
				WHEN wt.weight_slope = 0 AND gd.goal_type = 'maintain' THEN 'on_track'
				ELSE 'off_track'
			END as trajectory_status
		FROM weight_trend wt
		CROSS JOIN goal_data gd`

	var predictions models.ProgressPredictions
	var targetWeight sql.NullFloat64
	var weeklyTarget sql.NullFloat64
	var targetDate sql.NullTime
	var dailyChange, weeklyChange sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&dailyChange,
		&weeklyChange,
		&targetWeight,
		&predictions.GoalType,
		&weeklyTarget,
		&targetDate,
		&predictions.TrajectoryStatus,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get progress predictions: %w", err)
	}

	if dailyChange.Valid {
		predictions.DailyWeightChange = &dailyChange.Float64
	}
	if weeklyChange.Valid {
		predictions.WeeklyWeightChange = &weeklyChange.Float64
	}
	if targetWeight.Valid {
		predictions.TargetWeight = &targetWeight.Float64
	}
	if weeklyTarget.Valid {
		predictions.WeeklyTarget = &weeklyTarget.Float64
	}
	if targetDate.Valid {
		predictions.TargetDate = &targetDate.Time
	}

	// Calculate predicted weight after specified days
	if predictions.DailyWeightChange != nil {
		currentWeight, err := r.getCurrentWeight(ctx, userID)
		if err == nil {
			predictedWeight := currentWeight + (*predictions.DailyWeightChange * float64(daysToPredict))
			predictions.PredictedWeight = &predictedWeight
		}
	}

	return &predictions, nil
}

// Helper method to get current weight
func (r *ProgressAnalyticsRepository) getCurrentWeight(ctx context.Context, userID int64) (float64, error) {
	query := `
		SELECT weight 
		FROM body_measurements 
		WHERE user_id = $1 
		ORDER BY measurement_date DESC 
		LIMIT 1`

	var weight float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&weight)
	if err != nil {
		return 0, fmt.Errorf("failed to get current weight: %w", err)
	}

	return weight, nil
}

// GetConsistencyAnalytics analyzes user consistency in tracking
func (r *ProgressAnalyticsRepository) GetConsistencyAnalytics(ctx context.Context, userID int64, days int) (*models.ConsistencyAnalytics, error) {
	query := `
		WITH measurement_frequency AS (
			SELECT 
				DATE_TRUNC('week', measurement_date) as week,
				COUNT(*) as measurements_in_week
			FROM body_measurements
			WHERE user_id = $1 
				AND measurement_date >= NOW() - INTERVAL '%d days'
			GROUP BY DATE_TRUNC('week', measurement_date)
		),
		photo_frequency AS (
			SELECT 
				DATE_TRUNC('week', capture_date) as week,
				COUNT(*) as photos_in_week
			FROM progress_photos
			WHERE user_id = $1 
				AND capture_date >= NOW() - INTERVAL '%d days'
			GROUP BY DATE_TRUNC('week', capture_date)
		),
		ideal_weeks AS (
			SELECT GENERATE_SERIES(
				DATE_TRUNC('week', NOW() - INTERVAL '%d days'),
				DATE_TRUNC('week', NOW()),
				INTERVAL '1 week'
			)::date as week
		)
		SELECT 
			COUNT(DISTINCT mf.week) as weeks_with_measurements,
			COUNT(DISTINCT pf.week) as weeks_with_photos,
			COUNT(*) as total_weeks,
			AVG(mf.measurements_in_week) as avg_measurements_per_week,
			AVG(pf.photos_in_week) as avg_photos_per_week,
			MAX(mf.measurements_in_week) as max_measurements_in_week,
			COUNT(CASE WHEN mf.measurements_in_week >= 1 THEN 1 END) as weeks_with_any_measurements,
			COUNT(CASE WHEN pf.photos_in_week >= 1 THEN 1 END) as weeks_with_any_photos
		FROM ideal_weeks iw
		LEFT JOIN measurement_frequency mf ON iw.week = mf.week
		LEFT JOIN photo_frequency pf ON iw.week = pf.week`

	query = fmt.Sprintf(query, days, days, days)
	var analytics models.ConsistencyAnalytics
	var avgMeasurements, avgPhotos, maxMeasurements sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.WeeksWithMeasurements,
		&analytics.WeeksWithPhotos,
		&analytics.TotalWeeks,
		&avgMeasurements,
		&avgPhotos,
		&maxMeasurements,
		&analytics.WeeksWithAnyMeasurements,
		&analytics.WeeksWithAnyPhotos,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get consistency analytics: %w", err)
	}

	if avgMeasurements.Valid {
		analytics.AvgMeasurementsPerWeek = &avgMeasurements.Float64
	}
	if avgPhotos.Valid {
		analytics.AvgPhotosPerWeek = &avgPhotos.Float64
	}
	if maxMeasurements.Valid {
		analytics.MaxMeasurementsInWeek = &maxMeasurements.Float64
	}

	// Calculate consistency scores
	if analytics.TotalWeeks > 0 {
		analytics.MeasurementConsistency = float64(analytics.WeeksWithMeasurements) / float64(analytics.TotalWeeks) * 100
		analytics.PhotoConsistency = float64(analytics.WeeksWithPhotos) / float64(analytics.TotalWeeks) * 100
		analytics.OverallConsistency = (analytics.MeasurementConsistency + analytics.PhotoConsistency) / 2
	}

	return &analytics, nil
}

// GetAchievementAnalytics analyzes user achievements and patterns
func (r *ProgressAnalyticsRepository) GetAchievementAnalytics(ctx context.Context, userID int64) (*models.AchievementAnalytics, error) {
	query := `
		WITH achievements AS (
			SELECT 
				'milestone' as achievement_type,
				achieved_date as date,
				title as description,
				CASE WHEN priority = 'high' THEN 3 
					 WHEN priority = 'medium' THEN 2 
					 ELSE 1 END as points
			FROM milestones 
			WHERE user_id = $1 AND is_achieved = true
			UNION ALL
			SELECT 
				'goal' as achievement_type,
				achieved_date as date,
				CONCAT('Weight goal: ', goal_type) as description,
				CASE WHEN goal_type = 'lose' THEN 5 
					 WHEN goal_type = 'gain' THEN 5 
					 ELSE 3 END as points
			FROM weight_goals 
			WHERE user_id = $1 AND achieved_date IS NOT NULL
		),
		monthly_achievements AS (
			SELECT 
				DATE_TRUNC('month', date) as month,
				COUNT(*) as achievement_count,
				SUM(points) as total_points
			FROM achievements
			GROUP BY DATE_TRUNC('month', date)
		),
		streaks AS (
			SELECT 
				achievement_type,
				date,
				LEAD(date, 1) OVER (PARTITION BY achievement_type ORDER BY date) as next_date,
				EXTRACT(DAYS FROM LEAD(date, 1) OVER (PARTITION BY achievement_type ORDER BY date) - date) as days_diff
			FROM achievements
		)
		SELECT 
			COUNT(*) as total_achievements,
			COUNT(DISTINCT achievement_type) as achievement_types,
			SUM(points) as total_points,
			AVG(total_points) as avg_monthly_points,
			MAX(total_points) as max_monthly_points,
			MIN(date) as first_achievement,
			MAX(date) as last_achievement,
			COUNT(CASE WHEN achievement_type = 'milestone' THEN 1 END) as milestone_achievements,
			COUNT(CASE WHEN achievement_type = 'goal' THEN 1 END) as goal_achievements,
			COUNT(CASE WHEN days_diff <= 7 THEN 1 END) as consecutive_achievements
		FROM achievements a
		LEFT JOIN monthly_achievements ma ON true
		LEFT JOIN streaks s ON true`

	var analytics models.AchievementAnalytics
	var firstAchievement, lastAchievement sql.NullTime
	var avgMonthly, maxMonthly sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&analytics.TotalAchievements,
		&analytics.AchievementTypes,
		&analytics.TotalPoints,
		&avgMonthly,
		&maxMonthly,
		&firstAchievement,
		&lastAchievement,
		&analytics.MilestoneAchievements,
		&analytics.GoalAchievements,
		&analytics.ConsecutiveAchievements,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get achievement analytics: %w", err)
	}

	if firstAchievement.Valid {
		analytics.FirstAchievementDate = &firstAchievement.Time
	}
	if lastAchievement.Valid {
		analytics.LastAchievementDate = &lastAchievement.Time
	}
	if avgMonthly.Valid {
		analytics.AvgMonthlyPoints = &avgMonthly.Float64
	}
	if maxMonthly.Valid {
		analytics.MaxMonthlyPoints = &maxMonthly.Float64
	}

	return &analytics, nil
}
