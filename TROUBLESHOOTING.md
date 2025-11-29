# Troubleshooting Guide

## Overview

This guide provides solutions to common issues encountered while developing and deploying the DoctorHealthy nutrition platform.

## Quick Diagnostics

### Health Checks

```bash
# Backend health
curl http://localhost:8080/health

# Frontend health
curl http://localhost:3000

# Database connectivity
cd nutrition-platform/backend && make db-check

# Redis connectivity
redis-cli ping
```

### Log Locations

```bash
# Backend logs
tail -f nutrition-platform/backend/logs/app.log

# Frontend logs
# Check browser console for client-side errors

# System logs
tail -f /var/log/system.log  # Linux
tail -f /var/log/install.log  # macOS
```

## Backend Issues

### Port Already in Use

**Symptoms:**
```
Error: listen tcp :8080: bind: address already in use
```

**Solutions:**

**Option 1: Kill the process**
```bash
# Find process using port 8080
lsof -ti:8080

# Kill the process
lsof -ti:8080 | xargs kill -9

# Or use the script
./scripts/kill-port.sh 8080
```

**Option 2: Change port**
```bash
# Edit .env file
PORT=8081

# Or use environment variable
export PORT=8081 && make run
```

**Option 3: Find and kill manually**
```bash
# List all processes
ps aux | grep main.go

# Kill specific PID
kill -9 <PID>
```

### Database Connection Issues

**Symptoms:**
```
Error: database is locked
Error: unable to open database file
Error: no such table: users
```

**Solutions:**

**SQLite Issues:**
```bash
# Check database file exists
ls -la data/doctorhealthy.db

# Check permissions
chmod 664 data/doctorhealthy.db

# Check database integrity
sqlite3 data/doctorhealthy.db "PRAGMA integrity_check;"

# Reset database
make db-reset

# Recreate database
rm data/doctorhealthy.db
make migrate
make seed
```

**PostgreSQL Issues:**
```bash
# Check PostgreSQL is running
brew services list | grep postgresql
sudo systemctl status postgresql

# Start PostgreSQL
brew services start postgresql
sudo systemctl start postgresql

# Check connection
psql -h localhost -U postgres -d doctorhealthy_dev

# Reset database
dropdb doctorhealthy_dev
createdb doctorhealthy_dev
make migrate
```

### Migration Issues

**Symptoms:**
```
Error: migration failed
Error: table already exists
Error: no such table
```

**Solutions:**
```bash
# Check migration status
make migrate-status

# Force rerun migrations
make migrate-down
make migrate-up

# Reset all migrations
make db-reset

# Manually run specific migration
goose -dir "./migrations" up 001
```

### Module/Dependency Issues

**Symptoms:**
```
Error: module not found
Error: cannot find package
Error: version constraint failed
```

**Solutions:**
```bash
# Clean and reinstall dependencies
go clean -modcache
go mod download
go mod tidy

# Update specific module
go get -u github.com/some/module

# Clear vendor directory
rm -rf vendor/
make deps

# Reinitialize go modules
rm go.mod go.sum
go mod init doctorhealthy
go mod tidy
```

### Build/Compilation Issues

**Symptoms:**
```
Error: syntax error
Error: undefined: SomeFunction
Error: cannot use type as type
```

**Solutions:**
```bash
# Format code
make fmt
go fmt ./...

# Check for syntax errors
go build ./...

# Run linter
make lint

# Check specific package
go vet ./handlers/nutrition_data_handler.go

# Clear build cache
go clean -cache

# Rebuild
make build
```

### Redis Connection Issues

**Symptoms:**
```
Error: dial tcp: connection refused
Error: Redis server not available
```

**Solutions:**
```bash
# Check Redis is running
redis-cli ping

# Start Redis
brew services start redis
sudo systemctl start redis

# Check Redis config
redis-cli CONFIG GET bind

# Test connection
telnet localhost 6379

# Disable Redis (if not needed)
export REDIS_HOST="" && make run
```

## Frontend Issues

### Module Resolution Errors

**Symptoms:**
```
Error: Cannot find module 'next'
Error: Module not found: Can't resolve 'react'
Error: ENOENT: no such file or directory
```

**Solutions:**
```bash
# Clear node_modules
rm -rf node_modules package-lock.json

# Clear npm cache
npm cache clean --force

# Reinstall dependencies
npm install

# Update npm
npm install -g npm@latest

# Try with yarn
yarn install
```

### TypeScript Compilation Errors

**Symptoms:**
```
Error: Type 'string' is not assignable to type 'number'
Error: Property 'name' does not exist on type
Error: Cannot find module '../types/api'
```

**Solutions:**
```bash
# Check TypeScript configuration
cat tsconfig.json

# Clear TypeScript cache
rm -rf .next
rm .tsbuildinfo

# Type check
npm run type-check

# Update types
npm update @types/react @types/node

# Strict mode issues (temporarily disable)
# Add to tsconfig.json:
# "strict": false
```

### Next.js Development Server Issues

**Symptoms:**
```
Error: Port 3000 is already in use
Error: ENOENT: no such file or directory
Error: Module not found: Can't resolve 'fs'
```

**Solutions:**
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Clear Next.js cache
rm -rf .next

# Check Node.js version
node --version  # Should be 18+

# Clear npm cache
npm cache clean --force

# Rebuild
npm run build
npm run dev
```

### CSS/Styling Issues

**Symptoms:**
```
Error: CSS module not found
Tailwind classes not working
Styles not applying
```

**Solutions:**
```bash
# Check Tailwind config
cat tailwind.config.ts

# Rebuild CSS
npm run build

# Check PostCSS config
cat postcss.config.js

# Clear CSS cache
rm -rf .next/static/css

# Restart dev server
npm run dev
```

## Integration Issues

### CORS Issues

**Symptoms:**
```
Error: Access-Control-Allow-Origin missing
Error: CORS policy error
Network error in browser
```

**Solutions:**

**Backend CORS Configuration:**
```bash
# Check .env CORS settings
cat .env | grep CORS

# Update CORS origins
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# Restart backend
make run
```

**Frontend API Configuration:**
```bash
# Check frontend .env.local
cat .env.local | grep API

# Update API base URL
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### Authentication/JWT Issues

**Symptoms:**
```
Error: Invalid token
Error: Unauthorized
Error: JWT token expired
```

**Solutions:**
```bash
# Check JWT secret
cat .env | grep JWT_SECRET

# Generate new JWT secret
openssl rand -base64 32

# Update .env
JWT_SECRET=new-secret-key-here

# Restart backend
make run
```

### File Upload Issues

**Symptoms:**
```
Error: File too large
Error: Invalid file type
Error: Upload directory not found
```

**Solutions:**
```bash
# Check upload directory
ls -la uploads/

# Create upload directory
mkdir -p uploads
chmod 755 uploads/

# Check file size limit
cat .env | grep MAX_FILE_SIZE

# Update file size limit
MAX_FILE_SIZE=50MB
```

## Performance Issues

### Slow API Responses

**Symptoms:**
- API calls taking >5 seconds
- Frontend loading spinners persisting
- Database queries slow

**Diagnostics:**
```bash
# Check database queries
make db-debug

# Profile Go application
go tool pprof http://localhost:8080/debug/pprof/profile

# Check memory usage
ps aux | grep main.go

# Monitor API response times
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/api/v1/nutrition-data/recipes
```

**Solutions:**
```bash
# Enable Redis caching
export REDIS_HOST=localhost
make run

# Add database indexes
# Add to migration files:
# CREATE INDEX idx_recipes_cuisine ON recipes(cuisine);

# Optimize queries
# Use EXPLAIN QUERY PLAN
sqlite3 data/doctorhealthy.db "EXPLAIN QUERY PLAN SELECT * FROM recipes WHERE cuisine = 'Mediterranean';"
```

### Memory Leaks

**Symptoms:**
- Memory usage increasing over time
- Server crashes after hours
- Out of memory errors

**Diagnostics:**
```bash
# Monitor memory usage
top -p $(pgrep -f main.go)

# Check Go memory stats
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Monitor garbage collection
curl http://localhost:8080/debug/pprof/heap > /tmp/heap.prof
go tool pprof -text /tmp/heap.prof
```

**Solutions:**
```bash
# Restart server regularly
# Add to crontab:
# 0 */6 * * * pkill -f main.go && cd /path/to/backend && make run

# Optimize code for memory
# Review goroutine usage
curl http://localhost:8080/debug/pprof/goroutine?debug=1
```

## Environment-Specific Issues

### macOS Issues

**Symptoms:**
```
Error: command not found: brew
Error: xcode command line tools not found
Error: Go installation not found
```

**Solutions:**
```bash
# Install Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Xcode command line tools
xcode-select --install

# Install Go
brew install go

# Set Go path
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
source ~/.zshrc
```

### Linux Issues

**Symptoms:**
```
Error: permission denied
Error: package not found
Error: service not found
```

**Solutions:**
```bash
# Update package manager
sudo apt update && sudo apt upgrade  # Ubuntu/Debian
sudo yum update  # CentOS/RHEL

# Install dependencies
sudo apt install golang-go nodejs npm postgresql redis-server

# Fix permissions
sudo chown -R $USER:$USER /usr/local/go

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Windows Issues

**Symptoms:**
```
Error: command not found: make
Error: go: command not found
Error: path too long
```

**Solutions:**
```bash
# Install Chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install tools
choco install golang nodejs make

# Use Git Bash or WSL
# Recommended to use WSL2 for better Linux compatibility
```

## Testing Issues

### Test Failures

**Symptoms:**
```
FAIL: TestNutritionDataHandler_GetRecipes
Error: assertion failed
Error: timeout waiting for test
```

**Solutions:**
```bash
# Run specific test with verbose output
go test ./tests/integration -run TestNutritionDataHandler_GetRecipes -v

# Run with race detection
go test -race ./tests/integration

# Increase timeout
go test -timeout 30s ./tests/integration

# Check test database
make test-db-setup
make test-db-reset
```

### Coverage Issues

**Symptoms:**
```
Coverage: 0% of statements
Error: no test files
Error: coverage report generation failed
```

**Solutions:**
```bash
# Run tests with coverage
make test-coverage

# Check test files exist
ls -la tests/

# Run specific package coverage
go test -coverprofile=coverage.out ./handlers
go tool cover -html=coverage.out -o coverage.html

# Check test naming convention
# Tests must end with _test.go
```

## Production Issues

### Deployment Failures

**Symptoms:**
```
Error: 503 Service Unavailable
Error: 502 Bad Gateway
Error: Database connection failed
```

**Solutions:**
```bash
# Check service status
systemctl status doctorhealthy
systemctl status nginx

# Check logs
journalctl -u doctorhealthy -f
tail -f /var/log/nginx/error.log

# Restart services
systemctl restart doctorhealthy
systemctl restart nginx

# Check configuration
nginx -t
```

### SSL/TLS Issues

**Symptoms:**
```
Error: SSL certificate expired
Error: HTTPS connection failed
Error: certificate chain error
```

**Solutions:**
```bash
# Check certificate expiry
openssl x509 -in /path/to/cert.pem -noout -dates

# Renew certificate (Let's Encrypt)
certbot renew

# Check Nginx SSL config
nginx -t | grep -i ssl

# Restart Nginx
systemctl restart nginx
```

## Debugging Tools

### Backend Debugging

**Built-in Debugging:**
```bash
# Enable debug mode
export LOG_LEVEL=debug

# Enable pprof
curl http://localhost:8080/debug/pprof/

# Generate CPU profile
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Generate heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

**Delve Debugger:**
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug application
dlv debug main.go

# Debug test
dlv test ./tests/integration
```

### Frontend Debugging

**Browser DevTools:**
- **Elements**: Inspect HTML/CSS
- **Console**: Check JavaScript errors
- **Network**: Monitor API calls
- **Performance**: Analyze loading times
- **React DevTools**: Component state debugging

**VS Code Debugging:**
```json
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Next.js: debug",
      "type": "node-terminal",
      "request": "launch",
      "command": "npm run dev"
    }
  ]
}
```

## Getting Help

### Community Resources

- **GitHub Issues**: [Create issue](https://github.com/DrKhaled123/kiro-nutrition/issues)
- **Discussions**: [GitHub Discussions](https://github.com/DrKhaled123/kiro-nutrition/discussions)
- **Documentation**: Check `/docs` directory
- **API Docs**: http://localhost:8080/docs

### Support Channels

- **Email**: support@doctorhealthy.com
- **Discord**: [Join our community](https://discord.gg/doctorhealthy)
- **Stack Overflow**: Use tags `doctorhealthy`, `go`, `nextjs`

### Bug Reports

When reporting bugs, include:

1. **Environment**: OS, Go version, Node.js version
2. **Error Message**: Full error stack trace
3. **Steps to Reproduce**: Detailed reproduction steps
4. **Expected vs Actual**: What should happen vs what happens
5. **Logs**: Relevant log entries

### Feature Requests

For feature requests:

1. **Use Case**: Describe the problem you're solving
2. **Proposed Solution**: How you envision the feature working
3. **Alternatives Considered**: Other approaches you thought of
4. **Additional Context**: Any other relevant information

## Emergency Procedures

### Complete System Reset

If all else fails, perform a complete reset:

```bash
# Backup current data
make db-backup

# Stop all services
make stop
pkill -f "node\|go\|redis\|postgres"

# Clean all caches
go clean -modcache
npm cache clean --force
rm -rf .next node_modules
rm -rf vendor

# Reinstall everything
make setup
make migrate
make seed
make run
```

### Data Recovery

```bash
# Restore from backup
make db-restore backup_file.sql

# Check data integrity
sqlite3 data/doctorhealthy.db "PRAGMA integrity_check;"

# Verify key tables exist
sqlite3 data/doctorhealthy.db ".tables"
sqlite3 data/doctorhealthy.db "SELECT COUNT(*) FROM recipes;"
```

---

## Phase 1: Caching & Rate Limiting Issues

### Cache Not Working
**Symptoms:** X-Cache header always shows MISS
**Solutions:**
1. Check Redis is running: `redis-cli ping`
2. Verify REDIS_ADDR environment variable
3. Check cache middleware is enabled in main.go
4. Verify TTL is not too short

### Rate Limiting Too Strict
**Symptoms:** Getting 429 errors frequently
**Solutions:**
1. Check rate limit configuration in main.go
2. Verify rate limit window duration
3. Check if using Redis-backed limiter (more accurate)
4. Review rate limit headers: X-RateLimit-Limit, X-RateLimit-Remaining

### Performance Not Improved
**Symptoms:** Response times same with/without cache
**Solutions:**
1. Verify cache is actually hitting (check X-Cache: HIT)
2. Check cache TTL is appropriate
3. Verify endpoints are cacheable (GET requests only)
4. Check skip paths configuration

---

**Remember**: Most issues are configuration-related. Double-check environment variables, file permissions, and network connectivity before diving deep into debugging.
