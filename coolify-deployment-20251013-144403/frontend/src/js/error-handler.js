/**
 * Frontend Error Handler for Nutrition Platform PWA
 * Handles offline support, user feedback, and error recovery
 */

class ErrorHandler {
    constructor(config = {}) {
        this.config = {
            enableOfflineSupport: true,
            enableUserFeedback: true,
            enableMetrics: true,
            retryAttempts: 3,
            retryDelay: 1000,
            offlineStorageKey: 'nutrition_platform_offline_data',
            errorLogKey: 'nutrition_platform_error_log',
            maxErrorLogs: 100,
            ...config
        };

        this.isOnline = navigator.onLine;
        this.offlineQueue = [];
        this.errorLog = this.loadErrorLog();
        this.metrics = {
            totalErrors: 0,
            networkErrors: 0,
            apiErrors: 0,
            offlineOperations: 0,
            recoveredOperations: 0
        };

        this.init();
    }

    init() {
        this.setupEventListeners();
        this.setupServiceWorkerErrorHandling();
        this.setupGlobalErrorHandlers();
        this.loadOfflineQueue();
        
        if (this.config.enableMetrics) {
            this.startMetricsCollection();
        }
    }

    setupEventListeners() {
        // Network status monitoring
        window.addEventListener('online', () => {
            this.isOnline = true;
            this.handleOnlineRecovery();
        });

        window.addEventListener('offline', () => {
            this.isOnline = false;
            this.handleOfflineMode();
        });

        // Unhandled promise rejections
        window.addEventListener('unhandledrejection', (event) => {
            this.handleError(event.reason, 'unhandled_promise');
        });

        // Global error handler
        window.addEventListener('error', (event) => {
            this.handleError(event.error, 'global_error', {
                filename: event.filename,
                lineno: event.lineno,
                colno: event.colno
            });
        });
    }

    setupServiceWorkerErrorHandling() {
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.addEventListener('message', (event) => {
                if (event.data.type === 'ERROR') {
                    this.handleError(event.data.error, 'service_worker');
                }
            });
        }
    }

    setupGlobalErrorHandlers() {
        // Override fetch to add error handling
        const originalFetch = window.fetch;
        window.fetch = async (...args) => {
            try {
                const response = await originalFetch(...args);
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
                }
                return response;
            } catch (error) {
                return this.handleFetchError(error, args);
            }
        };
    }

    async handleFetchError(error, fetchArgs) {
        const [url, options = {}] = fetchArgs;
        
        this.metrics.networkErrors++;
        
        // If offline, queue the request
        if (!this.isOnline || error.name === 'TypeError') {
            return this.handleOfflineRequest(url, options, error);
        }

        // Retry logic for network errors
        if (this.shouldRetry(error)) {
            return this.retryRequest(url, options, error);
        }

        // Log error and show user feedback
        this.logError(error, 'fetch_error', { url, options });
        this.showUserFeedback('network_error', error.message);
        
        throw error;
    }

    async handleOfflineRequest(url, options, error) {
        if (!this.config.enableOfflineSupport) {
            throw error;
        }

        // Queue for later retry
        const queueItem = {
            id: Date.now() + Math.random(),
            url,
            options,
            timestamp: Date.now(),
            retryCount: 0
        };

        this.offlineQueue.push(queueItem);
        this.saveOfflineQueue();
        this.metrics.offlineOperations++;

        this.showUserFeedback('offline_queued', 'Request queued for when connection is restored');

        // Return cached data if available
        const cachedData = await this.getCachedData(url);
        if (cachedData) {
            return new Response(JSON.stringify(cachedData), {
                status: 200,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        throw new Error('No network connection and no cached data available');
    }

    async retryRequest(url, options, originalError, attempt = 1) {
        if (attempt > this.config.retryAttempts) {
            throw originalError;
        }

        await this.delay(this.config.retryDelay * attempt);

        try {
            const response = await fetch(url, options);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            this.metrics.recoveredOperations++;
            return response;
        } catch (error) {
            return this.retryRequest(url, options, originalError, attempt + 1);
        }
    }

    async handleOnlineRecovery() {
        this.showUserFeedback('online_recovered', 'Connection restored. Processing queued requests...');
        
        const queueCopy = [...this.offlineQueue];
        this.offlineQueue = [];
        
        for (const item of queueCopy) {
            try {
                await fetch(item.url, item.options);
                this.metrics.recoveredOperations++;
            } catch (error) {
                // Re-queue if still failing
                if (item.retryCount < this.config.retryAttempts) {
                    item.retryCount++;
                    this.offlineQueue.push(item);
                } else {
                    this.logError(error, 'offline_recovery_failed', item);
                }
            }
        }
        
        this.saveOfflineQueue();
        
        if (this.offlineQueue.length === 0) {
            this.showUserFeedback('all_recovered', 'All queued requests processed successfully');
        }
    }

    handleOfflineMode() {
        this.showUserFeedback('offline_mode', 'You are offline. Some features may be limited.');
    }

    handleError(error, type, context = {}) {
        this.metrics.totalErrors++;
        
        const errorInfo = {
            message: error.message || error,
            stack: error.stack,
            type,
            context,
            timestamp: Date.now(),
            userAgent: navigator.userAgent,
            url: window.location.href
        };

        this.logError(error, type, context);
        
        // Show appropriate user feedback
        this.showUserFeedback(type, this.getUserFriendlyMessage(error, type));
        
        // Send to analytics if available
        this.sendErrorToAnalytics(errorInfo);
    }

    logError(error, type, context = {}) {
        const logEntry = {
            id: Date.now() + Math.random(),
            timestamp: new Date().toISOString(),
            type,
            message: error.message || error,
            stack: error.stack,
            context,
            userAgent: navigator.userAgent,
            url: window.location.href
        };

        this.errorLog.push(logEntry);
        
        // Keep only the last N errors
        if (this.errorLog.length > this.config.maxErrorLogs) {
            this.errorLog = this.errorLog.slice(-this.config.maxErrorLogs);
        }
        
        this.saveErrorLog();
        console.error(`[${type}]`, error, context);
    }

    showUserFeedback(type, message) {
        if (!this.config.enableUserFeedback) return;

        const notification = this.createNotification(type, message);
        this.displayNotification(notification);
    }

    createNotification(type, message) {
        const typeConfig = {
            network_error: { icon: '🌐', color: '#f44336', duration: 5000 },
            offline_mode: { icon: '📱', color: '#ff9800', duration: 3000 },
            offline_queued: { icon: '⏳', color: '#2196f3', duration: 3000 },
            online_recovered: { icon: '✅', color: '#4caf50', duration: 3000 },
            all_recovered: { icon: '🎉', color: '#4caf50', duration: 4000 },
            api_error: { icon: '⚠️', color: '#f44336', duration: 5000 },
            validation_error: { icon: '📝', color: '#ff9800', duration: 4000 },
            default: { icon: 'ℹ️', color: '#2196f3', duration: 3000 }
        };

        const config = typeConfig[type] || typeConfig.default;
        
        return {
            id: Date.now() + Math.random(),
            type,
            message,
            ...config,
            timestamp: Date.now()
        };
    }

    displayNotification(notification) {
        // Create notification element
        const notificationEl = document.createElement('div');
        notificationEl.className = 'error-notification';
        notificationEl.innerHTML = `
            <div class="notification-content">
                <span class="notification-icon">${notification.icon}</span>
                <span class="notification-message">${notification.message}</span>
                <button class="notification-close" onclick="this.parentElement.parentElement.remove()">&times;</button>
            </div>
        `;
        
        notificationEl.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: ${notification.color};
            color: white;
            padding: 12px 16px;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            z-index: 10000;
            max-width: 400px;
            animation: slideIn 0.3s ease-out;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            font-size: 14px;
        `;

        // Add CSS animation
        if (!document.querySelector('#error-notification-styles')) {
            const styles = document.createElement('style');
            styles.id = 'error-notification-styles';
            styles.textContent = `
                @keyframes slideIn {
                    from { transform: translateX(100%); opacity: 0; }
                    to { transform: translateX(0); opacity: 1; }
                }
                .notification-content {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                }
                .notification-close {
                    background: none;
                    border: none;
                    color: white;
                    font-size: 18px;
                    cursor: pointer;
                    margin-left: auto;
                }
            `;
            document.head.appendChild(styles);
        }

        document.body.appendChild(notificationEl);
        
        // Auto-remove after duration
        setTimeout(() => {
            if (notificationEl.parentElement) {
                notificationEl.remove();
            }
        }, notification.duration);
    }

    getUserFriendlyMessage(error, type) {
        const messages = {
            network_error: 'Network connection problem. Please check your internet connection.',
            api_error: 'Server error occurred. Please try again later.',
            validation_error: 'Please check your input and try again.',
            offline_mode: 'You are currently offline. Some features may not be available.',
            timeout_error: 'Request timed out. Please try again.',
            default: 'An unexpected error occurred. Please try again.'
        };

        return messages[type] || messages.default;
    }

    async getCachedData(url) {
        try {
            const cache = await caches.open('nutrition-platform-cache');
            const response = await cache.match(url);
            return response ? await response.json() : null;
        } catch (error) {
            console.warn('Cache access failed:', error);
            return null;
        }
    }

    shouldRetry(error) {
        // Retry on network errors, timeouts, and 5xx server errors
        return error.name === 'TypeError' || 
               error.message.includes('timeout') ||
               (error.message.includes('HTTP 5'));
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    loadErrorLog() {
        try {
            const stored = localStorage.getItem(this.config.errorLogKey);
            return stored ? JSON.parse(stored) : [];
        } catch (error) {
            console.warn('Failed to load error log:', error);
            return [];
        }
    }

    saveErrorLog() {
        try {
            localStorage.setItem(this.config.errorLogKey, JSON.stringify(this.errorLog));
        } catch (error) {
            console.warn('Failed to save error log:', error);
        }
    }

    loadOfflineQueue() {
        try {
            const stored = localStorage.getItem(this.config.offlineStorageKey);
            this.offlineQueue = stored ? JSON.parse(stored) : [];
        } catch (error) {
            console.warn('Failed to load offline queue:', error);
            this.offlineQueue = [];
        }
    }

    saveOfflineQueue() {
        try {
            localStorage.setItem(this.config.offlineStorageKey, JSON.stringify(this.offlineQueue));
        } catch (error) {
            console.warn('Failed to save offline queue:', error);
        }
    }

    sendErrorToAnalytics(errorInfo) {
        // Send to analytics service if available
        if (window.gtag) {
            window.gtag('event', 'exception', {
                description: errorInfo.message,
                fatal: false,
                custom_map: {
                    error_type: errorInfo.type,
                    error_context: JSON.stringify(errorInfo.context)
                }
            });
        }
    }

    startMetricsCollection() {
        // Collect and report metrics periodically
        setInterval(() => {
            this.reportMetrics();
        }, 60000); // Every minute
    }

    reportMetrics() {
        console.log('Error Handler Metrics:', this.metrics);
        
        // Send to monitoring service if available
        if (window.navigator.sendBeacon) {
            const metricsData = JSON.stringify({
                timestamp: Date.now(),
                metrics: this.metrics,
                userAgent: navigator.userAgent
            });
            
            // Replace with your metrics endpoint
            // navigator.sendBeacon('/api/metrics', metricsData);
        }
    }

    // Public API methods
    
    getErrorLog() {
        return [...this.errorLog];
    }

    getMetrics() {
        return { ...this.metrics };
    }

    clearErrorLog() {
        this.errorLog = [];
        this.saveErrorLog();
    }

    clearOfflineQueue() {
        this.offlineQueue = [];
        this.saveOfflineQueue();
    }

    isOffline() {
        return !this.isOnline;
    }

    getOfflineQueueSize() {
        return this.offlineQueue.length;
    }

    // Enhanced validation methods
    validateClientData(clientData) {
        const errors = [];
        
        try {
            // Required fields validation
            if (!clientData.name || clientData.name.trim() === '') {
                errors.push('الاسم مطلوب - Name is required');
            }
            
            if (!clientData.age || isNaN(clientData.age) || clientData.age < 1 || clientData.age > 120) {
                errors.push('العمر يجب أن يكون بين 1 و 120 سنة - Age must be between 1 and 120');
            }
            
            if (!clientData.weight || isNaN(clientData.weight) || clientData.weight < 20 || clientData.weight > 500) {
                errors.push('الوزن يجب أن يكون بين 20 و 500 كيلو - Weight must be between 20 and 500 kg');
            }
            
            if (!clientData.height || isNaN(clientData.height) || clientData.height < 100 || clientData.height > 250) {
                errors.push('الطول يجب أن يكون بين 100 و 250 سم - Height must be between 100 and 250 cm');
            }
            
            if (!clientData.gender || !['male', 'female', 'ذكر', 'أنثى'].includes(clientData.gender)) {
                errors.push('الجنس مطلوب - Gender is required');
            }
            
            if (!clientData.activityLevel || !['sedentary', 'light', 'moderate', 'active', 'very_active'].includes(clientData.activityLevel)) {
                errors.push('مستوى النشاط مطلوب - Activity level is required');
            }
            
            if (!clientData.goal || !['lose_weight', 'maintain_weight', 'gain_weight', 'muscle_gain'].includes(clientData.goal)) {
                errors.push('الهدف مطلوب - Goal is required');
            }
            
        } catch (error) {
            this.logError({
                type: 'Validation Error',
                message: `Client data validation failed: ${error.message}`,
                data: clientData,
                timestamp: new Date().toISOString()
            });
            errors.push('خطأ في التحقق من البيانات - Data validation error');
        }
        
        return errors;
    }
    
    // Syntax validation for JSON data
    validateJSONSyntax(jsonString, context = 'Unknown') {
        try {
            JSON.parse(jsonString);
            return { valid: true, data: JSON.parse(jsonString) };
        } catch (error) {
            this.logError({
                type: 'JSON Syntax Error',
                message: `Invalid JSON in ${context}: ${error.message}`,
                data: jsonString.substring(0, 200) + '...',
                timestamp: new Date().toISOString()
            });
            return { 
                valid: false, 
                error: `خطأ في صيغة البيانات - JSON syntax error in ${context}`,
                details: error.message 
            };
        }
    }
    
    // Validate recipe data structure
    validateRecipeData(recipe) {
        const errors = [];
        
        try {
            if (!recipe.name || typeof recipe.name !== 'string') {
                errors.push('اسم الوصفة مطلوب - Recipe name is required');
            }
            
            if (!recipe.ingredients || !Array.isArray(recipe.ingredients) || recipe.ingredients.length === 0) {
                errors.push('المكونات مطلوبة - Ingredients are required');
            }
            
            if (!recipe.instructions || typeof recipe.instructions !== 'string') {
                errors.push('تعليمات التحضير مطلوبة - Instructions are required');
            }
            
            if (recipe.nutrition) {
                if (typeof recipe.nutrition !== 'object') {
                    errors.push('بيانات التغذية يجب أن تكون كائن - Nutrition data must be an object');
                } else {
                    const requiredNutrition = ['calories', 'protein', 'carbs', 'fat'];
                    requiredNutrition.forEach(nutrient => {
                        if (recipe.nutrition[nutrient] && isNaN(recipe.nutrition[nutrient])) {
                            errors.push(`${nutrient} يجب أن يكون رقم - ${nutrient} must be a number`);
                        }
                    });
                }
            }
            
        } catch (error) {
            this.logError({
                type: 'Recipe Validation Error',
                message: `Recipe validation failed: ${error.message}`,
                data: recipe,
                timestamp: new Date().toISOString()
            });
            errors.push('خطأ في التحقق من بيانات الوصفة - Recipe validation error');
        }
        
        return errors;
    }
    
    // Safe data access with fallbacks
    safeGet(obj, path, defaultValue = null) {
        try {
            const keys = path.split('.');
            let result = obj;
            
            for (const key of keys) {
                if (result === null || result === undefined || !(key in result)) {
                    return defaultValue;
                }
                result = result[key];
            }
            
            return result;
        } catch (error) {
            this.logError({
                type: 'Safe Access Error',
                message: `Failed to access path '${path}': ${error.message}`,
                timestamp: new Date().toISOString()
            });
            return defaultValue;
        }
    }
    
    // Enhanced error display with Arabic support
    showValidationErrors(errors, containerId = 'validationErrors') {
        const container = document.getElementById(containerId) || this.createValidationContainer(containerId);
        
        if (errors.length === 0) {
            container.style.display = 'none';
            return;
        }
        
        container.style.display = 'block';
        container.innerHTML = `
            <div class="alert alert-danger alert-dismissible fade show" role="alert">
                <div class="d-flex align-items-start">
                    <i class="fas fa-exclamation-triangle me-2 mt-1"></i>
                    <div class="flex-grow-1">
                        <h6 class="alert-heading mb-2">يرجى تصحيح الأخطاء التالية:</h6>
                        <ul class="mb-0">
                            ${errors.map(error => `<li>${error}</li>`).join('')}
                        </ul>
                    </div>
                    <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
                </div>
            </div>
        `;
        
        // Auto-scroll to errors
        container.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
    
    createValidationContainer(containerId) {
        const container = document.createElement('div');
        container.id = containerId;
        container.className = 'validation-errors-container mb-3';
        
        // Insert at the top of the main content area
        const mainContent = document.querySelector('.nutrition-content') || document.body;
        mainContent.insertBefore(container, mainContent.firstChild);
        
        return container;
    }
}

// Initialize global error handler
const errorHandler = new ErrorHandler({
    enableOfflineSupport: true,
    enableUserFeedback: true,
    enableMetrics: true,
    retryAttempts: 3,
    retryDelay: 1000
});

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ErrorHandler;
} else {
    window.ErrorHandler = ErrorHandler;
    window.errorHandler = errorHandler;
}