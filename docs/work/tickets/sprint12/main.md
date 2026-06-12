# Sprint 12: Production & K8s Readiness Audit

**Epic Goal:** Audit and harden the platform for Kubernetes deployment. Focus on extracting hardcoded configurations and ensuring the system is "environment-agnostic".

---

## 🔍 Hardcoded Configuration Audit
The following areas must be extracted into environment variables before deployment:

### 1. Backend Service URLs (gRPC/HTTP)
- [x] `AUTH_SERVICE_ADDR`: config struct reads from ENV; `localhost:50052` fallback removed.
- [x] `ELASTICSEARCH_URL`: config struct reads from ENV; `localhost:9200` fallback removed from trace-service.
- [x] `KAFKA_BROKERS`: config struct reads from ENV; no fallback was present in Go code.
- [x] `CASSANDRA_ADDR`: `127.0.0.1:9042` fallback removed from audit-service `cassandra/store.go`.

### 2. Identity Provider Config (OIDC)
- [x] `KRATOS_PUBLIC_URL`: configurable via ENV; no hardcode found in Go.
- [x] `HYDRA_ADMIN_URL`: `localhost:4445` fallback removed from auth-service `init.go`.
- [ ] `EXPECTED_ISSUER`: must match the public domain (e.g., `https://runtimeroasters.com/`). *(`.env` files correct for dev; K8s override needed at deploy time)*

### 3. Client App (Frontend)
- [ ] Unified Domain Path: UI must call `/api/*` instead of direct service ports.
- [ ] `NEXT_PUBLIC_GATEWAY_URL`: point to `https://<domain>/api`.

### 4. Additional Go Fallbacks Removed
- [x] `auth-service/grpc/handler.go`: removed `GetSystemStatus` and `SeedData` localhost fallback maps (7 services × 2 methods).
- [x] `socket-service/init.go`: removed `ValkeyAddr` (localhost:6379) and `InternalHost` (127.0.0.1) fallbacks.
- [x] `retail-service/init.go`: removed `InternalHost` (127.0.0.1) fallback.
- [x] `trace-service/init.go`: removed `ElasticsearchURL` and `RetailServiceURL` fallbacks.

---

## ✅ K8s Readiness Checklist

### 1. Application Hardening
- [x] **Probes**: Implement `/health` or `/ready` endpoints in all Go services. All services now expose `/health/live` and `/health/ready` (via `RegisterReadiness`).
- [x] **Shutdown**: Implement Graceful Shutdown (handling `SIGTERM`) to finish active requests/Kafka messages.
- [x] **Logging**: Ensure all logs are sent to `stdout` in JSON format (Zap production encoder).
- [x] **Non-Root**: Update Dockerfiles to use a non-root user (e.g., `nobody` or `nonroot`).

### 2. Networking (Unified Domain)
- [x] **Ingress Mapping**:
    - `/` ➔ `client-app`
    - `/api/*` ➔ `krakend` (API Gateway)
    - `/.well-known/*` ➔ `identity-proxy`
- [x] **CORS**: Standardize allowed origins to the production domain.

### 3. Database & State
- [x] **Migrations**: Move migration logic to K8s `InitContainers` or a dedicated Migration Job.
- [x] **Secrets**: No passwords in `.env` files. Mount secrets via K8s Secret volumes or CSI drivers.
- [x] **Idempotency**: Verify Kafka Consumers handle duplicate messages (Redelivery) safely.

### 4. Continuous Integration
- [x] **Build**: Multi-stage Docker builds to keep image size small.
- [x] **Linter**: Add `golangci-lint` to CI pipeline.
- [x] **Security**: Scan images for vulnerabilities (Trivy/Snyk).

---
*Status: Active Audit Phase.*
