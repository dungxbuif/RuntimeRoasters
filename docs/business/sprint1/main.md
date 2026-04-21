# Sprint 1: Complete Infrastructure — Kanban Board

**Status:** `IN_PROGRESS` | **Goal:** Toàn bộ infra sẵn sàng. Sprint 2 chỉ viết business logic.

---

## Kanban Board

| 🕒 To Do | 🚧 In Progress | ✅ Done |
| :--- | :--- | :--- |
| [RR-4: Proto + Wire + KrakenD + Client Shell](./RR-4.md) | | [RR-1: Infra Kick-off](./RR-1.md) |
| [RR-5: Docker Complete — SigNoz](./RR-5.md) | | [RR-2: Config & Logger](./RR-2.md) |
| [RR-6: Demo Service — Clean Arch + Full Slice](./RR-6.md) | | [RR-3: Base & Errs](./RR-3.md) |
| [RR-7: Client App — Control Plane](./RR-7.md) | | |
| [RR-8: Client App — Business UI Shell](./RR-8.md) | | |

---

## Sprint Goal

Sau Sprint 1, toàn bộ infrastructure chạy được end-to-end:
- `docker compose up -d` khởi động mọi thứ
- `GET localhost:8081/v1/demo/ping` trả về JSON qua KrakenD → grpc-gateway → gRPC handler
- Trace visible trong SigNoz tại `localhost:3000/control/observability`
- `apps/demo-service/` là canonical template — mọi service Sprint 2+ copy từ đây

---

## Dependency Order

```
RR-4 (proto + wire + krakend)
  └─ RR-4-1: buf toolchain
  └─ RR-4-2: grpc-gateway scaffold
  └─ RR-4-3: Google Wire
  └─ RR-4-4: KrakenD config + client-app shell
  └─ RR-4-5: GetDemoFarm handler (E2E verify)

RR-5 (docker complete) — có thể song song với RR-4

RR-6 (demo-service) — depends on RR-4, RR-5
  Full slice: domain → usecase → repo → delivery
  OTel: traces, metrics, logs, Kafka propagation

RR-7 (client control plane) — depends on RR-5, RR-6
  /control/services, /control/api-explorer, /control/observability

RR-8 (client business UI) — depends on RR-7
  Layout + placeholder sections
```
