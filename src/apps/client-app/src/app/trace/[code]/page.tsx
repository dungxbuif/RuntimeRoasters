'use client';

import { traceService } from '@/services/trace.service';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { CheckCircle2, Coffee, Factory, Loader2, MapPin, Package, Sprout, Store, Truck } from 'lucide-react';
import Link from 'next/link';
import { useParams } from 'next/navigation';

function milestoneIcon(status: string) {
  const value = status.toUpperCase();
  if (value.includes('HARVEST') || value.includes('FARM') || value.includes('ORIGIN')) return Sprout;
  if (value.includes('WAREHOUSE') || value.includes('PROCESS')) return Factory;
  if (value.includes('DELIVER') || value.includes('PICKUP') || value.includes('SHIP')) return Truck;
  if (value.includes('STOCK') || value.includes('BATCH')) return Package;
  if (value.includes('SERV')) return Coffee;
  return CheckCircle2;
}

export default function PublicTracePage() {
  const params = useParams();
  const traceCode = String(params?.code || '');
  const query = useQuery({ queryKey: ['public-cup-trace', traceCode], queryFn: () => traceService.getPublicTrace(traceCode), enabled: Boolean(traceCode), retry: false });

  if (query.isLoading) return <TraceState icon={Loader2} spin title="Tracing this cup" detail="Resolving its retail sale and upstream lineage." />;
  if (query.isError) {
    const missing = axios.isAxiosError(query.error) && query.error.response?.status === 404;
    return <TraceState icon={missing ? MapPin : Factory} title={missing ? 'Cup not found' : 'Trace temporarily unavailable'} detail={missing ? 'This trace code does not match a sold cup.' : 'The trace service could not assemble this journey. Please try again.'} />;
  }
  const trace = query.data!;

  return <main className="min-h-screen bg-slate-50 text-slate-950">
    <header className="border-b border-slate-200 bg-white"><div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-5"><Link href="/coffee" className="flex items-center gap-3"><div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-primary text-white"><Coffee className="h-5 w-5" /></div><div><div className="text-sm font-black uppercase tracking-tight">Runtime Roasters</div><div className="text-[9px] font-bold uppercase tracking-[0.2em] text-slate-400">Verified product journey</div></div></Link><div className="max-w-[48vw] text-right"><div className="text-[9px] font-black uppercase tracking-widest text-slate-400">Trace code</div><div className="truncate font-mono text-xs font-black text-primary">{trace.trace_code}</div></div></div></header>

    <div className="mx-auto max-w-5xl px-6 py-10 lg:py-14">
      <section className="grid gap-6 rounded-3xl border border-slate-200 bg-white p-7 shadow-sm md:grid-cols-[1fr_auto] md:items-center md:p-9">
        <div><div className="mb-3 flex items-center gap-2 text-[10px] font-black uppercase tracking-[0.18em] text-green-600"><CheckCircle2 className="h-4 w-4" />Public-safe trace verified</div><h1 className="text-4xl font-black tracking-[-0.05em]">Farm to Cup Journey</h1><p className="mt-3 text-sm leading-6 text-slate-500">{trace.product.name} · {trace.product.size}, served at {trace.store.name} in {trace.store.city}.</p></div>
        <div className="rounded-2xl bg-slate-950 px-6 py-5 text-white"><div className="text-[9px] font-black uppercase tracking-widest text-slate-400">Invoice</div><div className="mt-1 font-mono text-sm font-black">{trace.invoice_no}</div><div className="mt-3 text-[10px] text-slate-400">{new Date(trace.sale_item.sold_at).toLocaleString()}</div></div>
      </section>

      <div className="mt-8 grid gap-8 lg:grid-cols-[minmax(0,1fr)_280px]">
        <section><h2 className="mb-5 text-lg font-black">Journey milestones</h2><div className="relative"><div className="absolute bottom-6 left-6 top-6 w-px bg-slate-200" />
          <div className="space-y-4">{trace.milestones.map((milestone, index) => { const Icon = milestoneIcon(milestone.status); return <article key={`${milestone.status}-${milestone.timestamp}-${index}`} className="relative flex gap-5"><div className="z-10 flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-primary text-white shadow-lg shadow-primary/20"><Icon className="h-5 w-5" /></div><div className="flex-1 rounded-2xl border border-slate-200 bg-white p-5"><div className="flex flex-col justify-between gap-2 sm:flex-row sm:items-start"><div><h3 className="text-sm font-black">{milestone.title}</h3><p className="mt-1 text-xs leading-5 text-slate-500">{milestone.description}</p>{milestone.service && <div className="mt-2 text-[9px] font-black uppercase tracking-widest text-primary">{milestone.service}</div>}</div><time className="whitespace-nowrap text-[10px] font-bold text-slate-400">{new Date(milestone.timestamp).toLocaleString()}</time></div></div></article>; })}</div>
        </div></section>

        <aside className="space-y-4"><Info title="Sold cup" icon={Coffee} rows={[["Product ID", trace.product_id], ["SKU", trace.product.sku], ["Lot", trace.inventory_lot_id]]} /><Info title="Origin" icon={Sprout} rows={[["Harvest", trace.harvest_id], ["Batch", trace.origin_batch_id], ["Warehouse", trace.warehouse.id]]} /><Info title="Store" icon={Store} rows={[["Name", trace.store.name], ["City", trace.store.city], ["Address", trace.store.address]]} />
          {trace.service_participation.length > 0 && <div className="rounded-2xl border border-slate-200 bg-white p-5"><div className="text-[10px] font-black uppercase tracking-widest text-slate-400">Service participation</div><div className="mt-3 flex flex-wrap gap-2">{trace.service_participation.map(service => <span key={service} className="rounded-lg bg-slate-100 px-2.5 py-1 text-[10px] font-bold text-slate-600">{service}</span>)}</div></div>}
        </aside>
      </div>
    </div>
  </main>;
}

function Info({ title, icon: Icon, rows }: { title: string; icon: React.ElementType; rows: string[][] }) { return <div className="rounded-2xl border border-slate-200 bg-white p-5"><div className="mb-4 flex items-center gap-2"><Icon className="h-4 w-4 text-primary" /><h2 className="text-sm font-black">{title}</h2></div><div className="space-y-3">{rows.map(([label, value]) => <div key={label}><div className="text-[9px] font-black uppercase tracking-widest text-slate-400">{label}</div><div className="mt-0.5 break-all text-xs font-bold text-slate-700">{value || 'Not recorded'}</div></div>)}</div></div>; }

function TraceState({ icon: Icon, title, detail, spin = false }: { icon: React.ElementType; title: string; detail: string; spin?: boolean }) { return <main className="flex min-h-screen items-center justify-center bg-slate-50 p-6"><div className="max-w-md rounded-3xl border border-slate-200 bg-white p-10 text-center shadow-sm"><Icon className={`mx-auto h-10 w-10 text-primary ${spin ? 'animate-spin' : ''}`} /><h1 className="mt-5 text-2xl font-black">{title}</h1><p className="mt-2 text-sm leading-6 text-slate-500">{detail}</p><Link href="/coffee" className="mt-6 inline-flex rounded-xl bg-slate-950 px-5 py-3 text-[10px] font-black uppercase tracking-widest text-white">Return to coffee menu</Link></div></main>; }
