# 📋 Master Strategy Discussion: Sprints 4 — 7

**Status:** `🚧 IN DISCUSSION`
**Participants:** PO, TechLead, BA, Secretary (Agent)
**Objective:** Agree on the overall business flow and how to apply advanced techniques for the remaining phases of Runtime Roasters.

---

## 🏗️ 1. The Big Picture

After completing **Sprint 3 (Farm Management)**, we have established the "Roots" of the supply chain. The next phase is building the "Trunk" and "Branches," transforming the isolated Microservices into a unified entity.

### Roadmap for the next 4 phases:
1.  **Distributed Transactions (Sprint 4)**: Ensure consistent Order & Inventory reservation (Saga).
2.  **Real-time Logistics (Sprint 5)**: Dispatch vehicles and track GPS (Valkey Geo).
3.  **Data Transparency (Sprint 6)**: Traceability and Monitoring (Elasticsearch/Cassandra).
4.  **Experience & Resiliency (Sprint 7)**: Central Control Dashboard and Zero Trust (Chaos/mTLS).

---

## 📝 2. Discussion Log

*(The Secretary will update the content here based on feedback from the PO and the team)*

### 📝 2. Discussion Log

#### Key Decisions:
- **Mock Payment**: Temporarily use `mock-payment-service` in the early phases to focus on perfecting the Saga flow. The actual Stripe integration will be deferred to the final Sprints.
- **Epic-based Sprints**: Break down the roadmap. Each Sprint will focus entirely on a single technical/business Epic to ensure quality and focus.

#### New Epic Roadmap (Sprint 4 — 12):

1. **Sprint 4 (Epic: Security & Core Infra Hardening)**
   - Content: Deploy PgBouncer, complete E2E Auth Integration and Token Revocation (Logout).
   - Objective: The infrastructure and security foundation reach "Production-ready" standards.

2. **Sprint 5 (Epic: Distributed Order Orchestration)**
   - Content: Retail Service, Transactional Design (Selective Outbox).
   - Objective: Record orders and publish events atomically.

3. **Sprint 6 (Epic: Inventory Consistency & Saga Lite)**
   - Content: Warehouse Service, Saga Participant, Inbox Pattern (Idempotency T2), Mock Payment logic.
   - Objective: Complete the 2-step Saga flow (Retail-Warehouse) ensuring idempotency.

4. **Sprint 7 (Epic: Logistics & Real-time Delivery)**
   - Content: Logistics Service, Basic Shipment management.
   - Objective: Automatically coordinate shipping once the warehouse reserves goods.

5. **Sprint 8 (Epic: Geographic Intelligence)**
   - Content: Valkey Geo integration, GPS Simulator, Real-time Tracking UI.
   - Objective: Track driver locations on the map in real-time.

6. **Sprint 9 (Epic: System-wide Traceability & Search)**
   - Content: Trace Service, CQRS Read Model, Elasticsearch Integration.
   - Objective: Search and retrieve the complete history of coffee beans in < 100ms.

7. **Sprint 10 (Epic: Reliability & Observability)**
   - Content: Full OpenTelemetry instrumentation, Prometheus/Grafana, Cassandra Audit Log.
   - Objective: A fully transparent system with an immutable Audit log.

8. **Sprint 11 (Epic: Real-world Commerce Integration)**
   - Content: Build the actual Payment Service, integrate Stripe API, perfect the 3-step Saga flow.
   - Objective: Process actual payments and handle the automated Refund mechanism.

9. **Sprint 12 (Epic: Control Plane & The Grand Finale)**
   - Content: Monitor Service (SSE), React Flow System Mesh, Chaos Control Panel, mTLS.
   - Objective: Visualize the entire Microservices "universe" and demonstrate self-healing capabilities.

---

## ✅ 3. Decisions & Action Items

*(Record roadblocks after the discussion)*

#### Details for Topic 1: Sprint 4 - Security & Core Infra Hardening

**1. PgBouncer: Connection Pooling Strategy**
- **Technical (TechLead)**:
    - Use the `bitnami/pgbouncer` image.
    - Mode: `transaction` (suitable for Go microservices using short-lived connections).
    - Config: Limit `default_pool_size=20` to Postgres, but allow `max_client_conn=500` from microservices.
    - Update: All service `.env` files will change from port `54321` (Postgres direct) to `6432` (PgBouncer).
- **Business (BA)**: Ensure system stability when scaling up the number of product batches and order transactions.

**2. E2E Auth Integration: Verifying the "Two-Gate Model"**
- **Technical (TechLead)**: Control 3 checkpoints:
    - `KrakenD (Gate 1)`: Verify scope (e.g., `farm:read`).
    - `gRPC Interceptor`: Extract identity from metadata for inter-service calls.
    - `OTel Correlation`: Attach `user_id` to every Trace span on SigNoz.
- **Business (BA)**: Absolutely protect farm data, ensuring privacy and security between Farm Managers.

**3. Token Revocation: Secure Logout Flow**
- **Technical (TechLead)**: Deploy Distributed Blacklist.
    - Mechanism: Store revoked `jti` in Valkey with TTL.
    - Logic: Middleware checks Valkey `EXISTS`. Choose the **Fail-closed** approach (block requests if Valkey is down) to maximize security.
- **Business (BA)**: Provide an actual Logout feature for users, protecting accounts in case of compromise.
