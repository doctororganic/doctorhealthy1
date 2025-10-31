# 🚀 Auto-Factory Orchestrator - Deployment System

## Overview

This is a comprehensive automated build, test, fix, and deployment system for the Nutrition Platform. It includes parallel testing, real-time monitoring, SSH deployment, and automatic error fixing.

## 📁 Files

- **AUTO-FACTORY-ORCHESTRATOR.sh** - Main orchestration script
- **PARALLEL-TEST-RUNNER.sh** - Parallel test execution
- **REAL-TIME-MONITOR.sh** - Real-time application monitoring
- **SSH-DEPLOY.sh** - SSH-based deployment
- **AUTO-FIX-AGENT.sh** - Automatic error fixing

## 🎯 Quick Start

### 1. Move Scripts to Backend Directory

```bash
# From the models directory
mv *.sh ../
cd ..
chmod +x *.sh
```

### 2. Run Full Build & Test

```bash
./AUTO-FACTORY-ORCHESTRATOR.sh
```

### 3. Run Tests Only

```bash
./PARALLEL-TEST-RUNNER.sh
```

### 4. Monitor Application

```bash
./REAL-TIME-MONITOR.sh
```

### 5. Deploy to Remote Server

```bash
SSH_HOST=your-server.com SSH_USER=root ./SSH-DEPLOY.sh
```

## 📋 Features

### Auto-Factory Orchestrator
- ✅ Environment validation
- ✅ Parallel test execution
- ✅ Automatic error fixing
- ✅ Docker image building
- ✅ Integration testing
- ✅ Deployment package creation
- ✅ SSH deployment (optional)

### Parallel Test Runner
- ✅ Runs all test suites in parallel
- ✅ Coverage reports
- ✅ Detailed logging
- ✅ Real-time progress

### Real-Time Monitor
- ✅ Health check monitoring
- ✅ API endpoint testing
- ✅ Success rate tracking
- ✅ Auto-refresh display

### SSH Deploy
- ✅ Automated remote deployment
- ✅ Backup management
- ✅ Service management
- ✅ Health verification

## 🔧 Configuration

### Environment Variables

```bash
# SSH Deployment
export SSH_HOST=your-server.com
export SSH_USER=root
export SSH_PORT=22
export DEPLOY_PATH=/opt/nutrition-platform

# Monitoring
export API_URL=http://localhost:8080
export CHECK_INTERVAL=5
```

## 📊 Usage Examples

### Full CI/CD Pipeline

```bash
# 1. Fix any issues
./AUTO-FIX-AGENT.sh

# 2. Run full build and test
./AUTO-FACTORY-ORCHESTRATOR.sh

# 3. Deploy to production
SSH_HOST=prod.example.com ./SSH-DEPLOY.sh

# 4. Monitor deployment
API_URL=http://prod.example.com:8080 ./REAL-TIME-MONITOR.sh
```

### Development Workflow

```bash
# Run tests during development
./PARALLEL-TEST-RUNNER.sh

# Monitor local development
./REAL-TIME-MONITOR.sh
```

## 🐛 Troubleshooting

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

### Deployment Issues

```bash
# Check SSH connection
ssh $SSH_USER@$SSH_HOST "echo 'Connection OK'"

# Check remote service
ssh $SSH_USER@$SSH_HOST "systemctl status nutrition-platform"
```

## 📝 Logs

All logs are saved in `logs/` directory:
- `logs/orchestrator/` - Build and deployment logs
- `logs/tests/` - Test execution logs

## 🎨 Output Colors

- 🟢 Green - Success messages
- 🔵 Blue - Info messages
- 🟡 Yellow - Warnings
- 🔴 Red - Errors

## 🚀 Advanced Usage

### Custom Test Configuration

```bash
# Run specific test suite
cd backend
go test ./models/... -v -cover

# Run with race detection
go test -race ./...

# Run with timeout
go test -timeout 30s ./...
```

### Docker Deployment

```bash
# Build and run with Docker
docker build -t nutrition-platform .
docker run -p 8080:8080 nutrition-platform
```

### Continuous Integration

Add to your CI/CD pipeline:

```yaml
# .github/workflows/ci.yml
- name: Run Tests
  run: ./PARALLEL-TEST-RUNNER.sh

- name: Build and Deploy
  run: ./AUTO-FACTORY-ORCHESTRATOR.sh
```

## 📞 Support

For issues or questions, check the logs in `logs/` directory.
