"use client";

import React, { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { warehouseService, ProductionBatch, Inventory } from '@/services/warehouse.service';
import { Loader2, Package, TrendingUp, AlertTriangle, CheckCircle2, Factory, Database, ArrowRight } from 'lucide-react';

export default function WarehouseStockLogicPage() {
  const queryClient = useQueryClient();
  const [intakeWeight, setIntakeWeight] = useState<number>(0);
  const [selectedBatchId, setSelectedBatchId] = useState<string | null>(null);

  const { data: batches = [], isLoading: loadingBatches } = useQuery({
    queryKey: ['warehouse', 'batches'],
    queryFn: () => warehouseService.listBatches(),
    refetchInterval: 5000,
  });

  const { data: inventory = [], isLoading: loadingInventory } = useQuery({
    queryKey: ['warehouse', 'inventory'],
    queryFn: () => warehouseService.getInventory(),
    refetchInterval: 10000,
  });

  const updateIntakeMutation = useMutation({
    mutationFn: ({ id, weight }: { id: string, weight: number }) => warehouseService.updateIntake(id, weight),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouse', 'batches'] });
      setIntakeWeight(0);
      setSelectedBatchId(null);
    }
  });

  const selectedBatch = useMemo(() => 
    batches.find(b => b.id === selectedBatchId) || batches.find(b => b.status === 'RECEIVED'),
  [batches, selectedBatchId]);

  const originalWeight = selectedBatch?.intake_weight || 500;
  const deviation = Math.abs((originalWeight - intakeWeight) / originalWeight * 100);
  const showWarning = deviation > 2 && intakeWeight > 0;

  return (
    <div className="min-h-full p-0 flex flex-col gap-8 font-body selection:bg-primary-fixed selection:text-on-primary-fixed relative overflow-hidden bg-transparent">
      {/* Header Section */}
      <header className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6 relative z-10">
        <div>
          <h1 className="text-4xl font-black font-headline tracking-tighter text-slate-900 uppercase italic">Warehouse <span className="text-primary">Operations</span></h1>
          <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-1 italic">Inventory Management & Processing Lifecycle</p>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12 relative z-10 flex-1">
        {/* Left: Intake Command Center */}
        <div className="lg:col-span-2 space-y-8 flex flex-col min-h-0">
          <div className="bg-white p-10 rounded-[3rem] border border-slate-200 shadow-xl relative overflow-hidden group flex-shrink-0">
            <div className="flex justify-between items-center mb-8">
              <h3 className="font-headline font-black text-2xl uppercase tracking-tight text-slate-900 flex items-center gap-3">
                 <Package className="w-6 h-6 text-primary" />
                 1. Intake Verification
              </h3>
              {selectedBatch ? (
                <span className="bg-primary/10 text-primary px-4 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-primary/20">
                  Awaiting: Batch #{selectedBatch.batch_id}
                </span>
              ) : (
                <span className="bg-slate-100 text-slate-400 px-4 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-slate-200">
                  No Pending Intake
                </span>
              )}
            </div>

            {selectedBatch && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
                <div className="space-y-6">
                  <div className="bg-slate-50 p-6 rounded-3xl border border-slate-100">
                    <span className="text-[10px] font-black text-slate-400 uppercase block mb-2 tracking-widest opacity-60">Farm Declared Weight</span>
                    <div className="text-4xl font-headline font-black text-slate-900 tabular-nums">{originalWeight} kg</div>
                    <p className="text-xs text-slate-500 mt-2 font-bold uppercase italic">Origin: {selectedBatch.origin_code} ({selectedBatch.coffee_type})</p>
                  </div>

                  <div className="space-y-4">
                    <label className="text-[10px] font-black text-primary uppercase block tracking-widest ml-2">Actual Intake Weight (kg)</label>
                    <input 
                      type="number" 
                      value={intakeWeight || ''} 
                      onChange={(e) => setIntakeWeight(Number(e.target.value))}
                      placeholder="Enter scale value..."
                      className="w-full bg-white border-2 border-slate-200 rounded-2xl px-6 py-4 text-2xl font-headline font-black focus:border-primary transition-all outline-none"
                    />
                  </div>
                </div>

                <div className="flex flex-col justify-center items-center p-8 bg-slate-50/50 rounded-[2.5rem] border-2 border-dashed border-slate-200">
                  {intakeWeight > 0 ? (
                    <div className="text-center">
                      <div className={`text-6xl font-headline font-black mb-2 ${showWarning ? 'text-red-500 animate-bounce' : 'text-green-500'}`}>
                        {deviation.toFixed(2)}%
                      </div>
                      <span className="text-[10px] font-black uppercase tracking-widest text-slate-400">Weight Deviation</span>
                      
                      {showWarning && (
                        <div className="mt-6 p-4 bg-red-50 text-red-600 rounded-2xl border border-red-100 flex flex-col gap-2">
                          <div className="flex items-center gap-2 font-black text-[10px] uppercase">
                            <AlertTriangle className="w-3.5 h-3.5" />
                            Anomaly Threshold Exceeded
                          </div>
                          <p className="text-[9px] font-bold italic">Manual reason entry required in production.</p>
                        </div>
                      )}
                    </div>
                  ) : (
                    <div className="text-center opacity-30">
                      <div className="w-20 h-20 bg-slate-200 rounded-full flex items-center justify-center mx-auto mb-4">
                        <TrendingUp className="w-10 h-10 text-slate-400" />
                      </div>
                      <p className="text-[10px] font-black uppercase tracking-widest text-slate-500">Awaiting Scale Input</p>
                    </div>
                  )}
                </div>
              </div>
            )}

            <div className="mt-10 flex justify-end">
              <button 
                onClick={() => selectedBatch && updateIntakeMutation.mutate({ id: selectedBatch.id, weight: intakeWeight })}
                disabled={intakeWeight <= 0 || updateIntakeMutation.isPending || !selectedBatch}
                className={`px-10 py-4 rounded-2xl font-black text-xs uppercase tracking-widest shadow-2xl transition-all flex items-center gap-3 ${
                  intakeWeight > 0 ? 'bg-slate-900 text-white hover:bg-primary' : 'bg-slate-100 text-slate-400 opacity-50'
                }`}
              >
                {updateIntakeMutation.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Confirm Intake'}
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </div>

          {/* Active Production Batches */}
          <div className="bg-slate-900 text-white p-10 rounded-[3rem] shadow-2xl relative overflow-hidden flex-1 flex flex-col min-h-0">
             <div className="flex justify-between items-center mb-8 shrink-0">
                <h3 className="font-headline font-black text-2xl uppercase tracking-tighter italic flex items-center gap-3">
                  <Factory className="w-6 h-6 text-primary" />
                  2. Processing Queue
                </h3>
                <div className="flex gap-2">
                   <div className="w-3 h-3 bg-green-400 rounded-full animate-pulse"></div>
                   <span className="text-[10px] font-black uppercase opacity-60">{batches.filter(b => b.status === 'PROCESSING').length} Active</span>
                </div>
             </div>

             <div className="space-y-6 overflow-y-auto flex-1 pr-2 custom-scrollbar">
                {batches.filter(b => b.status !== 'STOCKED').length === 0 ? (
                  <div className="h-full flex flex-col items-center justify-center opacity-20 py-10">
                    <Database className="w-16 h-16 mb-4" />
                    <p className="text-xs font-black uppercase tracking-widest">No active processing</p>
                  </div>
                ) : (
                  batches.filter(b => b.status !== 'STOCKED').map((batch) => (
                    <div key={batch.id} className="bg-white/5 border border-white/10 p-6 rounded-[2rem] group hover:bg-white/10 transition-all cursor-pointer">
                      <div className="flex justify-between items-start mb-6">
                          <div>
                            <div className="text-primary text-[10px] font-black tracking-widest uppercase mb-1">Batch #{batch.batch_id}</div>
                            <h4 className="text-xl font-headline font-black italic">{batch.coffee_type} - {batch.origin_code}</h4>
                          </div>
                          <div className="text-right">
                            <div className="text-2xl font-headline font-black tabular-nums text-green-400">{batch.total_output_weight.toFixed(1)} kg</div>
                            <div className="text-[9px] font-black uppercase opacity-40 italic">Current Yield</div>
                          </div>
                      </div>

                      <div className="flex items-center gap-2 mb-6">
                          {batch.roast_runs?.map((run) => (
                            <div key={run.id} className="h-10 w-10 bg-white/10 rounded-xl flex items-center justify-center border border-white/20 hover:bg-primary transition-all">
                              <span className="text-[10px] font-black">{run.run_number}</span>
                            </div>
                          ))}
                          <button className="h-10 w-10 border-2 border-dashed border-white/20 rounded-xl flex items-center justify-center hover:bg-white/10 transition-all">
                            <span className="material-symbols-outlined !text-sm">add</span>
                          </button>
                      </div>

                      <div className="flex justify-between items-center">
                          <div className="flex gap-6 text-[10px] font-black uppercase opacity-60 italic">
                            <span>Input: {batch.total_input_weight}kg</span>
                            <span className={batch.weight_loss_percent > 15 ? 'text-red-400' : 'text-green-400'}>
                              Loss: {batch.weight_loss_percent.toFixed(1)}%
                            </span>
                          </div>
                          <div className="flex items-center gap-3">
                             <span className="text-[8px] font-black text-slate-500 uppercase tracking-widest">{batch.status}</span>
                             <button 
                                onClick={() => warehouseService.finalizeBatch(batch.id)}
                                className="bg-white text-slate-900 px-6 py-2 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-primary hover:text-white transition-all shadow-xl"
                             >
                                Finalize & Stock-in
                             </button>
                          </div>
                      </div>
                    </div>
                  ))
                )}
             </div>
          </div>
        </div>

        {/* Right: Real-time Stock Dashboard */}
        <div className="space-y-8 flex flex-col min-h-0">
           <div className="bg-white p-10 rounded-[3.5rem] border border-slate-200 shadow-2xl relative overflow-hidden flex flex-col h-full">
              <div className="flex items-center gap-4 mb-10 shrink-0">
                 <div className="w-16 h-16 rounded-3xl bg-primary flex items-center justify-center text-white shadow-lg shadow-primary/20">
                    <Database className="w-8 h-8" />
                 </div>
                 <div>
                    <h3 className="font-headline font-black text-2xl uppercase tracking-tighter italic text-slate-900">Finished Stock</h3>
                    <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Live Inventory Node</p>
                 </div>
              </div>

              <div className="space-y-6 flex-1 overflow-y-auto pr-2 custom-scrollbar">
                 {loadingInventory ? (
                   <div className="py-20 text-center animate-pulse">
                      <Loader2 className="w-8 h-8 animate-spin mx-auto text-slate-300 mb-2" />
                      <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Syncing Stock...</span>
                   </div>
                 ) : inventory.length === 0 ? (
                   <div className="py-20 text-center opacity-30">
                      <p className="text-[10px] font-black uppercase tracking-widest italic">No stock-in recorded</p>
                   </div>
                 ) : (
                   inventory.map((item) => (
                    <div key={item.id} className="bg-slate-50 p-6 rounded-3xl border border-slate-100 hover:border-primary/20 transition-all group">
                       <div className="flex justify-between items-end mb-3">
                          <div>
                            <span className="text-[8px] font-black tracking-widest uppercase text-slate-400 block mb-1">SKU: {item.sku}</span>
                            <span className="text-[11px] font-black text-slate-900 uppercase italic tracking-tighter">{item.coffee_type} ({item.origin_code})</span>
                          </div>
                          <span className="text-2xl font-headline font-black tabular-nums text-slate-900">{item.available_quantity.toFixed(1)} <span className="text-xs font-body opacity-40">kg</span></span>
                       </div>
                       <div className="w-full bg-slate-200 h-1.5 rounded-full overflow-hidden p-0.5">
                          <div className={`bg-primary h-full rounded-full transition-all duration-1000 group-hover:scale-x-105`} style={{ width: `${Math.min((item.available_quantity/2000)*100, 100)}%` }}></div>
                       </div>
                    </div>
                   ))
                 )}
              </div>

              <div className="mt-10 p-6 bg-slate-900 rounded-[2rem] text-white shrink-0">
                 <div className="flex items-center gap-3 mb-4">
                    <CheckCircle2 className="w-4 h-4 text-green-400 animate-pulse" />
                    <span className="text-[10px] font-black uppercase tracking-widest">Kafka Node: Active</span>
                 </div>
                 <p className="text-[9px] font-bold text-slate-400 leading-relaxed italic uppercase tracking-tighter">
                    Topic: <code className="text-primary font-mono lowercase tracking-normal">warehouse.stock.updated</code>
                 </p>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
}

