import React, { ReactNode } from 'react';

interface EntityCardProps {
  title: string;
  subtitle?: string;
  badge?: ReactNode;
  icon?: ReactNode;
  children: ReactNode;
  className?: string;
}

export function EntityCard({ title, subtitle, badge, icon, children, className = '' }: EntityCardProps) {
  return (
    <div className={`border border-slate-200 rounded-2xl p-4 shadow-sm hover:border-slate-300 transition-colors bg-white ${className}`}>
      <div className="flex justify-between items-start mb-4">
        <div>
          <h4 className="font-bold text-slate-900 text-sm flex items-center gap-2">
            {icon && <span className="text-slate-400">{icon}</span>}
            {title}
          </h4>
          {subtitle && <p className="text-xs text-slate-500 mt-1">{subtitle}</p>}
        </div>
        {badge && <div>{badge}</div>}
      </div>
      <div>{children}</div>
    </div>
  );
}
