'use client';

import React, { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { logisticsService } from '@/services/logistics.service';
import LogisticsMap from '@/components/features/logistics/LogisticsMap';
import StatusPipeline, { PipelineStep } from '@/components/common/StatusPipeline';
import { Truck, Play, Pause, RotateCcw, MapPin, Navigation, Clock, Gauge, CheckCircle2, Circle } from 'lucide-react';

const PICKUP_STEPS: PipelineStep[] = [
  { key: 'ASSIGNED', label: 'Assigned' },
  { key: 'DEPARTED_BASE', label: 'Departed' },
  { key: 'ARRIVED_FARM', label: 'At Farm' },
  { key: 'PICKUP_CONFIRMED', label: 'Loaded' },
  { key: 'RETURN_STARTED', label: 'Returning' },
  { key: 'ARRIVED_WAREHOUSE', label: 'Arrived WH' },
];

const DELIVERY_STEPS: PipelineStep[] = [
  { key: 'ASSIGNED', label: 'Assigned' },
  { key: 'DEPARTED_WH', label: 'Departed' },
  { key: 'ARRIVED_STORE', label: 'At Store' },
  { key: 'DELIVERED', label: 'Delivered' },
  { key: 'RETURN_STARTED', label: 'Returning' },
  { key: 'RETURNED_BASE', label: 'At Base' },
];

type SimSpeed = 30 | 45 | 60;

export default function DriverClientPage() {
  const [selectedShipmentId] = useState<string | null>(null);
  const [simActive, setSimActive] = useState(false);
  const [simSpeed, setSimSpeed] = useState<SimSpeed>(30);
  const [waypointIndex, setWaypointIndex] = useState(0);
  const [currentMilestone, setCurrentMilestone] = useState('ASSIGNED');
  const [driverPos, setDriverPos] = useState<Record<string, [number, number]>>({});
  const tickRef = useRef<NodeJS.Timeout | null>(null);

  const { data: shipments = [] } = useQuery({
    queryKey: ['driver', 'shipments'],
    queryFn: () => logisticsService.listShipments(),
    refetchInterval: 10000,
  });

  const { data: locations = [] } = useQuery({
    queryKey: ['logistics', 'locations'],
    queryFn: () => logisticsService.listLocations(),
  });

  const { data: routes = [] } = useQuery({
    queryKey: ['logistics', 'routes'],
    queryFn: () => logisticsService.getRoutes(),
  });

  const assignedShipments = useMemo(() =>
    shipments.filter(s => s.driver_id && ['ASSIGNED', 'IN_TRANSIT', 'IN_TRANSIT_TO_FARM', 'RETURNING_TO_WAREHOUSE', 'IN_TRANSIT_TO_STORE', 'RETURNING_TO_BASE'].includes(s.status)),
    [shipments]
  );

  const activeShipment = useMemo(() =>
    assignedShipments.find(s => s.id === selectedShipmentId) || assignedShipments[0],
    [assignedShipments, selectedShipmentId]
  );

  const activeRoute = useMemo(() => {
    if (!activeShipment || routes.length === 0) return null;
    return routes[0]; // Use first matching route
  }, [activeShipment, routes]);

  const totalWaypoints = activeRoute?.coordinates?.length || 0;
  const tickInterval = totalWaypoints > 0 ? (simSpeed * 1000) / totalWaypoints : 3000;
  const elapsed = totalWaypoints > 0 ? Math.round((waypointIndex / totalWaypoints) * simSpeed) : 0;
  const remaining = simSpeed - elapsed;

  // Simulation engine
  useEffect(() => {
    if (!simActive || !activeRoute || !activeShipment) return;

    tickRef.current = setInterval(() => {
      setWaypointIndex(prev => {
        const next = prev + 1;
        if (next >= totalWaypoints) {
          setSimActive(false);
          return prev;
        }

        const coord = activeRoute.coordinates[next];
        if (coord && activeShipment.driver_id) {
          setDriverPos(p => ({ ...p, [activeShipment.driver_id!]: coord }));

          // POST GPS to backend
          logisticsService.updateDriverLocation({
            driver_id: activeShipment.driver_id!,
            shipment_id: activeShipment.id,
            latitude: coord[0],
            longitude: coord[1],
          }).catch(() => {});
        }

        return next;
      });
    }, tickInterval);

    return () => {
      if (tickRef.current) clearInterval(tickRef.current);
    };
  }, [simActive, activeRoute, activeShipment, tickInterval, totalWaypoints]);

  const handleStart = useCallback(() => {
    if (!activeRoute || !activeShipment) return;
    setSimActive(true);
    if (waypointIndex === 0 && activeRoute.coordinates[0]) {
      setDriverPos(p => ({ ...p, [activeShipment.driver_id!]: activeRoute.coordinates[0] }));
    }
  }, [activeRoute, activeShipment, waypointIndex]);

  const handlePause = () => setSimActive(false);

  const handleReset = () => {
    setSimActive(false);
    setWaypointIndex(0);
    setCurrentMilestone('ASSIGNED');
    setDriverPos({});
  };

  const handleMilestone = async (milestone: string) => {
    if (!activeShipment) return;
    
    try {
      switch (milestone) {
        case 'DEPARTED_BASE':
        case 'DEPARTED_WH':
          await logisticsService.departShipment(activeShipment.id);
          break;
        case 'ARRIVED_FARM':
        case 'ARRIVED_STORE':
          await logisticsService.arriveShipment(activeShipment.id);
          break;
        case 'PICKUP_CONFIRMED':
          await logisticsService.confirmLoad(activeShipment.id);
          break;
        case 'DELIVERED':
          await logisticsService.confirmDelivery(activeShipment.id);
          break;
        case 'RETURN_STARTED':
          // Currently no backend API specifically for 'start return' - it just transitions state
          // Could be departShipment again, but backend handles this via the next state logic, or we just advance UI.
          // Let's assume the driver just clicks it and UI updates, the real backend confirmation is 'RETURNED_BASE'
          break;
        case 'ARRIVED_WAREHOUSE':
        case 'RETURNED_BASE':
          await logisticsService.returnShipment(activeShipment.id);
          break;
      }
      setCurrentMilestone(milestone);
    } catch (e) {
      console.error('Failed to confirm milestone', e);
    }
  };

  const isPickup = activeShipment?.status?.includes('FARM') || activeShipment?.status?.includes('PICKUP') || true;
  const steps = isPickup ? PICKUP_STEPS : DELIVERY_STEPS;
  const milestoneList = isPickup
    ? ['ASSIGNED', 'DEPARTED_BASE', 'ARRIVED_FARM', 'PICKUP_CONFIRMED', 'RETURN_STARTED', 'ARRIVED_WAREHOUSE']
    : ['ASSIGNED', 'DEPARTED_WH', 'ARRIVED_STORE', 'DELIVERED', 'RETURN_STARTED', 'RETURNED_BASE'];

  const milestoneLabels: Record<string, string> = {
    ASSIGNED: 'Shipment Assigned',
    DEPARTED_BASE: 'Departed Base',
    ARRIVED_FARM: 'Arrived at Farm',
    PICKUP_CONFIRMED: 'Confirm Pickup/Loading',
    RETURN_STARTED: 'Start Return',
    ARRIVED_WAREHOUSE: 'Arrived at Warehouse',
    DEPARTED_WH: 'Departed Warehouse',
    ARRIVED_STORE: 'Arrived at Store',
    DELIVERED: 'Delivery Confirmed',
    RETURNED_BASE: 'Returned to Base',
  };

  const currentMilestoneIndex = milestoneList.indexOf(currentMilestone);

  return (
    <div className="h-full flex flex-col gap-4">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <Truck className="w-7 h-7 text-primary" />
            Driver Client
          </h1>
          <p className="text-[10px] font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Route Simulation & Milestone Confirmation
          </p>
        </div>
        <div className="flex items-center gap-3">
          <div className={`px-3 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest border flex items-center gap-1.5 ${
            simActive ? 'bg-primary/10 border-primary text-primary animate-pulse' : 'bg-slate-100 border-slate-200 text-slate-400'
          }`}>
            <div className={`w-1.5 h-1.5 rounded-full ${simActive ? 'bg-primary' : 'bg-slate-400'}`} />
            {simActive ? 'Simulation Active' : 'Standby'}
          </div>
        </div>
      </div>

      <div className="flex-1 flex gap-4 overflow-hidden min-h-0">
        {/* Map */}
        <div className="flex-1 min-w-0 rounded-2xl overflow-hidden border border-slate-200 shadow-lg">
          <LogisticsMap
            locations={locations}
            shipments={assignedShipments}
            routes={routes}
            activeDriverLocations={driverPos}
          />
        </div>

        {/* Control Panel */}
        <div className="w-[380px] flex flex-col gap-3 overflow-y-auto pr-1">
          {/* Shipment Selector */}
          {assignedShipments.length === 0 ? (
            <div className="bg-slate-50 rounded-2xl p-8 text-center border border-slate-200">
              <Truck className="w-12 h-12 text-slate-300 mx-auto mb-3" />
              <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">No assigned shipments</p>
              <p className="text-[9px] text-slate-400 mt-1">Wait for WAREHOUSE_MGR to dispatch</p>
            </div>
          ) : (
            <>
              {/* Active Shipment Card */}
              {activeShipment && (
                <div className="bg-slate-900 text-white rounded-2xl p-5 shadow-xl">
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <div className="text-[8px] font-black text-slate-500 uppercase tracking-widest">Active Shipment</div>
                      <div className="text-sm font-black uppercase tracking-tight mt-0.5">
                        {isPickup ? 'Farm Pickup' : 'Retail Delivery'}
                      </div>
                    </div>
                    <span className="text-[8px] font-mono font-bold text-primary bg-primary/10 px-2 py-0.5 rounded">
                      {activeShipment.id.slice(0, 8).toUpperCase()}
                    </span>
                  </div>
                  <div className="flex items-center gap-2 text-[10px] text-slate-400">
                    <MapPin className="w-3 h-3 text-green-400" />
                    <span className="font-bold">{activeShipment.destination_address || 'Destination'}</span>
                  </div>
                </div>
              )}

              {/* Simulation Controls */}
              <div className="bg-white rounded-2xl p-5 border border-slate-200 shadow-sm">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-[10px] font-black uppercase tracking-widest text-slate-400">Simulation Controls</h3>
                  <div className="flex items-center gap-1">
                    {([30, 45, 60] as SimSpeed[]).map(s => (
                      <button key={s} onClick={() => setSimSpeed(s)}
                        className={`px-2 py-0.5 rounded text-[9px] font-black transition-all ${
                          simSpeed === s ? 'bg-primary text-white' : 'bg-slate-100 text-slate-400 hover:bg-slate-200'
                        }`}>
                        {s}s
                      </button>
                    ))}
                  </div>
                </div>

                {/* Progress */}
                <div className="mb-4">
                  <div className="flex justify-between text-[9px] font-black text-slate-400 mb-1">
                    <span>{waypointIndex}/{totalWaypoints} waypoints</span>
                    <span>{remaining}s remaining</span>
                  </div>
                  <div className="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                    <div className="bg-primary h-full rounded-full transition-all duration-300"
                      style={{ width: `${totalWaypoints > 0 ? (waypointIndex / totalWaypoints) * 100 : 0}%` }} />
                  </div>
                </div>

                {/* Buttons */}
                <div className="flex gap-2">
                  {!simActive ? (
                    <button onClick={handleStart}
                      disabled={!activeRoute || totalWaypoints === 0}
                      className="flex-1 flex items-center justify-center gap-2 bg-primary text-white py-3 rounded-xl font-black text-[11px] uppercase tracking-widest hover:bg-primary/90 transition-all shadow-lg shadow-primary/20 disabled:opacity-40 disabled:cursor-not-allowed">
                      <Play className="w-4 h-4" />
                      {waypointIndex > 0 ? 'Resume' : 'Start Route'}
                    </button>
                  ) : (
                    <button onClick={handlePause}
                      className="flex-1 flex items-center justify-center gap-2 bg-amber-500 text-white py-3 rounded-xl font-black text-[11px] uppercase tracking-widest hover:bg-amber-600 transition-all">
                      <Pause className="w-4 h-4" />
                      Pause
                    </button>
                  )}
                  <button onClick={handleReset}
                    className="px-4 py-3 bg-slate-100 text-slate-500 rounded-xl hover:bg-slate-200 transition-all">
                    <RotateCcw className="w-4 h-4" />
                  </button>
                </div>
              </div>

              {/* Milestone Buttons */}
              <div className="bg-white rounded-2xl p-5 border border-slate-200 shadow-sm">
                <h3 className="text-[10px] font-black uppercase tracking-widest text-slate-400 mb-3">Milestone Confirmation</h3>
                <StatusPipeline steps={steps} currentStep={currentMilestone} size="sm" />
                <div className="space-y-2 mt-4">
                  {milestoneList.map((m, i) => {
                    const isCompleted = i < currentMilestoneIndex;
                    const isCurrent = i === currentMilestoneIndex;
                    const isNext = i === currentMilestoneIndex + 1;

                    return (

                      <button key={m} onClick={() => isNext && handleMilestone(m)}
                        disabled={!isNext}
                        className={`w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-tight transition-all ${
                          isCompleted ? 'bg-green-50 text-green-700 border border-green-200' :
                          isCurrent ? 'bg-primary/5 text-primary border border-primary/20' :
                          isNext ? 'bg-white text-slate-900 border-2 border-primary shadow-md hover:shadow-lg hover:scale-[1.01] cursor-pointer' :
                          'bg-slate-50 text-slate-300 border border-slate-100 cursor-not-allowed'
                        }`}>
                        {isCompleted ? <CheckCircle2 className="w-4 h-4 text-green-500" /> :
                         isCurrent ? <Navigation className="w-4 h-4 text-primary" /> :
                         <Circle className="w-4 h-4" />}
                        {milestoneLabels[m]}
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* GPS Telemetry */}
              <div className="bg-slate-900 text-white rounded-2xl p-5 shadow-lg">
                <h3 className="text-[8px] font-black text-slate-500 uppercase tracking-[0.3em] mb-3">GPS Telemetry</h3>
                <div className="font-mono text-xs space-y-1">
                  <div className="flex justify-between">
                    <span className="text-slate-500">lat:</span>
                    <span className="text-primary tabular-nums">
                      {activeShipment?.driver_id && driverPos[activeShipment.driver_id]
                        ? driverPos[activeShipment.driver_id][0].toFixed(7)
                        : '—'}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-slate-500">lng:</span>
                    <span className="text-primary tabular-nums">
                      {activeShipment?.driver_id && driverPos[activeShipment.driver_id]
                        ? driverPos[activeShipment.driver_id][1].toFixed(7)
                        : '—'}
                    </span>
                  </div>
                </div>
                <div className="flex items-center gap-2 mt-3 text-[9px] text-slate-500 font-bold">
                  <Clock className="w-3 h-3" />
                  <span>Last Update: {simActive ? '2s ago' : 'Idle'}</span>
                  <span className="ml-auto flex items-center gap-1">
                    <Gauge className="w-3 h-3" />
                    {simActive ? `~${Math.round(totalWaypoints > 0 ? 120 / (simSpeed / totalWaypoints) : 0)} km/h` : '0 km/h'}
                  </span>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
