---
artifact_type: ticket
id: RR-URG-09A
status: open
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
- [ ] Update `src/app/page.tsx` to include the `ArchitectureTopology` component.
- [ ] Implement `src/app/coffee/page.tsx` as the Public QR Trace List.
- [ ] Verify no-auth API calls for these pages.

## Acceptance Criteria
- Public users can see the topology on the home page.
- Public users can browse sold cups at `/coffee`.
- No login is required for these views.
