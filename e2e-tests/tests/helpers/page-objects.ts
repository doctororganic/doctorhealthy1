import { Page, expect } from '@playwright/test';

/**
 * Page Object Models for Nutrition Platform
 * 
 * This file contains page object models for major pages and components
 * to maintain clean, maintainable, and reusable test code.
 */

// Base Page Object
export class BasePage {
  constructor(public page: Page) {}

  // Common navigation elements
  get navigationMenu() {
    return this.page.locator('nav[role="navigation"]');
  }

  get homeLink() {
    return this.page.locator('a[href="/"], nav a:has-text("Home")');
  }

  get dashboardLink() {
    return this.page.locator('a[href="/dashboard"], nav a:has-text("Dashboard")');
  }

  get profileLink() {
    return this.page.locator('a[href="/profile"], nav a:has-text("Profile")');
  }

  get nutritionLink() {
    return this.page.locator('a[href="/nutrition"], nav a:has-text("Nutrition")');
  }

  get workoutsLink() {
    return this.page.locator('a[href="/workouts"], nav a:has-text("Workouts")');
  }

  get languageToggle() {
    return this.page.locator('[data-testid="language-toggle"], button:has-text("العربية"), button:has-text("English")');
  }

  get logoutButton() {
    return this.page.locator('[data-testid="logout"], button:has-text("Logout"), button:has-text("تسجيل الخروج")');
  }

  // Common actions
  async navigateTo(url: string) {
    await this.page.goto(url);
    await this.page.waitForLoadState('networkidle');
  }

  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForTimeout(1000); // Wait for any JavaScript to execute
  }

  async takeScreenshot(name: string) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    await this.page.screenshot({ 
      path: `test-results/screenshots/${name}-${timestamp}.png`,
      fullPage: true 
    });
  }

  async checkAccessibility() {
    // Basic accessibility checks
    await expect(this.page.locator('h1')).toBeVisible();
    await expect(this.page.locator('main, [role="main"]')).toBeVisible();
  }

  async checkResponsive(viewport: { width: number; height: number }) {
    await this.page.setViewportSize(viewport);
    await this.waitForPageLoad();
  }
}

// Authentication Page Object
export class AuthPage extends BasePage {
  // Login form elements
  get loginForm() {
    return this.page.locator('form[data-testid="login-form"], form:has(input[type="email"])');
  }

  get emailInput() {
    return this.page.locator('input[name="email"], input[type="email"], input[id="email"]');
  }

  get passwordInput() {
    return this.page.locator('input[name="password"], input[type="password"], input[id="password"]');
  }

  get loginButton() {
    return this.page.locator('button[type="submit"], button:has-text("Login"), button:has-text("Sign In")');
  }

  get registerLink() {
    return this.page.locator('a[href*="register"], a:has-text("Register"), a:has-text("Sign Up")');
  }

  // Registration form elements
  get registrationForm() {
    return this.page.locator('form[data-testid="register-form"]');
  }

  get firstNameInput() {
    return this.page.locator('input[name="firstName"], input[name="first_name"], input[id="firstName"]');
  }

  get lastNameInput() {
    return this.page.locator('input[name="lastName"], input[name="last_name"], input[id="lastName"]');
  }

  get confirmPasswordInput() {
    return this.page.locator('input[name="confirmPassword"], input[name="confirm_password"]');
  }

  get dateOfBirthInput() {
    return this.page.locator('input[name="dateOfBirth"], input[name="date_of_birth"], input[type="date"]');
  }

  get genderSelect() {
    return this.page.locator('select[name="gender"], [data-testid="gender-select"]');
  }

  get registerButton() {
    return this.page.locator('button[type="submit"], button:has-text("Register"), button:has-text("Sign Up")');
  }

  // Error messages
  get errorMessage() {
    return this.page.locator('[data-testid="error"], .error, .alert-danger, [role="alert"]');
  }

  get successMessage() {
    return this.page.locator('[data-testid="success"], .success, .alert-success');
  }

  // Form validation
  async login(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.loginButton.click();
    await this.waitForPageLoad();
  }

  async register(userData: {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    dateOfBirth: string;
    gender: string;
  }) {
    await this.firstNameInput.fill(userData.firstName);
    await this.lastNameInput.fill(userData.lastName);
    await this.emailInput.fill(userData.email);
    await this.passwordInput.fill(userData.password);
    await this.confirmPasswordInput.fill(userData.password);
    await this.dateOfBirthInput.fill(userData.dateOfBirth);
    await this.genderSelect.selectOption({ label: userData.gender });
    await this.registerButton.click();
    await this.waitForPageLoad();
  }

  async expectLoginSuccess() {
    await expect(this.page).toHaveURL(/dashboard|home/);
    await expect(this.dashboardLink).toBeVisible();
  }

  async expectValidationError(message?: string) {
    await expect(this.errorMessage).toBeVisible();
    if (message) {
      await expect(this.errorMessage).toContainText(message);
    }
  }
}

// Dashboard Page Object
export class DashboardPage extends BasePage {
  get pageTitle() {
    return this.page.locator('h1, [data-testid="dashboard-title"]');
  }

  get welcomeMessage() {
    return this.page.locator('[data-testid="welcome"], h1:has-text("Welcome"), h1:has-text("مرحبا")');
  }

  get nutritionSummary() {
    return this.page.locator('[data-testid="nutrition-summary"], .nutrition-card');
  }

  get caloriesConsumed() {
    return this.page.locator('[data-testid="calories-consumed"], .calories-consumed');
  }

  get caloriesRemaining() {
    return this.page.locator('[data-testid="calories-remaining"], .calories-remaining');
  }

  get proteinIntake() {
    return this.page.locator('[data-testid="protein-intake"], .protein-intake');
  }

  get carbsIntake() {
    return this.page.locator('[data-testid="carbs-intake"], .carbs-intake');
  }

  get fatIntake() {
    return this.page.locator('[data-testid="fat-intake"], .fat-intake');
  }

  get recentMeals() {
    return this.page.locator('[data-testid="recent-meals"], .recent-meals, .meal-list');
  }

  get addMealButton() {
    return this.page.locator('[data-testid="add-meal"], button:has-text("Add Meal"), button:has-text("إضافة وجبة")');
  }

  get nutritionGoalsButton() {
    return this.page.locator('[data-testid="nutrition-goals"], button:has-text("Goals"), button:has-text("الأهداف")');
  }

  async navigateToDashboard() {
    await this.navigateTo('/dashboard');
  }

  async expectDashboardLoaded() {
    await expect(this.pageTitle).toBeVisible();
    await expect(this.nutritionSummary).toBeVisible();
  }

  async expectNutritionDataVisible() {
    await expect(this.caloriesConsumed).toBeVisible();
    await expect(this.caloriesRemaining).toBeVisible();
    await expect(this.proteinIntake).toBeVisible();
    await expect(this.carbsIntake).toBeVisible();
    await expect(this.fatIntake).toBeVisible();
  }
}

// Nutrition Page Object
export class NutritionPage extends BasePage {
  get pageTitle() {
    return this.page.locator('h1:has-text("Nutrition"), h1:has-text("التغذية")');
  }

  get foodSearchInput() {
    return this.page.locator('input[placeholder*="search"], input[data-testid="food-search"], #food-search');
  }

  get searchButton() {
    return this.page.locator('button:has-text("Search"), button:has-text("بحث")');
  }

  get foodResults() {
    return this.page.locator('[data-testid="food-results"], .food-list, .search-results');
  }

  get foodItems() {
    return this.page.locator('[data-testid="food-item"], .food-item');
  }

  get addFoodButton() {
    return this.page.locator('[data-testid="add-food"], button:has-text("Add Food"), button:has-text("إضافة طعام")');
  }

  get mealTypeSelect() {
    return this.page.locator('select[name="mealType"], [data-testid="meal-type"]');
  }

  get quantityInput() {
    return this.page.locator('input[name="quantity"], input[type="number"], [data-testid="quantity"]');
  }

  get saveMealButton() {
    return this.page.locator('[data-testid="save-meal"], button:has-text("Save"), button:has-text("حفظ")');
  }

  async searchFood(query: string) {
    await this.foodSearchInput.fill(query);
    await this.searchButton.click();
    await this.page.waitForTimeout(2000); // Wait for search results
  }

  async addFirstFoodToMeal(mealType: string, quantity: number) {
    await this.foodItems.first().click();
    await this.mealTypeSelect.selectOption({ label: mealType });
    await this.quantityInput.fill(quantity.toString());
    await this.saveMealButton.click();
    await this.waitForPageLoad();
  }

  async expectFoodResults() {
    await expect(this.foodResults).toBeVisible();
    await expect(this.foodItems.first()).toBeVisible();
  }
}

// Profile Page Object
export class ProfilePage extends BasePage {
  get pageTitle() {
    return this.page.locator('h1:has-text("Profile"), h1:has-text("الملف الشخصي")');
  }

  get profileForm() {
    return this.page.locator('form[data-testid="profile-form"], .profile-form');
  }

  get editProfileButton() {
    return this.page.locator('[data-testid="edit-profile"], button:has-text("Edit Profile")');
  }

  get saveProfileButton() {
    return this.page.locator('[data-testid="save-profile"], button:has-text("Save"), button:has-text("حفظ")');
  }

  get deleteAccountButton() {
    return this.page.locator('[data-testid="delete-account"], button:has-text("Delete Account")');
  }

  get changePasswordButton() {
    return this.page.locator('[data-testid="change-password"], button:has-text("Change Password")');
  }

  async navigateToProfile() {
    await this.navigateTo('/profile');
  }

  async expectProfileLoaded() {
    await expect(this.pageTitle).toBeVisible();
    await expect(this.profileForm).toBeVisible();
  }
}

// Settings Page Object
export class SettingsPage extends BasePage {
  get pageTitle() {
    return this.page.locator('h1:has-text("Settings"), h1:has-text("الإعدادات")');
  }

  get languageSelector() {
    return this.page.locator('select[name="language"], [data-testid="language-selector"]');
  }

  get themeToggle() {
    return this.page.locator('[data-testid="theme-toggle"], .theme-switch');
  }

  get notificationsToggle() {
    return this.page.locator('[data-testid="notifications-toggle"], input[type="checkbox"][name="notifications"]');
  }

  get saveSettingsButton() {
    return this.page.locator('[data-testid="save-settings"], button:has-text("Save Settings")');
  }

  async setLanguage(language: string) {
    await this.languageSelector.selectOption({ label: language });
  }

  async toggleTheme() {
    await this.themeToggle.click();
  }

  async toggleNotifications() {
    await this.notificationsToggle.check();
  }

  async saveSettings() {
    await this.saveSettingsButton.click();
    await this.waitForPageLoad();
  }
}
