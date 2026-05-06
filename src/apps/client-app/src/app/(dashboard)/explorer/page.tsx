'use client';

import dynamic from 'next/dynamic';
import 'swagger-ui-react/swagger-ui.css';
import { ENV } from '@/constants/env';

const SwaggerUI = dynamic(() => import('swagger-ui-react'), { ssr: false });

export default function ExplorerPage() {
  return (
    <div className="p-8 md:p-12 lg:p-16 relative z-10">
      <div className="mb-10">
        <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-primary-fixed/30 text-on-primary-fixed text-xs font-black uppercase tracking-[0.2em] mb-4 border border-primary/10 shadow-sm italic">
          <span className="material-symbols-outlined text-[14px]">dns</span>
          <span>Core Infrastructure</span>
        </div>
        <h2 className="text-5xl font-black font-headline tracking-tighter uppercase italic text-on-surface">
          API <span className="text-primary">Explorer</span>
        </h2>
        <p className="text-on-surface-variant font-medium text-lg mt-4 max-w-2xl leading-relaxed">
          Interact with the Runtime Roasters microservices mesh via the KrakenD Gateway. 
          Real-time contract validation for distributed supply chain events.
        </p>
      </div>
      
      <div className="border-[6px] border-surface-container-high rounded-[2rem] overflow-hidden shadow-2xl bg-white relative group transition-all duration-500 hover:border-primary/20">
        <div className="absolute inset-0 pointer-events-none border border-outline-variant/20 rounded-[1.8rem] z-20"></div>
        <div className="p-4 bg-surface-container-low border-b border-outline-variant/10 flex items-center justify-between">
           <div className="flex gap-1.5">
              <div className="w-2.5 h-2.5 rounded-full bg-error/30"></div>
              <div className="w-2.5 h-2.5 rounded-full bg-amber-500/30"></div>
              <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/30"></div>
           </div>
           <span className="font-mono text-[10px] font-black text-slate-400 uppercase tracking-widest italic">{ENV.SWAGGER_JSON_URL.replace('http://', '').replace('https://', '')}</span>
        </div>
        <div className="min-h-[600px] relative z-10">
          <SwaggerUI url={ENV.SWAGGER_JSON_URL} />
        </div>
      </div>
      
      {/* Decorative Graphics */}
      <div className="absolute bottom-0 left-0 w-64 h-64 bg-tertiary-fixed/10 rounded-full blur-[100px] pointer-events-none z-0"></div>
    </div>
  );
}
