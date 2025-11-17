import { test, expect } from '../helpers/test-fixtures';
import { AuthPage, DashboardPage, NutritionPage } from '../helpers/page-objects';

/**
 * Bilingual Support E2E Tests
 * 
 * This test suite verifies that the application properly supports both English and Arabic
 * including RTL layout, text direction, and translated content.
 */

test.describe('English Language Support', () => {
  let authPage: AuthPage;
  let dashboardPage: DashboardPage;
  let nutritionPage: NutritionPage;

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page);
    dashboardPage = new DashboardPage(page);
    nutritionPage = new NutritionPage(page);
  });

  test('should display English content correctly', async ({ page }) => {
    // Navigate to homepage
    await page.goto('/');
    await authPage.waitForPageLoad();

    // Check for English text content
    await expect(page.locator('body')).toContainText(/login|register|nutrition|dashboard/i);
    
    // Check LTR text direction
    const htmlDir = await page.getAttribute('html', 'dir');
    expect(htmlDir).toBe('ltr');
  });

  test('should show English error messages', async ({ page }) => {
    // Navigate to login page
    await page.goto('/auth/login');
    await authPage.waitForPageLoad();

    // Submit invalid login
    await authPage.login('invalid@example.com', 'wrongpassword');
    
    // Verify error message is in English
    await expect(authPage.errorMessage).toContainText(/invalid|incorrect|failed/i);
  });

  test('should handle English form labels and placeholders', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');
    await authPage.waitForPageLoad();

    // Check English form labels
    await expect(authPage.emailInput).toHaveAttribute('placeholder', /email/i);
    await expect(authPage.passwordInput).toHaveAttribute('placeholder', /password/i);
    await expect(authPage.firstNameInput).toHaveAttribute('placeholder', /first name/i);
    await expect(authPage.lastNameInput).toHaveAttribute('placeholder', /last name/i);
  });
});

test.describe('Arabic Language Support', () => {
  let authPage: AuthPage;
  let dashboardPage: DashboardPage;
  let nutritionPage: NutritionPage;

  test.beforeEach(async ({ arabicPage }) => {
    authPage = new AuthPage(arabicPage);
    dashboardPage = new DashboardPage(arabicPage);
    nutritionPage = new NutritionPage(arabicPage);
  });

  test('should display Arabic content correctly', async ({ arabicPage }) => {
    // Navigate to homepage
    await arabicPage.goto('/');
    await authPage.waitForPageLoad();

    // Check for Arabic text content
    await expect(arabicPage.locator('body')).toContainText(/تسجيل الدخول|تسجيل|التغذية|لوحة التحكم/i);
    
    // Check RTL text direction
    const htmlDir = await arabicPage.getAttribute('html', 'dir');
    expect(htmlDir).toBe('rtl');
  });

  test('should show Arabic error messages', async ({ arabicPage }) => {
    // Navigate to login page
    await arabicPage.goto('/auth/login');
    await authPage.waitForPageLoad();

    // Submit invalid login
    await authPage.login('invalid@example.com', 'wrongpassword');
    
    // Verify error message is in Arabic
    await expect(authPage.errorMessage).toContainText(/غير صالح|فشل|خطأ/i);
  });

  test('should handle Arabic form labels and placeholders', async ({ arabicPage }) => {
    // Navigate to registration page
    await arabicPage.goto('/auth/register');
    await authPage.waitForPageLoad();

    // Check Arabic form labels
    await expect(authPage.emailInput).toHaveAttribute('placeholder', /البريد الإلكتروني/i);
    await expect(authPage.passwordInput).toHaveAttribute('placeholder', /كلمة المرور/i);
    await expect(authPage.firstNameInput).toHaveAttribute('placeholder', /الاسم الأول/i);
    await expect(authPage.lastNameInput).toHaveAttribute('placeholder', /الاسم الأخير/i);
  });

  test('should maintain RTL layout throughout application', async ({ arabicPage }) => {
    // Navigate through different pages and check RTL consistency
    await arabicPage.goto('/');
    await authPage.waitForPageLoad();
    
    // Check homepage
    let htmlDir = await arabicPage.getAttribute('html', 'dir');
    expect(htmlDir).toBe('rtl');

    // Navigate to login
    await arabicPage.goto('/auth/login');
    await authPage.waitForPageLoad();
    htmlDir = await arabicPage.getAttribute('html', 'dir');
    expect(htmlDir).toBe('rtl');

    // Navigate to dashboard (if accessible)
    await arabicPage.goto('/dashboard');
    await authPage.waitForPageLoad();
    htmlDir = await arabicPage.getAttribute('html', 'dir');
    expect(htmlDir).toBe('rtl');
  });

  test('should handle Arabic user registration', async ({ arabicPage, testData }) => {
    const arabicUser = testData.users.find(u => u.language === 'ar');
    
    if (!arabicUser) {
      test.skip('No Arabic test user found');
    }

    // Navigate to registration page
    await arabicPage.goto('/auth/register');
    await authPage.waitForPageLoad();

    // Fill form with Arabic data
    await authPage.firstNameInput.fill(arabicUser.firstName);
    await authPage.lastNameInput.fill(arabicUser.lastName);
    await authPage.emailInput.fill(arabicUser.email);
    await authPage.passwordInput.fill(arabicUser.password);
    await authPage.confirmPasswordInput.fill(arabicUser.password);
    await authPage.dateOfBirthInput.fill(arabicUser.dateOfBirth);
    await authPage.genderSelect.selectOption({ label: 'ذكر' });
    await authPage.registerButton.click();
    await authPage.waitForPageLoad();

    // Verify successful registration
    await expect(arabicPage.locator('body')).toContainText(/مرحبا|أهلا بك/i);
  });
});

test.describe('Language Switching', () => {
  let authPage: AuthPage;

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page);
  });

  test('should switch from English to Arabic', async ({ page }) => {
    // Start with English
    await page.goto('/');
    await authPage.waitForPageLoad();
    
    // Verify initial LTR direction
    let htmlDir = await page.getAttribute('html', 'dir');
    expect(htmlDir).toBe('ltr');

    // Find and click language toggle
    const languageToggle = page.locator('[data-testid="language-toggle"], button:has-text("العربية")').first();
    if (await languageToggle.isVisible()) {
      await languageToggle.click();
      await page.waitForTimeout(2000);

      // Verify RTL direction after switch
      htmlDir = await page.getAttribute('html', 'dir');
      expect(htmlDir).toBe('rtl');
      
      // Verify Arabic content
      await expect(page.locator('body')).toContainText(/العربية|تسجيل/i);
    } else {
      test.skip('Language toggle not found');
    }
  });

  test('should switch from Arabic to English', async ({ arabicPage }) => {
    // Start with Arabic
    await arabicPage.goto('/');
    await authPage.waitForPageLoad();
    
    // Verify initial RTL direction
    let htmlDir = await arabicPage.getAttribute('html', 'dir');
    expect(htmlDir).toBe('rtl');

    // Find and click language toggle
    const languageToggle = arabicPage.locator('[data-testid="language-toggle"], button:has-text("English")').first();
    if (await languageToggle.isVisible()) {
      await languageToggle.click();
      await arabicPage.waitForTimeout(2000);

      // Verify LTR direction after switch
      htmlDir = await arabicPage.getAttribute('html', 'dir');
      expect(htmlDir).toBe('ltr');
      
      // Verify English content
      await expect(arabicPage.locator('body')).toContainText(/English|Login/i);
    } else {
      test.skip('Language toggle not found');
    }
  });

  test('should maintain user session when switching languages', async ({ page, testData }) => {
    const existingUser = testData.users.find(u => u.email === 'testuser@example.com');
    
    if (!existingUser) {
      test.skip('No existing test user found');
    }

    // Login in English
    await page.goto('/auth/login');
    await authPage.login(existingUser.email, existingUser.password);
    await authPage.waitForPageLoad();

    // Verify logged in state
    await expect(page.locator('body')).toContainText(/dashboard|welcome/i);

    // Switch language
    const languageToggle = page.locator('[data-testid="language-toggle"], button:has-text("العربية")').first();
    if (await languageToggle.isVisible()) {
      await languageToggle.click();
      await page.waitForTimeout(2000);

      // Verify still logged in with Arabic content
      await expect(page.locator('body')).toContainText(/لوحة التحكم|مرحبا/i);
    } else {
      test.skip('Language toggle not found');
    }
  });
});

test.describe('Bi-directional Content', () => {
  test('should handle mixed language content correctly', async ({ page }) => {
    await page.goto('/');
    await authPage.waitForPageLoad();

    // Test that elements with mixed English/Arabic content display correctly
    const mixedContentElements = await page.locator('[lang="en"], [lang="ar"]').count();
    
    // If there are mixed language elements, they should have proper lang attributes
    if (mixedContentElements > 0) {
      const elements = page.locator('[lang]');
      for (let i = 0; i < await elements.count(); i++) {
        const element = elements.nth(i);
        const lang = await element.getAttribute('lang');
        expect(lang).toMatch(/^(en|ar)$/);
      }
    }
  });

  test('should handle Arabic numbers correctly', async ({ arabicPage }) => {
    await arabicPage.goto('/dashboard');
    await authPage.waitForPageLoad();

    // Check if Arabic numbers are used where appropriate
    const numberElements = arabicPage.locator('[data-testid="calories"], [data-testid="protein"], [data-testid="carbs"]');
    
    if (await numberElements.count() > 0) {
      // Numbers should display correctly in Arabic context
      const firstNumber = await numberElements.first().textContent();
      expect(firstNumber).toMatch(/\d+/); // Should contain digits
    }
  });
});

test.describe('Font and Typography', () => {
  test('should use appropriate fonts for Arabic content', async ({ arabicPage }) => {
    await arabicPage.goto('/');
    await authPage.waitForPageLoad();

    // Check for Arabic-appropriate fonts
    const bodyFont = await arabicPage.evaluate(() => {
      const style = window.getComputedStyle(document.body);
      return style.fontFamily;
    });

    // Should include Arabic-friendly fonts
    expect(bodyFont.toLowerCase()).toMatch(/arial|tahoma|noto|cairo|amiri/i);
  });

  test('should maintain readability in both languages', async ({ page, arabicPage }) => {
    // Test English readability
    await page.goto('/');
    await authPage.waitForPageLoad();
    
    const englishFontSize = await page.evaluate(() => {
      const style = window.getComputedStyle(document.body);
      return parseInt(style.fontSize);
    });
    expect(englishFontSize).toBeGreaterThan(12); // Minimum readable size

    // Test Arabic readability
    await arabicPage.goto('/');
    await authPage.waitForPageLoad();
    
    const arabicFontSize = await arabicPage.evaluate(() => {
      const style = window.getComputedStyle(document.body);
      return parseInt(style.fontSize);
    });
    expect(arabicFontSize).toBeGreaterThan(12); // Minimum readable size
  });
});
