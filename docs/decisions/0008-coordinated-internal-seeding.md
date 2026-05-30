# ADR 0008: Coordinated Internal Seeding via Auth Service

## Status
Accepted

## Context
When bootstrapping the Runtime Roasters local distributed microservices monorepo for the first-run experience (FRX), we must seed demo data (Farms, Harvests, Warehouse Intakes, Batches, Orders, Payments, and Shipments) with realistic relationships.
To keep authorization and data integrity robust, database records must be associated with the user's real Kratos UUID (`created_by` or `driver_id` columns) rather than storing email strings directly in the database.
However, because Ory Kratos generates dynamic UUIDs at startup/registration, backend services cannot predict these UUIDs without querying.

## Decision
Instead of making Next.js `client-app` query all users and coordinate seeding calls to individual microservices (which leaks service internals to the frontend and degrades performance):
1. The Admin UI will trigger a single `POST /v1/system/seed` API request directed to `auth-service`.
2. `auth-service` (possessing admin Kratos access) will dynamically fetch Kratos identities, map their emails to their newly assigned dynamic UUIDs, and compile a `users_map` payload.
3. `auth-service` will internally propagate seed requests to downstream services (e.g., `logistics-service`, `retail-service`, `farm-service`) in parallel, passing the `users_map`.
4. Downstream services will use this map to dynamically look up the real Kratos UUIDs for fields like `created_by` or `driver_id`.

## Consequences
- **Security**: The Next.js frontend is insulated from internal service seeding orchestrations.
- **Robustness**: Database record mappings are strictly validated against Kratos UUIDs, preventing permissions issues with GormScoper.
- **Cohesion**: Cross-service mock data is linked dynamically using unique business keys (such as email, name, role) instead of hardcoded primary keys.
