'use client';

import { adminService, User } from '@/services/admin.service';
import React, { useState } from 'react';
import { testId, e2eSelectors } from '@/lib/utils/test-id';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { USER_ROLE_LABELS } from '@/constants/domain';

export default function AdminUsersPage() {
  const queryClient = useQueryClient();
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    role: 'FARM_MANAGER',
  });

  const { data: users = [], isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: () => adminService.listUsers(),
  });

  const createUserMutation = useMutation({
    mutationFn: (data: any) => adminService.createUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      setShowModal(false);
      setFormData({ email: '', password: '', role: 'FARM_MANAGER' });
    },
    onError: (err) => {
      console.error('Failed to create user', err);
      alert('Error creating user. Check console for details.');
    }
  });

  const handleCreateUser = (e: React.FormEvent) => {
    e.preventDefault();
    const name = formData.email.split('@')[0];
    createUserMutation.mutate({
      ...formData,
      name,
    });
  };

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-4xl font-black font-headline text-on-surface uppercase italic tracking-tighter">
            Users <span className="text-tertiary">& Teams</span>
          </h1>
          <p className="text-on-surface-variant font-medium">Provision and manage cluster identities.</p>
        </div>
        <button 
          onClick={() => setShowModal(true)}
          {...testId(e2eSelectors.CREATE_USER_BTN)}
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
          <table className="w-full text-left border-collapse" {...testId(e2eSelectors.USER_LIST_TABLE)}>
            <thead className="bg-surface-container-low border-b border-outline-variant/10">
              <tr>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">User</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">Role</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant">ID</th>
                <th className="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-on-surface-variant text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/5">
              {users.map((user: User) => (
                <tr key={user.id} className="hover:bg-surface-container-low/30 transition-colors" {...testId(`user-row-${user.email}`)}>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center font-black">{user.email[0].toUpperCase()}</div>
                      <div>
                        <p className="font-bold text-on-surface uppercase tracking-tight text-xs">{user.email.split('@')[0]}</p>
                        <p className="text-[9px] text-on-surface-variant/60 font-mono">{user.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full bg-surface-container-high text-[9px] font-black uppercase tracking-widest border border-outline-variant/10">
                      {USER_ROLE_LABELS[user.role] || user.role}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <code className="text-[9px] text-on-surface-variant/40 font-mono">{user.id}</code>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button className="text-on-surface-variant hover:text-primary transition-colors">
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
        <div className="fixed inset-0 bg-slate-950/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4" {...testId(e2eSelectors.USER_MODAL)}>
          <div className="bg-surface-container-lowest w-full max-w-md rounded-[2rem] border border-outline-variant/20 shadow-2xl p-8">
            <h2 className="text-2xl font-black font-headline uppercase italic tracking-tighter mb-6">Create <span className="text-tertiary">New User</span></h2>
            <form onSubmit={handleCreateUser} className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">Email Address</label>
                <input 
                  type="email"
                  required
                  {...testId(e2eSelectors.USER_EMAIL_INPUT)}
                  placeholder="manager@runtimeroasters.com"
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
                  {...testId(e2eSelectors.USER_PASS_INPUT)}
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors"
                  value={formData.password}
                  onChange={(e) => setFormData({...formData, password: e.target.value})}
                />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black uppercase tracking-widest ml-1 opacity-60">System Role</label>
                <select 
                  className="w-full bg-surface-container-low border border-outline-variant/20 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-tertiary transition-colors font-bold"
                  value={formData.role}
                  {...testId(e2eSelectors.USER_ROLE_SELECT)}
                  onChange={(e) => setFormData({...formData, role: e.target.value})}
                >
                  {Object.entries(USER_ROLE_LABELS).map(([id, label]) => (
                    <option key={id} value={id}>{label}</option>
                  ))}
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
                  disabled={createUserMutation.isPending}
                  {...testId(e2eSelectors.USER_SUBMIT_BTN)}
                  className="flex-1 bg-tertiary text-on-tertiary px-6 py-3 rounded-xl font-bold uppercase text-xs tracking-widest shadow-lg shadow-tertiary/20 disabled:opacity-50"
                >
                  {createUserMutation.isPending ? 'Processing...' : 'Create User'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
