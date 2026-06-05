# Code Standards

Project-specific code standards should be added by humans as the technology stack is selected.

## Default Rules

- Prefer existing local patterns over new abstractions.
- Keep changes scoped to the active ticket or bug.
- Avoid unrelated refactors.
- Add abstractions only when they remove real complexity or match an established pattern.
- Do not add dependencies without documenting the reason and impact.

## Dependency Rules

Adding or replacing a major dependency requires:

- Detail design approval
- ADR
- Update to relevant setup or deployment docs
- Test evidence

## Naming And Structure

- Go services follow `cmd/`, `internal/domain`, `internal/usecase`, `internal/infrastructure`, and `internal/delivery` boundaries.
- Raw errors must not cross from usecase to transport; map service errors through `pkg/errs`.
- Propagate `context.Context` from handler to usecase to repository to preserve identity and trace context.
- Logs must use context-aware structured logging so trace IDs remain searchable.
- Domain validation belongs close to domain entities or explicit validators, not scattered across handlers.
- Multi-table writes must use an explicit transaction boundary.

## Project-Specific Engineering Guardrails

# Engineering Guardrails For Agents

These rules are mandatory for all future coding agents and reviewers. They are derived from real defects found during RR-URG-01 trace end-to-end debugging.

## 1. Service Bootstrap Must Be Unified

Rule:

- Every backend service must initialize the shared operational baseline: logger, config, OpenTelemetry tracer/propagator, database instrumentation, auth/Casbin where applicable, Kafka producer/consumer wrappers, readiness checks, and graceful shutdown.
- Prefer `pkg/base.App` for HTTP/gRPC services.
- If a service cannot use `pkg/base.App` because it is worker-only, its `internal/app/init.go` must explicitly initialize the same baseline and document the exception.

Forbidden:

- A service that consumes or produces Kafka without `telemetry.InitTracer`.
- A service that opens GORM directly instead of `database.NewPostgres`.
- A worker service that silently skips tracing, config defaults, or readiness because it has no HTTP server.

Review checklist:

- `otel.SetTextMapPropagator(...)` is reachable through bootstrap.
- `database.NewPostgres(...)` is used for Postgres access.
- Kafka uses `pkg/kafka.NewProducer` and `pkg/kafka.NewConsumer`.
- Service startup evidence includes the service name and readiness or consumer startup log.

## 2. Migrations Must Match Real Existing Database State

Rule:

- Any schema change must be safe against the current local/demo database, not only a fresh database.
- Before adding/changing GORM models, inspect the existing table schema when the table may already exist.
- `AutoMigrate` is allowed for simple additive fields, but it must not be trusted for relationship/foreign-key changes across old schemas.
- If an existing table type differs from a new model, write an explicit migration or an `ensure*Table` function that is idempotent and avoids invalid foreign keys.

Forbidden:

- Adding a GORM relationship that creates an FK against an existing incompatible column type.
- Assuming `AutoMigrate` will repair old enum, uuid, varchar, jsonb, or FK mismatches.
- Marking a migration done without running it against the current demo DB.

Required verification:

```bash
docker exec rr-postgres psql -U user -d <service_db> -c "\d <table>"
GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/<service>/internal/...
```

## 3. Transactional Outbox Must Preserve Trace Context

Rule:

- If a request writes an outbox event and a later relay publishes it, the outbox row must persist W3C trace context.
- Store at least `traceparent`; store `tracestate` when present.
- The relay must extract that context before calling `producer.Publish`.

Forbidden:

- Publishing outbox events from `context.Background()` without restoring request trace context.
- Using only a custom `trace_id` metadata value while dropping W3C `traceparent`.
- Treating outbox as a business-only record when the event is part of an observable SAGA.

Required verification:

- Trigger one HTTP request.
- Verify downstream Kafka consumers and trace-service projection preserve the same trace ID.

## 4. Kafka Idempotency Must Not Depend Only On Offset

Rule:

- Consumer idempotency keys must prefer business-stable identifiers.
- Use `pkg/kafka.MessageID(msg)` or equivalent logic:
  1. `topic:event_id` when payload has `event_id`.
  2. `topic:key` when key is present.
  3. `topic-partition-offset` only as a legacy fallback.

Forbidden:

- Using only `fmt.Sprintf("%s-%d-%d", topic, partition, offset)` for inbox or trace dedupe.
- Assuming Kafka offsets survive local broker recreation, topic reset, or demo state reset.

Required verification:

- Unit tests must cover event ID, topic key, and offset fallback.
- Local demo reset must not cause new business events to be skipped as duplicates.

## 5. Kafka And Gateway Must Use Current Internal Topology

Rule:

- Local Kafka runs with the official Apache Kafka image and no ZooKeeper dependency. Do not reintroduce ZooKeeper.
- SigNoz/ClickHouse observability may have its own coordination service. Do not conflate it with Kafka topology.
- KrakenD must resolve JWKS through internal Docker service DNS when it runs in Docker.
- Protected gateway endpoints must be tested with a real UI-issued token, not only direct service calls.

Forbidden:

- Adding `KAFKA_ZOOKEEPER_CONNECT` or connecting Kafka to any ZooKeeper/coordinator service.
- Pointing KrakenD JWT validator at an externally protected JWKS proxy that requires headers KrakenD does not send.
- Declaring auth/gateway done after direct service-port tests only.

Required verification:

```bash
docker compose -f deployments/docker-compose.dev.yaml config --services
docker exec rr-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
npx playwright test e2e/manual-live-evidence.spec.ts --project=chromium
```

## 6. End-To-End Evidence Must Cover The Actual Flow

Rule:

- For SAGA or trace work, evidence must prove the full intended path, not a partial direct-service shortcut.
- The accepted minimum for RR-URG trace work is:
  - browser login,
  - KrakenD protected endpoint,
  - service handler,
  - outbox/Kafka producer,
  - Kafka consumer,
  - downstream event publish,
  - trace-service Postgres/Elasticsearch projection.

Forbidden:

- Replacing gateway tests with direct service-port calls unless the gateway issue is itself the reported blocker.
- Stopping after the first event appears in trace-service when downstream topics are required.
- Accepting multiple trace IDs for one intended SAGA when the ticket requires one connected trace.

Required evidence:

- Command output or saved markdown under `docs/work/tickets/.../evidence/`.
- Business IDs used in the run.
- The trace ID and the topics/events observed.

## 7. Process Hygiene During Live Debugging

Rule:

- Before live verification, ensure only one instance of each service consumer is running.
- When using `go run`, old compiled binaries can remain alive even if the shell session is gone; verify ports and process list.
- Restart all affected producers and consumers after changing shared Kafka, telemetry, or idempotency code.

Forbidden:

- Running evidence with mixed old/new service binaries in the same consumer groups.
- Leaving local Go/Next dev processes running after a task unless the user explicitly asked to keep them.
- Trusting stale logs after a restart-sensitive change.

Useful checks:

```bash
lsof -nP -iTCP:8082 -sTCP:LISTEN
lsof -nP -iTCP:8083 -sTCP:LISTEN
lsof -nP -iTCP:8084 -sTCP:LISTEN
lsof -nP -iTCP:8085 -sTCP:LISTEN
lsof -nP -iTCP:8086 -sTCP:LISTEN
lsof -nP -iTCP:8087 -sTCP:LISTEN
lsof -nP -iTCP:8088 -sTCP:LISTEN
```

## 8. Storage Boundary Rules Are Enforced In Code

Rule:

- Workflow-critical fields must be typed columns, indexed when queried or scoped.
- JSONB is for payload snapshots and flexible metadata, not for fields that drive auth, status transitions, dashboards, or joins.
- Elasticsearch is the traceability read model, not the operational source of truth.
- Cassandra is append-only audit/live history, not operational state.

Forbidden:

- Hiding `store_id`, `farm_id`, `warehouse_id`, `driver_id`, `shipment_id`, `order_id`, `status`, or timestamps only inside JSONB.
- Writing operational workflow state only to Elasticsearch or Cassandra.

## 9. Seed Data Must Have Explicit Sources

Rule:

- Demo/master data must live in explicit seed files, SQL migrations, or documented admin bootstrap payloads.
- Service startup seed is allowed only when it loads a checked-in seed file and writes idempotently.
- Usecases must not hide master data in hardcoded slices; they may accept seed data from a seed package and perform idempotent persistence.
- Data that requires privileged identity creation or cross-service role assignment must be triggered through an explicit `ADMIN` bootstrap flow and must be idempotent.

Forbidden:

- Adding new demo stores, vehicles, drivers, routes, warehouses, farms, or personas only as hardcoded literals inside a usecase.
- Making a seed depend on non-deterministic autogenerated IDs when other services or docs need to reference that record.
- Creating admin/manager/driver identities from a normal service startup path without explicit `ADMIN` intent.

Review checklist:

- Seed file path is documented in the runbook.
- Re-running the seed does not duplicate records.
- Required business IDs are stable and match manual test guides.
- Any account/persona bootstrap clearly states which role may execute it.

## 10. Ticket Close Criteria

A backend integration ticket is not done until:

- Relevant unit tests pass.
- Relevant service tests pass.
- A live/manual evidence command has been run when the ticket changes cross-service behavior.
- Docs are updated in the same turn when a rule, architecture decision, or demo flow changes.
- Known residual risk is explicitly named, not silently ignored.
