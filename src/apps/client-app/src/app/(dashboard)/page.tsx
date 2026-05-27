'use client';

import React, { useEffect } from 'react';
import { ArchitectureDiagramCanvas } from '@/components/features/architecture-topology/ArchitectureTopology';
import { AUTH_PARAMS } from '@/constants/auth';
import { APP_ROUTES } from '@/constants/routes';
import { useAuth } from '@/lib/auth';
import { storageService } from '@/services/storage.service';
import { Coffee, Key, LogOut, User, ArrowRight } from 'lucide-react';
import Link from 'next/link';
import { testId } from '@/lib/utils/test-id';

export default function DashboardPage() {
  const { isAuthenticated, login, logout, refreshSession } = useAuth();

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
              className="px-6 py-2.5 bg-slate-900 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:scale-[0.98] transition-all flex items-center gap-2"
            >
              <Key className="w-4 h-4" />
              Access Terminal
            </button>
          )}
        </div>
      </header>

      <div className="flex-1 flex overflow-hidden">
        <div className="flex-1 flex flex-col min-w-0">
          <div className="px-10 py-6 flex justify-between items-end bg-white/50 border-b border-outline-variant/5">
            <div>
              <h1 className="text-3xl font-black font-headline text-on-surface tracking-tighter uppercase italic leading-none">
                <span {...testId('public-topology-heading')}>Architecture <span className="text-primary">Topology</span></span>
              </h1>
              <p className="text-[10px] font-bold text-on-surface-variant/60 mt-2 uppercase tracking-widest leading-none">
                Left-to-right Runtime Roasters service map
              </p>
            </div>
          </div>

          <div className="flex-1 relative bg-white overflow-auto p-0">
            <ArchitectureDiagramCanvas />
          </div>
        </div>
      </div>
    </div>
  );
}
