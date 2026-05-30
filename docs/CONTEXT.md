# Session Context & Development State

Last Updated: 2026-05-30

## Phase 1: System Bootstrap & Seeding
**Status:** `✅ COMPLETED`

### Technical Achievements:
1.  **Login Flow Optimization:** Fixed the "double-click" stall by ensuring automatic Hydra flow initiation if Kratos authentication succeeds without a challenge. This resolves the state where a user had a Kratos session but no OAuth2 token.
2.  **Deterministic Seeding:** Implemented `SeedUsers` in `auth-service` to create all **26 master identities** (6 Farm Managers, 5 Store Managers, 3 Warehouse Managers, 11 Drivers, 1 Admin) defined in `MASTER_DATA.md` automatically.
3.  **Propagation Fixes:**
    *   Updated `auth-service` to use non-cancelable background context for seeding, preventing KrakenD timeouts (default 3s, increased to 60s) from killing the long-running identity creation process.
    *   Added `X-Internal-Secret` propagation to all downstream seeding calls to pass through security guards.
    *   Standardized ports in `.env`: Logistics (8085), Payment (8086), Trace (8087), Audit (8088), Warehouse (8089).
4.  **Data Integrity:** Updated seeding logic in all services to use **Upsert (Save)** instead of `FirstOrCreate`, ensuring master data updates correctly when forced even if partial records exist.
6.  **Frontend Authorization (RBAC):** Migrated hardcoded UI role checks to use robust `CasbinGuard`. Ensured `ADMIN` role honors Casbin rules from the database exactly as defined (e.g., Read-Only access to Domain operations, restricting Create/Delete actions). Fixed the Casbin DB policies to use correct `keyMatch` mapping syntax (e.g. `/v1/users/*`).

### Verification & Refactoring:
- **E2E Automation:** Relocated the verification script to `src/apps/client-app/e2e/bootstrap.setup.spec.ts` as a proper Playwright test and increased its timeout to support background seeding. Status: `PASSING`.
- **RBAC Magic Strings:** Eliminated hardcoded strings (e.g., `obj="/v1/warehouse/dispatch"`) from `CasbinGuard` by mapping them to `AUTH_RESOURCES` constants.
- **Linting:** Configured a custom ESLint `no-restricted-syntax` rule to catch and prevent future use of literal magic strings in `CasbinGuard` properties.
- **Data State:** Kratos has 26 identities; Farm (6), Retail (5), Logistics (seeded) are initialized.

## Current Context & System Shape
- **Mission**: Farm-to-Cup supply-chain system built as a Go microservices monorepo with Next.js control-plane UI.
- **Client App Port**: `http://localhost:3000`
- **KrakenD Gateway**: `http://localhost:8081`

## Next Steps
- [ ] Manual verification of Dashboard topology.
- [ ] Move to Phase 2: Logistics Real-time Tracking & Saga Simulator.

## Agent Constraints
- Always export the current session/task context, active changes, and outstanding tasks to `docs/CONTEXT.md` before concluding.
- Follow Clean Architecture patterns and central authorization designs.
- Strictly adhere to specified user/manager role behaviors.
