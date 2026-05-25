# Architecture

The application stack for **Runtime Roasters** is selected and implemented as a **Go-based Microservices Ecosystem** with **Apache Kafka** and **OpenTelemetry** integration.

For the definitive technical blueprint, service catalog, and infrastructure map, see:
👉 **[Master System Architecture Specification](./product/system-architecture.md)**

This document defines the generic architecture discovery rules, boundary rules, and dependency patterns that agents must follow when extending the system.

## Discovery Before Shape

Before proposing implementation shape, identify:

- **Product surfaces:** React/Next.js frontend, Go gRPC/REST APIs, Kafka consumers/producers.
- **Runtime stack:** Go 1.22+, Postgres (GORM), Kafka, Valkey, SigNoz (OTel).
- **Core domains:** Farm, Batch, Order, Warehouse, Logistics, Trace, Audit.
- **Boundary inputs:** gRPC/REST requests, Kafka events, environment variables, DB rows.
- **Validation ladder:** `Makefile` commands for linting, testing, and service simulation.

Record new stack choices or architectural changes in `docs/decisions/`.

## Default Layering (Clean Architecture)

Services follow a Hexagonal / Clean Architecture pattern:

```text
domain (Entities, Value Objects, Repository Interfaces)
  <- application (Use Cases, Saga Handlers)
      <- infrastructure (DB Adapters, Kafka Producer/Consumer, Valkey)
          <- interface (gRPC Server, Gin HTTP Handlers)
              <- app surfaces (KrakenD Gateway, Next.js Client)
```

## Dependency Rule

Inner layers must not depend on outer layers.

| Layer | May depend on | Must not depend on |
| --- | --- | --- |
| domain | nothing project-external except tiny pure utilities | framework, database, UI, provider, process/env |
| application | domain | framework, UI, provider, database concrete clients |
| infrastructure | domain, application | interface controllers or UI |
| interface | all backend layers | UI state or platform shell assumptions |
| app surfaces | API contracts and app-facing clients | domain internals directly |

## Parse-First Boundary Rule

Unknown data must be parsed at boundaries before it enters inner code.

Boundaries include:

- HTTP request bodies, params, and query strings.
- Session payloads and identity claims.
- Environment variables.
- Database rows returned from external clients.
- Platform shell payloads.
- Deep links, tokens, and signed URLs.
- Provider webhooks, events, and async payloads.

Target flow:

```text
unknown input
  -> parser
  -> typed DTO or command
  -> application use case
  -> domain object/value object
```

Inner layers should work with meaningful product types such as `UserId`,
`AccountId`, `WorkspaceId`, `Role`, `DateRange`, or domain-specific IDs,
rather than repeatedly validating raw strings.

## Command/Query Boundary

If the product has both reads and writes, keep command/query separation clear at
the code level even when the storage layer is simple:

- Commands mutate state and own audit side effects.
- Queries read state and format for consumers.
- Shared domain rules live in domain/application, not controllers.

## Observability Contract

The future server should emit one canonical JSON log line per request with:

- timestamp
- level
- request_id
- user_id when known
- action
- duration_ms
- status_code
- message

Audit logs are product records. Application logs are operational records. Do not
use one as a substitute for the other.
