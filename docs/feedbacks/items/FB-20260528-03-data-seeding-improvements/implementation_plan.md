# Implementation Plan: E2E Data Seeding Strategy

## 1. Expected Changes
The current seeding system uses random numbers for emails and lacks a fully completed end-to-end historical record, leaving the Traceability and Finance dashboards empty. 

We will build a **Logical Data Seeder** that creates deterministic users, physical nodes, and a fully completed "Farm-to-Cup" flow spanning back several days so that historical data is immediately available upon environment reset.

## 2. User & Resource Seeding (Deterministic)
We will replace `mgr.<timestamp>@runtimeroasters.com` with deterministic identities tied to the locations they manage.

| Role | Email | Assigned Entity |
| :--- | :--- | :--- |
| `ADMIN` | `admin@runtimeroasters.com` | Global |
| `FARM_MANAGER` | `mgr.caudat@runtimeroasters.com` | Farm: Cầu Đất (Arabica) |
| `FARM_MANAGER` | `mgr.bmt@runtimeroasters.com` | Farm: Buôn Ma Thuột (Robusta) |
| `WAREHOUSE_MGR` | `mgr.songthan@runtimeroasters.com` | Warehouse: KCN Sóng Thần |
| `STORE_MGR` | `mgr.hcm01@runtimeroasters.com` | Store: HCM District 1 |
| `STORE_MGR` | `mgr.hcm03@runtimeroasters.com` | Store: HCM District 3 |
| `DRIVER` | `driver.songthan01@runtimeroasters.com` | Auto-seeded to KCN Sóng Thần |

## 3. The Logical Historical Flow (Traceability Seed)
To populate `/dashboard/traceability` and `/dashboard/finance`, we will programmatically generate a **Completed E2E Saga** with staggered chronological timestamps to simulate a process that occurred over the past week.

**The Historical Timeline (Mocked for Traceability):**
- **T - 5 Days**: `FARM_MANAGER` (Cầu Đất) creates a Harvest. (Event: `farm.harvest.created`)
- **T - 4 Days**: `WAREHOUSE_MGR` dispatches a truck. Driver completes the pickup route. (Events: `logistics.pickup.completed`, `warehouse.intake.created`)
- **T - 3 Days**: Roastery processes the green beans. Weight drops by 15% due to moisture loss. A finished Roast Batch is created. (Event: `warehouse.batch.finalized`)
- **T - 2 Days**: Stock is updated and made available for sale. (Event: `warehouse.stock.updated`)
- **T - 1 Day**: `STORE_MGR` (HCM01) places a retail order for $500. (Events: `retail.order.created`, `payment.intent.created`, `payment.simulated_completed`)
- **T - 6 Hours**: Warehouse reserves stock and dispatches delivery. Driver completes the store delivery. (Events: `warehouse.stock.reserved`, `logistics.delivery.completed`)

## 4. Implementation Approach
- **Scripted Generation**: Create a dedicated `seed_e2e_history.go` script (or expand `reset-demo-state.sh`).
- **Direct Event Insertion**: To avoid race conditions and the complexity of real-time SAGA delays during environment startup, the script will directly insert the required historical domain events into the `trace-service` read models (PostgreSQL & Elasticsearch) and `finance` models, ensuring the data is instantly queryable without waiting for Kafka processing.
- **Active State**: The seed script will also populate at least 1 "Active" (In-Progress) order so the Admin can see live deliveries on the Logistics Map.
- **Documentation Update**: Because the demo scenario heavily impacts how reviewers evaluate the system, we will explicitly update the demo guides (e.g., `docs/product/GUIDE.md` or `README.md`) to reflect the new deterministic accounts (`mgr.caudat@...`) and the pre-populated historical flow.

## 5. Impacted Scope
- `deployments/reset-demo-state.sh`
- `src/apps/logistics-service/internal/seed/logistics.json` (for location coordinates)
- Potential new seed script inside `trace-service` or a generic `seeder` script.
- `docs/product/GUIDE.md` or relevant demo documentation files to sync the script with the manual instructions.

## 6. Required Validation
- Run `./deployments/reset-demo-state.sh`.
- Log in as `admin@runtimeroasters.com`.
- Verify `/dashboard/users` lists the deterministic emails above.
- Verify `/dashboard/traceability` contains a fully rendered timeline graph for the Cầu Đất -> Sóng Thần -> HCM01 batch.
- Verify `/dashboard/finance` shows completed historical revenue.
