# Sprint 1: Complete Infrastructure — Kanban Board

**Status:** `IN_PROGRESS` | **Goal:** Toàn bộ infra sẵn sàng. Sprint 2 chỉ viết business logic.

---

## Kanban Board

| 🕒 To Do | 🚧 In Progress | ✅ Done |
| :--- | :--- | :--- |
| [RR-4: Proto + Wire + KrakenD + Apps Shell](./RR-4.md) | [RR-4-1: buf toolchain](./RR-4/RR-4-1/ticket.md) | [RR-1: Infra Kick-off](./RR-1.md) |
| [RR-5: Docker Complete — Jaeger](./RR-5.md) | | [RR-2: Config & Logger](./RR-2.md) |
| [RR-6: Demo Service — Clean Arch + Full Slice](./RR-6.md) | | [RR-3: Base & Errs](./RR-3.md) |
| [RR-7: Control App — Jaeger Tracing](./RR-7.md) | | |
| [RR-8: Client App — Business UI + API Explorer](./RR-8.md) | | |

---

## Sprint Goal

Sau Sprint 1, toàn bộ infrastructure chạy được end-to-end:
- `docker compose up -d` khởi động mọi thứ (Postgres, Redis, Jaeger, KrakenD)
- `GET localhost:8081/v1/demo/ping` trả về JSON qua KrakenD → grpc-gateway → gRPC handler
- Trace visible trong Jaeger tại `localhost:16686` (Control Plane)
- API Explorer (Swagger) tích hợp trong `client-app`, load contract qua KrakenD
- `apps/demo-service/` là canonical template — mọi service Sprint 2+ copy từ đây

---

## Dependency Order

```
RR-4 (proto + wire + krakend + app shells)
  └─ RR-4-1: buf toolchain
  └─ RR-4-2: grpc-gateway scaffold
  └─ RR-4-3: Google Wire
  └─ RR-4-4: KrakenD config + apps shell (client)
  └─ RR-4-5: GetDemo handler (E2E verify)

RR-5 (docker complete) — Jaeger setup

RR-6 (demo-service) — Clean Arch + OTel

RR-7 (control-app: Jaeger) — Observability Dashboard

RR-8 (client-app: Business UI) — Layout + API Explorer
```
