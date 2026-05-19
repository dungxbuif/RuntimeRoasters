# 📘 Runtime Roasters Developer Guide

This guide provides a comprehensive technical breakdown of the **Runtime Roasters** platform. It serves as the primary technical reference for developers, covering architectural patterns, coding standards, and operational details.

---

## 🏗️ 1. Architectural Philosophy

Runtime Roasters is built as a highly available, event-driven microservices ecosystem. The system is designed so that scaling from 1 to N instances is purely an infrastructure operation with zero logic changes.

### Core Principles
- **Clean Architecture:** Strict separation of concerns (Domain ➔ UseCase ➔ Infrastructure ➔ Delivery).
- **Abstraction First:** All infrastructure (DB, Queue, Cache) is abstracted via Interfaces.
- **Calculated Consistency:** **Transactional Outbox** ensures critical events are never lost; **Inbox Pattern** ensures idempotency.
- **Dual Idempotency:** Two layers of protection: `Idempotency-Key` (HTTP) and `Transactional Inbox` (Kafka).
- **Observability by Design:** Distributed tracing (W3C) is propagated across Gateway, gRPC, and Kafka.
- **Zero-Trust Security:** Internal model with decentralized authorization (Casbin).

---

## 🌊 2. System Data Flows

### 🔐 Identity & SSO
1.  **Identity:** Managed by **Ory Kratos**. Handles user profiles, registration, and login.
2.  **Authentication:** Managed by **Ory Hydra**. Issues OIDC/OAuth2 tokens.
3.  **Gatekeeper:** KrakenD API Gateway validates JWT signatures and `scope` claims.
4.  **Authorization:** Each service runs a **Casbin Enforcer**. Policies are synced from `auth-service` via Kafka (`auth.policy.changed`) and gRPC snapshots.

### 📦 Order SAGA (Choreography)
1.  **Retail Service:** Creates Order (Pending) ➔ Publishes `OrderCreated` event.
2.  **Warehouse Service:** Receives `OrderCreated` ➔ Reserves Stock ➔ Publishes `StockReserved` (or `StockReservationFailed`).
3.  **Payment Service:** Receives `StockReserved` ➔ Processes Payment ➔ Publishes `PaymentCompleted`.
4.  **Retail Service:** Receives `PaymentCompleted` ➔ Updates Order to `Confirmed`.
5.  *Compensation:* If any step fails, services publish compensation events to roll back previous actions.

---

## 🛠️ 3. Development Workflow

### Project Structure (src/apps/)
Each service follows this layout:
- `cmd/`: Process entry point (Main).
- `internal/`:
    - `domain/`: Business entities and interfaces.
    - `usecase/`: Business logic.
    - `infrastructure/`: Implementations (DB, Kafka, Valkey).
    - `delivery/`: Transport handlers (gRPC, Gin).
- `config/`: Configuration mapping.

### 📦 3.1. Internal Framework (`src/pkg/`)
All services leverage a shared core framework to ensure consistency and reduce boilerplate:
- **`base`**: The service orchestrator. Handles gRPC/HTTP server lifecycle, graceful shutdown, and cross-cutting auth/security middlewares.
- **`database`**: GORM-based abstraction for PostgreSQL with built-in support for Transactional Outbox.
- **`kafka`**: High-level wrapper for Segmentio/Kafka-Go, implementing the Outbox/Inbox patterns.
- **`valkey`**: Client for Valkey (Redis alternative) used for distributed locking and real-time tracking.
- **`errs`**: Centralized error handling using RFC 9457 (Problem Details).
- **`telemetry`**: OpenTelemetry integration for tracing and metrics.
- **`logger`**: Structured JSON logging via Uber-Zap with Trace-ID propagation.

### 🧪 3.2. Testing Strategy
We employ a multi-layered testing strategy to ensure system reliability:
- **Unit Tests**: Focus on `domain` and `usecase` logic. Mocking is used for all external dependencies (DB, Kafka).
- **Integration Tests**: Verify the interaction between `infrastructure` implementations and real (Docker-based) dependencies.
- **E2E Tests**: (Future) Simulated supply chain runs across the full stack.
- **Contract Testing**: Protobuf definitions serve as the strict contract between services.

### Key Tools
- **Go Workspaces (`go.work`):** We use a Go workspace to manage multiple internal modules (e.g., `src` containing apps and pkg). This allows for seamless cross-module development and dependency resolution.
- **Taskfile:** Automated tasks for infra, proto generation, and seeding.
- **Air:** Hot-reload for Go development.
- **Buf:** Modern Protobuf management and code generation.
- **KrakenD:** Powerful API gateway for orchestration and security.

### Coding Standards
- **DRY:** Use shared logic in `src/pkg/` (logger, database, telemetry).
- **Type-Safety:** All internal communication must use gRPC/Protobuf.
- **Explicit over Implicit:** Manual dependency injection is preferred over magic containers.

---

## 🧭 4. Operational Reference

### Local Service URLs
- **Client App:** `http://localhost:3000`
- **KrakenD (Gateway):** `http://localhost:8081`
- **Kafka UI:** `http://localhost:8090`
- **Kibana (ES):** `http://localhost:5601`
- **Kratos Public:** `http://localhost:4433`
- **Hydra Public:** `http://localhost:4444`

### Kafka Topics
- `auth.policy.changed`: Broadcasts Casbin policy updates.
- `auth.user.events`: User lifecycle events.
- `farm.harvest.events`: Harvest records from plantations.
- `order.saga.events`: Main topic for order choreography.

---

## 🗺️ 5. Detailed Guides & Decision Records

For a deeper dive into specific system behaviors and the reasoning behind them:
- 🏛️ **[Master System Architecture Specification](./architecture/SYSTEM_ARCHITECTURE.md)**: The "Big Picture" blueprint.
- 🌊 **[System Architecture Flows](./architecture/FLOWS.md)**: Identity, Sagas, and CQRS patterns.
- 🛠️ **[Technical Knowledge Base](./technical/README.md)**: Exhaustive details on Config, Data Models, and API Contracts.
- 🎨 **[Frontend Development Guide](./technical/FRONTEND.md)**: Next.js patterns and UI/UX standards.
- 📦 **[Domain & Business Logic](./domain/02-PROCESSING_INVENTORY.md)**: Rules for coffee batch lifecycles and IDs.
- 📋 **[Product Requirements](./business/product-requirements.md)**: The business vision and feature specs.
- 📅 **[Sprint Roadmap](./business/sprint-planning.md)**: Development phases and task breakdown.
- 📜 **[Architecture Decisions (ADRs)](./architecture/adrs/)**: History of critical technical choices.

---
*This document is the "Living Constitution" of Runtime Roasters. Keep it updated as the architecture evolves.*
