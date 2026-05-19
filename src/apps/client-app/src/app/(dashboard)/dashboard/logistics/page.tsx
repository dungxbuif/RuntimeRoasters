'use client';

import React, { useEffect, useState, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { logisticsService } from '@/services/logistics.service';
import { Location, Shipment, RouteData } from '@/types/logistics';
import LogisticsMap from '@/components/features/logistics/LogisticsMap';
import { Activity, Gauge, MapPin, Truck, AlertTriangle, CheckCircle2 } from 'lucide-react';

export default function LogisticsDashboard() {
  const [activeDriverLocations, setActiveDriverLocations] = useState<Record<string, [number, number]>>({});
  const [simulationActive, setSimulationActive] = useState(true);

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

  const { data: routes = [] } = useQuery({
    queryKey: ['logistics', 'routes'],
    queryFn: () => logisticsService.getRoutes(),
  });

  // Client-side simulation logic
  useEffect(() => {
    if (!simulationActive || routes.length === 0) return;

    const interval = setInterval(() => {
      setActiveDriverLocations(prev => {
        const next = { ...prev };
        
        // Simulate a few drivers on routes
        routes.forEach((route, idx) => {
          const driverId = `driver-${idx + 1}`;
          const currentPos = prev[driverId];
          let nextPos: [number, number];
          
          if (!currentPos) {
            nextPos = route.coordinates[0];
          } else {
            // Find current index in coordinates
            const currentIndex = route.coordinates.findIndex(
              p => p[0] === currentPos[0] && p[1] === currentPos[1]
            );
            
            if (currentIndex === -1 || currentIndex === route.coordinates.length - 1) {
              nextPos = route.coordinates[0]; // Reset to start
            } else {
              nextPos = route.coordinates[currentIndex + 1];
            }
          }

          next[driverId] = nextPos;

          // BRIDGE TO BACKEND: Update driver location in Valkey
          // We only do this if simulation is active to avoid overwhelming backend when paused
          const shipmentId = shipments.find(s => s.driver_id === driverId)?.id || '';
          logisticsService.updateDriverLocation({
            driver_id: driverId,
            shipment_id: shipmentId,
            latitude: nextPos[0],
            longitude: nextPos[1]
          }).catch(err => console.error(`Failed to update location for ${driverId}`, err));
        });
        
        return next;
      });
    }, 3000);

    return () => clearInterval(interval);
  }, [routes, simulationActive]);

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
        
        <div className="flex gap-4">
          <button 
            onClick={() => setSimulationActive(!simulationActive)}
            className={`px-4 py-2 rounded-lg font-black text-[10px] uppercase tracking-widest border transition-all ${
              simulationActive 
                ? 'bg-primary/10 border-primary text-primary shadow-[0_0_15px_rgba(0,74,198,0.2)]' 
                : 'bg-slate-100 border-slate-200 text-slate-400'
            }`}
          >
            {simulationActive ? '● SIMULATION ACTIVE' : '○ SIMULATION PAUSED'}
          </button>
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
                    className={`p-4 rounded-lg border-l-4 transition-all hover:translate-x-1 cursor-pointer ${
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
    </div>
  );
}
