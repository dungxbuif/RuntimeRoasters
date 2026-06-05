---
artifact_type: user_stories
id: USER_STORIES
status: active
owner: shared
human_fields:
  - role
  - need
  - benefit
  - acceptance_criteria
ai_fields:
  - story_rows
  - trace_links
  - status_updates
shared_fields:
  - stories
updated: 2026-06-05
---

# User Stories

| ID | Role | Need | Benefit | Acceptance Criteria | Status |
| --- | --- | --- | --- | --- | --- |
| US-AUTH-001 | Admin | Initialize deterministic demo users/resources. | Demo can start from a clean environment. | Admin can bootstrap seed state; seeded accounts and assignments are deterministic. | active |
| US-FARM-001 | Farm Manager | Declare harvest for assigned farm. | Supply-chain story starts from real farm operation. | Farm Manager can create harvest; non-scoped users cannot. | active |
| US-WH-001 | Warehouse Manager | Dispatch pickup/delivery with driver and vehicle. | Physical movement is visible and role-owned. | Warehouse Manager can assign fleet; ADMIN does not perform daily dispatch. | active |
| US-DRIVER-001 | Driver | Use Driver Client to confirm route milestones. | Logistics state advances through real backend APIs. | Assigned driver can start simulation and post milestones/GPS. | active |
| US-STORE-001 | Store Manager | Create paid store order. | Retail demand triggers warehouse/payment/logistics SAGA. | Store Manager can create scoped order and see incoming delivery. | active |
| US-PUBLIC-QR-001 | Public visitor | Buy a demo cup without login and scan its QR. | Visitor can inspect Farm-to-Cup provenance. | UI creates `product_id`, backend persists sale, QR opens public trace. | ready |
| US-STORE-SALES-001 | Store Manager | View cups/items sold by selected store. | Store can audit sold cup IDs and QR trace links. | Sold-items UI lists invoice, product, sold time, and `product_id`. | ready |

## UI Coverage Notes

Legacy UI gap analysis has been reduced to active product-facing stories:

- Admin needs warehouse, store, vehicle, and driver resource management only as management functions; Admin should not perform daily harvest/order/dispatch work.
- Farm Manager needs harvest declaration plus downstream pickup visibility.
- Warehouse Manager needs inbound pickup queue, dispatch controls, receipt/intake, outbound delivery queue, and active shipment status.
- Store Manager needs store-scoped ordering, incoming delivery status, sold-cup history, and QR trace links.
- Driver needs assigned shipment list, route replay, GPS posting, and milestone confirmation.
- Public visitor needs sanitized no-auth trace pages for sold cups/items.

Stale mock-only UI claims from pre-Harness docs are not accepted behavior until backed by current tickets and validation evidence.
