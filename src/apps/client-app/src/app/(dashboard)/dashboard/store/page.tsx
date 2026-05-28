'use client';

import NotificationFeed, { NotificationItem } from '@/components/common/NotificationFeed';
import StatusPipeline, { PipelineStep } from '@/components/common/StatusPipeline';
import { DashboardLayout } from '@/components/ui/templates/DashboardLayout';
import { DataGrid } from '@/components/ui/templates/DataGrid';
import { APP_ROUTES } from '@/constants/routes';
import { Order, retailService } from '@/services/retail.service';
import { useQuery } from '@tanstack/react-query';
import { Clock, MapPin, Plus, ShoppingBag, Store, Truck } from 'lucide-react';
import Link from 'next/link';
import { useMemo, useState } from 'react';

const ORDER_STEPS: PipelineStep[] = [
  { key: 'CREATED', label: 'Created' },
  { key: 'PAYMENT_PENDING', label: 'Payment' },
  { key: 'RESERVED', label: 'Reserved' },
  { key: 'DISPATCH_REQUESTED', label: 'Dispatch' },
  { key: 'IN_TRANSIT', label: 'In Transit' },
  { key: 'DELIVERED', label: 'Delivered' },
  { key: 'COMPLETED', label: 'Completed' },
];

const STATUS_COLORS: Record<string, string> = {
  CREATED: 'bg-slate-100 text-slate-600 border-slate-200',
  PAYMENT_PENDING: 'bg-amber-50 text-amber-700 border-amber-200',
  PAYMENT_SIMULATED: 'bg-amber-50 text-amber-700 border-amber-200',
  RESERVED: 'bg-blue-50 text-blue-700 border-blue-200',
  DISPATCH_REQUESTED: 'bg-indigo-50 text-indigo-700 border-indigo-200',
  IN_TRANSIT: 'bg-primary/10 text-primary border-primary/20',
  DELIVERED: 'bg-green-50 text-green-700 border-green-200',
  COMPLETED: 'bg-green-50 text-green-700 border-green-200',
  FAILED: 'bg-red-50 text-red-700 border-red-200',
};

export default function StoreDashboardPage() {
  const [selectedStoreId, setSelectedStoreId] = useState<string>('');

  const { data: stores = [] } = useQuery({
    queryKey: ['stores'],
    queryFn: () => retailService.listStores(),
  });

  const effectiveStoreId = selectedStoreId || stores[0]?.id || '';

  const { data: allOrders = [] } = useQuery({
    queryKey: ['orders'],
    queryFn: () => retailService.listOrders(),
    refetchInterval: 5000,
  });

  const orders = useMemo(() => {
    if (!effectiveStoreId) return [];
    return allOrders.filter(o => o.store_id === effectiveStoreId);
  }, [allOrders, effectiveStoreId]);

  const activeOrder = orders.find(o => o.status === 'IN_TRANSIT' || o.status === 'DISPATCH_REQUESTED');

  const notifications: NotificationItem[] = useMemo(() => [
    { id: '1', icon: 'check_circle', message: `Stock reserved for ${activeOrder?.id?.slice(0, 8).toUpperCase() || 'order'}`, time: '2min ago', color: 'green' },
    { id: '2', icon: 'local_shipping', message: 'Driver dispatched for delivery', time: '5min ago', color: 'blue' },
    { id: '3', icon: 'payments', message: 'Payment pending for ORD-8832', time: '12min ago', color: 'amber' },
  ], [activeOrder]);

  const formatVND = (amount: number) =>
    new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount);

  const orderColumns = [
    {
      key: 'id',
      header: 'Order ID',
      render: (order: Order) => <span className="text-xs font-black font-mono text-slate-900">{order.id.slice(0, 8).toUpperCase()}</span>
    },
    {
      key: 'items',
      header: 'Items',
      render: (_order: Order) => <span className="text-xs font-bold text-slate-600">Retail Order</span>

    },
    {
      key: 'amount',
      header: 'Amount',
      render: (order: Order) => <span className="text-xs font-black text-slate-900 italic">{formatVND(order.total_amount || 0)}</span>
    },
    {
      key: 'status',
      header: 'Status',
      render: (order: Order) => (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-[9px] font-black uppercase border ${STATUS_COLORS[order.status] || STATUS_COLORS.CREATED}`}>
          {order.status.replace(/_/g, ' ')}
        </span>
      )
    },
    {
      key: 'created',
      header: 'Created',
      render: (order: Order) => (
        <span className="text-[10px] font-bold text-slate-400">
          {new Date(order.created_at).toLocaleString('vi-VN', { hour: '2-digit', minute: '2-digit' })}
        </span>
      )
    }
  ];

  return (
    <DashboardLayout
      title="Retail Dashboard"
      subtitle="Store-scoped Order Management & Delivery Tracking"
      icon={Store}
      actions={
        <>
          <select
            value={effectiveStoreId}
            onChange={e => setSelectedStoreId(e.target.value)}
            className="bg-white border border-slate-200 rounded-xl px-4 py-2 text-xs font-bold uppercase tracking-tight appearance-none cursor-pointer shadow-sm focus:border-primary outline-none"
          >
            {stores.length === 0 && <option value="">No assigned stores</option>}
            {stores.map(s => (
              <option key={s.id} value={s.id}>{s.name} ({s.location})</option>
            ))}
          </select>
          <Link href={APP_ROUTES.DASHBOARD.RETAIL_ORDERS}
            className="flex items-center gap-2 bg-primary text-white px-5 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-primary/90 transition-all shadow-lg shadow-primary/20">
            <Plus className="w-4 h-4" />
            New Order
          </Link>
        </>
      }
    >
      {stores.length === 0 ? (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <Store className="w-16 h-16 text-slate-300 mx-auto mb-4" />
            <h3 className="text-sm font-black text-slate-400 uppercase tracking-widest">No Stores Assigned</h3>
            <p className="text-xs text-slate-400 mt-2">Contact ADMIN to assign store_ids to your account</p>
          </div>
        </div>
      ) : (
        <>
          {/* Main Content */}
          <div className="flex-1 flex flex-col gap-6 min-w-0 overflow-y-auto">
            {/* Order Status Pipeline */}
            {activeOrder && (
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Order Status Pipeline</h3>
                  <span className="text-[9px] font-mono font-bold text-primary bg-primary/10 px-2 py-0.5 rounded">
                    {activeOrder.id.slice(0, 8).toUpperCase()}
                  </span>
                </div>
                <StatusPipeline steps={ORDER_STEPS} currentStep={activeOrder.status} />
              </div>
            )}

            {/* Active Orders Table */}
            <div className="flex-1 overflow-hidden flex flex-col">
              <DataGrid
                title="Active Orders"
                icon={ShoppingBag}
                badge={`${orders.length} ORDERS`}
                data={orders}
                columns={orderColumns}
                keyExtractor={(o) => o.id}
              />
            </div>
          </div>

          {/* Right Panel */}
          <div className="w-80 flex flex-col gap-4 overflow-y-auto">
            {/* Incoming Delivery Widget */}
            {activeOrder && (
              <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                <div className="px-5 py-3 bg-primary/5 border-b border-primary/10">
                  <h3 className="text-[10px] font-black uppercase tracking-widest text-primary flex items-center gap-2">
                    <Truck className="w-4 h-4" />
                    Incoming Delivery
                  </h3>
                </div>
                <div className="p-5">
                  {/* Mini Map Placeholder */}
                  <div className="w-full h-32 bg-slate-100 rounded-xl mb-4 flex items-center justify-center relative overflow-hidden">
                    <div className="absolute inset-0 opacity-10" style={{
                      backgroundImage: 'radial-gradient(circle at 70% 60%, #004ac6 2px, transparent 2px)',
                      backgroundSize: '20px 20px'
                    }} />
                    <div className="absolute" style={{ top: '40%', left: '65%' }}>
                      <div className="w-4 h-4 bg-primary rounded-full animate-pulse shadow-[0_0_12px_rgba(0,74,198,0.5)]" />
                    </div>
                    <MapPin className="w-8 h-8 text-slate-300 absolute" style={{ top: '30%', left: '75%' }} />
                    <span className="text-[8px] font-black text-slate-400 uppercase tracking-widest absolute bottom-2 left-3">Live Tracking</span>
                  </div>

                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                      <Truck className="w-4 h-4 text-primary" />
                    </div>
                    <div>
                      <div className="text-[10px] font-black text-slate-900">Driver Alpha</div>
                      <div className="text-[9px] text-slate-400 font-bold flex items-center gap-1">
                        <Clock className="w-3 h-3" /> ETA ~12min
                      </div>
                    </div>
                  </div>

                  <div className="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                    <div className="bg-primary h-full rounded-full w-[75%] transition-all" />
                  </div>
                  <div className="flex justify-between text-[8px] font-bold text-slate-400 mt-1">
                    <span>Warehouse</span>
                    <span>75%</span>
                    <span>Store</span>
                  </div>
                </div>
              </div>
            )}

            {/* Notifications */}
            <NotificationFeed items={notifications} />
          </div>
        </>
      )}
    </DashboardLayout>
  );
}
