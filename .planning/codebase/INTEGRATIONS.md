# External Integrations

**Analysis Date:** 2025-02-14

## APIs & External Services

**Payment Gateways:**
- Stripe - Used for credit card payments (simulated).
  - SDK/Client: `provider.simulatedGateway`
  - Auth: `STRIPE_WEBHOOK_SECRET`
- VNPay - Used for local Vietnamese payments (simulated).
  - SDK/Client: `provider.simulatedGateway`
  - Auth: `VNPAY_SECRET`

**Identity & Auth:**
- Ory Kratos - Identity management (registration, login, profile).
  - Client: `ory/kratos-client-go`
  - Auth: `KRATOS_ADMIN_URL`, `KRATOS_PUBLIC_URL`
- Ory Hydra - OAuth2/OIDC provider (token issuance).
  - Client: `ory/hydra-client-go`
  - Auth: `HYDRA_ADMIN_URL`, `HYDRA_PUBLIC_URL`

## Data Storage

**Databases:**
- PostgreSQL 16
  - Connection: `DATABASE_DSN`
  - Client: GORM
- Valkey 7.2 (Redis compatible)
  - Connection: `VALKEY_ADDR`
  - Client: `redis/go-redis`
  - Usage: Session storage, GPS tracking, Caching.
- Cassandra 4.1
  - Usage: Immutable audit logs.
- Elasticsearch 8.15.3
  - Usage: CQRS read model and search backend for Trace Service.

**File Storage:**
- Local filesystem only (noted for development).

**Caching:**
- Valkey (Redis) used for general caching and session management.

## Authentication & Identity

**Auth Provider:**
- Custom (RuntimeRoasters Auth Service) built on Ory Stack.
  - Implementation: `src/apps/auth-service/`
  - Uses Ory Kratos for identities and Ory Hydra for OAuth2.

## Monitoring & Observability

**Error Tracking:**
- OpenTelemetry (OTel) integration.

**Logs:**
- Zap Logger (`go.uber.org/zap`).
- Structured logging to stdout.

## CI/CD & Deployment

**Hosting:**
- Docker Compose (Development).

**CI Pipeline:**
- Not detected (GitHub Actions expected but not verified in root).

## Environment Configuration

**Required env vars:**
- `DATABASE_DSN`
- `KAFKA_BROKERS`
- `VALKEY_ADDR`
- `KRATOS_ADMIN_URL`
- `HYDRA_ADMIN_URL`
- `INTERNAL_SECRET`

**Secrets location:**
- `.env` files and environment variables.

## Webhooks & Callbacks

**Incoming:**
- `/v1/webhooks/stripe`: Stripe payment status updates (`src/apps/payment-service/internal/app/app.go`).
- `/v1/webhooks/vnpay`: VNPay payment status updates.

**Outgoing:**
- None detected.

## Message Broker (Kafka)

**Kafka Topics:**
- `retail.order.created`: Published when a new order is placed.
- `payment.intent.created`: Published when a payment session is initiated.
- `payment.completed`: Published when payment is successful.
- `payment.failed`: Published when payment fails.
- `payment.refunded`: Published when a refund is processed.
- `warehouse.stock.reserved`: Published when inventory is successfully reserved.
- `warehouse.stock.reservation_failed`: Published when inventory reservation fails.
- `warehouse.stock.updated`: Published on stock changes.
- `logistics.shipment.assigned`: Published when a shipment is assigned to a driver.
- `logistics.shipment.delivered`: Published when a shipment reaches its destination.
- `logistics.gps.updated`: Real-time driver location updates.
- `authz.policy.sync`: Internal topic for Casbin policy synchronization.

## Internal gRPC Services

**AuthService:**
- Exposed by: `auth-service`
- Used by: `farm-service`, `logistics-service`, `retail-service`, `warehouse-service` (via `pkg/base/security`).

**FarmService:**
- Exposed by: `farm-service`
- Used by: `client-app` (via gRPC-Gateway).

**DemoService:**
- Exposed by: `demo-service`
- Used by: Internal scripts/testing.

---

*Integration audit: 2025-02-14*
