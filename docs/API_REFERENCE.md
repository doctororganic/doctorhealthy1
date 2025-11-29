# Nutrition Platform API Reference

## Overview

The Nutrition Platform API provides comprehensive endpoints for nutrition data, fitness tracking, health monitoring, and personalized recommendations. This document covers all available endpoints, authentication, request/response formats, and usage examples.

**Base URL**: `https://api.nutrition-platform.com/api/v1`  
**Authentication**: JWT Bearer Token  
**Content-Type**: `application/json`  
**Rate Limiting**: 100 requests per minute per user  

## Table of Contents

1. [Authentication](#authentication)
2. [Nutrition Data Endpoints](#nutrition-data-endpoints)
3. [Health & Progress Endpoints](#health--progress-endpoints)
4. [Fitness Endpoints](#fitness-endpoints)
5. [User Management](#user-management)
6. [Error Handling](#error-handling)
7. [Rate Limiting](#rate-limiting)
8. [Caching](#caching)

## Authentication

### POST /auth/login
**Description**: Authenticate user and receive JWT token  
**Rate Limit**: 10 requests per minute  

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123",
      "email": "user@example.com",
      "name": "John Doe",
      "createdAt": "2023-01-01T00:00:00Z"
    },
    "expiresIn": 86400
  }
}
```

**Status Codes**:
- `200 OK`: Authentication successful
- `401 Unauthorized`: Invalid credentials
- `400 Bad Request`: Missing or invalid fields
- `429 Too Many Requests`: Rate limit exceeded

### POST /auth/register
**Description**: Register new user account  
**Rate Limit**: 5 requests per minute  

**Request Body**:
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "Jane Doe",
  "age": 30,
  "gender": "female",
  "height": 165,
  "weight": 60
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "456",
      "email": "newuser@example.com",
      "name": "Jane Doe"
    },
    "message": "User registered successfully"
  }
}
```

### POST /auth/refresh
**Description**: Refresh JWT token  
**Headers**: `Authorization: Bearer <token>`  

**Request Body**:
```json
{
  "refreshToken": "refresh_token_here"
}
```

## Nutrition Data Endpoints

### GET /nutrition-data/recipes
**Description**: Get recipes with pagination and filtering  
**Authentication**: Required  
**Cache**: 5 minutes  

**Query Parameters**:
- `page` (int, default: 1): Page number
- `limit` (int, default: 20, max: 100): Items per page
- `cuisine` (string, optional): Filter by cuisine type
- `dietType` (string, optional): Filter by diet type
- `maxCalories` (int, optional): Maximum calories per serving
- `minProtein` (int, optional): Minimum protein per serving
- `search` (string, optional): Search term for recipe names

**Example**:
```bash
curl "https://api.nutrition-platform.com/api/v1/nutrition-data/recipes?page=1&limit=20&cuisine=Italian&dietType=vegetarian"
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "items": [
      {
        "id": "recipe_001",
        "name": "Margherita Pizza",
        "cuisine": "Italian",
        "dietType": "vegetarian",
        "calories": 280,
        "protein": 12,
        "carbs": 35,
        "fat": 8,
        "ingredients": [
          {
            "name": "Pizza dough",
            "amount": "200g",
            "calories": 180
          },
          {
            "name": "Tomato sauce",
            "amount": "50g",
            "calories": 30
          }
        ],
        "instructions": [
          "Preheat oven to 220°C",
          "Roll out pizza dough",
          "Add tomato sauce and toppings"
        ],
        "prepTime": 20,
        "cookTime": 15,
        "servings": 2
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "totalPages": 8,
      "hasNext": true,
      "hasPrev": false
    }
  }
}
```

**Status Codes**:
- `200 OK`: Success
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Authentication required
- `429 Too Many Requests`: Rate limit exceeded

### GET /nutrition-data/workouts
**Description**: Get workout plans with filtering  
**Authentication**: Required  
**Cache**: 10 minutes  

**Query Parameters**:
- `page` (int, default: 1): Page number
- `limit` (int, default: 20): Items per page
- `difficulty` (string, optional): beginner, intermediate, advanced
- `duration` (int, optional): Workout duration in minutes
- `equipment` (string, optional): none, basic, full
- `muscleGroup` (string, optional): Targeted muscle group

**Response**:
```json
{
  "status": "success",
  "data": {
    "items": [
      {
        "id": "workout_001",
        "name": "Full Body Strength",
        "difficulty": "intermediate",
        "duration": 45,
        "equipment": "basic",
        "muscleGroups": ["chest", "back", "legs"],
        "exercises": [
          {
            "name": "Squats",
            "sets": 3,
            "reps": 12,
            "rest": 60,
            "description": "Stand with feet shoulder-width apart"
          }
        ],
        "caloriesBurned": 250
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 80,
      "totalPages": 4
    }
  }
}
```

### GET /nutrition-data/complaints
**Description**: Get nutrition-related complaints and solutions  
**Authentication**: Required  
**Cache**: 30 minutes  

**Query Parameters**:
- `page` (int, default: 1): Page number
- `limit` (int, default: 20): Items per page
- `category` (string, optional): digestive, energy, sleep, etc.

**Response**:
```json
{
  "status": "success",
  "data": {
    "items": [
      {
        "id": "complaint_001",
        "name": "Fatigue",
        "category": "energy",
        "description": "Persistent tiredness and lack of energy",
        "causes": ["poor nutrition", "lack of sleep", "stress"],
        "solutions": [
          {
            "type": "nutrition",
            "recommendation": "Increase iron-rich foods",
            "foods": ["spinach", "red meat", "lentils"]
          },
          {
            "type": "supplement",
            "recommendation": "Consider B-complex vitamins",
            "dosage": "50mg daily"
          }
        ]
      }
    ]
  }
}
```

## Health & Progress Endpoints

### GET /health
**Description**: System health check  
**Authentication**: None  
**Cache**: None  

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2023-12-01T12:00:00Z",
  "version": "1.2.0",
  "uptime": 86400,
  "checks": {
    "database": "healthy",
    "redis": "healthy",
    "external_apis": "healthy"
  }
}
```

### GET /health/detailed
**Description**: Detailed system health with metrics  
**Authentication**: Required  
**Cache**: 1 minute  

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2023-12-01T12:00:00Z",
  "metrics": {
    "cpu_usage": 45.2,
    "memory_usage": 67.8,
    "disk_usage": 23.1,
    "active_connections": 150,
    "requests_per_second": 25.5
  },
  "dependencies": {
    "database": {
      "status": "healthy",
      "response_time": 2
    },
    "redis": {
      "status": "healthy",
      "response_time": 1
    }
  }
}
```

### GET /progress/weight
**Description**: Get weight tracking data  
**Authentication**: Required  

**Query Parameters**:
- `startDate` (string, optional): ISO date string
- `endDate` (string, optional): ISO date string
- `limit` (int, default: 100): Maximum records

**Response**:
```json
{
  "status": "success",
  "data": {
    "entries": [
      {
        "id": "weight_001",
        "weight": 70.5,
        "bodyFat": 18.2,
        "muscleMass": 55.3,
        "date": "2023-12-01T08:00:00Z",
        "notes": "Morning weight after workout"
      }
    ],
    "summary": {
      "startWeight": 72.0,
      "currentWeight": 70.5,
      "totalChange": -1.5,
      "trend": "decreasing"
    }
  }
}
```

### POST /progress/weight
**Description**: Add weight entry  
**Authentication**: Required  

**Request Body**:
```json
{
  "weight": 70.5,
  "bodyFat": 18.2,
  "muscleMass": 55.3,
  "notes": "Morning weight after workout"
}
```

## Fitness Endpoints

### POST /fitness/workouts/generate
**Description**: Generate personalized workout plan  
**Authentication**: Required  
**Rate Limit**: 10 requests per hour  

**Request Body**:
```json
{
  "goals": ["muscle_gain", "fat_loss"],
  "preferences": {
    "duration": 45,
    "difficulty": "intermediate",
    "equipment": "basic",
    "muscleGroups": ["chest", "back", "legs"]
  },
  "limitations": {
    "injuries": ["knee"],
    "conditions": ["asthma"]
  }
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "workoutPlan": {
      "id": "generated_001",
      "name": "Personalized Full Body",
      "duration": 45,
      "difficulty": "intermediate",
      "exercises": [
        {
          "name": "Modified Squats",
          "sets": 3,
          "reps": 12,
          "adaptations": "knee-friendly modifications"
        }
      ],
      "schedule": {
        "frequency": "3 times per week",
        "restDays": ["Tuesday", "Thursday", "Sunday"]
      }
    }
  }
}
```

### GET /fitness/nutrition/calculator
**Description**: Calculate nutritional needs  
**Authentication**: Required  

**Query Parameters**:
- `age` (int, required): Age in years
- `gender` (string, required): male or female
- `height` (int, required): Height in cm
- `weight` (int, required): Weight in kg
- `activityLevel` (string, required): sedentary, light, moderate, active, very_active
- `goal` (string, required): lose_weight, maintain, gain_weight

**Response**:
```json
{
  "status": "success",
  "data": {
    "bmr": 1546,
    "tdee": 2154,
    "dailyCalories": {
      "lose_weight": 1654,
      "maintain": 2154,
      "gain_weight": 2654
    },
    "macros": {
      "protein": 120,
      "carbs": 200,
      "fats": 65
    },
    "recommendations": {
      "waterIntake": 2500,
      "fiberIntake": 30
    }
  }
}
```

## User Management

### GET /users/profile
**Description**: Get user profile  
**Authentication**: Required  

**Response**:
```json
{
  "status": "success",
  "data": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe",
    "age": 30,
    "gender": "male",
    "height": 175,
    "weight": 70,
    "preferences": {
      "dietaryRestrictions": ["none"],
      "fitnessGoals": ["muscle_gain"],
      "activityLevel": "moderate"
    },
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-12-01T10:00:00Z"
  }
}
```

### PUT /users/profile
**Description**: Update user profile  
**Authentication**: Required  

**Request Body**:
```json
{
  "name": "John Smith",
  "preferences": {
    "dietaryRestrictions": ["vegetarian"],
    "fitnessGoals": ["muscle_gain", "endurance"],
    "activityLevel": "active"
  }
}
```

## Error Handling

All API endpoints return consistent error responses:

```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": {
      "field": "age",
      "reason": "must be between 18 and 100"
    },
    "timestamp": "2023-12-01T12:00:00Z",
    "requestId": "req_123456"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|-------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Invalid request parameters |
| `AUTHENTICATION_ERROR` | 401 | Invalid or missing authentication |
| `AUTHORIZATION_ERROR` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `RATE_LIMIT_EXCEEDED` | 429 | Rate limit exceeded |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

## Rate Limiting

The API implements rate limiting to ensure fair usage:

- **Default Limit**: 100 requests per minute per user
- **Authentication**: 10 requests per minute (login attempts)
- **Workout Generation**: 10 requests per hour
- **Search Endpoints**: 200 requests per minute

### Rate Limit Headers

All responses include rate limiting headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1638360000
X-RateLimit-RetryAfter: 30
```

## Caching

The API implements intelligent caching to improve performance:

- **Nutrition Data**: 5-30 minutes based on endpoint
- **Health Checks**: 1 minute for detailed checks
- **User Data**: No caching (real-time)
- **Search Results**: 10 minutes

### Cache Headers

```
Cache-Control: public, max-age=300
X-Cache: HIT
ETag: "abc123"
Last-Modified: Wed, 01 Dec 2023 12:00:00 GMT
```

## SDKs and Libraries

### JavaScript/TypeScript
```bash
npm install @nutrition-platform/client
```

```typescript
import { NutritionClient } from '@nutrition-platform/client';

const client = new NutritionClient({
  baseUrl: 'https://api.nutrition-platform.com/api/v1',
  apiKey: 'your-api-key'
});

const recipes = await client.recipes.list({
  cuisine: 'Italian',
  dietType: 'vegetarian'
});
```

### Go
```bash
go get github.com/nutrition-platform/go-client
```

```go
import "github.com/nutrition-platform/go-client"

client := nutrition.NewClient("your-api-key")
recipes, err := client.Recipes.List(nutrition.RecipesParams{
    Cuisine: "Italian",
    DietType: "vegetarian",
})
```

## Testing

### Environment
- **Sandbox**: `https://sandbox-api.nutrition-platform.com/api/v1`
- **Production**: `https://api.nutrition-platform.com/api/v1`

### Test Credentials
- **Email**: `test@example.com`
- **Password**: `test123`

## Support

- **Documentation**: https://docs.nutrition-platform.com
- **Status Page**: https://status.nutrition-platform.com
- **Support Email**: api-support@nutrition-platform.com
- **GitHub Issues**: https://github.com/nutrition-platform/issues

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for detailed version history and API changes.
