'use client';

import React from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { financeService } from '@/services/finance.service';
import { CreditCard, DollarSign, CheckCircle2, Clock, AlertCircle, RefreshCcw, ExternalLink, LucideIcon, XCircle } from 'lucide-react';
import { Payment, PaymentStatus } from '@/types/finance';

const STATUS_CONFIG: Record<PaymentStatus, { icon: LucideIcon, color: string, label: string }> = {
  'CREATED': { icon: Clock, color: 'text-slate-400 bg-slate-100', label: 'Created' },
  'PENDING': { icon: Clock, color: 'text-amber-500 bg-amber-50', label: 'Pending' },
  'SUCCEEDED': { icon: CheckCircle2, color: 'text-green-500 bg-green-50', label: 'Succeeded' },
  'FAILED': { icon: AlertCircle, color: 'text-red-500 bg-red-50', label: 'Failed' },
  'CANCELLED': { icon: AlertCircle, color: 'text-slate-500 bg-slate-50', label: 'Cancelled' },
  'REFUNDED': { icon: RefreshCcw, color: 'text-blue-500 bg-blue-50', label: 'Refunded' },
};

export default function FinancePage() {
  const queryClient = useQueryClient();
  const { data: payments = [], isLoading } = useQuery({
    queryKey: ['payments'],
    queryFn: () => financeService.listPayments(),
    refetchInterval: 10000,
  });
  const webhookMutation = useMutation({
    mutationFn: ({ payment, outcome }: { payment: Payment; outcome: 'passed' | 'failed' }) =>
      financeService.sendStripeWebhook(payment, outcome),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['payments'] }),
  });

  const stats = {
    total: payments.length,
    succeeded: payments.filter(p => p.status === 'SUCCEEDED').length,
    totalVolume: payments
      .filter(p => p.status === 'SUCCEEDED')
      .reduce((sum, p) => sum + p.amount, 0),
  };

  return (
    <div className="h-full flex flex-col gap-8">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <CreditCard className="w-8 h-8 text-primary" />
            Financial Integrity
          </h1>
          <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Transaction Ledger & Payment Gateway Reconciliation
          </p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-[2rem] border border-slate-200 shadow-sm flex items-center gap-4">
           <div className="w-12 h-12 rounded-2xl bg-blue-50 flex items-center justify-center text-primary">
              <DollarSign className="w-6 h-6" />
           </div>
           <div>
              <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Volume</div>
              <div className="text-2xl font-black text-slate-900">
                {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(stats.totalVolume)}
              </div>
           </div>
        </div>
        <div className="bg-white p-6 rounded-[2rem] border border-slate-200 shadow-sm flex items-center gap-4">
           <div className="w-12 h-12 rounded-2xl bg-green-50 flex items-center justify-center text-green-500">
              <CheckCircle2 className="w-6 h-6" />
           </div>
           <div>
              <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Successful</div>
              <div className="text-2xl font-black text-slate-900">{stats.succeeded} / {stats.total}</div>
           </div>
        </div>
        <div className="bg-white p-6 rounded-[2rem] border border-slate-200 shadow-sm flex items-center gap-4">
           <div className="w-12 h-12 rounded-2xl bg-slate-900 flex items-center justify-center text-white">
              <CreditCard className="w-6 h-6" />
           </div>
           <div>
              <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Payment Providers</div>
              <div className="text-2xl font-black text-slate-900">STRIPE, VNPAY</div>
           </div>
        </div>
      </div>

      {/* Transaction Table */}
      <div className="bg-white rounded-[2.5rem] border border-slate-200 shadow-xl overflow-hidden flex flex-col min-h-0">
        <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
           <h2 className="text-xs font-black uppercase tracking-[0.2em] text-slate-500 flex items-center gap-2">
             Transaction History
           </h2>
        </div>
        
        <div className="flex-1 overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50/30 text-[10px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-100">
                <th className="px-6 py-4">Transaction ID</th>
                <th className="px-6 py-4">Order ID</th>
                <th className="px-6 py-4">Amount</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">Gateway ID</th>
                <th className="px-6 py-4">Date</th>
                <th className="px-6 py-4 text-right">Webhook</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-50">
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="px-6 py-20 text-center">
                    <Loader2 className="w-8 h-8 text-primary animate-spin mx-auto mb-2" />
                    <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Fetching Ledger...</span>
                  </td>
                </tr>
              ) : payments.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-6 py-20 text-center">
                    <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest italic opacity-50">No transactions recorded</div>
                  </td>
                </tr>
              ) : (
                payments.map((payment) => {
                  const status = STATUS_CONFIG[payment.status] || STATUS_CONFIG.PENDING;
                  const isPending = payment.status === 'PENDING' || payment.status === 'CREATED';
                  const isPosting = webhookMutation.isPending && webhookMutation.variables?.payment.id === payment.id;
                  const providerRef = payment.provider_ref || payment.stripe_payment_intent_id;
                  return (
                    <tr key={payment.id} className="hover:bg-slate-50/50 transition-colors group">
                      <td className="px-6 py-4">
                        <span className="text-xs font-black font-mono text-slate-900">{payment.id.slice(0, 13).toUpperCase()}</span>
                      </td>
                      <td className="px-6 py-4">
                        <span className="text-xs font-bold text-slate-500 uppercase tracking-tighter">{payment.order_id.slice(0, 8)}</span>
                      </td>
                      <td className="px-6 py-4">
                        <span className="text-xs font-black text-slate-900 italic">
                          {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(payment.amount)}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <div className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-tighter border ${status.color}`}>
                          <status.icon className="w-3 h-3" />
                          {status.label}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        {providerRef ? (
                          <div className="flex items-center gap-2 text-slate-400 group-hover:text-primary transition-colors cursor-pointer">
                            <span className="text-[10px] font-bold font-mono truncate max-w-[120px]">{providerRef}</span>
                            <ExternalLink className="w-3 h-3" />
                          </div>
                        ) : (
                          <span className="text-[10px] font-bold text-slate-300 italic">INTERNAL</span>
                        )}
                      </td>
                      <td className="px-6 py-4">
                        <span className="text-[10px] font-bold text-slate-400 uppercase">{new Date(payment.created_at).toLocaleString()}</span>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            type="button"
                            disabled={!isPending || isPosting}
                            onClick={() => webhookMutation.mutate({ payment, outcome: 'passed' })}
                            className="inline-flex h-8 items-center gap-1.5 rounded-md border border-green-200 bg-green-50 px-2.5 text-[10px] font-black uppercase text-green-700 transition-colors hover:bg-green-100 disabled:cursor-not-allowed disabled:opacity-40"
                          >
                            <CheckCircle2 className="h-3.5 w-3.5" />
                            Pass
                          </button>
                          <button
                            type="button"
                            disabled={!isPending || isPosting}
                            onClick={() => webhookMutation.mutate({ payment, outcome: 'failed' })}
                            className="inline-flex h-8 items-center gap-1.5 rounded-md border border-red-200 bg-red-50 px-2.5 text-[10px] font-black uppercase text-red-700 transition-colors hover:bg-red-100 disabled:cursor-not-allowed disabled:opacity-40"
                          >
                            <XCircle className="h-3.5 w-3.5" />
                            Fail
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function Loader2(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
  );
}
