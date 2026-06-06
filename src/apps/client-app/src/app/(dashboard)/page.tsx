'use client';

import React, { useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import { AUTH_PARAMS } from '@/constants/auth';
import { APP_ROUTES } from '@/constants/routes';
import { useAuth } from '@/lib/auth';
import { storageService } from '@/services/storage.service';
import { Coffee, Key, LogOut, User, ArrowRight, QrCode, ShoppingBag, Info, ShieldCheck, Globe } from 'lucide-react';
import Link from 'next/link';
import { testId } from '@/lib/utils/test-id';

const ArchitectureDiagramCanvas = dynamic(
  () => import('@/components/features/architecture-topology/ArchitectureTopology').then(mod => mod.ArchitectureDiagramCanvas),
  {
    ssr: false,
    loading: () => <div className="w-full h-full bg-slate-50 animate-pulse flex items-center justify-center font-black text-slate-400 italic">Initializing Topology Mesh...</div>
  }
);

export default function LandingPage() {
  const { isAuthenticated, login, logout, refreshSession } = useAuth();
  const [activeTab, setActiveTab] = useState<'topology' | 'registry'>('topology');

  useEffect(() => {
    const url = new URL(window.location.href);
    const token = url.searchParams.get(AUTH_PARAMS.ACCESS_TOKEN);
    if (!token) return;

    storageService.setAccessToken(token);
    url.searchParams.delete(AUTH_PARAMS.ACCESS_TOKEN);
    window.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`);

    refreshSession().finally(() => {
      window.location.href = APP_ROUTES.DASHBOARD.USERS;
    });
  }, [refreshSession]);

  return (
    <div className="flex flex-col h-screen bg-[#f7f9fb] font-body overflow-hidden">
      <header className="h-16 flex justify-between items-center px-8 border-b border-outline-variant/10 bg-white/80 backdrop-blur-xl sticky top-0 z-50">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white shadow-lg">
            <Coffee className="w-5 h-5" />
          </div>
          <h2 className="text-lg font-black font-headline text-on-surface uppercase tracking-tighter italic" {...testId('public-brand')}>
            Runtime Roasters
          </h2>
        </div>

        <nav className="hidden md:flex items-center gap-8">
          <Link href="/" className="text-xs font-black uppercase tracking-widest text-primary border-b-2 border-primary pb-1">Platform</Link>
          <Link href="/coffee" className="text-xs font-black uppercase tracking-widest text-slate-400 hover:text-primary transition-colors">Order Now</Link>
          <Link href="/trace" className="text-xs font-black uppercase tracking-widest text-slate-400 hover:text-primary transition-colors">Trace Registry</Link>
          <Link href="/api-docs" className="text-xs font-black uppercase tracking-widest text-slate-400 hover:text-primary transition-colors">Specs</Link>
        </nav>

        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <div className="flex items-center gap-3 bg-white p-1.5 pr-5 rounded-full border border-outline-variant/10 shadow-sm">
              <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary border border-primary/20 shadow-inner">
                <User className="w-5 h-5" />
              </div>
              <Link
                href={APP_ROUTES.DASHBOARD.USERS}
                className="text-xs font-black text-on-surface hover:text-primary transition-colors uppercase italic tracking-tighter flex items-center gap-1"
              >
                Dashboard <ArrowRight className="w-3 h-3" />
              </Link>
              <button onClick={() => logout()} className="ml-2 p-2 text-slate-400 hover:text-error transition-colors" title="Logout">
                <LogOut className="w-4 h-4" />
              </button>
            </div>
          ) : (
            <button
              onClick={() => login()}
              {...testId('public-access-terminal')}
              className="px-6 py-2.5 bg-slate-900 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:scale-[0.98] transition-all flex items-center gap-2 shadow-lg shadow-slate-900/20"
            >
              <Key className="w-4 h-4" />
              Access Terminal
            </button>
          )}
        </div>
      </header>

      <div className="flex-1 flex overflow-hidden">
        {/* Left Info Panel */}
        <div className="w-[400px] border-r border-outline-variant/5 bg-white overflow-y-auto p-10 flex flex-col gap-12">
           <section>
              <div className="flex items-center gap-2 text-primary mb-4">
                 <ShieldCheck className="w-5 h-5" />
                 <span className="text-[10px] font-black uppercase tracking-[0.3em]">Verified Transparency</span>
              </div>
              <h1 className="text-4xl font-black font-headline text-on-surface tracking-tighter uppercase italic leading-none mb-6">
                Data <span className="text-primary italic">Proof</span> Ecosystem
              </h1>
              <p className="text-sm font-medium text-slate-500 leading-relaxed">
                Experience 100% digital transparency. Every bean is tracked via real-time GPS and cryptographically linked through our distributed service mesh.
              </p>
           </section>

           <div className="space-y-4">
              <Link 
                href="/coffee"
                className="flex items-center justify-between group p-6 bg-slate-50 rounded-[2rem] border border-slate-100 hover:border-primary/30 transition-all shadow-sm"
              >
                 <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-white rounded-2xl flex items-center justify-center shadow-md group-hover:scale-110 transition-transform">
                       <ShoppingBag className="w-6 h-6 text-primary" />
                    </div>
                    <div>
                       <div className="text-[10px] font-black uppercase text-slate-400 tracking-widest">Step 1: Purchase</div>
                       <div className="text-sm font-black uppercase tracking-tight">Buy Coffee Kiosk</div>
                    </div>
                 </div>
                 <ArrowRight className="w-4 h-4 text-slate-300 group-hover:text-primary transition-colors" />
              </Link>

              <Link 
                href="/trace"
                className="flex items-center justify-between group p-6 bg-slate-50 rounded-[2rem] border border-slate-100 hover:border-primary/30 transition-all shadow-sm"
              >
                 <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-white rounded-2xl flex items-center justify-center shadow-md group-hover:scale-110 transition-transform">
                       <QrCode className="w-6 h-6 text-primary" />
                    </div>
                    <div>
                       <div className="text-[10px] font-black uppercase text-slate-400 tracking-widest">Step 2: Verify</div>
                       <div className="text-sm font-black uppercase tracking-tight">Public Trace Registry</div>
                    </div>
                 </div>
                 <ArrowRight className="w-4 h-4 text-slate-300 group-hover:text-primary transition-colors" />
              </Link>
           </div>

           <div className="mt-auto pt-10 border-t border-slate-100">
              <div className="flex items-center gap-2 mb-4">
                 <Globe className="w-4 h-4 text-slate-400" />
                 <span className="text-[9px] font-black uppercase tracking-widest text-slate-400">System Health</span>
              </div>
              <div className="grid grid-cols-2 gap-4">
                 <div className="bg-slate-50 rounded-xl p-4 border border-slate-100">
                    <div className="text-[8px] font-black text-slate-400 uppercase mb-1">Status</div>
                    <div className="flex items-center gap-1.5">
                       <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
                       <span className="text-[10px] font-bold text-slate-700">Operational</span>
                    </div>
                 </div>
                 <div className="bg-slate-50 rounded-xl p-4 border border-slate-100">
                    <div className="text-[8px] font-black text-slate-400 uppercase mb-1">Latency</div>
                    <div className="text-[10px] font-bold text-slate-700">42ms</div>
                 </div>
              </div>
           </div>
        </div>

        {/* Right Topology Panel */}
        <div className="flex-1 flex flex-col min-w-0 bg-white">
          <div className="px-10 py-6 flex justify-between items-center border-b border-outline-variant/5">
            <div>
              <h3 className="text-[10px] font-black text-on-surface-variant/60 uppercase tracking-[0.3em] leading-none">
                Live Distributed Service Map
              </h3>
            </div>
            <div className="flex gap-2">
               <button 
                 onClick={() => setActiveTab('topology')}
                 className={`px-4 py-2 rounded-lg text-[9px] font-black uppercase tracking-widest transition-all ${
                   activeTab === 'topology' ? 'bg-primary text-white shadow-lg' : 'bg-slate-100 text-slate-400 hover:bg-slate-200'
                 }`}
               >
                 Service Mesh
               </button>
            </div>
          </div>

          <div className="flex-1 relative overflow-hidden">
            <ArchitectureDiagramCanvas isPrivate={false} />
          </div>
        </div>
      </div>
    </div>
  );
}
