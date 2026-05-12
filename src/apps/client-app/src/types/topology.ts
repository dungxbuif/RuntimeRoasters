// ─── Topology Type Definitions ───────────────────────────────────────────────

export type NodeCategory =
  | 'client'
  | 'gateway'
  | 'identity'
  | 'service'
  | 'database'
  | 'cache'
  | 'queue';

export interface NodePosition {
  x: number; // percent of container width (0–100)
  y: number; // percent of container height (0–100)
}

export interface TopologyNode {
  id: string;
  label: string;
  subtitle: string;
  icon: string; // Material Symbols Outlined name
  position: NodePosition;
  category: NodeCategory;
}

export interface GroupFrame {
  id: string;
  label: string;
  // All in % of container
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface TopologyEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  style?: 'solid' | 'dashed';
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

export interface FlowHistoryEntry {
  scenarioId: string;
  scenarioName: string;
  timestamp: number;
}

export type AnimationState = 'idle' | 'playing' | 'paused' | 'done';

export interface TopologyState {
  activeScenarioId: string | null;
  animState: AnimationState;
  currentStepIndex: number;
  activeEdgeIds: string[];
  activeNodeIds: string[];
  history: FlowHistoryEntry[];
}
