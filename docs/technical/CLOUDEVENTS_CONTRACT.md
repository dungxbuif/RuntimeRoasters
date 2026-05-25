# CloudEvents Contract

Date: 2026-05-24

## Purpose

Runtime Roasters uses Kafka for domain events between services. All new domain events must be published as CloudEvents JSON using CloudEvents `specversion: "1.0"`.

This contract exists to keep SAGA, trace, audit, UI timeline, and demo evidence consistent across services.

## Envelope

Every Kafka domain message value must be a CloudEvent JSON object:

```json
{
  "specversion": "1.0",
  "id": "evt_01HY...",
  "type": "retail.order.created",
  "source": "/services/retail-service",
  "subject": "orders/ORDER-001",
  "time": "2026-05-24T10:00:00Z",
  "datacontenttype": "application/json",
  "correlationid": "ORDER-001",
  "causationid": "evt_previous",
  "traceid": "8991cbb271c04df9b233ef5f98d58a8c",
  "orderid": "ORDER-001",
  "storeid": "STORE-HCM-01",
  "data": {
    "order_id": "ORDER-001",
    "store_id": "STORE-HCM-01",
    "items": []
  }
}
```

Required CloudEvents attributes:
- `id`: event ID, globally unique enough for idempotency.
- `type`: canonical event type. In this project it matches the Kafka topic.
- `source`: service URI, for example `/services/payment-service`.
- `subject`: primary business entity path, for example `orders/{order_id}`.
- `time`: event occurrence time.
- `datacontenttype`: must be `application/json`.
- `data`: typed business payload.

Required extension:
- `correlationid`: durable business flow ID used to group related events.

Optional extensions:
- `causationid`: previous CloudEvent `id` that caused this event.
- `traceid`: OpenTelemetry trace ID, when a span exists.
- Business IDs: `orderid`, `paymentid`, `storeid`, `shipmentid`, `harvestid`, `batchid`, `farmid`, `warehouseid`, `driverid`, `vehicleid`, `orgid`.

Extension names must stay lowercase alphanumeric because CloudEvents extension names are intentionally constrained.

## ID Rules

Business IDs and OpenTelemetry trace IDs are different concepts:

- Business IDs are durable product/domain IDs: `order_id`, `shipment_id`, `harvest_id`, `batch_id`, `store_id`, `farm_id`, `warehouse_id`, `driver_id`, `vehicle_id`.
- `traceid` is a technical observability ID used to find spans in SigNoz/ClickHouse.
- Never use `traceid` as a business identifier.
- Trace-service and audit-service must query business history by business IDs, not by `traceid`.
- `traceid` may be stored alongside events for debugging and UI correlation links.
- Derived CloudEvents must preserve the incoming `traceid` extension when one exists. This keeps demo flows readable even when a direct Kafka/client simulator does not provide W3C `traceparent` transport headers.

Recommended `correlationid`:
- Retail paid-order flow: `order_id`.
- Farm harvest flow: `harvest_id`.
- Logistics-only location flow: `shipment_id` when present, otherwise `driver_id`.

Recommended `causationid`:
- Root events omit it.
- Derived events use the incoming CloudEvent `id`.

## Canonical Topics

Canonical topics:

| Event type / Kafka topic | Producer | Primary consumers |
| --- | --- | --- |
| `farm.harvest.created` | Farm | Warehouse, Trace, Audit |
| `retail.order.created` | Retail | Payment, Trace, Audit |
| `payment.intent.created` | Payment | Retail, Trace, Audit |
| `payment.simulated_completed` | Payment | Warehouse, Retail, Trace, Audit |
| `payment.completed` | Payment webhook/provider | Warehouse, Retail, Trace, Audit |
| `payment.failed` | Payment | Retail, Trace, Audit |
| `payment.refunded` | Payment | Retail, Trace, Audit |
| `warehouse.stock.reserved` | Warehouse | Logistics, Retail, Trace, Audit |
| `warehouse.stock.reservation_failed` | Warehouse | Payment, Retail, Trace, Audit |
| `warehouse.inventory.updated` | Warehouse | Logistics, Trace, Audit |
| `warehouse.pickup.requested` | Warehouse | Logistics, Farm UI, Trace, Audit |
| `warehouse.pickup.received` | Warehouse | Farm UI, Trace, Audit |
| `warehouse.dispatch.requested` | Warehouse | Logistics, Warehouse UI, Trace, Audit |
| `warehouse.intake.created` | Warehouse | Farm UI, Trace, Audit |
| `logistics.delivery.assigned` | Logistics | Retail, Trace, Audit |
| `logistics.delivery.departed` | Logistics | Retail, Warehouse UI, Trace, Audit |
| `logistics.delivery.arrived_at_store` | Logistics | Retail, Warehouse UI, Trace, Audit |
| `logistics.delivery.driver_confirmed` | Logistics | Retail, Warehouse UI, Trace, Audit |
| `logistics.delivery.completed` | Logistics | Retail, Trace, Audit |
| `logistics.pickup.assigned` | Logistics | Farm UI, Warehouse UI, Trace, Audit |
| `logistics.pickup.departed` | Logistics | Farm UI, Warehouse UI, Trace, Audit |
| `logistics.pickup.arrived_at_farm` | Logistics | Farm UI, Warehouse UI, Trace, Audit |
| `logistics.pickup.loading_confirmed` | Logistics | Farm UI, Warehouse UI, Trace, Audit |
| `logistics.pickup.return_started` | Logistics | Farm UI, Warehouse UI, Trace, Audit |
| `logistics.pickup.arrived_at_warehouse` | Logistics | Warehouse, Farm UI, Trace, Audit |
| `logistics.pickup.completed` | Logistics | Warehouse, Farm UI, Trace, Audit |
| `logistics.driver.return_started` | Logistics | Warehouse UI, Retail/Farm UI, Trace, Audit |
| `logistics.driver.return_completed` | Logistics | Warehouse UI, Retail/Farm UI, Trace, Audit |
| `logistics.driver.returned_to_base` | Logistics | Warehouse UI, Retail UI, Trace, Audit |
| `logistics.gps.updated` | Logistics | Trace, Audit |
| `logistics.shipment.status_changed` | Logistics | Trace, Audit, UI |
| `notification.created` | Notification | UI, Trace, Audit |
| `notification.acknowledged` | Notification | UI, Trace, Audit |
| `socket.broadcast.requested` | Backend services | Socket service, Trace, Audit |

RR-URG-02 defines the contract for all canonical topics above. Some producers/state machines are implemented by later urgent tickets, but trace/audit must already accept the full list.

Legacy topics must not be used by runtime code:
- `farm.harvest.events`
- `warehouse.stock.updated`
- `logistics.shipment.assigned`
- `logistics.shipment.delivered`

These legacy names may appear only in documentation sections that explicitly mark them as deprecated or forbidden. They must not appear in runtime code, service config, `.env.example`, seed scripts, tests, or demo evidence expectations.

## Payload Rules

`data` contains business payload only. Do not duplicate generic envelope fields in `data` unless an existing typed payload already uses them for compatibility.

Allowed examples:
- `order_id`, `store_id`, `items`, `total_amount`
- `payment_id`, `provider`, `provider_ref`, `amount`, `currency`
- `shipment_id`, `driver_id`, `lat`, `long`
- `harvest_id`, `coffee_type`, `origin_code`, `quantity`

Do not place authorization scope only inside JSONB payload. Services that need durable query or scoping must store typed columns such as `store_id`, `order_id`, or `shipment_id`.

## Service Responsibilities

Producers:
- Build events through `pkg/events.NewCloudEvent`.
- Set canonical `type` and publish to the same Kafka topic.
- Set all known business ID extensions.
- Preserve OTel propagation through Kafka headers.
- When producing a derived event from an incoming CloudEvent, pass through the incoming `traceid` extension as event metadata.

Consumers:
- Parse with `pkg/events.ParseCloudEvent`.
- Reject invalid CloudEvents instead of silently accepting legacy JSON.
- Use `DataAs[T]` for typed payloads.
- Use CloudEvents extensions for idempotency, trace, audit, and read model IDs.

Trace-service:
- Stores raw CloudEvent JSON.
- Uses CloudEvent `type` for timeline state.
- Stores typed business IDs into queryable columns.
- Stores `traceid` for observability correlation only.

Audit-service:
- Stores raw CloudEvent JSON immutably.
- Partitions by business ID extension, in priority order: `orderid`, `shipmentid`, `harvestid`, `batchid`, `storeid`, `farmid`, then event `id`.
- Hashes the raw CloudEvent payload plus Kafka metadata.

## Examples

Retail order:

```json
{
  "specversion": "1.0",
  "id": "evt-order-1",
  "type": "retail.order.created",
  "source": "/services/retail-service",
  "subject": "orders/order-1",
  "time": "2026-05-24T10:00:00Z",
  "datacontenttype": "application/json",
  "correlationid": "order-1",
  "traceid": "8991cbb271c04df9b233ef5f98d58a8c",
  "orderid": "order-1",
  "storeid": "store-1",
  "data": {
    "event_id": "evt-order-1",
    "order_id": "order-1",
    "store_id": "store-1",
    "items": [{ "sku": "SL-ARABICA-ROASTED", "quantity": 1 }],
    "total_amount": 100,
    "payment_method": "STRIPE",
    "occurred_at": "2026-05-24T10:00:00Z"
  }
}
```

Simulated payment completion:

```json
{
  "specversion": "1.0",
  "id": "evt-payment-1",
  "type": "payment.simulated_completed",
  "source": "/services/payment-service",
  "subject": "orders/order-1",
  "time": "2026-05-24T10:00:05Z",
  "datacontenttype": "application/json",
  "correlationid": "order-1",
  "causationid": "evt-order-1",
  "traceid": "8991cbb271c04df9b233ef5f98d58a8c",
  "orderid": "order-1",
  "paymentid": "payment-1",
  "storeid": "store-1",
  "data": {
    "event_id": "evt-payment-1",
    "payment_id": "payment-1",
    "order_id": "order-1",
    "store_id": "store-1",
    "provider": "STRIPE",
    "amount": 100,
    "currency": "USD",
    "occurred_at": "2026-05-24T10:00:05Z"
  }
}
```

Logistics GPS:

```json
{
  "specversion": "1.0",
  "id": "evt-gps-1",
  "type": "logistics.gps.updated",
  "source": "/services/logistics-service",
  "subject": "drivers/driver-1",
  "time": "2026-05-24T10:00:20Z",
  "datacontenttype": "application/json",
  "correlationid": "shipment-1",
  "traceid": "8991cbb271c04df9b233ef5f98d58a8c",
  "shipmentid": "shipment-1",
  "driverid": "driver-1",
  "storeid": "store-1",
  "data": {
    "event_id": "evt-gps-1",
    "driver_id": "driver-1",
    "shipment_id": "shipment-1",
    "store_id": "store-1",
    "lat": 10.7769,
    "long": 106.7009,
    "occurred_at": "2026-05-24T10:00:20Z"
  }
}
```

## Agent Checklist

When adding or changing an event:
- Use `pkg/events.NewCloudEvent`; do not hand-build the envelope.
- Add or reuse a canonical topic constant in `pkg/events/contracts.go`.
- Update producer, consumer, trace, audit, `.env.example`, and docs together.
- Include `correlationid` and all relevant business ID extensions.
- Keep OTel `traceid` separate from business IDs.
- Preserve incoming `traceid` on derived CloudEvents; do not replace it with a business ID.
- Add or update tests proving the event parses as CloudEvent and the consumer uses typed `data`.
