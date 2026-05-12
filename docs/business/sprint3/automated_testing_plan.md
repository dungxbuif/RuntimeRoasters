# Sprint 3: Automated Testing Plan (High-Fidelity)

Tài liệu này quy hoạch toàn bộ các kịch bản kiểm thử tự động (Unit & Integration) cho Sprint 3, đảm bảo hệ thống Nông trại hoạt động ổn định và bảo mật.

---

## 🏗️ 1. Backend Testing (Go)

### 1.1 Unit Tests (Dịch vụ)
- **Vị trí**: `src/apps/farm-service/internal/usecase/__tests__`
- **Mục tiêu**: Kiểm tra logic nghiệp vụ mà không cần DB thực.
- **Kịch bản**:
    - [ ] `CreateFarm`: Phải gán đúng `owner_id` và validate diện tích > 0.
    - [ ] `ListFarms`: Phải lọc đúng danh sách dựa trên Role (Admin thấy hết, Manager thấy của mình).
    - [ ] `SyncCasbin`: Đảm bảo gọi đúng Kratos API và nạp đủ quyền vào Casbin.

### 1.2 Integration Tests (Database)
- **Vị trí**: `src/apps/farm-service/internal/infrastructure/repository/__tests__`
- **Mục tiêu**: Kiểm tra câu lệnh SQL thực tế với GORM và Postgres (dùng Docker test container).
- **Kịch bản**:
    - [ ] Kiểm tra ràng buộc khóa ngoại (Foreign Keys).
    - [ ] Kiểm tra tính chính xác của cơ chế lọc `WHERE owner_id = ...`.

---

## 🎭 2. Frontend E2E Testing (Playwright)

Chúng ta sẽ chạy các bộ test mô phỏng hành vi người dùng thật trên trình duyệt.

### 2.1 Persona: System Admin
- **Flow**: Login -> Dashboard -> Create Manager User -> Navigate to Farm Registry -> Assign Farm to Manager -> Delete Farm.
- **Verify**: Data hiển thị đồng nhất, thông báo thành công (Toast) hiện ra.

### 2.2 Persona: Farm Manager
- **Flow**: Login -> Dashboard -> Create Self-Farm -> Verify only 1 farm visible.
- **Verify**: Không thể truy cập URL của trang Quản lý User (Chặn bởi RoleGuard).

### 2.3 Persona: Anonymous (Public)
- **Flow**: Truy cập `/` -> Chạy thử scenario OIDC Login -> Verify không thấy Sidebar quản trị.

---

## ⚙️ 3. Pipeline Thực thi (Automation)

Để chạy toàn bộ bộ test, sử dụng các lệnh sau:

```bash
# 1. Chạy Unit Test Backend
cd src/apps/farm-service && go test ./...

# 2. Chạy E2E Test Frontend (Yêu cầu infra đang chạy)
cd src/apps/client-app && npm run test:e2e
```

---

## ✅ 4. Danh sách Acceptance Criteria (Tự động)
- [ ] Code Coverage cho Farm Service > 80%.
- [ ] Playwright pass 100% các kịch bản Persona.
- [ ] Không có lỗi gRPC communication giữa Farm và Auth service.
