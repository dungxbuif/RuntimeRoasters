"use client";

import React, { useState } from 'react';

export default function WarehouseStockLogicPage() {
  const [intakeWeight, setIntakeWeight] = useState<number>(0);
  const [originalWeight] = useState<number>(500); // Mock original from Farm
  
  const deviation = Math.abs((originalWeight - intakeWeight) / originalWeight * 100);
  const showWarning = deviation > 2 && intakeWeight > 0;

  return (
    <div className="min-h-full p-8 font-body selection:bg-primary-fixed selection:text-on-primary-fixed relative overflow-hidden bg-surface">
      {/* Header Section */}
      <header className="mb-10 flex flex-col md:flex-row justify-between items-start md:items-end gap-6 relative z-10">
        <div>
          <div className="flex items-center gap-3 text-primary font-marker text-2xl mb-2 italic">
            <span>Runtime Roasters</span>
            <span className="material-symbols-outlined !text-sm">trending_flat</span>
            <span className="text-on-surface-variant opacity-60">Warehouse Service</span>
          </div>
          <h1 className="text-5xl font-headline font-black tracking-tighter text-on-surface uppercase italic">Warehouse Operations</h1>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12 relative z-10">
        {/* Left: Intake Command Center */}
        <div className="lg:col-span-2 space-y-8">
          <div className="bg-surface-container-lowest p-8 rounded-[2rem] border-2 border-primary/10 shadow-xl relative overflow-hidden group">
            <div className="absolute top-0 left-0 w-full h-2 bg-primary"></div>
            <div className="flex justify-between items-center mb-8">
              <h3 className="font-headline font-black text-2xl uppercase tracking-tight text-on-surface">1. Intake Verification</h3>
              <span className="bg-primary/10 text-primary px-4 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-primary/20">
                Awaiting: Harvest #H-2026-001
              </span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
              <div className="space-y-6">
                <div className="bg-surface-container-low p-6 rounded-2xl border border-outline-variant/30">
                  <span className="text-[10px] font-black text-on-surface-variant uppercase block mb-2 tracking-widest opacity-60">Farm Declared Weight</span>
                  <div className="text-4xl font-headline font-black text-on-surface tabular-nums">{originalWeight} kg</div>
                  <p className="text-xs text-on-surface-variant mt-2 font-bold uppercase italic">Origin: Cau Dat (Arabica)</p>
                </div>

                <div className="space-y-4">
                  <label className="text-[10px] font-black text-primary uppercase block tracking-widest">Actual Intake Weight (kg)</label>
                  <input 
                    type="number" 
                    value={intakeWeight || ''} 
                    onChange={(e) => setIntakeWeight(Number(e.target.value))}
                    placeholder="Enter scale value..."
                    className="w-full bg-surface-container-highest border-2 border-outline-variant rounded-2xl px-6 py-4 text-2xl font-headline font-black focus:border-primary transition-all outline-none"
                  />
                </div>
              </div>

              <div className="flex flex-col justify-center items-center p-8 bg-surface-container-low rounded-3xl border-2 border-dashed border-outline-variant">
                {intakeWeight > 0 ? (
                  <div className="text-center">
                    <div className={`text-6xl font-headline font-black mb-2 ${showWarning ? 'text-error animate-bounce' : 'text-tertiary'}`}>
                      {deviation.toFixed(2)}%
                    </div>
                    <span className="text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Weight Deviation</span>
                    
                    {showWarning && (
                      <div className="mt-6 p-4 bg-error-container text-on-error-container rounded-xl border border-error/20 flex flex-col gap-2">
                        <div className="flex items-center gap-2 font-black text-[10px] uppercase">
                          <span className="material-symbols-outlined !text-sm">warning</span>
                          Deviation Threshold Exceeded
                        </div>
                        <textarea 
                          placeholder="Please enter anomaly reason (Required)..."
                          className="w-full bg-white/50 rounded-lg p-3 text-xs font-bold border border-error/20 focus:outline-none"
                        />
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center opacity-30">
                    <span className="material-symbols-outlined !text-6xl mb-4">balance</span>
                    <p className="text-[10px] font-black uppercase tracking-widest">Awaiting Scale Input</p>
                  </div>
                )}
              </div>
            </div>

            <div className="mt-10 flex justify-end">
              <button 
                disabled={intakeWeight <= 0}
                className={`px-10 py-4 rounded-2xl font-black text-xs uppercase tracking-widest shadow-2xl transition-all flex items-center gap-3 ${
                  intakeWeight > 0 ? 'bg-primary text-white hover:scale-105' : 'bg-outline-variant text-on-surface-variant opacity-50'
                }`}
              >
                Confirm & Create Batch
                <span className="material-symbols-outlined !text-sm">send</span>
              </button>
            </div>
          </div>

          {/* Active Production Batches (1-N Visualization) */}
          <div className="bg-on-background text-white p-8 rounded-[2rem] shadow-2xl relative overflow-hidden">
             <div className="flex justify-between items-center mb-8">
                <h3 className="font-headline font-black text-2xl uppercase tracking-tighter italic">2. Active Processing (1-N)</h3>
                <div className="flex gap-2">
                   <div className="w-3 h-3 bg-emerald-400 rounded-full animate-pulse"></div>
                   <span className="text-[10px] font-black uppercase opacity-60">2 Batches Active</span>
                </div>
             </div>

             <div className="space-y-6">
                <div className="bg-white/5 border border-white/10 p-6 rounded-2xl group hover:bg-white/10 transition-all">
                   <div className="flex justify-between items-start mb-6">
                      <div>
                         <div className="text-primary text-[10px] font-black tracking-widest uppercase mb-1">Batch #RR-P-CD-20260514-834</div>
                         <h4 className="text-xl font-headline font-black italic">Arabica - Cau Dat</h4>
                      </div>
                      <div className="text-right">
                         <div className="text-2xl font-headline font-black tabular-nums text-emerald-400">45.0 kg</div>
                         <div className="text-[9px] font-black uppercase opacity-40 italic">Current Yield</div>
                      </div>
                   </div>

                   {/* Micro Roast Runs Visualization */}
                   <div className="flex items-center gap-2 mb-6">
                      {[
                        { id: 1, out: 10 },
                        { id: 2, out: 12.5 },
                        { id: 3, out: 11.8 },
                        { id: 4, out: 10.7 }
                      ].map((run) => (
                        <div key={run.id} className="h-10 w-10 bg-white/10 rounded-lg flex items-center justify-center border border-white/20 group/run hover:bg-primary transition-all cursor-help relative" title={`Run #${run.id}: ${run.out}kg`}>
                           <span className="text-[10px] font-black">{run.id}</span>
                        </div>
                      ))}
                      <button className="h-10 w-10 border-2 border-dashed border-white/20 rounded-lg flex items-center justify-center hover:bg-white/10 transition-all">
                        <span className="material-symbols-outlined !text-sm">add</span>
                      </button>
                   </div>

                   <div className="flex justify-between items-center">
                      <div className="flex gap-6 text-[10px] font-black uppercase opacity-60 italic">
                         <span>Input: 55.0kg</span>
                         <span className="text-emerald-400">Loss: 18.2%</span>
                      </div>
                      <button className="bg-emerald-500 text-white px-6 py-2 rounded-full font-black text-[10px] uppercase tracking-widest hover:scale-105 transition-all">
                        Finalize & Stock-in
                      </button>
                   </div>
                </div>
             </div>
          </div>
        </div>

        {/* Right: Real-time Stock Dashboard */}
        <div className="space-y-8">
           <div className="primary-gradient p-8 rounded-[2.5rem] text-on-primary shadow-2xl relative overflow-hidden flex flex-col h-full">
              <div className="flex items-center gap-4 mb-10">
                 <span className="material-symbols-outlined !text-5xl">inventory_2</span>
                 <div>
                    <h3 className="font-headline font-black text-2xl uppercase tracking-tighter italic">Finished Stock</h3>
                    <p className="text-[10px] font-black opacity-60 uppercase tracking-widest">Live Inventory Levels</p>
                 </div>
              </div>

              <div className="space-y-6 flex-1">
                 {[
                   { sku: 'CD-ARABICA-ROASTED', qty: 1450, color: 'bg-white' },
                   { sku: 'BMT-ROBUSTA-ROASTED', qty: 820, color: 'bg-emerald-400' },
                   { sku: 'GL-ROBUSTA-ROASTED', qty: 310, color: 'bg-white/50' }
                 ].map((item) => (
                    <div key={item.sku} className="bg-white/10 p-5 rounded-2xl border border-white/10">
                       <div className="flex justify-between items-end mb-3">
                          <span className="text-[9px] font-black tracking-widest uppercase opacity-70">{item.sku}</span>
                          <span className="text-2xl font-headline font-black tabular-nums">{item.qty} kg</span>
                       </div>
                       <div className="w-full bg-white/10 h-1.5 rounded-full overflow-hidden p-0.5">
                          <div className={`${item.color} h-full rounded-full`} style={{ width: `${(item.qty/2000)*100}%` }}></div>
                       </div>
                    </div>
                 ))}
              </div>

              <div className="mt-10 p-6 bg-black/20 rounded-2xl border border-white/10 backdrop-blur-xl">
                 <div className="flex items-center gap-3 mb-4">
                    <span className="material-symbols-outlined text-emerald-400 animate-spin-slow">autorenew</span>
                    <span className="text-[10px] font-black uppercase tracking-widest">Kafka Sync Active</span>
                 </div>
                 <p className="text-[10px] font-bold opacity-60 leading-relaxed italic">
                    All inventory updates are broadcasted via topic: <code className="bg-black/40 px-2 py-0.5 rounded text-emerald-400">warehouse.stock.updated</code>
                 </p>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
}
