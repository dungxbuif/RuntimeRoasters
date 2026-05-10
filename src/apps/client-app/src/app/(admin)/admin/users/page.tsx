'use client';

import React, { useState, useEffect } from 'react';
import { adminService, User } from '@/services/admin.service';

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    name: '',
    role: 'farm_manager',
  });

  const fetchUsers = async () => {
    try {
      const data = await adminService.listUsers();
      setUsers(data);
    } catch (err) {
      console.error('Failed to fetch users', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await adminService.createUser(formData);
      setShowModal(false);
      setFormData({ email: '', password: '', name: '', role: 'farm_manager' });
      fetchUsers();
    } catch (err) {
      alert('Failed to create user');
    }
  };

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-4xl font-black font-headline text-on-surface uppercase italic tracking-tighter">
            User <span className="text-tertiary">Management</span>
          </h1>
          <p className="text-on-surface-variant font-medium">Onboard new Farm Managers and staff.</p>
        </div>
        <button 
          onClick={() => setShowModal(true)}
          className="bg-tertiary text-on-tertiary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-tertiary/20 flex items-center gap-2"
        >
          <span className="material-symbols-outlined !text-sm">person_add</span>
          Create New User
        </button>
      </div>

      <div className="bg-surface-container-lowest rounded-3xl border border-outline-variant/10 overflow-hidden shadow-sm">
        {isLoading ? (
          <div className="p-12 text-center text-on-surface-variant animate-pulse font-bold uppercase tracking-widest text-xs">Loading Identities...</div>
        ) : (
          <table className="w-full text-left border-collapse">
            <thead className="bg-surface-container-low border-b border-outline-variant/10">
              <tr>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">User</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Role</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">ID</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/5">
              {users.map((user) => (
                <tr key={user.id} className="hover:bg-surface-container-low/30 transition-colors">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-surface-container-high rounded-lg flex items-center justify-center font-black">{user.name[0]}</div>
                      <div>
                        <p className="font-bold text-on-surface uppercase tracking-tight text-xs">{user.name}</p>
                        <p className="text-[10px] text-on-surface-variant">{user.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 bg-surface-container text-[10px] font-black rounded-full uppercase tracking-tighter border border-outline-variant/10">{user.role}</span>
                  </td>
                  <td className="px-6 py-4 font-mono text-[9px] text-on-surface-variant opacity-50">
                    {user.id}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button className="text-on-surface-variant hover:text-primary transition-colors">
                      <span className="material-symbols-outlined">more_vert</span>
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
            <h2 className="text-2xl font-black font-headline uppercase italic tracking-tighter mb-6">Create <span className="text-tertiary">New Identity</span></h2>
            <form onSubmit={handleCreateUser} className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Full Name</label>
                <input 
                  required
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Email Address</label>
                <input 
                  type="email"
                  required
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors"
                  value={formData.email}
                  onChange={(e) => setFormData({...formData, email: e.target.value})}
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Temporary Password</label>
                <input 
                  type="password"
                  required
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors"
                  value={formData.password}
                  onChange={(e) => setFormData({...formData, password: e.target.value})}
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">System Role</label>
                <select 
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors"
                  value={formData.role}
                  onChange={(e) => setFormData({...formData, role: e.target.value})}
                >
                  <option value="farm_manager">FARM_MANAGER</option>
                  <option value="farm_admin">FARM_ADMIN</option>
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
                  className="flex-1 bg-tertiary text-on-tertiary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-tertiary/20"
                >
                  Create User
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
