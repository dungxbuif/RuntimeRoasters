# Runtime Roasters Release Checklist

This checklist is the release gate for infrastructure, security, validation, and deployment readiness.

---

## 🏗 PHASE 1: Infrastructure & Config
Ensure a solid foundation before running Application logic.

### 1.1. Databases & Storage
- [ ] **Postgres**: Full migrations have been run (`identity_db`, `demo_db`, ...).
- [ ] **Valkey/Redis**: Configured for sessions, socket state, and any token blacklist usage.
- [ ] **Kafka**: Necessary business topics exist and Kafka remains KRaft-compatible.
- [ ] **Healthchecks**: Every infrastructure service must have a `healthcheck` configuration in Docker/K8s.

### 1.2. Identity & Auth (Core Security)
- [ ] `EXPECTED_ISSUER`: Matches the `iss` claim (e.g., `http://localhost:4444/`).
- [ ] `JWKS_URL`: Points correctly to the internal endpoint (`http://identity:4434/...`).
- [ ] `INTERNAL_SECRET`: Changed from the development value and synchronized across internal callers.
- [ ] **OIDC Client Secret**: Configured with a production-grade random value.
- [ ] **Socket internal API keys**: Rotated from demo values.
- [ ] **Webhook secrets**: Provider webhook keys are environment-specific.

### 1.3. Environment Variables
- [ ] **Centralized Config**: Every sensitive environment variable (`*_SECRET`, `*_PASSWORD`) must be managed via Secret Manager or `.env.local`.
- [ ] **Frontend Env**: `NEXT_PUBLIC_*` has been checked for correctness when building for production.

---

## 🔒 PHASE 2: Security & Resilience
Ensure the system is safe and self-healing.

### 2.1. Code Audit
- [ ] **Fail-Fast & Retry**: Code fetching startup data (like JWKS) has Exponential Backoff Retry.
- [ ] **Error Handling**: Uses RFC 9457 (Problem Details) standard, no internal error information leaks.
- [ ] **Input Validation**: Every API endpoint has validation for request body/params.
- [ ] **Transactional Outbox**: Critical events must be saved in the same transaction as business data.
- [ ] **Idempotency (Inbox)**: The system must have a mechanism to prevent duplicate processing (Inbox Pattern) for Kafka Consumers and Webhooks.

### 2.2. Network Security
- [ ] **CORS**: Only whitelist official domains.
- [ ] **API Gateway Mapping**: Every resource endpoint must use plural nouns (Plural: `/v1/users`, `/v1/farms`).
- [ ] **Gateway Config Sync**: Run `force-recreate` or `reload` on the gateway to ensure the latest routing map is loaded.
- [ ] **TLS/SSL**: Ensure HTTPS is configured for all public traffic.

---

## 🧪 PHASE 3: Feature Verification
Confirm business logic runs correctly in practice.

### 3.1. Auth Flow (OIDC)
- [ ] **Login/Logout**: Operates smoothly, clears session/cookies upon logout.
- [ ] **Seamless Consent**: Users are not asked for permissions again if already trusted.
- [ ] **Identity Context**: User identity is correctly passed into the UseCase layer.

### 3.2. Observability
- [ ] **Tracing**: Each request generates a Trace ID, User ID appears in Span Attributes.
- [ ] **Logging**: Logs in JSON format (Structured Logging) for easy querying.
- [ ] **Metrics**: Basic metrics (Request count, Latency, Error rate) have been exported.

---

## 🚀 PHASE 4: Deployment & Ops
Final steps before "Go Live".

### 4.1. Orchestration
- [ ] **Dependency Order**: Microservices have `depends_on` with `condition: service_healthy`.
- [ ] **Init Containers**: (K8s) Containers waiting for DB/Identity to be ready before the app runs.
- [ ] **Resource Limits**: CPU/Memory Requests & Limits configured for each container.

### 4.2. CI/CD
- [ ] **Unit Tests**: Passed 100% before merging.
- [ ] **Linter**: No serious linting errors remaining.
- [ ] **Image Security**: Docker images scanned for vulnerabilities.

*Last updated: 2026-05-13 by TechLead Agent*
