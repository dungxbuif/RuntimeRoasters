---
artifact_type: ticket
id: BL-003
status: verified
owner: ai
sprint: 12
backlog_rank: 3
human_fields:
  - priority
  - acceptance
ai_fields:
  - proposed_approach
  - impact
  - evidence
---

# BL-003: Saga Monitor Real-time WebSocket Updates

## Problem Statement

The Saga Monitor (`/dashboard/retail`) currently refetches all orders every 5 seconds via HTTP polling. This creates unnecessary load and causes up to 5-second lag for saga state transitions that are visible in the UI. The existing WebSocket infrastructure (`/v1/realtime/stream`) already delivers real-time saga-related events (order, payment, warehouse events) to authenticated users, but the Saga Monitor page does not consume it.

## Scope

**In scope:**
- Connect the Saga Monitor page to the existing authenticated WebSocket stream via `useSocket`
- Invalidate the `['retail', 'orders']` React Query cache on receipt of any order/saga-related event
- Add a live connection status indicator (dot) to the page header
- Reduce polling fallback from 5s to 30s (resilience for WebSocket outages)

**Out of scope:**
- New backend endpoints or Kafka consumers
- SSE protocol (WebSocket already in place)
- Changes to other pages

## Acceptance Criteria

- [x] When an order status transitions (pay/confirm), the Saga Monitor updates within 1–2 seconds without a manual refresh
- [x] A live/offline connection badge is visible in the Saga Monitor header
- [x] When WebSocket is unavailable, polling at 30s (5s if disconnected) continues as a fallback
- [x] No regression in existing unit or lint checks

## Small Task Exemption

Small task exemption: yes
Reason: pure frontend hook integration, no API/DB/auth/runtime/deployment changes
Impact checked: API=no, DB=no, Security=no, Runtime=no, Standards=no

## Impacted Areas

- `src/apps/client-app/src/app/(dashboard)/dashboard/retail/page.tsx` — primary change
- No backend, no KrakenD, no database changes

## Test Plan

1. Run `npm test` — must stay 11/11 green
2. Run `npm run lint` — must pass
3. Manual: trigger pay/confirm on Saga Monitor, observe update within 2s

## Evidence

```
npm test — 11/11 passed
npm run lint — 0 errors in retail/page.tsx (52 pre-existing errors elsewhere, unchanged)
```

**Changed file:** `src/apps/client-app/src/app/(dashboard)/dashboard/retail/page.tsx`
- Added `useSocket` hook call (authenticated WS stream)
- Added `useEffect` that invalidates `['retail','orders']` query on saga-topic events
- Reduced polling to 30s fallback (5s when disconnected)
- Added Live/Polling connection badge in header

## Trace Links

- Backlog: `docs/work/BACKLOG.md` BL-003
- Sprint: `docs/work/tickets/sprint12/main.md`
