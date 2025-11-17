import { test, expect } from '../helpers/test-fixtures';
import { AuthPage, DashboardPage, NutritionPage, ProfilePage } from '../helpers/page-objects';
import { faker } from '@faker-js/faker';

/**
 * Complete User Workflow E2E Tests
 * 
 * This test suite covers the complete user journey from registration to data generation
 * including authentication, navigation, and core functionality.
 */

test.describe('Complete User Workflow', () => {
  let authPage: AuthPage;
  let dashboardPage: DashboardPage;
  let nutritionPage: NutritionPage;
  let profilePage: ProfilePage;

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page);
    dashboardPage = new DashboardPage(page);
    nutritionPage = new NutritionPage(page);
    profilePage = new ProfilePage(page);
  });

  test('should complete full user journey from registration to data generation', async ({ page, testData }) => {
    const testUser = {
      email: faker.internet.email(),
      password: 'TestPassword123!',
      firstName: faker.name.firstName(),
      lastName: faker.name.lastName(),
      dateOfBirth: '1990-01-01',
      gender: 'Male'
    };

    // Step 1: Navigate to the application
    await page.goto('/');
    await authPage.waitForPageLoad();
    await expect(page).toHaveTitle(/Nutrition Platform/);

    // Step 2: User Registration
    await authPage.registerLink.click();
    await authPage.waitForPageLoad();
    
    // Fill registration form
    await authPage.register(testUser);
    
    // Verify successful registration and redirect to dashboard
    await authPage.expectLoginSuccess();
    await dashboardPage.expectDashboardLoaded();
    
    // Step 3: Complete Profile Setup
    await profilePage.navigateToProfile();
    await profilePage.expectProfileLoaded();
    
    // Step 4: Navigate to Dashboard and verify data
    await dashboardPage.navigateToDashboard();
    await dashboardPage.expectDashboardLoaded();
    await dashboardPage.expectNutritionDataVisible();
    
    // Step 5: Add Nutrition Data
    await nutritionPage.navigateTo('/nutrition');
    await nutritionPage.searchFood('chicken breast');
    await nutritionPage.expectFoodResults();
    
    // Add food to meal
    await nutritionPage.addFirstFoodToMeal('Lunch', 100);
    
    // Step 6: Verify data appears in dashboard
    await dashboardPage.navigateToDashboard();
    await dashboardPage.expectDashboardLoaded();
    
    // Step 7: Logout
    await authPage.logoutButton.click();
    await authPage.waitForPageLoad();
    
    // Verify logged out state
    await expect(page).toHaveURL(/\/(login|register)?$/);
  });

  test('should handle login with existing user', async ({ page, testData }) => {
    const existingUser = testData.users.find(u => u.email === 'testuser@example.com');
    
    if (!existingUser) {
      test.skip('No existing test user found');
    }

    // Navigate to login page
    await page.goto('/auth/login');
    await authPage.waitForPageLoad();

    // Login with existing credentials
    await authPage.login(existingUser.email, existingUser.password);
    
    // Verify successful login
    await authPage.expectLoginSuccess();
    await dashboardPage.expectDashboardLoaded();
  });

  test('should handle login with invalid credentials', async ({ page }) => {
    // Navigate to login page
    await page.goto('/auth/login');
    await authPage.waitForPageLoad();

    // Attempt login with invalid credentials
    await authPage.login('invalid@example.com', 'invalidpassword');
    
    // Verify error message
    await authPage.expectValidationError();
  });

  test('should handle form validation errors', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');
    await authPage.waitForPageLoad();

    // Submit form without filling required fields
    await authPage.registerButton.click();
    
    // Check for validation errors
    await expect(authPage.errorMessage).toBeVisible();
    
    // Test email validation
    await authPage.emailInput.fill('invalid-email');
    await authPage.passwordInput.fill('123');
    await authPage.registerButton.click();
    
    // Verify validation messages appear
    await authPage.expectValidationError();
  });

  test('should handle password mismatch', async ({ page }) => {
    const testUser = {
      email: faker.internet.email(),
      password: 'TestPassword123!',
      confirmPassword: 'DifferentPassword123!',
      firstName: faker.name.firstName(),
      lastName: faker.name.lastName(),
      dateOfBirth: '1990-01-01',
      gender: 'Male'
    };

    // Navigate to registration page
    await page.goto('/auth/register');
    await authPage.waitForPageLoad();

    // Fill form with mismatched passwords
    await authPage.firstNameInput.fill(testUser.firstName);
    await authPage.lastNameInput.fill(testUser.lastName);
    await authPage.emailInput.fill(testUser.email);
    await authPage.passwordInput.fill(testUser.password);
    await authPage.confirmPasswordInput.fill(testUser.confirmPassword);
    await authPage.registerButton.click();

    // Verify password mismatch error
    await authPage.expectValidationError('Passwords do not match');
  });
});

test.describe('Authentication Edge Cases', () => {
  let authPage: AuthPage;

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page);
  });

  test('should handle session expiration', async ({ page, testData }) => {
    const existingUser = testData.users.find(u => u.email === 'testuser@example.com');
    
    if (!existingUser) {
      test.skip('No existing test user found');
    }

    // Login
    await page.goto('/auth/login');
    await authPage.login(existingUser.email, existingUser.password);
    await authPage.expectLoginSuccess();

    // Clear localStorage to simulate session expiration
    await page.evaluate(() => {
      localStorage.clear();
    });

    // Try to access protected route
    await page.goto('/dashboard');
    
    // Should redirect to login
    await expect(page).toHaveURL(/login/);
  });

  test('should handle concurrent login attempts', async ({ browser, testData }) => {
    const existingUser = testData.users.find(u => u.email === 'testuser@example.com');
    
    if (!existingUser) {
      test.skip('No existing test user found');
    }

    // Create two contexts
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();
    
    const authPage1 = new AuthPage(page1);
    const authPage2 = new AuthPage(page2);

    // Login from both contexts
    await Promise.all([
      page1.goto('/auth/login'),
      page2.goto('/auth/login')
    ]);

    await Promise.all([
      authPage1.login(existingUser.email, existingUser.password),
      authPage2.login(existingUser.email, existingUser.password)
    ]);

    // Both should succeed
    await Promise.all([
      expect(page1).toHaveURL(/dashboard/),
      expect(page2).toHaveURL(/dashboard/)
    ]);

    // Cleanup
    await context1.close();
    await context2.close();
  });

  test('should handle rate limiting on login attempts', async ({ page }) => {
    // Attempt multiple failed logins
    for (let i = 0; i < 5; i++) {
      await page.goto('/auth/login');
      await authPage.login('invalid@example.com', 'wrongpassword');
      
      if (i < 4) {
        await page.waitForTimeout(1000);
      }
    }

    // On the 5th attempt, should show rate limiting message
    await authPage.expectValidationError('Too many login attempts');
  });
});
