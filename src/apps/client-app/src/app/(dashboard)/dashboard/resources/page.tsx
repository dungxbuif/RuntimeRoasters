'use client';

import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminService, CreateWarehouseRequest, CreateStoreRequest, User } from '@/services/admin.service';
import { Warehouse, Store, Truck, Users, Plus, Settings, UserPlus, X, Loader2 } from 'lucide-react';

type TabKey = 'warehouses' | 'stores' | 'vehicles' | 'drivers';

const TABS: { key: TabKey; label: string; icon: React.ElementType }[] = [
  { key: 'warehouses', label: 'Warehouses', icon: Warehouse },
  { key: 'stores', label: 'Stores', icon: Store },
  { key: 'vehicles', label: 'Vehicles', icon: Truck },
  { key: 'drivers', label: 'Drivers', icon: Users },
];

const STATUS_COLORS: Record<string, string> = {
  ACTIVE: 'bg-green-50 text-green-700 border-green-200',
  PENDING: 'bg-amber-50 text-amber-700 border-amber-200',
  IDLE: 'bg-slate-100 text-slate-600 border-slate-200',
  ASSIGNED: 'bg-blue-50 text-blue-700 border-blue-200',
  EN_ROUTE: 'bg-primary/10 text-primary border-primary/20',
  AVAILABLE: 'bg-green-50 text-green-700 border-green-200',
  DRIVING: 'bg-primary/10 text-primary border-primary/20',
  OFFLINE: 'bg-slate-100 text-slate-400 border-slate-200',
};

export default function ResourceManagementPage() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<TabKey>('warehouses');
  const [showAssignModal, setShowAssignModal] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);

  // Queries
  const { data: users = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => adminService.listUsers(),
  });

  const { data: warehouses = [], isLoading: isWHLoading } = useQuery({
    queryKey: ['warehouses'],
    queryFn: () => adminService.listWarehouses(),
    enabled: activeTab === 'warehouses',
  });

  const { data: stores = [], isLoading: isStoresLoading } = useQuery({
    queryKey: ['stores'],
    queryFn: () => adminService.listStores(),
    enabled: activeTab === 'stores',
  });

  const { data: managers = [] } = useQuery({
    queryKey: ['managers'],
    queryFn: () => adminService.listManagers(),
  });

  // Mutations
  const createWHMutation = useMutation({
    mutationFn: (data: CreateWarehouseRequest) => adminService.createWarehouse(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['warehouses'] });
      setShowCreateModal(false);
    }
  });

  const createStoreMutation = useMutation({
    mutationFn: (data: CreateStoreRequest) => adminService.createStore(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stores'] });
      setShowCreateModal(false);
    }
  });

  const StatusBadge = ({ status }: { status: string }) => (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-[9px] font-black uppercase border ${STATUS_COLORS[status] || STATUS_COLORS.PENDING}`}>
      {status}
    </span>
  );

  const ActionLink = ({ onClick, children }: { onClick: () => void; children: React.ReactNode }) => (
    <button onClick={onClick} className="text-[10px] font-black text-primary hover:text-primary/70 uppercase tracking-tight transition-colors">
      {children}
    </button>
  );

  return (
    <div className="h-full flex flex-col gap-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <Settings className="w-7 h-7 text-primary" />
            Resource Management
          </h1>
          <p className="text-[10px] font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Warehouses, Stores, Vehicles & Drivers
          </p>
        </div>
        <button 
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 bg-primary text-white px-5 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-primary/90 transition-all shadow-lg shadow-primary/20"
        >
          <Plus className="w-4 h-4" />
          Create {activeTab.slice(0, -1)}
        </button>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 bg-slate-100 rounded-xl p-1">
        {TABS.map(tab => (
          <button key={tab.key} onClick={() => setActiveTab(tab.key)}
            className={`flex-1 flex items-center justify-center gap-2 py-2.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all ${
              activeTab === tab.key
                ? 'bg-white text-slate-900 shadow-sm'
                : 'text-slate-400 hover:text-slate-600'
            }`}>
            <tab.icon className="w-4 h-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="flex-1 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden flex flex-col min-h-0">
        <div className="flex-1 overflow-x-auto">
          {isWHLoading || isStoresLoading ? (
            <div className="h-full flex items-center justify-center text-[10px] font-black uppercase tracking-widest text-slate-400 animate-pulse italic">
              Syncing Ledger...
            </div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-slate-50/50 text-[9px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-100">
                  {activeTab === 'warehouses' && (<><th className="px-6 py-3">Name</th><th className="px-6 py-3">Location</th><th className="px-6 py-3">Manager</th><th className="px-6 py-3">Capacity</th><th className="px-6 py-3">Status</th><th className="px-6 py-3 text-right">Actions</th></>)}
                  {activeTab === 'stores' && (<><th className="px-6 py-3">Name</th><th className="px-6 py-3">Location</th><th className="px-6 py-3">Manager</th><th className="px-6 py-3">Status</th><th className="px-6 py-3 text-right">Actions</th></>)}
                  {activeTab === 'vehicles' && (<><th className="px-6 py-3">Plate</th><th className="px-6 py-3">Type</th><th className="px-6 py-3">Capacity</th><th className="px-6 py-3">Status</th><th className="px-6 py-3 text-right">Actions</th></>)}
                  {activeTab === 'drivers' && (<><th className="px-6 py-3">Name</th><th className="px-6 py-3">Email</th><th className="px-6 py-3">Vehicle</th><th className="px-6 py-3">Status</th><th className="px-6 py-3 text-right">Actions</th></>)}
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-50">
                {activeTab === 'warehouses' && warehouses.map(wh => (
                  <tr key={wh.id} className="hover:bg-slate-50/50 transition-colors">
                    <td className="px-6 py-4"><span className="text-xs font-black text-slate-900">{wh.name}</span><br/><span className="text-[9px] font-mono text-slate-400">{wh.code}</span></td>
                    <td className="px-6 py-4 text-xs font-bold text-slate-600">{wh.location}</td>
                    <td className="px-6 py-4 text-[10px] font-bold text-slate-500">{wh.manager_email || 'Unassigned'}</td>
                    <td className="px-6 py-4 text-xs font-bold text-slate-600">{wh.capacity}</td>
                    <td className="px-6 py-4"><StatusBadge status={wh.status} /></td>
                    <td className="px-6 py-4 text-right space-x-3">
                      <ActionLink onClick={() => {}}>Edit</ActionLink>
                      <ActionLink onClick={() => setShowAssignModal(wh.id)}>Assign</ActionLink>
                    </td>
                  </tr>
                ))}
                {activeTab === 'stores' && stores.map(st => (
                  <tr key={st.id} className="hover:bg-slate-50/50 transition-colors">
                    <td className="px-6 py-4"><span className="text-xs font-black text-slate-900">{st.name}</span><br/><span className="text-[9px] font-mono text-slate-400">{st.city}</span></td>
                    <td className="px-6 py-4 text-xs font-bold text-slate-600">{st.address}</td>
                    <td className="px-6 py-4 text-[10px] font-bold text-slate-500">{st.manager_email || 'Unassigned'}</td>
                    <td className="px-6 py-4"><StatusBadge status={st.status} /></td>
                    <td className="px-6 py-4 text-right space-x-3">
                      <ActionLink onClick={() => {}}>Edit</ActionLink>
                      <ActionLink onClick={() => setShowAssignModal(st.id)}>Assign</ActionLink>
                    </td>
                  </tr>
                ))}
                {activeTab === 'warehouses' && warehouses.length === 0 && (
                  <tr><td colSpan={6} className="px-6 py-12 text-center text-[10px] font-bold uppercase tracking-widest text-slate-400 italic">No warehouses registered.</td></tr>
                )}
                {activeTab === 'stores' && stores.length === 0 && (
                  <tr><td colSpan={5} className="px-6 py-12 text-center text-[10px] font-bold uppercase tracking-widest text-slate-400 italic">No stores registered.</td></tr>
                )}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {/* Create Modal */}
      {showCreateModal && (
        <CreateResourceModal 
          type={activeTab} 
          onClose={() => setShowCreateModal(false)}
          onWHSubmit={(data: CreateWarehouseRequest) => createWHMutation.mutate(data)}
          onStoreSubmit={(data: CreateStoreRequest) => createStoreMutation.mutate(data)}
          managers={managers}
          isPending={createWHMutation.isPending || createStoreMutation.isPending}
        />
      )}

      {/* Assign Modal (existing) */}
      {showAssignModal && (
        <div className="fixed inset-0 bg-black/40 z-50 flex items-center justify-center" onClick={() => setShowAssignModal(null)}>
          <div className="bg-white rounded-3xl p-8 w-[420px] shadow-2xl" onClick={e => e.stopPropagation()}>
            <h3 className="text-lg font-black uppercase tracking-tighter text-slate-900 mb-6 flex items-center gap-2">
              <UserPlus className="w-5 h-5 text-primary" />
              Quick Assign
            </h3>
            <div className="space-y-4">
              <div>
                <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Select User</label>
                <select className="w-full bg-slate-50 border border-slate-200 rounded-xl py-3 px-4 text-sm font-bold focus:border-primary outline-none">
                  <option value="">Select a manager...</option>
                  {users.map((u: User) => (
                    <option key={u.id} value={u.id}>{u.email} ({u.role})</option>
                  ))}
                </select>
              </div>
              <div className="flex gap-3 pt-4">
                <button onClick={() => setShowAssignModal(null)}
                  className="flex-1 py-3 bg-slate-100 text-slate-500 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all">
                  Cancel
                </button>
                <button onClick={() => setShowAssignModal(null)}
                  className="flex-1 py-3 bg-primary text-white rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-primary/90 transition-all shadow-lg">
                  Confirm Assignment
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

interface CreateResourceModalProps {
  type: TabKey;
  onClose: () => void;
  onWHSubmit: (data: CreateWarehouseRequest) => void;
  onStoreSubmit: (data: CreateStoreRequest) => void;
  managers: User[];
  isPending: boolean;
}

function CreateResourceModal({ type, onClose, onWHSubmit, onStoreSubmit, managers, isPending }: CreateResourceModalProps) {
  const [formData, setFormData] = useState({
    name: '',
    code: '',
    location: '',
    city: '',
    address: '',
    manager_id: '',
    capacity: '5000 Tons',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const manager = managers.find((m: User) => m.id === formData.manager_id);
    if (type === 'warehouses') {
      onWHSubmit({
        name: formData.name,
        code: formData.code,
        location: formData.location,
        manager_id: formData.manager_id,
        manager_email: manager?.email || '',
        capacity: formData.capacity,
      });
    } else if (type === 'stores') {
      onStoreSubmit({
        name: formData.name,
        city: formData.city,
        address: formData.address,
        manager_id: formData.manager_id,
        manager_email: manager?.email || '',
      });
    }
  };

  return (
    <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[60] flex items-center justify-center p-4">
      <div className="bg-white w-full max-w-lg rounded-[2.5rem] border border-slate-200 shadow-2xl p-10 relative overflow-hidden">
        <button onClick={onClose} className="absolute top-8 right-8 text-slate-400 hover:text-slate-900 transition-colors">
          <X className="w-6 h-6" />
        </button>

        <h2 className="text-3xl font-black uppercase tracking-tighter italic text-slate-900 mb-8">
          Register <span className="text-primary">{type.slice(0, -1)}</span>
        </h2>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-1">
            <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Entity Name</label>
            <input required type="text" placeholder="e.g. South Hub A" 
              className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none transition-all"
              value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})}
            />
          </div>

          {type === 'warehouses' ? (
            <>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Code</label>
                  <input required type="text" placeholder="WH-001" 
                    className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none uppercase font-mono"
                    value={formData.code} onChange={e => setFormData({...formData, code: e.target.value})}
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Capacity</label>
                  <input required type="text" placeholder="5000 Tons" 
                    className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none"
                    value={formData.capacity} onChange={e => setFormData({...formData, capacity: e.target.value})}
                  />
                </div>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Address</label>
                <input required type="text" placeholder="District 7, HCM City" 
                  className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none"
                  value={formData.location} onChange={e => setFormData({...formData, location: e.target.value})}
                />
              </div>
            </>
          ) : (
            <>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">City</label>
                  <input required type="text" placeholder="Hanoi" 
                    className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none"
                    value={formData.city} onChange={e => setFormData({...formData, city: e.target.value})}
                  />
                </div>
                <div className="space-y-1">
                   {/* Empty for spacing */}
                </div>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Full Address</label>
                <input required type="text" placeholder="123 Ly Thai To, Hoan Kiem" 
                  className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none"
                  value={formData.address} onChange={e => setFormData({...formData, address: e.target.value})}
                />
              </div>
            </>
          )}

          <div className="space-y-1">
            <label className="text-[10px] font-black uppercase tracking-widest ml-1 text-slate-400">Assigned Manager</label>
            <select required 
              className="w-full bg-slate-50 border border-slate-200 rounded-2xl px-5 py-4 text-sm font-bold focus:border-primary outline-none appearance-none"
              value={formData.manager_id} onChange={e => setFormData({...formData, manager_id: e.target.value})}
            >
              <option value="">Select personnel...</option>
              {managers.map((m: User) => (
                <option key={m.id} value={m.id}>{m.email} ({m.role})</option>
              ))}
            </select>
          </div>

          <div className="pt-6">
            <button 
              type="submit"
              disabled={isPending}
              className="w-full bg-slate-900 text-white py-5 rounded-2xl font-black uppercase tracking-[0.2em] text-xs hover:bg-primary transition-all shadow-2xl flex items-center justify-center gap-3 disabled:opacity-50"
            >
              {isPending ? (
                <><Loader2 className="w-4 h-4 animate-spin" /> Provisioning...</>
              ) : (
                'Finalize Registration'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
