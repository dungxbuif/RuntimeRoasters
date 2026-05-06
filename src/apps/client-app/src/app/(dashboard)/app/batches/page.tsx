import React from 'react';

export default function TraceabilityJourneyPage() {
  return (
    <div className="min-h-full bg-surface text-on-surface p-8 md:p-12 lg:p-20 font-body selection:bg-tertiary-fixed selection:text-on-tertiary-fixed relative overflow-x-hidden">
      {/* Page Header */}
      <div className="max-w-4xl mx-auto mb-20 relative z-10">
        <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-secondary-container text-on-secondary-container text-sm font-medium mb-4 shadow-sm border border-outline-variant/10">
          <span className="material-symbols-outlined text-[18px]">history</span>
          <span>Event Sourcing Demo</span>
        </div>
        <h1 className="font-headline text-5xl md:text-6xl font-bold tracking-tight text-on-background mb-4 uppercase italic">
          Traceability <span className="text-primary">Time Machine</span>
        </h1>
        <p className="font-body text-lg text-on-surface-variant max-w-2xl leading-relaxed">
          Following the complete lifecycle of Coffee Batch <span className="font-mono bg-surface-container-highest px-2 py-0.5 rounded text-primary text-sm italic">#CF-882A</span> from soil to sip. Observe how architectural patterns manage state across distributed nodes.
        </p>
      </div>

      {/* Timeline Diagram Container */}
      <div className="max-w-5xl mx-auto relative pb-32">
        {/* Central Sketchy Timeline Spine */}
        <div className="absolute left-1/2 top-0 bottom-0 w-2 -translate-x-1/2 z-0 pointer-events-none hidden md:block opacity-40">
          <svg className="overflow-visible" height="100%" width="40">
            <path className="text-outline" d="M 20 0 Q 30 150 15 300 T 25 600 T 15 900 T 22 1200 T 18 1500" fill="none" stroke="currentColor" strokeDasharray="8 8" strokeWidth="3"></path>
          </svg>
        </div>

        <div className="space-y-32 relative z-10">
          {/* Step 1: Farm */}
          <div className="flex flex-col md:flex-row items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-end md:pr-12 relative">
              <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_12px_32px_-4px_rgba(25,28,30,0.06)] ring-1 ring-outline-variant/15 w-full max-w-sm hover:ring-tertiary-fixed transition-all duration-300 relative overflow-hidden group/card">
                <div className="absolute -top-10 -right-10 w-32 h-32 bg-primary-fixed rounded-full blur-2xl opacity-40 mix-blend-multiply"></div>
                <div className="flex items-center space-x-4 mb-4 relative z-10">
                  <div className="w-12 h-12 rounded-lg bg-primary/10 text-primary flex items-center justify-center border border-primary/20">
                    <span className="material-symbols-outlined text-2xl">eco</span>
                  </div>
                  <div>
                    <div className="text-xs font-bold text-outline uppercase tracking-wider">Node 01</div>
                    <h3 className="font-headline text-xl font-semibold text-on-surface uppercase">Farm & Harvest</h3>
                  </div>
                </div>
                <div className="space-y-2 font-body text-sm text-on-surface-variant relative z-10 font-bold uppercase tracking-tighter">
                  <div className="flex justify-between border-b border-outline-variant/20 pb-1">
                    <span className="text-outline">Origin</span>
                    <span className="font-medium text-on-surface">Finca El Paraiso, CO</span>
                  </div>
                  <div className="flex justify-between border-b border-outline-variant/20 pb-1">
                    <span className="text-outline">Elevation</span>
                    <span className="font-medium text-on-surface italic">1,750m</span>
                  </div>
                  <div className="flex justify-between pt-1">
                    <span className="text-outline">Yield</span>
                    <span className="font-medium text-on-surface italic text-emerald-600">450kg Cherry</span>
                  </div>
                </div>
              </div>
              <div className="hidden md:block absolute right-0 top-1/2 w-12 border-t-2 border-dashed border-outline/40 -translate-y-1/2"></div>
            </div>
            <div className="w-6 h-6 rounded-full bg-surface border-4 border-primary z-10 hidden md:block shadow-lg"></div>
            <div className="w-full md:w-5/12 md:pl-12 mt-6 md:mt-0 flex items-center">
              <div className="relative">
                <p className="font-marker text-3xl text-blue-600 -rotate-2 leading-relaxed max-w-xs drop-shadow-sm">
                  <span className="bg-primary-fixed/50 px-1 inline-block transform -skew-x-6 mb-1">Event Sourced:</span> <br/>
                  Initial &apos;BeanHarvested&apos; event appended to immutable ledger.
                </p>
              </div>
            </div>
          </div>

          {/* Step 2: Process */}
          <div className="flex flex-col md:flex-row-reverse items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-start md:pl-12 relative">
              <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_12px_32px_-4px_rgba(25,28,30,0.06)] ring-1 ring-outline-variant/15 w-full max-w-sm hover:ring-tertiary-fixed transition-all duration-300 relative overflow-hidden group/card">
                <div className="absolute -bottom-8 -left-8 w-24 h-24 bg-tertiary-container rounded-full blur-xl opacity-20 mix-blend-multiply"></div>
                <div className="flex items-center space-x-4 mb-4 relative z-10">
                  <div className="w-12 h-12 rounded-lg bg-tertiary/10 text-tertiary flex items-center justify-center border border-tertiary/20">
                    <span className="material-symbols-outlined text-2xl">water_drop</span>
                  </div>
                  <div>
                    <div className="text-xs font-bold text-outline uppercase tracking-wider">Node 02</div>
                    <h3 className="font-headline text-xl font-semibold text-on-surface uppercase">Wet Processing</h3>
                  </div>
                </div>
                <div className="space-y-2 font-body text-sm text-on-surface-variant relative z-10 font-bold uppercase tracking-tighter">
                  <div className="flex justify-between border-b border-outline-variant/20 pb-1">
                    <span className="text-outline">Method</span>
                    <span className="font-medium text-on-surface">Fully Washed</span>
                  </div>
                  <div className="flex justify-between border-b border-outline-variant/20 pb-1">
                    <span className="text-outline">Fermentation</span>
                    <span className="font-medium text-on-surface italic">24 Hours</span>
                  </div>
                  <div className="flex justify-between pt-1">
                    <span className="text-outline">Moisture</span>
                    <span className="font-medium text-tertiary italic">11.2% (Target hit)</span>
                  </div>
                </div>
              </div>
              <div className="hidden md:block absolute left-0 top-1/2 w-12 border-t-2 border-dashed border-outline/40 -translate-y-1/2"></div>
            </div>
            <div className="w-6 h-6 rounded-full bg-surface border-4 border-tertiary z-10 hidden md:block shadow-lg"></div>
            <div className="w-full md:w-5/12 md:pr-12 mt-6 md:mt-0 flex justify-end items-center text-right">
              <div className="relative">
                <p className="font-marker text-3xl text-emerald-600 rotate-1 leading-relaxed max-w-xs ml-auto drop-shadow-sm">
                  State Machine Transition:<br/>
                  <code className="font-mono text-xl bg-surface-variant px-2 rounded-lg font-black mt-2 inline-block">Raw &rarr; Processed</code>
                </p>
              </div>
            </div>
          </div>

          {/* Step 3: Warehouse */}
          <div className="flex flex-col md:flex-row items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-end md:pr-12 relative">
              <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_12px_32px_-4px_rgba(25,28,30,0.06)] ring-1 ring-outline-variant/15 w-full max-w-sm hover:ring-tertiary-fixed transition-all duration-300 relative overflow-hidden group/card">
                <div className="flex items-center space-x-4 mb-4 relative z-10">
                  <div className="w-12 h-12 rounded-lg bg-secondary/10 text-secondary flex items-center justify-center border border-secondary/20">
                    <span className="material-symbols-outlined text-2xl">warehouse</span>
                  </div>
                  <div>
                    <div className="text-xs font-bold text-outline uppercase tracking-wider">Node 03</div>
                    <h3 className="font-headline text-xl font-semibold text-on-surface uppercase">Dry Milling</h3>
                  </div>
                </div>
                <div className="space-y-2 font-body text-sm text-on-surface-variant relative z-10 font-bold uppercase tracking-tighter">
                  <div className="flex justify-between border-b border-outline-variant/20 pb-1">
                    <span className="text-outline">Screen Size</span>
                    <span className="font-medium text-on-surface italic">15+ (Excelso)</span>
                  </div>
                  <div className="flex justify-between pt-1">
                    <span className="text-outline">Defects</span>
                    <span className="font-medium text-on-surface italic">0 Primary / 3 Sec.</span>
                  </div>
                </div>
              </div>
              <div className="hidden md:block absolute right-0 top-1/2 w-12 border-t-2 border-dashed border-outline/40 -translate-y-1/2"></div>
            </div>
            <div className="w-6 h-6 rounded-full bg-surface border-4 border-secondary z-10 hidden md:block shadow-lg"></div>
            <div className="w-full md:w-5/12 md:pl-12 mt-6 md:mt-0 flex items-center">
              <div className="relative">
                <p className="font-marker text-3xl text-slate-700 -rotate-1 leading-relaxed max-w-xs drop-shadow-sm">
                  CQRS Pattern: <br/>
                  Inventory <span className="border-b-4 border-tertiary-fixed-dim italic px-1">read model</span> updated for global visibility.
                </p>
              </div>
            </div>
          </div>

          {/* Step 4: Logistics */}
          <div className="flex flex-col md:flex-row-reverse items-center justify-between group">
            <div className="w-full md:w-5/12 flex justify-start md:pl-12 relative">
              <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_12px_32px_-4px_rgba(25,28,30,0.06)] ring-1 ring-outline-variant/15 w-full max-w-sm hover:ring-tertiary-fixed transition-all duration-300 relative overflow-hidden group/card border-l-4 border-error">
                <div className="flex items-center space-x-4 mb-4 relative z-10">
                  <div className="w-12 h-12 rounded-lg bg-surface-variant text-on-surface-variant flex items-center justify-center border border-outline-variant/40">
                    <span className="material-symbols-outlined text-2xl">directions_boat</span>
                  </div>
                  <div>
                    <div className="text-xs font-bold text-outline uppercase tracking-wider">Node 04</div>
                    <h3 className="font-headline text-xl font-semibold text-on-surface uppercase tracking-tighter">Global Transit</h3>
                  </div>
                </div>
                <div className="space-y-2 font-body text-sm text-on-surface-variant relative z-10 font-black uppercase italic tracking-widest text-red-800">
                  <div className="flex justify-between border-b border-outline-variant/20 pb-1">
                    <span className="opacity-50">Port</span>
                    <span>Buenaventura</span>
                  </div>
                  <div className="flex justify-between pt-1">
                    <span className="opacity-50">Status</span>
                    <span className="text-error flex items-center space-x-1">
                      <span className="material-symbols-outlined text-[14px]">warning</span>
                      <span>Delayed</span>
                    </span>
                  </div>
                </div>
              </div>
              <div className="hidden md:block absolute left-0 top-1/2 w-12 border-t-2 border-dashed border-outline/40 -translate-y-1/2"></div>
            </div>
            <div className="w-6 h-6 rounded-full bg-error border-4 border-error-container z-10 hidden md:block shadow-xl shadow-error-container/50"></div>
            <div className="w-full md:w-5/12 md:pr-12 mt-6 md:mt-0 flex justify-end items-center text-right relative">
              <div className="relative z-10">
                <p className="font-marker text-3xl text-on-surface-variant rotate-2 leading-tight max-w-sm drop-shadow-sm font-black">
                  <strong className="text-on-background text-blue-700">Saga Pattern:</strong> Transport booked.<br/>
                  <span className="text-error border-b-[3px] border-error-container pb-0.5 px-1 bg-error/5">Saga Rollback triggered here</span> due to ship delay! Compensating transactions fired to refund client.
                </p>
              </div>
            </div>
          </div>

          {/* Step 5: Retail */}
          <div className="flex flex-col md:flex-row items-center justify-between group opacity-50 grayscale hover:grayscale-0 hover:opacity-100 transition-all duration-500">
            <div className="w-full md:w-5/12 flex justify-end md:pr-12 relative">
              <div className="bg-surface-container-lowest p-6 rounded-xl shadow-sm ring-1 ring-outline-variant/15 w-full max-w-sm relative overflow-hidden">
                <div className="flex items-center space-x-4 mb-4 relative z-10">
                  <div className="w-12 h-12 rounded-lg bg-surface-variant text-outline flex items-center justify-center border border-outline-variant/40">
                    <span className="material-symbols-outlined text-2xl">local_cafe</span>
                  </div>
                  <div>
                    <div className="text-xs font-bold text-outline uppercase tracking-wider">Node 05</div>
                    <h3 className="font-headline text-xl font-semibold text-on-surface uppercase">Roastery & Retail</h3>
                  </div>
                </div>
                <div className="space-y-2 font-body text-sm text-on-surface-variant relative z-10 text-center italic py-4">
                  Awaiting arrival...
                </div>
              </div>
              <div className="hidden md:block absolute right-0 top-1/2 w-12 border-t-2 border-dashed border-outline/20 -translate-y-1/2"></div>
            </div>
            <div className="w-6 h-6 rounded-full bg-surface border-4 border-outline-variant/50 z-10 hidden md:block shadow-sm"></div>
            <div className="w-full md:w-5/12 md:pl-12 mt-6 md:mt-0 flex items-center">
              <div className="relative">
                <p className="font-marker text-3xl text-outline -rotate-1 leading-relaxed max-w-xs drop-shadow-sm italic">
                  Distributed Tracing:<br/>
                  Waiting for end-to-end span completion...
                </p>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}
