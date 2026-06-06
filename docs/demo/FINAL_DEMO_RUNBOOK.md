# Runtime Roasters: Final Demo Runbook (v1.0)

This guide defines the step-by-step procedure for demonstrating the 100% transparent "Farm-to-Cup" journey.

## 1. Environment & Prerequisites

- **Base URL**: `http://localhost:3000`
- **Identity Provider**: Ory Kratos (Local Docker)
- **Data State**: Clean seed recommended (`task env:reset` or Admin Bootstrap).

### Demo Account Registry

| Persona | Email | Password | Role |
| :--- | :--- | :--- | :--- |
| **System Admin** | `admin@runtimeroasters.com` | `Hello@123` | `ADMIN` |
| **Farm Manager** | `mgr.farm.caudat@runtimeroasters.com` | `Hello@123` | `FARM_MANAGER` |
| **Warehouse Mgr**| `mgr.wh.hn@runtimeroasters.com` | `Hello@123` | `WAREHOUSE_MGR` |
| **Driver** | `driver.hn.01@runtimeroasters.com` | `Hello@123` | `DRIVER` |
| **Store Manager**| `mgr.hn.hoankiem@runtimeroasters.com`| `Hello@123` | `STORE_MGR` |

---

## 2. Window Layout (Recommended)

- **Window A (Left)**: Dashboard Operator (Login as Manager/Driver).
- **Window B (Right)**: System Observer (Root Page - Topology).
- **Mobile/Tab (Optional)**: Public QR Trace (`/trace`).

---

## 3. Script A: The Inbound Journey (Farm to Warehouse)

**Goal**: Move physical coffee from Dalat to Hanoi Warehouse.

1.  **Login as Farm Manager** (`mgr.farm.caudat@runtimeroasters.com`).
2.  Navigate to **Farm Ops > Harvests**.
3.  Click **"Declare New Harvest"**:
    - Quantity: `500.0 KG`
    - Variety: `ARABICA`
4.  **Observe**: A new row appears with status `NEW`.
5.  **Login as Warehouse Manager** (`mgr.wh.hn@runtimeroasters.com`) in a new tab/session.
6.  Navigate to **Warehouse Operations**.
7.  Find the `Harvest #...` in the **Inbound Queue**.
8.  Click **"Dispatch Pickup"**:
    - Select Driver: `HN Driver 01`
    - Select Vehicle: `RR-SAFE-01`
9.  **Observe**: Status changes to `DISPATCHED`. Progress bar shows "Driver moving to farm...".
10. **Login as Driver** (`driver.hn.01@runtimeroasters.com`).
11. Navigate to **Driver Client**:
    - Click **"Start Route"**.
    - Advance through milestones: `Departed` ➔ `At Farm` ➔ `Loaded` ➔ `Returning`.
12. **Switch back to Warehouse Manager**:
    - Observe status: `PICKED_UP` ➔ "Returning to warehouse...".
    - Once driver hits `Arrived WH`, click **"Confirm Receipt"**.
13. **Result**: Green beans are now in stock.

---

## 4. Script B: The Outbound Journey (Retail Fulfillment)

**Goal**: Fulfill a paid customer order in Hanoi.

1.  **Login as Store Manager** (`mgr.hn.hoankiem@runtimeroasters.com`).
2.  Navigate to **Retail > Orders**.
3.  Click **"Create Fulfillment"** (or use the Buy Coffee Kiosk at `/coffee` for a more public feel).
4.  If creating via Dashboard:
    - Add Item: `SKU-ARABICA-01` (Quantity: 10)
    - Click **"Broadcast Order"**.
5.  **Observe (Real-time)**:
    - Saga Visualizer shows `PAYMENT_COMPLETED` ➔ `STOCK_RESERVED`.
6.  **Login as Warehouse Manager**:
    - Find the order in **Outbound Queue**.
    - Click **"Dispatch Delivery"**.
7.  **Login as Driver**:
    - Select the Retail Delivery shipment.
    - Simulate movement to the store.
    - Click **"Delivery Confirmed"**.
8.  **Result**: Order transitions to `COMPLETED`. Trace code generated.

---

## 5. Script C: Public Proof & Topology

**Goal**: Demonstrate system health and cryptographic lineage.

1.  Navigate to `http://localhost:3000` (No login required).
2.  **Observe Architecture Topology**:
    - Verify services are green.
    - Watch event pulses during Script A or B (if live).
3.  Click **"Trace Registry"** or go to `/trace`.
4.  Find the latest sold cup (should match Script B).
5.  Click **"Verify Lineage"**.
6.  **Observe**:
    - Milestone Timeline: `Harvested` ➔ `Inbound` ➔ `Roasted` ➔ `Paid` ➔ `Served`.
    - 100% digital proof of Dalat origin for the Hanoi cup.

---

## 6. Recovery & Troubleshooting

- **Socket Stalls**: If the map or progress bars don't update, click the **Refresh** icon on the top right of the dashboard.
- **Access Denied**: Clear browser storage (`localStorage.clear()`) and re-login.
- **Data Desync**: Run `task env:reset` to return to a known baseline.

---
*Operational Constitution of Runtime Roasters.*
