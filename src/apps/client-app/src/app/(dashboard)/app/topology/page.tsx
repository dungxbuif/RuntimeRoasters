import React from 'react';

export default function SystemTopologyPage() {
  return (
    <div className="flex-1 p-8 md:p-12 lg:p-20 bg-surface relative overflow-hidden h-full min-h-screen font-headline">
      {/* Grid Background Pattern */}
      <div className="absolute inset-0 pointer-events-none opacity-[0.03]" style={{ backgroundImage: 'radial-gradient(#000 2px, transparent 2px)', backgroundSize: '32px 32px' }}></div>
      
      <div className="max-w-7xl mx-auto relative z-10">
        <div className="mb-12 flex justify-between items-end border-b-2 border-outline-variant/10 pb-8">
          <div>
            <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-primary-fixed/30 text-on-primary-fixed text-[10px] font-black uppercase tracking-[0.2em] mb-4 border border-primary/10 shadow-sm italic">
              <span className="material-symbols-outlined text-[14px]">account_tree</span>
              <span>L2 Network Visualization</span>
            </div>
            <h2 className="font-black text-5xl tracking-tighter text-on-surface uppercase italic">Network <span className="text-primary">Topology</span></h2>
            <p className="font-body text-on-surface-variant mt-4 font-medium tracking-tight max-w-xl leading-relaxed">Real-time supply chain data flow across distributed ingestion, processing, and persistence layers. v4.2.0-STABLE</p>
          </div>
          <div className="flex gap-8 mb-2">
            <div className="flex items-center gap-3">
              <span className="w-4 h-4 bg-tertiary border-2 border-surface inline-block shadow-lg"></span>
              <span className="text-[10px] font-black uppercase tracking-widest text-on-surface-variant italic">System Healthy</span>
            </div>
            <div className="flex items-center gap-3">
              <span className="w-4 h-4 bg-secondary-fixed border-2 border-surface inline-block shadow-lg"></span>
              <span className="text-[10px] font-black uppercase tracking-widest text-on-surface-variant italic">Nominal Load</span>
            </div>
          </div>
        </div>

        {/* Circuit Layout */}
        <div className="relative grid grid-cols-1 lg:grid-cols-3 gap-20 pt-16">
          
          {/* Layer 1: Ingestion (Left) */}
          <div className="flex flex-col gap-12 relative">
            <div className="absolute right-[-5rem] top-1/2 w-20 border-b-4 border-outline-variant/20 border-dashed hidden lg:block"></div>
            <h3 className="text-[10px] font-black uppercase tracking-[0.4em] text-on-surface-variant opacity-40 border-l-4 border-tertiary pl-4 mb-6 italic">L1: Ingestion Mesh</h3>
            
            {/* Node Card 1 */}
            <div className="bg-surface-container-lowest border-2 border-outline-variant/10 p-8 rounded-[2rem] shadow-[12px_12px_0px_0px_rgba(0,74,198,0.03)] relative group hover:border-primary/30 transition-all duration-500">
              <div className="absolute -right-3 top-1/2 w-6 h-6 bg-surface-container-high rounded-none transform translate-x-full -translate-y-1/2 z-20 border-2 border-outline-variant/10 group-hover:bg-primary group-hover:border-primary transition-colors"></div>
              <div className="flex justify-between items-start mb-8 italic">
                <div>
                  <h4 className="font-black text-2xl mb-1 text-on-surface tracking-tight uppercase">Telemetry API</h4>
                  <span className="font-mono text-[9px] bg-surface-container-high text-primary px-3 py-1 border border-outline-variant/20 font-black tracking-tighter">REST / GRPC_V2</span>
                </div>
                <span className="material-symbols-outlined text-primary !text-4xl drop-shadow-sm">sensors</span>
              </div>
              <div className="space-y-6">
                <div>
                  <div className="flex justify-between text-[10px] font-black mb-3 text-on-surface-variant uppercase tracking-widest">
                    <span>Throughput</span>
                    <span className="text-primary">4,200 req/s</span>
                  </div>
                  <div className="h-3 rounded-full w-full bg-surface-container-high overflow-hidden p-0.5 border border-outline-variant/5">
                    <div className="h-full bg-primary rounded-full w-[75%] shadow-[0_0_12px_rgba(0,74,198,0.4)] transition-all duration-1000"></div>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4 text-[10px] font-black">
                  <div className="bg-surface p-3 border border-outline-variant/10 flex justify-between rounded-xl italic">
                    <span className="text-on-surface-variant opacity-50 uppercase tracking-tighter">Lat</span>
                    <span className="text-on-surface tabular-nums">12ms</span>
                  </div>
                  <div className="bg-surface p-3 border border-outline-variant/10 flex justify-between rounded-xl italic">
                    <span className="text-on-surface-variant opacity-50 uppercase tracking-tighter">Err</span>
                    <span className="text-tertiary tabular-nums">0.01%</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Node Card 2 */}
            <div className="bg-surface-container-lowest border-2 border-outline-variant/10 p-8 rounded-[2rem] shadow-[12px_12px_0px_0px_rgba(0,0,0,0.03)] relative group hover:border-tertiary-fixed transition-all duration-500">
              <div className="absolute -right-3 top-1/2 w-6 h-6 bg-surface-container-high rounded-none transform translate-x-full -translate-y-1/2 z-20 border-2 border-outline-variant/10 group-hover:bg-tertiary transition-colors"></div>
              <div className="flex justify-between items-start mb-8 italic">
                <div>
                  <h4 className="font-black text-2xl mb-1 text-on-surface tracking-tight uppercase">IoT Gateway</h4>
                  <span className="font-mono text-[9px] bg-surface-container-high text-tertiary px-3 py-1 border border-outline-variant/20 font-black tracking-tighter">MQTT_OVER_TLS</span>
                </div>
                <span className="material-symbols-outlined text-tertiary !text-4xl">router</span>
              </div>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-3 text-[10px] font-black">
                  <div className="bg-tertiary/5 border border-tertiary/20 p-4 text-tertiary font-black border-l-4 border-l-tertiary text-center uppercase tracking-widest shadow-sm">
                    ACTIVE: 142
                  </div>
                  <div className="bg-surface p-4 border border-outline-variant/10 text-center text-on-surface-variant opacity-40 font-black uppercase tracking-widest italic">
                    STALE: 0
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Layer 2: Processing (Center) */}
          <div className="flex flex-col justify-center relative">
            <h3 className="text-[10px] font-black uppercase tracking-[0.4em] text-on-surface-variant opacity-40 border-b-4 border-secondary-fixed pb-2 mb-12 text-center italic">L2: Event Infrastructure</h3>
            <div className="bg-surface-container-lowest border-[6px] border-surface-container-high p-12 rounded-[3.5rem] shadow-[20px_20px_60px_rgba(0,0,0,0.05)] relative transform lg:scale-110 group hover:border-primary/20 transition-all duration-700">
              {/* Connection Hubs */}
              <div className="absolute -left-5 top-1/2 w-5 h-16 bg-surface-container-highest group-hover:bg-primary transform -translate-y-1/2 transition-colors rounded-l-xl"></div>
              <div className="absolute -right-5 top-1/2 w-5 h-16 bg-surface-container-highest group-hover:bg-primary transform -translate-y-1/2 transition-colors rounded-r-xl"></div>
              
              <div className="text-center mb-10">
                <span className="material-symbols-outlined text-primary animate-pulse mb-6 !text-7xl drop-shadow-xl font-black">hub</span>
                <h4 className="font-black text-4xl tracking-tighter text-on-surface uppercase italic">Kafka Cluster</h4>
                <span className="font-mono text-xs bg-on-background text-white px-5 py-2 inline-block mt-6 font-black skew-x-[-12deg] shadow-lg italic">CORE_BROKER_PROD_01</span>
              </div>
              
              <div className="space-y-4 mt-12">
                {[
                  { label: 'Live Topics', val: '24', color: 'text-on-surface' },
                  { label: 'Replication', val: '3x_SYNCED', color: 'text-tertiary' },
                  { label: 'Bus State', val: 'NOMINAL', color: 'text-tertiary' }
                ].map((item, i) => (
                  <div key={i} className="border border-outline-variant/10 p-5 flex justify-between items-center bg-surface-container-low/50 rounded-2xl transition-all hover:bg-white hover:shadow-md">
                    <span className="font-headline text-[10px] font-black text-on-surface-variant uppercase tracking-widest italic">{item.label}</span>
                    <span className={`font-mono text-xs font-black px-3 py-1 rounded-full bg-white shadow-sm border border-outline-variant/5 ${item.color}`}>{item.val}</span>
                  </div>
                ))}
              </div>
              
              {/* Hand-drawn annotation */}
              <div className="absolute -bottom-20 left-1/2 -translate-x-1/2 w-full text-center">
                 <p className="font-marker text-3xl text-primary rotate-[-2deg] opacity-60">Partitioned by Batch_ID</p>
              </div>
            </div>
          </div>

          {/* Layer 3: Persistence (Right) */}
          <div className="flex flex-col gap-12 relative">
            <div className="absolute left-[-5rem] top-1/2 w-20 border-b-4 border-outline-variant/20 border-dashed hidden lg:block"></div>
            <h3 className="text-[10px] font-black uppercase tracking-[0.4em] text-on-surface-variant opacity-40 border-r-4 border-primary-container pr-4 mb-6 text-right italic">L3: State Persistence</h3>
            
            {/* Node Card 3 */}
            <div className="bg-surface-container-lowest border-2 border-outline-variant/10 p-8 rounded-[2rem] shadow-[12px_12px_0px_0px_rgba(0,0,0,0.03)] relative group hover:border-primary transition-all duration-500">
              <div className="absolute -left-3 top-1/2 w-6 h-6 bg-surface-container-high rounded-none transform -translate-x-full -translate-y-1/2 z-20 border-2 border-outline-variant/10 group-hover:bg-primary transition-colors"></div>
              <div className="flex justify-between items-start mb-8 italic">
                <div>
                  <h4 className="font-black text-2xl mb-1 text-on-surface tracking-tight uppercase">Distributed DB</h4>
                  <span className="font-mono text-[9px] bg-surface-container-high text-primary px-3 py-1 border border-outline-variant/20 font-black tracking-tighter italic">CASSANDRA_V4</span>
                </div>
                <span className="material-symbols-outlined text-primary !text-4xl drop-shadow-sm">database</span>
              </div>
              <div className="space-y-6">
                <div>
                  <div className="flex justify-between text-[10px] font-black mb-3 text-on-surface-variant uppercase tracking-widest italic">
                    <span>Data Utilization</span>
                    <span className="text-primary font-black tracking-tighter text-xs">8.2 TB / 10 TB</span>
                  </div>
                  <div className="h-3 rounded-full w-full bg-surface-container-high overflow-hidden p-0.5 border border-outline-variant/5 shadow-inner">
                    <div className="h-full bg-primary rounded-full w-[82%] shadow-[0_0_12px_rgba(0,74,198,0.4)]"></div>
                  </div>
                </div>
                <div className="text-[11px] font-black border border-outline-variant/10 p-4 text-center bg-surface rounded-2xl text-on-surface-variant uppercase tracking-[0.2em] shadow-sm">
                  Real-time Writes: <span className="text-on-surface ml-3 tabular-nums italic">2.4k pts/s</span>
                </div>
              </div>
            </div>

            {/* Node Card 4 */}
            <div className="bg-surface-container-lowest border-2 border-outline-variant/10 p-8 rounded-[2rem] shadow-[12px_12px_0px_0px_rgba(0,0,0,0.03)] relative group hover:border-secondary transition-all duration-500">
              <div className="absolute -left-3 top-1/2 w-6 h-6 bg-surface-container-high rounded-none transform -translate-x-full -translate-y-1/2 z-20 border-2 border-outline-variant/10 group-hover:bg-secondary transition-colors"></div>
              <div className="flex justify-between items-start mb-8 italic">
                <div>
                  <h4 className="font-black text-2xl mb-1 text-on-surface tracking-tight uppercase">Machine Learning</h4>
                  <span className="font-mono text-[9px] bg-surface-container-high text-secondary px-3 py-1 border border-outline-variant/20 font-black tracking-tighter">PREDICTIVE_CORE_V1</span>
                </div>
                <span className="material-symbols-outlined text-secondary !text-4xl">memory</span>
              </div>
              <div className="border border-outline-variant/10 p-5 bg-surface flex items-center justify-between rounded-[1.5rem] shadow-sm italic transition-all group-hover:bg-white group-hover:shadow-lg">
                <span className="font-headline text-[11px] font-black text-on-surface-variant uppercase tracking-widest opacity-60">Model Drift</span>
                <span className="font-mono text-xs font-black text-secondary tabular-nums tracking-widest underline decoration-2 decoration-secondary-fixed">0.02% [STABLE]</span>
              </div>
            </div>
          </div>

        </div>
        
        {/* Footnote with rough stroke */}
        <div className="mt-32 pt-12 border-t-4 border-outline-variant/10 border-dashed flex justify-center items-center">
           <div className="bg-white px-10 py-6 rounded-[2rem] shadow-2xl border-2 border-outline-variant/10 transform hover:scale-105 transition-transform flex items-center gap-8">
              <div className="w-12 h-12 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-600 border border-emerald-200 shadow-inner italic font-black">!</div>
              <p className="font-marker text-3xl text-on-surface-variant rotate-[-1deg]">All nodes are currently synchronized with the Global Cluster Clock.</p>
           </div>
        </div>
      </div>
    </div>
  );
}
