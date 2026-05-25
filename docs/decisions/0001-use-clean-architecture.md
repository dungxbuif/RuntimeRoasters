# ADR 0001: Applying Clean Architecture (Consumer-Owned Interfaces)

## Status
**Accepted**

## Context
The RuntimeRoasters project started in Sprint 1 with the goal of building a Microservices platform that is easy to maintain, easy to scale, and especially easy to write Unit Tests for. Traditional frameworks (like MVC) often lead to cross-package circular dependencies, locking business logic to the database.

## Decision
Apply **Clean Architecture Style** for all Microservices.
Core differences:
- **Consumer-Owned Interfaces:** The `usecase` layer (consumer) will declare the interfaces it needs (e.g., `FarmRepository`). 
- The `infrastructure` layer (e.g., `postgres`) will implement these interfaces (leveraging Go's Duck Typing) instead of declaring interfaces in the `domain` or `infrastructure` layers.
- `domain/` contains only Pure Go Structs and does not import any package outside the Standard Library.

## Consequences
- **Positive:** The `usecase` layer is completely decoupled from data storage implementation. Easily use the `mockery` library to generate mocks for Unit Tests.
- **Negative:** Developers transitioning from other languages (Java, C#) may take time to adapt to interfaces not being in the same file as their implementations.

## References
- **Sprint:** Sprint 1 (Bootstrap Phase).
- See implementation details at: [Clean Architecture Concepts](../README.md#-clean-architecture-framework--runtimeroasters).
