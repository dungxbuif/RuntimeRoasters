# Technology Stack

**Analysis Date:** 2025-02-14

## Languages

**Primary:**
- Go 1.25.0 - Backend microservices in `src/apps/` and shared packages in `src/pkg/`.

**Secondary:**
- TypeScript/JavaScript - Frontend client application in `src/apps/client-app/`.
- Protobuf - API contracts and service definitions in `api/`.
- SQL - Database migrations and initialization in `deployments/init-db.sql`.

## Runtime

**Environment:**
- Docker / Docker Compose - Primary development and deployment environment.

**Package Manager:**
- Go Modules - Backend dependency management (`src/go.mod`).
- NPM/PNPM - Frontend dependency management (inferred from `src/apps/client-app/`).
- Lockfile: `src/go.sum` and `go.work.sum` are present.

## Frameworks

**Core:**
- Gin-Gonic v1.12.0 - HTTP web framework for REST APIs.
- gRPC v1.80.0 - High-performance RPC framework for service-to-service communication.
- GORM v1.25.x - ORM for database interactions.
- Next.js / React - Frontend framework for `client-app`.
- Casbin v3.10.0 - Authorization library (RBAC/ABAC).

**Testing:**
- testify v1.11.1 - Assertion and mocking library for Go.
- Go standard testing package.

**Build/Dev:**
- Taskfile - Task runner (`Taskfile.yml`).
- Buf - Tooling for Protobuf/gRPC management (`api/buf.yaml`).
- Makefile - Build automation in `src/Makefile`.

## Key Dependencies

**Critical:**
- `segmentio/kafka-go` v0.4.51 - Kafka client for Go.
- `redis/go-redis` v9.18.0 - Redis/Valkey client.
- `ory/kratos-client-go` & `ory/hydra-client-go` - Identity and OAuth2 management.
- `spf13/viper` v1.21.0 - Configuration management.

**Infrastructure:**
- OpenTelemetry (OTel) - Observability and tracing.
- KrakenD - API Gateway.

## Configuration

**Environment:**
- Environment variables via `.env` files (e.g., `src/apps/auth-service/.env`).
- Viper for loading and parsing configurations.

**Build:**
- `src/Makefile`
- `Taskfile.yml`
- `api/buf.gen.yaml` for Protobuf generation.

## Platform Requirements

**Development:**
- Docker & Docker Compose
- Go 1.25.0+
- Node.js (for frontend)

**Production:**
- Kubernetes (implied by microservices architecture and OTel setup).

---

*Stack analysis: 2025-02-14*
