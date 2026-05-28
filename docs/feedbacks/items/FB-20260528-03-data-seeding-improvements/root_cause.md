# Root Cause Analysis

**Symptoms**:
1. Farm manager emails use random numbers (e.g., `mgr.1779846686075@...`).
2. Only Farm Managers are seeded, leaving other roles empty.
3. Financial Integrity dashboard is empty on start.
4. Traceability page is empty on fresh start.

**Actual Cause**: 
The existing seeding logic (Go seed scripts and `deployments/logistics-seed.sql`) was built incrementally and uses random timestamps for uniqueness. It does not populate the full suite of operational users. Furthermore, because the system relies on an Event-Driven Saga architecture, the Elasticsearch read models (for Traceability) and PostgreSQL aggregates (for Finance) remain empty until a user manually performs a complete End-to-End order flow in the UI.

**Missing/Misleading Info**:
There is no existing script designed to mock historical Domain Events and inject them directly into the read models to simulate past activity.
