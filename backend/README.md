# NutriTrack Backend API

A secure, production-ready Node.js/Express API for nutrition tracking with TypeScript, Prisma, PostgreSQL, and Redis.

## 🚀 Features

- **Secure Authentication** with JWT tokens and refresh tokens
- **Rate Limiting** to prevent abuse and DDoS attacks
- **Input Validation** and sanitization for security
- **Comprehensive Error Handling** with proper logging
- **API Documentation** with Swagger/OpenAPI
- **Database Integration** with Prisma ORM
- **Caching** with Redis (with memory fallback)
- **Health Checks** for monitoring
- **CORS Configuration** for frontend integration
- **Security Headers** with Helmet
- **Compression** for better performance
- **Request Logging** with Morgan

## 📋 Prerequisites

- Node.js 18+ 
- PostgreSQL 12+
- Redis 6+ (optional)
- npm or yarn

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd nutritrack-backend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Configure your `.env` file with the required variables:
   ```env
   # Database
   DATABASE_URL="postgresql://username:password@localhost:5432/nutritrack"
   
   # JWT
   JWT_SECRET="your-super-secret-jwt-key"
   JWT_REFRESH_SECRET="your-super-secret-refresh-key"
   SESSION_SECRET="your-session-secret"
   
   # Redis (optional)
   REDIS_URL="redis://localhost:6379"
   
   # Server
   NODE_ENV="development"
   PORT="3001"
   CORS_ORIGIN="http://localhost:3000"
   ```

4. **Set up the database**
   ```bash
   # Generate Prisma client
   npm run generate
   
   # Run database migrations
   npm run migrate
   
   # Seed the database (optional)
   npm run seed
   ```

## 🏃‍♂️ Running the Application

### Development Mode
```bash
npm run dev
```

### Production Mode
```bash
# Build the application
npm run build

# Start the server
npm start
```

### Using Docker
```bash
# Build the Docker image
npm run docker:build

# Run with Docker Compose
npm run docker:run
```

## 📚 API Documentation

Once the server is running, you can access the API documentation at:
- **Swagger UI**: http://localhost:3001/api/docs
- **OpenAPI JSON**: http://localhost:3001/api/docs.json

## 🔗 API Endpoints

### Health Checks
- `GET /health` - Basic health check
- `GET /api/health` - Detailed health status
- `GET /api/health/ready` - Readiness probe
- `GET /api/health/live` - Liveness probe

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password
- `GET /api/auth/me` - Get current user profile

### API Base URL
- **Development**: http://localhost:3001/api
- **Production**: https://your-domain.com/api

## 🛡️ Security Features

### Authentication & Authorization
- JWT access tokens with configurable expiration
- Refresh tokens for extended sessions
- Role-based access control
- Password hashing with bcrypt

### Rate Limiting
- General API rate limiting (100 requests per 15 minutes)
- Strict rate limiting for sensitive endpoints (5 requests per 15 minutes)
- Password reset rate limiting (3 attempts per hour)
- File upload rate limiting (20 uploads per hour)

### Input Validation & Sanitization
- Request body validation with express-validator
- Joi schema validation for complex objects
- Input sanitization to prevent XSS attacks
- Content-Type validation
- Request size validation

### Security Headers
- Helmet.js for security headers
- CORS configuration
- HSTS (HTTP Strict Transport Security)
- Content Security Policy
- X-Frame-Options, X-Content-Type-Options

## 📊 Monitoring & Logging

### Logging
- Winston logger with multiple levels
- Log rotation in production
- Request logging with Morgan
- Error logging with context

### Health Monitoring
- Database connectivity checks
- Redis connectivity checks
- Memory and CPU usage monitoring
- Custom health endpoints for orchestration

## 🗄️ Database

### Schema
The application uses Prisma ORM with PostgreSQL. The database schema includes:
- Users and profiles
- Nutrition goals
- Meals and food logs
- Progress records
- Body measurements
- Custom foods
- Water intake
- User devices

### Migrations
```bash
# Create new migration
npx prisma migrate dev --name <migration-name>

# Apply migrations in production
npm run migrate:deploy

# Reset database
npx prisma migrate reset
```

## 🔄 Caching

### Redis Configuration
- Automatic Redis connection with reconnection logic
- Memory cache fallback when Redis is unavailable
- TTL-based cache expiration
- Pattern-based cache invalidation

### Cache Usage
- User session data
- Frequently accessed food data
- API responses with configurable TTL

## 🧪 Testing

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```

## 📝 Code Quality

```bash
# Run ESLint
npm run lint

# Fix ESLint issues
npm run lint:fix

# Format code with Prettier
npm run format
```

## 🚀 Deployment

### Environment Setup
1. Set production environment variables
2. Configure production database
3. Set up Redis cluster (optional)
4. Configure SSL certificates

### Production Build
```bash
npm run build
npm start
```

### Docker Deployment
```bash
# Build production image
docker build -t nutritrack-backend:latest .

# Run container
docker run -p 3001:3001 --env-file .env nutritrack-backend:latest
```

## 🔧 Configuration

### Environment Variables
See `.env.example` for all available configuration options.

### Key Configuration Areas
- Database connection
- Redis connection
- JWT secrets
- CORS settings
- Rate limiting
- File uploads
- Email settings
- Logging configuration

## 📈 Performance

### Optimization Features
- Response compression
- Database connection pooling
- Redis caching
- Request size limits
- Memory usage monitoring

### Monitoring
- Health check endpoints
- Request logging
- Error tracking
- Performance metrics

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check DATABASE_URL in .env
   - Ensure PostgreSQL is running
   - Verify database credentials

2. **Redis Connection Failed**
   - Check REDIS_URL in .env
   - Redis is optional - app will continue without it
   - Check Redis server status

3. **JWT Token Errors**
   - Verify JWT_SECRET and JWT_REFRESH_SECRET
   - Check token expiration
   - Ensure proper token format

4. **CORS Issues**
   - Verify CORS_ORIGIN in .env
   - Check frontend URL configuration

### Logs
- Development: Console output
- Production: Check log files in `logs/` directory

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and linting
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the API documentation
- Review the troubleshooting section
