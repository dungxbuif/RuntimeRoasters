import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";

export type BatchStatus = 'RECEIVED' | 'PROCESSING' | 'STOCKED';
export type QualityFlag = 'NORMAL' | 'QUALITY_WARNING';

export interface RoastRun {
  id: string;
  batch_id: string;
  run_number: number;
  input_weight: number;
  output_weight: number;
  loss_percent: number;
  roaster_name: string;
  created_at: string;
}

export interface ProductionBatch {
  id: string;
  harvest_id: string;
  batch_id: string;
  status: BatchStatus;
  coffee_type: string;
  origin_code: string;
  intake_weight: number;
  total_input_weight: number;
  total_output_weight: number;
  weight_loss_percent: number;
  quality_flag: QualityFlag;
  intake_note?: string;
  roast_runs?: RoastRun[];
  created_at: string;
  updated_at: string;
}

export interface Inventory {
  id: string;
  coffee_type: string;
  origin_code: string;
  sku: string;
  available_quantity: number;
  updated_at: string;
}

class WarehouseService {
  async listBatches(): Promise<ProductionBatch[]> {
    const res = await api.get(API_ENDPOINTS.WAREHOUSE.BATCHES);
    return res.data.batches || [];
  }

  async getBatch(id: string): Promise<ProductionBatch> {
    const res = await api.get(`${API_ENDPOINTS.WAREHOUSE.BATCHES}/${id}`);
    return res.data.batch;
  }

  async updateIntake(id: string, actualWeight: number, note?: string): Promise<ProductionBatch> {
    const res = await api.patch(`${API_ENDPOINTS.WAREHOUSE.BATCHES}/${id}/intake`, {
      actual_weight: actualWeight,
      note: note
    });
    return res.data.batch;
  }

  async addRoastRun(batchId: string, data: { input: number; output: number; roaster: string }): Promise<RoastRun> {
    const res = await api.post(`${API_ENDPOINTS.WAREHOUSE.BATCHES}/${batchId}/runs`, data);
    return res.data.run;
  }

  async finalizeBatch(id: string): Promise<void> {
    await api.post(`${API_ENDPOINTS.WAREHOUSE.BATCHES}/${id}/finalize`);
  }

  async getInventory(): Promise<Inventory[]> {
    const res = await api.get(API_ENDPOINTS.WAREHOUSE.INVENTORY);
    return res.data.inventory || [];
  }
}

export const warehouseService = new WarehouseService();
