import { expect, test } from '@playwright/test';

/**
 * SUPPLY CHAIN JOURNEY: FARM TO CUP
 * This test verifies the entire canonical flow of Runtime Roasters.
 * Roles: FARM_MANAGER -> WAREHOUSE_MGR -> DRIVER -> STORE_MGR -> CUSTOMER
 */

const MOCK_IDS = {
  FARM: 'farm-123',
  HARVEST: 'RR-H-CD-20260606-001',
  SHIPMENT: 'ship-inbound-001',
  ORDER: 'ORD-PAID-7788',
  TRACE: 'TRC-CUP-9900'
};

const now = new Date().toISOString();

test.describe('Supply Chain Journey: Farm to Cup', () => {
  
  test.beforeEach(async ({ page }) => {
    // Standard setup for all tests
    await page.addInitScript(() => {
      window.localStorage.setItem('rr_access_token', 'journey-token');
    });
  });

  test('Complete Supply Chain Lifecycle', async ({ page }) => {
    // --- STAGE 1: FARM HARVEST ---
    await mockRole(page, 'FARM_MANAGER');
    await page.goto('/dashboard/farm-ops/harvests', { waitUntil: 'networkidle' });
    await expect(page.getByRole('heading', { name: /Harvest Declaration/i })).toBeVisible();

    // Interaction: Declare Harvest
    await page.getByText('Declare New Harvest').click();
    await page.fill('input[type="number"]', '500');
    // Mock the POST call for harvest creation
    await page.route('**/v1/harvests', async (route: any) => {
        if (route.request().method() === 'POST') {
            await route.fulfill({ status: 201, contentType: 'application/json', body: JSON.stringify({ harvest: { id: 101 } }) });
        }
    });
    await page.getByRole('button', { name: /Record Harvest/i }).click();
    await expect(page.getByText('Harvest Declaration')).toBeVisible();

    // --- STAGE 2: INBOUND DISPATCH ---
    await mockRole(page, 'WAREHOUSE_MGR');
    await mockWarehouseData(page, {
        pickup_requests: [{
            id: MOCK_IDS.SHIPMENT,
            harvest_id: MOCK_IDS.HARVEST,
            status: 'REQUESTED',
            coffee_type: 'ARABICA',
            quantity: 500,
            created_at: now,
            updated_at: now
        }]
    });
    await page.goto('/dashboard/warehouse', { waitUntil: 'networkidle' });
    await expect(page.getByText('Inbound Queue')).toBeVisible();

    // Interaction: Dispatch Pickup
    await page.getByRole('button', { name: /Dispatch Pickup/i }).first().click();
    await expect(page.getByText('Assign Fleet')).toBeVisible();
    // Fill assignment
    await mockLogisticsInventory(page);
    await page.selectOption('select >> nth=0', { index: 1 }); 
    await page.selectOption('select >> nth=1', { index: 1 }); 
    await page.route('**/v1/logistics/shipments/*/assign', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ shipment: { id: MOCK_IDS.SHIPMENT } }) });
    });
    await page.getByRole('button', { name: /Finalize Dispatch/i }).click();

    // --- STAGE 3: DRIVER PICKUP SIMULATION ---
    await mockRole(page, 'DRIVER');
    await mockDriverData(page, {
        shipments: [{
            id: MOCK_IDS.SHIPMENT,
            driver_id: 'user-DRIVER',
            status: 'ASSIGNED',
            type: 'FARM_PICKUP',
            destination_address: 'Cầu Đất Farm Node',
            created_at: now,
            updated_at: now
        }]
    });
    await page.goto('/dashboard/driver', { waitUntil: 'networkidle' });

    // Wait for the active shipment to be recognized by the filtering logic
    await expect(page.getByText('Active Shipment')).toBeVisible({ timeout: 15000 });

    // Interaction: Start Journey
    const startBtn = page.locator('button').filter({ hasText: 'Start Journey' });
    await expect(startBtn).toBeVisible({ timeout: 10000 });
    await startBtn.click();
    await expect(page.getByText('Simulation Active')).toBeVisible();

    // --- STAGE 4: RETAIL FULFILLMENT SAGA ---
    await mockRole(page, 'STORE_MGR');
    await mockRetailData(page, {
        orders: [{
            id: MOCK_IDS.ORDER,
            status: 'PAYMENT_COMPLETED',
            total_amount: 150000,
            created_at: now,
            updated_at: now
        }]
    });
    await page.goto('/dashboard/retail', { waitUntil: 'networkidle' });
    await expect(page.getByRole('heading', { name: /Retail Paid Orders/i })).toBeVisible();
    await expect(page.getByText('Paid')).toBeVisible();

    // --- STAGE 5: PUBLIC PROOF (CUSTOMER) ---
    await page.goto('/trace', { waitUntil: 'networkidle' });
    await expect(page.getByRole('heading', { name: /Public Trace Registry/i })).toBeVisible();
    await mockTraceDetail(page, MOCK_IDS.TRACE);
    await page.goto(`/trace/${MOCK_IDS.TRACE}`, { waitUntil: 'networkidle' });
    await expect(page.getByText('Cup Lineage Verified')).toBeVisible();
  });

  });

  async function mockLogisticsInventory(page: any) {
    await page.route('**/v1/logistics/drivers', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ drivers: [{ id: 'd1', name: 'John Doe', status: 'IDLE' }] }) });
    });
    await page.route('**/v1/logistics/vehicles/available', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ vehicles: [{ id: 'v1', plate_number: '51A-12345', type: 'TRUCK' }] }) });
    });
  }


// --- HELPERS ---

async function mockRole(page: any, role: string) {
  await page.route('**/v1/auth/me', async (route: any) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        user: { id: `user-${role}`, email: `${role.toLowerCase()}@test.com`, role: role, name: `Demo ${role}` }
      }),
    });
  });

  await page.route('**/v1/auth/snapshot/frontend', async (route: any) => {
    // Provide a full set of permissions for the test
    const policies = [
      `p, ${role}, harvest, read|write`,
      `p, ${role}, farm, read|write`,
      `p, ${role}, /v1/.*, read|write`,
      `p, ${role}, /runtime.system.v1.SystemService/*, .*`
    ];
    
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ policies }),
    });
  });
  
  // Also mock farms for FARM_MANAGER
  if (role === 'FARM_MANAGER') {
      await page.route('**/v1/farms', async (route: any) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ farms: [{ id: 1, name: 'Cầu Đất Farm' }] }),
        });
      });
      await page.route('**/v1/farms/1/harvests', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ harvests: [] }) });
      });
  }
}

async function mockWarehouseData(page: any, data: any) {
    await page.route('**/v1/warehouse/pickup-requests', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(data) });
    });
    await page.route('**/v1/warehouse/dispatch-requests', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ dispatch_requests: [] }) });
    });
    await page.route('**/v1/warehouse/batches', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ batches: [] }) });
    });
    await page.route('**/v1/warehouse/inventory', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ inventory: [] }) });
    });
}

async function mockDriverData(page: any, data: any) {
    await page.route('**/v1/logistics/shipments', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(data) });
    });
    await page.route('**/v1/logistics/locations', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ locations: [] }) });
    });
    // Mock the static routes JSON file
    await page.route('**/data/routes.json', async (route: any) => {
        await route.fulfill({ 
          status: 200, 
          contentType: 'application/json', 
          body: JSON.stringify([{ id: 'r1', coordinates: [[10, 106], [10.1, 106.1]] }]) 
        });
    });
}

async function mockRetailData(page: any, data: any) {
    await page.route('**/v1/retail/orders', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(data) });
    });
    await page.route('**/v1/retail/stores', async (route: any) => {
        await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ stores: [] }) });
    });
}

async function mockTraceDetail(page: any, code: string) {
    await page.route(`**/v1/public/trace/${code}`, async (route: any) => {
        await route.fulfill({ 
            status: 200, 
            contentType: 'application/json', 
            body: JSON.stringify({ 
                trace: {
                    product_id: 'CUP-1',
                    trace_code: code,
                    product: { name: 'Espresso', size: 'M' },
                    store: { name: 'District 1 Store' },
                    milestones: [
                        { status: 'SOLD', title: 'Coffee Served', timestamp: new Date().toISOString() },
                        { status: 'HARVESTED', title: 'Harvested at Cầu Đất', timestamp: new Date().toISOString() }
                    ]
                } 
            }) 
        });
    });
}
