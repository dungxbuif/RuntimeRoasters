---
artifact_type: troubleshooting
id: TROUBLESHOOTING
status: active
owner: shared
updated: 2026-06-05
---

# Troubleshooting

## OIDC Callback Or Login Loop

Known behavior:

- OIDC callback may return `/?access_token=...`.
- Root page must store the token, remove it from URL, refresh session, and redirect dashboard users.
- `session_already_available` means Kratos already has a browser session. If the app lacks a local access token, restart Hydra authorization rather than creating another normal Kratos login flow.

## Auth/Casbin Snapshot Startup

- Business services require auth-service readiness for HTTP guard bootstrap.
- Farm-service must retry Casbin snapshot sync instead of crashing if auth-service is not ready.
- Generated gRPC method names include proto package names, for example `/runtime.farm.v1.FarmService/ListFarms`.

## Gateway 404 On Successful `/v1` Responses

Gin `NoRoute` plus grpc-gateway can accidentally leave a `404` status on successful `/v1` responses.

Expected fix point:

- `pkg/base.App.FinalizeRoutes` must set `200` before delegating to the gateway mux.

## KrakenD JWT Validation

- KrakenD runs inside Docker and validates JWTs against internal Hydra DNS: `http://hydra:4444/.well-known/jwks.json`.
- Do not point KrakenD at protected host JWKS unless the required internal header is sent.

## E2E Server Binding

Sandboxed commands may fail to bind `127.0.0.1:3000` with `EPERM`.

Expected approach:

- Use approved escalation for Playwright runs that need local server binding.
- Prefer production `next start` E2E when `next dev` is unstable.
- Keep Playwright timeouts fail-fast to avoid hanging.

## Kafka/OTel

- Kafka uses Apache Kafka without ZooKeeper.
- SigNoz ClickHouse coordination is not Kafka coordination.
- Outbox relays must preserve W3C trace context.
- Kafka inbox/dedupe must prefer `event_id` or `topic:key`; offsets are fallback only.
