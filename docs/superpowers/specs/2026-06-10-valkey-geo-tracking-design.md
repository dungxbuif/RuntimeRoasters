# Design Spec: Valkey Geo Tracking (BL-002)

**Topic:** Real-time Event-Driven GPS Ingestion & Unified Mapping
**Status:** Approved
**Date:** 2026-06-10

## 1. Goal
Implement a high-performance, WebSocket-based GPS tracking system that replaces HTTP polling. This system must serve three distinct contexts: Public (Root Domain), Professional (Dashboards), and Admin (Global Control).

## 2. Architecture
- **Ingress:** Driver App streams GPS updates via WebSockets to `socket-service`.
- **Relay:** `socket-service` publishes updates to Kafka (`events.TopicLogisticsGPSUpdated`).
- **Persistence:** `logistics-service` consumes Kafka to update PostgreSQL (history) and Valkey (current state via `GEOADD`).
- **Egress:** `socket-service` broadcasts updates back to clients via WebSockets.

## 3. Component Visibility Tiers
We will use a unified `LogisticsMap` component with a `mode` prop.

| Tier | Context | Data Shown | Auth |
|---|---|---|---|
| **Public** | Root Landing Page | Truck Name, License Plate, Real-time Position | None |
| **Professional** | Manager/Driver Dashboards | + Shipment ID, Driver Name, Route Geometry, ETA | Required |
| **Admin** | Admin System Health | + System Metrics (Latency, Throughput, Fleet Performance) | Required (Admin) |

## 4. Technical Changes

### Backend
1. **Socket Service:**
   - Implement `gps_update` message handler in `readPump`.
   - Implement logic to broadcast `topology.event` for GPS updates.
   - Filter payloads based on `Visibility: "public"` for non-authenticated connections.

2. **Logistics Service:**
   - Add Kafka consumer for `logistics.gps.updated`.
   - Refactor `UpdateDriverLocation` usecase to use Valkey `GEOADD`.
   - Add `ListDriverLocations` endpoint for initial state hydration.

### Frontend
1. **LogisticsMap Component:**
   - Update to handle the three visibility modes.
   - Integrate WebSocket listener for live updates.
2. **Root Page:**
   - Add "Live Delivery Network" section using `LogisticsMap` in `public` mode.
3. **Admin Dashboard:**
   - Add global logistics control view with deep metrics.

## 5. Security & Privacy
- Public stream MUST NOT contain `shipment_id`, `driver_id` (internal UUIDs), or sensitive route names.
- Only `DRIVER` role can push `gps_update` messages.
