# 🔍 Runtime Roasters — Missing UI Gap Analysis

> **Scope**: FE only. This document maps every role's spec-required screens against the current `client-app` implementation, identifies what's missing, and details simulation/demo buttons needed.

---

## 0. Infra Gaps (Affect All Roles)

Before building role-specific screens, these TypeScript/infrastructure gaps need fixing:

| Gap | Current | Required |
|:----|:--------|:---------|
| `UserRole` type | `'ADMIN' \| 'FARM_ADMIN' \| 'FARM_MANAGER' \| 'PROCESSOR' \| 'DRIVER' \| 'GUEST'` | Add `'WAREHOUSE_MGR'` and `'STORE_MGR'` |
| Sidebar nav | Only shows ADMIN sections + generic ops links. No role-aware nav groups for WAREHOUSE_MGR / STORE_MGR / DRIVER / PROCESSOR | Need role-conditional nav sections per persona |
| WebSocket/SSE client | No socket-service client exists | Need `useRealtimeEvents()` hook connected to `ws://localhost:8091/v1/realtime/...` |
| Notification system | No in-app notification component | Need toast/notification feed for async events (pickup requested, stock reserved, delivery arrived, etc.) |

---

## 1. ADMIN Role — ✅ Mostly Complete

### Existing Screens
| Route | Status | Notes |
|:------|:-------|:------|
| `/dashboard` (Intelligence Hub) | ✅ Exists | KPI cards, system health, anomaly watch |
| `/dashboard/users` | ✅ Exists | Create/manage user accounts |
| `/dashboard/farm-ops/registry` | ✅ Exists | Create/manage farms |
| `/dashboard/profile` | ✅ Exists | Account settings |
| `/dashboard/explorer` | ✅ Exists | System explorer (ADMIN only) |

### Missing for ADMIN
| Screen | Purpose | Priority |
|:-------|:--------|:---------|
| **Warehouse management** (CRUD) | Create warehouses, assign WAREHOUSE_MGR | 🔴 High |
| **Store management** (CRUD) | Create stores, assign STORE_MGR | 🔴 High |
| **Vehicle/Driver fleet** (CRUD) | Create vehicles, create driver records | 🔴 High |
| **ADMIN bootstrap** extended | Current SystemBootstrapModal exists but doesn't cover warehouse/store/vehicle/driver seed | 🟡 Medium |

---

## 2. FARM_MANAGER Role — 🟡 Partially Complete

### Existing Screens
| Route | Status | Notes |
|:------|:-------|:------|
| `/dashboard/farm-ops/registry` | ✅ Exists | Farm list (scoped by assignment) |
| `/dashboard/farm-ops/harvests` | ✅ Exists | Harvest list + create harvest |

### Missing for FARM_MANAGER
| Screen | Purpose | Simulation/Demo Button | Priority |
|:-------|:--------|:----------------------|:---------|
| **Harvest → Pickup status tracker** | After creating harvest, show pickup status: `PICKUP_REQUESTED → ASSIGNED → PICKED_UP → ARRIVED_WAREHOUSE → INTAKE_CREATED` | No button needed — watch-only via realtime | 🔴 High |
| **Active pickup map** | If farm has an active pickup, show driver progress on mini-map | Watch-only, fed by socket events | 🟡 Medium |
| **Realtime notifications** | "Driver assigned to your harvest", "Driver arrived at farm", "Pickup loading confirmed" | N/A — notification feed component | 🟡 Medium |

### Current Harvest Page Gap
The harvest page creates harvests but does **not** show the downstream pickup lifecycle. After creation, FARM_MANAGER has no visibility of what happens next.

---

## 3. WAREHOUSE_MGR Role — 🔴 Almost Entirely Missing

This is the **largest gap**. The current `/dashboard/warehouse` page shows batch intake/processing/inventory, but is missing the entire **inbound pickup dispatch** and **outbound retail dispatch** workflows, and has no role-scoping.

### Existing Screens
| Route | Status | Notes |
|:------|:-------|:------|
| `/dashboard/warehouse` | 🟡 Partial | Shows batch intake verification + processing queue + finished stock. Missing dispatch UIs |

### Missing for WAREHOUSE_MGR
| Screen/Section | Purpose | Simulation/Demo Button | Priority |
|:---------------|:--------|:----------------------|:---------|
| **Inbound Pickup Request Queue** | List of `warehouse.pickup.requested` events from farm harvests. WAREHOUSE_MGR sees pending pickup requests for assigned warehouse | N/A — list view | 🔴 Critical |
| **Dispatch Pickup button** | Assign vehicle + driver to a pickup request | **🎮 "Dispatch Pickup" button** → calls `POST logistics dispatch`, emits `warehouse.pickup.requested` | 🔴 Critical |
| **Active Inbound Shipments** | Track active farm→warehouse pickups with driver status | Watch via realtime map widget | 🔴 High |
| **Receipt/Intake creation** | After driver returns, create intake from returned pickup (currently partially exists) | Button to confirm receipt after `logistics.pickup.arrived_at_warehouse` | 🔴 High |
| **Outbound Dispatch Request Queue** | After stock reservation for retail orders, show pending dispatch requests | N/A — list view | 🔴 Critical |
| **Dispatch Retail Delivery button** | Assign vehicle + driver for store delivery after stock reservation | **🎮 "Dispatch Delivery" button** → calls logistics dispatch | 🔴 Critical |
| **Active Outbound Shipments** | Track active warehouse→store deliveries | Watch via realtime map widget | 🔴 High |
| **Warehouse selector** | If assigned to multiple warehouses, pick which one to manage | Dropdown/selector at top | 🟡 Medium |
| **Realtime notifications** | "New pickup request from K'Ho Farm", "Driver returned with harvest", "Stock reserved for Order #X" | Notification feed component | 🟡 Medium |

---

## 4. STORE_MGR Role — 🔴 Completely Missing (Dedicated UI)

The current retail pages (`/dashboard/retail` and `/dashboard/retail/orders`) are technical showcase UIs showing Saga orchestrator visualization and order creation, but they're **not scoped to STORE_MGR's assigned stores** and have no delivery tracking.

### Existing Screens
| Route | Status | Notes |
|:------|:-------|:------|
| `/dashboard/retail` | 🟡 Showcase only | Shows Saga Orchestrator diagram with hardcoded data. Not operational |
| `/dashboard/retail/orders` | 🟡 Partial | Create order form exists but needs store-scoping |

### Missing for STORE_MGR
| Screen/Section | Purpose | Simulation/Demo Button | Priority |
|:---------------|:--------|:----------------------|:---------|
| **Store Selector** (scoped) | Only show assigned stores from JWT `store_ids` | N/A — dropdown | 🔴 Critical |
| **Create Paid Order** (functional) | Existing form needs to be scoped to assigned store with real inventory SKUs | N/A — uses existing form | 🔴 Critical |
| **Order Status Tracker** | After order creation, show status pipeline: `CREATED → PAYMENT_PENDING → RESERVED → DISPATCH_REQUESTED → IN_TRANSIT → DELIVERED → COMPLETED` | N/A — watch-only via realtime | 🔴 Critical |
| **Incoming Delivery Tracker** | Show active delivery with driver GPS position on mini-map for assigned store | Watch-only, fed by socket events | 🔴 High |
| **Order History** | List of all orders for assigned stores with status filtering | N/A — list/table view | 🟡 Medium |
| **Realtime notifications** | "Stock reserved for your order", "Driver dispatched", "Delivery arriving" | Notification feed component | 🟡 Medium |

---

## 5. DRIVER Role — 🔴 Completely Missing

This is the **most simulation-heavy** role. No Driver Client UI exists at all.

### Missing for DRIVER
| Screen/Section | Purpose | Simulation/Demo Button | Priority |
|:---------------|:--------|:----------------------|:---------|
| **Driver Client main screen** | Show assigned shipments (only those assigned to this driver). Scoped by JWT driver identity | N/A — list view | 🔴 Critical |
| **Assigned Shipment Detail** | Show shipment info: origin, destination, cargo details, route preview on map | N/A — detail view | 🔴 Critical |
| **🎮 "Start Route" button** | Begin route simulation: browser loads seeded route points, animates vehicle icon, posts GPS/status updates to backend every 2-3 seconds for 30-45 seconds | **This is the core demo simulation button**. Starts a timed animation loop that: 1) loads route coordinates, 2) iterates through them at demo pace, 3) calls `POST /v1/logistics/drivers/location` at each tick | 🔴 Critical |
| **Route animation map** | Full-screen map showing the driver's vehicle moving along the seeded road path. Vehicle icon snaps to route points | Leaflet/Mapbox map with animated marker | 🔴 Critical |
| **🎮 Milestone Buttons** | Confirm shipment milestones. Different buttons depending on shipment type: | | 🔴 Critical |
| | **FARM_PICKUP flow:** | | |
| | • "Departed Base" → marks departure | **🎮 Button** | |
| | • "Arrived at Farm" → marks farm arrival | **🎮 Button** | |
| | • "Pickup Confirmed / Loading" → marks cargo loaded | **🎮 Button** | |
| | • "Return Started" → begins return leg | **🎮 Button** (starts return route animation) | |
| | • "Arrived at Warehouse" → marks return complete | **🎮 Button** | |
| | **RETAIL_DELIVERY flow:** | | |
| | • "Departed Warehouse" | **🎮 Button** | |
| | • "Arrived at Store" | **🎮 Button** | |
| | • "Delivery Confirmed" | **🎮 Button** | |
| | • "Return Started" → begins return leg | **🎮 Button** (starts return route animation) | |
| | • "Returned to Base" | **🎮 Button** | |
| **Simulation status indicator** | Show current simulation state: paused, active, completed. Show elapsed time and progress bar | Auto-updated during simulation | 🟡 Medium |
| **Simulation speed control** | Allow adjusting sim speed (e.g., 30s / 45s / 60s for the route) | Slider or preset buttons | 🟢 Low |

### Driver Simulation UX Flow (30-45 seconds total)
```
WAREHOUSE_MGR clicks "Dispatch" on pickup/delivery request
  └→ DRIVER opens assigned shipment
      └→ DRIVER clicks "Start Route" 
          └→ Browser loads seeded route coords from API/static data
          └→ Map shows animated vehicle icon moving along road path
          └→ Every 2-3s: POST GPS coordinates to backend
          └→ Backend validates assignment, stores in Valkey GEO + Postgres
          └→ Backend emits Kafka events → Socket service → All watching dashboards update
          └→ At destination: DRIVER clicks milestone button
          └→ Return leg starts (mandatory)
          └→ At base: DRIVER clicks final "Returned" button
```

---

## 6. PROCESSOR Role — 🟡 Partially Covered

Processing is partially covered in the current warehouse page, but not scoped to a PROCESSOR role.

### Missing for PROCESSOR
| Screen/Section | Purpose | Priority |
|:---------------|:--------|:---------|
| **Processing Queue** (role-scoped) | Show only processing batches assigned to this processor's warehouse | 🟡 Medium |
| **Start/Update Processing** | Button to transition batch status: `RECEIVED → HULLING → ROASTING → STOCKED` | 🟡 Medium |
| **Weight Loss Calculator** | Already partially exists. May need role-scoping only | 🟢 Low |

---

## 7. Finance/Payment — 🟡 Partially Complete

### Existing
| Route | Status | Notes |
|:------|:-------|:------|
| `/dashboard/finance` | ✅ Exists | Shows payment list with **Pass / Fail** webhook buttons |

### Analysis
The Finance page already has the key **🎮 "Pass" and "Fail" buttons** for simulating Stripe webhook outcomes. This is one of the most complete simulation UIs.

### Missing
| Screen/Section | Purpose | Priority |
|:---------------|:--------|:---------|
| **Role scoping** | Who should access Finance? Currently unguarded. Likely ADMIN or a finance role | 🟡 Medium |
| **Realtime status update** | After Pass/Fail, show downstream effects (stock reserved, order status changed) | 🟡 Medium |

---

## 8. Public Screens — 🟡 Partially Complete

### Existing
| Route | Status | Notes |
|:------|:-------|:------|
| `/` (Root) | ✅ Exists | Architecture Topology + Flow Control panel |
| `/dashboard/traceability` | ✅ Exists | Basic trace view |

### Missing for Public
| Screen/Section | Purpose | Priority |
|:---------------|:--------|:---------|
| **Public QR Trace Showcase** | `/trace/public/:trace_code` — public product journey page. Show farm→warehouse→retail sanitized trace | 🔴 High |
| **"Generate Demo QR Codes" button** | Button on public trace showcase page that fetches seeded product list and renders QR codes | 🟡 Medium |
| **QR code scanner/reader** | Click-to-open trace from QR | 🟡 Medium |

---

## 9. Realtime Components (Cross-Role)

These components don't exist yet but are needed by multiple role screens:

| Component | Used By | Priority |
|:----------|:--------|:---------|
| `useRealtimeEvents()` hook | All dashboard screens for live updates | 🔴 Critical |
| `NotificationFeed` component | Sidebar or header notification bell | 🔴 High |
| `MiniMapWidget` | FARM_MANAGER (pickup progress), STORE_MGR (delivery progress), WAREHOUSE_MGR (inbound/outbound tracking) | 🔴 High |
| `StatusPipeline` component | Reusable status stepper showing progress through milestones (e.g., CREATED→ASSIGNED→IN_TRANSIT→DELIVERED) | 🔴 High |
| `ShipmentTracker` component | Card showing active shipment with driver name, route, ETA, current status | 🟡 Medium |

---

## 10. Summary Gap Matrix

| Role | Existing Pages | Missing Pages | Sim Buttons Needed | Realtime UI |
|:-----|:---------------|:--------------|:-------------------|:------------|
| **ADMIN** | 5 ✅ | 3 (warehouse/store/fleet CRUD) | None (management only) | Health monitoring |
| **FARM_MANAGER** | 2 ✅ | 1 (pickup status tracker) | None (watch-only) | Pickup status stream |
| **WAREHOUSE_MGR** | 1 🟡 | 6 (dispatch, receipt, queues) | 🎮 Dispatch Pickup, 🎮 Dispatch Delivery, 🎮 Confirm Receipt | Full realtime |
| **STORE_MGR** | 0 ❌ | 4 (store selector, order, tracker, history) | None (uses existing order form) | Delivery status stream |
| **DRIVER** | 0 ❌ | 5 (client, detail, map, milestones, sim) | 🎮 Start Route, 🎮 6-8 Milestone buttons, 🎮 Speed control | Full GPS stream |
| **PROCESSOR** | 0.5 🟡 | 1 (role-scoped processing) | 🎮 Start/Finalize Processing | Batch status |
| **Public** | 2 ✅ | 2 (QR trace page, QR generator) | 🎮 Generate QR Codes | Sanitized topology |

---

## 11. Recommended Build Priority

```
Phase 1: Infrastructure (enables all roles)
  ├─ Fix UserRole type (add WAREHOUSE_MGR, STORE_MGR)
  ├─ useRealtimeEvents() hook
  ├─ NotificationFeed component
  ├─ StatusPipeline component
  └─ Role-aware Sidebar navigation

Phase 2: DRIVER Client (highest demo impact)
  ├─ /dashboard/driver — Assigned shipment list
  ├─ /dashboard/driver/[shipmentId] — Detail + map + simulation
  ├─ Route animation engine (timed GPS posting)
  └─ Milestone confirmation buttons

Phase 3: WAREHOUSE_MGR Operational Screens
  ├─ Inbound pickup request queue
  ├─ Dispatch pickup (vehicle/driver assignment)
  ├─ Active inbound/outbound shipment tracking
  ├─ Receipt/intake from returned pickup
  └─ Outbound retail dispatch queue + button

Phase 4: STORE_MGR Dashboard
  ├─ Store-scoped order creation
  ├─ Order status pipeline tracker
  ├─ Incoming delivery tracker (mini-map)
  └─ Order history

Phase 5: ADMIN CRUD Extensions
  ├─ Warehouse CRUD + assign manager
  ├─ Store CRUD + assign store manager
  └─ Vehicle/Driver fleet CRUD

Phase 6: Public QR Trace
  ├─ /trace/public/:trace_code page
  ├─ QR code generator button
  └─ Sanitized product journey view

Phase 7: Polish
  ├─ PROCESSOR role-scoped processing queue
  ├─ Finance role-scoping
  └─ Cross-role realtime notification polish
```

> [!IMPORTANT]
> **Phase 2 (DRIVER Client)** is the most visually impactful for demo — it's the "wow factor" with map animation, vehicle movement, and milestone buttons. It also generates the realtime events that all other dashboards consume.

> [!NOTE]
> The current Logistics page (`/dashboard/logistics`) already has basic simulation logic (client-side GPS ticking) but it runs from the **Logistics dashboard** perspective, not from an authenticated **DRIVER** perspective. The spec requires driver-authenticated simulation posting to backend, not just client-side map animation.
---
# Missing UI Implementation Walkthrough

I have completed the implementation of all missing UI screens identified in our earlier Gap Analysis. These screens are now built with high-fidelity mock data and layout, ready to be wired up to the real backend APIs when needed.

## Changes Made

### 1. Infrastructure Foundation
- **Role Types**: Added `WAREHOUSE_MGR` and `STORE_MGR` to the `UserRole` type.
- **Routes**: Added new routes in `constants/routes.ts` (`/dashboard/driver`, `/dashboard/store`, `/dashboard/resources`, `/trace`).
- **Sidebar Navigation**: Completely rebuilt the `Sidebar.tsx` to support role-aware sections using `RoleGuard`. Now, each role sees exactly what they need (e.g., `DRIVER` sees "Driver Client", `STORE_MGR` sees "Retail Operations").
- **Shared Components**: Created `StatusPipeline` (a horizontal milestone stepper) and `NotificationFeed` (a right-rail activity feed), which are heavily used in all the new dashboards.

### 2. Driver Client (`/dashboard/driver`)
- Built the canonical demo path screen.
- Features a Route Simulation engine with speed controls (30s, 45s, 60s).
- GPS Telemetry panel that reads out current mock coordinates.
- Milestone confirmation pipeline (Assigned → Departed → At Farm → Loaded → Returning → Arrived).

### 3. Warehouse Manager (`/dashboard/warehouse`)
- Completely overhauled the existing page into a **3-Column Dispatch Layout**.
- **Inbound Queue (Left)**: Manages incoming coffee from farms, with buttons to "Dispatch Pickup" and "Confirm Receipt".
- **Processing & Inventory (Center)**: Shows active roast batches and live finished stock levels.
- **Outbound Queue (Right)**: Manages deliveries to retail stores, with a "Dispatch Delivery" button and live driver tracking widgets.

### 4. Store Manager (`/dashboard/store`)
- New dashboard tailored for the retail side.
- Store selector at the top (to filter orders for a specific store).
- Table of active orders and an embedded pipeline tracking the status of the current active delivery.
- "Incoming Delivery" mini-map widget to visualize ETA of the driver.

### 5. Farm Manager Update (`/dashboard/farm-ops/harvests`)
- Enhanced the existing Harvest ledger.
- Added a `StatusPipeline` under each harvest row so the Farm Manager can see downstream pickup logistics (e.g., when a truck arrives to take the harvest).
- Added a "Pickup Activity" notification feed to the page.

### 6. Admin Resource Management (`/dashboard/resources`)
- Built a unified CRUD dashboard for system resources.
- Tabbed interface to manage: **Warehouses**, **Stores**, **Vehicles**, and **Drivers**.
- Includes a "Quick Assign" modal mock to map drivers to vehicles or managers to stores.

### 7. Public QR Trace (`/trace/[code]`)
- Unauthenticated, mobile-friendly view designed for end-consumers.
- Renders a clean vertical timeline ("Farm to Cup Journey") that displays every step (Harvested → Picked Up → Processed → Stocked → Ordered → Delivered → Served).
- Features a mock QR code scanning target area.

## Next Steps
- You can now preview all of these screens by navigating to them in the Next.js app running locally.
- To test the role-based navigation, you may need to update your local mock user's role in the `auth-service` or the `useAuth` hook mock.
- The next major technical phase would be to connect these UIs to the real backend endpoints (e.g., building out the missing gRPC handlers for Logistics and updating the frontend services).
