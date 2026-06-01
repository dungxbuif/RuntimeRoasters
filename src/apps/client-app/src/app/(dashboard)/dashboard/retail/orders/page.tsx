'use client';

import React, { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { retailService, CreateOrderRequest } from '@/services/retail.service';
import { Store, ShoppingBag, Plus, Trash2, Send, Loader2, AlertCircle } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { v4 as uuidv4 } from 'uuid';
import { CasbinGuard } from '@/lib/auth';

export default function CreateOrderPage() {
  const router = useRouter();
  const [selectedStore, setSelectedStore] = useState('');
  const [items, setItems] = useState<{ sku: string; quantity: number }[]>([
    { sku: 'COFFEE-ARABICA-001', quantity: 10 }
  ]);

  const { data: stores = [] } = useQuery({
    queryKey: ['stores'],
    queryFn: () => retailService.listStores(),
  });

  const mutation = useMutation({
    mutationFn: (data: { request: CreateOrderRequest; key: string }) => 
      retailService.createOrder(data.request, data.key),
    onSuccess: (order) => {
      router.push(`/dashboard/retail?orderId=${order.id}`);
    }
  });

  const addItem = () => {
    setItems([...items, { sku: '', quantity: 1 }]);
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const updateItem = (index: number, field: 'sku' | 'quantity', value: string | number) => {
    const newItems = [...items];
    const item = newItems[index];
    if (field === 'sku' && typeof value === 'string') {
      item.sku = value;
    } else if (field === 'quantity' && typeof value === 'number') {
      item.quantity = value;
    }
    setItems(newItems);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedStore || items.length === 0) return;
    
    mutation.mutate({
      request: {
        store_id: selectedStore,
        items: items
      },
      key: uuidv4()
    });
  };

  return (
    <CasbinGuard obj={AUTH_RESOURCES.ORDER} act={AUTH_ACTIONS.WRITE} fallback={<div className="p-12 text-center font-black uppercase italic text-slate-400">Access Denied: Retail Managers only.</div>}>
      <div className="h-full max-w-4xl mx-auto flex flex-col gap-8">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <ShoppingBag className="w-8 h-8 text-primary" />
            Supply Chain Order
          </h1>
          <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Initiate Saga Choreography for Store Replenishment
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-8 bg-white p-10 rounded-[3rem] border border-slate-200 shadow-2xl relative overflow-hidden">
        {/* Store Selection */}
        <div className="space-y-4">
          <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.3em] flex items-center gap-2">
            <Store className="w-3.5 h-3.5" />
            Receiving Store Node
          </label>
          <select 
            value={selectedStore}
            onChange={(e) => setSelectedStore(e.target.value)}
            required
            className="w-full bg-slate-50 border-2 border-slate-100 rounded-2xl py-4 px-6 text-sm font-bold uppercase tracking-tight focus:border-primary outline-none transition-all appearance-none cursor-pointer"
          >
            <option value="">Select a store location...</option>
            {stores.map(store => (
              <option key={store.id} value={store.id}>{store.name} ({store.location})</option>
            ))}
          </select>
        </div>

        {/* Order Items */}
        <div className="space-y-6">
          <div className="flex justify-between items-center">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.3em] flex items-center gap-2">
              <Plus className="w-3.5 h-3.5" />
              Manifest Inventory
            </label>
            <button 
              type="button"
              onClick={addItem}
              className="text-[9px] font-black text-primary uppercase tracking-widest border-b-2 border-primary/20 hover:border-primary pb-0.5 transition-all"
            >
              + Add Line Item
            </button>
          </div>

          <div className="space-y-3">
            {items.map((item, index) => (
              <div key={index} className="flex gap-4 items-center animate-in slide-in-from-left-4 duration-300">
                <input 
                  type="text"
                  placeholder="SKU Code (e.g. COFFEE-001)"
                  value={item.sku}
                  onChange={(e) => updateItem(index, 'sku', e.target.value)}
                  required
                  className="flex-1 bg-slate-50 border border-slate-100 rounded-xl py-3 px-4 text-xs font-bold font-mono focus:border-primary outline-none"
                />
                <input 
                  type="number"
                  min="1"
                  value={item.quantity}
                  onChange={(e) => updateItem(index, 'quantity', parseInt(e.target.value))}
                  required
                  className="w-24 bg-slate-50 border border-slate-100 rounded-xl py-3 px-4 text-xs font-black text-center focus:border-primary outline-none"
                />
                <button 
                  type="button"
                  onClick={() => removeItem(index)}
                  className="p-3 text-slate-300 hover:text-red-500 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>
        </div>

        {/* Submit Section */}
        <div className="pt-8 border-t border-slate-100 flex flex-col gap-6">
           <div className="bg-slate-900 rounded-[2rem] p-6 text-white flex items-center justify-between group">
              <div className="flex items-center gap-4">
                 <div className="w-10 h-10 rounded-xl bg-primary/20 flex items-center justify-center text-primary group-hover:scale-110 transition-transform">
                    <Send className="w-5 h-5" />
                 </div>
                 <div>
                    <div className="text-[8px] font-black text-slate-500 uppercase tracking-widest mb-1">Execution Mode</div>
                    <div className="text-xs font-black italic tracking-tighter">Distributed Transaction (Saga)</div>
                 </div>
              </div>
              <CasbinGuard obj={AUTH_RESOURCES.ORDER} act={AUTH_ACTIONS.WRITE}>
                <button 
                  type="submit"
                  disabled={mutation.isPending}
                  className="bg-white text-slate-900 px-8 py-3 rounded-xl text-[10px] font-black uppercase tracking-[0.2em] hover:bg-primary hover:text-white transition-all shadow-xl disabled:opacity-50 disabled:grayscale"
                >
                  {mutation.isPending ? (
                    <span className="flex items-center gap-2">
                      <Loader2 className="w-3 h-3 animate-spin" />
                      Processing...
                    </span>
                  ) : 'Broadcast Order'}
                </button>
              </CasbinGuard>
           </div>

           {mutation.isError && (
             <div className="flex items-center gap-3 p-4 bg-red-50 border border-red-100 rounded-2xl text-red-600 text-[10px] font-black uppercase tracking-widest italic animate-bounce">
                <AlertCircle className="w-4 h-4" />
                Saga Initialization Failed: {(mutation.error as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Network Error'}
             </div>
           )}
        </div>
        
        {/* Background Sketch */}
        <div className="absolute top-0 right-0 -translate-y-1/2 translate-x-1/2 w-64 h-64 bg-primary/5 rounded-full blur-3xl pointer-events-none"></div>
      </form>

      {/* Info Panel */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 opacity-60 grayscale hover:grayscale-0 hover:opacity-100 transition-all">
         <div className="bg-slate-50 border border-slate-200 p-6 rounded-3xl">
            <h4 className="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-2">Protocol: Atomicity</h4>
            <p className="text-[10px] font-bold text-slate-600 leading-relaxed italic">
              All or nothing. If any service in the chain fails, automatic compensations will restore state consistency.
            </p>
         </div>
         <div className="bg-slate-50 border border-slate-200 p-6 rounded-3xl">
            <h4 className="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-2">Protocol: Idempotency</h4>
            <p className="text-[10px] font-bold text-slate-600 leading-relaxed italic">
              Guaranteed exactly-once processing using unique execution keys per transaction broadcast.
            </p>
         </div>
      </div>
    </div>
    </CasbinGuard>
  );
}
