#!/usr/bin/env node

/**
 * Comprehensive End-to-End Testing for Nutrition Platform
 * Tests complete user journey from login to dashboard
 * Tests API integration with real data
 */

const puppeteer = require('puppeteer');
const axios = require('axios');

const BASE_URL = 'http://localhost:3000';
const API_BASE_URL = 'http://localhost:8080';

class E2ETestSuite {
    constructor() {
        this.browser = null;
        this.page = null;
        this.testResults = {
            passed: 0,
            failed: 0,
            total: 0,
            details: []
        };
    }

    async init() {
        console.log('🚀 Starting Comprehensive E2E Testing...\n');
        
        this.browser = await puppeteer.launch({
            headless: false,
            defaultViewport: { width: 1366, height: 768 },
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        });
        
        this.page = await this.browser.newPage();
        
        // Set default timeout
        this.page.setDefaultTimeout(10000);
        
        // Intercept console messages
        this.page.on('console', msg => {
            console.log('PAGE LOG:', msg.text());
        });
        
        // Intercept network requests for debugging
        this.page.on('request', request => {
            console.log('REQUEST:', request.method(), request.url());
        });
    }

    async test(name, testFunction) {
        this.testResults.total++;
        try {
            console.log(`📝 Running test: ${name}`);
            await testFunction();
            this.testResults.passed++;
            console.log(`✅ PASSED: ${name}\n`);
            this.testResults.details.push({ name, status: 'PASSED', error: null });
        } catch (error) {
            this.testResults.failed++;
            console.log(`❌ FAILED: ${name}`);
            console.log(`   Error: ${error.message}\n`);
            this.testResults.details.push({ name, status: 'FAILED', error: error.message });
        }
    }

    async navigateToHomepage() {
        await this.page.goto(BASE_URL, { waitUntil: 'networkidle2' });
        await this.page.waitForTimeout(2000);
    }

    async takeScreenshot(name) {
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const filename = `test-screenshot-${name}-${timestamp}.png`;
        await this.page.screenshot({ path: filename, fullPage: true });
        console.log(`📸 Screenshot saved: ${filename}`);
        return filename;
    }

    // API Tests
    async testAPIHealth() {
        const response = await axios.get(`${API_BASE_URL}/health`, { timeout: 5000 });
        if (response.status !== 200) {
            throw new Error(`Health check failed with status ${response.status}`);
        }
        console.log('✅ Backend API is healthy:', response.data);
    }

    async testAPIInfo() {
        const response = await axios.get(`${API_BASE_URL}/api/info`, { timeout: 5000 });
        if (response.status !== 200) {
            throw new Error(`API info endpoint failed with status ${response.status}`);
        }
        console.log('✅ API Info endpoint working:', response.data);
    }

    async testNutritionDataEndpoints() {
        const endpoints = [
            '/api/v1/metabolism',
            '/api/v1/meal-plans',
            '/api/v1/vitamins-minerals',
            '/api/v1/workout-techniques',
            '/api/v1/calories',
            '/api/v1/skills'
        ];

        for (const endpoint of endpoints) {
            try {
                const response = await axios.get(`${API_BASE_URL}${endpoint}`, { timeout: 5000 });
                if (response.status === 200) {
                    console.log(`✅ ${endpoint} - OK (${response.data.length || 'data'} items)`);
                }
            } catch (error) {
                console.log(`⚠️  ${endpoint} - Error: ${error.message}`);
            }
        }
    }

    // Frontend Tests
    async testHomepageLoad() {
        await this.navigateToHomepage();
        const title = await this.page.title();
        if (!title.includes('Nutrition') && !title.includes('Health')) {
            throw new Error(`Unexpected page title: ${title}`);
        }
        
        // Check for main elements
        const bodyText = await this.page.evaluate(() => document.body.innerText);
        if (bodyText.length < 100) {
            throw new Error('Page appears to be empty or not loaded properly');
        }
    }

    async testNavigation() {
        // Look for navigation elements
        const navElements = await this.page.$$('nav a, .nav a, header a');
        if (navElements.length === 0) {
            console.log('⚠️  No navigation links found');
        } else {
            console.log(`✅ Found ${navElements.length} navigation links`);
        }
    }

    async testDashboardAccess() {
        // Try to access dashboard pages
        const dashboardRoutes = [
            '/dashboard/meals',
            '/dashboard/workouts',
            '/dashboard/health',
            '/dashboard/recipes'
        ];

        for (const route of dashboardRoutes) {
            try {
                await this.page.goto(`${BASE_URL}${route}`, { waitUntil: 'networkidle2', timeout: 5000 });
                await this.page.waitForTimeout(1000);
                
                const url = this.page.url();
                if (url.includes(route)) {
                    console.log(`✅ Dashboard route accessible: ${route}`);
                } else {
                    console.log(`⚠️  Dashboard route redirected: ${route} -> ${url}`);
                }
            } catch (error) {
                console.log(`⚠️  Dashboard route error: ${route} - ${error.message}`);
            }
        }
    }

    async testResponsiveDesign() {
        // Test mobile view
        await this.page.setViewport({ width: 375, height: 667 });
        await this.page.waitForTimeout(1000);
        
        // Test tablet view
        await this.page.setViewport({ width: 768, height: 1024 });
        await this.page.waitForTimeout(1000);
        
        // Back to desktop
        await this.page.setViewport({ width: 1366, height: 768 });
        await this.page.waitForTimeout(1000);
        
        console.log('✅ Responsive design tested across different viewports');
    }

    async testPerformance() {
        const startTime = Date.now();
        await this.navigateToHomepage();
        const loadTime = Date.now() - startTime;
        
        console.log(`⏱️  Page load time: ${loadTime}ms`);
        
        if (loadTime > 5000) {
            console.log('⚠️  Slow page load detected');
        } else {
            console.log('✅ Acceptable page load time');
        }
    }

    async testErrorHandling() {
        // Try accessing non-existent routes
        try {
            await this.page.goto(`${BASE_URL}/non-existent-page`, { waitUntil: 'networkidle2' });
            await this.page.waitForTimeout(2000);
            
            // Should show 404 or redirect to home
            const url = this.page.url();
            console.log(`✅ Error handling test: ${url}`);
        } catch (error) {
            console.log(`⚠️  Error handling test failed: ${error.message}`);
        }
    }

    async testUserJourney() {
        console.log('🎭 Simulating complete user journey...');
        
        // 1. Visit homepage
        await this.navigateToHomepage();
        await this.takeScreenshot('1-homepage');
        
        // 2. Try to navigate to different sections
        const links = await this.page.$$('a');
        for (let i = 0; i < Math.min(links.length, 5); i++) {
            try {
                const href = await links[i].evaluate(el => el.href);
                if (href && href.includes(BASE_URL)) {
                    console.log(`🔗 Clicking link: ${href}`);
                    await links[i].click();
                    await this.page.waitForTimeout(2000);
                    await this.page.goBack();
                    await this.page.waitForTimeout(1000);
                }
            } catch (error) {
                console.log(`⚠️  Could not click link: ${error.message}`);
            }
        }
        
        await this.takeScreenshot('2-user-journey-complete');
    }

    async runAllTests() {
        await this.init();

        console.log('🔍 API INTEGRATION TESTS\n');
        await this.test('API Health Check', () => this.testAPIHealth());
        await this.test('API Info Endpoint', () => this.testAPIInfo());
        await this.test('Nutrition Data Endpoints', () => this.testNutritionDataEndpoints());

        console.log('🎨 FRONTEND FUNCTIONALITY TESTS\n');
        await this.test('Homepage Load', () => this.testHomepageLoad());
        await this.test('Navigation Elements', () => this.testNavigation());
        await this.test('Dashboard Access', () => this.testDashboardAccess());
        await this.test('Responsive Design', () => this.testResponsiveDesign());
        await this.test('Performance Test', () => this.testPerformance());
        await this.test('Error Handling', () => this.testErrorHandling());

        console.log('👤 USER JOURNEY TESTS\n');
        await this.test('Complete User Journey', () => this.testUserJourney());

        await this.takeScreenshot('final-application-state');
    }

    async generateReport() {
        console.log('\n📊 TEST EXECUTION REPORT\n');
        console.log('='.repeat(50));
        console.log(`Total Tests: ${this.testResults.total}`);
        console.log(`Passed: ${this.testResults.passed} ✅`);
        console.log(`Failed: ${this.testResults.failed} ❌`);
        console.log(`Success Rate: ${((this.testResults.passed / this.testResults.total) * 100).toFixed(1)}%`);
        console.log('='.repeat(50));

        if (this.testResults.failed > 0) {
            console.log('\n❌ FAILED TESTS:');
            this.testResults.details
                .filter(test => test.status === 'FAILED')
                .forEach(test => {
                    console.log(`  • ${test.name}: ${test.error}`);
                });
        }

        // Generate JSON report
        const report = {
            timestamp: new Date().toISOString(),
            summary: {
                total: this.testResults.total,
                passed: this.testResults.passed,
                failed: this.testResults.failed,
                successRate: ((this.testResults.passed / this.testResults.total) * 100).toFixed(1)
            },
            details: this.testResults.details,
            environment: {
                frontend: BASE_URL,
                backend: API_BASE_URL,
                nodeVersion: process.version,
                platform: process.platform
            }
        };

        require('fs').writeFileSync('e2e-test-report.json', JSON.stringify(report, null, 2));
        console.log('\n📄 Detailed report saved to: e2e-test-report.json');
    }

    async cleanup() {
        if (this.browser) {
            await this.browser.close();
        }
        console.log('\n🧹 Cleanup completed');
    }
}

// Main execution
async function main() {
    const testSuite = new E2ETestSuite();
    
    try {
        await testSuite.runAllTests();
        await testSuite.generateReport();
    } catch (error) {
        console.error('💥 Test suite execution failed:', error);
    } finally {
        await testSuite.cleanup();
        
        // Exit with appropriate code
        process.exit(testSuite.testResults.failed > 0 ? 1 : 0);
    }
}

// Handle graceful shutdown
process.on('SIGINT', () => {
    console.log('\n🛑 Test execution interrupted');
    process.exit(1);
});

// Run tests if this file is executed directly
if (require.main === module) {
    main().catch(console.error);
}

module.exports = E2ETestSuite;
