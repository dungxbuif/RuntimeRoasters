'use client';

import React, { useState } from 'react';
import { adminService, CreateFarmRequest, User } from '@/services/admin.service';
import { farmService, Farm } from '@/services/farm.service';
import { testId, e2eSelectors } from '@/lib/utils/test-id';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/lib/auth';
import { FARM_LOCATIONS, COFFEE_TYPES } from '@/constants/domain';

export default function AdminFarmsPage() {
  const queryClient = useQueryClient();
  const { user: currentUser } = useAuth();
  const [showModal, setShowModal] = useState(false);
  const [editingFarm, setEditingFarm] = useState<Farm | null>(null);
  const [formData, setFormData] = useState<CreateFarmRequest>({
    name: '',
    location: 'CAU_DAT',
    area: 0,
    farm_type: 'ARABICA',
    owner_id: '',
  });

  // Queries
  const { data: farms = [], isLoading: isFarmsLoading } = useQuery({
    queryKey: ['farms'],
    queryFn: () => farmService.listFarms(),
  });

  const { data: managers = [], isLoading: isManagersLoading } = useQuery({
    queryKey: ['managers'],
    queryFn: async () => {
      const data = await adminService.listManagers();
      // Set default owner_id when managers load
      if (data.length > 0 && !formData.owner_id) {
        setFormData(prev => ({ ...prev, owner_id: data[0].id }));
      }
      return data;
    },
  });

  // Mutations
  const createFarmMutation = useMutation({
    mutationFn: (data: CreateFarmRequest) => adminService.createFarm(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['farms'] });
      closeModal();
    },
    onError: () => {
      alert('Failed to create farm. Ensure you have ownership permissions.');
    }
  });

  const updateFarmMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: CreateFarmRequest }) => farmService.updateFarm(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['farms'] });
      closeModal();
    },
    onError: () => {
      alert('Failed to update farm. Ensure you have ownership permissions.');
    }
  });

  const deleteFarmMutation = useMutation({
    mutationFn: (id: number) => farmService.deleteFarm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['farms'] });
    },
    onError: () => {
      alert('Failed to delete farm. Ensure you have ownership permissions.');
    }
  });

  const resetForm = () => {
    setEditingFarm(null);
    setFormData({
      name: '',
      location: 'CAU_DAT',
      area: 0,
      farm_type: 'ARABICA',
      owner_id: managers[0]?.id || '',
    });
  };

  const closeModal = () => {
    setShowModal(false);
    resetForm();
  };

  const openCreateModal = () => {
    resetForm();
    setShowModal(true);
  };

  const openEditModal = (farm: Farm) => {
    setEditingFarm(farm);
    setFormData({
      name: farm.name,
      location: farm.location,
      area: farm.area,
      farm_type: farm.farm_type,
      owner_id: farm.owner_id,
    });
    setShowModal(true);
  };

  const handleSubmitFarm = (e: React.FormEvent) => {
    e.preventDefault();
    if (editingFarm) {
      updateFarmMutation.mutate({ id: editingFarm.id, data: formData });
      return;
    }
    createFarmMutation.mutate(formData);
  };

  const handleDeleteFarm = (id: number) => {
    if (!confirm('Are you sure you want to decommission this farm node?')) return;
    deleteFarmMutation.mutate(id);
  };

  const isLoading = isFarmsLoading || isManagersLoading;
  const isManager = currentUser?.role === 'FARM_MANAGER';

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-4xl font-black font-headline text-on-surface uppercase italic tracking-tighter">
            Farm <span className="text-primary">Registry</span>
          </h1>
          <p className="text-on-surface-variant font-medium">Manage and provision origin coffee nodes.</p>
        </div>
        <button 
          onClick={openCreateModal}
          {...testId(e2eSelectors.CREATE_FARM_BTN)}
          className="bg-primary text-on-primary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-primary/20 flex items-center gap-2"
        >
          <span className="material-symbols-outlined !text-sm">add</span>
          Create New Farm
        </button>
      </div>

      <div className="bg-surface-container-lowest rounded-3xl border border-outline-variant/10 overflow-hidden shadow-sm">
        {isLoading ? (
          <div className="p-12 text-center text-on-surface-variant animate-pulse font-bold uppercase tracking-widest text-xs">Syncing Registry...</div>
        ) : (
          <table className="w-full text-left border-collapse" {...testId(e2eSelectors.FARM_LIST_TABLE)}>
            <thead className="bg-surface-container-low border-b border-outline-variant/10">
              <tr>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Farm Name</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Location</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Assigned Manager (Owner ID)</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/5">
              {farms.map((farm: Farm) => (
                <tr key={farm.id} className="hover:bg-surface-container-low/30 transition-colors" {...testId(`farm-row-${farm.name}`)}>
                  <td className="px-6 py-4">
                    <p className="font-black text-on-surface uppercase tracking-tight text-xs">{farm.name}</p>
                  </td>
                  <td className="px-6 py-4 text-[10px] font-medium text-on-surface-variant italic">
                    {FARM_LOCATIONS.find(l => l.code === farm.location)?.name || farm.location}
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <div className="w-5 h-5 bg-secondary-container/50 rounded flex items-center justify-center text-on-secondary-container">
                        <span className="material-symbols-outlined !text-[10px]">person</span>
                      </div>
                      <code className="text-[9px] text-on-surface-variant/40 font-mono">{farm.owner_id}</code>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex justify-end gap-2">
                      <button
                        onClick={() => openEditModal(farm)}
                        className="text-on-surface-variant hover:text-primary p-2 rounded-lg transition-colors"
                        {...testId(e2eSelectors.FARM_EDIT_BTN)}
                      >
                        <span className="material-symbols-outlined !text-sm">edit</span>
                      </button>
                      
                      {/* Only show delete if NOT a farm manager, or per business rule: managers can't delete */}
                      {!isManager && (
                        <button 
                          onClick={() => handleDeleteFarm(farm.id)}
                          className="text-on-surface-variant hover:text-error p-2 rounded-lg transition-colors"
                          {...testId(e2eSelectors.FARM_DELETE_BTN)}
                        >
                          <span className="material-symbols-outlined !text-sm">delete</span>
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-slate-950/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4" {...testId(e2eSelectors.FARM_MODAL)}>
          <div className="bg-surface-container-lowest w-full max-w-md rounded-[2rem] border border-outline-variant/20 shadow-2xl p-8">
            <h2 className="text-2xl font-black font-headline uppercase italic tracking-tighter mb-6">
              {editingFarm ? (
                <>
                  Update <span className="text-primary">Farm</span>
                </>
              ) : (
                <>
                  Provision <span className="text-primary">New Farm</span>
                </>
              )}
            </h2>
            <form onSubmit={handleSubmitFarm} className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Farm Name</label>
                <input 
                  required
                  {...testId(e2eSelectors.FARM_NAME_INPUT)}
                  placeholder="e.g. Arabica Highlands Alpha"
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Area (Hectares)</label>
                  <input 
                    type="number"
                    required
                    step="0.1"
                    {...testId(e2eSelectors.FARM_AREA_INPUT)}
                    className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors"
                    value={formData.area}
                    onChange={(e) => setFormData({...formData, area: parseFloat(e.target.value)})}
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Coffee Type</label>
                  <select 
                    className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors"
                    value={formData.farm_type}
                    onChange={(e) => setFormData({...formData, farm_type: e.target.value})}
                  >
                    {COFFEE_TYPES.map(ct => (
                      <option key={ct.code} value={ct.code}>{ct.name}</option>
                    ))}
                  </select>
                </div>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Regional Location</label>
                <select 
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors"
                  value={formData.location}
                  onChange={(e) => setFormData({...formData, location: e.target.value})}
                >
                  {FARM_LOCATIONS.map(loc => (
                    <option key={loc.code} value={loc.code}>{loc.name}</option>
                  ))}
                </select>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Assign Manager</label>
                <select 
                  required
                  {...testId(e2eSelectors.FARM_OWNER_SELECT)}
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors font-bold"
                  value={formData.owner_id}
                  onChange={(e) => setFormData({...formData, owner_id: e.target.value})}
                >
                  {managers.map((m: User) => (
                    <option key={m.id} value={m.id}>{m.name} ({m.email})</option>
                  ))}
                  {managers.length === 0 && <option disabled>No Managers Available</option>}
                </select>
              </div>
              <div className="flex gap-4 pt-4">
                <button 
                  type="button"
                  onClick={closeModal}
                  className="flex-1 px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest border border-outline-variant/20 hover:bg-surface-container-low transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  disabled={createFarmMutation.isPending || updateFarmMutation.isPending}
                  {...testId(e2eSelectors.FARM_SUBMIT_BTN)}
                  className="flex-1 bg-primary text-on-primary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-primary/20 disabled:opacity-50"
                >
                  {editingFarm
                    ? (updateFarmMutation.isPending ? 'Updating...' : 'Update Farm')
                    : (createFarmMutation.isPending ? 'Provisioning...' : 'Provision Farm')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
