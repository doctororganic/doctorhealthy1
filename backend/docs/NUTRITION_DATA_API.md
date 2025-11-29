# Nutrition Data API Documentation

## Overview

The Nutrition Data API provides access to comprehensive nutrition, workout, health complaint, metabolism, and drug-nutrition interaction data. All endpoints support pagination, filtering, searching, and sorting capabilities.

**Base URL**: `/api/v1/nutrition-data`

**API Version**: `1.0.0`

---

## Table of Contents

1. [Authentication](#authentication)
2. [Common Query Parameters](#common-query-parameters)
3. [Endpoints](#endpoints)
   - [Recipes](#recipes)
   - [Workouts](#workouts)
   - [Complaints](#complaints)
   - [Metabolism](#metabolism)
   - [Drugs-Nutrition](#drugs-nutrition)
   - [Answer Generation](#answer-generation)
4. [Validation Endpoints](#validation-endpoints)
5. [Error Handling](#error-handling)
6. [Rate Limiting](#rate-limiting)
7. [Examples](#examples)

---

## Authentication

Most endpoints require authentication via API key or JWT token. Include authentication in the request header:

```http
Authorization: Bearer YOUR_API_KEY_OR_TOKEN
```

**Note**: Some endpoints may be publicly accessible. Check individual endpoint documentation.

---

## Common Query Parameters

All GET endpoints support the following query parameters:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | integer | No | 1 | Page number (must be > 0) |
| `limit` | integer | No | 20 | Results per page (1-100) |
| `search` | string | No | - | Search query string |
| `filter` | string | No | - | Filter criteria (format: `field:value,field2:value2`) |
| `sort` | string | No | - | Sort field name |
| `order` | string | No | asc | Sort order (`asc` or `desc`) |
| `fields` | string | No | - | Field selection (comma-separated) |

### Filter Format

Filters use the format: `field:value` with multiple filters separated by commas.

**Example**: `filter=origin:mediterranean,calories_min:1300`

### Sort Order

- `asc` - Ascending order (default)
- `desc` - Descending order

---

## Endpoints

### Recipes

#### Get Recipes

Retrieve recipe data with meal plans, calorie levels, and weekly plans.

**Endpoint**: `GET /api/v1/nutrition-data/recipes`

**Query Parameters**:
- `page` (integer, optional): Page number
- `limit` (integer, optional): Results per page (1-100)
- `search` (string, optional): Search query
- `filter` (string, optional): Filter criteria
  - `origin` - Recipe origin (e.g., `mediterranean`)
  - `calories_min` - Minimum calories
  - `calories_max` - Maximum calories
- `sort` (string, optional): Sort field
- `order` (string, optional): Sort order (`asc` or `desc`)

**Response**:
```json
{
  "status": "success",
  "data": {
    "diet_name": "Mediterranean Diet",
    "principles": ["Fresh vegetables", "Olive oil", "Whole grains"],
    "calorie_levels": [
      {
        "calories": 1300,
        "goal": "weight_loss",
        "weekly_plan": {
          "Saturday": {
            "breakfast": "...",
            "lunch": "...",
            "dinner": "...",
            "snacks": "..."
          }
        }
      }
    ],
    "origin": "Mediterranean"
  },
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1,
    "has_next": false,
    "has_prev": false
  },
  "filters": {}
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/nutrition-data/recipes?page=1&limit=20&filter=origin:mediterranean"
```

**JavaScript Example**:
```javascript
const response = await fetch('/api/v1/nutrition-data/recipes?page=1&limit=20', {
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY'
  }
});
const data = await response.json();
```

**Python Example**:
```python
import requests

response = requests.get(
    'http://localhost:8080/api/v1/nutrition-data/recipes',
    params={'page': 1, 'limit': 20},
    headers={'Authorization': 'Bearer YOUR_API_KEY'}
)
data = response.json()
```

---

### Workouts

#### Get Workouts

Retrieve workout plans with exercises, training schedules, and goals.

**Endpoint**: `GET /api/v1/nutrition-data/workouts`

**Query Parameters**:
- `page` (integer, optional): Page number
- `limit` (integer, optional): Results per page (1-100)
- `search` (string, optional): Search query
- `filter` (string, optional): Filter criteria
  - `goal` - Workout goal (e.g., `weight_loss`, `muscle_gain`)
  - `training_days_per_week` - Training days (1-7)
  - `experience_level` - Experience level (e.g., `beginner`, `intermediate`)
- `sort` (string, optional): Sort field
- `order` (string, optional): Sort order (`asc` or `desc`)

**Response**:
```json
{
  "status": "success",
  "data": {
    "api_version": "1.0",
    "goal": "weight_loss",
    "training_days_per_week": 4,
    "weekly_plan": {
      "Day 1": {
        "exercises": [
          {
            "name": {
              "en": "Push-ups",
              "ar": "تمارين الضغط"
            },
            "sets": 3,
            "reps": 12
          }
        ]
      }
    },
    "scientific_references": []
  },
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1,
    "has_next": false,
    "has_prev": false
  },
  "filters": {}
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/nutrition-data/workouts?filter=goal:weight_loss&training_days_per_week:4"
```

---

### Complaints

#### Get All Complaints

Retrieve health complaint cases with recommendations.

**Endpoint**: `GET /api/v1/nutrition-data/complaints`

**Query Parameters**:
- `page` (integer, optional): Page number
- `limit` (integer, optional): Results per page (1-100)
- `search` (string, optional): Search query
- `filter` (string, optional): Filter criteria
  - `condition` - Health condition name
- `sort` (string, optional): Sort field
- `order` (string, optional): Sort order (`asc` or `desc`)

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "diabetes_type_2",
      "condition_en": "Type 2 Diabetes",
      "condition_ar": "داء السكري من النوع الثاني",
      "recommendations": {
        "nutrition": ["Low glycemic index foods", "Portion control"],
        "exercise": ["Regular aerobic exercise", "Strength training"],
        "medications": ["Metformin", "Insulin therapy"]
      },
      "enhanced_recommendations": {}
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  },
  "filters": {}
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/nutrition-data/complaints?page=1&limit=20&search=diabetes"
```

#### Get Complaint by ID

Retrieve a specific health complaint case by ID.

**Endpoint**: `GET /api/v1/nutrition-data/complaints/:id`

**Path Parameters**:
- `id` (string, required): Complaint case ID

**Response**:
```json
{
  "status": "success",
  "data": {
    "id": "diabetes_type_2",
    "condition_en": "Type 2 Diabetes",
    "condition_ar": "داء السكري من النوع الثاني",
    "recommendations": {
      "nutrition": ["Low glycemic index foods"],
      "exercise": ["Regular aerobic exercise"],
      "medications": ["Metformin"]
    }
  }
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/nutrition-data/complaints/diabetes_type_2"
```

**Error Response** (404):
```json
{
  "error": "Complaint not found",
  "message": "No complaint found with ID: invalid_id"
}
```

---

### Metabolism

#### Get Metabolism Guide

Retrieve comprehensive metabolism guide with sections and explanations.

**Endpoint**: `GET /api/v1/nutrition-data/metabolism`

**Query Parameters**: None (returns complete guide)

**Response**:
```json
{
  "status": "success",
  "data": {
    "metabolism_guide": {
      "title": {
        "en": "Metabolism Guide",
        "ar": "دليل الأيض"
      },
      "sections": [
        {
          "section_id": "metabolism_basics",
          "title": {
            "en": "Metabolism Basics",
            "ar": "أساسيات الأيض"
          },
          "content": {
            "en": "Metabolism is the process...",
            "ar": "الأيض هو العملية..."
          }
        }
      ]
    }
  }
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/nutrition-data/metabolism"
```

---

### Drugs-Nutrition

#### Get Drugs-Nutrition Interactions

Retrieve drug-nutrition interaction data and recommendations.

**Endpoint**: `GET /api/v1/nutrition-data/drugs-nutrition`

**Query Parameters**: None (returns complete data)

**Response**:
```json
{
  "status": "success",
  "data": {
    "supportedLanguages": ["en", "ar"],
    "nutritionalRecommendations": {
      "general": {
        "en": "General recommendations...",
        "ar": "التوصيات العامة..."
      },
      "interactions": {},
      "timing": {},
      "supplements": {}
    }
  }
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/nutrition-data/drugs-nutrition"
```

---

### Answer Generation

#### Generate Answer

Generate AI-powered contextual answers based on user queries using multiple data sources.

**Endpoint**: `POST /api/v1/nutrition-data/generate-answer`

**Request Body**:
```json
{
  "query": "What recipes are good for weight loss?",
  "data_types": ["recipes", "complaints"],
  "user_id": "optional_user_id"
}
```

**Request Fields**:
- `query` (string, required): User query/question
- `data_types` (array, required): Data types to search (options: `recipes`, `workouts`, `complaints`, `metabolism`, `drugs`)
- `user_id` (string, optional): User ID for personalization

**Response**:
```json
{
  "status": "success",
  "query": "What recipes are good for weight loss?",
  "answer": "Based on your query, here are recipe recommendations for weight loss:\n\n1. Mediterranean Diet Plan (1300 calories)\n- Focus on fresh vegetables and lean proteins\n- Includes weekly meal plans\n\n2. Low-Calorie Meal Options\n- Portion-controlled meals\n- High-fiber foods\n\nThese recipes are suitable for weight loss and provide approximately 1300-1500 calories per serving.",
  "sources": {
    "recipes": [...],
    "complaints": [...]
  },
  "quality_score": 8.5
}
```

**Example Request**:
```bash
curl -X POST "http://localhost:8080/api/v1/nutrition-data/generate-answer" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "query": "What recipes are good for diabetes?",
    "data_types": ["recipes", "complaints"]
  }'
```

**JavaScript Example**:
```javascript
const response = await fetch('/api/v1/nutrition-data/generate-answer', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_API_KEY'
  },
  body: JSON.stringify({
    query: 'What workouts are good for beginners?',
    data_types: ['workouts']
  })
});
const data = await response.json();
```

**Python Example**:
```python
import requests

response = requests.post(
    'http://localhost:8080/api/v1/nutrition-data/generate-answer',
    json={
        'query': 'What recipes are good for weight loss?',
        'data_types': ['recipes', 'complaints']
    },
    headers={'Authorization': 'Bearer YOUR_API_KEY'}
)
data = response.json()
```

**Query Examples**:
- Recipe queries: "recipes for weight loss", "healthy meal plans"
- Workout queries: "workout plan for beginners", "exercises for weight loss"
- Health queries: "advice for diabetes", "nutrition for high blood pressure"
- Drug queries: "interactions with metformin", "foods to avoid with medication"
- Metabolism queries: "how to boost metabolism", "metabolism explanation"

---

## Validation Endpoints

### Base URL: `/api/v1/validation`

### Validate All Files

Validate all nutrition data files and generate quality reports.

**Endpoint**: `GET /api/v1/validation/all`

**Response**:
```json
{
  "status": "success",
  "valid_count": 5,
  "invalid_count": 0,
  "total_files": 5,
  "results": [
    {
      "file": "qwen-recipes.json",
      "valid": true,
      "errors": [],
      "warnings": [],
      "quality": {
        "completeness": 95.0,
        "consistency": 90.0,
        "accuracy": 92.0,
        "uniqueness": 98.0,
        "overall": 93.75,
        "grade": "A"
      },
      "stats": {
        "total_records": 1,
        "fields_validated": 15
      }
    }
  ]
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/validation/all"
```

### Validate Specific File

Validate a specific nutrition data file.

**Endpoint**: `GET /api/v1/validation/file/:filename`

**Path Parameters**:
- `filename` (string, required): Name of the file to validate
  - Options: `qwen-recipes.json`, `qwen-workouts.json`, `complaints.json`, `metabolism.json`, `drugs-and-nutrition.json`

**Response**:
```json
{
  "status": "success",
  "result": {
    "file": "qwen-recipes.json",
    "valid": true,
    "errors": [],
    "warnings": [],
    "quality": {
      "completeness": 95.0,
      "consistency": 90.0,
      "accuracy": 92.0,
      "uniqueness": 98.0,
      "overall": 93.75,
      "grade": "A"
    },
    "suggestions": []
  }
}
```

**Example Request**:
```bash
curl "http://localhost:8080/api/v1/validation/file/qwen-recipes.json"
```

**Error Response** (400):
```json
{
  "error": "Filename is required"
}
```

---

## Error Handling

### Error Response Format

All errors follow this format:

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters or request body |
| 401 | Unauthorized - Missing or invalid authentication |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

### Common Errors

#### 400 Bad Request

**Invalid Query Parameters**:
```json
{
  "error": "Invalid parameters",
  "message": "invalid page parameter: must be positive integer"
}
```

**Missing Required Fields**:
```json
{
  "error": "Invalid request",
  "message": "Query is required"
}
```

#### 404 Not Found

**Resource Not Found**:
```json
{
  "error": "Complaint not found",
  "message": "No complaint found with ID: invalid_id"
}
```

#### 500 Internal Server Error

**Server Error**:
```json
{
  "error": "Failed to load recipes",
  "message": "File not found: qwen-recipes.json"
}
```

---

## Rate Limiting

API requests are rate-limited to ensure fair usage:

- **General API**: 100 requests per second per IP
- **Authentication endpoints**: 10 requests per second per IP

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

When rate limit is exceeded:

```json
{
  "error": "Rate limit exceeded",
  "message": "Too many requests. Please try again later."
}
```

---

## Examples

### Complete Workflow Example

```bash
# 1. Validate all data files
curl "http://localhost:8080/api/v1/validation/all"

# 2. Search for recipes
curl "http://localhost:8080/api/v1/nutrition-data/recipes?search=mediterranean&page=1&limit=10"

# 3. Get specific complaint
curl "http://localhost:8080/api/v1/nutrition-data/complaints/diabetes_type_2"

# 4. Generate answer
curl -X POST "http://localhost:8080/api/v1/nutrition-data/generate-answer" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What are good recipes and workouts for diabetes?",
    "data_types": ["recipes", "workouts", "complaints"]
  }'
```

### Advanced Filtering Example

```bash
# Filter recipes by origin and calorie range
curl "http://localhost:8080/api/v1/nutrition-data/recipes?filter=origin:mediterranean,calories_min:1300,calories_max:1500"

# Filter workouts by goal and training days
curl "http://localhost:8080/api/v1/nutrition-data/workouts?filter=goal:weight_loss,training_days_per_week:4"

# Search and paginate complaints
curl "http://localhost:8080/api/v1/nutrition-data/complaints?search=diabetes&page=2&limit=10&sort=id&order=asc"
```

### JavaScript SDK Example

```javascript
class NutritionDataAPI {
  constructor(baseURL, apiKey) {
    this.baseURL = baseURL;
    this.apiKey = apiKey;
  }

  async getRecipes(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const response = await fetch(`${this.baseURL}/nutrition-data/recipes?${queryString}`, {
      headers: { 'Authorization': `Bearer ${this.apiKey}` }
    });
    return response.json();
  }

  async generateAnswer(query, dataTypes) {
    const response = await fetch(`${this.baseURL}/nutrition-data/generate-answer`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiKey}`
      },
      body: JSON.stringify({ query, data_types: dataTypes })
    });
    return response.json();
  }
}

// Usage
const api = new NutritionDataAPI('http://localhost:8080/api/v1', 'YOUR_API_KEY');
const recipes = await api.getRecipes({ page: 1, limit: 20 });
const answer = await api.generateAnswer('recipes for weight loss', ['recipes']);
```

### Python SDK Example

```python
import requests

class NutritionDataAPI:
    def __init__(self, base_url, api_key):
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {'Authorization': f'Bearer {api_key}'}
    
    def get_recipes(self, **params):
        response = requests.get(
            f'{self.base_url}/nutrition-data/recipes',
            params=params,
            headers=self.headers
        )
        return response.json()
    
    def generate_answer(self, query, data_types):
        response = requests.post(
            f'{self.base_url}/nutrition-data/generate-answer',
            json={'query': query, 'data_types': data_types},
            headers={**self.headers, 'Content-Type': 'application/json'}
        )
        return response.json()

# Usage
api = NutritionDataAPI('http://localhost:8080/api/v1', 'YOUR_API_KEY')
recipes = api.get_recipes(page=1, limit=20)
answer = api.generate_answer('recipes for weight loss', ['recipes'])
```

---

## Support

For issues, questions, or feature requests, please contact the API support team or refer to the project documentation.

**API Version**: 1.0.0  
**Last Updated**: 2024  
**Documentation Version**: 1.0

