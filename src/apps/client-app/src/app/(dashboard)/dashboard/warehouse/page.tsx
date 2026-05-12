import React from 'react';

export default function WarehouseStockLogicPage() {
  return (
    <div className="min-h-full p-8 font-body selection:bg-primary-fixed selection:text-on-primary-fixed relative overflow-hidden bg-surface">
      {/* Header Section */}
      <header className="mb-10 flex flex-col md:flex-row justify-between items-start md:items-end gap-6 relative z-10">
        <div>
          <div className="flex items-center gap-3 text-primary font-marker text-2xl mb-2 italic">
            <span>Runtime Roasters</span>
            <span className="material-symbols-outlined !text-sm">trending_flat</span>
            <span className="text-on-surface-variant opacity-60">Warehouse Service</span>
          </div>
          <h1 className="text-5xl font-headline font-black tracking-tighter text-on-surface uppercase italic">Inventory Command</h1>
        </div>
        <div className="flex gap-4">
          <button className="bg-on-background text-white px-6 py-3 rounded-xl font-black text-xs uppercase tracking-widest shadow-2xl hover:scale-95 transition-all flex items-center gap-3 group border border-slate-700">
            <span className="material-symbols-outlined !text-sm group-hover:rotate-90 transition-transform">terminal</span>
            Access DB Shell
          </button>
        </div>
      </header>

      {/* Bento Grid: Regional Warehouses */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12 relative z-10">
        <div className="col-span-1 md:col-span-2 bg-surface-container-lowest p-8 rounded-[2rem] border border-outline-variant/15 shadow-sm relative overflow-hidden group">
          <div className="absolute top-0 left-0 w-2 h-full bg-primary shadow-[0_0_15px_rgba(0,74,198,0.3)]"></div>
          <div className="flex justify-between items-start mb-10 relative z-10">
            <div>
              <h3 className="font-headline font-black text-2xl mb-1 uppercase tracking-tight text-on-surface">Stock Levels: North America (US-EAST)</h3>
              <p className="text-on-surface-variant font-medium text-sm">Real-time inventory distribution across SKU clusters.</p>
            </div>
            <span className="bg-tertiary-fixed text-on-tertiary-fixed px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-tertiary/20 shadow-sm animate-pulse">
              Live_Sync: Active
            </span>
          </div>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6 relative z-10">
            {[
              { label: 'Arabica Dark', val: '12.4k', trend: '↑ 12%', color: 'bg-primary', w: '75%', trendColor: 'text-tertiary' },
              { label: 'Robusta Blend', val: '8.2k', trend: '↓ 4%', color: 'bg-primary', w: '45%', trendColor: 'text-error' },
              { label: 'Decaf Gold', val: '3.1k', trend: '--', color: 'bg-outline-variant', w: '20%', trendColor: 'text-on-surface-variant' },
              { label: 'Cold Brew Kit', val: '15.9k', trend: '↑ 24%', color: 'bg-tertiary', w: '90%', trendColor: 'text-tertiary' }
            ].map((item, i) => (
              <div key={i} className="bg-surface-container-low p-5 rounded-2xl border border-outline-variant/5 group/card hover:border-primary/20 transition-colors">
                <span className="text-[10px] font-black text-on-surface-variant uppercase block mb-3 tracking-widest opacity-60">{item.label}</span>
                <div className="flex items-end gap-2 mb-4">
                  <span className="text-3xl font-headline font-black text-on-surface tabular-nums">{item.val}</span>
                  <span className={`${item.trendColor} text-[10px] font-black mb-1 italic`}>{item.trend}</span>
                </div>
                <div className="w-full bg-surface-container-highest h-1.5 rounded-full overflow-hidden">
                  <div className={`${item.color} h-full group-hover/card:scale-x-105 transition-transform origin-left`} style={{ width: item.w }}></div>
                </div>
              </div>
            ))}
          </div>
          
          <div className="mt-10 flex justify-end">
            <div className="font-marker text-primary text-2xl rotate-[-2deg] flex items-center gap-3 opacity-70">
              <span className="material-symbols-outlined scale-x-[-1] animate-bounce">subdirectory_arrow_right</span>
              Replenishment expected 14:00 UTC
            </div>
          </div>
        </div>

        <div className="primary-gradient p-8 rounded-[2.5rem] text-on-primary relative flex flex-col justify-between shadow-2xl overflow-hidden group">
          <div className="absolute top-0 right-0 -mr-16 -mt-16 w-64 h-64 bg-white/10 rounded-full blur-3xl group-hover:bg-white/20 transition-all"></div>
          <div className="relative z-10">
            <span className="material-symbols-outlined !text-6xl mb-6 drop-shadow-xl">database</span>
            <h3 className="font-headline font-black text-3xl mb-4 uppercase tracking-tighter italic">Persistence Layer</h3>
            <div className="space-y-4 font-mono text-[10px] font-bold tracking-widest opacity-80 uppercase leading-relaxed">
              <p className="flex justify-between border-b border-white/20 pb-1"><span>Backing Store:</span> <span>PostgreSQL 16</span></p>
              <p className="flex justify-between border-b border-white/20 pb-1"><span>Isolation:</span> <span>Repeatable Read</span></p>
              <p className="flex justify-between border-b border-white/20 pb-1"><span>Status:</span> <span className="text-emerald-300">Strongly Consistent</span></p>
            </div>
          </div>
          <div className="mt-10 relative z-10">
            <div className="flex justify-between items-center text-[10px] font-black mb-2 tracking-[0.2em]">
              <span>ACID COMPLIANCE</span>
              <span className="italic">100%</span>
            </div>
            <div className="w-full bg-white/20 h-3 rounded-full overflow-hidden border border-white/10 p-0.5">
              <div className="bg-emerald-400 h-full w-full rounded-full shadow-[0_0_15px_rgba(52,211,153,0.8)]"></div>
            </div>
          </div>
        </div>
      </div>

      {/* Saga Visualizer Section */}
      <section className="mb-12 relative z-10">
        <div className="flex flex-col md:flex-row items-end justify-between mb-8 gap-6">
          <div>
            <h2 className="text-4xl font-headline font-black text-on-surface uppercase tracking-tighter italic">Saga Flow: Stock Reservation</h2>
            <p className="text-on-surface-variant max-w-2xl font-medium mt-2">Visualizing the Choreographed Saga flow between Retail and Warehouse during the checkout lifecycle.</p>
          </div>
          <button className="bg-surface-container-high text-on-surface-variant px-8 py-3 rounded-full font-black text-[10px] uppercase tracking-widest hover:bg-white transition-all border border-outline-variant/30 shadow-sm">
            View Live Log Events
          </button>
        </div>

        <div className="bg-surface-container-lowest rounded-[3rem] p-16 relative border border-outline-variant/15 shadow-2xl overflow-hidden">
          <div className="grid grid-cols-3 gap-20 relative z-10">
            {/* Column 1: Retail Service */}
            <div className="flex flex-col items-center">
              <div className="w-24 h-24 bg-primary-fixed rounded-[2rem] flex items-center justify-center mb-6 shadow-xl border-2 border-primary/20 group hover:scale-110 transition-transform cursor-pointer">
                <span className="material-symbols-outlined !text-4xl text-primary">shopping_cart</span>
              </div>
              <span className="font-black text-xs uppercase tracking-widest text-on-surface italic">Retail Service</span>
              <div className="w-1.5 h-80 bg-surface-container-high flex flex-col items-center mt-6 relative rounded-full">
                <div className="w-8 h-40 bg-primary/10 border-x-4 border-primary/30 mt-6 rounded-sm relative overflow-hidden">
                   <div className="absolute inset-0 bg-primary/5 animate-pulse"></div>
                </div>
              </div>
            </div>

            {/* Column 2: Event Bus */}
            <div className="flex flex-col items-center">
              <div className="w-24 h-24 bg-on-background rounded-full flex items-center justify-center mb-6 border-4 border-dashed border-outline-variant/50 shadow-2xl relative group hover:rotate-90 transition-transform duration-1000">
                <span className="material-symbols-outlined !text-4xl text-white">hub</span>
              </div>
              <span className="font-black text-xs uppercase tracking-widest text-on-surface-variant italic opacity-60">Kafka Broker</span>
              <div className="w-1.5 h-80 bg-surface-container-high relative mt-6 rounded-full">
                {/* Message Path 1 */}
                <div className="absolute top-16 left-[-160px] w-[320px] h-1 bg-primary shadow-[0_0_15px_rgba(0,74,198,0.4)] animate-dash">
                  <div className="absolute right-0 top-[-6px] w-3 h-3 border-t-4 border-r-4 border-primary rotate-45"></div>
                  <span className="absolute top-[-35px] left-1/2 -translate-x-1/2 text-[9px] font-black bg-on-background text-white px-4 py-1.5 rounded-full uppercase tracking-tighter whitespace-nowrap shadow-lg">RESERVE_STOCK_CMD</span>
                </div>
                {/* Message Path 2 */}
                <div className="absolute top-48 left-[-160px] w-[320px] h-1 bg-tertiary shadow-[0_0_15px_rgba(0,98,66,0.4)]">
                  <div className="absolute left-0 top-[-6px] w-3 h-3 border-b-4 border-l-4 border-tertiary rotate-45"></div>
                  <span className="absolute bottom-[-35px] left-1/2 -translate-x-1/2 text-[9px] font-black bg-tertiary text-white px-4 py-1.5 rounded-full uppercase tracking-tighter whitespace-nowrap shadow-lg">STOCK_RESERVED_EVENT</span>
                </div>
              </div>
            </div>

            {/* Column 3: Warehouse Service */}
            <div className="flex flex-col items-center">
              <div className="w-24 h-24 bg-tertiary-container rounded-[2rem] flex items-center justify-center mb-6 shadow-xl border-2 border-tertiary/20 group hover:scale-110 transition-transform cursor-pointer">
                <span className="material-symbols-outlined !text-4xl text-tertiary">inventory_2</span>
              </div>
              <span className="font-black text-xs uppercase tracking-widest text-on-surface italic">Warehouse API</span>
              <div className="w-1.5 h-80 bg-surface-container-high flex flex-col items-center mt-6 rounded-full">
                <div className="w-8 h-32 bg-tertiary/10 border-x-4 border-tertiary/30 mt-24 rounded-sm relative overflow-hidden">
                  <div className="absolute inset-0 bg-tertiary/5 animate-pulse"></div>
                </div>
              </div>
            </div>
          </div>

          {/* Floating Callout */}
          <div className="absolute top-1/2 right-12 translate-y-[-50%] max-w-sm p-10 bg-white/80 backdrop-blur-2xl rounded-[2.5rem] border-2 border-white shadow-[0_30px_60px_-12px_rgba(0,0,0,0.15)] transform hover:-translate-x-2 transition-transform">
            <h4 className="font-headline font-black text-primary mb-6 flex items-center gap-3 uppercase tracking-tighter italic text-xl">
              <span className="material-symbols-outlined !text-2xl animate-spin-slow">autorenew</span>
              Atomic Local Tx
            </h4>
            <p className="text-sm text-on-surface-variant font-bold leading-relaxed tracking-tight">
              Warehouse service executes a local transaction to decrement <code className="bg-surface-container px-2 py-0.5 rounded text-primary italic font-black text-xs">avail_qty</code> and increment <code className="bg-surface-container px-2 py-0.5 rounded text-tertiary italic font-black text-xs">resv_qty</code>. 
              <br/><br/>
              Consistent read models are generated via <strong className="text-on-surface italic underline decoration-tertiary decoration-2">Debezium connectors</strong> polling Postgres binlogs.
            </p>
          </div>

          {/* Marker Note */}
          <div className="absolute top-[75%] left-[10%] font-marker text-on-surface-variant transform -rotate-3 hover:text-primary transition-colors cursor-help group opacity-40 hover:opacity-100">
            <span className="flex items-center gap-2 text-3xl tracking-tighter">
              Validate Postgres isolation levels
              <span className="material-symbols-outlined animate-bounce">trending_flat</span>
            </span>
          </div>
          
          <div className="absolute inset-0 pointer-events-none opacity-[0.02]" style={{ backgroundImage: 'radial-gradient(#000 1px, transparent 1px)', backgroundSize: '32px 32px' }}></div>
        </div>
      </section>
    </div>
  );
}
