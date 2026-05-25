# RR-URG-06: Realtime Notifications And Socket/SSE Broadcasts

## Priority

P2. Depends on stable event contracts and enough backend events to broadcast.

## Problem

The demo needs two-screen visibility. Dashboards must update while the actor performs actions, but socket/SSE must not become the source of truth.

## Scope

- Add persisted notifications or notification module.
- Add socket/SSE realtime service/module.
- Support public sanitized topology stream.
- Support authenticated role-scoped private streams.
- Add display delay/pacing for demo.

## Implementation Details

### 1. Notification Model

Fields:

- `id`
- `scope`
- `target_role`
- `target_user_id`
- `farm_id`
- `warehouse_id`
- `store_id`
- `shipment_id`
- `entity_type`
- `entity_id`
- `type`
- `title`
- `message`
- `status`: `UNREAD`, `READ`, `ACKED`, `RESOLVED`
- `created_at`
- `acknowledged_at`

Endpoints:

- `GET /v1/notifications`
- `POST /v1/notifications/:id/ack`
- `POST /v1/notifications/:id/resolve`

Checklist:

- [ ] Notifications persisted.
- [ ] Notifications scoped by role/entity.
- [ ] Pickup request notification.
- [ ] Outbound dispatch notification.
- [ ] Driver assignment notification.
- [ ] Delivery/return notification.

### 2. Socket/SSE Stream

Private stream:

- requires JWT.
- uses role/entity scope.
- sends only events the user can see.

Public stream:

- no auth allowed only for sanitized root ArchitectureTopology/public trace visualization.
- must not expose private IDs, payment secrets, user data, or operational action payloads.

Candidate endpoints:

- `GET /v1/realtime/stream`
- `GET /v1/realtime/public/topology`

Checklist:

- [ ] Private stream validates JWT.
- [ ] Private stream scopes by role/entity.
- [ ] Public stream sanitized.
- [ ] Disconnects handled.
- [ ] Backpressure simple but safe.

### 3. Events To Broadcast

- `notification.created`
- `warehouse.pickup.requested`
- `logistics.pickup.assigned`
- `logistics.gps.updated`
- `logistics.shipment.status_changed`
- `logistics.pickup.arrived_at_warehouse`
- `warehouse.intake.created`
- `warehouse.dispatch.requested`
- `logistics.delivery.assigned`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

### 4. Demo Pacing

Config:

- `REALTIME_DEMO_DELAY_MS=500`
- `LOGISTICS_SIM_TICK_SECONDS=3`

Checklist:

- [ ] Delay applies only to broadcast/display, not persisted state.
- [ ] Demo viewer can see each transition.
- [ ] Fast-forward action emits audited event.

## Acceptance Criteria

1. Warehouse dashboard updates when harvest creates pickup request.
2. Driver Client receives assigned shipment.
3. Logistics map receives GPS/status updates.
4. Retail dashboard receives incoming delivery updates.
5. Public ArchitectureTopology receives sanitized service/event updates.
6. Private streams require auth.
7. Socket/SSE failure does not corrupt backend state.

## Test Checklist

- [ ] Unit: notification scoping.
- [ ] Unit: socket event sanitizer.
- [ ] Integration: authenticated stream receives assigned event.
- [ ] Integration: unauthorized stream denied.
- [ ] Manual: two-browser demo updates in order.
