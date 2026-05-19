# ADR 0002: Choosing SigNoz over traditional stack (Prometheus/Jaeger/Loki)

**Status:** 🔴 REVERTED (Removed - 2026-05-10)

## Context
The project initially chose SigNoz as a centralized Observability platform to replace individual containers for Jaeger, Prometheus, etc.

## Decision
Discontinue the use of SigNoz and ClickHouse for now to reduce local resource usage for developers. The system will revert to lighter solutions (e.g., standalone Jaeger) or focus on completing business logic first.

## Consequences
- Must update `docker-compose.dev.yaml` and telemetry configuration in the code.
- Reduces RAM/CPU load on dev machines (ClickHouse consumes significant resources).
