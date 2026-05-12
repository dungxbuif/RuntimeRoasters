'use client';

import React, { useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface TraceLogEntry {
  id: string;
  timestamp: string;
  step: number;
  description: string;
}

interface TraceLogProps {
  entries: TraceLogEntry[];
}

export const TraceLog = ({ entries }: TraceLogProps) => {
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [entries]);

  return (
    <div className="w-full h-full bg-slate-950 flex flex-col overflow-hidden shadow-[inset_0_4px_20px_rgba(0,0,0,0.5)]">
      {/* Console Header */}
      <div className="px-6 py-2 border-b border-white/10 bg-slate-900/50 flex justify-between items-center shrink-0">
        <div className="flex gap-4 items-center">
          <div className="flex gap-1.5">
            <div className="w-2.5 h-2.5 rounded-full bg-[#ff5f56]"></div>
            <div className="w-2.5 h-2.5 rounded-full bg-[#ffbd2e]"></div>
            <div className="w-2.5 h-2.5 rounded-full bg-[#27c93f]"></div>
          </div>
          <div className="h-4 w-[1px] bg-white/10 mx-1"></div>
          <div className="flex gap-2 items-center">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_8px_#10b981]"></div>
            <span className="text-[9px] font-black uppercase tracking-[0.2em] text-slate-400">Live_Trace_Stream_v4.2</span>
          </div>
        </div>
        <div className="flex gap-6 items-center">
          <span className="text-[8px] font-mono text-slate-600 font-bold uppercase tracking-widest">PROPAGATION: ACTIVE</span>
          <span className="text-[8px] font-mono text-slate-700">TTY_NODE_01</span>
        </div>
      </div>
      
      {/* Console Body */}
      <div 
        ref={scrollRef}
        className="flex-1 overflow-y-auto p-6 px-10 space-y-2.5 font-mono text-[10px] custom-scrollbar selection:bg-primary/30"
      >
        <AnimatePresence initial={false}>
          {entries.map((entry) => (
            <motion.div 
              key={entry.id}
              initial={{ opacity: 0, x: -5 }}
              animate={{ opacity: 1, x: 0 }}
              className="flex gap-8 group border-l border-white/5 pl-4 hover:border-primary/40 transition-colors relative"
            >
              <div className="absolute left-[-1px] top-0 bottom-0 w-[2px] bg-primary scale-y-0 group-hover:scale-y-100 transition-transform origin-top"></div>
              <span className="text-slate-600 shrink-0 font-bold w-20">[{entry.timestamp}]</span>
              <span className="text-primary font-black shrink-0 tracking-tighter w-16">STEP_{entry.step}</span>
              <span className="text-slate-300 italic group-hover:text-white transition-colors leading-relaxed">{entry.description}</span>
            </motion.div>
          ))}
          {entries.length === 0 && (
            <div className="flex items-center gap-3 text-slate-700 italic py-4">
              <span className="animate-pulse font-black text-primary w-4">&gt;_</span>
              <span className="uppercase text-[9px] tracking-[0.2em] font-black opacity-30">Listening for inbound cluster telemetry...</span>
            </div>
          )}
        </AnimatePresence>
      </div>

      {/* Console Footer */}
      <div className="px-6 py-1.5 border-t border-white/5 bg-slate-900/30 flex justify-between items-center shrink-0">
         <div className="flex gap-4">
            <span className="text-[7px] font-bold text-slate-700 uppercase tracking-widest">UTF-8</span>
            <span className="text-[7px] font-bold text-slate-700 uppercase tracking-widest">Line: {entries.length}</span>
         </div>
         <span className="text-[7px] font-bold text-slate-800 uppercase">RuntimeRoasters Terminal</span>
      </div>
    </div>
  );
};
