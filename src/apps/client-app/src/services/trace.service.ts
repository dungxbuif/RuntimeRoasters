import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";
import { TraceDocument, TraceEvent } from "@/types/trace";

export interface PublicTraceItem {
  product_id: string;
  trace_code: string;
  product_name: string;
  store_name: string;
  invoice_no: string;
  sold_at: string;
  public_url: string;
}

export interface PublicTraceMilestone {
  status: string;
  title: string;
  description: string;
  timestamp: string;
  service?: string;
}

export interface PublicCupTrace {
  product_id: string;
  trace_code: string;
  sale_item: { id: string; sold_at: string };
  invoice_no: string;
  product: { menu_item_id: string; name: string; sku: string; size: string };
  store: { id: string; name: string; city: string; address: string };
  farm: { harvest_id: string };
  warehouse: { id: string };
  origin_batch_id: string;
  harvest_id: string;
  inventory_lot_id: string;
  milestones: PublicTraceMilestone[];
  service_participation: string[];
  trace_ids: string[];
  public_safe: boolean;
}

class TraceService {
  async getPublicTrace(code: string): Promise<PublicCupTrace> {
    const res = await api.get(`/v1/public/trace/${encodeURIComponent(code)}`);
    return res.data.trace;
  }

  async listPublicTraceItems(limit = 25): Promise<PublicTraceItem[]> {
    const res = await api.get('/v1/public/trace-items', { params: { limit } });
    return res.data.trace_items || [];
  }
  async getTraceEvents(id: string): Promise<TraceEvent[]> {
    const res = await api.get(API_ENDPOINTS.TRACE.BY_ID(id));
    return res.data.events || [];
  }

  async getTraceDocument(id: string): Promise<TraceDocument | null> {
    try {
      const res = await api.get(API_ENDPOINTS.TRACE.DOCUMENT(id));
      return res.data.document || null;
    } catch (error) {
      console.error("Failed to fetch trace document", error);
      return null;
    }
  }
}

export const traceService = new TraceService();
