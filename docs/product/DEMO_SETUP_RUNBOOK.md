# Runtime Roasters Demo Setup Runbook

This is the single operational runbook for installing, starting, resetting, seeding, and opening the local demo playground.

## 1. Prerequisites

Required:

```text
Docker Desktop / Docker Compose
Go 1.22+
Node.js 20+
Task CLI
Air
```

Install Task:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

Install Air:

```bash
go install github.com/air-verse/air@latest
```

Make sure Go binaries are on `PATH`:

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
task --version
air -v
```

## 2. Environment Files

Service-local `.env` files are already present for local dev in this workspace. If a fresh clone only has examples, create them from examples:

```bash
cp deployments/.env.example deployments/.env
cp src/apps/auth-service/.env.example src/apps/auth-service/.env
cp src/apps/farm-service/.env.example src/apps/farm-service/.env
cp src/apps/retail-service/.env.example src/apps/retail-service/.env
cp src/apps/logistics-service/.env.example src/apps/logistics-service/.env
cp src/apps/payment-service/.env.example src/apps/payment-service/.env
cp src/apps/trace-service/.env.example src/apps/trace-service/.env
cp src/apps/audit-service/.env.example src/apps/audit-service/.env
cp src/apps/warehouse-service/.env.example src/apps/warehouse-service/.env
cp src/apps/client-app/.env.example src/apps/client-app/.env.local
```

## 3. Start Order

Use this order for a clean local demo:

```bash
task infra
```

Expected infrastructure:

- Postgres
- Kafka
- Kafka UI
- Valkey
- Elasticsearch
- Kibana
- Cassandra
- Kratos
- Hydra
- Identity proxy
- KrakenD
- SigNoz + ClickHouse
- OTel collector

Check containers:

```bash
docker compose -f deployments/docker-compose.dev.yaml ps
```

Reset demo business state:

```bash
bash deployments/reset-demo-state.sh
```

Seed identity/OIDC bootstrap data:

```bash
task seed
```

Start backend services with Air hot reload:

```bash
task be
```

Start frontend:

```bash
task fe
```

Or start backend + frontend together:

```bash
task dev
```

Do not run the same service twice. If VS Code Debug is running a service, do not also start that service through `task be`.

## 4. Playground Ports

| Surface | URL / Port | Purpose |
| --- | --- | --- |
| Client App root | `http://localhost:3000/` | Public Architecture Topology showcase. No login required. |
| Dashboard | `http://localhost:3000/dashboard` | Authenticated control-plane UI. |
| API docs page | `http://localhost:3000/api-docs` | Client-side Swagger/OpenAPI viewer. |
| KrakenD gateway | `http://localhost:8081` | Browser/API entrypoint. |
| Kafka UI | `http://localhost:8090` | Inspect topics/messages. |
| SigNoz | `http://localhost:3301` | OpenTelemetry traces. |
| Kibana | `http://localhost:5601` | Elasticsearch UI. |
| Elasticsearch | `http://localhost:9200` | Trace/search read model backend. |
| Kratos public | `http://localhost:4433` | Identity public/proxy endpoint. |
| Hydra public | `http://localhost:4444` | OAuth2/OIDC public endpoint. |
| Postgres | `localhost:54321` | Shared local Postgres host port. |
| Valkey | `localhost:6379` | Cache/GPS state. |
| Cassandra | `localhost:9042` | Audit storage. |
| OTLP gRPC / HTTP | `localhost:4317` / `localhost:4318` | Service trace export. |

Backend service ports:

| Service | HTTP | gRPC |
| --- | --- | --- |
| Auth | `8082` | `50052` |
| Farm | `8083` | `50053` |
| Retail | `8084` | `50054` |
| Logistics | `8085` | `50055` |
| Payment | `8086` | `50056` |
| Trace | `8087` | `50057` |
| Audit | `8088` | `50058` |
| Warehouse | `8089` | `50059` |

## 5. Seed And Demo Data Sources

Executable seed/reset files:

| File / Source | What It Seeds Or Resets | When It Runs |
| --- | --- | --- |
| `deployments/init-db.sql` | Creates Postgres databases: `identity_db`, `hydra_db`, `auth_db`, `farm_db`, `warehouse_db`, `retail_db`, `logistics_db`, `payment_db`, `trace_db`, `audit_db`. | First Postgres volume init. |
| `deployments/seed.sh` | Seeds Kratos admin identity and Hydra `client-app`. | `task seed` or compose `seeder`. |
| `deployments/kratos/seed-admin.json` | Admin identity: `admin@runtimeroasters.com` / `Hello@123` / `ADMIN`. | Read by `deployments/seed.sh`. |
| `deployments/hydra/client-app.json` | OAuth2 client `client-app`, secret `client-secret`, callback `http://localhost:3000/api/auth/callback`. | Read by `deployments/seed.sh`. |
| `deployments/reset-demo-state.sh` | Deletes canonical Kafka demo topics, truncates business tables, deletes trace index, drops Cassandra audit keyspace, restarts KrakenD. | Manual reset before demo. |
| `src/apps/retail-service/internal/seed/stores.json` | Retail stores: Hoan Kiem, Cau Giay, District 1, District 7, Hai Chau. | Retail service startup; loaded by embedded seed package and inserted idempotently. |
| `src/apps/logistics-service/internal/seed/logistics.json` | Demo vehicles, UUID drivers, warehouse/farm/store locations. | Logistics service startup; loaded by embedded seed package and inserted idempotently. |
| `src/apps/client-app/public/data/routes.json` | Browser-side route geometry for logistics visualization/simulation. | Client app runtime. |
| `src/apps/logistics-service/testdata/routes.json` | Service-side/static route dataset used for logistics testing/demo references. | Test/demo reference. |
| `src/apps/auth-service/internal/infrastructure/casbin/default_policies.csv` | Default Casbin role policies. | Auth service startup. |

Reference/demo-planning seed files:

| File | Status |
| --- | --- |
| `deployments/logistics-seed.sql` | Older reference SQL for logistics locations/drivers. Current runtime seed is `src/apps/logistics-service/internal/seed/logistics.json`. Do not treat this SQL file as the primary seed. |
| `docs/product/GUIDE.md` account/persona tables | Demo contract/reference. Not all listed identities are currently guaranteed by `deployments/seed.sh`. |

## 6. Current Seeded Accounts

Guaranteed by current executable seeder:

| Email | Password | Role | Source |
| --- | --- | --- | --- |
| `admin@runtimeroasters.com` | `Hello@123` | `ADMIN` | `deployments/kratos/seed-admin.json` |

Documented demo personas used by sprint guides:

| Email | Password | Role | Notes |
| --- | --- | --- | --- |
| `manager.caudat@runtimeroasters.com` | `Hello@123` | `FARM_MANAGER` | Needed for harvest UI demo. Create via Admin/User flow if absent. |
| `warehouse.hn@runtimeroasters.com` | `Hello@123` | `WAREHOUSE_MGR` | Needs `warehouse_ids = ["WAREHOUSE-HN-001"]`. Create via Admin/User flow if absent. |
| `driver@runtimeroasters.com` | `Hello@123` | `DRIVER` | Logistics DB seeds a driver row with `user_id = driver@runtimeroasters.com`; Kratos identity may still need creation if absent. |
| `mgr.hn.hoankiem@runtimeroasters.com` | `Hello@123` | `STORE_MGR` | Needs `store_ids = ["11111111-1111-1111-1111-111111111101"]`. Create via Admin/User flow if absent. |

If a documented persona cannot log in, first log in as admin and create/patch that identity with the role/scope above. This mismatch is a known seed-data gap to close in RR-URG-08.

## 7. What Reset Does And Does Not Do

`deployments/reset-demo-state.sh` resets demo runtime state:

- Kafka canonical demo topics are deleted if present.
- Retail orders/outbox/inbox are truncated.
- Payment rows/inbox are truncated.
- Warehouse pickup/intake/batch/inventory/inbox tables are truncated.
- Logistics shipments/drivers/vehicles/locations/inbox/dedupe tables are truncated.
- Trace events/documents are truncated.
- Audit logs are truncated.
- Elasticsearch `coffee_traceability` index is deleted.
- Cassandra keyspace `runtime_roasters_audit` is dropped.
- KrakenD is restarted.

Additional clean cache step:

```bash
docker exec rr-valkey valkey-cli FLUSHDB
```

After reset, start backend services again. Retail and logistics services seed their master data from their `internal/seed/*.json` files on startup.

## 8. Health Checks

Run after `task infra`:

```bash
curl -s -i http://127.0.0.1:8081/__health
curl -s -i http://127.0.0.1:4433/health/ready
curl -s -i http://127.0.0.1:4444/health/ready
curl -s -i http://127.0.0.1:3301/api/v1/health
docker exec rr-postgres pg_isready -U user -d postgres
docker exec rr-valkey valkey-cli ping
docker exec rr-cassandra cqlsh -e 'DESCRIBE KEYSPACES' 127.0.0.1 9042
docker exec rr-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

Expected:

- HTTP checks return `200`.
- Postgres says `accepting connections`.
- Valkey returns `PONG`.
- Cassandra prints keyspaces.
- Kafka prints topic names.
- Elasticsearch can be `yellow` in single-node local dev.

## 9. Start Playing

Recommended first local playthrough:

1. Open public topology:

   ```text
   http://localhost:3000/
   ```

   Expected: Architecture Topology is visible without login.

2. Open dashboard:

   ```text
   http://localhost:3000/dashboard
   ```

3. Login as:

   ```text
   admin@runtimeroasters.com
   Hello@123
   ```

4. If role demo accounts are not present, create the personas listed in this runbook.

5. Run the Sprint 1-3 manual guide:

   ```text
   docs/stories/history/sprint-emergency-final-demo/RR-URG-01-03-manual-ui-test-guide.md
   ```

6. Watch support UIs:

   ```text
   Kafka UI: http://localhost:8090
   SigNoz:   http://localhost:3301
   Kibana:   http://localhost:5601
   ```

## 10. Useful Commands

Start infra:

```bash
task infra
```

Seed admin/OIDC:

```bash
task seed
```

Reset demo business data:

```bash
bash deployments/reset-demo-state.sh
docker exec rr-valkey valkey-cli FLUSHDB
```

Start backend:

```bash
task be
```

Start frontend:

```bash
task fe
```

Start everything:

```bash
task dev
```

Stop containers:

```bash
docker compose -f deployments/docker-compose.dev.yaml down
```

Check port conflicts:

```bash
lsof -nP -iTCP:8082 -sTCP:LISTEN
lsof -nP -iTCP:8083 -sTCP:LISTEN
lsof -nP -iTCP:8084 -sTCP:LISTEN
lsof -nP -iTCP:8085 -sTCP:LISTEN
lsof -nP -iTCP:8086 -sTCP:LISTEN
lsof -nP -iTCP:8087 -sTCP:LISTEN
lsof -nP -iTCP:8088 -sTCP:LISTEN
lsof -nP -iTCP:8089 -sTCP:LISTEN
```
