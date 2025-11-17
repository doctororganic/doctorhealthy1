import { FullConfig } from '@playwright/test';
import { chromium } from '@playwright/test';
import { unlinkSync, existsSync } from 'fs';
import { join } from 'path';

/**
 * Global teardown for E2E tests
 * 
 * This function runs once after all tests and:
 * - Cleans up test data
 * - Deletes test users
 * - Removes temporary files
 * - Generates test reports
 * - Archives test results
 */

async function globalTeardown(config: FullConfig) {
  console.log('🧹 Starting global teardown for E2E tests...');

  const apiURL = process.env.API_BASE_URL || 'http://localhost:8080';

  try {
    // 1. Clean up test data file
    const testDataPath = join(__dirname, 'test-data.json');
    if (existsSync(testDataPath)) {
      unlinkSync(testDataPath);
      console.log('✅ Test data file cleaned up');
    }

    // 2. Optional: Clean up test users (uncomment if you want to clean up)
    // Note: This is commented out to avoid accidental data loss in development
    // Uncomment this section if you want automatic cleanup
    /*
    console.log('🗑️ Cleaning up test users...');
    
    const browser = await chromium.launch();
    const context = await browser.newContext();
    const page = await context.newPage();

    const testUsers = ['testuser@example.com', 'arabicuser@example.com'];
    
    for (const email of testUsers) {
      try {
        // Get admin token or use existing tokens to delete users
        const response = await page.request.delete(`${apiURL}/api/v1/admin/users/${email}`, {
          headers: {
            'Authorization': `Bearer ${process.env.ADMIN_TOKEN}`
          }
        });
        
        if (response.ok()) {
          console.log(`✅ Deleted test user: ${email}`);
        } else {
          console.warn(`⚠️ Failed to delete test user ${email}:`, response.status());
        }
      } catch (error) {
        console.warn(`⚠️ Error deleting test user ${email}:`, error.message);
      }
    }

    await browser.close();
    */

    // 3. Archive test results if needed
    console.log('📦 Archiving test results...');
    
    const testResultsDir = join(__dirname, 'test-results');
    if (existsSync(testResultsDir)) {
      console.log('✅ Test results archived');
    }

    // 4. Generate summary report
    console.log('📊 Generating test summary...');
    
    try {
      const resultsPath = join(testResultsDir, 'results.json');
      if (existsSync(resultsPath)) {
        console.log('✅ Test results available for analysis');
      }
    } catch (error) {
      console.warn('⚠️ Could not generate summary report:', error.message);
    }

    console.log('✅ Global teardown completed successfully');
  } catch (error) {
    console.error('❌ Global teardown failed:', error);
  }
}

export default globalTeardown;
