# Project Port Allocation Registry

Tài liệu này quy hoạch toàn bộ Port cho các dịch vụ trong hệ thống Runtime Roasters để tránh xung đột và dễ dàng quản lý trong môi trường Docker/Local.

## 1. Gateway & Auth Infrastructure

| Service | HTTP Port | gRPC Port | Admin/Other Port |
| :--- | :--- | :--- | :--- |
| **KrakenD API Gateway** | `8080` | - | `8090` (Debug) |
| **Ory Hydra (OIDC)** | `4445` (Public) | - | `4444` (Admin) |
| **Ory Kratos (Identity)** | `4433` (Public) | - | `4434` (Admin) |
| **MailSlurp/Mailhog** | `1025` (SMTP) | - | `8025` (UI) |

## 2. Business Microservices (Go Apps)

Dải Port gRPC bắt đầu từ `50051`. Dải Port HTTP bắt đầu từ `8081`.

| Service | HTTP Port | gRPC Port | Description |
| :--- | :--- | :--- | :--- |
| **Auth Service** | `8081` | `50051` | Session & RBAC Sync |
| **Farm Service** | `8082` | `50052` | Harvesting & Farm Mgmt |
| **Warehouse Service** | `8083` | `50053` | Processing & Inventory |
| **Retail Service** | `8084` | `50054` | Ordering & POS |
| **Logistics Service** | `8085` | `50055` | Delivery & Tracking |
| **Payment Service** | `8086` | `50056` | Stripe Integration |
| **Trace Service** | `8087` | `50057` | CQRS & Traceability |
| **Audit Service** | `8088` | `50058` | Immutable Logs |

## 3. Infrastructure & Databases

| Component | Port | Description |
| :--- | :--- | :--- |
| **PostgreSQL** | `5432` | Shared instance (Multi-DB) |
| **Valkey (Redis)** | `6379` | Cache & Distributed Lock |
| **Kafka** | `9092` | Message Broker |
| **Zookeeper** | `2181` | Kafka Orchestration |
| **SigNoz (UI)** | `3301` | Observability Dashboard |
| **Elasticsearch** | `9200` | Trace Search Engine |
| **Cassandra** | `9042` | Audit Log Storage |

## 4. Frontend Applications

| App | Port | Description |
| :--- | :--- | :--- |
| **Client App (Next.js)** | `3000` | Dashboard & Storefront |
| **Admin Panel** | `3001` | Internal Management |

---
**Ghi chú:** Khi triển khai Docker Compose, các service sẽ gọi nhau qua Service Name (Internal DNS), nhưng Port mapping ra Host máy tính cá nhân sẽ tuân thủ bảng này.
