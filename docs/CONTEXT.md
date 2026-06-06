---
artifact_type: context
id: CONTEXT
status: active
owner: shared
---
# Session Context & Development State

Last Updated: 2026-06-06

## Current Git Snapshot
**Checked:** 2026-06-06
**Active branch:** `rr-urg-07-trace-public-qr`
**HEAD:** `7663268 RR-URG-07 start trace public QR`

### Active Dirty Changes
- RR-URG-07A Retail schema and seed implementation:
  - rewrote the development `000001_init.up.sql` without business inserts.
  - added `menus`, 42 `menu_items` from `SAMPLE_MENU.md`, inventory lots,
    sales, one-cup sale items, stock movements, and materialized availability.
  - seeded 5 stores, 10 lineage lots, 25 sold cups, and 210 store-menu
    availability rows deterministically.
  - added ledger reconciliation and insert-if-absent demo semantics so re-seed
    does not reset user runtime movements.
  - recorded ADR-0010 and reconciled SDD, ERD, master data, tickets, validation,
    traceability, and release notes.
- `docs/architecture/SDD/README.md` and `docs/architecture/README.md`: added the missing system-level Software Design Document, consolidating system boundaries, service ownership, primary flows, data/security/runtime design, failure handling, verification, and links to canonical master docs.
- `src/apps/client-app/src/constants/casbin.ts`: frontend Casbin matcher now allows direct role policy matching (`r.sub == p.sub`) as well as inherited role matching (`g(r.sub, p.sub)`), fixing frontend unit failures where role policies such as `STORE_MGR` did not evaluate.
- `src/apps/socket-service/internal/app/app_test.go`: added socket-service route integration coverage for private stream ticket requirement, public WebSocket access, ticketed private WebSocket access, and one-time ticket rejection.
- `src/apps/client-app/e2e/operational_realtime_awareness.spec.ts`: added flow-oriented Playwright coverage for Warehouse/Store dashboard notification rendering and role-scope isolation with mocked gateway auth/API contracts; passed through production Next runtime.
- `src/apps/client-app/package.json`: added `start:e2e` and `test:e2e:realtime-notifications` scripts.
- `src/apps/client-app/playwright.config.ts`: added fail-fast E2E guardrails: production webServer, `127.0.0.1` baseURL, 30s test timeout, 10s navigation timeout, 5s action/expect timeout, line+HTML reporters.
- `src/apps/socket-service/internal/platform/doc.go` and `src/apps/socket-service/internal/platform/platform_live_test.go`: added opt-in live platform test for Kafka -> socket-service -> Redis/Valkey notification persistence.
- `src/apps/warehouse-service/internal/app/app.go`: removed duplicate stub registration for `/v1/warehouse/dispatch-requests` routes so warehouse-service can boot for platform verification.
- `docs/work/tickets/sprint10/RR-URG-06/ticket.md`: updated verification evidence and remaining manual visual demo gap.
- Harness v1 docs reconciliation:
  - Audited new Harness reference repo at `/private/tmp/harness-new`.
  - Exported old docs from `4404600^` into `/private/tmp/runtime-roasters-old-docs`.
  - Added reviewed RuntimeRoasters-specific Harness core files instead of copying raw templates: `docs/work/FEEDBACK_LOG.md`, `docs/work/TRACEABILITY.md`, `docs/requirements/REQUIREMENTS.md`, `docs/requirements/USER_STORIES.md`, engineering setup/local/troubleshooting docs, work/ticket/bug/phase indexes/templates, `docs/templates/FEEDBACK.md`, `docs/templates/SDD.md`, and `docs/releases/RELEASE_NOTES_TEMPLATE.md`.
  - Updated stale migration paths from pre-Harness product/story/feedback locations to Harness v1 locations.
  - Split RR-URG-07 into 07A-07E with Harness frontmatter and trace metadata.
  - Consolidated legacy non-template docs into Harness targets: requirements, master data, UI design, operations, canonical standards, roadmap, and validation matrix; removed duplicate legacy files after review.

### Outstanding Verification
- RR-URG-07A targeted Retail tests passed:
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/retail-service/...`
- RR-URG-07A fresh PostgreSQL migration passed on `retail_07a_test`.
- Opt-in live PostgreSQL service seed test remains pending because localhost
  network access was denied by the execution sandbox.
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
- Harness docs structure audit passed:
  - No required docs from the cloned Harness reference are missing.
  - Stale path scan for pre-Harness product/story/feedback locations returned no active matches.
- Manual visual review still pending:
  - Full two-browser role-session demo with real logged-in users has not been run.

### Next Task
- Next ticket family is RR-URG-07: Public Sold-Cup QR Trace.
- RR-URG-07 was split into:
  - `RR-URG-07A`: retail sale schema and seed data.
  - `RR-URG-07B`: retail demo sale APIs.
  - `RR-URG-07C`: trace-service public sold-cup query.
  - `RR-URG-07D`: client sale and QR UI.
  - `RR-URG-07E`: evidence and platform verification.
- Current discussion corrected RR-URG-07 demo seed/QR model: QR semantics should represent a UI-issued `product_id` for a retail sold cup/item.
- Identifier vocabulary: `menu_item_id` identifies a chain-wide menu entry; `product_id` identifies one sold cup/item and is the public trace token.
- Current code/system drift identified before implementation:
  - `src/apps/client-app/src/app/trace/[code]/page.tsx` is still static/mock (`DEMO_JOURNEY`).
  - trace-service `TraceDocument` and Elasticsearch read model do not yet model `trace_code` or public sold-unit lookup.
  - trace-service public routes currently expose public topology only.
- Needed RR-URG-07 implementation direction:
  - Keep cup purchase simple: no payment and no live delivery workflow in the buy-click.
  - UI creates `product_id`; backend validates/persists it and never derives origin from the QR token alone.
  - Add chain-wide retail menu items, retail inventory lots, retail sales/invoices, sale items, and stock movements.
  - Seed store inventory lots with upstream farm/batch/warehouse lineage.
  - Add trace-service hybrid lookup by `product_id`/`trace_code`, combining retail sold-cup details with upstream trace data.
  - Replace static public trace UI with API-backed fetch/rendering and add sold-items UI.

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
- **Harness Migration:** Upgrading to Harness v1 (Pure Markdown) - structure reconciled against cloned Harness reference; content was rewritten/merged from RuntimeRoasters docs rather than blindly copied.
- **Next Up:** finish the live PostgreSQL seed check, then RR-URG-07B Retail Demo Sale APIs.
- RR-URG-07A is in review at
  `docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/technical_design.md`.
- `docs/architecture/SDD/MASTER_DATA.md` is the canonical contract for seed
  data and domain enums; `docs/requirements/SAMPLE_MENU.md` owns menu content.
- Farm role contract simplified by ADR-0009: `FARM_MANAGER` is the only farm operator role; the duplicate role was removed from Kratos, Casbin, farm-service, frontend types/options, tests, and docs.
- FARM_MANAGER-only verification:
  - repository-wide removed-role scan returned no matches.
  - Kratos identity schema JSON parsed successfully.
  - targeted Casbin/shared-scoper/farm repository tests passed.
  - frontend unit tests passed 9/9 and lint passed.
  - full farm-service test command remains blocked by an existing test mock missing `FarmRepository.Count`; unrelated to the role removal.
  - frontend production build could not be re-run because PID `37366` was already running `next build` and held `.next/lock` for more than 60 seconds; the process was not killed.
  - existing external/Kratos identities carrying the removed role value must be migrated to `FARM_MANAGER` before token issuance.
