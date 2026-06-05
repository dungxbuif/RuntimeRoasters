---
artifact_type: validation_matrix
id: VALIDATION_MATRIX
status: active
owner: shared
---
# Validation Matrix

| Requirement | Phase | Ticket/Bug | Contract/Behavior | Unit | Integration | E2E | Platform | Status | Evidence |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| RBAC | M1 | ADMIN-01 | Admin restricted from domain writes | yes | no | yes | no | implemented | verified |
| Logistics Map | M2 | LOG-01 | Map route filtering and light mode | no | no | yes | no | implemented | verified |
| Order Saga | M2 | SAGA-01 | Inventory reservation flow | yes | yes | no | no | implemented | verified |
| Public Sold-Cup QR Trace | M2 | RR-URG-07A..07E | UI-issued sold-item `product_id` resolves to sanitized origin trace | planned | planned | planned | planned | planned | pending |
| Fresh Setup And Admin Bootstrap | M1 | Bootstrap | Reset, admin login, first-run seed, resource assignment visibility | partial | partial | pending | partial | partial | Legacy TC-1.* consolidated; needs current evidence refresh |
| Upstream Farm To Warehouse | M2 | Farm/Warehouse/Logistics | Harvest declaration, pickup dispatch, driver route, warehouse intake | pending | pending | pending | pending | planned | Legacy TC-2.* consolidated |
| Downstream Warehouse To Retail | M2 | Retail/Payment/Warehouse/Logistics | Store order, payment pass/fail, stock reservation, dispatch, return to base | partial | partial | pending | pending | partial | Legacy TC-3.* consolidated |
| Observability Trace Integrity | M2 | Trace/Observability | One trace ID across gateway, services, Kafka, trace read models | pending | pending | pending | pending | planned | Legacy TC-4.* consolidated |
| Realtime Awareness | M2 | RR-URG-06 | Public topology and private role-scoped notification streams | yes | yes | yes | yes | implemented | RR-URG-06 evidence |

## Manual UAT Flow Checklist

Use these scenario-level checks when a ticket changes the corresponding flow:

- Admin bootstrap: reset environment, login as `admin@runtimeroasters.com`, initialize data, verify users/resources/assignments.
- Farm to warehouse: Farm Manager declares harvest; Warehouse Manager dispatches pickup; Driver completes route and milestones; warehouse intake is created.
- Warehouse to retail: Warehouse creates/finalizes batch; Store Manager creates order; Finance pass/fail controls payment; Warehouse reserves and dispatches; Driver delivers and returns; order completes or rolls back.
- Public provenance: public QR/trace page renders sanitized farm-to-cup journey for a sold cup/item.

The checklist is guidance only; pass/fail proof belongs in rows above or ticket evidence files.
