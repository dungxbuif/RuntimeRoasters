# System Bootstrap & Seeding Protocol

This document defines the process for initializing the Runtime Roasters environment with a consistent, demo-ready state.

> For the distinction between deterministic seeding, backend-simulated history, browser-simulated GPS, and real runtime flows, see [FB-20260528-08 Demo Flow Contract](../../feedbacks/items/FB-20260528-08-demo-flow-contract/DEMO_FLOW_CONTRACT.md).

## 1. First-Run Experience (FRX)

When the system is started with empty databases, the Admin Dashboard provides a guided bootstrap process.

- **Trigger**: The `SystemBootstrapModal.tsx` component in `client-app` automatically detects the unseeded state.
- **Check Mechanism**: Calls `GET /v1/logistics/system/status` and `GET /v1/retail/system/status`. If `seeded: false` is returned, the modal appears.
- **Action**: Clicking "Initialize DB" calls the `seedAll` method in `systemService`, which triggers the `/system/seed` POST endpoints on backend services.

## 2. Deterministic Identity Mapping

To ensure a predictable demo environment, we use the following fixed identities. All accounts use a standard default password (refer to local `.env`).

| Role | Email | Assigned Entity |
| :--- | :--- | :--- |
| `ADMIN` | `admin@runtimeroasters.com` | Global System Admin |
| `FARM_MANAGER` | `mgr.caudat@runtimeroasters.com` | Cầu Đất Farm (Arabica) |
| `FARM_MANAGER` | `mgr.bmt@runtimeroasters.com` | Buôn Ma Thuột Farm (Robusta) |
| `WAREHOUSE_MGR` | `mgr.songthan@runtimeroasters.com` | KCN Sóng Thần Warehouse |
| `STORE_MGR` | `mgr.hcm01@runtimeroasters.com` | Store: District 1, HCM |
| `STORE_MGR` | `mgr.hn01@runtimeroasters.com` | Store: Hoàn Kiếm, Hà Nội |
| `DRIVER` | `driver.songthan01@runtimeroasters.com` | Assigned to Sóng Thần Warehouse |

## 3. Historical Data Simulation (Saga Seeding)

Real-time Saga flows can take several minutes to complete across the distributed system. To provide immediate visualization on dashboards, the seeder simulates a historical timeline by directly inserting records into read models.

### Timeline Simulation (7-Day Span)
- **Day 1**: Harvest created at Cầu Đất.
- **Day 2**: Logistics pickup from Farm to Warehouse.
- **Day 3**: Roasting Batch finalized at Roastery.
- **Day 4**: Stock updated and made available.
- **Day 5**: Retail Order placed by Store Manager.
- **Day 6**: Payment successfully processed.
- **Day 7**: Final delivery to Retail Store.

### Target Read Models
- **Trace Service**: `trace_events` (PostgreSQL) and `coffee_traceability` (Elasticsearch).
- **Finance**: `payments` table (PostgreSQL) in `payment-db`.
- **Logistics**: `shipment_logs` in `logistics-db`.

## 4. Resetting State

To perform a clean demo, developers should use the provided script:

```bash
# Deletes all data, resets Kafka offsets, and prepares for a fresh bootstrap
./deployments/reset-demo-state.sh
```

After running this script, log in as `admin@runtimeroasters.com` to trigger the "Initialize DB" modal.
