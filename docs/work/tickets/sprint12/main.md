# Sprint 12: Production & K8s Readiness Audit

**Epic Goal:** Audit and harden the platform for Kubernetes deployment. Focus on extracting hardcoded configurations and ensuring the system is "environment-agnostic".

---

## 🔍 Hardcoded Configuration Audit
The following areas must be extracted into environment variables before deployment:

### 1. Backend Service URLs (gRPC/HTTP)
- [ ] `AUTH_SERVICE_ADDR`: currently often hardcoded as `localhost:50052`.
- [ ] `ELASTICSEARCH_URL`: currently hardcoded as `http://localhost:9200`.
- [ ] `KAFKA_BROKERS`: currently hardcoded as `localhost:9094`.
- [ ] `CASSANDRA_ADDR`: currently hardcoded as `127.0.0.1:9042`.

### 2. Identity Provider Config (OIDC)
- [ ] `KRATOS_PUBLIC_URL`: must be configurable via ENV for unified domain routing.
- [ ] `HYDRA_ADMIN_URL`: often hardcoded as `http://localhost:4445`.
- [ ] `EXPECTED_ISSUER`: must match the public domain (e.g., `https://runtimeroasters.com/`).

### 3. Client App (Frontend)
- [ ] Unified Domain Path: UI must call `/api/*` instead of direct service ports.
- [ ] `NEXT_PUBLIC_GATEWAY_URL`: point to `https://<domain>/api`.

---

## ✅ K8s Readiness Checklist

### 1. Application Hardening
- [ ] **Probes**: Implement `/health` or `/ready` endpoints in all Go services.
- [ ] **Shutdown**: Implement Graceful Shutdown (handling `SIGTERM`) to finish active requests/Kafka messages.
- [ ] **Logging**: Ensure all logs are sent to `stdout` in JSON format (Zap production encoder).
- [ ] **Non-Root**: Update Dockerfiles to use a non-root user (e.g., `nobody` or `nonroot`).

### 2. Networking (Unified Domain)
- [ ] **Ingress Mapping**:
    - `/` ➔ `client-app`
    - `/api/*` ➔ `krakend` (API Gateway)
    - `/.well-known/*` ➔ `identity-proxy`
- [ ] **CORS**: Standardize allowed origins to the production domain.

### 3. Database & State
- [ ] **Migrations**: Move migration logic to K8s `InitContainers` or a dedicated Migration Job.
- [ ] **Secrets**: No passwords in `.env` files. Mount secrets via K8s Secret volumes or CSI drivers.
- [ ] **Idempotency**: Verify Kafka Consumers handle duplicate messages (Redelivery) safely.

### 4. Continuous Integration
- [ ] **Build**: Multi-stage Docker builds to keep image size small.
- [ ] **Security**: Scan images for vulnerabilities (Trivy/Snyk).

---
*Status: Active Audit Phase.*
