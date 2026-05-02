'use client';

import Link from "next/link";
import { usePathname } from "next/navigation";

export default function Sidebar() {
  const pathname = usePathname();

  const mainLinks = [
    { href: "/", label: "System Overview", icon: "account_tree" },
    { href: "/app/topology", label: "Microservices", icon: "schema" },
    { href: "/app/batches", label: "Data Flow", icon: "dataset" },
    { href: "/explorer", label: "Infrastructure", icon: "dns" },
    { href: "/app/resiliency", label: "Observability", icon: "monitoring" },
  ];

  const supplyChainLinks = [
    { href: "/app/farms", label: "Farm Origin", icon: "agriculture" },
    { href: "/app/logistics", label: "Real-time Transit", icon: "local_shipping" },
    { href: "/app/warehouse", label: "Warehouse Stock", icon: "warehouse" },
    { href: "/app/retail", label: "Retail Saga", icon: "store" },
  ];

  return (
    <aside className="flex flex-col h-full p-4 space-y-2 fixed left-0 top-16 z-40 bg-[#f2f4f6] w-64 border-r border-outline-variant/10">
      <div className="flex items-center gap-3 px-4 py-6 mb-2">
        <div className="w-10 h-10 bg-primary-container rounded-xl flex items-center justify-center text-on-primary-container shadow-lg">
          <span className="material-symbols-outlined">coffee</span>
        </div>
        <div>
          <h2 className="text-sm font-bold font-headline text-on-surface uppercase tracking-tighter italic">Runtime Roasters</h2>
          <p className="text-[10px] text-on-surface-variant font-black uppercase tracking-widest opacity-60">Global Origin v4.1</p>
        </div>
      </div>
      
      <nav className="flex-1 space-y-1 overflow-y-auto custom-scrollbar pr-1">
        {mainLinks.map((link) => (
          <SidebarLink key={link.href} {...link} active={pathname === link.href} />
        ))}
        
        <div className="px-4 py-6 border-t border-outline-variant/10 mt-6">
          <p className="text-[9px] font-black text-slate-400 uppercase tracking-[0.3em] mb-6 italic opacity-70">Supply Chain Ops</p>
          <div className="space-y-1">
            {supplyChainLinks.map((link) => (
              <SidebarLink key={link.href} {...link} active={pathname === link.href} small />
            ))}
          </div>
        </div>

        <div className="pt-6 px-2">
          <button className="w-full bg-primary text-on-primary py-4 px-4 rounded-xl font-bold flex items-center justify-center gap-2 shadow-xl shadow-primary/20 hover:scale-[0.98] transition-all uppercase text-[10px] tracking-widest italic">
            <span className="material-symbols-outlined !text-sm font-black">add</span>
            New Component
          </button>
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
