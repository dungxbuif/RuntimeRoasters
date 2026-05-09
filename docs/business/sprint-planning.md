# Lộ trình Phát triển RuntimeRoasters (Portfolio Edition)

**Mục tiêu cốt lõi:** Showcase Senior Backend — kiến trúc phân tán, OTel, Clean Architecture, DI, API Gateway.

**Chiến lược:**
- 1 Solo Developer (Full-stack/Backend focus)
- Sprint 1 = **Infrastructure only** — mọi boilerplate, base code, infra setup phải hoàn tất
- Sprint 2+ = **Business logic only** — chỉ viết domain/usecase/repo, không setup thêm gì
- Ticket style: BDD AC only. Technical design do Dev tự tạo khi implement.

---

## Sprint 1 — Complete Infrastructure

**Goal:** Sau Sprint 1, mọi thứ đều chạy được end-to-end. Sprint 2 chỉ viết business logic.

**Demo Service** (`apps/demo-service/`) = canonical template cho tất cả services sau. Farm Service bắt đầu sạch từ Sprint 2 bằng cách copy từ đây.

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-1](./sprint1/RR-1/ticket.md) | Infra Kick-off — Monorepo, Go Workspaces | ✅ Done |
| [RR-2](./sprint1/RR-2/ticket.md) | Config & Logger — Viper, Zap, BaseConfig | ✅ Done |
| [RR-3](./sprint1/RR-3/ticket.md) | Base & Errs — RFC 9457, app lifecycle | ✅ Done |
| [RR-4](./sprint1/RR-4/ticket.md) | Proto + Wire + KrakenD + Client Shell | 🕒 To Do |
| [RR-5](./sprint1/RR-5/ticket.md) | Docker Complete — SigNoz + Full Compose | 🕒 To Do |
| [RR-6](./sprint1/RR-6/ticket.md) | Demo Service — Clean Arch + Full Stack Slice | 🕒 To Do |
| [RR-7](./sprint1/RR-7/ticket.md) | Client App — Control Plane | 🕒 To Do |
| [RR-8](./sprint1/RR-8/ticket.md) | Client App — Business UI Shell | 🕒 To Do |

Chi tiết: [`docs/business/sprint1/main.md`](./sprint1/main.md)

---

## Sprint 2 — Security & Access Control

**Goal:** Thiết lập nền tảng AuthN/AuthZ toàn diện. Mọi API call phải được bảo mật trước khi viết business logic.

| Ticket | Summary |
| :--- | :--- |
| [RR-9](./sprint2/RR-9/ticket.md) | Identity Server Infrastructure (Ory Kratos) |
| [RR-10](./sprint2/RR-10/ticket.md) | Client-Side Auth & Login Flow |
| [RR-11](./sprint2/RR-11/ticket.md) | Backend Security Core - JWT Validation |
| [RR-12](./sprint2/RR-12/ticket.md) | Fine-grained Authorization (Casbin) |
| [RR-13](./sprint2/RR-13/ticket.md) | End-to-End Secure Integration |
| [RR-14](./sprint2/RR-14/ticket.md) | Token Revocation & Global Logout |

Chi tiết: [`docs/business/sprint2/main.md`](./sprint2/main.md)

---

## Sprint 3 — Farm Service Business Logic

**Goal:** Farm Service đầy đủ CRUD. Copy boilerplate từ `apps/demo-service/`, chỉ viết business logic.

| Ticket | Summary |
| :--- | :--- |
| [RR-15](./sprint3/RR-15/ticket.md) | Farm Service Bootstrap & Domain |
| [RR-16](./sprint3/RR-16/ticket.md) | Farm Repository & Database |
| [RR-17](./sprint3/RR-17/ticket.md) | Farm UseCase & API (gRPC/REST) |
| [RR-18](./sprint3/RR-18/ticket.md) | Farm Control Plane (UI) |

Chi tiết: [`docs/business/sprint3/main.md`](./sprint3/main.md)

---

## Sprint 4 — Distributed Transactions (Saga)

**Goal:** Retail Service + Warehouse Service + Saga Choreography.

| Ticket | Summary |
| :--- | :--- |
| RR-19 | Retail Service & Outbox Pattern |
| RR-20 | Warehouse Service & Saga Participant |
| RR-21 | Saga Rollback & Compensations |
| RR-22 | Trace Service & CQRS |

Chi tiết: [`docs/business/sprint4/main.md`](./sprint4/main.md)
