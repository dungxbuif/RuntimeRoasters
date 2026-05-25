import { test, expect } from '@playwright/test';

/**
 * TC-1.0: Fresh System Setup & ADMIN Bootstrap
 * Covers TC-1.2, TC-1.3, TC-1.7, and RBAC Restrictions
 */
test.describe('Flow 1.0: ADMIN System Bootstrap & RBAC', () => {
  const ADMIN_EMAIL = 'admin@runtimeroasters.com';
  const ADMIN_PASS = 'Hello@123';
  const FARM_MANAGER_EMAIL = 'mgr.hn.hoankiem@runtimeroasters.com'; // Standardized seed email

  test.beforeEach(async ({ page }) => {
    // Clear cookies/storage to start fresh
    await page.context().clearCookies();
  });

  test('TC-1.2 & TC-1.3: Admin Seeding and Intelligence Hub', async ({ page }) => {
    await page.goto('http://localhost:3000/login');
    
    // Login as ADMIN
    await page.fill('input[name="identifier"]', ADMIN_EMAIL);
    await page.fill('input[name="password"]', ADMIN_PASS);
    await page.click('button[type="submit"]');

    // Wait for redirect - can land on /dashboard or /dashboard/users
    await page.waitForURL(/.*dashboard/, { timeout: 15000 });
    
    // Force bootstrap if modal doesn't show up automatically
    await page.goto('http://localhost:3000/dashboard?bootstrap=true');
    await expect(page.locator('h2:has-text("System Bootstrap")')).toBeVisible({ timeout: 10000 });

    // Trigger Seeding
    await page.click('button:has-text("Initialize Data")');
    await expect(page.locator('text=System data initialized successfully!')).toBeVisible({ timeout: 30000 });
    
    // Navigate to Intelligence Hub
    await page.goto('http://localhost:3000/dashboard');
    await expect(page.locator('h1:has-text("Intelligence Hub")')).toBeVisible();
    
    // Verify Admin-only KPI
    await expect(page.locator('text=Pilot Farms')).toBeVisible();
    await expect(page.locator('text=6')).toBeVisible();
  });

  test('RBAC-1.1: Farm Manager Restriction Test', async ({ page }) => {
    await page.goto('http://localhost:3000/login');
    
    // Login as FARM_MANAGER
    await page.fill('input[name="identifier"]', FARM_MANAGER_EMAIL);
    await page.fill('input[name="password"]', ADMIN_PASS);
    await page.click('button[type="submit"]');

    await page.waitForURL(/.*dashboard/, { timeout: 15000 });

    // Attempt to access ADMIN-only page via direct navigation
    await page.goto('http://localhost:3000/dashboard/users');
    
    // We expect them NOT to see the User table (it might redirect or show empty/guard)
    await expect(page.locator('table')).not.toBeVisible();
  });

  test('UI-Visibility: Admin vs Manager Sidebar', async ({ page }) => {
    // 1. Check Admin Sidebar
    await page.goto('http://localhost:3000/login');
    await page.fill('input[name="identifier"]', ADMIN_EMAIL);
    await page.fill('input[name="password"]', ADMIN_PASS);
    await page.click('button[type="submit"]');
    
    await page.waitForURL(/.*dashboard/, { timeout: 15000 });
    await expect(page.locator('text=Manage Users')).toBeVisible();
    await expect(page.locator('text=System Explorer')).toBeVisible();

    // 2. Check Manager Sidebar
    await page.click('[data-e2e="logout-btn"]');
    await page.waitForURL(/.*login/);
    
    await page.fill('input[name="identifier"]', FARM_MANAGER_EMAIL);
    await page.fill('input[name="password"]', ADMIN_PASS);
    await page.click('button[type="submit"]');

    await page.waitForURL(/.*dashboard/, { timeout: 15000 });

    // Manager should NOT see Manage Users or System Explorer
    await expect(page.locator('text=Manage Users')).not.toBeVisible();
    await expect(page.locator('text=System Explorer')).not.toBeVisible();
    
    // But SHOULD see Manage Farms
    await expect(page.locator('text=Manage Farms')).toBeVisible();
  });
});
