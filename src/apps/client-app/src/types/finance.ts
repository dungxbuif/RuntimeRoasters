export type PaymentStatus = 'CREATED' | 'PENDING' | 'SUCCEEDED' | 'FAILED' | 'CANCELLED' | 'REFUNDED';

export interface Payment {
  id: string;
  order_id: string;
  store_id?: string;
  provider?: string;
  provider_ref?: string;
  stripe_payment_intent_id?: string;
  amount: number;
  currency: string;
  status: PaymentStatus;
  refund_ref?: string;
  checkout_url?: string;
  simulated?: boolean;
  stripe_refund_id?: string;
  created_at: string;
  updated_at: string;
}

export interface StripeWebhookDemoKey {
  provider: string;
  key_name: string;
  signing_secret: string;
  webhook_url: string;
  algorithm: 'HMAC-SHA256';
}

export interface AuditLog {
  id: string;
  partition: string;
  event_type: string;
  aggregate_id: string;
  payload: string;
  hash: string;
  previous_hash: string;
  occurred_at: string;
}
