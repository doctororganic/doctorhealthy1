import { test as base, Page, expect, devices } from '@playwright/test';
import { readFileSync } from 'fs';
import { join } from 'path';

// Test data interface
interface TestData {
  users: Array<{
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    dateOfBirth: string;
    gender: string;
    language: string;
    tokens?: {
      accessToken: string;
      refreshToken: string;
    };
  }>;
  meals: Array<{
    name: string;
    type: string;
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
  }>;
  nutritionGoals: Array<{
    name: string;
    targetCalories: number;
    targetProtein: number;
    targetCarbs: number;
    targetFat: number;
  }>;
}

// Extended test fixtures
export const test = base.extend<{
  page: Page;
  testData: TestData;
  authenticatedPage: Page;
  arabicPage: Page;
  mobilePage: Page;
}>({
  // Load test data
  testData: async ({}, use) => {
    try {
      const testDataPath = join(__dirname, '../../test-data.json');
      const testData = JSON.parse(readFileSync(testDataPath, 'utf-8'));
      await use(testData);
    } catch (error) {
      console.warn('⚠️ Could not load test data, using defaults');
      const defaultData: TestData = {
        users: [
          {
            email: 'testuser@example.com',
            password: 'TestPassword123!',
            firstName: 'Test',
            lastName: 'User',
            dateOfBirth: '1990-01-01',
            gender: 'male',
            language: 'en'
          }
        ],
        meals: [],
        nutritionGoals: []
      };
      await use(defaultData);
    }
  },

  // Authenticated page fixture
  authenticatedPage: async ({ page, testData, context }, use) => {
    const testUser = testData.users.find(u => u.email === 'testuser@example.com');
    
    if (testUser?.tokens?.accessToken) {
      // Set auth token in localStorage
      await page.addInitScript(() => {
        window.localStorage.setItem('accessToken', 'TEST_TOKEN');
        window.localStorage.setItem('refreshToken', 'TEST_REFRESH_TOKEN');
      });
      
      // Add auth headers
      await page.route('**/*', route => {
        const headers = route.request().headers();
        headers['authorization'] = `Bearer ${testUser.tokens.accessToken}`;
        route.continue({ headers });
      });
    }
    
    await use(page);
  },

  // Arabic page fixture
  arabicPage: async ({ page, context }, use) => {
    const arabicContext = await context.browser()?.newContext({
      locale: 'ar-SA',
      extraHTTPHeaders: {
        'Accept-Language': 'ar-SA',
      }
    });
    
    const arabicPage = await arabicContext?.newPage();
    if (arabicPage) {
      await use(arabicPage);
      await arabicPage.close();
      await arabicContext?.close();
    } else {
      throw new Error('Could not create Arabic page context');
    }
  },

  // Mobile page fixture
  mobilePage: async ({ page, context }, use) => {
    const mobileContext = await context.browser()?.newContext({
      ...devices['iPhone 13'],
    });
    
    const mobilePage = await mobileContext?.newPage();
    if (mobilePage) {
      await use(mobilePage);
      await mobilePage.close();
      await mobileContext?.close();
    } else {
      throw new Error('Could not create mobile page context');
    }
  }
});

// Export expect for use in tests
export { expect };

// Re-export all Playwright test types
export type { Page, BrowserContext, Browser } from '@playwright/test';
