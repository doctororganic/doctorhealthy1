# 🚀 Quick Start Guide - Two Track Development

## Overview

Two independent tracks that can be developed **in parallel**:

- **Track 1:** Frontend User Experience (Search + Calculator)
- **Track 2:** Backend Performance & Security (Caching + Rate Limiting)

---

## 📋 Track 1: Frontend (4-5 hours)

### Files to Create:
1. `frontend-nextjs/src/hooks/useSearch.ts` - Search hook with debouncing
2. `frontend-nextjs/src/components/search/AdvancedSearch.tsx` - Search UI component
3. `frontend-nextjs/src/utils/nutritionCalculations.ts` - BMR/TDEE calculations
4. `frontend-nextjs/src/components/nutrition/NutritionCalculator.tsx` - Calculator UI
5. `frontend-nextjs/src/app/(dashboard)/calculator/page.tsx` - Calculator page

### Quick Commands:
```bash
cd frontend-nextjs

# Create files (copy from NEXT_PHASE_TWO_TRACK_PLAN.md)
# Test search
npm run dev
# Navigate to /recipes and test search

# Test calculator
# Navigate to /calculator and test calculations
```

### Key Features:
- ✅ Real-time search with debouncing
- ✅ Filter by cuisine, diet type, halal
- ✅ BMR/TDEE calculations
- ✅ Macro target calculations
- ✅ BMI calculator

---

## 🔒 Track 2: Backend (3-4 hours)

### Files to Create/Modify:
1. `backend/middleware/cache_middleware.go` - Response caching
2. `backend/middleware/enhanced_rate_limiter.go` - User-based rate limiting
3. `backend/middleware/security_headers.go` - Security headers (if not exists)
4. `backend/main.go` - Integrate middleware

### Quick Commands:
```bash
cd backend

# Set Redis (optional - cache works without it)
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=

# Build and test
go build ./...
go run main.go

# Test caching
curl -v http://localhost:8080/api/v1/nutrition-data/recipes?limit=5
# Check X-Cache header

# Test rate limiting
for i in {1..110}; do curl http://localhost:8080/api/v1/nutrition-data/recipes; done
# Should get 429 after limit
```

### Key Features:
- ✅ Response caching (Redis or in-memory)
- ✅ User-based rate limiting
- ✅ Security headers (CSP, XSS protection, etc.)
- ✅ Configurable TTL per endpoint

---

## 🎯 Testing Checklist

### Track 1:
- [ ] Search works on recipes page
- [ ] Filters apply correctly
- [ ] Calculator shows accurate BMR/TDEE
- [ ] Macro targets are correct
- [ ] Components are responsive

### Track 2:
- [ ] Cache middleware enabled
- [ ] X-Cache header present
- [ ] Rate limiting works
- [ ] Security headers present
- [ ] No errors in logs

---

## 📝 Next Steps After Completion

1. **Deploy Track 1:**
   ```bash
   cd frontend-nextjs
   npm run build
   # Deploy to production
   ```

2. **Deploy Track 2:**
   ```bash
   cd backend
   go build
   # Set environment variables
   # Deploy to production
   ```

3. **Monitor:**
   - Check cache hit rates
   - Monitor rate limit violations
   - Verify security headers with security scanner

---

## 🔗 Full Details

See `NEXT_PHASE_TWO_TRACK_PLAN.md` for:
- Complete code examples
- Detailed explanations
- Testing procedures
- Troubleshooting tips

---

**Total Time:** 7-9 hours (can be done in parallel)
**Priority:** HIGH
**Status:** Ready to Start ✅

