# RR-URG-01 Trace E2E Evidence

Date: 2026-05-23

Scope: live browser login through Client App, authenticated Gateway calls through KrakenD, Kafka-backed paid-order SAGA projection into trace-service, harvest event consumption by warehouse, and SigNoz/ClickHouse waterfall evidence.

## Playwright Evidence

Command:

```bash
npx playwright test e2e/manual-live-evidence.spec.ts --project=chromium
```

Result:

```text
1 passed (17.2s)
```

Evidence payload:

```json
{
  "traceparent": "00-8991cbb271c04df9b233ef5f98d58a8c-7c5e21782914443d-01",
  "expectedTraceID": "8991cbb271c04df9b233ef5f98d58a8c",
  "subject": "60855b7f-8dd2-4272-9f7d-08387c32f826",
  "role": "ADMIN",
  "farmID": "8",
  "harvestID": "15",
  "storeID": "11111111-1111-1111-1111-111111111105",
  "orderID": "68669cad-06fb-4b6b-8583-6d77801b727b",
  "traceDocumentFound": true,
  "traceIDs": ["8991cbb271c04df9b233ef5f98d58a8c"],
  "topics": [
    "retail.order.created",
    "payment.intent.created",
    "payment.completed",
    "warehouse.stock.reserved",
    "logistics.shipment.assigned"
  ]
}
```

## Trace-Service DB Evidence

Command:

```bash
docker exec rr-postgres psql -U user -d trace_db -c "SELECT topic, trace_id FROM trace_events WHERE trace_id = '8991cbb271c04df9b233ef5f98d58a8c' ORDER BY occurred_at;"
```

Result:

```text
retail.order.created        | 8991cbb271c04df9b233ef5f98d58a8c
payment.intent.created      | 8991cbb271c04df9b233ef5f98d58a8c
payment.completed           | 8991cbb271c04df9b233ef5f98d58a8c
warehouse.stock.reserved    | 8991cbb271c04df9b233ef5f98d58a8c
logistics.shipment.assigned | 8991cbb271c04df9b233ef5f98d58a8c
```

## Elasticsearch Evidence

Trace document exists in `coffee_traceability` for order `68669cad-06fb-4b6b-8583-6d77801b727b` and includes:

- `trace_ids`: `["8991cbb271c04df9b233ef5f98d58a8c"]`
- timeline entries for order created, payment intent, payment completed, stock reserved, and shipment assigned.

## SigNoz / ClickHouse Evidence

SigNoz UI:

- `http://localhost:3301`
- Screenshot: `docs/business/sprint-emergency-final-demo/evidence/signoz-ui.png`

ClickHouse query:

```sql
SELECT trace_id, serviceName, name, kind_string, db_name, db_operation
FROM signoz_traces.distributed_signoz_index_v3
WHERE trace_id = '8991cbb271c04df9b233ef5f98d58a8c'
ORDER BY timestamp
LIMIT 100;
```

Result includes:

```text
farm-service       POST                                   Server
farm-service       runtime.farm.v1.FarmService/CreateFarm Client
farm-service       insert farms                           Client insert
retail-service     GET /v1/stores                         Server
retail-service     POST /v1/orders                        Server
retail-service     insert orders                          Client insert
retail-service     kafka.produce retail.order.created     Producer
payment-service    kafka.consume retail.order.created     Consumer
payment-service    kafka.produce payment.completed        Producer
warehouse-service  kafka.consume payment.completed        Consumer
warehouse-service  update inventories                     Client update
warehouse-service  kafka.produce warehouse.stock.reserved Producer
logistics-service  kafka.consume warehouse.stock.reserved Consumer
logistics-service  insert shipments                       Client insert
logistics-service  kafka.produce logistics.shipment.assigned Producer
trace-service      kafka.consume logistics.shipment.assigned Consumer
```

Aggregate check:

```text
spans: 108
services: 7
```

## Harvest/Warehouse Evidence

Playwright also created harvest `15`. The warehouse consumer accepted the event path during the same live run.

## Kafka Evidence

Command:

```bash
docker exec rr-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

Result excerpt:

```text
farm.harvest.events
retail.order.created
payment.intent.created
payment.completed
warehouse.stock.reserved
logistics.shipment.assigned
```

## Verdict

PASS. One authenticated browser-driven order flow produced a single trace ID through KrakenD, retail-service, Kafka, payment-service, warehouse-service, logistics-service, trace-service projection, Elasticsearch, and SigNoz/ClickHouse. The same run also verified harvest event consumption into warehouse.
