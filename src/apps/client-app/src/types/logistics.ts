export type LocationType = 'FARM' | 'ROASTERY' | 'RETAILER';

export interface Location {
  id: string;
  name: string;
  type: LocationType;
  lat: number;
  lng: number;
}

export type ShipmentStatus = 
  | 'PENDING' 
  | 'ASSIGNED' 
  | 'PICKED_UP' 
  | 'IN_TRANSIT' 
  | 'DELIVERED' 
  | 'IN_TRANSIT_TO_FARM'
  | 'ARRIVED_AT_FARM'
  | 'RETURNING_TO_WAREHOUSE'
  | 'ARRIVED_WAREHOUSE'
  | 'IN_TRANSIT_TO_STORE'
  | 'ARRIVED_AT_STORE'
  | 'RETURNING_TO_BASE'
  | 'RETURNED_TO_BASE'
  | 'COMPLETED'
  | 'FAILED' 
  | 'CANCELLED';

export interface Shipment {
  id: string;
  type?: 'FARM_PICKUP' | 'RETAIL_DELIVERY' | 'RETURN_TO_BASE';
  order_id: string;
  harvest_id?: string;
  farm_id?: string;
  warehouse_id?: string;
  vehicle_id?: string;
  driver_id: string | null;
  status: ShipmentStatus;
  route_id?: string;
  current_leg?: 'OUTBOUND' | 'RETURN';
  destination_address: string;
  created_at: string;
  updated_at: string;
  delivered_at: string | null;
  returned_at?: string | null;
}

export interface RouteData {
  id: string;
  name: string;
  from: string;
  to: string;
  coordinates: [number, number][]; // [lat, lng][]
}

export interface DriverLocationUpdate {
  driver_id?: string;
  shipment_id?: string;
  latitude: number;
  longitude: number;
  heading?: number;
  speed?: number;
  route_index?: number;
  status?: string;
  occurred_at?: string;
}

export interface Driver {
  id: string;
  user_id?: string;
  name: string;
  phone: string;
  vehicle_id?: string;
  status: string;
  current_shipment_id?: string;
  is_available: boolean;
}

export interface Vehicle {
  id: string;
  plate_number: string;
  type: string;
  capacity_kg: number;
  home_warehouse_id?: string;
  status: string;
  current_latitude?: number;
  current_longitude?: number;
}
