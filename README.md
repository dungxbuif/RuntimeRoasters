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

👉 **[Master System Architecture Specification](./docs/architecture/SYSTEM_ARCHITECTURE.md)**: The definitive technical guide to the platform's design.

### 🗺️ Service Port Map

| Service | REST Port | gRPC Port | Description |
| :--- | :--- | :--- | :--- |
| **KrakenD** | 8081 | - | API Gateway (Entry Point) |
| **Auth** | 8082 | 50052 | Identity, Session & Policy Management |
| **Farm** | 8083 | 50053 | Plantation & Harvest Management |
| **Retail** | 8084 | 50054 | Orders & Customer Transactions |
| **Warehouse**| 8085 | 50055 | Inventory & Stock Reservation |
| **Logistics**| 8087 | 50057 | Driver Tracking & Delivery |
| **Payment**  | 8088 | 50058 | Simulated Payment Gateway (Stripe/VNPay) |
| **Trace**    | 8089 | 50059 | Supply Chain Traceability (CQRS) |
| **Audit**    | 8091 | 50061 | Immutable Event Logging |

### 📖 Detailed Guides
- 📘 **[Developer Guide](./docs/DEVELOPER_GUIDE.md)**: Deep dive into patterns, coding standards, and service interactions.
- 🛠️ **[Technical Knowledge Base](./docs/technical/README.md)**: Exhaustive details on Config, Data Models, and API Contracts.
- 🏢 **[Business Specification](./docs/domain/BUSINESS_SPECIFICATION.md)**: Master domain logic and handover-ready requirements.
- 🌊 **[Architecture Flows](./docs/architecture/FLOWS.md)**: Visualizing identity and data consistency.
- 🎨 **[Frontend Development Guide](./docs/technical/FRONTEND.md)**: Next.js patterns and UI/UX standards.
- 📋 **[Product Requirements](./docs/business/product-requirements.md)**: Business vision and feature specifications.
- 📅 **[Sprint Roadmap](./docs/business/sprint-planning.md)**: Development phases and execution history.
- 📜 **[Architecture Decisions (ADRs)](./architecture/adrs/)**: History of critical technical choices.

---
*Developed by [Dung Bui](https://github.com/dungxbuif) for Engineering Showcase.*
