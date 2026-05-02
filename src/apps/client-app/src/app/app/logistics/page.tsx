import React from 'react';

export default function LogisticsRealTimeTransitPage() {
  return (
    <div className="h-full relative flex font-body bg-slate-950 overflow-hidden">
      {/* Map Background (Simulated) */}
      <div className="absolute inset-0 z-0 bg-slate-900 overflow-hidden">
        <div 
          className="absolute inset-0 opacity-20 bg-[url('https://images.unsplash.com/photo-1524661135-423995f22d0b?auto=format&fit=crop&q=80')] bg-cover bg-center" 
          style={{ filter: 'grayscale(100%) contrast(120%) brightness(50%)' }}
        ></div>
        <div className="absolute inset-0 bg-gradient-to-b from-slate-950/80 via-transparent to-slate-950/90"></div>
        
        {/* Simulated Glowing Paths and Nodes */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none" xmlns="http://www.w3.org/2000/svg">
          <path 
            className="animate-dash" 
            d="M 300 400 C 400 300, 600 500, 800 350" 
            fill="none" 
            stroke="rgba(34, 197, 94, 0.3)" 
            strokeDasharray="4 4" 
            strokeWidth="2"
          ></path>
          <path d="M 300 400 C 400 300, 600 500, 800 350" fill="none" stroke="#22c55e" strokeWidth="1"></path>
          <circle cx="300" cy="400" fill="#22c55e" r="4" className="shadow-[0_0_10px_#22c55e]"></circle>
          <circle cx="800" cy="350" fill="#00daf3" r="6" className="shadow-[0_0_15px_#00daf3]"></circle>
          <circle cx="500" cy="200" fill="#f59e0b" r="4" className="shadow-[0_0_10px_#f59e0b]"></circle>
          <g transform="translate(550, 410)">
            <circle cx="0" cy="0" fill="#0f172a" r="8" stroke="#22c55e" strokeWidth="2"></circle>
            <text fill="#22c55e" fontSize="10" textAnchor="middle" x="0" y="3.5" className="material-symbols-outlined">local_shipping</text>
          </g>
        </svg>
      </div>

      {/* Floating Command Panels */}
      <div className="relative z-10 w-full h-full flex p-6 gap-6 overflow-hidden">
        {/* Left Column: Active Shipments */}
        <div className="w-[420px] flex flex-col gap-4 h-full">
          <div className="glass-card rounded-xl p-6 border border-slate-800 shadow-2xl flex flex-col h-full overflow-hidden bg-slate-900/60 backdrop-blur-xl">
            <div className="flex items-center justify-between mb-6 pb-4 border-b border-slate-800">
              <h2 className="font-headline text-xs text-slate-400 uppercase tracking-[0.2em] flex items-center gap-2 font-bold">
                <span className="material-symbols-outlined text-blue-400 text-lg">radar</span>
                Transit Telemetry
              </h2>
              <span className="text-[10px] bg-slate-800 px-2 py-1 rounded text-slate-300 tabular-nums font-black tracking-widest italic border border-slate-700">4 EN ROUTE</span>
            </div>
            
            <div className="flex-1 overflow-y-auto pr-2 space-y-4">
              {/* Shipment Card 1 */}
              <div className="bg-slate-950/40 rounded-lg p-4 border-l-4 border-primary hover:bg-slate-800/80 transition-all cursor-pointer group">
                <div className="flex justify-between items-start mb-3 text-white">
                  <div className="flex items-center gap-2">
                    <span className="material-symbols-outlined text-primary text-lg group-hover:scale-110 transition-transform font-black">local_shipping</span>
                    <span className="font-headline text-sm font-black tracking-tight uppercase italic">SHP-8924-A</span>
                  </div>
                  <span className="text-[9px] bg-primary/10 text-primary px-2 py-0.5 rounded-full uppercase font-bold tracking-tighter border border-primary/20 italic">En Route</span>
                </div>
                <div className="grid grid-cols-2 gap-3 text-[10px] font-bold uppercase tracking-tighter">
                  <div className="bg-slate-950/80 rounded p-2 flex flex-col border border-slate-800">
                    <span className="text-slate-500 text-[8px] mb-1">Origin</span>
                    <span className="text-slate-300 truncate">Farm Node 04</span>
                  </div>
                  <div className="bg-slate-950/80 rounded p-2 flex flex-col border border-slate-800">
                    <span className="text-slate-500 text-[8px] mb-1">Destination</span>
                    <span className="text-slate-300 truncate">WH-Alpha</span>
                  </div>
                </div>
                <div className="mt-3 flex items-center justify-between bg-slate-950/80 rounded p-2.5 border border-slate-800 group-hover:border-primary/30 transition-colors text-white">
                  <div className="flex items-center gap-2">
                    <span className="material-symbols-outlined text-tertiary-fixed text-[16px]">thermostat</span>
                    <span className="text-tertiary-fixed font-black tabular-nums text-xs italic">18.4°C</span>
                  </div>
                  <div className="flex items-center gap-2 text-slate-500 font-black">
                    <span className="material-symbols-outlined text-[16px]">speed</span>
                    <span className="tabular-nums text-[10px]">64 km/h</span>
                  </div>
                </div>
              </div>

              {/* Shipment Card 2 */}
              <div className="bg-slate-950/40 rounded-lg p-4 border-l-4 border-error hover:bg-slate-800/80 transition-all cursor-pointer text-white">
                <div className="flex justify-between items-start mb-3">
                  <div className="flex items-center gap-2">
                    <span className="material-symbols-outlined text-error text-lg font-black">warning</span>
                    <span className="font-headline text-sm font-black tracking-tight uppercase italic">SHP-7712-B</span>
                  </div>
                  <span className="text-[9px] bg-error/10 text-error px-2 py-0.5 rounded-full uppercase font-bold tracking-tighter border border-error/20 italic">Delayed</span>
                </div>
                <div className="mt-3 flex items-center justify-between bg-slate-950/80 rounded p-2.5 border border-error/20 shadow-[0_0_10px_rgba(186,26,26,0.1)]">
                  <div className="flex items-center gap-2 font-black">
                    <span className="material-symbols-outlined text-error text-[16px] animate-pulse">thermostat</span>
                    <span className="text-error tabular-nums text-xs italic">24.1°C</span>
                  </div>
                  <span className="text-[8px] text-error uppercase font-black italic tracking-widest">Delta +2.1°</span>
                </div>
              </div>

              {/* Shipment Card 3 */}
              <div className="bg-slate-950/40 rounded-lg p-4 border-l-4 border-slate-700 opacity-60 text-slate-400 italic">
                <div className="flex justify-between items-start mb-3">
                  <div className="flex items-center gap-2 font-black">
                    <span className="material-symbols-outlined text-slate-500 text-lg">done_all</span>
                    <span className="font-headline text-sm tracking-tight uppercase">SHP-6641-X</span>
                  </div>
                  <span className="text-[9px] bg-slate-800 text-slate-500 px-2 py-0.5 rounded-full uppercase font-bold tracking-tighter">Delivered</span>
                </div>
                <div className="text-[10px] flex items-center gap-1 font-bold">
                  <span className="material-symbols-outlined text-[14px]">schedule</span> Offloaded 14m ago
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Micro-stats & Map Overlay Controls */}
        <div className="flex-1 flex flex-col justify-end items-end gap-6 pointer-events-none">
          {/* Map Controls */}
          <div className="glass-card rounded-lg flex flex-col p-1.5 border border-slate-800 shadow-2xl pointer-events-auto bg-slate-900/80">
            <button className="p-2.5 text-slate-500 hover:text-white hover:bg-slate-800 rounded transition-colors group">
              <span className="material-symbols-outlined !text-lg group-active:scale-90 transition-transform">add</span>
            </button>
            <div className="h-px w-full bg-slate-800 my-1"></div>
            <button className="p-2.5 text-slate-500 hover:text-white hover:bg-slate-800 rounded transition-colors group">
              <span className="material-symbols-outlined !text-lg group-active:scale-90 transition-transform">remove</span>
            </button>
            <div className="h-px w-full bg-slate-800 my-1"></div>
            <button className="p-2.5 text-primary bg-primary/10 rounded transition-colors shadow-inner">
              <span className="material-symbols-outlined !text-lg font-black">my_location</span>
            </button>
          </div>
          
          {/* Warehouse Stock Mini-View */}
          <div className="glass-card rounded-xl p-8 border border-slate-800 shadow-2xl w-[360px] pointer-events-auto transform hover:translate-y-[-4px] transition-transform bg-slate-900/90 text-white">
            <h3 className="font-headline text-[10px] font-black text-slate-500 uppercase tracking-[0.3em] mb-8 flex items-center gap-3">
              <span className="material-symbols-outlined text-primary text-lg font-black">inventory_2</span>
              Node Capacities
            </h3>
            <div className="space-y-8 font-black uppercase tracking-tighter italic">
              <div>
                <div className="flex justify-between text-[11px] mb-2">
                  <span className="text-slate-300">WH-Alpha (Primary)</span>
                  <span className="text-primary tabular-nums">82% [LOAD]</span>
                </div>
                <div className="w-full bg-slate-950 rounded-full h-2 overflow-hidden flex border border-slate-800 shadow-inner">
                  <div className="bg-primary h-full shadow-[0_0_12px_rgba(0,74,198,0.6)]" style={{ width: '82%' }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between text-[11px] mb-2">
                  <span className="text-slate-300">WH-Beta (Transit)</span>
                  <span className="text-error tabular-nums">94% [CRIT]</span>
                </div>
                <div className="w-full bg-slate-950 rounded-full h-2 overflow-hidden flex border border-slate-800 shadow-inner">
                  <div className="bg-error h-full shadow-[0_0_12px_rgba(186,26,26,0.6)] animate-pulse" style={{ width: '94%' }}></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
