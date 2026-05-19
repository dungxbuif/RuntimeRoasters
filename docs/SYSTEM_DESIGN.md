# Runtime Roasters — System Design & Architecture Deep Dive

This document provides a comprehensive technical breakdown of the **Runtime Roasters** platform. It serves as the primary technical reference for the system's distributed patterns, architectural decisions, and data flows.

---

## 🏗️ 1. High-Level Architecture & Blueprint

Runtime Roasters is built as a highly available, event-driven microservices ecosystem. The system is designed so that scaling from 1 to N instances is purely an infrastructure operation with zero logic changes.

### Core Principles (The Constitution)
- **Abstraction First:** All infrastructure (DB, Queue, Cache) is abstracted via Interfaces.
- **Calculated Consistency:** **Transactional Outbox** ensures critical events are never lost; **Inbox Pattern** ensures idempotency.
- **Dual Idempotency:** Two layers of protection: `Idempotency-Key` (HTTP) and `Transactional Inbox` (Kafka).
- **Observability by Design:** Distributed tracing (W3C) is propagated across Gateway, gRPC, and Kafka.
- **Fail-Closed Security:** Internal Zero Trust model with decentralized authorization.

### Infrastructure Overview
```mermaid
graph TB
    subgraph "External World"
        Browser([Browser / Mobile])
    end

    subgraph "Client App"
        CA[client-app :3000]
    end

    subgraph "API Gateway (HA)"
        GW[KrakenD :8080]
    end

    subgraph "Microservices (Stateless)"
        AS[auth-service :8081]
        FS[farm-service :8082]
        WS[warehouse-service :8083]
        RS[retail-service :8084]
        LS[logistics-service :8085]
        TS[trace-service :8086]
        PS[payment-service :8087]
        AUS[audit-service :8088]
    end

    subgraph "Data Persistence (HA)"
        PG[(PostgreSQL :5432)]
        RD[(Valkey/Redis :6379)]
        ES[(Elasticsearch :9200)]
        CS[(Cassandra :9042)]
    end

    subgraph "Message Broker (HA)"
        KF[Kafka Cluster :9092]
    end

    Browser -->|:3000| CA
    CA -->|REST :8080| GW
    GW -->|gRPC+JWT| AS
    GW -->|gRPC+JWT| FS
    GW -->|gRPC+JWT| WS
    GW -->|gRPC+JWT| RS
    GW -->|gRPC+JWT| LS
    
    FS -.->|Outbox/Inbox| KF
    WS -.->|Outbox/Inbox| KF
    RS -.->|Outbox/Inbox| KF
    LS -.->|Outbox/Inbox| KF
    KF -.-> TS
    KF -.-> AUS
```

---

## 🌊 2. System Data Flows

Detailed interaction diagrams and step-by-step logic for core business processes.

*   👉 **[Identity & SSO Flow](./architecture/flows/identity-authentication.md)**: OIDC, Kratos, and Hydra integration.
*   👉 **[System-Wide Data Flows](./architecture/flows/system-data-flow.md)**: Sagas, CQRS, and Transactional Outbox.
*   👉 **[User Management Flow](./architecture/flows/user-management-creation.md)**: Admin-only creation and propagation.

---

## 📜 3. Architectural Decision Records (ADR Log)

A record of the critical technical decisions that shaped the platform.

| ID | Title | Status | Link |
| :-- | :--- | :--- | :--- |
| **0001** | **Clean Architecture** | ✅ Accepted | [View ADR](./architecture/adrs/0001-use-clean-architecture.md) |
| **0002** | **SigNoz Reversion** | 🔴 Reverted | [View ADR](./architecture/adrs/0002-use-signoz.md) |
| **0003** | **Two-Gate AuthZ** | ✅ Accepted | [View ADR](./architecture/adrs/0003-two-gate-authz-casbin.md) |
| **0004** | **GORM Migration** | ✅ Accepted | [View ADR](./architecture/adrs/0004-migrate-sqlx-to-gorm.md) |

---

## 🛠️ 4. Technical Standards & Patterns

### Clean Architecture Layers
- **Domain:** Pure entities + domain errors. Zero external imports.
- **UseCase:** Business logic + Interface definitions.
- **Infrastructure:** Implementations (GORM, Kafka, Valkey).
- **Delivery:** Transport handlers (gRPC, Gin).

### Security Model (Zero Trust)
1. **API Gateway:** KrakenD verifies JWT signatures and `scope` claims.
2. **Service Level:** Casbin Enforcer checks `role` vs `resource:action`.
3. **Data Level:** Repositories enforce ownership scopes (e.g., `WHERE store_id = ?`).

### Observability
- **Tracing:** W3C Trace-ID propagation across HTTP -> gRPC -> Kafka.
- **Logging:** Zap Logger with automatic `trace_id` injection.
- **Metrics:** Prometheus exporters integrated via `pkg/telemetry`.

---
*Last Updated: May 2026*
