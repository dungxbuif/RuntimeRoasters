import api from "@/lib/axios";

export interface Farm {
  id: string;
  name: string;
  location: string;
  area: number;
  farm_type: string;
  owner_id: string;
  created_at: number;
  updated_at: number;
}

class FarmService {
  async listFarms(): Promise<Farm[]> {
    const res = await api.get("/v1/farms");
    return res.data.farms || [];
  }
}

export const farmService = new FarmService();
