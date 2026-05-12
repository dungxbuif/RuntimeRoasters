'use client';

import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';

export interface FlowScenario {
  id: string;
  name: string;
  description: string;
}

interface FlowControlPanelProps {
  scenarios: FlowScenario[];
  activeScenarioId: string | null;
  onSelectScenario: (id: string) => void;
  onPlay: () => void;
  onPause: () => void;
  onReset: () => void;
  isPlaying: boolean;
  progress: number; // 0 to 100
  history: Array<{ id: string; name: string; timestamp: string }>;
}

export const FlowControlPanel = ({
  scenarios,
  activeScenarioId,
  onSelectScenario,
  onPlay,
  onPause,
  onReset,
  isPlaying,
  progress,
  history
}: FlowControlPanelProps) => {
  return (
    <div className="w-full h-full bg-slate-50 flex flex-col p-6 space-y-8 overflow-y-auto">
      <div>
        <h3 className="text-[10px] font-black uppercase tracking-[0.3em] text-slate-400 mb-6 flex items-center gap-2">
          <span className="w-2 h-2 bg-primary rounded-full animate-pulse"></span>
          Scenario Controller
        </h3>
        
        <div className="space-y-3">
          {scenarios.map((s) => (
            <button
              key={s.id}
              onClick={() => onSelectScenario(s.id)}
              className={`w-full text-left p-4 rounded-xl border transition-all ${
                activeScenarioId === s.id 
                  ? 'bg-white border-primary shadow-[10px_10px_0px_0px_rgba(0,74,198,0.05)]' 
                  : 'bg-white border-outline-variant/10 hover:border-outline-variant/30'
              }`}
            >
              <p className={`text-xs font-black uppercase tracking-tight ${activeScenarioId === s.id ? 'text-primary' : 'text-on-surface'}`}>
                {s.name}
              </p>
              <p className="text-[10px] text-on-surface-variant mt-1 line-clamp-2 leading-relaxed">{s.description}</p>
            </button>
          ))}
        </div>
      </div>

      <div className="bg-white rounded-2xl p-6 border border-outline-variant/20 shadow-sm">
        <h4 className="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-4 italic">Execution State</h4>
        
        <div className="flex gap-4 mb-6">
          <button 
            onClick={isPlaying ? onPause : onPlay}
            disabled={!activeScenarioId}
            className={`flex-1 py-3 rounded-xl flex items-center justify-center gap-2 transition-all font-black uppercase tracking-widest text-[10px] ${
              !activeScenarioId ? 'opacity-30 grayscale cursor-not-allowed' :
              isPlaying ? 'bg-amber-50 text-amber-600 border border-amber-200' : 'bg-primary text-white shadow-lg shadow-primary/20 hover:scale-[0.98]'
            }`}
          >
            <span className="material-symbols-outlined !text-sm">
              {isPlaying ? 'pause' : 'play_arrow'}
            </span>
            <span>
              {isPlaying ? 'Pause' : 'Execute'}
            </span>
          </button>
          <button 
            onClick={onReset}
            className="w-12 h-12 rounded-xl border border-outline-variant/10 flex items-center justify-center hover:bg-slate-100 transition-colors text-slate-400"
          >
            <span className="material-symbols-outlined !text-sm">restart_alt</span>
          </button>
        </div>

        <div className="space-y-2">
          <div className="flex justify-between text-[8px] font-black text-slate-400 uppercase tracking-widest">
            <span>Flow Completion</span>
            <span className="text-primary">{Math.round(progress)}%</span>
          </div>
          <div className="h-1.5 w-full bg-slate-100 rounded-full overflow-hidden p-0.5">
            <motion.div 
              className="h-full bg-primary rounded-full shadow-[0_0_8px_rgba(0,74,198,0.3)]" 
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
            />
          </div>
        </div>
      </div>

      <div className="flex-1 min-h-0 flex flex-col">
        <h4 className="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-4 italic">Telemetry History</h4>
        <div className="flex-1 overflow-y-auto pr-2 custom-scrollbar">
          <div className="space-y-3">
            <AnimatePresence initial={false}>
              {history.map((item, i) => (
                <motion.div 
                  key={`${item.id}-${i}`}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="flex gap-3 text-[10px] border-l-2 border-slate-200 pl-4 py-1"
                >
                  <span className="text-slate-400 tabular-nums font-mono">{item.timestamp}</span>
                  <span className="font-bold text-on-surface uppercase tracking-tighter truncate">{item.name}</span>
                </motion.div>
              ))}
              {history.length === 0 && (
                <p className="text-[9px] text-slate-400 italic">No executions logged.</p>
              )}
            </AnimatePresence>
          </div>
        </div>
      </div>
    </div>
  );
};
