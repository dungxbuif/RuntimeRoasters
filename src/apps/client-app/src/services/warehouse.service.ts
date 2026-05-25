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
  warehouse_id?: string;
  sku: string;
  available_quantity: number;
  updated_at: string;
}

export interface Intake {
  id: string;
  harvest_id: string;
  pickup_id?: string;
  warehouse_id?: string;
  coffee_type: string;
  origin_code: string;
  quantity: number;
  status: string;
  batch_id?: string;
  created_at: string;
  updated_at: string;
}

export interface PickupRequest {
  id: string;
  harvest_id: string;
  farm_id: string;
  warehouse_id: string;
  origin_location_id: string;
  quantity: number;
  coffee_type: string;
  origin_code: string;
  status: string;
  notification_id?: string;
  shipment_id?: string;
  created_at: string;
  updated_at: string;
  dispatched_at?: string;
  received_at?: string;
}

export interface DispatchRequest {
  id: string;
  status: string;
}

class WarehouseService {
  async listIntakes(): Promise<Intake[]> {
    const res = await api.get(API_ENDPOINTS.WAREHOUSE.INTAKES);
    return res.data.intakes || [];
  }

  async listPickupRequests(): Promise<PickupRequest[]> {
    const res = await api.get(API_ENDPOINTS.WAREHOUSE.PICKUP_REQUESTS);
    return res.data.pickup_requests || [];
  }

  async dispatchPickupRequest(id: string): Promise<PickupRequest> {
    const res = await api.post(`${API_ENDPOINTS.WAREHOUSE.PICKUP_REQUESTS}/${id}/dispatch`);
    return res.data.pickup_request;
  }

  async receivePickupRequest(id: string): Promise<Intake> {
    const res = await api.post(`${API_ENDPOINTS.WAREHOUSE.PICKUP_REQUESTS}/${id}/receive`);
    return res.data.intake;
  }

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

  async startProcessing(id: string): Promise<void> {
    await api.post(`${API_ENDPOINTS.WAREHOUSE.BATCHES}/${id}/process`);
  }

  async getInventory(): Promise<Inventory[]> {
    const res = await api.get(API_ENDPOINTS.WAREHOUSE.INVENTORY);
    return res.data.inventory || [];
  }

  async listDispatchRequests(): Promise<DispatchRequest[]> {
    const res = await api.get(API_ENDPOINTS.WAREHOUSE.DISPATCH_REQUESTS);
    return res.data.dispatch_requests || [];
  }

  async dispatchRequest(id: string): Promise<DispatchRequest> {
    const res = await api.post(`${API_ENDPOINTS.WAREHOUSE.DISPATCH_REQUESTS}/${id}/dispatch`);
    return res.data;
  }
}

export const warehouseService = new WarehouseService();
