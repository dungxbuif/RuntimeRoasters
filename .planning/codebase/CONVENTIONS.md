# Coding Conventions

**Analysis Date:** 2025-02-13

## Naming Patterns

**Files:**
- Go files: `snake_case.go` (e.g., `farm_usecase.go`, `problem.go`)
- Test files: `[name]_test.go` (e.g., `farm_usecase_test.go`)
- Documentation: `kebab-case.md` or `UPPERCASE.md` (e.g., `technical_design.md`, `STACK.md`)

**Functions:**
- Go: `PascalCase` for exported functions (e.g., `NewFarmUsecase`), `camelCase` for unexported.
- Handlers: Usually `[Action][Entity]` (e.g., `CreateFarm`, `GetFarmByID`).

**Variables:**
- Go: `camelCase` for local variables, `PascalCase` for exported package-level variables or struct fields.
- Context: Always named `ctx`.
- Logger: Usually named `log` or `ctxLog`.

**Types:**
- Structs/Interfaces: `PascalCase` (e.g., `FarmUsecase`, `Repository`).

## Code Style

**Formatting:**
- `gofmt` and `goimports` are standard.
- Linting uses `golangci-lint` (as seen in `src/Makefile`).

**Linting:**
- Tool: `golangci-lint`
- Key rules: Standard Go recommendations, `context` propagation, and error checking.

## Import Organization

**Order:**
1. Standard library imports
2. Third-party imports
3. Local project imports (`github.com/dungxbuif/RuntimeRoasters/...`)

**Path Aliases:**
- The project uses the full module path: `github.com/dungxbuif/RuntimeRoasters/...`

## Error Handling

**Patterns:**
- **Sentinel Errors:** Defined in `src/pkg/errs/problem.go` (e.g., `ErrNotFound`, `ErrValidation`).
- **Standardized Responses:** Uses RFC 9457 (Problem Details for HTTP APIs) via `pkg/errs`.
- **gRPC Mapping:** Errors are mapped to gRPC status codes using `errs.ToGRPCError(err)` in handlers.
- **No Raw Errors:** Never return raw errors from UseCase to Transport layer; always wrap or use sentinel errors.

## Logging

**Framework:** `zap` (Uber's structured logger)

**Patterns:**
- **Traceability:** Every log entry should include a `trace_id`.
- **Context-Aware:** Use `logger.FromContext(ctx)` to automatically attach TraceID.
- **Levels:** `Info` for standard operations, `Warn` for client errors (4xx), `Error` for internal server errors (5xx).

## Comments

**When to Comment:**
- Exported functions, types, and constants should have a descriptive comment.
- Complex logic or business rules that aren't self-explanatory.

**JSDoc/TSDoc:**
- Frontend (`client-app`) follows standard Next.js/React patterns.

## Function Design

**Size:** Functions should be focused and ideally small. UseCase methods typically coordinate multiple repository calls.

**Parameters:**
- First parameter is almost always `context.Context`.
- Use structs for complex parameter lists.

**Return Values:**
- Usually `(result, error)`.

## Module Design

**Exports:**
- Use interfaces for dependency injection (defined in the same package as the consumer or in a dedicated domain layer).
- Exported constructors `New[TypeName]` return the implementation struct or interface.

**Barrel Files:**
- Not typically used in Go; Go uses package-level exports.
- Frontend uses standard ESM exports.

## Project Management Conventions

**Jira-in-Markdown:**
- Every task must have a directory in `docs/business/sprintX/RR-x/`.
- `ticket.md`: Business requirements (BA/PO role).
- `technical_design.md`: Technical implementation plan (Dev/Tech Lead role).
- Use `subtickets/` for large Epics.

---

*Convention analysis: 2025-02-13*
