# 🌊 Technical Reference: Core Sequences & Rationale

This document provides a deep dive into the most complex technical flows of the **Runtime Roasters** platform. It explains not only *how* data moves but *why* specific distributed patterns were chosen.

---

## 1. The Distributed Order SAGA (Choreography)

We use **Choreographed Sagas** for order fulfillment. Unlike Orchestrated Sagas, there is no central "master" service; instead, each service knows what to do when it hears a specific event.

### Sequence Diagram
```mermaid
sequenceDiagram
    participant Retail as Retail Service
    participant WH as Warehouse Service
    participant Payment as Payment Service
    participant Kafka as Kafka Cluster

    Retail->>Retail: DB Transaction: Save Order (Pending) + Outbox
    Retail->>Kafka: Publish: OrderCreated
    Kafka-->>WH: Consume: OrderCreated
    WH->>WH: DB Transaction: Reserve Stock + Outbox
    WH->>Kafka: Publish: StockReserved
    Kafka-->>Payment: Consume: StockReserved
    Payment->>Payment: DB Transaction: Process Payment + Outbox
    Payment->>Kafka: Publish: PaymentCompleted
    Kafka-->>Retail: Consume: PaymentCompleted
    Retail->>Retail: DB Transaction: Mark Order as SUCCESS
```

### Technical Rationale: Why Choreography?
- **Decoupling**: The Retail service doesn't need to know that the Warehouse or Payment services even exist. It only cares about publishing and consuming events.
- **Scalability**: Services can be scaled independently based on their load (e.g., the Warehouse service might need more replicas during a harvest season).

---

## 2. Distributed Authorization Sync

To achieve **Zero-Trust** security with high performance, we use a distributed Casbin model.

### Sequence Diagram
```mermaid
sequenceDiagram
    participant Auth as Auth Service
    participant Service as Any Service (e.g. Farm)
    participant Kafka as Kafka Cluster

    Note over Auth, Service: Stage 1: Bootstrap
    Service->>Auth: gRPC: GetFullSnapshot (On Startup)
    Auth-->>Service: Return all Casbin Policies
    Service->>Service: Load Policies into Memory

    Note over Auth, Service: Stage 2: Live Updates
    Auth->>Auth: Admin adds new permission
    Auth->>Kafka: Publish: auth.policy.changed
    Kafka-->>Service: Consume: auth.policy.changed
    Service->>Service: Hot-reload Enforcer in memory
```

### Technical Rationale: Why Distributed Auth?
- **Low Latency**: Authorization checks are performed locally in memory at the service level (ns vs ms). No gRPC round-trip is required for every request.
- **Resilience**: If the Auth Service or Kafka is temporarily down, the services can still function using their cached policies.

---

## 3. CQRS Read-Model Projection

The **Trace Service** provides a fast, searchable view of the supply chain history by projecting data into Elasticsearch.

### Sequence Diagram
```mermaid
sequenceDiagram
    participant Farm as Farm Service
    participant Kafka as Kafka Cluster
    participant Trace as Trace Service
    participant ES as Elasticsearch

    Farm->>Kafka: Publish: HarvestCreated
    Kafka-->>Trace: Consume: HarvestCreated
    Trace->>Trace: Enrich with Farm metadata
    Trace->>ES: Upsert Document: coffee_trace_idx
    Note right of ES: Fast Full-text Search Enabled
```

### Technical Rationale: Why CQRS?
- **Performance**: Relational databases (Postgres) are great for ACID transactions but slow for complex traceability joins. ES is optimized for high-speed retrieval.
- **Isolation**: Heavy search queries on the Dashboard don't impact the performance of the core business services.

---

## 4. The Journey of a Bean (End-to-End Lifecycle)

This diagram tracks the production-demo lifecycle of a coffee batch across the entire supply chain. The important change from earlier sprint shortcuts is that harvest does not become warehouse intake immediately; it must pass through logistics pickup first.

```mermaid
sequenceDiagram
    autonumber
    participant FM as Farm Manager
    participant FS as Farm Service
    participant WH as Warehouse Service
    participant WM as Warehouse Manager
    participant LG as Logistics Service
    participant DR as Driver Client
    participant RS as Retail Service
    participant SM as Store Manager
    participant RT as Realtime/Socket
    participant EU as End User

    FM->>FS: Declare Harvest (Batch-H)
    FS-->>WH: Event: farm.harvest.created
    WH-->>RT: Notify pickup request
    RT-->>WM: Warehouse dashboard update
    WM->>WH: Dispatch vehicle/driver
    WH-->>LG: Event: warehouse.pickup.requested
    LG-->>DR: Assigned pickup shipment
    DR->>LG: Start route simulation + GPS ticks
    LG-->>RT: logistics.gps.updated/status_changed
    DR->>LG: Confirm pickup/loading
    DR->>LG: Confirm return to warehouse
    LG-->>WH: Event: logistics.pickup.arrived_at_warehouse
    WH->>WH: Create Intake from returned pickup
    WH->>WH: Processing (Roasting)
    WH->>WH: Transform Batch-H -> Batch-S (Stock)
    SM->>RS: Create paid order
    RS-->>WH: Event: retail.order.created
    RS->>RS: Payment success or simulated payment success
    WH->>WH: Reserve stock with inventory lock
    WH-->>RT: Notify outbound dispatch
    WM->>WH: Dispatch retail delivery
    WH-->>LG: Event: warehouse.dispatch.requested
    DR->>LG: Start delivery simulation + GPS ticks
    DR->>LG: Confirm delivered
    LG-->>RS: Event: logistics.delivery.completed
    DR->>LG: Return to warehouse/base
    LG-->>RT: logistics.driver.returned_to_base
    RS->>RS: Mark demand completed
    EU->>EU: Scans QR Code
    EU->>EU: View Traceability Report (Farm -> WH -> Retail)
```

### UI/Role Notes

- `FARM_MANAGER` creates harvest and watches pickup status.
- `WAREHOUSE_MGR` dispatches vehicles/drivers and controls warehouse receipt/processing.
- `DRIVER` owns route simulation and milestone confirmation in the main demo flow.
- `STORE_MGR` creates paid retail orders and watches incoming delivery status.
- Root Client App `ArchitectureTopology` may be public/no-auth with sanitized topology/demo data.
- Private realtime dashboard streams should be authenticated and role-scoped.

---

## 5. Realtime Demo Stream

The production-demo UI should have enough delay/pacing for viewers to understand the flow.

```mermaid
sequenceDiagram
    autonumber
    participant UIA as Actor Dashboard
    participant UIB as Trace/Map/Topology
    participant Sock as Socket/SSE Service
    participant Kafka as Kafka
    participant Svc as Domain Service

    UIA->>Svc: User action with JWT
    Svc->>Svc: Validate role and record scope
    Svc->>Kafka: Publish domain event
    Kafka-->>Sock: Consume event
    Sock-->>UIA: Role-scoped update
    Sock-->>UIB: Sanitized or role-scoped visual update
```

The socket service is display transport only. Persisted backend state and domain events remain source of truth.

---

## 6. Zero-Consent Identity Flow (Seamless SSO)

To optimize internal operations, we implement a **Zero-Consent** OIDC flow for trusted internal dashboard clients.

```mermaid
sequenceDiagram
    participant User
    participant NextJS as Client App
    participant Hydra as Ory Hydra
    participant Kratos as Ory Kratos

    User->>NextJS: Access Dashboard
    NextJS->>Hydra: /oauth2/auth (Start Flow)
    Hydra->>NextJS: /login (Challenge)
    NextJS->>Kratos: /self-service/login
    User->>Kratos: Credentials
    Kratos->>Hydra: Accept Login (API)
    Hydra-->>User: Auto-Consent (Internal Skip)
    Hydra->>NextJS: Callback + Authorization Code
    NextJS->>Hydra: Exchange Code for JWT
    NextJS-->>User: Authenticated Dashboard
```

### Technical Rationale: Why Zero-Consent?
- **User Experience**: Removes redundant "Allow access" screens for internal employees.
- **Security**: The "skip" logic is only applied to **Trusted Client IDs** verified by the Monitor/Auth service.

---
*Technical reference for Runtime Roasters Architectural Decision Patterns.*
