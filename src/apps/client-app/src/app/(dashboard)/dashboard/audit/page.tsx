'use client';

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { financeService } from '@/services/finance.service';
import { ShieldCheck, History, Database, Search, Fingerprint, Lock, ShieldAlert } from 'lucide-react';

export default function AuditPage() {
  const [partition, setPartition] = useState('default');

  const { data: logs = [], isLoading } = useQuery({
    queryKey: ['audit', partition],
    queryFn: () => financeService.listAuditLogs(partition),
    refetchInterval: 15000,
  });

  return (
    <div className="h-full flex flex-col gap-8">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-black uppercase tracking-tighter text-slate-900 flex items-center gap-3 italic">
            <ShieldCheck className="w-8 h-8 text-primary" />
            Immutable Audit Trail
          </h1>
          <p className="text-xs font-bold text-slate-500 uppercase tracking-widest mt-1 italic">
            Tamper-proof Ledger via Hash Chaining & Cassandra Storage
          </p>
        </div>

        <div className="flex bg-white p-1 rounded-2xl border-2 border-slate-200">
           {['default', 'system', 'security'].map(p => (
             <button
                key={p}
                onClick={() => setPartition(p)}
                className={`px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all ${
                  partition === p ? 'bg-slate-900 text-white shadow-lg' : 'text-slate-400 hover:text-slate-600'
                }`}
             >
               {p}
             </button>
           ))}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8 flex-1 min-h-0 overflow-hidden">
        {/* Left: Audit Info */}
        <div className="lg:col-span-1 space-y-6">
           <div className="bg-slate-900 rounded-[2rem] p-8 text-white relative overflow-hidden group">
              <Database className="absolute -right-4 -bottom-4 w-32 h-32 opacity-10 group-hover:rotate-12 transition-transform duration-700" />
              <div className="relative z-10">
                 <h3 className="text-[10px] font-black text-slate-400 uppercase tracking-[0.3em] mb-6">Storage Node</h3>
                 <div className="text-3xl font-black italic tracking-tighter mb-2">Apache Cassandra</div>
                 <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest leading-relaxed">
                   Optimized for high-velocity immutable event logging.
                 </p>
              </div>
           </div>

           <div className="bg-white rounded-[2rem] p-8 border border-slate-200 space-y-6">
              <div className="flex items-center gap-4">
                 <div className="w-10 h-10 rounded-xl bg-green-50 flex items-center justify-center text-green-500">
                    <Lock className="w-5 h-5" />
                 </div>
                 <div>
                    <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Integrity</div>
                    <div className="text-sm font-black text-slate-900 uppercase italic">Hash Chained</div>
                 </div>
              </div>
              <div className="flex items-center gap-4">
                 <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center text-primary">
                    <Fingerprint className="w-5 h-5" />
                 </div>
                 <div>
                    <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Authentication</div>
                    <div className="text-sm font-black text-slate-900 uppercase italic">Digital Signatures</div>
                 </div>
              </div>
              <div className="flex items-center gap-4">
                 <div className="w-10 h-10 rounded-xl bg-amber-50 flex items-center justify-center text-amber-500">
                    <ShieldAlert className="w-5 h-5" />
                 </div>
                 <div>
                    <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Compliance</div>
                    <div className="text-sm font-black text-slate-900 uppercase italic">SOX-Ready</div>
                 </div>
              </div>
           </div>
        </div>

        {/* Right: Audit Logs */}
        <div className="lg:col-span-3 bg-white rounded-[2.5rem] border border-slate-200 shadow-2xl overflow-hidden flex flex-col">
           <div className="p-6 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
              <div className="flex items-center gap-2">
                 <History className="w-4 h-4 text-slate-400" />
                 <h2 className="text-xs font-black uppercase tracking-[0.2em] text-slate-500">Event Stream: {partition}</h2>
              </div>
              <div className="relative">
                 <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400" />
                 <input 
                    type="text" 
                    placeholder="Search logs..." 
                    className="bg-white border border-slate-200 rounded-full py-1.5 pl-8 pr-4 text-[10px] font-bold outline-none focus:border-primary w-48"
                 />
              </div>
           </div>

           <div className="flex-1 overflow-y-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="bg-slate-50/30 text-[9px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-100">
                    <th className="px-6 py-4">Event ID</th>
                    <th className="px-6 py-4">Type</th>
                    <th className="px-6 py-4">Aggregate ID</th>
                    <th className="px-6 py-4">Hash Proof</th>
                    <th className="px-6 py-4">Occurred At</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-50">
                   {isLoading ? (
                     <tr>
                        <td colSpan={5} className="px-6 py-20 text-center">
                           <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
                           <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Querying Cassandra Node...</span>
                        </td>
                     </tr>
                   ) : logs.length === 0 ? (
                     <tr>
                        <td colSpan={5} className="px-6 py-20 text-center">
                           <div className="text-[10px] font-black text-slate-400 uppercase tracking-widest italic opacity-50">No logs found in this partition</div>
                        </td>
                     </tr>
                   ) : (
                     logs.map(log => (
                       <tr key={log.id} className="hover:bg-slate-50/80 transition-all cursor-pointer group">
                          <td className="px-6 py-4">
                             <div className="flex items-center gap-2">
                                <span className="w-2 h-2 rounded-full bg-primary opacity-30"></span>
                                <span className="text-xs font-black font-mono text-slate-900">{log.id.slice(0, 13).toUpperCase()}</span>
                             </div>
                          </td>
                          <td className="px-6 py-4">
                             <span className="text-[10px] font-black text-slate-600 bg-slate-100 px-2 py-0.5 rounded border border-slate-200 uppercase tracking-tighter">
                                {log.event_type}
                             </span>
                          </td>
                          <td className="px-6 py-4">
                             <span className="text-[10px] font-bold text-slate-500 font-mono">{log.aggregate_id}</span>
                          </td>
                          <td className="px-6 py-4">
                             <div className="flex items-center gap-2 group-hover:text-primary transition-colors">
                                <span className="text-[8px] font-mono text-slate-400 bg-slate-50 px-1.5 py-0.5 rounded truncate max-w-[80px]">{log.hash}</span>
                                <Lock className="w-3 h-3 opacity-20" />
                             </div>
                          </td>
                          <td className="px-6 py-4">
                             <span className="text-[10px] font-bold text-slate-400 uppercase">{new Date(log.occurred_at).toLocaleString()}</span>
                          </td>
                       </tr>
                     ))
                   )}
                </tbody>
              </table>
           </div>
        </div>
      </div>
    </div>
  );
}
