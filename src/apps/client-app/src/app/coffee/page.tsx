'use client';

import { expandSaleItems } from '@/features/retail-sale/sale';
import { PublicMenuItem, retailService, Sale } from '@/services/retail.service';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Check, Coffee, Loader2, MapPin, Minus, Plus, ReceiptText, ShoppingBag, X } from 'lucide-react';
import Link from 'next/link';
import { QRCodeSVG } from 'qrcode.react';
import { useMemo, useState } from 'react';

const formatVND = (amount: number) => new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount);

export default function CoffeeKioskPage() {
  const [storeId, setStoreId] = useState('');
  const [quantities, setQuantities] = useState<Record<string, number>>({});
  const [invoice, setInvoice] = useState<Sale | null>(null);

  const storesQuery = useQuery({ queryKey: ['public-retail-stores'], queryFn: () => retailService.listPublicStores() });
  const stores = storesQuery.data || [];
  const effectiveStoreId = storeId || stores[0]?.id || '';
  const selectedStore = stores.find(store => store.id === effectiveStoreId);
  const menuQuery = useQuery({
    queryKey: ['public-store-menu', effectiveStoreId],
    queryFn: () => retailService.listPublicStoreMenu(effectiveStoreId),
    enabled: Boolean(effectiveStoreId),
  });
  const menuItems = useMemo(() => (menuQuery.data || []).flatMap(menu => menu.items), [menuQuery.data]);
  const selectedItems = menuItems.filter(item => (quantities[item.id] || 0) > 0);
  const cupCount = selectedItems.reduce((sum, item) => sum + (quantities[item.id] || 0), 0);
  const total = selectedItems.reduce((sum, item) => sum + item.price * (quantities[item.id] || 0), 0);

  const saleMutation = useMutation({
    mutationFn: async () => {
      if (!selectedStore) throw new Error('Choose a store before buying coffee.');
      const items = expandSaleItems(selectedItems.map(item => ({ menuItemId: item.id, quantity: quantities[item.id] })), selectedStore.code || selectedStore.id);
      return retailService.createPublicSale({ store_id: selectedStore.id, items }, crypto.randomUUID());
    },
    onSuccess: sale => {
      setInvoice(sale);
      setQuantities({});
      void menuQuery.refetch();
    },
  });

  function changeQuantity(item: PublicMenuItem, delta: number) {
    setQuantities(current => {
      const next = Math.max(0, Math.min(item.available_units || 0, (current[item.id] || 0) + delta));
      return { ...current, [item.id]: next };
    });
  }

  return (
    <main className="min-h-screen bg-slate-50 text-slate-950">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-5 lg:px-10">
          <div className="flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-primary text-white shadow-lg shadow-primary/20"><Coffee className="h-5 w-5" /></div>
            <div><div className="text-sm font-black uppercase tracking-tight">Runtime Roasters</div><div className="text-[10px] font-bold uppercase tracking-[0.22em] text-slate-400">Public coffee kiosk</div></div>
          </div>
          <Link href="/" className="text-xs font-bold text-slate-500 transition hover:text-primary">Back to platform</Link>
        </div>
      </header>

      <div className="mx-auto grid max-w-7xl gap-8 px-6 py-8 lg:grid-cols-[minmax(0,1fr)_360px] lg:px-10 lg:py-12">
        <section>
          <div className="mb-8 flex flex-col justify-between gap-5 md:flex-row md:items-end">
            <div><h1 className="text-4xl font-black tracking-[-0.05em]">Coffee menu</h1><p className="mt-2 max-w-xl text-sm leading-6 text-slate-500">Choose a store, select each cup, and receive a unique traceable QR on your invoice.</p></div>
            <label className="min-w-72">
              <span className="mb-2 block text-[10px] font-black uppercase tracking-[0.18em] text-slate-400">Choose your store</span>
              <div className="relative"><MapPin className="pointer-events-none absolute left-4 top-3.5 h-4 w-4 text-primary" /><select value={effectiveStoreId} onChange={event => { setStoreId(event.target.value); setQuantities({}); }} className="h-11 w-full appearance-none rounded-xl border border-slate-200 bg-white pl-11 pr-4 text-sm font-bold outline-none focus:border-primary">
                {stores.map(store => <option key={store.id} value={store.id}>{store.name} · {store.city}</option>)}
              </select></div>
            </label>
          </div>

          {storesQuery.isLoading || menuQuery.isLoading ? <Loading /> : menuQuery.isError ? <Message title="Menu unavailable" detail="The store menu could not be loaded. Please try again." /> : (
            <div className="space-y-8">
              {(menuQuery.data || []).map(menu => <div key={menu.id}>
                <div className="mb-4 flex items-baseline justify-between border-b border-slate-200 pb-3"><div><h2 className="text-lg font-black">{menu.name}</h2><p className="text-xs text-slate-400">{menu.description}</p></div><span className="text-[10px] font-black uppercase tracking-widest text-slate-400">{menu.items.length} choices</span></div>
                <div className="grid gap-3 md:grid-cols-2">
                  {menu.items.map(item => {
                    const quantity = quantities[item.id] || 0;
                    const available = Boolean(item.available && item.available_units);
                    return <article key={item.id} className={`flex min-h-36 flex-col justify-between rounded-2xl border bg-white p-5 transition ${quantity ? 'border-primary shadow-lg shadow-primary/10' : 'border-slate-200'} ${!available ? 'opacity-55' : 'hover:border-slate-300 hover:shadow-sm'}`}>
                      <div className="flex justify-between gap-4"><div><div className="text-[10px] font-black uppercase tracking-[0.18em] text-primary">{item.category_code} · {item.size}</div><h3 className="mt-1 text-base font-black">{item.name}</h3><p className="mt-1 line-clamp-2 text-xs leading-5 text-slate-500">{item.description}</p></div><div className="text-right"><div className="text-sm font-black">{formatVND(item.price)}</div><div className="mt-1 text-[10px] font-bold text-slate-400">{item.available_units || 0} available</div></div></div>
                      <div className="mt-4 flex items-center justify-between"><span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">{item.coffee_type}</span><div className="flex items-center gap-3"><button aria-label={`Remove ${item.name}`} onClick={() => changeQuantity(item, -1)} disabled={!quantity} className="flex h-8 w-8 items-center justify-center rounded-lg border border-slate-200 disabled:opacity-30"><Minus className="h-3.5 w-3.5" /></button><span className="w-5 text-center text-sm font-black">{quantity}</span><button aria-label={`Add ${item.name}`} onClick={() => changeQuantity(item, 1)} disabled={!available || quantity >= (item.available_units || 0)} className="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-950 text-white disabled:opacity-30"><Plus className="h-3.5 w-3.5" /></button></div></div>
                    </article>;
                  })}
                </div>
              </div>)}
            </div>
          )}
        </section>

        <aside className="lg:sticky lg:top-8 lg:h-fit">
          <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-xl shadow-slate-200/50">
            <div className="flex items-center gap-3 border-b border-slate-100 pb-5"><div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-950 text-white"><ShoppingBag className="h-4 w-4" /></div><div><h2 className="font-black">Your order</h2><p className="text-xs text-slate-400">{selectedStore?.name || 'Choose a store'}</p></div></div>
            <div className="min-h-40 py-5">{selectedItems.length === 0 ? <div className="flex h-32 flex-col items-center justify-center text-center text-slate-400"><Coffee className="mb-3 h-7 w-7 opacity-40" /><p className="text-xs font-bold">Add a coffee to begin</p></div> : <div className="space-y-4">{selectedItems.map(item => <div key={item.id} className="flex justify-between gap-4"><div><div className="text-sm font-bold">{item.name}</div><div className="text-[11px] text-slate-400">{item.size} × {quantities[item.id]}</div></div><div className="text-sm font-black">{formatVND(item.price * quantities[item.id])}</div></div>)}</div>}</div>
            <div className="border-t border-slate-100 pt-5"><div className="mb-5 flex items-end justify-between"><div><div className="text-[10px] font-black uppercase tracking-widest text-slate-400">Total · {cupCount} cups</div><div className="mt-1 text-2xl font-black tracking-tight">{formatVND(total)}</div></div><ReceiptText className="h-6 w-6 text-primary" /></div>
              <button onClick={() => saleMutation.mutate()} disabled={!cupCount || saleMutation.isPending} className="flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-primary text-xs font-black uppercase tracking-[0.16em] text-white shadow-lg shadow-primary/20 transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-40">{saleMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Check className="h-4 w-4" />}Buy coffee</button>
              {saleMutation.isError && <p className="mt-3 text-xs font-medium text-red-600">{saleMutation.error.message}</p>}
            </div>
          </div>
        </aside>
      </div>
      {invoice && <InvoiceModal sale={invoice} onClose={() => setInvoice(null)} />}
    </main>
  );
}

function InvoiceModal({ sale, onClose }: { sale: Sale; onClose: () => void }) {
  return <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/55 p-4 backdrop-blur-sm" role="dialog" aria-modal="true" aria-label="Invoice">
    <div className="max-h-[92vh] w-full max-w-2xl overflow-y-auto rounded-3xl bg-white shadow-2xl"><div className="flex items-start justify-between border-b border-slate-100 p-6"><div><div className="text-[10px] font-black uppercase tracking-[0.2em] text-green-600">Purchase completed</div><h2 className="mt-1 text-2xl font-black">Invoice {sale.invoice_no}</h2><p className="mt-1 text-xs text-slate-400">Each QR identifies one physical cup.</p></div><button aria-label="Close invoice" onClick={onClose} className="rounded-xl border border-slate-200 p-2"><X className="h-4 w-4" /></button></div>
      <div className={`grid gap-4 p-6 ${sale.items.length > 1 ? 'sm:grid-cols-2' : 'mx-auto max-w-sm'}`}>{sale.items.map(item => { const traceURL = `${window.location.origin}/trace/${encodeURIComponent(item.trace_code)}`; return <div key={item.id} className="rounded-2xl border border-slate-200 p-5"><QRCodeSVG value={traceURL} size={144} className="mx-auto" /><div className="mt-4 text-center"><div className="font-black">{item.product_name} · {item.size}</div><div className="mt-1 break-all font-mono text-[10px] font-bold text-primary">{item.product_id}</div><Link href={`/trace/${encodeURIComponent(item.trace_code)}`} className="mt-4 inline-flex rounded-lg bg-slate-950 px-4 py-2 text-[10px] font-black uppercase tracking-widest text-white">Open trace</Link></div></div>; })}</div>
      <div className="flex items-center justify-between border-t border-slate-100 bg-slate-50 px-6 py-5"><span className="text-xs font-bold text-slate-500">{sale.items.length} traceable cups</span><span className="text-lg font-black">{formatVND(sale.total_amount)}</span></div>
    </div>
  </div>;
}

function Loading() { return <div className="flex min-h-80 items-center justify-center gap-3 text-sm font-bold text-slate-400"><Loader2 className="h-5 w-5 animate-spin" />Loading coffee menu</div>; }
function Message({ title, detail }: { title: string; detail: string }) { return <div className="rounded-2xl border border-slate-200 bg-white p-10 text-center"><h2 className="font-black">{title}</h2><p className="mt-2 text-sm text-slate-500">{detail}</p></div>; }
