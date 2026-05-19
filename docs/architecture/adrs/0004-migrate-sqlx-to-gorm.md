# ADR 0004: Migrating from sqlx to GORM

## Status
**Accepted**

## Context
In Sprint 1 and early Sprint 2, the project used `sqlx` to interact with the PostgreSQL database. `sqlx` offered high performance and flexibility for writing raw SQL.
However, entering Sprint 3 (Farm Service - RR-16) with more complex features, managing relationships, automating basic CRUD operations, and handling transactions with raw SQL became repetitive and prone to human error.

## Decision
Completely replace `sqlx` with **GORM**.
- Update the wrapper in `pkg/database/postgres.go` to initialize connections through `gorm.DB`.
- Refactor Repositories (e.g., FarmRepository) to leverage GORM's chaining syntax (`db.Where().First()`, `db.Create()`).
- Wrap database transaction management logic with a `WithTx` function using `gorm.Transaction()`.

## Consequences
- **Positive:** Increases Developer Velocity by eliminating boilerplate SQL for basic operations. Simplifies handling complex data relationships (Has-Many, Belongs-To) without manual JOINs.
- **Negative:** Introduces some overhead compared to raw SQL due to GORM's use of reflection. The team must understand GORM's Preload mechanism to avoid N+1 Query issues.

## References
- **Sprint:** Cuối Sprint 2 / Đầu Sprint 3.
- **Ticket:** RR-16 (Technical Design: Farm Repository).
- **Discussion:** Developer request to switch ORM ("Switch to GORM, I will use GORM").
