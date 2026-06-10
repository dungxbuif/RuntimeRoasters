import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";

class RealtimeService {
  async getTicket(): Promise<string> {
    const res = await api.post(API_ENDPOINTS.REALTIME.TICKETS);
    return res.data.ticket;
  }
}

export const realtimeService = new RealtimeService();
