'use client';

import { NotificationItem } from '@/components/common/NotificationFeed';
import StatusPipeline from '@/components/common/StatusPipeline';
import { DashboardLayout } from '@/components/ui/templates/DashboardLayout';
import { logisticsService } from '@/services/logistics.service';
import { warehouseService } from '@/services/warehouse.service';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { CheckCircle2, MapPin, Truck, Warehouse, UserPlus, X, Loader2 } from 'lucide-react';
import { useMemo, useState } from 'react';
import { AUTH_ACTIONS, AUTH_RESOURCES } from '@/constants/resources';
import { CasbinGuard } from '@/lib/auth';

export default function WarehouseOperationsPage() {
  const queryClient = useQueryClient();
  const [selectedWarehouseId, setSelectedWarehouseId] = useState<string>('wh-hn-001');
  const [assignModal, setAssignModal] = useState<{ type: 'pickup' | 'delivery', id: string } | null>(null);
  const [assignmentData, setAssignmentData] = useState({ driver_id: '', vehicle_id: '' });

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

  const dispatchMutation = useMutation({
    mutationFn: (args: { type: 'pickup' | 'delivery', id: string, driver_id: string, vehicle_id: string }) => {
      if (args.type === 'pickup') return warehouseService.dispatchPickupRequest(args.id, args.driver_id, args.vehicle_id);
      return warehouseService.dispatchRequest(args.id);
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

  const activeProcessingBatches = useMemo(() => batches.filter(b => b.status !== 'STOCKED'), [batches]);

  const notifications: NotificationItem[] = useMemo(() => [
    { id: '1', icon: 'check_circle', message: `K'Ho Coffee Farm • 500kg Arabica • Store Q1`, time: '13:49', color: 'blue' },
    { id: '2', icon: 'local_shipping', message: 'Cau Dat Arabica • 800kg Typica • Carder Alpha', time: '20:26', color: 'blue' },
    { id: '3', icon: 'payments', message: `K'Ho Coffee Farm • 500kg Hoan Kiem`, time: '20:07', color: 'blue' },
    { id: '4', icon: 'payments', message: `Aeroco Coffee • 1200kg Robusta • Store Hoan Kiem`, time: '08:09', color: 'green' },
  ], []);

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

  return (
    <DashboardLayout
      title="Warehouse Operations"
      subtitle="Inventory Management & Dispatch Control"
      icon={Warehouse}
      actions={
        <div className="flex flex-col items-end gap-1">
          <label className="text-[9px] font-black uppercase tracking-widest text-slate-400">Warehouse</label>
          <select
            value={selectedWarehouseId}
            onChange={e => setSelectedWarehouseId(e.target.value)}
            className="bg-white border border-slate-200 rounded-lg px-4 py-2 text-xs font-bold uppercase tracking-tight appearance-none cursor-pointer shadow-sm outline-none"
          >
            <option value="wh-hn-001">Warehouse Hanoi (HN-001)</option>
            <option value="wh-hcm-001">Warehouse HCM (HCM-001)</option>
          </select>
        </div>
      }
    >
      {/* 3-Column Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 flex-1 min-h-0 w-full">
        
        {/* LEFT: INBOUND QUEUE */}
        <div className="flex flex-col gap-4">
          <h2 className="text-xl font-black text-slate-900 flex items-center justify-between">
            Inbound Queue
          </h2>
          <div className="bg-white rounded-[2rem] border border-slate-200 shadow-sm flex-1 p-5 flex flex-col gap-4 overflow-y-auto custom-scrollbar">
            <div className="flex justify-between items-center pb-2 border-b border-slate-100">
               <span className="text-[10px] font-black uppercase tracking-widest text-slate-500">Inbound Pickup Requests</span>
               <span className="bg-amber-500 text-white px-2 py-0.5 rounded text-[9px] font-black">3 NEW</span>
            </div>
            
            {pickupRequests.map(req => (
              <div key={req.id} className="border border-slate-200 rounded-2xl p-4 shadow-sm hover:border-slate-300 transition-colors">
                 <h4 className="font-bold text-slate-900 text-sm mb-2">{req.origin_code} • {req.quantity}kg {req.coffee_type}</h4>
                 <div className="flex flex-wrap items-center gap-2 mb-4">
                    <span className={`text-[9px] font-black px-2 py-1 rounded-full uppercase bg-amber-100 text-amber-700`}>{req.status.replace('_', ' ')}</span>
                 </div>
                 
                 {req.status === 'CREATED' && (
                   <CasbinGuard obj={AUTH_RESOURCES.WAREHOUSE_DISPATCH} act={AUTH_ACTIONS.WRITE}>
                    <button onClick={() => handleDispatchClick('pickup', req.id)} className="w-full bg-primary text-white rounded-xl py-2.5 text-[10px] font-black uppercase tracking-widest flex items-center justify-center gap-2 hover:bg-primary/90 transition-all">
                      <Truck className="w-4 h-4" /> Dispatch Pickup
                    </button>
                   </CasbinGuard>
                 )}
                 {(req.status === 'DISPATCHED' || req.status === 'LOADED') && (
                   <div>
                     <p className="text-xs text-slate-600 mb-2 font-medium flex items-center gap-2"><Truck className="w-3 h-3"/> En Route</p>
                     <div className="w-full bg-slate-100 h-1.5 rounded-full overflow-hidden"><div className="bg-blue-500 h-full w-[45%]" /></div>
                   </div>
                 )}
                 {(req.status === 'LOADED' || req.status === 'DISPATCHED') && (
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
               <span className="bg-red-500 text-white px-2 py-0.5 rounded text-[9px] font-black">2 PENDING</span>
            </div>
            
            {dispatchRequests.map(req => (
              <div key={req.id} className="border border-slate-200 rounded-2xl p-4 shadow-sm hover:border-slate-300 transition-colors">
                 <h4 className="font-bold text-slate-900 text-sm mb-2">Order #{req.id.slice(-6)}</h4>
                 <div className="flex flex-wrap items-center gap-2 mb-4">
                    <span className={`text-[9px] font-black px-2 py-1 rounded-full uppercase bg-slate-200 text-slate-700`}>{req.status.replace('_', ' ')}</span>
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
                  {drivers.map((d: any) => (
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
                  {vehicles.map((v: any) => (
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
