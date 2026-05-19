export type PaymentStatus = 'CREATED' | 'PENDING' | 'SUCCEEDED' | 'FAILED' | 'CANCELLED' | 'REFUNDED';

export interface Payment {
  id: string;
  order_id: string;
  stripe_payment_intent_id?: string;
  amount: number;
  currency: string;
  status: PaymentStatus;
  stripe_refund_id?: string;
  created_at: string;
  updated_at: string;
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
