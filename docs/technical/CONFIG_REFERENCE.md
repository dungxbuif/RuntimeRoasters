# ⚙️ Technical Reference: Configuration & Environment

Runtime Roasters follows the **12-Factor App** methodology for configuration. The system is designed to be **Fail-Fast**: if a required environment variable is missing or invalid, the service will panic immediately during the bootstrap phase.

---

## 1. Core Configuration Mechanism (`pkg/config`)

All services use a centralized configuration loader located in `src/pkg/config`. This loader leverages **Viper** and **Reflect** to provide:
- **Type-safe mapping**: Environment variables are mapped directly to Go structs.
- **Fail-fast validation**: Every field in the `Config` struct is validated. String fields cannot be empty, and integers cannot be zero.
- **Multi-path loading**: Services search for `.env` files in parent directories (up to 4 levels) to support both local development and Docker environments.

### The "No Fallback" Rule
We intentionally avoid setting default values in the code. This ensures that the environment (Docker Compose, K8s, or Local) is the absolute **Source of Truth**. If a variable like `DATABASE_URL` is missing, the app will not start.

---

## 2. Global Environment Variables (`.env`)

These variables are typically shared across the entire infrastructure.

| Variable | Description | Example |
| :--- | :--- | :--- |
| `APP_ENV` | Environment name (dev, staging, prod) | `dev` |
| `POSTGRES_PORT` | Port for the shared Postgres instance | `54321` |
| `VALKEY_ADDR` | Host and port for Valkey/Redis | `localhost:6379` |
| `KAFKA_BROKERS` | List of Kafka brokers (comma-separated) | `localhost:9094` |
| `INTERNAL_SECRET` | Shared secret for internal service auth | `dev-secret-xxxx` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OpenTelemetry collector endpoint | `localhost:4317` |
| `NEXT_PUBLIC_SIGNOZ_URL` | Public SigNoz UI URL for demo/ops links | `http://localhost:3301` |

---

## 3. Service-Specific Variables

Each service extends the `BaseConfig` with its own specific needs (mostly Kafka topics and third-party keys).

### Example: Retail Service
- `KAFKA_ORDER_CREATED_TOPIC`: Topic for publishing new orders.
- `AUTH_SERVICE_ADDR`: gRPC address of the Auth service for policy checks.
- `JWKS_URL`: URL to fetch Hydra's public keys for JWT validation.

---

## 4. Development Startup Sequence

To ensure all dependencies are ready, services should be started in this order:

1.  **Infrastructure (`task infra`)**: Starts Postgres, Kafka, Valkey, Kratos, Hydra, SigNoz, ClickHouse, and the OTel Collector.
2.  **Seeding (`task seed`)**: Provisions the Admin user and OAuth2 clients.
3.  **Auth Service**: Must be up first as it provides Casbin policies to others.
4.  **Core Business Services**: Farm, Warehouse, Retail, etc.
5.  **Intelligence Services**: Trace, Audit (these are purely consumers).
6.  **Frontend (`task fe`)**: The final entry point.

---
*Technical reference for Runtime Roasters System Administration.*
