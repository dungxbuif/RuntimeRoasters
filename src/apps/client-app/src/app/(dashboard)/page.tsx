'use client';

import React, { useEffect, useState, Suspense, useRef } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useDemoPing } from '@/hooks/useDemo';
import { HeroTitle, MicroLabel } from '@/components/ui/atoms';
import { Terminal, CheckCircle2, User, LogOut, Key } from 'lucide-react';
import { storageService } from '@/services/storage.service';
import { isBrowser } from '@/lib/utils/window';
import { AUTH_PARAMS, OIDC_CONFIG } from '@/constants/auth';
import { APP_ROUTES } from '@/constants/routes';
import { ENV } from '@/constants/env';

function DashboardContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const initialized = useRef(false);
  
  const { data, isLoading, error, refetch } = useDemoPing(!!token);

  useEffect(() => {
    if (isBrowser && !initialized.current) {
      initialized.current = true;
      const storedToken = storageService.getAccessToken();
      const urlToken = searchParams.get(AUTH_PARAMS.ACCESS_TOKEN);

      console.log('%c[DASHBOARD]', 'color:#34d399;font-weight:bold', 'init', {
        hasUrlToken: !!urlToken,
        hasStoredToken: !!storedToken,
      });

      if (urlToken) {
        console.log('%c[DASHBOARD]', 'color:#34d399;font-weight:bold', 'token from URL → storing and cleaning URL');
        storageService.setAccessToken(urlToken);
        setTimeout(() => setToken(urlToken), 0);
        router.replace(APP_ROUTES.HOME);
      } else if (storedToken) {
        console.log('%c[DASHBOARD]', 'color:#34d399;font-weight:bold', 'token from storage → authenticated');
        setTimeout(() => setToken(storedToken), 0);
      } else {
        const authUrl = new URL(`${ENV.HYDRA_PUBLIC_URL}/oauth2/auth`);
        authUrl.searchParams.append('client_id', OIDC_CONFIG.CLIENT_ID);
        authUrl.searchParams.append('response_type', 'code');
        authUrl.searchParams.append('scope', 'openid offline_access');
        authUrl.searchParams.append('redirect_uri', `${ENV.APP_URL}/api/auth/callback`);
        authUrl.searchParams.append('state', crypto.randomUUID());
        console.log('%c[DASHBOARD]', 'color:#34d399;font-weight:bold', 'no token → redirecting to Hydra OAuth2', authUrl.toString());
        window.location.href = authUrl.toString();
      }
    }
  }, [searchParams, router, token]);

  const handleLogout = async () => {
    storageService.clearAll();
    setToken(null);
    try {
      const { default: kratos } = await import('@/lib/ory/kratos');
      const { data } = await kratos.createBrowserLogoutFlow();
      window.location.href = data.logout_url;
    } catch (e) {
      router.push(APP_ROUTES.LOGIN);
    }
  };

  return (
    <div className="p-8 md:p-12 lg:p-16 space-y-12 bg-[radial-gradient(circle_at_bottom_left,_var(--tw-gradient-stops))] from-slate-50 via-white to-emerald-50/20 min-h-screen">
      <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <HeroTitle>System <span className="text-primary italic">Control</span></HeroTitle>
          <p className="text-on-surface-variant mt-2 font-medium text-lg max-w-2xl">
            Identity-aware terminal for Runtime Roasters cluster.
          </p>
        </div>
        
        {token ? (
          <div className="flex items-center gap-4 bg-white p-2 pr-6 rounded-full border border-outline-variant/20 shadow-xl shadow-primary/5">
            <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center text-primary border border-primary/20">
              <User className="w-6 h-6" />
            </div>
            <div>
              <p className="text-[10px] font-black uppercase tracking-widest text-slate-400 leading-none mb-1">Authenticated As</p>
              <p className="text-sm font-black font-headline italic uppercase tracking-tight text-on-surface">Admin Operator</p>
            </div>
            <button 
              onClick={handleLogout}
              className="ml-4 p-2 text-slate-400 hover:text-error transition-colors"
              title="Logout"
            >
              <LogOut className="w-5 h-5" />
            </button>
          </div>
        ) : (
          <button 
            onClick={() => {
              const authUrl = new URL(`${ENV.HYDRA_PUBLIC_URL}/oauth2/auth`);
              authUrl.searchParams.append('client_id', OIDC_CONFIG.CLIENT_ID);
              authUrl.searchParams.append('response_type', 'code');
              authUrl.searchParams.append('scope', 'openid offline_access');
              authUrl.searchParams.append('redirect_uri', `${ENV.APP_URL}/api/auth/callback`);
              authUrl.searchParams.append('state', crypto.randomUUID());
              window.location.href = authUrl.toString();
            }}
            className="px-8 py-4 bg-slate-900 text-white rounded-2xl font-black uppercase tracking-widest hover:scale-[0.98] transition-all flex items-center gap-3"
          >
            <Key className="w-5 h-5" />
            Access Terminal
          </button>
        )}
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        <div className="bg-white p-10 rounded-[3rem] border border-outline-variant/15 shadow-2xl relative overflow-hidden group">
          <div className="absolute top-0 left-0 w-2 h-full bg-emerald-500"></div>
          <MicroLabel className="mb-6">Security Context</MicroLabel>
          <h3 className="text-2xl font-black font-headline uppercase italic mb-8 text-on-surface">OIDC Verification</h3>
          
          <div className="space-y-6">
            <div className="p-6 rounded-3xl bg-slate-50 border border-slate-100">
              <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mb-4 italic">Active Access Token</p>
              <div className="font-mono text-[10px] break-all bg-white p-4 rounded-xl border border-slate-200 text-slate-600 max-h-32 overflow-y-auto">
                {token || 'No active token found. Please login.'}
              </div>
            </div>

            <button 
              onClick={() => refetch()}
              disabled={isLoading || !token}
              className="w-full py-5 bg-emerald-600 text-white rounded-2xl font-black uppercase tracking-widest hover:scale-[0.98] disabled:opacity-30 disabled:grayscale transition-all flex items-center justify-center gap-3 shadow-xl shadow-emerald-500/20"
            >
              <Terminal className="w-5 h-5" />
              {isLoading ? 'Decrypting Response...' : 'Call Secured API'}
            </button>

            <div className="bg-slate-900 rounded-3xl p-8 font-mono text-xs overflow-hidden border-4 border-slate-800 shadow-2xl">
              <div className="flex justify-between text-slate-500 mb-4 border-b border-slate-800 pb-4 uppercase tracking-tighter">
                <span>Response Output</span>
                <span className="text-emerald-400 font-black italic">{data ? '200_OK' : error ? '401_UNAUTHORIZED' : 'WAITING'}</span>
              </div>
              {data ? (
                <pre className="text-emerald-400 leading-relaxed">
                  {JSON.stringify(data, null, 2)}
                </pre>
              ) : error ? (
                <pre className="text-error italic leading-relaxed">
                  {JSON.stringify((error as { response?: { data: unknown } }).response?.data || "Unauthorized: Token invalid or missing.", null, 2)}
                </pre>
              ) : (
                <pre className="text-slate-600 italic">Waiting for secure handshake...</pre>
              )}
            </div>
          </div>
        </div>

        <div className="space-y-8">
          <div className="bg-slate-900 p-10 rounded-[3rem] border border-slate-800 shadow-2xl relative overflow-hidden text-white">
            <div className="absolute -top-24 -right-24 w-64 h-64 bg-primary/10 rounded-full blur-[80px]"></div>
            <MicroLabel className="mb-6 text-slate-400">System Logs</MicroLabel>
            <h3 className="text-2xl font-black font-headline uppercase italic mb-8">Service Mesh Health</h3>
            
            <div className="space-y-6">
              {[
                { label: 'Identity Proxy', status: 'Healthy', color: 'text-emerald-400' },
                { label: 'OIDC Provider (Hydra)', status: 'Active', color: 'text-emerald-400' },
                { label: 'Identity Mgr (Kratos)', status: 'Running', color: 'text-emerald-400' },
                { label: 'Internal JWKS', status: 'Protected', color: 'text-primary' },
              ].map((item, i) => (
                <div key={i} className="flex justify-between items-center py-3 border-b border-white/5">
                  <span className="text-xs font-bold uppercase tracking-widest text-slate-400">{item.label}</span>
                  <span className={`text-xs font-black uppercase italic ${item.color}`}>{item.status}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-surface-container-low p-10 rounded-[3rem] border border-outline-variant/10 shadow-sm">
            <MicroLabel className="mb-6">Observability</MicroLabel>
            <div className="flex gap-4 group text-tertiary">
              <div className="w-12 h-12 rounded-2xl bg-white flex items-center justify-center border border-tertiary/20 shrink-0 shadow-xl shadow-tertiary/5">
                <CheckCircle2 className="w-6 h-6" />
              </div>
              <div>
                <p className="text-sm font-black leading-tight mb-1 uppercase italic tracking-tight">Trace Propagation Active</p>
                <p className="text-xs font-medium text-on-surface-variant leading-relaxed">
                  Every request generates a unique Trace-ID. View full Gantt charts in <a href={ENV.JAEGER_URL} target="_blank" className="underline font-bold hover:text-primary transition-colors">Jaeger UI</a>.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function Dashboard() {
  return (
    <Suspense fallback={
      <div className="flex flex-col items-center justify-center min-h-screen bg-[radial-gradient(circle_at_bottom_left,_var(--tw-gradient-stops))] from-slate-50 via-white to-emerald-50/20 gap-6">
        <div className="w-16 h-16 border-4 border-primary/10 border-t-primary rounded-full animate-spin shadow-xl shadow-primary/5"></div>
        <div className="font-black uppercase italic tracking-[0.3em] text-primary/60 text-xs animate-pulse">
          Initializing Terminal UI...
        </div>
      </div>
    }>
      <DashboardContent />
    </Suspense>
  );
}
