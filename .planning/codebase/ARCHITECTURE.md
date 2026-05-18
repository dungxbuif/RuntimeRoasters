<!-- refreshed: 2025-02-13 -->
# Architecture

**Analysis Date:** 2025-02-13

## System Overview

```text
┌─────────────────────────────────────────────────────────────┐
│                      Client App / Frontend                   │
│         `src/apps/client-app`                                │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway (KrakenD)                   │
│         `deployments/krakend`                                │
├──────────────────┬──────────────────┬───────────────────────┤
│   auth-service   │  farm-service    │ warehouse-service     │
│ `src/apps/auth-*`│ `src/apps/farm-*`│ `src/apps/warehouse-*`│
└────────┬─────────┴────────┬─────────┴──────────┬────────────┘
         │                  │                     │
         ▼                  ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Shared Infrastructure                     │
│         `src/pkg/` (database, kafka, telemetry, etc.)        │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│  PostgreSQL / Valkey / Kafka / Ory Kratos                    │
│  `deployments/docker-compose.dev.yaml`                       │
└─────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility | File |
|-----------|----------------|------|
| API Gateway | Entry point, routing, JWT validation (Gate 1) | `deployments/krakend/krakend.json` |
| Microservices | Domain logic, data processing, async events | `src/apps/farm-service/cmd/main.go` |
| Auth Service | Centralized role management (Casbin snapshot) | `src/apps/auth-service/cmd/main.go` |
| Client App | UI for user interaction | `src/apps/client-app/src/app/page.tsx` |
| Common Libs | Database, logging, telemetry, kafka clients | `src/pkg/base/app.go` |

## Pattern Overview

**Overall:** Microservices with Clean Architecture

**Key Characteristics:**
- **Calculated Consistency:** Uses Transactional Outbox for critical event publishing.
- **Observability by Design:** Built-in W3C tracing, metrics, and log correlation via OpenTelemetry.
- **Fail-Closed Security:** Three-Gate Authorization (KrakenD -> Interceptor -> Repository scope).
- **Stateless Services:** Designed for horizontal scalability.

## Layers

**Domain Layer:**
- Purpose: Core entities and domain-specific errors. No external dependencies.
- Location: `src/apps/farm-service/internal/domain/`
- Contains: Structs, enums, pure functions.
- Depends on: Nothing outside the project.
- Used by: UseCase, Infrastructure, Delivery.

**UseCase Layer:**
- Purpose: Application business rules, coordinates domain entities and repositories.
- Location: `src/apps/farm-service/internal/usecase/`
- Contains: Interactors, Repository Interfaces.
- Depends on: Domain layer.
- Used by: Delivery layer.

**Infrastructure Layer:**
- Purpose: Database implementations, external API clients, message brokers.
- Location: `src/apps/farm-service/internal/infrastructure/`
- Contains: GORM Repositories, Kafka Publishers, Outbox Relays.
- Depends on: UseCase, Domain, `src/pkg/*`.
- Used by: Main (Composition Root).

**Delivery Layer:**
- Purpose: Exposing the application to external consumers (gRPC, HTTP).
- Location: `src/apps/farm-service/internal/delivery/grpc/`
- Contains: gRPC handlers, DTO mapping.
- Depends on: UseCase layer.
- Used by: Framework routing (`src/pkg/base`).

## Data Flow

### Primary Request Path (gRPC Gateway)

1. Client request hits Gateway (`deployments/krakend`)
2. Request routed to Service Gateway / HTTP fallback (`src/pkg/base/app.go`)
3. gRPC Interceptor validates AuthZ (`src/pkg/base/casbin/transport/grpc/interceptor.go`)
4. Delivery Handler receives request (`src/apps/farm-service/internal/delivery/grpc/farm_handler.go`)
5. UseCase executes business logic (`src/apps/farm-service/internal/usecase/farm_usecase.go`)
6. Repository persists data (`src/apps/farm-service/internal/infrastructure/repository/farm_repository.go`)

### Asynchronous Event Publishing (Transactional Outbox)

1. UseCase saves entity AND `OutboxEvent` in one DB transaction (`src/apps/farm-service/internal/usecase/harvest_usecase.go`).
2. Outbox Relay worker polls/listens for new events (`src/apps/farm-service/internal/infrastructure/event/outbox_relay.go`).
3. Publisher sends event to Kafka (`src/apps/farm-service/internal/infrastructure/event/publisher.go`).
4. Relay marks Outbox event as processed.

## Key Abstractions

**Base App Framework:**
- Purpose: Standardized microservice bootstrapping (Gin, gRPC, OTel, Zap, Health).
- Examples: `src/pkg/base/app.go`
- Pattern: Builder / Composition.

**Auth Enforcer:**
- Purpose: In-memory Casbin enforcer with background sync.
- Examples: `src/pkg/base/casbin/resilient_reader.go`
- Pattern: Polling/Snapshot replication.

**Transactional Outbox/Inbox:**
- Purpose: Reliable messaging without dual-write issues.
- Examples: `src/apps/farm-service/internal/domain/outbox_event.go`
- Pattern: Outbox Pattern.

## Entry Points

**Microservice Bootstrap:**
- Location: `src/apps/farm-service/cmd/main.go` -> `src/apps/farm-service/internal/app/init.go`
- Triggers: Application startup.
- Responsibilities: Manual Dependency Injection, connecting to DB/Kafka, starting base app.

## Architectural Constraints

- **Dependency Rule:** Inner layers (Domain) MUST NOT import outer layers (Infrastructure/Delivery).
- **Internal Communication:** Exclusively gRPC between services. No HTTP internal calls.
- **Manual DI:** No reflection-based DI frameworks (like Uber Dig). Services use manual composition in `init.go`.
- **Three-Gate Authorization:** 
  1. Identity via KrakenD
  2. RPC Method authz via Casbin gRPC Interceptor
  3. Data-level authz via Scoped Repositories (Gorm Scopes)

## Anti-Patterns

### Direct Database Access from Delivery Layer

**What happens:** Delivery handler directly calls DB functions.
**Why it's wrong:** Bypasses business logic, violates Clean Architecture, makes unit testing difficult.
**Do this instead:** Define an interface in `usecase`, implement it in `infrastructure/repository`, and inject it into the handler via the usecase.

### Dual-Write to DB and Kafka

**What happens:** Calling DB save and Kafka publish synchronously in the same function.
**Why it's wrong:** If Kafka is down, the DB transaction succeeds but the event is lost, causing system inconsistency.
**Do this instead:** Use the Transactional Outbox pattern (`src/apps/farm-service/internal/infrastructure/event/outbox_relay.go`).

## Error Handling

**Strategy:** Centralized RFC 9457 (Problem Details).

**Patterns:**
- Errors generated in domain/usecase use standard types (`src/pkg/errs`).
- Delivery layer translates them to gRPC/HTTP status codes via middlewares/interceptors.
- `src/pkg/errs/renderer.go` formats the output.

## Cross-Cutting Concerns

**Logging:** Structured logging using Uber Zap (`src/pkg/logger`). Enriched with Trace IDs.
**Validation:** Domain entities contain `Validate()` methods (`src/apps/farm-service/internal/domain/farm.go`).
**Authentication:** Ory Kratos/Hydra via OAuth2/OIDC. Interceptors extract Identity from context (`src/pkg/base/identity`).

---

*Architecture analysis: 2025-02-13*