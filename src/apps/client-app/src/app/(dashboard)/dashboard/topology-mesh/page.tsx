'use client';

import React from 'react';
import { HighFidelityTopology } from '@/components/features/chaos-topology/high-fidelity/HighFidelityTopology';

export default function SystemTopologyPage() {
  return (
    <div className="flex-1 p-8 md:p-12 lg:p-20 bg-[#f7f9fb] relative overflow-hidden h-screen font-headline">
      {/* Grid Background Pattern */}
      <div className="absolute inset-0 pointer-events-none opacity-[0.03]" style={{ backgroundImage: 'radial-gradient(#000 2px, transparent 2px)', backgroundSize: '32px 32px' }}></div>
      
      <div className="max-w-7xl mx-auto h-full flex flex-col relative z-10">
        <div className="mb-12 flex justify-between items-end border-b-2 border-outline-variant/10 pb-8">
          <div>
            <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-primary-fixed/30 text-on-primary-fixed text-[10px] font-black uppercase tracking-[0.2em] mb-4 border border-primary/10 shadow-sm italic">
              <span className="material-symbols-outlined text-[14px]">account_tree</span>
              <span>L2 Network Visualization</span>
            </div>
            <h2 className="font-black text-5xl tracking-tighter text-on-surface uppercase italic">Network <span className="text-primary">Topology</span></h2>
            <p className="font-body text-on-surface-variant mt-4 font-medium tracking-tight max-w-xl leading-relaxed">High-fidelity microservices architecture orchestration map. v4.2.0-STABLE</p>
          </div>
        </div>

        <div className="flex-1 min-h-0 overflow-auto">
          <div className="min-w-[1400px] min-h-[800px] relative h-full">
            <HighFidelityTopology />
          </div>
        </div>
      </div>
    </div>
  );
}
