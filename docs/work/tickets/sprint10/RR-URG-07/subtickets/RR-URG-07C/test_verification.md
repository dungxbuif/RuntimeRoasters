---
artifact_type: test_verification
id: TEST-RR-URG-07C
status: complete
owner: ai
updated: 2026-06-06
---

# RR-URG-07C Test Verification

- Pass: `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/trace-service/...`
- Pass: public composer/sanitizer, dependency error, 404, sold-item list, and
  no-auth Gin route tests.
- Pass: shared backend suite and `go vet` for Retail, Trace, and events.
- Pass: live Retail and Trace services on ports `18084`/`18087` against fresh
  PostgreSQL databases.
- Pass: `GET /v1/public/trace/RR-CUP-HK-0001` returned a public-safe seeded
  trace with invoice, product, store, lot, batch, harvest, and warehouse IDs.
- Pass: Kafka consumer persisted `retail.sale.completed` in trace PostgreSQL
  after configuration was reconciled.
- Unknown product behavior is covered by unit/HTTP integration and frontend E2E.

