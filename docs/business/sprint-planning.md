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
| RR-1 | Infra Kick-off — Monorepo, Go Workspaces | ✅ Done |
| RR-2 | Config & Logger — Viper, Zap, BaseConfig | ✅ Done |
| RR-3 | Base & Errs — RFC 9457, app lifecycle | ✅ Done |
| RR-4 | Proto + Wire + KrakenD + Client Shell | 🕒 To Do |
| RR-5 | Docker Complete — SigNoz + Full Compose | 🕒 To Do |
| RR-6 | Demo Service — Clean Arch + Full Stack Slice | 🕒 To Do |
| RR-7 | Client App — Control Plane | 🕒 To Do |
| RR-8 | Client App — Business UI Shell | 🕒 To Do |

Chi tiết: [`docs/business/sprint1/main.md`](./sprint1/main.md)

---

## Sprint 2 — Security & Access Control

**Goal:** Thiết lập nền tảng AuthN/AuthZ toàn diện. Mọi API call phải được bảo mật trước khi viết business logic.

| Ticket | Summary |
| :--- | :--- |
| RR-9 | Identity Server Infrastructure (Ory Kratos) |
| RR-10 | Client-Side Auth & Login Flow |
| RR-11 | Backend Security Core - JWT Validation |
| RR-12 | Fine-grained Authorization (Casbin) |
| RR-13 | End-to-End Secure Integration |

Chi tiết: [`docs/business/sprint2/main.md`](./sprint2/main.md)

---

## Sprint 3 — Farm Service Business Logic

**Goal:** Farm Service đầy đủ CRUD. Copy boilerplate từ `apps/demo-service/`, chỉ viết business logic.

| Ticket | Summary |
| :--- | :--- |
| S3-1 | Farm Service — Bootstrap từ demo-service template |
| S3-2 | Farm Domain — Entities & DB Schema |
| S3-3 | Farm Repository — Full CRUD (sqlx) |
| S3-4 | Farm UseCase — Business Rules |
| S3-5 | Farm Delivery — gRPC + REST handlers |
| S3-6 | Farm UI — List + Create Farm (`/app/farms`) |
| S3-7 | Farm UI — Batch Management |

---

## Sprint 4 — Distributed Transactions (Saga)

**Goal:** Retail Service + Warehouse Service + Saga Choreography.

| Ticket | Summary |
| :--- | :--- |
| S3-1 | Retail Service — POST /orders, Outbox Pattern |
| S3-2 | Warehouse Service — Reserve inventory, Saga participant |
| S3-3 | Saga Rollback — Compensating actions |
| S3-4 | Trace Service — Kafka consumer → Elasticsearch CQRS |

---

## Sprint 4 — Security & Auth

**Goal:** Production-grade security.

| Ticket | Summary |
| :--- | :--- |
| S4-1 | Ory Kratos — Identity, JWT, JWKS |
| S4-2 | In-service JWT validation (in-memory JWKS) |
| S4-3 | Casbin — RBAC/ABAC per service |
| S4-4 | mTLS — gRPC internal certificates |
