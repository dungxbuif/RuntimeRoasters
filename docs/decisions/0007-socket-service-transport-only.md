# ADR 0007: Socket Service Is Transport Only

## Status

Accepted.

## Context

Sprint 6/7 needs realtime public topology and private dashboard streams without
creating another durable source of truth. Trace-service already owns PostgreSQL
and Elasticsearch read models for queryable demo history.

## Decision

Add `socket-service` as a DB-free WebSocket fanout service.

- It uses Gorilla WebSocket.
- It consumes Kafka traceable topics plus `socket.broadcast.requested`.
- It accepts internal API-key pushes at `POST /internal/v1/socket/events`.
- It stores only ephemeral socket session, presence, and latest-event hints in
  Valkey with TTL.
- It does not own durable history, topology config, or business state.
- Trace-service owns topology config/history query APIs.

## Consequences

- Public topology can use push from socket-service and pull from trace-service.
- Private streams require JWT/Casbin plus event-level role/entity filtering.
- Multi-replica socket fanout still needs a later design; Valkey session state
  does not by itself solve cross-replica broadcast delivery.
