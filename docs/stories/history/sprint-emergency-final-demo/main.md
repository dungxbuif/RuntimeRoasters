# Emergency Sprint: Final Demo Flow Recovery

**Sprint Type:** Emergency integration sprint  
**Source of Truth:** `MISING_IMPLEMENTATION_SPEC.md`  
**Goal:** Restore the missing end-to-end demo flow and make it testable from traceId propagation through public QR traceability.

---

## 1. Execution Order

This sprint must be executed in dependency order. Do not start UI-heavy demo work before the traceId and event backbone are verified.

| Order | Ticket | Summary | Why First/Next |
| :--- | :--- | :--- | :--- |
| 1 | [RR-URG-01](./RR-URG-01-traceid-verification.md) | Verify and fix traceId propagation end-to-end | Foundation for trace UI, debugging, and QR trace demo |
| 2 | [RR-URG-02](./RR-URG-02-event-contracts-storage-boundaries.md) | Standardize event contracts, IDs, and storage boundaries | Prevents downstream services/UI from using incompatible payloads |
| 3 | [RR-URG-03](./RR-URG-03-warehouse-http-and-pickup.md) | Add Warehouse HTTP APIs and harvest-to-pickup flow | Replaces broken Harvest -> Intake shortcut |
| 4 | [RR-URG-04](./RR-URG-04-logistics-driver-return.md) | Implement shipment legs, vehicles, Driver Client updates, mandatory return | Physical logistics backbone |
| 5 | [RR-URG-05](./RR-URG-05-paid-order-fulfillment.md) | Paid order -> payment -> reservation -> delivery -> return | Downstream retail flow |
| 6 | [RR-URG-06](./RR-URG-06-realtime-notifications.md) | Socket/SSE and persisted notifications | Makes two-screen demo visible |
| 7 | [RR-URG-07](./RR-URG-07-trace-service-public-qr.md) | Trace-service ES read model, trace-history, public QR demo | Public showcase proof |
| 8 | [RR-URG-08](./RR-URG-08-role-auth-demo-data.md) | Role scoping, demo accounts, seeded flow data | Ensures demo users can execute only intended actions |
| 9 | [RR-URG-09](./RR-URG-09-client-demo-flows.md) | Client dashboards and Driver Client simulation UI | Final visible experience |
| 10 | [RR-URG-10](./RR-URG-10-demo-runbook-e2e.md) | Demo guide and E2E verification checklist | Final acceptance run |

---

## 2. Context Map For New Implementers

Read this sprint in the following order if you are new to the project:

1. `MISING_IMPLEMENTATION_SPEC.md` for the full business context and investigation findings.
2. `docs/product/domain/README.md` for canonical flow and role principle.
3. `docs/product/domain/README.md` for who can press which UI button.
4. `docs/product/TECH.md` for Postgres/JSONB/Elasticsearch/Cassandra/Valkey boundaries.
5. This sprint main file for execution order.
6. RR-URG tickets in numeric order.

### Why This Emergency Sprint Exists

Previous sprint slices implemented farm, warehouse, retail, payment, logistics, trace, and audit separately. The missing context is the physical flow between them:

- harvest currently jumps too quickly into warehouse intake.
- logistics does not model pickup/return legs clearly.
- driver simulation is not intentionally tied to authenticated driver actions.
- retail final demo must be paid order, not replenishment.
- public QR trace must show real prepared trace data.
- traceId continuity is not proven across the whole system.

### Canonical Demo Story

Farm side:

1. `FARM_MANAGER` creates harvest.
2. Warehouse creates pickup request, not intake.
3. `WAREHOUSE_MGR` dispatches vehicle/driver.
4. `DRIVER` logs in, starts route simulation, confirms pickup/loading.
5. `DRIVER` returns to warehouse/base.
6. Warehouse creates intake.
7. Warehouse processing/inventory continues.

Retail side:

1. `STORE_MGR` creates paid order.
2. Payment succeeds or simulated payment succeeds.
3. Warehouse reserves inventory.
4. `WAREHOUSE_MGR` dispatches vehicle/driver.
5. `DRIVER` delivers to store.
6. `DRIVER` returns to warehouse/base.
7. Retail order becomes `COMPLETED`.

Public trace side:

1. Public user opens no-auth QR trace page.
2. Client shows QR codes for seeded products.
3. QR opens a real trace document from trace-service/Elasticsearch.

---

## 3. Non-Negotiable Decisions

- Retail demand for final demo is **paid order**.
- Payment may be simulated, but warehouse reservation/delivery must start from paid or simulated-paid order.
- Driver return-to-base is mandatory for farm pickup and retail delivery.
- Retail order `COMPLETED` waits for driver return-to-base to be recorded.
- Trace-history belongs to `trace-service`; do not introduce a separate monitor microservice in this sprint.
- SigNoz + ClickHouse is the observability backend for OpenTelemetry waterfall evidence.
- Kafka remains Apache Kafka without ZooKeeper. The SigNoz ClickHouse coordinator is a separate observability dependency and must not be reused by Kafka.
- Elasticsearch is the business traceability read model.
- Cassandra is append-only audit and optional trace-service live history.
- PostgreSQL is source of truth for operational state.
- PostgreSQL JSONB is metadata/payload support, not a replacement for typed workflow/scope columns.
- Public root/QR trace pages may be no-auth only with sanitized data.
- Private dashboards, Driver Client, and private realtime streams require auth and record-level scoping.

---

## 4. Ticket Ownership Boundaries

| Concern | Primary Ticket | Supporting Tickets |
| :--- | :--- | :--- |
| traceId across gateway/services/Kafka/DB/logs | RR-URG-01 | RR-URG-02, RR-URG-07 |
| canonical event names and payloads | RR-URG-02 | all backend tickets |
| farm harvest creates pickup request | RR-URG-03 | RR-URG-04, RR-URG-06 |
| warehouse receive creates intake | RR-URG-03 | RR-URG-04 |
| driver GPS and milestone confirmations | RR-URG-04 | RR-URG-09 |
| mandatory driver return | RR-URG-04 | RR-URG-05, RR-URG-07 |
| paid order SAGA | RR-URG-05 | RR-URG-02, RR-URG-03, RR-URG-04 |
| realtime dashboard updates | RR-URG-06 | RR-URG-09 |
| public QR trace | RR-URG-07 | RR-URG-08, RR-URG-09 |
| auth scope and demo seed data | RR-URG-08 | all UI/backend tickets |
| frontend dashboards and Driver Client | RR-URG-09 | RR-URG-03 to RR-URG-08 |
| final user guide | RR-URG-10 | all tickets |

---

## 5. Definition Of Done

- [x] RR-URG-01: A single paid order or harvest demo produces connected `trace_id` across gateway, services, DB spans, Kafka publish/consume, and downstream handlers.
- [x] RR-URG-01: SigNoz at `http://localhost:3301` shows OTel spans for the same `trace_id` used by trace-service evidence.
- Farm harvest creates pickup request, not intake.
- Warehouse dispatches driver/vehicle for farm pickup.
- Driver Client simulates route and posts GPS/status updates while logged in as `DRIVER`.
- Driver confirms pickup and mandatory return-to-base.
- Warehouse creates intake only after returned pickup.
- Processing/inventory can continue from received intake.
- Paid order reserves inventory with concurrency guard.
- Warehouse dispatches retail delivery.
- Driver confirms delivery and mandatory return-to-base.
- Trace-service projects real journey to Elasticsearch.
- Public QR trace page opens real prepared trace document.
- Role matrix is enforced for `ADMIN`, `FARM_MANAGER`, `WAREHOUSE_MGR`, `STORE_MGR`, `DRIVER`, and optional `PROCESSOR`.
- Demo runbook lists accounts, windows, buttons, expected realtime visuals, and recovery steps.

---

## 6. Master Checklist

### TraceId Foundation

- [x] KrakenD forwards configured W3C `traceparent`.
- [x] Gin/gRPC handlers start child spans with the incoming context.
- [x] Kafka producer injects `traceparent` into message headers.
- [x] Kafka consumer extracts `traceparent` before business handling.
- [x] GORM/Postgres tracing plugin is installed under shared DB bootstrap.
- [x] Logs include `trace_id` from OTel context.
- [x] One E2E command proves a single trace across services.

### Backend Flow

- [ ] Harvest event creates pickup request and notification.
- [ ] Direct `HarvestCreated -> Intake` shortcut removed from normal flow.
- [ ] Warehouse HTTP APIs match frontend/KrakenD paths.
- [ ] Logistics supports shipment type, legs, vehicle, driver, route, confirmations, mandatory return.
- [ ] Driver location update validates assigned shipment.
- [ ] Paid order flow includes payment success/simulated success before reservation.
- [ ] Inventory reservation uses Valkey/Redlock or equivalent lock.
- [ ] Trace-service projects all physical logistics milestones.

### Frontend And Demo

- [ ] Farm dashboard creates harvest and shows pickup status.
- [ ] Warehouse dashboard shows pickup/outbound queues and dispatch actions.
- [ ] Driver Client shows assigned shipment, Start button, route animation, and milestone buttons.
- [ ] Logistics dashboard shows live shipment/GPS state from backend.
- [ ] Retail dashboard creates paid order and shows incoming delivery.
- [ ] Public ArchitectureTopology remains no-auth with sanitized data.
- [ ] Public QR trace page shows real seeded trace documents.

### Security And Roles

- [ ] `ADMIN` setup/assignment/overview only by default.
- [ ] `FARM_MANAGER` scoped to assigned farms.
- [ ] `WAREHOUSE_MGR` scoped to assigned warehouses.
- [ ] `STORE_MGR` scoped to `store_ids`.
- [ ] `DRIVER` scoped to assigned shipments.
- [ ] Private socket/SSE streams are authenticated and role-scoped.

### Final Demo

- [ ] Demo accounts verified.
- [ ] Seed data includes at least 2 public trace products.
- [ ] Two-screen demo script works.
- [ ] Recovery/fast-forward steps are documented.
- [ ] Full backend tests pass or known unrelated failures are documented.
