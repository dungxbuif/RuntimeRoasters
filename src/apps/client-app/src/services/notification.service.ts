import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";

export interface RuntimeNotification {
  id: string;
  target_role?: string;
  target_user_id?: string;
  farm_id?: string;
  warehouse_id?: string;
  store_id?: string;
  shipment_id?: string;
  entity_type: string;
  entity_id: string;
  type: string;
  title: string;
  message: string;
  severity: "info" | "warning" | "error";
  status: "UNREAD" | "ACKED" | "RESOLVED";
  created_at: string;
  acknowledged_at?: string;
}

class NotificationService {
  async listNotifications(): Promise<RuntimeNotification[]> {
    const res = await api.get(API_ENDPOINTS.REALTIME.NOTIFICATIONS);
    return res.data.notifications || [];
  }

  async acknowledge(id: string): Promise<RuntimeNotification> {
    const res = await api.post(`${API_ENDPOINTS.REALTIME.NOTIFICATIONS}/${id}/ack`);
    return res.data.notification;
  }

  async resolve(id: string): Promise<RuntimeNotification> {
    const res = await api.post(`${API_ENDPOINTS.REALTIME.NOTIFICATIONS}/${id}/resolve`);
    return res.data.notification;
  }
}

export const notificationService = new NotificationService();
