import React from 'react';

export default function RetailSagaOrchestratorPage() {
  return (
    <div className="min-h-full bg-surface text-on-surface p-8 md:p-12 lg:p-16 font-body selection:bg-primary-fixed selection:text-on-primary-fixed relative overflow-x-hidden">
      {/* SVG Filters for Sketchy Effects */}
      <svg className="absolute hidden" height="0" width="0">
        <filter id="retail-rough-edge">
          <feTurbulence baseFrequency="0.05" numOctaves="3" result="noise" type="fractalNoise"></feTurbulence>
          <feDisplacementMap in="SourceGraphic" in2="noise" scale="3"></feDisplacementMap>
        </filter>
      </svg>

      <div className="max-w-7xl mx-auto space-y-10 relative z-10">
        {/* Hero Section */}
        <section className="flex justify-between items-end border-b border-outline-variant/30 pb-8">
          <div>
            <h1 className="text-4xl font-black font-headline tracking-tighter text-on-surface uppercase italic">Saga Orchestrator</h1>
            <p className="text-on-surface-variant mt-2 font-medium tracking-widest text-xs uppercase">Distributed Transaction Monitor // Retail Core</p>
          </div>
          <div className="flex gap-4">
            <span className="px-3 py-1 bg-error-container text-on-error-container rounded-full text-[10px] font-black flex items-center gap-1 border border-error/20 shadow-sm">
              <span className="material-symbols-outlined !text-sm animate-pulse">warning</span> ROLLBACK_ACTIVE
            </span>
          </div>
        </section>

        {/* Hero Stats */}
        <section className="grid grid-cols-1 md:grid-cols-4 gap-6">
          {[
            { label: 'Active Sagas', val: '1,284', trend: '+12%', trendColor: 'text-tertiary' },
            { label: 'Failure Rate', val: '0.42%', trend: 'Stable', trendColor: 'text-on-surface-variant' },
            { label: 'Avg. Latency', val: '342ms', trend: '+14ms', trendColor: 'text-error' },
            { label: 'Rollback Count', val: '18', trend: 'Last 24h', trendColor: 'text-on-surface-variant' }
          ].map((stat, i) => (
            <div key={i} className="bg-surface-container-lowest p-6 rounded-2xl border border-outline-variant/10 shadow-sm flex flex-col justify-between h-32 hover:shadow-md transition-shadow">
              <span className="text-on-surface-variant font-headline text-[10px] font-black uppercase tracking-[0.2em] opacity-60">{stat.label}</span>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-black tracking-tighter text-on-surface italic">{stat.val}</span>
                <span className={`${stat.trendColor} text-[10px] font-black uppercase tracking-tighter italic`}>{stat.trend}</span>
              </div>
            </div>
          ))}
        </section>

        {/* The Saga Visualizer */}
        <section className="bg-surface-container-lowest rounded-[2rem] p-12 border border-outline-variant/15 shadow-xl relative overflow-hidden">
          <div className="flex justify-between items-start mb-16 relative z-10">
            <div>
              <h3 className="text-2xl font-black font-headline text-on-surface mb-2 tracking-tight uppercase italic">Order #RR-9942 Transaction Flow</h3>
              <p className="text-on-surface-variant max-w-xl text-sm font-medium leading-relaxed">
                Real-time trace of the distributed transaction. Visualizing the Choreography-based Saga with active compensation logic.
              </p>
            </div>
          </div>

          {/* Diagram Canvas */}
          <div className="relative min-h-[400px] w-full flex items-center justify-center">
            {/* SVG Connector Lines */}
            <svg className="absolute inset-0 w-full h-full pointer-events-none opacity-30" xmlns="http://www.w3.org/2000/svg">
              <path className="marker-path" d="M 220 180 L 380 180" fill="none" stroke="#004ac6" strokeWidth="2.5" strokeDasharray="5 3"></path>
              <path className="marker-path" d="M 580 180 L 740 180" fill="none" stroke="#004ac6" strokeWidth="2.5" strokeDasharray="5 3"></path>
              <path className="marker-path" d="M 840 220 Q 840 320 500 320" fill="none" stroke="#ba1a1a" strokeDasharray="6" strokeWidth="2" style={{ filter: 'url(#retail-rough-edge)' }}></path>
              <path className="marker-path" d="M 400 320 L 120 320 Q 80 320 80 220" fill="none" stroke="#ba1a1a" strokeDasharray="6" strokeWidth="2" style={{ filter: 'url(#retail-rough-edge)' }}></path>
            </svg>

            {/* Service Nodes */}
            <div className="flex justify-between w-full items-center gap-12 relative z-10 px-12">
              {/* Node 1 */}
              <div className="flex flex-col items-center gap-4">
                <div className="w-48 h-44 bg-surface-container-low rounded-3xl p-6 flex flex-col items-center justify-center text-center gap-4 relative border-2 border-tertiary-fixed shadow-lg hover:scale-105 transition-transform cursor-pointer">
                  <div className="w-14 h-14 bg-tertiary-fixed rounded-full flex items-center justify-center text-on-tertiary-fixed shadow-inner">
                    <span className="material-symbols-outlined !text-3xl">payments</span>
                  </div>
                  <div>
                    <h4 className="font-black text-xs text-on-surface uppercase tracking-tighter">Payment Service</h4>
                    <p className="text-[9px] text-tertiary font-mono font-black mt-1 uppercase italic tracking-widest">STATUS: COMMITTED</p>
                  </div>
                  <div className="absolute -top-4 -right-4 w-10 h-10 bg-tertiary-fixed rounded-full flex items-center justify-center text-on-tertiary-fixed border-4 border-surface shadow-lg">
                    <span className="material-symbols-outlined !text-lg font-black">check</span>
                  </div>
                </div>
              </div>
              {/* Node 2 */}
              <div className="flex flex-col items-center gap-4">
                <div className="w-48 h-44 bg-surface-container-low rounded-3xl p-6 flex flex-col items-center justify-center text-center gap-4 relative border-2 border-tertiary-fixed shadow-lg hover:scale-105 transition-transform cursor-pointer">
                  <div className="w-14 h-14 bg-tertiary-fixed rounded-full flex items-center justify-center text-on-tertiary-fixed shadow-inner">
                    <span className="material-symbols-outlined !text-3xl">inventory_2</span>
                  </div>
                  <div>
                    <h4 className="font-black text-xs text-on-surface uppercase tracking-tighter">Warehouse (WMS)</h4>
                    <p className="text-[9px] text-tertiary font-mono font-black mt-1 uppercase italic tracking-widest">ALLOCATED_OK</p>
                  </div>
                  <div className="absolute -top-4 -right-4 w-10 h-10 bg-tertiary-fixed rounded-full flex items-center justify-center text-on-tertiary-fixed border-4 border-surface shadow-lg">
                    <span className="material-symbols-outlined !text-lg font-black">check</span>
                  </div>
                </div>
              </div>
              {/* Node 3 */}
              <div className="flex flex-col items-center gap-4">
                <div className="w-48 h-44 bg-error-container rounded-3xl p-6 flex flex-col items-center justify-center text-center gap-4 relative border-2 border-error shadow-2xl hover:scale-105 transition-transform cursor-pointer">
                  <div className="w-14 h-14 bg-error rounded-full flex items-center justify-center text-on-error shadow-xl animate-pulse">
                    <span className="material-symbols-outlined !text-3xl">local_shipping</span>
                  </div>
                  <div>
                    <h4 className="font-black text-xs text-on-error-container uppercase tracking-tighter">Logistics API</h4>
                    <p className="text-[9px] text-error font-mono font-black mt-1 uppercase tracking-tighter italic">FAILURE: 404_ADDR</p>
                  </div>
                  <div className="absolute -top-4 -right-4 w-10 h-10 bg-error rounded-full flex items-center justify-center text-on-error border-4 border-surface shadow-lg">
                    <span className="material-symbols-outlined !text-lg font-black">close</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Marker Annotations */}
            <div className="absolute top-0 left-1/4 transform -translate-x-1/2">
              <div className="font-marker text-2xl text-primary rotate-[-4deg] max-w-[200px] drop-shadow-sm font-black">
                Phase 1: Event-driven payment verification
              </div>
            </div>
            <div className="absolute bottom-10 left-1/2 transform -translate-x-1/2">
              <div className="bg-error-container border-2 border-error/20 p-6 rounded-2xl shadow-2xl scale-110">
                <div className="font-marker text-3xl text-error font-black mb-1 tracking-tighter">Saga Rollback Triggered!</div>
                <p className="font-marker text-lg text-on-error-container font-medium">Reversing locally committed transactions...</p>
              </div>
            </div>
          </div>
          
          <div className="absolute inset-0 pointer-events-none opacity-[0.02]" style={{ backgroundImage: 'radial-gradient(#000 1px, transparent 1px)', backgroundSize: '32px 32px' }}></div>
        </section>

        {/* Bottom Bento Logs */}
        <section className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 bg-surface-container-lowest p-8 rounded-3xl border border-outline-variant/15 shadow-sm">
            <h3 className="font-black font-headline text-sm uppercase tracking-[0.2em] text-on-surface-variant opacity-60 mb-8 flex items-center gap-2 italic">
              <span className="material-symbols-outlined text-primary !text-xl">list_alt</span>
              Live Compensation Log
            </h3>
            <div className="space-y-4">
              {[
                { type: 'REFUND_TX_INITIATED', id: '8829-PAY', time: '14:22:01.04', status: 'tertiary' },
                { type: 'STOCK_RESTORE_COMPLETE', id: 'SKU: ETH-L-01', time: '14:21:58.22', status: 'tertiary' },
                { type: 'SHIPMENT_REJECTED', id: 'INVALID_POSTAL', time: '14:21:55.10', status: 'error' }
              ].map((log, i) => (
                <div key={i} className={`flex items-center justify-between p-5 rounded-2xl border border-outline-variant/10 ${log.status === 'error' ? 'bg-error-container/30 border-error/10' : 'bg-surface-container-low/50'}`}>
                  <div className="flex items-center gap-4">
                    <span className={`w-2.5 h-2.5 rounded-full ${log.status === 'error' ? 'bg-error animate-pulse' : 'bg-tertiary-fixed-dim shadow-[0_0_8px_rgba(78,222,163,0.5)]'}`}></span>
                    <span className={`text-[10px] font-mono font-black ${log.status === 'error' ? 'text-error' : 'text-on-surface'}`}>{log.type}</span>
                    <span className="text-[10px] text-on-surface-variant font-black uppercase tracking-widest opacity-40">{log.id}</span>
                  </div>
                  <span className="text-[10px] text-on-surface-variant font-mono tabular-nums opacity-60">{log.time}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-on-background text-white p-10 rounded-[2.5rem] relative overflow-hidden shadow-2xl border-t-8 border-primary">
            <div className="relative z-10 h-full flex flex-col">
              <h4 className="font-black font-headline text-xs uppercase tracking-[0.3em] text-slate-500 mb-10 italic">Orchestrator Config</h4>
              <ul className="space-y-8 flex-1">
                {[
                  { label: 'Timeout Policy', val: '5000ms' },
                  { label: 'Retry Strategy', val: 'Exp. Backoff' },
                  { label: 'Max Retries', val: '3' },
                  { label: 'Consistency', val: 'Eventual' }
                ].map((conf, i) => (
                  <li key={i} className="flex justify-between items-center border-b border-slate-800 pb-3 group">
                    <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest group-hover:text-slate-300 transition-colors">{conf.label}</span>
                    <span className="font-mono text-[10px] font-black text-primary-fixed-dim skew-x-[-10deg]">{conf.val}</span>
                  </li>
                ))}
              </ul>
              <button className="mt-12 w-full py-4 bg-white text-on-background font-black text-[10px] uppercase tracking-[0.2em] rounded-2xl hover:bg-primary-fixed hover:scale-95 transition-all shadow-xl flex items-center justify-center gap-3">
                <span className="material-symbols-outlined !text-lg">settings_suggest</span>
                Update Saga Logic
              </button>
            </div>
            <div className="absolute -bottom-20 -right-20 w-64 h-64 bg-primary/10 rounded-full blur-[100px]"></div>
          </div>
        </section>
      </div>
    </div>
  );
}
