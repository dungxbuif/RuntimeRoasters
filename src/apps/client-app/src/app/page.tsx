'use client';

import React, { useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import { Coffee, Truck, ArrowRight, ShieldCheck, Zap, Globe } from 'lucide-react';
import Link from 'next/link';
import { useSocket } from '@/hooks/useSocket';
import { useQuery } from '@tanstack/react-query';
import { logisticsService } from '@/services/logistics.service';

const LogisticsMap = dynamic(() => import('@/components/features/logistics/LogisticsMap'), {
  ssr: false,
  loading: () => <div className="w-full h-full bg-slate-900 animate-pulse flex items-center justify-center font-black text-slate-500 italic">Initializing Global Relay...</div>
});

export default function LandingPage() {
  const [activeDriverLocations, setActiveDriverLocations] = useState<Record<string, [number, number]>>({});

  const { data: initialLocations = [] } = useQuery({
    queryKey: ['public', 'driver-locations'],
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

  const { lastMessage } = useSocket(undefined, true);

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

  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col">
      {/* Navigation */}
      <nav className="border-b border-white/5 bg-slate-950/50 backdrop-blur-xl sticky top-0 z-[2000]">
        <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/20">
              <Coffee className="w-5 h-5 text-white" />
            </div>
            <div>
              <div className="text-sm font-black uppercase tracking-tight">Runtime Roasters</div>
              <div className="text-[9px] font-bold text-slate-500 uppercase tracking-widest">Next-Gen Supply Chain</div>
            </div>
          </div>
          <div className="flex items-center gap-8">
            <Link href="/course" className="text-[10px] font-black uppercase tracking-widest text-slate-400 hover:text-white transition-colors">Interactive Tour</Link>
            <Link href="/coffee" className="text-[10px] font-black uppercase tracking-widest text-slate-400 hover:text-white transition-colors">Coffee Kiosk</Link>
            <Link href="/trace" className="text-[10px] font-black uppercase tracking-widest text-slate-400 hover:text-white transition-colors">Traceability</Link>
            <Link href="/login" className="bg-white text-slate-950 px-5 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-primary hover:text-white transition-all shadow-xl">
              Operator Login
            </Link>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative h-[80vh] flex flex-col">
        <div className="absolute inset-0 z-0">
          <LogisticsMap 
            locations={[]}
            shipments={[]}
            routes={[]}
            activeDriverLocations={activeDriverLocations}
            mode="public"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/20 to-transparent pointer-events-none" />
          <div className="absolute inset-0 bg-gradient-to-r from-slate-950 via-transparent to-transparent pointer-events-none" />
        </div>

        <div className="relative z-10 max-w-7xl mx-auto px-6 pt-32 flex-1 flex flex-col justify-center">
          <div className="max-w-2xl space-y-8">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-[10px] font-black uppercase tracking-widest italic animate-pulse">
              <Globe className="w-3 h-3" /> Live Delivery Network Active
            </div>
            <h1 className="text-7xl font-black font-headline leading-[0.9] tracking-tighter uppercase italic">
              Real-time <br />
              <span className="text-primary">Transparency.</span>
            </h1>
            <p className="text-lg text-slate-400 font-medium max-w-lg">
              Witness the heartbeat of our supply chain. Every truck, every harvest, 
              every cup of coffee tracked in real-time across the Vietnamese highlands.
            </p>
            <div className="flex gap-4">
              <Link href="/coffee" className="bg-primary text-white px-8 py-4 rounded-2xl font-black uppercase tracking-widest hover:scale-105 transition-all flex items-center gap-3 shadow-2xl shadow-primary/30">
                Order a Traceable Cup <ArrowRight className="w-4 h-4" />
              </Link>
              <Link href="/course" className="bg-white/5 border border-white/10 text-white px-8 py-4 rounded-2xl font-black uppercase tracking-widest hover:bg-white/10 transition-all flex items-center gap-3 shadow-2xl">
                Interactive Tour <ArrowRight className="w-4 h-4" />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Feature Grid */}
      <section className="bg-slate-950 py-32 border-t border-white/5">
        <div className="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-3 gap-12">
          <div className="space-y-6">
            <div className="w-12 h-12 bg-white/5 rounded-2xl flex items-center justify-center border border-white/10">
              <ShieldCheck className="w-6 h-6 text-primary" />
            </div>
            <h3 className="text-xl font-black uppercase italic">Immutable Proof</h3>
            <p className="text-slate-500 text-sm leading-relaxed">
              Every stage of the journey is cryptographically signed and stored in our distributed read-model, providing absolute proof of origin.
            </p>
          </div>
          <div className="space-y-6">
            <div className="w-12 h-12 bg-white/5 rounded-2xl flex items-center justify-center border border-white/10">
              <Zap className="w-6 h-6 text-amber-500" />
            </div>
            <h3 className="text-xl font-black uppercase italic">Valkey GEO Edge</h3>
            <p className="text-slate-500 text-sm leading-relaxed">
              High-performance geospatial indexing at the edge ensures sub-millisecond latency for live tracking and routing.
            </p>
          </div>
          <div className="space-y-6">
            <div className="w-12 h-12 bg-white/5 rounded-2xl flex items-center justify-center border border-white/10">
              <Truck className="w-6 h-6 text-blue-500" />
            </div>
            <h3 className="text-xl font-black uppercase italic">Fleet Autonomy</h3>
            <p className="text-slate-500 text-sm leading-relaxed">
              Decentralized driver ingestion enables real-time synchronization between farms, roasteries, and retail outlets.
            </p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="mt-auto border-t border-white/5 py-12 bg-slate-950">
        <div className="max-w-7xl mx-auto px-6 flex justify-between items-center text-[10px] font-bold text-slate-600 uppercase tracking-widest">
          <p>© 2026 Runtime Roasters. All rights reserved.</p>
          <div className="flex gap-8">
            <Link href="/api-docs" className="hover:text-white transition-colors">API Reference</Link>
            <Link href="https://github.com" className="hover:text-white transition-colors">GitHub</Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
