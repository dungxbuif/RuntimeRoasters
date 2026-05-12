'use client';

import React from 'react';

interface HandDrawnNodeProps {
  label: string;
  subtitle?: string;
  icon: string;
  category: 'primary' | 'secondary' | 'error' | 'infra';
  isActive?: boolean;
  className?: string;
}

export const HandDrawnNode = ({ label, subtitle, icon, category, isActive, className }: HandDrawnNodeProps) => {
  const styles = {
    primary: 'bg-primary-fixed border-primary text-on-primary-fixed',
    secondary: 'bg-secondary-container border-secondary text-on-secondary-container',
    error: 'bg-error-container border-error text-on-error-container opacity-80 border-dashed',
    infra: 'bg-surface-container-low border-outline text-on-surface-variant'
  };

  return (
    <div 
      className={`w-32 h-24 border-2 rounded-lg flex flex-col items-center justify-center z-20 transition-all duration-500 ${styles[category]} ${isActive ? 'scale-110 shadow-xl' : ''} ${className}`}
      style={{ filter: 'url(#rough-edge)' }}
    >
      <span className="material-symbols-outlined mb-1">{icon}</span>
      <span className="font-body text-[11px] font-bold text-center px-2 leading-tight uppercase tracking-tighter">{label}</span>
      {subtitle && <span className="text-[8px] opacity-60 uppercase font-black tracking-widest mt-0.5">{subtitle}</span>}
      
      {category === 'error' && (
        <span className="text-[8px] text-error font-black uppercase mt-1 bg-error/10 px-2 py-0.5 rounded leading-none">Failed</span>
      )}
    </div>
  );
};
