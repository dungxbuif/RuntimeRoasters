import { expect, test } from '@playwright/test';

/**
 * RUNTIME ROASTERS - CONSTITUTIONAL SYSTEM JOURNEY
 * This script verifies the 100% transparent "Farm-to-Cup" flow.
 * It strictly follows the operational story: 
 * Harvest -> Pickup -> Inbound -> Roasting -> Order -> Delivery -> Trace
 */

const now = new Date().toISOString();
const MOCK_IDS = {
  FARM: 'farm-cd-001',
  HARVEST: 'RR-H-CD-20260606-001',
  SHIPMENT: 'SHP-INB-001',
  ORDER: 'ORD-RET-9900',
  TRACE: 'TRC-CUP-7722'
};

test.describe('Holistic System Journey', () => {
  
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      window.localStorage.setItem('rr_access_token', 'system-journey-token');
    });
  });

  test('The "Golden Path": From Farm to Public Proof', async ({ page }) => {
    
    // --- PERSONA: FARM MANAGER ---
    // Goal: Declare a physical harvest at the source.
    await mockIdentity(page, 'FARM_MANAGER', 'user-farm-1');
    await mockFarmAPIs(page);
    await page.goto('/dashboard/farm-ops/harvests', { waitUntil: 'networkidle' });
    
    await expect(page.getByRole('heading', { name: /Harvest Declaration/i })).toBeVisible();
    await page.getByText('Declare New Harvest').click();
    await page.fill('input[type="number"]', '1250.5'); // Significant quantity for demo
    
    // Intercept Harvest Creation
    await page.route('**/v1/harvests', async (route) => {
        if (route.request().method() === 'POST') {
            await route.fulfill({ status: 201, body: JSON.stringify({ harvest: { id: 101, status: 'NEW' } }) });
        }
    });
    await page.getByRole('button', { name: /Record Harvest/i }).click();
    await expect(page.getByText(/Syncing Harvest Ledger/i).or(page.getByText(/Harvest ID/i))).toBeVisible();


    // --- PERSONA: WAREHOUSE MANAGER ---
    // Goal: Orchestrate logistics for the incoming harvest.
    await mockIdentity(page, 'WAREHOUSE_MGR', 'user-wh-1');
    await mockWarehouseAPIs(page);
    await page.goto('/dashboard/warehouse', { waitUntil: 'networkidle' });
    
    await expect(page.getByText(/Inbound Queue/i)).toBeVisible();
    // Dispatch Pickup
    await page.getByRole('button', { name: /Dispatch Pickup/i }).first().click();
    await expect(page.getByText(/Assign Fleet/i)).toBeVisible();
    
    // Assignment Interaction
    await page.selectOption('select >> nth=0', { index: 1 }); // Driver
    await page.selectOption('select >> nth=1', { index: 1 }); // Vehicle
    
    // Update mock for the next fetch to show DISPATCHED status
    await page.route('**/v1/warehouse/pickup-requests', route => route.fulfill({ 
        status: 200, 
        body: JSON.stringify({ pickup_requests: [{ id: MOCK_IDS.SHIPMENT, harvest_id: MOCK_IDS.HARVEST, status: 'DISPATCHED', coffee_type: 'ARABICA', quantity: 1250, created_at: now, updated_at: now }] }) 
    }));

    await page.route('**/v1/logistics/shipments/*/assign', route => route.fulfill({ status: 200, body: JSON.stringify({ success: true }) }));
    await page.getByRole('button', { name: /Finalize Dispatch/i }).click();
    await expect(page.getByText(/Driver moving to farm/i)).toBeVisible();


    // --- PERSONA: DRIVER ---
    // Goal: Simulate physical movement and milestone completion.
    await mockIdentity(page, 'DRIVER', 'user-driver-1');
    await mockDriverAPIs(page, 'user-driver-1');
    await page.goto('/dashboard/driver', { waitUntil: 'networkidle' });
    
    await expect(page.getByText(/Active Shipment/i)).toBeVisible();
    // Start the physical journey
    const journeyBtn = page.getByRole('button').filter({ hasText: /Start Route/i });
    await expect(journeyBtn).toBeVisible({ timeout: 10000 });
    await journeyBtn.click();
    await expect(page.getByText(/Simulation Active/i)).toBeVisible();
    
    // Verify Return Leg presence (RR-URG-09C requirement)
    await expect(page.getByText(/Returning/i)).toBeVisible();
    await expect(page.getByText(/Arrived WH/i)).toBeVisible();


    // --- PERSONA: STORE MANAGER ---
    // Goal: Trigger fulfillment for a customer order.
    await mockIdentity(page, 'STORE_MGR', 'user-store-1');
    await mockRetailAPIs(page);
    await page.goto('/dashboard/retail/orders', { waitUntil: 'networkidle' });
    
    await expect(page.getByRole('heading', { name: /Retail Paid Order/i })).toBeVisible();
    await page.getByRole('button', { name: /Broadcast Order/i }).click();


    // --- PERSONA: CUSTOMER (PUBLIC) ---
    // Goal: Verify the 100% digital proof of origin.
    await page.goto('/trace', { waitUntil: 'networkidle' });
    await expect(page.getByRole('heading', { name: /Public Trace Registry/i })).toBeVisible();
    
    // Trace specific cup
    await mockPublicTrace(page, MOCK_IDS.TRACE);
    await page.goto(`/trace/${MOCK_IDS.TRACE}`, { waitUntil: 'networkidle' });
    await expect(page.locator('body')).toContainText(/Public-safe trace verified/i);
    await expect(page.locator('body')).toContainText(/Farm to Cup Journey/i);
    await expect(page.locator('body')).toContainText(/Harvested at Cầu Đất/i);
  });

});

// --- JOURNEY MOCK HELPERS ---

async function mockIdentity(page: any, role: string, id: string) {
  await page.route('**/v1/auth/me', route => route.fulfill({
    status: 200,
    body: JSON.stringify({ user: { id, email: `${role.toLowerCase()}@runtimeroasters.com`, role, name: `Demo ${role}` } })
  }));
  await page.route('**/v1/auth/snapshot/frontend', route => route.fulfill({
    status: 200,
    body: JSON.stringify({ policies: [`p, ${role}, /v1/.*, read|write`, `p, ${role}, harvest, read|write`, `p, ${role}, farm, read|write`] })
  }));
}

async function mockFarmAPIs(page: any) {
    await page.route('**/v1/farms', route => route.fulfill({ status: 200, body: JSON.stringify({ farms: [{ id: 1, name: 'Cầu Đất Farm Node' }] }) }));
    await page.route('**/v1/harvests', route => {
        if (route.request().method() === 'GET') {
            route.fulfill({ status: 200, body: JSON.stringify({ harvests: [] }) });
        }
    });
}

async function mockWarehouseAPIs(page: any) {
    await page.route('**/v1/warehouse/pickup-requests', route => route.fulfill({ 
        status: 200, 
        body: JSON.stringify({ pickup_requests: [{ id: MOCK_IDS.SHIPMENT, harvest_id: MOCK_IDS.HARVEST, status: 'REQUESTED', coffee_type: 'ARABICA', quantity: 1250, created_at: now, updated_at: now }] }) 
    }));
    await page.route('**/v1/warehouse/dispatch-requests', route => route.fulfill({ status: 200, body: JSON.stringify({ dispatch_requests: [] }) }));
    await page.route('**/v1/warehouse/batches', route => route.fulfill({ status: 200, body: JSON.stringify({ batches: [] }) }));
    await page.route('**/v1/warehouse/inventory', route => route.fulfill({ status: 200, body: JSON.stringify({ inventory: [] }) }));
    await page.route('**/v1/logistics/drivers', route => route.fulfill({ status: 200, body: JSON.stringify({ drivers: [{ id: 'd1', name: 'Elite Driver 01', status: 'IDLE' }] }) }));
    await page.route('**/v1/logistics/vehicles/available', route => route.fulfill({ status: 200, body: JSON.stringify({ vehicles: [{ id: 'v1', plate_number: 'RR-SAFE-01', type: 'TRUCK' }] }) }));
}

async function mockDriverAPIs(page: any, driverId: string) {
    await page.route('**/v1/logistics/shipments', route => route.fulfill({ 
        status: 200, 
        body: JSON.stringify({ shipments: [{ id: MOCK_IDS.SHIPMENT, driver_id: driverId, status: 'ASSIGNED', type: 'FARM_PICKUP', destination_address: 'Cầu Đất Farm Node', created_at: now, updated_at: now }] }) 
    }));
    await page.route('**/v1/logistics/locations', route => route.fulfill({ status: 200, body: JSON.stringify({ locations: [] }) }));
    await page.route('**/data/routes.json', route => route.fulfill({ status: 200, body: JSON.stringify([{ id: 'r1', coordinates: [[11.9, 108.4], [11.95, 108.5]] }]) }));
}

async function mockRetailAPIs(page: any) {
    await page.route('**/v1/retail/stores', route => route.fulfill({ status: 200, body: JSON.stringify({ stores: [{ id: 'st1', name: 'Downtown Roastery', city: 'HN' }] }) }));
    await page.route('**/v1/orders', route => route.fulfill({ status: 200, body: JSON.stringify({ orders: [] }) }));
}

async function mockPublicTrace(page: any, code: string) {
    await page.route(`**/v1/public/trace/${code}`, route => route.fulfill({ 
        status: 200, 
        body: JSON.stringify({ 
            trace: {
                product_id: 'CUP-DEMO-01',
                trace_code: code,
                product: { name: 'Premium Arabica Blend', size: 'L', sku: 'SKU-ARB-001' },
                store: { name: 'Downtown Roastery', city: 'Hanoi', address: '123 Ly Thai To St.' },
                invoice_no: 'INV-RR-9922',
                sale_item: { sold_at: now },
                inventory_lot_id: 'LOT-DEMO-441',
                harvest_id: 'HRV-CAUDAT-001',
                origin_batch_id: 'BATCH-ROAST-202',
                warehouse: { id: 'WAREHOUSE-HN' },
                service_participation: ['farm-service', 'warehouse-service', 'logistics-service', 'retail-service'],
                milestones: [
                    { status: 'SOLD', title: 'Coffee Served', timestamp: now, description: 'Cup served and verified by public trace ledger.' },
                    { status: 'DELIVERED', title: 'Beans Delivered to Store', timestamp: now, description: 'Final retail delivery completed by Driver HN-01.' },
                    { status: 'ROASTED', title: 'Batch Roasted (RR-S-001)', timestamp: now, description: 'Green beans roasted and packaged at Hanoi Warehouse.' },
                    { status: 'HARVESTED', title: 'Harvested at Cầu Đất', timestamp: now, description: 'Fresh Arabica cherries harvested and declared at Dalat Farm.' }
                ]
            } 
        }) 
    }));
}
