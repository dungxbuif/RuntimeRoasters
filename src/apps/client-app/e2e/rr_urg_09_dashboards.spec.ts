import { expect, test } from '@playwright/test';

const now = new Date().toISOString();

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    window.localStorage.setItem('rr_access_token', 'e2e-token-09');
  });
});

test('public entry points are accessible without login', async ({ page }) => {
  // 1. Root Topology
  await page.goto('/');
  await expect(page.locator('h1')).toContainText(/Architecture Topology/i);
  await expect(page.getByRole('link', { name: /Trace Registry/i })).toBeVisible();

  // 2. Trace Registry
  await page.goto('/trace');
  await expect(page.locator('h1')).toContainText(/Public Trace Registry/i);
});

test('driver simulation UI includes return leg controls', async ({ page }) => {
  await mockDriverSession(page);
  await page.goto('/dashboard/driver');

  await expect(page.getByRole('heading', { name: /Driver Client/i })).toBeVisible();
  
  // We should see assigned shipment
  await expect(page.getByText(/Active Shipment/i)).toBeVisible();
  
  // Check milestones for return leg
  await expect(page.getByText(/Start Return/i)).toBeVisible();
  await expect(page.getByText(/At Base/i).or(page.getByText(/Arrived WH/i))).toBeVisible();
});

test('retail dashboard uses paid order terminology', async ({ page }) => {
  await mockRetailSession(page);
  await page.goto('/dashboard/retail/orders');

  await expect(page.getByRole('heading', { name: /Retail Paid Order/i })).toBeVisible();
  await expect(page.getByText(/Fulfillment Store Node/i)).toBeVisible();
});

async function mockDriverSession(page: import('@playwright/test').Page) {
  await page.route('**/v1/auth/me', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        user: { id: 'driver-1', email: 'driver@runtimeroasters.com', role: 'DRIVER', name: 'Demo Driver' },
      }),
    });
  });

  await page.route('**/v1/auth/snapshot/frontend', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ policies: [`p, DRIVER, /v1/logistics/.*, read|write`] }),
    });
  });

  await page.route('**/v1/logistics/shipments', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        shipments: [{
          id: 'ship-09',
          order_id: 'ord-09',
          driver_id: 'driver-1',
          status: 'ASSIGNED',
          type: 'RETAIL_DELIVERY',
          destination_address: 'Test Store Address',
          created_at: now,
          updated_at: now
        }]
      }),
    });
  });

  await page.route('**/v1/logistics/locations', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ locations: [] }) });
  });

  await page.route('**/v1/logistics/routes', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ routes: [{ id: 'rt-1', coordinates: [[10, 106], [10.1, 106.1]] }] }),
    });
  });
}

async function mockRetailSession(page: import('@playwright/test').Page) {
  await page.route('**/v1/auth/me', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        user: { id: 'mgr-1', email: 'mgr@runtimeroasters.com', role: 'STORE_MGR', name: 'Store Manager' },
      }),
    });
  });

  await page.route('**/v1/auth/snapshot/frontend', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ policies: [`p, STORE_MGR, /v1/retail/.*, read|write`] }),
    });
  });

  await page.route('**/v1/retail/stores', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ stores: [{ id: 'st-1', name: 'Store 1', location: 'HN' }] }),
    });
  });
}
