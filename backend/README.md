# Nutrition Platform Backend

A comprehensive nutrition and training platform backend API built with Go, Echo framework, and PostgreSQL.

## Features

### 🔐 **Security & Authentication**
- ✅ **Secure API Authentication** - API key-based authentication with scopes and rate limiting
- ✅ **Request Signing** - HMAC-SHA256 request signing for sensitive operations
- ✅ **Rate Limiting** - Configurable rate limiting per API key (100 requests/hour default)
- ✅ **Security Headers** - CORS, CSP, XSS protection, and other security headers

### 📊 **Monitoring & Analytics**
- ✅ **Real-time Analytics** - API usage monitoring and analytics
- ✅ **Health Checks** - Application and database health monitoring
- ✅ **Logging & Monitoring** - Structured logging with Prometheus metrics
- ✅ **Usage Tracking** - Comprehensive API usage tracking and reporting

### 🍽️ **Nutrition & Food Management**
- ✅ **Recipe Management** - Create, search, and manage recipes with nutritional information
- ✅ **Food Database** - Comprehensive food database with nutritional data
- ✅ **Nutrition Analysis** - Analyze nutritional content of foods and meals
- ✅ **Meal Planning** - Generate personalized meal plans based on goals and restrictions
- ✅ **Dietary Restrictions** - Support for halal, kosher, vegetarian, vegan, and allergy filters

### 🏥 **Health & Medical**
- ✅ **Health Conditions** - Database of diseases and health conditions with recommendations
- ✅ **Health Complaints** - Track and manage user health complaints and symptoms
- ✅ **Health Assessment** - Comprehensive health risk assessment and recommendations
- ✅ **Injury Management** - Track injuries and recovery recommendations
- ✅ **Symptom Checker** - Basic symptom checking with urgency assessment

### 💊 **Medications & Supplements**
- ✅ **Medication Database** - Comprehensive medication information and interactions
- ✅ **Drug Interactions** - Check for medication, food, and supplement interactions
- ✅ **Supplement Tracking** - Track vitamins, minerals, and supplement usage
- ✅ **Nutritional Effects** - Track how medications affect nutritional needs

### 🏋️ **Fitness & Workouts**
- ✅ **Workout Programs** - Comprehensive workout programs for different fitness levels
- ✅ **Exercise Database** - Detailed exercise information with instructions and modifications
- ✅ **Workout Tracking** - Track completed workouts and progress
- ✅ **Fitness Analytics** - Analyze workout patterns and progress over time

### 🎯 **Personalization & AI Recommendations**
- ✅ **AI-Powered Nutrition Plan Recommendations** - Evidence-based plan selection using scientific data
- ✅ **Nutritional Plans** - Detailed nutritional plans based on health conditions
- ✅ **Goal-based Recommendations** - Personalized recommendations based on user goals
- ✅ **Progress Tracking** - Track progress towards health and fitness goals
- ✅ **Risk Assessment** - Comprehensive health risk scoring and mitigation
- ✅ **Plan Comparison** - Compare multiple nutrition plans with scoring and rationale

### 🛠️ **Development & Deployment**
- ✅ **Comprehensive Testing** - Unit, integration, and security tests (70% coverage minimum)
- ✅ **Production Ready** - Docker, monitoring, and deployment scripts
- ✅ **Database Migrations** - Automated database schema management
- ✅ **Multi-language Support** - English and Arabic language support

## Quick Start

### Development

1. **Clone and setup:**
   ```bash
   cd nutrition-platform/backend
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run database migrations:**
   ```bash
   go run cmd/migrate/main.go -direction up
   ```

4. **Seed initial data:**
   ```bash
   go run cmd/seed/main.go
   ```

5. **Start development server:**
   ```bash
   make dev
   # or
   air
   ```

### Production Deployment

1. **Using Docker Compose:**
   ```bash
   docker-compose -f docker-compose.production.yml up -d
   ```

2. **Manual deployment:**
   ```bash
   chmod +x scripts/deploy-production.sh
   ./scripts/deploy-production.sh
   ```

## API Documentation

### Authentication

All API endpoints (except public ones) require an API key:

```bash
# Using Authorization header
curl -H "Authorization: Bearer nk_your_api_key_here" \
     https://api.example.com/api/v1/users

# Using X-API-Key header
curl -H "X-API-Key: nk_your_api_key_here" \
     https://api.example.com/api/v1/users
```

### API Key Management

```bash
# Create API key
POST /api/v1/api-keys
{
  "name": "My App Key",
  "scopes": ["nutrition", "read_only"],
  "rate_limit": 100
}

# List API keys
GET /api/v1/api-keys

# Get API key statistics
GET /api/v1/api-keys/{id}/stats?days=30

# Revoke API key
DELETE /api/v1/api-keys/{id}
```

### Core Endpoints

```bash
# Health check
GET /health

# API information
GET /api/info

# Nutrition analysis
POST /api/nutrition/analyze
{
  "food": "apple",
  "quantity": 100,
  "unit": "g",
  "checkHalal": true
}

# Generate meal plan
POST /api/generate-meal-plan
{
  "age": 25,
  "gender": "male",
  "height": 175,
  "weight": 70,
  "activityLevel": "moderate",
  "goal": "maintain"
}
```

### Recipe Management

```bash
# Create recipe
POST /api/v1/recipes
{
  "name": "Healthy Chicken Salad",
  "cuisine": "Mediterranean",
  "difficulty_level": "easy",
  "prep_time_minutes": 15,
  "ingredients": [...],
  "instructions": [...],
  "dietary_tags": ["gluten-free", "high-protein"]
}

# Search recipes
GET /api/v1/recipes?cuisine=Mediterranean&difficulty=easy&is_halal=true

# Get recipe by ID
GET /api/v1/recipes/{id}

# Rate recipe
POST /api/v1/recipes/{id}/rate
{
  "rating": 5,
  "review": "Delicious and healthy!"
}
```

### Health Management

```bash
# Health assessment
POST /api/v1/health/assessment
{
  "age": 30,
  "gender": "female",
  "height": 165,
  "weight": 60,
  "health_conditions": ["hypertension"],
  "current_medications": ["lisinopril"],
  "activity_level": "moderate"
}

# Report health complaint
POST /api/v1/health/complaints
{
  "complaint_type": "headache",
  "severity": 6,
  "symptoms": ["throbbing pain", "sensitivity to light"],
  "duration_days": 2
}

# Symptom checker
GET /api/v1/health/symptom-checker?symptoms=fever,cough,fatigue

# Get health conditions
GET /api/v1/health/conditions?category=cardiovascular
```

### Medication & Supplement Management

```bash
# Add user medication
POST /api/v1/medications/user
{
  "medication_id": "med_123",
  "dosage": "10mg",
  "frequency": "twice daily",
  "start_date": "2024-01-01"
}

# Check drug interactions
POST /api/v1/medications/interactions
{
  "medications": ["metformin", "lisinopril"],
  "supplements": ["vitamin_d"]
}

# Add supplement
POST /api/v1/supplements/user
{
  "supplement_name": "Vitamin D3",
  "dosage": "1000 IU",
  "frequency": "daily"
}
```

### Workout Management

```bash
# Get workout programs
GET /api/v1/workouts/programs?fitness_level=beginner&program_type=strength

# Log workout session
POST /api/v1/workouts/sessions
{
  "workout_program_id": "program_123",
  "duration_minutes": 45,
  "exercises_completed": [...],
  "perceived_exertion": 7
}

# Get workout analytics
GET /api/v1/workouts/analytics?days=30
```

### AI-Powered Nutrition Plan Recommendations

```bash
# Get personalized nutrition plan recommendations
POST /api/v1/nutrition-plans/recommendations
{
  "age": 35,
  "gender": "female",
  "height": 165,
  "weight": 70,
  "activity_level": "moderate",
  "health_conditions": ["hypertension", "prediabetes"],
  "health_goals": ["weight_loss", "heart_health"],
  "current_medications": ["lisinopril"],
  "smoking_status": "never",
  "alcohol_consumption": "light",
  "sleep_hours_per_night": 7,
  "stress_level": 6,
  "exercise_frequency": 3
}

# Response includes:
# - Top 5 recommended plans with scores (0-100)
# - Scientific evidence for each recommendation
# - Personalized rationale and benefits
# - Macro distribution calculations
# - Implementation considerations
# - Medical approval requirements

# Quick nutrition assessment
GET /api/v1/nutrition-plans/quick-assessment?age=30&height=170&weight=65&gender=male&goal=weight_loss

# Compare multiple nutrition plans
POST /api/v1/nutrition-plans/comparison
{
  "plan_types": ["mediterranean", "ketogenic", "dash", "low_carb"],
  "health_assessment": { ... }
}

# Get all available plan types
GET /api/v1/nutrition-plans/types

# Get detailed information about a specific plan
GET /api/v1/nutrition-plans/types/mediterranean

# Create personalized nutrition plan
POST /api/v1/nutrition-plans/personalized
{
  "plan_type": "mediterranean",
  "health_assessment": { ... },
  "duration_weeks": 12
}
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Check coverage threshold (70% minimum)
make check-coverage

# Run security audit
make security-audit
```

## Configuration

Key environment variables:

```bash
# Server
SERVER_PORT=8080
ENVIRONMENT=production

# Database
DB_HOST=localhost
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=your_secure_password

# Security (CRITICAL - Use strong random values)
JWT_SECRET=your_jwt_secret_key_here
API_KEY_SECRET=your_api_key_secret_here
ENCRYPTION_KEY=your_32_character_encryption_key

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
```

## API Key Scopes

- `read_only` - Read access to all endpoints
- `read_write` - Read and write access
- `admin` - Full administrative access
- `nutrition` - Access to nutrition endpoints
- `workouts` - Access to workout endpoints
- `meals` - Access to meal planning endpoints
- `health` - Access to health endpoints
- `supplements` - Access to supplement endpoints

## Monitoring

- **Prometheus Metrics**: `http://localhost:9090/metrics`
- **Grafana Dashboard**: `http://localhost:3001`
- **Health Check**: `http://localhost:8080/health`

## Security Features

- API key authentication with HMAC-SHA256 hashing
- Request signing for sensitive operations
- Rate limiting per API key
- CORS protection
- Security headers (CSP, HSTS, etc.)
- Input validation and sanitization
- SQL injection protection
- XSS protection

## Database Schema

The application uses PostgreSQL in production and SQLite for development/testing. Key tables:

- `users` - User accounts and profiles
- `api_keys` - API key management
- `api_key_usage` - Usage tracking and analytics
- `foods` - Food database with nutritional information
- `exercises` - Exercise database
- `meal_plans` - User meal plans
- `workout_plans` - User workout plans

## Development Tools

- **Hot Reload**: Air for development hot reloading
- **Migrations**: golang-migrate for database migrations
- **Testing**: testify for testing framework
- **Linting**: golangci-lint for code quality
- **Security**: gosec for security scanning

## Deployment

### Docker

```bash
# Build image
docker build -t nutrition-platform-backend .

# Run container
docker run -p 8080:8080 nutrition-platform-backend
```

### Systemd Service

```bash
# Install service
sudo cp scripts/nutrition-platform-backend.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable nutrition-platform-backend
sudo systemctl start nutrition-platform-backend
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run security audit: `make security-audit`
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support or questions:
- Check the [Production Deployment Checklist](../PRODUCTION_DEPLOYMENT_CHECKLIST.md)
- Review the API documentation
- Check application logs: `/var/log/nutrition-platform/`