import { expect, test } from '@playwright/test';
import { ADMIN_PASS, loginViaUi, logoutViaUi } from './helpers';

const TEST_MANAGER_EMAIL = `mgr.${Date.now()}@runtimeroasters.com`;
const TEST_FARM_NAME = `Hectare Node ${Date.now()}`;
const UPDATED_FARM_NAME = `${TEST_FARM_NAME} Updated`;

test.describe('Sprint 3: Farm Service & Vertical Slice', () => {
  test('admin provisions manager and farm, updates it, and manager sees only assigned farm and cannot delete', async ({ browser, page }) => {
    await loginViaUi(page);

    await page.goto('/dashboard/users');
    await page.locator('[data-e2e="create-user-btn"]').click();
    await page.locator('[data-e2e="user-email-input"]').fill(TEST_MANAGER_EMAIL);
    await page.locator('[data-e2e="user-password-input"]').fill(ADMIN_PASS);
    await page.locator('[data-e2e="user-role-select"]').selectOption('FARM_MANAGER');
    await page.locator('[data-e2e="user-submit-btn"]').click();
    await expect(page.locator(`[data-e2e="user-row-${TEST_MANAGER_EMAIL}"]`)).toBeVisible({ timeout: 15000 });

    await page.goto('/dashboard/farm-ops/registry');
    await page.locator('[data-e2e="create-farm-btn"]').click();
    await page.locator('[data-e2e="farm-name-input"]').fill(TEST_FARM_NAME);
    await page.locator('[data-e2e="farm-area-input"]').fill('120.5');
    const managerLabel = `${TEST_MANAGER_EMAIL.split('@')[0]} (${TEST_MANAGER_EMAIL})`;
    await page.locator('[data-e2e="farm-owner-select"]').selectOption({ label: managerLabel });
    await page.locator('[data-e2e="farm-submit-btn"]').click();
    const adminFarmRow = page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`);
    await expect(adminFarmRow).toBeVisible({ timeout: 15000 });

    await adminFarmRow.locator('[data-e2e="farm-edit-btn"]').click();
    await expect(page.locator('[data-e2e="farm-modal"]')).toBeVisible();
    await page.locator('[data-e2e="farm-name-input"]').fill(UPDATED_FARM_NAME);
    await page.locator('[data-e2e="farm-area-input"]').fill('132.75');
    await page.locator('[data-e2e="farm-submit-btn"]').click();

    const updatedAdminFarmRow = page.locator(`[data-e2e="farm-row-${UPDATED_FARM_NAME}"]`);
    await expect(updatedAdminFarmRow).toBeVisible({ timeout: 15000 });

    const managerContext = await browser.newContext();
    const managerPage = await managerContext.newPage();

    await loginViaUi(managerPage, TEST_MANAGER_EMAIL, ADMIN_PASS);
    await managerPage.goto('/dashboard/farm-ops/registry');

    const managerFarmRow = managerPage.locator(`[data-e2e="farm-row-${UPDATED_FARM_NAME}"]`);
    await expect(managerFarmRow).toBeVisible({ timeout: 15000 });
    await expect(managerFarmRow.locator('[data-e2e="farm-delete-btn"]')).toHaveCount(0);
    await managerContext.close();

    page.once('dialog', async (dialog) => {
      await dialog.accept();
    });
    await updatedAdminFarmRow.locator('[data-e2e="farm-delete-btn"]').click();
    await expect(updatedAdminFarmRow).toHaveCount(0, { timeout: 15000 });

    await logoutViaUi(page);
  });
});
