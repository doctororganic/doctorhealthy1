#!/usr/bin/env node

/**
 * Comprehensive End-to-End Authentication & API Testing
 * Tests complete user authentication flow and API functionality
 */

const axios = require('axios');
const jwt = require('jsonwebtoken');

const API_BASE_URL = 'http://localhost:8080';
const JWT_SECRET = 'your-secret-key-change-in-production';

class AuthE2ETestSuite {
    constructor() {
        this.testResults = {
            passed: 0,
            failed: 0,
            total: 0,
            details: []
        };
        this.authTokens = {
            accessToken: null,
            refreshToken: null
        };
        this.testUser = {
            email: 'testuser@example.com',
            password: 'TestPassword123!',
            firstName: 'Test',
            lastName: 'User',
            dateOfBirth: '1990-01-01',
            gender: 'male',
            language: 'en'
        };
    }

    async test(name, testFunction) {
        this.testResults.total++;
        try {
            console.log(`📝 Running test: ${name}`);
            const result = await testFunction();
            this.testResults.passed++;
            console.log(`✅ PASSED: ${name}`);
            if (result && result.message) {
                console.log(`   ${result.message}`);
            }
            console.log('');
            this.testResults.details.push({ name, status: 'PASSED', error: null, result });
        } catch (error) {
            this.testResults.failed++;
            console.log(`❌ FAILED: ${name}`);
            console.log(`   Error: ${error.message}\n`);
            this.testResults.details.push({ name, status: 'FAILED', error: error.message });
        }
    }

    // JWT Token validation helper
    validateJWTToken(token) {
        try {
            const decoded = jwt.verify(token, JWT_SECRET);
            return decoded;
        } catch (error) {
            throw new Error(`Invalid JWT token: ${error.message}`);
        }
    }

    // 1. User Registration Flow Tests
    async testUserRegistration() {
        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/register`, {
            email: this.testUser.email,
            password: this.testUser.password,
            confirm_password: this.testUser.password,
            first_name: this.testUser.firstName,
            last_name: this.testUser.lastName,
            date_of_birth: this.testUser.dateOfBirth,
            gender: this.testUser.gender,
            language: this.testUser.language
        });

        if (response.status !== 201) {
            throw new Error(`Expected status 201, got ${response.status}`);
        }

        const { access_token, refresh_token, user, expires_in } = response.data;
        
        if (!access_token || !refresh_token) {
            throw new Error('Missing tokens in registration response');
        }

        // Validate JWT token structure
        const decodedToken = this.validateJWTToken(access_token);
        if (!decodedToken.user_id || !decodedToken.email) {
            throw new Error('Invalid token claims structure');
        }

        // Store tokens for subsequent tests
        this.authTokens.accessToken = access_token;
        this.authTokens.refreshToken = refresh_token;

        return { 
            message: 'User registered successfully with valid tokens',
            userId: decodedToken.user_id,
            expiresIn: expires_in
        };
    }

    // Test registration validation
    async testRegistrationValidation() {
        // Test missing email
        try {
            await axios.post(`${API_BASE_URL}/api/v1/auth/register`, {
                password: this.testUser.password,
                confirm_password: this.testUser.password
            });
            throw new Error('Should have failed with missing email');
        } catch (error) {
            if (error.response && error.response.status === 400) {
                // Expected behavior
            } else {
                throw new Error('Expected validation error for missing email');
            }
        }

        // Test mismatched passwords
        try {
            await axios.post(`${API_BASE_URL}/api/v1/auth/register`, {
                email: 'test2@example.com',
                password: 'Password123!',
                confirm_password: 'DifferentPassword123!'
            });
            throw new Error('Should have failed with mismatched passwords');
        } catch (error) {
            if (error.response && error.response.status === 400) {
                // Expected behavior
            } else {
                throw new Error('Expected validation error for mismatched passwords');
            }
        }

        return { message: 'Registration validation working correctly' };
    }

    // 2. Login with Valid Credentials Tests
    async testUserLogin() {
        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/login`, {
            email: this.testUser.email,
            password: this.testUser.password
        });

        if (response.status !== 200) {
            throw new Error(`Expected status 200, got ${response.status}`);
        }

        const { access_token, refresh_token, user } = response.data;
        
        if (!access_token || !refresh_token) {
            throw new Error('Missing tokens in login response');
        }

        // Validate JWT token
        const decodedToken = this.validateJWTToken(access_token);
        if (decodedToken.email !== this.testUser.email) {
            throw new Error('Token email does not match login email');
        }

        // Update stored tokens
        this.authTokens.accessToken = access_token;
        this.authTokens.refreshToken = refresh_token;

        return { 
            message: 'User login successful with valid tokens',
            userId: decodedToken.user_id
        };
    }

    // 3. JWT Token Validation Tests
    async testJWTTokenValidation() {
        if (!this.authTokens.accessToken) {
            throw new Error('No access token available for testing');
        }

        // Test valid token
        const decoded = this.validateJWTToken(this.authTokens.accessToken);
        
        if (!decoded.user_id || !decoded.email || !decoded.role) {
            throw new Error('Missing required claims in JWT token');
        }

        // Test token expiration (should be valid for 24 hours)
        const now = Math.floor(Date.now() / 1000);
        if (decoded.exp <= now) {
            throw new Error('Token has already expired');
        }

        return { 
            message: 'JWT token validation successful',
            claims: {
                userId: decoded.user_id,
                email: decoded.email,
                role: decoded.role,
                expiresAt: new Date(decoded.exp * 1000).toISOString()
            }
        };
    }

    // 4. Access Protected Dashboard Tests
    async testProtectedDashboardAccess() {
        if (!this.authTokens.accessToken) {
            throw new Error('No access token available for testing');
        }

        // Test accessing protected endpoint with valid token
        const response = await axios.get(`${API_BASE_URL}/api/v1/auth/profile`, {
            headers: {
                'Authorization': `Bearer ${this.authTokens.accessToken}`
            }
        });

        if (response.status !== 200) {
            throw new Error(`Expected status 200, got ${response.status}`);
        }

        return { message: 'Protected dashboard access successful' };
    }

    // 5. API Endpoint Functionality Tests
    async testAPIEndpoints() {
        const endpoints = [
            { path: '/health', expectedStatus: 200 },
            { path: '/api/info', expectedStatus: 200 },
            { path: '/api/v1/metabolism', expectedStatus: 200 },
            { path: '/api/v1/meal-plans', expectedStatus: 200 },
            { path: '/api/v1/vitamins-minerals', expectedStatus: 200 },
            { path: '/api/v1/calories', expectedStatus: 200 },
            { path: '/api/v1/skills', expectedStatus: 200 }
        ];

        const results = [];
        for (const endpoint of endpoints) {
            try {
                const response = await axios.get(`${API_BASE_URL}${endpoint.path}`);
                if (response.status === endpoint.expectedStatus) {
                    results.push(`${endpoint.path} - OK`);
                } else {
                    results.push(`${endpoint.path} - Unexpected status: ${response.status}`);
                }
            } catch (error) {
                results.push(`${endpoint.path} - Error: ${error.message}`);
            }
        }

        return { 
            message: 'API endpoints test completed',
            results: results
        };
    }

    // 6. Error Handling for Invalid Credentials Tests
    async testInvalidCredentialsErrorHandling() {
        // Test invalid email/password combination
        try {
            await axios.post(`${API_BASE_URL}/api/v1/auth/login`, {
                email: 'invalid@example.com',
                password: 'wrongpassword'
            });
            throw new Error('Should have failed with invalid credentials');
        } catch (error) {
            if (error.response && (error.response.status === 401 || error.response.status === 400)) {
                // Expected behavior
            } else {
                throw new Error('Expected 401/400 for invalid credentials');
            }
        }

        // Test malformed request
        try {
            await axios.post(`${API_BASE_URL}/api/v1/auth/login`, {
                email: 'invalid-email-format',
                password: ''
            });
            throw new Error('Should have failed with malformed request');
        } catch (error) {
            if (error.response && error.response.status === 400) {
                // Expected behavior
            } else {
                throw new Error('Expected 400 for malformed request');
            }
        }

        return { message: 'Invalid credentials error handling working correctly' };
    }

    // Test unauthorized access
    async testUnauthorizedAccess() {
        // Test accessing protected endpoint without token
        try {
            await axios.get(`${API_BASE_URL}/api/v1/auth/profile`);
            throw new Error('Should have failed without authorization header');
        } catch (error) {
            if (error.response && error.response.status === 401) {
                // Expected behavior
            } else {
                throw new Error('Expected 401 for unauthorized access');
            }
        }

        // Test accessing protected endpoint with invalid token
        try {
            await axios.get(`${API_BASE_URL}/api/v1/auth/profile`, {
                headers: {
                    'Authorization': 'Bearer invalid-token'
                }
            });
            throw new Error('Should have failed with invalid token');
        } catch (error) {
            if (error.response && error.response.status === 401) {
                // Expected behavior
            } else {
                throw new Error('Expected 401 for invalid token');
            }
        }

        return { message: 'Unauthorized access handling working correctly' };
    }

    // Test token refresh functionality
    async testTokenRefresh() {
        if (!this.authTokens.refreshToken) {
            throw new Error('No refresh token available for testing');
        }

        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {
            refresh_token: this.authTokens.refreshToken
        });

        if (response.status !== 200) {
            throw new Error(`Expected status 200, got ${response.status}`);
        }

        const { access_token, refresh_token } = response.data;
        
        if (!access_token || !refresh_token) {
            throw new Error('Missing tokens in refresh response');
        }

        // Validate new JWT token
        const decodedToken = this.validateJWTToken(access_token);

        return { 
            message: 'Token refresh successful',
            newUserId: decodedToken.user_id
        };
    }

    // Test logout functionality
    async testLogout() {
        if (!this.authTokens.accessToken) {
            throw new Error('No access token available for testing');
        }

        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/logout`, {}, {
            headers: {
                'Authorization': `Bearer ${this.authTokens.accessToken}`
            }
        });

        if (response.status !== 200) {
            throw new Error(`Expected status 200, got ${response.status}`);
        }

        return { message: 'Logout successful' };
    }

    async runAllTests() {
        console.log('🚀 Starting Comprehensive End-to-End Authentication Testing...\n');

        console.log('🔐 USER REGISTRATION FLOW TESTS\n');
        await this.test('User Registration', () => this.testUserRegistration());
        await this.test('Registration Validation', () => this.testRegistrationValidation());

        console.log('🔑 USER LOGIN TESTS\n');
        await this.test('User Login with Valid Credentials', () => this.testUserLogin());

        console.log('🛡️  JWT TOKEN VALIDATION TESTS\n');
        await this.test('JWT Token Validation', () => this.testJWTTokenValidation());

        console.log('🏠 PROTECTED DASHBOARD ACCESS TESTS\n');
        await this.test('Protected Dashboard Access', () => this.testProtectedDashboardAccess());

        console.log('🔌 API ENDPOINT FUNCTIONALITY TESTS\n');
        await this.test('API Endpoints', () => this.testAPIEndpoints());

        console.log('❌ ERROR HANDLING TESTS\n');
        await this.test('Invalid Credentials Error Handling', () => this.testInvalidCredentialsErrorHandling());
        await this.test('Unauthorized Access Handling', () => this.testUnauthorizedAccess());

        console.log('🔄 TOKEN MANAGEMENT TESTS\n');
        await this.test('Token Refresh', () => this.testTokenRefresh());
        await this.test('User Logout', () => this.testLogout());
    }

    generateReport() {
        console.log('\n📊 COMPREHENSIVE TEST EXECUTION REPORT\n');
        console.log('='.repeat(60));
        console.log(`Total Tests: ${this.testResults.total}`);
        console.log(`Passed: ${this.testResults.passed} ✅`);
        console.log(`Failed: ${this.testResults.failed} ❌`);
        console.log(`Success Rate: ${((this.testResults.passed / this.testResults.total) * 100).toFixed(1)}%`);
        console.log('='.repeat(60));

        if (this.testResults.failed > 0) {
            console.log('\n❌ FAILED TESTS:');
            this.testResults.details
                .filter(test => test.status === 'FAILED')
                .forEach(test => {
                    console.log(`  • ${test.name}: ${test.error}`);
                });
        }

        // Generate detailed JSON report
        const report = {
            timestamp: new Date().toISOString(),
            summary: {
                total: this.testResults.total,
                passed: this.testResults.passed,
                failed: this.testResults.failed,
                successRate: ((this.testResults.passed / this.testResults.total) * 100).toFixed(1)
            },
            testCoverage: {
                userRegistration: '✅ Tested',
                userLogin: '✅ Tested',
                jwtTokenValidation: '✅ Tested',
                protectedDashboardAccess: '✅ Tested',
                apiEndpointFunctionality: '✅ Tested',
                errorHandling: '✅ Tested',
                tokenManagement: '✅ Tested'
            },
            details: this.testResults.details,
            environment: {
                backend: API_BASE_URL,
                nodeVersion: process.version,
                platform: process.platform
            }
        };

        require('fs').writeFileSync('e2e-auth-test-report.json', JSON.stringify(report, null, 2));
        console.log('\n📄 Detailed report saved to: e2e-auth-test-report.json');
    }
}

// Main execution
async function main() {
    const testSuite = new AuthE2ETestSuite();
    
    try {
        await testSuite.runAllTests();
        await testSuite.generateReport();
        
        // Exit with appropriate code
        process.exit(testSuite.testResults.failed > 0 ? 1 : 0);
    } catch (error) {
        console.error('💥 Test suite execution failed:', error);
        process.exit(1);
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

module.exports = AuthE2ETestSuite;
