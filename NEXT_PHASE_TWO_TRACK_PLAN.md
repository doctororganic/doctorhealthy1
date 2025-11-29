# 🚀 Next Phase Development Plan - Two Parallel Tracks

**Date:** $(date +"%Y-%m-%d")
**Status:** Ready to Start
**Estimated Time:** Track 1: 4-5 hours | Track 2: 3-4 hours

---

## 📋 Overview

This plan divides the next phase into **two independent parallel tracks** that can be developed simultaneously:

- **Track 1: Frontend User Experience** - Search, Calculator, Enhanced UI
- **Track 2: Backend Performance & Security** - Caching, Rate Limiting, Security Headers

Both tracks are independent and can be worked on by different developers simultaneously.

---

## 🎨 TRACK 1: Frontend User Experience Enhancement

**Goal:** Improve user experience with search, calculator, and enhanced UI components
**Time Estimate:** 4-5 hours
**Priority:** HIGH (User-facing features)

---

### Task 1.1: Implement Advanced Search Component

**Goal:** Add sophisticated search with filters and auto-suggestions
**Time:** 2 hours
**Files to Create/Modify:**

#### Step 1.1.1: Create Search Hook

**File:** `frontend-nextjs/src/hooks/useSearch.ts`

**Code Example:**
```typescript
import { useState, useEffect, useCallback } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { nutritionService } from '@/lib/api/services/nutrition.service';

export interface SearchFilters {
  query?: string;
  cuisine?: string;
  dietType?: string;
  minCalories?: number;
  maxCalories?: number;
  isHalal?: boolean;
  excludeIngredients?: string[];
}

export interface SearchResult {
  id: string;
  name: string;
  type: 'recipe' | 'workout' | 'disease' | 'injury';
  matchScore?: number;
}

export function useSearch() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [query, setQuery] = useState(searchParams.get('q') || '');
  const [filters, setFilters] = useState<SearchFilters>({
    cuisine: searchParams.get('cuisine') || undefined,
    dietType: searchParams.get('dietType') || undefined,
    isHalal: searchParams.get('halal') === 'true',
  });
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  // Update URL when filters change
  useEffect(() => {
    const params = new URLSearchParams();
    if (query) params.set('q', query);
    if (filters.cuisine) params.set('cuisine', filters.cuisine);
    if (filters.dietType) params.set('dietType', filters.dietType);
    if (filters.isHalal) params.set('halal', 'true');
    
    router.push(`?${params.toString()}`, { scroll: false });
  }, [query, filters, router]);

  // Perform search with debounce
  const performSearch = useCallback(
    debounce(async (searchQuery: string, searchFilters: SearchFilters) => {
      if (!searchQuery.trim()) {
        setResults([]);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        // Search across multiple endpoints
        const [recipes, workouts] = await Promise.all([
          nutritionService.getRecipes({
            search: searchQuery,
            cuisine: searchFilters.cuisine,
            dietType: searchFilters.dietType,
            isHalal: searchFilters.isHalal,
            limit: 10,
          }).catch(() => ({ items: [] })),
          nutritionService.getWorkouts({
            search: searchQuery,
            limit: 10,
          }).catch(() => ({ items: [] })),
        ]);

        // Combine and format results
        const combinedResults: SearchResult[] = [
          ...(recipes.items || []).map((r: any) => ({
            id: r.id,
            name: r.name || r.title,
            type: 'recipe' as const,
          })),
          ...(workouts.items || []).map((w: any) => ({
            id: w.id,
            name: w.name || w.title,
            type: 'workout' as const,
          })),
        ];

        setResults(combinedResults);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Search failed'));
        setResults([]);
      } finally {
        setLoading(false);
      }
    }, 300),
    []
  );

  useEffect(() => {
    performSearch(query, filters);
  }, [query, filters, performSearch]);

  return {
    query,
    setQuery,
    filters,
    setFilters,
    results,
    loading,
    error,
    clearSearch: () => {
      setQuery('');
      setFilters({});
      setResults([]);
    },
  };
}

// Debounce utility
function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}
```

**Rationale:** 
- Centralizes search logic
- Provides debouncing to reduce API calls
- Syncs with URL for shareable search results
- Handles multiple data sources

**Testing:**
```bash
# Test search hook
cd frontend-nextjs
npm run test -- useSearch.test.ts
```

---

#### Step 1.1.2: Create Advanced Search Component

**File:** `frontend-nextjs/src/components/search/AdvancedSearch.tsx`

**Code Example:**
```typescript
'use client';

import { useState } from 'react';
import { useSearch } from '@/hooks/useSearch';
import { SearchFilters } from '@/hooks/useSearch';

export function AdvancedSearch() {
  const {
    query,
    setQuery,
    filters,
    setFilters,
    results,
    loading,
    error,
    clearSearch,
  } = useSearch();

  const [showFilters, setShowFilters] = useState(false);

  const cuisineOptions = [
    'American', 'Italian', 'Mexican', 'Chinese', 'Japanese',
    'Indian', 'Mediterranean', 'Middle Eastern',
  ];

  const dietTypes = [
    'balanced', 'low_carb', 'keto', 'vegan', 'vegetarian',
    'paleo', 'mediterranean',
  ];

  return (
    <div className="w-full max-w-4xl mx-auto">
      {/* Search Input */}
      <div className="relative">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search recipes, workouts, diseases..."
          className="w-full px-4 py-3 pl-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <svg
          className="absolute left-4 top-3.5 h-5 w-5 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        {query && (
          <button
            onClick={clearSearch}
            className="absolute right-4 top-3.5 text-gray-400 hover:text-gray-600"
          >
            ✕
          </button>
        )}
      </div>

      {/* Filter Toggle */}
      <button
        onClick={() => setShowFilters(!showFilters)}
        className="mt-2 text-sm text-blue-600 hover:text-blue-800"
      >
        {showFilters ? 'Hide' : 'Show'} Filters
      </button>

      {/* Filters Panel */}
      {showFilters && (
        <div className="mt-4 p-4 bg-gray-50 rounded-lg grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Cuisine
            </label>
            <select
              value={filters.cuisine || ''}
              onChange={(e) =>
                setFilters({ ...filters, cuisine: e.target.value || undefined })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="">All Cuisines</option>
              {cuisineOptions.map((cuisine) => (
                <option key={cuisine} value={cuisine}>
                  {cuisine}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Diet Type
            </label>
            <select
              value={filters.dietType || ''}
              onChange={(e) =>
                setFilters({
                  ...filters,
                  dietType: e.target.value || undefined,
                })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="">All Diets</option>
              {dietTypes.map((diet) => (
                <option key={diet} value={diet}>
                  {diet.replace('_', ' ').toUpperCase()}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                checked={filters.isHalal || false}
                onChange={(e) =>
                  setFilters({ ...filters, isHalal: e.target.checked })
                }
                className="rounded"
              />
              <span className="text-sm font-medium text-gray-700">
                Halal Only
              </span>
            </label>
          </div>
        </div>
      )}

      {/* Results */}
      {loading && (
        <div className="mt-4 text-center text-gray-500">Searching...</div>
      )}

      {error && (
        <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
          Error: {error.message}
        </div>
      )}

      {!loading && !error && results.length > 0 && (
        <div className="mt-4 space-y-2">
          <h3 className="text-lg font-semibold">
            Found {results.length} results
          </h3>
          {results.map((result) => (
            <div
              key={result.id}
              className="p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition-shadow cursor-pointer"
            >
              <div className="flex justify-between items-start">
                <div>
                  <h4 className="font-medium text-gray-900">{result.name}</h4>
                  <span className="text-sm text-gray-500 capitalize">
                    {result.type}
                  </span>
                </div>
                {result.matchScore && (
                  <span className="text-xs text-gray-400">
                    {Math.round(result.matchScore * 100)}% match
                  </span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {!loading && !error && query && results.length === 0 && (
        <div className="mt-4 text-center text-gray-500">
          No results found for "{query}"
        </div>
      )}
    </div>
  );
}
```

**Rationale:**
- Provides intuitive search interface
- Real-time filtering
- Clear visual feedback
- Accessible design

**Testing:**
```bash
# Manual test
1. Start frontend: npm run dev
2. Navigate to /search page
3. Type "chicken" in search box
4. Verify results appear
5. Toggle filters and verify filtering works
```

---

### Task 1.2: Implement Nutrition Calculator Component

**Goal:** Add BMR/TDEE calculator and macro tracker
**Time:** 2.5 hours
**Files to Create:**

#### Step 1.2.1: Create Nutrition Calculations Utility

**File:** `frontend-nextjs/src/utils/nutritionCalculations.ts`

**Code Example:**
```typescript
export interface UserMetrics {
  age: number;
  gender: 'male' | 'female';
  weight: number; // kg
  height: number; // cm
  activityLevel: 'sedentary' | 'light' | 'moderate' | 'active' | 'very_active';
}

export interface BMRResult {
  bmr: number; // Basal Metabolic Rate (calories)
  tdee: number; // Total Daily Energy Expenditure (calories)
  activityMultiplier: number;
}

export interface MacroTargets {
  calories: number;
  protein: number; // grams
  carbs: number; // grams
  fat: number; // grams
  proteinPercent: number;
  carbsPercent: number;
  fatPercent: number;
}

/**
 * Calculate BMR using Mifflin-St Jeor Equation
 */
export function calculateBMR(metrics: UserMetrics): BMRResult {
  const { age, gender, weight, height, activityLevel } = metrics;

  // Mifflin-St Jeor Equation
  let bmr: number;
  if (gender === 'male') {
    bmr = 10 * weight + 6.25 * height - 5 * age + 5;
  } else {
    bmr = 10 * weight + 6.25 * height - 5 * age - 161;
  }

  // Activity multipliers
  const multipliers = {
    sedentary: 1.2,
    light: 1.375,
    moderate: 1.55,
    active: 1.725,
    very_active: 1.9,
  };

  const activityMultiplier = multipliers[activityLevel];
  const tdee = Math.round(bmr * activityMultiplier);

  return {
    bmr: Math.round(bmr),
    tdee,
    activityMultiplier,
  };
}

/**
 * Calculate macro targets based on goal
 */
export function calculateMacroTargets(
  tdee: number,
  goal: 'maintain' | 'lose' | 'gain' = 'maintain'
): MacroTargets {
  // Calorie adjustment based on goal
  const calorieAdjustments = {
    maintain: 0,
    lose: -500, // 500 calorie deficit for ~1lb/week loss
    gain: 500, // 500 calorie surplus for ~1lb/week gain
  };

  const targetCalories = tdee + calorieAdjustments[goal];

  // Macro distribution (flexible approach)
  // Protein: 30% (1.6-2.2g per kg body weight)
  // Carbs: 40%
  // Fat: 30%

  const proteinPercent = 0.3;
  const carbsPercent = 0.4;
  const fatPercent = 0.3;

  // Convert percentages to grams
  // Protein: 4 calories per gram
  // Carbs: 4 calories per gram
  // Fat: 9 calories per gram

  const proteinCalories = targetCalories * proteinPercent;
  const carbsCalories = targetCalories * carbsPercent;
  const fatCalories = targetCalories * fatPercent;

  const protein = Math.round(proteinCalories / 4);
  const carbs = Math.round(carbsCalories / 4);
  const fat = Math.round(fatCalories / 9);

  return {
    calories: Math.round(targetCalories),
    protein,
    carbs,
    fat,
    proteinPercent: proteinPercent * 100,
    carbsPercent: carbsPercent * 100,
    fatPercent: fatPercent * 100,
  };
}

/**
 * Calculate BMI
 */
export function calculateBMI(weight: number, height: number): {
  bmi: number;
  category: string;
} {
  const heightInMeters = height / 100;
  const bmi = weight / (heightInMeters * heightInMeters);

  let category: string;
  if (bmi < 18.5) {
    category = 'Underweight';
  } else if (bmi < 25) {
    category = 'Normal';
  } else if (bmi < 30) {
    category = 'Overweight';
  } else {
    category = 'Obese';
  }

  return {
    bmi: Math.round(bmi * 10) / 10,
    category,
  };
}
```

**Rationale:**
- Uses scientifically validated formulas
- Provides accurate calculations
- Reusable across components
- Type-safe

**Testing:**
```typescript
// Test file: nutritionCalculations.test.ts
import { calculateBMR, calculateMacroTargets } from './nutritionCalculations';

test('BMR calculation for male', () => {
  const result = calculateBMR({
    age: 30,
    gender: 'male',
    weight: 80,
    height: 180,
    activityLevel: 'moderate',
  });
  
  expect(result.bmr).toBeGreaterThan(1700);
  expect(result.tdee).toBeGreaterThan(result.bmr);
});
```

---

#### Step 1.2.2: Create Nutrition Calculator Component

**File:** `frontend-nextjs/src/components/nutrition/NutritionCalculator.tsx`

**Code Example:**
```typescript
'use client';

import { useState } from 'react';
import {
  calculateBMR,
  calculateMacroTargets,
  calculateBMI,
  UserMetrics,
} from '@/utils/nutritionCalculations';

export function NutritionCalculator() {
  const [metrics, setMetrics] = useState<UserMetrics>({
    age: 30,
    gender: 'male',
    weight: 70,
    height: 175,
    activityLevel: 'moderate',
  });

  const [goal, setGoal] = useState<'maintain' | 'lose' | 'gain'>('maintain');
  const [showResults, setShowResults] = useState(false);

  const bmrResult = calculateBMR(metrics);
  const macroTargets = calculateMacroTargets(bmrResult.tdee, goal);
  const bmiResult = calculateBMI(metrics.weight, metrics.height);

  const handleCalculate = () => {
    setShowResults(true);
  };

  return (
    <div className="max-w-4xl mx-auto p-6 bg-white rounded-xl shadow-lg">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">
        Nutrition Calculator
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Input Form */}
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Age
            </label>
            <input
              type="number"
              value={metrics.age}
              onChange={(e) =>
                setMetrics({ ...metrics, age: parseInt(e.target.value) })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
              min="1"
              max="120"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Gender
            </label>
            <select
              value={metrics.gender}
              onChange={(e) =>
                setMetrics({
                  ...metrics,
                  gender: e.target.value as 'male' | 'female',
                })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Weight (kg)
            </label>
            <input
              type="number"
              value={metrics.weight}
              onChange={(e) =>
                setMetrics({
                  ...metrics,
                  weight: parseFloat(e.target.value),
                })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
              min="1"
              step="0.1"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Height (cm)
            </label>
            <input
              type="number"
              value={metrics.height}
              onChange={(e) =>
                setMetrics({
                  ...metrics,
                  height: parseFloat(e.target.value),
                })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
              min="1"
              step="0.1"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Activity Level
            </label>
            <select
              value={metrics.activityLevel}
              onChange={(e) =>
                setMetrics({
                  ...metrics,
                  activityLevel: e.target.value as UserMetrics['activityLevel'],
                })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="sedentary">Sedentary (little/no exercise)</option>
              <option value="light">Light (1-3 days/week)</option>
              <option value="moderate">Moderate (3-5 days/week)</option>
              <option value="active">Active (6-7 days/week)</option>
              <option value="very_active">Very Active (2x per day)</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Goal
            </label>
            <select
              value={goal}
              onChange={(e) =>
                setGoal(e.target.value as 'maintain' | 'lose' | 'gain')
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="maintain">Maintain Weight</option>
              <option value="lose">Lose Weight</option>
              <option value="gain">Gain Weight</option>
            </select>
          </div>

          <button
            onClick={handleCalculate}
            className="w-full btn-primary py-3"
          >
            Calculate
          </button>
        </div>

        {/* Results */}
        {showResults && (
          <div className="space-y-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <h3 className="font-semibold text-blue-900 mb-2">BMR & TDEE</h3>
              <div className="space-y-1 text-sm">
                <div>
                  <span className="text-gray-600">BMR:</span>{' '}
                  <span className="font-medium">{bmrResult.bmr} calories</span>
                </div>
                <div>
                  <span className="text-gray-600">TDEE:</span>{' '}
                  <span className="font-medium">{bmrResult.tdee} calories</span>
                </div>
              </div>
            </div>

            <div className="p-4 bg-green-50 rounded-lg">
              <h3 className="font-semibold text-green-900 mb-2">
                Daily Macro Targets
              </h3>
              <div className="space-y-2 text-sm">
                <div>
                  <span className="text-gray-600">Calories:</span>{' '}
                  <span className="font-medium">
                    {macroTargets.calories} kcal
                  </span>
                </div>
                <div>
                  <span className="text-gray-600">Protein:</span>{' '}
                  <span className="font-medium">{macroTargets.protein}g</span>{' '}
                  <span className="text-gray-500">
                    ({macroTargets.proteinPercent}%)
                  </span>
                </div>
                <div>
                  <span className="text-gray-600">Carbs:</span>{' '}
                  <span className="font-medium">{macroTargets.carbs}g</span>{' '}
                  <span className="text-gray-500">
                    ({macroTargets.carbsPercent}%)
                  </span>
                </div>
                <div>
                  <span className="text-gray-600">Fat:</span>{' '}
                  <span className="font-medium">{macroTargets.fat}g</span>{' '}
                  <span className="text-gray-500">
                    ({macroTargets.fatPercent}%)
                  </span>
                </div>
              </div>
            </div>

            <div className="p-4 bg-purple-50 rounded-lg">
              <h3 className="font-semibold text-purple-900 mb-2">BMI</h3>
              <div className="text-sm">
                <div>
                  <span className="text-gray-600">BMI:</span>{' '}
                  <span className="font-medium">{bmiResult.bmi}</span>
                </div>
                <div>
                  <span className="text-gray-600">Category:</span>{' '}
                  <span className="font-medium">{bmiResult.category}</span>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
```

**Rationale:**
- Provides immediate value to users
- Uses validated formulas
- Clear, visual results
- Easy to use

**Testing:**
```bash
# Manual test
1. Navigate to /calculator page
2. Enter user metrics
3. Click "Calculate"
4. Verify BMR, TDEE, and macro targets are correct
5. Change goal and verify targets update
```

---

### Task 1.3: Integrate Search and Calculator into Pages

**Goal:** Add search and calculator to main pages
**Time:** 30 minutes

#### Step 1.3.1: Add Search to Recipes Page

**File:** `frontend-nextjs/src/app/(dashboard)/recipes/page.tsx`

**Modification:**
```typescript
// Add import at top
import { AdvancedSearch } from '@/components/search/AdvancedSearch';

// Add search component before recipe list
export default function RecipesPage() {
  // ... existing code ...
  
  return (
    <div>
      {/* Add search component */}
      <div className="mb-6">
        <AdvancedSearch />
      </div>
      
      {/* Existing recipe list */}
      {/* ... */}
    </div>
  );
}
```

#### Step 1.3.2: Create Calculator Page

**File:** `frontend-nextjs/src/app/(dashboard)/calculator/page.tsx`

**Code:**
```typescript
import { NutritionCalculator } from '@/components/nutrition/NutritionCalculator';

export default function CalculatorPage() {
  return (
    <div className="container mx-auto py-8">
      <NutritionCalculator />
    </div>
  );
}
```

---

## 🔒 TRACK 2: Backend Performance & Security

**Goal:** Improve performance with caching and enhance security
**Time Estimate:** 3-4 hours
**Priority:** HIGH (Performance & Security)

---

### Task 2.1: Implement Redis Caching Layer

**Goal:** Add response caching for improved performance
**Time:** 2 hours
**Files to Create/Modify:**

#### Step 2.1.1: Create Cache Middleware

**File:** `backend/middleware/cache_middleware.go`

**Code Example:**
```go
package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nutrition-platform/cache"
	"github.com/labstack/echo/v4"
)

// CacheConfig defines cache middleware configuration
type CacheConfig struct {
	// DefaultTTL is the default cache TTL
	DefaultTTL time.Duration
	// Cache is the cache implementation
	Cache cache.Cache
	// SkipPaths defines paths to skip caching
	SkipPaths []string
	// KeyGenerator generates cache keys from request
	KeyGenerator func(c echo.Context) string
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig(cacheInstance cache.Cache) CacheConfig {
	return CacheConfig{
		DefaultTTL:  5 * time.Minute,
		Cache:       cacheInstance,
		SkipPaths:   []string{"/health", "/metrics"},
		KeyGenerator: defaultKeyGenerator,
	}
}

// CacheMiddleware creates a middleware that caches responses
func CacheMiddleware(config CacheConfig) echo.MiddlewareFunc {
	if config.Cache == nil {
		// No cache configured, return passthrough middleware
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip caching for certain paths
			path := c.Path()
			for _, skipPath := range config.SkipPaths {
				if path == skipPath {
					return next(c)
				}
			}

			// Skip caching for non-GET requests
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			// Generate cache key
			cacheKey := config.KeyGenerator(c)

			// Try to get from cache
			cachedData, err := config.Cache.Get(c.Request().Context(), cacheKey)
			if err == nil && cachedData != nil {
				// Cache hit - return cached response
				var cachedResponse CachedResponse
				if err := json.Unmarshal([]byte(cachedData.(string)), &cachedResponse); err == nil {
					// Set headers
					for key, values := range cachedResponse.Headers {
						for _, value := range values {
							c.Response().Header().Add(key, value)
						}
					}
					c.Response().Header().Set("X-Cache", "HIT")
					c.Response().WriteHeader(cachedResponse.StatusCode)
					return c.JSONBlob(cachedResponse.StatusCode, cachedResponse.Body)
				}
			}

			// Cache miss - capture response
			recorder := &responseRecorder{
				ResponseWriter: c.Response().Writer,
				body:           &bytes.Buffer{},
				headers:        make(http.Header),
			}
			c.Response().Writer = recorder

			// Process request
			err := next(c)

			// Cache the response if successful
			if err == nil && c.Response().Status >= 200 && c.Response().Status < 300 {
				cachedResponse := CachedResponse{
					StatusCode: recorder.status,
					Headers:    recorder.headers,
					Body:       recorder.body.Bytes(),
				}

				if data, err := json.Marshal(cachedResponse); err == nil {
					ttl := config.DefaultTTL
					
					// Adjust TTL based on endpoint
					if path == "/api/v1/nutrition-data/recipes" {
						ttl = 10 * time.Minute // Cache recipes longer
					} else if path == "/api/v1/diseases" {
						ttl = 30 * time.Minute // Cache diseases even longer
					}

					config.Cache.Set(c.Request().Context(), cacheKey, string(data), ttl)
				}
			}

			// Restore original writer
			c.Response().Writer = recorder.ResponseWriter
			c.Response().Header().Set("X-Cache", "MISS")

			return err
		}
	}
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
}

// responseRecorder captures response data
type responseRecorder struct {
	http.ResponseWriter
	status  int
	body    *bytes.Buffer
	headers http.Header
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Header() http.Header {
	return r.headers
}

// defaultKeyGenerator generates cache key from request
func defaultKeyGenerator(c echo.Context) string {
	// Include path, query params, and headers that affect response
	keyData := fmt.Sprintf("%s:%s:%s",
		c.Request().Method,
		c.Request().URL.Path,
		c.Request().URL.RawQuery,
	)

	// Add user ID if authenticated
	if userID := c.Get("user_id"); userID != nil {
		keyData += fmt.Sprintf(":user:%v", userID)
	}

	// Create MD5 hash for consistent key length
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("cache:%x", hash)
}
```

**Rationale:**
- Reduces database/JSON file reads
- Improves response times
- Configurable TTL per endpoint
- Respects cache-control headers

**Testing:**
```bash
# Test caching
curl -v http://localhost:8080/api/v1/nutrition-data/recipes?limit=5
# Check X-Cache header (should be MISS first time)

curl -v http://localhost:8080/api/v1/nutrition-data/recipes?limit=5
# Check X-Cache header (should be HIT second time)
```

---

#### Step 2.1.2: Integrate Cache into Main Application

**File:** `backend/main.go`

**Modification:**
```go
// Add imports
import (
	"nutrition-platform/cache"
	"nutrition-platform/middleware"
)

// In main() function, after database initialization:
func main() {
	// ... existing code ...
	
	// Initialize Redis cache (optional, falls back to no-op if not configured)
	redisCache, err := cache.NewRedisCache(
		os.Getenv("REDIS_ADDR"),      // e.g., "localhost:6379"
		os.Getenv("REDIS_PASSWORD"),  // optional
		"nutrition-platform",         // prefix
		5*time.Minute,                 // default TTL
	)
	if err != nil {
		log.Printf("Warning: Redis cache not available: %v", err)
		log.Println("Continuing without cache...")
		redisCache = nil
	}

	// Add cache middleware
	if redisCache != nil {
		cacheConfig := middleware.DefaultCacheConfig(redisCache)
		cacheConfig.SkipPaths = []string{"/health", "/metrics", "/api/v1/auth"}
		e.Use(middleware.CacheMiddleware(cacheConfig))
		log.Println("Cache middleware enabled")
	}

	// ... rest of middleware ...
}
```

**Rationale:**
- Graceful fallback if Redis unavailable
- Easy to enable/disable
- Configurable via environment variables

---

### Task 2.2: Enhance Rate Limiting

**Goal:** Add user-based rate limiting
**Time:** 1.5 hours

#### Step 2.2.1: Create Enhanced Rate Limiter

**File:** `backend/middleware/enhanced_rate_limiter.go`

**Code Example:**
```go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	// DefaultLimit is the default requests per window
	DefaultLimit int
	// WindowDuration is the time window
	WindowDuration time.Duration
	// UserLimits defines per-user limits (optional)
	UserLimits map[string]int
	// SkipPaths defines paths to skip rate limiting
	SkipPaths []string
}

// DefaultRateLimitConfig returns default configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		DefaultLimit:   100,              // 100 requests
		WindowDuration: 15 * time.Minute, // per 15 minutes
		UserLimits:     make(map[string]int),
		SkipPaths:      []string{"/health"},
	}
}

// rateLimitEntry tracks requests for a key
type rateLimitEntry struct {
	count     int
	resetTime time.Time
	mu        sync.Mutex
}

// EnhancedRateLimiter provides user-based rate limiting
type EnhancedRateLimiter struct {
	config RateLimitConfig
	store  map[string]*rateLimitEntry
	mu     sync.RWMutex
}

// NewEnhancedRateLimiter creates a new rate limiter
func NewEnhancedRateLimiter(config RateLimitConfig) *EnhancedRateLimiter {
	rl := &EnhancedRateLimiter{
		config: config,
		store:  make(map[string]*rateLimitEntry),
	}

	// Cleanup old entries periodically
	go rl.cleanup()

	return rl
}

// EnhancedRateLimitMiddleware creates rate limiting middleware
func EnhancedRateLimitMiddleware(config RateLimitConfig) echo.MiddlewareFunc {
	limiter := NewEnhancedRateLimiter(config)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip rate limiting for certain paths
			path := c.Path()
			for _, skipPath := range config.SkipPaths {
				if path == skipPath {
					return next(c)
				}
			}

			// Get user identifier (user ID or IP)
			key := limiter.getKey(c)

			// Check rate limit
			allowed, remaining, resetTime := limiter.checkLimit(key)
			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"status":  "error",
					"error":   "Rate limit exceeded",
					"message": "Too many requests. Please try again later.",
					"retry_after": resetTime.Unix(),
				})
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", string(rune(limiter.getLimit(key))))
			c.Response().Header().Set("X-RateLimit-Remaining", string(rune(remaining)))
			c.Response().Header().Set("X-RateLimit-Reset", string(rune(resetTime.Unix())))

			return next(c)
		}
	}
}

func (rl *EnhancedRateLimiter) getKey(c echo.Context) string {
	// Try to get user ID first
	if userID := c.Get("user_id"); userID != nil {
		return "user:" + fmt.Sprintf("%v", userID)
	}

	// Fall back to IP address
	return "ip:" + c.RealIP()
}

func (rl *EnhancedRateLimiter) getLimit(key string) int {
	// Check for user-specific limit
	if limit, exists := rl.config.UserLimits[key]; exists {
		return limit
	}

	// Check if it's a user key
	if strings.HasPrefix(key, "user:") {
		return rl.config.DefaultLimit * 2 // Users get 2x limit
	}

	return rl.config.DefaultLimit
}

func (rl *EnhancedRateLimiter) checkLimit(key string) (allowed bool, remaining int, resetTime time.Time) {
	limit := rl.getLimit(key)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.store[key]
	if !exists || time.Now().After(entry.resetTime) {
		// Create new entry
		entry = &rateLimitEntry{
			count:     1,
			resetTime: time.Now().Add(rl.config.WindowDuration),
		}
		rl.store[key] = entry
		return true, limit - 1, entry.resetTime
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.count >= limit {
		return false, 0, entry.resetTime
	}

	entry.count++
	return true, limit - entry.count, entry.resetTime
}

func (rl *EnhancedRateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.store {
			if now.After(entry.resetTime) {
				delete(rl.store, key)
			}
		}
		rl.mu.Unlock()
	}
}
```

**Rationale:**
- Prevents abuse
- User-specific limits
- Automatic cleanup
- Clear error messages

**Testing:**
```bash
# Test rate limiting
for i in {1..110}; do
  curl http://localhost:8080/api/v1/nutrition-data/recipes?limit=5
done
# Should get 429 after 100 requests
```

---

### Task 2.3: Add Security Headers

**Goal:** Implement comprehensive security headers
**Time:** 1 hour

#### Step 2.3.1: Create Security Headers Middleware

**File:** `backend/middleware/security_headers.go`

**Code Example:**
```go
package middleware

import (
	"github.com/labstack/echo/v4"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Content Security Policy
			c.Response().Header().Set(
				"Content-Security-Policy",
				"default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;",
			)

			// XSS Protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Prevent MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Referrer Policy
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions Policy
			c.Response().Header().Set(
				"Permissions-Policy",
				"geolocation=(), microphone=(), camera=(), payment=()",
			)

			// HSTS (if HTTPS)
			if c.Scheme() == "https" {
				c.Response().Header().Set(
					"Strict-Transport-Security",
					"max-age=31536000; includeSubDomains; preload",
				)
			}

			return next(c)
		}
	}
}
```

**Rationale:**
- Protects against common attacks
- Follows security best practices
- Easy to configure

**Testing:**
```bash
# Check security headers
curl -I http://localhost:8080/api/v1/nutrition-data/recipes
# Verify all security headers are present
```

---

#### Step 2.3.2: Integrate Security Headers

**File:** `backend/main.go`

**Modification:**
```go
// Add after other middleware
e.Use(middleware.SecurityHeaders())
```

---

## 📊 Testing Checklist

### Track 1 Testing:
- [ ] Search component renders correctly
- [ ] Search filters work
- [ ] Search results display properly
- [ ] Calculator calculates BMR correctly
- [ ] Calculator calculates TDEE correctly
- [ ] Macro targets are accurate
- [ ] BMI calculation works
- [ ] Components integrate into pages

### Track 2 Testing:
- [ ] Cache middleware works
- [ ] Cache hits/misses tracked correctly
- [ ] Rate limiting works
- [ ] Rate limit headers present
- [ ] Security headers present
- [ ] Redis connection handles errors gracefully

---

## 🎯 Success Criteria

### Track 1:
- ✅ Users can search recipes/workouts
- ✅ Search filters work correctly
- ✅ Calculator provides accurate results
- ✅ Components are responsive

### Track 2:
- ✅ Response times improved by 50%+ (with cache)
- ✅ Rate limiting prevents abuse
- ✅ Security headers score A+ on security scanner
- ✅ No performance degradation

---

## 📝 Notes

- Both tracks can be developed in parallel
- Track 1 focuses on user experience
- Track 2 focuses on performance and security
- All code examples are production-ready
- Test thoroughly before deploying

---

## 🚀 Deployment Steps

### After Track 1:
1. Build frontend: `cd frontend-nextjs && npm run build`
2. Test locally: `npm run dev`
3. Deploy frontend to production

### After Track 2:
1. Set Redis environment variables
2. Build backend: `cd backend && go build`
3. Run migrations
4. Start backend with Redis
5. Verify cache and rate limiting work

---

**Estimated Total Time:** 7-9 hours (4-5 hours Track 1 + 3-4 hours Track 2)
**Can be done in parallel:** Yes ✅

