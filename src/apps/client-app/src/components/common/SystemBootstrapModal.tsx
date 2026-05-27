'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/lib/auth';
import { systemService } from '@/services/system.service';
import { Database, Loader2, CheckCircle2, AlertCircle } from 'lucide-react';

export default function SystemBootstrapModal() {
  const { user } = useAuth();
  const [show, setShow] = useState(false);
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [message, setMessage] = useState('');

  const checkStatus = async () => {
    try {
      const result = await systemService.checkAllStatus();
      console.debug('[SystemBootstrapModal] Status check result:', result);
      if (!result.allSeeded) {
        console.debug('[SystemBootstrapModal] System not seeded, showing modal');
        setShow(true);
      } else {
        console.debug('[SystemBootstrapModal] System already seeded, keeping modal hidden');
      }
    } catch (err) {
      console.error('[SystemBootstrapModal] Failed to check status:', err);
    }
  };

  useEffect(() => {
    console.debug('[SystemBootstrapModal] User changed:', user?.email, 'Role:', user?.role);
    
    // Allow force trigger via URL for testing
    const params = new URLSearchParams(window.location.search);
    if (params.get('bootstrap') === 'true' && user?.role === 'ADMIN') {
      window.setTimeout(() => setShow(true), 0);
      return;
    }

    if (user?.role === 'ADMIN') {
      console.debug('[SystemBootstrapModal] User is ADMIN, checking system status...');
      window.setTimeout(() => {
        void checkStatus();
      }, 0);
    } else {
      window.setTimeout(() => setShow(false), 0);
    }
  }, [user]);

  const handleSeed = async () => {
    setLoading(true);
    setStatus('loading');
    try {
      await systemService.seedAll();
      setStatus('success');
      setMessage('System data initialized successfully!');
      setTimeout(() => setShow(false), 2000);
    } catch {
      setStatus('error');
      setMessage('Failed to initialize system data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (!show) return null;

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div className="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden border border-outline-variant/20 animate-in fade-in zoom-in duration-300">
        <div className="p-8">
          <div className="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center text-primary mb-6 mx-auto">
            <Database className="w-8 h-8" />
          </div>
          
          <h2 className="text-2xl font-black text-center text-on-surface uppercase tracking-tight mb-2">
            System Bootstrap
          </h2>
          <p className="text-slate-500 text-center text-sm font-medium mb-8">
            The system detects that the initial data has not been seeded. Would you like to initialize the core data now?
          </p>

          {status === 'idle' && (
            <button
              onClick={handleSeed}
              disabled={loading}
              className="w-full py-4 bg-primary text-white rounded-xl font-black uppercase tracking-widest hover:bg-primary-fixed transition-all flex items-center justify-center gap-3 shadow-lg shadow-primary/20"
            >
              {loading ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Initialize Data'}
            </button>
          )}

          {status === 'loading' && (
            <div className="flex flex-col items-center gap-4 py-4">
              <Loader2 className="w-10 h-10 animate-spin text-primary" />
              <p className="text-xs font-black uppercase tracking-widest text-primary animate-pulse">Seeding microservices...</p>
            </div>
          )}

          {status === 'success' && (
            <div className="flex flex-col items-center gap-4 py-4 text-tertiary">
              <CheckCircle2 className="w-12 h-12" />
              <p className="text-sm font-bold">{message}</p>
            </div>
          )}

          {status === 'error' && (
            <div className="space-y-4">
              <div className="flex flex-col items-center gap-4 text-error">
                <AlertCircle className="w-12 h-12" />
                <p className="text-sm font-bold">{message}</p>
              </div>
              <button
                onClick={() => setStatus('idle')}
                className="w-full py-3 border-2 border-slate-200 text-slate-500 rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-slate-50 transition-all"
              >
                Try Again
              </button>
            </div>
          )}
        </div>
        
        <div className="px-8 py-4 bg-slate-50 border-t border-outline-variant/10 flex justify-center">
          <p className="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] italic">
            ADMIN PRIVILEGES REQUIRED
          </p>
        </div>
      </div>
    </div>
  );
}
