# Production Deployment Guide

## 📋 Overview

This guide provides comprehensive instructions for deploying the Nutrition Platform to production environments with enterprise-grade reliability, security, and monitoring.

## 🎯 Prerequisites

### Infrastructure Requirements

#### **Minimum Production Environment**
- **CPU**: 4 cores minimum, 8 cores recommended
- **Memory**: 8GB minimum, 16GB recommended
- **Storage**: 100GB SSD minimum, 500GB recommended
- **Network**: 1Gbps connection, redundant paths

#### **Database Requirements**
- **PostgreSQL**: Version 13+ with replication support
- **Redis**: Version 6+ with cluster support
- **Backup Storage**: 1TB+ for automated backups
- **Connection Pooling**: pgBouncer or similar

#### **Security Requirements**
- **SSL/TLS**: Valid certificates for all endpoints
- **Firewall**: Configured application firewall
- **WAF**: Web Application Firewall (recommended)
- **DDoS Protection**: Cloud-based DDoS mitigation

### Software Dependencies

```bash
# Required Tools
- Docker 20.10+
- Docker Compose 2.0+
- Kubernetes 1.24+ (optional)
- Helm 3.0+ (if using Kubernetes)
- kubectl 1.24+ (if using Kubernetes)

# Development Tools
- Go 1.21+
- Node.js 18+
- PostgreSQL client tools
- Redis CLI tools
```

## 🚀 Deployment Strategies

### **1. Blue-Green Deployment (Recommended)**

Deploy to production with zero downtime using blue-green strategy:

```bash
# Deploy to production with blue-green strategy
./scripts/production-deploy.sh production --strategy blue-green
```

**Benefits:**
- Zero downtime
- Instant rollback capability
- Comprehensive health checks
- Traffic splitting support

### **2. Rolling Deployment**

Gradual deployment with health checks:

```bash
# Rolling deployment with health checks
./scripts/production-deploy.sh production --strategy rolling
```

**Benefits:**
- Resource efficient
- Gradual traffic shift
- Automatic rollback on failure

### **3. Canary Deployment**

Deploy to subset of users first:

```bash
# Canary deployment (10% traffic)
./scripts/production-deploy.sh production --strategy canary --traffic-percentage 10
```

**Benefits:**
- Risk mitigation
- Real user testing
- Gradual rollout

## 📁 Environment Configuration

### **Production Environment Setup**

1. **Create Production Configuration**
```bash
# Copy production template
cp config/production.env.example config/production.env

# Edit with production values
nano config/production.env
```

2. **Configure Database**
```bash
# PostgreSQL Configuration
DATABASE_URL="postgresql://user:password@primary-db:5432/nutrition_db?sslmode=require"
DATABASE_READ_REPLICA_URL="postgresql://user:password@replica-db:5432/nutrition_db?sslmode=require"

# Connection Pool Settings
DB_MAX_OPEN_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=10
DB_CONNECTION_MAX_LIFETIME=5m
```

3. **Configure Redis**
```bash
# Redis Configuration
REDIS_ADDR="redis-cluster:6379"
REDIS_PASSWORD="your-redis-password"
REDIS_DB=0

# Redis Cluster (if applicable)
REDIS_CLUSTER_ENABLED=true
REDIS_CLUSTER_NODES="redis-1:6379,redis-2:6379,redis-3:6379"
```

4. **Security Configuration**
```bash
# JWT Configuration
JWT_SECRET="your-super-secure-jwt-secret-key"
JWT_EXPIRY=24h

# CORS Configuration
CORS_ORIGINS="https://yourdomain.com,https://app.yourdomain.com"
CORS_METHODS="GET,POST,PUT,DELETE,OPTIONS"
CORS_HEADERS="Content-Type,Authorization"

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
RATE_LIMIT_BURST=200
```

## 🐳 Container Deployment

### **Docker Compose Deployment**

1. **Create Production Compose File**
```yaml
# docker-compose.production.yml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - backend
      - frontend

  backend:
    image: nutrition-platform/backend:latest
    environment:
      - ENV=production
    env_file:
      - config/production.env
    depends_on:
      - postgres
      - redis
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G

  frontend:
    image: nutrition-platform/frontend:latest
    environment:
      - NODE_ENV=production
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1'
          memory: 1G

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 2G

volumes:
  postgres_data:
  redis_data:
```

2. **Deploy Containers**
```bash
# Deploy to production
docker-compose -f docker-compose.production.yml up -d

# Scale services as needed
docker-compose -f docker-compose.production.yml up -d --scale backend=5
```

### **Kubernetes Deployment**

1. **Create Kubernetes Manifests**
```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: nutrition-platform

---
# k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: nutrition-platform
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: nutrition-platform/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        envFrom:
        - secretRef:
            name: nutrition-secrets
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
# k8s/backend-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
  namespace: nutrition-platform
spec:
  selector:
    app: backend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

---
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nutrition-ingress
  namespace: nutrition-platform
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    - app.yourdomain.com
    secretName: nutrition-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: backend-service
            port:
              number: 80
  - host: app.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

2. **Deploy to Kubernetes**
```bash
# Apply all manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n nutrition-platform
kubectl get services -n nutrition-platform
kubectl get ingress -n nutrition-platform
```

## 🔒 Security Configuration

### **SSL/TLS Setup**

1. **Generate SSL Certificates**
```bash
# Using Let's Encrypt
certbot certonly --standalone -d yourdomain.com -d api.yourdomain.com

# Copy certificates to nginx
cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ./ssl/
cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ./ssl/
```

2. **Configure Nginx SSL**
```nginx
# nginx/nginx.conf
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/ssl/certs/fullchain.pem;
    ssl_certificate_key /etc/ssl/certs/privkey.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    location / {
        proxy_pass http://backend-service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### **Firewall Configuration**

```bash
# UFW Configuration (Ubuntu)
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw enable

# iptables Rules
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -j DROP
```

## 📊 Monitoring & Observability

### **Health Checks**

The application provides comprehensive health check endpoints:

```bash
# Liveness Probe
GET /health/live
# Response: {"status": "alive", "timestamp": "2023-12-01T10:00:00Z"}

# Readiness Probe
GET /health/ready
# Response: {"status": "healthy", "ready": true, ...}

# Detailed Health
GET /health
# Response: Comprehensive health status with component details
```

### **Metrics Collection**

1. **Prometheus Metrics**
```bash
# Metrics endpoint
GET /metrics
# Response: Prometheus-formatted metrics

# Custom metrics
http_requests_total{method="GET",endpoint="/api/users",status_code="200"} 1234
http_request_duration_seconds{method="GET",endpoint="/api/users"} 0.123
db_connections_active 15
cache_hits_total 5678
```

2. **Grafana Dashboard**
```json
{
  "dashboard": {
    "title": "Nutrition Platform",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

### **Alerting Rules**

```yaml
# prometheus/alerts.yml
groups:
- name: nutrition-platform
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status_code=~"5.."}[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} errors per second"

  - alert: DatabaseConnectionsHigh
    expr: db_connections_active > 20
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High database connections"
      description: "Database has {{ $value }} active connections"

  - alert: MemoryUsageHigh
    expr: memory_usage_bytes / 1024 / 1024 / 1024 > 8
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage"
      description: "Memory usage is {{ $value }}GB"
```

## 🔄 Backup & Recovery

### **Database Backup**

1. **Automated Backups**
```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/nutrition_db_$DATE.sql"

# Create backup
pg_dump $DATABASE_URL > $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

# Upload to cloud storage (AWS S3 example)
aws s3 cp $BACKUP_FILE.gz s3://your-backup-bucket/

# Clean old backups (keep 30 days)
find $BACKUP_DIR -name "*.gz" -mtime +30 -delete
```

2. **Restore Database**
```bash
# Restore from backup
gunzip -c backup_file.sql.gz | psql $DATABASE_URL

# Or using restore script
./scripts/restore-database.sh backup_file.sql.gz
```

### **Disaster Recovery**

1. **Recovery Procedures**
```bash
# 1. Stop application
docker-compose down

# 2. Restore database
./scripts/restore-database.sh latest_backup.sql.gz

# 3. Verify data integrity
./scripts/verify-database.sh

# 4. Start application
docker-compose up -d

# 5. Health check
curl -f https://api.yourdomain.com/health/ready
```

## 🚨 Incident Response

### **Common Issues**

1. **High CPU Usage**
```bash
# Check CPU usage
docker stats

# Scale application
docker-compose up -d --scale backend=5

# Check for memory leaks
go tool pprof http://localhost:8080/debug/pprof/heap
```

2. **Database Connection Issues**
```bash
# Check database connections
docker exec postgres psql -U $POSTGRES_USER -c "SELECT * FROM pg_stat_activity;"

# Reset connections
docker-compose restart postgres
```

3. **Redis Issues**
```bash
# Check Redis status
docker exec redis redis-cli ping

# Check Redis memory
docker exec redis redis-cli info memory

# Clear Redis cache (if needed)
docker exec redis redis-cli FLUSHALL
```

### **Rollback Procedures**

1. **Quick Rollback**
```bash
# Rollback to previous version
./scripts/production-deploy.sh --rollback v1.2.3 production
```

2. **Manual Rollback**
```bash
# Stop current deployment
docker-compose down

# Switch to previous image
docker pull nutrition-platform/backend:v1.2.3
docker pull nutrition-platform/frontend:v1.2.3

# Update compose file with previous tags
# Then redeploy
docker-compose up -d
```

## 📋 Deployment Checklist

### **Pre-Deployment Checklist**

- [ ] Environment configuration validated
- [ ] Database backups created
- [ ] SSL certificates valid
- [ ] Security scan passed
- [ ] Performance tests passed
- [ ] Health checks configured
- [ ] Monitoring setup verified
- [ ] Rollback procedures tested
- [ ] Team notified of deployment
- [ ] Maintenance window scheduled (if needed)

### **Post-Deployment Checklist**

- [ ] Health checks passing
- [ ] Monitoring alerts verified
- [ ] Performance metrics normal
- [ ] Error rates within thresholds
- [ ] User feedback collected
- [ ] Documentation updated
- [ ] Team notified of completion
- [ ] Success metrics recorded

## 🔧 Troubleshooting

### **Common Deployment Issues**

1. **Container Startup Failures**
```bash
# Check container logs
docker-compose logs backend
docker-compose logs frontend

# Check resource usage
docker stats

# Verify configuration
docker-compose config
```

2. **Network Connectivity Issues**
```bash
# Check service discovery
docker network ls
docker network inspect nutrition-platform_default

# Test connectivity
docker exec backend ping postgres
docker exec backend ping redis
```

3. **SSL Certificate Issues**
```bash
# Check certificate validity
openssl x509 -in /etc/ssl/certs/fullchain.pem -text -noout

# Test SSL configuration
openssl s_client -connect yourdomain.com:443
```

## 📚 Additional Resources

- [API Documentation](./API_REFERENCE.md)
- [Troubleshooting Guide](../TROUBLESHOOTING.md)
- [Development Setup](../DEVELOPMENT_SETUP.md)
- [Monitoring Configuration](./MONITORING_SETUP.md)

## 🆘 Support

For production deployment support:

1. **Emergency Contact**: [devops-team@yourcompany.com]
2. **Slack Channel**: #production-support
3. **Documentation**: [Internal Wiki]
4. **Monitoring Dashboard**: [Grafana URL]

---

**Note**: This guide is regularly updated. Always check the latest version before deploying to production.
