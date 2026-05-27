# ☕ Runtime Roasters — Farm-to-Cup Coffee Supply Chain Platform

## 🌟 Overview
**Runtime Roasters** is a production-grade **Microservices** ecosystem built with **Golang**. It digitizes the entire coffee lifecycle—from plantation harvest to retail delivery—ensuring absolute transparency, traceability, and efficiency in the supply chain.

This project serves as a **technical showcase** for modern distributed system patterns, demonstrating how to build resilient, scalable, and observable applications.

### 🚀 Technical Highlights
*   **Event-Driven Architecture:** Core communication via **Apache Kafka** with **Transactional Outbox/Inbox** patterns for guaranteed consistency.
*   **Advanced Distributed Patterns:** Implementation of **Choreographed Sagas** for multi-service transactions (Order ➔ Payment ➔ Warehouse).
*   **Identity & Security:** Zero-trust model using **Ory Kratos** (Identity), **Ory Hydra** (OAuth2/OIDC), and decentralized **Casbin** authorization.
*   **High Performance Storage:** Polyglot persistence using **PostgreSQL** (ACID), **Elasticsearch** (Search/CQRS), **Cassandra** (Audit Logs), and **Valkey** (Caching/Tracking).
*   **Real-time Capabilities:** Geo-accurate driver tracking and automated logistics routing.
*   **Full Observability:** Distributed tracing (OpenTelemetry), structured logging, and health monitoring.

---

## 🛠️ Development Setup

### 📋 Prerequisites
Ensure you have the following installed:
- **Go 1.22+**
- **Docker & Docker Compose**
- **Node.js 20+** (for Frontend)
- **Task** (Optional, but recommended: `go install github.com/go-task/task/v3/cmd/task@latest`)
- **Air** (For hot-reload: `go install github.com/air-verse/air@latest`)

### 🏃 Quick Start
Follow these steps to get the system running locally:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-repo/RuntimeRoasters.git
    cd RuntimeRoasters
    ```

2.  **Environment Configuration:**
    ```bash
    cp .env.example .env
    ```

3.  **Start Infrastructure (Database, Kafka, Identity):**
    ```bash
    task infra
    ```
    *Wait for all containers to be healthy (approx. 30-60s).*

4.  **Seed Initial Data (Admin User & OAuth2 Clients):**
    ```bash
    task seed
    ```

5.  **Generate Protobuf & gRPC Code:**
    ```bash
    task proto
    ```

6.  **Start Services:**
    - To start the full backend + frontend: `task dev`
    - To start only backend services: `task be`
    - Individual services can be started via `air` in their respective `src/apps/` folders.

---

### 🏗️ Architecture & Documentation

The system follows **Clean Architecture** principles to keep business logic isolated from infrastructure concerns.

👉 **[Master System Architecture Specification](./docs/product/TECH.md)**: The definitive technical guide to the platform's design.

### 🗺️ Service Port Map

| Service | REST Port | gRPC Port | Description |
| :--- | :--- | :--- | :--- |
| **KrakenD** | 8081 | - | API Gateway (Entry Point) |
| **Auth** | 8082 | 50052 | Identity, Session & Policy Management |
| **Farm** | 8083 | 50053 | Plantation & Harvest Management |
| **Retail** | 8084 | 50054 | Orders & Customer Transactions |
| **Logistics** | 8085 | 50055 | Driver Tracking & Delivery |
| **Payment** | 8086 | 50056 | Simulated Payment Gateway (Stripe/VNPay) |
| **Trace** | 8087 | 50057 | Supply Chain Traceability (CQRS) |
| **Audit** | 8088 | 50058 | Immutable Event Logging |
| **Warehouse** | 8089 | 50059 | Inventory, Pickup, Processing & Stock Reservation |
| **Socket** | 8091 | 50060 | WebSocket realtime topology/demo fanout |

Infrastructure UIs and endpoints:

| Component | Port | Description |
| :--- | :--- | :--- |
| **Client App** | 3000 | Next.js control-plane UI |
| **Kafka UI** | 8090 | Kafka topic/browser UI |
| **SigNoz** | 3301 | Observability UI backed by ClickHouse |
| **Kibana** | 5601 | Elasticsearch UI |
| **Postgres** | 54321 | Shared local Postgres host port |
| **Elasticsearch** | 9200 | Trace/search read model |
| **Cassandra** | 9042 | Audit log storage |
| **Valkey** | 6379 | Cache/GPS tracking |
| **OTLP** | 4317 / 4318 | OpenTelemetry gRPC/HTTP |
| **Kratos Public** | 4433 | Identity public/proxy endpoint |
| **Hydra Public** | 4444 | OAuth2/OIDC public endpoint |

### 📖 Detailed Guides
- 📘 **[Developer Guide](./docs/product/GUIDE.md)**: Deep dive into patterns, coding standards, and service interactions.
- ▶️ **[Demo Setup Runbook](./docs/product/DEMO_SETUP_RUNBOOK.md)**: Install, start order, seed/reset data, ports, and playground entrypoints.
- 🛠️ **[Technical Specification](./docs/product/TECH.md)**: Architecture, contracts, storage, messaging, and infrastructure.
- 📜 **[Business Specification](./docs/product/domain/README.md)**: BA-facing flow, roles, UI actions, and business rules.
- 🌊 **[Product/System Specification](./docs/product/SPEC.md)**: Product scope, business context, and reference requirements.
- 🎨 **[UI/UX Design Notes](./docs/product/ui-ux/DESIGN.md)**: Next.js dashboard visual language and UX standards.
- 📅 **[Sprint Roadmap](./docs/stories/ROADMAP.md)**: Development phases and execution history.
- 📜 **[Architecture Decisions (ADRs)](./docs/decisions/)**: History of critical technical choices.

---
*Developed by [Dung Bui](https://github.com/dungxbuif) for Engineering Showcase.*
