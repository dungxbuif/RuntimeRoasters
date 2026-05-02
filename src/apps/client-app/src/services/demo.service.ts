import api from "@/lib/axios";
import { DemoPingResponse } from "@/types/api";

export const demoService = {
  ping: async (): Promise<DemoPingResponse> => {
    const response = await api.get<DemoPingResponse>('/v1/demo/ping');
    return response.data;
  },
};
