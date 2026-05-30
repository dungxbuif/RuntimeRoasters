# Research Notes

## Local Evidence Reviewed

- `docs/HARNESS.md`
- `docs/FEATURE_INTAKE.md`
- `docs/ARCHITECTURE.md`
- `docs/FEEDBACK_WORKFLOW.md`
- `docs/product/GUIDE.md`
- `docs/product/TECH.md`
- `docs/product/SPEC.md`
- `docs/product/standards/BOOTSTRAP.md`
- `src/apps/client-app/src/services/logistics.service.ts`
- `src/apps/client-app/src/app/(dashboard)/dashboard/logistics/page.tsx`
- `src/apps/client-app/src/app/(dashboard)/dashboard/driver/page.tsx`
- `src/apps/client-app/src/app/(dashboard)/dashboard/traceability/page.tsx`
- `src/apps/client-app/src/app/(dashboard)/dashboard/finance/page.tsx`
- `src/apps/logistics-service/internal/app/app.go`
- `src/apps/logistics-service/internal/usecase/service.go`
- `src/apps/trace-service/internal/usecase/service.go`
- `src/apps/trace-service/internal/search/elasticsearch.go`

## Findings

- Logistics UI fetches locations and shipments from backend APIs, but route coordinates come from `/data/routes.json`.
- Logistics dashboard has fallback demo shipments when the API returns no shipment data.
- Driver client posts simulated GPS points to the backend and calls real shipment milestone APIs.
- Logistics backend validates driver/shipment scope, stores accepted driver location in Valkey GEO data, updates vehicle coordinates, and publishes `logistics.gps.updated`.
- Trace service parses CloudEvents, extracts business IDs, writes `trace_events`, rebuilds `trace_documents`, and upserts Elasticsearch documents.
- Traceability UI currently requires the user to enter a known Batch ID or Shipment ID before a document is loaded.
- Finance UI reads real payment rows and sends signed demo Stripe-compatible webhook payloads through the configured payment API.
- Bootstrap docs already define historical data simulation through direct read-model insertion for immediate dashboard readiness.

## Harness Note

`scripts/harness query matrix` initially failed because `harness.db` did not exist. `scripts/harness init` created the local durable database before recording the intake.
