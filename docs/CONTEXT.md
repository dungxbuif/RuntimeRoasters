---
artifact_type: context
id: CONTEXT
status: active
owner: shared
---
# Session Context & Development State

Last Updated: 2026-06-05

## Current Git Snapshot
**Checked:** 2026-06-05
**Active branch:** `rr-urg-07-trace-public-qr`
**HEAD:** `7663268 RR-URG-07 start trace public QR`

### Active Dirty Changes
- `src/apps/client-app/src/constants/casbin.ts`: frontend Casbin matcher now allows direct role policy matching (`r.sub == p.sub`) as well as inherited role matching (`g(r.sub, p.sub)`), fixing frontend unit failures where role policies such as `STORE_MGR` did not evaluate.
- `src/apps/socket-service/internal/app/app_test.go`: added socket-service route integration coverage for private stream ticket requirement, public WebSocket access, ticketed private WebSocket access, and one-time ticket rejection.
- `src/apps/client-app/e2e/operational_realtime_awareness.spec.ts`: added flow-oriented Playwright coverage for Warehouse/Store dashboard notification rendering and role-scope isolation with mocked gateway auth/API contracts; passed through production Next runtime.
- `src/apps/client-app/package.json`: added `start:e2e` and `test:e2e:realtime-notifications` scripts.
- `src/apps/client-app/playwright.config.ts`: added fail-fast E2E guardrails: production webServer, `127.0.0.1` baseURL, 30s test timeout, 10s navigation timeout, 5s action/expect timeout, line+HTML reporters.
- `src/apps/socket-service/internal/platform/doc.go` and `src/apps/socket-service/internal/platform/platform_live_test.go`: added opt-in live platform test for Kafka -> socket-service -> Redis/Valkey notification persistence.
- `src/apps/warehouse-service/internal/app/app.go`: removed duplicate stub registration for `/v1/warehouse/dispatch-requests` routes so warehouse-service can boot for platform verification.
- `docs/work/tickets/history/sprint-emergency-final-demo/RR-URG-06-realtime-notifications.md`: updated verification evidence and remaining manual visual demo gap.

### Outstanding Verification
- RR-URG-06 targeted backend passed:
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/socket-service/... ./apps/auth-service/internal/infrastructure/casbin/...`
- Frontend checks passed:
  - `npm test`, `npm run lint`, `npm run build`
- KrakenD config parse passed:
  - `node -e "JSON.parse(require('fs').readFileSync('deployments/krakend/krakend.json','utf8')); console.log('krakend json ok')"`
- Playwright browser dependency was installed:
  - `npx playwright install chromium`
- Operational realtime awareness E2E passed:
  - `npm run test:e2e:realtime-notifications` passed 3 Playwright production-runtime flow tests in 9.2s.
  - Covered flows: warehouse outbound dispatch notification, store incoming delivery notification, and store-manager isolation.
- Platform live test passed:
  - Kafka `warehouse.dispatch.requested` event was consumed by running socket-service and persisted as a role/warehouse-scoped Redis notification.
- Harness matrix updated:
  - RR-URG-06: unit yes, integration yes, e2e yes, platform yes.
- Manual visual review still pending:
  - Full two-browser role-session demo with real logged-in users has not been run.

### Next Task
- Next ticket is RR-URG-07: Trace-Service Read Model, Live History, And Public QR Trace Demo.
- Current discussion corrected RR-URG-07 demo seed/QR model: QR semantics should represent the retail sold unit/receipt item, such as a cup/order item sold by a store.
- Current code/system drift identified before implementation:
  - `src/apps/client-app/src/app/trace/[code]/page.tsx` is still static/mock (`DEMO_JOURNEY`).
  - trace-service `TraceDocument` and Elasticsearch read model do not yet model `trace_code` or public sold-unit lookup.
  - trace-service public routes currently expose public topology only.
- Needed RR-URG-07 implementation direction:
  - Add chain-wide retail products/menu, retail inventory lots, retail sales/invoices.
  - Seed store inventory lots with upstream farm/batch/warehouse lineage.
  - Add a sold-item public trace read model keyed by UI-issued `cup_id` exposed as `trace_code`.
  - Replace static public trace UI with API-backed fetch/rendering.

## Phase 1: System Bootstrap & Seeding
**Status:** `✅ COMPLETED`

## Phase 2: Logistics Real-time Tracking & Saga Simulator
**Status:** `✅ COMPLETED`

### Technical Achievements:
1.  **Identity & RBAC**: Unified `INTERNAL_SECRET`, strict `ADMIN` read-only domain access, upgraded frontend Casbin matcher.
2.  **Resource Management**: Admin UI for registering Warehouses and Retail Stores, Fleet Assignment Modal in Warehouse Ops.
3.  **Performance & Stability**: `next/dynamic` for heavy visual components, Gorm `AutoMigrate` fixes, idempotent seeding.

## Current Context & System Shape
- **Client App Port**: `http://localhost:3000`
- **KrakenD Gateway**: `http://localhost:8081`
- **Internal Shared Secret**: Standardized across all `.env` files.

## Next Steps
- [ ] Implement actual Valkey Geo tracking in `logistics-service`.
- [ ] Enhance Saga Monitor with real-time SSE updates.
- [ ] Move to Phase 3: System Reliability & Chaos Testing.

## Active Work (Harness v1)
- **Harness Migration:** Upgrading to Harness v1 (Pure Markdown) - **COMPLETED**.
- **Next Up:** RR-URG-07 (Trace-Service Read Model, Live History, And Public QR Trace Demo).
