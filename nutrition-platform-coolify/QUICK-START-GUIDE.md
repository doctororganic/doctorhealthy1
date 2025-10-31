# 🚀 Quick Start Guide - Nutrition Platform

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Available Scripts](#available-scripts)
5. [Deployment](#deployment)
6. [Monitoring](#monitoring)
7. [Troubleshooting](#troubleshooting)

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Node.js 18+ (for frontend)
- SSH access to deployment server (optional)

## Installation

### 1. Move Scripts to Backend Directory

```bash
cd backend/models
mv *.sh ../
cd ..
chmod +x *.sh
```

### 2. Run Master Control

```bash
./MASTER-CONTROL.sh
```

This opens an interactive menu with all available options.

## Quick Start

### Option 1: Interactive Menu

```bash
./MASTER-CONTROL.sh
```

Select option 1 for complete setup.

### Option 2: Command Line

```bash
# Complete setup
./COMPLETE-SETUP.sh

# Or step by step:
./AUTO-FIX-AGENT.sh              # Fix issues
./PARALLEL-TEST-RUNNER.sh        # Run tests
./AUTO-FACTORY-ORCHESTRATOR.sh   # Build & package
```

## Available Scripts

### 🎯 Core Scripts

| Script | Description |
|--------|-------------|
| `MASTER-CONTROL.sh` | Interactive control panel |
| `COMPLETE-SETUP.sh` | One-command complete setup |
| `AUTO-FACTORY-ORCHESTRATOR.sh` | Full build & test pipeline |

### 🧪 Testing Scripts

| Script | Description |
|--------|-------------|
| `PARALLEL-TEST-RUNNER.sh` | Run all tests in parallel |
| `LOAD-TEST.sh` | Performance testing |
| `SECURITY-SCAN.sh` | Security vulnerability scan |

### 🚀 Deployment Scripts

| Script | Description |
|--------|-------------|
| `SSH-DEPLOY.sh` | Deploy to remote server |
| `DOCKER-COMPOSE-GENERATOR.sh` | Generate Docker Compose |
| `FRONTEND-BUILDER.sh` | Build React frontend |

### 🔧 Utility Scripts

| Script | Description |
|--------|-------------|
| `AUTO-FIX-AGENT.sh` | Automatic issue fixing |
| `REAL-TIME-MONITOR.sh` | Real-time monitoring |

## Deployment

### Local Development

```bash
# Start backend
cd backend
./bin/server

# In another terminal, monitor
./REAL-TIME-MONITOR.sh
```

### Production Deployment

```bash
# Option 1: Interactive
./MASTER-CONTROL.sh
# Select option 4

# Option 2: Command line
SSH_HOST=your-server.com SSH_USER=root ./SSH-DEPLOY.sh
```

### Docker Deployment

```bash
# Generate docker-compose
./DOCKER-COMPOSE-GENERATOR.sh

# Start services
docker-compose -f docker-compose.production.yml up -d
```

## Monitoring

### Real-Time Monitoring

```bash
# Local
./REAL-TIME-MONITOR.sh

# Remote
API_URL=http://your-server.com:8080 ./REAL-TIME-MONITOR.sh
```

### Load Testing

```bash
# Default (10 users, 100 requests each)
./LOAD-TEST.sh

# Custom
CONCURRENT_USERS=50 REQUESTS_PER_USER=200 ./LOAD-TEST.sh
```

## Troubleshooting

### Tests Failing

```bash
./AUTO-FIX-AGENT.sh
./PARALLEL-TEST-RUNNER.sh
```

### Build Errors

```bash
cd backend
go mod tidy
go clean -cache
go build -o bin/server ./cmd/server
```

### Docker Issues

```bash
docker system prune -f
docker-compose down -v
docker-compose up --build
```

### SSH Deployment Issues

```bash
# Test connection
ssh $SSH_USER@$SSH_HOST "echo 'OK'"

# Check logs
ssh $SSH_USER@$SSH_HOST "journalctl -u nutrition-platform -n 50"
```

## Environment Variables

### SSH Deployment

```bash
export SSH_HOST=your-server.com
export SSH_USER=root
export SSH_PORT=22
export DEPLOY_PATH=/opt/nutrition-platform
```

### Monitoring

```bash
export API_URL=http://localhost:8080
export CHECK_INTERVAL=5
```

### Load Testing

```bash
export CONCURRENT_USERS=10
export REQUESTS_PER_USER=100
```

## CI/CD Integration

### GitHub Actions

Copy `CI-CD-PIPELINE.yml` to `.github/workflows/`:

```bash
mkdir -p .github/workflows
cp backend/models/CI-CD-PIPELINE.yml .github/workflows/
```

### GitLab CI

```yaml
stages:
  - test
  - build
  - deploy

test:
  script:
    - cd backend/models
    - ./PARALLEL-TEST-RUNNER.sh

build:
  script:
    - cd backend/models
    - ./AUTO-FACTORY-ORCHESTRATOR.sh

deploy:
  script:
    - cd backend/models
    - ./SSH-DEPLOY.sh
```

## Advanced Usage

### Custom Test Configuration

```bash
# Run specific tests
cd backend
go test ./models/... -v

# With coverage
go test ./... -cover -coverprofile=coverage.out

# With race detection
go test -race ./...
```

### Docker Multi-Stage Build

```bash
# Build optimized image
docker build --target production -t nutrition-platform:prod .

# Build with cache
docker build --cache-from nutrition-platform:latest .
```

### Parallel Deployment

```bash
# Deploy to multiple servers
for server in server1 server2 server3; do
    SSH_HOST=$server ./SSH-DEPLOY.sh &
done
wait
```

## Support

- 📖 Full documentation: `DEPLOYMENT-README.md`
- 🐛 Issues: Check logs in `logs/` directory
- 💬 Questions: Review script comments

## Next Steps

1. ✅ Run complete setup
2. ✅ Test locally
3. ✅ Deploy to staging
4. ✅ Run load tests
5. ✅ Deploy to production
6. ✅ Monitor and optimize

Happy coding! 🎉
