# 🚀 Comprehensive Plan: Fix Bugs & Build Next.js Frontend with Node.js Integration

## 📋 Executive Summary

This document provides a complete implementation plan to fix critical bugs in the existing nutrition platform and build a modern Next.js frontend with proper Node.js backend integration. The plan addresses all identified issues and follows best practices from the reference documents.

## 🐛 Phase 1: Critical Bug Fixes (Immediate)

### 1.1 Fix Node.js Backend Critical Issues

#### Issue 1: Missing prom-client import
**File:** `nutrition-platform/production-nodejs/server.js:95`
```javascript
// Fix: Add missing import at the top of server.js
const { register } = require('prom-client');
```

#### Issue 2: Missing logs directory
**File:** `nutrition-platform/production-nodejs/server.js`
```javascript
// Fix: Add after imports
const fs = require('fs');
const path = require('path');

// Ensure logs directory exists
const logsDir = path.join(__dirname, 'logs');
if (!fs.existsSync(logsDir)) {
  fs.mkdirSync(logsDir, { recursive: true });
}
```

#### Issue 3: Redis connection issue
**File:** `nutrition-platform/production-nodejs/services/redisClient.js:41`
```javascript
// Fix: Remove the connect() call - ioredis connects automatically
// Remove this line:
await this.client.connect();
```

#### Issue 4: Frontend API endpoint mismatch
**File:** `nutrition-platform/frontend/src/js/app.js:4`
```javascript
// Fix: Change from /api/v1 to /api
this.apiBaseUrl = 'http://localhost:8080/api';
```

### 1.2 Create Quick Fix Script

```bash
#!/bin/bash
# nutrition-platform/fix-critical-bugs.sh

echo "🔧 Fixing critical bugs..."

# Fix server.js
sed -i '' 's/const { express } = require("express");/const { express } = require("express");\nconst { register } = require("prom-client");/' production-nodejs/server.js

# Add logs directory creation
sed -i '' '/const app = express();/i\
const fs = require("fs");\
const path = require("path");\
\
const logsDir = path.join(__dirname, "logs");\
if (!fs.existsSync(logsDir)) {\
  fs.mkdirSync(logsDir, { recursive: true });\
}\
' production-nodejs/server.js

# Fix Redis client
sed -i '' '/await this.client.connect();/d' production-nodejs/services/redisClient.js

# Fix frontend API URL
sed -i '' 's|http://localhost:8080/api/v1|http://localhost:8080/api|' frontend/src/js/app.js

echo "✅ Critical bugs fixed!"
```

## 🏗️ Phase 2: Next.js Frontend Architecture

### 2.1 Project Structure

```
nutrition-platform/
├── frontend-nextjs/              # New Next.js frontend
│   ├── src/
│   │   ├── app/                  # App Router
│   │   │   ├── (auth)/           # Route groups
│   │   │   │   ├── login/
│   │   │   │   │   └── page.tsx
│   │   │   │   └── register/
│   │   │   │       └── page.tsx
│   │   │   ├── (dashboard)/      # Protected routes
│   │   │   │   ├── dashboard/
│   │   │   │   │   └── page.tsx
│   │   │   │   ├── nutrition/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── analyze/
│   │   │   │   │       └── page.tsx
│   │   │   │   └── profile/
│   │   │   │       └── page.tsx
│   │   │   ├── api/              # API routes (proxy)
│   │   │   │   └── auth/
│   │   │   │       └── [...nextauth]/
│   │   │   ├── globals.css
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── components/
│   │   │   ├── ui/              # Reusable UI
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   └── index.ts
│   │   │   ├── forms/            # Form components
│   │   │   │   ├── LoginForm.tsx
│   │   │   │   ├── NutritionForm.tsx
│   │   │   │   └── RegisterForm.tsx
│   │   │   ├── features/         # Feature-specific
│   │   │   │   ├── NutritionAnalysis/
│   │   │   │   ├── Dashboard/
│   │   │   │   └── UserProfile/
│   │   │   └── providers/        # Context providers
│   │   │       ├── AuthProvider.tsx
│   │   │       ├── ThemeProvider.tsx
│   │   │       └── ErrorBoundary.tsx
│   │   ├── lib/
│   │   │   ├── api/             # API client
│   │   │   │   ├── client.ts
│   │   │   │   └── services/
│   │   │   │       ├── nutrition.service.ts
│   │   │   │       └── auth.service.ts
│   │   │   ├── logger/          # Logging setup
│   │   │   │   └── index.ts
│   │   │   ├── utils/           # Utilities
│   │   │   │   ├── validations.ts
│   │   │   │   ├── helpers.ts
│   │   │   │   └── constants.ts
│   │   │   ├── validations/     # Zod schemas
│   │   │   │   └── nutrition.schema.ts
│   │   │   └── env.ts           # Environment config
│   │   ├── types/               # TypeScript types
│   │   │   ├── api.ts
│   │   │   ├── auth.ts
│   │   │   └── nutrition.ts
│   │   ├── hooks/               # Custom hooks
│   │   │   ├── useAuth.ts
│   │   │   ├── useNutrition.ts
│   │   │   └── useLocalStorage.ts
│   │   └── middleware.ts        # Next.js middleware
│   ├── public/
│   │   ├── icons/
│   │   └── images/
│   ├── tests/
│   │   ├── __mocks__/
│   │   ├── components/
│   │   └── pages/
│   ├── .env.local
│   ├── .env.example
│   ├── next.config.js
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   └── docker-compose.dev.yml
├── production-nodejs/             # Existing backend (fixed)
├── docker-compose.nextjs.yml     # New docker config
├── nginx.nextjs.conf             # Updated nginx config
└── DEPLOYMENT_GUIDE.md           # Updated deployment guide
```

### 2.2 Next.js Configuration

```typescript
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

### 2.3 Environment Configuration

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

### 2.4 TypeScript Configuration

```json
// frontend-nextjs/tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

## 🔧 Phase 3: Enhanced Node.js Backend

### 3.1 Improved Backend Structure

```typescript
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

### 3.2 Enhanced API Service

```typescript
// production-nodejs/services/apiService.js
const logger = require('../utils/logger');
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

## 🔐 Phase 4: Authentication & API Integration

### 4.1 Next.js API Client with Logging

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

### 4.2 Nutrition Service Integration

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

## 📊 Phase 5: Logging & Error Handling

### 5.1 Frontend Logger Setup

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

### 5.2 Error Boundary with Logging

```typescript
// frontend-nextjs/src/components/providers/ErrorBoundary.tsx
'use client';

import { Component, ErrorInfo, ReactNode } from 'react';
import { loggers } from '@/lib/logger';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    loggers.error.error({
      error: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
      timestamp: new Date().toISOString(),
    }, 'React Error Boundary caught error');
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-6">
            <div className="flex items-center justify-center w-12 h-12 mx-auto bg-red-100 rounded-full">
              <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <h1 className="mt-4 text-center text-xl font-semibold text-gray-900">Something went wrong</h1>
            <p className="mt-2 text-center text-gray-600">We've been notified and are working on a fix.</p>
            <button
              onClick={() => window.location.reload()}
              className="mt-6 w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors"
            >
              Reload Page
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
```

## 🎨 Phase 6: Frontend Components

### 6.1 Nutrition Analysis Component

```typescript
// frontend-nextjs/src/components/features/NutritionAnalysis/NutritionForm.tsx
'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { nutritionService, NutritionAnalysisRequest } from '@/lib/api/services/nutrition.service';
import { loggers } from '@/lib/logger';

const formSchema = z.object({
  food: z.string().min(1, 'Food name is required'),
  quantity: z.number().min(0.1, 'Quantity must be at least 0.1'),
  unit: z.enum(['g', 'kg', 'oz', 'lb']),
  checkHalal: z.boolean().default(false),
});

type FormData = z.infer<typeof formSchema>;

export function NutritionForm() {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<NutritionAnalysisResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      food: '',
      quantity: 100,
      unit: 'g',
      checkHalal: true,
    },
  });

  const onSubmit = async (data: FormData) => {
    setIsLoading(true);
    setError(null);
    setResult(null);

    try {
      loggers.nutrition.info('Submitting nutrition analysis', { food: data.food });
      
      const analysisResult = await nutritionService.analyzeNutrition(data);
      setResult(analysisResult);
      
      loggers.nutrition.info('Nutrition analysis completed successfully', {
        requestId: analysisResult.requestId,
        calories: analysisResult.calories
      });
      
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to analyze nutrition';
      setError(errorMessage);
      
      loggers.nutrition.error('Nutrition analysis failed', {
        food: data.food,
        error: errorMessage
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Nutrition Analysis</h2>
      
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <label htmlFor="food" className="block text-sm font-medium text-gray-700">
            Food Item
          </label>
          <input
            {...register('food')}
            type="text"
            id="food"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            placeholder="e.g., apple, chicken breast"
          />
          {errors.food && (
            <p className="mt-1 text-sm text-red-600">{errors.food.message}</p>
          )}
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="quantity" className="block text-sm font-medium text-gray-700">
              Quantity
            </label>
            <input
              {...register('quantity', { valueAsNumber: true })}
              type="number"
              id="quantity"
              step="0.1"
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            />
            {errors.quantity && (
              <p className="mt-1 text-sm text-red-600">{errors.quantity.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="unit" className="block text-sm font-medium text-gray-700">
              Unit
            </label>
            <select
              {...register('unit')}
              id="unit"
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 p-2"
            >
              <option value="g">Grams (g)</option>
              <option value="kg">Kilograms (kg)</option>
              <option value="oz">Ounces (oz)</option>
              <option value="lb">Pounds (lb)</option>
            </select>
            {errors.unit && (
              <p className="mt-1 text-sm text-red-600">{errors.unit.message}</p>
            )}
          </div>
        </div>

        <div className="flex items-center">
          <input
            {...register('checkHalal')}
            type="checkbox"
            id="checkHalal"
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label htmlFor="checkHalal" className="ml-2 block text-sm text-gray-900">
            Check Halal status
          </label>
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? 'Analyzing...' : 'Analyze Nutrition'}
        </button>
      </form>

      {error && (
        <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      {result && (
        <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-md">
          <h3 className="text-lg font-semibold text-green-900 mb-4">Analysis Results</h3>
          <div className="space-y-2 text-sm">
            <p><strong>Food:</strong> {result.food} ({result.quantity}{result.unit})</p>
            <p><strong>Calories:</strong> {result.calories.toFixed(1)} kcal</p>
            <p><strong>Protein:</strong> {result.protein.toFixed(1)}g</p>
            <p><strong>Carbs:</strong> {result.carbs.toFixed(1)}g</p>
            <p><strong>Fat:</strong> {result.fat.toFixed(1)}g</p>
            <p><strong>Fiber:</strong> {result.fiber.toFixed(1)}g</p>
            <p><strong>Halal Status:</strong> {result.isHalal ? '✅ Halal' : '❌ Not Halal/Unknown'}</p>
            <p className="text-xs text-gray-500 mt-2">
              Processing time: {result.processingTime}ms (Request ID: {result.requestId})
            </p>
          </div>
          <button
            onClick={() => {
              setResult(null);
              reset();
            }}
            className="mt-4 w-full bg-gray-600 text-white py-2 px-4 rounded-md hover:bg-gray-700 transition-colors"
          >
            Analyze Another Item
          </button>
        </div>
      )}
    </div>
  );
}
```

## 🐳 Phase 7: Docker Configuration

### 7.1 Updated Docker Compose for Next.js + Node.js

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

### 7.2 Nginx Configuration

```nginx
# nutrition-platform/nginx.nextjs.conf
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging format
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=general:10m rate=30r/s;

    upstream backend {
        server backend:8080;
    }

    upstream frontend {
        server frontend:3000;
    }

    # HTTP server (redirect to HTTPS)
    server {
        listen 80;
        server_name _;
        return 301 https://$server_name$request_uri;
    }

    # HTTPS server
    server {
        listen 443 ssl http2;
        server_name localhost;

        # SSL configuration
        ssl_certificate /etc/ssl/cert.pem;
        ssl_certificate_key /etc/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # Security headers
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";

        # API routes
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
            
            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Frontend routes
        location / {
            limit_req zone=general burst=50 nodelay;
            
            proxy_pass http://frontend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
            
            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Health check
        location /health {
            proxy_pass http://backend/health;
            access_log off;
        }

        # Static assets caching
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
            proxy_pass http://frontend;
        }
    }
}
```

## 🧪 Phase 8: Testing Strategy

### 8.1 Backend Testing

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

### 8.2 Frontend Testing

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

  it('shows error message when API fails', async () => {
    const user = userEvent.setup();
    const errorMessage = 'Failed to analyze nutrition';

    mockedNutritionService.analyzeNutrition.mockRejectedValue(new Error(errorMessage));

    render(<NutritionForm />);

    await user.type(screen.getByLabelText('Food Item'), 'apple');
    await user.type(screen.getByLabelText('Quantity'), '100');
    await user.click(screen.getByRole('button', { name: 'Analyze Nutrition' }));

    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });
});
```

## 🚀 Phase 9: Implementation Steps

### 9.1 Step-by-Step Implementation Guide

1. **Fix Critical Bugs (Day 1)**
   ```bash
   cd nutrition-platform
   chmod +x fix-critical-bugs.sh
   ./fix-critical-bugs.sh
   ```

2. **Setup Next.js Project (Day 1-2)**
   ```bash
   cd nutrition-platform
   npx create-next-app@latest frontend-nextjs --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
   cd frontend-nextjs
   npm install axios zod react-hook-form @hookform/resolvers pino
   ```

3. **Implement Backend Enhancements (Day 2-3)**
   - Enhanced logging
   - API service improvements
   - Error handling

4. **Build Frontend Components (Day 3-5)**
   - API client setup
   - Nutrition analysis form
   - Dashboard components
   - Authentication flow

5. **Setup Docker & Deployment (Day 5-6)**
   - Configure Docker Compose
   - Setup Nginx reverse proxy
   - Test local deployment

6. **Testing & Validation (Day 6-7)**
   - Unit tests
   - Integration tests
   - E2E testing

7. **Production Deployment (Day 7+)**
   - Environment configuration
   - SSL setup
   - Monitoring configuration

## 📋 Phase 10: Success Criteria

### 10.1 Technical Requirements

- [ ] All critical bugs fixed
- [ ] Next.js frontend functional
- [ ] API integration working
- [ ] Proper error handling
- [ ] Comprehensive logging
- [ ] Authentication flow
- [ ] Docker deployment ready
- [ ] Nginx reverse proxy configured
- [ ] Tests passing
- [ ] Production-ready configuration

### 10.2 Performance Requirements

- [ ] Page load time < 2 seconds
- [ ] API response time < 500ms
- [ ] 99.9% uptime
- [ ] Proper caching implemented
- [ ] Bundle size optimized

### 10.3 Security Requirements

- [ ] HTTPS enforced
- [ ] Security headers configured
- [ ] Input validation on all endpoints
- [ ] Rate limiting implemented
- [ ] Sensitive data redacted from logs
- [ ] CORS properly configured

## 🎯 Conclusion

This comprehensive plan provides a roadmap to transform your nutrition platform into a modern, robust application using Next.js and Node.js. The implementation follows best practices from the reference documents and addresses all identified issues.

The plan emphasizes:
- **Immediate bug fixes** to stabilize the existing system
- **Modern frontend architecture** with Next.js App Router
- **Robust API integration** with proper error handling
- **Comprehensive logging** for debugging and monitoring
- **Production-ready deployment** with Docker and Nginx
- **Thorough testing** to ensure reliability

By following this plan, you'll have a nutrition platform that's not only bug-free but also scalable, maintainable, and ready for production deployment.