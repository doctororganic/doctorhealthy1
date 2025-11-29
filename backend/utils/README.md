# Utility Functions for Order 2

This directory contains utility functions that will be helpful for implementing Order 2 tasks.

## Files

### `pagination.go`
Helper functions for pagination:
- `CalculatePagination()` - Calculates pagination metadata
- `GetOffset()` - Calculates SQL offset
- `ValidatePaginationParams()` - Validates pagination parameters

### `filtering.go`
Helper functions for filtering:
- `ParseFilterString()` - Parses filter strings like "field:value"
- `ParseMultipleFilters()` - Parses multiple filter strings
- `BuildWhereClause()` - Builds SQL WHERE clause from filters

### `search.go`
Helper functions for search:
- `BuildSearchClause()` - Builds SQL LIKE clause for searching
- `CalculateRelevanceScore()` - Calculates relevance score for results
- `SanitizeSearchQuery()` - Sanitizes search queries

### `sorting.go`
Helper functions for sorting:
- `ParseSortString()` - Parses sort strings like "field:asc"
- `BuildOrderByClause()` - Builds SQL ORDER BY clause
- `ValidateSortField()` - Validates sort fields

### `response.go`
Helper functions for API responses:
- `SuccessResponse()` - Sends success response
- `SuccessResponseWithPagination()` - Sends success response with pagination
- `ErrorResponse()` - Sends error response
- Various convenience functions for common HTTP status codes

### `validation.go`
Helper functions for validation:
- `ValidateQueryParams()` - Validates page and limit parameters
- `ValidateSortOrder()` - Validates sort order
- `ValidateFieldSelection()` - Validates field selection
- `ExtractQueryParams()` - Extracts and validates all query parameters

## Usage Examples

### Pagination
```go
import "nutrition-platform/utils"

page, limit := utils.ValidatePaginationParams(1, 20)
offset := utils.GetOffset(page, limit)
meta := utils.CalculatePagination(page, limit, total)
```

### Filtering
```go
filters, _ := utils.ParseMultipleFilters([]string{"cuisine:mediterranean", "calories:1300-1500"})
whereClause, args, _ := utils.BuildWhereClause(filters, allowedFields)
```

### Search
```go
searchQuery := utils.SearchQuery{
    Query:    "diabetes",
    Fields:   []string{"condition_en", "condition_ar"},
    Language: "both",
}
clause, args := utils.BuildSearchClause(searchQuery, "")
```

### Sorting
```go
sortConfig := utils.ParseSortString("name:desc")
orderBy, _ := utils.BuildOrderByClause(sortConfig, allowedFields, "id")
```

### Response
```go
utils.SuccessResponseWithPagination(c, data, pagination, filters)
```

