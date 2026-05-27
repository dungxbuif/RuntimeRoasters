'use client';

import Link from "next/link";
import { usePathname } from "next/navigation";

import { APP_ROUTES } from "@/constants/routes";
import { useAuth } from "@/lib/auth";
import { PermissionGuard } from "@/lib/auth/casbin";

export default function Sidebar() {
  const pathname = usePathname();
  const { user } = useAuth();

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
        {/* ADMIN: Identity & Access */}
        <PermissionGuard action="read" resource="sidebar_users">
          <NavSection label="Identity & Access">
            <SidebarLink href={APP_ROUTES.DASHBOARD.USERS} label="Manage Users" icon="group_add" active={pathname === APP_ROUTES.DASHBOARD.USERS} />
            <SidebarLink href={APP_ROUTES.DASHBOARD.RESOURCES} label="Resources" icon="domain" active={pathname === APP_ROUTES.DASHBOARD.RESOURCES} />
          </NavSection>
        </PermissionGuard>

        {/* ADMIN / FARM_ADMIN / FARM_MANAGER: Farm Operations */}
        <PermissionGuard action="read" resource="sidebar_farms">
          <NavSection label="Farm Operations">
            <SidebarLink href={APP_ROUTES.DASHBOARD.FARMS} label="Farm Registry" icon="eco" active={pathname === APP_ROUTES.DASHBOARD.FARMS} />
            <SidebarLink href={APP_ROUTES.DASHBOARD.HARVESTS} label="Harvests" icon="grass" active={pathname === APP_ROUTES.DASHBOARD.HARVESTS} />
          </NavSection>
        </PermissionGuard>

        {/* WAREHOUSE_MGR / PROCESSOR: Warehouse Operations */}
        <PermissionGuard action="read" resource="sidebar_warehouse">
          <NavSection label="Warehouse Operations">
            <SidebarLink href={APP_ROUTES.DASHBOARD.WAREHOUSE} label="Warehouse Ops" icon="warehouse" active={pathname === APP_ROUTES.DASHBOARD.WAREHOUSE} />
            <SidebarLink href={APP_ROUTES.DASHBOARD.BATCHES} label="Batch Lifecycle" icon="science" active={pathname === APP_ROUTES.DASHBOARD.BATCHES} />
          </NavSection>
        </PermissionGuard>

        {/* STORE_MGR: Retail Operations */}
        <PermissionGuard action="read" resource="sidebar_store">
          <NavSection label="Retail Operations">
            <SidebarLink href={APP_ROUTES.DASHBOARD.STORE} label="Store Dashboard" icon="storefront" active={pathname === APP_ROUTES.DASHBOARD.STORE} />
            <SidebarLink href={APP_ROUTES.DASHBOARD.RETAIL_ORDERS} label="Create Order" icon="add_shopping_cart" active={pathname === APP_ROUTES.DASHBOARD.RETAIL_ORDERS} />
          </NavSection>
        </PermissionGuard>

        {/* DRIVER: Driver Client */}
        <PermissionGuard action="read" resource="sidebar_driver">
          <NavSection label="Driver Client">
            <SidebarLink href={APP_ROUTES.DASHBOARD.DRIVER} label="My Shipments" icon="local_shipping" active={pathname.startsWith(APP_ROUTES.DASHBOARD.DRIVER)} />
          </NavSection>
        </PermissionGuard>

        {/* All authenticated: Supply Chain Ops */}
        <NavSection label="Supply Chain Ops">
          <SidebarLink href="/dashboard" label="Intelligence Hub" icon="analytics" active={pathname === '/dashboard'} small />
          <SidebarLink href={APP_ROUTES.DASHBOARD.LOGISTICS} label="Logistics Map" icon="map" active={pathname === APP_ROUTES.DASHBOARD.LOGISTICS} small />
          <SidebarLink href={APP_ROUTES.DASHBOARD.TRACEABILITY} label="Provenance Trace" icon="qr_code" active={pathname === APP_ROUTES.DASHBOARD.TRACEABILITY} small />
          <SidebarLink href={APP_ROUTES.DASHBOARD.FINANCE} label="Finance" icon="payments" active={pathname === APP_ROUTES.DASHBOARD.FINANCE} small />
        </NavSection>

        {/* ADMIN only: System Intelligence */}
        <PermissionGuard action="read" resource="sidebar_dashboard">
          <NavSection label="System Intelligence">
            <SidebarLink href={APP_ROUTES.DASHBOARD.RETAIL} label="Saga Monitor" icon="device_hub" active={pathname === APP_ROUTES.DASHBOARD.RETAIL} small />
            <SidebarLink href={APP_ROUTES.DASHBOARD.EXPLORER} label="System Explorer" icon="dns" active={pathname === APP_ROUTES.DASHBOARD.EXPLORER} small />
            <SidebarLink href={APP_ROUTES.DASHBOARD.AUDIT} label="Audit Logs" icon="shield" active={pathname === APP_ROUTES.DASHBOARD.AUDIT} small />
          </NavSection>
        </PermissionGuard>

        {/* Account + Showcase link */}
        <NavSection label="Account">
          <SidebarLink href={APP_ROUTES.DASHBOARD.PROFILE} label="Account Settings" icon="manage_accounts" active={pathname === APP_ROUTES.DASHBOARD.PROFILE} small />
        </NavSection>

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

function NavSection({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="px-1 mb-6 pt-4 border-t border-outline-variant/10 first:border-t-0 first:pt-0">
      <p className="text-[9px] font-black text-slate-400 uppercase tracking-[0.3em] mb-3 px-4 italic opacity-70">{label}</p>
      <div className="space-y-1">{children}</div>
    </div>
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
