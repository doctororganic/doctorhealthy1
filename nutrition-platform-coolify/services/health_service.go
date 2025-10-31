package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"nutrition-platform/models"
)

// HealthService handles health-related operations
type HealthService struct {
	db *sql.DB
}

// NewHealthService creates a new health service
func NewHealthService(db *sql.DB) *HealthService {
	return &HealthService{db: db}
}

// CreateHealthComplaint creates a new health complaint
func (s *HealthService) CreateHealthComplaint(userID string, req *models.CreateHealthComplaintRequest) (*models.UserHealthComplaint, error) {
	complaint := &models.UserHealthComplaint{
		ID:                       generateUUID(),
		UserID:                   userID,
		ComplaintType:            req.ComplaintType,
		Severity:                 req.Severity,
		Description:              req.Description,
		Symptoms:                 req.Symptoms,
		DurationDays:             req.DurationDays,
		Frequency:                req.Frequency,
		Triggers:                 req.Triggers,
		CurrentMedications:       req.CurrentMedications,
		ReportedAt:               time.Now(),
		Status:                   "active",
		MedicalAttentionRequired: req.MedicalAttentionRequired,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	if err := s.storeHealthComplaint(complaint); err != nil {
		return nil, fmt.Errorf("failed to store health complaint: %w", err)
	}

	return complaint, nil
}

// CreateUserInjury creates a new user injury record
func (s *HealthService) CreateUserInjury(userID string, req *models.CreateUserInjuryRequest) (*models.UserInjury, error) {
	injury := &models.UserInjury{
		ID:                       generateUUID(),
		UserID:                   userID,
		InjuryID:                 req.InjuryID,
		CustomInjuryName:         req.CustomInjuryName,
		Severity:                 req.Severity,
		InjuryDate:               req.InjuryDate,
		Description:              req.Description,
		TreatmentReceived:        req.TreatmentReceived,
		CurrentStatus:            req.CurrentStatus,
		AffectsExercise:          req.AffectsExercise,
		ExerciseLimitations:      req.ExerciseLimitations,
		MedicalClearanceRequired: req.MedicalClearanceRequired,
		ExpectedRecoveryDate:     req.ExpectedRecoveryDate,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	if err := s.storeUserInjury(injury); err != nil {
		return nil, fmt.Errorf("failed to store user injury: %w", err)
	}

	return injury, nil
}

// GetHealthConditions retrieves health conditions with optional filtering
func (s *HealthService) GetHealthConditions(category string, page, limit int) ([]models.HealthCondition, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, category)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("%s", conditions[0])
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT id, name, name_ar, category, icd_10_code, description, description_ar,
		       symptoms, risk_factors, complications, dietary_recommendations,
		       exercise_recommendations, lifestyle_modifications, severity_levels,
		       is_chronic, requires_medical_supervision, created_at, updated_at
		FROM health_conditions %s
		ORDER BY name
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query health conditions: %w", err)
	}
	defer rows.Close()

	var healthConditions []models.HealthCondition
	for rows.Next() {
		var condition models.HealthCondition
		var symptomsJSON, riskFactorsJSON, complicationsJSON []byte
		var dietaryRecommendationsJSON, exerciseRecommendationsJSON, lifestyleModificationsJSON []byte
		var severityLevelsJSON []byte

		err := rows.Scan(
			&condition.ID, &condition.Name, &condition.NameAr, &condition.Category,
			&condition.ICD10Code, &condition.Description, &condition.DescriptionAr,
			&symptomsJSON, &riskFactorsJSON, &complicationsJSON,
			&dietaryRecommendationsJSON, &exerciseRecommendationsJSON,
			&lifestyleModificationsJSON, &severityLevelsJSON,
			&condition.IsChronic, &condition.RequiresMedicalSupervision,
			&condition.CreatedAt, &condition.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan health condition: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(symptomsJSON, &condition.Symptoms)
		json.Unmarshal(riskFactorsJSON, &condition.RiskFactors)
		json.Unmarshal(complicationsJSON, &condition.Complications)
		json.Unmarshal(dietaryRecommendationsJSON, &condition.DietaryRecommendations)
		json.Unmarshal(exerciseRecommendationsJSON, &condition.ExerciseRecommendations)
		json.Unmarshal(lifestyleModificationsJSON, &condition.LifestyleModifications)
		json.Unmarshal(severityLevelsJSON, &condition.SeverityLevels)

		healthConditions = append(healthConditions, condition)
	}

	return healthConditions, nil
}

// GetUserHealthComplaints retrieves user's health complaints
func (s *HealthService) GetUserHealthComplaints(userID string, status string) ([]models.UserHealthComplaint, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := "WHERE " + fmt.Sprintf("%s", conditions[0])
	if len(conditions) > 1 {
		whereClause += " AND " + fmt.Sprintf("%s", conditions[1])
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, complaint_type, severity, description, symptoms,
		       duration_days, frequency, triggers, current_medications,
		       reported_at, status, medical_attention_required, created_at, updated_at
		FROM user_health_complaints %s
		ORDER BY reported_at DESC
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user health complaints: %w", err)
	}
	defer rows.Close()

	var complaints []models.UserHealthComplaint
	for rows.Next() {
		var complaint models.UserHealthComplaint
		var symptomsJSON, triggersJSON, currentMedicationsJSON []byte

		err := rows.Scan(
			&complaint.ID, &complaint.UserID, &complaint.ComplaintType,
			&complaint.Severity, &complaint.Description, &symptomsJSON,
			&complaint.DurationDays, &complaint.Frequency, &triggersJSON,
			&currentMedicationsJSON, &complaint.ReportedAt, &complaint.Status,
			&complaint.MedicalAttentionRequired, &complaint.CreatedAt, &complaint.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan health complaint: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(symptomsJSON, &complaint.Symptoms)
		json.Unmarshal(triggersJSON, &complaint.Triggers)
		json.Unmarshal(currentMedicationsJSON, &complaint.CurrentMedications)

		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

// GetUserInjuries retrieves user's injuries
func (s *HealthService) GetUserInjuries(userID string, status string) ([]models.UserInjury, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("current_status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := "WHERE " + fmt.Sprintf("%s", conditions[0])
	if len(conditions) > 1 {
		whereClause += " AND " + fmt.Sprintf("%s", conditions[1])
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, injury_id, custom_injury_name, severity, injury_date,
		       description, treatment_received, current_status, affects_exercise,
		       exercise_limitations, medical_clearance_required, expected_recovery_date,
		       created_at, updated_at
		FROM user_injuries %s
		ORDER BY injury_date DESC
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user injuries: %w", err)
	}
	defer rows.Close()

	var injuries []models.UserInjury
	for rows.Next() {
		var injury models.UserInjury
		var exerciseLimitationsJSON []byte

		err := rows.Scan(
			&injury.ID, &injury.UserID, &injury.InjuryID, &injury.CustomInjuryName,
			&injury.Severity, &injury.InjuryDate, &injury.Description,
			&injury.TreatmentReceived, &injury.CurrentStatus, &injury.AffectsExercise,
			&exerciseLimitationsJSON, &injury.MedicalClearanceRequired,
			&injury.ExpectedRecoveryDate, &injury.CreatedAt, &injury.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user injury: %w", err)
		}

		// Parse JSON fields
		json.Unmarshal(exerciseLimitationsJSON, &injury.ExerciseLimitations)

		injuries = append(injuries, injury)
	}

	return injuries, nil
}

// PerformHealthAssessment performs a comprehensive health assessment
func (s *HealthService) PerformHealthAssessment(req *models.HealthAssessmentRequest) (*models.HealthAssessmentResponse, error) {
	assessment := &models.HealthAssessmentResponse{
		GeneratedAt: time.Now(),
	}

	// Calculate BMI
	heightM := req.Height / 100 // Convert cm to meters
	assessment.BMI = req.Weight / (heightM * heightM)
	assessment.BMICategory = s.getBMICategory(assessment.BMI)

	// Calculate BMR using Mifflin-St Jeor Equation
	if req.Gender == "male" {
		assessment.BMR = 10*req.Weight + 6.25*req.Height - 5*float64(req.Age) + 5
	} else {
		assessment.BMR = 10*req.Weight + 6.25*req.Height - 5*float64(req.Age) - 161
	}

	// Calculate TDEE
	activityMultipliers := map[string]float64{
		"sedentary":   1.2,
		"light":       1.375,
		"moderate":    1.55,
		"active":      1.725,
		"very_active": 1.9,
	}
	assessment.TDEE = assessment.BMR * activityMultipliers[req.ActivityLevel]

	// Assess health risk factors
	assessment.HealthRiskFactors = s.assessHealthRiskFactors(req)

	// Generate recommendations
	assessment.Recommendations = s.generateHealthRecommendations(req, assessment)

	// Calculate nutritional needs
	assessment.NutritionalNeeds = s.calculateNutritionalNeeds(req, assessment)

	// Generate exercise recommendations
	assessment.ExerciseRecommendations = s.generateExerciseRecommendations(req)

	// Generate lifestyle recommendations
	assessment.LifestyleRecommendations = s.generateLifestyleRecommendations(req)

	// Identify red flags
	assessment.RedFlags = s.identifyRedFlags(req)

	// Generate follow-up recommendations
	assessment.FollowUpRecommendations = s.generateFollowUpRecommendations(req)

	// Calculate overall risk score
	assessment.RiskScore = s.calculateRiskScore(req, assessment)

	return assessment, nil
}

// Helper methods

func (s *HealthService) getBMICategory(bmi float64) string {
	if bmi < 18.5 {
		return "Underweight"
	} else if bmi < 25 {
		return "Normal weight"
	} else if bmi < 30 {
		return "Overweight"
	} else {
		return "Obese"
	}
}

func (s *HealthService) assessHealthRiskFactors(req *models.HealthAssessmentRequest) []string {
	var riskFactors []string

	// Age-related risks
	if req.Age > 65 {
		riskFactors = append(riskFactors, "Advanced age increases risk of chronic diseases")
	}

	// BMI-related risks
	heightM := req.Height / 100
	bmi := req.Weight / (heightM * heightM)
	if bmi >= 30 {
		riskFactors = append(riskFactors, "Obesity increases risk of diabetes, heart disease, and other conditions")
	} else if bmi < 18.5 {
		riskFactors = append(riskFactors, "Underweight may indicate nutritional deficiencies")
	}

	// Lifestyle risks
	if req.SmokingStatus == "current" {
		riskFactors = append(riskFactors, "Smoking significantly increases risk of cancer, heart disease, and lung disease")
	}

	if req.AlcoholConsumption == "heavy" {
		riskFactors = append(riskFactors, "Heavy alcohol consumption increases risk of liver disease and other health issues")
	}

	if req.SleepHoursPerNight < 6 {
		riskFactors = append(riskFactors, "Insufficient sleep increases risk of various health problems")
	}

	if req.StressLevel >= 8 {
		riskFactors = append(riskFactors, "High stress levels can negatively impact physical and mental health")
	}

	if req.ExerciseFrequency < 3 {
		riskFactors = append(riskFactors, "Insufficient physical activity increases risk of chronic diseases")
	}

	return riskFactors
}

func (s *HealthService) generateHealthRecommendations(req *models.HealthAssessmentRequest, assessment *models.HealthAssessmentResponse) []models.HealthRecommendation {
	var recommendations []models.HealthRecommendation

	// BMI-based recommendations
	if assessment.BMI >= 30 {
		recommendations = append(recommendations, models.HealthRecommendation{
			Category:    "Weight Management",
			Priority:    "high",
			Title:       "Weight Loss Program",
			Description: "Implement a structured weight loss program to reduce health risks",
			ActionItems: []string{
				"Create a caloric deficit of 500-750 calories per day",
				"Increase physical activity to 150+ minutes per week",
				"Focus on whole foods and portion control",
				"Consider consulting with a registered dietitian",
			},
			Timeline: "6-12 months",
		})
	}

	// Exercise recommendations
	if req.ExerciseFrequency < 3 {
		recommendations = append(recommendations, models.HealthRecommendation{
			Category:    "Physical Activity",
			Priority:    "high",
			Title:       "Increase Physical Activity",
			Description: "Regular exercise is crucial for overall health and disease prevention",
			ActionItems: []string{
				"Aim for at least 150 minutes of moderate-intensity exercise per week",
				"Include both cardiovascular and strength training exercises",
				"Start gradually and progressively increase intensity",
				"Find activities you enjoy to maintain consistency",
			},
			Timeline: "Start immediately, build over 4-6 weeks",
		})
	}

	// Sleep recommendations
	if req.SleepHoursPerNight < 7 {
		recommendations = append(recommendations, models.HealthRecommendation{
			Category:    "Sleep Health",
			Priority:    "medium",
			Title:       "Improve Sleep Quality and Duration",
			Description: "Adequate sleep is essential for physical and mental health",
			ActionItems: []string{
				"Aim for 7-9 hours of sleep per night",
				"Establish a consistent sleep schedule",
				"Create a relaxing bedtime routine",
				"Limit screen time before bed",
			},
			Timeline: "2-4 weeks to establish new habits",
		})
	}

	return recommendations
}

func (s *HealthService) calculateNutritionalNeeds(req *models.HealthAssessmentRequest, assessment *models.HealthAssessmentResponse) models.NutritionalNeeds {
	needs := models.NutritionalNeeds{
		DailyCalories: int(assessment.TDEE),
		Vitamins:      make(map[string]float64),
		Minerals:      make(map[string]float64),
	}

	// Calculate macronutrient needs
	needs.Protein = req.Weight * 1.2 // 1.2g per kg body weight (minimum)
	if req.ExerciseFrequency >= 4 {
		needs.Protein = req.Weight * 1.6 // Higher for active individuals
	}

	// Carbohydrates: 45-65% of total calories
	needs.Carbohydrates = float64(needs.DailyCalories) * 0.55 / 4 // 55% of calories, 4 cal/g

	// Fat: 20-35% of total calories
	needs.Fat = float64(needs.DailyCalories) * 0.30 / 9 // 30% of calories, 9 cal/g

	// Fiber: 25-35g per day
	needs.Fiber = 30

	// Water: 35ml per kg body weight
	needs.Water = req.Weight * 0.035

	// Basic vitamin and mineral needs (simplified)
	needs.Vitamins["Vitamin C"] = 90    // mg
	needs.Vitamins["Vitamin D"] = 20    // mcg
	needs.Vitamins["Vitamin B12"] = 2.4 // mcg

	needs.Minerals["Calcium"] = 1000 // mg
	needs.Minerals["Iron"] = 18      // mg (for women)
	if req.Gender == "male" {
		needs.Minerals["Iron"] = 8 // mg (for men)
	}
	needs.Minerals["Magnesium"] = 400 // mg

	return needs
}

func (s *HealthService) generateExerciseRecommendations(req *models.HealthAssessmentRequest) []models.ExerciseRecommendation {
	var recommendations []models.ExerciseRecommendation

	// Cardiovascular exercise
	recommendations = append(recommendations, models.ExerciseRecommendation{
		Type:      "cardio",
		Frequency: 5,  // times per week
		Duration:  30, // minutes
		Intensity: "moderate",
		SpecificExercises: []string{
			"Brisk walking",
			"Swimming",
			"Cycling",
			"Dancing",
		},
		Precautions: []string{
			"Start slowly if you're new to exercise",
			"Stay hydrated",
			"Stop if you experience chest pain or severe shortness of breath",
		},
		ProgressionPlan: "Increase duration by 5 minutes every 2 weeks",
	})

	// Strength training
	recommendations = append(recommendations, models.ExerciseRecommendation{
		Type:      "strength",
		Frequency: 2,  // times per week
		Duration:  45, // minutes
		Intensity: "moderate",
		SpecificExercises: []string{
			"Bodyweight squats",
			"Push-ups",
			"Planks",
			"Resistance band exercises",
		},
		Precautions: []string{
			"Focus on proper form over heavy weights",
			"Allow rest days between strength sessions",
			"Warm up before lifting",
		},
		ProgressionPlan: "Increase resistance or repetitions every 2-3 weeks",
	})

	return recommendations
}

func (s *HealthService) generateLifestyleRecommendations(req *models.HealthAssessmentRequest) []string {
	var recommendations []string

	if req.SmokingStatus == "current" {
		recommendations = append(recommendations, "Quit smoking - consider nicotine replacement therapy or counseling")
	}

	if req.AlcoholConsumption == "heavy" {
		recommendations = append(recommendations, "Reduce alcohol consumption to moderate levels (1-2 drinks per day max)")
	}

	if req.StressLevel >= 7 {
		recommendations = append(recommendations, "Implement stress management techniques such as meditation, yoga, or deep breathing exercises")
	}

	if req.WaterIntakeLiters < 2 {
		recommendations = append(recommendations, "Increase daily water intake to at least 2-3 liters per day")
	}

	recommendations = append(recommendations, "Maintain regular meal times and avoid skipping meals")
	recommendations = append(recommendations, "Limit processed foods and increase consumption of whole foods")
	recommendations = append(recommendations, "Practice good hygiene and get regular health screenings")

	return recommendations
}

func (s *HealthService) identifyRedFlags(req *models.HealthAssessmentRequest) []string {
	var redFlags []string

	// Severe symptoms that require immediate attention
	for _, symptom := range req.CurrentSymptoms {
		switch symptom {
		case "chest pain", "severe headache", "difficulty breathing", "severe abdominal pain":
			redFlags = append(redFlags, fmt.Sprintf("Seek immediate medical attention for: %s", symptom))
		}
	}

	// Extreme BMI values
	heightM := req.Height / 100
	bmi := req.Weight / (heightM * heightM)
	if bmi < 16 || bmi > 40 {
		redFlags = append(redFlags, "Extreme BMI value requires medical evaluation")
	}

	// Multiple health conditions
	if len(req.HealthConditions) >= 3 {
		redFlags = append(redFlags, "Multiple health conditions require coordinated medical care")
	}

	return redFlags
}

func (s *HealthService) generateFollowUpRecommendations(req *models.HealthAssessmentRequest) []string {
	var followUp []string

	followUp = append(followUp, "Schedule annual physical examination")
	followUp = append(followUp, "Monitor blood pressure monthly")
	followUp = append(followUp, "Track weight weekly")

	if req.Age > 40 {
		followUp = append(followUp, "Annual blood work including lipid panel and glucose")
	}

	if req.Age > 50 {
		followUp = append(followUp, "Consider colonoscopy screening")
		if req.Gender == "female" {
			followUp = append(followUp, "Annual mammogram screening")
		}
	}

	return followUp
}

func (s *HealthService) calculateRiskScore(req *models.HealthAssessmentRequest, assessment *models.HealthAssessmentResponse) int {
	score := 0

	// Age factor
	if req.Age > 65 {
		score += 20
	} else if req.Age > 50 {
		score += 10
	}

	// BMI factor
	if assessment.BMI >= 30 {
		score += 15
	} else if assessment.BMI < 18.5 {
		score += 10
	}

	// Lifestyle factors
	if req.SmokingStatus == "current" {
		score += 25
	}
	if req.AlcoholConsumption == "heavy" {
		score += 15
	}
	if req.ExerciseFrequency < 2 {
		score += 10
	}
	if req.SleepHoursPerNight < 6 {
		score += 10
	}
	if req.StressLevel >= 8 {
		score += 10
	}

	// Health conditions
	score += len(req.HealthConditions) * 5

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

func (s *HealthService) storeHealthComplaint(complaint *models.UserHealthComplaint) error {
	symptomsJSON, _ := json.Marshal(complaint.Symptoms)
	triggersJSON, _ := json.Marshal(complaint.Triggers)
	currentMedicationsJSON, _ := json.Marshal(complaint.CurrentMedications)

	query := `
		INSERT INTO user_health_complaints (
			id, user_id, complaint_type, severity, description, symptoms,
			duration_days, frequency, triggers, current_medications,
			reported_at, status, medical_attention_required, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := s.db.Exec(query,
		complaint.ID, complaint.UserID, complaint.ComplaintType, complaint.Severity,
		complaint.Description, symptomsJSON, complaint.DurationDays, complaint.Frequency,
		triggersJSON, currentMedicationsJSON, complaint.ReportedAt, complaint.Status,
		complaint.MedicalAttentionRequired, complaint.CreatedAt, complaint.UpdatedAt,
	)

	return err
}

func (s *HealthService) storeUserInjury(injury *models.UserInjury) error {
	exerciseLimitationsJSON, _ := json.Marshal(injury.ExerciseLimitations)

	query := `
		INSERT INTO user_injuries (
			id, user_id, injury_id, custom_injury_name, severity, injury_date,
			description, treatment_received, current_status, affects_exercise,
			exercise_limitations, medical_clearance_required, expected_recovery_date,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := s.db.Exec(query,
		injury.ID, injury.UserID, injury.InjuryID, injury.CustomInjuryName,
		injury.Severity, injury.InjuryDate, injury.Description, injury.TreatmentReceived,
		injury.CurrentStatus, injury.AffectsExercise, exerciseLimitationsJSON,
		injury.MedicalClearanceRequired, injury.ExpectedRecoveryDate,
		injury.CreatedAt, injury.UpdatedAt,
	)

	return err
}
