import React from "react";
import { cn } from "@/lib/utils/cn";

export const PulseDot = ({ color = 'emerald' }: { color?: 'emerald' | 'error' | 'amber' | 'primary' }) => {
  const colors = {
    emerald: "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)]",
    error: "bg-error shadow-[0_0_8px_#ba1a1a]",
    amber: "bg-amber-500 shadow-[0_0_8px_rgba(245,158,11,0.6)]",
    primary: "bg-primary shadow-[0_0_8px_rgba(0,74,198,0.6)]"
  };

  return <div className={cn("w-2 h-2 rounded-full animate-pulse", colors[color])}></div>;
};

interface StatusPillProps {
  children: React.ReactNode;
  variant?: 'success' | 'warning' | 'error' | 'neutral';
  className?: string;
}

export const StatusPill = ({ children, variant = 'neutral', className }: StatusPillProps) => {
  const variants = {
    success: "bg-tertiary-fixed text-on-tertiary-fixed border-tertiary/20",
    warning: "bg-amber-100 text-amber-800 border-amber-200",
    error: "bg-error-container text-on-error-container border-error/10",
    neutral: "bg-surface-container-high text-on-surface-variant border-outline-variant/10"
  };

  return (
    <span className={cn(
      "px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border shadow-sm inline-flex items-center gap-1.5",
      variants[variant],
      className
    )}>
      {children}
    </span>
  );
};

export const Tag = ({ children, className }: { children: React.ReactNode, className?: string }) => (
  <span className={cn(
    "font-mono text-[9px] bg-surface-container-high text-on-surface-variant px-2 py-0.5 border border-outline-variant/20 font-black tracking-tighter italic rounded",
    className
  )}>
    {children}
  </span>
);
