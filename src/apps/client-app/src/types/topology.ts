export type TopologyTone = 'touchpoint' | 'gateway' | 'identity' | 'service' | 'db' | 'infra';

export interface TopologyNode {
  id: string;
  label: string;
  subtitle: string;
  icon: string;
  tone?: TopologyTone;
  position?: { x: number; y: number };
  category?: string;
}

export interface TopologyEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  pattern?: string;
  style?: 'solid' | 'dashed';
}

export interface TopologyFlow {
  id: string;
  name: string;
  description: string;
  node_ids: string[];
  edge_ids: string[];
  last_seen_at?: string | null;
}

export interface TopologyConfig {
  nodes: TopologyNode[];
  edges: TopologyEdge[];
  flows: TopologyFlow[];
}

export interface TopologyHistoryEntry {
  event_id: string;
  topic: string;
  status: string;
  flow_id: string;
  node_id: string;
  edge_id: string;
  pattern: string;
  source_service: string;
  visibility: 'public' | 'private';
  trace_id?: string;
  occurred_at: string;
  display_payload?: Record<string, unknown>;
}

export interface GroupFrame {
  id: string;
  label: string;
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface FlowStep {
  edgeId: string;
  stepNumber: number;
  description: string;
  durationMs: number;
}

export interface FlowScenario {
  id: string;
  name: string;
  description: string;
  steps: FlowStep[];
}
