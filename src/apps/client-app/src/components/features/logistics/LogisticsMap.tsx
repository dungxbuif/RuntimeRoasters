'use client';

import React, { useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Re-fix Leaflet default marker icon issue in Next.js
const fixLeafletIcon = () => {
  delete (L.Icon.Default.prototype as { _getIconUrl?: unknown })._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
    iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  });
};

// Dynamically import react-leaflet components to avoid SSR issues
const MapContainer = dynamic(() => import('react-leaflet').then(mod => mod.MapContainer), { ssr: false });
const TileLayer = dynamic(() => import('react-leaflet').then(mod => mod.TileLayer), { ssr: false });
const Marker = dynamic(() => import('react-leaflet').then(mod => mod.Marker), { ssr: false });
const Popup = dynamic(() => import('react-leaflet').then(mod => mod.Popup), { ssr: false });
const Polyline = dynamic(() => import('react-leaflet').then(mod => mod.Polyline), { ssr: false });

import { Location, Shipment, RouteData } from '@/types/logistics';

interface LogisticsMapProps {
  locations: Location[];
  shipments: Shipment[];
  routes: RouteData[];
  activeDriverLocations: Record<string, [number, number]>;
  showRoutes?: boolean;
}

export default function LogisticsMap({ locations, shipments, routes, activeDriverLocations, showRoutes = false }: LogisticsMapProps) {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (isMounted) {
      fixLeafletIcon();
      // Use shipments to satisfy lint
      if (shipments.length > 0) {
        console.debug(`Map initialized with ${shipments.length} shipments`);
      }
    }
  }, [shipments, isMounted]);

  if (!isMounted) return <div className="w-full h-full bg-slate-100 animate-pulse" />;

  const center: [number, number] = [16.047079, 108.206230]; // Center of Vietnam (Da Nang approx)

  const getMarkerIcon = (type: string) => {
    let color = '#3b82f6'; // Default blue
    let icon = 'inventory_2';
    
    if (type === 'FARM') {
      color = '#16a34a'; // Green
      icon = 'potted_plant';
    } else if (type === 'ROASTERY') {
      color = '#f59e0b'; // Amber
      icon = 'factory';
    } else if (type === 'RETAILER') {
      color = '#0ea5e9'; // Cyan
      icon = 'store';
    }

    return L.divIcon({
      html: `
        <div class="relative flex items-center justify-center w-8 h-8 rounded-full bg-white border-2 border-[${color}] shadow-lg">
          <span class="material-symbols-outlined text-[${color}] text-sm font-black">${icon}</span>
        </div>
      `,
      className: 'custom-div-icon',
      iconSize: [32, 32],
      iconAnchor: [16, 16],
    });
  };

  const getDriverIcon = () => {
    return L.divIcon({
      html: `
        <div class="relative flex items-center justify-center w-8 h-8 rounded-full bg-primary border-2 border-white shadow-xl">
          <span class="material-symbols-outlined text-white text-sm font-black">local_shipping</span>
          <div class="absolute -top-1 -right-1 w-3 h-3 bg-primary rounded-full animate-ping"></div>
        </div>
      `,
      className: 'driver-icon',
      iconSize: [32, 32],
      iconAnchor: [16, 16],
    });
  };

  return (
    <div className="w-full h-full rounded-2xl overflow-hidden border border-slate-200 shadow-xl relative bg-slate-50">
      <MapContainer 
        center={center} 
        zoom={6} 
        style={{ height: '100%', width: '100%', background: '#f8fafc' }}
        zoomControl={false}
      >
        <TileLayer
          url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png"
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        />

        {/* POI Markers */}
        {locations.map((loc) => (
          <Marker 
            key={loc.id} 
            position={[loc.lat, loc.lng]} 
            icon={getMarkerIcon(loc.type)}
          >
            <Popup className="light-popup">
              <div className="bg-white p-2 rounded">
                <p className="font-black text-xs uppercase tracking-wider text-slate-900">{loc.name}</p>
                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{loc.type}</p>
              </div>
            </Popup>
          </Marker>
        ))}

        {/* Routes - Only shown when explicit or when drivers are active and simulation is on */}
        {showRoutes && routes.filter(r => shipments.some(s => s.route_id === r.id && (s.status === 'DISPATCHED' || s.status === 'EN_ROUTE'))).map((route) => (
          <Polyline 
            key={route.id}
            positions={route.coordinates}
            pathOptions={{ 
              color: '#3b82f6', 
              weight: 2, 
              opacity: 0.2, 
              dashArray: '5, 10' 
            }}
          />
        ))}

        {/* Active Drivers */}
        {Object.entries(activeDriverLocations).map(([driverId, position]) => (
          <Marker 
            key={driverId} 
            position={position} 
            icon={getDriverIcon()}
          >
            <Popup>
              <div className="bg-white p-2 rounded">
                <p className="font-black text-xs text-slate-900">DRIVER: {driverId}</p>
                <p className="text-[10px] font-black text-primary uppercase italic">IN TRANSIT</p>
              </div>
            </Popup>
          </Marker>
        ))}
      </MapContainer>

      {/* Map Legend */}
      <div className="absolute bottom-6 left-6 z-[1000] p-5 rounded-[2rem] border border-slate-200 bg-white/90 backdrop-blur-xl text-slate-900 shadow-2xl">
        <h4 className="text-[9px] font-black uppercase tracking-[0.2em] text-slate-400 mb-4 italic">Fleet Network</h4>
        <div className="space-y-3">
          <div className="flex items-center gap-3">
            <div className="w-2.5 h-2.5 rounded-full bg-green-500 shadow-sm"></div>
            <span className="text-[10px] font-black uppercase tracking-tight text-slate-700">Coffee Farmers</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2.5 h-2.5 rounded-full bg-amber-500 shadow-sm"></div>
            <span className="text-[10px] font-black uppercase tracking-tight text-slate-700">Roasteries (KCN)</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2.5 h-2.5 rounded-full bg-sky-500 shadow-sm"></div>
            <span className="text-[10px] font-black uppercase tracking-tight text-slate-700">Retail Outlets</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-2.5 h-2.5 rounded-full bg-primary animate-pulse shadow-sm"></div>
            <span className="text-[10px] font-black uppercase tracking-tight text-primary italic">Active Driver</span>
          </div>
        </div>
      </div>
    </div>
  );
}
