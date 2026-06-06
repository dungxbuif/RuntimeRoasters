'use client';

import { traceService, PublicTraceItem } from '@/services/trace.service';
import { useQuery } from '@tanstack/react-query';
import { Coffee, Loader2, QrCode, Search, ExternalLink, ArrowRight } from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';
import { QRCodeSVG } from 'qrcode.react';

export default function PublicTraceListPage() {
  const [search, setSearch] = useState('');
  
  const { data: items = [], isLoading } = useQuery({
    queryKey: ['public-trace-items'],
    queryFn: () => traceService.listPublicTraceItems(50),
    refetchInterval: 30000,
  });

  const filteredItems = items.filter(item => 
    item.product_name.toLowerCase().includes(search.toLowerCase()) ||
    item.invoice_no.toLowerCase().includes(search.toLowerCase()) ||
    item.trace_code.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <main className="min-h-screen bg-slate-50 text-slate-950">
      <header className="border-b border-slate-200 bg-white sticky top-0 z-50">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4 lg:px-10">
          <Link href="/" className="flex items-center gap-3 group">
            <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-primary text-white shadow-lg group-hover:scale-105 transition-transform">
              <Coffee className="h-5 w-5" />
            </div>
            <div>
              <div className="text-sm font-black uppercase tracking-tight">Runtime Roasters</div>
              <div className="text-[9px] font-bold uppercase tracking-[0.2em] text-slate-400">Verified Journey</div>
            </div>
          </Link>
          <nav className="hidden md:flex items-center gap-8">
            <Link href="/" className="text-xs font-black uppercase tracking-widest text-slate-400 hover:text-primary transition-colors">Home</Link>
            <Link href="/coffee" className="text-xs font-black uppercase tracking-widest text-slate-400 hover:text-primary transition-colors">Buy Coffee</Link>
            <Link href="/trace" className="text-xs font-black uppercase tracking-widest text-primary border-b-2 border-primary pb-1">Traceability</Link>
          </nav>
          <Link href="/login" className="px-5 py-2 bg-slate-950 text-white rounded-xl font-black uppercase tracking-widest text-[10px] hover:bg-primary transition-all shadow-lg shadow-slate-900/10">
            Terminal Login
          </Link>
        </div>
      </header>

      <div className="mx-auto max-w-7xl px-6 py-10 lg:px-10 lg:py-16">
        <div className="mb-12">
          <div className="flex flex-col md:flex-row md:items-end justify-between gap-6">
            <div>
              <h1 className="text-5xl font-black tracking-tight leading-none mb-4">
                Public <span className="text-primary">Trace</span> Registry
              </h1>
              <p className="text-slate-500 max-w-2xl font-medium leading-relaxed">
                Browse recently sold cups and verify their complete farm-to-cup journey. 
                Every cup is uniquely identified and cryptographically linked to its origin.
              </p>
            </div>
            <div className="relative group">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400 group-focus-within:text-primary transition-colors" />
              <input 
                type="text" 
                placeholder="Search by invoice or cup name..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full md:w-80 bg-white border border-slate-200 rounded-2xl py-4 pl-12 pr-6 text-sm font-bold shadow-sm focus:border-primary focus:ring-4 focus:ring-primary/5 outline-none transition-all"
              />
            </div>
          </div>
        </div>

        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-24 gap-4 text-slate-400">
            <Loader2 className="h-10 w-10 animate-spin" />
            <p className="text-[10px] font-black uppercase tracking-widest italic">Decrypting Registry...</p>
          </div>
        ) : filteredItems.length === 0 ? (
          <div className="bg-white rounded-[3rem] border border-slate-200 p-20 text-center shadow-sm">
            <QrCode className="h-16 w-16 text-slate-200 mx-auto mb-6" />
            <h2 className="text-2xl font-black mb-2">No traceable cups found</h2>
            <p className="text-slate-500">Try adjusting your search or buy a fresh cup at the kiosk.</p>
            <Link href="/coffee" className="mt-8 inline-flex items-center gap-2 bg-primary text-white px-8 py-3 rounded-2xl font-black uppercase tracking-widest text-[11px] shadow-lg shadow-primary/20 hover:scale-[1.02] transition-all">
              Go to Coffee Kiosk <ArrowRight className="h-4 w-4" />
            </Link>
          </div>
        ) : (
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {filteredItems.map((item) => (
              <TraceCard key={item.trace_code} item={item} />
            ))}
          </div>
        )}
      </div>

      <footer className="border-t border-slate-200 bg-white py-12 mt-12">
        <div className="mx-auto max-w-7xl px-6 lg:px-10 flex flex-col md:flex-row justify-between items-center gap-8">
          <div className="flex items-center gap-3 opacity-50 grayscale">
            <Coffee className="h-6 w-6" />
            <span className="text-xs font-black uppercase tracking-widest">Runtime Roasters &copy; 2026</span>
          </div>
          <div className="flex gap-8 text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
            <Link href="/" className="hover:text-primary">Architecture</Link>
            <Link href="/coffee" className="hover:text-primary">Menu</Link>
            <Link href="/api-docs" className="hover:text-primary">Public API</Link>
          </div>
        </div>
      </footer>
    </main>
  );
}

function TraceCard({ item }: { item: PublicTraceItem }) {
  const traceURL = `/trace/${encodeURIComponent(item.trace_code)}`;
  
  return (
    <div className="group bg-white border border-slate-200 rounded-[2.5rem] p-6 shadow-sm hover:shadow-xl hover:border-primary/20 transition-all duration-500 flex flex-col h-full">
      <div className="mb-6 relative aspect-square bg-slate-50 rounded-3xl flex items-center justify-center overflow-hidden border border-slate-100 group-hover:bg-primary/5 transition-colors">
        <div className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center z-10">
          <Link 
            href={traceURL}
            className="bg-slate-950 text-white p-4 rounded-2xl shadow-2xl scale-90 group-hover:scale-100 transition-transform"
          >
            <ExternalLink className="h-6 w-6" />
          </Link>
        </div>
        <QRCodeSVG 
          value={`${typeof window !== 'undefined' ? window.location.origin : ''}${traceURL}`} 
          size={160} 
          className="relative z-0 group-hover:scale-95 transition-transform duration-500 opacity-80 group-hover:opacity-100"
          includeMargin
        />
      </div>

      <div className="flex-1">
        <div className="flex justify-between items-start gap-4 mb-2">
          <h3 className="text-lg font-black leading-tight group-hover:text-primary transition-colors">{item.product_name}</h3>
          <span className="shrink-0 text-[10px] font-black text-slate-400 bg-slate-100 px-2 py-0.5 rounded uppercase">{item.store_name}</span>
        </div>
        <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-4">Invoice {item.invoice_no}</p>
        
        <div className="space-y-1.5 mb-6">
          <div className="text-[9px] font-black uppercase text-slate-400">Trace Code</div>
          <code className="block text-[10px] font-mono font-bold text-slate-600 truncate bg-slate-50 px-2 py-1.5 rounded-lg border border-slate-100 italic">
            {item.trace_code}
          </code>
        </div>
      </div>

      <div className="flex items-center justify-between pt-6 border-t border-slate-50">
        <span className="text-[10px] font-bold text-slate-400 italic">
          {new Date(item.sold_at).toLocaleString([], { dateStyle: 'short', timeStyle: 'short' })}
        </span>
        <Link 
          href={traceURL}
          className="text-[10px] font-black uppercase tracking-widest text-primary flex items-center gap-1.5 group/link"
        >
          Verify Lineage <ArrowRight className="h-3 w-3 group-hover/link:translate-x-1 transition-transform" />
        </Link>
      </div>
    </div>
  );
}
