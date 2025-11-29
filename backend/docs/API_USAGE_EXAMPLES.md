# API Usage Examples

## Overview

This document provides comprehensive examples of how to use the DoctorHealthy nutrition platform API endpoints.

## Base URL

```
Development: http://localhost:8080
Production: https://api.doctorhealthy.com
```

## Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/v1/protected-endpoint
```

## Nutrition Data Endpoints

### Get Recipes with Pagination

**Request:**
```bash
curl "http://localhost:8080/api/v1/nutrition-data/recipes?page=1&limit=20"
```

**Response Format:**
```json
{
  "status": "success",
  "items": [
    {
      "id": 1,
      "name": "Grilled Chicken Salad",
      "calories": 350,
      "protein": 35,
      "carbs": 15,
      "fat": 18,
      "ingredients": ["chicken breast", "lettuce", "tomatoes", "olive oil"],
      "instructions": ["Grill chicken", "Chop vegetables", "Mix together"],
      "cuisine": "Mediterranean",
      "dietType": "high-protein",
      "prepTime": 25,
      "cookTime": 15,
      "servings": 2,
      "isHalal": true
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### Search Recipes by Filters

**Request:**
```bash
curl "http://localhost:8080/api/v1/nutrition-data/recipes?cuisine=Mediterranean&dietType=high-protein&maxCalories=400"
```

### Get Workouts

**Request:**
```bash
curl "http://localhost:8080/api/v1/nutrition-data/workouts?page=1&limit=10&goal=weight-loss&experience_level=beginner"
```

**Response Format:**
```json
{
  "status": "success",
  "items": [
    {
      "id": 1,
      "goal": "weight-loss",
      "experience_level": "beginner",
      "plan": {
        "week": [
          {
            "day": "Monday",
            "exercises": [
              {
                "name": "Push-ups",
                "sets": 3,
                "reps": 10,
                "rest": 60
              }
            ]
          }
        ]
      },
      "title": {
        "en": "Beginner Weight Loss Program",
        "ar": "برنامج إنقاص الوزن للمبتدئين"
      },
      "description": {
        "en": "A comprehensive 4-week program for beginners",
        "ar": "برنامج شامل لمدة 4 أسابيع للمبتدئين"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

## Frontend Integration Examples

### Using React Hook

```typescript
import { useRecipes } from '../hooks/useNutritionData';
import { LoadingSkeleton, ErrorDisplay, Pagination } from '../components/ui';

function RecipesPage() {
  const { 
    data, 
    loading, 
    error, 
    pagination, 
    goToPage 
  } = useRecipes({ 
    page: 1, 
    limit: 20,
    cuisine: 'Mediterranean',
    dietType: 'high-protein'
  });
  
  if (loading) return <LoadingSkeleton count={5} />;
  if (error) return <ErrorDisplay error={error} onRetry={() => window.location.reload()} />;
  
  return (
    <div className="container mx-auto px-4">
      <h1 className="text-2xl font-bold mb-6">Healthy Recipes</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {data?.items.map((recipe: Recipe) => (
          <RecipeCard key={recipe.id} recipe={recipe} />
        ))}
      </div>
      
      {pagination && (
        <Pagination 
          page={pagination.page}
          totalPages={pagination.total_pages}
          onPageChange={goToPage}
          className="mt-8"
        />
      )}
    </div>
  );
}

// RecipeCard component
interface RecipeCardProps {
  recipe: Recipe;
}

function RecipeCard({ recipe }: RecipeCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
      <h3 className="text-lg font-semibold mb-2">{recipe.name}</h3>
      
      <div className="flex items-center gap-4 text-sm text-gray-600 mb-4">
        <span className="flex items-center gap-1">
          🔥 {recipe.calories} cal
        </span>
        <span className="flex items-center gap-1">
          🥩 {recipe.protein}g protein
        </span>
        <span className="flex items-center gap-1">
          🍞 {recipe.carbs}g carbs
        </span>
        <span className="flex items-center gap-1">
          🥑 {recipe.fat}g fat
        </span>
      </div>
      
      <div className="flex items-center justify-between text-sm">
        <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded">
          {recipe.cuisine}
        </span>
        <span className="bg-green-100 text-green-800 px-2 py-1 rounded">
          {recipe.dietType}
        </span>
      </div>
      
      <div className="mt-4 text-sm text-gray-500">
        ⏱️ Prep: {recipe.prepTime}min | Cook: {recipe.cookTime}min
        {recipe.isHalal && <span className="ml-2 text-green-600">✓ Halal</span>}
      </div>
    </div>
  );
}
```

### Advanced Search Integration

```typescript
import { useSearch } from '../hooks/useSearch';
import { SearchFilters, AdvancedSearch } from '../components/search';

function AdvancedRecipeSearch() {
  const {
    query,
    filters,
    results,
    loading,
    error,
    setQuery,
    setFilters,
    search,
    clearFilters
  } = useSearch({
    endpoint: '/api/v1/nutrition-data/recipes',
    initialFilters: {
      maxCalories: 500,
      dietType: '',
      cuisine: '',
      isHalal: false
    }
  });

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-6">Recipe Search</h1>
        
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-1">
            <SearchFilters
              filters={filters}
              onFiltersChange={setFilters}
              onClear={clearFilters}
            />
          </div>
          
          <div className="lg:col-span-2">
            <AdvancedSearch
              query={query}
              onQueryChange={setQuery}
              onSearch={search}
              loading={loading}
              results={results}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
```

## Error Handling

### Standard Error Response

All endpoints return consistent error responses:

```json
{
  "status": "error",
  "error": "Invalid pagination parameters",
  "details": {
    "field": "page",
    "message": "Page must be a positive integer"
  },
  "timestamp": 1640995200
}
```

### Frontend Error Handling

```typescript
import { ErrorDisplay } from '../components/ui/ErrorDisplay';

function ApiExample() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch('/api/v1/nutrition-data/recipes?page=1&limit=20');
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Request failed');
      }
      
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      {loading && <LoadingSkeleton />}
      {error && <ErrorDisplay error={error} onRetry={fetchData} />}
      {data && <RecipeList recipes={data.items} />}
    </div>
  );
}
```

## Caching

### Cache Headers

API responses include cache headers:

```http
X-Cache: HIT
X-Cache-Key: /api/v1/nutrition-data/recipes?page=1&limit=20
X-Cache-TTL: 5m0s
```

### Client-Side Caching

```typescript
// Simple client-side cache implementation
class ApiCache {
  private cache = new Map<string, { data: any; timestamp: number; ttl: number }>();

  async get(key: string, fetcher: () => Promise<any>, ttl: number = 300000): Promise<any> {
    const cached = this.cache.get(key);
    const now = Date.now();

    if (cached && (now - cached.timestamp) < cached.ttl) {
      return cached.data;
    }

    const data = await fetcher();
    this.cache.set(key, { data, timestamp: now, ttl });
    return data;
  }
}

const apiCache = new ApiCache();

// Usage
const recipes = await apiCache.get(
  'recipes?page=1&limit=20',
  () => fetch('/api/v1/nutrition-data/recipes?page=1&limit=20').then(r => r.json()),
  300000 // 5 minutes
);
```

## Rate Limiting

API endpoints are rate-limited to prevent abuse:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995800
```

### Handling Rate Limits

```typescript
const handleRateLimit = async (apiCall: () => Promise<any>) => {
  try {
    return await apiCall();
  } catch (error) {
    if (error.status === 429) {
      const resetTime = parseInt(error.headers.get('X-RateLimit-Reset') || '0');
      const waitTime = Math.max(0, resetTime * 1000 - Date.now());
      
      await new Promise(resolve => setTimeout(resolve, waitTime));
      return await apiCall(); // Retry after waiting
    }
    throw error;
  }
};
```

## WebSocket Examples

### Real-time Updates

```typescript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to real-time updates');
  
  // Subscribe to nutrition data updates
  ws.send(JSON.stringify({
    type: 'subscribe',
    channel: 'nutrition-updates',
    filters: {
      dataType: 'recipes'
    }
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch (message.type) {
    case 'recipe-updated':
      updateRecipeInUI(message.data);
      break;
    case 'new-workout':
      addWorkoutToList(message.data);
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket connection closed');
  // Implement reconnection logic
};
```

## Testing Examples

### Unit Testing API Calls

```typescript
import { renderHook, act } from '@testing-library/react-hooks';
import { useRecipes } from '../hooks/useNutritionData';

// Mock fetch
global.fetch = jest.fn();

describe('useRecipes', () => {
  beforeEach(() => {
    fetch.mockClear();
  });

  test('should fetch recipes successfully', async () => {
    const mockRecipes = {
      status: 'success',
      items: [
        { id: 1, name: 'Test Recipe', calories: 300 }
      ],
      pagination: { page: 1, limit: 20, total: 1, total_pages: 1 }
    };

    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockRecipes)
    });

    const { result, waitForNextUpdate } = renderHook(() => useRecipes());

    expect(result.current.loading).toBe(true);

    await waitForNextUpdate();

    expect(result.current.loading).toBe(false);
    expect(result.current.data).toEqual(mockRecipes);
    expect(result.current.error).toBe(null);
  });
});
```

## Performance Optimization

### Debounced Search

```typescript
import { useDebouncedCallback } from '../hooks/useDebouncedCallback';

function SearchComponent() {
  const [query, setQuery] = useState('');
  const debouncedSearch = useDebouncedCallback(
    (searchQuery) => {
      // Perform API search
      searchRecipes(searchQuery);
    },
    500 // 500ms delay
  );

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    debouncedSearch(value);
  };

  return (
    <input
      type="text"
      value={query}
      onChange={handleInputChange}
      placeholder="Search recipes..."
    />
  );
}
```

### Infinite Scroll

```typescript
function InfiniteRecipeList() {
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const { data, loading, error } = useRecipes({ page, limit: 20 });

  const loadMore = useCallback(() => {
    if (!loading && hasMore) {
      setPage(prev => prev + 1);
    }
  }, [loading, hasMore]);

  useEffect(() => {
    if (data?.pagination) {
      setHasMore(data.pagination.page < data.pagination.total_pages);
    }
  }, [data]);

  return (
    <InfiniteScroll
      dataLength={data?.items?.length || 0}
      next={loadMore}
      hasMore={hasMore}
      loader={<LoadingSkeleton count={3} />}
    >
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {data?.items.map((recipe: Recipe) => (
          <RecipeCard key={recipe.id} recipe={recipe} />
        ))}
      </div>
    </InfiniteScroll>
  );
}
```

## Best Practices

1. **Always handle loading and error states**
2. **Implement proper pagination for large datasets**
3. **Use caching to reduce API calls**
4. **Implement retry logic for failed requests**
5. **Use TypeScript for type safety**
6. **Debounce search inputs to reduce API calls**
7. **Implement proper error boundaries**
8. **Use WebSockets for real-time updates**
9. **Monitor API usage and respect rate limits**
10. **Test API integrations thoroughly**

## Support

For API support and questions:
- Documentation: [API Docs](http://localhost:8080/docs)
- Health Check: [Health Status](http://localhost:8080/health)
- Support: support@doctorhealthy.com
