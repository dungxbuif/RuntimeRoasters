'use client';

import React from 'react';
import { HandDrawnNode } from './HandDrawnNode';
import { HandDrawnEdge } from './HandDrawnEdge';
import { motion } from 'framer-motion';

interface HighFidelityTopologyProps {
  activeStep?: number;
}

export const HighFidelityTopology = ({ activeStep = -1 }: HighFidelityTopologyProps) => {
  return (
    <div className="w-full h-full relative excalidraw-bg border-2 border-outline-variant/10 rounded-[3rem] overflow-hidden bg-white shadow-inner flex items-center justify-center">
      {/* SVG Filters & Definitions */}
      <svg className="absolute w-0 h-0">
        <defs>
          <filter id="rough-edge">
            <feTurbulence baseFrequency="0.04" numOctaves="3" result="noise" type="fractalNoise"></feTurbulence>
            <feDisplacementMap in="SourceGraphic" in2="noise" scale="3" xChannelSelector="R" yChannelSelector="G"></feDisplacementMap>
          </filter>
          <marker id="arrowhead" markerHeight="7" markerWidth="10" orient="auto" refX="9" refY="3.5">
            <polygon fill="#737686" points="0 0, 10 3.5, 0 7"></polygon>
          </marker>
          <marker id="arrowhead-error" markerHeight="7" markerWidth="10" orient="auto" refX="9" refY="3.5">
            <polygon fill="#ba1a1a" points="0 0, 10 3.5, 0 7"></polygon>
          </marker>
          <marker id="arrowhead-primary" markerHeight="7" markerWidth="10" orient="auto" refX="9" refY="3.5">
            <polygon fill="#004ac6" points="0 0, 10 3.5, 0 7"></polygon>
          </marker>
        </defs>
      </svg>

      {/* Diagram Canvas */}
      <div className="relative w-full max-w-5xl aspect-video scale-110 translate-x-[-2rem]">
        
        {/* L1: Inbound */}
        <div className="absolute top-1/2 left-0 -translate-y-1/2 flex flex-col gap-32">
          <HandDrawnNode label="Browser UI" subtitle="Next.js Client" icon="web" category="primary" />
        </div>

        {/* Edges L1 -> L2 */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none" style={{ overflow: 'visible' }}>
          <HandDrawnEdge d="M 128 281 L 280 281" isActive={activeStep >= 0} label="HTTPS" />
        </svg>

        {/* L2: Gateway */}
        <div className="absolute top-1/2 left-[300px] -translate-y-1/2">
           <HandDrawnNode label="KrakenD GW" subtitle="API Gateway" icon="door_open" category="primary" isActive={activeStep >= 0} />
        </div>

        {/* Edges L2 -> L3 (Identity) */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none" style={{ overflow: 'visible' }}>
           {/* To Hydra */}
           <HandDrawnEdge d="M 428 250 Q 450 120 580 120" isActive={activeStep >= 1} />
           {/* To Kratos */}
           <HandDrawnEdge d="M 428 310 Q 450 440 580 440" isActive={activeStep >= 2} />
        </svg>

        {/* L3: Identity Services */}
        <div className="absolute top-[80px] left-[600px]">
           <HandDrawnNode label="Ory Hydra" subtitle="OAuth2/OIDC" icon="lock" category="secondary" isActive={activeStep >= 1} />
        </div>
        <div className="absolute bottom-[80px] left-[600px]">
           <HandDrawnNode label="Ory Kratos" subtitle="Identity Store" icon="person" category="secondary" isActive={activeStep >= 2} />
        </div>

        {/* Edges L3 -> L4 (Backend) */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none" style={{ overflow: 'visible' }}>
            <HandDrawnEdge d="M 728 120 L 880 120" isActive={activeStep >= 3} />
            <HandDrawnEdge d="M 728 440 L 880 440" isActive={activeStep >= 3} />
        </svg>

        {/* L4: Business Services */}
        <div className="absolute top-[80px] left-[900px]">
           <HandDrawnNode label="Auth Service" subtitle="RBAC Logic" icon="shield" category="infra" isActive={activeStep >= 3} />
        </div>
        <div className="absolute bottom-[80px] left-[900px]">
           <HandDrawnNode label="Farm Service" subtitle="Domain Logic" icon="agriculture" category="infra" isActive={activeStep >= 4} />
        </div>

        {/* Annotation */}
        <motion.div 
           initial={{ opacity: 0 }}
           animate={{ opacity: 1 }}
           className="absolute top-[20%] right-[-10%] font-annotation text-primary text-xl rotate-[5deg] max-w-[180px]"
        >
           Event-driven sync deferred to Phase 2
           <svg height="40" width="40" className="mt-2">
              <path d="M 0 0 Q 20 20 10 35" fill="none" stroke="#004ac6" strokeWidth="2" markerEnd="url(#arrowhead-primary)" style={{ filter: 'url(#rough-edge)' }} />
           </svg>
        </motion.div>

      </div>

      {/* Industrial Label */}
      <div className="absolute bottom-10 left-12">
         <p className="text-[10px] font-black uppercase tracking-[0.5em] text-slate-300 italic">Architecture_Narrator_V4.2</p>
      </div>
    </div>
  );
};
