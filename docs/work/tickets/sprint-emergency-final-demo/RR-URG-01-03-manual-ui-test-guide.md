# RR-URG-01 To RR-URG-03 Manual UI Test Guide

Purpose: one consolidated guide for manually verifying the emergency sprint scope up to RR-URG-03: trace propagation, canonical event contracts/storage projections, and the harvest-to-warehouse pickup/intake/batch flow.

## Scope Covered

- RR-URG-01: end-to-end `trace_id` propagation and SigNoz/trace-service visibility.
- RR-URG-02: paid-order SAGA events project into payment, warehouse, logistics, trace, and audit stores.
- RR-URG-03: harvest creates a warehouse pickup request first; intake and inventory happen only after dispatch and returned pickup receipt.

RR-URG-04 driver-return UI is not part of this guide.

## Accounts

Default password for every account below:

```text
Hello@123
```

| Flow | Account | Role | Use For | Expected Scope |
| --- | --- | --- | --- | --- |
| Public topology | no login | public | Open root architecture showcase. | Root page is intentionally public. |
| Admin overview | `admin@runtimeroasters.com` | `ADMIN` | View global overview and management screens. | Aggregate/global view. |
| Farm harvest | `manager.caudat@runtimeroasters.com` | `FARM_MANAGER` | Create harvest from UI. | Assigned farm data only. |
| Warehouse ops | `warehouse.hn@runtimeroasters.com` | `WAREHOUSE_MGR` | View warehouse stock/batches; get token for warehouse actions. | `WAREHOUSE-HN-001`. |
| Retail paid order | `mgr.hn.hoankiem@runtimeroasters.com` | `STORE_MGR` | Create paid/replenishment order from UI. | Store `11111111-1111-1111-1111-111111111101`. |

## Important Current UI Notes

- `/` is public and should show Architecture Topology without login.
- `/dashboard/farm-ops/harvests` has UI controls for creating harvests.
- `/dashboard/retail/orders` has UI controls for creating retail orders.
- `/dashboard/warehouse` shows batch/inventory state, but RR-URG-03 pickup queue dispatch/receive controls are not yet exposed as dedicated UI buttons. Use curl for dispatch/receive in this guide.
- `/dashboard/logistics` currently has a map/monitor screen. Its simulation button only posts GPS for real backend shipments with seeded UUID drivers; fallback visual mock rows are display-only and must not publish GPS updates.
- `/dashboard/traceability` can search a known business ID after trace-service has projected events.

## Step 0: Start Infrastructure

From repository root:

```bash
docker compose -f deployments/docker-compose.dev.yaml up -d
docker compose -f deployments/docker-compose.dev.yaml ps
```

Expected:

- `rr-postgres`, `rr-kafka`, `rr-valkey`, `rr-elasticsearch`, `rr-cassandra`, `rr-signoz-clickhouse`, `rr-signoz-zookeeper-1`, `rr-signoz`, `rr-otel-collector`, `rr-kratos`, and `rr-hydra` are running.
- Single-node Elasticsearch may be `yellow`; that is acceptable.

Quick health checks:

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
- Postgres returns `accepting connections`.
- Valkey returns `PONG`.
- Cassandra prints keyspaces.
- Kafka prints topic names.

## Step 1: Reset Demo State

Use this before a clean manual run:

```bash
bash deployments/reset-demo-state.sh
```

Expected:

```text
RuntimeRoasters demo state reset.
```

## Step 2: Start Backend Services

Option A: VS Code debug mode:

1. Open RuntimeRoasters in VS Code.
2. Open **Run and Debug**.
3. Select **Debug All Backend**.
4. Press **F5**.

Expected:

- Auth, Farm, Retail, Logistics, Payment, Trace, Audit, and Warehouse start.
- No port conflict errors.

Option B: terminal dev mode with Air hot reload:

```bash
task be
```

If `task` is not installed:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

If `air` is not installed:

```bash
go install github.com/air-verse/air@latest
```

Make sure Go binaries are on your shell `PATH`:

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
air -v
```

Manual fallback without Task, open one terminal per service and run Air directly:

```bash
cd src/apps/auth-service
air -c .air.toml
```

```bash
cd src/apps/farm-service
air -c .air.toml
```

```bash
cd src/apps/retail-service
air -c .air.toml
```

```bash
cd src/apps/logistics-service
air -c .air.toml
```

```bash
cd src/apps/payment-service
air -c .air.toml
```

```bash
cd src/apps/trace-service
air -c .air.toml
```

```bash
cd src/apps/audit-service
air -c .air.toml
```

```bash
cd src/apps/warehouse-service
air -c .air.toml
```

Expected service ports:

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

## Step 3: Start Client App

```bash
cd src/apps/client-app
npm run dev
```

Open:

```text
http://localhost:3000
```

Expected:

- Root page loads without login.
- Architecture Topology is visible.
- No auth wall appears on `/`.

## Step 4: Login Pattern

For each role-specific test:

1. Open `http://localhost:3000`.
2. Click the login/dashboard action.
3. Enter the email and password from the account table.
4. Submit the login form.
5. After redirect, verify the dashboard sidebar is visible.

To switch accounts reliably:

- Use a separate browser profile/incognito window per account, or clear site data for `localhost`.
- Do not keep stale local storage tokens between role tests.

Optional token check from browser console after login:

```js
localStorage.getItem('rr_access_token')
```

Expected:

- Value is a non-empty JWT.

## Step 5: RR-URG-01 Public Topology And Trace Observability Smoke Test

Browser:

1. Open `http://localhost:3000/`.
2. Do not log in.
3. Confirm the architecture topology canvas is visible.
4. Open `http://localhost:3301`.

Expected UI result:

- Root page is accessible without auth.
- SigNoz UI loads.

Runtime checks:

```bash
curl -s -i http://127.0.0.1:8081/__health
docker logs rr-otel-collector --tail 80
```

Expected:

- KrakenD health returns `200`.
- OTel collector logs do not show a persistent `nop` pipeline or repeated export failures.

## Step 6: RR-URG-02 Paid Order SAGA From UI

Browser A: login as:

```text
mgr.hn.hoankiem@runtimeroasters.com
Hello@123
```

Then:

1. Open sidebar item **Market Orders** or direct URL:

   ```text
   http://localhost:3000/dashboard/retail/orders
   ```

2. In **Receiving Store Node**, select **Hoan Kiem Store**.
3. In **Manifest Inventory**, use:

   ```text
   SKU: SL-ARABICA-ROASTED
   Quantity: 1
   ```

   If that SKU is unavailable in your local inventory run, use the default UI SKU but expect warehouse reservation may fail until inventory exists.

4. Click **Broadcast Order**.
5. The UI redirects to **Saga Monitor**.

Expected UI result:

- No `Saga Initialization Failed` error.
- You land on `/dashboard/retail?...`.

Find the latest order in DB:

```bash
ORDER_ID=$(docker exec rr-postgres psql -U user -d retail_db -t -A -c "SELECT id FROM orders ORDER BY created_at DESC LIMIT 1;")
echo "${ORDER_ID}"
```

Verify payment:

```bash
docker exec rr-postgres psql -U user -d payment_db -c "SELECT order_id, store_id, status, provider, amount FROM payments WHERE order_id = '${ORDER_ID}';"
```

Expected:

- One payment row.
- `status = SUCCEEDED`.

Verify logistics:

```bash
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT id, order_id, destination_store_id, driver_id, status FROM shipments WHERE order_id = '${ORDER_ID}';"
```

Expected:

- One shipment row.
- `destination_store_id = 11111111-1111-1111-1111-111111111101`.
- `driver_id` is populated.

Verify trace projection:

```bash
docker exec rr-postgres psql -U user -d trace_db -c "SELECT topic, order_id, store_id, shipment_id, trace_id FROM trace_events WHERE order_id = '${ORDER_ID}' ORDER BY occurred_at;"
```

Expected topics:

```text
retail.order.created
payment.intent.created
payment.simulated_completed
warehouse.stock.reserved
logistics.delivery.assigned
```

Expected data:

- `order_id` equals `${ORDER_ID}`.
- `store_id = 11111111-1111-1111-1111-111111111101`.
- `trace_id` is populated.

Verify audit projection:

```bash
docker exec rr-postgres psql -U user -d audit_db -c "SELECT topic, partition_key, store_id FROM audit_logs WHERE partition_key = '${ORDER_ID}' ORDER BY occurred_at;"
```

Expected:

- Same canonical topics are visible in audit.
- `partition_key = ${ORDER_ID}`.

Optional UI traceability:

1. Login as any account that can access the dashboard.
2. Open **Provenance Trace** or:

   ```text
   http://localhost:3000/dashboard/traceability
   ```

3. Enter `${ORDER_ID}`.
4. Click **Trace**.

Expected:

- If trace-service has already built the document for this ID, the page shows **Ledger Verified** and a timeline.
- If the page says **Trace Failed**, use the DB query above as RR-URG-02 evidence and retry after a few seconds.

## Step 7: RR-URG-03 Harvest Creates Pickup Request, Not Intake

Browser B: login as:

```text
manager.caudat@runtimeroasters.com
Hello@123
```

Then:

1. Open sidebar item **Harvest Declaration** or direct URL:

   ```text
   http://localhost:3000/dashboard/farm-ops/harvests
   ```

2. Click **Declare New Harvest**.
3. Select an available **Origin Farm**.
4. Select:

   ```text
   Coffee Variety: Arabica
   Yield: 100
   Harvest Date: today
   ```

5. Click **Record Harvest**.

Expected UI result:

- Modal closes.
- Harvest appears in the harvest ledger table.
- Status appears as a new/active harvest state.

Find the latest harvest:

```bash
HARVEST_ID=$(docker exec rr-postgres psql -U user -d farm_db -t -A -c "SELECT id FROM harvests ORDER BY created_at DESC LIMIT 1;")
echo "${HARVEST_ID}"
```

Verify warehouse pickup request:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT id, harvest_id, farm_id, warehouse_id, quantity, status FROM pickup_requests WHERE harvest_id = '${HARVEST_ID}';"
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT count(*) AS intakes FROM intakes WHERE harvest_id = '${HARVEST_ID}';"
```

Expected:

- Exactly one pickup request exists.
- `status = REQUESTED`.
- Intake count is `0`.

This proves the RR-URG-03 behavior change: harvest no longer creates intake immediately.

## Step 8: RR-URG-03 Dispatch Pickup

Browser C: login as:

```text
warehouse.hn@runtimeroasters.com
Hello@123
```

Then open:

```text
http://localhost:3000/dashboard/warehouse
```

Expected UI result before receipt:

- Warehouse page loads.
- Finished stock and processing queue are visible.
- The new pickup request is not yet shown as intake/batch because dispatch/receipt is still pending.

Get warehouse token from browser console:

```js
localStorage.getItem('rr_access_token')
```

Set it in terminal:

```bash
TOKEN='<paste rr_access_token here>'
PICKUP_ID=$(docker exec rr-postgres psql -U user -d warehouse_db -t -A -c "SELECT id FROM pickup_requests WHERE harvest_id = '${HARVEST_ID}' LIMIT 1;")
```

Dispatch through KrakenD:

```bash
curl -i -X POST "http://localhost:8081/v1/warehouse/pickup-requests/${PICKUP_ID}/dispatch" -H "Authorization: Bearer ${TOKEN}"
```

Expected:

- HTTP `200`.
- Pickup status becomes `DISPATCHED`.

Verify:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT harvest_id, status, dispatched_at FROM pickup_requests WHERE id = '${PICKUP_ID}';"
```

## Step 9: RR-URG-03 Simulate Returned Pickup

Until RR-URG-04 UI is used, simulate logistics returning to warehouse by publishing the canonical event:

```bash
EVENT_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
printf '%s\n' "SHIP-RR-URG-03-MANUAL|{\"specversion\":\"1.0\",\"id\":\"evt-pickup-arrived-${HARVEST_ID}\",\"source\":\"/services/logistics-service\",\"type\":\"logistics.pickup.arrived_at_warehouse\",\"subject\":\"shipments/SHIP-RR-URG-03-MANUAL\",\"time\":\"${EVENT_TIME}\",\"datacontenttype\":\"application/json\",\"correlationid\":\"${HARVEST_ID}\",\"causationid\":\"${PICKUP_ID}\",\"shipmentid\":\"SHIP-RR-URG-03-MANUAL\",\"harvestid\":\"${HARVEST_ID}\",\"farmid\":\"FARM-CAUDAT-001\",\"warehouseid\":\"WAREHOUSE-HN-001\",\"driverid\":\"DRIVER-DEMO-001\",\"vehicleid\":\"VEHICLE-DEMO-001\",\"data\":{\"event_id\":\"evt-pickup-arrived-${HARVEST_ID}\",\"shipment_id\":\"SHIP-RR-URG-03-MANUAL\",\"harvest_id\":\"${HARVEST_ID}\",\"farm_id\":\"FARM-CAUDAT-001\",\"warehouse_id\":\"WAREHOUSE-HN-001\",\"driver_id\":\"DRIVER-DEMO-001\",\"vehicle_id\":\"VEHICLE-DEMO-001\",\"status\":\"ARRIVED_AT_WAREHOUSE\",\"occurred_at\":\"${EVENT_TIME}\"}}" | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic logistics.pickup.arrived_at_warehouse --property parse.key=true --property key.separator='|'
```

Verify:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT harvest_id, shipment_id, status FROM pickup_requests WHERE id = '${PICKUP_ID}';"
```

Expected:

- `status = ARRIVED_WAREHOUSE`.
- `shipment_id = SHIP-RR-URG-03-MANUAL`.

## Step 10: RR-URG-03 Receive Pickup And Create Intake

Receive through KrakenD:

```bash
curl -i -X POST "http://localhost:8081/v1/warehouse/pickup-requests/${PICKUP_ID}/receive" -H "Authorization: Bearer ${TOKEN}"
```

Verify intake:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT id, harvest_id, pickup_id, warehouse_id, quantity, status FROM intakes WHERE harvest_id = '${HARVEST_ID}';"
```

Expected:

- HTTP `200`.
- One intake exists.
- Intake status is `UNASSIGNED`.
- Re-running the receive curl does not create a second intake.

Idempotency check:

```bash
curl -i -X POST "http://localhost:8081/v1/warehouse/pickup-requests/${PICKUP_ID}/receive" -H "Authorization: Bearer ${TOKEN}"
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT count(*) AS intakes FROM intakes WHERE harvest_id = '${HARVEST_ID}';"
```

Expected:

- Intake count remains `1`.

## Step 11: RR-URG-03 Create Batch, Process, Finalize, Stock In

Create batch:

```bash
INTAKE_ID=$(docker exec rr-postgres psql -U user -d warehouse_db -t -A -c "SELECT id FROM intakes WHERE harvest_id = '${HARVEST_ID}' LIMIT 1;")
curl -s -X POST "http://localhost:8081/v1/warehouse/batches" -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" -d "{\"intake_ids\":[\"${INTAKE_ID}\"]}"
BATCH_ID=$(docker exec rr-postgres psql -U user -d warehouse_db -t -A -c "SELECT id FROM production_batches ORDER BY created_at DESC LIMIT 1;")
echo "${BATCH_ID}"
```

Process:

```bash
curl -i -X POST "http://localhost:8081/v1/warehouse/batches/${BATCH_ID}/process" -H "Authorization: Bearer ${TOKEN}"
```

Wait for processing simulation:

```bash
sleep 8
```

Finalize:

```bash
curl -i -X POST "http://localhost:8081/v1/warehouse/batches/${BATCH_ID}/finalize" -H "Authorization: Bearer ${TOKEN}"
```

Verify DB:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT batch_id, warehouse_id, status, total_input_weight, total_output_weight FROM production_batches WHERE id = '${BATCH_ID}';"
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT sku, warehouse_id, available_quantity FROM inventories ORDER BY updated_at DESC;"
```

Expected:

- Batch status is `STOCKED`.
- Inventory contains finished roasted SKU such as `SL-ARABICA-ROASTED`.
- `available_quantity > 0`.

Browser:

1. Refresh `http://localhost:3000/dashboard/warehouse`.
2. Watch **Finished Stock**.

Expected UI result:

- Inventory item appears or quantity increases.
- Processing queue no longer shows the finalized batch as active.

## Step 12: Optional Traceability UI For Harvest/Warehouse Flow

Open:

```text
http://localhost:3000/dashboard/traceability
```

Search one of:

```text
${HARVEST_ID}
${BATCH_ID}
SHIP-RR-URG-03-MANUAL
```

Expected:

- If trace document exists, page shows **Ledger Verified** and event count.
- If not yet available, verify via DB:

```bash
docker exec rr-postgres psql -U user -d trace_db -c "SELECT topic, harvest_id, warehouse_id, shipment_id, trace_id FROM trace_events WHERE harvest_id = '${HARVEST_ID}' OR shipment_id = 'SHIP-RR-URG-03-MANUAL' ORDER BY occurred_at;"
```

## Pass Checklist

- [ ] Docker infrastructure is running.
- [ ] Backend services are running once each; no duplicate service instances on the same ports.
- [ ] Client app opens at `http://localhost:3000`.
- [ ] Root Architecture Topology is visible without login.
- [ ] STORE_MGR can create a retail order from **Market Orders**.
- [ ] Paid order creates payment, warehouse reservation, logistics shipment, trace rows, and audit rows.
- [ ] FARM_MANAGER can create a harvest from **Harvest Declaration**.
- [ ] Harvest creates pickup request with `REQUESTED` status.
- [ ] Harvest does not create intake immediately.
- [ ] WAREHOUSE_MGR can dispatch pickup via gateway API.
- [ ] Returned pickup event changes pickup to `ARRIVED_WAREHOUSE`.
- [ ] Receive creates exactly one intake.
- [ ] Batch process/finalize creates finished inventory.
- [ ] Warehouse UI shows finished stock after refresh.
- [ ] SigNoz is available at `http://localhost:3301`.

## Troubleshooting

If login loops:

- Clear browser site data for `localhost`.
- Start again from `http://localhost:3000`.

If protected gateway call returns `401`:

- Re-login and refresh `TOKEN` from `localStorage.getItem('rr_access_token')`.

If gateway call returns `403`:

- Confirm the account role matches the action.
- Warehouse actions require `warehouse.hn@runtimeroasters.com`.

If a service does not start:

- Check port conflicts:

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

If trace rows are missing:

- Confirm trace-service is running.
- Confirm Kafka topics contain events.
- Check trace-service logs for consumer errors.

If paid order does not reserve stock:

- Run the warehouse intake/batch/finalize flow first to create inventory.
- Retry the order using SKU `SL-ARABICA-ROASTED`.
