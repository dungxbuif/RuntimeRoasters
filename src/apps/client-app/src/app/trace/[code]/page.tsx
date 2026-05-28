'use client';

import React from 'react';
import { useParams } from 'next/navigation';
import { Coffee, Sprout, Truck, Factory, Package, ShoppingCart, CheckCircle2, MapPin } from 'lucide-react';

interface JourneyStep {
  icon: React.ElementType;
  label: string;
  detail: string;
  date: string;
  extra?: string;
  color: string;
}

// Demo trace data
const DEMO_JOURNEY: JourneyStep[] = [
  { icon: Sprout, label: 'Harvested', detail: "K'Ho Coffee Farm, Lạc Dương, Đà Lạt", date: 'May 10, 2026', extra: '500kg Arabica Heirloom', color: 'bg-green-500' },
  { icon: Truck, label: 'Picked Up', detail: 'Driver pickup from farm', date: 'May 11, 2026', color: 'bg-blue-500' },
  { icon: Factory, label: 'Warehouse Received', detail: 'Warehouse Hanoi (HN-001)', date: 'May 11, 2026', color: 'bg-blue-500' },
  { icon: Coffee, label: 'Processed', detail: 'Roasted to Medium-Dark', date: 'May 12, 2026', extra: 'Weight: 420kg (16% loss)', color: 'bg-amber-500' },
  { icon: Package, label: 'Stocked', detail: 'SKU: ARABICA-DALAT-M', date: 'May 12, 2026', color: 'bg-blue-500' },
  { icon: ShoppingCart, label: 'Ordered', detail: 'Store Hoàn Kiếm, Hanoi', date: 'May 13, 2026', color: 'bg-blue-500' },
  { icon: Truck, label: 'Delivered', detail: 'Driver delivery completed', date: 'May 13, 2026', color: 'bg-green-500' },
  { icon: CheckCircle2, label: 'Served', detail: 'Ready for customer', date: 'May 14, 2026', color: 'bg-green-600' },
];

export default function PublicTracePage() {
  const params = useParams();
  const traceCode = (params?.code as string) || 'RR-S-CD-20260510-0001';

  return (
    <div className="min-h-screen bg-[#fafbfc] font-body">
      {/* Header */}
      <header className="bg-white border-b border-slate-100 px-8 py-5">
        <div className="max-w-3xl mx-auto flex justify-between items-center">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-primary rounded-lg flex items-center justify-center text-white shadow-md">
              <Coffee className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-sm font-black text-slate-900 uppercase tracking-tighter italic">Runtime Roasters</h2>
              <p className="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Product Journey</p>
            </div>
          </div>
          <div className="text-right">
            <div className="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Trace Code</div>
            <div className="text-xs font-black font-mono text-primary">{traceCode}</div>
          </div>
        </div>
      </header>

      {/* Journey Timeline */}
      <main className="max-w-3xl mx-auto px-8 py-12">
        <div className="text-center mb-12">
          <h1 className="text-3xl font-black uppercase tracking-tighter text-slate-900 italic">
            Farm to <span className="text-primary">Cup</span> Journey
          </h1>
          <p className="text-sm text-slate-500 mt-2 font-medium">
            Complete traceability for every coffee product
          </p>
        </div>

        {/* Timeline */}
        <div className="relative">
          {/* Vertical line */}
          <div className="absolute left-6 top-4 bottom-4 w-0.5 bg-slate-200" />

          <div className="space-y-0">
            {DEMO_JOURNEY.map((step, i) => {
              const isLast = i === DEMO_JOURNEY.length - 1;
              return (
                <div key={step.label} className="relative flex gap-6 pb-8 group">
                  {/* Node */}
                  <div className="relative z-10 flex-shrink-0">
                    <div className={`w-12 h-12 rounded-2xl ${step.color} flex items-center justify-center text-white shadow-lg group-hover:scale-110 transition-transform`}>
                      <step.icon className="w-5 h-5" />
                    </div>
                  </div>

                  {/* Content */}
                  <div className={`flex-1 bg-white rounded-2xl p-5 border border-slate-100 shadow-sm group-hover:shadow-md group-hover:border-slate-200 transition-all ${isLast ? 'border-green-200 bg-green-50/30' : ''}`}>
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="text-sm font-black text-slate-900 uppercase tracking-tight">{step.label}</h3>
                        <p className="text-xs text-slate-600 mt-0.5 font-medium">{step.detail}</p>
                        {step.extra && (
                          <p className="text-[10px] text-slate-400 font-bold mt-1 italic">{step.extra}</p>
                        )}
                      </div>
                      <span className="text-[10px] font-bold text-slate-400 uppercase whitespace-nowrap ml-4">{step.date}</span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* QR Section */}
        <div className="mt-12 bg-white rounded-3xl border border-slate-200 p-8 text-center shadow-sm">
          <div className="w-32 h-32 bg-slate-100 rounded-2xl mx-auto mb-4 flex items-center justify-center border-2 border-dashed border-slate-300">
            <div className="text-center">
              <MapPin className="w-8 h-8 text-slate-400 mx-auto mb-1" />
              <span className="text-[8px] font-black text-slate-400 uppercase tracking-widest">QR Code</span>
            </div>
          </div>
          <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Scan to verify authenticity</p>
        </div>

        {/* Service Architecture Badge */}
        <div className="mt-8 flex justify-center items-center gap-8 py-6 opacity-60">
          <span className="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Powered by</span>
          {[
            { icon: Sprout, label: 'Farm' },
            { icon: Factory, label: 'Warehouse' },
            { icon: Truck, label: 'Logistics' },
            { icon: ShoppingCart, label: 'Retail' },
          ].map(s => (
            <div key={s.label} className="flex flex-col items-center gap-1">
              <s.icon className="w-5 h-5 text-primary" />
              <span className="text-[8px] font-black text-slate-500 uppercase">{s.label}</span>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
