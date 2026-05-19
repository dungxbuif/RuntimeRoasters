import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";
import { Location, Shipment, DriverLocationUpdate, RouteData } from "@/types/logistics";

class LogisticsService {
  async listLocations(): Promise<Location[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.LOCATIONS);
    return res.data.locations || [];
  }

  async listShipments(): Promise<Shipment[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.SHIPMENTS);
    return res.data.shipments || [];
  }

  async updateDriverLocation(data: DriverLocationUpdate): Promise<void> {
    await api.post(API_ENDPOINTS.LOGISTICS.UPDATE_LOCATION, data);
  }

  // Helper to fetch static routes for simulation
  async getRoutes(): Promise<RouteData[]> {
    // In a real app, this might be an API call. 
    // Here we might fetch a static JSON or use a mock.
    try {
      const res = await fetch('/data/routes.json');
      if (!res.ok) return [];
      return await res.json();
    } catch (e) {
      console.error("Failed to fetch routes", e);
      return [];
    }
  }
}

export const logisticsService = new LogisticsService();
