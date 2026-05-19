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
  | 'FAILED' 
  | 'CANCELLED';

export interface Shipment {
  id: string;
  order_id: string;
  driver_id: string | null;
  status: ShipmentStatus;
  destination_address: string;
  created_at: string;
  updated_at: string;
  delivered_at: string | null;
}

export interface RouteData {
  id: string;
  name: string;
  from: string;
  to: string;
  coordinates: [number, number][]; // [lat, lng][]
}

export interface DriverLocationUpdate {
  driver_id: string;
  shipment_id?: string;
  latitude: number;
  longitude: number;
}
