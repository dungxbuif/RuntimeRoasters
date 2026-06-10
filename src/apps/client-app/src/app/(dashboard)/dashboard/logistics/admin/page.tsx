'use client';

import React, { useEffect, useState, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { logisticsService } from '@/services/logistics.service';
import { Shipment } from '@/types/logistics';
import dynamic from 'next/dynamic';
import { Activity, Gauge, Truck, AlertTriangle, ShieldCheck, Zap, Globe, BarChart3, Database } from 'lucide-react';
import { useSocket } from '@/hooks/useSocket';

const LogisticsMap = dynamic(() => import('@/components/features/logistics/LogisticsMap'), {
  ssr: false,
  loading: () => <div className="w-full h-full bg-slate-900 animate-pulse flex items-center justify-center font-black text-slate-500 italic">Synchronizing Fleet Intelligence...</div>
});

export default function LogisticsAdminDashboard() {
  const [activeDriverLocations, setActiveDriverLocations] = useState<Record<string, [number, number]>>({});

  const { data: locations = [] } = useQuery({
    queryKey: ['logistics', 'locations'],
    queryFn: () => logisticsService.listLocations(),
  });

  const { data: shipments = [] } = useQuery({
    queryKey: ['logistics', 'shipments'],
    queryFn: () => logisticsService.listShipments(),
    refetchInterval: 10000,
  });

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

  const { lastMessage, isConnected } = useSocket();

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

  const activeShipments = shipments.filter(s => s.status === 'IN_TRANSIT' || s.status.includes('IN_TRANSIT')).length;

  return (
    <div className="h-full flex flex-col gap-6">
      {/* Admin Header */}
      <div className="flex justify-between items-end">
        <div>
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-red-500/10 border border-red-500/20 text-red-500 text-[10px] font-black uppercase tracking-widest mb-4 italic">
            <ShieldCheck className="w-3 h-3" /> Fleet Command & Control
          </div>
          <h1 className="text-4xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <Globe className="w-10 h-10 text-primary" />
            Global Logistics Ops
          </h1>
        </div>
        
        <div className="flex gap-3">
          <div className={`px-4 py-2 rounded-xl border flex items-center gap-3 transition-all ${
            isConnected ? 'bg-emerald-50 border-emerald-200 text-emerald-600' : 'bg-red-50 border-red-200 text-red-600'
          }`}>
            <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-emerald-500 animate-pulse' : 'bg-red-500'}`} />
            <span className="text-[10px] font-black uppercase tracking-widest">
              {isConnected ? 'Real-time Stream: Active' : 'Relay Disconnected'}
            </span>
          </div>
        </div>
      </div>

      {/* Admin Metrics Grid */}
      <div className="grid grid-cols-4 gap-6">
        <MetricCard icon={<Truck />} label="Active Fleet" value={activeShipments} subValue="Trucks in Motion" color="blue" />
        <MetricCard icon={<Zap />} label="Avg Latency" value="14ms" subValue="Valkey GEO Edge" color="amber" />
        <MetricCard icon={<Activity />} label="Throughput" value="8.4" subValue="Events / Second" color="emerald" />
        <MetricCard icon={<Database />} label="Nodes" value="12" subValue="Active Instances" color="slate" />
      </div>

      <div className="flex-1 flex gap-6 overflow-hidden min-h-0">
        {/* Central Map */}
        <div className="flex-1 min-w-0">
          <LogisticsMap 
            locations={locations}
            shipments={shipments}
            routes={[]}
            activeDriverLocations={activeDriverLocations}
            mode="admin"
            showRoutes={true}
          />
        </div>

        {/* Admin Sidebar */}
        <div className="w-96 flex flex-col gap-6">
          <div className="bg-slate-900 text-white rounded-3xl p-8 shadow-2xl flex flex-col gap-8">
            <h2 className="text-[10px] font-black uppercase tracking-[0.3em] text-slate-500 flex items-center gap-2">
              <BarChart3 className="w-4 h-4 text-primary" />
              System Intelligence
            </h2>
            
            <div className="space-y-8">
              <SystemMetric label="Kafka Ingestion Rate" value="98.2%" color="bg-primary" />
              <SystemMetric label="PostgreSQL Sync" value="100%" color="bg-emerald-500" />
              <SystemMetric label="Valkey Memory Usage" value="1.2GB" color="bg-amber-500" />
            </div>

            <div className="mt-auto pt-8 border-t border-white/5">
              <div className="p-4 bg-white/5 rounded-2xl border border-white/10">
                <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Operational Note</p>
                <p className="text-xs font-medium text-slate-300">
                  All geospatial data is being streamed via persistent WebSockets to minimize client-side overhead.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function MetricCard({ icon, label, value, subValue, color }: any) {
  const colors: any = {
    blue: 'text-blue-500 bg-blue-50 border-blue-100',
    amber: 'text-amber-500 bg-amber-50 border-amber-100',
    emerald: 'text-emerald-500 bg-emerald-50 border-emerald-100',
    slate: 'text-slate-500 bg-slate-50 border-slate-100',
  };
  
  return (
    <div className="bg-white p-6 rounded-3xl border border-slate-200 shadow-sm flex items-center gap-6">
      <div className={`w-14 h-14 rounded-2xl flex items-center justify-center border ${colors[color]}`}>
        {React.cloneElement(icon, { size: 24 })}
      </div>
      <div>
        <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">{label}</p>
        <h3 className="text-2xl font-black text-slate-900 leading-tight">{value}</h3>
        <p className="text-[9px] font-bold text-slate-400 uppercase italic">{subValue}</p>
      </div>
    </div>
  );
}

function SystemMetric({ label, value, color }: any) {
  return (
    <div className="space-y-2">
      <div className="flex justify-between items-end">
        <span className="text-[10px] font-black uppercase text-slate-400">{label}</span>
        <span className="text-xs font-black italic">{value}</span>
      </div>
      <div className="w-full bg-white/5 h-1.5 rounded-full overflow-hidden">
        <div className={`h-full ${color}`} style={{ width: value.includes('%') ? value : '70%' }} />
      </div>
    </div>
  );
}
