/**
 * Code Validation and API Health Check System
 * Provides comprehensive validation for JavaScript code and API endpoints
 */

class CodeValidator {
    constructor() {
        this.validationResults = {
            syntax: [],
            apis: [],
            performance: [],
            security: []
        };
        this.apiEndpoints = [
            '/api/recipes',
            '/api/workouts',
            '/api/diseases',
            '/api/health-check'
        ];
    }

    /**
     * Validate JavaScript syntax and common issues
     */
    validateJavaScriptSyntax(code, filename = 'unknown') {
        const issues = [];
        
        try {
            // Basic syntax validation using Function constructor
            new Function(code);
            issues.push({
                type: 'success',
                message: `Syntax validation passed for ${filename}`,
                severity: 'info'
            });
        } catch (error) {
            issues.push({
                type: 'syntax_error',
                message: `Syntax error in ${filename}: ${error.message}`,
                severity: 'error',
                line: this.extractLineNumber(error.message)
            });
        }

        // Check for common issues
        issues.push(...this.checkCommonIssues(code, filename));
        
        this.validationResults.syntax.push(...issues);
        return issues;
    }

    /**
     * Check for common JavaScript issues
     */
    checkCommonIssues(code, filename) {
        const issues = [];
        const lines = code.split('\n');

        lines.forEach((line, index) => {
            const lineNum = index + 1;
            
            // Check for console.log statements (should be removed in production)
            if (line.includes('console.log') && !line.includes('//')) {
                issues.push({
                    type: 'console_log',
                    message: `Console.log found in ${filename} at line ${lineNum}`,
                    severity: 'warning',
                    line: lineNum
                });
            }

            // Check for missing semicolons
            if (line.trim().match(/^(var|let|const|return).*[^;{}]$/)) {
                issues.push({
                    type: 'missing_semicolon',
                    message: `Possible missing semicolon in ${filename} at line ${lineNum}`,
                    severity: 'warning',
                    line: lineNum
                });
            }

            // Check for undefined variables
            if (line.includes('undefined') && !line.includes('typeof')) {
                issues.push({
                    type: 'undefined_usage',
                    message: `Potential undefined usage in ${filename} at line ${lineNum}`,
                    severity: 'warning',
                    line: lineNum
                });
            }

            // Check for eval usage (security risk)
            if (line.includes('eval(')) {
                issues.push({
                    type: 'security_risk',
                    message: `eval() usage detected in ${filename} at line ${lineNum} - security risk`,
                    severity: 'error',
                    line: lineNum
                });
            }
        });

        return issues;
    }

    /**
     * Validate API endpoints health
     */
    async validateAPIHealth() {
        const results = [];
        
        for (const endpoint of this.apiEndpoints) {
            try {
                const startTime = performance.now();
                const response = await fetch(endpoint, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    timeout: 5000
                });
                const endTime = performance.now();
                const responseTime = endTime - startTime;

                const result = {
                    endpoint,
                    status: response.status,
                    statusText: response.statusText,
                    responseTime: Math.round(responseTime),
                    healthy: response.ok,
                    timestamp: new Date().toISOString()
                };

                if (response.ok) {
                    result.message = `API ${endpoint} is healthy (${responseTime.toFixed(2)}ms)`;
                    result.severity = 'success';
                } else {
                    result.message = `API ${endpoint} returned ${response.status}: ${response.statusText}`;
                    result.severity = 'error';
                }

                results.push(result);
            } catch (error) {
                results.push({
                    endpoint,
                    status: 0,
                    statusText: 'Network Error',
                    responseTime: 0,
                    healthy: false,
                    message: `API ${endpoint} is unreachable: ${error.message}`,
                    severity: 'error',
                    timestamp: new Date().toISOString()
                });
            }
        }

        this.validationResults.apis = results;
        return results;
    }

    /**
     * Check data file integrity
     */
    async validateDataFiles() {
        const dataFiles = [
            '/data/nutrition_optimized.json',
            '/data/workouts.json',
            '/data/exercises.json',
            '/data/medical-plans.json'
        ];

        const results = [];

        for (const file of dataFiles) {
            try {
                const response = await fetch(file);
                if (response.ok) {
                    const data = await response.json();
                    results.push({
                        file,
                        valid: true,
                        message: `Data file ${file} is valid JSON`,
                        severity: 'success',
                        recordCount: Array.isArray(data) ? data.length : Object.keys(data).length
                    });
                } else {
                    results.push({
                        file,
                        valid: false,
                        message: `Data file ${file} not accessible: ${response.status}`,
                        severity: 'error'
                    });
                }
            } catch (error) {
                results.push({
                    file,
                    valid: false,
                    message: `Data file ${file} validation failed: ${error.message}`,
                    severity: 'error'
                });
            }
        }

        return results;
    }

    /**
     * Performance validation
     */
    validatePerformance() {
        const results = [];
        
        // Check page load time
        if (window.performance && window.performance.timing) {
            const loadTime = window.performance.timing.loadEventEnd - window.performance.timing.navigationStart;
            
            if (loadTime > 3000) {
                results.push({
                    type: 'performance',
                    message: `Page load time is slow: ${loadTime}ms`,
                    severity: 'warning',
                    value: loadTime
                });
            } else {
                results.push({
                    type: 'performance',
                    message: `Page load time is acceptable: ${loadTime}ms`,
                    severity: 'success',
                    value: loadTime
                });
            }
        }

        // Check memory usage
        if (window.performance && window.performance.memory) {
            const memoryUsage = window.performance.memory.usedJSHeapSize / 1024 / 1024;
            
            if (memoryUsage > 50) {
                results.push({
                    type: 'memory',
                    message: `High memory usage: ${memoryUsage.toFixed(2)}MB`,
                    severity: 'warning',
                    value: memoryUsage
                });
            } else {
                results.push({
                    type: 'memory',
                    message: `Memory usage is normal: ${memoryUsage.toFixed(2)}MB`,
                    severity: 'success',
                    value: memoryUsage
                });
            }
        }

        this.validationResults.performance = results;
        return results;
    }

    /**
     * Security validation
     */
    validateSecurity() {
        const results = [];

        // Check for HTTPS
        if (location.protocol !== 'https:' && location.hostname !== 'localhost') {
            results.push({
                type: 'security',
                message: 'Site is not using HTTPS',
                severity: 'warning'
            });
        } else {
            results.push({
                type: 'security',
                message: 'HTTPS is properly configured',
                severity: 'success'
            });
        }

        // Check for mixed content
        const scripts = document.querySelectorAll('script[src]');
        scripts.forEach(script => {
            if (script.src.startsWith('http://') && location.protocol === 'https:') {
                results.push({
                    type: 'mixed_content',
                    message: `Mixed content detected: ${script.src}`,
                    severity: 'error'
                });
            }
        });

        this.validationResults.security = results;
        return results;
    }

    /**
     * Run comprehensive validation
     */
    async runFullValidation() {
        console.log('🔍 Starting comprehensive validation...');
        
        const results = {
            timestamp: new Date().toISOString(),
            performance: this.validatePerformance(),
            security: this.validateSecurity(),
            apis: await this.validateAPIHealth(),
            dataFiles: await this.validateDataFiles()
        };

        // Generate summary
        const summary = this.generateValidationSummary(results);
        
        console.log('✅ Validation completed:', summary);
        return { results, summary };
    }

    /**
     * Generate validation summary
     */
    generateValidationSummary(results) {
        const summary = {
            total: 0,
            passed: 0,
            warnings: 0,
            errors: 0,
            categories: {}
        };

        Object.keys(results).forEach(category => {
            if (category === 'timestamp') return;
            
            const categoryResults = results[category];
            const categoryStats = {
                total: categoryResults.length,
                passed: categoryResults.filter(r => r.severity === 'success').length,
                warnings: categoryResults.filter(r => r.severity === 'warning').length,
                errors: categoryResults.filter(r => r.severity === 'error').length
            };

            summary.categories[category] = categoryStats;
            summary.total += categoryStats.total;
            summary.passed += categoryStats.passed;
            summary.warnings += categoryStats.warnings;
            summary.errors += categoryStats.errors;
        });

        summary.healthScore = summary.total > 0 ? 
            Math.round(((summary.passed + (summary.warnings * 0.5)) / summary.total) * 100) : 100;

        return summary;
    }

    /**
     * Extract line number from error message
     */
    extractLineNumber(errorMessage) {
        const match = errorMessage.match(/line (\d+)/i);
        return match ? parseInt(match[1]) : null;
    }

    /**
     * Display validation results in UI
     */
    displayResults(containerId = 'validation-results') {
        const container = document.getElementById(containerId);
        if (!container) {
            console.warn(`Container ${containerId} not found`);
            return;
        }

        const html = this.generateResultsHTML();
        container.innerHTML = html;
    }

    /**
     * Generate HTML for validation results
     */
    generateResultsHTML() {
        const { results, summary } = this.validationResults;
        
        return `
            <div class="validation-dashboard">
                <div class="validation-header">
                    <h3>System Validation Dashboard</h3>
                    <div class="health-score ${this.getHealthScoreClass(summary.healthScore)}">
                        Health Score: ${summary.healthScore}%
                    </div>
                </div>
                
                <div class="validation-summary">
                    <div class="stat-card success">
                        <span class="stat-number">${summary.passed}</span>
                        <span class="stat-label">Passed</span>
                    </div>
                    <div class="stat-card warning">
                        <span class="stat-number">${summary.warnings}</span>
                        <span class="stat-label">Warnings</span>
                    </div>
                    <div class="stat-card error">
                        <span class="stat-number">${summary.errors}</span>
                        <span class="stat-label">Errors</span>
                    </div>
                </div>
                
                <div class="validation-details">
                    ${Object.keys(summary.categories).map(category => 
                        this.generateCategoryHTML(category, results[category])
                    ).join('')}
                </div>
            </div>
        `;
    }

    /**
     * Generate HTML for a validation category
     */
    generateCategoryHTML(category, results) {
        return `
            <div class="validation-category">
                <h4>${category.charAt(0).toUpperCase() + category.slice(1)} Validation</h4>
                <div class="validation-items">
                    ${results.map(result => `
                        <div class="validation-item ${result.severity}">
                            <span class="validation-icon">${this.getSeverityIcon(result.severity)}</span>
                            <span class="validation-message">${result.message}</span>
                            ${result.responseTime ? `<span class="validation-time">${result.responseTime}ms</span>` : ''}
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    /**
     * Get CSS class for health score
     */
    getHealthScoreClass(score) {
        if (score >= 90) return 'excellent';
        if (score >= 75) return 'good';
        if (score >= 60) return 'fair';
        return 'poor';
    }

    /**
     * Get icon for severity level
     */
    getSeverityIcon(severity) {
        const icons = {
            success: '✅',
            warning: '⚠️',
            error: '❌',
            info: 'ℹ️'
        };
        return icons[severity] || 'ℹ️';
    }
}

// Global instance
window.codeValidator = new CodeValidator();

// Auto-run validation on page load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        setTimeout(() => window.codeValidator.runFullValidation(), 1000);
    });
} else {
    setTimeout(() => window.codeValidator.runFullValidation(), 1000);
}

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = CodeValidator;
}