import React from "react";
import { cn } from "@/lib/utils/cn";
import { MicroLabel } from "../atoms/Typography";

interface MetricCardProps {
  label: string;
  value: string;
  trend?: string;
  trendDirection?: 'up' | 'down' | 'stable';
  icon?: string;
  className?: string;
}

export const MetricCard = ({ label, value, trend, trendDirection = 'stable', icon, className }: MetricCardProps) => {
  const trendColors = {
    up: "text-emerald-500",
    down: "text-error",
    stable: "text-on-surface-variant opacity-60"
  };

  return (
    <div className={cn(
      "bg-surface-container-lowest p-10 rounded-[2.5rem] border border-outline-variant/15 shadow-[0_20px_50px_-12px_rgba(0,0,0,0.05)] hover:shadow-xl transition-all duration-500 relative overflow-hidden group",
      className
    )}>
      {icon && (
        <div className="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-opacity">
          <span className="material-symbols-outlined !text-7xl">{icon}</span>
        </div>
      )}
      <MicroLabel className="mb-6 italic">{label}</MicroLabel>
      <div className="flex items-baseline gap-4 relative z-10">
        <p className="text-6xl font-headline font-black italic text-on-surface tracking-tighter">{value}</p>
        {trend && (
          <span className={cn("text-[10px] font-black uppercase italic", trendColors[trendDirection])}>
            {trend}
          </span>
        )}
      </div>
    </div>
  );
};
