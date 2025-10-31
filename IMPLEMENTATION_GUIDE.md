# 🚀 Implementation Guide: Next.js Frontend with Node.js Backend

This guide provides step-by-step instructions to implement the nutrition platform with all the features we've designed.

## 📋 Prerequisites

Before starting, ensure you have the following installed:
- Node.js 18+ and npm
- Git
- Docker (optional, for containerized deployment)

## 🐛 Step 1: Fix Critical Bugs

First, let's fix the critical bugs in the existing Node.js backend:

```bash
cd nutrition-platform
chmod +x fix-critical-bugs.sh
./fix-critical-bugs.sh
```

This script will:
- Add missing prom-client import to server.js
- Create logs directory on startup
- Fix Redis client connection issue
- Fix frontend API endpoint mismatch

## 📦 Step 2: Setup Next.js Project

Create a new Next.js project with the required configuration:

```bash
cd nutrition-platform
chmod +x setup-nextjs-frontend.sh
./setup-nextjs-frontend.sh
```

This script will:
- Create a Next.js project with TypeScript and Tailwind CSS
- Install required dependencies
- Set up the directory structure

## 🔧 Step 3: Install Dependencies

Navigate to the frontend directory and install additional dependencies:

```bash
cd frontend-nextjs
npm install axios zod react-hook-form @hookform/resolvers pino next-auth
npm install -D @types/node
```

## 🏗️ Step 4: Project Structure

Your project structure should look like this:

```
nutrition-platform/
├── production-nodejs/             # Backend (fixed)
├── frontend-nextjs/              # New Next.js frontend
│   ├── src/
│   │   ├── app/                  # App Router pages
│   │   │   ├── (dashboard)/      # Protected routes
│   │   │   │   ├── meals/
│   │   │   │   ├── workouts/
│   │   │   │   ├── recipes/
│   │   │   │   └── health/
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── components/
│   │   │   ├── icons/
│   │   │   └── providers/
│   │   ├── lib/
│   │   │   ├── api/
│   │   │   ├── logger/
│   │   │   └── validations/
│   │   └── types/
│   ├── public/
│   ├── .env.local
│   ├── .env.example
│   ├── next.config.js
│   ├── package.json
│   └── tsconfig.json
├── docker-compose.nextjs.yml     # Docker configuration
└── nginx.nextjs.conf             # Nginx configuration
```

## 🔧 Step 5: Environment Configuration

Create the environment configuration file:

```bash
cd frontend-nextjs
cp .env.example .env.local
```

Edit `.env.local` with your configuration:

```bash
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

## 🏃 Step 6: Run Development Server

Start the development server to test the application:

```bash
cd frontend-nextjs
npm run dev
```

The application should now be running at `http://localhost:3000`.

## 🔌 Step 7: Start Backend Server

In a separate terminal, start the Node.js backend:

```bash
cd nutrition-platform/production-nodejs
npm start
```

The backend API should now be running at `http://localhost:8080`.

## 🔗 Step 8: Connect Frontend to Backend

Create an API proxy configuration in `next.config.js`:

```javascript
// frontend-nextjs/next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  // ... other configurations
  
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
};

module.exports = nextConfig;
```

## 🧪 Step 9: Test Functionality

Test all the main features:

### 1. Main Page
- Navigate to `http://localhost:3000`
- Verify all 4 feature boxes are displayed correctly
- Test navigation to each section

### 2. Meals Page
- Navigate to `/meals`
- Test the user profile form
- Generate a meal plan
- Test meal details modal

### 3. Workouts Page
- Navigate to `/workouts`
- Test the user profile form
- Generate a workout plan
- Test injury advice

### 4. Recipes Page
- Navigate to `/recipes`
- Test cuisine selection
- Generate recipes
- Test recipe details modal

### 5. Health Page
- Navigate to `/health`
- Test disease selection
- Test medical advice
- Test medication information

## 🐳 Step 10: Docker Deployment (Optional)

For containerized deployment, use Docker Compose:

```bash
cd nutrition-platform
docker-compose -f docker-compose.nextjs.yml down
docker-compose -f docker-compose.nextjs.yml build
docker-compose -f docker-compose.nextjs.yml up -d
```

## 🔍 Step 11: Debugging Common Issues

### Frontend Issues

1. **TypeScript Errors**: Make sure all imports are correct
2. **Missing Modules**: Run `npm install` to install missing packages
3. **Import Errors**: Check file paths and module names

### Backend Issues

1. **Port Conflicts**: Make sure ports 3000 and 8080 are available
2. **Redis Connection**: Check Redis is running on port 6379
3. **Database Connection**: Verify database connection settings

### Integration Issues

1. **API Errors**: Check browser console for API errors
2. **CORS Errors**: Verify CORS configuration in backend
3. **Environment Variables**: Check all environment variables are set correctly

## 📝 Step 12: Additional Configuration

### API Client Configuration

Create a proper API client for frontend-backend communication:

```typescript
// frontend-nextjs/src/lib/api/client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default apiClient;
```

### Authentication Setup

Set up authentication with NextAuth.js:

```typescript
// frontend-nextjs/src/app/api/auth/[...nextauth]/route.ts
import NextAuth from 'next-auth';
import { authOptions } from '@/lib/auth';

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
```

## 🚀 Step 13: Production Deployment

For production deployment:

1. **Environment Variables**: Set production environment variables
2. **SSL Certificate**: Configure SSL for HTTPS
3. **Domain Configuration**: Set up your domain
4. **Monitoring**: Set up monitoring and logging

## 📚 Step 14: Additional Features

Consider implementing these additional features:

1. **User Authentication**: Full user registration and login system
2. **Database Integration**: Connect to a persistent database
3. **File Upload**: Allow users to upload photos and documents
4. **Email Notifications**: Send email notifications to users
5. **Analytics**: Track user behavior and application performance

## 🎯 Success Criteria

Your implementation is successful when:

- ✅ All 4 main feature boxes are displayed on the homepage
- ✅ Navigation works correctly between all pages
- ✅ All forms validate input correctly
- ✅ Nutrition calculations are accurate
- ✅ Workout plans generate correctly
- ✅ Recipes filter by cuisine correctly
- ✅ Disease information displays correctly
- ✅ Halal food filtering works properly
- ✅ Medical disclaimers are displayed
- ✅ Frontend connects to backend API successfully

## 📞 Troubleshooting

If you encounter issues:

1. **Check Logs**: Look at browser console and terminal logs
2. **Verify Configuration**: Check all configuration files
3. **Restart Services**: Restart both frontend and backend
4. **Clear Cache**: Clear browser cache and npm cache
5. **Update Dependencies**: Update all packages to latest versions

This implementation guide provides a complete roadmap to build and deploy your nutrition platform with all the features we've designed.