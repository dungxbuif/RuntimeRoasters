import { expect, test } from '@playwright/test';

test('public coffee sale creates per-cup identity, invoice QR, and trace page', async ({ page }) => {
  let submittedProductID = '';
  await page.route('**/v1/public/retail/stores', route => route.fulfill({ json: { stores: [{ id: 'store-1', code: 'STORE-HN-HK', name: 'Hoan Kiem', city: 'Hanoi', address: '2 Ly Thai To', status: 'ACTIVE' }] } }));
  await page.route('**/v1/public/retail/stores/store-1/menu-items', route => route.fulfill({ json: { menus: [{ id: 'menu-1', code: 'DEFAULT', name: 'Runtime Coffee Menu', description: 'Traceable coffee', items: [{ id: 'item-1', menu_id: 'menu-1', code: 'ESP-S', product_group_code: 'ESP', category_code: 'ESPRESSO', size: 'S', name: 'Espresso', description: 'Double shot', stock_sku: 'BEAN-ARABICA', coffee_type: 'ARABICA', price: 45000, available: true, available_units: 3 }] }] } }));
  await page.route('**/v1/public/retail/sales', async route => {
    const request = route.request().postDataJSON();
    expect(request.items).toHaveLength(1);
    expect(request.items[0].menu_item_id).toBe('item-1');
    expect(request.items[0].product_id).toMatch(/^RR-CUP-HNHK-/);
    submittedProductID = request.items[0].product_id;
    await route.fulfill({ status: 201, json: { sale: { id: 'sale-1', store_id: 'store-1', invoice_no: 'INV-001', total_amount: 45000, sold_at: '2026-06-06T12:00:00Z', items: [{ id: 'sale-item-1', product_id: submittedProductID, trace_code: submittedProductID, product_name: 'Espresso', size: 'S', unit_price: 45000, public_url: `/trace/${submittedProductID}`, sold_at: '2026-06-06T12:00:00Z' }] } } });
  });
  await page.route('**/v1/public/trace/**', route => route.fulfill({ json: { trace: { product_id: submittedProductID, trace_code: submittedProductID, sale_item: { id: 'sale-item-1', sold_at: '2026-06-06T12:00:00Z' }, invoice_no: 'INV-001', product: { menu_item_id: 'item-1', name: 'Espresso', sku: 'BEAN-ARABICA', size: 'S' }, store: { id: 'store-1', name: 'Hoan Kiem', city: 'Hanoi', address: '2 Ly Thai To' }, farm: { harvest_id: 'harvest-1' }, warehouse: { id: 'warehouse-1' }, origin_batch_id: 'batch-1', harvest_id: 'harvest-1', inventory_lot_id: 'lot-1', milestones: [{ status: 'ORIGIN_LINKED', title: 'Origin lineage linked', description: 'Harvest harvest-1, batch batch-1', timestamp: '2026-06-05T12:00:00Z', service: 'farm-service' }, { status: 'SERVED', title: 'Coffee served', description: 'Espresso at Hoan Kiem', timestamp: '2026-06-06T12:00:00Z', service: 'retail-service' }], service_participation: ['farm-service', 'retail-service'], trace_ids: ['trace-1'], public_safe: true } } }));

  await page.goto('/coffee');
  await expect(page.getByRole('heading', { name: 'Coffee menu' })).toBeVisible();
  await page.screenshot({ path: '../../../docs/work/tickets/sprint10/evidence/rr-urg-07/public-kiosk.png', fullPage: true });
  await page.getByRole('button', { name: 'Add Espresso' }).click();
  await page.getByRole('button', { name: 'Buy coffee' }).click();
  await expect(page.getByRole('dialog', { name: 'Invoice' })).toBeVisible();
  await expect(page.getByText('INV-001')).toBeVisible();
  await expect(page.getByRole('dialog').getByRole('img')).toBeVisible();
  await page.screenshot({ path: '../../../docs/work/tickets/sprint10/evidence/rr-urg-07/invoice-qr.png', fullPage: true });
  await page.getByRole('link', { name: 'Open trace' }).click();
  await expect(page).toHaveURL(/\/trace\/RR-CUP-HNHK-/);
  await expect(page.getByRole('heading', { name: 'Farm to Cup Journey' })).toBeVisible();
  await expect(page.getByText('Origin lineage linked')).toBeVisible();
  await expect(page.getByText('Coffee served')).toBeVisible();
  await page.screenshot({ path: '../../../docs/work/tickets/sprint10/evidence/rr-urg-07/public-trace.png', fullPage: true });
});

test('unknown sold cup renders a public not-found state', async ({ page }) => {
  await page.route('**/v1/public/trace/missing-cup', route => route.fulfill({ status: 404, json: { message: 'sold item not found' } }));
  await page.goto('/trace/missing-cup');
  await expect(page.getByRole('heading', { name: 'Cup not found' })).toBeVisible();
});
