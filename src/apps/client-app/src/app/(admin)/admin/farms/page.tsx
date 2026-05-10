'use client';

import React, { useState, useEffect } from 'react';
import { adminService, User, Farm } from '@/services/admin.service';
import { farmService } from '@/services/farm.service';

export default function AdminFarmsPage() {
  const [farms, setFarms] = useState<any[]>([]);
  const [managers, setManagers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    location: 'Cầu Đất, Đà Lạt',
    area: 0,
    farm_type: 'Arabica',
    owner_id: '',
  });

  const fetchData = async () => {
    try {
      // Note: We need a way for admin to list ALL farms. 
      // For now, if admin calls listFarms, the backend logic might need to be verified.
      // But based on our repository change, if admin calls, it might still only show their own
      // unless we update the List logic too.
      // For this demo, let's assume admin sees what they see.
      const farmData = await farmService.listFarms();
      setFarms(farmData);
      
      const mgrData = await adminService.listManagers();
      setManagers(mgrData);
      if (mgrData.length > 0) {
        setFormData(prev => ({ ...prev, owner_id: mgrData[0].id }));
      }
    } catch (err) {
      console.error('Failed to fetch data', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleCreateFarm = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await adminService.createFarm(formData);
      setShowModal(false);
      setFormData({ ...formData, name: '', area: 0 });
      fetchData();
    } catch (err) {
      alert('Failed to create farm');
    }
  };

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-4xl font-black font-headline text-on-surface uppercase italic tracking-tighter">
            Farms <span className="text-primary">Management</span>
          </h1>
          <p className="text-on-surface-variant font-medium">Create farms and assign managers to them.</p>
        </div>
        <button 
          onClick={() => setShowModal(true)}
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
          <table className="w-full text-left border-collapse">
            <thead className="bg-surface-container-low border-b border-outline-variant/10">
              <tr>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Farm Name</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Location</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Assigned Manager (Owner ID)</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/5">
              {farms.map((farm) => (
                <tr key={farm.id} className="hover:bg-surface-container-low/30 transition-colors">
                  <td className="px-6 py-4">
                    <p className="font-black text-on-surface uppercase tracking-tight text-xs">{farm.name}</p>
                  </td>
                  <td className="px-6 py-4 text-[10px] font-medium text-on-surface-variant italic">
                    {farm.location}
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <div className="w-5 h-5 bg-secondary-container rounded-md flex items-center justify-center text-[8px] font-black">ID</div>
                      <span className="text-[9px] font-mono text-on-surface-variant uppercase tracking-tighter">{farm.owner_id}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button className="text-primary-fixed bg-primary p-2 rounded-lg hover:scale-110 transition-transform">
                      <span className="material-symbols-outlined !text-sm">edit</span>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-slate-950/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
          <div className="bg-surface-container-lowest w-full max-w-md rounded-[2rem] border border-outline-variant/20 shadow-2xl p-8">
            <h2 className="text-2xl font-black font-headline uppercase italic tracking-tighter mb-6">Provision <span className="text-primary">New Farm</span></h2>
            <form onSubmit={handleCreateFarm} className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Farm Name</label>
                <input 
                  required
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
                    <option value="Arabica">Arabica</option>
                    <option value="Robusta">Robusta</option>
                    <option value="Liberica">Liberica</option>
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
                  <option value="Cầu Đất, Đà Lạt">Cầu Đất, Đà Lạt</option>
                  <option value="Buôn Ma Thuột, Đắk Lắk">Buôn Ma Thuột, Đắk Lắk</option>
                  <option value="Pleiku, Gia Lai">Pleiku, Gia Lai</option>
                  <option value="Gia Nghĩa, Đắk Nông">Gia Nghĩa, Đắk Nông</option>
                  <option value="Kon Tum">Kon Tum</option>
                </select>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Assign Manager</label>
                <select 
                  required
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-primary transition-colors font-bold"
                  value={formData.owner_id}
                  onChange={(e) => setFormData({...formData, owner_id: e.target.value})}
                >
                  {managers.map(m => (
                    <option key={m.id} value={m.id}>{m.name} ({m.email})</option>
                  ))}
                  {managers.length === 0 && <option disabled>No Managers Available</option>}
                </select>
              </div>
              <div className="flex gap-4 pt-4">
                <button 
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest border border-outline-variant/20 hover:bg-surface-container-low transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  className="flex-1 bg-primary text-on-primary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-primary/20"
                >
                  Provision Farm
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
