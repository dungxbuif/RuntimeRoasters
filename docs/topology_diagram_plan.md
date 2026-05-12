# Control Plane Visualization — TopologyDiagram Component

## Background

Refactor the `(dashboard)/page.tsx` landing page from a static "Showcase Terminal" into an interactive **Topology Dashboard** with a full-screen microservices diagram. The diagram will use custom components with **Framer Motion** animations (already installed) instead of React Flow (not installed, avoiding new dependency).

The existing design system uses a **light Industrial Premium** palette with Material Design 3 tokens (blue primary `#004ac6`, emerald tertiary `#006242`, soft grays). Typography uses JetBrains Mono / Roboto Mono / Gaegu / Indie Flower. All existing patterns (OOP Services, No Magic Strings, Atomic Design) are preserved.

---

## Landing Page Wireframe

### Full Page Layout

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│  HEADER (sticky, h-16, glass bg)                                                │
│  ┌─────────────────────────────────┐                    ┌─────────────────────┐ │
│  │ ☕ RUNTIME ROASTERS              │                    │ 🔑 Access Terminal  │ │
│  │    (logo + brand, blue-600)     │                    │  (or User avatar    │ │
│  │                                 │                    │   if logged in)     │ │
│  └─────────────────────────────────┘                    └─────────────────────┘ │
├──────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  HERO ZONE (compact, ~80px)                                                      │
│  ┌──────────────────────────────────────────────────────────────────────────────┐│
│  │  System Topography                              We are here! ↙ (handwriting)││
│  │  「Visualize how user intent flows through       (Indie Flower annotation)   ││
│  │    the RuntimeRoasters infrastructure stack.」                               ││
│  └──────────────────────────────────────────────────────────────────────────────┘│
│                                                                                  │
│  ┌───────────────────────────────────────────────────┬──────────────────────────┐│
│  │                                                   │                          ││
│  │           TOPOLOGY DIAGRAM (~85vh)                │  FLOW CONTROL PANEL      ││
│  │           (relative container)                    │  (w-72, right side)      ││
│  │                                                   │                          ││
│  │    ┌───────┐                                      │  ┌────────────────────┐  ││
│  │    │🌐 UI  │              TOUCHPOINTS             │  │ DEMO SCENARIOS     │  ││
│  │    │Browser│                                      │  │                    │  ││
│  │    └───┬───┘                                      │  │ ● OIDC Login Flow  │  ││
│  │        │                                          │  │ ○ Casbin Sync Flow │  ││
│  │        ▼                                          │  │ ○ Order Saga Flow  │  ││
│  │   ┌─────────┐                                     │  │                    │  ││
│  │   │🚪KrakenD│          GATEWAY                    │  └────────────────────┘  ││
│  │   │ Gateway │                                     │                          ││
│  │   └──┬──┬───┘                                     │  ┌────────────────────┐  ││
│  │      │  │                                         │  │ PLAYBACK           │  ││
│  │      │  └──────────────┐                          │  │ [▶ Play] [⏸] [↺]  │  ││
│  │      ▼                 ▼                          │  │ Step 3/6 ████░░    │  ││
│  │  ┌────────┐      ┌──────────┐     IDENTITY        │  └────────────────────┘  ││
│  │  │🔐 Hydra│─────▶│👤 Kratos │                     │                          ││
│  │  │  OIDC  │      │ Identity │                     │  ┌────────────────────┐  ││
│  │  └───┬────┘      └──────────┘                     │  │ HISTORY            │  ││
│  │      │                                            │  │                    │  ││
│  │      ▼                                            │  │ 22:45 OIDC Login   │  ││
│  │  ┌─────────┐  ┌──────────┐  ┌───────────┐         │  │ 22:40 Casbin Sync  │  ││
│  │  │🛡️ Auth  │  │🌿 Farm   │  │📦Warehouse│ SERVICES│  │ 22:35 OIDC Login   │  ││
│  │  │ Service │  │ Service  │  │  Service  │         │  └────────────────────┘  ││
│  │  └──┬──────┘  └────┬─────┘  └─────┬─────┘         │                          ││
│  │     │              │              │               │                          ││
│  │     ▼              ▼              ▼               │                          ││
│  │  ┌──────┐   ┌─────────┐    ┌────────┐   INFRA     │                          ││
│  │  │🗄️ PG │   │📨 Kafka │    │⚡Redis │             │                          ││
│  │  │ SQL  │   │ Cluster │    │ Cache  │             │                          ││
│  │  └──────┘   └─────────┘    └────────┘             │                          ││
│  │                                                   │                          ││
│  │  ┌──────────────────────────────────────────┐     │                          ││
│  │  │ TRACE LOG (bottom, h-32, terminal style) │     │                          ││
│  │  │ ❶ 22:51:03 Browser → KrakenD: POST /auth │     │                          ││
│  │  │ ❷ 22:51:03 KrakenD → Hydra: Forward req  │     │                          ││
│  │  │ ❸ 22:51:04 Hydra → Kratos: Verify creds  │     │                          ││
│  │  │ ▌ (cursor blink)                          │     │                          ││
│  │  └──────────────────────────────────────────┘     │                          ││
│  └───────────────────────────────────────────────────┴──────────────────────────┘│
│                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

### Node Position Map (% coordinates in diagram container)

```
  Layer 0: TOUCHPOINT
  ──────────────────────────────────────────────
         Browser [50%, 5%]

  Layer 1: GATEWAY  
  ──────────────────────────────────────────────
         KrakenD [50%, 20%]

  Layer 2: IDENTITY
  ──────────────────────────────────────────────
    Hydra [30%, 35%]     Kratos [55%, 35%]

  Layer 3: SERVICES
  ──────────────────────────────────────────────
    Auth [18%, 52%]   Farm [43%, 52%]   Warehouse [68%, 52%]

  Layer 4: INFRASTRUCTURE
  ──────────────────────────────────────────────
    PostgreSQL [18%, 72%]  Kafka [43%, 72%]  Redis [68%, 72%]
```

### Edge Connections

| Edge ID | From → To | Label |
|:---|:---|:---|
| `e-browser-gw` | Browser → KrakenD | REST :8081 |
| `e-gw-hydra` | KrakenD → Hydra | OAuth2 |
| `e-gw-auth` | KrakenD → Auth Service | gRPC+JWT |
| `e-gw-farm` | KrakenD → Farm Service | gRPC+JWT |
| `e-hydra-kratos` | Hydra → Kratos | Login/Consent |
| `e-auth-pg` | Auth Service → PostgreSQL | SQL |
| `e-farm-pg` | Farm Service → PostgreSQL | SQL |
| `e-warehouse-pg` | Warehouse → PostgreSQL | SQL |
| `e-auth-kafka` | Auth Service → Kafka | Outbox |
| `e-warehouse-kafka` | Warehouse → Kafka | Outbox/Inbox |
| `e-auth-redis` | Auth Service → Redis | Blacklist |
| `e-farm-kafka` | Farm Service → Kafka | Outbox |

### Demo Scenario: "OIDC Login Flow" (6 steps)

```
 Step ❶  Browser → KrakenD
         "User clicks Login. Browser sends GET /oauth2/auth to Gateway."

 Step ❷  KrakenD → Hydra
         "Gateway forwards OAuth2 Authorization request to Hydra."

 Step ❸  Hydra → Kratos
         "Hydra delegates credential verification to Kratos Identity."

 Step ❹  Kratos → Auth Service
         "Kratos triggers login-accept via Auth Service gRPC endpoint."

 Step ❺  Auth Service → PostgreSQL
         "Auth Service persists session and writes audit log to DB."

 Step ❻  Auth Service → Redis
         "JWT token signature cached in Redis for fast revocation checks."
```

### Visual Style Reference

| Element | Style |
|:---|:---|
| **Node (touchpoint)** | White card, `secondary-container` icon bg, rounded-xl, subtle shadow |
| **Node (gateway)** | White card, `primary/20` border glow, `primary-fixed` bg tint |
| **Node (service)** | White card, `tertiary/10` left accent bar, hover→lift |
| **Node (infra)** | Smaller card, `surface-container` bg, DB stack lines for PostgreSQL |
| **Node (active)** | Pulsing `primary` border glow + scale(1.02) |
| **Edge (default)** | Dashed, `outline-variant` color, 1.5px |
| **Edge (active)** | Solid `primary` color, 2.5px, animated dash-offset flow |
| **Step badge** | `primary-container` circle, white number, appears on edge midpoint |
| **Trace Log** | `slate-900` bg, monospace font, terminal aesthetic |
| **Control Panel** | `surface-container-low` bg, rounded-2xl, cards inside |

---

## Proposed Changes

### Data Layer

#### [NEW] [topology.types.ts](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/types/topology.ts)

Define all TypeScript interfaces:
- `TopologyNode` — id, label, subtitle, icon (Material Symbol name), position `{x, y}`, category (`touchpoint | gateway | service | infra`), color token
- `TopologyEdge` — id, source, target, label
- `FlowScenario` — id, name, description, steps array of `FlowStep`
- `FlowStep` — edgeId, stepNumber, description, durationMs
- `FlowHistoryEntry` — scenario id, timestamp

#### [NEW] [topology.data.ts](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/constants/topology.ts)

Mock data constants mapped to the actual RuntimeRoasters architecture:
- **Nodes** (10): Browser/UI, KrakenD Gateway, Ory Hydra, Ory Kratos, Auth Service, Farm Service, Warehouse Service, Kafka Cluster, PostgreSQL, Redis/Valkey
- **Edges**: Static connections between all nodes
- **1 Demo Scenario**: "OIDC Login Flow" — 6 steps tracing `Browser → KrakenD → Hydra → Kratos → Auth Service → PostgreSQL`, matching the existing system blueprint

---

### Component Layer

#### [NEW] [TopologyNode.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/components/features/topology/TopologyNode.tsx)

Single node component:
- Renders a card with Material Symbol icon, label, subtitle
- Color coding by `category` using existing theme tokens (blue for gateway, emerald for services, gray for infra)
- `framer-motion` entrance animation (`scale` + `opacity` on mount)
- Glow/pulse effect when node is "active" in a running flow
- Positioned absolutely within the diagram container via `position: {x, y}` percent-based coordinates

#### [NEW] [TopologyEdge.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/components/features/topology/TopologyEdge.tsx)

SVG path edge component:
- Draws a curved path between two node positions using quadratic bezier
- Default state: dashed, muted (`outline-variant` color)
- Active state: solid primary color with animated dash-offset (flowing dots effect)
- Step number badge — a numbered circle that appears on the midpoint of the edge during animation
- Uses `framer-motion` for opacity/pathLength animation

#### [NEW] [TraceLog.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/components/features/topology/TraceLog.tsx)

Small floating panel (bottom-right of diagram):
- Shows current step description as a terminal-style log line
- Auto-scrolls as new steps appear
- Uses `font-headline` (JetBrains Mono) for monospace aesthetic
- Each log entry has step number badge + timestamp + description

#### [NEW] [FlowControlPanel.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/components/features/topology/FlowControlPanel.tsx)

Sidebar panel (right side, ~280px width):
- **Scenario Selector**: List of available demo scenarios with radio-button selection
- **Playback Controls**: Play / Pause / Reset buttons
- **History Timeline**: List of recently viewed scenarios with timestamp, click to replay
- Uses `localStorage` via existing `storageService` for history persistence

#### [NEW] [TopologyDiagram.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/components/features/topology/TopologyDiagram.tsx)

Orchestrator component:
- Renders all nodes and edges on a relative-positioned container
- Manages flow animation state via `useReducer`: `idle | playing | paused`
- Steps through `FlowStep[]` on interval (~1.5s per step), activating edges sequentially
- Passes active step data down to edges, nodes, and trace log
- SVG layer for edges sits underneath the node layer

#### [NEW] [index.ts](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/components/features/topology/index.ts)

Barrel export for `TopologyDiagram`.

---

### Page Layer

#### [MODIFY] [page.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/app/(dashboard)/page.tsx)

Complete rewrite of the dashboard page:
- Remove existing static cards (OIDC Verification, Service Mesh Health, Observability)
- Replace with minimalist layout: thin header zone with `HeroTitle` ("System Topography") + subtitle, then full-width `TopologyDiagram` taking ~85vh
- Keep auth logic (token detection, login/logout) in the header area but make it compact

#### [MODIFY] [layout.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/app/(dashboard)/layout.tsx)

Remove the `<Sidebar />` import and usage for the dashboard route group. The landing page should be **sidebar-free** (minimalist). The sidebar remains exclusively for `(manage)` routes.

---

### Styling

#### [MODIFY] [globals.css](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/app/globals.css)

Add utility classes:
- `.topology-node-glow` — box-shadow animation for active nodes
- `@keyframes flow-dash` — animated stroke-dashoffset for edge flow effect
- `@keyframes pulse-glow` — subtle pulse for active node borders

---

## Design Decisions

> [!IMPORTANT]
> **No React Flow dependency** — Using custom SVG + Framer Motion to avoid adding a heavy dependency. The existing `framer-motion@12.38.0` is already installed and sufficient for all animations.

> [!IMPORTANT]
> **Light theme** — Following the user's explicit request for light-leaning color scheme. The existing Material Design 3 tokens (`surface: #f7f9fb`, `primary: #004ac6`, `tertiary: #006242`) provide excellent contrast on light backgrounds.

> [!IMPORTANT]
> **Percent-based positioning** — Nodes use `%` coordinates relative to the diagram container. This makes the diagram inherently responsive without complex layout calculations.

---

## Verification Plan

### Automated Tests
- `npx tsc --noEmit` — TypeScript compilation check
- `npm run lint` — ESLint pass

### Manual Verification
- `npm run dev` → Navigate to `/` → Verify diagram renders with all 10 nodes
- Click "OIDC Login Flow" scenario → Verify edges animate step-by-step with numbered badges
- Verify trace log updates in sync with animation
- Verify history is saved and replay works
- Verify `/manage/*` routes still have sidebar
