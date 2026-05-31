import { test, expect } from '@playwright/test';

test('Initial System Seeding via Bootstrap Modal', async ({ page }) => {
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

  // Navigate to Dashboard and wait for Modal
  await page.waitForURL('**/dashboard**');
  
  // The modal might take a second to appear as it checks status
  const seedBtn = page.locator('button:has-text("Initialize Data")');
  await expect(seedBtn).toBeVisible({ timeout: 15000 });
  
  await seedBtn.click();
  
  // Wait for success message
  await expect(page.locator('text=System data initialized successfully!')).toBeVisible({ timeout: 60000 });
  
  console.log('✅ System seeded successfully.');
});
