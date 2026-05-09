# RR-12: Technical Design — Fine-grained Authorization (Casbin)

**Branch:** `feat/RR-12`
**Status:** `DRAFT`
**Author:** Tech Lead
**Date:** 2026-05-09

---

## 1. Context & Goal

Sau khi **RR-11** hoàn thiện việc xác thực danh tính (Authentication), hệ thống đã biết "Ai đang gọi". Nhiệm vụ của **RR-12** là xác định "Họ được làm gì" (Authorization).

Mục tiêu cốt lõi:
- Xây dựng cơ chế phân quyền tập trung tại tầng `pkg/base`.
- Áp dụng mô hình **Hybrid RBAC+ABAC**: RBAC tại biên (Gateway/Middleware) và ABAC (Data Ownership) tại tầng Repository.
- Đảm bảo tính mở rộng cao, không bị thắt cổ chai khi số lượng resource tăng lên hàng triệu.

---

## 2. Architecture: The Two-Gate Model

Chúng ta chia quyền truy cập thành 2 cửa khẩu kiểm soát:

### Gate 1: RBAC (Role-Based Access Control)
- **Công cụ:** Casbin Enforcer.
- **Vị trí:** `Gin Middleware` & `gRPC Interceptor`.
- **Logic:** Kiểm tra `(User.Role, API_Method, Action)`. 
- **Mục đích:** Chặn các truy cập sai vai trò ngay từ vòng gửi xe (ví dụ: `driver` không được phép gọi hàm `CreateFarm`).

### Gate 2: Data Scoping (ABAC / Ownership)
- **Công cụ:** SQL `WHERE` clauses.
- **Vị trí:** `Repository Layer`.
- **Logic:** Ép thêm điều kiện `owner_id = caller_id` vào mọi câu lệnh Query/Command.
- **Mục đích:** Đảm bảo dù `farmer` có quyền gọi hàm `ListFarms`, họ cũng chỉ thấy được Nông trại của chính mình.

---

## 3. Casbin Configuration

### 3.1 Model Definition (`model.conf`)
Sử dụng mô hình RBAC có phân cấp (Hierarchy).

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
# Admin có toàn quyền, các role khác check theo chính sách và phân cấp role
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))
```

### 3.2 Policy Format (`policy.csv`)
Sử dụng gRPC Method Name làm định danh Resource (`obj`) để thống nhất cho cả gRPC và Gateway.

```csv
# Roles Hierarchy
g, admin, farmer
g, admin, processor

# Permissions
p, farmer, /farm.v1.FarmService/CreateFarm, write
p, farmer, /farm.v1.FarmService/ListFarms, read
p, processor, /farm.v1.FarmService/ListFarms, read
```

---

## 4. Package Structure (`pkg/base/casbin`)

```text
src/pkg/base/casbin/
├── config.go           # Cấu hình đường dẫn model/policy hoặc Adapter DB
├── enforcer.go         # Khởi tạo Casbin Enforcer (Singleton-ish)
├── middleware_gin.go   # Gin AuthZ Middleware
└── interceptor_grpc.go # gRPC AuthZ Interceptor
```

---

## 5. Integration Workflow

1. **Extraction:** Middleware lấy `identity.Claims` từ context (đã được RR-11 inject).
2. **Evaluation:** Gọi `enforcer.Enforce(claims.Role, fullMethod, action)`.
3. **Action mapping:**
    - `GET` -> `read`
    - `POST` / `PUT` / `PATCH` -> `write`
    - `DELETE` -> `delete`
4. **Decision:** 
    - `true` -> `context.Next()`
    - `false` -> Return `errs.ErrForbidden` (403 Forbidden).

---

## 6. Data Scoping & ABAC Implementation

### 6.1 Scoping Claims Helper
- **Package:** `pkg/base/casbin`
- **Logic:** Một helper tập trung chuyển đổi kết quả từ Casbin policy (ví dụ: `region: dak_lak`) thành các typed SQL filter structs.
- **Mục đích:** Tránh việc viết SQL raw rải rác trong các handler, đảm bảo tính nhất quán khi áp dụng bộ lọc dữ liệu tại tầng Repository.

### 6.2 Two-Gate Context
Hệ thống sử dụng cơ chế bảo mật 2 lớp độc lập:
- **Gate 1 (App Scopes):** Được đảm nhiệm bởi KrakenD API Gateway (RR-10). Kiểm tra các scope ứng dụng trên JWT.
- **Gate 2 (User Roles):** Được đảm nhiệm bởi Casbin (RR-12) tại mức service. Kiểm tra vai trò của người dùng trên từng resource cụ thể.

---

## 7. Guidelines cho Developer

1. **Tầng UseCase:** Tuyệt đối không viết `if role == "..."`. Hãy giả định rằng khi request vào đến đây, user đã có quyền thực hiện hành động đó.
2. **Tầng Repository:** Luôn nhận `callerID` làm tham số bắt buộc cho các hàm Query.
3. **Thêm quyền mới:** Chỉ cần cập nhật file `policy.csv` (Dev) hoặc bảng `casbin_rule` (Prod), không cần sửa code service.

---

## 8. Implementation Tips

- **Testing:** Sử dụng [Casbin Online Editor](https://casbin.org/editor/) để debug model và matchers cực nhanh trước khi cập nhật vào code.
- **Admin Bypass:** Luôn đặt điều kiện kiểm tra `admin` (`g(r.sub, "admin")`) ở đầu matcher để đảm bảo quản trị viên không bị khóa quyền do lỗi cấu hình chính sách.
- **Middleware Ordering:** Đảm bảo Casbin Middleware được đăng ký chạy **SAU** JWT Validation Middleware để đảm bảo thông tin định danh (Identity) đã sẵn sàng trong context.

---

## 9. Risks & Mitigation

| Rủi ro | Giải pháp |
|---|---|
| **Performance:** Enforce mỗi request làm tăng latency. | Casbin Enforcer giữ policy trong memory, check cực nhanh. Với số lượng role nhỏ, latency < 1ms. |
| **Stale Policy:** Khi đổi quyền trong DB, service không nhận được ngay. | Sử dụng cơ chế `Watcher` của Casbin (Redis Pub/Sub) để notify các pod refresh cache khi policy thay đổi. |
| **Ambiguous Actions:** Không biết map API nào vào `read` hay `write`. | Thống nhất chuẩn đặt tên Method và mapping tại mức Middleware. |
