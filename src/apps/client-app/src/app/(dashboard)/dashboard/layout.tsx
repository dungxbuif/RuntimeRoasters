'use client';

import Sidebar from "@/components/common/Sidebar";
import { AuthGuard, useAuth } from "@/lib/auth";
import { LogOut, User } from "lucide-react";

export default function ManagementLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, logout } = useAuth();

  return (
    <AuthGuard>
      <div className="flex min-h-screen bg-[#f7f9fb]">
        <Sidebar />
        <div className="flex-1 flex flex-col md:ml-64">
          <header className="h-16 flex justify-between items-center px-8 border-b border-outline-variant/10 bg-white/70 backdrop-blur-xl sticky top-0 z-30" data-e2e="dashboard-header">
            <div className="flex items-center gap-2">
              <span className="text-[10px] font-black uppercase tracking-[0.3em] text-slate-400 italic">Cluster Node</span>
              <span className="text-[10px] font-mono text-primary font-bold">/DASHBOARD</span>
            </div>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-3 bg-slate-100/50 p-1 pr-4 rounded-full border border-outline-variant/10">
                <div className="w-8 h-8 rounded-full bg-white flex items-center justify-center text-primary shadow-sm border border-outline-variant/10">
                  <User className="w-4 h-4" />
                </div>
                <div className="text-left hidden sm:block">
                  <p className="text-[9px] font-bold text-on-surface uppercase tracking-tight truncate max-w-[100px]">
                    {user?.email.split('@')[0]}
                  </p>
                </div>
                <button 
                  onClick={() => logout()}
                  className="ml-2 p-1 text-slate-400 hover:text-error transition-colors"
                  title="Logout"
                  data-e2e="logout-btn"
                >
                  <LogOut className="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </header>
          <main className="flex-1 p-8 lg:p-12 relative overflow-y-auto">
            {children}
          </main>
        </div>
      </div>
    </AuthGuard>
  );
}
