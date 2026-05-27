'use client';

import React from 'react';

export interface NotificationItem {
  id: string;
  icon: string;
  message: string;
  time: string;
  color?: 'blue' | 'green' | 'amber' | 'red';
}

interface NotificationFeedProps {
  items: NotificationItem[];
  title?: string;
}

const COLOR_MAP = {
  blue: 'bg-primary',
  green: 'bg-green-500',
  amber: 'bg-amber-500',
  red: 'bg-red-500',
};

export default function NotificationFeed({ items, title = 'Notifications' }: NotificationFeedProps) {
  return (
    <div className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
      <div className="px-6 py-4 border-b border-slate-100 bg-slate-50/50">
        <h3 className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 flex items-center gap-2">
          <span className="material-symbols-outlined !text-sm text-primary">notifications</span>
          {title}
        </h3>
      </div>
      <div className="divide-y divide-slate-50 max-h-[280px] overflow-y-auto">
        {items.length === 0 ? (
          <div className="px-6 py-8 text-center">
            <p className="text-[10px] font-black text-slate-300 uppercase tracking-widest italic">No notifications</p>
          </div>
        ) : (
          items.map((item) => (
            <div key={item.id} className="flex items-start gap-3 px-6 py-3 hover:bg-slate-50/50 transition-colors">
              <div className={`w-2 h-2 rounded-full mt-1.5 flex-shrink-0 ${COLOR_MAP[item.color || 'blue']}`} />
              <div className="min-w-0 flex-1">
                <p className="text-[11px] font-bold text-slate-700 leading-snug">{item.message}</p>
                <p className="text-[9px] font-bold text-slate-400 mt-0.5 uppercase">{item.time}</p>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
