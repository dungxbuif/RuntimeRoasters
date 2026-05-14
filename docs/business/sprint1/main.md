# Sprint 1: Infrastructure Foundation & Baseline

**Trạng thái:** ✅ Hoàn thành (Completed)
**Mục tiêu:** Thiết lập nền tảng hạ tầng kỹ thuật (Infrastructure) vững chắc, đảm bảo môi trường phát triển đồng nhất cho toàn bộ hệ sinh thái Microservices.

---

## 📋 Trạng thái Ticket (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-1](./RR-1/ticket.md) | [Tech] Infrastructure Kick-off (Docker Compose) | ✅ Done | DevOps |
| [RR-2](./RR-2/ticket.md) | [Tech] Base Framework: Config & Logger | ✅ Done | Tech Lead |
| [RR-3](./RR-3/ticket.md) | [Tech] Common Packages: Base & Errs | ✅ Done | Backend |
| [RR-4](./RR-4/ticket.md) | [Epic] Service Toolchain & API Gateway (KrakenD) | ✅ Done | Tech Lead |
| [RR-5](./RR-5/ticket.md) | [Tech] Observability Stack (Jaeger Setup) | ✅ Done | DevOps |
| [RR-6](./RR-6/ticket.md) | [Tech] Demo Service: Clean Architecture Template | ✅ Done | Backend |
| [RR-7](./RR-7/ticket.md) | [Tech] Control Plane UI: Jaeger Integration | ✅ Done | Frontend |
| [RR-8](./RR-8/ticket.md) | [BA] Client App: Core Shell & API Explorer | ✅ Done | Product Owner |

---

## 💡 Tầm nhìn Kỹ thuật (Technical Vision)
- **Zero-Config Onboarding:** Một lệnh `docker compose up` khởi động toàn bộ vũ trụ RuntimeRoasters.
- **Canonical Template:** `demo-service` là mẫu mực để các service sau copy theo (Clean Arch, DI, OTel).
- **Early Observability:** Tracing phải hoạt động ngay từ ngày đầu tiên để debug luồng gRPC.

## 📊 Kết quả đạt được (Sprint Result)
- Toàn bộ hạ tầng (Postgres, Valkey, Kafka, Jaeger) đã sẵn sàng.
- Luồng gRPC-Gateway qua KrakenD hoạt động ổn định.
- Framework chung (`pkg/`) đã được kiểm chứng qua Demo Service.
