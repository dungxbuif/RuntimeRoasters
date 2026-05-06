# Sprint 2: Security & Access Control — Kanban Board

**Status:** `PLANNING` | **Goal:** Thiết lập nền tảng AuthN/AuthZ toàn diện. Mọi API call phải được bảo mật trước khi viết business logic.

---

## 🎯 Sprint Goal

> **Người dùng có thể đăng nhập qua Identity Server, nhận JWT và các Microservices có thể tự động xác thực/phân quyền (AuthN/AuthZ) độc lập bằng Casbin mà không làm tăng độ trễ hệ thống.**

---

## 📋 Kanban Board

| 🕒 To Do                                                    | 🚧 In Progress | ✅ Done |
| :---------------------------------------------------------- | :------------- | :------ |
| [RR-9: Identity Server Infrastructure](./RR-9/RR-9.md)           |                |         |
| [RR-10: Client-Side Auth & Login Flow](./RR-10.md)          |                |         |
| [RR-11: Backend Security Core - JWT Validation](./RR-11.md) |                |         |
| [RR-12: Fine-grained Authorization (Casbin)](./RR-12.md)    |                |         |
| [RR-13: End-to-End Secure Integration](./RR-13.md)          |                |         |

---

## 🛤️ Dependency Flow

```mermaid
graph TD
    RR9[RR-9: Identity Server] --> RR10[RR-10: Client Login]
    RR9 --> RR11[RR-11: Backend JWT Middleware]
    RR11 --> RR12[RR-12: Casbin Authorization]
    RR10 --> RR13[RR-13: E2E Verification]
    RR12 --> RR13
```

---

## 🛠️ Technical Stack & Prep

- **Identity Server:** Ory Kratos.
- **Backend Core:** `pkg/base/auth` (JWT v5) + `pkg/base/casbin`.
- **Visualization:** SigNoz API integration prep (ensuring `user.id` is in spans).
- **Documentation:** [Security Architecture](../architecture/security.md), [Control Plane Visualization Strategy](../architecture/control-plane-visualization.md).
