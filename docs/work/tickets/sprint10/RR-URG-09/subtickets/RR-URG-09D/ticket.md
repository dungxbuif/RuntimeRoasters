---
artifact_type: ticket
id: RR-URG-09D
status: open
owner: ai
priority: medium
trace:
  backlog_item: BL-001
  parent_ticket: RR-URG-09
---

# RR-URG-09D: Retail Dashboard & Real-time Integration

## Problem
Retail dashboard needs better store-scoping and real-time saga feedback.

## Scope
- Retail Dashboard: JWT store-scoping, "Paid Order" terminology.
- Real-time: Ensure socket updates trigger UI refreshes across dashboards.

## Tasks
- [ ] Implement store-scoping for `STORE_MGR`.
- [ ] Update terminology and saga feedback in Retail UI.
- [ ] Verify WebSocket connectivity and fanout.
