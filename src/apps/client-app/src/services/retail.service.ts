import api from "@/lib/axios";
import { financeService } from "@/services/finance.service";

export interface Store {
  id: string;
  code?: string;
  name: string;
  city?: string;
  address?: string;
  location?: string;
  status?: string;
  is_active?: boolean;
}

export interface PublicMenuItem {
  id: string;
  menu_id: string;
  code: string;
  product_group_code: string;
  category_code: string;
  size: string;
  name: string;
  description?: string;
  stock_sku: string;
  coffee_type: string;
  price: number;
  image_url?: string;
  available?: boolean;
  available_units?: number;
}

export interface PublicMenu {
  id: string;
  code: string;
  name: string;
  description?: string;
  items: PublicMenuItem[];
}

export interface SaleItem {
  id: string;
  product_id: string;
  trace_code: string;
  product_name: string;
  size: string;
  unit_price: number;
  public_url: string;
  sold_at: string;
  invoice_no?: string;
}

export interface Sale {
  id: string;
  store_id: string;
  invoice_no: string;
  total_amount: number;
  sold_at: string;
  items: SaleItem[];
}

export interface CreateSaleRequest {
  store_id: string;
  items: Array<{ menu_item_id: string; product_id: string }>;
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

  async listPublicStoreMenu(storeId: string): Promise<PublicMenu[]> {
    const res = await api.get(`/v1/public/retail/stores/${encodeURIComponent(storeId)}/menu-items`);
    return res.data.menus || [];
  }

  async listPublicStores(): Promise<Store[]> {
    const res = await api.get('/v1/public/retail/stores');
    return res.data.stores || [];
  }

  async createPublicSale(data: CreateSaleRequest, idempotencyKey: string): Promise<Sale> {
    const res = await api.post('/v1/public/retail/sales', data, { headers: { 'Idempotency-Key': idempotencyKey } });
    return res.data.sale;
  }

  async listSoldItems(storeId: string): Promise<SaleItem[]> {
    const res = await api.get(`/v1/retail/stores/${encodeURIComponent(storeId)}/sold-items`);
    return res.data.sold_items || [];
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
