import { test, expect } from '@playwright/test';

test('Architecture Topology UI Layout', async ({ page }) => {
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

  // Navigate to Dashboard and close Modal if exists
  await page.waitForURL('**/dashboard**');
  await page.evaluate(() => {
    const modal = document.querySelector('.fixed.inset-0');
    if (modal) modal.remove();
  });

  // Navigate to Topology Mesh
  await page.goto('http://localhost:3000/dashboard/topology-mesh');

  // Wait for react flow to render
  await page.waitForSelector('.react-flow', { timeout: 15000 });

  // Check if group nodes are present
  const frontendGroup = page.locator('.react-flow__node-groupNode:has-text("Client / Edge")');
  await expect(frontendGroup).toBeVisible();

  const servicesGroup = page.locator('.react-flow__node-groupNode:has-text("Microservices")');
  await expect(servicesGroup).toBeVisible();
});
