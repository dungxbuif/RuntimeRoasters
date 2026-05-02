'use client';

import React from 'react';
import { useDemoPing } from '@/hooks/useDemo';
import { HeroTitle, MicroLabel } from '@/components/ui/atoms';
import { Terminal, Activity, CheckCircle2 } from 'lucide-react';

export default function Dashboard() {
  const { data, isLoading, error, refetch } = useDemoPing(false);

  return (
    <div className="p-8 md:p-12 lg:p-16 space-y-12">
      <header>
        <HeroTitle>System <span className="text-tertiary">Verification</span></HeroTitle>
        <p className="text-on-surface-variant mt-4 font-medium text-lg max-w-2xl leading-relaxed">
          Verify the end-to-end integration using the production-ready service layer.
        </p>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        <div className="bg-white p-10 rounded-[2.5rem] border border-outline-variant/15 shadow-xl relative overflow-hidden group">
          <div className="absolute top-0 left-0 w-2 h-full bg-primary"></div>
          <MicroLabel className="mb-6">Test Execution</MicroLabel>
          <h3 className="text-2xl font-black font-headline uppercase italic mb-8 text-on-surface">Service Ping-Pong</h3>
          
          <div className="space-y-6">
            <button 
              onClick={() => refetch()}
              disabled={isLoading}
              className="w-full py-4 bg-primary text-white rounded-2xl font-black uppercase tracking-widest hover:scale-[0.98] transition-all flex items-center justify-center gap-3 shadow-lg shadow-primary/20"
            >
              <Terminal className="w-5 h-5" />
              {isLoading ? 'Executing Request...' : 'Trigger API Call'}
            </button>

            <div className="bg-slate-900 rounded-2xl p-6 font-mono text-xs overflow-hidden border-2 border-slate-800 shadow-inner">
              <div className="flex justify-between text-slate-500 mb-4 border-b border-slate-800 pb-2 uppercase tracking-tighter">
                <span>Response Output</span>
                <span className="text-tertiary-fixed font-black italic">status_200</span>
              </div>
              {data ? (
                <pre className="text-emerald-400">
                  {JSON.stringify(data, null, 2)}
                </pre>
              ) : error ? (
                <pre className="text-error italic">
                  ERROR: System unreachable.
                </pre>
              ) : (
                <pre className="text-slate-600 italic">Waiting for trigger...</pre>
              )}
            </div>
          </div>
        </div>

        <div className="bg-surface-container-low p-10 rounded-[2.5rem] border border-outline-variant/10 shadow-sm relative overflow-hidden">
          <MicroLabel className="mb-6">Observability Guide</MicroLabel>
          <h3 className="text-2xl font-black font-headline uppercase italic mb-8 text-on-surface">Verify OpenTelemetry</h3>
          <div className="space-y-8">
            <div className="flex gap-4 group text-tertiary">
              <div className="w-10 h-10 rounded-full bg-white flex items-center justify-center border border-tertiary shrink-0 shadow-lg shadow-tertiary/10">
                <CheckCircle2 className="w-5 h-5" />
              </div>
              <p className="text-sm font-black leading-relaxed">
                Check <a href="http://localhost:16686" target="_blank" className="underline">Jaeger UI</a> for traces.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
