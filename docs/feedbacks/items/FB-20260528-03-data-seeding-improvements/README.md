# Data Seeding Improvements

## Status: TESTING

## Description
This ticket consolidates multiple data-seeding related feedbacks:
1. **Seed FarmManager email**: Data seed is currently using random numbers for `FARM_MANAGER` (e.g., `mgr.1779846686075@...`). We need deterministic emails (e.g., `mgr.caudat@...`).
2. **Seed Users & Roles**: Currently only `FARM_MANAGER` is seeded. Need to seed all roles (Admin, Warehouse, Retail, Driver).
3. **Data Seed & Financial Integrity**: Traceability and Finance dashboards are empty on fresh start.
4. **Traceability Page Empty**: `/dashboard/traceability` should show a completed "Farm-to-Cup" path.

## Update (2026-05-29)
- Implementation for historical seeder (7 days) completed.
- Deterministic Kratos accounts standardized.
- Integrated into `task env:reset`.
- Ready for manual verification.

**Unified Goal**: Implement a comprehensive Data Seeding strategy that pre-populates deterministic users, resources, and a chronological historical flow (Saga) so all dashboards (Traceability, Finance, Logistics) are demo-ready upon running `scripts/reset-env.sh`. Update demo documentation accordingly.
