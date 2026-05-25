'use client';

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { traceService } from '@/services/trace.service';
import TraceTimeline from '@/components/features/traceability/TraceTimeline';
import { Search, QrCode, Loader2, Info } from 'lucide-react';

export default function TraceabilityPage() {
  const [searchId, setSearchId] = useState('');
  const [activeId, setActiveId] = useState<string | null>(null);

  const { data: document, isLoading, isError } = useQuery({
    queryKey: ['trace', activeId],
    queryFn: () => activeId ? traceService.getTraceDocument(activeId) : null,
    enabled: !!activeId,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchId.trim()) {
      setActiveId(searchId.trim());
    }
  };

  return (
    <div className="h-full flex flex-col gap-8">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <QrCode className="w-8 h-8 text-primary" />
            Origin Traceability
          </h1>
          <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            360° Supply Chain Provenance via CQRS Read Model
          </p>
        </div>

        <form onSubmit={handleSearch} className="w-full md:w-96 relative group">
          <input
            type="text"
            value={searchId}
            onChange={(e) => setSearchId(e.target.value)}
            placeholder="Scan QR or Enter Batch ID..."
            className="w-full bg-white border-2 border-slate-200 rounded-2xl py-3 px-12 text-sm font-bold uppercase tracking-tight focus:border-primary transition-all outline-none"
          />
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400 group-focus-within:text-primary transition-colors" />
          <button 
            type="submit"
            className="absolute right-3 top-1/2 -translate-y-1/2 bg-slate-900 text-white px-3 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-primary transition-colors"
          >
            Trace
          </button>
        </form>
      </div>

      <div className="flex-1 bg-slate-50/50 rounded-[3rem] border-2 border-dashed border-slate-200 p-8 md:p-12 overflow-y-auto">
        {!activeId ? (
          <div className="h-full flex flex-col items-center justify-center text-center max-w-md mx-auto py-20 grayscale opacity-40">
            <div className="w-24 h-24 rounded-full bg-slate-200 flex items-center justify-center mb-6">
              <QrCode className="w-12 h-12 text-slate-400" />
            </div>
            <h2 className="text-xl font-black uppercase tracking-tighter text-slate-900 mb-2">Ready to Trace</h2>
            <p className="text-xs font-bold text-slate-500 uppercase tracking-widest leading-relaxed">
              Enter a Batch ID or Shipment ID to visualize the full Farm-to-Cup history of your coffee.
            </p>
          </div>
        ) : isLoading ? (
          <div className="h-full flex flex-col items-center justify-center py-20">
            <Loader2 className="w-12 h-12 text-primary animate-spin mb-4" />
            <p className="text-xs font-black uppercase tracking-widest text-primary italic">Reconstructing Timeline...</p>
          </div>
        ) : isError || !document ? (
          <div className="h-full flex flex-col items-center justify-center text-center max-w-md mx-auto py-20">
            <div className="w-20 h-20 rounded-full bg-red-100 flex items-center justify-center mb-6">
              <Info className="w-10 h-10 text-red-500" />
            </div>
            <h2 className="text-xl font-black uppercase tracking-tighter text-red-600 mb-2">Trace Failed</h2>
            <p className="text-xs font-bold text-slate-500 uppercase tracking-widest leading-relaxed">
              ID <span className="text-slate-900">&quot;{activeId}&quot;</span> not found in the provenance ledger. Please check the ID and try again.
            </p>
            <button 
              onClick={() => setActiveId(null)}
              className="mt-8 text-xs font-black text-primary uppercase tracking-widest border-b-2 border-primary pb-1"
            >
              Clear Search
            </button>
          </div>
        ) : (
          <div>
            <div className="mb-20 flex flex-col md:flex-row items-end justify-between gap-6 border-b-2 border-slate-100 pb-12">
              <div>
                <span className="text-[10px] font-black text-primary uppercase tracking-[0.3em] bg-primary/10 px-3 py-1 rounded-full italic">Ledger Verified</span>
                <h2 className="text-5xl font-black text-slate-900 mt-4 tracking-tighter uppercase italic">{document.entity_id}</h2>
                <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-2">Aggregated across {document.events.length} system events</p>
              </div>
              <div className="bg-white p-6 rounded-3xl border border-slate-200 shadow-xl flex gap-8">
                 <div className="flex flex-col">
                    <span className="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Status</span>
                    <span className="text-xs font-black text-green-500 uppercase italic">Immutable</span>
                 </div>
                 <div className="flex flex-col">
                    <span className="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Consistency</span>
                    <span className="text-xs font-black text-slate-900 uppercase italic">Strong</span>
                 </div>
              </div>
            </div>
            <TraceTimeline document={document} />
          </div>
        )}
      </div>
    </div>
  );
}
