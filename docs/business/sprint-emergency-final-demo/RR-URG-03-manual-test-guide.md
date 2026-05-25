# RR-URG-03 Manual Test Guide

Purpose: verify Warehouse HTTP APIs and the new `harvest -> pickup request -> returned pickup -> intake -> batch -> inventory` flow.

## Accounts

| Screen | Account | Password | Role | Expected scope |
| --- | --- | --- | --- | --- |
| Warehouse dashboard | `warehouse.hn@runtimeroasters.com` | `Hello@123` | `WAREHOUSE_MGR` | Warehouse `WAREHOUSE-HN-001`. |
| Admin overview | `admin@runtimeroasters.com` | `Hello@123` | `ADMIN` | Aggregate view only; not the normal operator. |

If the seeded warehouse account is missing in your local state, create or patch a Kratos identity with `role = WAREHOUSE_MGR` and `warehouse_ids = ["WAREHOUSE-HN-001"]`.

## Test Data

| Field | Value |
| --- | --- |
| `harvest_id` | `HARVEST-RR-URG-03-001` |
| `farm_id` | `FARM-CAUDAT-001` |
| `warehouse_id` | `WAREHOUSE-HN-001` |
| `shipment_id` | `SHIP-RR-URG-03-001` |
| `driver_id` | `DRIVER-DEMO-001` |
| `vehicle_id` | `VEHICLE-DEMO-001` |
| `coffee_type` | `ARABICA` |
| `origin_code` | `SL` |
| `quantity` | `100` |

## Step 1: Start Infrastructure

```bash
docker compose -f deployments/docker-compose.dev.yaml up -d
```

Expected:

- `rr-kafka`, `rr-postgres`, `rr-valkey`, `rr-elasticsearch`, `rr-cassandra`, `rr-clickhouse`, `rr-signoz`, and `rr-otel-collector` are running.

## Step 2: Clean Demo State

```bash
bash deployments/reset-demo-state.sh
```

Expected:

- Command ends with `RuntimeRoasters demo state reset.`
- Warehouse tables `pickup_requests`, `intakes`, `production_batches`, `roast_runs`, `inventories`, and `inbox_events` are empty.

## Step 3: Start Services

Terminal A:

```bash
cd src/apps/auth-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Terminal B:

```bash
cd src/apps/warehouse-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs include consumers for `farm.harvest.created`, `payment.simulated_completed`, `payment.completed`, and `logistics.pickup.arrived_at_warehouse`.
- HTTP is listening on `:8089`.

## Step 3.1: Get A Warehouse Manager Token

Open the Client App, login as `warehouse.hn@runtimeroasters.com`, then read the token from browser local storage:

```bash
cd src/apps/client-app
npm run dev
```

In the browser console after login:

```js
localStorage.getItem('rr_access_token')
```

Set it in your terminal:

```bash
TOKEN='<paste rr_access_token here>'
```

Expected:

- `TOKEN` is a non-empty JWT.
- The JWT contains role `WAREHOUSE_MGR` and `warehouse_ids` including `WAREHOUSE-HN-001`.

## Step 4: Publish Harvest CloudEvent

From repository root:

```bash
EVENT_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
printf '%s\n' 'HARVEST-RR-URG-03-001|{"specversion":"1.0","id":"evt-harvest-rr-urg-03","source":"/services/farm-service","type":"farm.harvest.created","subject":"harvests/HARVEST-RR-URG-03-001","time":"'"${EVENT_TIME}"'","datacontenttype":"application/json","correlationid":"HARVEST-RR-URG-03-001","harvestid":"HARVEST-RR-URG-03-001","farmid":"FARM-CAUDAT-001","warehouseid":"WAREHOUSE-HN-001","data":{"harvest_id":"HARVEST-RR-URG-03-001","coffee_type":"ARABICA","origin_code":"SL","quantity":100}}' | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic farm.harvest.created --property parse.key=true --property key.separator='|'
```

Expected database evidence:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT harvest_id, farm_id, warehouse_id, quantity, status FROM pickup_requests WHERE harvest_id = 'HARVEST-RR-URG-03-001';"
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT count(*) AS intakes FROM intakes WHERE harvest_id = 'HARVEST-RR-URG-03-001';"
```

- One pickup request exists with `status = REQUESTED`.
- Intake count is `0`.

## Step 5: Dispatch Pickup

Get the pickup id:

```bash
PICKUP_ID=$(docker exec rr-postgres psql -U user -d warehouse_db -t -A -c "SELECT id FROM pickup_requests WHERE harvest_id = 'HARVEST-RR-URG-03-001' LIMIT 1;")
curl -i -X POST "http://localhost:8089/v1/warehouse/pickup-requests/${PICKUP_ID}/dispatch" -H "Authorization: Bearer ${TOKEN}"
```

Expected:

- HTTP `200`.
- `pickup_request.status = DISPATCHED`.

## Step 6: Simulate Driver Return Event

```bash
EVENT_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
printf '%s\n' 'SHIP-RR-URG-03-001|{"specversion":"1.0","id":"evt-pickup-arrived-rr-urg-03","source":"/services/logistics-service","type":"logistics.pickup.arrived_at_warehouse","subject":"shipments/SHIP-RR-URG-03-001","time":"'"${EVENT_TIME}"'","datacontenttype":"application/json","correlationid":"HARVEST-RR-URG-03-001","causationid":"evt-pickup-dispatched-rr-urg-03","shipmentid":"SHIP-RR-URG-03-001","harvestid":"HARVEST-RR-URG-03-001","farmid":"FARM-CAUDAT-001","warehouseid":"WAREHOUSE-HN-001","driverid":"DRIVER-DEMO-001","vehicleid":"VEHICLE-DEMO-001","data":{"event_id":"evt-pickup-arrived-rr-urg-03","shipment_id":"SHIP-RR-URG-03-001","harvest_id":"HARVEST-RR-URG-03-001","farm_id":"FARM-CAUDAT-001","warehouse_id":"WAREHOUSE-HN-001","driver_id":"DRIVER-DEMO-001","vehicle_id":"VEHICLE-DEMO-001","status":"ARRIVED_AT_WAREHOUSE","occurred_at":"'"${EVENT_TIME}"'"}}' | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic logistics.pickup.arrived_at_warehouse --property parse.key=true --property key.separator='|'
```

Expected:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT harvest_id, shipment_id, status FROM pickup_requests WHERE harvest_id = 'HARVEST-RR-URG-03-001';"
```

- `status = ARRIVED_WAREHOUSE`.
- `shipment_id = SHIP-RR-URG-03-001`.

## Step 7: Receive Pickup And Create Intake

```bash
curl -i -X POST "http://localhost:8089/v1/warehouse/pickup-requests/${PICKUP_ID}/receive" -H "Authorization: Bearer ${TOKEN}"
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT harvest_id, pickup_id, warehouse_id, quantity, status FROM intakes WHERE harvest_id = 'HARVEST-RR-URG-03-001';"
```

Expected:

- HTTP `200`.
- One intake exists with `status = UNASSIGNED`.
- Re-running the receive curl does not create a second intake.

## Step 8: Create Batch, Process, Finalize

```bash
INTAKE_ID=$(docker exec rr-postgres psql -U user -d warehouse_db -t -A -c "SELECT id FROM intakes WHERE harvest_id = 'HARVEST-RR-URG-03-001' LIMIT 1;")
curl -s -X POST "http://localhost:8089/v1/warehouse/batches" -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" -d "{\"intake_ids\":[\"${INTAKE_ID}\"]}"
BATCH_ID=$(docker exec rr-postgres psql -U user -d warehouse_db -t -A -c "SELECT id FROM production_batches ORDER BY created_at DESC LIMIT 1;")
curl -i -X POST "http://localhost:8089/v1/warehouse/batches/${BATCH_ID}/process" -H "Authorization: Bearer ${TOKEN}"
sleep 8
curl -i -X POST "http://localhost:8089/v1/warehouse/batches/${BATCH_ID}/finalize" -H "Authorization: Bearer ${TOKEN}"
```

Expected:

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT batch_id, warehouse_id, status, total_input_weight, total_output_weight FROM production_batches WHERE id = '${BATCH_ID}';"
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT sku, warehouse_id, available_quantity FROM inventories;"
```

- Batch status is `STOCKED`.
- Inventory contains `SL-ARABICA-ROASTED` with positive `available_quantity`.

## Verification Checklist

- [ ] Infrastructure is running and reset completed.
- [ ] Auth-service and warehouse-service are running.
- [ ] Warehouse manager token is available and contains `WAREHOUSE_MGR` plus `warehouse_ids`.
- [ ] Harvest event creates exactly one `pickup_requests` row.
- [ ] Harvest event does not create an intake immediately.
- [ ] Dispatch endpoint returns `200` and moves pickup to `DISPATCHED`.
- [ ] Driver return event moves pickup to `ARRIVED_WAREHOUSE`.
- [ ] Receive endpoint creates exactly one intake and is idempotent.
- [ ] Batch creation consumes the received intake.
- [ ] Processing simulation moves batch to ready state.
- [ ] Finalize creates finished inventory for `SL-ARABICA-ROASTED`.
