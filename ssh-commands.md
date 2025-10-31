# 🔧 SSH Commands for Server Management

## Server Details
- **IP:** 128.140.111.171
- **Username:** root
- **Password:** Khaled55400214.
- **Domain:** super.doctorhealthy1.com

## 📱 Basic SSH Access
```bash
ssh root@128.140.111.171
```

## 📋 Check Deployment Logs
```bash
# View all application logs
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose logs -f'

# View specific service logs
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose logs -f backend'
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose logs -f postgres'
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose logs -f redis'
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose logs -f nginx'
```

## 📊 Check Service Status
```bash
# Check all containers status
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose ps'

# Check system resources
ssh root@128.140.111.171 'htop'

# Check disk space
ssh root@128.140.111.171 'df -h'

# Check memory usage
ssh root@128.140.111.171 'free -h'
```

## 🔄 Service Management
```bash
# Restart all services
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose restart'

# Restart specific service
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose restart backend'

# Stop all services
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose down'

# Start all services
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose up -d'

# Rebuild and restart
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose up -d --build'
```

## 🏥 Health Checks
```bash
# Test application health
ssh root@128.140.111.171 'curl -f http://localhost:8080/health'

# Test API info
ssh root@128.140.111.171 'curl -f http://localhost:8080/api/info'

# Test nutrition analysis
ssh root@128.140.111.171 'curl -X POST http://localhost:8080/api/nutrition/analyze -H "Content-Type: application/json" -d "{\"food\": \"apple\", \"quantity\": 100, \"unit\": \"g\", \"checkHalal\": true}"'
```

## 🗄️ Database Management
```bash
# Connect to PostgreSQL
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose exec postgres psql -U nutrition_user -d nutrition_platform'

# Run database migrations
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1/backend && go run cmd/migrate/main.go -direction up'

# Seed database
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1/backend && go run cmd/seed/main.go'

# Check database tables
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose exec postgres psql -U nutrition_user -d nutrition_platform -c "\dt"'
```

## 🔍 Troubleshooting
```bash
# Check if deployment directory exists
ssh root@128.140.111.171 'ls -la /opt/trae-new-healthy1'

# Check Docker installation
ssh root@128.140.111.171 'docker --version && docker-compose --version'

# Check Go installation
ssh root@128.140.111.171 'go version'

# Check running processes
ssh root@128.140.111.171 'ps aux | grep docker'

# Check network connectivity
ssh root@128.140.111.171 'netstat -tlnp | grep :8080'
```

## 🔒 SSL Certificate Management
```bash
# Check SSL certificate status
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose exec certbot certbot certificates'

# Renew SSL certificate
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose run --rm certbot'

# Check nginx configuration
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose exec nginx nginx -t'
```

## 📁 File Management
```bash
# View application files
ssh root@128.140.111.171 'ls -la /opt/trae-new-healthy1'

# View environment configuration
ssh root@128.140.111.171 'cat /opt/trae-new-healthy1/.env.production'

# View docker-compose configuration
ssh root@128.140.111.171 'cat /opt/trae-new-healthy1/docker-compose.yml'

# Check application data
ssh root@128.140.111.171 'ls -la /opt/trae-new-healthy1/data'
```

## 🚨 Emergency Commands
```bash
# Force restart all services
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose down --remove-orphans && docker-compose up -d --build'

# Clean up Docker system
ssh root@128.140.111.171 'docker system prune -f'

# Check system logs
ssh root@128.140.111.171 'journalctl -u docker -f'

# Reboot server (last resort)
ssh root@128.140.111.171 'reboot'
```

## 📊 Monitoring Commands
```bash
# Real-time container stats
ssh root@128.140.111.171 'docker stats'

# Monitor application logs in real-time
ssh root@128.140.111.171 'cd /opt/trae-new-healthy1 && docker-compose logs -f --tail=50'

# Check application performance
ssh root@128.140.111.171 'curl -s http://localhost:8080/metrics'
```

---

## 💡 Quick Tips:
1. **Always navigate to `/opt/trae-new-healthy1`** before running docker-compose commands
2. **Use `-f` flag with logs** to follow real-time updates
3. **Check container status first** with `docker-compose ps` before troubleshooting
4. **Wait 30-60 seconds** after restart commands for services to fully start
5. **Use `Ctrl+C`** to exit log viewing mode