/**
 * API Health Monitoring System
 * Provides real-time monitoring and health checks for all API endpoints
 */

class APIHealthMonitor {
    constructor() {
        this.endpoints = {
            // Backend API endpoints
            backend: {
                baseUrl: 'http://localhost:8080',
                endpoints: [
                    { path: '/api/health', method: 'GET', critical: true },
                    { path: '/api/recipes', method: 'GET', critical: true },
                    { path: '/api/workouts', method: 'GET', critical: true },
                    { path: '/api/diseases', method: 'GET', critical: true },
                    { path: '/api/nutrition/calculate', method: 'POST', critical: true },
                    { path: '/api/plans/generate', method: 'POST', critical: true }
                ]
            },
            // Frontend static resources
            frontend: {
                baseUrl: window.location.origin,
                endpoints: [
                    { path: '/data/nutrition_optimized.json', method: 'GET', critical: true },
                    { path: '/data/workouts.json', method: 'GET', critical: true },
                    { path: '/data/exercises.json', method: 'GET', critical: false },
                    { path: '/data/medical-plans.json', method: 'GET', critical: false }
                ]
            }
        };
        
        this.healthStatus = {
            overall: 'unknown',
            backend: 'unknown',
            frontend: 'unknown',
            lastCheck: null,
            details: {}
        };
        
        this.monitoringInterval = null;
        this.alertThresholds = {
            responseTime: 2000, // 2 seconds
            errorRate: 0.1, // 10%
            consecutiveFailures: 3
        };
        
        this.metrics = {
            requests: 0,
            successes: 0,
            failures: 0,
            totalResponseTime: 0,
            consecutiveFailures: 0
        };
    }

    /**
     * Start continuous health monitoring
     */
    startMonitoring(intervalMs = 30000) {
        console.log('🔍 Starting API health monitoring...');
        
        // Initial check
        this.checkAllEndpoints();
        
        // Set up periodic monitoring
        this.monitoringInterval = setInterval(() => {
            this.checkAllEndpoints();
        }, intervalMs);
        
        return this;
    }

    /**
     * Stop health monitoring
     */
    stopMonitoring() {
        if (this.monitoringInterval) {
            clearInterval(this.monitoringInterval);
            this.monitoringInterval = null;
            console.log('⏹️ API health monitoring stopped');
        }
        return this;
    }

    /**
     * Check all configured endpoints
     */
    async checkAllEndpoints() {
        const results = {
            timestamp: new Date().toISOString(),
            backend: await this.checkEndpointGroup('backend'),
            frontend: await this.checkEndpointGroup('frontend')
        };

        this.updateHealthStatus(results);
        this.updateMetrics(results);
        this.triggerHealthEvents(results);
        
        return results;
    }

    /**
     * Check a group of endpoints
     */
    async checkEndpointGroup(groupName) {
        const group = this.endpoints[groupName];
        if (!group) {
            throw new Error(`Unknown endpoint group: ${groupName}`);
        }

        const results = {
            groupName,
            baseUrl: group.baseUrl,
            healthy: true,
            endpoints: [],
            summary: {
                total: group.endpoints.length,
                healthy: 0,
                unhealthy: 0,
                critical_failures: 0,
                avg_response_time: 0
            }
        };

        const promises = group.endpoints.map(endpoint => 
            this.checkSingleEndpoint(group.baseUrl, endpoint)
        );

        try {
            results.endpoints = await Promise.all(promises);
            
            // Calculate summary
            let totalResponseTime = 0;
            results.endpoints.forEach(result => {
                if (result.healthy) {
                    results.summary.healthy++;
                } else {
                    results.summary.unhealthy++;
                    if (result.critical) {
                        results.summary.critical_failures++;
                    }
                }
                totalResponseTime += result.responseTime || 0;
            });
            
            results.summary.avg_response_time = Math.round(
                totalResponseTime / results.endpoints.length
            );
            
            // Group is healthy if no critical failures
            results.healthy = results.summary.critical_failures === 0;
            
        } catch (error) {
            console.error(`Error checking ${groupName} endpoints:`, error);
            results.healthy = false;
            results.error = error.message;
        }

        return results;
    }

    /**
     * Check a single endpoint
     */
    async checkSingleEndpoint(baseUrl, endpoint) {
        const url = `${baseUrl}${endpoint.path}`;
        const startTime = performance.now();
        
        const result = {
            url,
            path: endpoint.path,
            method: endpoint.method,
            critical: endpoint.critical,
            healthy: false,
            status: null,
            statusText: null,
            responseTime: null,
            error: null,
            timestamp: new Date().toISOString()
        };

        try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 10000); // 10s timeout
            
            const response = await fetch(url, {
                method: endpoint.method,
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                signal: controller.signal,
                // Add sample data for POST requests
                ...(endpoint.method === 'POST' && {
                    body: JSON.stringify(this.getSampleRequestData(endpoint.path))
                })
            });
            
            clearTimeout(timeoutId);
            const endTime = performance.now();
            
            result.status = response.status;
            result.statusText = response.statusText;
            result.responseTime = Math.round(endTime - startTime);
            result.healthy = response.ok;
            
            // Additional checks for specific endpoints
            if (response.ok && endpoint.path.includes('.json')) {
                try {
                    const data = await response.json();
                    result.dataValid = Array.isArray(data) || typeof data === 'object';
                    result.recordCount = Array.isArray(data) ? data.length : Object.keys(data).length;
                } catch (jsonError) {
                    result.healthy = false;
                    result.error = 'Invalid JSON response';
                }
            }
            
        } catch (error) {
            const endTime = performance.now();
            result.responseTime = Math.round(endTime - startTime);
            result.error = error.name === 'AbortError' ? 'Request timeout' : error.message;
            result.healthy = false;
        }

        return result;
    }

    /**
     * Get sample request data for POST endpoints
     */
    getSampleRequestData(path) {
        const sampleData = {
            '/api/nutrition/calculate': {
                age: 25,
                weight: 70,
                height: 175,
                gender: 'male',
                activityLevel: 'moderate',
                goal: 'maintain'
            },
            '/api/plans/generate': {
                clientData: {
                    age: 25,
                    weight: 70,
                    height: 175
                },
                preferences: {
                    cuisine: 'mediterranean',
                    dietType: 'balanced'
                }
            }
        };
        
        return sampleData[path] || {};
    }

    /**
     * Update overall health status
     */
    updateHealthStatus(results) {
        this.healthStatus.lastCheck = results.timestamp;
        this.healthStatus.backend = results.backend.healthy ? 'healthy' : 'unhealthy';
        this.healthStatus.frontend = results.frontend.healthy ? 'healthy' : 'unhealthy';
        
        // Overall status
        if (results.backend.healthy && results.frontend.healthy) {
            this.healthStatus.overall = 'healthy';
        } else if (results.backend.summary.critical_failures > 0 || results.frontend.summary.critical_failures > 0) {
            this.healthStatus.overall = 'critical';
        } else {
            this.healthStatus.overall = 'degraded';
        }
        
        this.healthStatus.details = results;
    }

    /**
     * Update performance metrics
     */
    updateMetrics(results) {
        const allEndpoints = [...results.backend.endpoints, ...results.frontend.endpoints];
        
        allEndpoints.forEach(endpoint => {
            this.metrics.requests++;
            
            if (endpoint.healthy) {
                this.metrics.successes++;
                this.metrics.consecutiveFailures = 0;
            } else {
                this.metrics.failures++;
                this.metrics.consecutiveFailures++;
            }
            
            if (endpoint.responseTime) {
                this.metrics.totalResponseTime += endpoint.responseTime;
            }
        });
    }

    /**
     * Trigger health-related events
     */
    triggerHealthEvents(results) {
        // Dispatch custom events for health status changes
        const event = new CustomEvent('apiHealthUpdate', {
            detail: {
                status: this.healthStatus,
                results,
                metrics: this.getMetricsSummary()
            }
        });
        
        window.dispatchEvent(event);
        
        // Check for alerts
        this.checkAlerts(results);
    }

    /**
     * Check for alert conditions
     */
    checkAlerts(results) {
        const alerts = [];
        
        // Check consecutive failures
        if (this.metrics.consecutiveFailures >= this.alertThresholds.consecutiveFailures) {
            alerts.push({
                type: 'consecutive_failures',
                severity: 'critical',
                message: `${this.metrics.consecutiveFailures} consecutive API failures detected`
            });
        }
        
        // Check response times
        const allEndpoints = [...results.backend.endpoints, ...results.frontend.endpoints];
        const slowEndpoints = allEndpoints.filter(ep => 
            ep.responseTime && ep.responseTime > this.alertThresholds.responseTime
        );
        
        if (slowEndpoints.length > 0) {
            alerts.push({
                type: 'slow_response',
                severity: 'warning',
                message: `${slowEndpoints.length} endpoints responding slowly`,
                details: slowEndpoints.map(ep => `${ep.path}: ${ep.responseTime}ms`)
            });
        }
        
        // Check error rate
        const errorRate = this.metrics.requests > 0 ? this.metrics.failures / this.metrics.requests : 0;
        if (errorRate > this.alertThresholds.errorRate) {
            alerts.push({
                type: 'high_error_rate',
                severity: 'warning',
                message: `High error rate: ${(errorRate * 100).toFixed(1)}%`
            });
        }
        
        // Trigger alert events
        alerts.forEach(alert => {
            const alertEvent = new CustomEvent('apiHealthAlert', {
                detail: alert
            });
            window.dispatchEvent(alertEvent);
            
            console.warn('🚨 API Health Alert:', alert);
        });
    }

    /**
     * Get metrics summary
     */
    getMetricsSummary() {
        const avgResponseTime = this.metrics.requests > 0 ? 
            Math.round(this.metrics.totalResponseTime / this.metrics.requests) : 0;
        
        const successRate = this.metrics.requests > 0 ? 
            (this.metrics.successes / this.metrics.requests) * 100 : 100;
        
        return {
            totalRequests: this.metrics.requests,
            successRate: Math.round(successRate * 100) / 100,
            averageResponseTime: avgResponseTime,
            consecutiveFailures: this.metrics.consecutiveFailures
        };
    }

    /**
     * Get current health status
     */
    getHealthStatus() {
        return {
            ...this.healthStatus,
            metrics: this.getMetricsSummary()
        };
    }

    /**
     * Generate health report
     */
    generateHealthReport() {
        const status = this.getHealthStatus();
        const timestamp = new Date().toISOString();
        
        return {
            timestamp,
            summary: {
                overall_status: status.overall,
                backend_status: status.backend,
                frontend_status: status.frontend,
                last_check: status.lastCheck
            },
            metrics: status.metrics,
            details: status.details,
            recommendations: this.generateRecommendations(status)
        };
    }

    /**
     * Generate recommendations based on health status
     */
    generateRecommendations(status) {
        const recommendations = [];
        
        if (status.overall === 'critical') {
            recommendations.push({
                priority: 'high',
                message: 'Critical API failures detected. Check backend services immediately.',
                action: 'investigate_backend'
            });
        }
        
        if (status.metrics.averageResponseTime > this.alertThresholds.responseTime) {
            recommendations.push({
                priority: 'medium',
                message: 'API response times are slow. Consider optimizing backend performance.',
                action: 'optimize_performance'
            });
        }
        
        if (status.metrics.successRate < 95) {
            recommendations.push({
                priority: 'medium',
                message: 'API success rate is below 95%. Review error logs and improve reliability.',
                action: 'improve_reliability'
            });
        }
        
        return recommendations;
    }

    /**
     * Display health dashboard
     */
    displayHealthDashboard(containerId = 'health-dashboard') {
        const container = document.getElementById(containerId);
        if (!container) {
            console.warn(`Container ${containerId} not found`);
            return;
        }
        
        const status = this.getHealthStatus();
        const html = this.generateDashboardHTML(status);
        container.innerHTML = html;
    }

    /**
     * Generate dashboard HTML
     */
    generateDashboardHTML(status) {
        const statusClass = {
            healthy: 'success',
            degraded: 'warning',
            critical: 'error',
            unknown: 'info'
        }[status.overall] || 'info';
        
        return `
            <div class="health-dashboard">
                <div class="dashboard-header">
                    <h3>API Health Dashboard</h3>
                    <div class="overall-status ${statusClass}">
                        <span class="status-indicator"></span>
                        <span class="status-text">${status.overall.toUpperCase()}</span>
                    </div>
                </div>
                
                <div class="health-metrics">
                    <div class="metric-card">
                        <div class="metric-value">${status.metrics.successRate}%</div>
                        <div class="metric-label">Success Rate</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value">${status.metrics.averageResponseTime}ms</div>
                        <div class="metric-label">Avg Response Time</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value">${status.metrics.totalRequests}</div>
                        <div class="metric-label">Total Requests</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value">${status.metrics.consecutiveFailures}</div>
                        <div class="metric-label">Consecutive Failures</div>
                    </div>
                </div>
                
                <div class="service-status">
                    <div class="service-group">
                        <h4>Backend Services</h4>
                        <div class="service-indicator ${status.backend === 'healthy' ? 'healthy' : 'unhealthy'}">
                            ${status.backend === 'healthy' ? '✅' : '❌'} ${status.backend}
                        </div>
                    </div>
                    <div class="service-group">
                        <h4>Frontend Resources</h4>
                        <div class="service-indicator ${status.frontend === 'healthy' ? 'healthy' : 'unhealthy'}">
                            ${status.frontend === 'healthy' ? '✅' : '❌'} ${status.frontend}
                        </div>
                    </div>
                </div>
                
                <div class="last-check">
                    Last Check: ${status.lastCheck ? new Date(status.lastCheck).toLocaleString() : 'Never'}
                </div>
            </div>
        `;
    }
}

// Global instance
window.apiHealthMonitor = new APIHealthMonitor();

// Auto-start monitoring when page loads
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        window.apiHealthMonitor.startMonitoring();
    });
} else {
    window.apiHealthMonitor.startMonitoring();
}

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = APIHealthMonitor;
}