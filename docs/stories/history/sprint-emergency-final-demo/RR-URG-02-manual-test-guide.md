# RR-URG-02 Manual Test Guide

Purpose: verify event contracts, business IDs, `traceid` propagation, trace projection, audit projection, and storage boundaries before continuing later urgent tickets.

## Accounts

RR-URG-02 contract verification is terminal-first and does not require login because it publishes CloudEvents directly to Kafka.

Use these accounts only when opening UI screens while watching the result:

| Screen | Account | Password | Role | Expected scope |
| --- | --- | --- | --- | --- |
| Admin overview | `admin@runtimeroasters.com` | `Hello@123` | `ADMIN` | Can view global overview. |
| Hanoi retail dashboard | `mgr.hn.hoankiem@runtimeroasters.com` | `Hello@123` | `STORE_MGR` | Store `11111111-1111-1111-1111-111111111101`. |
| Driver/logistics screen | `driver@runtimeroasters.com` | `Hello@123` | `DRIVER` | Driver simulation account for later tickets. |

Expected for this ticket: terminal evidence is authoritative. UI checks are optional because realtime/socket fanout is not RR-URG-02.

## Test Data

Use a new `order_id` for every rerun. If you reuse an `order_id`, idempotency should prevent duplicate processing.

Default one-time data:

| Field | Value |
| --- | --- |
| `order_id` | `11111111-2222-4333-8444-555555555556` |
| `store_id` | `11111111-1111-1111-1111-111111111101` |
| `traceid` | `22222222222222222222222222222222` |
| `sku` | `SL-ARABICA-ROASTED` |
| `quantity` | `1` |
| `amount` | `100` |

Contract-only future topic data:

| Field | Value |
| --- | --- |
| `shipment_id` | `33333333-4444-4555-8666-777777777777` |
| `harvest_id` | `HARVEST-RR-URG-02-001` |
| `farm_id` | `FARM-CAUDAT-001` |
| `warehouse_id` | `WAREHOUSE-HN-001` |
| `driver_id` | `22222222-2222-2222-2222-222222222201` |
| `vehicle_id` | `VEHICLE-DEMO-001` |

## Step 1: Start Infrastructure

From repository root:

```bash
docker compose -f deployments/docker-compose.dev.yaml up -d
```

Expected:

- `rr-kafka`, `rr-postgres`, `rr-valkey`, `rr-elasticsearch`, `rr-cassandra`, `rr-clickhouse`, `rr-signoz`, and `rr-otel-collector` are running or healthy.

Check:

```bash
docker compose -f deployments/docker-compose.dev.yaml ps
```

## Step 2: Optional Clean Demo State

Run this if you want the fixed IDs in this guide to work from a clean DB:

```bash
bash deployments/reset-demo-state.sh
```

Expected:

- Command ends with `RuntimeRoasters demo state reset.`

If you do not reset, change `order_id`, `event_id`, and `shipment_id` before publishing.

## Step 3: Start Required Services

Open five terminals.

Terminal A:

```bash
cd src/apps/payment-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs include `consumer listening` for `retail.order.created`.

Terminal B:

```bash
cd src/apps/warehouse-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs include `consumer listening` for `payment.simulated_completed`.
- Logs include `consumer listening` for `payment.completed`.

Terminal C:

```bash
cd src/apps/logistics-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs include `consumer listening` for `warehouse.stock.reserved`.

Terminal D:

```bash
cd src/apps/trace-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs include canonical topic consumers such as `retail.order.created`, `logistics.pickup.arrived_at_warehouse`, `notification.created`, and `socket.broadcast.requested`.

Terminal E:

```bash
cd src/apps/audit-service
set -a; source .env; set +a; GOCACHE=/private/tmp/runtime-roasters-go-cache go run ./cmd
```

Expected:

- Logs include canonical topic consumers such as `retail.order.created`, `logistics.pickup.arrived_at_warehouse`, `notification.created`, and `socket.broadcast.requested`.

Note: if `auth-service` is not running, you may see `Bootstrap auth sync failed` warnings. That does not block Kafka contract verification.

## Step 4: Publish Paid Order CloudEvent

From repository root:

```bash
ORDER_ID=11111111-2222-4333-8444-555555555556
EVENT_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
printf '%s\n' "${ORDER_ID}|{\"specversion\":\"1.0\",\"id\":\"evt-${ORDER_ID}\",\"source\":\"/services/retail-service\",\"type\":\"retail.order.created\",\"subject\":\"orders/${ORDER_ID}\",\"time\":\"${EVENT_TIME}\",\"datacontenttype\":\"application/json\",\"correlationid\":\"${ORDER_ID}\",\"orderid\":\"${ORDER_ID}\",\"storeid\":\"11111111-1111-1111-1111-111111111101\",\"traceid\":\"22222222222222222222222222222222\",\"data\":{\"event_id\":\"evt-${ORDER_ID}\",\"order_id\":\"${ORDER_ID}\",\"store_id\":\"11111111-1111-1111-1111-111111111101\",\"items\":[{\"sku\":\"SL-ARABICA-ROASTED\",\"quantity\":1}],\"total_amount\":100,\"payment_method\":\"STRIPE\",\"occurred_at\":\"${EVENT_TIME}\"}}" | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic retail.order.created --property parse.key=true --property key.separator='|'
```

Expected:

- Command exits `0`.
- Kafka may print `--property is deprecated`; this warning is acceptable.
- `EVENT_TIME` must be current. Payment-service intentionally ignores old order events when backfill is disabled.

## Step 5: Verify Payment

```bash
docker exec rr-postgres psql -U user -d payment_db -c "SELECT order_id, store_id, status, provider, amount FROM payments WHERE order_id = '11111111-2222-4333-8444-555555555556';"
```

Expected:

```text
order_id                              | store_id                              | status    | provider | amount
11111111-2222-4333-8444-555555555556 | 11111111-1111-1111-1111-111111111101 | SUCCEEDED | STRIPE   | 100.00
```

## Step 6: Verify Warehouse Reservation

```bash
docker exec rr-postgres psql -U user -d warehouse_db -c "SELECT message_id, event_type FROM inbox_events WHERE message_id = '11111111-2222-4333-8444-555555555556' OR event_type = 'payment.simulated_completed' ORDER BY processed_at DESC LIMIT 5;"
```

Expected:

- One row with `event_type = payment.simulated_completed`.
- One row with `message_id = 11111111-2222-4333-8444-555555555556` and `event_type = order_reserved`.

## Step 7: Verify Logistics Assignment

```bash
docker exec rr-postgres psql -U user -d logistics_db -c "SELECT id, order_id, destination_store_id, driver_id, status FROM shipments WHERE order_id = '11111111-2222-4333-8444-555555555556';"
```

Expected:

- Exactly one shipment.
- `destination_store_id = 11111111-1111-1111-1111-111111111101`.
- `driver_id` is not empty.
- `status = ASSIGNED`.

## Step 8: Verify Trace Projection

```bash
docker exec rr-postgres psql -U user -d trace_db -c "SELECT topic, order_id, store_id, shipment_id, trace_id FROM trace_events WHERE order_id = '11111111-2222-4333-8444-555555555556' ORDER BY occurred_at;"
```

Expected exactly these canonical topics:

```text
retail.order.created
payment.intent.created
payment.simulated_completed
warehouse.stock.reserved
logistics.delivery.assigned
```

Expected for every row:

- `order_id = 11111111-2222-4333-8444-555555555556`
- `store_id = 11111111-1111-1111-1111-111111111101`
- `trace_id = 22222222222222222222222222222222`

## Step 9: Verify Audit Projection

```bash
docker exec rr-postgres psql -U user -d audit_db -c "SELECT topic, partition_key, store_id FROM audit_logs WHERE partition_key = '11111111-2222-4333-8444-555555555556' ORDER BY occurred_at;"
```

Expected:

- Same five canonical topics as trace projection.
- `partition_key = 11111111-2222-4333-8444-555555555556`.
- `store_id = 11111111-1111-1111-1111-111111111101`.

## Step 10: Verify Future Topic Contract Projection

This does not execute the future pickup workflow. It verifies RR-URG-02's contract promise: trace/audit can accept canonical future topics today.

Publish a pickup milestone:

```bash
printf '%s\n' '33333333-4444-4555-8666-777777777777|{"specversion":"1.0","id":"evt-pickup-arrived-rr-urg-02","source":"/services/logistics-service","type":"logistics.pickup.arrived_at_warehouse","subject":"shipments/33333333-4444-4555-8666-777777777777","time":"2026-05-24T07:25:20Z","datacontenttype":"application/json","correlationid":"HARVEST-RR-URG-02-001","causationid":"evt-pickup-loading-rr-urg-02","traceid":"22222222222222222222222222222222","shipmentid":"33333333-4444-4555-8666-777777777777","harvestid":"HARVEST-RR-URG-02-001","farmid":"FARM-CAUDAT-001","warehouseid":"WAREHOUSE-HN-001","driverid":"22222222-2222-2222-2222-222222222201","vehicleid":"VEHICLE-DEMO-001","data":{"event_id":"evt-pickup-arrived-rr-urg-02","shipment_id":"33333333-4444-4555-8666-777777777777","harvest_id":"HARVEST-RR-URG-02-001","farm_id":"FARM-CAUDAT-001","warehouse_id":"WAREHOUSE-HN-001","driver_id":"22222222-2222-2222-2222-222222222201","vehicle_id":"VEHICLE-DEMO-001","status":"ARRIVED_AT_WAREHOUSE","occurred_at":"2026-05-24T07:25:20Z"}}' | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic logistics.pickup.arrived_at_warehouse --property parse.key=true --property key.separator='|'
```

Verify trace:

```bash
docker exec rr-postgres psql -U user -d trace_db -c "SELECT topic, shipment_id, harvest_id, farm_id, warehouse_id, driver_id, vehicle_id, trace_id FROM trace_events WHERE shipment_id = '33333333-4444-4555-8666-777777777777' ORDER BY occurred_at;"
```

Expected:

- One row with `topic = logistics.pickup.arrived_at_warehouse`.
- All business IDs are populated.
- `trace_id = 22222222222222222222222222222222`.

Verify audit:

```bash
docker exec rr-postgres psql -U user -d audit_db -c "SELECT topic, partition_key FROM audit_logs WHERE partition_key = '33333333-4444-4555-8666-777777777777' ORDER BY occurred_at;"
```

Expected:

- One row with `topic = logistics.pickup.arrived_at_warehouse`.
- `partition_key = 33333333-4444-4555-8666-777777777777`.

## Step 11: Verify Socket Contract Projection

Publish a socket broadcast request:

```bash
printf '%s\n' 'socket-rr-urg-02|{"specversion":"1.0","id":"evt-socket-rr-urg-02","source":"/services/notification-service","type":"socket.broadcast.requested","subject":"channels/warehouse","time":"2026-05-24T07:26:20Z","datacontenttype":"application/json","correlationid":"HARVEST-RR-URG-02-001","traceid":"22222222222222222222222222222222","warehouseid":"WAREHOUSE-HN-001","data":{"event_id":"evt-socket-rr-urg-02","channel":"warehouse","role":"WAREHOUSE_MGR","warehouse_id":"WAREHOUSE-HN-001","event_type":"logistics.pickup.arrived_at_warehouse","payload":{"shipment_id":"33333333-4444-4555-8666-777777777777"},"occurred_at":"2026-05-24T07:26:20Z"}}' | docker exec -i rr-kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic socket.broadcast.requested --property parse.key=true --property key.separator='|'
```

Verify trace:

```bash
docker exec rr-postgres psql -U user -d trace_db -c "SELECT topic, warehouse_id, trace_id FROM trace_events WHERE topic = 'socket.broadcast.requested' AND trace_id = '22222222222222222222222222222222' ORDER BY occurred_at DESC LIMIT 1;"
```

Expected:

- One row with `topic = socket.broadcast.requested`.
- `warehouse_id = WAREHOUSE-HN-001`.

## Step 12: Verify No Runtime Legacy Topics

```bash
rg "farm\\.harvest\\.events|warehouse\\.stock\\.updated|logistics\\.shipment\\.assigned|logistics\\.shipment\\.delivered" src deployments AGENT.md
```

Expected:

- No output.

Historical docs may still mention these names only as deprecated/reference context.

## RR-URG-02 Pass Checklist

- [ ] Paid order CloudEvent is accepted.
- [ ] Payment row is `SUCCEEDED`.
- [ ] Warehouse writes `payment.simulated_completed` and `order_reserved`.
- [ ] Logistics creates one `ASSIGNED` shipment.
- [ ] Trace DB shows five current paid-order canonical topics.
- [ ] Audit DB shows five current paid-order canonical topics.
- [ ] Future pickup topic projects into trace and audit.
- [ ] Socket broadcast request topic projects into trace.
- [ ] All projected events carry business IDs separately from `trace_id`.
- [ ] Runtime source/config has no legacy topic usage.
