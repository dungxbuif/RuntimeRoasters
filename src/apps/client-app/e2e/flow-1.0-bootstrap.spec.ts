import { test, expect } from '@playwright/test';
import { loginViaUi, logoutViaUi } from './helpers';

/**
 * TC-1.0: Fresh System Setup & ADMIN Bootstrap
 * Covers TC-1.2, TC-1.3, TC-1.7, and RBAC Restrictions
 */
test.describe('Flow 1.0: ADMIN System Bootstrap & RBAC', () => {
  test.beforeEach(async ({ page }) => {
    // Clear cookies/storage to start fresh
    await page.context().clearCookies();
  });

  test('TC-1.2 & TC-1.3: Admin Seeding and Intelligence Hub', async ({ page }) => {
    await loginViaUi(page);
    await page.goto('http://localhost:3000/dashboard');
    await expect(page.locator('h1:has-text("Intelligence Hub")')).toBeVisible({ timeout: 15000 });
    
    await expect(page.locator('text=Pilot Farms')).toBeVisible();
  });

  test('RBAC-1.1: anonymous users are redirected away from admin-only pages', async ({ page }) => {
    await page.goto('http://localhost:3000/dashboard/users');
    await expect(page).toHaveURL(/\/login/, { timeout: 30000 });
  });

  test('UI-Visibility: Admin vs Manager Sidebar', async ({ page }) => {
    await loginViaUi(page);
    await page.goto('http://localhost:3000/dashboard');
    await expect(page.getByRole('heading', { name: /Intelligence Hub/i })).toBeVisible();
    await expect(page.locator('text=Provenance Trace')).toBeVisible();
    await expect(page.locator('text=System Explorer')).not.toBeVisible();

    await logoutViaUi(page);
    await page.goto('http://localhost:3000/dashboard');
    await expect(page).toHaveURL(/\/login/, { timeout: 30000 });
  });
});
