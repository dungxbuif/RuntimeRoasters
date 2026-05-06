---
name: runtime-roasters-reviewer
description: A specialized code review skill for the RuntimeRoasters project, focusing on Clean Architecture v4, Go best practices, and microservices security.
license: MIT
metadata:
  author: RuntimeRoasters Tech Lead
  version: "1.0.0"
---

# 🕵️ Code Review Skill: RuntimeRoasters Special

This skill defines the rigorous review process for the RuntimeRoasters ecosystem. It ensures that all contributions align with the project's high standards for architectural purity, security, and reliability.

## 🏛️ Clean Architecture v4 Standards

The most critical rule: **Interfaces belong to the consumer.**

### 1. Layer Purity
-   **Domain Layer**: Must be pure Go. No imports from `infrastructure`, `usecase`, or external frameworks (except `time` or basic `math`).
-   **UseCase Layer**: Declares its own Repository and Service interfaces. Imports only `domain`.
-   **Infrastructure Layer**: Implements UseCase interfaces. Imports `usecase`, `domain`, and `pkg/*`.
-   **Composition Root**: All DI wiring MUST happen in `cmd/main.go` using Google Wire.

### 2. Dependency Rule
-   Dependencies point **INWARDS** only.
-   `infrastructure` → `usecase` → `domain`
-   If you see `domain` importing `usecase`, it is a **BLOCKING** architectural violation.

## 🐹 Idiomatic Go Excellence

### 1. Error Handling
-   Use `pkg/errs` sentinel errors (`ErrNotFound`, `ErrConflict`, etc.).
-   Wrap errors with `%w` to preserve context.
-   Return RFC 9457 Problem Details for all API failures.

### 2. Concurrency & Context
-   `context.Context` is the first parameter of every I/O or blocking function.
-   Always check for `ctx.Done()` in long-running loops.
-   No goroutine leaks: every `go func()` must have a clear exit strategy or be managed by a worker pool.

## 🔒 Security & Identity (Ory Stack)

### 1. Authentication
-   Extract identity from `X-User-ID` and `X-Role` headers (injected by Gateway/Middleware).
-   Never trust user-provided IDs in the request body if they conflict with headers.

### 2. Data Protection
-   Sanitize all inputs before database queries.
-   Use `sqlx` parameterized queries to prevent SQL injection.
-   Log masking: Ensure sensitive data (tokens, passwords, emails) are not logged.

## 📡 Distributed Reliability

### 1. Transactional Outbox
-   Operations that mutate state AND emit events MUST use `pkg/database.WithTx`.
-   Verify that outbox events are written to the `outbox` table in the same transaction.

### 2. Idempotency
-   Check for `Idempotency-Key` implementation for POST/PUT requests.
-   Ensure event consumers implement the `Inbox` pattern to prevent duplicate processing.

## 🧪 Testing Requirements
-   **Unit Tests**: Table-driven tests for all business logic.
-   **Mocks**: Use `mockery` to generate mocks for interfaces defined in `usecase`.
-   **Coverage**: Target >80% for `domain` and `usecase` layers.
