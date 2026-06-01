import { expect, test } from '@playwright/test';

const now = '2026-06-05T00:00:00.000Z';

test.setTimeout(15000);

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    window.localStorage.setItem('rr_access_token', 'e2e-token');
  });
});

test('warehouse fulfillment flow surfaces outbound dispatch notifications to warehouse managers', async ({ page }) => {
  await mockGateway(page, {
    role: 'WAREHOUSE_MGR',
    notifications: [
      {
        id: 'notif-dispatch-1',
        target_role: 'WAREHOUSE_MGR',
        warehouse_id: 'warehouse-1',
        entity_type: 'order',
        entity_id: 'order-1',
        type: 'warehouse.dispatch.requested',
        title: 'Outbound dispatch requested',
        message: 'A paid order has reserved stock and is ready for delivery dispatch.',
        severity: 'info',
        status: 'UNREAD',
        created_at: now,
      },
    ],
  });

  await page.goto('/dashboard/warehouse', { waitUntil: 'domcontentloaded', timeout: 10000 });

  await expect(page.getByRole('heading', { name: /Warehouse Operations/i })).toBeVisible();
  await expect(page.getByText('A paid order has reserved stock and is ready for delivery dispatch.')).toBeVisible();
});

test('store delivery flow surfaces incoming delivery notifications to store managers', async ({ page }) => {
  await mockGateway(page, {
    role: 'STORE_MGR',
    notifications: [
      {
        id: 'notif-delivery-1',
        target_role: 'STORE_MGR',
        store_id: 'store-1',
        shipment_id: 'shipment-1',
        entity_type: 'shipment',
        entity_id: 'shipment-1',
        type: 'logistics.delivery.completed',
        title: 'Delivery status updated',
        message: 'A retail delivery status changed.',
        severity: 'warning',
        status: 'UNREAD',
        created_at: now,
      },
    ],
  });

  await page.goto('/dashboard/store', { waitUntil: 'domcontentloaded', timeout: 10000 });

  await expect(page.getByRole('heading', { name: /Retail Dashboard/i })).toBeVisible();
  await expect(page.getByText('A retail delivery status changed.')).toBeVisible();
});

test('role-scoped notification flow hides warehouse operations from store managers', async ({ page }) => {
  await mockGateway(page, {
    role: 'STORE_MGR',
    notifications: [],
  });

  await page.goto('/dashboard/store', { waitUntil: 'domcontentloaded', timeout: 10000 });

  await expect(page.getByRole('heading', { name: /Retail Dashboard/i })).toBeVisible();
  await expect(page.getByText('A paid order has reserved stock and is ready for delivery dispatch.')).toHaveCount(0);
  await expect(page.getByText('No notifications')).toBeVisible();
});

async function mockGateway(page: import('@playwright/test').Page, options: {
  role: 'WAREHOUSE_MGR' | 'STORE_MGR';
  notifications: unknown[];
}) {
  await page.route('**/v1/auth/me', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        user: {
          id: `${options.role.toLowerCase()}-1`,
          email: `${options.role.toLowerCase()}@runtimeroasters.com`,
          role: options.role,
          name: options.role,
        },
      }),
    });
  });

  await page.route('**/v1/auth/snapshot/frontend', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        policies: [
          `p, ${options.role}, /v1/realtime/notifications, read`,
          `p, ${options.role}, /v1/realtime/notifications/.*, write`,
          `p, ${options.role}, /v1/warehouse/.*, read|write`,
          `p, ${options.role}, /v1/logistics/.*, read|write`,
          `p, ${options.role}, /v1/orders, read|write`,
          `p, ${options.role}, /v1/stores, read`,
        ],
      }),
    });
  });

  await page.route('**/v1/realtime/notifications', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ notifications: options.notifications }),
    });
  });

  await page.route('**/v1/warehouse/pickup-requests', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ pickup_requests: [] }) });
  });
  await page.route('**/v1/warehouse/dispatch-requests', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ dispatch_requests: [] }) });
  });
  await page.route('**/v1/warehouse/batches', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ batches: [] }) });
  });
  await page.route('**/v1/warehouse/inventory', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ inventory: [] }) });
  });
  await page.route('**/v1/logistics/drivers', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ drivers: [] }) });
  });
  await page.route('**/v1/logistics/vehicles/available', async route => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ vehicles: [] }) });
  });
  await page.route('**/v1/stores', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ stores: [{ id: 'store-1', name: 'Demo Store', location: 'Saigon', is_active: true }] }),
    });
  });
  await page.route('**/v1/orders', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        orders: [{ id: 'order-1', store_id: 'store-1', status: 'IN_TRANSIT', total_amount: 150000, created_at: now, updated_at: now }],
      }),
    });
  });
}
