'use client';

import React from 'react';

type NodeTone = 'core' | 'healthy' | 'warning' | 'danger' | 'neutral';

interface MeshNode {
  id: string;
  title: string;
  subtitle: string;
  icon: string;
  zone: string;
  tone: NodeTone;
  metric: string;
  position: string;
}

const nodes: MeshNode[] = [
  {
    id: 'browser',
    title: 'Client App',
    subtitle: 'Next.js Browser',
    icon: 'web_asset',
    zone: 'Edge',
    tone: 'core',
    metric: '3000',
    position: 'left-[4%] top-[10%]',
  },
  {
    id: 'gateway',
    title: 'KrakenD Gateway',
    subtitle: 'Unified API Surface',
    icon: 'hub',
    zone: 'Ingress',
    tone: 'healthy',
    metric: '8081',
    position: 'left-[28%] top-[8%]',
  },
  {
    id: 'hydra',
    title: 'Ory Hydra',
    subtitle: 'OIDC Provider',
    icon: 'lock',
    zone: 'Identity',
    tone: 'warning',
    metric: '4444',
    position: 'left-[53%] top-[5%]',
  },
  {
    id: 'kratos',
    title: 'Ory Kratos',
    subtitle: 'Identity Store',
    icon: 'fingerprint',
    zone: 'Identity',
    tone: 'warning',
    metric: '4433',
    position: 'left-[75%] top-[17%]',
  },
  {
    id: 'auth',
    title: 'Auth Service',
    subtitle: 'RBAC + Login Accept',
    icon: 'admin_panel_settings',
    zone: 'Service',
    tone: 'core',
    metric: '8082',
    position: 'left-[35%] top-[42%]',
  },
  {
    id: 'farm',
    title: 'Farm Service',
    subtitle: 'Farms + Harvests',
    icon: 'agriculture',
    zone: 'Service',
    tone: 'healthy',
    metric: '8083',
    position: 'left-[63%] top-[46%]',
  },
  {
    id: 'postgres',
    title: 'PostgreSQL',
    subtitle: 'State + Casbin Rules',
    icon: 'database',
    zone: 'Persistence',
    tone: 'neutral',
    metric: '54321',
    position: 'left-[21%] top-[73%]',
  },
  {
    id: 'kafka',
    title: 'Kafka Bus',
    subtitle: 'Domain Events',
    icon: 'sync_alt',
    zone: 'Events',
    tone: 'danger',
    metric: '9094',
    position: 'left-[51%] top-[75%]',
  },
  {
    id: 'valkey',
    title: 'Valkey Cache',
    subtitle: 'Runtime Cache',
    icon: 'bolt',
    zone: 'Cache',
    tone: 'healthy',
    metric: '6379',
    position: 'left-[77%] top-[70%]',
  },
];

const toneClasses: Record<NodeTone, { node: string; icon: string; badge: string; dot: string }> = {
  core: {
    node: 'bg-on-background border-slate-700 text-white shadow-2xl',
    icon: 'bg-primary-fixed text-on-primary-fixed',
    badge: 'bg-primary-fixed/20 text-primary-fixed border-primary-fixed/30',
    dot: 'bg-primary-fixed',
  },
  healthy: {
    node: 'bg-surface-container-low border-tertiary-fixed text-on-surface shadow-lg',
    icon: 'bg-tertiary-fixed text-on-tertiary-fixed',
    badge: 'bg-tertiary-fixed/20 text-tertiary border-tertiary/20',
    dot: 'bg-tertiary-fixed',
  },
  warning: {
    node: 'bg-amber-50 border-amber-300 text-amber-950 shadow-lg',
    icon: 'bg-amber-400 text-amber-950',
    badge: 'bg-amber-100 text-amber-800 border-amber-300',
    dot: 'bg-amber-400',
  },
  danger: {
    node: 'bg-error-container/80 border-error text-on-error-container shadow-2xl',
    icon: 'bg-error text-on-error',
    badge: 'bg-error text-on-error border-error',
    dot: 'bg-error',
  },
  neutral: {
    node: 'bg-surface-container-lowest border-outline-variant/20 text-on-surface shadow-md',
    icon: 'bg-surface-container-high text-on-surface-variant',
    badge: 'bg-surface-container-low text-on-surface-variant border-outline-variant/20',
    dot: 'bg-slate-400',
  },
};

const links = [
  { d: 'M 190 95 L 348 95', color: '#004ac6' },
  { d: 'M 504 95 L 650 95', color: '#006242' },
  { d: 'M 787 111 Q 878 118 918 174', color: '#d89b00' },
  { d: 'M 418 176 Q 430 281 462 358', color: '#004ac6' },
  { d: 'M 653 178 Q 620 285 612 386', color: '#d89b00' },
  { d: 'M 572 454 L 760 472', color: '#006242' },
  { d: 'M 485 526 Q 374 622 340 691', color: '#64748b' },
  { d: 'M 635 550 Q 664 650 660 710', color: '#ba1a1a' },
  { d: 'M 827 552 Q 885 640 918 691', color: '#006242' },
];

function MeshNodeCard({ node }: { node: MeshNode }) {
  const tone = toneClasses[node.tone];

  return (
    <div className={`absolute ${node.position} w-56 h-40 rounded-3xl border-2 p-5 flex flex-col justify-between transition-transform hover:scale-[1.03] ${tone.node}`}>
      <div className="flex items-start justify-between gap-3">
        <div className={`w-12 h-12 rounded-full flex items-center justify-center shadow-inner ${tone.icon}`}>
          <span className="material-symbols-outlined !text-2xl">{node.icon}</span>
        </div>
        <span className={`px-2.5 py-1 rounded-full border text-[9px] font-black uppercase tracking-widest ${tone.badge}`}>
          {node.zone}
        </span>
      </div>

      <div>
        <h3 className="font-headline text-sm font-black uppercase tracking-tight italic leading-tight">
          {node.title}
        </h3>
        <p className="mt-1 text-[10px] font-bold uppercase tracking-tight opacity-60">
          {node.subtitle}
        </p>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className={`w-2 h-2 rounded-full ${tone.dot} ${node.tone === 'danger' ? 'animate-pulse' : ''}`}></span>
          <span className="font-mono text-[10px] font-black uppercase opacity-70">Live</span>
        </div>
        <span className="font-mono text-[10px] font-black opacity-50">:{node.metric}</span>
      </div>
    </div>
  );
}

export default function SystemTopologyPage() {
  return (
    <div className="min-h-full bg-surface text-on-surface p-8 md:p-12 lg:p-16 font-body relative overflow-x-hidden">
      <div className="max-w-7xl mx-auto space-y-8">
        <section className="flex flex-col lg:flex-row lg:items-end justify-between gap-6 border-b border-outline-variant/30 pb-8">
          <div>
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary-fixed/30 text-on-primary-fixed text-[10px] font-black uppercase tracking-[0.2em] mb-4 border border-primary/10 shadow-sm italic">
              <span className="material-symbols-outlined !text-sm">account_tree</span>
              Runtime Mesh
            </div>
            <h1 className="text-4xl md:text-5xl font-black font-headline tracking-tighter uppercase italic">
              Architecture <span className="text-primary">Topology</span>
            </h1>
            <p className="text-on-surface-variant mt-3 font-medium tracking-widest text-xs uppercase">
              Browser, gateway, identity, services, persistence, events
            </p>
          </div>

          <div className="grid grid-cols-3 gap-3 w-full lg:w-auto">
            {[
              { label: 'Services', value: '09', color: 'text-primary' },
              { label: 'Healthy', value: '07', color: 'text-tertiary' },
              { label: 'Watch', value: '02', color: 'text-error' },
            ].map((stat) => (
              <div key={stat.label} className="bg-surface-container-lowest border border-outline-variant/10 rounded-2xl px-5 py-4 shadow-sm min-w-28">
                <div className="text-[9px] font-black uppercase tracking-[0.2em] text-on-surface-variant/50">{stat.label}</div>
                <div className={`mt-1 text-2xl font-black italic tracking-tighter ${stat.color}`}>{stat.value}</div>
              </div>
            ))}
          </div>
        </section>

        <section className="bg-surface-container-lowest rounded-[2rem] border border-outline-variant/15 shadow-xl relative overflow-hidden">
          <div className="absolute inset-0 pointer-events-none opacity-[0.025]" style={{ backgroundImage: 'radial-gradient(#000 1px, transparent 1px)', backgroundSize: '32px 32px' }}></div>
          <div className="p-4 bg-surface-container-low border-b border-outline-variant/10 flex items-center justify-between relative z-10">
            <div className="flex gap-1.5">
              <div className="w-2.5 h-2.5 rounded-full bg-error/30"></div>
              <div className="w-2.5 h-2.5 rounded-full bg-amber-500/30"></div>
              <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/30"></div>
            </div>
            <span className="font-mono text-[10px] font-black text-slate-400 uppercase tracking-widest italic">
              localhost:3000 / RuntimeRoasters
            </span>
          </div>

          <div className="relative overflow-auto">
            <div className="relative min-w-[1120px] h-[880px]">
              <svg className="absolute inset-0 w-full h-full pointer-events-none" viewBox="0 0 1120 880" aria-hidden="true">
                <defs>
                  <marker id="mesh-arrow" markerWidth="10" markerHeight="10" refX="8" refY="3" orient="auto" markerUnits="strokeWidth">
                    <path d="M0,0 L0,6 L9,3 z" fill="currentColor"></path>
                  </marker>
                </defs>
                {links.map((link, index) => (
                  <path
                    key={index}
                    d={link.d}
                    fill="none"
                    stroke={link.color}
                    strokeWidth="3"
                    strokeDasharray="8 6"
                    markerEnd="url(#mesh-arrow)"
                    className="opacity-55"
                  />
                ))}
              </svg>

              <div className="absolute left-[7%] bottom-[7%] bg-on-background text-white px-8 py-5 rounded-3xl shadow-2xl border-t-4 border-primary w-80">
                <div className="flex items-center justify-between mb-4">
                  <span className="font-headline text-[10px] font-black uppercase tracking-[0.25em] text-slate-500 italic">Trace Window</span>
                  <span className="w-2.5 h-2.5 bg-tertiary-fixed rounded-full animate-pulse"></span>
                </div>
                <div className="space-y-3">
                  {['OIDC challenge accepted', 'JWT role claim mapped', 'Casbin policy checked', 'Outbox event queued'].map((item) => (
                    <div key={item} className="flex items-center gap-3 border-b border-slate-800 pb-2 last:border-0">
                      <span className="material-symbols-outlined !text-sm text-primary-fixed">check_circle</span>
                      <span className="font-mono text-[10px] font-black uppercase tracking-tight text-slate-300">{item}</span>
                    </div>
                  ))}
                </div>
              </div>

              {nodes.map((node) => (
                <MeshNodeCard key={node.id} node={node} />
              ))}
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
