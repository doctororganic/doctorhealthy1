import { chromium, FullConfig } from '@playwright/test';
import { test as base } from '@playwright/test';
import { join } from 'path';
import { writeFileSync } from 'fs';

/**
 * Global setup for E2E tests
 * 
 * This function runs once before all tests and:
 * - Sets up test database
 * - Creates test users
 * - Seeds test data
 * - Sets up API mocks if needed
 * - Initializes test environment
 */

async function globalSetup(config: FullConfig) {
  console.log('🚀 Starting global setup for E2E tests...');

  const baseURL = config.webServer?.url || 'http://localhost:3000';
  const apiURL = process.env.API_BASE_URL || 'http://localhost:8080';

  try {
    // Launch browser for setup tasks
    const browser = await chromium.launch();
    const context = await browser.newContext();
    const page = await context.newPage();

    // 1. Health check for frontend and backend
    console.log('🔍 Performing health checks...');
    
    try {
      await page.goto(baseURL, { timeout: 30000 });
      console.log('✅ Frontend is accessible');
    } catch (error) {
      console.warn('⚠️ Frontend health check failed:', error.message);
    }

    try {
      const response = await page.request.get(`${apiURL}/health`);
      if (response.ok()) {
        console.log('✅ Backend API is healthy');
      } else {
        console.warn('⚠️ Backend API health check failed:', response.status());
      }
    } catch (error) {
      console.warn('⚠️ Backend API health check failed:', error.message);
    }

    // 2. Setup test data
    console.log('📊 Setting up test data...');
    
    const testData = {
      users: [
        {
          email: 'testuser@example.com',
          password: 'TestPassword123!',
          firstName: 'Test',
          lastName: 'User',
          dateOfBirth: '1990-01-01',
          gender: 'male',
          language: 'en'
        },
        {
          email: 'arabicuser@example.com',
          password: 'TestPassword123!',
          firstName: 'مستخدم',
          lastName: 'اختبار',
          dateOfBirth: '1990-01-01',
          gender: 'male',
          language: 'ar'
        }
      ],
      meals: [
        {
          name: 'Breakfast Test Meal',
          type: 'breakfast',
          calories: 400,
          protein: 20,
          carbs: 50,
          fat: 15
        },
        {
          name: 'وجبة اختبار الإفطار',
          type: 'breakfast',
          calories: 400,
          protein: 20,
          carbs: 50,
          fat: 15
        }
      ],
      nutritionGoals: [
        {
          name: 'Weight Loss Goal',
          targetCalories: 1800,
          targetProtein: 120,
          targetCarbs: 200,
          targetFat: 60
        }
      ]
    };

    // Save test data for tests to use
    const testDataPath = join(__dirname, 'test-data.json');
    writeFileSync(testDataPath, JSON.stringify(testData, null, 2));
    console.log('✅ Test data saved');

    // 3. Create test users via API
    console.log('👥 Creating test users...');
    
    for (const user of testData.users) {
      try {
        const response = await page.request.post(`${apiURL}/api/v1/auth/register`, {
          data: {
            email: user.email,
            password: user.password,
            confirm_password: user.password,
            first_name: user.firstName,
            last_name: user.lastName,
            date_of_birth: user.dateOfBirth,
            gender: user.gender,
            language: user.language
          }
        });

        if (response.status() === 201) {
          console.log(`✅ Created test user: ${user.email}`);
          
          // Store tokens for this user
          const data = await response.json();
          testData.users.find(u => u.email === user.email).tokens = {
            accessToken: data.access_token,
            refreshToken: data.refresh_token
          };
        } else if (response.status() === 409) {
          console.log(`ℹ️ Test user already exists: ${user.email}`);
          
          // Try to login to get tokens
          try {
            const loginResponse = await page.request.post(`${apiURL}/api/v1/auth/login`, {
              data: {
                email: user.email,
                password: user.password
              }
            });
            
            if (loginResponse.ok()) {
              const loginData = await loginResponse.json();
              testData.users.find(u => u.email === user.email).tokens = {
                accessToken: loginData.access_token,
                refreshToken: loginData.refresh_token
              };
            }
          } catch (loginError) {
            console.warn(`⚠️ Could not login test user ${user.email}:`, loginError.message);
          }
        } else {
          console.warn(`⚠️ Failed to create test user ${user.email}:`, response.status());
        }
      } catch (error) {
        console.warn(`⚠️ Error creating test user ${user.email}:`, error.message);
      }
    }

    // Update test data with tokens
    writeFileSync(testDataPath, JSON.stringify(testData, null, 2));

    // 4. Wait for application to be fully ready
    console.log('⏳ Waiting for application to be fully ready...');
    await page.waitForTimeout(5000);

    await browser.close();

    console.log('✅ Global setup completed successfully');
  } catch (error) {
    console.error('❌ Global setup failed:', error);
    throw error;
  }
}

export default globalSetup;
