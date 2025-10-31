# 🚀 Running the Implementation

This guide provides the exact commands to run to implement your nutrition platform with Next.js frontend and Node.js backend.

## 📋 Prerequisites

Make sure you have:
- Node.js 18+ installed
- Docker and Docker Compose installed
- Git installed
- Command line/terminal access

## 🐛 Step 1: Fix Critical Bugs

First, let's fix the critical bugs in the existing Node.js backend:

```bash
# Navigate to nutrition platform directory
cd nutrition-platform

# Make the script executable
chmod +x fix-critical-bugs.sh

# Run the bug fix script
./fix-critical-bugs.sh
```

### What This Does:
- ✅ Adds missing prom-client import to server.js
- ✅ Creates logs directory on startup
- ✅ Fixes Redis client connection issue
- ✅ Fixes frontend API endpoint mismatch

### Verify Bug Fixes:
```bash
# Check if logs directory was created
ls -la production-nodejs/logs/

# Check if prom-client import was added
grep -n "prom-client" production-nodejs/server.js

# Check if Redis connection was fixed
grep -n "await this.client.connect" production-nodejs/services/redisClient.js
# Should return no results if fixed correctly
```

## 📦 Step 2: Setup Next.js Frontend

Now, let's create the Next.js frontend:

```bash
# Make the setup script executable
chmod +x setup-nextjs-frontend.sh

# Run the setup script
./setup-nextjs-frontend.sh
```

### What This Does:
- ✅ Creates Next.js project with TypeScript and Tailwind CSS
- ✅ Installs all required dependencies
- ✅ Sets up the directory structure
- ✅ Creates environment configuration
- ✅ Sets up API client and logger

### Verify Frontend Setup:
```bash
# Navigate to frontend directory
cd frontend-nextjs

# Check if dependencies are installed
ls node_modules/axios
ls node_modules/zod
ls node_modules/react-hook-form

# Check if directories were created
ls -la src/components/
ls -la src/lib/
ls -la src/types/

# Check if configuration files were created
ls -la .env.local
ls -la next.config.js
```

## 🔧 Step 3: Install Dependencies

If the setup script didn't complete successfully, install dependencies manually:

```bash
cd frontend-nextjs

# Install additional dependencies
npm install axios zod react-hook-form @hookform/resolvers pino next-auth

# Install dev dependencies
npm install -D @types/node

# Verify installation
npm list axios zod react-hook-form
```

## 🏃 Step 4: Run Development Servers

### Start Backend Server:
```bash
# Navigate to backend directory
cd nutrition-platform/production-nodejs

# Start the backend server
npm start
```

### Start Frontend Server:
```bash
# Open a new terminal window
cd nutrition-platform/frontend-nextjs

# Start the frontend development server
npm run dev
```

### Verify Servers Are Running:
```bash
# Check backend health
curl http://localhost:8080/health

# Check frontend is accessible
curl -I http://localhost:3000

# Or open in browser
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

## 🧪 Step 5: Test Functionality

### Test Main Page:
1. Open `http://localhost:3000` in your browser
2. Verify all 4 feature boxes are displayed:
   - Meals and Body Enhancing
   - Workouts and Injuries
   - Recipes and Review
   - Diseases and Healthy-Lifestyle
3. Click on each box to verify navigation works

### Test Meals Page:
1. Navigate to `/meals`
2. Fill out the user profile form with sample data:
   - Name: John Doe
   - Age: 30
   - Weight: 70kg
   - Height: 170cm
   - Activity Level: Moderate
   - Goal: Maintain Weight
   - Metabolic Rate: Medium
3. Click "Generate Meal Plan"
4. Verify BMI calculation is accurate
5. Verify nutrition calculations are correct
6. Click on a meal to view details

### Test Workouts Page:
1. Navigate to `/workouts`
2. Fill out the user profile form with sample data
3. Select a workout goal and location
4. Click "Generate Workout Plan"
5. Verify exercise recommendations are displayed
6. Check alternative exercises are available

### Test Recipes Page:
1. Navigate to `/recipes`
2. Select a cuisine (e.g., Italian)
3. Set diet type (e.g., Balanced)
4. Click "Generate Recipes"
5. Verify recipes are displayed with nutrition information
6. Check halal badges are present

### Test Health Page:
1. Navigate to `/health`
2. Select a disease (e.g., Diabetes)
3. Add medications if needed
4. Click "Get Health Advice"
5. Verify disease information is displayed
6. Check medical disclaimers are present

## 🔗 Step 6: Connect Frontend to Backend

### Test API Connection:
```bash
# Test API proxy from frontend
curl http://localhost:3000/api/health

# Test direct API call
curl http://localhost:8080/api/health
```

### Verify API Integration:
1. Check browser console for any API errors
2. Verify API requests are logged correctly
3. Test error handling by triggering an API error

## 🐳 Step 7: Docker Deployment (Optional)

For production deployment using Docker:

```bash
# Navigate to root directory
cd nutrition-platform

# Build and start services with Docker Compose
docker-compose -f docker-compose.nextjs.yml down
docker-compose -f docker-compose.nextjs.yml build
docker-compose -f docker-compose.nextjs.yml up -d
```

### Verify Docker Deployment:
```bash
# Check if containers are running
docker-compose -f docker-compose.nextjs.yml ps

# Check logs
docker-compose -f docker-compose.nextjs.yml logs -f backend
docker-compose -f docker-compose.nextjs.yml logs -f frontend

# Check if application is accessible
curl -I http://localhost
```

## 🔍 Step 8: Troubleshooting Common Issues

### Backend Issues:
```bash
# Check if backend is running
ps aux | grep node

# Check if port 8080 is available
lsof -i :8080

# Check backend logs
cd production-nodejs
npm start

# Check Redis connection
redis-cli ping
```

### Frontend Issues:
```bash
# Check if frontend is running
ps aux | grep next

# Check if port 3000 is available
lsof -i :3000

# Check frontend logs
cd frontend-nextjs
npm run dev
```

### Integration Issues:
```bash
# Check environment variables
cat frontend-nextjs/.env.local

# Check API configuration
cat frontend-nextjs/next.config.js

# Test API endpoint directly
curl -v http://localhost:8080/api/health
```

## 📊 Step 9: Monitor Application

### Check Application Health:
```bash
# Backend health
curl http://localhost:8080/health

# Frontend health
curl -I http://localhost:3000

# Database health (if using)
curl http://localhost:5432/health
```

### Monitor Logs:
```bash
# Backend logs
cd production-nodejs
tail -f logs/combined.log

# Frontend logs
cd frontend-nextjs
npm run dev
```

## 📋 Step 10: Complete Implementation Checklist

Use the CHECKLIST.md to verify all aspects of implementation:

```bash
# Open the checklist
cat nutrition-platform/CHECKLIST.md
```

### Check Key Requirements:
- ✅ All 4 main feature boxes are displayed
- ✅ Color scheme is correct (white background, pale yellow, green/blue accents)
- ✅ Nutrition calculations work accurately
- ✅ Workout plans generate correctly
- ✅ Recipes filter by cuisine correctly
- ✅ Disease information displays properly
- ✅ Halal food filtering works correctly
- ✅ Medical disclaimers are displayed
- ✅ Frontend connects to backend API successfully

## 🎯 Step 11: Production Deployment

For production deployment:

### 1. Update Environment Variables:
```bash
cd frontend-nextjs
cp .env.example .env.production
# Edit .env.production with production values
```

### 2. Build for Production:
```bash
cd frontend-nextjs
npm run build
```

### 3. Deploy with Docker:
```bash
cd nutrition-platform
docker-compose -f docker-compose.nextjs.yml down
docker-compose -f docker-compose.nextjs.yml build --no-cache
docker-compose -f docker-compose.nextjs.yml up -d
```

### 4. Verify Production Deployment:
```bash
# Check if application is running
docker-compose -f docker-compose.nextjs.yml ps

# Check if application is accessible
curl -I https://yourdomain.com
```

## 📞 Step 12: Support

If you encounter any issues:

### Check Logs:
- Browser console for frontend errors
- Terminal logs for backend errors
- Docker logs for container issues

### Common Solutions:
1. **Port Conflicts**: Kill processes using ports 3000 and 8080
2. **Dependencies**: Run `npm install` to install missing packages
3. **Permissions**: Run `chmod +x` on scripts if permission denied
4. **Environment Variables**: Check all environment variables are set correctly

### Get Help:
1. Review the implementation guide: `cat nutrition-platform/IMPLEMENTATION_GUIDE.md`
2. Check the checklist: `cat nutrition-platform/CHECKLIST.md`
3. Review the code comments for additional information

---

## 🎉 Success Criteria

Your implementation is successful when:

- ✅ All scripts run without errors
- ✅ Both frontend and backend servers start successfully
- ✅ All 4 main pages are accessible and functional
- ✅ Navigation between pages works correctly
- ✅ Forms validate input properly
- ✅ Calculations are accurate
- ✅ API integration works correctly
- ✅ Docker deployment works (if attempted)

---

## 📚 Additional Resources

### Documentation:
- `nutrition-platform/IMPLEMENTATION_GUIDE.md` - Detailed implementation guide
- `nutrition-platform/CHECKLIST.md` - Implementation checklist
- `nutrition-platform/COMPREHENSIVE_IMPLEMENTATION_PLAN.md` - Complete plan document
- `nutrition-platform/IMPLEMENTATION_SCRIPTS.md` - Code examples and scripts

### Scripts:
- `nutrition-platform/fix-critical-bugs.sh` - Bug fix script
- `nutrition-platform/setup-nextjs-frontend.sh` - Frontend setup script

This run implementation guide provides all the exact commands needed to successfully implement your nutrition platform with Next.js frontend and Node.js backend integration.