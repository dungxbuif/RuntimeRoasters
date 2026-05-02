import React from 'react';

export default function FarmOriginDashboard() {
  return (
    <div className="min-h-full p-8 font-body relative overflow-x-hidden">
      {/* Hero Header */}
      <div className="mb-12 relative">
        <h1 className="text-5xl font-black font-headline text-on-surface tracking-tighter mb-4 uppercase italic">
          Farm & Origin <span className="text-primary">Dashboard</span>
        </h1>
        <p className="text-on-surface-variant max-w-2xl leading-relaxed font-medium">
          Real-time telemetry and architectural orchestration for the Runtime Roasters global supply chain. Tracking every bean from high-altitude plantation nodes to the event-driven cloud core.
        </p>
        <div className="absolute -right-4 -top-4 opacity-10">
          <span className="material-symbols-outlined text-[120px]">agriculture</span>
        </div>
      </div>

      {/* Dashboard Grid */}
      <div className="grid grid-cols-12 gap-8">
        {/* Plantation Node Map */}
        <div className="col-span-12 lg:col-span-8 bg-surface-container-lowest rounded-2xl p-8 border border-outline-variant/15 relative overflow-hidden shadow-sm hover:shadow-md transition-shadow">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h3 className="text-xl font-bold font-headline text-on-surface uppercase tracking-tight italic">Plantation Nodes</h3>
              <p className="text-sm text-on-surface-variant font-medium">Active Telemetry Map</p>
            </div>
            <div className="flex gap-2">
              <span className="px-3 py-1 bg-tertiary-fixed text-on-tertiary-fixed text-xs font-black rounded-full flex items-center gap-1 uppercase tracking-widest border border-tertiary/20">
                <span className="w-2 h-2 bg-tertiary rounded-full animate-pulse"></span> 14 Active Nodes
              </span>
            </div>
          </div>
          <div className="w-full h-[400px] rounded-xl bg-surface-container overflow-hidden relative border border-outline-variant/30 group">
            <img 
              alt="Plantation Map Location" 
              className="w-full h-full object-cover mix-blend-overlay opacity-40 grayscale group-hover:grayscale-0 transition-all duration-700" 
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuDyvk3bKTKao-vImJYm7bEm8YMpNFkX2MF5WOooNn5QVU7uE2Do1cTC6f8Fvn6vG57xKZ85l1KnV46ObhJtMJX4m7FyFnris9IiYBnaywNwTNUXBO91imsBkrkE76t6ZH4mzYn0G0ClZpj7I6yE3AGYPn1Xioq9tVeMJ0xO30tq7r6OKHjMkMCbVCoBB6vP_fJrAND3_HI_iNwIL6biAk30fk0CdqLQaAasmesHpeIiPayTpftr0sZ_zcHoLt5VPjAb5ryG4FliYScM" 
            />
            {/* Map Annotations */}
            <div className="absolute top-10 left-10 p-5 bg-white/70 backdrop-blur-xl rounded-xl border border-primary/20 max-w-[220px] shadow-2xl transform hover:-rotate-1 transition-transform">
              <p className="font-marker text-primary text-2xl rotate-[-2deg] mb-1">Node RR-ETH-04</p>
              <div className="h-1 bg-primary w-12 mb-3"></div>
              <p className="text-[10px] uppercase font-black text-on-surface-variant tracking-[0.2em] mb-1">Elevation</p>
              <p className="text-2xl font-headline font-black italic tracking-tighter">1,850m</p>
            </div>
            <div className="absolute bottom-12 right-12 flex flex-col items-end transform hover:scale-105 transition-transform">
              <p className="font-marker text-secondary text-2xl rotate-[3deg] mb-2 drop-shadow-sm">Multiple varieties detected</p>
              <svg className="w-32 h-16 opacity-60" viewBox="0 0 100 50">
                <path className="marker-line" d="M10,10 Q50,40 90,10" fill="none" stroke="#515f74" strokeWidth="2.5" strokeDasharray="5 3"></path>
              </svg>
            </div>
          </div>
        </div>

        {/* Harvest Batches Metrics */}
        <div className="col-span-12 lg:col-span-4 space-y-8">
          <div className="bg-surface-container-low rounded-3xl p-8 border border-outline-variant/15 shadow-sm">
            <h3 className="text-xl font-black font-headline mb-6 uppercase italic tracking-tighter border-b border-outline-variant/10 pb-4">Active Harvests</h3>
            <div className="space-y-6">
              {[
                { id: 'HB-9022', name: 'Yirgacheffe Grade 1', status: '82% Processed', color: 'text-primary', bg: 'bg-primary-fixed', elev: '1900m', moist: '11.4%', var: 'Heirloom' },
                { id: 'HB-8814', name: 'Guji Washed', status: '100% Processed', color: 'text-tertiary', bg: 'bg-tertiary-fixed', elev: '2100m', moist: '10.8%', var: 'Kurume' }
              ].map((batch, i) => (
                <div key={i} className="bg-surface-container-lowest p-6 rounded-2xl border border-outline-variant/10 hover:shadow-xl hover:border-primary/20 transition-all cursor-pointer group">
                  <div className="flex justify-between items-start mb-4">
                    <span className={`text-[10px] font-black ${batch.color} px-3 py-1 ${batch.bg} rounded-full uppercase tracking-widest`}>#{batch.id}</span>
                    <span className="text-[10px] text-on-surface-variant font-black uppercase tracking-tighter italic">{batch.status}</span>
                  </div>
                  <h4 className="font-black text-lg text-on-surface mb-4 uppercase tracking-tight">{batch.name}</h4>
                  <div className="grid grid-cols-3 gap-3">
                    <div className="text-center p-2.5 bg-surface rounded-xl border border-outline-variant/5">
                      <p className="text-[8px] text-on-surface-variant uppercase font-black">Elev.</p>
                      <p className="text-xs font-black italic">{batch.elev}</p>
                    </div>
                    <div className="text-center p-2.5 bg-surface rounded-xl border border-outline-variant/5 text-emerald-600">
                      <p className="text-[8px] text-on-surface-variant uppercase font-black">Moist.</p>
                      <p className="text-xs font-black italic">{batch.moist}</p>
                    </div>
                    <div className="text-center p-2.5 bg-surface rounded-xl border border-outline-variant/5">
                      <p className="text-[8px] text-on-surface-variant uppercase font-black">Var.</p>
                      <p className="text-xs font-black italic">{batch.var}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-primary p-8 rounded-3xl relative overflow-hidden group shadow-2xl hover:scale-[1.02] transition-transform cursor-help">
            <div className="relative z-10 text-white">
              <span className="material-symbols-outlined text-4xl mb-4 animate-bounce">auto_awesome</span>
              <h4 className="font-black text-2xl font-headline mb-2 uppercase italic tracking-tighter">Origin AI Insight</h4>
              <p className="text-sm font-medium opacity-90 leading-relaxed">Ideal humidity levels predicted for the next 72 hours across the highlands.</p>
            </div>
            <div className="absolute -bottom-10 -right-10 text-white/10 transition-transform group-hover:scale-125 duration-1000">
              <span className="material-symbols-outlined text-[160px]">psychology</span>
            </div>
          </div>
        </div>

        {/* Transactional Outbox Diagram */}
        <div className="col-span-12 bg-surface-container-lowest rounded-[2rem] p-12 border border-outline-variant/20 relative shadow-inner overflow-hidden">
          <div className="flex items-center gap-6 mb-12 relative z-10">
            <div className="w-16 h-16 rounded-2xl bg-secondary-container flex items-center justify-center text-on-secondary-container shadow-xl">
              <span className="material-symbols-outlined !text-3xl">hub</span>
            </div>
            <div>
              <h3 className="text-3xl font-black font-headline tracking-tighter uppercase italic">Transactional Outbox Pattern</h3>
              <p className="text-on-surface-variant font-medium">Ensuring eventual consistency from Farm to Cloud</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-5 gap-8 items-center relative z-10">
            {/* 1. Service */}
            <div className="bg-surface p-8 rounded-2xl border-4 border-primary border-dashed relative group">
              <h5 className="text-[10px] font-black uppercase tracking-widest mb-6 flex items-center gap-2 text-primary">
                <span className="material-symbols-outlined text-sm">dns</span> Origin Service
              </h5>
              <div className="space-y-3">
                <div className="bg-white p-3 rounded-lg shadow-sm text-[10px] border border-outline-variant/30 font-bold tracking-tighter italic">POST /harvest-event</div>
                <div className="bg-primary-fixed p-3 rounded-lg shadow-sm text-[10px] border border-primary/20 font-black text-primary uppercase tracking-tighter animate-pulse">DB Transaction Start</div>
              </div>
              <div className="absolute -right-6 top-1/2 -translate-y-1/2 z-10 hidden md:block">
                <span className="material-symbols-outlined text-primary text-5xl animate-pulse">trending_flat</span>
              </div>
            </div>

            {/* 2. Database */}
            <div className="md:col-span-2 bg-surface-container-low p-10 rounded-[2rem] border border-outline-variant relative shadow-2xl">
              <h5 className="text-[10px] font-black uppercase tracking-[0.3em] text-on-surface-variant mb-8 text-center opacity-60">PostgreSQL Cloud Instance</h5>
              <div className="flex gap-6">
                <div className="flex-1 bg-white p-6 rounded-2xl border-2 border-outline-variant shadow-inner opacity-60 grayscale scale-95">
                  <p className="text-[9px] font-black text-center mb-4 uppercase tracking-widest">Harvest Table</p>
                  <div className="h-2 bg-slate-100 rounded-full mb-2"></div>
                  <div className="h-2 bg-slate-100 rounded-full mb-2 w-3/4"></div>
                  <div className="h-2 bg-primary/20 rounded-full w-1/2"></div>
                </div>
                <div className="flex-1 bg-white p-6 rounded-2xl border-[3px] border-primary shadow-[0_20px_50px_rgba(0,74,198,0.15)] transform scale-105 relative">
                  <div className="absolute -top-3 -right-3 w-8 h-8 bg-primary rounded-full flex items-center justify-center text-white shadow-lg">
                    <span className="material-symbols-outlined !text-xs font-black">check</span>
                  </div>
                  <p className="text-[9px] font-black text-center mb-4 uppercase tracking-widest text-primary italic">Outbox Table</p>
                  <div className="h-2.5 bg-slate-100 rounded-full mb-2"></div>
                  <div className="h-2.5 bg-slate-100 rounded-full mb-2 w-[90%]"></div>
                  <div className="h-2.5 bg-primary rounded-full"></div>
                </div>
              </div>
              <p className="font-marker text-primary text-xl absolute -bottom-10 left-1/2 -translate-x-1/2 w-full text-center tracking-tighter">Atomic write guarantees data safety!</p>
            </div>

            {/* 3. Relay */}
            <div className="flex flex-col items-center gap-3 group">
              <div className="w-20 h-20 rounded-full bg-tertiary-fixed flex items-center justify-center border-4 border-tertiary shadow-xl group-hover:rotate-180 transition-transform duration-700">
                <span className="material-symbols-outlined text-tertiary !text-4xl animate-spin-slow">autorenew</span>
              </div>
              <p className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-800">Relay Worker</p>
              <p className="text-[9px] text-on-surface-variant font-bold italic text-center px-4">Polling 100ms</p>
            </div>

            {/* 4. Kafka */}
            <div className="bg-slate-950 p-8 rounded-2xl border-2 border-slate-700 shadow-[0_30px_60px_-12px_rgba(0,0,0,0.5)] relative">
              <div className="flex items-center gap-3 mb-6 border-b border-slate-800 pb-4">
                <div className="w-2.5 h-2.5 rounded-full bg-error animate-pulse shadow-[0_0_8px_#ba1a1a]"></div>
                <h5 className="text-xs font-black text-white uppercase tracking-[0.3em] italic">Apache Kafka</h5>
              </div>
              <div className="space-y-2">
                <div className="bg-slate-900 p-3 text-white text-[10px] rounded-lg flex justify-between border border-slate-800 font-mono">
                  <span className="opacity-50">TOPIC: origin.tel</span>
                  <span className="text-tertiary-fixed font-black tracking-widest">READY</span>
                </div>
                <div className="bg-slate-900 p-3 text-white text-[10px] rounded-lg border border-slate-800 font-mono tracking-tighter overflow-hidden whitespace-nowrap">
                   Msg: {"{batch_id: \"HB-9022\"}"}
                </div>
              </div>
            </div>
          </div>
          
          {/* Background Grid Pattern */}
          <div className="absolute inset-0 pointer-events-none opacity-[0.03]" style={{ backgroundImage: 'radial-gradient(#000 1px, transparent 1px)', backgroundSize: '32px 32px' }}></div>
        </div>
      </div>
    </div>
  );
}
