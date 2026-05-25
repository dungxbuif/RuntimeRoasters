import api from "@/lib/axios";
import { API_ENDPOINTS } from "@/constants/api";
import { Location, Shipment, DriverLocationUpdate, RouteData, Driver, Vehicle } from "@/types/logistics";

class LogisticsService {
  async listLocations(): Promise<Location[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.LOCATIONS);
    return res.data.locations || [];
  }

  async listShipments(): Promise<Shipment[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.SHIPMENTS);
    return res.data.shipments || [];
  }

  async getShipment(id: string): Promise<Shipment> {
    const res = await api.get(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}`);
    return res.data.shipment;
  }

  async assignShipment(id: string, data: { driver_id?: string; vehicle_id?: string }): Promise<Shipment> {
    const res = await api.post(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}/assign`, data);
    return res.data.shipment;
  }

  async departShipment(id: string): Promise<Shipment> {
    const res = await api.post(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}/depart`);
    return res.data.shipment;
  }

  async arriveShipment(id: string): Promise<Shipment> {
    const res = await api.post(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}/arrive`);
    return res.data.shipment;
  }

  async confirmLoad(id: string): Promise<Shipment> {
    const res = await api.post(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}/confirm-load`);
    return res.data.shipment;
  }

  async confirmDelivery(id: string): Promise<Shipment> {
    const res = await api.post(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}/confirm-delivery`);
    return res.data.shipment;
  }

  async returnShipment(id: string): Promise<Shipment> {
    const res = await api.post(`${API_ENDPOINTS.LOGISTICS.SHIPMENTS}/${id}/return`);
    return res.data.shipment;
  }

  async updateDriverLocation(data: DriverLocationUpdate): Promise<void> {
    await api.post(API_ENDPOINTS.LOGISTICS.UPDATE_LOCATION, data);
  }

  async listDrivers(): Promise<Driver[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.DRIVERS);
    return res.data.drivers || [];
  }

  async listVehicles(): Promise<Vehicle[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.VEHICLES);
    return res.data.vehicles || [];
  }

  async listAvailableVehicles(): Promise<Vehicle[]> {
    const res = await api.get(API_ENDPOINTS.LOGISTICS.AVAILABLE_VEHICLES);
    return res.data.vehicles || [];
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
