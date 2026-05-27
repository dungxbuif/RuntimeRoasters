import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";
import { Payment, AuditLog, StripeWebhookDemoKey } from "@/types/finance";

type StripeWebhookOutcome = "passed" | "failed";

class FinanceService {
  async listPayments(): Promise<Payment[]> {
    const res = await api.get(API_ENDPOINTS.PAYMENT.PAYMENTS);
    return res.data.payments || [];
  }

  async getPaymentByOrder(orderId: string): Promise<Payment | null> {
    try {
      const res = await api.get(API_ENDPOINTS.PAYMENT.BY_ORDER(orderId));
      return res.data.payment || null;
    } catch {
      return null;
    }
  }

  async getStripeWebhookDemoKey(): Promise<StripeWebhookDemoKey> {
    const res = await api.get(API_ENDPOINTS.PAYMENT.STRIPE_WEBHOOK_KEY);
    return res.data;
  }

  async sendStripeWebhook(payment: Payment, outcome: StripeWebhookOutcome): Promise<void> {
    const key = await this.getStripeWebhookDemoKey();
    const timestamp = Math.floor(Date.now() / 1000);
    const providerRef = payment.provider_ref || payment.stripe_payment_intent_id || payment.id;
    const eventType = outcome === "passed" ? "payment_intent.succeeded" : "payment_intent.payment_failed";
    const payload = {
      id: `evt_rr_demo_${outcome}_${crypto.randomUUID()}`,
      object: "event",
      type: eventType,
      livemode: false,
      created: timestamp,
      data: {
        object: {
          id: providerRef,
          object: "payment_intent",
          amount: Math.round(Number(payment.amount || 0) * 100),
          currency: (payment.currency || "usd").toLowerCase(),
          status: outcome === "passed" ? "succeeded" : "requires_payment_method",
          ...(outcome === "failed"
            ? {
                last_payment_error: {
                  type: "card_error",
                  code: "card_declined",
                  message: "Demo card declined",
                },
              }
            : {}),
          metadata: {
            order_id: payment.order_id,
            payment_id: payment.id,
            store_id: payment.store_id || "",
          },
        },
      },
    };
    const rawBody = JSON.stringify(payload);
    const signature = await hmacSha256Hex(key.signing_secret, `${timestamp}.${rawBody}`);
    await api.post(key.webhook_url || API_ENDPOINTS.PAYMENT.STRIPE_WEBHOOK, rawBody, {
      headers: {
        "Content-Type": "application/json",
        "Stripe-Signature": `t=${timestamp},v1=${signature}`,
      },
      transformRequest: [(data) => data],
    });
  }

  async listAuditLogs(partition: string = "default"): Promise<AuditLog[]> {
    const res = await api.get(API_ENDPOINTS.AUDIT.BY_PARTITION(partition));
    return res.data.logs || [];
  }
}

export const financeService = new FinanceService();

async function hmacSha256Hex(secret: string, message: string): Promise<string> {
  const encoder = new TextEncoder();
  const key = await crypto.subtle.importKey(
    "raw",
    encoder.encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"]
  );
  const signature = await crypto.subtle.sign("HMAC", key, encoder.encode(message));
  return [...new Uint8Array(signature)].map((byte) => byte.toString(16).padStart(2, "0")).join("");
}
