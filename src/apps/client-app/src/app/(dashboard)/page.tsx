'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/lib/auth';
import { HighFidelityTopology } from '@/components/features/chaos-topology/high-fidelity/HighFidelityTopology';
import { FlowControlPanel, FlowScenario } from '@/components/features/chaos-topology/FlowControlPanel';
import { TraceLog } from '@/components/features/chaos-topology/TraceLog';
import { LogOut, User, Key, Coffee, ArrowRight } from 'lucide-react';
import Link from 'next/link';
import { APP_ROUTES } from '@/constants/routes';

const SCENARIOS: FlowScenario[] = [
  { 
    id: 'oidc-login', 
    name: 'OIDC Login Flow', 
    description: 'Trace how identity is verified across KrakenD, Hydra, and Kratos.' 
  },
  { 
    id: 'outbox-sync', 
    name: 'Transactional Outbox', 
    description: 'Simulate data persistence and eventual consistency via Kafka relay.' 
  },
  { 
    id: 'resilient-casbin', 
    name: 'Resilient Casbin Sync', 
    description: 'Bootstrapping and self-healing of distributed authorization rules.' 
  }
];

const OIDC_LOGS = [
  'User initiates login challenge via Browser HTTPS.',
  'KrakenD forwards OAuth2 discovery to Ory Hydra.',
  'Hydra delegates identity lookup to Ory Kratos cluster.',
  'Auth Service validates session via internal gRPC channel.',
  'Farm Service allows access based on synced Casbin rules.'
];

export default function DashboardPage() {
  const { isAuthenticated, login, logout } = useAuth();
  
  // Animation State
  const [activeScenario, setActiveScenario] = useState<string | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentStep, setCurrentStep] = useState(-1);
  const [traceLogs, setTraceLogs] = useState<{ id: string; timestamp: string; step: number; description: string }[]>([]);
  const [history, setHistory] = useState<{ id: string; name: string; timestamp: string }[]>([]);

  const resetAnimation = useCallback(() => {
    setIsPlaying(false);
    setCurrentStep(-1);
    setTraceLogs([]);
  }, []);

  const playScenario = () => {
    if (!activeScenario) return;
    setIsPlaying(true);
  };

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isPlaying && activeScenario === 'oidc-login') {
      interval = setInterval(() => {
        setCurrentStep(prev => {
          const next = prev + 1;
          if (next >= OIDC_LOGS.length) {
            setIsPlaying(false);
            setHistory(h => [{ 
              id: Date.now().toString(), 
              name: 'OIDC Login Flow', 
              timestamp: new Date().toLocaleTimeString() 
            }, ...h.slice(0, 4)]);
            return prev;
          }

          // Update Trace Log
          setTraceLogs(logs => [...logs, {
            id: Date.now().toString(),
            timestamp: new Date().toLocaleTimeString(),
            step: next + 1,
            description: OIDC_LOGS[next]
          }]);

          return next;
        });
      }, 1500);
    }
    return () => clearInterval(interval);
  }, [isPlaying, activeScenario]);

  const progress = activeScenario === 'oidc-login' 
    ? ((currentStep + 1) / OIDC_LOGS.length) * 100 
    : 0;

  return (
    <div className="flex flex-col h-screen bg-[#f7f9fb] font-body overflow-hidden">
      {/* HEADER */}
      <header className="h-16 flex justify-between items-center px-8 border-b border-outline-variant/10 bg-white/70 backdrop-blur-xl sticky top-0 z-50">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white shadow-lg">
            <Coffee className="w-5 h-5" />
          </div>
          <h2 className="text-lg font-black font-headline text-on-surface uppercase tracking-tighter italic">Runtime Roasters</h2>
        </div>

        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <div className="flex items-center gap-3 bg-white/40 p-1.5 pr-5 rounded-full border border-outline-variant/10 shadow-sm backdrop-blur-md">
              <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary border border-primary/20 shadow-inner">
                <User className="w-5 h-5" />
              </div>
              <div className="text-left">
                <p className="text-[8px] font-black text-slate-500 uppercase leading-none mb-1">Authenticated</p>
                <Link 
                  href={APP_ROUTES.DASHBOARD.USERS}
                  className="text-xs font-black text-on-surface hover:text-primary transition-colors uppercase italic tracking-tighter truncate max-w-[120px] flex items-center gap-1"
                >
                  Dashboard <ArrowRight className="w-3 h-3" />
                </Link>
              </div>
              <button 
                onClick={() => logout()}
                className="ml-4 p-2 text-slate-400 hover:text-error transition-colors"
                title="Logout"
              >
                <LogOut className="w-4 h-4" />
              </button>
            </div>
          ) : (
            <button 
              onClick={() => login()}
              className="px-6 py-2.5 bg-slate-900 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:scale-[0.98] transition-all flex items-center gap-2"
            >
              <Key className="w-4 h-4" />
              Access Terminal
            </button>
          )}
        </div>
      </header>

      {/* MAIN CONTENT AREA */}
      <div className="flex-1 flex overflow-hidden">
        
        {/* LEFT COLUMN: HERO + DIAGRAM + TRACE */}
        <div className="flex-1 flex flex-col min-w-0">
          
          {/* HERO ZONE */}
          <div className="px-10 py-6 flex justify-between items-end bg-white/30 border-b border-outline-variant/5">
            <div>
              <h1 className="text-3xl font-black font-headline text-on-surface tracking-tighter uppercase italic leading-none">
                System <span className="text-primary">Topography</span>
              </h1>
              <p className="text-[10px] font-bold text-on-surface-variant/60 mt-2 uppercase tracking-widest leading-none">
                Interactive Architecture Showcase • Alpha Cluster v4.2
              </p>
            </div>
            <div className="relative group cursor-default hidden md:block">
               <p className="font-marker text-2xl text-primary rotate-[-4deg] opacity-60">
                  Cluster Map ↙
               </p>
            </div>
          </div>

          {/* DIAGRAM AREA */}
          <div className="flex-1 relative bg-white/5 overflow-auto flex items-center justify-center">
              <HighFidelityTopology activeStep={currentStep} />
          </div>

          {/* TRACE LOG AREA */}
          <div className="h-48 shrink-0">
            <TraceLog entries={traceLogs} />
          </div>
        </div>

        {/* RIGHT COLUMN: FLOW CONTROL PANEL */}
        <div className="w-80 h-full border-l border-outline-variant/10 bg-white">
          <FlowControlPanel 
            scenarios={SCENARIOS}
            activeScenarioId={activeScenario}
            onSelectScenario={(id) => { resetAnimation(); setActiveScenario(id); }}
            onPlay={playScenario}
            onPause={() => setIsPlaying(false)}
            onReset={resetAnimation}
            isPlaying={isPlaying}
            progress={progress}
            history={history}
          />
        </div>

      </div>
    </div>
  );
}
