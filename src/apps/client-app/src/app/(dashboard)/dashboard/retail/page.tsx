'use client';

import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { retailService, Order } from '@/services/retail.service';
import { ShoppingBag, CreditCard, CheckCircle2, Truck, Package, AlertCircle, RefreshCw, Loader2, Search, X } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import StatusPipeline, { PipelineStep } from '@/components/common/StatusPipeline';
import { AUTH_ACTIONS, AUTH_RESOURCES } from '@/constants/resources';
import { CasbinGuard } from '@/lib/auth';

const SAGA_STEPS: PipelineStep[] = [
  { key: 'PENDING', label: 'Created' },
  { key: 'PAYMENT_PENDING', label: 'Payment' },
  { key: 'PAYMENT_COMPLETED', label: 'Paid' },
  { key: 'RESERVED', label: 'Reserved' },
  { key: 'DISPATCH_REQUESTED', label: 'Awaiting' },
  { key: 'SHIPPING', label: 'Shipping' },
  { key: 'DELIVERED', label: 'Delivered' },
  { key: 'COMPLETED', label: 'Completed' },
];

const STATUS_CONFIG: Record<string, { label: string, color: string, icon: LucideIcon }> = {
  PENDING: { label: 'Awaiting Payment', color: 'bg-amber-50 text-amber-700 border-amber-200', icon: CreditCard },
  PAYMENT_PENDING: { label: 'Processing Payment', color: 'bg-blue-50 text-blue-700 border-blue-200', icon: Loader2 },
  PAYMENT_COMPLETED: { label: 'Payment Received', color: 'bg-green-50 text-green-700 border-green-200', icon: CheckCircle2 },
  RESERVED: { label: 'Inventory Reserved', color: 'bg-purple-50 text-purple-700 border-purple-200', icon: Package },
  DISPATCH_REQUESTED: { label: 'Awaiting Dispatch', color: 'bg-indigo-50 text-indigo-700 border-indigo-200', icon: Truck },
  SHIPPING: { label: 'In Transit', color: 'bg-primary/10 text-primary border-primary/20', icon: Truck },
  DELIVERED: { label: 'Delivered', color: 'bg-teal-50 text-teal-700 border-teal-200', icon: CheckCircle2 },
  COMPLETED: { label: 'Saga Completed', color: 'bg-slate-900 text-white border-slate-800', icon: CheckCircle2 },
  REJECTED: { label: 'Rejected', color: 'bg-error-container text-on-error-container border-error/20', icon: X },
  FAILED: { label: 'Saga Failed', color: 'bg-error text-white border-error', icon: AlertCircle },
};

export default function RetailOperationsPage() {
  const queryClient = useQueryClient();
  const [filter, setFilter] = useState('');

  const { data: orders = [], isLoading } = useQuery({
    queryKey: ['retail', 'orders'],
    queryFn: () => retailService.listOrders(),
    refetchInterval: 5000,
  });

  const payMutation = useMutation({
    mutationFn: (orderId: string) => retailService.simulatePayment(orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['retail', 'orders'] });
    }
  });

  const confirmMutation = useMutation({
    mutationFn: (orderId: string) => retailService.confirmOrder(orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['retail', 'orders'] });
    }
  });

  const filteredOrders = orders.filter(o =>
    o.id.toLowerCase().includes(filter.toLowerCase()) ||
    o.status.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div className="h-full flex flex-col gap-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <ShoppingBag className="w-7 h-7 text-primary" />
            Retail Operations
          </h1>
          <p className="text-[10px] font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Distributed Order Saga Orchestrator
          </p>
        </div>
        <div className="flex gap-4">
           <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <input
                type="text"
                placeholder="Search orders..."
                className="bg-white border border-slate-200 rounded-xl py-2 pl-10 pr-4 text-xs font-bold focus:border-primary outline-none w-64 shadow-sm"
                value={filter}
                onChange={e => setFilter(e.target.value)}
              />
           </div>
           <button
             onClick={() => queryClient.invalidateQueries({ queryKey: ['retail', 'orders'] })}
             className="p-2.5 bg-white border border-slate-200 rounded-xl hover:bg-slate-50 transition-colors shadow-sm"
           >
             <RefreshCw className={`w-4 h-4 text-slate-600 ${isLoading ? 'animate-spin' : ''}`} />
           </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden flex flex-col min-h-0">
        <div className="flex-1 overflow-y-auto">
          {isLoading ? (
            <div className="h-full flex flex-col items-center justify-center gap-4 text-slate-400">
               <Loader2 className="w-8 h-8 animate-spin" />
               <span className="text-[10px] font-black uppercase tracking-widest italic">Synchronizing Transaction Ledger...</span>
            </div>
          ) : filteredOrders.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center gap-4 text-slate-300">
               <ShoppingBag className="w-16 h-16 opacity-20" />
               <span className="text-[10px] font-black uppercase tracking-widest italic">No orders in current cycle.</span>
            </div>
          ) : (
            <div className="divide-y divide-slate-100">
              {filteredOrders.map(order => (
                <OrderRow
                  key={order.id}
                  order={order}
                  onPay={() => payMutation.mutate(order.id)}
                  isPaying={payMutation.isPending && payMutation.variables === order.id}
                  onConfirm={() => confirmMutation.mutate(order.id)}
                  isConfirming={confirmMutation.isPending && confirmMutation.variables === order.id}
                />
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Quick Summary */}
      <div className="flex gap-8 px-8 py-4 bg-slate-900 rounded-2xl shadow-xl">
         <div className="flex flex-col">
            <span className="text-[8px] font-black text-slate-500 uppercase tracking-widest mb-1">Total Volume</span>
            <span className="text-xl font-black text-white italic tracking-tighter">${orders.reduce((acc, o) => acc + o.total_amount, 0).toLocaleString()}</span>
         </div>
         <div className="w-px h-10 bg-slate-800 my-auto"></div>
         <div className="flex flex-col">
            <span className="text-[8px] font-black text-slate-500 uppercase tracking-widest mb-1">Active Sagas</span>
            <span className="text-xl font-black text-primary italic tracking-tighter">{orders.filter(o => !['COMPLETED', 'REJECTED', 'FAILED'].includes(o.status)).length}</span>
         </div>
         <div className="w-px h-10 bg-slate-800 my-auto"></div>
         <div className="flex flex-col">
            <span className="text-[8px] font-black text-slate-500 uppercase tracking-widest mb-1">Health Rate</span>
            <span className="text-xl font-black text-green-400 italic tracking-tighter">99.8%</span>
         </div>
      </div>
    </div>
  );
}

function OrderRow({ order, onPay, isPaying, onConfirm, isConfirming }: { order: Order, onPay: () => void, isPaying: boolean, onConfirm: () => void, isConfirming: boolean }) {
  const config = STATUS_CONFIG[order.status] || STATUS_CONFIG.PENDING;
  const StatusIcon = config.icon;

  return (
    <div className="p-8 hover:bg-slate-50/50 transition-colors group">
      <div className="flex justify-between items-start mb-8">
        <div className="flex gap-6 items-center">
          <div className="w-14 h-14 bg-slate-100 rounded-2xl flex items-center justify-center group-hover:bg-white group-hover:shadow-md transition-all">
            <ShoppingBag className="w-6 h-6 text-slate-400 group-hover:text-primary" />
          </div>
          <div>
            <div className="flex items-center gap-3">
              <h3 className="text-sm font-black text-slate-900 uppercase tracking-tight">Order #{order.id.slice(0, 8).toUpperCase()}</h3>
              <span className={`px-2.5 py-0.5 rounded-full text-[9px] font-black uppercase border flex items-center gap-1.5 ${config.color}`}>
                <StatusIcon className={`w-3 h-3 ${order.status === 'PAYMENT_PENDING' ? 'animate-spin' : ''}`} />
                {config.label}
              </span>
            </div>
            <div className="flex items-center gap-4 mt-1.5">
               <span className="text-[10px] font-bold text-slate-400 uppercase tracking-widest flex items-center gap-1">
                 <CreditCard className="w-3 h-3" /> ${order.total_amount.toLocaleString()}
               </span>
               <span className="text-[10px] font-bold text-slate-300 uppercase tracking-widest">
                 {new Date(order.created_at).toLocaleString()}
               </span>
            </div>
          </div>
        </div>

        <div className="flex gap-3">
          {order.status === 'PENDING' && (
            <CasbinGuard obj={AUTH_RESOURCES.ORDER} act={AUTH_ACTIONS.WRITE}>
              <button
                onClick={onPay}
                disabled={isPaying}
                className="bg-primary text-white px-6 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg shadow-primary/20 hover:bg-primary/90 transition-all disabled:opacity-50 flex items-center gap-2"
              >
                {isPaying ? <Loader2 className="w-3 h-3 animate-spin" /> : <CreditCard className="w-3 h-3" />}
                Pay via Stripe
              </button>
            </CasbinGuard>
          )}
          {order.status === 'DELIVERED' && (
            <CasbinGuard obj={AUTH_RESOURCES.ORDER} act={AUTH_ACTIONS.WRITE}>
              <button
                onClick={onConfirm}
                disabled={isConfirming}
                className="bg-green-600 text-white px-6 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg shadow-green-600/20 hover:bg-green-700 transition-all disabled:opacity-50 flex items-center gap-2"
              >
                {isConfirming ? <Loader2 className="w-3 h-3 animate-spin" /> : <CheckCircle2 className="w-3 h-3" />}
                Confirm Receipt
              </button>
            </CasbinGuard>
          )}
          {['COMPLETED'].includes(order.status) && (
             <div className="bg-green-50 text-green-600 p-2.5 rounded-xl border border-green-100">
               <CheckCircle2 className="w-5 h-5" />
             </div>
          )}
        </div>
      </div>

      <div className="pl-20">
        <StatusPipeline steps={SAGA_STEPS} currentStep={order.status} size="sm" />
      </div>
    </div>
  );
}
