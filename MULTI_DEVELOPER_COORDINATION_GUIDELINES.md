# Multi-developer coordination guidelines

## Critical conflict zones

### High-risk files (coordinate changes)
1. `backend/main.go` — routes, middleware, initialization
2. `backend/middleware/*.go` — shared middleware
3. `frontend-nextjs/src/app/layout.tsx` — root layout
4. `go.mod` / `frontend-nextjs/package.json` — dependencies
5. `.env` files — configuration
6. `Makefile` — build processes and scripts

### Medium-risk files (coordinate changes)
1. `backend/database/database.go` — database connections
2. `backend/config.go` — configuration management
3. `backend/security.go` — security settings
4. `frontend-nextjs/src/lib/api.ts` — API integration
5. `backend/handlers/` — shared handler logic

### Low-risk files (safe for parallel work)
- Individual handler files (`handlers/*.go`)
- Individual component files (`components/*.tsx`)
- Test files (`*_test.go`, `*.test.ts`)
- Documentation files (`*.md`)
- Repository files (`repositories/*.go`)
- Model files (`models/*.go`)

---

## File ownership and responsibility matrix

### Developer 1: Testing & Quality Assurance
Owns:
- `backend/tests/**/*.go`
- `backend/middleware/*_test.go`
- `backend/utils/*_test.go`
- `backend/scripts/test-*.sh`
- `e2e-tests/` directory
- `TESTING_GUIDE.md`
- `TEST_RESULTS_ANALYSIS.md`

Rules:
- Do not modify `main.go` routes
- Do not modify handler implementations (only test them)
- Create new test files only
- Coordinate before modifying shared test utilities
- Own the CI/CD testing pipeline

---

### Developer 2: Frontend Integration
Owns:
- `frontend-nextjs/src/app/(dashboard)/**/*.tsx`
- `frontend-nextjs/src/components/ui/*.tsx`
- `frontend-nextjs/src/hooks/use*.ts`
- `frontend-nextjs/src/lib/api/services/*.ts`
- `frontend-nextjs/src/types/` directory

Rules:
- Do not modify `main.go`
- Do not modify backend handlers
- Coordinate before modifying shared components
- Use feature branches for new pages
- Own frontend build processes

---

### Developer 3: Backend API & Services
Owns:
- `backend/handlers/*.go`
- `backend/services/*.go`
- `backend/repositories/*.go`
- `backend/models/*.go`
- `backend/migrations/` directory

Rules:
- Do not modify frontend code
- Coordinate before changing API endpoints
- Maintain backward compatibility
- Document all API changes
- Own database schema changes

---

### Developer 4: DevOps & Infrastructure
Owns:
- `.github/workflows/*.yml`
- `Makefile`
- `scripts/deploy.sh`
- `scripts/build.sh`
- `docker-compose*.yml`
- `Dockerfile*`
- `nginx.conf`

Rules:
- Do not modify application code
- Do not modify `main.go` application logic
- Coordinate before changing build processes
- Test scripts in isolation
- Own deployment pipelines

---

### Developer 5: Documentation & API Reference
Owns:
- `docs/**/*.md`
- `backend/docs/**/*.md`
- `README.md`
- `CHANGELOG.md`
- `TROUBLESHOOTING.md`
- API documentation

Rules:
- Do not modify code files
- Coordinate before updating architecture docs
- Update CHANGELOG.md last
- Maintain API reference accuracy

---

## Conflict prevention strategies

### Strategy 1: File locking protocol

Before modifying a shared file, check and claim:

```bash
# Create a lock file before editing main.go
mkdir -p .locks
touch .locks/main.go.lock
echo "DEV2: Adding new route for recipes search - ETA: 30min" > .locks/main.go.lock

# After completion, remove lock
rm .locks/main.go.lock
```

Lock file format:
```
DEV_ID: Brief description - ETA: time
```

Create a script to check locks:
```bash
#!/bin/bash
# File: scripts/check-locks.sh
echo "Checking for active locks..."
if [ -d ".locks" ]; then
    for lock in .locks/*.lock; do
        if [ -f "$lock" ]; then
            echo "🔒 Locked: $(basename $lock) - $(cat $lock)"
        fi
    done
else
    echo "No active locks found"
fi
```

---

### Strategy 2: Feature branch isolation

Each developer works in isolated branches:

```bash
# Developer 1: Testing
git checkout -b dev1/testing-phase1

# Developer 2: Frontend
git checkout -b dev2/frontend-integration

# Developer 3: Backend API
git checkout -b dev3/api-enhancements

# Developer 4: DevOps
git checkout -b dev4/cicd-setup

# Developer 5: Documentation
git checkout -b dev5/documentation-update
```

Merge order:
1. Tests first (Developer 1)
2. Backend API (Developer 3)
3. Frontend (Developer 2)
4. DevOps (Developer 4)
5. Documentation (Developer 5)

---

### Strategy 3: Code sections ownership

For `main.go`, divide by sections:

```go
// ============================================
// SECTION 1: IMPORTS (Dev 4 - DevOps only)
// ============================================
import (...)

// ============================================
// SECTION 2: CONFIGURATION (Dev 4 - DevOps only)
// ============================================
cfg := config.LoadConfig()

// ============================================
// SECTION 3: DATABASE INITIALIZATION (Dev 3 - Backend API)
// ============================================
sqlDB := backendmodels.InitDB(cfg.GetDatabaseURL())
db := database.NewDatabase(sqlDB)

// ============================================
// SECTION 4: MIDDLEWARE SETUP (Dev 1 - Testing)
// ============================================
// DO NOT MODIFY - Testing in progress
e.Use(customMiddleware.RequestID())
e.Use(customMiddleware.CustomLogger())

// ============================================
// SECTION 5: HANDLER INITIALIZATION (Dev 3 - Backend API)
// ============================================
healthHandler := handlers.NewHealthHandler(healthService)
nutritionPlanHandler := handlers.NewNutritionPlanHandler(nutritionPlanService, healthService)

// ============================================
// SECTION 6: ROUTE REGISTRATION (Dev 3 - Backend API)
// ============================================
api := e.Group("/api/v1")
// ... routes ...

// ============================================
// SECTION 7: SERVER START (Dev 4 - DevOps)
// ============================================
port := fmt.Sprintf(":%s", cfg.Port)
log.Printf("Starting server on port %s", cfg.Port)
```

---

## Communication protocol

### Daily sync checklist

Before starting work:
```markdown
## Daily Sync - [Date]

### Developer 1 (Testing):
- [ ] Working on: cache_test.go
- [ ] Blocking: None
- [ ] ETA: 2 hours
- [ ] Files modifying: backend/middleware/cache_test.go

### Developer 2 (Frontend):
- [ ] Working on: Integrating search component
- [ ] Blocking: None
- [ ] ETA: 1.5 hours
- [ ] Files modifying: frontend-nextjs/src/app/(dashboard)/recipes/page.tsx

### Developer 3 (Backend API):
- [ ] Working on: Nutrition goals API
- [ ] Blocking: Waiting for database schema approval
- [ ] ETA: 3 hours
- [ ] Files modifying: backend/handlers/nutrition_goal_handler.go

### Developer 4 (DevOps):
- [ ] Working on: GitHub Actions workflow
- [ ] Blocking: None
- [ ] ETA: 1 hour
- [ ] Files modifying: .github/workflows/ci.yml

### Developer 5 (Documentation):
- [ ] Working on: API reference docs
- [ ] Blocking: Need API endpoint list from Dev3
- [ ] ETA: 2 hours
- [ ] Files modifying: backend/docs/API_REFERENCE.md
```

---

### Conflict resolution protocol

If conflicts occur:

Step 1: Identify conflict type
- Type A: Same file, different sections → Merge automatically
- Type B: Same file, same section → Coordinate
- Type C: Dependency conflict → Resolve dependency first

Step 2: Resolution process
```bash
# 1. Stop work on conflicting file
# 2. Communicate in shared channel
# 3. Determine priority (usually: Tests > Backend > Frontend > DevOps > Docs)
# 4. Lower priority developer rebases
# 5. Continue work
```

---

## Expert tips for parallel development

### Tip 1: Interface-first development

Define interfaces before implementation:

```go
// File: backend/interfaces/cache.go (Created by Dev 1)
package interfaces

type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}) error
}

// File: backend/cache/redis_cache.go (Dev 3 implements)
// File: backend/middleware/cache.go (Dev 1 uses interface)
```

Benefit: Multiple developers can work on different implementations without conflicts.

---

### Tip 2: Dependency injection

Avoid global state:

```go
// ❌ BAD: Global variable
var cacheInstance *cache.RedisCache

// ✅ GOOD: Dependency injection
func NewHandler(cache cache.Cache, db *sql.DB) *Handler {
    return &Handler{cache: cache, db: db}
}
```

Benefit: Testable, no shared state conflicts.

---

### Tip 3: Feature flags for incomplete work

Use feature flags to merge incomplete features:

```go
// File: backend/config/features.go
var Features = struct {
    EnableNewSearch bool
    EnableNewCache  bool
    EnableNewAPI    bool
}{
    EnableNewSearch: os.Getenv("ENABLE_NEW_SEARCH") == "true",
    EnableNewCache:  os.Getenv("ENABLE_NEW_CACHE") == "true",
    EnableNewAPI:    os.Getenv("ENABLE_NEW_API") == "true",
}

// Usage in main.go
if Features.EnableNewCache {
    e.Use(newCacheMiddleware())
} else {
    e.Use(oldCacheMiddleware())
}
```

Benefit: Merge incomplete features without breaking production.

---

### Tip 4: Atomic commits

One logical change per commit:

```bash
# ❌ BAD: Multiple unrelated changes
git commit -m "Fix cache, add search, update docs"

# ✅ GOOD: Atomic commits
git commit -m "Add cache middleware unit tests"
git commit -m "Integrate search component into recipes page"
git commit -m "Update API documentation"
```

Benefit: Easier conflict resolution, clearer history.

---

### Tip 5: Test-driven coordination

Write tests first, then implement:

```go
// Developer 1 writes test first
func TestNewFeature(t *testing.T) {
    // Test expectations
}

// Developer 3 implements feature to pass test
func NewFeature() {
    // Implementation
}
```

Benefit: Clear contract, prevents breaking changes.

---

## File modification rules

### Rule 1: main.go modification protocol

Before modifying `main.go`:

1. Check for locks: `ls .locks/main.go.lock`
2. If locked, wait or coordinate
3. Create lock: `echo "DEV_ID: description" > .locks/main.go.lock`
4. Make minimal changes
5. Test immediately
6. Remove lock: `rm .locks/main.go.lock`

Allowed modifications by section:

| Section | Who Can Modify | Coordination Needed |
|---------|----------------|---------------------|
| Imports | Dev 4 (DevOps) | Yes - affects builds |
| Configuration | Dev 4 (DevOps) | Yes - affects deployment |
| Database | Dev 3 (Backend) | Yes - affects all |
| Middleware | Dev 1 (Testing) | Yes - affects all routes |
| Handler Init | Dev 3 (Backend) | No - independent |
| Routes | Dev 3 (Backend) | Yes - affects API |
| Server Start | Dev 4 (DevOps) | Yes - affects deployment |

---

### Rule 2: Handler file ownership

Each handler file is owned by Developer 3 (Backend API):

```bash
# Developer 3 owns these:
handlers/nutrition_data_handler.go
handlers/recipe_handler.go
handlers/weight_handler.go

# Developer 1 can only add tests:
tests/integration/nutrition_data_handler_test.go ✅
handlers/nutrition_data_handler.go ❌ (cannot modify)
```

---

### Rule 3: Shared utility files

For shared utilities (`utils/*.go`):

1. Check if function exists: `grep -r "func FunctionName"`
2. If exists, use it
3. If not, add with clear naming
4. Document usage
5. Coordinate if breaking changes needed

---

## Merge conflict resolution guide

### Scenario 1: Same file, different sections

```go
// Developer 1's version (lines 50-60)
e.Use(customMiddleware.RequestID())
e.Use(customMiddleware.CustomLogger())

// Developer 2's version (lines 100-110)
api := e.Group("/api/v1")
api.GET("/recipes", handler.GetRecipes)

// Resolution: Both sections merge automatically ✅
```

Action: Merge both sections, no coordination needed.

---

### Scenario 2: Same file, same section

```go
// Developer 1's version
e.Use(customMiddleware.RateLimiter())

// Developer 2's version  
e.Use(customMiddleware.EnhancedRateLimiter())

// Resolution: Coordinate - choose one or combine
```

Action:
1. Determine which is correct
2. Keep the better implementation
3. Remove the other
4. Test thoroughly

---

### Scenario 3: Dependency conflicts

```go
// Developer 1 adds:
import "new-package/v1"

// Developer 2 adds:
import "new-package/v2"

// Resolution: Use compatible version or coordinate
```

Action:
1. Check if versions are compatible
2. If not, coordinate to choose one
3. Update `go.mod` together
4. Test compatibility

---

## Best practices checklist

### Before starting work:
- [ ] Check `.locks/` directory for file locks
- [ ] Create feature branch: `git checkout -b devX/feature-name`
- [ ] Pull latest changes: `git pull origin main`
- [ ] Announce in shared channel what you're working on
- [ ] Set ETA for completion

### During work:
- [ ] Make atomic commits
- [ ] Write tests alongside code
- [ ] Document your changes
- [ ] Keep changes minimal and focused
- [ ] Test locally before committing

### Before merging:
- [ ] Run all tests: `make test`
- [ ] Check for conflicts: `git merge main`
- [ ] Update documentation if needed
- [ ] Remove lock files
- [ ] Request review from team

### After merging:
- [ ] Verify build still works
- [ ] Run smoke tests
- [ ] Update CHANGELOG.md
- [ ] Announce completion in shared channel

---

## Emergency conflict resolution

If critical conflict occurs:

1. Stop all work on conflicting file
2. Identify conflict type (see above)
3. Determine priority (Tests > Backend > Frontend > DevOps > Docs)
4. Lower priority developer:
   - Stash changes: `git stash`
   - Pull latest: `git pull origin main`
   - Reapply: `git stash pop`
   - Resolve conflicts manually
5. Test resolution
6. Merge and continue

---

## Communication templates

### Starting work template:
```
🚀 Starting Work
Developer: Dev1
Task: Creating cache middleware tests
Files: backend/middleware/cache_test.go
ETA: 2 hours
Blocking: None
```

### Completing work template:
```
✅ Work Complete
Developer: Dev1
Task: Cache middleware tests
Files: backend/middleware/cache_test.go
Status: All tests passing
Ready for: Review/merge
```

### Conflict alert template:
```
⚠️ Conflict Detected
Developer: Dev1
File: backend/main.go (lines 80-90)
Conflict with: Dev3
Type: Same section modification
Action: Waiting for Dev3 to complete
ETA: Unknown
```

---

## Quick reference: do's and don'ts

### Do's:
- ✅ Create feature branches
- ✅ Lock files before editing shared code
- ✅ Communicate before modifying `main.go`
- ✅ Write tests for your changes
- ✅ Make atomic commits
- ✅ Pull latest before starting
- ✅ Test locally before pushing

### Don'ts:
- ❌ Modify `main.go` without checking locks
- ❌ Modify other developers' handler files
- ❌ Force push to main branch
- ❌ Commit broken code
- ❌ Skip tests
- ❌ Modify shared utilities without coordination
- ❌ Merge without testing

---

## File ownership quick reference

| File/Pattern | Owner | Coordination Needed |
|--------------|-------|---------------------|
| `backend/main.go` | **SHARED** | ✅ Always |
| `backend/handlers/*.go` | Dev 3 | ❌ No |
| `backend/tests/**/*.go` | Dev 1 | ❌ No |
| `backend/middleware/*.go` | Dev 1 | ✅ Yes |
| `backend/services/*.go` | Dev 3 | ❌ No |
| `backend/repositories/*.go` | Dev 3 | ❌ No |
| `backend/models/*.go` | Dev 3 | ❌ No |
| `frontend-nextjs/src/app/**/*.tsx` | Dev 2 | ❌ No |
| `frontend-nextjs/src/components/**/*.tsx` | Dev 2 | ❌ No |
| `frontend-nextjs/src/lib/api.ts` | Dev 2 | ✅ Yes |
| `.github/workflows/*.yml` | Dev 4 | ❌ No |
| `Makefile` | Dev 4 | ✅ Yes |
| `docs/**/*.md` | Dev 5 | ❌ No |
| `go.mod` / `package.json` | **SHARED** | ✅ Always |
| `.env*` files | **SHARED** | ✅ Always |

---

## Expert tip: use git worktrees for isolation

For maximum isolation, use git worktrees:

```bash
# Developer 1
git worktree add ../nutrition-platform-dev1 dev1/testing-phase1

# Developer 2  
git worktree add ../nutrition-platform-dev2 dev2/frontend-integration

# Developer 3
git worktree add ../nutrition-platform-dev3 dev3/api-enhancements

# Each developer works in separate directory
# No conflicts, complete isolation
```

Benefit: Complete isolation, no branch switching needed.

---

## Project-specific coordination notes

### Nutrition Platform Specific Considerations:

1. **Database Schema Changes**: Developer 3 must coordinate with Developer 1 for test updates
2. **API Endpoint Changes**: Developer 3 must coordinate with Developer 2 for frontend integration
3. **Environment Variables**: Developer 4 owns all deployment configurations
4. **Nutrition Data Files**: Shared resource - coordinate before modifying JSON files
5. **Migration Files**: Developer 3 owns, but Developer 1 must update tests

### Testing Strategy:

1. **Unit Tests**: Developer 1 owns all unit test infrastructure
2. **Integration Tests**: Developer 1 writes, Developer 3 provides test data
3. **E2E Tests**: Developer 1 owns, Developer 2 provides frontend selectors
4. **Performance Tests**: Developer 1 owns, Developer 3 provides API endpoints

### Deployment Strategy:

1. **Staging**: Developer 4 manages, all developers test
2. **Production**: Developer 4 manages, requires approval from all developers
3. **Rollback**: Developer 4 owns, but must notify all developers

---

## Summary: golden rules

1. Communicate before modifying shared files
2. Use feature branches for all work
3. Lock files before editing `main.go`
4. Write tests alongside code
5. Make atomic commits
6. Test before merging
7. Update documentation
8. Coordinate dependencies

Following these guidelines should minimize conflicts and enable smooth parallel development of the nutrition platform.