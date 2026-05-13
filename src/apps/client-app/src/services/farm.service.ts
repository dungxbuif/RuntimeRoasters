import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";

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

export type CoffeeType = 'ARABICA' | 'ROBUSTA' | 'CHERRY' | 'CULI';
export type HarvestStatus = 'NEW' | 'PROCESSING' | 'COMPLETED';

export interface Harvest {
  id: string;
  farm_id: string;
  coffee_type: CoffeeType;
  quantity: number;
  harvest_date: string;
  status: HarvestStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateHarvestInput {
  farm_id: string;
  coffee_type: CoffeeType;
  quantity: number;
  harvest_date: string;
}

class FarmService {
  async listFarms(): Promise<Farm[]> {
    const res = await api.get(API_ENDPOINTS.FARM.FARMS);
    return res.data.farms || [];
  }

  async updateFarm(id: string, data: Partial<Farm>): Promise<Farm> {
    const res = await api.put(`${API_ENDPOINTS.FARM.FARMS}/${id}`, data);
    return res.data.farm;
  }

  async deleteFarm(id: string): Promise<void> {
    await api.delete(`${API_ENDPOINTS.FARM.FARMS}/${id}`);
  }

  // Harvest methods
  async createHarvest(data: CreateHarvestInput): Promise<Harvest> {
    const res = await api.post(API_ENDPOINTS.FARM.HARVESTS, data);
    return res.data.harvest;
  }

  async listHarvests(): Promise<Harvest[]> {
    const res = await api.get(API_ENDPOINTS.FARM.HARVESTS);
    return res.data.harvests || [];
  }
}

export const farmService = new FarmService();
