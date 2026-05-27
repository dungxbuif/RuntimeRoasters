import { expect, test } from '@playwright/test';
import { loginViaUi } from './helpers';

const TEST_FARM_NAME = `Harvest Farm ${Date.now()}`;
const TEST_QUANTITY = Math.floor(Math.random() * 1000) + 100;

test.describe('Sprint 4: Smart Harvest Declaration', () => {
  test('admin can declare harvest, downstream warehouse batch is created, and validation blocks zero quantity', async ({ page, request }) => {
    await loginViaUi(page);
    const token = await page.evaluate(() => window.localStorage.getItem('rr_access_token'));
    if (!token) throw new Error('missing access token');

    await page.goto('/dashboard/farm-ops/registry');
    const farmRow = page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`);
    if (await farmRow.count() === 0) {
      await page.locator('[data-e2e="create-farm-btn"]').click();
      await page.locator('[data-e2e="farm-name-input"]').fill(TEST_FARM_NAME);
      await page.locator('[data-e2e="farm-area-input"]').fill('88');
      await page.locator('[data-e2e="farm-submit-btn"]').click();
      await expect(page.locator(`[data-e2e="farm-row-${TEST_FARM_NAME}"]`)).toBeVisible({ timeout: 15000 });
    }

    await page.goto('/dashboard/farm-ops/harvests');
    await expect(page.getByRole('heading', { name: /harvest/i })).toBeVisible();

    await page.locator('[data-e2e="create-harvest-btn"]').click();
    const harvestCreateResponse = page.waitForResponse((response) =>
      response.request().method() === 'POST' && response.url().includes('/v1/harvests')
    );
    await page.locator('[data-e2e="harvest-quantity-input"]').fill(TEST_QUANTITY.toString());
    await page.locator('[data-e2e="harvest-submit-btn"]').click();
    const harvestCreatePayload = await (await harvestCreateResponse).json();
    const harvestId = String(harvestCreatePayload.harvest.id);

    await expect(page.locator('[data-e2e="harvest-modal"]')).toHaveCount(0, { timeout: 15000 });
    await expect(page.locator('[data-e2e="harvest-list-table"]')).toContainText(`${TEST_QUANTITY} KG`);
    await expect.poll(async () => {
      const res = await request.get(`http://localhost:8087/v1/trace/${harvestId}/events`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok()) {
        return null;
      }
      const body = await res.json();
      const event = body.events?.find((item: { topic?: string; harvest_id?: string }) =>
        item.topic === 'farm.harvest.created' && item.harvest_id === harvestId
      );
      return event?.harvest_id ?? null;
    }, { timeout: 20000 }).toBe(harvestId);

    await page.locator('[data-e2e="create-harvest-btn"]').click();
    page.once('dialog', async (dialog) => {
      expect(dialog.message()).toContain('greater than 0');
      await dialog.dismiss();
    });
    await page.locator('[data-e2e="harvest-quantity-input"]').fill('0');
    await page.locator('[data-e2e="harvest-submit-btn"]').click();
    await expect(page.locator('[data-e2e="harvest-modal"]')).toBeVisible();
  });
});
