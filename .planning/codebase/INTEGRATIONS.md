# External Integrations

**Analysis Date:** 2025-02-13

## APIs & External Services

**Identity & Auth:**
- Ory Kratos - Identity Management (Registration, Login, User Profiles)
  - SDK/Client: `github.com/ory/client-go`
  - Internal endpoint: `http://kratos:4434/admin`
- Ory Hydra - OAuth2 & OIDC Provider (Token issuance, introspection)
  - SDK/Client: `github.com/ory/hydra-client-go/v2`
  - Public endpoint: `http://localhost:4444`
  - Admin endpoint: `http://localhost:4445`

**API Gateway:**
- KrakenD - Stateless API Gateway
  - Config: `deployments/krakend/krakend.json`
  - Port: `8081`

## Data Storage

**Databases:**
- PostgreSQL 16
  - Connection: `DATABASE_URL` (e.g., `postgres://user:password@postgres:5432/dbname`)
  - Client: `gorm.io/gorm` wrapped in `src/pkg/database/postgres.go`

**Messaging:**
- Apache Kafka 7.6
  - Brokers: `KAFKA_BROKERS` (e.g., `localhost:9094`)
  - Client: `github.com/segmentio/kafka-go` wrapped in `src/pkg/kafka/`
  - Usage: Saga coordination, AuthZ policy sync, token revocation.

**Caching & Transient Storage:**
- Valkey 7.2 (Redis-compatible)
  - Connection: `VALKEY_ADDR` (e.g., `localhost:6379`)
  - Client: `github.com/redis/go-redis/v9` wrapped in `src/pkg/valkey/`
  - Usage: Session storage, Token Blacklist, Idempotency, GPS Tracking.

## Authentication & Identity

**Auth Provider:**
- Custom Identity Stack (Ory Hydra + Kratos)
  - Implementation: OIDC flow for users, OAuth2 Client Credentials for M2M.
  - Proxy: Nginx-based `identity` service protecting Hydra's JWKS.

## Monitoring & Observability

**Distributed Tracing:**
- OpenTelemetry (OTEL)
  - Collector: `otel-collector` at `localhost:4317` (gRPC)
  - Exporter: `otlptracegrpc`

**Logs:**
- Structured Logging with Zap (`src/pkg/logger/`)
- Injected TraceIDs/SpanIDs for log correlation.

## CI/CD & Deployment

**Hosting:**
- Containerized (Docker Compose for local dev, potentially Kubernetes for production).

**CI Pipeline:**
- Not detected in root (checking for `.github/workflows` or similar).

## Environment Configuration

**Required env vars:**
- `DATABASE_URL`: Postgres connection string.
- `VALKEY_ADDR`: Valkey server address.
- `KAFKA_BROKERS`: List of Kafka brokers.
- `INTERNAL_SECRET`: Shared secret for internal service communication.
- `KRATOS_ADMIN_URL`: URL to Kratos Admin API.

**Secrets location:**
- `.env` files (not committed).
- Docker Compose environment variables.

## Webhooks & Callbacks

**Incoming:**
- Hydra Consent/Login endpoints: `http://localhost:3000/consent`, `http://localhost:3000/login`

**Outgoing:**
- None detected (Kafka used for internal events instead).

---

*Integration audit: 2025-02-13*
