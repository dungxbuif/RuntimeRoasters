import { expect, test } from '@playwright/test';

test.describe('Sprint 1: Infrastructure Foundation & Baseline', () => {
  test('public showcase loads without management shell', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('[data-e2e="public-brand"]')).toBeVisible();
    await expect(page.locator('[data-e2e="public-access-terminal"]')).toBeVisible();
    await expect(page.locator('[data-e2e="public-topology-heading"]')).toBeVisible();
    await expect(page.locator('[data-e2e="dashboard-header"]')).toHaveCount(0);
  });

  test('protected dashboard redirects anonymous users to login', async ({ page }) => {
    await page.goto('/dashboard/users');
    await expect(page).toHaveURL(/\/login/, { timeout: 30000 });
    await expect(page.locator('[data-e2e="login-submit"]')).toBeVisible();
  });
});
