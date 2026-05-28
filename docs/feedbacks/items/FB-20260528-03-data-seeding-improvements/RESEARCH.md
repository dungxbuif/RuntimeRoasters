# Phase 13: Data Seeding Improvements - Research

## Summary
The current seeding mechanism is fragmented across infrastructure scripts (`seed.sh`), service-level Go code (`SeedData` methods), and JSON files. While core infrastructure (Kratos/Hydra) and static resources (Vehicles, Locations) are partially seeded, the system lacks deterministic identities for all roles and a historical "Farm-to-Cup" Saga flow, resulting in empty Traceability and Finance dashboards on initialization.

## Key Findings

### 1. Existing Scripts & Database Initialization
- **`deployments/reset-demo-state.sh`**: Focuses on cleanup. It deletes 30+ Kafka topics, truncates tables across 11 databases (Retail, Payment, Warehouse, Logistics, Trace, Audit), and clears Elasticsearch (`coffee_traceability` index) and Cassandra.
- **`deployments/seed.sh`**: Handles infrastructure-level seeding. It provisions the `admin@runtimeroasters.com` user in Ory Kratos and registers the OAuth2 client in Ory Hydra.
- **`deployments/init-db.sql`**: Only creates the 11 microservice databases; it does not contain data insertion logic.

### 2. Identity & Manager Seeding
- **Current State**: Only the `ADMIN` role is consistently seeded. Manager accounts (Farm, Warehouse, Store) are often generated with random timestamps (e.g., `mgr.1779846686075@...`) during E2E tests or manual intake.
- **Deterministic Alignment**: `src/apps/retail-service/internal/seed/stores.json` already contains some deterministic emails like `mgr.hn.hoankiem@runtimeroasters.com`. These should be standardized across the system (e.g., `mgr.caudat@runtimeroasters.com` for the Cầu Đất farm).

### 3. Logistics & Resource Seeding
- **`src/apps/logistics-service/internal/seed/logistics.json`**: This is the source of truth for physical coordinates. It defines:
  - **Farms**: Cầu Đất (Arabica), BMT (Robusta), Pleiku (Chư Sê).
  - **Warehouses**: Hanoi (HN-001), Sóng Thần (HCM-001), Hòa Khánh (DN-001).
  - **Retailers**: Hoàn Kiếm, Cầu Giấy, District 1, District 7, Hải Châu.
- **Service Logic**: `logistics-service` and `retail-service` both have `internal/seed/seed.go` files that embed these JSON files and provide `SeedData` methods accessible via GRPC (SystemHandler).

### 4. Historical Data & Traceability
- **The Gap**: The `/dashboard/traceability` (Trace Service) and `/dashboard/finance` (Payment Service) pages are empty because no Saga flows have been completed in a fresh environment.
- **Recommendation**: To avoid the complexity of running real-time Sagas (which involve Kafka delays), historical data should be seeded via **Direct Event Insertion**.
- **Trace Service Storage**: 
  - **PostgreSQL (`trace_db`)**: Stores `trace_events` and `trace_documents`.
  - **Elasticsearch**: Index `coffee_traceability` for graph queries.
- **Payment Service Storage**: PostgreSQL (`payment_db`) table `payments`.

### 5. Deterministic Identity Mapping
The following mapping was identified for replacement of random `mgr.*` emails:
- `FARM_MANAGER`: `mgr.caudat@runtimeroasters.com`, `mgr.bmt@runtimeroasters.com`
- `WAREHOUSE_MGR`: `mgr.songthan@runtimeroasters.com`
- `STORE_MGR`: `mgr.hcm01@runtimeroasters.com`, `mgr.hn.hoankiem@runtimeroasters.com`
- `DRIVER`: `driver.songthan01@runtimeroasters.com`

## Proposed Implementation Path
1. **Consolidated Seeder**: Create a Go-based seeder (or expand `demo-service`) that orchestrates the `SeedData` GRPC calls across all services.
2. **Historical Timeline**: Programmatically generate events with staggered `occurred_at` timestamps (T-5 days to T-6 hours) to simulate a complete Harvest -> Roasting -> Delivery cycle.
3. **Traceability Sync**: Ensure events inserted into PostgreSQL are also indexed into Elasticsearch to make the Traceability graph immediately visible.
