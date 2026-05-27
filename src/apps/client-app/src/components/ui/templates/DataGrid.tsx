import React, { ReactNode } from 'react';
import { LucideIcon } from 'lucide-react';

interface Column<T> {
  key: string;
  header: string;
  render: (item: T) => ReactNode;
}

interface DataGridProps<T> {
  title: string;
  icon?: LucideIcon;
  badge?: string;
  data: T[];
  columns: Column<T>[];
  keyExtractor: (item: T) => string;
}

export function DataGrid<T>({ title, icon: Icon, badge, data, columns, keyExtractor }: DataGridProps<T>) {
  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden flex flex-col h-full">
      <div className="px-6 py-4 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
        <h3 className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 flex items-center gap-2">
          {Icon && <Icon className="w-4 h-4 text-primary" />}
          {title}
        </h3>
        {badge && (
          <span className="text-[9px] font-bold bg-slate-900 text-white px-2 py-0.5 rounded italic">
            {badge}
          </span>
        )}
      </div>
      <div className="overflow-x-auto flex-1">
        <table className="w-full text-left">
          <thead>
            <tr className="bg-slate-50/30 text-[9px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-100">
              {columns.map(col => (
                <th key={col.key} className="px-6 py-3">{col.header}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-50">
            {data.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-8 text-center text-xs text-slate-400 italic">
                  No data available
                </td>
              </tr>
            ) : (
              data.map(item => (
                <tr key={keyExtractor(item)} className="hover:bg-slate-50/50 transition-colors">
                  {columns.map(col => (
                    <td key={col.key} className="px-6 py-4">
                      {col.render(item)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
