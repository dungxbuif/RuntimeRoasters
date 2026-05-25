'use client';

import React from 'react';
import { useAuth } from '@/lib/auth';
import { useQuery } from '@tanstack/react-query';
import { farmService } from '@/services/farm.service';
import { adminService } from '@/services/admin.service';
import { Coffee, Users, MapPin, Truck, AlertTriangle, TrendingUp } from 'lucide-react';
import { testId } from '@/lib/utils/test-id';

export default function AdminDashboardPage() {
  const { user } = useAuth();

  const { data: farms = [] } = useQuery({
    queryKey: ['farms'],
    queryFn: () => farmService.listFarms(),
  });

  const { data: users = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => adminService.listUsers(),
  });

  const stats = [
    { label: 'Pilot Farms', value: farms.length, icon: MapPin, color: 'text-primary' },
    { label: 'Active Personnel', value: users.length, icon: Users, color: 'text-tertiary' },
    { label: 'Retail Stores', value: 5, icon: Coffee, color: 'text-orange-500' },
    { label: 'Fleet Status', value: '3 Active', icon: Truck, color: 'text-emerald-500' },
  ];

  return (
    <div className="space-y-10">
      <div>
        <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-100 text-slate-500 text-[10px] font-black uppercase tracking-widest mb-4 border border-slate-200 shadow-sm italic">
          <span className="material-symbols-outlined text-[14px]">query_stats</span>
          <span>Global Visibility Engine</span>
        </div>
        <h1 className="text-5xl font-black font-headline text-on-surface uppercase tracking-tighter italic">
          Intelligence <span className="text-primary">Hub</span>
        </h1>
        <p className="text-on-surface-variant font-medium text-lg mt-2">
          System-wide oversight for the Runtime Roasters supply chain.
        </p>
      </div>

      {/* KPI Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => (
          <div key={stat.label} className="bg-white p-6 rounded-[2rem] border border-outline-variant/10 shadow-sm hover:shadow-xl transition-all group overflow-hidden relative">
            <div className={`absolute -right-4 -top-4 opacity-5 group-hover:scale-110 transition-transform duration-500 ${stat.color}`}>
              <stat.icon size={120} />
            </div>
            <div className="relative z-10">
              <div className={`w-12 h-12 rounded-2xl flex items-center justify-center mb-4 bg-slate-50 border border-outline-variant/5 shadow-inner ${stat.color}`}>
                <stat.icon size={20} />
              </div>
              <p className="text-[10px] font-black text-on-surface-variant uppercase tracking-[0.2em]">{stat.label}</p>
              <h3 className="text-3xl font-black font-headline text-on-surface mt-1">{stat.value}</h3>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* System Health Card */}
        <div className="lg:col-span-2 bg-white rounded-[2.5rem] border border-outline-variant/10 shadow-sm p-10 overflow-hidden relative">
           <div className="flex justify-between items-start mb-8">
              <div>
                <h3 className="text-2xl font-black font-headline uppercase italic tracking-tighter">Chain <span className="text-primary">Integrity</span></h3>
                <p className="text-xs font-bold text-slate-400 uppercase tracking-widest mt-1">Real-time infrastructure pulse</p>
              </div>
              <div className="flex gap-2">
                <div className="px-3 py-1 bg-emerald-50 text-emerald-600 rounded-full text-[9px] font-black uppercase tracking-widest border border-emerald-100 flex items-center gap-1.5 animate-pulse">
                  <div className="w-1 h-1 rounded-full bg-emerald-500"></div>
                  Systems Nominal
                </div>
              </div>
           </div>
           
           <div className="space-y-6">
              {[
                { name: 'Identity Service (Ory Stack)', status: 'Healthy', latency: '24ms' },
                { name: 'Kafka Cluster (Event Bus)', status: 'Healthy', latency: '12ms' },
                { name: 'Traceability Read Model', status: 'Healthy', latency: '45ms' },
              ].map((svc) => (
                <div key={svc.name} className="flex justify-between items-center p-4 bg-slate-50/50 rounded-2xl border border-outline-variant/5">
                  <div className="flex items-center gap-4">
                    <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
                    <span className="text-sm font-bold text-on-surface uppercase tracking-tight">{svc.name}</span>
                  </div>
                  <div className="text-right">
                    <span className="text-[10px] font-mono text-slate-400">{svc.latency}</span>
                  </div>
                </div>
              ))}
           </div>
        </div>

        {/* Alerts / Anomaly Feed */}
        <div className="bg-slate-900 text-white rounded-[2.5rem] p-10 shadow-2xl relative overflow-hidden">
           <div className="absolute top-0 right-0 p-4 opacity-10">
              <TrendingUp size={160} />
           </div>
           <h3 className="text-xl font-black font-headline uppercase italic tracking-tighter mb-8 flex items-center gap-3">
              <AlertTriangle className="text-amber-500" />
              <span>Anomaly <span className="text-primary">Watch</span></span>
           </h3>
           
           <div className="space-y-6">
              <div className="p-5 bg-white/5 rounded-2xl border border-white/10">
                <p className="text-[10px] font-black text-amber-500 uppercase tracking-widest mb-1">Observation Required</p>
                <p className="text-xs font-medium text-slate-300">No anomalies detected in the last 24 hours. System is operating within 12-20% roasting loss parameters.</p>
              </div>
              <div className="text-center pt-4">
                <button className="text-[10px] font-black uppercase tracking-widest text-slate-500 hover:text-white transition-colors">
                  View Full Audit Log
                </button>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
}
