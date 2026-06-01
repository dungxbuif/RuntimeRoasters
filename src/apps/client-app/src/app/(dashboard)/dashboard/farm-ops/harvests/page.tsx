'use client';

import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { farmService, Harvest, CreateHarvestInput, CoffeeType } from '@/services/farm.service';
import { COFFEE_TYPES, HARVEST_STATUSES } from '@/constants/domain';
import { e2eSelectors, testId } from '@/lib/utils/test-id';
import { CasbinGuard } from '@/lib/auth';
import StatusPipeline, { PipelineStep } from '@/components/common/StatusPipeline';
import NotificationFeed from '@/components/common/NotificationFeed';

const PICKUP_STEPS: PipelineStep[] = [
  { key: 'CREATED', label: 'Created' },
  { key: 'PICKUP_REQUESTED', label: 'Pickup Req' },
  { key: 'PICKUP_ASSIGNED', label: 'Assigned' },
  { key: 'PICKED_UP', label: 'Picked Up' },
  { key: 'ARRIVED_WAREHOUSE', label: 'At WH' },
  { key: 'INTAKE_CREATED', label: 'Intake' },
];

function getPickupStatus(harvestStatus: string): string {
  const mapping: Record<string, string> = {
    'NEW': 'CREATED',
    'PENDING': 'PICKUP_REQUESTED',
    'ASSIGNED': 'PICKUP_ASSIGNED',
    'PICKED_UP': 'PICKED_UP',
    'ARRIVED_WAREHOUSE': 'ARRIVED_WAREHOUSE',
    'INTAKE_CREATED': 'INTAKE_CREATED',
    'COMPLETED': 'INTAKE_CREATED',
  };
  return mapping[harvestStatus] || 'CREATED';
}

type ApiError = {
  response?: { data?: { message?: string } };
  message?: string;
};

export default function HarvestsPage() {
  const queryClient = useQueryClient();
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState<CreateHarvestInput>({
    farm_id: 0,
    coffee_type: 'ARABICA',
    quantity: 0,
    harvest_date: new Date().toISOString().split('T')[0],
  });

  // Queries
  const { data: harvests = [], isLoading: isHarvestsLoading } = useQuery({
    queryKey: ['harvests'],
    queryFn: () => farmService.listHarvests(),
  });

  const { data: farms = [], isLoading: isFarmsLoading } = useQuery({
    queryKey: ['farms'],
    queryFn: async () => {
      const data = await farmService.listFarms();
      if (data.length > 0 && !formData.farm_id) {
        setFormData(prev => ({ ...prev, farm_id: data[0].id }));
      }
      return data;
    },
  });

  // Mutations
  const createHarvestMutation = useMutation({
    mutationFn: (data: CreateHarvestInput) => farmService.createHarvest(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['harvests'] });
      setShowModal(false);
      setFormData({
        farm_id: farms[0]?.id || 0,
        coffee_type: 'ARABICA',
        quantity: 0,
        harvest_date: new Date().toISOString().split('T')[0],
      });
    },
    onError: (error: ApiError) => {
      alert(`Failed to declare harvest: ${error.response?.data?.message || error.message}`);
    }
  });

  const handleCreateHarvest = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.quantity <= 0) {
      alert('Quantity must be greater than 0');
      return;
    }
    createHarvestMutation.mutate(formData);
  };

  const isLoading = isHarvestsLoading || isFarmsLoading;

  return (
    <CasbinGuard obj={AUTH_RESOURCES.HARVEST} act={AUTH_ACTIONS.READ} fallback={<div className="p-12 text-center font-black uppercase italic text-slate-400">Access Denied: Farm personnel only.</div>}>
      <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-4xl font-black font-headline text-on-surface uppercase italic tracking-tighter">
            Harvest <span className="text-tertiary">Declaration</span>
          </h1>
          <p className="text-on-surface-variant font-medium">Record harvests and broadcast to <span className="text-primary font-black italic">Warehouse Service</span> via Kafka.</p>
        </div>
        <CasbinGuard obj={AUTH_RESOURCES.HARVEST} act={AUTH_ACTIONS.WRITE}>
          <button 
            onClick={() => setShowModal(true)}
            {...testId(e2eSelectors.CREATE_HARVEST_BTN)}
            className="bg-tertiary text-on-tertiary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-tertiary/20 flex items-center gap-2"
          >
            <span className="material-symbols-outlined !text-sm">add_circle</span>
            Declare New Harvest
          </button>
        </CasbinGuard>
      </div>

      <div className="bg-surface-container-lowest rounded-3xl border border-outline-variant/10 overflow-hidden shadow-sm">
        {isLoading ? (
          <div className="p-12 text-center text-on-surface-variant animate-pulse font-bold uppercase tracking-widest text-xs">Syncing Harvest Ledger...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse" {...testId(e2eSelectors.HARVEST_LIST_TABLE)}>
              <thead className="bg-surface-container-low border-b border-outline-variant/10">
                <tr>
                  <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Harvest ID</th>
                  <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Farm</th>
                  <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Type</th>
                  <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Quantity</th>
                  <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Date</th>
                  <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/5">
                {harvests.map((harvest: Harvest) => (
                  <tr key={harvest.id} className="hover:bg-surface-container-low/30 transition-colors">
                    <td className="px-6 py-4">
                      <code className="text-[10px] text-on-surface-variant/60 font-mono">#{harvest.id.toString().padStart(4, '0')}</code>
                    </td>
                    <td className="px-6 py-4">
                      <p className="font-black text-on-surface uppercase tracking-tight text-xs">
                        {farms.find(f => f.id === harvest.farm_id)?.name || 'Unknown Farm'}
                      </p>
                      <p className="text-[9px] text-on-surface-variant opacity-50 font-mono">{harvest.farm_id}</p>
                    </td>
                    <td className="px-6 py-4">
                      <span className="px-2 py-1 rounded-md bg-secondary-container/30 text-on-secondary-container text-[10px] font-bold uppercase italic">
                        {harvest.coffee_type}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <p className="font-black text-on-surface text-xs italic">{harvest.quantity} KG</p>
                    </td>
                    <td className="px-6 py-4 text-[10px] font-medium text-on-surface-variant italic">
                      {new Date(harvest.harvest_date).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-1.5 mb-2">
                        <div className={`w-2 h-2 rounded-full ${
                          HARVEST_STATUSES[harvest.status as keyof typeof HARVEST_STATUSES]?.color || 'bg-slate-300'
                        } ${harvest.status === 'NEW' ? 'animate-pulse' : ''}`}></div>
                        <span className="text-[10px] font-black uppercase tracking-widest text-on-surface-variant italic">
                          {HARVEST_STATUSES[harvest.status as keyof typeof HARVEST_STATUSES]?.label || harvest.status}
                        </span>
                      </div>
                      <div className="max-w-[300px]">
                        <StatusPipeline steps={PICKUP_STEPS} currentStep={getPickupStatus(harvest.status)} size="sm" />
                      </div>
                    </td>
                  </tr>
                ))}
                {harvests.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-6 py-12 text-center text-on-surface-variant font-bold uppercase tracking-widest text-[10px] italic">
                      No harvest records found in the ledger.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Pickup Notifications */}
      <div className="max-w-md">
        <NotificationFeed
          title="Pickup Activity"
          items={[
            { id: '1', icon: 'truck', message: 'Driver assigned to latest harvest', time: '5min ago', color: 'blue' },
            { id: '2', icon: 'warehouse', message: 'Warehouse created pickup request', time: '15min ago', color: 'amber' },
            { id: '3', icon: 'check', message: 'Previous harvest intake completed', time: '1h ago', color: 'green' },
          ]}
        />
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-slate-950/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4" {...testId(e2eSelectors.HARVEST_MODAL)}>
          <div className="bg-surface-container-lowest w-full max-w-md rounded-[2rem] border border-outline-variant/20 shadow-2xl p-8 relative overflow-hidden">
            {/* Background Accent */}
            <div className="absolute -top-24 -right-24 w-48 h-48 bg-tertiary/10 rounded-full blur-3xl"></div>
            
            <h2 className="text-2xl font-black font-headline uppercase italic tracking-tighter mb-6 relative z-10">
              New <span className="text-tertiary">Harvest Entry</span>
            </h2>
            
            <form onSubmit={handleCreateHarvest} className="space-y-4 relative z-10">
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Origin Farm</label>
                <select 
                  required
                  {...testId(e2eSelectors.HARVEST_FARM_SELECT)}
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors font-bold uppercase tracking-tight"
                  value={formData.farm_id || ''}
                  onChange={(e) => setFormData({...formData, farm_id: Number(e.target.value)})}
                >
                  {farms.map(farm => (
                    <option key={farm.id} value={farm.id}>{farm.name}</option>
                  ))}
                  {farms.length === 0 && <option disabled>No Farms Available</option>}
                </select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Coffee Variety</label>
                <select 
                  {...testId(e2eSelectors.HARVEST_TYPE_SELECT)}
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors font-bold"
                  value={formData.coffee_type}
                    onChange={(e) => setFormData({...formData, coffee_type: e.target.value as CoffeeType})}
                  >
                    {COFFEE_TYPES.map(ct => (
                      <option key={ct.code} value={ct.code}>{ct.name}</option>
                    ))}
                  </select>
                </div>
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Yield (KG)</label>
                <input 
                  type="number"
                  required
                  step="0.1"
                  min="0.1"
                  placeholder="0.0"
                  {...testId(e2eSelectors.HARVEST_QUANTITY_INPUT)}
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors font-bold italic"
                  value={formData.quantity || ''}
                    onChange={(e) => setFormData({...formData, quantity: parseFloat(e.target.value)})}
                  />
                </div>
              </div>

              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Harvest Date</label>
                <input 
                  type="date"
                  required
                  {...testId(e2eSelectors.HARVEST_DATE_INPUT)}
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors font-bold"
                  value={formData.harvest_date}
                  onChange={(e) => setFormData({...formData, harvest_date: e.target.value})}
                />
              </div>

              <div className="flex gap-4 pt-4">
                <button 
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest border border-outline-variant/20 hover:bg-surface-container-low transition-colors"
                >
                  Discard
                </button>
                <button 
                  type="submit"
                  disabled={createHarvestMutation.isPending}
                  {...testId(e2eSelectors.HARVEST_SUBMIT_BTN)}
                  className="flex-1 bg-tertiary text-on-tertiary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-tertiary/20 disabled:opacity-50 flex items-center justify-center gap-2"
                >
                  {createHarvestMutation.isPending ? (
                    <>
                      <span className="w-4 h-4 border-2 border-on-tertiary/30 border-t-on-tertiary rounded-full animate-spin"></span>
                      Recording...
                    </>
                  ) : 'Record Harvest'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
    </CasbinGuard>
  );
}
