'use client';

import React, { useEffect, useState, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { logisticsService } from '@/services/logistics.service';
import { Shipment } from '@/types/logistics';
import dynamic from 'next/dynamic';
import { Activity, Gauge, Truck, AlertTriangle, CheckCircle2, MapPin, X } from 'lucide-react';
import Link from 'next/link';
import { useSocket } from '@/hooks/useSocket';

const LogisticsMap = dynamic(() => import('@/components/features/logistics/LogisticsMap'), {
  ssr: false,
  loading: () => <div className="w-full h-full bg-slate-100 animate-pulse flex items-center justify-center font-black text-slate-400 italic">Synchronizing Geospatial Node...</div>
});

export default function LogisticsDashboard() {
  const [activeDriverLocations, setActiveDriverLocations] = useState<Record<string, [number, number]>>({});
  const [selectedShipment, setSelectedShipment] = useState<Shipment | null>(null);

  const { data: initialLocations = [] } = useQuery({
    queryKey: ['logistics', 'driver-locations'],
    queryFn: () => logisticsService.listDriverLocations(),
  });

  useEffect(() => {
    if (initialLocations.length > 0) {
      const locs: Record<string, [number, number]> = {};
      initialLocations.forEach(l => {
        locs[l.driver_id] = [l.latitude, l.longitude];
      });
      setActiveDriverLocations(locs);
    }
  }, [initialLocations]);

  const { lastMessage } = useSocket();

  useEffect(() => {
    if (lastMessage?.type === 'topology.event') {
       const event = lastMessage.data;
       if (event.topic === 'logistics.gps.updated') {
         const payload = event.payload;
         setActiveDriverLocations(prev => ({
           ...prev,
           [payload.driver_id]: [payload.latitude, payload.longitude]
         }));
       }
    }
  }, [lastMessage]);

  const { data: locations = [] } = useQuery({
    queryKey: ['logistics', 'locations'],
    queryFn: () => logisticsService.listLocations(),
  });

  const { data: shipmentsData = [] } = useQuery({
    queryKey: ['logistics', 'shipments'],
    queryFn: () => logisticsService.listShipments(),
    refetchInterval: 5000,
  });

  const shipments = useMemo(() => {
    if (shipmentsData.length > 0) return shipmentsData;
    // Mock data for demo if API returns empty
    return [
      { id: 'shp-001', order_id: 'ORD-7721', driver_id: 'driver-1', status: 'IN_TRANSIT', destination_address: 'Xưởng rang Sóng Thần', created_at: new Date().toISOString(), updated_at: new Date().toISOString(), delivered_at: null },
      { id: 'shp-002', order_id: 'ORD-8832', driver_id: 'driver-2', status: 'IN_TRANSIT', destination_address: 'Cửa hàng Quận 1', created_at: new Date().toISOString(), updated_at: new Date().toISOString(), delivered_at: null },
      { id: 'shp-003', order_id: 'ORD-9943', driver_id: 'driver-3', status: 'PENDING', destination_address: 'Xưởng rang Hòa Khánh', created_at: new Date().toISOString(), updated_at: new Date().toISOString(), delivered_at: null },
    ] as Shipment[];
  }, [shipmentsData]);

  const simulationTargets = useMemo(() => {
    return shipmentsData.filter(shipment =>
      shipment.id &&
      shipment.driver_id &&
      isUuid(shipment.driver_id) &&
      ['IN_TRANSIT', 'IN_TRANSIT_TO_FARM', 'RETURNING_TO_WAREHOUSE', 'IN_TRANSIT_TO_STORE', 'RETURNING_TO_BASE'].includes(shipment.status)
    );
  }, [shipmentsData]);

  const { data: routes = [] } = useQuery({
    queryKey: ['logistics', 'routes'],
    queryFn: () => logisticsService.getRoutes(),
  });

  return (
    <div className="h-full flex flex-col gap-6">
      {/* Header Info */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <Truck className="w-8 h-8 text-primary" />
            Real-time Logistics Control
          </h1>
          <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Monitoring {shipments.length} active shipments across Vietnam network
          </p>
        </div>
      </div>

      <div className="flex-1 flex gap-6 overflow-hidden min-h-0">
        {/* Left: Map Visualization */}
        <div className="flex-1 min-w-0">
          <LogisticsMap 
            locations={locations}
            shipments={shipments}
            routes={routes}
            activeDriverLocations={activeDriverLocations}
            showRoutes={true}
          />
        </div>

        {/* Right: Shipment Telemetry */}
        <div className="w-96 flex flex-col gap-4 overflow-hidden">
          <div className="glass-card flex flex-col h-full bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
            <div className="p-4 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
              <h2 className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 flex items-center gap-2">
                <Activity className="w-4 h-4 text-blue-500" />
                Transit Telemetry
              </h2>
              <span className="text-[9px] font-bold bg-slate-900 text-white px-2 py-0.5 rounded italic">
                {shipments.filter(s => s.status === 'IN_TRANSIT').length} EN ROUTE
              </span>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {shipments.length === 0 ? (
                <div className="h-full flex flex-col items-center justify-center text-slate-400 opacity-50 grayscale py-20">
                  <Truck className="w-12 h-12 mb-2 stroke-[1]" />
                  <p className="text-[10px] font-bold uppercase tracking-widest">No active shipments</p>
                </div>
              ) : (
                shipments.map(shipment => (
                  <div 
                    key={shipment.id}
                    onClick={() => setSelectedShipment(shipment)}
                    className={`p-4 rounded-lg border-l-4 transition-all hover:translate-x-1 cursor-pointer ${
                      selectedShipment?.id === shipment.id ? 'ring-2 ring-primary ring-offset-2 shadow-md' : ''
                    } ${
                      shipment.status === 'IN_TRANSIT' 
                        ? 'bg-blue-50/50 border-primary' 
                        : shipment.status === 'FAILED'
                        ? 'bg-red-50/50 border-error'
                        : 'bg-slate-50 border-slate-300'
                    }`}
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div className="flex items-center gap-2">
                        <span className="text-[10px] font-black font-mono text-slate-900">
                          {shipment.id.slice(0, 8).toUpperCase()}
                        </span>
                        {shipment.status === 'DELIVERED' && <CheckCircle2 className="w-3 h-3 text-green-500" />}
                        {shipment.status === 'FAILED' && <AlertTriangle className="w-3 h-3 text-red-500 animate-pulse" />}
                      </div>
                      <span className={`text-[8px] font-black uppercase px-2 py-0.5 rounded-full border ${
                        shipment.status === 'IN_TRANSIT' 
                          ? 'bg-primary/10 border-primary/20 text-primary' 
                          : 'bg-slate-200 border-slate-300 text-slate-500'
                      }`}>
                        {shipment.status}
                      </span>
                    </div>

                    <div className="grid grid-cols-2 gap-2 text-[9px] font-bold text-slate-500 uppercase tracking-tighter">
                      <div className="flex flex-col">
                        <span className="text-[7px] text-slate-400">Order</span>
                        <span className="text-slate-700 truncate">{shipment.order_id.slice(0, 8)}</span>
                      </div>
                      <div className="flex flex-col">
                        <span className="text-[7px] text-slate-400">Dest</span>
                        <span className="text-slate-700 truncate">{shipment.destination_address || '---'}</span>
                      </div>
                    </div>

                    {shipment.status === 'IN_TRANSIT' && (
                      <div className="mt-3 pt-3 border-t border-blue-100 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="flex items-center gap-1 text-blue-600 font-black italic">
                            <Gauge className="w-3 h-3" />
                            <span className="tabular-nums">64 km/h</span>
                          </div>
                        </div>
                        <div className="flex -space-x-2">
                          <div className="w-6 h-6 rounded-full bg-primary border-2 border-white flex items-center justify-center text-[8px] text-white font-bold">
                            {shipment.driver_id?.slice(-1) || 'D'}
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                ))
              )}
            </div>

            {/* Micro-stats */}
            <div className="p-4 bg-slate-900 text-white rounded-t-2xl mt-auto">
              <h3 className="text-[8px] font-black text-slate-500 uppercase tracking-[0.3em] mb-4">Network Health</h3>
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between text-[9px] font-black mb-1 italic">
                    <span className="text-slate-400">Avg Transit Time</span>
                    <span className="text-primary">4h 12m</span>
                  </div>
                  <div className="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                    <div className="bg-primary h-full w-[65%]"></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between text-[9px] font-black mb-1 italic">
                    <span className="text-slate-400">Success Rate</span>
                    <span className="text-green-400">99.8%</span>
                  </div>
                  <div className="w-full bg-slate-800 h-1 rounded-full overflow-hidden">
                    <div className="bg-green-400 h-full w-[99%]"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Shipment Detail Drawer */}
      {selectedShipment && (
        <div className="fixed inset-y-0 right-0 w-96 bg-white shadow-2xl z-[100] border-l border-slate-200 flex flex-col animate-in slide-in-from-right duration-300">
           <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
              <div>
                <h3 className="text-sm font-black uppercase tracking-tight">Shipment Details</h3>
                <code className="text-[10px] font-mono text-slate-400">ID: {selectedShipment.id.slice(0, 8).toUpperCase()}</code>
              </div>
              <button 
                onClick={() => setSelectedShipment(null)}
                className="p-2 hover:bg-slate-200 rounded-full transition-colors"
              >
                <X className="w-4 h-4 text-slate-400" />
              </button>
           </div>
           
           <div className="flex-1 overflow-y-auto p-6 space-y-8">
              {/* Status Section */}
              <div>
                 <span className="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-4 block">Current Status</span>
                 <div className="flex items-center gap-3">
                    <div className={`w-10 h-10 rounded-2xl flex items-center justify-center ${
                      selectedShipment.status === 'DELIVERED' ? 'bg-green-100 text-green-600' : 'bg-blue-100 text-blue-600'
                    }`}>
                       <Truck className="w-5 h-5" />
                    </div>
                    <div>
                       <div className="text-sm font-black uppercase italic">{selectedShipment.status}</div>
                       <div className="text-[10px] font-bold text-slate-400 italic">Last updated: {new Date(selectedShipment.updated_at).toLocaleTimeString()}</div>
                    </div>
                 </div>
              </div>

              {/* Destination Section */}
              <div className="bg-slate-50 rounded-2xl p-5 border border-slate-100">
                 <span className="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-3 block">Destination</span>
                 <div className="flex items-start gap-3">
                    <MapPin className="w-4 h-4 text-primary mt-0.5" />
                    <p className="text-xs font-bold text-slate-700 leading-relaxed">{selectedShipment.destination_address || 'Unspecified Address'}</p>
                 </div>
              </div>

              {/* Assignment Section */}
              <div className="grid grid-cols-2 gap-4">
                 <div className="bg-white border border-slate-200 rounded-xl p-4">
                    <span className="text-[8px] font-black uppercase text-slate-400 mb-2 block">Driver</span>
                    <div className="text-[10px] font-bold text-slate-800">{selectedShipment.driver_id ? `Driver #${selectedShipment.driver_id.slice(-4)}` : 'Unassigned'}</div>
                 </div>
                 <div className="bg-white border border-slate-200 rounded-xl p-4">
                    <span className="text-[8px] font-black uppercase text-slate-400 mb-2 block">Vehicle</span>
                    <div className="text-[10px] font-bold text-slate-800">{selectedShipment.vehicle_id ? `Unit #${selectedShipment.vehicle_id.slice(-4)}` : 'Unassigned'}</div>
                 </div>
              </div>

              {/* Telemetry Section */}
              {activeDriverLocations[selectedShipment.driver_id || ''] && (
                <div className="space-y-3">
                   <span className="text-[9px] font-black uppercase tracking-widest text-slate-400 block">Live Telemetry</span>
                   <div className="font-mono text-[10px] bg-slate-900 text-primary p-4 rounded-xl space-y-1">
                      <div className="flex justify-between">
                         <span className="opacity-50 text-white">LAT</span>
                         <span>{activeDriverLocations[selectedShipment.driver_id!][0].toFixed(6)}</span>
                      </div>
                      <div className="flex justify-between">
                         <span className="opacity-50 text-white">LNG</span>
                         <span>{activeDriverLocations[selectedShipment.driver_id!][1].toFixed(6)}</span>
                      </div>
                   </div>
                </div>
              )}
           </div>

           <div className="p-6 bg-slate-50 border-t border-slate-100">
              <button 
                onClick={() => setSelectedShipment(null)}
                className="w-full bg-slate-900 text-white py-3 rounded-xl font-black uppercase text-[10px] tracking-widest flex items-center justify-center gap-2 hover:bg-primary transition-colors shadow-xl"
              >
                Close Details <X className="w-3 h-3" />
              </button>
           </div>
        </div>
      )}
    </div>
  );
}

function isUuid(value: string) {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(value);
}
