# Order Approval And Fulfillment

## Current Implementation State

There is no complex human approval workflow in the current demo. "Approval" means the order passes event-driven fulfillment gates.

## Fulfillment Gates

| Gate | Meaning |
| --- | --- |
| Order created | `STORE_MGR` creates an order and retail emits `retail.order.created`. |
| Payment intent | Payment service creates a demo provider intent. |
| Payment success | Signed demo webhook marks the payment successful and emits payment completion. |
| Stock reserved | Warehouse reserves stock only after the successful payment gate. |
| Warehouse dispatch | Warehouse/logistics assigns delivery resources. |
| Driver delivery | Driver confirms delivery milestone. |
| Return-to-base completion | Driver completes return leg after pickup or delivery. |

## Simulated Pieces

- Payment provider settlement is simulated/demo.
- Finance Pass/Fail buttons create signed Stripe-compatible demo webhook payloads.
- Historical seeded dashboards may include direct read-model rows for immediate demo state.

## Real Runtime Pieces

- Retail order creation.
- Payment rows and webhook business validation.
- Warehouse stock reservation logic.
- Logistics assignment and shipment status transitions.
- Kafka events for downstream trace/audit/topology.
- Compensation/refund behavior where implemented for failure paths.

## Omitted Pieces

- Manual multi-step human order approval.
- Full accounting general ledger.
- Real external payment clearing.
- Store/warehouse/farm dual confirmation mode.
- Fraud review and settlement reconciliation.

## Role Responsibilities

| Role | Current demo responsibility |
| --- | --- |
| `STORE_MGR` | Creates retail orders and watches store-scoped fulfillment. |
| `WAREHOUSE_MGR` | Handles stock, dispatch, pickup, intake, and warehouse-scoped shipment work. |
| `DRIVER` | Executes assigned shipments, posts GPS, and confirms milestones. |
| `ADMIN` | Initializes demo state and reviews aggregate operational surfaces. |

## Recommended Demo Script

1. Log in as `STORE_MGR` and create an order.
2. Open `/dashboard/finance` in a separate profile and click Pass for the pending payment.
3. Confirm payment changes to succeeded.
4. Verify warehouse reservation/dispatch state.
5. Log in as `DRIVER`, run route simulation, and confirm delivery.
6. Search traceability by known order/shipment ID.

## Review Language

Use: "Order approval is modeled as fulfillment gates in an event-driven saga."

Avoid: "The current demo has a human approval workflow."
