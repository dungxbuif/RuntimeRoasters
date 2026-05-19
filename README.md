# Runtime Roasters — Farm-to-Cup Coffee Supply Chain Platform

[English](#english) | [Tiếng Việt](#tiếng-việt)

---

<a name="english"></a>

## ☕ The Project
**Runtime Roasters** is a production-grade **Microservices** ecosystem built in `Golang`. It digitizes the entire coffee lifecycle—from plantation harvest to retail delivery—ensuring absolute transparency and traceability.

### 🚀 Technical Showcase (The "Wow" Factor)
This project implements sophisticated distributed patterns used by top-tier engineering teams:

*   **Choreographed Sagas:** Automated multi-service transactions (Order ➔ Payment ➔ Warehouse) with failure compensation.
*   **Transactional Outbox & Inbox:** Guarantees "Exactly-Once" processing and eventual consistency via Kafka.
*   **CQRS Implementation:** Segregating high-volume writes (Postgres) from high-speed traceability reads (Elasticsearch).
*   **Real-time Geo-Tracking:** Live GPS processing using Valkey and OSRM road-accurate routing.
*   **Zero-Trust Security:** mTLS, OIDC (Ory Kratos/Hydra), and decentralized Casbin authorization.
*   **Full Observability:** W3C-compliant distributed tracing across all services (OpenTelemetry + Jaeger).

---

### 📖 Documentation Deep Dive
For a detailed technical evaluation, please refer to:
👉 **[SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md)**: Architecture blueprints, Data flows, and ADRs (Decision Records).

---

### 🏗️ Tech Stack
*   **Backends:** Go 1.22+ (Clean Architecture, gRPC, Gin)
*   **Frontends:** Next.js 15 (TypeScript, Tailwind, Leaflet)
*   **Messaging:** Apache Kafka (Event-driven core)
*   **Identity:** Ory Kratos & Ory Hydra (SSO/Identity)
*   **Storage:** PostgreSQL, Valkey (Redis alternative), Elasticsearch, Cassandra
*   **Infrastructure:** KrakenD (API Gateway), Docker Compose, OpenTelemetry

---

### ⏱️ Quick Start (Development)
```bash
# 1. Provision the cluster
docker-compose -f deployments/docker-compose.dev.yaml up -d

# 2. Seed the data
./deployments/seed.sh

# 3. Access the Dashboard
# Open http://localhost:3000
```

---

<a name="tiếng-việt"></a>

## 🇻🇳 Tiếng Việt

**Runtime Roasters** là một nền tảng quản lý chuỗi cung ứng cà phê "từ nông trường đến tách ly". Dự án được xây dựng 100% bằng `Golang` nhằm trình diễn các kỹ thuật kiến trúc Microservices hiện đại nhất.

### Điểm nhấn Kỹ thuật:
- **Saga Pattern:** Quản lý giao dịch liên dịch vụ tự động.
- **Transactional Outbox:** Đảm bảo dữ liệu không bao giờ bị mất giữa các Service.
- **CQRS:** Tách biệt luồng ghi (Postgres) và luồng đọc truy xuất nhanh (Elasticsearch).
- **Real-time Tracking:** Giám sát vị trí tài xế thời gian thực trên bản đồ.

---
*Developed for Portfolio & Engineering Showcase.*
