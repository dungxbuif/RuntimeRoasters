import React from 'react';

export default function ResiliencyFlowPage() {
  return (
    <div className="h-full flex flex-col font-body bg-surface text-on-surface overflow-hidden">
      {/* Header Section */}
      <div className="p-8 lg:p-12 bg-white shadow-sm border-b border-outline-variant/10">
        <div className="max-w-6xl mx-auto flex flex-col md:flex-row justify-between items-start md:items-end gap-6">
          <div>
            <div className="flex items-center gap-2 mb-4">
              <span className="bg-error-container text-on-error-container text-[10px] font-black px-3 py-1 rounded-full uppercase tracking-widest border border-error/10">Analysis Mode</span>
              <span className="text-on-surface-variant text-[10px] font-bold uppercase tracking-widest italic opacity-40">Resiliency_Engine: v2.4</span>
            </div>
            <h1 className="font-headline text-5xl font-black tracking-tighter uppercase italic">Resiliency <span className="text-primary">Flow</span></h1>
            <p className="text-on-surface-variant text-lg max-w-2xl font-medium mt-4 leading-relaxed">
              Visualizing infrastructure failure scenarios and automated fallback patterns within the distributed system core.
            </p>
          </div>
          <div className="flex gap-4">
            <button className="px-8 py-4 bg-on-background text-white font-black text-xs uppercase tracking-widest rounded-2xl hover:scale-95 transition-all shadow-2xl flex items-center gap-3 italic">
              <span className="material-symbols-outlined !text-sm">restart_alt</span>
              Reset System State
            </button>
          </div>
        </div>
      </div>

      {/* Content Grid */}
      <div className="flex-1 p-8 lg:p-12 overflow-y-auto">
        <div className="max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12">
          
          {/* Controls Panel */}
          <div className="lg:col-span-4 space-y-8">
            <div className="bg-surface-container-lowest rounded-[2rem] p-10 border border-outline-variant/15 shadow-sm relative overflow-hidden group">
              <div className="absolute top-0 left-0 w-2 h-full bg-error"></div>
              <h3 className="font-headline font-black text-xs uppercase tracking-[0.3em] text-on-surface-variant opacity-40 mb-10 flex items-center gap-3 italic">
                <span className="material-symbols-outlined text-error !text-xl animate-pulse">warning</span>
                Simulated Faults
              </h3>
              
              <div className="space-y-6">
                {[
                  { label: 'Primary DB', sub: 'PostgreSQL Cluster', icon: 'database', color: 'text-primary', active: false },
                  { label: 'Kafka Broker', sub: 'OUTAGE_DETECTED', icon: 'hub', color: 'text-error', active: true, fail: true },
                  { label: 'Payment API', sub: 'External Gateway', icon: 'api', color: 'text-secondary', active: false }
                ].map((fault, i) => (
                  <div key={i} className={`flex items-center justify-between p-5 rounded-2xl border transition-all cursor-pointer group/item shadow-sm ${fault.fail ? 'bg-error-container/20 border-error/30' : 'bg-surface border-outline-variant/5 hover:border-primary/30'}`}>
                    <div className="flex items-center gap-5">
                      <div className={`w-12 h-12 rounded-full flex items-center justify-center shadow-inner group-hover/item:scale-110 transition-transform ${fault.fail ? 'bg-error text-white animate-pulse' : 'bg-surface-container-high text-on-surface-variant opacity-60'}`}>
                        <span className="material-symbols-outlined !text-2xl">{fault.icon}</span>
                      </div>
                      <div>
                        <div className={`font-black text-xs uppercase tracking-tight ${fault.fail ? 'text-error' : 'text-on-surface'}`}>{fault.label}</div>
                        <div className={`text-[10px] font-bold ${fault.fail ? 'text-error/60 italic' : 'text-on-surface-variant opacity-40'}`}>{fault.sub}</div>
                      </div>
                    </div>
                    <div className={`w-12 h-6 rounded-full relative p-1 transition-colors ${fault.active ? 'bg-error' : 'bg-outline-variant/30'}`}>
                      <div className={`w-4 h-4 bg-white rounded-full shadow-md absolute transition-all ${fault.active ? 'right-1' : 'left-1'}`}></div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Handwritten note element */}
            <div className="bg-amber-50 rounded-[2rem] p-10 border-2 border-amber-200 shadow-xl transform rotate-1 relative overflow-hidden">
              <div className="absolute top-0 right-0 p-6 opacity-10">
                <span className="material-symbols-outlined !text-[100px]">sticky_note_2</span>
              </div>
              <h4 className="font-marker text-4xl text-amber-700 mb-4 tracking-tighter italic">Engine Status</h4>
              <p className="font-body text-base text-amber-900 font-bold leading-relaxed tracking-tight">
                Kafka cluster is currently unresponsive. Observing fallback behavior triggered by the Outbox Pattern. Data remains consistent in local persistence layer.
              </p>
              <div className="mt-6 flex items-center gap-3">
                <div className="w-2.5 h-2.5 rounded-full bg-error animate-ping"></div>
                <span className="font-mono text-[10px] font-black text-error uppercase tracking-widest italic">Emergency Relay Active</span>
              </div>
            </div>
          </div>

          {/* Resiliency Diagram Area */}
          <div className="lg:col-span-8 bg-surface-container-lowest rounded-[3rem] p-16 border border-outline-variant/15 shadow-xl relative min-h-[650px] flex items-center justify-center overflow-hidden">
            {/* Diagram Background Grid */}
            <div className="absolute inset-0 pointer-events-none opacity-[0.02]" style={{ backgroundImage: 'radial-gradient(#000 1px, transparent 1px)', backgroundSize: '40px 40px' }}></div>

            <div className="relative w-full max-w-2xl aspect-video scale-110">
              {/* Primary Service Node */}
              <div className="absolute top-1/2 left-0 -translate-y-1/2 w-48 h-36 bg-on-background border-4 border-slate-700 rounded-[2.5rem] flex flex-col items-center justify-center z-20 shadow-2xl transition-transform hover:rotate-1 italic">
                <span className="material-symbols-outlined text-primary-fixed !text-5xl mb-3 drop-shadow-xl">shopping_cart</span>
                <span className="font-headline text-xs font-black text-white uppercase tracking-widest">Order Service</span>
                <div className="mt-3 flex gap-1.5 opacity-30">
                  {[1, 2, 3].map(i => <div key={i} className="w-3 h-1 bg-primary-fixed rounded-full"></div>)}
                </div>
              </div>

              {/* FAILED PATH VISUAL */}
              <svg className="absolute top-1/2 left-[140px] w-64 h-20 -translate-y-1/2 z-10 overflow-visible opacity-60">
                <path className="marker-path" d="M 0 10 L 220 10" fill="none" stroke="#ba1a1a" strokeWidth="4" strokeDasharray="10 8" style={{ filter: 'url(#rough-filter)' }}></path>
                <g transform="translate(110, 10)" className="text-error drop-shadow-xl">
                  <line x1="-15" y1="-15" x2="15" y2="15" stroke="currentColor" strokeWidth="6" strokeLinecap="round"></line>
                  <line x1="15" y1="-15" x2="-15" y2="15" stroke="currentColor" strokeWidth="6" strokeLinecap="round"></line>
                </g>
              </svg>

              {/* Target Service Node (FAILED) */}
              <div className="absolute top-0 right-0 w-48 h-36 bg-error-container/20 border-4 border-error border-dashed rounded-[2.5rem] flex flex-col items-center justify-center z-20 opacity-80 shadow-2xl italic">
                <span className="material-symbols-outlined text-error !text-5xl mb-3">hub</span>
                <span className="font-headline text-xs font-black text-error uppercase tracking-widest">Kafka Bus</span>
                <span className="text-[10px] font-black bg-error text-white px-4 py-1 rounded-full uppercase mt-3 tracking-tighter shadow-lg">Failed</span>
              </div>

              {/* FALLBACK PATH VISUAL */}
              <svg className="absolute top-1/2 left-[160px] w-64 h-56 z-10 overflow-visible">
                <path className="marker-path" d="M 0 10 Q 60 180 200 180" fill="none" stroke="#006242" strokeWidth="5" strokeDasharray="12 6" style={{ filter: 'url(#rough-filter)' }}></path>
                <path d="M 190 172 L 205 180 L 190 188" fill="none" stroke="#006242" strokeWidth="4" strokeLinecap="round" strokeLinejoin="round"></path>
              </svg>

              {/* Fallback Persistence Node */}
              <div className="absolute bottom-0 right-0 w-48 h-36 bg-tertiary-fixed/30 border-4 border-tertiary rounded-[2.5rem] flex flex-col items-center justify-center z-20 shadow-2xl group transition-all hover:scale-105 italic">
                <div className="absolute -top-4 -right-4 w-12 h-12 bg-tertiary rounded-full flex items-center justify-center text-white shadow-xl group-hover:rotate-12 transition-transform border-4 border-surface">
                  <span className="material-symbols-outlined !text-xl font-black">check</span>
                </div>
                <span className="material-symbols-outlined text-tertiary !text-5xl mb-3 drop-shadow-lg">table_rows</span>
                <span className="font-headline text-xs font-black text-on-tertiary-fixed-variant uppercase tracking-widest text-center px-6 leading-tight italic">Outbox Store</span>
              </div>

              {/* Annotations */}
              <div className="absolute top-[20%] left-[35%] font-marker text-error text-4xl -rotate-12 drop-shadow-sm font-black whitespace-nowrap opacity-80">
                Connection Refused!
              </div>
              <div className="absolute bottom-[20%] left-[25%] font-marker text-tertiary text-4xl rotate-2 drop-shadow-sm flex items-center gap-6 font-black">
                <span className="material-symbols-outlined !text-5xl animate-bounce">subdirectory_arrow_right</span>
                Fallback pattern
              </div>

              {/* Status Badge */}
              <div className="absolute -bottom-24 left-1/2 -translate-x-1/2 bg-on-background text-white px-10 py-5 rounded-[2rem] shadow-2xl border-2 border-primary/30 flex items-center gap-6 group hover:bg-slate-800 transition-colors italic">
                <div className="w-3.5 h-3.5 bg-tertiary-fixed rounded-full animate-pulse shadow-[0_0_12px_rgba(78,222,163,0.8)]"></div>
                <span className="font-black text-sm uppercase tracking-[0.2em]">Resiliency: <span className="text-tertiary-fixed">Active Fallback</span></span>
                <span className="material-symbols-outlined text-slate-500 text-xl opacity-40">info</span>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}
