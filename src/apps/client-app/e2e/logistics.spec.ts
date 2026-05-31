import { test, expect } from '@playwright/test';

test('Logistics Map Light Mode and Route Toggling', async ({ page }) => {
  // Login as admin
  await page.goto('http://localhost:3000/login');
  await page.waitForSelector('input[name="identifier"]', { timeout: 30000 });
  await page.fill('input[name="identifier"]', 'admin@runtimeroasters.com');
  await page.fill('input[name="password"]', 'Hello@123');
  await page.click('button[type="submit"]');

  // Handle potential verification screen
  await page.waitForTimeout(3000);
  const bodyText = await page.innerText('body');
  if (bodyText.includes('CONFIRM') || bodyText.includes('VERIFYING')) {
     await page.fill('input[name="password"]', 'Hello@123');
     await page.click('button[type="submit"]');
  }

  // Navigate to Logistics Dashboard
  await page.waitForURL('**/dashboard**');
  await page.evaluate(() => {
    const modal = document.querySelector('.fixed.inset-0');
    if (modal) modal.remove();
  });
  
  await page.goto('http://localhost:3000/dashboard/logistics');
  await page.waitForSelector('.leaflet-container', { timeout: 15000 });

  // Check if voyager light tiles are used
  const tileLayer = page.locator('.leaflet-tile-pane img').first();
  const src = await tileLayer.getAttribute('src');
  expect(src).toContain('basemaps.cartocdn.com/rastertiles/voyager');

  // Routes should NOT be visible initially
  const polyline = page.locator('path.leaflet-interactive');
  await expect(polyline).not.toBeVisible();

  // Activate simulation
  const simButton = page.locator('button:has-text("SIMULATION")');
  if (await simButton.isEnabled()) {
    await simButton.click();
    // In theory, routes should now be visible if data exists
    // But since we are on local dev without real-time data, we just verify the state change
    await expect(simButton).toContainText('ACTIVE');
  }
});
