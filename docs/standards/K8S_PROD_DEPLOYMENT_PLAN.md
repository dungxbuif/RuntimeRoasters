# Runtime Roasters - Production Kubernetes Deployment Plan

## Summary

Deploy Runtime Roasters to Kubernetes with the current architecture intact:
Next.js client, KrakenD gateway, stateless Go services, Ory Kratos/Hydra identity,
Kafka event mesh, Postgres service databases, Valkey socket/session state,
Elasticsearch trace read model, Cassandra audit storage, and SigNoz/OpenTelemetry
observability.

Default packaging choice: create an umbrella Helm chart under
`deployments/k8s/helm/runtime-roasters` with environment values files for dev,
staging, and production. Production should prefer managed infrastructure for
Postgres, Kafka, Elasticsearch, and Redis/Valkey. Cassandra and SigNoz may be
managed or self-hosted depending on the target platform.

## Production Checklist

### Images

- Build one immutable image per runtime component:
  `client-app`, `krakend`, `auth-service`, `farm-service`, `retail-service`,
  `warehouse-service`, `logistics-service`, `payment-service`, `trace-service`,
  `audit-service`, and `socket-service`.
- Tag images by git SHA or release version. Do not deploy `latest` to production.
- Run vulnerability scanning before promotion.

### Config And Secrets

- Store non-sensitive config in ConfigMaps.
- Store passwords, client secrets, webhook secrets, internal secrets, and socket
  API key hashes in Kubernetes Secrets or an external secret manager.
- Rotate all development secrets before production:
  `INTERNAL_SECRET`, Hydra client secret, database passwords, Stripe demo webhook
  key, and socket internal API keys.
- Keep issuer/JWKS consistent:
  public issuer uses the production public domain, while internal JWKS should use
  service DNS or the protected identity proxy.

### Data And Migrations

- Run database migrations as Kubernetes Jobs before application rollout.
- Keep migrations service-owned and idempotent.
- Do not use GORM AutoMigrate as production compatibility proof.
- Production bootstrap should create only the minimum admin/OIDC data.
- Demo seed data belongs in staging/demo environments, not production.

### Network

- Public ingress exposes only:
  `client-app`, `krakend`, Kratos public endpoint, and Hydra public endpoint if
  required by browser OIDC flow.
- Go service HTTP/gRPC ports remain internal ClusterIP services.
- WebSocket traffic routes through ingress/KrakenD to `socket-service`.
- TLS is mandatory for all public traffic.
- CORS allowlist must contain only official production domains.

### Runtime

- Every service has liveness and readiness probes.
- Services with startup dependencies use startup probes or init containers where
  appropriate.
- Resource requests and limits are set for every pod.
- Add PodDisruptionBudgets for critical surfaces:
  gateway, client app, auth, retail, payment, trace, and socket.
- Use HPA only after baseline metrics are available.

### Observability

- Deploy OTel Collector as the in-cluster telemetry target.
- All services export OTLP to collector DNS, not localhost.
- SigNoz or the selected observability backend receives traces, metrics, and logs.
- Minimum dashboards:
  request rate, error rate, p95 latency, Kafka consumer lag, DB connection pool,
  webhook failures, and socket connection count.

## Kubernetes Implementation Plan

### Helm Layout

- Create umbrella chart: `deployments/k8s/helm/runtime-roasters`.
- Add shared templates for Go services:
  Deployment, HTTP Service, gRPC Service, probes, resources, env, and envFrom.
- Add component templates for:
  client app, KrakenD, Ory Kratos/Hydra, migration jobs, and OTel collector.
- Add values files:
  `values-dev.yaml`, `values-staging.yaml`, and `values-prod.yaml`.

### Workloads

- Deploy stateless apps as Deployments.
- Deploy migrations and identity bootstrap as Jobs.
- Deploy public endpoints through Ingress.
- Deploy internal service discovery through ClusterIP Services.
- Mount KrakenD, Ory, and OTel configuration from ConfigMaps.
- Inject secrets through Kubernetes Secrets or external secret integration.

### Stateful Dependencies

- Production default:
  managed Postgres, Kafka, Redis/Valkey, and Elasticsearch.
- Staging/demo fallback:
  self-hosted Helm charts with persistent volumes.
- Kafka must remain KRaft-compatible. Do not reintroduce ZooKeeper for Kafka.

### Deployment Order

1. Create namespace, ConfigMaps, and Secrets.
2. Verify external datastore connectivity.
3. Run Kratos/Hydra migrations and identity/OIDC bootstrap.
4. Run business service migrations.
5. Deploy `auth-service`.
6. Deploy business/support services:
   farm, retail, warehouse, logistics, payment, trace, audit, and socket.
7. Deploy KrakenD.
8. Deploy client app.
9. Run smoke tests and staging E2E gates before traffic promotion.

## Verification Gate

### Static Gate

- `go test ./apps/... ./pkg/... ./runtime/...`
- `npx tsc --noEmit`
- `npm run lint`
- `npm test`
- `jq empty deployments/krakend/krakend.json`
- `helm template` for every environment values file.
- Kubernetes schema validation for rendered manifests.

### Smoke Gate

- Gateway health returns `200`.
- Unauthenticated protected gateway request returns `401`.
- OIDC login/callback works on production domain.
- `GET /v1/auth/me` returns the expected authenticated identity.
- `GET /v1/traces/public/topology/config` works without auth.
- `POST /v1/realtime/tickets` requires JWT.
- Public and private WebSocket routes connect through ingress.

### Business E2E Gate

- Admin bootstrap path works.
- Farm harvest emits Kafka event and warehouse pickup appears.
- Paid order creates payment intent.
- Stripe-compatible webhook pass emits `payment.completed`.
- Warehouse reserves stock only after `payment.completed`.
- Stripe-compatible webhook fail emits `payment.failed` and does not create a
  stock reservation.
- Trace document route `/v1/trace/{id}/document` works through KrakenD.
- One order flow keeps one trace ID across gateway, services, Kafka,
  trace-service, and the observability backend.

## Assumptions

- First Kubernetes target is staging, then production promotion after smoke and
  E2E pass.
- Helm umbrella chart is the default deployment packaging.
- Managed infrastructure is preferred for production.
- Topology config is not runtime-admin editable.
- Realtime uses WebSocket only. SSE remains out of scope.
