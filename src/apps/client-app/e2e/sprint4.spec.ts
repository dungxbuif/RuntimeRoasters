import { test, expect } from '@playwright/test';

/**
 * Sprint 4 E2E Test: Smart Harvest Declaration
 * Rules:
 * - Admin logs in
 * - Admin navigates to Harvest Declaration
 * - Admin records a new harvest
 * - Verification: Harvest appears in ledger
 * - Verification: Validation prevents invalid (0 quantity) harvest
 */

const ADMIN_EMAIL = 'admin@runtimeroasters.com';
const ADMIN_PASS = 'Hello@123';
const TEST_QUANTITY = Math.floor(Math.random() * 1000) + 1;

const SELECTORS = {
  LOGIN_SUBMIT: '[data-e2e="login-submit"]',
  AUTH_LOADER: '[data-e2e="auth-loader"]',
  DASHBOARD_HEADER: '[data-e2e="dashboard-header"]',
  
  // Harvest Page
  DECLARE_HARVEST_BTN: 'button:has-text("Declare New Harvest")',
  HARVEST_FARM_SELECT: 'select', // First select is farm
  HARVEST_TYPE_SELECT: 'select:nth-of-type(2)', // Second select is type
  HARVEST_QUANTITY_INPUT: 'input[type="number"]',
  HARVEST_SUBMIT_BTN: 'button:has-text("Record Harvest")',
  HARVEST_LEDGER_TABLE: 'table',
};

test.describe('Sprint 4: Smart Harvest Declaration', () => {

  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies();
  });

  test('Should declare harvest successfully', async ({ page }) => {
    await test.step('Step 1: Admin Login', async () => {
      await page.goto('/login');
      await page.fill('input[type="email"], input[name="identifier"]', ADMIN_EMAIL);
      await page.fill('input[type="password"], input[name="password"]', ADMIN_PASS);
      await page.click(SELECTORS.LOGIN_SUBMIT);
      
      await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 15000 });
      await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
    });

    await test.step('Step 2: Navigate to Harvest Declaration', async () => {
      await page.goto('/dashboard/farm-ops/harvests');
      await page.waitForSelector(SELECTORS.AUTH_LOADER, { state: 'hidden' });
      await expect(page.locator('h1')).toContainText('Harvest Declaration');
    });

    await test.step('Step 3: Record New Harvest', async () => {
      await page.click(SELECTORS.DECLARE_HARVEST_BTN);
      
      // Select first farm if not selected
      const farmSelect = page.locator('select').first();
      await farmSelect.waitFor({ state: 'visible' });
      
      await page.fill(SELECTORS.HARVEST_QUANTITY_INPUT, TEST_QUANTITY.toString());
      
      // Record
      await page.click(SELECTORS.HARVEST_SUBMIT_BTN);
      
      // Modal should close
      await expect(page.locator(SELECTORS.HARVEST_SUBMIT_BTN)).not.toBeVisible();
      
      // Should see the quantity in the table
      await expect(page.locator(SELECTORS.HARVEST_LEDGER_TABLE)).toContainText(`${TEST_QUANTITY} KG`);
    });

    await test.step('Step 4: Verify Validation (Zero Quantity)', async () => {
      await page.click(SELECTORS.DECLARE_HARVEST_BTN);
      await page.fill(SELECTORS.HARVEST_QUANTITY_INPUT, '0');
      
      // Intercept alert
      page.once('dialog', async dialog => {
        expect(dialog.message()).toContain('greater than 0');
        await dialog.dismiss();
      });
      
      await page.click(SELECTORS.HARVEST_SUBMIT_BTN);
      
      // Modal should still be open
      await expect(page.locator(SELECTORS.HARVEST_SUBMIT_BTN)).toBeVisible();
    });
  });
});
