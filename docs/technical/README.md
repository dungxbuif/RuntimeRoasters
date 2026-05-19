# 🛠️ Runtime Roasters Technical Knowledge Base

This directory contains deep-dive technical specifications for the Runtime Roasters platform. These documents are intended for senior engineers, architects, and system administrators.

## 📁 Technical Master Reference

1.  **⚙️ [Configuration & Environment](./CONFIG_REFERENCE.md)**
    *   Fail-fast mechanism, 12-factor app config, and global environment reference.
2.  **🗄️ [Data Models & Persistence](./DATABASE_MODELS.md)**
    *   GORM models, ER diagrams, Transactional Outbox/Inbox schemas, and storage boundaries for Postgres JSONB, Elasticsearch, Cassandra, and Valkey.
3.  **🔗 [API Contracts & Messaging](./API_CONTRACTS.md)**
    *   gRPC service definitions, REST-to-gRPC mapping, and Kafka event schemas.
4.  **🌊 [Core Sequences & Rationale](./FLOW_SEQUENCES.md)**
    *   Sequence diagrams and the "Why" behind Sagas, CQRS, and Distributed AuthZ.
5.  **🎨 [Frontend Development Guide](./FRONTEND.md)**
    *   Next.js patterns, UI/UX standards, and OIDC integration.
6.  **🚛 [Logistics Simulation Guide](./LOGISTICS_SIMULATION.md)**
    *   Driver Client route simulation, realtime GPS updates, and OSRM route generation.
7.  **🐝 [Kafka Engineering Conventions](./KAFKA.md)**
    *   Partitioning, idempotency, and high-availability configuration.

## 📏 Engineering Standards
- **[Code Style & Canonical Template](./standards/CODE_STYLE.md)**
- **[Modular Bootstrap Pattern](./standards/BOOTSTRAP.md)**

---
*Maintained by the Runtime Roasters Engineering Team.*
