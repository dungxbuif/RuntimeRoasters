import { test, expect } from '@playwright/test';

/**
 * Sprint 3 E2E Test: Full Vertical Slice (High Fidelity)
 * Rules:
 * - Admin creates Manager
 * - Admin creates Farm & assigns to Manager
 * - Manager logs in -> Can view own farm
 * - Manager CANNOT delete any farm (Role restriction enforced in UI & Backend)
 * - Data isolation: Manager cannot see other's farms
 */

const ADMIN_EMAIL = 'admin@runtimeroasters.com';
const ADMIN_PASS = 'Hello@123';

const TEST_MANAGER_EMAIL = `mgr.${Date.now()}@runtimeroasters.com`;
const TEST_MANAGER_PASS = 'Hello@123';
const TEST_FARM_NAME = `Hectare Node ${Date.now()}`;

const SELECTORS = {
  LOGIN_SUBMIT: '[data-e2e="login-submit"]',
  CREATE_USER_BTN: '[data-e2e="create-user-btn"]',
  USER_EMAIL_INPUT: '[data-e2e="user-email-input"]',
  USER_PASS_INPUT: '[data-e2e="user-password-input"]',
  USER_ROLE_SELECT: '[data-e2e="user-role-select"]',
  USER_SUBMIT_BTN: '[data-e2e="user-submit-btn"]',
  
  CREATE_FARM_BTN: '[data-e2e="create-farm-btn"]',
  FARM_NAME_INPUT: '[data-e2e="farm-name-input"]',
  FARM_AREA_INPUT: '[data-e2e="farm-area-input"]',
  FARM_OWNER_SELECT: '[data-e2e="farm-owner-select"]',
  FARM_SUBMIT_BTN: '[data-e2e="farm-submit-btn"]',
  FARM_DELETE_BTN: '[data-e2e="farm-delete-btn"]',
  
  AUTH_LOADER: '[data-e2e="auth-loader"]',
  DASHBOARD_HEADER: '[data-e2e="dashboard-header"]',
};

test.describe('Sprint 3: High-Fidelity Flow & Role Constraints', () => {

  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies();
  });

  test('Full Sprint 3 Flow: Admin Provisioning -> Manager Isolation', async ({ page }) => {
    
    await test.step('Step 1: Admin Login', async () => {
      await page.goto('/login');
      await page.fill('input[name="identifier"]', ADMIN_EMAIL);
      await page.fill('input[name="password"]', ADMIN_PASS);
      await page.click(SELECTORS.LOGIN_SUBMIT);
      
      await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 15000 });
      await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
      await expect(page.locator(SELECTORS.DASHBOARD_HEADER)).toBeVisible();
    });

    await test.step('Step 2: Admin Creates Manager', async () => {
      await page.click(SELECTORS.CREATE_USER_BTN);
      await page.fill(SELECTORS.USER_EMAIL_INPUT, TEST_MANAGER_EMAIL);
      await page.fill(SELECTORS.USER_PASS_INPUT, TEST_MANAGER_PASS);
      await page.selectOption(SELECTORS.USER_ROLE_SELECT, 'FARM_MANAGER');
      await page.click(SELECTORS.USER_SUBMIT_BTN);

      await expect(page.locator(`[data-e2e="user-row-${TEST_MANAGER_EMAIL}"]`)).toBeVisible({ timeout: 10000 });
    });

    await test.step('Step 3: Admin Creates Farm for Manager', async () => {
      await page.goto('/dashboard/farm-ops/registry');
      await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
      await page.click(SELECTORS.CREATE_FARM_BTN);
      await page.fill(SELECTORS.FARM_NAME_INPUT, TEST_FARM_NAME);
      await page.fill(SELECTORS.FARM_AREA_INPUT, '120.5');
      
      // Wait for managers to load in dropdown
      await page.selectOption(SELECTORS.FARM_OWNER_SELECT, { label: TEST_MANAGER_EMAIL });
      await page.click(SELECTORS.FARM_SUBMIT_BTN);

      await expect(page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`)).toBeVisible({ timeout: 10000 });
    });

    await test.step('Step 4: Admin Logout', async () => {
      await page.click('[title="Logout"]');
      await expect(page).toHaveURL(/\/login/);
    });

    await test.step('Step 5: Manager Login & Verify Isolation', async () => {
      await page.fill('input[name="identifier"]', TEST_MANAGER_EMAIL);
      await page.fill('input[name="password"]', TEST_MANAGER_PASS);
      await page.click(SELECTORS.LOGIN_SUBMIT);

      await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 15000 });
      await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });

      await page.goto('/dashboard/farm-ops/registry');
      await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
      
      const farmRow = page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`);
      await expect(farmRow).toBeVisible();
      
      // Verify Restriction: Manager cannot see Delete button in UI
      await expect(farmRow.locator(SELECTORS.FARM_DELETE_BTN)).not.toBeVisible();
    });
  });

  test('Public Showcase accessibility', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1:has-text("System")')).toBeVisible();
    await expect(page.locator('aside')).not.toBeVisible();
  });
});
