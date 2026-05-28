# Data Seeding Improvements

## Status: IN_PROGRESS

## Description
This ticket consolidates multiple data-seeding related feedbacks:
1. **Seed FarmManager email**: Data seed is currently using random numbers for `FARM_MANAGER` (e.g. `mgr.1779846686075@runtimeroasters.com`). Desired outcome is to use deterministic emails based on the assigned farm.
2. **Seed Users & Roles**: Currently, only `FARM_MANAGER` roles are seeded. Need to seed all roles (`ADMIN`, `WAREHOUSE_MGR`, `STORE_MGR`, etc.) to make the system ready for a full demo.
3. **Data Seed & Financial Integrity**: The system lacks a comprehensive initial data seed, particularly for Financial Integrity, which is completely empty upon initialization.
4. **Traceability Page Empty**: The `/dashboard/traceability` page is completely empty on fresh start. Requires a fully completed sample dataset (historical Saga flow) so the Traceability graph is visible immediately.

**Unified Goal**: Implement a comprehensive Data Seeding strategy that pre-populates deterministic users, resources, and a chronological historical flow (Saga) so all dashboards (Traceability, Finance, Logistics) are demo-ready upon running `reset-demo-state.sh`. Update demo documentation accordingly.
