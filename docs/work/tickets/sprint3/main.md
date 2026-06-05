# Sprint 3: Farm Service & Vertical Slice

**Status:** ✅ Completed
**Goal:** Finalize the Farm Management Service with full business features, applying a Vertical Slice approach combined with Clean Architecture to ensure readiness for the harvesting process.

---

## 📋 Ticket Status (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-15](./RR-15/ticket.md) | [Tech] Farm Service: Scaffolding & Composition Root | ✅ Done | Tech Lead |
| [RR-16](./RR-16/ticket.md) | [BA] Farm Catalog Management (Create & List) | ✅ Done | Product Owner |
| [RR-17](./RR-17/ticket.md) | [Tech] Farm CRUD: Update & Delete Logic | ✅ Done | Backend |
| [RR-18](./RR-18/ticket.md) | [BA] Land Parcel & Growing Region Planning | ✅ Done | Farm Manager |

---

## 💡 Business Vision
- **Digital Farm Twin:** Every real-world farm must be accurately reflected in the system with parameters such as area, altitude, and primary coffee type.
- **Resource Management:** Help managers grasp the production capacity of each raw material region (Cau Dat, Buon Ma Thuot, etc.).
- **Foundation for Traceability:** Farm information is the "root" of the entire subsequent traceability chain.

## 📊 Sprint Result
- Successfully initialized the Farm Service with all Clean Architecture layers.
- Completed the Farm Management API suite (CRUD).
- Integrated Casbin authorization: Only Farm Managers/Farm Admins can manage their respective farm data.
