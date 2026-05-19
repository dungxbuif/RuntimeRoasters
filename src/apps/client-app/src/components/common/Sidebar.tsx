'use client';

import Link from "next/link";
import { usePathname } from "next/navigation";

import { APP_ROUTES } from "@/constants/routes";
import { RoleGuard } from "@/lib/auth";

export default function Sidebar() {
  const pathname = usePathname();

  const managementLinks = [
    { href: APP_ROUTES.DASHBOARD.USERS, label: "Manage Users", icon: "group_add" },
    { href: APP_ROUTES.DASHBOARD.FARMS, label: "Manage Farms", icon: "admin_panel_settings" },
    { href: APP_ROUTES.DASHBOARD.PROFILE, label: "Account Settings", icon: "manage_accounts" },
  ];

  const operationalLinks = [
    { href: APP_ROUTES.DASHBOARD.HARVESTS, label: "Harvest Declaration", icon: "eco" },
    { href: APP_ROUTES.DASHBOARD.FARM_TELEMETRY, label: "Farm Telemetry", icon: "agriculture" },
    { href: APP_ROUTES.DASHBOARD.LOGISTICS, label: "Transit Monitor", icon: "local_shipping" },
    { href: APP_ROUTES.DASHBOARD.WAREHOUSE, label: "Stock Analytics", icon: "warehouse" },
    { href: APP_ROUTES.DASHBOARD.RETAIL, label: "Saga Monitor", icon: "analytics" },
    { href: APP_ROUTES.DASHBOARD.RETAIL_ORDERS, label: "Market Orders", icon: "shopping_cart" },
    { href: APP_ROUTES.DASHBOARD.TRACEABILITY, label: "Provenance Trace", icon: "qr_code" },
  ];

  const diagnosticLinks = [
    { href: APP_ROUTES.DASHBOARD.TOPOLOGY, label: "Cluster Mesh", icon: "schema" },
    { href: APP_ROUTES.DASHBOARD.BATCHES, label: "Stream Flow", icon: "dataset" },
    { href: APP_ROUTES.DASHBOARD.FINANCE, label: "Finance Ledger", icon: "payments" },
    { href: APP_ROUTES.DASHBOARD.AUDIT, label: "Immutable Audit", icon: "shield_check" },
    { href: APP_ROUTES.DASHBOARD.EXPLORER, label: "Nodes & Pods", icon: "dns" },
    { href: APP_ROUTES.DASHBOARD.RESILIENCY, label: "Fault Logs", icon: "monitoring" },
  ];

  return (
    <aside className="flex flex-col h-full p-4 space-y-2 fixed left-0 top-0 z-40 bg-[#f2f4f6] w-64 border-r border-outline-variant/10 shadow-2xl shadow-black/20">
      <div className="flex items-center gap-3 px-4 py-8 mb-4">
        <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center text-white shadow-lg">
          <span className="material-symbols-outlined">coffee</span>
        </div>
        <div>
          <h2 className="text-sm font-black font-headline text-on-surface uppercase tracking-tighter italic text-primary">Runtime Roasters</h2>
          <p className="text-[10px] text-slate-500 font-black uppercase tracking-[0.2em] opacity-80 mt-1">Management Portal</p>
        </div>
      </div>
      
      <nav className="flex-1 space-y-1 overflow-y-auto custom-scrollbar pr-1">
        {/* Priority Management Section - PROTECTED */}
        <RoleGuard roles={['ADMIN', 'FARM_ADMIN', 'FARM_MANAGER']}>
          <div className="px-1 mb-8">
            <p className="text-[9px] font-black text-slate-500 uppercase tracking-[0.3em] mb-4 px-4 italic opacity-70">Identity & Access</p>
            <div className="space-y-1">
              {managementLinks.map((link) => (
                <SidebarLink key={link.href} {...link} active={pathname === link.href} />
              ))}
            </div>
          </div>
        </RoleGuard>

        {/* Operational Section - PROTECTED */}
        <div className="px-1 mb-8 pt-4 border-t border-outline-variant/10">
          <p className="text-[9px] font-black text-slate-400 uppercase tracking-[0.3em] mb-4 px-4 italic opacity-70">Supply Chain Ops</p>
          <div className="space-y-1">
            {operationalLinks.map((link) => (
              <SidebarLink key={link.href} {...link} active={pathname === link.href} small />
            ))}
          </div>
        </div>

        {/* System & Diagnostics Section - PROTECTED */}
        <div className="px-1 pt-4 border-t border-outline-variant/10">
          <p className="text-[9px] font-black text-slate-400 uppercase tracking-[0.3em] mb-4 px-4 italic opacity-70">System Intelligence</p>
          <div className="space-y-1">
            {diagnosticLinks.map((link) => (
              <SidebarLink key={link.href} {...link} active={pathname === link.href} small />
            ))}
          </div>
        </div>

        {/* Home/Showcase Link */}
        <div className="mt-10 px-4">
          <Link 
            href={APP_ROUTES.HOME}
            className="flex items-center gap-3 text-primary hover:text-primary-fixed text-[10px] font-black uppercase tracking-[0.2em] italic transition-colors"
          >
            <span className="material-symbols-outlined !text-sm">arrow_back</span>
            Back to Showcase
          </Link>
        </div>
      </nav>

      <div className="mt-auto space-y-1 pb-10 border-t border-outline-variant/10 pt-6">
        <Link className="flex items-center gap-3 px-4 py-2 text-slate-500 hover:text-primary text-[11px] font-black uppercase tracking-widest transition-colors italic" href="#">
          <span className="material-symbols-outlined text-sm">auto_stories</span>
          Docs
        </Link>
        <Link className="flex items-center gap-3 px-4 py-2 text-slate-500 hover:text-primary text-[11px] font-black uppercase tracking-widest transition-colors italic" href="#">
          <span className="material-symbols-outlined text-sm">help_outline</span>
          Support
        </Link>
      </div>
    </aside>
  );
}

function SidebarLink({ href, icon, label, active, small = false }: { href: string; icon: string; label: string; active: boolean; small?: boolean }) {
  return (
    <Link 
      href={href} 
      className={`flex items-center gap-4 px-4 ${small ? 'py-2.5' : 'py-3.5'} rounded-xl transition-all group relative overflow-hidden ${
        active 
          ? 'bg-white text-tertiary shadow-[4px_4px_0px_0px_rgba(0,98,66,0.08)] border border-tertiary/20 font-black' 
          : 'text-on-surface-variant/70 hover:bg-white/40 font-bold'
      }`}
    >
      <span className={`material-symbols-outlined ${small ? 'text-lg' : 'text-xl'} ${active ? 'opacity-100 text-tertiary font-black' : 'opacity-40 group-hover:opacity-100 group-hover:text-tertiary'} transition-all`}>
        {icon}
      </span>
      <span className={`${small ? 'text-[11px]' : 'text-[13px]'} uppercase tracking-tight ${active ? 'text-on-surface' : 'group-hover:text-on-surface'} transition-colors whitespace-nowrap`}>
        {label}
      </span>
    </Link>
  );
}
