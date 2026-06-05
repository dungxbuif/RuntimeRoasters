# RR-URG-10: Demo Runbook And End-To-End Verification

## Priority

P3 final ticket. Depends on all implementation tickets.

## Problem

The final demo needs a reliable operator guide. The guide must tell the user which account logs in, what button to click, which screen to watch, what should happen, and how to recover when a demo step stalls.

## Scope

- Create user-facing demo guide.
- Document accounts and routes.
- Document two-screen demo setup.
- Document expected realtime visuals.
- Document public QR trace demo.
- Add final E2E checklists.

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

Checklist:

- [ ] Guide lists exact account email/password.
- [ ] Guide lists exact route/page URLs.
- [ ] Guide lists exact buttons/actions.
- [ ] Guide lists expected events/status changes.
- [ ] Guide lists where to watch realtime.

### 2. Script A: Farm Pickup

Windows:

- Window 1: Farm Manager.
- Window 2: Warehouse Manager.
- Window 3: Driver Client.
- Optional observer: Logistics map or Trace page.

Steps:

1. Farm Manager creates harvest.
2. Warehouse sees pickup request.
3. Warehouse dispatches driver.
4. Driver starts route to farm.
5. Driver confirms pickup/loading.
6. Driver returns to warehouse/base.
7. Warehouse creates intake.
8. Processing can start.

Checklist:

- [ ] Pickup request visible.
- [ ] GPS marker moves.
- [ ] Pickup milestone visible.
- [ ] Return milestone visible.
- [ ] Intake created.
- [ ] Trace updated.

### 3. Script B: Paid Order Delivery

Windows:

- Window 1: Store Manager.
- Window 2: Warehouse Manager.
- Window 3: Driver Client.
- Optional observer: Logistics map or Trace page.

Steps:

1. Store Manager creates paid order.
2. Payment succeeds or simulated payment succeeds.
3. Warehouse reserves stock.
4. Warehouse dispatches retail delivery.
5. Driver starts route to store.
6. Driver confirms delivery.
7. Driver returns to warehouse/base.
8. Order completes according to final completion policy.

Checklist:

- [ ] Paid order created.
- [ ] Payment status success/simulated success.
- [ ] Stock reserved.
- [ ] Dispatch request visible.
- [ ] Delivery route visible.
- [ ] Return route visible.
- [ ] Order terminal state reached.

### 4. Script C: Public ArchitectureTopology

Steps:

1. Open root Client App without login.
2. Verify topology renders.
3. Trigger a demo flow.
4. Watch sanitized live service/event highlights if enabled.

Checklist:

- [ ] No auth required.
- [ ] No private payload exposed.
- [ ] Service highlights match active flow.

### 5. Script D: Public QR Trace

Steps:

1. Open public trace showcase page without login.
2. Click generate/show QR codes.
3. Click or scan one QR.
4. Verify trace page loads a real trace document.

Expected trace sections:

- farm/harvest.
- pickup and driver return.
- warehouse intake.
- processing/inventory.
- paid order.
- delivery and driver return.
- service/component detail if OTel exists.

Checklist:

- [ ] QR list visible.
- [ ] Public trace page loads without login.
- [ ] Trace is real sold-cup data.
- [ ] No private user/payment secrets.

### 6. Final Technical Verification

Required before demo:

- [ ] `go test ./...` or documented targeted test set.
- [ ] Frontend build passes.
- [ ] E2E suite passes or known gaps documented.
- [ ] TraceId sample recorded.
- [ ] Seed reset verified.
- [ ] Demo runbook followed from clean state.

## Acceptance Criteria

1. A new reviewer can follow the runbook without engineering help.
2. Every step has expected UI result.
3. Every role uses the correct account.
4. QR trace demo proves real traceability.
5. TraceId continuity evidence is included.
6. Known limitations and recovery steps are documented.
