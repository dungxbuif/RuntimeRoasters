# Session Context & Development State

Last Updated: 2026-06-05

## Current Git Snapshot
**Checked:** 2026-06-05
**Active branch:** `rr-urg-06-realtime-notifications`
**HEAD:** `54faa90 RR-URG-06 realtime notifications`

### Active Dirty Changes
- `src/apps/client-app/src/constants/casbin.ts`: frontend Casbin matcher now allows direct role policy matching (`r.sub == p.sub`) as well as inherited role matching (`g(r.sub, p.sub)`), fixing frontend unit failures where role policies such as `STORE_MGR` did not evaluate.
- `src/apps/socket-service/internal/app/app_test.go`: added socket-service route integration coverage for private stream ticket requirement, public WebSocket access, ticketed private WebSocket access, and one-time ticket rejection.
- `src/apps/client-app/e2e/operational_realtime_awareness.spec.ts`: added flow-oriented Playwright coverage for Warehouse/Store dashboard notification rendering and role-scope isolation with mocked gateway auth/API contracts; passed through production Next runtime.
- `src/apps/client-app/package.json`: added `start:e2e` and `test:e2e:realtime-notifications` scripts.
- `src/apps/client-app/playwright.config.ts`: added fail-fast E2E guardrails: production webServer, `127.0.0.1` baseURL, 30s test timeout, 10s navigation timeout, 5s action/expect timeout, line+HTML reporters.
- `src/apps/socket-service/internal/platform/doc.go` and `src/apps/socket-service/internal/platform/platform_live_test.go`: added opt-in live platform test for Kafka -> socket-service -> Redis/Valkey notification persistence.
- `src/apps/warehouse-service/internal/app/app.go`: removed duplicate stub registration for `/v1/warehouse/dispatch-requests` routes so warehouse-service can boot for platform verification.
- `docs/stories/history/sprint-emergency-final-demo/RR-URG-06-realtime-notifications.md`: updated verification evidence and remaining manual visual demo gap.
- `docs/CONTEXT.md`: updated with this RR-URG-06 verification context.

### Outstanding Verification
- RR-URG-06 targeted backend passed:
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/socket-service/... ./apps/auth-service/internal/infrastructure/casbin/...`
- Frontend checks passed:
  - `npm test`
  - `npm run lint`
  - `npm run build`
- KrakenD config parse passed:
  - `node -e "JSON.parse(require('fs').readFileSync('deployments/krakend/krakend.json','utf8')); console.log('krakend json ok')"`
- Playwright browser dependency was installed:
  - `npx playwright install chromium`
- Operational realtime awareness E2E passed:
  - `npm run test:e2e:realtime-notifications` passed 3 Playwright production-runtime flow tests in 9.2s when run with approved escalation to allow binding `127.0.0.1:3000`.
  - Covered flows: warehouse outbound dispatch notification, store incoming delivery notification, and store-manager isolation from warehouse operations notifications.
  - The non-escalated run failed fast at webServer readiness because sandbox blocked `next start -H 127.0.0.1 -p 3000` with `listen EPERM`.
- Platform live test passed:
  - Started required local infra/services, fixed warehouse duplicate route boot panic, and ran `RUN_PLATFORM_TESTS=1 VALKEY_ADDR=127.0.0.1:6379 KAFKA_BROKERS=127.0.0.1:9094 go test -tags=platform ./apps/socket-service/internal/platform -run TestWarehouseDispatchEventCreatesLiveNotification -count=1 -v`.
  - Result: Kafka `warehouse.dispatch.requested` event was consumed by running socket-service and persisted as a role/warehouse-scoped Redis notification.
- Harness matrix updated:
  - RR-URG-06: unit yes, integration yes, e2e yes, platform yes.
- Manual visual review still pending:
  - Full two-browser role-session demo with real logged-in users has not been run.
  - `npm run dev` is still not suitable for E2E evidence because it emitted repeated `Can't resolve 'tailwindcss' in '/Users/dungxbuif/workspace/RuntimeRoasters/src/apps'`; production `next start` via Playwright webServer is the current E2E route.

### Next Task
- RR-URG-06 verification is complete in Harness for unit, integration, E2E, and platform evidence.
- Continue RR-URG-07 trace-service/public QR work after committing or carrying forward the RR-URG-06 verification/test updates.
- Current recent commits:
  - `54faa90 RR-URG-06 realtime notifications`
  - `55700d5 RR-URG-05 paid order fulfillment`

## Phase 1: System Bootstrap & Seeding
**Status:** `✅ COMPLETED`

## Phase 2: Logistics Real-time Tracking & Saga Simulator
**Status:** `✅ COMPLETED`

### Technical Achievements:
1.  **Identity & RBAC (Admin & Policy Gaps)**:
    *   **Unified `INTERNAL_SECRET`**: Standardized the internal shared secret across all 10+ microservices and infrastructure (Hydra, Identity Proxy). Fixed the `401 Unauthorized` issues during inter-service seeding propagation.
    *   **Strict RBAC**: Restricted `ADMIN` to read-only access for domain write actions (Harvests, Orders) while maintaining full system management.
    *   **Regex Matcher**: Upgraded frontend Casbin matcher to support `regexMatch`, enabling complex path policies (e.g., `/v1/warehouse/.*`).
    *   **Page-Level Protection**: Wrapped `Harvests` and `Create Order` pages with `CasbinGuard` to prevent unauthorized direct URL access.

2.  **Resource Management & Fleet Assignment**:
    *   **Admin UI for Infrastructure**: Implemented a **Resource Management** page for Administrators to register new Warehouses and Retail Stores with real API integration.
    *   **Fleet Assignment Modal**: Added a modal in **Warehouse Ops** for Warehouse Managers to explicitly assign Drivers and Vehicles during the dispatch flow.
    *   **Logistics Map Audit**: Converted the map to **Light Mode** (Voyager tiles) and implemented route filtering (only show routes with active shipments).

3.  **Performance & Technical Stability**:
    *   **Performance Optimization**: Implemented `next/dynamic` for heavy visual components (`LogisticsMap`, `ArchitectureDiagramCanvas`) to reduce initial bundle size.
    *   **Migration Robustness**: Fixed Gorm `AutoMigrate` failures related to constraint naming and foreign key dependency order in `warehouse-service` and `retail-service`.
    *   **Idempotent Seeding**: Upgraded seeder script to handle nullable `order_id` in shipments and implemented `FirstOrCreate` for resource creation to avoid 500 errors on retries.
    *   **Gateway (KrakenD)**: Added missing endpoints for listing orders, creating stores, and managing warehouses. Fixed CORS and allowed methods configuration.

### Verification & Testing:
- **E2E Tests**:
    - `admin_restrictions.spec.ts`: `PASSING` (Verified Admin restricted from domain writes).
    - `warehouse_ops.spec.ts`: `PASSING` (Verified Warehouse Manager fleet assignment modal).
    - `logistics.spec.ts`: `PASSING` (Verified map light mode and route visibility).
- **Technical Documentation**: Detailed system flows for Logistics, Traceability (CQRS), and Order SAGA documented in `docs/TECH.md`.

## Current Context & System Shape
- **Client App Port**: `http://localhost:3000`
- **KrakenD Gateway**: `http://localhost:8081`
- **Internal Shared Secret**: Standardized across all `.env` files.

## Next Steps
- [ ] Implement actual Valkey Geo tracking in `logistics-service`.
- [ ] Enhance Saga Monitor with real-time SSE updates.
- [ ] Move to Phase 3: System Reliability & Chaos Testing.

## Agent Constraints
- Always export the current session/task context, active changes, and outstanding tasks to `docs/CONTEXT.md` before concluding.
- Follow Clean Architecture patterns and central authorization designs.
- Strictly adhere to specified user/manager role behaviors.
