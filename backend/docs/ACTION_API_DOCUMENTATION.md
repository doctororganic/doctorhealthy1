# Action-Oriented API Documentation

**Version**: 1.0.0  
**Base URL**: `/api/v1/actions`  
**Authentication**: All endpoints require JWT Bearer token

---

## Overview

This API follows an action-oriented architecture where users interact with buttons/actions rather than direct CRUD operations. All endpoints are user-facing and handle business logic internally.

---

## 🔐 Authentication

All endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

---

## 📊 Progress Tracking Actions

### 1. Track Measurement
**Action**: User clicks "Log Measurement" button

**Endpoint**: `POST /api/v1/actions/track-measurement`

**Request Body**:
```json
{
  "measurement_type": "waist",
  "value": 85.5,
  "unit": "cm",
  "date": "2024-01-15",
  "notes": "Morning measurement"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Measurement logged successfully",
  "data": {
    "id": 1,
    "user_id": 123,
    "measurement_type": "waist",
    "value": 85.5,
    "unit": "cm",
    "date": "2024-01-15",
    "notes": "Morning measurement",
    "created_at": "2024-01-15T08:00:00Z"
  }
}
```

**Valid Measurement Types**:
- `weight`, `height`, `body_fat`, `body_fat_percentage`, `muscle_mass`
- `waist`, `chest`, `hips`, `neck`
- `left_bicep`, `right_bicep`, `left_forearm`, `right_forearm`
- `left_thigh`, `right_thigh`, `left_calf`, `right_calf`

---

### 2. Get Progress Summary
**Action**: User clicks "View Progress" button

**Endpoint**: `GET /api/v1/actions/progress-summary?days=30`

**Query Parameters**:
- `days` (optional): Number of days to analyze (default: 30)

**Response**:
```json
{
  "status": "success",
  "data": {
    "period_days": 30,
    "start_date": "2023-12-15",
    "end_date": "2024-01-15",
    "stats": {
      "weight": {
        "current": 75.5,
        "previous": 77.0,
        "change": -1.5,
        "change_percent": -1.95,
        "trend": "decreasing"
      },
      "waist": {
        "current": 85.5,
        "previous": 88.0,
        "change": -2.5,
        "change_percent": -2.84,
        "trend": "decreasing"
      }
    }
  }
}
```

---

### 3. Get Measurement History
**Action**: User views measurement history

**Endpoint**: `GET /api/v1/actions/measurement-history?type=waist&page=1&limit=20`

**Query Parameters**:
- `type` (optional): Filter by measurement type
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)
- `start_date` (optional): Start date filter (YYYY-MM-DD)
- `end_date` (optional): End date filter (YYYY-MM-DD)

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "measurement_type": "waist",
      "value": 85.5,
      "unit": "cm",
      "date": "2024-01-15",
      "notes": "Morning measurement"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "totalPages": 3
  }
}
```

---

### 4. Get Progress Charts
**Action**: User views progress charts

**Endpoint**: `GET /api/v1/actions/progress-charts?days=90&type=weight`

**Query Parameters**:
- `days` (optional): Number of days (default: 30)
- `type` (optional): Measurement type to chart

**Response**:
```json
{
  "status": "success",
  "data": {
    "chart_data": [
      {
        "date": "2024-01-01",
        "value": 77.0
      },
      {
        "date": "2024-01-15",
        "value": 75.5
      }
    ],
    "trend": "decreasing",
    "average_change": -0.1
  }
}
```

---

### 5. Compare Measurements
**Action**: User clicks "Compare" button

**Endpoint**: `POST /api/v1/actions/compare-measurements`

**Request Body**:
```json
{
  "start_date": "2024-01-01",
  "end_date": "2024-01-15"
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "start_date": "2024-01-01",
    "end_date": "2024-01-15",
    "comparisons": {
      "weight": {
        "start": 77.0,
        "end": 75.5,
        "change": -1.5,
        "change_percent": -1.95
      }
    }
  }
}
```

---

### 6. Upload Progress Photo
**Action**: User clicks "Upload Photo" button

**Endpoint**: `POST /api/v1/actions/upload-progress-photo`

**Request Body**:
```json
{
  "photo_url": "https://example.com/photos/photo.jpg",
  "thumbnail_url": "https://example.com/photos/thumb.jpg",
  "date": "2024-01-15",
  "weight": 75.5,
  "notes": "Front view"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Progress photo uploaded successfully",
  "data": {
    "id": 1,
    "photo_url": "https://example.com/photos/photo.jpg",
    "date": "2024-01-15",
    "weight": 75.5
  }
}
```

---

### 7. Get Photo History
**Action**: User views photo gallery

**Endpoint**: `GET /api/v1/actions/photo-history?page=1&limit=20`

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "photo_url": "https://example.com/photos/photo.jpg",
      "date": "2024-01-15",
      "weight": 75.5
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 10,
    "totalPages": 1
  }
}
```

---

## 🍎 Nutrition Actions

### 1. Generate Meal Plan
**Action**: User clicks "Generate Meal Plan" button

**Endpoint**: `POST /api/v1/actions/generate-meal-plan`

**Request Body**:
```json
{
  "goal": "weight_loss",
  "target_calories": 2000,
  "duration": 7,
  "preferences": ["vegetarian", "low-carb"],
  "restrictions": ["gluten-free"]
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Meal plan generated successfully",
  "data": {
    "user_id": "123",
    "goal": "weight_loss",
    "target_calories": 2000,
    "duration_days": 7,
    "preferences": ["vegetarian", "low-carb"],
    "restrictions": ["gluten-free"],
    "generated_at": "2024-01-15T10:00:00Z",
    "meals": []
  }
}
```

---

### 2. Log Meal
**Action**: User clicks "Log Meal" button

**Endpoint**: `POST /api/v1/actions/log-meal`

**Request Body**:
```json
{
  "food_id": 123,
  "recipe_id": null,
  "meal_type": "breakfast",
  "quantity": 1.5,
  "unit": "serving",
  "date": "2024-01-15",
  "notes": "Delicious breakfast"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Meal logged successfully",
  "data": {
    "food_id": 123,
    "meal_type": "breakfast",
    "quantity": 1.5,
    "unit": "serving",
    "date": "2024-01-15",
    "logged_at": "2024-01-15T08:00:00Z"
  }
}
```

**Valid Meal Types**: `breakfast`, `lunch`, `dinner`, `snack`

---

### 3. Get Nutrition Summary
**Action**: User clicks "View Nutrition Summary" button

**Endpoint**: `GET /api/v1/actions/nutrition-summary?days=7`

**Query Parameters**:
- `days` (optional): Number of days (default: 7)

**Response**:
```json
{
  "status": "success",
  "data": {
    "period_days": 7,
    "start_date": "2024-01-08",
    "end_date": "2024-01-15",
    "totals": {
      "calories": 14000,
      "protein": 700,
      "carbs": 1750,
      "fat": 350,
      "fiber": 140
    },
    "daily_averages": {
      "calories": 2000,
      "protein": 100,
      "carbs": 250,
      "fat": 50
    },
    "meals_logged": 21
  }
}
```

---

### 4. Get Meal Recommendations
**Action**: User clicks "Get Recommendations" button

**Endpoint**: `GET /api/v1/actions/meal-recommendations?meal_type=breakfast&calories=500`

**Query Parameters**:
- `meal_type` (optional): Type of meal (default: breakfast)
- `calories` (optional): Maximum calories

**Response**:
```json
{
  "status": "success",
  "meal_type": "breakfast",
  "max_calories": 500,
  "recommendations": [
    {
      "id": "rec_1",
      "name": "Healthy Breakfast Option",
      "meal_type": "breakfast",
      "calories": 450,
      "description": "A balanced meal recommendation"
    }
  ]
}
```

---

## 💪 Fitness Actions

### 1. Generate Workout
**Action**: User clicks "Generate Workout" button

**Endpoint**: `POST /api/v1/actions/generate-workout`

**Request Body**:
```json
{
  "goal": "weight_loss",
  "duration": 30,
  "difficulty": "intermediate",
  "equipment": ["dumbbells", "mat"],
  "muscle_groups": ["chest", "arms"],
  "restrictions": ["knee_injury"]
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Workout plan generated successfully",
  "data": {
    "goal": "weight_loss",
    "duration": 30,
    "difficulty": "intermediate",
    "equipment": ["dumbbells", "mat"],
    "muscle_groups": ["chest", "arms"],
    "restrictions": ["knee_injury"],
    "generated_at": "2024-01-15T10:00:00Z",
    "exercises": []
  }
}
```

**Valid Goals**: `weight_loss`, `muscle_gain`, `endurance`, `flexibility`, `general_fitness`  
**Valid Difficulty**: `beginner`, `intermediate`, `advanced`

---

### 2. Log Workout
**Action**: User clicks "Log Workout" button

**Endpoint**: `POST /api/v1/actions/log-workout`

**Request Body**:
```json
{
  "workout_plan_id": 123,
  "exercises": [
    {
      "exercise_id": 1,
      "sets": 3,
      "reps": 12,
      "weight": 20.5,
      "duration": 0
    }
  ],
  "duration": 30,
  "date": "2024-01-15",
  "notes": "Great workout!"
}
```

**Response**:
```json
{
  "status": "success",
  "message": "Workout logged successfully",
  "data": {
    "workout_plan_id": 123,
    "exercises": [...],
    "duration": 30,
    "date": "2024-01-15",
    "logged_at": "2024-01-15T18:00:00Z"
  }
}
```

---

### 3. Get Fitness Summary
**Action**: User clicks "View Fitness Summary" button

**Endpoint**: `GET /api/v1/actions/fitness-summary?days=30`

**Query Parameters**:
- `days` (optional): Number of days (default: 30)

**Response**:
```json
{
  "status": "success",
  "data": {
    "period_days": 30,
    "start_date": "2023-12-15",
    "end_date": "2024-01-15",
    "totals": {
      "workouts_completed": 20,
      "total_duration": 600,
      "total_calories": 12000
    },
    "averages": {
      "workouts_per_week": 5,
      "avg_duration": 30,
      "avg_calories": 600
    },
    "most_trained_muscles": ["chest", "arms", "legs"],
    "favorite_exercises": ["Push-ups", "Squats"]
  }
}
```

---

### 4. Get Workout Recommendations
**Action**: User clicks "Get Recommendations" button

**Endpoint**: `GET /api/v1/actions/workout-recommendations?goal=weight_loss&duration=30`

**Query Parameters**:
- `goal` (optional): Workout goal (default: general_fitness)
- `duration` (optional): Duration in minutes (default: 30)

**Response**:
```json
{
  "status": "success",
  "goal": "weight_loss",
  "duration": 30,
  "recommendations": [
    {
      "id": "workout_rec_1",
      "name": "Recommended Workout",
      "goal": "weight_loss",
      "duration": 30,
      "description": "A personalized workout recommendation"
    }
  ]
}
```

---

## 🔄 Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message describing what went wrong"
}
```

**Common HTTP Status Codes**:
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request format or parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## 📝 Notes

1. **Date Format**: All dates should be in `YYYY-MM-DD` format
2. **Pagination**: Default page size is 20, maximum is 100
3. **Authentication**: All endpoints require valid JWT token
4. **User Isolation**: Users can only access their own data
5. **Validation**: All input is validated server-side

---

## 🚀 Frontend Integration Example

```typescript
// Example: Log a measurement
async function logMeasurement(type: string, value: number) {
  try {
    const response = await fetch('/api/v1/actions/track-measurement', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify({
        measurement_type: type,
        value: value,
        unit: 'cm'
      })
    });
    
    const result = await response.json();
    
    if (result.status === 'success') {
      showSuccessMessage('Measurement logged!');
      return result.data;
    } else {
      showError(result.error);
    }
  } catch (error) {
    showError('Failed to log measurement');
  }
}
```

---

**Last Updated**: $(date)

