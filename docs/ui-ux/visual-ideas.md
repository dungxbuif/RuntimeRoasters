# Thiết Kế Trực Quan UI: RuntimeRoasters Dashboard

## 1. Concept Tổng Thể

Client App là một codebase Next.js 15 phục vụ **hai nhóm người dùng** với hai pattern hoàn toàn khác nhau:

| Khu vực | Đối tượng | Pattern | Route |
| :--- | :--- | :--- | :--- |
| **Control Plane** | Admin / Dev | BFF (Route Handler proxy) | `/control/*` |
| **Business UI** | User / Operator | Direct (Browser → KrakenD) | `/app/*` |

---

## 2. Control Plane (Admin Dashboard)

Khu vực dành cho admin và developer để giám sát hệ thống. Mọi route `/control/*` được gate bởi `middleware.ts`.

### `/control/services` — Service Health Board

Grid hiển thị trạng thái real-time của từng service:

```
┌─────────────────────────────────────────────────┐
│  ● demo-service    Running   p50: 12ms           │
│  ● farm-service    Running   p50: 8ms            │
│  ● krakend         Running   p50: 2ms            │
│  ○ warehouse       Degraded  last seen: 2m ago   │
└─────────────────────────────────────────────────┘
```

Data source: BFF Route Handler gọi `GET /health/ready` của từng service, aggregate lại.

### `/control/api-explorer` — Swagger UI

Embed `swagger-ui-react` load `demo.swagger.json` (và các service swagger khác). Execute trực tiếp qua KrakenD `:8081`. Không cần container Swagger UI riêng.

### `/control/observability` — SigNoz Embed

SigNoz UI nhúng qua iframe — toàn bộ tính năng: Traces, Metrics, Logs, Dashboards.

```
┌─────────────────────────────────────────────────┐
│  /control/observability                          │
│  ┌───────────────────────────────────────────┐  │
│  │  <iframe src="/signoz/" />                │  │
│  │  SigNoz UI — Traces / Metrics / Logs      │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

Proxy flow option: `/signoz/*` → `next.config.ts` rewrite → `SIGNOZ_INTERNAL_URL` (internal Docker DNS). For local demo, SigNoz is also exposed directly at `http://localhost:3301`; admin-only embedding can be added later when the control page is implemented.

---

## 3. Business UI (Operations Dashboard)

Khu vực cho end-user và operator nghiệp vụ. Browser gọi thẳng KrakenD `:8081` — không qua BFF.

### Layout: Split-Screen (Sprint 2+)

```
┌────────────────────────────────────────────────────────────┐
│  Sidebar           │  Main Content Area                     │
│  ─────────         │  ──────────────────────────────────── │
│  Supply Chain      │  [Page content]                        │
│  > Farms           │                                        │
│  > Batches         │  Live Data Flow Visualization          │
│  > Logistics       │  (React Flow — Sprint 2+)              │
│  > Warehouse       │                                        │
│  > Retail          │                                        │
└────────────────────────────────────────────────────────────┘
```

### Sprint 1 Scope (RR-8)

- Layout chuẩn với Sidebar navigation
- Dashboard landing page với service status overview
- Placeholder sections: Supply Chain Overview, Recent Events, Quick Actions
- Ready để tích hợp Farm/Batch features từ Sprint 2

### Sprint 2+ — Farm & Batch Features

`/app/farms` — Danh sách nông trại, CRUD
`/app/batches` — Quản lý mẻ thu hoạch, timeline

### Sprint 3+ — Live Flow Visualization

Khi Saga pattern và Kafka hoàn thiện, thêm:

**React Flow Live Map** — Visualize luồng sự kiện theo thời gian thực:
- Nodes: các service (Farm, Process, Warehouse, Retail, Trace, Audit)
- Animated edges: event đang chạy (màu theo loại: xanh=success, đỏ=error, vàng=processing)
- Data source: Monitor Service subscribe Kafka → broadcast SSE → frontend

**Visual Metaphors (ngôn ngữ cà phê):**
- Hạt xanh: event từ Farm Service
- Hạt nâu: event từ Processing Service
- Xe tải: GPS ping từ Logistics
- Nhịp đập (pulse): Kafka cluster activity

**Chaos Control Panel** (demo scenarios):
- Kịch bản A: Ngắt Kafka → xem Outbox retry
- Kịch bản B: Hết hàng → xem Saga rollback
- Kịch bản C: QR trace → CQRS timeline replay

---

## 4. Navigation & Auth

```
/ (landing)
├── /login
├── /control/*  [admin-only — middleware.ts gate]
│   ├── /control/services
│   ├── /control/api-explorer
│   └── /control/observability
└── /app/*  [authenticated users]
    ├── /app/farms
    ├── /app/batches
    ├── /app/logistics
    └── /app/retail
```
