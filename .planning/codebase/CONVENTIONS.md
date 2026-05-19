# Coding Conventions

**Analysis Date:** 2025-05-14

## Naming Patterns

**Files:**
- Snake case for Go files: `farm_usecase.go`, `farm_handler.go`.
- Kebab case for directories (apps): `farm-service`, `auth-service`.
- Single word for package names: `usecase`, `domain`, `delivery`, `infrastructure`.

**Functions:**
- PascalCase for exported functions: `NewFarmUsecase`, `CreateFarm`.
- camelCase for unexported functions: `mapToProto`.

**Variables:**
- camelCase for local variables and struct fields: `farmReq`, `ownerID`.
- Short names for receivers: `u` for `farmUsecase`, `h` for `FarmHandler`.

**Types:**
- PascalCase for exported types: `FarmUsecase`, `FarmRepository`.
- camelCase for unexported implementation types: `farmUsecase`.

## Code Style

**Formatting:**
- `gofmt` (implied by Go standards).
- `goimports` for import management.

**Linting:**
- `golangci-lint` is used as per `src/Makefile`.
- Standard rules likely apply as no custom configuration was found in the root.

## Import Organization

**Order:**
1. Standard library imports (e.g., `context`, `errors`).
2. Internal repository imports (e.g., `github.com/dungxbuif/RuntimeRoasters/apps/...`).
3. Third-party library imports (e.g., `go.uber.org/zap`, `github.com/stretchr/testify`).

**Path Aliases:**
- `farmv1` for `github.com/dungxbuif/RuntimeRoasters/runtime/farm/v1`.
- `svcconfig` for local service config packages.

## Error Handling

**Patterns:**
- Use RFC 9457 (Problem Details) for API responses: `src/pkg/errs/problem.go`.
- Define sentinel errors in `pkg/errs` or domain layer: `ErrNotFound`, `ErrInvalidFarmName`.
- Wrap errors with `%w` to preserve context: `fmt.Errorf("%w: %v", errs.ErrValidation, err)`.
- Use `errs.ToGRPCError(err)` in delivery layer to map internal errors to gRPC status codes.

## Logging

**Framework:** `zap` (via `go.uber.org/zap`).

**Patterns:**
- Use context-aware logging to include TraceIDs: `logger.FromContext(ctx)`.
- Log level based on environment: `production` uses json, `development` uses colored console output.
- Avoid logging secrets or tokens.

## Comments

**When to Comment:**
- Exported functions, types, and constants should have descriptive comments.
- Complex logic or architectural decisions should be documented inline.

**JSDoc/TSDoc:**
- Not applicable (Go-focused backend). Go doc comments are used.

## Function Design

**Size:** Functions are generally small and focused on a single responsibility.

**Parameters:** `context.Context` is always the first parameter for I/O bound or cross-cutting concern functions.

**Return Values:** Typically returns `(result, error)`.

## Module Design

**Exports:** Interfaces are exported to allow mocking and decoupling. Implementation structs are often unexported.

**Barrel Files:** Not used in Go. Package-level organization is used instead.

**Architecture Layers (Clean Architecture):**
- **Domain (`internal/domain`):** Pure entities and business logic. No external imports.
- **UseCase (`internal/usecase`):** Business use cases. Defines interfaces for repositories and external services.
- **Delivery (`internal/delivery`):** Transport layer (gRPC, REST handlers). Maps transport models to domain models.
- **Infrastructure (`internal/infrastructure`):** Implementation details (database repositories, event publishers).

---

*Convention analysis: 2025-05-14*
