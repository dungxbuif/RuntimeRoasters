'use client';

import React from 'react';

interface LoginCardProps {
  children: React.ReactNode;
  title: string;
  subtitle: string;
}

export const LoginCard: React.FC<LoginCardProps> = ({ children, title, subtitle }) => {
  return (
      <div className="bg-white/80 backdrop-blur-xl shadow-[0_32px_64px_-16px_rgba(0,0,0,0.1)] rounded-[2.5rem] overflow-hidden relative group">
        <div className="p-10 md:p-12">
          <header className="mb-10 text-center">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-primary/5 rounded-2xl mb-6 border border-primary/10 shadow-inner">
              <span className="material-symbols-outlined text-primary text-3xl font-black">coffee</span>
            </div>
            <h1 className="text-3xl font-black font-headline uppercase italic tracking-tighter text-on-surface mb-2">
              {title}
            </h1>
            <p className="text-sm font-medium text-on-surface-variant/70 uppercase tracking-widest italic">
              {subtitle}
            </p>
          </header>

          <main>
            {children}
          </main>

        </div>

        {/* Decorative corner element */}
        <div className="absolute -bottom-10 -right-10 w-32 h-32 bg-primary/5 rounded-full blur-3xl pointer-events-none group-hover:bg-primary/10 transition-colors"></div>
      </div>
  );
};
