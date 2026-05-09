# [RR-12] Fine-grained Authorization (Casbin)

- **Summary:** Tích hợp Casbin với mô hình Hybrid RBAC+ABAC: Casbin kiểm tra Role (Gate 2), Repository filter quyền sở hữu dữ liệu bằng SQL. Hỗ trợ Role Hierarchy và Data Scopes.
- **Priority:** `MEDIUM`
- **Type:** Infrastructure
- **ADR Reference:** `ADR-2026-05-04-AUTH` — Session 5, 6, 7, 8

---

## 📖 User Story
> As a security officer, I want fine-grained access control that enforces both role-based permissions at the service boundary and data ownership at the repository level — so that a user can only read/write data they are authorized for, even if they hold a high-privilege role in another domain.

## 🔍 Acceptance Criteria

### Scenario 1: RBAC Middleware — Role Check (Gate 2)
- **Given:** A request passes JWT validation (RR-11) and carries a `role` claim.
- **When:** The Casbin middleware evaluates the request.
- **Then:** It checks `(role, resource, action)` against the loaded policy.
- **And:** Authorized requests pass through; unauthorized requests return `403 Forbidden`.

### Scenario 2: Least Privilege — Role Separation
- **Given:** Casbin policy is loaded with the defined role set.
- **When:** A user with role `processor` attempts to access a `warehouse_mgr`-only endpoint.
- **Then:** It returns `403 Forbidden` (roles are not interchangeable even if domains overlap).

  **Defined Roles:** `admin`, `farmer`, `processor`, `warehouse_mgr`, `driver`, `store_mgr`, `guest`.

### Scenario 3: Role Hierarchy (Admin Inherits All)
- **Given:** A user with role `admin`.
- **When:** They access any endpoint that any other role can access.
- **Then:** Casbin's `g(sub, role)` hierarchy grants access automatically without duplicating policies.

### Scenario 4: Unauthorized Action — Wrong Role
- **Given:** A user with role `farmer` attempts `DELETE /v1/demo`.
- **When:** The Casbin enforcer checks the request.
- **Then:** It returns `403 Forbidden`.

### Scenario 5: Hybrid ABAC — Data Ownership Enforcement
- **Given:** A user with role `farmer` requests a list of farms.
- **When:** The UseCase calls the Repository.
- **Then:** The Repository applies ownership SQL filtering: `WHERE user_id = $1` (or `WHERE org_id = $1`).
- **And:** The user cannot see another farmer's data even if their role would allow the action type.

### Scenario 6: Data Scopes — Scoping Claims from Casbin
- **Given:** A user's Casbin policy defines a data scope (e.g., `region: dak_lak`).
- **When:** The middleware runs the enforcer.
- **Then:** The enforcer returns Scoping Claims to the UseCase.
- **And:** The Repository uses these claims to generate a SQL `WHERE region = $1` clause via a centralized Helper, not ad-hoc per handler.
