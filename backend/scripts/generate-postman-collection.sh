#!/bin/bash
# Generate Postman Collection from API endpoints
# Usage: ./scripts/generate-postman-collection.sh

cat > postman_collection.json << 'EOF'
{
  "info": {
    "name": "Nutrition Platform API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080",
      "type": "string"
    },
    {
      "key": "token",
      "value": "",
      "type": "string"
    }
  ],
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"password123\",\n  \"first_name\": \"Test\",\n  \"last_name\": \"User\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/v1/auth/register",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "auth", "register"]
            }
          }
        },
        {
          "name": "Login",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "if (pm.response.code === 200) {",
                  "    var jsonData = pm.response.json();",
                  "    pm.collectionVariables.set('token', jsonData.data.access_token);",
                  "}"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"password123\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/v1/auth/login",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "auth", "login"]
            }
          }
        }
      ]
    },
    {
      "name": "Actions",
      "item": [
        {
          "name": "Track Measurement",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"waist\": 85.5,\n  \"measurement_date\": \"2025-01-15\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/v1/actions/track-measurement",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "actions", "track-measurement"]
            }
          }
        },
        {
          "name": "Progress Summary",
          "request": {
            "method": "GET",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              }
            ],
            "url": {
              "raw": "{{base_url}}/api/v1/actions/progress-summary?days=30",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "actions", "progress-summary"],
              "query": [
                {
                  "key": "days",
                  "value": "30"
                }
              ]
            }
          }
        },
        {
          "name": "Generate Meal Plan",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"goal\": \"weight_loss\",\n  \"target_calories\": 2000,\n  \"duration\": 7\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/v1/actions/generate-meal-plan",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "actions", "generate-meal-plan"]
            }
          }
        },
        {
          "name": "Generate Workout",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"goal\": \"weight_loss\",\n  \"duration\": 30,\n  \"difficulty\": \"intermediate\"\n}"
            },
            "url": {
              "raw": "{{base_url}}/api/v1/actions/generate-workout",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "actions", "generate-workout"]
            }
          }
        }
      ]
    }
  ]
}
EOF

echo "✅ Postman collection generated: postman_collection.json"
echo "Import this file into Postman to test all endpoints!"