# ADR 0003: Authorization using Two-Gate Hybrid model (Casbin + Data Scoping)

## Status
**Accepted**

## Context
In multi-tenant and microservices systems, authorization is not just about "Does this user have the right to call this API?" (RBAC), but also "Does this user own this data?" (ABAC). If all this logic is embedded in the code, UseCases will bloat with `if-else` statements.

## Decision
Apply the **Two-Gate Hybrid AuthZ** model:
1. **Gate 1 (Edge/Middleware):** Use **Casbin** to check RBAC (Role-Based Access Control). Block invalid requests at Gin Middleware or gRPC Interceptors. Use gRPC Method Name as the Resource identifier in the policy file for consistency.
2. **Gate 2 (Database Layer):** Apply **Data Scoping** in the Repository layer. All SQL queries are forced to include a `WHERE owner_id = $1` clause to protect data at the row level (Row-level security).

Additionally, Casbin policy will be synchronized using a **Resilient Sync** architecture (Snapshot via gRPC + Live Update via Kafka).

## Consequences
- **Positive:** The Business Logic layer (UseCase) remains completely clean, containing no authorization logic. Multiple layers of security (Defense in depth).
- **Negative:** Requires developers to always remember to pass `ownerID` to Repository functions. The authorization synchronization structure (Casbin Sync) increases infrastructure complexity.

## References
- **Sprint:** Sprint 2.
- **Ticket:** RR-12 (Fine-grained Authorization).
- See details at: [Resilient AuthZ Sync](../README.md#resilient-authz-sync-architecture).
