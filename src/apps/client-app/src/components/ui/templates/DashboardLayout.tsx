import React, { ReactNode } from 'react';
import { LucideIcon } from 'lucide-react';

interface DashboardLayoutProps {
  title: string;
  subtitle: string;
  icon?: LucideIcon;
  actions?: ReactNode;
  children: ReactNode;
}

export function DashboardLayout({ title, subtitle, icon: Icon, actions, children }: DashboardLayoutProps) {
  return (
    <div className="h-full flex flex-col gap-6 font-body">
      {/* Header */}
      <div className="flex justify-between items-center bg-slate-100 p-6 rounded-[2rem] border border-slate-200">
        <div>
          <h1 className="text-3xl font-black uppercase tracking-tighter text-slate-900 italic flex items-center gap-3">
            {Icon && <Icon className="w-8 h-8 text-primary" />}
            {title}
          </h1>
          <p className="text-[12px] font-bold text-slate-500 uppercase tracking-widest mt-1">
            {subtitle}
          </p>
        </div>
        {actions && (
          <div className="flex items-center gap-3">
            {actions}
          </div>
        )}
      </div>

      {/* Main Content */}
      <div className="flex-1 flex gap-6 overflow-hidden min-h-0">
        {children}
      </div>
    </div>
  );
}
