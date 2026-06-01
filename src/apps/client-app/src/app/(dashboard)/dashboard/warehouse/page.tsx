'use client';

import { NotificationItem } from '@/components/common/NotificationFeed';
import StatusPipeline from '@/components/common/StatusPipeline';
import { DashboardLayout } from '@/components/ui/templates/DashboardLayout';
import { logisticsService } from '@/services/logistics.service';
import { warehouseService } from '@/services/warehouse.service';
import { DispatchRequest, PickupRequest } from '@/services/warehouse.service';
import { Driver, Vehicle } from '@/types/logistics';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { CheckCircle2, MapPin, Truck, Warehouse, UserPlus, X, Loader2, Package } from 'lucide-react';
import { useMemo, useState } from 'react';
import { AUTH_ACTIONS, AUTH_RESOURCES } from '@/constants/resources';
import { CasbinGuard } from '@/lib/auth';

export default function WarehouseOperationsPage() {
  const queryClient = useQueryClient();
  const [assignModal, setAssignModal] = useState<{ type: 'pickup' | 'delivery', id: string } | null>(null);
  const [assignmentData, setAssignmentData] = useState({ driver_id: '', vehicle_id: '' });

  // Queries
  const { data: batches = [] } = useQuery({
    queryKey: ['warehouse', 'batches'],
    queryFn: () => warehouseService.listBatches(),
    refetchInterval: 5000,
  });

  const { data: inventory = [] } = useQuery({
    queryKey: ['warehouse', 'inventory'],
    queryFn: () => warehouseService.getInventory(),
    refetchInterval: 10000,
  });

  const { data: pickupRequests = [] } = useQuery({
    queryKey: ['warehouse', 'pickupRequests'],
    queryFn: () => warehouseService.listPickupRequests(),
    refetchInterval: 5000,
  });

  const { data: dispatchRequests = [] } = useQuery({
    queryKey: ['warehouse', 'dispatchRequests'],
    queryFn: () => warehouseService.listDispatchRequests(),
    refetchInterval: 5000,
  });

  const { data: drivers = [] } = useQuery({
    queryKey: ['logistics', 'drivers'],
    queryFn: () => logisticsService.listDrivers(),
  });

  const { data: vehicles = [] } = useQuery({
    queryKey: ['logistics', 'vehicles'],
    queryFn: () => logisticsService.listAvailableVehicles(),
  });

  // Mutations
  const dispatchMutation = useMutation<PickupRequest | DispatchRequest, Error, { type: 'pickup' | 'delivery', id: string, driver_id: string, vehicle_id: string }>({
    mutationFn: (args) => {
      if (args.type === 'pickup') return warehouseService.dispatchPickupRequest(args.id, args.driver_id, args.vehicle_id);
      return warehouseService.dispatchRequest(args.id, args.driver_id, args.vehicle_id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouse'] });
      setAssignModal(null);
      setAssignmentData({ driver_id: '', vehicle_id: '' });
    }
  });

  const receiveMutation = useMutation({
    mutationFn: (id: string) => warehouseService.receivePickupRequest(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouse'] });
    }
  });

  // Derived state
  const activeProcessingBatches = useMemo(() =>
    batches.filter(b => ['RECEIVED', 'HULLING', 'ROASTING'].includes(b.status)),
    [batches]
  );

  const handleDispatchClick = (type: 'pickup' | 'delivery', id: string) => {
    setAssignModal({ type, id });
  };

  const handleFinalizeDispatch = () => {
    if (!assignModal) return;
    dispatchMutation.mutate({
      ...assignModal,
      ...assignmentData
    });
  };

  const handleConfirmReceipt = (id: string) => {
    receiveMutation.mutate(id);
  };

  const notifications: NotificationItem[] = [
    { id: '1', icon: 'inventory_2', message: 'Low stock warning: SKU-AR-001', time: '10m ago', color: 'blue' },
    { id: '2', icon: 'local_shipping', message: 'Driver arriving in 5 mins', time: '2m ago', color: 'green' },
  ];

  return (
    <DashboardLayout
      title="Warehouse Operations"
      subtitle="Fulfillment, Inventory & Processing"
      icon={Warehouse}
    >
      <div className="h-full grid grid-cols-1 lg:grid-cols-3 gap-8">

        {/* LEFT: INBOUND QUEUE */}
        <div className="flex flex-col gap-4">
          <h2 className="text-xl font-black text-slate-900 flex items-center justify-between">
            Inbound Queue
            <span className="text-[10px] font-black bg-primary text-white px-2 py-0.5 rounded-full">{pickupRequests.length}</span>
          </h2>
          <div className="bg-white rounded-[2rem] border border-slate-200 shadow-sm flex-1 p-5 flex flex-col gap-4 overflow-y-auto custom-scrollbar">
            {pickupRequests.map(req => (
              <div key={req.id} className="border border-slate-200 rounded-2xl p-4 shadow-sm hover:border-slate-300 transition-colors">
                 <div className="flex justify-between items-start mb-2">
                    <h4 className="font-bold text-slate-900 text-sm">Harvest #{req.harvest_id.slice(-6)}</h4>
                    <span className={`text-[9px] font-black px-2 py-0.5 rounded-full uppercase ${
                      req.status === 'REQUESTED' ? 'bg-amber-100 text-amber-700' : 'bg-blue-100 text-blue-700'
                    }`}>{req.status}</span>
                 </div>
                 <div className="flex flex-wrap items-center gap-2 mb-4">
                    <span className="text-[9px] font-black bg-slate-100 text-slate-500 px-2 py-1 rounded-full uppercase">{req.coffee_type}</span>
                    <span className="text-[9px] font-bold text-slate-400 italic">{req.quantity} KG</span>
                 </div>

                 {req.status === 'REQUESTED' && (
                   <CasbinGuard obj={AUTH_RESOURCES.WAREHOUSE_DISPATCH} act={AUTH_ACTIONS.WRITE}>
                    <button onClick={() => handleDispatchClick('pickup', req.id)} className="w-full bg-primary text-white rounded-xl py-2.5 text-[10px] font-black uppercase tracking-widest flex items-center justify-center gap-2 hover:bg-primary/90 transition-all">
                      <Truck className="w-4 h-4" /> Dispatch Pickup
                    </button>
                   </CasbinGuard>
                 )}
                 {req.status === 'DISPATCHED' && (
                   <div className="space-y-2">
                     <p className="text-[9px] font-black text-primary uppercase tracking-widest flex items-center gap-1 italic">
                       <Loader2 className="w-3 h-3 animate-spin" /> Driver en route to farm...
                     </p>
                     <div className="w-full bg-slate-100 h-1.5 rounded-full overflow-hidden"><div className="bg-blue-500 h-full w-[45%]" /></div>
                   </div>
                 )}
                 {(req.status === 'ARRIVED_WAREHOUSE' || req.status === 'PICKED_UP') && (
                   <CasbinGuard obj={AUTH_RESOURCES.WAREHOUSE_RECEIVE} act={AUTH_ACTIONS.WRITE}>
                    <button onClick={() => handleConfirmReceipt(req.id)} className="w-full mt-3 bg-green-600 text-white rounded-xl py-2.5 text-[10px] font-black uppercase tracking-widest flex items-center justify-center gap-2 hover:bg-green-700 transition-all">
                      <CheckCircle2 className="w-4 h-4" /> Confirm Receipt
                    </button>
                   </CasbinGuard>
                 )}
              </div>
            ))}
          </div>
        </div>

        {/* MIDDLE: PROCESSING & INVENTORY */}
        <div className="flex flex-col gap-4">
          <h2 className="text-xl font-black text-slate-900 flex items-center justify-between">
            Processing & Inventory
          </h2>
          <div className="bg-white rounded-[2rem] border border-slate-200 shadow-sm flex-1 p-5 flex flex-col gap-6 overflow-y-auto custom-scrollbar">

             {/* Processing Section */}
             <div>
               <div className="flex justify-between items-center pb-2 border-b border-slate-100 mb-4">
                  <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">Processing Queue</span>
               </div>

               {activeProcessingBatches.length === 0 ? (
                 <div className="text-center py-6 border border-slate-100 rounded-2xl bg-slate-50/50">
                    <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest italic">No active batches</p>
                 </div>
               ) : (
                 <div className="space-y-3">
                   {activeProcessingBatches.map(batch => (
                     <div key={batch.id} className="border border-slate-200 rounded-xl p-3">
                        <div className="flex justify-between mb-2">
                           <span className="font-bold text-xs">Batch #{batch.batch_id}</span>
                        </div>
                        <StatusPipeline
                          steps={[
                            {key: 'RECEIVED', label: 'Received'},
                            {key: 'HULLING', label: 'Hulling'},
                            {key: 'ROASTING', label: 'Roasting'},
                            {key: 'STOCKED', label: 'Stocked'}
                          ]}
                          currentStep={batch.status}
                          size="sm"
                        />
                     </div>
                   ))}
                 </div>
               )}
             </div>

             {/* Inventory Section */}
             <div>
               <div className="flex justify-between items-center pb-2 border-b border-slate-100 mb-4">
                  <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">Finished Stock</span>
               </div>

               {inventory.length === 0 ? (
                 <div className="text-center py-6 border border-slate-100 rounded-2xl bg-slate-50/50">
                    <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest italic">No finished stock</p>
                 </div>
               ) : (
                 <div className="space-y-4">
                    {inventory.map(item => (
                      <div key={item.id}>
                        <div className="flex justify-between text-[10px] font-black uppercase mb-1">
                           <span className="text-slate-900">{item.sku}</span>
                           <span className="text-slate-500">{item.available_quantity} Available</span>
                        </div>
                        <div className="w-full bg-slate-200 h-4 rounded overflow-hidden flex text-[8px] font-black text-white items-center">
                           <div className="bg-primary h-full px-1 flex items-center justify-end" style={{ width: `${Math.min(item.available_quantity/2000*100, 80)}%` }}>{item.available_quantity}</div>
                           <div className="bg-slate-300 h-full px-1 flex items-center text-slate-500" style={{ width: '20%' }}>Reserved</div>
                        </div>
                      </div>
                    ))}
                 </div>
               )}
             </div>

          </div>
        </div>

        {/* RIGHT: OUTBOUND QUEUE */}
        <div className="flex flex-col gap-4">
          <h2 className="text-xl font-black text-slate-900 flex items-center justify-between">
            Outbound Queue
          </h2>
          <div className="bg-white rounded-[2rem] border border-slate-200 shadow-sm flex-1 p-5 flex flex-col gap-4 overflow-y-auto custom-scrollbar">
            <div className="flex justify-between items-center pb-2 border-b border-slate-100">
               <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">Outbound Dispatch Queue</span>
               <span className="bg-primary/10 text-primary px-2 py-0.5 rounded text-[9px] font-black">{dispatchRequests.filter(r => r.status === 'STOCK_RESERVED').length} PENDING</span>
            </div>

            {dispatchRequests.map(req => (
              <div key={req.id} className="border border-slate-200 rounded-2xl p-4 shadow-sm hover:border-slate-300 transition-colors">
                 <div className="flex justify-between items-start mb-2">
                    <h4 className="font-bold text-slate-900 text-sm">Order #{req.order_id.slice(0, 8)}</h4>
                    <span className="text-[9px] font-black bg-amber-50 text-amber-700 px-2 py-0.5 rounded-md border border-amber-100">{req.status}</span>
                 </div>
                 <div className="flex flex-wrap items-center gap-2 mb-4 text-[10px] font-bold text-slate-500">
                    <MapPin className="w-3 h-3" />
                    <span>Store: {req.store_id}</span>
                 </div>

                 {req.status === 'STOCK_RESERVED' && (
                   <CasbinGuard obj={AUTH_RESOURCES.WAREHOUSE_DISPATCH} act={AUTH_ACTIONS.WRITE}>
                    <button onClick={() => handleDispatchClick('delivery', req.id)} className="w-full bg-primary text-white rounded-xl py-2.5 text-[10px] font-black uppercase tracking-widest flex items-center justify-center gap-2 hover:bg-primary/90 transition-all">
                      <Truck className="w-4 h-4" /> Dispatch Delivery
                    </button>
                   </CasbinGuard>
                 )}
                 {req.status === 'DISPATCHED' && (
                   <div className="bg-slate-50 rounded-xl p-3 border border-slate-100">
                     <p className="text-[10px] font-black uppercase tracking-widest text-primary mb-2 flex items-center gap-2"><MapPin className="w-3 h-3"/> Live Tracking</p>
                     <div className="w-full h-20 bg-slate-200 rounded-lg relative overflow-hidden flex items-center justify-center">
                        <div className="absolute inset-0" style={{ backgroundImage: 'radial-gradient(circle at center, #cbd5e1 1px, transparent 1px)', backgroundSize: '10px 10px' }} />
                        <div className="w-3 h-3 bg-primary rounded-full animate-pulse relative z-10 shadow-[0_0_8px_rgba(0,74,198,0.6)]" />
                     </div>
                   </div>
                 )}
              </div>
            ))}
            {dispatchRequests.length === 0 && (
                <div className="h-full flex flex-col items-center justify-center text-slate-300 py-12">
                   <Package className="w-12 h-12 opacity-20 mb-2" />
                   <p className="text-[10px] font-black uppercase italic">No pending dispatches.</p>
                </div>
            )}
          </div>
        </div>

      </div>

      {/* Footer Notifications Strip */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-4 mt-2">
         <div className="flex gap-8 overflow-x-auto custom-scrollbar items-center">
            {notifications.map((n, i) => (
              <div key={i} className="flex items-center gap-2 whitespace-nowrap min-w-0">
                 <div className={`w-2 h-2 rounded-full ${n.color === 'blue' ? 'bg-primary' : 'bg-green-500'}`} />
                 <span className="text-[10px] font-mono text-slate-400">{n.time}</span>
                 <span className="text-xs font-bold text-slate-700">{n.message}</span>
              </div>
            ))}
         </div>
      </div>

      {/* Assignment Modal */}
      {assignModal && (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
          <div className="bg-white w-full max-w-md rounded-[2rem] border border-slate-200 shadow-2xl p-8 relative">
            <button onClick={() => setAssignModal(null)} className="absolute top-6 right-6 text-slate-400 hover:text-slate-900">
              <X className="w-6 h-6" />
            </button>

            <h3 className="text-2xl font-black uppercase tracking-tighter text-slate-900 mb-6 flex items-center gap-2 italic">
              <UserPlus className="w-6 h-6 text-primary" />
              Assign <span className="text-primary">Fleet</span>
            </h3>

            <div className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Select Driver</label>
                <select
                  className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm font-bold outline-none focus:border-primary appearance-none"
                  value={assignmentData.driver_id}
                  onChange={e => setAssignmentData({...assignmentData, driver_id: e.target.value})}
                >
                  <option value="">Choose a driver...</option>
                  {drivers.map((d: Driver) => (
                    <option key={d.id} value={d.id}>{d.name} ({d.status})</option>
                  ))}
                </select>
              </div>

              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Select Vehicle</label>
                <select
                  className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm font-bold outline-none focus:border-primary appearance-none"
                  value={assignmentData.vehicle_id}
                  onChange={e => setAssignmentData({...assignmentData, vehicle_id: e.target.value})}
                >
                  <option value="">Choose a vehicle...</option>
                  {vehicles.map((v: Vehicle) => (
                    <option key={v.id} value={v.id}>{v.plate_number} - {v.type}</option>
                  ))}
                </select>
              </div>

              <div className="pt-4">
                <button
                  onClick={handleFinalizeDispatch}
                  disabled={!assignmentData.driver_id || !assignmentData.vehicle_id || dispatchMutation.isPending}
                  className="w-full bg-slate-900 text-white py-4 rounded-2xl font-black uppercase tracking-[0.2em] text-xs hover:bg-primary transition-all disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {dispatchMutation.isPending ? <><Loader2 className="w-4 h-4 animate-spin" /> Dispatching...</> : 'Finalize Dispatch'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  );
}
