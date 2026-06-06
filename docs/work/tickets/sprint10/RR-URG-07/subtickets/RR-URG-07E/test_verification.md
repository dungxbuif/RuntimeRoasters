---
artifact_type: test_verification
id: TEST-RR-URG-07E
status: complete
owner: ai
updated: 2026-06-06
---

# RR-URG-07E Final Verification

## Acceptance Mapping

| Acceptance | Evidence | Result |
| --- | --- | --- |
| Seeded menu, lots, and sold cups | Fresh Retail migration + seed twice | pass |
| Anonymous sale without payment | Real `POST /v1/public/retail/sales` returned `201` | pass |
| UI-issued product ID persisted | `RR-CUP-PLATFORM-070E-2` in sale item and public response | pass |
| Inventory and movement updated | PostgreSQL showed `quantity_delta=-18` and updated lot balance | pass |
| Event reaches trace service | Outbox `COMPLETED`; trace DB stored `retail.sale.completed` | pass |
| Public trace resolves same cup | Real Trace HTTP returned `200` for `RR-CUP-PLATFORM-070E-2` | pass |
| Invoice QR and public UI | Playwright 2/2 plus three screenshots | pass |
| Concurrent no-oversell | Live PostgreSQL produced one success, one conflict, zero balance | pass |

## Commands

- Backend: `go test` and `go vet` for Retail, Trace, and events: pass.
- Frontend: `npm test`, `npm run lint`, `npm run build`: pass.
- E2E: `npm run test:e2e:rr-urg-07`: pass, 2 tests.
- Gateway: `jq empty deployments/krakend/krakend.json`: pass.
- Hygiene: `git diff --check`: pass.

## Platform Notes

- PostgreSQL 16 databases: `rr_urg_07e_retail`, `rr_urg_07e_trace`.
- Kafka and Hydra were started from the development Compose stack.
- Elasticsearch was deliberately disabled to prove PostgreSQL fallback.
- Temporary Retail/Trace service processes were stopped after evidence.

