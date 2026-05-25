import { expect, test } from '@playwright/test';
import { ADMIN_EMAIL, ADMIN_PASS, loginViaUi } from './helpers';

type ApiResponse<T> = {
  status: number;
  body: T;
};

test.describe('Manual live evidence', () => {
  test('admin logs in and triggers harvest plus paid-order event flows', async ({ page, request }) => {
    test.setTimeout(120_000);

    await loginViaUi(page, ADMIN_EMAIL, ADMIN_PASS);

    const auth = await page.evaluate(() => {
      const token = window.localStorage.getItem('rr_access_token');
      if (!token) {
        throw new Error('missing rr_access_token after UI login');
      }
      const [, payload] = token.split('.');
      const claims = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
      return { token, claims };
    });

    const traceparent = `00-${crypto.randomUUID().replace(/-/g, '')}-${crypto.randomUUID().replace(/-/g, '').slice(0, 16)}-01`;
    const api = async <T>(path: string, init: { method?: string; data?: unknown; headers?: Record<string, string> } = {}): Promise<ApiResponse<T>> => {
      const response = await request.fetch(`http://localhost:8081${path}`, {
        method: init.method ?? 'GET',
        data: init.data,
        headers: {
          Authorization: `Bearer ${auth.token}`,
          traceparent,
          ...(init.headers ?? {}),
        },
      });
      const text = await response.text();
      let body: T;
      try {
        body = text ? JSON.parse(text) : ({} as T);
      } catch {
        body = { raw: text } as T;
      }
      if (!response.ok()) {
        throw new Error(`${init.method ?? 'GET'} ${path} failed ${response.status()}: ${text}`);
      }
      return { status: response.status(), body };
    };

    const farmName = `Evidence Farm ${Date.now()}`;
    const farm = await api<{ farm: { id: number } }>('/v1/farms', {
      method: 'POST',
      data: {
          name: farmName,
          location: 'CAU_DAT',
          area: 12.5,
          farm_type: 'ARABICA',
          owner_id: 'admin@runtimeroasters.com',
      },
    });

    const harvest = await api<{ harvest: { id: number } }>('/v1/harvests', {
      method: 'POST',
      data: {
          farm_id: farm.body.farm.id,
          coffee_type: 'ARABICA',
          quantity: 42,
          notes: 'manual live trace evidence',
      },
    });

    const stores = await api<{ stores: Array<{ id: string }> }>('/v1/stores');
    const storeID = stores.body.stores[0]?.id;
    if (!storeID) {
      throw new Error('no retail store seeded');
    }

    const order = await api<{ order: { id: string } }>('/v1/orders', {
      method: 'POST',
      headers: { 'Idempotency-Key': `manual-evidence-${Date.now()}` },
      data: {
          store_id: storeID,
          items: [{ sku: 'SKU-ARABICA-250G', quantity: 1 }],
          total_amount: 19.99,
          payment_method: 'STRIPE',
      },
    });

    const expectedTraceID = traceparent.split('-')[1];
    const readTraceDocument = async () => {
      const response = await request.fetch(`http://localhost:8087/v1/trace/${order.body.order.id}/document`, {
        headers: {
          Authorization: `Bearer ${auth.token}`,
          traceparent,
        },
      });
      if (!response.ok()) {
        return null;
      }
      const json = await response.json();
      return json.document ?? json;
    };

    await expect
      .poll(async () => {
        const doc = await readTraceDocument();
        if (!doc) {
          return null;
        }
        const topics = (doc.events ?? []).map((event: { topic: string }) => event.topic);
        const traceIDs = doc.trace_ids ?? [];
        if (
          traceIDs.includes(expectedTraceID) &&
          topics.includes('retail.order.created') &&
          topics.includes('payment.intent.created') &&
          topics.includes('payment.simulated_completed') &&
          topics.includes('warehouse.stock.reserved')
        ) {
          return doc;
        }
        return null;
      }, { timeout: 60_000, intervals: [1_000, 2_000, 5_000] })
      .not.toBeNull();

    const finalDocument = await readTraceDocument();
    if (!finalDocument) {
      throw new Error('trace document disappeared after polling');
    }
    const topics = (finalDocument.events ?? []).map((event: { topic: string }) => event.topic);
    const traceIDs = finalDocument.trace_ids ?? [];

    const evidence = {
      traceparent,
      expectedTraceID,
      subject: auth.claims.sub,
      role: auth.claims.role ?? auth.claims.ext?.role,
      farmID: farm.body.farm.id,
      harvestID: harvest.body.harvest.id,
      storeID,
      orderID: order.body.order.id,
      traceDocumentFound: true,
      traceIDs,
      topics,
    };

    console.log('MANUAL_LIVE_EVIDENCE', JSON.stringify(evidence));
    expect(evidence.traceparent).toMatch(/^00-[a-f0-9]{32}-[a-f0-9]{16}-01$/);
    expect(evidence.harvestID).toBeTruthy();
    expect(evidence.orderID).toBeTruthy();
    expect(traceIDs).toContain(expectedTraceID);
    expect(traceIDs).toEqual([expectedTraceID]);
    expect(topics).toContain('retail.order.created');
    expect(topics).toContain('payment.intent.created');
    expect(topics).toContain('payment.simulated_completed');
    expect(topics).toContain('warehouse.stock.reserved');
  });
});
