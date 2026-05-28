---
phase: 13-data-seeding-improvements
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - src/apps/retail-service/internal/seed/stores.json
  - src/apps/logistics-service/internal/seed/logistics.json
  - src/apps/retail-service/internal/usecase/service.go
  - src/apps/logistics-service/internal/usecase/service.go
  - deployments/seed.sh
autonomous: true
requirements: [SEED-01, SEED-02, SEED-06]

must_haves:
  truths:
    - "Store, Farm and Logistics managers have deterministic emails tied to their locations."
    - "Resource IDs for farms, warehouses, and stores are consistent across services."
    - "New deterministic identities are seeded into Ory Kratos for login."
  artifacts:
    - path: "src/apps/retail-service/internal/seed/stores.json"
      provides: "Deterministic store metadata"
    - path: "src/apps/logistics-service/internal/seed/logistics.json"
      provides: "Deterministic logistics nodes"
  key_links:
    - "stores.json -> retail-service"
    - "logistics.json -> logistics-service"
    - "seed.sh -> Ory Kratos"
---

<objective>
Standardize resource identities and manager accounts across the system to ensure deterministic state initialization, including Ory Kratos identity provisioning.
</objective>

<tasks>
<task type="auto">
  <name>Task 1: Update Seed JSON Files</name>
  <files>src/apps/retail-service/internal/seed/stores.json, src/apps/logistics-service/internal/seed/logistics.json</files>
  <action>
    Replace random manager emails and IDs with deterministic ones:
    - Farm Managers: mgr.caudat@runtimeroasters.com (Farm: CAUDAT), mgr.bmt@runtimeroasters.com (Farm: BMT)
    - Warehouse Manager: mgr.songthan@runtimeroasters.com (Warehouse: SONGTHAN)
    - Store Managers: mgr.hcm01@runtimeroasters.com (Store: Q1), mgr.hn01@runtimeroasters.com (Store: HOANKIEM)
    - Driver: driver.songthan01@runtimeroasters.com (Auto-seeded to KCN Sóng Thần)
    Ensure IDs match between logistics.json and stores.json for shared entities.
  </action>
  <verify>grep "driver.songthan01" src/apps/logistics-service/internal/seed/logistics.json</verify>
  <done>JSON files contain the new deterministic mapping including the DRIVER role.</done>
</task>

<task type="auto">
  <name>Task 2: Align SeedData Logic</name>
  <files>src/apps/retail-service/internal/usecase/service.go, src/apps/logistics-service/internal/usecase/service.go</files>
  <action>
    Update the internal/seed/seed.go or usecase logic to ensure that calling SeedData over gRPC correctly processes the updated JSON files and handles UPSERT logic to prevent duplicates.
  </action>
  <verify>go test ./src/apps/retail-service/... ./src/apps/logistics-service/...</verify>
  <done>Services correctly ingest deterministic seed data.</done>
</task>

<task type="auto">
  <name>Task 3: Seed Identities in Ory Kratos</name>
  <files>deployments/seed.sh</files>
  <action>
    Update seed.sh to include provisioning calls for all new deterministic accounts in Ory Kratos. Use a standard default password for all seeded managers/drivers.
  </action>
  <verify>grep "mgr.caudat" deployments/seed.sh</verify>
  <done>Kratos identity seeding is configured for all required roles.</done>
</task>
</tasks>

---

---
phase: 13-data-seeding-improvements
plan: 02
type: execute
wave: 2
depends_on: ["13-01"]
files_modified:
  - src/apps/demo-service/cmd/seeder/main.go
  - src/apps/demo-service/internal/seeder/history.go
autonomous: true
requirements: [SEED-03]

must_haves:
  truths:
    - "Historical Traceability data is visible on the dashboard immediately after reset."
    - "Historical revenue is visible on the Finance dashboard."
  artifacts:
    - path: "src/apps/demo-service/cmd/seeder/main.go"
      provides: "Main entry point for E2E seeder"
  key_links:
    - "seeder -> trace_db (PostgreSQL)"
    - "seeder -> trace_db (Elasticsearch)"
    - "seeder -> payment_db (PostgreSQL)"
---

<objective>
Implement a historical data seeder that directly populates the Trace and Payment databases with a completed "Farm-to-Cup" cycle.
</objective>

<tasks>
<task type="auto">
  <name>Task 1: Scaffold E2E Seeder</name>
  <files>src/apps/demo-service/cmd/seeder/main.go</files>
  <action>
    Create a new Go entry point that connects to trace_db and payment_db. 
    Define the historical timeline logic constants:
    - T-5 days: Harvest created
    - T-4 days: Intake created
    - T-3 days: Roast Batch finalized
    - T-2 days: Order created
    - T-1 day: Payment completed
    - T-6 hours: Delivery completed
  </action>
  <verify>grep "T-5 days" src/apps/demo-service/cmd/seeder/main.go || grep "Time.Now().Add(-5 \* 24 \* time.Hour)" src/apps/demo-service/cmd/seeder/main.go</verify>
  <done>Seeder project structure and timeline logic are ready.</done>
</task>

<task type="auto">
  <name>Task 2: Implement Direct DB Insertion</name>
  <files>src/apps/demo-service/internal/seeder/history.go</files>
  <action>
    Implement logic to insert TraceEvents and TraceDocuments into PostgreSQL (trace_db) and index them into Elasticsearch (coffee_traceability).
    Insert completed Payment records into PostgreSQL (payment_db).
    Use the deterministic IDs from Plan 01 to ensure the graph is connected.
  </action>
  <verify>go run src/apps/demo-service/cmd/seeder/main.go && docker exec postgres psql -U user -d trace_db -c "SELECT count(*) FROM trace_events;"</verify>
  <done>Seeder can generate and insert historical data into all required tables and indices.</done>
</task>
</tasks>

---

---
phase: 13-data-seeding-improvements
plan: 03
type: execute
wave: 3
depends_on: ["13-02"]
files_modified:
  - deployments/reset-demo-state.sh
  - docs/product/GUIDE.md
  - README.md
autonomous: true
requirements: [SEED-04, SEED-05]

must_haves:
  truths:
    - "Environment reset flow is fully automated and includes all new seeders."
    - "Documentation accurately reflects the seeded demo state."
  artifacts:
    - path: "deployments/reset-demo-state.sh"
      provides: "Unified demo reset orchestration"
  key_links:
    - "reset-demo-state.sh -> logistics/retail SeedData"
    - "reset-demo-state.sh -> demo-service seeder"
---

<objective>
Integrate the comprehensive seeding flow into the environment reset script and update documentation.
</objective>

<tasks>
<task type="auto">
  <name>Task 1: Update Reset Script</name>
  <files>deployments/reset-demo-state.sh</files>
  <action>
    Modify reset-demo-state.sh to:
    1. Call the SeedData gRPC endpoints for logistics and retail services.
    2. Execute the new E2E history seeder (go run src/apps/demo-service/cmd/seeder/main.go).
    Ensure the script handles service readiness wait times.
  </action>
  <verify>./deployments/reset-demo-state.sh</verify>
  <done>Environment reset automatically triggers comprehensive seeding.</done>
</task>

<task type="auto">
  <name>Task 2: Update Documentation</name>
  <files>docs/product/GUIDE.md, README.md</files>
  <action>
    Update the "Login" and "Demo Scenario" sections to reflect the new deterministic account credentials (e.g., mgr.caudat@runtimeroasters.com) and describe the pre-populated historical data available on the dashboards.
  </action>
  <verify>Manual check of markdown files</verify>
  <done>Documentation is in sync with the new seeding state.</done>
</task>
</tasks>
