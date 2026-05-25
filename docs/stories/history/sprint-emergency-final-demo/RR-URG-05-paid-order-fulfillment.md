# RR-URG-05: Paid Order Fulfillment To Warehouse Reservation And Delivery

## Priority

P1. Depends on RR-URG-01, RR-URG-02, and inventory/warehouse readiness.

## Problem

Final demo uses paid retail order, not replenishment. The flow must show payment success or simulated success before warehouse reservation, dispatch, driver delivery, and mandatory return.

## Scope

- Align retail flow to paid order.
- Keep payment provider simulation allowed.
- Ensure warehouse reservation uses inventory lock.
- Trigger outbound dispatch after reservation.
- Complete order according to delivery/return policy.

## Implementation Details

### 1. Retail Paid Order

Retail order fields:

- `id`
- `store_id`
- `items`
- `total_amount`
- `payment_status`
- `fulfillment_status`
- `shipment_id`
- `trace_code` or public trace relation
- `correlation_id`

Statuses:

- `CREATED`
- `PAYMENT_PENDING`
- `PAYMENT_COMPLETED`
- `RESERVED`
- `DISPATCH_REQUESTED`
- `IN_TRANSIT`
- `DELIVERED`
- `RETURN_RECORDED`
- `COMPLETED`
- `FAILED`

Checklist:

- [ ] UI label uses paid order/checkout wording.
- [ ] `STORE_MGR` can create paid order only for assigned stores.
- [ ] Missing `store_ids` fails closed.
- [ ] Order event carries `store_id`, `order_id`, `correlation_id`, `trace_id`.

### 2. Payment

Payment can be:

- real provider webhook, or
- simulated payment success in production-demo.

Checklist:

- [ ] Payment success event is idempotent.
- [ ] Simulated success uses same downstream contract as real success.
- [ ] Payment webhook HMAC protection remains intact.
- [ ] Failed payment does not reserve inventory.

### 3. Warehouse Reservation

Reservation must:

- consume paid/successful order event.
- lock inventory by `warehouse_id` + `sku` or inventory lot.
- reserve available stock idempotently.
- publish `warehouse.stock.reserved` or `warehouse.stock.reservation_failed`.

Checklist:

- [ ] Valkey/Redlock or equivalent lock used.
- [ ] Concurrent orders cannot oversell.
- [ ] Reservation idempotent by event/order ID.
- [ ] Failure publishes compensating event.

### 4. Outbound Dispatch

After `warehouse.stock.reserved`:

- create outbound dispatch request.
- notify Warehouse dashboard.
- Warehouse Manager dispatches vehicle/driver.
- Logistics creates `RETAIL_DELIVERY`.

Checklist:

- [ ] Dispatch request tied to order/store/warehouse.
- [ ] Warehouse Manager scoped by warehouse.
- [ ] Store Manager can watch incoming delivery only for assigned store.

### 5. Completion Policy

Decision:

- Driver delivery confirmation marks delivery reached.
- Driver return-to-base is mandatory.
- Retail order `COMPLETED` waits until `logistics.driver.returned_to_base`.
- UI may show `DELIVERED` after store delivery, but this is not terminal.

Checklist:

- [ ] Delivery event updates order to `DELIVERED`.
- [ ] Return event updates order to `COMPLETED`.
- [ ] Trace timeline includes both delivery and return.

### 6. Order SAGA State Machine

This is the canonical paid order SAGA for the emergency sprint.

| Step | Event/API | Producer | Consumer/Owner | Order State | Notes |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | `POST /v1/retail/orders` | Retail UI/API | Retail Service | `CREATED` | `STORE_MGR` must be scoped to `store_id`. |
| 2 | `retail.order.created` | Retail Service | Payment Service | `PAYMENT_PENDING` | Carries `order_id`, `store_id`, `correlation_id`, `trace_id`. |
| 3 | `payment.completed` or `payment.simulated_completed` | Payment Service | Warehouse Service | `PAYMENT_COMPLETED` | Payment webhook HMAC remains protected if real provider used. |
| 4 | `warehouse.stock.reserved` | Warehouse Service | Retail + Warehouse UI + Logistics trigger | `RESERVED` | Reservation uses inventory lock and idempotency. |
| 5 | `warehouse.dispatch.requested` | Warehouse Service | Warehouse UI | `DISPATCH_REQUESTED` | Warehouse Manager assigns vehicle/driver. |
| 6 | `logistics.delivery.assigned` | Logistics Service | Retail/Warehouse/Driver UI | `IN_TRANSIT` after driver departs | Driver Client owns route simulation. |
| 7 | `logistics.delivery.driver_confirmed` | Logistics Service | Retail Service + Trace Service | `DELIVERED` | Not terminal because return is mandatory. |
| 8 | `logistics.driver.returned_to_base` | Logistics Service | Retail Service + Warehouse/Logistics UI | `COMPLETED` | Driver/vehicle become available. |

Failure/compensation events:

- `payment.failed`: order `FAILED`, no reservation.
- `warehouse.stock.reservation_failed`: order `FAILED` or `BACKORDERED` depending on final UI copy.
- `logistics.delivery.failed`: order `DELIVERY_FAILED`, warehouse/ops notification required.

### 7. Warehouse And Logistics Responsibility Split

| Responsibility | Owner |
| :--- | :--- |
| Paid order creation | Retail Service |
| Payment success/simulation | Payment Service |
| Inventory reservation and lock | Warehouse Service |
| Outbound dispatch request | Warehouse Service |
| Vehicle/driver assignment | Warehouse Manager action through Warehouse/Logistics API |
| GPS route simulation | Driver Client + Logistics Service validation |
| Delivery confirmation | Driver through Logistics Service |
| Return-to-base confirmation | Driver through Logistics Service |
| Final order completion | Retail Service after return event |
| Public trace projection | Trace Service |

## Acceptance Criteria

1. Paid order is the only final-demo retail demand path.
2. Payment/simulated payment success occurs before reservation.
3. Reservation uses concurrency guard.
4. Warehouse dispatch request is created after successful reservation.
5. Driver delivers and returns to base.
6. Retail order reaches terminal state according to documented completion policy.

## Test Checklist

- [ ] Unit: `STORE_MGR` cannot create order for unassigned store.
- [ ] Unit: payment duplicate event idempotent.
- [ ] Unit: reservation lock prevents oversell.
- [ ] Unit: reservation failure emits failure event.
- [ ] Integration: paid order -> payment success -> reserved -> dispatched.
- [ ] Integration: delivery -> return -> order completed.
