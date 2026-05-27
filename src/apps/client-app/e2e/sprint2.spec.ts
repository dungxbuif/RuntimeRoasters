import { expect, test } from '@playwright/test';
import { loginViaUi, logoutViaUi } from './helpers';

test.describe('Sprint 2: Identity & Access Control', () => {
  test('admin can login and logout through OIDC flow', async ({ page }) => {
    await loginViaUi(page);
    await page.goto('/dashboard');

    await expect(page.getByRole('heading', { name: /Intelligence Hub/i })).toBeVisible();
    await expect(page.getByText('Provenance Trace')).toBeVisible();

    await logoutViaUi(page);
    await page.goto('/dashboard/users');
    await expect(page).toHaveURL(/\/login/, { timeout: 30000 });
  });
});
