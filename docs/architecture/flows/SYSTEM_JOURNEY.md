# Runtime Roasters: Constitutional System Journey

This document defines the canonical "Farm-to-Cup" operational flow. The system's integrity is verified by adhering to this holistic journey.

## The "Golden Path" Lifecycle

### 1. Origin (Farm Manager)
- **Action**: Declare a physical coffee harvest.
- **Verification**: Harvest is recorded in the ledger and a Pickup Request is automatically broadcasted.
- **Evidence**: `Harvest Declaration` dashboard reflects the new batch.

### 2. Logistics Orchestration (Warehouse Manager)
- **Action**: Dispatch a driver and vehicle for the pending harvest pickup.
- **Verification**: Fleet is assigned and the shipment transitions to `DISPATCHED`.
- **Evidence**: `Inbound Queue` shows live driver progress.

### 3. Physical Movement (Driver)
- **Action**: Simulate the route from base to farm and back to warehouse.
- **Verification**: Milestones (Arrived, Loaded, Returning) are confirmed.
- **Evidence**: `Driver Client` completes the mandatory return-to-base leg.

### 4. Retail Fulfillment (Store Manager)
- **Action**: Initialize a "Paid Order" fulfillment saga.
- **Verification**: Saga coordinates inventory reservation and outbound delivery.
- **Evidence**: `Retail Paid Orders` dashboard visualizes the saga steps (Payment -> Reserved -> Shipping).

### 5. Data Proof (Public Customer)
- **Action**: Scan a cup's QR code to verify its origin.
- **Verification**: Trace service assembles the complete hybrid lineage.
- **Evidence**: `Public Trace Registry` displays 100% transparent history from the specific harvest and roast batch.

---
*Verified via automated E2E journey script: `src/apps/client-app/e2e/flows/system_journey.spec.ts`*
