import { API_ENDPOINTS } from '@/constants/api';
import { ENV } from '@/constants/env';
import api from '@/lib/axios';
import { TopologyConfig, TopologyHistoryEntry } from '@/types/topology';

class TopologyService {
  async getPublicConfig(): Promise<TopologyConfig> {
    const res = await api.get(API_ENDPOINTS.TRACE.PUBLIC_TOPOLOGY_CONFIG);
    return res.data;
  }

  async getPublicHistory(flowId?: string): Promise<TopologyHistoryEntry[]> {
    const res = await api.get(API_ENDPOINTS.TRACE.PUBLIC_TOPOLOGY_HISTORY, {
      params: { flow_id: flowId || undefined, limit: 80 },
    });
    return res.data.events || [];
  }

  async getTicket(): Promise<string> {
    const res = await api.post(API_ENDPOINTS.REALTIME.TICKETS);
    return res.data.ticket;
  }

  publicSocketURL(flowId?: string): string {
    const params = new URLSearchParams();
    if (flowId) params.set('flow_id', flowId);
    return `${ENV.SOCKET_URL}/v1/realtime/public/topology/ws${params.size ? `?${params}` : ''}`;
  }

  privateSocketURL(ticket: string, flowId?: string): string {
    const params = new URLSearchParams();
    params.set('ticket', ticket);
    if (flowId) params.set('flow_id', flowId);
    return `${ENV.SOCKET_URL}/v1/realtime/stream?${params}`;
  }
}

export const topologyService = new TopologyService();
