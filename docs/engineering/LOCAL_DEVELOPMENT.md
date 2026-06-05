---
artifact_type: local_development
id: LOCAL_DEVELOPMENT
status: active
owner: shared
updated: 2026-06-05
---

# Local Development

## Key Ports

- Client App: `http://localhost:3000`
- KrakenD Gateway: `http://localhost:8081`
- Auth Service: HTTP `8082`, gRPC `50052`
- Farm Service: HTTP `8083`, gRPC `50053`
- Retail Service: HTTP `8084`, gRPC `50054`
- Logistics Service: HTTP `8085`, gRPC `50055`
- Payment Service: HTTP `8086`, gRPC `50056`
- Trace Service: HTTP `8087`, gRPC `50057`
- Audit Service: HTTP `8088`, gRPC `50058`
- Warehouse Service: HTTP `8089`, gRPC `50059`
- Socket Service: HTTP `8091`, gRPC `50060`
- Kafka UI: `http://localhost:8090`
- Elasticsearch: `http://localhost:9200`
- SigNoz UI: `http://localhost:3301`
- OTLP: gRPC `4317`, HTTP `4318`
- Postgres: `localhost:54321`
- Cassandra: `localhost:9042`
- Kratos Public: `http://localhost:4433`
- Hydra Public: `http://localhost:4444`

## Common Checks

Backend targeted tests:

```bash
cd src
GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./pkg/base/auth/... ./pkg/base/casbin/... ./pkg/base/security/... ./apps/auth-service/... ./apps/retail-service/... ./apps/payment-service/... ./apps/logistics-service/... ./apps/trace-service/... ./apps/audit-service/... ./apps/warehouse-service/...
```

Frontend checks:

```bash
cd src/apps/client-app
npm run lint
npm run build
```

Gateway/config smoke checks:

```bash
curl -s -i http://127.0.0.1:8082/health/live
curl -s -I http://127.0.0.1:3000/
```

## Demo Simulation Boundaries

Use `docs/requirements/REQUIREMENTS.md` before adding shortcuts or demo-only data.

## Harness Workflow

For non-small work, create or update a ticket in `docs/work/tickets/` and a detail design before implementation.
