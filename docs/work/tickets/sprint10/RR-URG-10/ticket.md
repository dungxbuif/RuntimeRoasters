# RR-URG-10: Demo Runbook And End-To-End Verification

## Priority

P3 final ticket. Depends on all implementation tickets.

## Problem

The final demo needs a reliable operator guide. The guide must tell the user which account logs in, what button to click, which screen to watch, what should happen, and how to recover when a demo step stalls.

## Status: done

## Completion

Ticket completed. Final Demo Runbook is available at [FINAL_DEMO_RUNBOOK.md](../../../demo/FINAL_DEMO_RUNBOOK.md).

## Tasks
- [x] Create user-facing demo guide.
- [x] Document accounts and routes.
- [x] Document two-screen demo setup.
- [x] Document expected realtime visuals.
- [x] Document public QR trace demo.
- [x] Add final E2E checklists.

## Implementation Details

### 1. Demo Guide File

Create:

- `docs/demo/FINAL_DEMO_RUNBOOK.md` or equivalent.

Required sections:

- Environment prerequisites.
- Demo account table.
- Seed/reset command.
- Window layout.
- Script A: Farm -> Warehouse pickup -> Intake -> Processing.
- Script B: Paid order -> Reservation -> Delivery -> Return.
- Script C: Public ArchitectureTopology.
- Script D: Public QR trace.
- Troubleshooting and recovery.

### 2. Script A: Farm Pickup
(Steps implemented and verified in Runbook)

### 3. Script B: Paid Order Delivery
(Steps implemented and verified in Runbook)

### 4. Script C: Public ArchitectureTopology
(Steps implemented and verified in Runbook)

### 5. Script D: Public QR Trace
(Steps implemented and verified in Runbook)

### 6. Final Technical Verification
- [x] Backend tests pass.
- [x] Frontend build verified.
- [x] E2E System Journey verified.

## Acceptance Criteria

1. A new reviewer can follow the runbook without engineering help.
2. Every step has expected UI result.
3. Every role uses the correct account.
4. QR trace demo proves real traceability.
5. TraceId continuity evidence is included.
6. Known limitations and recovery steps are documented.
