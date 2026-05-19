# Sprint 2: Identity & Access Control

**Status:** ✅ Completed
**Goal:** Build a centralized Identity and Authorization system, protecting all APIs and providing standard OIDC Login/Logout flows.

---

## 📋 Ticket Status (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-9](./RR-9/ticket.md) | [Tech] Identity Infrastructure (Ory Kratos & Hydra) | ✅ Done | Security Eng |
| [RR-10](./RR-10/ticket.md) | [BA] Client Auth Flow: Login/Consent UI | ✅ Done | Product Owner |
| [RR-11](./RR-11/ticket.md) | [Tech] JWT Validation & Propagation | ✅ Done | Tech Lead |
| [RR-12](./RR-12/ticket.md) | [Epic] Casbin RBAC Authorization Engine | ✅ Done | Tech Lead |
| [RR-13](./RR-13/ticket.md) | [Tech] E2E Integration: Secure Service-to-Service | ✅ Done | Backend |
| [RR-14](./RR-14/ticket.md) | [Tech] Token Revocation & Valkey Blacklist | ✅ Done | Tech Lead |

---

## 💡 Business Vision
- **Zero Trust Architecture:** Every request entering the system must be strictly identified and authorized.
- **Seamless Experience:** Users only need to log in once (SSO) to access all services within the ecosystem.
- **Fine-grained Control:** Authorization down to specific actions (e.g., Farm Manager can only create Harvests, not approve Payments).

## 📊 Sprint Result
- Successfully integrated the Ory suite (Kratos/Hydra) for Identity management.
- Implemented **Three-Gate Security** (Gateway -> JWT -> Casbin).
- The system supports instantaneous Token revocation via Valkey.
