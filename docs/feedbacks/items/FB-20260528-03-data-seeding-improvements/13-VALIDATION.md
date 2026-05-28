# Phase 13: Data Seeding Improvements - Validation

## Validation Architecture

The validation for Phase 13 will be performed at three levels to ensure that the data seeding is deterministic, functional, and demo-ready.

### 1. Deterministic Identity Verification
- **Target:** Ory Kratos & System Profiles.
- **Protocol:** Run `./deployments/reset-demo-state.sh` and verify that all pre-defined accounts can log in and have the correct profile data (Name, Role, Location).
- **Tooling:** Kratos CLI and gRPC System Handler tests.

### 2. Historical Data Integrity
- **Target:** Trace PostgreSQL, Elasticsearch, and Payment PostgreSQL.
- **Protocol:** Check row counts in `trace_events`, `trace_documents`, and `payments` tables immediately after environment reset. Verify that Elasticsearch contains the same events as PostgreSQL for graph rendering.
- **Tooling:** `psql`, `curl` (for ES), and `go run src/apps/demo-service/cmd/seeder/main.go --verify`.

### 3. Frontend Dashboard Readiness
- **Target:** Admin Intelligence Hub, Traceability Graph, and Finance View.
- **Protocol:** Access the dashboard as `admin@runtimeroasters.com`. Confirm that the Traceability graph renders a complete "Farm-to-Cup" path and that Finance charts show historical revenue from the past 7 days.
- **Tooling:** E2E Playwright tests and manual UI audit.

## Verification Scenarios

| Scenario | Input | Expected Output | Status |
| :--- | :--- | :--- | :--- |
| **System Bootstrap** | `./deployments/reset-demo-state.sh` | Zero errors, all 11 DBs seeded, 100+ events in Trace DB. | `PENDING` |
| **Login Verification** | `mgr.caudat@runtimeroasters.com` | Success login, role: `FARM_MANAGER`, view: Cầu Đất Farm. | `PENDING` |
| **Graph Visualization** | Access `/dashboard/traceability` | Graph shows Harvest -> Roastery -> Store delivery path. | `PENDING` |
| **Financial Audit** | Access `/dashboard/finance` | Revenue > $0 from historical seeded retail orders. | `PENDING` |
