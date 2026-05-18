'use client';

import React, { useMemo } from 'react';
import dagre from 'dagre';
import ReactFlow, {
  Controls,
  Edge,
  Handle,
  MarkerType,
  Node,
  NodeProps,
  Position,
} from 'reactflow';
import 'reactflow/dist/style.css';

type AppNodeData = {
  title: string;
  subtitle: string;
  icon: string;
  tone: 'touchpoint' | 'gateway' | 'service' | 'db' | 'infra';
  width: number;
  height: number;
};

const toneStyles: Record<AppNodeData['tone'], { card: string; iconWrap: string; icon: string }> = {
  touchpoint: { card: 'bg-white border-slate-300', iconWrap: 'bg-slate-100', icon: 'text-slate-700' },
  gateway: { card: 'bg-blue-50 border-blue-400', iconWrap: 'bg-blue-600', icon: 'text-white' },
  service: { card: 'bg-violet-50 border-violet-300', iconWrap: 'bg-violet-100', icon: 'text-violet-700' },
  db: { card: 'bg-cyan-50 border-cyan-300', iconWrap: 'bg-cyan-100', icon: 'text-cyan-700' },
  infra: { card: 'bg-emerald-50 border-emerald-300', iconWrap: 'bg-emerald-100', icon: 'text-emerald-700' },
};

function ArchitectureNode({ data }: NodeProps<AppNodeData>) {
  const tone = toneStyles[data.tone];
  return (
    <div className={`relative h-full w-full rounded-xl border p-3 shadow-sm ${tone.card}`}>
      <Handle type="target" position={Position.Left} style={{ width: 8, height: 8, opacity: 0 }} />
      <Handle type="source" position={Position.Right} style={{ width: 8, height: 8, opacity: 0 }} />
      <div className="flex items-start gap-2.5">
        <div className={`h-8 w-8 rounded-md flex items-center justify-center ${tone.iconWrap}`}>
          <span className={`material-symbols-outlined text-[18px] ${tone.icon}`}>{data.icon}</span>
        </div>
        <div>
          <div className="text-[12px] font-semibold leading-tight text-slate-900">{data.title}</div>
          <div className="mt-1 text-[10px] text-slate-600 leading-tight">{data.subtitle}</div>
        </div>
      </div>
    </div>
  );
}

const nodeTypes = { architectureNode: ArchitectureNode };

const baseNodes: Node[] = [
  {
    id: 'client',
    type: 'architectureNode',
    data: { title: 'Client App', subtitle: 'Next.js', icon: 'devices', tone: 'touchpoint', width: 186, height: 82 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'kratos',
    type: 'architectureNode',
    data: { title: 'Ory Kratos', subtitle: 'Identity', icon: 'fingerprint', tone: 'touchpoint', width: 126, height: 50 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'hydra',
    type: 'architectureNode',
    data: { title: 'Ory Hydra', subtitle: 'OAuth2 / OIDC', icon: 'lock', tone: 'touchpoint', width: 126, height: 50 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'gateway',
    type: 'architectureNode',
    data: { title: 'KrakenD Gateway', subtitle: 'API Gateway', icon: 'router', tone: 'gateway', width: 198, height: 84 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'auth',
    type: 'architectureNode',
    data: { title: 'Auth Service', subtitle: 'Authorization', icon: 'passkey', tone: 'service', width: 132, height: 60 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'farm',
    type: 'architectureNode',
    data: { title: 'Farm Service', subtitle: 'Farm Manager', icon: 'agriculture', tone: 'service', width: 132, height: 60 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'warehouse',
    type: 'architectureNode',
    data: { title: 'Warehouse Service', subtitle: 'worker', icon: 'warehouse', tone: 'service', width: 132, height: 60 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'authdb',
    type: 'architectureNode',
    data: { title: 'PostgreSQL', subtitle: 'auth_db', icon: 'database', tone: 'db', width: 118, height: 54 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'farmdb',
    type: 'architectureNode',
    data: { title: 'PostgreSQL', subtitle: 'farm_db', icon: 'database', tone: 'db', width: 118, height: 54 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'waredb',
    type: 'architectureNode',
    data: { title: 'PostgreSQL', subtitle: 'warehouse_db', icon: 'database', tone: 'db', width: 118, height: 54 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'kafka',
    type: 'architectureNode',
    data: { title: 'Kafka', subtitle: 'Event Stream', icon: 'hub', tone: 'infra', width: 164, height: 72 },
    position: { x: 0, y: 0 },
  },
  {
    id: 'valkey',
    type: 'architectureNode',
    data: { title: 'Valkey/Redis', subtitle: 'Cache', icon: 'bolt', tone: 'infra', width: 164, height: 72 },
    position: { x: 0, y: 0 },
  },
];

const baseEdges: Edge[] = [
  { id: 'e1', source: 'client', target: 'gateway' },
  { id: 'e4', source: 'gateway', target: 'auth' },
  { id: 'e5', source: 'gateway', target: 'farm' },
  { id: 'e6', source: 'gateway', target: 'warehouse' },
  { id: 'e7', source: 'auth', target: 'authdb' },
  { id: 'e8', source: 'farm', target: 'farmdb' },
  { id: 'e9', source: 'warehouse', target: 'waredb' },
  { id: 'e10', source: 'micro-anchor', target: 'infra-anchor' },
];

function layout(nodes: Node[], edges: Edge[]): Node[] {
  const g = new dagre.graphlib.Graph();
  g.setDefaultEdgeLabel(() => ({}));
  g.setGraph({ rankdir: 'LR', ranksep: 96, nodesep: 52, marginx: 28, marginy: 24 });

  nodes.forEach((n) => {
    const data = n.data as AppNodeData;
    g.setNode(n.id, { width: data.width, height: data.height });
  });
  edges.forEach((e) => g.setEdge(e.source, e.target));

  dagre.layout(g);

  return nodes.map((n) => {
    const p = g.node(n.id);
    const data = n.data as AppNodeData;
    return {
      ...n,
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
      position: { x: p.x - data.width / 2, y: p.y - data.height / 2 },
      style: { width: data.width, height: data.height },
    };
  });
}

function bounds(nodes: Node[]) {
  const xs = nodes.map((n) => n.position.x);
  const ys = nodes.map((n) => n.position.y);
  const x2 = nodes.map((n) => n.position.x + Number(n.style?.width ?? 0));
  const y2 = nodes.map((n) => n.position.y + Number(n.style?.height ?? 0));
  return {
    x: Math.min(...xs),
    y: Math.min(...ys),
    xMax: Math.max(...x2),
    yMax: Math.max(...y2),
  };
}

function rectOf(node: Node) {
  return {
    x: node.position.x,
    y: node.position.y,
    w: Number(node.style?.width ?? 0),
    h: Number(node.style?.height ?? 0),
  };
}

function buildFramedNodes(layouted: Node[]): Node[] {
  const adjusted = layouted.map((n) => ({ ...n, position: { ...n.position } }));
  const byId = new Map(adjusted.map((n) => [n.id, n]));
  const clientIds = ['client'];
  const microIds = ['kratos', 'hydra', 'auth', 'farm', 'warehouse', 'authdb', 'farmdb', 'waredb'];
  const infraIds = ['kafka', 'valkey'];

  const pick = (ids: string[]) => ids.map((id) => byId.get(id)).filter(Boolean) as Node[];
  const microNodes = pick(microIds);
  const gatewayNode = byId.get('gateway');

  const padX = 32;
  const padY = 30;

  const m = bounds(microNodes);

  // Place Ory block right above KrakenD.
  const kratos = byId.get('kratos');
  const hydra = byId.get('hydra');
  if (gatewayNode && kratos && hydra) {
    const gRect = rectOf(gatewayNode);
    const kRect = rectOf(kratos);
    const hRect = rectOf(hydra);
    const startX = gRect.x + gRect.w / 2 - kRect.w / 2;
    const topY = gRect.y - (kRect.h + hRect.h + 30);
    kratos.position.x = startX;
    kratos.position.y = topY;
    hydra.position.x = startX;
    hydra.position.y = topY + kRect.h + 12;
    byId.set('kratos', kratos);
    byId.set('hydra', hydra);
  }

  // Keep Kafka + Valkey on one vertical column, right side of microservices.
  const kafka = byId.get('kafka');
  const valkey = byId.get('valkey');
  if (kafka && valkey) {
    const kafkaH = Number(kafka.style?.height ?? 0);
    const colCenterY = m.y + (m.yMax - m.y) / 2;
    const baseX = m.xMax + 180;
    kafka.position.x = baseX;
    kafka.position.y = colCenterY - kafkaH - 16;
    valkey.position.x = baseX;
    valkey.position.y = colCenterY + 16;
    // Keep pick() references updated for later bounds calc.
    byId.set('kafka', kafka);
    byId.set('valkey', valkey);
  }

  const i = bounds(pick(infraIds));
  const mAfter = bounds(pick(microIds));
  const clientAfter = bounds(pick(clientIds));
  const g = gatewayNode ? rectOf(gatewayNode) : null;
  const microLeft = g ? g.x + g.w / 2 : m.x - padX;
  const microWidthAfter = mAfter.xMax - microLeft + padX + 56;
  const clientTopRaw = clientAfter.y - padY;
  const microTopRaw = mAfter.y - padY - 20;
  const infraTopRaw = i.y - padY;
  const commonTop = Math.min(clientTopRaw, microTopRaw, infraTopRaw);
  const microTop = commonTop;
  const microHeight = mAfter.yMax - commonTop + padY + 20;
  const clientHeight = Math.max(clientAfter.yMax - commonTop + 14, 210);
  const infraHeight = Math.max(i.yMax - commonTop + 14, 230);

  const frameNodes: Node[] = [
    {
      id: 'frame-client',
      type: 'group',
      data: { label: 'Client Zone' },
      position: { x: clientAfter.x - padX, y: commonTop },
      zIndex: -1,
      style: {
        width: clientAfter.xMax - clientAfter.x + padX * 2,
        height: clientHeight,
        border: '1px dashed #94a3b8',
        borderRadius: 12,
        background: 'transparent',
      },
      draggable: false,
      selectable: false,
    },
    {
      id: 'frame-micro',
      type: 'group',
      data: { label: 'Microservices' },
      position: { x: microLeft, y: microTop },
      zIndex: -1,
      style: {
        width: microWidthAfter,
        height: microHeight,
        border: '1px dashed #60a5fa',
        borderRadius: 12,
        background: 'transparent',
      },
      draggable: false,
      selectable: false,
    },
    {
      id: 'frame-infra',
      type: 'group',
      data: { label: 'Kafka / Redis' },
      position: { x: i.x - padX, y: commonTop },
      zIndex: -1,
      style: {
        width: i.xMax - i.x + padX * 2,
        height: infraHeight,
        border: '1px dashed #34d399',
        borderRadius: 12,
        background: 'transparent',
      },
      draggable: false,
      selectable: false,
    },
  ];

  const pairDefs = [
    { id: 'pair-auth', label: 'Auth + DB', a: 'auth', b: 'authdb' },
    { id: 'pair-farm', label: 'Farm + DB', a: 'farm', b: 'farmdb' },
    { id: 'pair-warehouse', label: 'Warehouse + DB', a: 'warehouse', b: 'waredb' },
  ] as const;

  const pairFrames: Node[] = pairDefs
    .map((pair) => {
      const n1 = byId.get(pair.a);
      const n2 = byId.get(pair.b);
      if (!n1 || !n2) return null;
      const b = bounds([n1, n2]);
      return {
        id: pair.id,
        type: 'group',
        data: { label: pair.label },
        position: { x: b.x - 12, y: b.y - 10 },
        zIndex: -1,
        style: {
          width: b.xMax - b.x + 24,
          height: b.yMax - b.y + 20,
          border: '1px solid #dbeafe',
          borderRadius: 10,
          background: 'rgba(239,246,255,0.32)',
        },
        draggable: false,
        selectable: false,
      } as Node;
    })
    .filter(Boolean) as Node[];

  const oryFrame = (() => {
    const k = byId.get('kratos');
    const h = byId.get('hydra');
    if (!k || !h) return null;
    const b = bounds([k, h]);
    return {
      id: 'pair-ory',
      type: 'group',
      data: { label: 'ORY (Exposed)' },
      position: { x: b.x - 8, y: b.y - 8 },
      zIndex: -1,
      style: {
        width: b.xMax - b.x + 16,
        height: b.yMax - b.y + 16,
        border: '1px solid #fde68a',
        borderRadius: 10,
        background: 'rgba(254,249,195,0.32)',
      },
      draggable: false,
      selectable: false,
    } as Node;
  })();

  const anchorNodes: Node[] = [
    {
      id: 'micro-anchor',
      position: { x: microLeft + microWidthAfter - 2, y: microTop + microHeight / 2 - 2 },
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
      style: { width: 4, height: 4, opacity: 0, border: '0px', background: 'transparent' },
      draggable: false,
      selectable: false,
      deletable: false,
    },
    {
      id: 'infra-anchor',
      position: { x: i.x - padX + 2, y: i.y - padY + (i.yMax - i.y + padY * 2) / 2 - 2 },
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
      style: { width: 4, height: 4, opacity: 0, border: '0px', background: 'transparent' },
      draggable: false,
      selectable: false,
      deletable: false,
    },
  ];

  return [...frameNodes, ...(oryFrame ? [oryFrame] : []), ...pairFrames, ...anchorNodes, ...adjusted];
}

export function ArchitectureDiagramCanvas() {
  const nodes = useMemo(() => buildFramedNodes(layout(baseNodes, baseEdges)), []);
  const edges = useMemo(
    () =>
      baseEdges.map((e) => ({
        ...e,
        type: 'smoothstep',
        animated: true,
        zIndex: 50,
        style: { stroke: '#334155', strokeWidth: 2.2 },
        markerEnd: { type: MarkerType.ArrowClosed, color: '#334155', width: 18, height: 18 },
      })),
    []
  );

  return (
    <div className="h-full w-full">
      <ReactFlow
        fitView
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        defaultMarkerColor="#334155"
        defaultEdgeOptions={{ zIndex: 50 }}
      >
        <Controls />
      </ReactFlow>
    </div>
  );
}

export function ArchitectureTopology() {
  return <ArchitectureDiagramCanvas />;
}
