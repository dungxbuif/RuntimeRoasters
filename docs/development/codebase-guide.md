# 🛠️ Codebase & Development Guide

Tài liệu này hướng dẫn cách triển khai code trong hệ thống **MạchNguyên**, đảm bảo tính đồng nhất giữa 10+ Go services.

---

## 1. Clean Architecture Standard
Mỗi service tuân thủ cấu trúc 4 lớp. Dữ liệu chỉ được phép đi từ ngoài vào trong.

### 🔵 Lớp Domain (`internal/domain`)
*   **Chứa:** Business Entities (Struct), Interface của Repository và UseCase.
*   **Quy tắc:** Tuyệt đối không import bất kỳ thư viện bên ngoài nào (ngoại trừ chuẩn Go). Đây là lớp "thuần khiết" nhất.
*   **Ví dụ:** `type Farm struct { ... }`, `type FarmRepository interface { ... }`.

### 🟢 Lớp UseCase (`internal/usecase`)
*   **Chứa:** Logic nghiệp vụ chính (Business Logic).
*   **Quy tắc:** Chỉ phụ thuộc vào lớp Domain. Điều phối dữ liệu giữa các Repository.
*   **Ví dụ:** Hàm `CreateFarm` sẽ kiểm tra tên farm có trùng không trước khi gọi Repository để lưu.

### 🟡 Lớp Repository (`internal/repository`)
*   **Chứa:** Implementation của DB (SQL queries, Redis commands).
*   **Quy tắc:** Nhận vào Domain Entity và trả về Domain Entity. Mọi logic về Database (GORM, sqlx) nằm ở đây.

### 🔴 Lớp Delivery (`internal/delivery`)
*   **Chứa:** HTTP Handlers (Gin) hoặc gRPC Servers.
*   **Quy tắc:** Chịu trách nhiệm Parse Request (JSON/Proto), gọi UseCase, và Format Response (RFC 7807).

---

## 2. Quy tắc "Không Hardcode" với `pkg/config`
Mọi service phải bắt đầu bằng việc load cấu hình:
```go
// Ví dụ trong main.go
var cfg FarmConfig
err := config.LoadConfig(".", "farm", &cfg)
```
Tất cả các hằng số về hạ tầng (Port, DB_URL, Kafka_Topic) phải nằm trong file `.env`.

---

## 3. Quản lý Lỗi (Error Handling)
Hệ thống sử dụng chuẩn **RFC 7807 (Problem Details)**.
*   Tất cả các lỗi nghiệp vụ được định nghĩa tại `internal/domain/errors.go`.
*   Middleware tại `pkg/errs` sẽ tự động bắt các lỗi này và format thành JSON chuẩn cho UI.

---

## 4. Quy trình thêm một Service mới
1.  Tạo thư mục trong `apps/`.
2.  Định nghĩa Protobuf tại `api/proto/`.
3.  Implement lớp Domain (Entity & Interface).
4.  Implement lớp Repository (DB interaction).
5.  Implement lớp UseCase (Business logic).
6.  Dựng Delivery (HTTP/gRPC server).
7.  Đăng ký routing tại API Gateway.
