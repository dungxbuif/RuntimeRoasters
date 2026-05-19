'use client';

import React from 'react';
import { TraceDocument } from '@/types/trace';

export interface TraceTimelineProps {
  document: TraceDocument;
}

const TOPIC_CONFIG: Record<string, { icon: string, color: string, label: string }> = {
  'farm.harvest.events': { 
    icon: 'agriculture', 
    color: 'primary', 
    label: 'Harvest Origin' 
  },
  'process.batch.completed': { 
    icon: 'sync', 
    color: 'emerald-500', 
    label: 'Processing' 
  },
  'warehouse.stock.reserved': { 
    icon: 'inventory_2', 
    color: 'tertiary', 
    label: 'Stock Allocated' 
  },
  'retail.order.created': { 
    icon: 'payments', 
    color: 'secondary', 
    label: 'Order Placed' 
  },
  'logistics.shipment.assigned': { 
    icon: 'local_shipping', 
    color: 'amber-500', 
    label: 'Shipment Assigned' 
  },
  'logistics.shipment.delivered': { 
    icon: 'check_circle', 
    color: 'success', 
    label: 'Delivered' 
  }
};

export default function TraceTimeline({ document }: TraceTimelineProps) {
  const sortedEvents = [...document.events].sort((a, b) => 
    new Date(a.occurred_at).getTime() - new Date(b.occurred_at).getTime()
  );

  return (
    <div className="max-w-5xl mx-auto relative pb-32">
      <div className="absolute left-1/2 top-0 bottom-0 w-2 -translate-x-1/2 z-0 pointer-events-none hidden md:block opacity-20 bg-primary/20 rounded-full"></div>

      <div className="space-y-32 relative z-10">
        {sortedEvents.map((event, index) => {
          const config = TOPIC_CONFIG[event.topic] || { 
            icon: 'info', 
            color: 'slate-400', 
            label: 'System Event' 
          };
          const isEven = index % 2 === 0;
          let payload: any = {};
          try {
            payload = JSON.parse(event.payload);
          } catch (e) {
            console.error("Failed to parse event payload", e);
          }

          return (
            <div key={event.id} className={`flex flex-col md:flex-row items-center justify-between group ${isEven ? '' : 'md:flex-row-reverse'}`}>
              <div className={`w-full md:w-5/12 flex ${isEven ? 'justify-end md:pr-12' : 'justify-start md:pl-12'} relative`}>
                <div className={`bg-white p-8 rounded-[2rem] shadow-xl border border-slate-200 w-full max-w-sm hover:border-primary transition-all duration-500 relative overflow-hidden`}>
                  <div className="flex items-center space-x-4 mb-6">
                    <div className="w-14 h-14 rounded-2xl bg-slate-100 flex items-center justify-center border-2 border-slate-200">
                      <span className="material-symbols-outlined text-3xl">{config.icon}</span>
                    </div>
                    <div>
                      <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest opacity-50">{config.label}</div>
                      <h3 className="font-headline text-xl font-black text-slate-800 uppercase italic truncate">
                        {event.topic.split('.').slice(-1)[0]}
                      </h3>
                    </div>
                  </div>
                  
                  <div className="space-y-3 font-body text-xs text-slate-600 font-bold uppercase tracking-widest italic opacity-80">
                    {Object.entries(payload).map(([key, value]) => (
                      typeof value !== 'object' && (
                        <div key={key} className="flex justify-between border-b border-slate-100 pb-2">
                          <span className="opacity-60">{key.replace(/_/g, ' ')}</span>
                          <span className="text-slate-900 truncate ml-4">{String(value)}</span>
                        </div>
                      )
                    ))}
                    <div className="flex justify-between pt-1 opacity-40">
                      <span>Time</span>
                      <span>{new Date(event.occurred_at).toLocaleString()}</span>
                    </div>
                  </div>
                </div>
              </div>
              
              <div className={`w-10 h-10 rounded-full bg-white border-[6px] border-primary z-10 hidden md:block shadow-2xl ${index === sortedEvents.length - 1 ? 'animate-pulse' : ''}`}></div>
              
              <div className={`w-full md:w-5/12 ${isEven ? 'md:pl-12' : 'md:pr-12'} mt-6 md:mt-0 ${!isEven && 'flex justify-end'}`}>
                 <div className={`font-marker text-3xl text-slate-700 ${isEven ? '-rotate-2' : 'rotate-1 text-right'} leading-tight`}>
                    <span className={`bg-slate-100 px-2 py-1 inline-block mb-2`}>{config.label}:</span><br/>
                    {event.topic}
                 </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
