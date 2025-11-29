# Action API Documentation

## Overview

The Nutrition Platform Action API provides user-facing endpoints that users interact with through buttons and actions in the UI. These endpoints are designed to be intuitive and action-oriented, following the principle that users should be able to perform common actions with simple, clear API calls.

## Base URL

```
http://localhost:8080/api/v1/actions
```

## Authentication

All action endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Progress Tracking Actions

### Track Measurement
Logs a new body measurement for the user.

**Endpoint:** `POST /track-measurement`

**Request Body:**
```json
{
  "measurement_date": "2024-01-15T10:00:00Z",
  "weight": 75.5,
  "height": 175.0,
  "body_fat_percentage": 15.5,
  "muscle_mass": 60.0,
  "waist": 85.0,
  "chest": 95.0,
  "left_bicep": 35.0,
  "right_bicep": 35.5,
  "left_forearm": 28.0,
  "right_forearm": 28.5,
  "left_thigh": 55.0,
  "right_thigh": 55.5,
  "left_calf": 38.0,
  "right_calf": 38.5,
  "neck": 40.0,
  "hips": 95.0,
  "notes": "Morning measurements after workout"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Measurement logged successfully",
  "data": {
    "id": 123,
    "user_id": 1,
    "measurement_date": "2024-01-15T10:00:00Z",
    "weight": 75.5,
    "waist": 85.0,
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

### Get Progress Summary
Retrieves a summary of the user's progress over a specified period.

**Endpoint:** `GET /progress-summary?days=30`

**Query Parameters:**
- `days` (optional): Number of days to include in summary (default: 30)

**Response:**
```json
{
  "status": "success",
  "data": {
    "period_days": 30,
    "start_date": "2023-12-16",
    "end_date": "2024-01-15",
    "weight_change": -2.5,
    "measurements_count": 15,
    "latest_measurements": {
      "weight": 75.5,
      "waist": 85.0,
      "body_fat_percentage": 15.5
    },
    "trends": {
      "weight": "decreasing",
      "waist": "stable",
      "body_fat_percentage": "decreasing"
    }
  }
}
```

### Get Measurement History
Retrieves paginated history of user's measurements.

**Endpoint:** `GET /measurement-history?page=1&limit=20&start_date=2024-01-01&end_date=2024-12-31`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)
- `start_date` (optional): Filter by start date (YYYY-MM-DD)
- `end_date` (optional): Filter by end date (YYYY-MM-DD)

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 123,
      "measurement_date": "2024-01-15T10:00:00Z",
      "weight": 75.5,
      "waist": 85.0,
      "notes": "Morning measurements"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  }
}
```

### Get Progress Charts
Retrieves chart data for progress visualization.

**Endpoint:** `GET /progress-charts?days=90`

**Query Parameters:**
- `days` (optional): Number of days to include (default: 30)

**Response:**
```json
{
  "status": "success",
  "data": {
    "weight_chart": [
      {"date": "2024-01-01", "value": 78.0},
      {"date": "2024-01-08", "value": 77.5}
    ],
    "waist_chart": [
      {"date": "2024-01-01", "value": 86.0},
      {"date": "2024-01-08", "value": 85.5}
    ]
  }
}
```

### Compare Measurements
Compares measurements between two dates.

**Endpoint:** `POST /compare-measurements`

**Request Body:**
```json
{
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "start_measurements": {
      "date": "2024-01-01",
      "weight": 78.0,
      "waist": 86.0
    },
    "end_measurements": {
      "date": "2024-01-31",
      "weight": 75.5,
      "waist": 85.0
    },
    "changes": {
      "weight": -2.5,
      "waist": -1.0
    },
    "percentage_changes": {
      "weight": -3.21,
      "waist": -1.16
    }
  }
}
```

### Upload Progress Photo
Uploads a progress photo with optional metadata.

**Endpoint:** `POST /upload-progress-photo`

**Request Body:**
```json
{
  "photo_url": "https://example.com/photos/progress1.jpg",
  "thumbnail_url": "https://example.com/photos/progress1_thumb.jpg",
  "date": "2024-01-15",
  "weight": 75.5,
  "notes": "Front pose after 1 month"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Progress photo uploaded successfully",
  "data": {
    "id": 456,
    "user_id": 1,
    "photo_url": "https://example.com/photos/progress1.jpg",
    "thumbnail_url": "https://example.com/photos/progress1_thumb.jpg",
    "date": "2024-01-15",
    "weight": 75.5,
    "notes": "Front pose after 1 month",
    "uploaded_at": "2024-01-15T10:00:00Z"
  }
}
```

### Get Photo History
Retrieves paginated history of progress photos.

**Endpoint:** `GET /photo-history?page=1&limit=20`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 456,
      "photo_url": "https://example.com/photos/progress1.jpg",
      "thumbnail_url": "https://example.com/photos/progress1_thumb.jpg",
      "date": "2024-01-15",
      "weight": 75.5,
      "notes": "Front pose after 1 month",
      "uploaded_at": "2024-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 25,
    "totalPages": 2
  }
}
```

## Nutrition Actions

### Generate Meal Plan
Generates a personalized meal plan based on user goals and preferences.

**Endpoint:** `POST /generate-meal-plan`

**Request Body:**
```json
{
  "goal": "weight_loss",
  "target_calories": 2000,
  "duration": 7,
  "preferences": ["low_carb", "high_protein"],
  "restrictions": ["gluten_free", "dairy_free"]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Meal plan generated successfully",
  "data": {
    "id": 789,
    "user_id": 1,
    "goal": "weight_loss",
    "target_calories": 2000,
    "duration_days": 7,
    "preferences": ["low_carb", "high_protein"],
    "restrictions": ["gluten_free", "dairy_free"],
    "generated_at": "2024-01-15T10:00:00Z",
    "meals": [
      {
        "day": 1,
        "breakfast": {
          "name": "Protein Oatmeal",
          "calories": 450,
          "protein": 30,
          "carbs": 40,
          "fat": 12
        },
        "lunch": {
          "name": "Grilled Chicken Salad",
          "calories": 550,
          "protein": 40,
          "carbs": 30,
          "fat": 20
        },
        "dinner": {
          "name": "Baked Salmon with Vegetables",
          "calories": 600,
          "protein": 45,
          "carbs": 35,
          "fat": 25
        }
      }
    ]
  }
}
```

### Log Meal
Logs a consumed meal for the user.

**Endpoint:** `POST /log-meal`

**Request Body:**
```json
{
  "food_id": 123,
  "recipe_id": 456,
  "meal_type": "breakfast",
  "quantity": 200,
  "unit": "grams",
  "date": "2024-01-15",
  "notes": "Added extra protein powder"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Meal logged successfully",
  "data": {
    "id": 1011,
    "user_id": 1,
    "food_id": 123,
    "recipe_id": 456,
    "meal_type": "breakfast",
    "quantity": 200,
    "unit": "grams",
    "date": "2024-01-15",
    "notes": "Added extra protein powder",
    "logged_at": "2024-01-15T10:00:00Z"
  }
}
```

### Get Nutrition Summary
Retrieves nutrition summary over a specified period.

**Endpoint:** `GET /nutrition-summary?days=7`

**Query Parameters:**
- `days` (optional): Number of days to include (default: 7)

**Response:**
```json
{
  "status": "success",
  "data": {
    "period_days": 7,
    "start_date": "2024-01-08",
    "end_date": "2024-01-15",
    "totals": {
      "calories": 14000,
      "protein": 800,
      "carbs": 1200,
      "fat": 500,
      "fiber": 150
    },
    "daily_averages": {
      "calories": 2000,
      "protein": 114,
      "carbs": 171,
      "fat": 71
    },
    "meals_logged": 21,
    "goals_met": {
      "calories": true,
      "protein": true,
      "carbs": false
    }
  }
}
```

### Get Meal Recommendations
Gets meal recommendations based on meal type and calorie target.

**Endpoint:** `GET /meal-recommendations?meal_type=breakfast&calories=500`

**Query Parameters:**
- `meal_type` (optional): Type of meal (breakfast, lunch, dinner, snack)
- `calories` (optional): Target calories for the meal

**Response:**
```json
{
  "status": "success",
  "meal_type": "breakfast",
  "max_calories": 500,
  "recommendations": [
    {
      "id": "rec_1",
      "name": "Greek Yogurt Parfait",
      "meal_type": "breakfast",
      "calories": 450,
      "protein": 25,
      "carbs": 50,
      "fat": 15,
      "description": "High protein breakfast with berries and granola",
      "ingredients": ["greek yogurt", "mixed berries", "granola", "honey"],
      "prep_time": 5,
      "difficulty": "easy"
    }
  ]
}
```

## Fitness Actions

### Generate Workout
Generates a personalized workout plan.

**Endpoint:** `POST /generate-workout`

**Request Body:**
```json
{
  "goal": "weight_loss",
  "duration": 30,
  "difficulty": "beginner",
  "equipment": ["dumbbells", "resistance_bands"],
  "muscle_groups": ["legs", "core"],
  "restrictions": ["no_impact", "knee_friendly"]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Workout plan generated successfully",
  "data": {
    "id": 1213,
    "user_id": 1,
    "goal": "weight_loss",
    "duration": 30,
    "difficulty": "beginner",
    "equipment": ["dumbbells", "resistance_bands"],
    "muscle_groups": ["legs", "core"],
    "restrictions": ["no_impact", "knee_friendly"],
    "generated_at": "2024-01-15T10:00:00Z",
    "exercises": [
      {
        "name": "Bodyweight Squats",
        "sets": 3,
        "reps": 15,
        "rest_time": 60,
        "muscle_groups": ["legs", "glutes"],
        "equipment": ["bodyweight"],
        "difficulty": "beginner"
      },
      {
        "name": "Plank",
        "sets": 3,
        "duration": 30,
        "rest_time": 45,
        "muscle_groups": ["core"],
        "equipment": ["bodyweight"],
        "difficulty": "beginner"
      }
    ]
  }
}
```

### Log Workout
Logs a completed workout session.

**Endpoint:** `POST /log-workout`

**Request Body:**
```json
{
  "workout_plan_id": 1213,
  "exercises": [
    {
      "exercise_id": 1,
      "sets": 3,
      "reps": 15,
      "weight": 0,
      "duration": null,
      "notes": "Felt good, proper form"
    },
    {
      "exercise_id": 2,
      "sets": 3,
      "reps": null,
      "weight": null,
      "duration": 30,
      "notes": "Held plank for full duration"
    }
  ],
  "duration": 45,
  "date": "2024-01-15",
  "notes": "Great workout! Feeling stronger."
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Workout logged successfully",
  "data": {
    "id": 1414,
    "user_id": 1,
    "workout_plan_id": 1213,
    "exercises": [...],
    "duration": 45,
    "date": "2024-01-15",
    "notes": "Great workout! Feeling stronger.",
    "logged_at": "2024-01-15T10:00:00Z"
  }
}
```

### Get Fitness Summary
Retrieves fitness summary over a specified period.

**Endpoint:** `GET /fitness-summary?days=30`

**Query Parameters:**
- `days` (optional): Number of days to include (default: 30)

**Response:**
```json
{
  "status": "success",
  "data": {
    "period_days": 30,
    "start_date": "2023-12-16",
    "end_date": "2024-01-15",
    "totals": {
      "workouts_completed": 12,
      "total_duration": 540,
      "total_calories": 2400
    },
    "averages": {
      "workouts_per_week": 2.8,
      "avg_duration": 45,
      "avg_calories": 200
    },
    "most_trained_muscles": ["legs", "core", "chest"],
    "favorite_exercises": ["squats", "plank", "push-ups"]
  }
}
```

### Get Workout Recommendations
Gets workout recommendations based on goals and preferences.

**Endpoint:** `GET /workout-recommendations?goal=weight_loss&duration=30`

**Query Parameters:**
- `goal` (optional): Fitness goal (weight_loss, muscle_gain, endurance, flexibility)
- `duration` (optional): Workout duration in minutes

**Response:**
```json
{
  "status": "success",
  "goal": "weight_loss",
  "duration": 30,
  "recommendations": [
    {
      "id": "workout_rec_1",
      "name": "HIIT Full Body Burn",
      "goal": "weight_loss",
      "duration": 30,
      "difficulty": "intermediate",
      "description": "High-intensity interval training for maximum calorie burn",
      "equipment": ["bodyweight", "timer"],
      "muscle_groups": ["full_body"],
      "estimated_calories": 300
    }
  ]
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional error details"
  }
}
```

### Common Error Codes

- `400 BAD_REQUEST`: Invalid request data or validation errors
- `401 UNAUTHORIZED`: Authentication required or invalid token
- `403 FORBIDDEN`: Access denied (user doesn't have permission)
- `404 NOT_FOUND`: Resource not found
- `500 INTERNAL_ERROR`: Server error

## Rate Limiting

API endpoints are rate-limited to prevent abuse:
- 100 requests per minute per user
- 1000 requests per hour per user

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Total requests allowed
- `X-RateLimit-Remaining`: Requests remaining in current window
- `X-RateLimit-Reset`: Time when rate limit resets (Unix timestamp)

## Pagination

Endpoints that return lists support pagination:
- `page`: Page number (starts at 1)
- `limit`: Items per page (max 100)

Pagination metadata is included in responses:
```json
{
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  }
}
```

## Testing

### Quick Test Script

```bash
# Test authentication
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.access_token')

# Test track measurement
curl -X POST http://localhost:8080/api/v1/actions/track-measurement \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "weight": 75.5,
    "waist": 85.0,
    "notes": "Test measurement"
  }'

# Test progress summary
curl -X GET "http://localhost:8080/api/v1/actions/progress-summary?days=30" \
  -H "Authorization: Bearer $TOKEN"
```

## SDK/Client Libraries

### JavaScript/TypeScript

```typescript
import { actionApi } from './lib/actionApi';

// Track measurement
const result = await actionApi.progress.trackMeasurement({
  weight: 75.5,
  waist: 85.0,
  notes: "Morning measurements"
});

// Get progress summary
const summary = await actionApi.progress.getProgressSummary(30);

// Generate meal plan
const mealPlan = await actionApi.nutrition.generateMealPlan({
  goal: "weight_loss",
  target_calories: 2000,
  duration: 7
});
```

### React Components

```typescript
import { LogMeasurement } from './components/LogMeasurement';

function ProgressPage() {
  const handleMeasurementSuccess = (result) => {
    console.log('Measurement logged:', result);
  };

  const handleMeasurementError = (error) => {
    console.error('Failed to log measurement:', error);
  };

  return (
    <div>
      <h1>Progress Tracking</h1>
      <LogMeasurement 
        onSuccess={handleMeasurementSuccess}
        onError={handleMeasurementError}
      />
    </div>
  );
}
```

## Best Practices

1. **Authentication**: Always include a valid JWT token in requests
2. **Error Handling**: Implement proper error handling for all API calls
3. **Rate Limiting**: Respect rate limits and implement exponential backoff
4. **Validation**: Validate input data on the client side before sending
5. **Pagination**: Use pagination for large datasets
6. **Caching**: Cache frequently accessed data to reduce API calls

## Support

For API support and questions:
- Documentation: [Link to comprehensive docs]
- Issues: [Link to issue tracker]
- Email: support@nutrition-platform.com
