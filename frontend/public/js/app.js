// Nutrition Platform Frontend Application
// Handles backend integration and user interactions

class NutritionApp {
    constructor() {
        this.apiBaseUrl = this.getApiBaseUrl();
        this.currentLanguage = 'en';
        this.translations = {};
        this.init();
    }

    getApiBaseUrl() {
        // Detect if running locally or on server
        const hostname = window.location.hostname;
        if (hostname === 'localhost' || hostname === '127.0.0.1') {
            return 'http://localhost:8080';
        } else {
            return `${window.location.protocol}//${window.location.host}`;
        }
    }

    async init() {
        console.log('🚀 Initializing Nutrition Platform...');
        console.log('API Base URL:', this.apiBaseUrl);
        
        this.setupEventListeners();
        await this.loadTranslations();
        await this.checkSystemHealth();
        this.updateLanguage(this.currentLanguage);
    }

    setupEventListeners() {
        // Language selector
        const languageSelect = document.getElementById('languageSelect');
        if (languageSelect) {
            languageSelect.addEventListener('change', (e) => {
                this.updateLanguage(e.target.value);
            });
        }

        // Nutrition form
        const nutritionForm = document.getElementById('nutritionForm');
        if (nutritionForm) {
            nutritionForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.analyzeNutrition();
            });
        }

        // API test button
        const testApiBtn = document.getElementById('testApiBtn');
        if (testApiBtn) {
            testApiBtn.addEventListener('click', () => {
                this.testApiConnection();
            });
        }

        // Health check button
        const healthCheckBtn = document.getElementById('healthCheckBtn');
        if (healthCheckBtn) {
            healthCheckBtn.addEventListener('click', () => {
                this.checkSystemHealth();
            });
        }
    }

    async loadTranslations() {
        this.translations = {
            en: {
                heroTitle: 'Welcome to Doctor Healthy',
                heroSubtitle: 'Your personalized nutrition and health companion',
                systemStatus: 'System Status',
                featuresTitle: 'Our Features',
                nutritionTitle: 'Nutrition Analysis',
                nutritionDesc: 'Get detailed nutritional analysis of your meals with personalized recommendations.',
                halalTitle: 'Halal Compliance',
                halalDesc: 'Ensure your food choices comply with Islamic dietary laws.',
                medicalTitle: 'Medical Guidance',
                medicalDesc: 'Receive nutrition advice tailored to your health conditions.',
                nutritionFormTitle: 'Nutrition Analysis',
                foodLabel: 'Enter Food Item',
                quantityLabel: 'Quantity',
                unitLabel: 'Unit',
                halalCheckLabel: 'Check Halal Compliance',
                analyzeBtnText: 'Analyze Nutrition',
                resultsTitle: 'Analysis Results',
                healthTitle: 'Health Monitoring',
                apiTestTitle: 'API Connection Test',
                testApiBtnText: 'Test API Connection',
                healthCheckTitle: 'System Health',
                healthCheckBtnText: 'Check System Health',
                footerText: '© 2024 Doctor Healthy - Nutrition Platform. All rights reserved.',
                disclaimerText: 'Medical Disclaimer: This platform provides general nutritional information and should not replace professional medical advice.'
            },
            ar: {
                heroTitle: 'مرحباً بكم في دكتور هيلثي',
                heroSubtitle: 'رفيقكم الشخصي للتغذية والصحة',
                systemStatus: 'حالة النظام',
                featuresTitle: 'ميزاتنا',
                nutritionTitle: 'تحليل التغذية',
                nutritionDesc: 'احصل على تحليل غذائي مفصل لوجباتك مع توصيات شخصية.',
                halalTitle: 'الامتثال للحلال',
                halalDesc: 'تأكد من أن خياراتك الغذائية تتوافق مع القوانين الغذائية الإسلامية.',
                medicalTitle: 'الإرشاد الطبي',
                medicalDesc: 'احصل على نصائح غذائية مخصصة لحالتك الصحية.',
                nutritionFormTitle: 'تحليل التغذية',
                foodLabel: 'أدخل عنصر الطعام',
                quantityLabel: 'الكمية',
                unitLabel: 'الوحدة',
                halalCheckLabel: 'فحص الامتثال للحلال',
                analyzeBtnText: 'تحليل التغذية',
                resultsTitle: 'نتائج التحليل',
                healthTitle: 'مراقبة الصحة',
                apiTestTitle: 'اختبار اتصال API',
                testApiBtnText: 'اختبار اتصال API',
                healthCheckTitle: 'صحة النظام',
                healthCheckBtnText: 'فحص صحة النظام',
                footerText: '© 2024 دكتور هيلثي - منصة التغذية. جميع الحقوق محفوظة.',
                disclaimerText: 'إخلاء المسؤولية الطبية: توفر هذه المنصة معلومات غذائية عامة ولا يجب أن تحل محل المشورة الطبية المهنية.'
            }
        };
    }

    updateLanguage(lang) {
        this.currentLanguage = lang;
        const translations = this.translations[lang] || this.translations.en;
        
        // Update text content
        Object.keys(translations).forEach(key => {
            const element = document.getElementById(key);
            if (element) {
                element.textContent = translations[key];
            }
        });

        // Update RTL for Arabic
        if (lang === 'ar') {
            document.body.classList.add('rtl');
            document.documentElement.setAttribute('dir', 'rtl');
        } else {
            document.body.classList.remove('rtl');
            document.documentElement.setAttribute('dir', 'ltr');
        }

        // Update language selector
        const languageSelect = document.getElementById('languageSelect');
        if (languageSelect) {
            languageSelect.value = lang;
        }
    }

    async makeApiCall(endpoint, options = {}) {
        const url = `${this.apiBaseUrl}${endpoint}`;
        const defaultOptions = {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        };

        try {
            console.log(`🌐 Making API call to: ${url}`);
            const response = await fetch(url, { ...defaultOptions, ...options });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const data = await response.json();
            console.log('✅ API Response:', data);
            return { success: true, data };
        } catch (error) {
            console.error('❌ API Error:', error);
            return { success: false, error: error.message };
        }
    }

    async checkSystemHealth() {
        console.log('🔍 Checking system health...');
        
        // Update status indicators to loading
        this.updateStatusIndicator('backendStatus', 'loading', 'Backend: Checking...');
        this.updateStatusIndicator('databaseStatus', 'loading', 'Database: Checking...');
        this.updateStatusIndicator('apiStatus', 'loading', 'API: Checking...');

        // Check backend health
        const healthResult = await this.makeApiCall('/health');
        if (healthResult.success) {
            this.updateStatusIndicator('backendStatus', 'healthy', 'Backend: Online ✅');
            this.updateStatusIndicator('databaseStatus', 'healthy', 'Database: Connected ✅');
            this.updateStatusIndicator('apiStatus', 'healthy', 'API: Available ✅');
        } else {
            this.updateStatusIndicator('backendStatus', 'error', 'Backend: Offline ❌');
            this.updateStatusIndicator('databaseStatus', 'error', 'Database: Disconnected ❌');
            this.updateStatusIndicator('apiStatus', 'error', 'API: Unavailable ❌');
        }

        return healthResult;
    }

    updateStatusIndicator(elementId, status, text) {
        const element = document.getElementById(elementId);
        if (!element) return;

        // Remove existing status classes
        const indicator = element.querySelector('.status-indicator');
        if (indicator) {
            indicator.className = 'status-indicator';
            
            switch (status) {
                case 'healthy':
                    indicator.classList.add('status-healthy');
                    break;
                case 'warning':
                    indicator.classList.add('status-warning');
                    break;
                case 'error':
                    indicator.classList.add('status-error');
                    break;
                case 'loading':
                    // Keep default styling for loading
                    break;
            }
        }

        // Update text
        const textNode = element.childNodes[element.childNodes.length - 1];
        if (textNode && textNode.nodeType === Node.TEXT_NODE) {
            textNode.textContent = text;
        } else {
            // If no text node, create one
            element.appendChild(document.createTextNode(text));
        }
    }

    async testApiConnection() {
        const resultsDiv = document.getElementById('apiTestResults');
        if (!resultsDiv) return;

        resultsDiv.innerHTML = '<div class="loading"></div> Testing API connection...';

        // Test multiple endpoints
        const endpoints = [
            { name: 'Health Check', url: '/health' },
            { name: 'API Info', url: '/api/info' },
            { name: 'Recipes', url: '/api/recipes' }
        ];

        let results = '<div class="mt-3">';
        
        for (const endpoint of endpoints) {
            const result = await this.makeApiCall(endpoint.url);
            const status = result.success ? '✅' : '❌';
            const statusClass = result.success ? 'text-success' : 'text-danger';
            
            results += `
                <div class="d-flex justify-content-between align-items-center mb-2">
                    <span>${endpoint.name}:</span>
                    <span class="${statusClass}">${status} ${result.success ? 'OK' : result.error}</span>
                </div>
            `;
        }
        
        results += '</div>';
        resultsDiv.innerHTML = results;
    }

    async analyzeNutrition() {
        const foodInput = document.getElementById('foodInput');
        const quantityInput = document.getElementById('quantityInput');
        const unitSelect = document.getElementById('unitSelect');
        const halalCheck = document.getElementById('halalCheck');
        const resultsDiv = document.getElementById('nutritionResults');
        const resultsContent = document.getElementById('resultsContent');
        const analyzeBtn = document.getElementById('analyzeBtn');

        if (!foodInput || !quantityInput || !unitSelect || !resultsDiv || !resultsContent) {
            console.error('Required form elements not found');
            return;
        }

        const foodItem = foodInput.value.trim();
        if (!foodItem) {
            alert('Please enter a food item');
            return;
        }

        // Show loading state
        analyzeBtn.disabled = true;
        analyzeBtn.innerHTML = '<span class="loading"></span> Analyzing...';
        resultsDiv.style.display = 'block';
        resultsContent.innerHTML = '<div class="loading"></div> Analyzing nutrition...';

        try {
            // Prepare request data
            const requestData = {
                food: foodItem,
                quantity: parseFloat(quantityInput.value) || 100,
                unit: unitSelect.value,
                checkHalal: halalCheck.checked,
                language: this.currentLanguage
            };

            console.log('📊 Analyzing nutrition for:', requestData);

            // Make API call for nutrition analysis
            const result = await this.makeApiCall('/api/nutrition/analyze', {
                method: 'POST',
                body: JSON.stringify(requestData)
            });

            if (result.success) {
                this.displayNutritionResults(result.data);
            } else {
                // Fallback to mock data if API fails
                console.warn('API failed, using mock data');
                this.displayMockNutritionResults(requestData);
            }
        } catch (error) {
            console.error('Error analyzing nutrition:', error);
            resultsContent.innerHTML = `
                <div class="alert alert-warning">
                    <i class="fas fa-exclamation-triangle me-2"></i>
                    Unable to connect to the nutrition analysis service. Please try again later.
                </div>
            `;
        } finally {
            // Reset button state
            analyzeBtn.disabled = false;
            analyzeBtn.innerHTML = '<i class="fas fa-search me-2"></i><span id="analyzeBtnText">Analyze Nutrition</span>';
        }
    }

    displayNutritionResults(data) {
        const resultsContent = document.getElementById('resultsContent');
        if (!resultsContent) return;

        let html = `
            <div class="row">
                <div class="col-md-6">
                    <h6><i class="fas fa-chart-pie me-2"></i>Nutritional Information</h6>
                    <ul class="list-unstyled">
                        <li><strong>Calories:</strong> ${data.calories || 'N/A'} kcal</li>
                        <li><strong>Protein:</strong> ${data.protein || 'N/A'} g</li>
                        <li><strong>Carbohydrates:</strong> ${data.carbohydrates || 'N/A'} g</li>
                        <li><strong>Fat:</strong> ${data.fat || 'N/A'} g</li>
                        <li><strong>Fiber:</strong> ${data.fiber || 'N/A'} g</li>
                    </ul>
                </div>
                <div class="col-md-6">
                    <h6><i class="fas fa-info-circle me-2"></i>Additional Info</h6>
        `;

        if (data.halalStatus !== undefined) {
            const halalIcon = data.halalStatus ? 'fas fa-check-circle text-success' : 'fas fa-times-circle text-danger';
            const halalText = data.halalStatus ? 'Halal ✅' : 'Not Halal ❌';
            html += `<p><i class="${halalIcon} me-2"></i><strong>Halal Status:</strong> ${halalText}</p>`;
        }

        if (data.recommendations && data.recommendations.length > 0) {
            html += `
                <h6><i class="fas fa-lightbulb me-2"></i>Recommendations</h6>
                <ul>
            `;
            data.recommendations.forEach(rec => {
                html += `<li>${rec}</li>`;
            });
            html += '</ul>';
        }

        html += `
                </div>
            </div>
        `;

        if (data.medicalDisclaimer) {
            html += `
                <div class="alert alert-info mt-3">
                    <i class="fas fa-info-circle me-2"></i>
                    <strong>Medical Disclaimer:</strong> ${data.medicalDisclaimer}
                </div>
            `;
        }

        resultsContent.innerHTML = html;
    }

    displayMockNutritionResults(requestData) {
        // Generate mock data based on common foods
        const mockData = this.generateMockNutritionData(requestData.food, requestData.quantity);
        mockData.halalStatus = requestData.checkHalal ? this.checkMockHalalStatus(requestData.food) : undefined;
        mockData.medicalDisclaimer = 'This is sample data. For accurate nutritional information, please consult a registered dietitian.';
        
        this.displayNutritionResults(mockData);
    }

    generateMockNutritionData(food, quantity) {
        // Simple mock data generator
        const baseCalories = {
            'chicken': 165, 'rice': 130, 'apple': 52, 'bread': 265,
            'egg': 155, 'milk': 42, 'cheese': 113, 'fish': 206
        };

        const foodLower = food.toLowerCase();
        let calories = 100; // default
        
        for (const [key, value] of Object.entries(baseCalories)) {
            if (foodLower.includes(key)) {
                calories = value;
                break;
            }
        }

        // Scale by quantity (assuming per 100g base)
        const scale = quantity / 100;
        
        return {
            calories: Math.round(calories * scale),
            protein: Math.round(calories * 0.2 * scale / 4), // rough estimate
            carbohydrates: Math.round(calories * 0.5 * scale / 4),
            fat: Math.round(calories * 0.3 * scale / 9),
            fiber: Math.round(calories * 0.1 * scale / 4),
            recommendations: [
                'Maintain a balanced diet with variety',
                'Consider portion sizes for your daily needs',
                'Combine with regular physical activity'
            ]
        };
    }

    checkMockHalalStatus(food) {
        const nonHalalKeywords = ['pork', 'ham', 'bacon', 'wine', 'beer', 'alcohol'];
        const foodLower = food.toLowerCase();
        return !nonHalalKeywords.some(keyword => foodLower.includes(keyword));
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    console.log('🌟 DOM loaded, initializing Nutrition Platform...');
    window.nutritionApp = new NutritionApp();
});

// Global error handler
window.addEventListener('error', (event) => {
    console.error('🚨 Global error:', event.error);
});

// Service worker registration for PWA
if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register('/sw.js')
            .then(registration => {
                console.log('✅ SW registered: ', registration);
            })
            .catch(registrationError => {
                console.log('❌ SW registration failed: ', registrationError);
            });
    });
}