# Valkey Geo Tracking Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement real-time GPS tracking using WebSockets, Kafka, and Valkey GEO.

**Architecture:** Edge ingestion via WebSockets in `socket-service`, event processing in `logistics-service` with Valkey `GEOADD`, and real-time broadcast back through `socket-service`.

**Tech Stack:** Go (Gin, Kafka-go, Go-redis), Next.js (TypeScript, WebSocket).

---

### Task 1: Socket Service Ingress (GPS WebSocket)

**Files:**
- Modify: `src/apps/socket-service/internal/usecase/service.go`
- Test: `src/apps/socket-service/internal/usecase/service_test.go`

- [ ] **Step 1: Update readPump to handle gps_update**
Implement the logic to parse incoming JSON messages. If type is `gps_update`, verify `DRIVER` role and publish to Kafka `events.TopicLogisticsGPSUpdated`.

- [ ] **Step 2: Add integration test for WS GPS ingestion**
Verify that sending a WS message results in a Kafka publish call.

- [ ] **Step 3: Commit**
```bash
git add src/apps/socket-service/internal/usecase/service.go
git commit -m "feat(socket): add WebSocket GPS ingestion to Kafka"
```

### Task 2: Logistics Service Consumer & Valkey GEO

**Files:**
- Modify: `src/apps/logistics-service/internal/usecase/service.go`
- Modify: `src/apps/logistics-service/internal/app/init.go`
- Test: `src/apps/logistics-service/internal/usecase/service_test.go`

- [ ] **Step 1: Implement HandleGPSUpdatedEvent**
Consume the Kafka event, perform `GEOADD` on Valkey, and update PostgreSQL `vehicles` coordinates.

- [ ] **Step 2: Add ListDriverLocations API**
Add a method to fetch all current driver locations from Valkey for initial map hydration.

- [ ] **Step 3: Register Kafka Consumer in init.go**
Ensure the worker is listening to `events.TopicLogisticsGPSUpdated`.

- [ ] **Step 4: Commit**
```bash
git add src/apps/logistics-service/internal/usecase/service.go src/apps/logistics-service/internal/app/init.go
git commit -m "feat(logistics): consume GPS events and update Valkey GEO"
```

### Task 3: Unified Map Component (Frontend)

**Files:**
- Modify: `src/apps/client-app/src/components/features/logistics/LogisticsMap.tsx`
- Modify: `src/apps/client-app/src/app/(dashboard)/dashboard/logistics/page.tsx`

- [ ] **Step 1: Update LogisticsMap for multiple modes**
Accept `mode` prop (`public`, `pro`, `admin`) and handle visibility logic for labels/IDs.

- [ ] **Step 2: Integrate WebSocket in Logistics Dashboard**
Connect to `socket-service` and listen for `topology.event` with topic `logistics.gps.updated`.

- [ ] **Step 3: Commit**
```bash
git add src/apps/client-app/src/components/features/logistics/LogisticsMap.tsx
git commit -m "feat(ui): unified map component with real-time WebSocket support"
```

### Task 4: Public & Admin Views

**Files:**
- Modify: `src/apps/client-app/src/app/page.tsx` (or landing page)
- Create: `src/apps/client-app/src/app/(dashboard)/dashboard/logistics/admin/page.tsx`

- [ ] **Step 1: Add Live Network to Root Page**
Embed `LogisticsMap` in `public` mode on the root landing page.

- [ ] **Step 2: Add Admin Metrics Dashboard**
Create a new admin view with deep fleet metrics.

- [ ] **Step 3: Commit**
```bash
git add src/apps/client-app/src/app/page.tsx
git commit -m "feat(ui): add public live network and admin metrics"
```
