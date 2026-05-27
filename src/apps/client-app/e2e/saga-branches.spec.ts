import { APIRequestContext, expect, test } from '@playwright/test';
import { ADMIN_EMAIL, ADMIN_PASS, loginViaUi } from './helpers';

// Helper for HMAC SHA-256 for Stripe Webhook Signature
async function hmacSha256Hex(secret: string, data: string): Promise<string> {
  const enc = new TextEncoder();
  const key = await crypto.subtle.importKey('raw', enc.encode(secret), { name: 'HMAC', hash: 'SHA-256' }, false, ['sign']);
  const signature = await crypto.subtle.sign('HMAC', key, enc.encode(data));
  return Array.from(new Uint8Array(signature))
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}

type ApiResponse<T> = {
  status: number;
  body: T;
};

type ApiInit = {
  method?: string;
  data?: unknown;
  headers?: Record<string, string>;
};

type EntityResponse<K extends string> = Record<K, { id: string; status?: string; payment_intent?: string }>;
type PaymentByOrderResponse = {
  payment: {
    id: string;
    order_id: string;
    provider_ref: string;
    status: string;
    store_id: string;
  };
};

const SEEDED_STORE_ID = '11111111-1111-1111-1111-111111111101';

test.describe('100% SAGA Business Domain Automation', () => {
  let auth: { token: string };
  let traceparent: string;
  let api: <T>(path: string, init?: ApiInit) => Promise<ApiResponse<T>>;
  let webhookKey: string;
  let webhookUrl: string;

  test.beforeEach(async ({ page, request }) => {
    test.setTimeout(120_000);
    await loginViaUi(page, ADMIN_EMAIL, ADMIN_PASS);

    auth = await page.evaluate(() => {
      const token = window.localStorage.getItem('rr_access_token');
      if (!token) throw new Error('missing rr_access_token');
      return { token };
    });

    traceparent = `00-${crypto.randomUUID().replace(/-/g, '')}-${crypto.randomUUID().replace(/-/g, '').slice(0, 16)}-01`;
    api = async <T>(path: string, init: ApiInit = {}): Promise<ApiResponse<T>> => {
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
      if (!response.ok()) throw new Error(`${init.method ?? 'GET'} ${path} failed ${response.status()}: ${text}`);
      const body: T = text ? JSON.parse(text) : ({} as T);
      return { status: response.status(), body };
    };

    const webhookRes = await api<{ signing_secret: string; webhook_url: string }>('/v1/payments/demo/stripe-webhook-key');
    webhookKey = webhookRes.body.signing_secret;
    webhookUrl = webhookRes.body.webhook_url;
  });

  const triggerWebhook = async (request: APIRequestContext, type: string, providerRef: string, orderId: string, paymentId: string, storeId: string) => {
    const webhookPayload = JSON.stringify({
      id: `evt_rr_evidence_${crypto.randomUUID()}`,
      object: 'event',
      type: type,
      livemode: false,
      created: Math.floor(Date.now() / 1000),
      data: {
        object: {
          id: providerRef,
          object: 'payment_intent',
          amount: 1999,
          currency: 'usd',
          status: type === 'payment_intent.succeeded' ? 'succeeded' : 'failed',
          metadata: { order_id: orderId, payment_id: paymentId, store_id: storeId },
        },
      },
    });
    const webhookTimestamp = Math.floor(Date.now() / 1000).toString();
    const webhookSignature = await hmacSha256Hex(webhookKey, `${webhookTimestamp}.${webhookPayload}`);
    await request.fetch(`http://localhost:8081${webhookUrl}`, {
      method: 'POST',
      data: webhookPayload,
      headers: {
        Authorization: `Bearer ${auth.token}`,
        'Content-Type': 'application/json',
        'Stripe-Signature': `t=${webhookTimestamp},v1=${webhookSignature}`,
        traceparent,
      },
    });
  };

  test('Branch 1: SAGA Rollback - Payment Failed', async ({ request }) => {
    test.setTimeout(120_000);
    // 1. Create Order against canonical seeded store.
    const order = await api<EntityResponse<'order'>>('/v1/orders', {
      method: 'POST',
      data: { store_id: SEEDED_STORE_ID, total_amount: 1999, items: [{ sku: 'COFFEE-ARABICA-001', quantity: 1, unit_price: 1999 }] }
    });

    // Extract Provider Ref
    const paymentRef: { current?: PaymentByOrderResponse['payment'] } = {};
    await expect.poll(async () => {
      try {
        const res = await api<PaymentByOrderResponse>(`/v1/payments/orders/${order.body.order.id}`);
        paymentRef.current = res.body.payment;
        return paymentRef.current.provider_ref;
      } catch {
        return '';
      }
    }, { timeout: 30000 }).not.toBe('');
    if (!paymentRef.current) throw new Error('missing payment');
    const payment = paymentRef.current;
    const providerRef = payment.provider_ref;

    // 2. Trigger Payment Failed Webhook
    await triggerWebhook(request, 'payment_intent.payment_failed', providerRef, order.body.order.id, payment.id, SEEDED_STORE_ID);

    // 3. Verify Order is REJECTED
    await expect.poll(async () => {
      const res = await api<EntityResponse<'order'>>(`/v1/orders/${order.body.order.id}`);
      return res.body.order.status;
    }, { timeout: 30000 }).toBe('REJECTED');
  });

  test('Branch 2: SAGA Rollback - Out of Stock', async ({ request }) => {
    test.setTimeout(120_000);
    // Create order with impossibly huge quantity to force Out Of Stock
    const order = await api<EntityResponse<'order'>>('/v1/orders', {
      method: 'POST',
      data: { store_id: SEEDED_STORE_ID, total_amount: 999999, items: [{ sku: 'COFFEE-ARABICA-001', quantity: 999999, unit_price: 1 }] }
    });

    const paymentRef: { current?: PaymentByOrderResponse['payment'] } = {};
    await expect.poll(async () => {
      try {
        const res = await api<PaymentByOrderResponse>(`/v1/payments/orders/${order.body.order.id}`);
        paymentRef.current = res.body.payment;
        return paymentRef.current.provider_ref;
      } catch {
        return '';
      }
    }, { timeout: 30000 }).not.toBe('');
    if (!paymentRef.current) throw new Error('missing payment');
    const payment = paymentRef.current;
    const providerRef = payment.provider_ref;

    // Trigger Payment SUCCEEDED
    await triggerWebhook(request, 'payment_intent.succeeded', providerRef, order.body.order.id, payment.id, SEEDED_STORE_ID);

    // Verify Order becomes REJECTED eventually because Warehouse refuses it (Out of Stock)
    await expect.poll(async () => {
      const res = await api<EntityResponse<'order'>>(`/v1/orders/${order.body.order.id}`);
      return res.body.order.status;
    }, { timeout: 30000 }).toBe('REJECTED');
  });

  test('Edge Case 3: Create Harvest with Non-existent Farm ID', async () => {
    // Attempt to create a harvest on farm id 999999 (which does not exist)
    // Farm-service should fail with 404/500 not found
    await expect(
      api('/v1/harvests', {
        method: 'POST',
        data: {
          farm_id: 999999,
          coffee_type: 'ARABICA',
          quantity: 10,
          notes: 'invalid farm edge case',
        },
      })
    ).rejects.toThrow();
  });

  test('Edge Case 4: Create Order with Non-existent SKU', async ({ request }) => {
    // Attempt to create an order with an invalid SKU.
    // The order should be created as PENDING, but then warehouses will reject it,
    // and the SAGA flow will eventually mark the order as REJECTED.
    const order = await api<EntityResponse<'order'>>('/v1/orders', {
      method: 'POST',
      data: {
        store_id: SEEDED_STORE_ID,
        total_amount: 199,
        items: [{ sku: 'SKU-COMPLETELY-INVALID-FAKE-SKU-999', quantity: 1, unit_price: 199 }],
      },
    });

    const paymentRef: { current?: PaymentByOrderResponse['payment'] } = {};
    await expect.poll(async () => {
      try {
        const res = await api<PaymentByOrderResponse>(`/v1/payments/orders/${order.body.order.id}`);
        paymentRef.current = res.body.payment;
        return paymentRef.current.provider_ref;
      } catch {
        return '';
      }
    }, { timeout: 30000 }).not.toBe('');
    if (!paymentRef.current) throw new Error('missing payment');
    const payment = paymentRef.current;
    const providerRef = payment.provider_ref;

    // Trigger Payment SUCCEEDED
    await triggerWebhook(request, 'payment_intent.succeeded', providerRef, order.body.order.id, payment.id, SEEDED_STORE_ID);

    await expect.poll(async () => {
      const res = await api<EntityResponse<'order'>>(`/v1/orders/${order.body.order.id}`);
      return res.body.order.status;
    }, { timeout: 30000 }).toBe('REJECTED');
  });

  test('Edge Case 5: Unauthorized Request without valid Bearer Token', async ({ request }) => {
    // Attempt to call a protected POST endpoint without an Authorization header
    const response = await request.fetch('http://localhost:8081/v1/orders', {
      method: 'POST',
      headers: {}, // No Authorization header
      data: {
        store_id: SEEDED_STORE_ID,
        total_amount: 199,
        items: [{ sku: 'COFFEE-ARABICA-001', quantity: 1, unit_price: 199 }],
      },
    });
    expect(response.status()).toBe(401);
  });
});
