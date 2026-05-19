# Business Logic: Retail Order Fulfillment

This document defines the downstream paid retail order flow from store demand to warehouse reservation and driver delivery.

## 1. Purpose

Paid retail orders trigger warehouse reservation and delivery. The production-demo flow may use simulated payment providers, but it should still model the business sequence as paid order -> reservation -> dispatch -> delivery -> driver return.

## 2. Canonical Flow

1. `STORE_MGR` logs in.
2. User opens assigned Retail dashboard.
3. User creates a paid retail order for an assigned store.
4. Retail service publishes order event.
5. Payment succeeds through real webhook or simulated payment success.
6. Warehouse reserves finished inventory with concurrency guard.
7. Warehouse dashboard shows outbound dispatch request.
8. `WAREHOUSE_MGR` assigns vehicle/driver.
9. `DRIVER` opens assigned shipment in Driver Client.
10. Driver starts route simulation; browser posts GPS/status updates to backend.
11. Driver confirms store arrival and delivery completion.
12. Retail order reaches delivery-completed status after driver delivery confirmation.
13. Driver return-to-base is mandatory and must be recorded.

## 3. Roles

| Role | Retail Visibility | Allowed UI Actions |
| :--- | :--- | :--- |
| `ADMIN` | Aggregate store/order counts and health | Create stores, assign store managers, view overview. |
| `STORE_MGR` | Assigned stores through `store_ids` | Create paid retail order, view incoming delivery status. |
| `WAREHOUSE_MGR` | Outbound dispatch queue for assigned warehouse | Assign vehicle/driver after stock reservation. |
| `DRIVER` | Assigned delivery shipments only | Start route simulation, update GPS, confirm arrival/delivery/return. |

Store-side receipt confirmation is optional strict mode. It is not required for the main production-demo flow.

## 4. UI Behavior

Retail dashboard should show:

- Assigned store selector.
- Create paid order action.
- Order/payment/reservation status.
- Incoming delivery notification.
- Delivery progress from Driver Client updates.
- Completed state after driver confirmation.

Retail dashboard should not directly update driver GPS or shipment route state.

## 5. State Rules

Retail order statuses:

- `CREATED`
- `PAYMENT_PENDING` or `PAYMENT_SIMULATED`
- `RESERVED`
- `DISPATCH_REQUESTED`
- `IN_TRANSIT`
- `DELIVERED`
- `COMPLETED`
- `FAILED`

Completion rule for demo:

- Main demo: driver delivery confirmation completes delivery.
- Driver return-to-base is still required after delivery.
- Strict future mode: store-side receipt confirmation may be required before completion.

## 6. Paid Order Decision

Decision: final demo uses paid order.

| Term | Meaning | When To Use |
| :--- | :--- | :--- |
| Paid order | A demand created from a customer/store order after payment success or simulated payment success. | Required for final demo. |
| Replenishment request | A store restock request independent from a direct customer payment. | Out of scope for final demo unless reopened later. |

Implementation guidance:

- The downstream flow is: paid order -> payment success or simulated success -> warehouse reservation -> dispatch -> driver delivery -> driver return.
- The UI should label the action as paid order or checkout/order placement, not generic replenishment.
- Replenishment can reuse the downstream pipeline in the future, but it is not part of this ticket.

## 7. Events

Retail:

- `retail.order.created`
- `retail.delivery.received`

Warehouse/logistics:

- `warehouse.stock.reserved`
- `warehouse.stock.reservation_failed`
- `warehouse.dispatch.requested`
- `logistics.delivery.assigned`
- `logistics.delivery.departed`
- `logistics.delivery.arrived_at_store`
- `logistics.delivery.driver_confirmed`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

## 8. Authorization

`STORE_MGR` can only see and create demand for assigned stores. Missing `store_ids` should fail closed: no store data and no create action.
