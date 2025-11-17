import { defineConfig, devices } from '@playwright/test';
import path from 'path';

/**
 * Comprehensive Playwright Configuration for Nutrition Platform E2E Tests
 * 
 * Features:
 * - Multiple device configurations (Desktop, Tablet, Mobile)
 * - Bilingual support testing (English/Arabic)
 * - Cross-browser testing (Chrome, Firefox, Safari, Edge)
 * - CI/CD integration ready
 * - Visual regression testing
 * - Network interception for mocking
 * - Accessibility testing
 */
export default defineConfig({
  // Global test configuration
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { outputFolder: 'playwright-report' }],
    ['json', { outputFile: 'test-results/results.json' }],
    ['junit', { outputFile: 'test-results/results.xml' }],
    ['line'],
    ['list']
  ],

  // Global settings for all tests
  use: {
    // Base URL for the application
    baseURL: process.env.BASE_URL || 'http://localhost:3000',

    // API base URL
    extraHTTPHeaders: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },

    // Trace configuration for debugging
    trace: 'retain-on-failure',

    // Screenshot configuration
    screenshot: 'only-on-failure',

    // Video configuration
    video: 'retain-on-failure',

    // Test timeout
    actionTimeout: 30000,
    navigationTimeout: 60000,

    // Locale for testing
    locale: 'en-US',

    // Timezone
    timezoneId: 'UTC',

    // User agent
    userAgent: 'Nutrition-Platform-E2E-Tests/1.0',

    // Ignore HTTPS errors for development
    ignoreHTTPSErrors: !process.env.PROD,
  },

  // Configure projects for different browsers and devices
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },

    // Tablet tests
    {
      name: 'tablet',
      use: { ...devices['iPad Pro'] },
    },

    // Mobile tests
    {
      name: 'mobile',
      use: { ...devices['iPhone 13'] },
    },

    // Arabic locale tests
    {
      name: 'arabic-desktop',
      use: { 
        ...devices['Desktop Chrome'],
        locale: 'ar-SA',
        extraHTTPHeaders: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Accept-Language': 'ar-SA',
        },
      },
    },

    {
      name: 'arabic-mobile',
      use: { 
        ...devices['iPhone 13'],
        locale: 'ar-SA',
        extraHTTPHeaders: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'Accept-Language': 'ar-SA',
        },
      },
    },

    // Responsive design tests
    {
      name: 'responsive-tests',
      use: { 
        ...devices['Desktop Chrome'],
        viewport: { width: 1920, height: 1080 },
      },
    },

    // Accessibility tests
    {
      name: 'accessibility',
      use: { 
        ...devices['Desktop Chrome'],
        viewport: { width: 1280, height: 720 },
      },
    },

    // Visual regression tests
    {
      name: 'visual-regression',
      use: { 
        ...devices['Desktop Chrome'],
        viewport: { width: 1280, height: 720 },
      },
    },

    // Dark mode tests
    {
      name: 'dark-mode',
      use: { 
        ...devices['Desktop Chrome'],
        colorScheme: 'dark',
        viewport: { width: 1280, height: 720 },
      },
    },
  ],

  // Web server configuration (for local development)
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120000,
  },

  // Global setup and teardown
  globalSetup: path.join(__dirname, 'global-setup.ts'),
  globalTeardown: path.join(__dirname, 'global-teardown.ts'),

  // Output directory
  outputDir: 'test-results',

  // Test timeout
  timeout: 120000,

  // Expect timeout
  expect: {
    timeout: 15000,
  },

  // Metadata
  metadata: {
    'Test Environment': process.env.NODE_ENV || 'test',
    'Application': 'Nutrition Platform',
    'Test Suite': 'E2E Integration Tests',
    'Version': '1.0.0',
  },
});
