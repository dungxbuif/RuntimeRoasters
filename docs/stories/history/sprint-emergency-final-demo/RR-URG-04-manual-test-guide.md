# RR-URG-04 Manual Test Guide

Purpose: verify Logistics shipment legs, Driver Client backend actions, GPS validation, and mandatory return-to-base.

## Accounts

| Screen | Account | Password | Role | Expected scope |
| --- | --- | --- | --- | --- |
| Driver client | `driver@runtimeroasters.com` | `Hello@123` | `DRIVER` | Can update only assigned shipment. |
| Warehouse dashboard | `warehouse.hn@runtimeroasters.com` | `Hello@123` | `WAREHOUSE_MGR` | Warehouse `WAREHOUSE-HN-001`. |

## Test Data

| Field | Value |
| --- | --- |
| `harvest_id` | `HARVEST-RR-URG-04-001` |
| `farm_id` | `FARM-CAUDAT-001` |
| `warehouse_id` | `WAREHOUSE-HN-001` |
| `driver_id` | `22222222-2222-2222-2222-222222222201` |
| `vehicle_id` | `VEHICLE-DEMO-001` |

## Step 1: Start Infrastructure And Clean State

```bash
docker compose -f deployments/docker-compose.dev.yaml up -d
bash deployments/reset-demo-state.sh
```

Expected:

- Reset ends with `RuntimeRoasters demo state reset.`
- Logistics tables `shipments`, `drivers`, `vehicles`, `locations`, and `processed_kafka_messages` are empty before logistics-service starts and reseeds demo data.

## Step 2: Start Services

Terminal A:

```bash
cd src/apps/auth-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Terminal B:

```bash
cd src/apps/logistics-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs show consumers for `warehouse.stock.reserved`, `warehouse.inventory.updated`, and `warehouse.pickup.requested`.
- HTTP is listening on `:8085`.
- Seeded demo drivers, vehicles, and locations exist.

## Step 3: Get A Driver Token

Open Client App, login as `driver@runtimeroasters.com`, then read browser local storage:

```bash
cd src/apps/client-app
npm run dev
```

Browser console:

```js
localStorage.getItem('rr_access_token')
```

Set it in terminal:

```bash
DRIVER_TOKEN='<paste rr_access_token here>'
```

Expected:

- JWT is non-empty.
- JWT role is `DRIVER`.

## Step 4: Create Farm Pickup Shipment From Warehouse Dispatch Event

Publish only the dispatched pickup event. `REQUESTED` must not create a shipment.

```bash
EVENT_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
printf '%s\n' 'PICKUP-RR-URG-04-001|{"specversion":"1.0","id":"evt-pickup-dispatched-rr-urg-04","source":"/services/warehouse-service","type":"warehouse.pickup.requested","subject":"pickups/PICKUP-RR-URG-04-001","time":"'"${EVENT_TIME}"'","datacontenttype":"application/json","correlationid":"HARVEST-RR-URG-04-001","harvestid":"HARVEST-RR-URG-04-001","farmid":"FARM-CAUDAT-001","warehouseid":"WAREHOUSE-HN-001","data":{"event_id":"evt-pickup-dispatched-rr-urg-04","pickup_id":"PICKUP-RR-URG-04-001","harvest_id":"HARVEST-RR-URG-04-001","farm_id":"FARM-CAUDAT-001","warehouse_id":"WAREHOUSE-HN-001","quantity":100,"status":"DISPATCHED","occurred_at":"'"${EVENT_TIME}"'"}}' | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic warehouse.pickup.requested --property parse.key=true --property key.separator='|'
```

Expected:

```bash
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT id, type, status, current_leg, harvest_id, farm_id, warehouse_id, driver_id, vehicle_id FROM shipments WHERE harvest_id = 'HARVEST-RR-URG-04-001';"
```

- One shipment exists.
- `type = FARM_PICKUP`.
- `status = ASSIGNED`.
- `current_leg = OUTBOUND`.
- `driver_id` and `vehicle_id` are not empty.

## Step 5: Driver Outbound To Farm

```bash
SHIPMENT_ID=$(docker exec rr-postgres psql -U user -d logistics_db -t -A -c "SELECT id FROM shipments WHERE harvest_id = 'HARVEST-RR-URG-04-001' LIMIT 1;")
curl -i -X POST "http://localhost:8085/v1/logistics/shipments/${SHIPMENT_ID}/depart" -H "Authorization: Bearer ${DRIVER_TOKEN}"
curl -i -X POST "http://localhost:8085/v1/logistics/drivers/location" -H "Authorization: Bearer ${DRIVER_TOKEN}" -H "Content-Type: application/json" -d "{\"shipment_id\":\"${SHIPMENT_ID}\",\"lat\":21.0285,\"lng\":105.8542,\"route_index\":1,\"status\":\"IN_TRANSIT_TO_FARM\",\"occurred_at\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}"
curl -i -X POST "http://localhost:8085/v1/logistics/shipments/${SHIPMENT_ID}/arrive" -H "Authorization: Bearer ${DRIVER_TOKEN}"
curl -i -X POST "http://localhost:8085/v1/logistics/shipments/${SHIPMENT_ID}/confirm-load" -H "Authorization: Bearer ${DRIVER_TOKEN}"
```

Expected:

```bash
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT status, current_leg, departed_at, arrived_at_origin_at, loaded_at FROM shipments WHERE id = '${SHIPMENT_ID}';"
```

- `status = PICKED_UP`.
- `current_leg = OUTBOUND`.
- Driver and vehicle are still busy, not available.

## Step 6: Mandatory Return To Warehouse

```bash
curl -i -X POST "http://localhost:8085/v1/logistics/shipments/${SHIPMENT_ID}/return" -H "Authorization: Bearer ${DRIVER_TOKEN}"
curl -i -X POST "http://localhost:8085/v1/logistics/drivers/location" -H "Authorization: Bearer ${DRIVER_TOKEN}" -H "Content-Type: application/json" -d "{\"shipment_id\":\"${SHIPMENT_ID}\",\"lat\":21.0286,\"lng\":105.8543,\"route_index\":2,\"status\":\"RETURNING_TO_WAREHOUSE\",\"occurred_at\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}"
curl -i -X POST "http://localhost:8085/v1/logistics/shipments/${SHIPMENT_ID}/arrive" -H "Authorization: Bearer ${DRIVER_TOKEN}"
```

Expected:

```bash
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT status, current_leg, return_started_at, returned_at FROM shipments WHERE id = '${SHIPMENT_ID}';"
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT id, status, is_available, current_shipment_id FROM drivers WHERE id = '22222222-2222-2222-2222-222222222201';"
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT id, status FROM vehicles WHERE id = 'VEHICLE-DEMO-001';"
```

- Shipment `status = ARRIVED_WAREHOUSE`.
- Driver `status = IDLE`, `is_available = true`, `current_shipment_id` empty.
- Vehicle `status = IDLE`.
- Warehouse-service can consume `logistics.pickup.arrived_at_warehouse` and enable warehouse receipt/intake in RR-URG-03.

## Negative Checks

1. Try GPS without token:

```bash
curl -i -X POST "http://localhost:8085/v1/logistics/drivers/location" -H "Content-Type: application/json" -d "{\"shipment_id\":\"${SHIPMENT_ID}\",\"lat\":21,\"lng\":105}"
```

Expected: `401`.

2. Try returning before pickup/loading on a fresh assigned shipment:

```bash
curl -i -X POST "http://localhost:8085/v1/logistics/shipments/${SHIPMENT_ID}/return" -H "Authorization: Bearer ${DRIVER_TOKEN}"
```

Expected after the shipment already returned: `400` because transition is invalid from `ARRIVED_WAREHOUSE`.

## Verification Checklist

- [ ] `warehouse.pickup.requested` with `status = DISPATCHED` creates one `FARM_PICKUP` shipment.
- [ ] Shipment gets driver and vehicle assignment.
- [ ] Driver can depart, arrive at farm, and confirm loading.
- [ ] Driver can post GPS only for assigned shipment.
- [ ] Driver/vehicle stay busy before return leg.
- [ ] Return action moves shipment to `RETURNING_TO_WAREHOUSE`.
- [ ] Final arrive action moves shipment to `ARRIVED_WAREHOUSE`.
- [ ] Driver and vehicle become available only after return.
- [ ] No-token GPS request is rejected.
