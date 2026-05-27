'use client';

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { adminService } from '@/services/admin.service';
import { Warehouse, Store, Truck, Users, Plus, Settings, UserPlus } from 'lucide-react';

type TabKey = 'warehouses' | 'stores' | 'vehicles' | 'drivers';

const TABS: { key: TabKey; label: string; icon: React.ElementType }[] = [
  { key: 'warehouses', label: 'Warehouses', icon: Warehouse },
  { key: 'stores', label: 'Stores', icon: Store },
  { key: 'vehicles', label: 'Vehicles', icon: Truck },
  { key: 'drivers', label: 'Drivers', icon: Users },
];

// Demo data — replace with real API calls
const DEMO_WAREHOUSES = [
  { id: 'wh-hn-001', name: 'Warehouse Hanoi', code: 'HN-001', location: 'Hòa Lạc, Hanoi', manager: 'warehouse.hn@runtimeroasters.com', capacity: '5000 tons', status: 'ACTIVE' },
  { id: 'wh-hcm-001', name: 'Warehouse HCM', code: 'HCM-001', location: 'Sóng Thần, HCM', manager: 'Unassigned', capacity: '8000 tons', status: 'PENDING' },
  { id: 'wh-dn-001', name: 'Warehouse Da Nang', code: 'DN-001', location: 'Hòa Khánh, Đà Nẵng', manager: 'warehouse.dn@runtimeroasters.com', capacity: '3000 tons', status: 'ACTIVE' },
];

const DEMO_STORES = [
  { id: 'st-01', name: 'Runtime Roasters Hoàn Kiếm', location: 'Hanoi', manager: 'mgr.hn.hoankiem@runtimeroasters.com', status: 'ACTIVE' },
  { id: 'st-02', name: 'Runtime Roasters Quận 1', location: 'HCM', manager: 'mgr.hcm.q1@runtimeroasters.com', status: 'ACTIVE' },
  { id: 'st-03', name: 'Runtime Roasters Hải Châu', location: 'Đà Nẵng', manager: 'mgr.dn.haichau@runtimeroasters.com', status: 'ACTIVE' },
];

const DEMO_VEHICLES = [
  { id: 'v-01', plate: '30A-12345', type: 'Refrigerated Truck', capacity: '2 tons', status: 'IDLE' },
  { id: 'v-02', plate: '51B-67890', type: 'Van', capacity: '500 kg', status: 'ASSIGNED' },
  { id: 'v-03', plate: '43C-11111', type: 'Pickup', capacity: '1 ton', status: 'EN_ROUTE' },
];

const DEMO_DRIVERS = [
  { id: 'd-01', name: 'Driver Alpha', email: 'driver@runtimeroasters.com', vehicle: '30A-12345', status: 'AVAILABLE' },
  { id: 'd-02', name: 'Driver Beta', email: 'driver.beta@runtimeroasters.com', vehicle: 'Unassigned', status: 'DRIVING' },
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
  const [activeTab, setActiveTab] = useState<TabKey>('warehouses');
  const [showAssignModal, setShowAssignModal] = useState<string | null>(null);

  const { data: users = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => adminService.listUsers(),
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
        <button className="flex items-center gap-2 bg-primary text-white px-5 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-primary/90 transition-all shadow-lg shadow-primary/20">
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
              {activeTab === 'warehouses' && DEMO_WAREHOUSES.map(wh => (
                <tr key={wh.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4"><span className="text-xs font-black text-slate-900">{wh.name}</span><br/><span className="text-[9px] font-mono text-slate-400">{wh.code}</span></td>
                  <td className="px-6 py-4 text-xs font-bold text-slate-600">{wh.location}</td>
                  <td className="px-6 py-4 text-[10px] font-bold text-slate-500">{wh.manager}</td>
                  <td className="px-6 py-4 text-xs font-bold text-slate-600">{wh.capacity}</td>
                  <td className="px-6 py-4"><StatusBadge status={wh.status} /></td>
                  <td className="px-6 py-4 text-right space-x-3">
                    <ActionLink onClick={() => {}}>Edit</ActionLink>
                    <ActionLink onClick={() => setShowAssignModal(wh.id)}>Assign</ActionLink>
                  </td>
                </tr>
              ))}
              {activeTab === 'stores' && DEMO_STORES.map(st => (
                <tr key={st.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 text-xs font-black text-slate-900">{st.name}</td>
                  <td className="px-6 py-4 text-xs font-bold text-slate-600">{st.location}</td>
                  <td className="px-6 py-4 text-[10px] font-bold text-slate-500">{st.manager}</td>
                  <td className="px-6 py-4"><StatusBadge status={st.status} /></td>
                  <td className="px-6 py-4 text-right space-x-3">
                    <ActionLink onClick={() => {}}>Edit</ActionLink>
                    <ActionLink onClick={() => setShowAssignModal(st.id)}>Assign</ActionLink>
                  </td>
                </tr>
              ))}
              {activeTab === 'vehicles' && DEMO_VEHICLES.map(v => (
                <tr key={v.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 text-xs font-black font-mono text-slate-900">{v.plate}</td>
                  <td className="px-6 py-4 text-xs font-bold text-slate-600">{v.type}</td>
                  <td className="px-6 py-4 text-xs font-bold text-slate-600">{v.capacity}</td>
                  <td className="px-6 py-4"><StatusBadge status={v.status} /></td>
                  <td className="px-6 py-4 text-right"><ActionLink onClick={() => {}}>Edit</ActionLink></td>
                </tr>
              ))}
              {activeTab === 'drivers' && DEMO_DRIVERS.map(d => (
                <tr key={d.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 text-xs font-black text-slate-900">{d.name}</td>
                  <td className="px-6 py-4 text-[10px] font-bold text-slate-500">{d.email}</td>
                  <td className="px-6 py-4 text-xs font-bold text-slate-600">{d.vehicle}</td>
                  <td className="px-6 py-4"><StatusBadge status={d.status} /></td>
                  <td className="px-6 py-4 text-right"><ActionLink onClick={() => {}}>Edit</ActionLink></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Summary Strip */}
      <div className="flex gap-6 px-6 py-3 bg-slate-50 rounded-xl border border-slate-100">
        {[
          { icon: Warehouse, label: 'Warehouses', count: DEMO_WAREHOUSES.length },
          { icon: Store, label: 'Stores', count: DEMO_STORES.length },
          { icon: Truck, label: 'Vehicles', count: DEMO_VEHICLES.length },
          { icon: Users, label: 'Drivers', count: DEMO_DRIVERS.length },
        ].map(s => (
          <div key={s.label} className="flex items-center gap-2">
            <s.icon className="w-4 h-4 text-primary" />
            <span className="text-[10px] font-black text-slate-900">{s.count}</span>
            <span className="text-[10px] font-bold text-slate-400 uppercase">{s.label}</span>
          </div>
        ))}
      </div>

      {/* Assign Modal */}
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
                  {users.map((u: { id: string; email: string; role: string }) => (
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
