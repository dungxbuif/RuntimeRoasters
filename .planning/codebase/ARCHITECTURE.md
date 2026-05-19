<!-- refreshed: 2025-05-15 -->
# Architecture

**Analysis Date:** 2025-05-15

## System Overview

```text
┌─────────────────────────────────────────────────────────────┐
│                      Edge & Client Layer                    │
│      Next.js Client App / KrakenD API Gateway               │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Identity & Security                      │
│        Ory Kratos (Identity) / Ory Hydra (OAuth2)           │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Core Business Services                   │
│  `src/apps/farm-service`     `src/apps/retail-service`      │
│  `src/apps/warehouse-service` `src/apps/logistics-service`   │
└────────┬─────────────────┬──────────────────┬───────────────┘
         │                 │                  │
         ▼                 ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                    Support & Intelligence                   │
│  `src/apps/trace-service` (CQRS)                            │
│  `src/apps/audit-service` (Immutable Logs)                  │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  Data & Messaging Backbone                  │
│  Kafka (Events) / Postgres (ACID) / Elasticsearch (Search)  │
│  Cassandra (Audit) / Valkey (Cache)                         │
└─────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility | File |
|-----------|----------------|------|
| **API Gateway** | Entry point, JWT validation, rate limiting, request muxing. | `deployments/krakend/krakend.json` |
| **Auth Service** | Casbin policy authority and user proxy. | `src/apps/auth-service` |
| **Farm Service** | Manages plantations, coffee batches, and harvest records. | `src/apps/farm-service` |
| **Retail Service** | Handles the Order SAGA, customer profiles, and storefront logic. | `src/apps/retail-service` |
| **Warehouse Service** | Inventory management and real-time stock reservation. | `src/apps/warehouse-service` |
| **Logistics Service** | Real-time driver GPS tracking and route calculation. | `src/apps/logistics-service` |
| **Trace Service** | The "Read Model" for supply chain traceability using CQRS. | `src/apps/trace-service` |
| **Audit Service** | Immutable append-only log of significant system events. | `src/apps/audit-service` |

## Pattern Overview

**Overall:** Microservices + Clean Architecture (Hexagonal)

**Key Characteristics:**
- **Event-Driven Core:** Asynchronous communication via Apache Kafka for eventual consistency.
- **Transactional Outbox:** Ensures data consistency between database updates and Kafka publishing.
- **Polyglot Persistence:** Different databases optimized for specific workloads (Postgres, Cassandra, ES, Valkey).

## Layers

**Service Layer (Clean Architecture):**
- **Purpose:** Decouples business logic from infrastructure and external interfaces.
- **Location:** `src/apps/[service]/internal/`
- **Contains:** 
  - `domain/`: Business entities and repository interfaces.
  - `usecase/`: Application logic and business rules.
  - `app/`: Dependency injection and application bootstrapping.
- **Depends on:** Nothing (Domain), Domain (Usecase).
- **Used by:** External transport layers (gRPC/HTTP).

## Data Flow

### Primary Request Path (Synchronous)

1. **Client Request:** Next.js app sends REST request to KrakenD.
2. **Gateway:** KrakenD validates JWT via Hydra and checks RBAC. (`deployments/krakend/krakend.json`)
3. **Muxing:** KrakenD forwards request to backend service via gRPC.
4. **Service App:** `init.go` handles gRPC server and routes to `usecase`. (`src/apps/[service]/internal/app/init.go`)
5. **Usecase:** Business logic executes, interacting with `domain` models. (`src/apps/[service]/internal/usecase/service.go`)
6. **Persistence:** GORM saves state to Postgres. (`src/apps/[service]/internal/domain/models.go`)

### Event-Driven Path (Asynchronous)

1. **Transaction:** Service updates DB and writes to `outbox_events` in one transaction.
2. **Relay:** A background worker (or DB trigger) picks up outbox events and publishes to Kafka.
3. **Consumption:** Downstream services (e.g., `trace-service`) consume from Kafka.
4. **Projection:** `trace-service` updates its read model in Elasticsearch. (`src/apps/trace-service/internal/search/elasticsearch.go`)

**State Management:**
- Stateless services; all persistent state is in databases or Kafka. Session data managed by Ory Kratos.

## Key Abstractions

**Transactional Outbox:**
- Purpose: Guarantees "at-least-once" delivery of events.
- Examples: `src/pkg/events/contracts.go`
- Pattern: Local table + background publisher.

**Two-Gate Authorization:**
- Purpose: Secure the system at the edge and at the service level.
- Examples: `deployments/krakend/krakend.json` (Gate 1), `src/pkg/base/casbin/` (Gate 2).
- Pattern: KrakenD (RBAC) + Casbin Enforcer (Fine-grained AuthZ).

## Entry Points

**gRPC Server:**
- Location: `src/apps/[service]/internal/app/init.go`
- Triggers: Inbound gRPC calls from Gateway or other services.
- Responsibilities: Server lifecycle, interceptors (telemetry, auth).

**Main CMD:**
- Location: `src/apps/[service]/cmd/main.go`
- Triggers: Service startup.
- Responsibilities: Config loading, logging setup, app initialization.

## Architectural Constraints

- **Single Module:** All code resides in one Go module `github.com/dungxbuif/RuntimeRoasters` located in `src/`.
- **Statelessness:** Services must not store session state locally.
- **Observability:** Every service must use OTel for tracing. (`src/pkg/telemetry/`)

## Anti-Patterns

### Circular Dependencies

**What happens:** Attempting to import `service A` in `service B` and vice versa.
**Why it's wrong:** Breaks Go compilation and creates tight coupling.
**Do this instead:** Use Kafka events for decoupling or move shared logic to `src/pkg/`.

## Error Handling

**Strategy:** Centralized error handling using a custom error package.

**Patterns:**
- **Domain Errors:** Defined in `src/pkg/errs/`.
- **gRPC Interceptors:** Translate internal errors to gRPC status codes. (`src/pkg/base/auth/transport/grpc/interceptor.go`)

## Cross-Cutting Concerns

**Logging:** Structured JSON logging using Zap. (`src/pkg/logger/`)
**Validation:** Request validation using `go-playground/validator`.
**Authentication:** JWT-based identity propagation. (`src/pkg/base/identity/`)

---

*Architecture analysis: 2025-05-15*
