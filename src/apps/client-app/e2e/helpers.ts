import { expect, Page } from '@playwright/test';

export const ADMIN_EMAIL = 'admin@runtimeroasters.com';
export const ADMIN_PASS = 'Hello@123';
export const MANAGER_EMAIL = 'manager.caudat@runtimeroasters.com';
export const MANAGER_PASS = 'Hello@123';

export async function loginViaUi(page: Page, email = ADMIN_EMAIL, password = ADMIN_PASS) {
  await page.goto('/');
  await page.locator('[data-e2e="public-access-terminal"]').click();
  await expect(page).toHaveURL(/\/login/, { timeout: 30000 });

  for (let attempt = 0; attempt < 3; attempt += 1) {
    const identifier = page.locator('[data-e2e="input-identifier"]');
    const passwordInput = page.locator('[data-e2e="input-password"]');
    await passwordInput.waitFor({ state: 'visible', timeout: 15000 });

    if (await identifier.count()) {
      await identifier.fill(email);
    }
    await passwordInput.fill(password);

    await page.locator('[data-e2e="login-submit"]').click();

    try {
      await expect(page).toHaveURL(/(\/dashboard\/users|\/\?access_token=)/, { timeout: 15000 });
      break;
    } catch {
      if (!(await page.locator('[data-e2e="login-submit"]').count())) {
        break;
      }
    }
  }

  if (/\?access_token=/.test(page.url()) || page.url() === 'http://localhost:3000/') {
    await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 30000 });
  } else {
    await expect(page).toHaveURL(/\/dashboard\/users/, { timeout: 30000 });
  }
  await expect(page.locator('[data-e2e="auth-loader"]')).toBeHidden({ timeout: 30000 });
  await expect(page.locator('[data-e2e="dashboard-header"]')).toBeVisible();
}

export async function logoutViaUi(page: Page) {
  await page.locator('[data-e2e="logout-btn"]').click();
  await expect(page).toHaveURL(/\/login/, { timeout: 30000 });
  await page.context().clearCookies();
  try {
    await page.evaluate(() => {
      localStorage.clear();
      sessionStorage.clear();
    });
  } catch {}
}
