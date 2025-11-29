#!/usr/bin/env bash
set -euo pipefail
BASE_URL="${1:-http://localhost:8080}"

check() {
  local name="$1" url="$2"
  local code time
  time=$(curl -o /dev/null -s -w '%{time_total}' "$url")
  code=$(curl -o /dev/null -s -w '%{http_code}' "$url")
  printf "%-40s %4s  %6ss\n" "$name" "$code" "$time"
}

echo "DoctorHealthy Comprehensive Smoke Tests -> $BASE_URL"
echo "=================================================="

# Health checks
echo ""
echo "Health Checks:"
check "GET /health" "$BASE_URL/health"

# Nutrition Data Endpoints
echo ""
echo "Nutrition Data Endpoints:"
check "GET /recipes" "$BASE_URL/api/v1/nutrition-data/recipes?limit=5"
check "GET /workouts" "$BASE_URL/api/v1/nutrition-data/workouts?limit=5"
check "GET /complaints" "$BASE_URL/api/v1/nutrition-data/complaints?limit=5"
check "GET /metabolism" "$BASE_URL/api/v1/nutrition-data/metabolism?limit=5"
check "GET /drugs-nutrition" "$BASE_URL/api/v1/nutrition-data/drugs-nutrition?limit=5"

# Disease/Injury/Vitamins Endpoints
echo ""
echo "Disease/Injury/Vitamins:"
check "GET /diseases" "$BASE_URL/api/v1/diseases?limit=5"
check "GET /injuries" "$BASE_URL/api/v1/injuries?limit=5"
check "GET /vitamins" "$BASE_URL/api/v1/vitamins-minerals/vitamins?limit=5"

# Pagination Tests
echo ""
echo "Pagination Tests:"
check "GET /recipes (page=2)" "$BASE_URL/api/v1/nutrition-data/recipes?page=2&limit=5"
check "GET /workouts (page=2)" "$BASE_URL/api/v1/nutrition-data/workouts?page=2&limit=5"

# Error Tests
echo ""
echo "Error Handling Tests:"
check "GET /recipes (invalid page)" "$BASE_URL/api/v1/nutrition-data/recipes?page=-1"
check "GET /recipes (invalid limit)" "$BASE_URL/api/v1/nutrition-data/recipes?limit=1000"

echo ""
echo "=================================================="
echo "Smoke tests complete!"
