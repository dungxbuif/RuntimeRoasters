# [RR-16] Flow 1: Create & List Farm

- **Summary:** Triển khai luồng tính năng tạo mới và hiển thị danh sách Nông trại (Farm).
- **Priority:** `HIGH`
- **Type:** Feature

---

## 🔍 Acceptance Criteria

### Scenario 1: Tạo mới Nông trại thành công
- **Given:** Tôi là một Farmer đã đăng nhập và đang ở trang "Thêm Nông Trại".
- **When:** Tôi nhập đầy đủ thông tin (Tên, Địa chỉ, Diện tích) và nhấn "Lưu".
- **Then:** Hệ thống lưu dữ liệu vào Postgres và hiển thị thông báo thành công.

### Scenario 2: Hiển thị danh sách Nông trại
- **Given:** Tôi đã có một số nông trại trong hệ thống.
- **When:** Tôi truy cập trang danh sách nông trại.
- **Then:** Tôi phải thấy danh sách các nông trại của mình với đầy đủ thông tin cơ bản.

---

## 👨‍💻 Developer Implementation Guide

### 1. Domain Model (`internal/domain/farm.go`)
- **Entity:** `Farm` struct với các trường: `ID` (UUID), `Name`, `Location`, `Area` (float64), `OwnerID` (string).
- **Technique:** `OwnerID` là khóa quan trọng để thực hiện Data Scoping (ABAC).

### 2. Repository (`internal/infrastructure/postgres/farm_repo.go`)
- **Functions:**
  - `Create(ctx, *Farm) error`: Chèn record vào Postgres.
  - `ListByOwner(ctx, ownerID string) ([]*Farm, error)`: TRUY VẤN bắt buộc kèm `WHERE owner_id = $1`.

### 3. UseCase (`internal/usecase/farm_usecase.go`)
- **Input/Output:** Nhận DTO (Request struct), trả về DTO (Response struct) hoặc Domain Entity.
- **`CreateFarm` Logic:**
  1. Lấy `CurrentUserID` từ `identity.FromContext(ctx)`.
  2. Validate: `Area > 0`.
  3. Gán `OwnerID = CurrentUserID`.
  4. Gọi Repo `Create`.

### 4. gRPC Handler (`internal/delivery/grpc/handler.go`)
- **Technique:** Gọi UseCase và map kết quả sang Protobuf messages. Phải log trace ID bằng `logger.FromContext(ctx)`.
