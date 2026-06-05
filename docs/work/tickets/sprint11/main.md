# Sprint 11: Observability, Security & Project Polish

**Epic Goal:** Comprehensive system monitoring, Zero Trust security, performance tuning, and final project polish.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-31](./RR-31.md) | [Tech] Audit Service (Cassandra) | 🕒 To Do | Tech Lead |
| [RR-32](./RR-32.md) | [Tech] System Health Dashboards (SigNoz) | 🕒 To Do | DevOps |
| [RR-33](./RR-33.md) | [Tech] mTLS Zero Trust | 🕒 To Do | Security Eng |
| [RR-34](./RR-34.md) | [Tech] Token Revocation & Security Audit | 🕒 To Do | Tech Lead |
| [RR-35](./RR-35.md) | [Tech] Load Testing & Performance Tuning | 🕒 To Do | SRE |
| [RR-36](./RR-36.md) | [Tech] Final Bug Fixes & Refactoring | 🕒 To Do | Team |
| [RR-37](./RR-37.md) | [BA] Final Demo Preparation | 🕒 To Do | PO |
| [RR-38](./RR-38.md) | [Tech] Documentation & Handover | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Cassandra:** Store massive amounts of Audit Log data with extremely fast write capabilities.
- **SigNoz/OTEL:** Collect Traces and Metrics to monitor the performance of each service.
- **mTLS:** Encrypt communication between Microservices.
- **Valkey:** Used to manage the revoked token list (Blacklist).
- **Chaos Engineering:** Test system resilience when a service dies.
- **Code Quality:** Review và refactor những phần code còn chưa tối ưu.
- **Documentation:** Hoàn thiện README, API Docs (Swagger/Buf).
- **Handover:** Chuẩn bị tài liệu vận hành.
