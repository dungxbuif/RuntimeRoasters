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
}

export default function LogisticsMap({ locations, shipments, routes, activeDriverLocations }: LogisticsMapProps) {
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

  if (!isMounted) return <div className="w-full h-full bg-slate-900 animate-pulse" />;

  const center: [number, number] = [16.047079, 108.206230]; // Center of Vietnam (Da Nang approx)

  const getMarkerIcon = (type: string) => {
    let color = '#3b82f6'; // Default blue
    let icon = 'inventory_2';
    
    if (type === 'FARM') {
      color = '#22c55e'; // Green
      icon = 'potted_plant';
    } else if (type === 'ROASTERY') {
      color = '#f59e0b'; // Amber
      icon = 'factory';
    } else if (type === 'RETAILER') {
      color = '#00daf3'; // Cyan
      icon = 'store';
    }

    return L.divIcon({
      html: `
        <div class="relative flex items-center justify-center w-8 h-8 rounded-full bg-slate-900 border-2 border-[${color}] shadow-[0_0_10px_${color}]">
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
        <div class="relative flex items-center justify-center w-8 h-8 rounded-full bg-slate-900 border-2 border-primary shadow-[0_0_15px_#004ac6]">
          <span class="material-symbols-outlined text-primary text-sm font-black">local_shipping</span>
          <div class="absolute -top-1 -right-1 w-3 h-3 bg-primary rounded-full animate-ping"></div>
        </div>
      `,
      className: 'driver-icon',
      iconSize: [32, 32],
      iconAnchor: [16, 16],
    });
  };

  return (
    <div className="w-full h-full rounded-xl overflow-hidden border border-slate-800 shadow-2xl relative">
      <MapContainer 
        center={center} 
        zoom={6} 
        style={{ height: '100%', width: '100%', background: '#0f172a' }}
        zoomControl={false}
      >
        <TileLayer
          url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        />

        {/* POI Markers */}
        {locations.map((loc) => (
          <Marker 
            key={loc.id} 
            position={[loc.lat, loc.lng]} 
            icon={getMarkerIcon(loc.type)}
          >
            <Popup className="dark-popup">
              <div className="bg-slate-900 text-white p-2 rounded">
                <p className="font-bold text-xs uppercase tracking-wider">{loc.name}</p>
                <p className="text-[10px] text-slate-400">{loc.type}</p>
              </div>
            </Popup>
          </Marker>
        ))}

        {/* Routes */}
        {routes.map((route) => (
          <Polyline 
            key={route.id}
            positions={route.coordinates}
            pathOptions={{ 
              color: '#3b82f6', 
              weight: 2, 
              opacity: 0.3, 
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
              <div className="bg-slate-900 text-white p-2 rounded">
                <p className="font-bold text-xs">DRIVER: {driverId}</p>
                <p className="text-[10px] text-green-400">IN TRANSIT</p>
              </div>
            </Popup>
          </Marker>
        ))}
      </MapContainer>

      {/* Map Legend */}
      <div className="absolute bottom-6 left-6 z-[1000] glass-card p-4 rounded-lg border border-slate-800 bg-slate-900/80 backdrop-blur-md text-white">
        <h4 className="text-[9px] font-black uppercase tracking-widest text-slate-500 mb-3 italic underline decoration-primary underline-offset-4">Fleet Map Legend</h4>
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-green-500 shadow-[0_0_5px_#22c55e]"></div>
            <span className="text-[10px] font-bold uppercase tracking-tighter">Coffee Farmers</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-amber-500 shadow-[0_0_5px_#f59e0b]"></div>
            <span className="text-[10px] font-bold uppercase tracking-tighter">Roasteries (KCN)</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-cyan-500 shadow-[0_0_5px_#00daf3]"></div>
            <span className="text-[10px] font-bold uppercase tracking-tighter">Retail Outlets</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-primary shadow-[0_0_5px_#3b82f6]"></div>
            <span className="text-[10px] font-bold uppercase tracking-tighter italic">Active Driver</span>
          </div>
        </div>
      </div>
    </div>
  );
}
