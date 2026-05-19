# 🛠️ Runtime Roasters Engineering Log

This document tracks technical challenges, bug fixes, and significant implementation milestones. It serves as a historical record of the system's evolution.

---

## 🐞 Bugs Found & Fixed

### Auth & Security
- **Login Accept 404:** Fixed an issue where `POST /v1/auth/login/accept` returned 404 because Gin's `NoRoute` handler interfered with the gRPC gateway.
- **OIDC Callback Loop:** Fixed the application landing page to properly capture the access token, strip it from the URL, and redirect to the dashboard.
- **Casbin Method Mapping:** Resolved a 403 error in `farm-service` by correctly mapping gRPC methods to Casbin actions (Get/List ➔ `read`, Delete ➔ `delete`, others ➔ `write`).
- **JWKS Protection:** Implemented an Nginx proxy (`rr-identity`) to protect Hydra's JWKS endpoint using an internal secret, ensuring only authorized services can fetch public keys.

### Service Interactions
- **Auth Polling Fallback:** Added a background bootstrap mechanism in `farm-service` to resiliently fetch Casbin snapshots if `auth-service` is not ready at startup.
- **Transactional Consistency:** Fixed a race condition in the Outbox worker where events could be published twice; added database-level locking during the "Mark as Processed" phase.
- **Elasticsearch Mapping:** Updated the `trace-service` read model to use correct keyword mappings for UUID fields, enabling accurate filtering by `order_id`.

---

## 🚀 Key Milestones

### Sprint 6-10: The SAGA Flow
Successfully connected the end-to-end coffee order flow:
1.  **Retail Service:** Initiates Order.
2.  **Warehouse Service:** Reserves coffee bean stock.
3.  **Payment Service:** Simulates Stripe/VNPay processing.
4.  **Logistics Service:** Assigns a driver and calculates real-time route.
5.  **Trace Service:** Fills the Elasticsearch read model for the "Track My Order" feature.
6.  **Audit Service:** Writes immutable logs to Cassandra for compliance.

### Documentation & Knowledge Base (Current)
Executed a comprehensive overhaul of the project's documentation to move from "Service-centric" to "Domain-centric" knowledge:
- Created the **Master System Architecture Specification** as the definitive technical source.
- Compiled the **Business & Domain Specification** (handover-ready) to define core supply chain rules (Roasting loss, FIFO, Batching).
- Consolidated all architectural flows into a single visual reference.
- Synchronized the `.planning/codebase/` map with the current implementation state.

### Multi-Region Design
Implemented W3C Trace-ID propagation across all transport layers (HTTP ➔ gRPC ➔ Kafka), allowing a single request to be tracked as it hops through the entire ecosystem.

---
*Last Updated: May 2026*
