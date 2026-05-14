import React from 'react';

export default function TraceabilityJourneyPage() {
  return (
    <div className="min-h-full bg-surface text-on-surface p-8 md:p-12 lg:p-20 font-body selection:bg-tertiary-fixed selection:text-on-tertiary-fixed relative overflow-x-hidden">
      {/* Page Header */}
      <div className="max-w-4xl mx-auto mb-20 relative z-10">
        <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-secondary-container text-on-secondary-container text-sm font-medium mb-4 shadow-sm border border-outline-variant/10">
          <span className="material-symbols-outlined text-[18px]">history</span>
          <span>1-N Traceability Logic</span>
        </div>
        <h1 className="font-headline text-5xl md:text-6xl font-bold tracking-tight text-on-background mb-4 uppercase italic">
          Batch <span className="text-primary">Time Machine</span>
        </h1>
        <p className="font-body text-lg text-on-surface-variant max-w-2xl leading-relaxed font-bold uppercase tracking-tighter italic">
          Visualizing the <span className="text-primary underline decoration-2">Aggregation Pattern</span>: How 1 Production Batch maps to multiple Roast Runs.
        </p>
      </div>

      {/* Timeline Diagram Container */}
      <div className="max-w-5xl mx-auto relative pb-32">
        {/* Central Spine */}
        <div className="absolute left-1/2 top-0 bottom-0 w-2 -translate-x-1/2 z-0 pointer-events-none hidden md:block opacity-20 bg-primary/20 rounded-full"></div>

        <div className="space-y-32 relative z-10">
          {/* Step 1: Farm Origin */}
          <div className="flex flex-col md:flex-row items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-end md:pr-12 relative">
              <div className="bg-surface-container-lowest p-8 rounded-[2rem] shadow-xl border border-outline-variant/15 w-full max-w-sm hover:border-primary transition-all duration-500 relative overflow-hidden">
                <div className="flex items-center space-x-4 mb-6">
                  <div className="w-14 h-14 rounded-2xl bg-primary/10 text-primary flex items-center justify-center border-2 border-primary/20">
                    <span className="material-symbols-outlined text-3xl">agriculture</span>
                  </div>
                  <div>
                    <div className="text-[10px] font-black text-outline uppercase tracking-widest opacity-50">Origin Node</div>
                    <h3 className="font-headline text-xl font-black text-on-surface uppercase italic">Harvest #H-2026-001</h3>
                  </div>
                </div>
                <div className="space-y-3 font-body text-xs text-on-surface-variant font-bold uppercase tracking-widest italic opacity-80">
                  <div className="flex justify-between border-b border-outline-variant/10 pb-2">
                    <span>Farm</span>
                    <span className="text-on-surface">Cau Dat Highlands</span>
                  </div>
                  <div className="flex justify-between pt-1">
                    <span>Initial Weight</span>
                    <span className="text-primary">500.00 kg</span>
                  </div>
                </div>
              </div>
            </div>
            <div className="w-10 h-10 rounded-full bg-surface border-[6px] border-primary z-10 hidden md:block shadow-2xl animate-pulse"></div>
            <div className="w-full md:w-5/12 md:pl-12 mt-6 md:mt-0">
               <div className="font-marker text-3xl text-primary -rotate-2 leading-tight">
                  <span className="bg-primary/10 px-2 py-1 inline-block mb-2">Event Published:</span><br/>
                  Kafka broadcasted harvest to Warehouse Service.
               </div>
            </div>
          </div>

          {/* Step 2: Aggregated Processing (The 1-N Focus) */}
          <div className="flex flex-col md:flex-row-reverse items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-start md:pl-12 relative">
              <div className="bg-on-background text-white p-8 rounded-[2.5rem] shadow-2xl w-full max-w-md relative overflow-hidden border border-white/10">
                <div className="flex items-center space-x-4 mb-8">
                  <div className="w-14 h-14 rounded-2xl bg-white/10 text-white flex items-center justify-center border border-white/20 animate-spin-slow">
                    <span className="material-symbols-outlined text-3xl">sync</span>
                  </div>
                  <div>
                    <div className="text-[10px] font-black opacity-50 uppercase tracking-widest">Processing Node</div>
                    <h3 className="font-headline text-xl font-black uppercase italic">Batch #RR-P-CD-001</h3>
                  </div>
                </div>
                
                <div className="grid grid-cols-2 gap-4 mb-8">
                   {[1, 2, 3, 4].map(run => (
                      <div key={run} className="bg-white/5 border border-white/10 p-4 rounded-xl flex items-center justify-between hover:bg-white/10 transition-all cursor-help">
                         <span className="text-[10px] font-black opacity-40">RUN #{run}</span>
                         <span className="font-headline text-sm font-black tabular-nums text-emerald-400">~12kg</span>
                      </div>
                   ))}
                </div>

                <div className="p-4 bg-white/5 rounded-2xl border border-white/10 flex justify-between items-center italic">
                   <span className="text-[10px] font-black uppercase opacity-60">Aggr. Output</span>
                   <span className="text-2xl font-headline font-black text-emerald-400 tabular-nums">45.0 kg</span>
                </div>
              </div>
            </div>
            <div className="w-10 h-10 rounded-full bg-surface border-[6px] border-emerald-500 z-10 hidden md:block shadow-2xl"></div>
            <div className="w-full md:w-5/12 md:pr-12 mt-6 md:mt-0 flex justify-end">
               <div className="font-marker text-3xl text-emerald-600 rotate-1 text-right leading-tight">
                  <span className="bg-emerald-100 px-2 py-1 inline-block mb-2">1-N Aggregation:</span><br/>
                  Combining multiple small roast runs into one commercial batch.
               </div>
            </div>
          </div>

          {/* Step 3: Stock-In */}
          <div className="flex flex-col md:flex-row items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-end md:pr-12 relative">
              <div className="bg-surface-container-lowest p-8 rounded-[2rem] shadow-xl border border-outline-variant/15 w-full max-w-sm hover:border-tertiary transition-all duration-500 overflow-hidden relative">
                <div className="absolute top-0 right-0 w-24 h-24 bg-tertiary/5 rounded-full -mr-10 -mt-10 blur-2xl"></div>
                <div className="flex items-center space-x-4 mb-6">
                  <div className="w-14 h-14 rounded-2xl bg-tertiary/10 text-tertiary flex items-center justify-center border-2 border-tertiary/20">
                    <span className="material-symbols-outlined text-3xl">inventory_2</span>
                  </div>
                  <div>
                    <div className="text-[10px] font-black text-outline uppercase tracking-widest opacity-50">Stock Node</div>
                    <h3 className="font-headline text-xl font-black text-on-surface uppercase italic">PROD-CD-001</h3>
                  </div>
                </div>
                <div className="space-y-3 font-body text-xs text-on-surface-variant font-bold uppercase tracking-widest italic opacity-80">
                  <div className="flex justify-between border-b border-outline-variant/10 pb-2">
                    <span>SKU</span>
                    <span className="text-on-surface font-black">CD-ARABICA-ROASTED</span>
                  </div>
                  <div className="flex justify-between pt-1">
                    <span>Final Qty</span>
                    <span className="text-tertiary font-black">45.0 kg</span>
                  </div>
                </div>
              </div>
            </div>
            <div className="w-10 h-10 rounded-full bg-surface border-[6px] border-tertiary z-10 hidden md:block shadow-2xl"></div>
            <div className="w-full md:w-5/12 md:pl-12 mt-6 md:mt-0">
               <div className="font-marker text-3xl text-slate-700 -rotate-1 leading-tight">
                  <span className="bg-tertiary/10 px-2 py-1 inline-block mb-2">Inventory Finalized:</span><br/>
                  Batch locked & SKU quantity incremented in Global Stock.
               </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
