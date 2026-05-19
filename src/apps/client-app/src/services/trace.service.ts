import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";
import { TraceDocument, TraceEvent } from "@/types/trace";

class TraceService {
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
