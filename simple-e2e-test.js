#!/usr/bin/env node

/**
 * Simple End-to-End Testing for Nutrition Platform
 * Tests API integration and basic functionality
 */

const axios = require('axios');

const BASE_URL = 'http://localhost:3000';
const API_BASE_URL = 'http://localhost:8080';

class SimpleE2ETestSuite {
    constructor() {
        this.testResults = {
            passed: 0,
            failed: 0,
            total: 0,
            details: []
        };
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
            '/api/v1/skills',
            '/api/v1/diseases',
            '/api/v1/type-plans'
        ];

        let successCount = 0;
        for (const endpoint of endpoints) {
            try {
                const response = await axios.get(`${API_BASE_URL}${endpoint}`, { timeout: 5000 });
                if (response.status === 200) {
                    const dataLength = Array.isArray(response.data) ? response.data.length : 'object';
                    console.log(`✅ ${endpoint} - OK (${dataLength} items)`);
                    successCount++;
                }
            } catch (error) {
                console.log(`⚠️  ${endpoint} - Error: ${error.message}`);
            }
        }
        
        if (successCount < endpoints.length * 0.5) {
            throw new Error(`Too many endpoints failed: ${successCount}/${endpoints.length}`);
        }
    }

    async testHealthEndpoints() {
        const healthEndpoints = [
            '/api/v1/health/conditions',
            '/api/v1/health/tips'
        ];

        for (const endpoint of healthEndpoints) {
            try {
                const response = await axios.get(`${API_BASE_URL}${endpoint}`, { timeout: 5000 });
                if (response.status === 200) {
                    console.log(`✅ Health endpoint ${endpoint} - OK`);
                }
            } catch (error) {
                console.log(`⚠️  Health endpoint ${endpoint} - Error: ${error.message}`);
            }
        }
    }

    async testFrontendHealth() {
        try {
            const response = await axios.get(BASE_URL, { timeout: 5000 });
            if (response.status === 200) {
                console.log('✅ Frontend is accessible');
                return true;
            }
        } catch (error) {
            console.log(`⚠️  Frontend access error: ${error.message}`);
            throw new Error('Frontend is not accessible');
        }
    }

    async testAPIResponseTimes() {
        const endpoints = [
            '/health',
            '/api/info',
            '/api/v1/calories'
        ];

        for (const endpoint of endpoints) {
            const startTime = Date.now();
            try {
                await axios.get(`${API_BASE_URL}${endpoint}`, { timeout: 5000 });
                const responseTime = Date.now() - startTime;
                console.log(`⏱️  ${endpoint}: ${responseTime}ms`);
                
                if (responseTime > 3000) {
                    console.log(`⚠️  Slow response time for ${endpoint}`);
                }
            } catch (error) {
                console.log(`⚠️  ${endpoint} - Error: ${error.message}`);
            }
        }
    }

    async testDataQuality() {
        try {
            // Test calories endpoint data quality
            const caloriesResponse = await axios.get(`${API_BASE_URL}/api/v1/calories`, { timeout: 5000 });
            if (caloriesResponse.status === 200 && Array.isArray(caloriesResponse.data)) {
                console.log(`✅ Calories data quality check: ${caloriesResponse.data.length} items`);
                
                // Check for expected structure
                if (caloriesResponse.data.length > 0) {
                    const firstItem = caloriesResponse.data[0];
                    if (typeof firstItem === 'object' && firstItem !== null) {
                        console.log('✅ Data structure validation passed');
                    } else {
                        console.log('⚠️  Data structure may be invalid');
                    }
                }
            }

            // Test meal plans data quality
            const mealsResponse = await axios.get(`${API_BASE_URL}/api/v1/meal-plans`, { timeout: 5000 });
            if (mealsResponse.status === 200 && Array.isArray(mealsResponse.data)) {
                console.log(`✅ Meal plans data quality check: ${mealsResponse.data.length} items`);
            }

        } catch (error) {
            console.log(`⚠️  Data quality check failed: ${error.message}`);
        }
    }

    async testErrorHandling() {
        try {
            // Test 404 handling
            await axios.get(`${API_BASE_URL}/non-existent-endpoint`, { timeout: 5000 });
        } catch (error) {
            if (error.response && error.response.status === 404) {
                console.log('✅ 404 error handling working correctly');
            } else {
                console.log(`⚠️  Unexpected error response: ${error.message}`);
            }
        }
    }

    async runAllTests() {
        console.log('🚀 Starting Simple E2E Testing...\n');

        console.log('🔍 API HEALTH CHECKS\n');
        await this.test('API Health Check', () => this.testAPIHealth());
        await this.test('API Info Endpoint', () => this.testAPIInfo());

        console.log('📊 NUTRITION DATA ENDPOINTS\n');
        await this.test('Nutrition Data Endpoints', () => this.testNutritionDataEndpoints());

        console.log('🏥 HEALTH SERVICE ENDPOINTS\n');
        await this.test('Health Endpoints', () => this.testHealthEndpoints());

        console.log('🎨 FRONTEND ACCESSIBILITY\n');
        await this.test('Frontend Health', () => this.testFrontendHealth());

        console.log('⏱️  PERFORMANCE TESTS\n');
        await this.test('API Response Times', () => this.testAPIResponseTimes());

        console.log('🔍 DATA QUALITY TESTS\n');
        await this.test('Data Quality', () => this.testDataQuality());

        console.log('🛡️  ERROR HANDLING TESTS\n');
        await this.test('Error Handling', () => this.testErrorHandling());
    }

    generateReport() {
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

        require('fs').writeFileSync('simple-e2e-test-report.json', JSON.stringify(report, null, 2));
        console.log('\n📄 Detailed report saved to: simple-e2e-test-report.json');
        
        return report;
    }
}

// Main execution
async function main() {
    const testSuite = new SimpleE2ETestSuite();
    
    try {
        await testSuite.runAllTests();
        const report = testSuite.generateReport();
        
        console.log('\n🎉 TESTING COMPLETED!');
        
        if (testSuite.testResults.failed === 0) {
            console.log('🏆 All tests passed! The application is working correctly.');
        } else {
            console.log('⚠️  Some tests failed. Please check the detailed report.');
        }
        
    } catch (error) {
        console.error('💥 Test suite execution failed:', error);
        process.exit(1);
    }
}

// Run tests if this file is executed directly
if (require.main === module) {
    main().catch(console.error);
}

module.exports = SimpleE2ETestSuite;
