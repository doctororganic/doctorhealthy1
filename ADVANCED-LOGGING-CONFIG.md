# 📊 Advanced Logging Configuration Guide

## 🔗 Coolify Dashboard Link

**🌐 Access URL:** https://api.doctorhealthy1.com

## 📋 Overview

This guide provides advanced logging configuration for your nutrition platform deployment, including structured logging, log aggregation, and monitoring setup.

## 🔧 Advanced Logging Configuration

### 1. **Structured Logging Setup**

Create this logging configuration for production:

```bash
# Create advanced logging configuration
cat > advanced-logging.yml << 'EOF'
version: "3.8"

logging:
  formatters:
    json:
      format: json
      class: python.logging.formatter.JSONFormatter
      
    detailed:
      format: "[%(asctime)s] %(levelname)s [%(name)s:%(lineno)s] - %(message)s"
      datefmt: "%Y-%m-%d %H:%M:%S"
      
    simple:
      format: "%(levelname)s - %(message)s"

  handlers:
    console:
      class: logging.StreamHandler
      level: INFO
      formatter: json
      stream: ext://sys.stdout
      
    file:
      class: logging.handlers.RotatingFileHandler
      level: DEBUG
      formatter: detailed
      filename: /app/logs/app.log
      maxBytes: 10485760  # 10MB
      backupCount: 5
      
    error_file:
      class: logging.handlers.RotatingFileHandler
      level: ERROR
      formatter: detailed
      filename: /app/logs/error.log
      maxBytes: 10485760  # 10MB
      backupCount: 5
      
    security:
      class: logging.handlers.RotatingFileHandler
      level: INFO
      formatter: json
      filename: /app/logs/security.log
      maxBytes: 10485760  # 10MB
      backupCount: 10

  loggers:
    root:
      level: INFO
      handlers: [console, file]
      
    django:
      level: INFO
      handlers: [console, file]
      propagate: false
      
    django.request:
      level: WARNING
      handlers: [error_file]
      propagate: false
      
    app:
      level: DEBUG
      handlers: [console, file]
      propagate: false
      
    security:
      level: INFO
      handlers: [console, security]
      propagate: false
EOF
```

### 2. **Docker Logging Configuration**

Update your Docker Compose file with advanced logging:

```yaml
version: "3.8"

services:
  app:
    build: .
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "5"
        labels: "service=app,environment=production"
    volumes:
      - ./logs:/app/logs
    environment:
      - LOG_LEVEL=INFO
      - LOG_FORMAT=json
      - LOG_FILE=/app/logs/app.log

  nginx:
    image: nginx:alpine
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "5"
        labels: "service=nginx,environment=production"

  postgres:
    image: postgres:15
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=postgres,environment=production"

  redis:
    image: redis:7-alpine
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=redis,environment=production"
```

### 3. **Application Logging Configuration**

Create this logging configuration for your application:

```javascript
// config/logging.js
const winston = require('winston');
const path = require('path');

// Define log levels
const levels = {
  error: 0,
  warn: 1,
  info: 2,
  http: 3,
  debug: 4,
};

// Define colors for console output
const colors = {
  error: 'red',
  warn: 'yellow',
  info: 'green',
  http: 'magenta',
  debug: 'blue',
};

// Create Winston logger
const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.combine(
    winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  defaultMeta: { service: 'nutrition-platform' },
  transports: [
    // Console transport
    new winston.transports.Console({
      format: winston.format.combine(
        winston.format.colorize({ all: true, colors }),
        winston.format.simple()
      )
    }),
    
    // File transport
    new winston.transports.File({
      filename: path.join(__dirname, '../logs/app.log'),
      maxsize: 10485760, // 10MB
      maxFiles: 5
    }),
    
    // Error file transport
    new winston.transports.File({
      filename: path.join(__dirname, '../logs/error.log'),
      level: 'error',
      maxsize: 10485760, // 10MB
      maxFiles: 5
    }),
    
    // Security log transport
    new winston.transports.File({
      filename: path.join(__dirname, '../logs/security.log'),
      level: 'info',
      maxsize: 10485760, // 10MB
      maxFiles: 10
    })
  ]
});

// Create security logger
const securityLogger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
    winston.format.json()
  ),
  defaultMeta: { service: 'nutrition-platform-security' },
  transports: [
    new winston.transports.File({
      filename: path.join(__dirname, '../logs/security.log'),
      maxsize: 10485760, // 10MB
      maxFiles: 10
    })
  ]
});

// Create performance logger
const performanceLogger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
    winston.format.json()
  ),
  defaultMeta: { service: 'nutrition-platform-performance' },
  transports: [
    new winston.transports.File({
      filename: path.join(__dirname, '../logs/performance.log'),
      maxsize: 10485760, // 10MB
      maxFiles: 5
    })
  ]
});

module.exports = {
  logger,
  securityLogger,
  performanceLogger
};
```

### 4. **API Request Logging Middleware**

Create middleware for API request logging:

```javascript
// middleware/logging.js
const { logger } = require('../config/logging');

// Request ID generator
const generateRequestId = () => {
  return Math.random().toString(36).substring(2, 15);
};

// Request logging middleware
const requestLogger = (req, res, next) => {
  const requestId = generateRequestId();
  const startTime = Date.now();
  
  // Add request ID to request headers
  req.requestId = requestId;
  res.setHeader('X-Request-ID', requestId);
  
  // Log request
  logger.info('API Request', {
    requestId,
    method: req.method,
    url: req.url,
    userAgent: req.get('User-Agent'),
    ip: req.ip,
    timestamp: new Date().toISOString()
  });
  
  // Override res.end to log response
  const originalEnd = res.end;
  res.end = function(chunk, encoding) {
    const responseTime = Date.now() - startTime;
    
    // Log response
    logger.info('API Response', {
      requestId,
      method: req.method,
      url: req.url,
      statusCode: res.statusCode,
      responseTime: `${responseTime}ms`,
      timestamp: new Date().toISOString()
    });
    
    originalEnd.call(this, chunk, encoding);
  };
  
  next();
};

module.exports = { requestLogger, generateRequestId };
```

### 5. **Security Event Logging**

Create security event logging:

```javascript
// utils/securityLogger.js
const { securityLogger } = require('../config/logging');

// Log security events
const logSecurityEvent = (event, details = {}) => {
  securityLogger.info('Security Event', {
    event,
    details,
    timestamp: new Date().toISOString(),
    severity: details.severity || 'medium'
  });
  
  // For critical events, also send alerts
  if (details.severity === 'critical') {
    console.error('🚨 CRITICAL SECURITY EVENT:', event, details);
    // TODO: Send to monitoring service
  }
};

// Authentication events
const logAuthEvent = (event, userId, details = {}) => {
  logSecurityEvent(`AUTH_${event}`, {
    userId,
    ...details,
    category: 'authentication'
  });
};

// Authorization events
const logAuthzEvent = (event, userId, resource, details = {}) => {
  logSecurityEvent(`AUTHZ_${event}`, {
    userId,
    resource,
    ...details,
    category: 'authorization'
  });
};

// Data access events
const logDataAccess = (event, userId, resource, details = {}) => {
  logSecurityEvent(`DATA_${event}`, {
    userId,
    resource,
    ...details,
    category: 'data_access'
  });
};

module.exports = {
  logSecurityEvent,
  logAuthEvent,
  logAuthzEvent,
  logDataAccess
};
```

### 6. **Performance Logging**

Create performance monitoring logging:

```javascript
// utils/performanceLogger.js
const { performanceLogger } = require('../config/logging');

// Log performance metrics
const logPerformanceMetric = (metric, value, unit = 'ms', details = {}) => {
  performanceLogger.info('Performance Metric', {
    metric,
    value,
    unit,
    ...details,
    timestamp: new Date().toISOString()
  });
};

// Log API response time
const logApiResponseTime = (endpoint, method, responseTime, details = {}) => {
  logPerformanceMetric(`api_response_${method}_${endpoint}`, responseTime, 'ms', {
    endpoint,
    method,
    ...details,
    category: 'api_performance'
  });
};

// Log database query time
const logDbQueryTime = (query, executionTime, details = {}) => {
  logPerformanceMetric('db_query', executionTime, 'ms', {
    query: query.substring(0, 100) + (query.length > 100 ? '...' : ''),
    ...details,
    category: 'database_performance'
  });
};

// Log memory usage
const logMemoryUsage = (type, usage, details = {}) => {
  logPerformanceMetric(`memory_${type}`, usage, 'MB', {
    type,
    ...details,
    category: 'resource_usage'
  });
};

module.exports = {
  logPerformanceMetric,
  logApiResponseTime,
  logDbQueryTime,
  logMemoryUsage
};
```

### 7. **Docker Compose with Logging**

Create Docker Compose with comprehensive logging:

```yaml
version: "3.8"

services:
  app:
    build: .
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "5"
        labels: "service=app,environment=production"
    volumes:
      - ./logs:/app/logs
      - ./config/logging.yml:/app/config/logging.yml
    environment:
      - LOG_LEVEL=INFO
      - LOG_FORMAT=json
      - LOG_FILE=/app/logs/app.log
      - SECURITY_LOG_LEVEL=INFO
      - PERFORMANCE_LOG_LEVEL=INFO
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=postgres,environment=production"
    environment:
      - POSTGRES_DB=nutrition_platform
      - POSTGRES_USER=nutrition_user
      - POSTGRES_PASSWORD=ac287cc0e30f54afad53c6dc7e02fd0cccad979d62b75d75d97b1ede12daf8d5
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./logs/postgres:/var/log/postgresql

  redis:
    image: redis:7-alpine
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=redis,environment=production"
    command: redis-server --requirepass f606b2d16d6697e666ce78a8685574d042df15484ca8f18f39f2e67bf38dc09a
    volumes:
      - redis_data:/data
      - ./logs/redis:/var/log/redis

  nginx:
    image: nginx:alpine
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=nginx,environment=production"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d
      - ./logs/nginx:/var/log/nginx
    depends_on:
      - app

volumes:
  postgres_data:
  redis_data:
```

### 8. **Log Aggregation Setup**

Create log aggregation configuration:

```bash
# Create log aggregation script
cat > setup-log-aggregation.sh << 'EOF'
#!/bin/bash

# Setup log aggregation for production
echo "🔧 Setting up log aggregation..."

# Create fluentd configuration
cat > fluentd.conf << 'FLUENT'
<source>
  @type tail
  path /app/logs/*.log
  pos_file /var/log/fluentd/app.log.pos
  tag app.*
  format json
  time_format %Y-%m-%dT%H:%M:%S.%NZ
</source>

# Security events
<match app.**>
  @type rewrite_tag_filter
  <rule>
    key category
    pattern ^security$
    add_prefix security
  </rule>
</match>

# Performance metrics
<match app.**>
  @type rewrite_tag_filter
  <rule>
    key category
    pattern ^(api_performance|database_performance|resource_usage)$
    add_prefix performance
  </rule>
</match>

# Add metadata
<filter app.**>
  @type record_transformer
  <record>
    hostname #{Socket.gethostname}
    environment production
    service nutrition-platform
  </record>
</filter>

# Send to log aggregation service
<match **>
  @type elasticsearch
  host elasticsearch
  port 9200
  index_name nutrition-platform
  type_name _doc
</match>
FLUENT

# Create fluentd Docker compose
cat > docker-compose.logging.yml << 'EOF'
version: "3.8"

services:
  fluentd:
    image: fluent/fluentd:v1.14-debian-1
    volumes:
      - ./fluentd.conf:/fluent/etc/fluent.conf
      - ./logs:/var/log/app
      - ./logs:/var/log/fluentd
    ports:
      - "24224:24224"
    environment:
      - FLUENTD_CONF=fluent.conf

  elasticsearch:
    image: elasticsearch:7.9.2
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  kibana:
    image: kibana:7.9.2
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch

volumes:
  es_data:
EOF

echo "✅ Log aggregation configuration created!"
echo "📊 Start with: docker-compose -f docker-compose.logging.yml up -d"
EOF

chmod +x setup-log-aggregation.sh
```

## 📊 Monitoring Dashboard Configuration

### 1. **Grafana Dashboard**

Create Grafana configuration for monitoring:

```yaml
# grafana/provisioning/datasources/elasticsearch.yml
apiVersion: 1

datasources:
  - name: Elasticsearch
    type: elasticsearch
    access: proxy
    url: http://elasticsearch:9200
    database: nutrition-platform
    jsonData:
      timeField: "@timestamp"
      esVersion: 7.0.0
      logMessageField: message
      logLevelField: level
```

### 2. **Prometheus Metrics**

Create Prometheus configuration:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'nutrition-platform'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s
```

## 🚀 Deployment Instructions

### 1. **Setup Logging**
```bash
# Create log directories
mkdir -p nutrition-platform/logs/{app,security,performance,postgres,redis,nginx}

# Setup log aggregation
cd nutrition-platform
./setup-log-aggregation.sh
```

### 2. **Deploy with Logging**
```bash
# Deploy with advanced logging
docker-compose -f docker-compose.yml -f docker-compose.logging.yml up -d
```

### 3. **Access Logs**
```bash
# View application logs
docker-compose logs -f app

# View aggregated logs
curl http://localhost:9200/nutrition-platform/_search?pretty

# Access Kibana dashboard
open http://localhost:5601
```

## 📋 Log Management Best Practices

### 1. **Log Retention**
- Application logs: 30 days
- Security logs: 1 year
- Performance logs: 90 days
- Error logs: 90 days

### 2. **Log Levels**
- Production: INFO, WARN, ERROR
- Development: DEBUG, INFO, WARN, ERROR
- Security: All security events

### 3. **Log Rotation**
- Maximum file size: 10MB
- Maximum files: 5-10
- Compress old logs after rotation

### 4. **Log Privacy**
- No sensitive data in logs
- Hash personal identifiers
- Anonymize IP addresses

## 🔍 Log Analysis Queries

### 1. **Security Events**
```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"category": "security"}}
      ]
    }
  }
}
```

### 2. **Performance Metrics**
```json
{
  "query": {
    "bool": {
      "must": [
        {"range": {"value": {"gte": 1000}}}
      ]
    }
  }
}
```

### 3. **Error Analysis**
```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"level": "error"}}
      ]
    }
  }
}
```

---
**Last Updated:** October 13, 2025  
**Logging Status:** ✅ CONFIGURED  
**Monitoring Level:** 🔒 ENTERPRISE GRADE