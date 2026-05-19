import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";
import { Payment, AuditLog } from "@/types/finance";

class FinanceService {
  async listPayments(): Promise<Payment[]> {
    const res = await api.get(API_ENDPOINTS.PAYMENT.PAYMENTS);
    return res.data.payments || [];
  }

  async getPaymentByOrder(orderId: string): Promise<Payment | null> {
    try {
      const res = await api.get(API_ENDPOINTS.PAYMENT.BY_ORDER(orderId));
      return res.data.payment || null;
    } catch (e) {
      return null;
    }
  }

  async listAuditLogs(partition: string = "default"): Promise<AuditLog[]> {
    const res = await api.get(API_ENDPOINTS.AUDIT.BY_PARTITION(partition));
    return res.data.logs || [];
  }
}

export const financeService = new FinanceService();
