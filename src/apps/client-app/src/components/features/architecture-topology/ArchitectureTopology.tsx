'use client';

import React, { useEffect, useMemo, useState } from 'react';
import dagre from 'dagre';
import ReactFlow, { Controls, Edge, Handle, MarkerType, Node, NodeProps, Position } from 'reactflow';
import 'reactflow/dist/style.css';
import { topologyService } from '@/services/topology.service';
import { TopologyConfig, TopologyHistoryEntry, TopologyNode, TopologyTone } from '@/types/topology';

type AppNodeData = Omit<TopologyNode, 'tone'> & { tone: TopologyTone; active: boolean; recent: boolean; width: number; height: number };

const toneStyles: Record<string, { card: string; iconWrap: string; icon: string }> = {
  touchpoint: { card: 'bg-white border-slate-300', iconWrap: 'bg-slate-100', icon: 'text-slate-700' },
  gateway: { card: 'bg-blue-50 border-blue-400', iconWrap: 'bg-blue-600', icon: 'text-white' },
  identity: { card: 'bg-amber-50 border-amber-300', iconWrap: 'bg-amber-100', icon: 'text-amber-700' },
  service: { card: 'bg-violet-50 border-violet-300', iconWrap: 'bg-violet-100', icon: 'text-violet-700' },
  db: { card: 'bg-cyan-50 border-cyan-300', iconWrap: 'bg-cyan-100', icon: 'text-cyan-700' },
  infra: { card: 'bg-emerald-50 border-emerald-300', iconWrap: 'bg-emerald-100', icon: 'text-emerald-700' },
};

function ArchitectureNode({ data }: NodeProps<AppNodeData>) {
  const tone = toneStyles[data.tone] || toneStyles.service;
  return (
    <div className={`relative h-full w-full rounded-lg border p-3 shadow-sm ${tone.card} ${data.recent ? 'ring-2 ring-green-500' : data.active ? 'ring-2 ring-slate-400' : 'opacity-55'}`}>
      <Handle type="target" position={Position.Left} style={{ width: 8, height: 8, opacity: 0 }} />
      <Handle type="source" position={Position.Right} style={{ width: 8, height: 8, opacity: 0 }} />
      <div className="flex items-start gap-2.5">
        <div className={`flex h-8 w-8 items-center justify-center rounded-md ${tone.iconWrap}`}>
          <span className={`material-symbols-outlined text-[18px] ${tone.icon}`}>{data.icon}</span>
        </div>
        <div className="min-w-0">
          <div className="text-[12px] font-semibold leading-tight text-slate-900">{data.label}</div>
          <div className="mt-1 text-[10px] leading-tight text-slate-600">{data.subtitle}</div>
        </div>
      </div>
    </div>
  );
}

const nodeTypes = { architectureNode: ArchitectureNode };

const fallbackConfig: TopologyConfig = {
  nodes: [
    { id: 'client.web', label: 'Client App', subtitle: 'Next.js', icon: 'devices', tone: 'touchpoint' },
    { id: 'gateway.krakend', label: 'KrakenD', subtitle: 'Gateway', icon: 'router', tone: 'gateway' },
    { id: 'service.trace', label: 'Trace Service', subtitle: 'Topology config/history', icon: 'timeline', tone: 'service' },
    { id: 'service.socket', label: 'Socket Service', subtitle: 'Realtime fanout', icon: 'sensors', tone: 'service' },
    { id: 'infra.kafka', label: 'Kafka', subtitle: 'Events', icon: 'hub', tone: 'infra' },
  ],
  edges: [
    { id: 'edge.client.gateway', source: 'client.web', target: 'gateway.krakend', label: 'REST' },
    { id: 'edge.gateway.trace', source: 'gateway.krakend', target: 'service.trace', label: 'query' },
    { id: 'edge.client.socket', source: 'client.web', target: 'service.socket', label: 'ws' },
    { id: 'edge.kafka.socket', source: 'infra.kafka', target: 'service.socket', label: 'consume' },
  ],
  flows: [
    { id: 'flow.realtime.socket-push-pull', name: 'Socket Push/Pull', description: 'Fallback realtime topology.', node_ids: ['client.web', 'service.socket', 'infra.kafka'], edge_ids: ['edge.client.socket', 'edge.kafka.socket'] },
  ],
};

function layout(nodes: Node[], edges: Edge[]): Node[] {
  const g = new dagre.graphlib.Graph();
  g.setDefaultEdgeLabel(() => ({}));
  g.setGraph({ rankdir: 'LR', ranksep: 90, nodesep: 42, marginx: 26, marginy: 24 });
  nodes.forEach((node) => g.setNode(node.id, { width: Number(node.style?.width ?? 172), height: Number(node.style?.height ?? 76) }));
  edges.forEach((edge) => g.setEdge(edge.source, edge.target));
  dagre.layout(g);
  return nodes.map((node) => {
    const point = g.node(node.id) || { x: 0, y: 0 };
    const width = Number(node.style?.width ?? 172);
    const height = Number(node.style?.height ?? 76);
    return { ...node, sourcePosition: Position.Right, targetPosition: Position.Left, position: { x: point.x - width / 2, y: point.y - height / 2 } };
  });
}

function buildGraph(config: TopologyConfig, activeFlowId: string, history: TopologyHistoryEntry[]) {
  const activeFlow = config.flows.find((flow) => flow.id === activeFlowId) || config.flows[0];
  const activeNodes = new Set(activeFlow?.node_ids || config.nodes.map((node) => node.id));
  const activeEdges = new Set(activeFlow?.edge_ids || config.edges.map((edge) => edge.id));
  const recent = history[0];
  const nodes: Node[] = config.nodes.map((node) => ({
    id: node.id,
    type: 'architectureNode',
    data: { ...node, tone: node.tone || 'service', active: activeNodes.has(node.id), recent: recent?.node_id === node.id, width: 174, height: 76 },
    position: { x: 0, y: 0 },
    style: { width: 174, height: 76 },
  }));
  const edges: Edge[] = config.edges
    .filter((edge) => activeEdges.has(edge.id) || activeNodes.has(edge.source) || activeNodes.has(edge.target))
    .map((edge) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      label: edge.label,
      type: 'smoothstep',
      animated: activeEdges.has(edge.id) || recent?.edge_id === edge.id,
      style: { stroke: recent?.edge_id === edge.id ? '#16a34a' : activeEdges.has(edge.id) ? '#334155' : '#cbd5e1', strokeWidth: recent?.edge_id === edge.id ? 3 : 2 },
      markerEnd: { type: MarkerType.ArrowClosed, color: recent?.edge_id === edge.id ? '#16a34a' : '#334155', width: 18, height: 18 },
    }));
  return { nodes: layout(nodes, edges), edges, activeFlow };
}

export function ArchitectureDiagramCanvas({ isPrivate = false }: { isPrivate?: boolean }) {
  const [config, setConfig] = useState<TopologyConfig>(fallbackConfig);
  const [activeFlowId, setActiveFlowId] = useState(fallbackConfig.flows[0].id);
  const [history, setHistory] = useState<TopologyHistoryEntry[]>([]);
  const [liveState, setLiveState] = useState<'connecting' | 'live' | 'polling'>('connecting');

  useEffect(() => {
    let cancelled = false;
    topologyService.getPublicConfig()
      .then((next) => {
        if (cancelled) return;
        setConfig(next);
        setActiveFlowId((current) => next.flows.some((flow) => flow.id === current) ? current : next.flows[0]?.id || current);
      })
      .catch(() => undefined);
    return () => { cancelled = true; };
  }, []);

  useEffect(() => {
    let cancelled = false;
    const load = () => topologyService.getPublicHistory(activeFlowId).then((events) => {
      if (!cancelled) setHistory(events);
    }).catch(() => undefined);
    load();
    const interval = window.setInterval(load, 7000);
    return () => {
      cancelled = true;
      window.clearInterval(interval);
    };
  }, [activeFlowId]);

  useEffect(() => {
    let cancelled = false;
    let socket: WebSocket | null = null;

    const connect = async () => {
      try {
        let url: string;
        if (isPrivate) {
          const ticket = await topologyService.getTicket();
          if (cancelled) return;
          url = topologyService.privateSocketURL(ticket, activeFlowId);
        } else {
          url = topologyService.publicSocketURL(activeFlowId);
        }

        socket = new WebSocket(url);
        socket.onopen = () => setLiveState('live');
        socket.onerror = () => setLiveState('polling');
        socket.onclose = () => {
          if (!cancelled) setLiveState((state) => state === 'live' ? 'polling' : state);
        };
        socket.onmessage = (message) => {
          try {
            const envelope = JSON.parse(message.data);
            if (envelope.type === 'topology.event') {
              setHistory((current) => [envelope.data as TopologyHistoryEntry, ...current].slice(0, 80));
            }
          } catch {
            setLiveState('polling');
          }
        };
      } catch {
        if (!cancelled) setLiveState('polling');
      }
    };

    connect();

    return () => {
      cancelled = true;
      if (socket) socket.close();
    };
  }, [activeFlowId, isPrivate]);

  const { nodes, edges, activeFlow } = useMemo(() => buildGraph(config, activeFlowId, history), [config, activeFlowId, history]);

  return (
    <div className="flex h-full min-h-[520px] w-full bg-white">
      <div className="w-72 shrink-0 border-r border-slate-200 bg-slate-50 p-3">
        <div className="mb-3 flex items-center justify-between">
          <div className="text-[10px] font-black uppercase tracking-widest text-slate-500">Flow</div>
          <div className={`rounded px-2 py-1 text-[9px] font-black uppercase ${liveState === 'live' ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'}`}>
            {liveState}
          </div>
        </div>
        <div className="space-y-2">
          {config.flows.map((flow) => (
            <button
              key={flow.id}
              type="button"
              onClick={() => setActiveFlowId(flow.id)}
              className={`w-full rounded-md border p-3 text-left transition-colors ${activeFlowId === flow.id ? 'border-slate-900 bg-white' : 'border-slate-200 bg-white/60 hover:bg-white'}`}
            >
              <div className="text-xs font-black text-slate-900">{flow.name}</div>
              <div className="mt-1 text-[10px] leading-snug text-slate-500">{flow.description}</div>
            </button>
          ))}
        </div>
      </div>
      <div className="min-w-0 flex-1">
        <div className="border-b border-slate-200 px-5 py-3">
          <div className="text-sm font-black text-slate-900">{activeFlow?.name}</div>
          <div className="text-[10px] font-bold uppercase tracking-widest text-slate-400">{history.length} events observed</div>
        </div>
        <div className="h-[calc(100%-52px)]">
          <ReactFlow fitView nodes={nodes} edges={edges} nodeTypes={nodeTypes} defaultEdgeOptions={{ zIndex: 50 }}>
            <Controls />
          </ReactFlow>
        </div>
      </div>
      <div className="w-80 shrink-0 border-l border-slate-200 bg-slate-50 p-3">
        <div className="mb-3 text-[10px] font-black uppercase tracking-widest text-slate-500">History</div>
        <div className="space-y-2 overflow-y-auto pr-1">
          {history.length === 0 ? (
            <div className="rounded-md border border-dashed border-slate-300 bg-white p-4 text-[10px] font-bold uppercase tracking-widest text-slate-400">No history yet</div>
          ) : history.slice(0, 20).map((event) => (
            <div key={`${event.event_id}-${event.occurred_at}`} className="rounded-md border border-slate-200 bg-white p-3">
              <div className="text-[10px] font-black text-slate-900">{event.topic}</div>
              <div className="mt-1 text-[10px] text-slate-500">{event.pattern} · {event.source_service}</div>
              <div className="mt-2 text-[9px] font-bold uppercase text-slate-400">{new Date(event.occurred_at).toLocaleTimeString()}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export function ArchitectureTopology() {
  return <ArchitectureDiagramCanvas />;
}
