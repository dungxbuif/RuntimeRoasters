# [RR-1] Infrastructure Kick-off

**User Story:**
As a **DevOps Engineer**, I want to set up a centralized infrastructure configuration so that the entire development team can launch all necessary components with a single command.

**Business Context:**
Consistency in the development environment (Local Consistency) is a vital factor in minimizing "Works on my machine" errors. A solid infrastructure helps accelerate progress in subsequent Sprints.

---

## 🔍 Scope
Build a `docker-compose.yaml` system in the `deployments/` directory, including:
- **Persistence:** PostgreSQL (with automated database initialization scripts), Valkey (Redis), Cassandra.
- **Messaging:** Kafka Cluster & Kafka UI.
- **Observability:** Jaeger, Elasticsearch.
- **API Gateway:** KrakenD.

---

## 🛠️ Tech Rules
- Use **Docker Networks** to isolate traffic between layers.
- **Volume Persistence:** Must mount volumes to the host to prevent data loss when containers restart.
- **Healthchecks:** Every container must have a healthcheck script to ensure correct startup order.

---

## ✅ Acceptance Criteria (AC)
1. **End-to-End Startup:** Running `docker compose up -d` without errors, all services are `healthy`.
2. **Database Auto-init:** Postgres automatically creates `farm_db`, `auth_db`, and `warehouse_db` upon first startup.
3. **UI Accessibility:** Kafka UI (:8080) and Jaeger UI (:16686) are accessible from a browser.
