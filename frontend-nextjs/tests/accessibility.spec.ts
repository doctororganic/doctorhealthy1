import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Accessibility Tests (WCAG 2.1 Level AA)', () => {
  
  test('Homepage should have no accessibility violations', async ({ page }) => {
    await page.goto('/');
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();
    
    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Dashboard should have no accessibility violations', async ({ page }) => {
    await page.goto('/dashboard');
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa'])
      .analyze();
    
    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Meals page should have proper ARIA labels', async ({ page }) => {
    await page.goto('/dashboard/meals');
    
    // Check for ARIA landmarks
    const main = await page.locator('main[role="main"]');
    await expect(main).toBeVisible();
    
    // Check for proper heading hierarchy
    const h1 = await page.locator('h1');
    await expect(h1).toHaveCount(1);
    
    // Check for keyboard navigation
    await page.keyboard.press('Tab');
    const focusedElement = await page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });

  test('Forms should have proper labels', async ({ page }) => {
    await page.goto('/dashboard/meals/new');
    
    // All inputs should have labels
    const inputs = await page.locator('input');
    const count = await inputs.count();
    
    for (let i = 0; i < count; i++) {
      const input = inputs.nth(i);
      const id = await input.getAttribute('id');
      const label = await page.locator(`label[for="${id}"]`);
      await expect(label).toBeVisible();
    }
  });

  test('Color contrast should meet WCAG AA standards', async ({ page }) => {
    await page.goto('/');
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2aa'])
      .include(['body'])
      .analyze();
    
    const contrastViolations = accessibilityScanResults.violations.filter(
      v => v.id === 'color-contrast'
    );
    
    expect(contrastViolations).toEqual([]);
  });

  test('Images should have alt text', async ({ page }) => {
    await page.goto('/');
    
    const images = await page.locator('img');
    const count = await images.count();
    
    for (let i = 0; i < count; i++) {
      const img = images.nth(i);
      const alt = await img.getAttribute('alt');
      expect(alt).toBeTruthy();
      expect(alt).not.toBe('');
    }
  });

  test('Keyboard navigation should work', async ({ page }) => {
    await page.goto('/dashboard');
    
    // Tab through all interactive elements
    const interactiveElements = await page.locator('a, button, input, select, textarea');
    const count = await interactiveElements.count();
    
    for (let i = 0; i < count; i++) {
      await page.keyboard.press('Tab');
      const focused = await page.locator(':focus');
      await expect(focused).toBeVisible();
    }
  });

  test('Skip to main content link should exist', async ({ page }) => {
    await page.goto('/');
    
    const skipLink = await page.locator('a[href="#main-content"]');
    await expect(skipLink).toBeInViewport();
    
    // Should be visible on focus
    await skipLink.focus();
    await expect(skipLink).toBeVisible();
  });

  test('Focus indicators should be visible', async ({ page }) => {
    await page.goto('/dashboard');
    
    const button = await page.locator('button').first();
    await button.focus();
    
    // Check for visible focus indicator
    const outline = await button.evaluate(el => {
      const styles = window.getComputedStyle(el);
      return styles.outline || styles.boxShadow;
    });
    
    expect(outline).not.toBe('none');
    expect(outline).not.toBe('');
  });

  test('Screen reader announcements should work', async ({ page }) => {
    await page.goto('/dashboard/meals');
    
    // Check for aria-live regions
    const liveRegion = await page.locator('[aria-live="polite"]');
    await expect(liveRegion).toBeAttached();
  });

  test('Modal dialogs should trap focus', async ({ page }) => {
    await page.goto('/dashboard');
    
    // Open modal
    await page.click('button:has-text("Add Meal")');
    
    // Check for proper ARIA attributes
    const modal = await page.locator('[role="dialog"]');
    await expect(modal).toBeVisible();
    await expect(modal).toHaveAttribute('aria-modal', 'true');
    
    // Focus should be trapped
    await page.keyboard.press('Tab');
    const focused = await page.locator(':focus');
    const isInsideModal = await focused.evaluate((el, modalEl) => {
      return modalEl.contains(el);
    }, await modal.elementHandle());
    
    expect(isInsideModal).toBe(true);
  });

  test('Error messages should be announced', async ({ page }) => {
    await page.goto('/dashboard/meals/new');
    
    // Submit form with errors
    await page.click('button[type="submit"]');
    
    // Check for aria-describedby on error fields
    const errorInput = await page.locator('input[aria-invalid="true"]').first();
    const describedBy = await errorInput.getAttribute('aria-describedby');
    expect(describedBy).toBeTruthy();
    
    // Error message should exist
    const errorMessage = await page.locator(`#${describedBy}`);
    await expect(errorMessage).toBeVisible();
  });

  test('Loading states should be announced', async ({ page }) => {
    await page.goto('/dashboard/meals');
    
    // Trigger loading state
    await page.click('button:has-text("Load More")');
    
    // Check for loading indicator with proper ARIA
    const loader = await page.locator('[role="status"][aria-live="polite"]');
    await expect(loader).toBeVisible();
    await expect(loader).toContainText(/loading/i);
  });
});

test.describe('Responsive Design Tests', () => {
  
  test('Mobile viewport should be accessible', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2aa'])
      .analyze();
    
    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Tablet viewport should be accessible', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/dashboard');
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2aa'])
      .analyze();
    
    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Touch targets should be at least 44x44px', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/dashboard');
    
    const buttons = await page.locator('button, a');
    const count = await buttons.count();
    
    for (let i = 0; i < count; i++) {
      const button = buttons.nth(i);
      const box = await button.boundingBox();
      
      if (box) {
        expect(box.width).toBeGreaterThanOrEqual(44);
        expect(box.height).toBeGreaterThanOrEqual(44);
      }
    }
  });
});
