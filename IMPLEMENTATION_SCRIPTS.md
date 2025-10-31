# 🛠️ Implementation Scripts & Code Examples

This document contains all the necessary scripts and code examples to implement the bug fixes and build the Next.js frontend with Node.js backend integration.

## 🐛 Critical Bug Fixes

### 1.1 Fix Critical Bugs Script

```bash
#!/bin/bash
# nutrition-platform/fix-critical-bugs.sh

echo "🔧 Fixing critical bugs in nutrition platform..."

# Check if we're in the right directory
if [ ! -d "production-nodejs" ] || [ ! -d "frontend" ]; then
    echo "❌ Error: Please run this script from the nutrition-platform root directory"
    exit 1
fi

# Create backup of original files
echo "📦 Creating backups..."
cp production-nodejs/server.js production-nodejs/server.js.backup
cp production-nodejs/services/redisClient.js production-nodejs/services/redisClient.js.backup
cp frontend/src/js/app.js frontend/src/js/app.js.backup

# Fix 1: Add missing prom-client import to server.js
echo "🔧 Fixing prom-client import in server.js..."
sed -i '' '/const express = require/a\
const { register } = require("prom-client");\
\
const express = require("express");' production-nodejs/server.js

# Fix 2: Add logs directory creation to server.js
echo "🔧 Adding logs directory creation to server.js..."
sed -i '' '/const helmet = require/a\
const fs = require("fs");\
const path = require("path");\
\
// Ensure logs directory exists\
const logsDir = path.join(__dirname, "logs");\
if (!fs.existsSync(logsDir)) {\
  fs.mkdirSync(logsDir, { recursive: true });\
}\
\
const helmet = require("helmet");' production-nodejs/server.js

# Fix 3: Fix Redis client connection issue
echo "🔧 Fixing Redis client connection issue..."
sed -i '' '/await this.client.connect();/d' production-nodejs/services/redisClient.js

# Fix 4: Fix frontend API endpoint
echo "🔧 Fixing frontend API endpoint..."
sed -i '' 's|http://localhost:8080/api/v1|http://localhost:8080/api|' frontend/src/js/app.js

echo "✅ Critical bugs fixed successfully!"
echo "📦 Backups created with .backup extension"
echo ""
echo "🚀 To test the fixes:"
echo "   cd production-nodejs && npm start"
echo "   # In another terminal:"
echo "   cd frontend && python -m http.server 8081"
echo ""
echo "🔄 To restore original files:"
echo "   cp production-nodejs/server.js.backup production-nodejs/server.js"
echo "   cp production-nodejs/services/redisClient.js.backup production-nodejs/services/redisClient.js"
echo "   cp frontend/src/js/app.js.backup frontend/src/js/app.js"
```

### 1.2 Enhanced Backend Files

#### Enhanced Logger (production-nodejs/src/config/logger.js)

```javascript
// production-nodejs/src/config/logger.js - Enhanced logging
const winston = require('winston');
const path = require('path');

// Ensure logs directory exists
const logsDir = path.join(__dirname, '../../logs');
if (!require('fs').existsSync(logsDir)) {
  require('fs').mkdirSync(logsDir, { recursive: true });
}

// Custom format for structured logging
const logFormat = winston.format.combine(
  winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
  winston.format.errors({ stack: true }),
  winston.format.printf(({ timestamp, level, message, ...meta }) => {
    let log = `${timestamp} [${level.toUpperCase()}]: ${message}`;
    if (Object.keys(meta).length > 0) {
      log += ` ${JSON.stringify(meta)}`;
    }
    return log;
  })
);

// Create logger instance
const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: logFormat,
  transports: [
    // Console transport for development
    new winston.transports.Console({
      format: winston.format.combine(
        winston.format.colorize(),
        logFormat
      )
    }),
    // File transport for production
    new winston.transports.File({
      filename: path.join(logsDir, 'error.log'),
      level: 'error',
      format: logFormat
    }),
    new winston.transports.File({
      filename: path.join(logsDir, 'combined.log'),
      format: logFormat
    })
  ]
});

// Stream for Morgan middleware
logger.stream = {
  write: (message) => {
    logger.info(message.trim());
  }
};

module.exports = logger;
```

#### Enhanced API Service (production-nodejs/services/apiService.js)

```javascript
// production-nodejs/services/apiService.js
const logger = require('../config/logger');
const nutritionService = require('./nutritionService');

class ApiService {
  constructor() {
    this.requestTimeout = 30000;
    this.maxRetries = 3;
  }

  async analyzeNutrition(food, quantity, unit, checkHalal = false) {
    const requestId = this.generateRequestId();
    const startTime = Date.now();

    try {
      logger.info('Nutrition analysis started', {
        requestId,
        food,
        quantity,
        unit,
        checkHalal
      });

      // Validate inputs
      this.validateNutritionInput({ food, quantity, unit });

      // Perform analysis
      const result = await nutritionService.analyze(food, quantity, unit, checkHalal);
      
      const duration = Date.now() - startTime;
      
      logger.info('Nutrition analysis completed', {
        requestId,
        duration: `${duration}ms`,
        food,
        result: {
          calories: result.calories,
          protein: result.protein,
          isHalal: result.isHalal
        }
      });

      return {
        ...result,
        requestId,
        processingTime: duration
      };

    } catch (error) {
      const duration = Date.now() - startTime;
      
      logger.error('Nutrition analysis failed', {
        requestId,
        duration: `${duration}ms`,
        food,
        error: error.message,
        stack: error.stack
      });

      throw new Error(`Failed to analyze nutrition: ${error.message}`);
    }
  }

  validateNutritionInput({ food, quantity, unit }) {
    if (!food || typeof food !== 'string' || food.trim().length === 0) {
      throw new Error('Food name is required and must be a non-empty string');
    }

    if (!quantity || isNaN(quantity) || quantity <= 0 || quantity > 100000) {
      throw new Error('Quantity must be a number between 0.1 and 100000');
    }

    const validUnits = ['g', 'kg', 'oz', 'lb'];
    if (!unit || !validUnits.includes(unit)) {
      throw new Error(`Unit must be one of: ${validUnits.join(', ')}`);
    }
  }

  generateRequestId() {
    return Math.random().toString(36).substring(2, 15);
  }

  async healthCheck() {
    try {
      const result = await nutritionService.getAvailableFoods();
      return {
        status: 'healthy',
        timestamp: new Date().toISOString(),
        uptime: process.uptime(),
        availableFoods: result.length,
        memory: process.memoryUsage(),
        cpu: process.cpuUsage()
      };
    } catch (error) {
      logger.error('Health check failed', { error: error.message });
      throw error;
    }
  }
}

module.exports = new ApiService();
```

## 🏗️ Next.js Frontend Setup

### 2.1 Next.js Project Setup Script

```bash
#!/bin/bash
# nutrition-platform/setup-nextjs-frontend.sh

echo "🚀 Setting up Next.js frontend for nutrition platform..."

# Check if we're in the right directory
if [ ! -d "production-nodejs" ]; then
    echo "❌ Error: Please run this script from the nutrition-platform root directory"
    exit 1
fi

# Create Next.js app
echo "📦 Creating Next.js application..."
npx create-next-app@latest frontend-nextjs \
  --typescript \
  --tailwind \
  --eslint \
  --app \
  --src-dir \
  --import-alias "@/*"

echo "📦 Installing additional dependencies..."
cd frontend-nextjs

# Install required packages
npm install axios zod react-hook-form @hookform/resolvers pino next-auth

# Install dev dependencies
npm install -D @types/node

echo "📁 Creating directory structure..."
mkdir -p src/components/ui src/components/forms src/components/features src/components/providers
mkdir -p src/lib/api/services src/lib/logger src/lib/utils src/lib/validations
mkdir -p src/types src/hooks src/tests

echo "✅ Next.js frontend setup completed!"
echo ""
echo "🚀 Next steps:"
echo "   cd frontend-nextjs"
echo "   npm run dev"
echo ""
echo "📝 After setup, implement the files from IMPLEMENTATION_SCRIPTS.md"
```

### 2.2 Next.js Configuration Files

#### next.config.js

```javascript
// frontend-nextjs/next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,
  
  // Environment variables
  env: {
    API_VERSION: 'v1',
  },

  // API proxy for development
  async rewrites() {
    if (process.env.NODE_ENV === 'development') {
      return [
        {
          source: '/api/:path*',
          destination: 'http://localhost:8080/api/:path*',
        },
      ];
    }
    return [];
  },

  // Security headers
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'X-DNS-Prefetch-Control',
            value: 'on',
          },
          {
            key: 'Strict-Transport-Security',
            value: 'max-age=63072000; includeSubDomains; preload',
          },
          {
            key: 'X-Frame-Options',
            value: 'SAMEORIGIN',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'X-XSS-Protection',
            value: '1; mode=block',
          },
        ],
      },
    ];
  },

  // Image optimization
  images: {
    domains: ['localhost'],
    formats: ['image/avif', 'image/webp'],
  },

  // Experimental features
  experimental: {
    serverActions: true,
    optimizeCss: true,
  },
};

module.exports = nextConfig;
```

#### Environment Configuration (.env.example)

```bash
# frontend-nextjs/.env.example
NODE_ENV=development
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Backend API configuration
API_URL=http://localhost:8080/api
API_TIMEOUT=30000

# Logging
NEXT_PUBLIC_LOG_LEVEL=debug
LOG_LEVEL=info

# Authentication
NEXTAUTH_SECRET=your-super-secret-key-min-32-chars-long
NEXTAUTH_URL=http://localhost:3000

# Feature flags
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_DEBUG=true
```

## 🔐 API Integration & Authentication

### 3.1 API Client Setup

#### API Client (frontend-nextjs/src/lib/api/client.ts)

```typescript
// frontend-nextjs/src/lib/api/client.ts
import axios, { AxiosError, AxiosInstance } from 'axios';
import { config } from '../env';
import { logger } from '../logger';

interface ApiError {
  message: string;
  code?: string;
  status?: number;
  errors?: Record<string, string[]>;
}

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: config.api.url,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true,
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        const requestId = crypto.randomUUID();
        config.headers['X-Request-ID'] = requestId;

        logger.info('API Request', {
          requestId,
          method: config.method?.toUpperCase(),
          url: config.url,
          data: this.sanitizeLogData(config.data),
        });

        return config;
      },
      (error) => {
        logger.error('Request setup error', { error: error.message });
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => {
        logger.info('API Response Success', {
          requestId: response.config.headers['X-Request-ID'],
          method: response.config.method?.toUpperCase(),
          url: response.config.url,
          status: response.status,
        });

        return response;
      },
      async (error: AxiosError<ApiError>) => {
        logger.error('API Response Error', {
          requestId: error.config?.headers['X-Request-ID'],
          method: error.config?.method?.toUpperCase(),
          url: error.config?.url,
          status: error.response?.status,
          errorMessage: error.response?.data?.message || error.message,
        });

        return Promise.reject(this.normalizeError(error));
      }
    );
  }

  private normalizeError(error: AxiosError<ApiError>): ApiError {
    const data = error.response?.data;
    return {
      message: data?.message || error.message || 'An unexpected error occurred',
      code: data?.code || 'UNKNOWN_ERROR',
      status: error.response?.status,
      errors: data?.errors,
    };
  }

  private sanitizeLogData(data: any): any {
    if (!data) return data;

    const sensitiveFields = ['password', 'token', 'secret', 'apiKey'];
    const sanitized = { ...data };

    for (const field of sensitiveFields) {
      if (field in sanitized) {
        sanitized[field] = '***REDACTED***';
      }
    }

    return sanitized;
  }

  async get<T>(url: string, params?: Record<string, any>): Promise<T> {
    const response = await this.client.get<T>(url, { params });
    return response.data;
  }

  async post<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.post<T>(url, data);
    return response.data;
  }
}

export const apiClient = new ApiClient();
```

#### Environment Configuration (frontend-nextjs/src/lib/env.ts)

```typescript
// frontend-nextjs/src/lib/env.ts
import { z } from 'zod';

const envSchema = z.object({
  NODE_ENV: z.enum(['development', 'production', 'test']),
  NEXT_PUBLIC_APP_URL: z.string().url(),
  NEXT_PUBLIC_API_URL: z.string().url(),
  API_URL: z.string().url().optional(),
  NEXT_PUBLIC_LOG_LEVEL: z.enum(['trace', 'debug', 'info', 'warn', 'error']),
  LOG_LEVEL: z.enum(['trace', 'debug', 'info', 'warn', 'error']).optional(),
  NEXTAUTH_SECRET: z.string().min(32),
  NEXTAUTH_URL: z.string().url(),
});

const processEnv = {
  NODE_ENV: process.env.NODE_ENV,
  NEXT_PUBLIC_APP_URL: process.env.NEXT_PUBLIC_APP_URL,
  NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
  API_URL: process.env.API_URL,
  NEXT_PUBLIC_LOG_LEVEL: process.env.NEXT_PUBLIC_LOG_LEVEL,
  LOG_LEVEL: process.env.LOG_LEVEL,
  NEXTAUTH_SECRET: process.env.NEXTAUTH_SECRET,
  NEXTAUTH_URL: process.env.NEXTAUTH_URL,
};

export const env = envSchema.parse(processEnv);

export const config = {
  app: {
    url: env.NEXT_PUBLIC_APP_URL,
    env: env.NODE_ENV,
  },
  api: {
    url: env.NEXT_PUBLIC_API_URL,
    serverUrl: env.API_URL || env.NEXT_PUBLIC_API_URL,
  },
  logging: {
    level: env.NEXT_PUBLIC_LOG_LEVEL,
  },
  auth: {
    secret: env.NEXTAUTH_SECRET,
    url: env.NEXTAUTH_URL,
  },
} as const;
```

### 3.2 Nutrition Service

```typescript
// frontend-nextjs/src/lib/api/services/nutrition.service.ts
import { apiClient } from '../client';
import { logger } from '../../logger';
import { z } from 'zod';

export const NutritionAnalysisSchema = z.object({
  food: z.string().min(1, 'Food name is required'),
  quantity: z.number().min(0.1, 'Quantity must be at least 0.1').max(100000, 'Quantity cannot exceed 100000'),
  unit: z.enum(['g', 'kg', 'oz', 'lb']),
  checkHalal: z.boolean().optional(),
});

export type NutritionAnalysisRequest = z.infer<typeof NutritionAnalysisSchema>;

export interface NutritionAnalysisResult {
  food: string;
  quantity: number;
  unit: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
  sugar: number;
  sodium: number;
  calcium: number;
  iron: number;
  vitaminC: number;
  vitaminA: number;
  isHalal: boolean | null;
  status: string;
  message: string;
  timestamp: string;
  requestId: string;
  processingTime: number;
}

export const nutritionService = {
  async analyzeNutrition(request: NutritionAnalysisRequest): Promise<NutritionAnalysisResult> {
    // Validate input with Zod
    const validatedRequest = NutritionAnalysisSchema.parse(request);
    
    logger.info('Analyzing nutrition', { 
      food: validatedRequest.food, 
      quantity: validatedRequest.quantity 
    });

    try {
      const result = await apiClient.post<NutritionAnalysisResult>('/nutrition/analyze', validatedRequest);
      
      logger.info('Nutrition analysis completed', {
        requestId: result.requestId,
        food: result.food,
        calories: result.calories,
        processingTime: result.processingTime
      });

      return result;
    } catch (error: any) {
      logger.error('Nutrition analysis failed', {
        food: validatedRequest.food,
        error: error.message
      });
      throw error;
    }
  },

  async getAvailableFoods(): Promise<string[]> {
    try {
      const result = await apiClient.get<{ foods: string[] }>('/nutrition/foods');
      return result.foods;
    } catch (error: any) {
      logger.error('Failed to get available foods', { error: error.message });
      throw error;
    }
  },
};
```

## 📊 Logging Setup

### 4.1 Frontend Logger

```typescript
// frontend-nextjs/src/lib/logger/index.ts
import pino from 'pino';

// Browser logger
export const logger = typeof window !== 'undefined' 
  ? pino({
      level: process.env.NEXT_PUBLIC_LOG_LEVEL || 'info',
      browser: {
        asObject: true,
        transmit: {
          level: 'error',
          send: async (level, logEvent) => {
            // Send errors to backend logging endpoint
            try {
              await fetch('/api/logs', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                  level,
                  messages: logEvent.messages,
                  bindings: logEvent.bindings,
                  timestamp: new Date().toISOString(),
                }),
              });
            } catch (error) {
              console.error('Failed to send log to server:', error);
            }
          },
        },
      },
    })
  : pino({
      level: process.env.LOG_LEVEL || 'info',
      transport: {
        target: 'pino-pretty',
        options: {
          colorize: true,
          translateTime: 'SYS:standard',
          ignore: 'pid,hostname',
        },
      },
    });

// Create child loggers for different contexts
export const createLogger = (context: string) => {
  return logger.child({ context });
};

// Structured logging helpers
export const loggers = {
  api: createLogger('api'),
  auth: createLogger('auth'),
  ui: createLogger('ui'),
  error: createLogger('error'),
  nutrition: createLogger('nutrition'),
};
```

## 🐳 Docker Configuration

### 5.1 Docker Compose for Next.js + Node.js

```yaml
# nutrition-platform/docker-compose.nextjs.yml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: nutrition_postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${DB_NAME:-nutrition_platform}
      POSTGRES_USER: ${DB_USER:-nutrition_user}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-nutrition_password_2024}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - nutrition_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-nutrition_user} -d ${DB_NAME:-nutrition_platform}"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: nutrition_redis
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD:-redis_password_2024}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - nutrition_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Node.js Backend API
  backend:
    build:
      context: ./production-nodejs
      dockerfile: Dockerfile
    container_name: nutrition_backend
    restart: unless-stopped
    environment:
      - NODE_ENV=production
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=${DB_NAME:-nutrition_platform}
      - DB_USER=${DB_USER:-nutrition_user}
      - DB_PASSWORD=${DB_PASSWORD:-nutrition_password_2024}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD:-redis_password_2024}
      - JWT_SECRET=${JWT_SECRET:-super_secret_jwt_key_2024}
      - LOG_LEVEL=info
      - CORS_ORIGINS=http://localhost:3000,http://frontend:3000
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "8080:8080"
    networks:
      - nutrition_network
    volumes:
      - ./production-nodejs/logs:/app/logs
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Next.js Frontend
  frontend:
    build:
      context: ./frontend-nextjs
      dockerfile: Dockerfile
      args:
        - NEXT_PUBLIC_API_URL=http://backend:8080/api
    container_name: nutrition_frontend
    restart: unless-stopped
    environment:
      - NODE_ENV=production
      - NEXT_PUBLIC_API_URL=http://backend:8080/api
      - NEXTAUTH_SECRET=${NEXTAUTH_SECRET:-nextauth_secret_2024}
      - NEXTAUTH_URL=http://localhost:3000
    ports:
      - "3000:3000"
    depends_on:
      - backend
    networks:
      - nutrition_network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Nginx Reverse Proxy
  nginx:
    image: nginx:alpine
    container_name: nutrition_nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.nextjs.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/ssl:ro
    depends_on:
      - frontend
      - backend
    networks:
      - nutrition_network

volumes:
  postgres_data:
  redis_data:

networks:
  nutrition_network:
    driver: bridge
```

### 5.2 Backend Dockerfile

```dockerfile
# production-nodejs/Dockerfile
FROM node:18-alpine AS base

# Install dependencies
FROM base AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# Create logs directory
RUN mkdir -p logs

# Copy source code
FROM base AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

# Production image
FROM base AS production
WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Copy production dependencies
COPY --from=deps --chown=nodejs:nodejs /app/node_modules ./node_modules

# Copy application code
COPY --chown=nodejs:nodejs . .

# Create logs directory
RUN mkdir -p logs && chown nodejs:nodejs logs

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD node -e "require('http').get('http://localhost:8080/health', (res) => process.exit(res.statusCode === 200 ? 0 : 1))"

# Start the application
CMD ["node", "server.js"]
```

### 5.3 Frontend Dockerfile

```dockerfile
# frontend-nextjs/Dockerfile
FROM node:18-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Set build-time environment variables
ENV NEXT_TELEMETRY_DISABLED 1
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN npm run build

# Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
```

## 🧪 Testing Examples

### 6.1 Backend Tests

```javascript
// production-nodejs/tests/api.test.js
const request = require('supertest');
const app = require('../server');
const { expect } = require('chai');

describe('Nutrition API', () => {
  describe('POST /api/nutrition/analyze', () => {
    it('should analyze nutrition for valid food', async () => {
      const response = await request(app)
        .post('/api/nutrition/analyze')
        .send({
          food: 'apple',
          quantity: 100,
          unit: 'g',
          checkHalal: true
        })
        .expect(200);

      expect(response.body).to.have.property('status', 'success');
      expect(response.body).to.have.property('food', 'apple');
      expect(response.body).to.have.property('calories');
      expect(response.body).to.have.property('isHalal');
      expect(response.body).to.have.property('requestId');
    });

    it('should return 400 for invalid input', async () => {
      const response = await request(app)
        .post('/api/nutrition/analyze')
        .send({
          food: '',
          quantity: -10,
          unit: 'invalid'
        })
        .expect(400);

      expect(response.body).to.have.property('status', 'error');
    });
  });

  describe('GET /health', () => {
    it('should return health status', async () => {
      const response = await request(app)
        .get('/health')
        .expect(200);

      expect(response.body).to.have.property('status', 'healthy');
      expect(response.body).to.have.property('uptime');
    });
  });
});
```

### 6.2 Frontend Tests

```typescript
// frontend-nextjs/src/components/features/NutritionAnalysis/__tests__/NutritionForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { NutritionForm } from '../NutritionForm';
import { nutritionService } from '@/lib/api/services/nutrition.service';

// Mock the nutrition service
jest.mock('@/lib/api/services/nutrition.service');
const mockedNutritionService = nutritionService as jest.Mocked<typeof nutritionService>;

describe('NutritionForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders form correctly', () => {
    render(<NutritionForm />);
    
    expect(screen.getByLabelText('Food Item')).toBeInTheDocument();
    expect(screen.getByLabelText('Quantity')).toBeInTheDocument();
    expect(screen.getByLabelText('Unit')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Analyze Nutrition' })).toBeInTheDocument();
  });

  it('submits form successfully', async () => {
    const user = userEvent.setup();
    const mockResult = {
      food: 'apple',
      quantity: 100,
      unit: 'g',
      calories: 52,
      protein: 0.3,
      carbs: 14,
      fat: 0.2,
      fiber: 2.4,
      sugar: 10.4,
      isHalal: true,
      status: 'success',
      message: 'Nutrition analysis completed successfully',
      requestId: 'test-request-id',
      processingTime: 45,
    };

    mockedNutritionService.analyzeNutrition.mockResolvedValue(mockResult);

    render(<NutritionForm />);

    await user.type(screen.getByLabelText('Food Item'), 'apple');
    await user.type(screen.getByLabelText('Quantity'), '100');
    await user.click(screen.getByRole('button', { name: 'Analyze Nutrition' }));

    await waitFor(() => {
      expect(screen.getByText('Analysis Results')).toBeInTheDocument();
      expect(screen.getByText('Food: apple (100g)')).toBeInTheDocument();
      expect(screen.getByText('Calories: 52.0 kcal')).toBeInTheDocument();
    });

    expect(mockedNutritionService.analyzeNutrition).toHaveBeenCalledWith({
      food: 'apple',
      quantity: 100,
      unit: 'g',
      checkHalal: true,
    });
  });
});
```

## 🚀 Deployment Script

### 7.1 Complete Deployment Script

```bash
#!/bin/bash
# nutrition-platform/deploy-nextjs-nodejs.sh

echo "🚀 Deploying Nutrition Platform with Next.js + Node.js..."

# Check if we're in the right directory
if [ ! -d "production-nodejs" ] || [ ! -d "frontend-nextjs" ]; then
    echo "❌ Error: Please run this script from the nutrition-platform root directory"
    exit 1
fi

# Create environment file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file..."
    cat > .env << EOF
# Database Configuration
DB_NAME=nutrition_platform
DB_USER=nutrition_user
DB_PASSWORD=nutrition_password_2024

# Redis Configuration
REDIS_PASSWORD=redis_password_2024

# JWT Configuration
JWT_SECRET=super_secret_jwt_key_2024

# NextAuth Configuration
NEXTAUTH_SECRET=nextauth_secret_2024

# CORS Configuration
CORS_ORIGINS=http://localhost:3000
EOF
fi

# Fix critical bugs
echo "🔧 Fixing critical bugs..."
chmod +x fix-critical-bugs.sh
./fix-critical-bugs.sh

# Build and start services
echo "🐳 Building and starting services..."
docker-compose -f docker-compose.nextjs.yml down
docker-compose -f docker-compose.nextjs.yml build
docker-compose -f docker-compose.nextjs.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
docker-compose -f docker-compose.nextjs.yml ps

# Test API endpoint
echo "🧪 Testing API endpoint..."
curl -f http://localhost/api/health || echo "❌ API health check failed"

# Test frontend
echo "🧪 Testing frontend..."
curl -f http://localhost:3000 || echo "❌ Frontend health check failed"

echo ""
echo "✅ Deployment completed!"
echo ""
echo "🌐 Access your application at:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost/api"
echo "   Health Check: http://localhost/health"
echo ""
echo "📊 To view logs:"
echo "   docker-compose -f docker-compose.nextjs.yml logs -f backend"
echo "   docker-compose -f docker-compose.nextjs.yml logs -f frontend"
echo ""
echo "🛑 To stop services:"
echo "   docker-compose -f docker-compose.nextjs.yml down"
```

## 📋 Implementation Checklist

### Backend Fixes
- [ ] Fix missing prom-client import in server.js
- [ ] Create logs directory on startup
- [ ] Fix Redis client connection issue
- [ ] Implement enhanced logging with structured logs
- [ ] Add comprehensive error handling
- [ ] Implement API service with validation

### Frontend Setup
- [ ] Create Next.js project with TypeScript
- [ ] Install required dependencies
- [ ] Set up API client with interceptors
- [ ] Implement nutrition service with Zod validation
- [ ] Create logger with context-specific loggers
- [ ] Build nutrition analysis form component
- [ ] Set up error boundary with logging

### Integration
- [ ] Configure API proxy for development
- [ ] Set up environment variables
- [ ] Implement proper error handling
- [ ] Add request/response logging
- [ ] Set up authentication flow
- [ ] Test API integration

### Docker & Deployment
- [ ] Create Dockerfiles for both frontend and backend
- [ ] Configure Docker Compose with all services
- [ ] Set up Nginx reverse proxy
- [ ] Configure health checks
- [ ] Set up proper networking
- [ ] Test complete deployment

### Testing
- [ ] Write unit tests for API endpoints
- [ ] Write component tests for frontend
- [ ] Test integration between frontend and backend
- [ ] Test error handling scenarios
- [ ] Test logging functionality

This document provides all the necessary code and scripts to implement a complete Next.js + Node.js nutrition platform with proper bug fixes, logging, and deployment configuration.