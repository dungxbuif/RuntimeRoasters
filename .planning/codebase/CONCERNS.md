# Codebase Concerns

**Analysis Date:** 2025-02-13

## Tech Debt

**Auth Service - Hydra v2 Integration:**
- Issue: Custom claims handling for Hydra v2 is incomplete. The code contains a TODO indicating that custom claims might need to be passed differently in v2, and currently only accepts the login with the subject.
- Files: `src/apps/auth-service/internal/usecase/user_usecase.go`
- Impact: Difficulty in passing rich session context (like user roles or permissions) directly into JWT ID/Access tokens, potentially requiring downstream services to perform extra database lookups.
- Fix approach: Update the Hydra SDK usage to correctly implement session claims for v2 as per Ory Hydra documentation.

**Auth Service - Casbin/Kratos Sync:**
- Issue: `SyncCasbinWithKratos` iterates over all identities fetched from Kratos to synchronize roles.
- Files: `src/apps/auth-service/internal/usecase/user_usecase.go`
- Impact: This approach does not scale. As the number of users grows, this operation will become increasingly slow and could lead to timeouts or memory issues during synchronization.
- Fix approach: Implement a more efficient sync mechanism, such as event-driven updates (listening to Kratos webhooks for identity changes) or incremental synchronization.

**Inconsistent ID Types:**
- Issue: `farm-service` uses `uint64` (autoincrement) for IDs, while `warehouse-service` uses `UUID` strings.
- Files: `src/apps/farm-service/internal/domain/farm.go`, `src/apps/warehouse-service/internal/domain/inventory.go`
- Impact: Minor cognitive load for developers and inconsistency in API design across the microservices.
- Fix approach: Standardize on one ID format (preferably UUIDs for distributed systems) across all services.

## Complexity Hotspots

**Frontend Topology Component:**
- Issue: `ArchitectureTopology.tsx` is a large component (over 400 lines) handling complex D3.js logic, state management, and UI rendering.
- Files: `src/apps/client-app/src/components/features/architecture-topology/ArchitectureTopology.tsx`
- Impact: High risk of regressions when modifying topology logic; difficult to test in isolation.
- Fix approach: Refactor the component into smaller sub-components and extract the D3 simulation logic into a custom hook or utility.

**Warehouse Batch ID Generation:**
- Issue: Batch IDs are generated using a simple random integer (`rand.Intn(999)`).
- Files: `src/apps/warehouse-service/internal/usecase/intake.go`
- Impact: Potential for ID collisions in a high-concurrency production environment where multiple batches are created for the same origin on the same day.
- Fix approach: Use a more robust sequence generator, a database-backed sequence, or include higher precision timestamps/UUIDs in the batch ID.

## Security Risks

**Development Infrastructure Configuration:**
- Risk: Use of default passwords (`password`) and disabled SSL for database connections in development configurations.
- Files: `deployments/docker-compose.dev.yaml`
- Current mitigation: These are restricted to the development environment.
- Recommendations: Ensure production deployment manifests (K8s/Terraform) use secure secrets management (e.g., HashiCorp Vault, AWS Secrets Manager) and enforce SSL.

**Security Test Execution:**
- Risk: E2E Security tests that verify the "Two-Gate" security model (KrakenD + Casbin) skip by default if KrakenD is not reachable.
- Files: `src/apps/farm-service/__tests__/e2e/security_test.go`, `src/apps/demo-service/__tests__/e2e/security_test.go`
- Current mitigation: Developers must manually ensure the infrastructure is running to execute these tests.
- Recommendations: Integrate these tests into a formal CI/CD pipeline where the infrastructure is guaranteed to be available (e.g., using Testcontainers).

## Performance Bottlenecks

**Database Transactions in Kafka Consumers:**
- Problem: Every message processed by the `warehouse-service` worker starts a new database transaction.
- Files: `src/apps/warehouse-service/internal/usecase/intake.go`
- Cause: High overhead of transaction management for every message when processing high-volume streams.
- Improvement path: Implement batch processing of Kafka messages or optimize the idempotency check to avoid full transactions where possible.

**Database-Level Policy Scoping:**
- Problem: `GormScoper` applies Casbin policies to every GORM query.
- Files: `src/pkg/base/casbin/scoper.go`, `src/apps/farm-service/internal/infrastructure/repository/farm_repository.go`
- Cause: Complex authorization policies might result in complex SQL WHERE clauses, potentially slowing down queries as the policy set grows.
- Improvement path: Monitor query performance and ensure appropriate indexes are in place for the fields used in scoping (like `owner_id`).

## Test Coverage Gaps

**Warehouse Service:**
- What's not tested: Core business logic for batch intake, idempotency (Inbox pattern), and inventory management.
- Files: `src/apps/warehouse-service/internal/usecase/`
- Risk: Critical business errors in the warehouse flow could go unnoticed, leading to data corruption or inventory mismatch.
- Priority: High

**Frontend Unit Tests:**
- What's not tested: Component rendering, utility functions, and complex UI logic in the Next.js application.
- Files: `src/apps/client-app/src/`
- Risk: UI regressions are likely as the application grows; reliance on slow E2E tests for verification.
- Priority: Medium

---

*Concerns audit: 2025-02-13*
