'use client';

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { traceService } from '@/services/trace.service';
import TraceTimeline from '@/components/features/traceability/TraceTimeline';

export default function TraceabilityJourneyPage() {
  const [searchId, setSearchId] = useState('');
  const [activeId, setActiveId] = useState('');

  const { data: document, isLoading, error } = useQuery({
    queryKey: ['trace', activeId],
    queryFn: () => traceService.getTraceDocument(activeId),
    enabled: !!activeId,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setActiveId(searchId);
  };

  return (
    <div className="min-h-full bg-surface text-on-surface p-8 md:p-12 lg:p-20 font-body selection:bg-tertiary-fixed selection:text-on-tertiary-fixed relative overflow-x-hidden">
      {/* Page Header */}
      <div className="max-w-4xl mx-auto mb-20 relative z-10">
        <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-secondary-container text-on-secondary-container text-sm font-medium mb-4 shadow-sm border border-outline-variant/10">
          <span className="material-symbols-outlined text-[18px]">history</span>
          <span>CQRS Traceability Engine</span>
        </div>
        <h1 className="font-headline text-5xl md:text-6xl font-bold tracking-tight text-on-background mb-4 uppercase italic">
          Batch <span className="text-primary">Time Machine</span>
        </h1>
        <p className="font-body text-lg text-on-surface-variant max-w-2xl leading-relaxed font-bold uppercase tracking-tighter italic mb-10">
          Visualizing the full journey of a coffee bean from farm to cup.
        </p>

        {/* Search Bar */}
        <form onSubmit={handleSearch} className="flex gap-4 max-w-xl">
          <div className="relative flex-1">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 material-symbols-outlined text-slate-400">search</span>
            <input 
              type="text" 
              value={searchId}
              onChange={(e) => setSearchId(e.target.value)}
              placeholder="Enter Batch ID, Order ID or Shipment ID..."
              className="w-full pl-12 pr-4 py-4 bg-white border-2 border-slate-200 rounded-2xl focus:border-primary outline-none transition-all font-bold text-slate-700 uppercase italic tracking-wider"
            />
          </div>
          <button 
            type="submit"
            className="px-8 bg-primary text-white rounded-2xl font-black uppercase italic tracking-widest hover:scale-95 transition-all shadow-lg shadow-primary/20"
          >
            Trace
          </button>
        </form>
      </div>

      {isLoading && (
        <div className="flex flex-col items-center justify-center py-20 space-y-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
          <p className="font-marker text-2xl text-slate-400">Rebuilding journey from Kafka events...</p>
        </div>
      )}

      {error && (
        <div className="max-w-xl mx-auto p-8 bg-error-container text-on-error-container rounded-3xl border-2 border-error/20 text-center">
          <span className="material-symbols-outlined text-5xl mb-4">error</span>
          <h3 className="text-2xl font-headline font-black uppercase italic mb-2">Trace Failed</h3>
          <p className="font-bold opacity-80">Could not retrieve traceability data for this ID.</p>
        </div>
      )}

      {document && <TraceTimeline document={document} />}

      {!activeId && !isLoading && (
        <div className="max-w-2xl mx-auto text-center py-20 border-4 border-dashed border-slate-200 rounded-[3rem] opacity-40">
           <span className="material-symbols-outlined text-8xl mb-6">query_stats</span>
           <h3 className="text-3xl font-headline font-black uppercase italic">No Active Trace</h3>
           <p className="font-marker text-xl mt-4">Enter an ID above to visualize the distributed transaction flow.</p>
        </div>
      )}
    </div>
  );
}
