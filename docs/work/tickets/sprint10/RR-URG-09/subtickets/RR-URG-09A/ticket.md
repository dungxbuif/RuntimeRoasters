---
artifact_type: ticket
id: RR-URG-09A
status: done
owner: ai
priority: high
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  parent_ticket: RR-URG-09
---

# RR-URG-09A: Public-Safe Topology & QR Trace List

## Problem
Public users need to see the system's live health (topology) and find cups to trace without logging in.

## Scope
- Integrate Public-Safe Architecture Topology into the landing page (root).
- Create a Public QR Trace List page (`/coffee`) to show seeded sold cups.
- Ensure no-auth access to these components.

## Tasks
- [x] Update `src/app/page.tsx` to include the `ArchitectureTopology` component.
- [x] Implement `src/app/trace/page.tsx` as the Public QR Trace Registry.
- [x] Verify no-auth API calls for these pages.

## Acceptance Criteria
- [x] Public users can see the topology on the home page.
- [x] Public users can browse sold cups at `/trace`.
- [x] No login is required for these views.
