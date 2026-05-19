# RuntimeRoasters Development Roadmap (Master Plan)

**Core Goal:** Build an exemplary Microservices system — distributed architecture, Event-Driven, Saga Pattern, Observability (Priority on Manual DI & Selective Outbox).

---

## 🚀 Roadmap Overview (Sprints 4-12)

| Sprint | Goal | Key Deliverables (Business + Tech) |
| :--- | :--- | :--- |
| **Emergency** | **Final Demo Flow Recovery** | TraceId verification, end-to-end Farm -> Warehouse -> Paid Order -> Retail logistics, public QR trace demo. |
| **S4** | **The Resilient Farm** | Transactional Outbox, Reliable Harvesting API, CloudEvents Blueprint. |
| **S5** | **Warehouse & Processing** | Unified Warehouse Service, Coffee Processing State Machine, Inventory & Saga Participant. |
| **S6** | **Retail & Order Saga** | Retail Service, Saga Orchestrator (Choreography), Idempotent Ordering. |
| **S7** | **Real-time Logistics** | Logistics Service, Driver Trips, Valkey GEO Real-time Tracking. |
| **S8** | **Financial Integrity** | Payment Service, Stripe Integration, Webhook HMAC, Saga Rollbacks. |
| **S9** | **Traceability (CQRS)** | Trace Service, Elasticsearch Read-Model, 360-degree Bean Journey View. |
| **S10** | **Observability & Audit** | Audit Service (Cassandra), System Health Dashboards. |
| **S11** | **Hardening & Finale** | mTLS Zero Trust, Token Revocation, PgBouncer, Chaos Engineering. |
| **S12** | **Project Polish** | Final technical debt cleanup, Load testing, and Final Demo. |

---

## 📋 Sprint Details (Progress)

### [x] Sprint 1-3: Foundation & Farm Core (Completed)
Completion of basic infrastructure and CRUD APIs for the Farm Service.

### [x] Sprint 4: The Resilient Farm (Completed)
- **Goal:** Ensure harvest data is never lost and clean up Sprint 3 technical debt.
- **Tickets:** [RR-4.0], [RR-4.1], [RR-4.2].
- Details: [`docs/business/sprint4/main.md`](./sprint4/main.md)

### [ ] Emergency Sprint: Final Demo Flow Recovery
- **Goal:** Recover missing cross-sprint end-to-end demo flow before production-demo review.
- **Critical First Ticket:** Verify traceId propagation across the whole system.
- Details: [`docs/business/sprint-emergency-final-demo/main.md`](./sprint-emergency-final-demo/main.md)

### [ ] Sprint 5: Warehouse & Processing
- **Goal:** Consolidate factory and warehouse processes. Handle beans from post-harvest intake, through roasting to finished goods storage and inventory reservation for Saga.
- Details: [`docs/business/sprint5/main.md`](./sprint5/main.md)

### [ ] Sprint 6: Retail & Order Saga
- **Goal:** Automated ordering and supply chain coordination (Choreography Saga).
- Details: [`docs/business/sprint6/main.md`](./sprint6/main.md)

### [ ] Sprint 7: Real-time Logistics
- **Goal:** Real-time transport vehicle tracking on a map.
- Details: [`docs/business/sprint7/main.md`](./sprint7/main.md)

### [ ] Sprint 8: Financial Integrity
- **Goal:** Real payments and automated refunds.
- Details: [`docs/business/sprint8/main.md`](./sprint8/main.md)

### [ ] Sprint 9: Traceability (CQRS)
- **Goal:** High-speed origin traceability.
- Details: [`docs/business/sprint9/main.md`](./sprint9/main.md)

### [ ] Sprint 10: Observability & Audit
- **Goal:** Comprehensive monitoring and immutable audit logs.
- Details: [`docs/business/sprint10/main.md`](./sprint10/main.md)

### [ ] Sprint 11: Hardening & Grand Finale
- **Goal:** mTLS security, handling outstanding technical debt (PgBouncer, Token Revocation).
- Details: [`docs/business/sprint11/main.md`](./sprint11/main.md)

### [ ] Sprint 12: Project Polish
- **Goal:** Load testing, code cleanup, and overall Demo.
- Details: [`docs/business/sprint12/main.md`](./sprint12/main.md)
