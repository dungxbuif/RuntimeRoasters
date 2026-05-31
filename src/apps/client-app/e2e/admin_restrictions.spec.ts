import { test, expect } from '@playwright/test';

test('Admin restricted from write actions in Farm/Retail', async ({ page }) => {
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

  // Navigate to Dashboard
  await page.waitForURL('**/dashboard**');
  await page.evaluate(() => {
    const modal = document.querySelector('.fixed.inset-0');
    if (modal) modal.remove();
  });

  // 1. Check Harvests Page - "Declare New Harvest" button should NOT be visible
  await page.goto('http://localhost:3000/dashboard/farm-ops/harvests');
  await page.waitForTimeout(2000);
  const declareBtn = page.locator('button:has-text("Declare New Harvest")');
  await expect(declareBtn).not.toBeVisible();

  // 2. Check Orders Page - "Broadcast Order" button should NOT be visible
  // Note: Sidebar link might also be hidden, but let's check the page directly
  await page.goto('http://localhost:3000/dashboard/retail/orders');
  await page.waitForTimeout(2000);
  const orderBtn = page.locator('button:has-text("Broadcast Order")');
  await expect(orderBtn).not.toBeVisible();
});

test('Admin can create Warehouse and Store', async ({ page }) => {
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

    // Navigate to Dashboard
    await page.waitForURL('**/dashboard**');
    
    // Close Bootstrap Modal or any overlays if they appear
    await page.waitForTimeout(3000); 
    await page.evaluate(() => {
      const overlays = document.querySelectorAll('div.fixed.inset-0');
      overlays.forEach(el => el.remove());
    });

    await page.goto('http://localhost:3000/dashboard/resources');
    await page.waitForTimeout(2000);
    await page.evaluate(() => {
      const overlays = document.querySelectorAll('div.fixed.inset-0');
      overlays.forEach(el => el.remove());
    });
    await page.waitForSelector('.material-symbols-outlined', { timeout: 15000 }); // Wait for any icon to load

    // 1. Create Warehouse
    await page.click('button:has-text("Warehouses")');
    await page.click('button:has-text("Create Warehouse")');
    await page.fill('input[placeholder="e.g. South Hub A"]', 'Test Warehouse');
    await page.fill('input[placeholder="WH-001"]', 'WH-TEST-01');
    await page.fill('input[placeholder="5000 Tons"]', '1000 Tons');
    await page.fill('input[placeholder="District 7, HCM City"]', 'Test Address');
    
    // Select manager
    await page.selectOption('select:has-text("Select personnel...")', { index: 1 });
    
    await page.click('button:has-text("Finalize Registration")');
    
    // Wait for modal to disappear (indicates onSuccess finished)
    await expect(page.locator('text=Finalize Registration')).not.toBeVisible({ timeout: 10000 });

    // Check if added to table
    await expect(page.locator('td span:has-text("Test Warehouse")').first()).toBeVisible({ timeout: 10000 });

    // 2. Create Store
    await page.click('button:has-text("Stores")');
    await page.click('button:has-text("Create Store")');
    await page.fill('input[placeholder="e.g. South Hub A"]', 'Test Store');
    await page.fill('input[placeholder="Hanoi"]', 'Hanoi');
    await page.fill('input[placeholder="123 Ly Thai To, Hoan Kiem"]', 'Test Store Address');
    
    // Select manager
    await page.selectOption('select:has-text("Select personnel...")', { index: 1 });
    
    await page.click('button:has-text("Finalize Registration")');
    
    // Check if added to table
    await expect(page.locator('td span:has-text("Test Store")').first()).toBeVisible({ timeout: 10000 });
});
