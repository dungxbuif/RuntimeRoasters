import api from "@/lib/axios";
import { financeService } from "@/services/finance.service";

export interface Store {
  id: string;
  name: string;
  location: string;
  is_active: boolean;
}

export interface CreateOrderRequest {
  store_id: string;
  items: {
    sku: string;
    quantity: number;
  }[];
}

export interface Order {
  id: string;
  store_id: string;
  status: string;
  total_amount: number;
  created_at: string;
  updated_at: string;
}

class RetailService {
  async listStores(): Promise<Store[]> {
    const res = await api.get('/v1/stores');
    return res.data.stores || [];
  }

  async createOrder(data: CreateOrderRequest, idempotencyKey?: string): Promise<Order> {
    const res = await api.post('/v1/orders', data, {
      headers: idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : {}
    });
    return res.data.order;
  }

  async getOrder(id: string): Promise<Order> {
    const res = await api.get(`/v1/orders/${id}`);
    return res.data.order;
  }

  async listOrders(): Promise<Order[]> {
    const res = await api.get(`/v1/orders`);
    return res.data.orders || [];
  }

  async simulatePayment(orderId: string): Promise<void> {
    const payment = await financeService.getPaymentByOrder(orderId);
    if (!payment) {
      throw new Error(`No payment intent found for order ${orderId}`);
    }
    await financeService.sendStripeWebhook(payment, "passed");
  }

  async confirmOrder(id: string): Promise<void> {
    await api.post(`/v1/orders/${id}/confirm`);
  }
}

export const retailService = new RetailService();
