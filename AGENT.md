# PROJECT AGENT CONTEXT: Runtime Roasters

This is the working context for AI agents in this repository. Treat it as the first file to read before changing code.

Last updated: 2026-05-19

## Mission

Runtime Roasters is a production-grade Farm-to-Cup supply-chain showcase built as a Go microservices monorepo with a Next.js control-plane UI. The project demonstrates distributed systems patterns: OIDC login, centralized authorization, gRPC/REST gatewaying, Kafka events, transactional outbox, service-local enforcement, and dashboard visualization.

## Current App Shape

- Source root is `src/`.
- Contracts live in `api/runtime/...` and generated Go lives in `src/runtime/...`.
- Services live under `src/apps/`.
- Shared backend packages live under `src/pkg/`.
- Frontend app is `src/apps/client-app`.
- Local infrastructure lives in `deployments/docker-compose.dev.yaml`.

Local ports:

- client-app: `http://localhost:3000`
- KrakenD gateway: `http://localhost:8081`
- auth-service HTTP/gRPC: `8082` / `50052`
- farm-service HTTP/gRPC: `8083` / `50053`
- retail-service HTTP/gRPC: `8084` / `50054`
- logistics-service HTTP/gRPC: `8085` / `50055`
- payment-service HTTP/gRPC: `8086` / `50056`
- trace-service HTTP/gRPC: `8087` / `50057`
- audit-service HTTP/gRPC: `8088` / `50058`
- Ory Kratos public/admin through identity proxy: `4433` / `4434`
- Ory Hydra public/admin: `4444` / `4445`
- Kafka UI: `http://localhost:8090`
- Postgres host port: `54321`
- Elasticsearch: `http://localhost:9200`
- Cassandra: `localhost:9042`

## Architecture Decisions

- Clean Architecture is mandatory:
  - `domain`: pure business entities and constants.
  - `usecase`: business orchestration and repository interfaces.
  - `infrastructure`: DB, Kafka, Ory, external adapters.
  - `delivery`: gRPC/HTTP handlers.
  - `internal/app/init.go`: manual dependency wiring. Do not add DI frameworks.
- Internal service APIs are gRPC-first. Browser traffic goes through Next.js and KrakenD REST endpoints.
- Authorization is centralized-management, distributed-enforcement:
  - Gate 1 is KrakenD JWT validation on protected REST endpoints.
  - Gate 2 is service-local JWT verification plus Casbin enforcement.
  - `auth-service` is the centralized Casbin writer and owns `auth_db.casbin_rule`.
  - business services enforce locally using an in-memory Casbin reader from `pkg/base/security`.
  - readers bootstrap with gRPC `AuthService.GetFullSnapshot`.
  - readers listen to Kafka topic `auth.policy.changed` for live policy refresh.
  - polling is fallback only, not the primary sync mechanism.
- Payment integrations are demo/simulated:
  - Stripe and VNPay provider adapters validate the demo contract only.
  - `/v1/webhooks/stripe` and `/v1/webhooks/vnpay` stay public at the route-auth layer because they represent provider callbacks, but provider signature/business validation still runs in payment-service.
- Retail order SAGA is the current backend happy path:
  - retail creates orders and emits `retail.order.created`.
  - warehouse reserves stock and emits stock events.
  - payment simulates Stripe/VNPay and emits payment events.
  - logistics assigns/delivers shipments.
  - trace-service builds Postgres plus Elasticsearch read models.
  - audit-service writes immutable event audit logs to Cassandra with Postgres fallback.
- Current farm role model:
  - `FARMER` has been removed.
  - Farm operations use `FARM_MANAGER`.
  - Broad farm administration uses `FARM_ADMIN`.
  - System-wide role is `ADMIN`.
- Topology UI:
  - The root page `/` must render the new Architecture Topology.
  - `/dashboard/topology-mesh` reuses the same topology component.
  - Do not restore the old `HighFidelityTopology`/FlowControl topology as the main architecture diagram.

## Current Known Fixes And Lessons

- Gin `NoRoute` plus grpc-gateway can accidentally leave a `404` status on successful `/v1` responses. `pkg/base.App.FinalizeRoutes` must set `200` before delegating to `gwMux`.
- OIDC callback returns `/?access_token=...`; root page must store the token, remove it from the URL, refresh session, and redirect to dashboard users.
- Kratos `session_already_available` means the browser already has a Kratos session. If the app has no local access token, restart Hydra authorization instead of creating another normal Kratos login flow.
- Generated gRPC full method names include the proto package, for example `/runtime.farm.v1.FarmService/ListFarms`. Casbin policies must match those names.
- gRPC authz action mapping is:
  - `Get*` / `List*` -> `read`
  - `Delete*` -> `delete`
  - everything else -> `write`
- Data-level `GormScoper` evaluates role permissions and scopes non-admin results by JWT subject ownership.
- `GET /v1/harvests` is not currently a proto route. The frontend harvest list should use `GET /v1/farms/{id}/harvests`.
- farm-service must not crash when auth-service is not yet ready. Its resilient reader background bootstrap should retry snapshot sync.
- REST business services from retail/payment/logistics/trace/audit require auth-service at startup because their HTTP guards bootstrap JWKS and the Casbin snapshot.
- Kafka consumers default to latest-offset startup for demo flows. Use `deployments/reset-demo-state.sh` before clean local SAGA demos.

## Frontend Rules

- Match existing dashboard visual language: compact operational UI, node cards, restrained colors, dense but scannable controls.
- Avoid marketing landing pages. Build the usable screen directly.
- Do not put cards inside cards. Use cards for discrete items, modals, and framed tools only.
- Root `/` is an architecture topology screen, with minimal access/dashboard action.
- Keep text within containers on mobile and desktop.

## Verification Commands

Use these after touching relevant areas:

```bash
cd src
GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./pkg/base/casbin/... ./apps/auth-service/internal/usecase ./apps/farm-service/internal/...
```

```bash
cd src
GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./pkg/base/auth/... ./pkg/base/casbin/... ./pkg/base/security/... ./apps/auth-service/... ./apps/retail-service/... ./apps/payment-service/... ./apps/logistics-service/... ./apps/trace-service/... ./apps/audit-service/... ./apps/warehouse-service/...
```

```bash
cd src/apps/client-app
npm run lint
```

Useful runtime checks:

```bash
curl -s http://127.0.0.1:8081/v1/users
curl -s -i http://127.0.0.1:8082/health/live
curl -s -i http://127.0.0.1:8083/health/live
curl -s -i http://127.0.0.1:8081/v1/orders/test-id
curl -s -I http://127.0.0.1:3000/
```

Expected `/v1/users` local demo users are `ADMIN` and `FARM_MANAGER` only.
Expected unauthenticated protected gateway calls return `401`.

## Documentation Pointers

- Latest implementation notes: `docs/engineering/current-context.md`
- Architecture overview: `docs/architecture/README.md`
- AuthZ ADR: `docs/architecture/adrs/0003-two-gate-authz-casbin.md`
- Manual test script: `TEST_CASES.md`
- Demo users: `docs/business/demo-identities.md`

## Agent Rules

- Do not revert user changes or unrelated dirty work.
- Prefer `rg` for searching.
- Prefer `apply_patch` for edits.
- Keep changes scoped to the task.
- When changing behavior, update docs in the same turn if the decision affects app context.
- When services are running for manual browser testing, keep them running and report URLs/status instead of killing them at the end.
