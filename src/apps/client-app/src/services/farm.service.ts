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

  async updateFarm(id: string, data: Partial<Farm>): Promise<Farm> {
    const res = await api.put(`/v1/farms/${id}`, data);
    return res.data.farm;
  }

  async deleteFarm(id: string): Promise<void> {
    await api.delete(`/v1/farms/${id}`);
  }
}

export const farmService = new FarmService();
