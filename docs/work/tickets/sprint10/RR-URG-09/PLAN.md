# Implementation Plan - RR-URG-09: Client Dashboards & Public Trace

## Overview
This plan covers the frontend refinements and new public pages required for the final demo. The work is split into 5 subtickets for focused execution.

## Subtickets

### RR-URG-09A: Public-Safe Topology & QR Trace List (Phase 1)
- **Goal**: Provide entry points for public users.
- **Tasks**:
  - Update root page (`/`) to integrate Architecture Topology.
  - Create `/trace` page to list recently sold cups with QR codes.
  - Add navigation between root, coffee kiosk, and trace list.

### RR-URG-09B: Farm & Warehouse Dashboard Refinements (Phase 2)
- **Goal**: Operational polish for managers.
- **Tasks**:
  - Hide harvest creation if no farms assigned.
  - Add visual feedback for pickup/dispatch lifecycle.
  - Scoped warehouse actions.

### RR-URG-09C: Driver Client & Logistics Simulation Improvements
- **Goal**: Full simulation lifecycle.
- **Tasks**:
  - Implement Return Leg in `DriverClient`.
  - Add Shipment Detail Drawer in `LogisticsDashboard`.
  - Display GPS telemetry and status history.

### RR-URG-09D: Retail Dashboard & Real-time Integration
- **Goal**: Final supply chain link.
- **Tasks**:
  - Store-scoped order list for `STORE_MGR`.
  - WebSocket integration for live status updates without refresh.
  - Consistent "Paid Order" terminology.

### RR-URG-09E: Final Verification & E2E Testing
- **Goal**: Quality assurance.
- **Tasks**:
  - Comprehensive Playwright E2E script.
  - Final documentation and demo evidence.

## Technical Approach
- **Auth**: Use `CasbinGuard` and `useAuth` hook for scoping.
- **Data**: Fetch from public endpoints for guest views.
- **Styling**: Tailwind CSS with consistent brand colors.
- **Real-time**: Socket service subscription based on resource IDs.

## Acceptance Criteria
- Full journey (Farm -> Warehouse -> Logistics -> Retail -> Trace) is actionable and visible.
- Public views are safe and informative.
- RBAC is correctly enforced at the UI level.
