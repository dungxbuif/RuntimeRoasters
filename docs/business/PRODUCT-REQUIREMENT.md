# Runtime Roasters — Product Requirements Document (PRD)

| Field            | Value                                                    |
| :--------------- | :------------------------------------------------------- |
| **Project**      | Runtime Roasters                                         |
| **Author**       | PM (AI-assisted)                                         |
| **Status**       | Draft v1.2                                               |
| **Last Updated** | 2026-04-16                                               |
| **Type**         | Technical Showcase                           |
| **Developer**    | Solo                                                     |
| **Deployment**   | Docker on Cloud VPS (HA-ready design)                    |

---

## 1. Executive Summary

**Runtime Roasters** is a `Showcase` project simulating a coffee supply chain management platform from farm to retail cup (Farm-to-Cup). The primary objective is to demonstrate the ability to design and implement a production-grade distributed `Microservices` architecture, using 100% `Golang` for the `Backend` and `ReactJS` for the `Frontend`.

The project is **not** intended to solve a real-world commercial problem but instead focuses on:
- Showcasing complex `Design Patterns` (Saga, CQRS, Outbox, Event Sourcing).
- Proving the capability to handle `Distributed Systems`.
- Providing a visual demonstration of the architecture through a "Control Plane Visualization" `Dashboard`.

**Evaluation Audience:** Employers, Technical Reviewers, TAs (Technical Architects).

---

## 2. Goals & Non-Goals

### 2.1 Goals
| #  | Goal                                                                                   | Measurement                                         |
| :- | :------------------------------------------------------------------------------------- | :----------------------------------------------------- |
| G1 | Showcase `Microservices` architecture with 9+ independent services communicating via `Kafka`/`gRPC` | All services running and communicating on `Docker Compose` |
| G2 | Successfully implement the `Saga Pattern` (Choreography) with `Compensating Actions`    | Demo rollback scenarios when `Warehouse` is out of stock + Auto-Refund via `Stripe` |
| G3 | Implement `Transactional Outbox` + `Inbox` (Idempotency)                               | Demo disconnecting `Kafka`, ensuring data remains consistent after recovery |
| G4 | Implement `CQRS` with `PostgreSQL` (Write) + `Elasticsearch` (Read)                     | Coffee traceability retrieval in < 100ms               |
| G5 | Visual architecture dashboard (Control Plane Visualization)                            | Viewers see real-time data flow between services       |
| G6 | Design for `High Availability` (HA-ready)                                              | Architecture can scale horizontally without refactoring |
| G7 | Integrate `Payment Gateway` (Stripe) with full `Webhook Security` + `Idempotency`      | Demo B2B payments and auto-refunds on Saga failure     |

### 2.2 Non-Goals
- ❌ No native mobile app (responsive web only).
- ❌ No processing of GPS data from real hardware devices (using simulated data).
- ❌ No complex financial/accounting management system (basic payment flow only).
- ❌ No multi-tenant targeting (single-tenant demo only).
- ❌ No real money processing in the demo environment (using `Stripe Test Mode`).

---

## 3. User Roles & Personas

| Role         | System Code  | Short Description                                       | Primary Permissions                                          |
| :----------- | :----------- | :------------------------------------------------------- | :----------------------------------------------------------- |
| **Farm Manager** | `FARM_MANAGER` | Manages coffee farms, declares harvests                | CRUD farms, create harvest batches, view history             |
| **Processor** | `PROCESSOR`  | Manages processing factory, roasting                    | Receive raw beans, create roast batches, issue `Batch ID`, packaging |
| **Driver**   | `DRIVER`     | Drives transport vehicles, updates GPS                  | Accept trips, update delivery status, send GPS coordinates    |
| **Store Manager** | `STORE_MGR` | Manages retail stores                                  | View inventory, create supply requests, receive goods        |
| **Admin**    | `ADMIN`      | System Administrator                                   | View full `Dashboard`, `Audit log`, manage users             |
| **Customer** | `end_user` | End consumer (no login required)                       | Scan QR, view traceability                                   |

---

## 4. Business Flows

### 4.1 Main Flow: Farm-to-Cup Pipeline

```
[Farm Manager]  [Processor]       [Warehouse]    [Payment]     [Logistics]    [Retail]      [End User]
   │                │                  │              │              │             │              │
   ├─ Harvest ────►│                  │              │              │             │              │
   │ (Harvest       │                  │              │              │             │              │
   │  Created)      ├─ Roasting ─────►│              │              │             │              │
   │                │ (BatchProcessed) │              │              │             │              │
   │                │                  ├─ Intake ────►│             │             │              │
   │                │                  │(InventoryAdd)│              │             │              │
   │                │                  │              │              │             ├─ Order       │
   │                │                  │              │              │             │(OrderCreated)│
   │                │                  │              │◄─────────────┤─────────────┤              │
   │                │                  │              │ Payment      │             │              │
   │                │                  │              │(PaymentIntent│             │              │
   │                │                  │              │  Created)    │             │              │
   │                │                  │              │──► Stripe ──►│             │              │
   │                │                  │              │  (Webhook)   │             │              │
   │                │                  │◄─────────────┤──────────────┤             │              │
   │                │                  │ Reserve stock│              │             │              │
   │                │                  │(StockReserved│              │             │              │
   │                │                  │              │              │◄────────────┤              │
   │                │                  │              │              │ Dispatch    │              │
   │                │                  │              │              ├─ GPS ──────►│              │
   │                │                  │              │              ├─ Delivery ─►│              │
   │                │                  │              │              │             │              │
   │                │                  │              │              │             │  Scan QR ────┤
   │                │                  │              │              │             │  Traceability│
```

### 4.2 Saga Flow: Supply Order (with Payment)

This is the most complex flow, implementing the `Saga Pattern` (Choreography) combined with a `Payment Gateway`:

| Step | Service          | Action                                | Event Emitted                | Failure → Compensation                           |
| :--- | :--------------- | :------------------------------------ | :--------------------------- | :----------------------------------------------- |
| 1    | `Retail`         | Store creates supply order            | `SupplyOrderCreated`         | —                                                |
| 2    | `Payment`        | Create `PaymentIntent` on Stripe      | `PaymentIntentCreated`       | —                                                |
| 3    | Frontend         | Display Stripe payment form           | —                            | User cancels → `PaymentCancelled` → Cancel order |
| 4    | `Payment`        | Receive Stripe Webhook, confirm payment| `PaymentCompleted`           | `PaymentFailed` → Cancel order                   |
| 5    | `Warehouse`      | Check & reserve stock                 | `StockReserved`              | `StockReserveFailed` → Refund Stripe + Cancel order |
| 6    | `Logistics`      | Find vehicle & create trip            | `ShipmentAssigned`           | `ShipmentFailed` → Release stock + Refund Stripe  |
| 7    | `Logistics`      | Driver accepts and transports         | `ShipmentInTransit`          | —                                                |
| 8    | `Logistics`      | Successful delivery                   | `ShipmentDelivered`          | —                                                |
| 9    | `Warehouse`      | Confirm stock deduction               | `StockDeducted`              | —                                                |
| 10   | `Retail`         | Receive goods, update inventory       | `SupplyOrderCompleted`       | —                                                |

**Compensating Actions (Rollback):**
- If Step 5 fails (out of stock) → `Payment` receives `StockReserveFailed`, calls `Stripe Refund API`, emits `PaymentRefunded`. `Retail` receives it → order changes to `REJECTED`.
- If Step 6 fails (no vehicle available) → `Warehouse` receives `ShipmentFailed`, releases stock. `Payment` receives the event, calls `Stripe Refund API`.
- **Principle:** If money was deducted on Stripe, there must be a path to refund it. Never leave money "hanging."

### 4.3 Payment Flow Details

#### 4.3.1 Payment Sequence

```
  [Frontend/POS]         [Payment Service]           [Stripe]             [Kafka]
       │                        │                       │                    │
       │── POST /orders ───────►│                       │                    │
       │   (SupplyOrderCreated  │                       │                    │
       │    from Kafka)         │                       │                    │
       │                        ├── Create PaymentIntent►│                   │
       │                        │◄── client_secret ──────┤                   │
       │◄── client_secret ──────┤                       │                    │
       │                        │                       │                    │
       │── Stripe.js submit ───►│───────────────────────►│                   │
       │   (card form)          │                       │                    │
       │                        │                       │                    │
       │                        │◄── Webhook POST ──────┤                   │
       │                        │   (payment_intent.     │                   │
       │                        │    succeeded)          │                   │
       │                        │                       │                    │
       │                        ├── HMAC verify ────────►│ (validate sig)    │
       │                        ├── Check Inbox ────────►│ (idempotency)     │
       │                        ├── Save + Outbox ──────►│                   │
       │                        │                       │  ──► PaymentCompleted
       │                        │                       │                    │
```

#### 4.3.2 Payment State Machine

```
                    ┌──────────────────────────────────────────┐
                    │                                          │
  ┌─────────┐   PaymentIntent   ┌──────────────┐              │
  │ CREATED ├──── created ─────►│   PENDING    │              │
  └─────────┘                   └──────┬───────┘              │
                                       │                      │
                          ┌────────────┼────────────┐         │
                     Webhook OK    Webhook FAIL   User cancel  │
                          │            │            │         │
                   ┌──────▼──────┐ ┌───▼──────┐ ┌──▼──────┐  │
                   │  SUCCEEDED  │ │  FAILED  │ │CANCELLED│  │
                   └──────┬──────┘ └──────────┘ └─────────┘  │
                          │                                   │
                   Saga Failure                                │
                   (Stock/Ship fail)                           │
                          │                                   │
                   ┌──────▼──────┐                            │
                   │  REFUNDED   ├────────────────────────────┘
                   └─────────────┘
```

**`Payment` Statuses:**

| Status       | Description                                                 |
| :----------- | :---------------------------------------------------------- |
| `CREATED`    | Order just created, waiting for `PaymentIntent` creation    |
| `PENDING`    | `PaymentIntent` created on Stripe, waiting for user payment |
| `SUCCEEDED`  | Stripe confirms successful payment (via Webhook)            |
| `FAILED`     | Stripe reports payment failure (card declined, etc.)        |
| `CANCELLED`  | User cancels payment before completion                      |
| `REFUNDED`   | Refunded via Stripe Refund API (due to Saga compensation)   |

#### 4.3.3 Payment Service Database Schema

```sql
-- Main Table: Record payments
CREATE TABLE payments (
    id                       UUID PRIMARY KEY,
    order_id                 UUID NOT NULL UNIQUE,
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    amount                   DECIMAL(12,2) NOT NULL,
    currency                 VARCHAR(3) DEFAULT 'VND',
    status                   VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    stripe_refund_id         VARCHAR(255),
    created_at               TIMESTAMPTZ DEFAULT NOW(),
    updated_at               TIMESTAMPTZ DEFAULT NOW()
);

-- Inbox: Anti-duplication for Webhook processing (Idempotency)
CREATE TABLE inbox_stripe_events (
    event_id    VARCHAR(255) PRIMARY KEY,  -- Stripe event ID (evt_xxx)
    event_type  VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);

-- Outbox: Ensure events reach Kafka
CREATE TABLE outbox_events (
    id            UUID PRIMARY KEY,
    aggregate_id  UUID NOT NULL,
    event_type    VARCHAR(100) NOT NULL,
    payload       JSONB NOT NULL,
    published     BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);
```

#### 4.3.4 Webhook Security (HMAC Validation)

Every Webhook from Stripe must pass through `HMAC` middleware before processing:

1. Stripe sends `POST /api/v1/webhooks/stripe` with the `Stripe-Signature` header.
2. The `Payment Service` uses the `Webhook Secret` to hash the received payload using `crypto/hmac` (Go).
3. Compares the calculated hash with `Stripe-Signature` → if they match, process; otherwise, reject (`HTTP 403`).
4. **Principle:** Never trust any request to `/webhooks/*` without `HMAC verification`.

> **Future Expansion:** The `Payment Service` architecture is designed following the `Payment Gateway Abstraction`. Stripe is the first implementation. Other providers (VNPay, MoMo, PayOS) can be added by implementing the same interface without affecting the Saga flow.

### 4.4 QR Traceability

When an `end_user` scans the QR code on the coffee packaging, the system returns:

| Information           | Source Service  | Example                                    |
| :-------------------- | :-------------- | :---------------------------------------- |
| Farm Name             | `Farm`          | "Son La Farm - Anh Minh"                  |
| Cultivation Area      | `Farm`          | "Plot A3, Son La, altitude 1200m"         |
| Coffee Variety        | `Farm`          | "Arabica Catimor"                         |
| Harvest Date          | `Farm`          | "2026-01-15"                              |
| Processing Method     | `Process`       | "Washed Process"                          |
| Roasting Date         | `Process`       | "2026-02-01"                              |
| Roast Level           | `Process`       | "Medium Roast"                            |
| Batch ID              | `Process`       | `RR-SL-2026-0215-A3`                      |
| Transport Route       | `Logistics`     | "Son La → Hanoi (320km, 6h)"              |
| Shop Intake Date      | `Retail`        | "2026-02-03"                              |

**Batch ID format:** `RR-{REGION_CODE}-{YEAR}-{MMDD}-{LOT}`

---

## 5. System Architecture Overview

### 5.1 Architectural Foundations

#### Codebase: Monorepo
The entire project resides in a single repository with a shared structure:
- **`api/`** — Shared `gRPC` `Proto` definitions (contracts between services).
- **`pkg/`** — Shared libraries: `Database Wrapper`, `Middleware` (`Auth`, `Tracing`, `Logging`), `Kafka Producer/Consumer` base.
- **`services/`** — Each Microservice is an independent module.

#### Service Structure: Clean Architecture (go-clean-arch v4)
Each service applies `Clean Architecture` with 3 distinct layers:

| Layer             | Content                                                                       |
| :---------------- | :---------------------------------------------------------------------------- |
| **Domain**        | Entities, Value Objects, Repository Interfaces. Completely framework-agnostic. |
| **UseCase**       | Business logic, Saga step handlers, Command/Query handlers.                   |
| **Infrastructure**| `Gin` HTTP handlers, `gRPC` servers, `PostgreSQL` repos, `Kafka` producers.   |

#### Infrastructure: Polyglot Persistence
Multi-modal storage, using the right tool for each data type:

| Database          | Role                                                          |
| :---------------- | :------------------------------------------------------------ |
| `PostgreSQL`      | Core `ACID` transactions — `Source of Truth` for each Microservice |
| `Apache Cassandra`| Perpetual raw Event storage (`Audit` / `Event Sourcing`)      |
| `Elasticsearch`   | High-speed lookup, `CQRS Read Model` for traceability         |
| `Valkey`          | `Cache`, `Distributed Lock`, real-time `GPS` coordinates       |

#### Deployment: Docker on Proxmox
Deployed 100% via `Docker Compose`. The host environment is a self-managed `Proxmox` server. All infrastructure (DBs, Kafka, Services) runs as containers, easily migratable to a Cloud VPS.

### 5.2 Service Inventory

| #  | Service                | Protocol                    | Database              | Kafka Role        |
| :- | :--------------------- | :-------------------------- | :-------------------- | :---------------- |
| 1  | `API Gateway`          | HTTP `:8081`                | —                     | —                 |
| 2  | `Identity Service`     | HTTP `:4433`                | Ory Kratos            | —                 |
| 3  | `Webhook Service`      | HTTP `:8092` / gRPC `:50062`| PostgreSQL (Inbox)    | Producer only     |
| 4  | `Farm Service`         | HTTP `:8083` / gRPC `:50053`| PostgreSQL            | Producer          |
| 5  | `Process Service`      | HTTP `:8084` / gRPC `:50054`| PostgreSQL            | Producer/Consumer |
| 6  | `Warehouse Service`    | HTTP `:8085` / gRPC `:50055`| PostgreSQL            | Producer/Consumer |
| 7  | `Retail Service`       | HTTP `:8086` / gRPC `:50056`| PostgreSQL            | Producer/Consumer |
| 8  | `Logistics Service`    | HTTP `:8087` / gRPC `:50057`| Valkey                 | Producer/Consumer |
| 9  | `Payment Service`      | HTTP `:8088` / gRPC `:50058`| PostgreSQL            | Producer/Consumer |
| 10 | `Trace Service`        | HTTP `:8089` / gRPC `:50059`| Elasticsearch         | Consumer          |
| 11 | `Audit Service`        | HTTP `:8091` / gRPC `:50061`| Apache Cassandra      | Consumer          |
| 12 | `Monitor Service`      | HTTP `:8090` (SSE/WS)       | —                     | Consumer          |

#### `Webhook Service` (Ingress Gateway) — Details

The `Webhook Service` is a specialized gateway, **completely separate** from the `API Gateway`, responsible for receiving all external data streams:

| Source           | Endpoint                         | Processing                                               |
| :--------------- | :------------------------------- | :------------------------------------------------------- |
| `Stripe`         | `POST /webhooks/stripe`          | `HMAC-SHA256` verify → `Inbox` dedup → publish to Kafka  |
| `VNPay`          | `POST /webhooks/vnpay`           | `HMAC-SHA512` verify → `Inbox` dedup → publish to Kafka  |
| `IoT Device`     | `POST /webhooks/iot/gps`         | Token auth → normalization → publish `logistics.gps.updated` |

**Principle:** The `Webhook Service` **does not contain business logic**. Its sole task is: signature verification → dedup using `Inbox` → payload normalization → publish to `Kafka`. Business logic is handled in the corresponding service.

### 5.3 Kafka Topic Design (Draft)

| Topic                              | Producer(s)      | Consumer(s)                         |
| :--------------------------------- | :--------------- | :---------------------------------- |
| `auth.user.events`                 | Auth             | Farm, Audit, Gateway (Cache)        |
| `farm.harvest.created`             | Farm             | Processing, Trace, Audit            |
| `process.batch.completed`          | Processing       | Warehouse, Trace, Audit             |
| `retail.order.created`             | Retail           | Payment, Trace, Audit               |
| `payment.intent.created`           | Payment          | Retail, Trace, Audit                |
| `payment.completed`                | Payment          | Warehouse, Retail, Trace, Audit     |
| `payment.failed`                   | Payment          | Retail, Trace, Audit                |
| `payment.refunded`                 | Payment          | Retail, Trace, Audit                |
| `warehouse.stock.reserved`         | Warehouse        | Logistics, Payment, Trace, Audit    |
| `warehouse.stock.reserve-failed`   | Warehouse        | Payment, Retail, Trace, Audit       |
| `warehouse.stock.released`         | Warehouse        | Retail, Trace, Audit                |
| `logistics.shipment.assigned`      | Logistics        | Warehouse, Retail, Trace, Audit     |
| `logistics.shipment.failed`        | Logistics        | Warehouse, Payment, Retail, Trace, Audit |
| `logistics.shipment.delivered`     | Logistics        | Warehouse, Retail, Trace, Audit     |
| `logistics.gps.updated`            | Logistics        | Monitor, Trace                      |
| `retail.order.completed`           | Retail           | Trace, Audit                        |

### 5.4 HA-Ready Design Principles

Although the demo is deployed on a single-node `Docker Compose`, the architecture is pre-designed to scale:

| Component           | HA Strategy                                                      |
| :------------------ | :--------------------------------------------------------------- |
| `Go Services`       | Stateless — horizontal scale by adding container replicas         |
| `PostgreSQL`        | Each service has its own DB (database-per-service) → independent scaling |
| `Kafka`             | Multi-partition topics, consumer groups for parallel processing   |
| `Valkey`            | Supports Cluster mode (Sentinel/Cluster) for GPS data            |
| `Elasticsearch`     | Shard/Replica strategy for read-model                            |
| `API Gateway`       | Stateless, can be placed behind a Load Balancer                  |

### 5.5 Technical Mechanisms

#### A. Decentralized Authorization

The `API Gateway` only handles **Authentication**: validating the `JWT` and extracting `Roles` from the token, then passing them down to services via `HTTP Headers` (`X-User-ID`, `X-User-Roles`).

Each Microservice integrates `Casbin` into its internal `Middleware` layer. Based on its own `policy.csv` file, the service **self-authorizes access** to each API endpoint — independent of the Gateway, increasing autonomy and reducing central load.

```
[Client] ──► [API Gateway] ──► JWT validate + extract Roles ──► Header: X-Roles=STORE_MGR
                                                                       │
                                                               [Retail Service]
                                                               Casbin Middleware
                                                               policy.csv: STORE_MGR CAN POST /orders
                                                                       │
                                                               ✅ Allow / ❌ Deny
```

#### B. Dual Idempotency

The system protects against duplication at **two independent levels**:

| Level | Request Type | Mechanism | Storage | TTL |
| :--- | :----------- | :----- | :------ | :-- |
| **Level 1** (Synchronous) | HTTP REST API | Header `Idempotency-Key` stored in `Valkey` | `Valkey` | 24h |
| **Level 2** (Asynchronous) | Kafka Consumer + Webhook | `Inbox Pattern` writing `event_id` to `PostgreSQL` | `PostgreSQL` | Permanent |

- **Level 1:** Client sends `Idempotency-Key: <uuid>` in the header. `API Gateway` middleware checks `Valkey`. If the key exists → return cached response, no re-processing.
- **Level 2:** Consumer (Kafka/Webhook) extracts `event_id`, opens a `Transaction`: checks the `inbox_events` table → if present → rollback and skip → if absent → save and process.

#### C. Configuration Management

Uses a combination of `.env` files and the `viper` library (Go):
- Each service has its own configuration file (`config.yaml` + `.env` override).
- `viper` automatically reads environment variables and supports hot-reloading.
- `Docker Compose` passes `Secrets` and `DB Host` via the `environment` block.
- No hardcoded credentials in the source code.

#### D. Distributed Tracing

`OpenTelemetry` + `SigNoz/ClickHouse`. Workflow:

1. `API Gateway` assigns a unique `Trace-ID` to each incoming request.
2. The `Trace-ID` is injected into:
   - `gRPC Context` (metadata) when calling internal services.
   - `HTTP Headers` (`traceparent`) when calling external services.
   - `Kafka Message Headers` when publishing events.
3. Each service creates a child `Span`, linked to the root `Trace-ID`.
4. `SigNoz UI` renders a waterfall/timeline displaying the entire transaction across multiple services.

> **Showcase value:** A `Trace-ID` from the store's order placement to Saga completion (through Payment → Warehouse → Logistics) is visualized intuitively on SigNoz — very impressive for TA reviewers.

---

## 6. Phased Delivery Plan

### Phase 1: Foundation
> **Goal:** Build the Monorepo skeleton, CI pipeline, and the first 3 core services.

| Deliverable                                 | Pattern showcase                     |
| :------------------------------------------ | :----------------------------------- |
| Project scaffold (Monorepo + shared libs)   | `DDD` project structure              |
| `API Gateway` (Gin + JWT parse + routing)   | `API Gateway Pattern`, `Rate Limiting` |
| `Identity Service` (Ory Kratos integration) | `OAuth2`/`OIDC`                      |
| `Farm Service` (Farm CRUD + harvest)        | `Transactional Design (Outbox)`      |
| Kafka + PostgreSQL infra (Docker Compose)   | `Event-Driven` base                  |
| Casbin middleware (shared lib)              | `Decentralized Authorization`        |
| gRPC proto definitions (shared)            | `Protocol Buffers`                   |

### Phase 2: Supply Chain + Payment
> **Goal:** Complete the Farm → Process → Warehouse pipeline, integrate Stripe, and full Saga flow.

| Deliverable                                    | Pattern showcase                         |
| :--------------------------------------------- | :--------------------------------------- |
| `Processing Service` (roasting, Batch ID)      | `Event-Driven Consumer/Producer`         |
| `Warehouse Service` (inventory, reserve/release)| `Saga Participant`                       |
| `Retail Service` (supply order placement)      | `Saga Orchestrator`                      |
| `Payment Service` (Stripe PaymentIntent + Webhook) | `Webhook HMAC`, `Inbox Pattern`      |
| Payment Gateway Abstraction (interface-based)   | `Strategy Pattern` (Stripe, VNPay, etc.) |
| Stripe Webhook + HMAC validation middleware     | `Zero Trust` (external data)             |
| Full Saga flow + Auto-Refund Compensation      | `Saga Pattern` (Choreography)           |
| `Inbox Pattern` (Idempotency on Consumer + Webhook) | `Transactional Inbox`             |

### Phase 3: Logistics & Real-time
> **Goal:** GPS tracking, transportation, and real-time data.

| Deliverable                                   | Pattern showcase                   |
| :-------------------------------------------- | :--------------------------------- |
| `Logistics Service` (dispatch, GPS tracking)  | `Geo-spatial` (Valkey GEO commands)|
| GPS simulator (fake DRIVER coordinates)        | Real-time data pipeline            |
| Valkey integration (cache + distributed lock)  | `Distributed Lock`, `Cache-aside`  |
| mTLS for gRPC between services                | `Zero Trust Architecture`          |

### Phase 4: Observability & Traceability
> **Goal:** CQRS read-model, audit trail, and distributed tracing.

| Deliverable                                    | Pattern showcase                  |
| :--------------------------------------------- | :-------------------------------- |
| `Traceability Service` (Kafka → Elasticsearch) | `CQRS` (Read-model projection)    |
| QR code generation + traceability lookup       | `CQRS Query` endpoint             |
| `Audit Service` (Kafka → Cassandra, hash chain)| `Event Sourcing` (Lite)           |
| `Hash Chaining` for data integrity             | `Data Integrity Pattern`          |
| OpenTelemetry + SigNoz integration             | `Distributed Tracing`             |
| Prometheus + Grafana dashboards                | `Observability Stack`             |

### Phase 5: Control Plane Visualization Dashboard (Frontend)
> **Goal:** ReactJS Dashboard combining Operational UI + System Visualization.

| Deliverable                                     | Tech showcase                    |
| :---------------------------------------------- | :------------------------------- |
| `Monitor Service` (Kafka → SSE/WebSocket)       | Real-time event broadcasting     |
| ReactJS app (Vite + React Flow + Framer Motion) | Modern frontend stack            |
| Split-screen: App View + System View            | `Service Mesh Visualization`     |
| Isometric 3D service map with animated edges    | `React Flow` animated edges      |
| Visual metaphors (Green Bean, Roasted Bean, etc.)| Domain-specific UI language      |
| Chaos Control panel (kill Kafka, set stock = 0) | `Resiliency` demonstration       |
| CQRS Time Machine (time slider retrieval)       | CQRS visual demo                 |

---

## 7. Non-Functional Requirements

| Category         | Requirement                                                         |
| :--------------- | :------------------------------------------------------------------ |
| **Performance**  | CQRS read query (traceability lookup) < 100ms                      |
| **Performance**  | GPS update latency < 500ms (from Logistics → Monitor → Dashboard)  |
| **Availability** | HA-ready design: Stateless services, database-per-service           |
| **Scalability**  | Kafka multi-partition, consumer group ready                          |
| **Security**     | JWT Authentication (Ory Kratos), Casbin Authorization per-service   |
| **Security**     | mTLS for all gRPC internal communication                            |
| **Security**     | HMAC signature validation for Stripe Webhook (anti-spoofing)        |
| **Security**     | Stripe Test Mode only — no real money processing in the demo        |
| **Idempotency**  | `Inbox Pattern` for Stripe Webhook (at-least-once → exactly-once)   |
| **Integrity**    | Hash chaining on Audit log (Cassandra) anti-tampering               |
| **Observability**| Distributed tracing (OpenTelemetry + SigNoz) on every request       |
| **Observability**| Prometheus metrics + Grafana dashboard for each service             |
| **Deployment**   | Full Docker Compose for local dev                                   |
| **Deployment**   | Docker images ready for Cloud VPS deployment                        |

---

## 8. Technical Decisions

| Decision                     | Choice                      | Reason                                                             |
| :--------------------------- | :-------------------------- | :----------------------------------------------------------------- |
| Backend language             | `Go` 1.22+                  | Performance, concurrency, ecosystem for microservices              |
| Code organization            | `Monorepo`                  | Share `proto`, `pkg` libs; easy cross-service change management    |
| Service architecture         | `Clean Architecture` v4     | Separate Domain/UseCase/Infra — testable, replaceable adapters      |
| HTTP framework               | `Gin`                       | Lightweight, high-performance REST                                 |
| Internal RPC                 | `gRPC` + `Protobuf`         | Type-safe, high-speed, contracts defined in `api/`                 |
| Message broker               | `Apache Kafka`              | Industry standard for event-driven architecture                    |
| Relational DB                | `PostgreSQL` v15            | ACID, mature, database-per-service                                 |
| Document DB                  | `Apache Cassandra`          | Wide-column store, flexible schema for audit logs                  |
| Search engine                | `Elasticsearch`             | Full-text search + CQRS read-model                                 |
| Cache / Real-time            | `Valkey`                    | Valkey alternative, GEO commands for GPS, Idempotency-Key store    |
| Payment gateway              | `Stripe` (Test Mode)        | Industry standard, excellent API docs, Webhook support             |
| Payment abstraction          | `Strategy + Factory Pattern`| `PaymentProvider` interface + `ProviderFactory` → easy to add VNPay |
| Identity                     | `Ory Kratos`                | Open-source, self-hosted identity management                       |
| Authorization                | `Casbin` + `policy.csv`     | Embeddable RBAC, decentralized per-service, Gateway only authn    |
| Config management            | `viper` + `.env`            | Flexible, hot-reload, no hardcoded secrets, Docker-friendly        |
| Dependency Injection         | Manual DI                   | Avoid magic, easy to debug, Composition Root at main.go            |
| Reliability                  | Selective Outbox            | Applied to critical flows to ensure consistency                    |
| Frontend                     | `ReactJS` (Vite)            | Combined Operational UI + Control Plane Visualization Dashboard    |
| Visualization                | `React Flow`                | Node-based UI for service mesh visualization                       |
| Animation                    | `Framer Motion`             | Micro-animations, glow effects                                     |
| State management             | `Zustand`                   | Lightweight state for real-time WebSocket data                     |
| Tracing                      | `OpenTelemetry` + `SigNoz`  | Trace-ID from Gateway, propagate via gRPC/Kafka/HTTP headers       |
| Metrics                      | `Prometheus` + `Grafana`    | Industry standard monitoring stack                                 |
| Host infrastructure          | `Proxmox`                   | Self-hosted hypervisor, VM-based Docker environment                |
| Deployment                   | `Docker` + `Docker Compose` | Containerized, Cloud VPS ready                                     |

---

## 9. Risks & Mitigations

| Risk                                          | Impact | Mitigation                                              |
| :--------------------------------------------- | :----- | :------------------------------------------------------ |
| Solo developer → scope creep                   | High   | Strict phased delivery, MVP-first mindset                |
| Kafka learning curve                           | Medium | Start with single-partition, scale later                 |
| Too many databases to manage                   | Medium | Docker Compose manages entire infrastructure             |
| Visualization Dashboard too complex            | High   | Phase 5 — only build after backend is stable             |
| HA design cannot be verified on single-node    | Low    | Document HA strategy, verify via architecture review     |
| Stripe Webhook replay/spoofing                 | High   | HMAC validation + Inbox idempotency pattern              |
| Money "hanging" when Saga fails post-payment   | High   | Auto-refund compensation, strict Payment state machine    |
| Stripe API version changes                     | Low    | Abstract via interface, easy to swap provider            |

---

## 10. Success Criteria

The project is considered a **success** when:

- [ ] All 9+ services run synchronously on `Docker Compose` without errors.
- [ ] `Saga Rollback` scenario demonstrated (out of stock → auto-refund Stripe → order reversed).
- [ ] `Outbox Pattern` demonstrated (kill Kafka → recover → events flush successfully).
- [ ] `Inbox Pattern` demonstrated (duplicate Webhook sent → processed only once).
- [ ] Successful Stripe payment (Test Mode) and refund when Saga fails.
- [ ] QR scan returns full Farm-to-Cup history of a `Batch ID`.
- [ ] Control Plane Visualization Dashboard displays real-time data flow between services.
- [ ] Reviewers/TAs can read the code and understand each `Pattern` applied.
- [ ] System can be deployed to a Cloud VPS using `docker-compose up -d`.

---

## Appendix A: Glossary

| Term                   | Definition                                                                          |
| :--------------------- | :---------------------------------------------------------------------------------- |
| `Batch ID`             | Unique identifier for a finished coffee batch: `RR-{REGION}-{YEAR}-{MMDD}-{LOT}`    |
| `Outbox Pattern`       | Writing events to DB in the same transaction as business data, worker pushes later  |
| `Inbox Pattern`        | Storing processed events to avoid duplicate processing (idempotency)               |
| `Saga`                 | Sequence of distributed transactions with compensating actions for rollback          |
| `CQRS`                 | Separation of write (Command) and read (Query) into two different systems            |
| `Clean Architecture`   | 3-layer architecture (Domain/UseCase/Infra) — framework-agnostic, testable           |
| `Polyglot Persistence` | Using multiple types of databases, each suited to specific data characteristics      |
| `Hash Chaining`        | Each audit record contains the hash of the previous record, creating an immutable chain |
| `Control Plane Visualization` | Real-time dashboard monitoring the entire system architecture               |
| `PaymentIntent`        | Stripe object representing a payment transaction waiting for processing              |
| `Webhook Service`      | Specialized Ingress Gateway receiving and validating external data (Stripe, IoT)     |
| `HMAC`                 | Hash-based Message Authentication Code — digital signature for integrity            |
| `Dual Idempotency`     | 2-level protection: Valkey for HTTP (sync) + Inbox Pattern for Kafka/Webhook (async) |
| `Idempotency-Key`      | UUID created by the client, sent in HTTP header to prevent duplicate execution      |
| `ProviderFactory`      | Factory creating Payment adapters (Stripe/VNPay) based on config                    |
| `Compensating Action`  | Rollback action (e.g., Refund) when a Saga step fails                               |
| `Proxmox`              | Open-source self-hosted hypervisor platform, running Docker VM                      |
| `Trace-ID`             | Unique ID generated at the Gateway, propagated throughout the system for debugging   |

---

## Appendix B: Reference Documents

| Document                                                                  | Location                                        |
| :------------------------------------------------------------------------ | :---------------------------------------------- |
| System Architecture                                                       | `docs/architecture/system-architecture.md`       |
| Database Schema Design                                                    | `docs/architecture/database-schema.md`           |
| UI/UX Visual Ideas (Visualization Dashboard)                               | `docs/ui-ux/visual-ideas.md`                     |
| REST API Specifications                                                   | `docs/api/rest-api.md`                           |
| gRPC Contract Definitions                                                 | `docs/api/grpc-contracts.md`                     |
| Deployment & Setup Guide                                                  | `docs/deployment/setup-guide.md`                 |
