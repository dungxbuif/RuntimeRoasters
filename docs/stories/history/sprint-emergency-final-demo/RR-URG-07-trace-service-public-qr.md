# RR-URG-07: Trace-Service Read Model, Live History, And Public QR Trace Demo

## Priority

P2. Depends on RR-URG-01 and RR-URG-02. Public QR demo depends on prepared data from RR-URG-08.

## Problem

The final showcase needs a public page where users generate/open QR codes for prepared products and see real traceability data. This must be backed by trace-service and Elasticsearch, not static UI mock data.

## Scope

- Extend trace-service event projection.
- Store business trace documents in Elasticsearch.
- Optionally store short-lived trace-service live history in Cassandra.
- Add public-safe trace lookup by `trace_code`.
- Prepare QR trace API contract.

## Implementation Details

### 1. Trace Document Model

Trace document should include:

- `trace_code`
- `origin_batch_id`
- `harvest_id`
- `intake_id`
- `production_batch_id`
- `inventory_lot_id`
- `order_id`
- `shipment_ids`
- `farm`
- `warehouse`
- `store`
- `product`
- `milestones`
- `service_participation`
- `trace_ids`
- `correlation_id`
- `public_safe`

Milestones:

- harvest created
- pickup requested
- pickup assigned
- driver departed
- arrived at farm
- pickup/loading confirmed
- return started
- arrived warehouse
- intake created
- processing started/finalized
- inventory stocked
- paid order created
- payment completed/simulated
- stock reserved
- delivery assigned
- arrived store
- delivered
- driver returned to base

Checklist:

- [ ] Elasticsearch mapping updated.
- [ ] Trace projection handles every milestone event.
- [ ] Public-safe fields separated from internal fields.
- [ ] Trace document query by `trace_code`.
- [ ] Trace document query by business IDs for authenticated users.

### 2. Trace-Service Live History

Decision:

- Keep live history in trace-service.
- Do not create a separate monitor-service.

If implemented:

- trace-service receives OTel spans or consumes selected events.
- Cassandra table stores short-lived `trace_history`.
- TTL recommended: 48 hours.
- SSE endpoint streams sanitized/public or role-scoped/private updates.

Candidate endpoints:

- `GET /v1/trace/stream?trace_id=...`
- `GET /v1/trace/stream?correlation_id=...`

Checklist:

- [ ] Live history table owned by trace-service.
- [ ] Audit Cassandra tables remain separate.
- [ ] TTL configured.
- [ ] Public stream sanitized.

### 3. Public QR APIs

Candidate endpoints:

- `GET /v1/public/trace-products`
- `GET /v1/public/trace/:trace_code`

Client behavior:

- Public page calls `trace-products`.
- Client renders QR codes locally from returned public URLs.
- QR opens `/trace/public/:trace_code`.
- Public trace page calls trace-service and renders real trace document.

Checklist:

- [ ] Public product list contains prepared seed products.
- [ ] QR URL stable.
- [ ] No auth needed for public-safe trace.
- [ ] Private data removed.

### 4. OTel Enrichment

Use OTel for technical detail, not as the main business story:

- services touched.
- operation names.
- latency.
- connected trace IDs.

Checklist:

- [ ] UI can show business journey without OTel.
- [ ] UI can enrich journey with OTel if available.
- [ ] Missing OTel data does not break public trace page.

## Acceptance Criteria

1. Trace-service projects full Farm -> Warehouse -> Retail journey.
2. Elasticsearch can return trace document by public `trace_code`.
3. Public QR page opens real trace document.
4. Public trace does not expose sensitive data.
5. Trace document includes driver return milestones.
6. OTel/service participation appears when available.

## Test Checklist

- [ ] Unit: event projection for each milestone.
- [ ] Unit: public sanitizer.
- [ ] Integration: seed product trace projected into Elasticsearch.
- [ ] Integration: public trace endpoint returns expected journey.
- [ ] Manual: QR page opens trace without login.
- [ ] Manual: trace page shows real prepared data, not mock data.
