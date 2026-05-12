// ─── Topology Mock Data (Horizontal Layout) ───────────────────────────────────
// Horizontal flow:  Client → [Gateway + Services frame] → [per-service: DB / Cache / Kafka]
//
// Coordinate system: x%, y% of the diagram container
// Approximate container: ~1400px wide × 680px tall

import type { FlowScenario, GroupFrame, TopologyEdge, TopologyNode } from '@/types/topology';

// ─── Group Frames ─────────────────────────────────────────────────────────────

export const TOPOLOGY_FRAMES: GroupFrame[] = [
  {
    id: 'frame-internal',
    label: 'Internal Infrastructure',
    x: 19,
    y: 3,
    width: 59,
    height: 94,
  },
  {
    id: 'frame-infra',
    label: 'Persistence & Messaging',
    x: 79,
    y: 3,
    width: 19,
    height: 94,
  },
];

// ─── Nodes ────────────────────────────────────────────────────────────────────
// Horizontal zones:
//   5%  → Client
//   28% → Gateway (KrakenD)
//   47% → Identity top row (Hydra, Kratos)
//   47% → Services mid/lower rows (Auth, Farm, Warehouse)
//   82% → Per-service infra (PG, Valkey, Kafka)

export const TOPOLOGY_NODES: TopologyNode[] = [
  // ── Client ──────────────────────────────────────────────────
  {
    id: 'browser',
    label: 'Browser',
    subtitle: 'React Client',
    icon: 'devices',
    position: { x: 9, y: 47 },
    category: 'client',
  },

  // ── Gateway ─────────────────────────────────────────────────
  {
    id: 'krakend',
    label: 'KrakenD',
    subtitle: 'API Gateway :8081',
    icon: 'router',
    position: { x: 31, y: 47 },
    category: 'gateway',
  },

  // ── Identity ────────────────────────────────────────────────
  {
    id: 'hydra',
    label: 'Ory Hydra',
    subtitle: 'OAuth2 / OIDC',
    icon: 'passkey',
    position: { x: 54, y: 18 },
    category: 'identity',
  },
  {
    id: 'kratos',
    label: 'Ory Kratos',
    subtitle: 'Identity Mgmt',
    icon: 'manage_accounts',
    position: { x: 68, y: 18 },
    category: 'identity',
  },

  // ── Services ────────────────────────────────────────────────
  {
    id: 'auth-service',
    label: 'Auth Service',
    subtitle: ':8082',
    icon: 'shield',
    position: { x: 54, y: 42 },
    category: 'service',
  },
  {
    id: 'farm-service',
    label: 'Farm Service',
    subtitle: ':8083',
    icon: 'agriculture',
    position: { x: 54, y: 62 },
    category: 'service',
  },
  {
    id: 'warehouse-service',
    label: 'Warehouse',
    subtitle: 'Inventory Svc',
    icon: 'warehouse',
    position: { x: 54, y: 82 },
    category: 'service',
  },

  // ── Per-service Infra ────────────────────────────────────────
  // Auth infra
  {
    id: 'pg-auth',
    label: 'PostgreSQL',
    subtitle: 'auth_db',
    icon: 'database',
    position: { x: 83, y: 35 },
    category: 'database',
  },
  {
    id: 'redis',
    label: 'Valkey',
    subtitle: 'JWT blacklist',
    icon: 'memory',
    position: { x: 83, y: 50 },
    category: 'cache',
  },

  // Farm infra
  {
    id: 'pg-farm',
    label: 'PostgreSQL',
    subtitle: 'farm_db',
    icon: 'database',
    position: { x: 83, y: 62 },
    category: 'database',
  },
  {
    id: 'kafka-farm',
    label: 'Kafka',
    subtitle: 'farm.created',
    icon: 'dynamic_feed',
    position: { x: 92, y: 62 },
    category: 'queue',
  },

  // Warehouse infra
  {
    id: 'pg-warehouse',
    label: 'PostgreSQL',
    subtitle: 'warehouse_db',
    icon: 'database',
    position: { x: 83, y: 82 },
    category: 'database',
  },
  {
    id: 'kafka-warehouse',
    label: 'Kafka',
    subtitle: 'inventory.*',
    icon: 'dynamic_feed',
    position: { x: 92, y: 82 },
    category: 'queue',
  },
];

// ─── Edges ────────────────────────────────────────────────────────────────────

export const TOPOLOGY_EDGES: TopologyEdge[] = [
  // Client → Gateway
  { id: 'e-browser-gw', source: 'browser', target: 'krakend', label: 'REST :8081' },

  // Gateway → Identity & Services
  { id: 'e-gw-hydra', source: 'krakend', target: 'hydra', label: 'OAuth2' },
  { id: 'e-gw-auth', source: 'krakend', target: 'auth-service', label: 'gRPC' },
  { id: 'e-gw-farm', source: 'krakend', target: 'farm-service', label: 'gRPC' },
  { id: 'e-gw-warehouse', source: 'krakend', target: 'warehouse-service', label: 'gRPC' },

  // Identity internal
  { id: 'e-hydra-kratos', source: 'hydra', target: 'kratos', label: 'login/consent', style: 'dashed' },
  { id: 'e-kratos-auth', source: 'kratos', target: 'auth-service', label: 'accept', style: 'dashed' },

  // Auth → infra
  { id: 'e-auth-pg', source: 'auth-service', target: 'pg-auth', label: 'SQL' },
  { id: 'e-auth-redis', source: 'auth-service', target: 'redis', label: 'blacklist' },

  // Farm → infra
  { id: 'e-farm-pg', source: 'farm-service', target: 'pg-farm', label: 'SQL' },
  { id: 'e-farm-kafka', source: 'farm-service', target: 'kafka-farm', label: 'outbox' },

  // Warehouse → infra
  { id: 'e-warehouse-pg', source: 'warehouse-service', target: 'pg-warehouse', label: 'SQL' },
  { id: 'e-warehouse-kafka', source: 'warehouse-service', target: 'kafka-warehouse', label: 'outbox' },
];

// ─── Flow Scenarios ───────────────────────────────────────────────────────────

export const FLOW_SCENARIOS: FlowScenario[] = [
  {
    id: 'oidc-login',
    name: 'OIDC Login Flow',
    description: 'Authorization Code: Browser → KrakenD → Hydra → Kratos → Auth → PG + Valkey',
    steps: [
      { edgeId: 'e-browser-gw', stepNumber: 1, description: 'User clicks Login → Browser sends GET /oauth2/auth to KrakenD.', durationMs: 1500 },
      { edgeId: 'e-gw-hydra', stepNumber: 2, description: 'KrakenD forwards OAuth2 Authorization request to Ory Hydra.', durationMs: 1500 },
      { edgeId: 'e-hydra-kratos', stepNumber: 3, description: 'Hydra delegates credential verification to Ory Kratos.', durationMs: 1500 },
      { edgeId: 'e-kratos-auth', stepNumber: 4, description: 'Kratos triggers login-accept via Auth Service gRPC endpoint.', durationMs: 1500 },
      { edgeId: 'e-auth-pg', stepNumber: 5, description: 'Auth Service persists session + audit log to PostgreSQL (auth_db).', durationMs: 1500 },
      { edgeId: 'e-auth-redis', stepNumber: 6, description: 'JWT fingerprint cached in Valkey for fast revocation checks. Done ✓', durationMs: 1500 },
    ],
  },
  {
    id: 'farm-create',
    name: 'Create Farm Flow',
    description: 'Farm CRUD: Browser → KrakenD → Farm Service → PG + Kafka Outbox',
    steps: [
      { edgeId: 'e-browser-gw', stepNumber: 1, description: 'Client sends POST /v1/farms with JWT Bearer token to Gateway.', durationMs: 1500 },
      { edgeId: 'e-gw-farm', stepNumber: 2, description: 'KrakenD verifies JWT scope, routes gRPC call to Farm Service.', durationMs: 1500 },
      { edgeId: 'e-farm-pg', stepNumber: 3, description: 'Farm Service writes Farm + Outbox event in one DB transaction.', durationMs: 1500 },
      { edgeId: 'e-farm-kafka', stepNumber: 4, description: 'Outbox Relay publishes farm.created to Kafka. Saga starts.', durationMs: 1500 },
    ],
  },
  {
    id: 'warehouse-saga',
    name: 'Warehouse Saga',
    description: 'Order saga: Browser → KrakenD → Warehouse → PG + Kafka',
    steps: [
      { edgeId: 'e-browser-gw', stepNumber: 1, description: 'Client places order — arrives at KrakenD Gateway.', durationMs: 1500 },
      { edgeId: 'e-gw-warehouse', stepNumber: 2, description: 'Gateway routes to Warehouse Service for stock reservation.', durationMs: 1500 },
      { edgeId: 'e-warehouse-pg', stepNumber: 3, description: 'Warehouse atomically reserves stock + writes Outbox event.', durationMs: 1500 },
      { edgeId: 'e-warehouse-kafka', stepNumber: 4, description: 'Outbox Relay publishes InventoryReserved → Kafka. Saga continues.', durationMs: 1500 },
    ],
  },
];

export const TOPOLOGY_HISTORY_KEY = 'rr_topology_history' as const;
