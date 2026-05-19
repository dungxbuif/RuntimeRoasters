# Codebase Concerns

**Analysis Date:** 2025-05-15

## Tech Debt

**Centralized Authorization Logic:**
- Issue: A gap exists in the Casbin matcher where Gate 2 (service-level RBAC) might incorrectly block role policies. Current matchers in `auth-service` and `farm-service` use `keyMatch` and `regexMatch` but may not handle nested role inheritance or specific route patterns consistently.
- Files: `src/apps/auth-service/internal/infrastructure/casbin/model.conf`, `src/apps/farm-service/configs/rbac_model.conf`, `src/pkg/base/casbin/transport/grpc/interceptor.go`
- Impact: Authorization failures for valid roles or unauthorized access if regex patterns are too broad.
- Fix approach: Standardize the Casbin matcher and validate against comprehensive integration tests using real JWT claims.

**Auth Service Hydra Integration:**
- Issue: Custom claims for Hydra v2 Go SDK are pending implementation.
- Files: `src/apps/auth-service/internal/usecase/user_usecase.go`
- Impact: Potential missing identity metadata in issued tokens.
- Fix approach: Update the `AcceptHydraLogin` logic to correctly map domain claims to Hydra's expected format.

**Large Usecase Bloat:**
- Issue: Several service usecases are growing large, exceeding 300 lines with multiple responsibilities (logic + event publishing).
- Files: `src/apps/payment-service/internal/usecase/service.go`, `src/apps/trace-service/internal/usecase/service.go`
- Impact: Increased maintenance difficulty and harder unit testing.
- Fix approach: Split usecases into smaller, single-purpose interactors or domain services.

## Security Considerations

**Record-Level Scoping Consistency:**
- Risk: While `payment-service` implements store-based scoping, other services like `farm-service` or `retail-service` may lack consistent record-level checks, relying only on role-level gRPC method access.
- Files: `src/apps/payment-service/internal/usecase/service.go`, `src/pkg/base/identity/scope.go`
- Current mitigation: `StoreScopeFromContext` helper in `identity` package.
- Recommendations: Implement GORM Scopers across all store-scoped entities to automatically inject `WHERE store_id IN (...)` clauses.

**Casbin Policy Distribution:**
- Risk: Policies are synced from `auth-service` via gRPC snapshots and Kafka. If the Kafka relay fails, services fall back to polling, which may introduce a delay in revoking permissions.
- Files: `src/pkg/base/casbin/reader.go`
- Current mitigation: Background synchronization with retry and polling fallback.
- Recommendations: Add a "Last Synced" metric to Prometheus/SigNoz to alert on stale authorization state.

## Performance Bottlenecks

**Snapshot-based Auth Sync:**
- Problem: `GetFullSnapshot` returns the entire policy set. As the user base and role assignments grow, this payload will increase in size and processing time.
- Files: `src/apps/auth-service/internal/delivery/grpc/handler.go`, `src/pkg/base/casbin/reader.go`
- Cause: Full state transfer instead of incremental updates.
- Improvement path: Implement incremental policy updates via Kafka for high-frequency changes, reserving full snapshots for bootstrap only.

## Fragile Areas

**Duplicate/Stale Service Directories:**
- Files: `src/apps/warehouse-service copy/`
- Why fragile: Contains duplicate code with space in import paths, causing build inconsistencies and developer confusion.
- Safe modification: Remove the directory and ensure `warehouse-service` is the single source of truth.
- Test coverage: Unknown.

**Manual Action Mapping:**
- Files: `src/pkg/base/casbin/transport/grpc/interceptor.go`
- Why fragile: Hardcoded mapping of gRPC method names (e.g., `/Get*` to `read`) is prone to errors if method naming conventions are not strictly followed.
- Safe modification: Use proto annotations to explicitly define the required action/resource for each RPC.

## Test Coverage Gaps

**Integration Tests for Auth Flows:**
- What's not tested: Full end-to-end flow from login challenge to service-level enforcement with specific roles like `STORE_MGR`.
- Files: `src/apps/auth-service/`, `src/apps/farm-service/`
- Risk: Regressions in security policy changes might not be detected until production.
- Priority: High

---

*Concerns audit: 2025-05-15*
