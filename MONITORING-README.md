# 📊 Complete Monitoring System - Loki + Grafana + Prometheus

## 🎯 Overview

A **production-ready monitoring stack** for the Nutrition Platform featuring:

- **🔍 Loki**: Log aggregation and querying
- **📊 Grafana**: Visualization and dashboards
- **📈 Prometheus**: Metrics collection and alerting
- **🚨 Alertmanager**: Notification management
- **🍎 Enhanced Node.js App**: Comprehensive metrics and monitoring

## 🚀 Quick Start

### One-Command Deployment

```bash
# Deploy complete monitoring system
./deploy-monitoring.sh deploy

# Or use docker-compose directly
docker-compose -f docker-compose.monitoring.yml up -d
```

### Access Points

| Service | URL | Purpose |
|---------|-----|---------|
| **Grafana** | http://localhost:3001 | Dashboards & Visualization |
| **Prometheus** | http://localhost:9090 | Metrics Collection |
| **Loki** | http://localhost:3100 | Log Aggregation |
| **Alertmanager** | http://localhost:9093 | Alert Management |
| **Nutrition App** | http://localhost:8080 | Main Application |

## 📋 What's Included

### ✅ **Complete Monitoring Stack**

#### **🔍 Log Aggregation (Loki)**
- **Real-time log streaming** from all containers
- **Structured log querying** with labels and filters
- **Log retention policies** with configurable storage
- **Docker container discovery** for automatic log collection

#### **📊 Visualization (Grafana)**
- **Pre-built dashboards** for nutrition platform metrics
- **Real-time charts** for request rates, response times, errors
- **Log visualization** with filtering and search
- **Alert management** interface

#### **📈 Metrics Collection (Prometheus)**
- **HTTP request metrics** (rate, duration, status codes)
- **Memory usage tracking** (heap, RSS, external)
- **Nutrition analysis metrics** (success rates, food types)
- **Error tracking** (types, endpoints, frequencies)

#### **🚨 Alerting (Alertmanager)**
- **Configurable alert rules** for critical events
- **Multiple notification channels** (webhook, email, Slack)
- **Alert grouping and suppression** to reduce noise
- **Alert history and management**

### ✅ **Enhanced Application Features**

#### **🔧 Advanced Metrics Collection**
```javascript
✅ HTTP request/response tracking
✅ Memory usage monitoring
✅ Nutrition analysis metrics
✅ Error rate tracking
✅ Response time histograms
✅ Custom business metrics
```

#### **📊 Prometheus Integration**
- **Standard metrics endpoint** at `/metrics`
- **Custom nutrition-specific metrics**
- **Performance monitoring** and alerting
- **Real-time metric updates**

#### **🔍 Enhanced Logging**
- **Structured JSON logging** with Winston
- **Request correlation IDs** for tracing
- **Performance metrics** in logs
- **Error context and stack traces**

## 🎨 **Pre-Built Dashboards**

### **📊 Nutrition Platform Dashboard**
- **Request Rate**: Real-time HTTP request tracking
- **Response Times**: 95th percentile latency monitoring
- **Error Rates**: Error tracking and alerting
- **Memory Usage**: Heap and system memory monitoring
- **Nutrition Metrics**: Analysis success rates and food types

### **📋 Log Dashboard**
- **Real-time log streaming** from all services
- **Log level filtering** (ERROR, WARN, INFO, DEBUG)
- **Container-specific logs** with filtering
- **Search and correlation** across log entries

## 🔧 **Configuration Files**

| File | Purpose | Location |
|------|---------|----------|
| `loki-config.yaml` | Loki server configuration | `monitoring/` |
| `promtail-config.yaml` | Log shipping configuration | `monitoring/` |
| `prometheus.yml` | Metrics collection config | `monitoring/` |
| `alertmanager.yml` | Alert notification config | `monitoring/` |
| `grafana-datasources.yaml` | Grafana data source config | `monitoring/` |
| `dashboard.json` | Pre-built dashboard | `monitoring/` |

## 🚨 **Alerting Rules**

### **Pre-Configured Alerts**
- **High Error Rate**: >5% error rate over 5 minutes
- **Memory Usage**: >80% memory utilization
- **Response Time**: >3 second average response time
- **Service Health**: Container health check failures

### **Custom Alert Configuration**
```yaml
# Example: High memory usage alert
- alert: HighMemoryUsage
  expr: nutrition_memory_usage_bytes > 100000000  # 100MB
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "High memory usage detected"
    description: "Memory usage is above 100MB for more than 2 minutes"
```

## 📊 **Metrics Collected**

### **Application Metrics**
- `nutrition_http_requests_total`: Total HTTP requests by method/route/status
- `nutrition_http_request_duration_seconds`: Request duration histogram
- `nutrition_memory_usage_bytes`: Memory usage by type
- `nutrition_analysis_total`: Nutrition analysis by food type and status
- `nutrition_errors_total`: Error tracking by type and endpoint

### **System Metrics**
- Container CPU and memory usage
- Network I/O and connections
- Disk usage and file system metrics
- Service health and uptime

## 🎯 **Live Demo Features**

### **🧪 Interactive Testing**
1. **Open Grafana**: http://localhost:3001 (admin/admin)
2. **Import Dashboard**: Use `monitoring/dashboard.json`
3. **Use Nutrition App**: http://localhost:8080
4. **Watch Live Metrics**: See real-time updates in dashboards

### **📈 Real-Time Visualization**
- **Request spikes** when using the nutrition analyzer
- **Memory trends** as the application processes requests
- **Error tracking** if any issues occur
- **Log streaming** showing application activity

## 🔧 **Customization**

### **Adding New Metrics**
```javascript
// In monitoringService.js
const customMetric = new promClient.Counter({
  name: 'nutrition_custom_metric',
  help: 'Description of custom metric',
  labelNames: ['label1', 'label2']
});

// Record metric
customMetric.inc({ label1: 'value1', label2: 'value2' });
```

### **Creating Custom Dashboards**
1. **Open Grafana** at http://localhost:3001
2. **Create new dashboard** or import existing
3. **Add panels** for your specific metrics
4. **Configure alerts** for critical thresholds

### **Log Filtering**
```javascript
// Filter logs by level
{container="nutrition-platform-monitored"} |= "ERROR"

// Filter by request ID
{container="nutrition-platform-monitored"} | json | requestId = "req_123"

// Filter by food type
{container="nutrition-platform-monitored"} | json | food = "apple"
```

## 🚀 **Production Deployment**

### **Environment Variables**
```bash
# Set these for production deployment
export ALLOWED_ORIGINS=https://yourdomain.com,http://localhost:3000
export LOG_LEVEL=warn
export PROMETHEUS_ENABLED=true
```

### **Docker Compose Override**
```yaml
# For production, create docker-compose.override.yml
version: '3.8'
services:
  nutrition-app:
    environment:
      - NODE_ENV=production
      - LOG_LEVEL=warn
    deploy:
      resources:
        limits:
          memory: 512M
```

## 🎉 **What You Get**

### **✅ Immediate Visual Feedback**
- **Real-time dashboards** showing live application metrics
- **Interactive charts** with zoom and filtering
- **Log streaming** with search and correlation
- **Alert notifications** for critical events

### **✅ Production-Ready Features**
- **Comprehensive monitoring** of all application aspects
- **Professional alerting** with notification management
- **Historical data** for trend analysis
- **Performance insights** for optimization

### **✅ Developer Experience**
- **Easy debugging** with correlated logs and metrics
- **Performance monitoring** for optimization
- **Error tracking** with detailed context
- **Capacity planning** with usage analytics

## 🎯 **Next Steps**

1. **🚀 Deploy the monitoring system**:
   ```bash
   ./deploy-monitoring.sh deploy
   ```

2. **🌐 Access Grafana** at http://localhost:3001

3. **📊 Import the dashboard** from `monitoring/dashboard.json`

4. **🧪 Test with nutrition app** at http://localhost:8080

5. **📈 Watch live metrics** as you use the application!

## 💡 **Benefits**

- **🔍 Complete observability** into application behavior
- **⚡ Fast debugging** with correlated logs and metrics
- **📊 Performance insights** for optimization
- **🚨 Proactive alerting** for issue prevention
- **📈 Capacity planning** with usage analytics

**Your nutrition platform now has enterprise-grade monitoring with beautiful visualizations!** 🎉📊