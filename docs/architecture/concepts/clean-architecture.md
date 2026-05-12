# 🏗️ Clean Architecture Framework — RuntimeRoasters

Tài liệu này định nghĩa tiêu chuẩn codebase và cách tổ chức thư mục áp dụng cho mọi microservice trong dự án.

---

## 1. Triết lý Kiến trúc (The Dependency Rule)

Triết lý cốt lõi là **Dependencies point INWARDS**. Tầng bên trong không được biết về sự tồn tại của tầng bên ngoài.

```
domain/         ← Pure entities + domain errors ONLY. Zero imports from this project.
    ↑
usecase/        ← Declares its OWN repository interfaces. Imports domain/ only.
    ↑
infrastructure/ ← Implements usecase interfaces. Imports usecase/ + domain/ + pkg/*

main.go         ← Composition Root. Only place that knows all concrete types.
```

---

## 2. Cấu trúc Thư mục & Ánh xạ (Directory Mapping)

Dưới đây là cấu trúc chuẩn cho một microservice (ví dụ: `farm-service`) và sự tương ứng với các tầng lý thuyết của Clean Architecture:

### 2.1. Tầng Domain (`internal/domain/`)
- **Lý thuyết:** Lõi của ứng dụng (Entities, Value Objects). Chứa logic nghiệp vụ thuần túy không đổi.
- **Nội dung:** 
    - `farm.go`: Định nghĩa struct `Farm` và các hàm validate nghiệp vụ.
    - `errors.go`: Các lỗi nghiệp vụ đặc thù của domain.

### 2.2. Tầng UseCase (`internal/usecase/`)
- **Lý thuyết:** Điều phối luồng dữ liệu (Application Rules). Chứa các kịch bản sử dụng hệ thống.
- **Nội dung:** 
    - `farm_usecase.go`: Implementation của logic nghiệp vụ (Create, Update, List).
    - `repository.go`: **QUAN TRỌNG:** Tầng này định nghĩa (Interface) những gì nó cần từ hạ tầng.

### 2.3. Tầng Infrastructure (`internal/infrastructure/`)
- **Lý thuyết:** Các chi tiết thực thi (Frameworks & Drivers). Database, External Services, Message Broker.
- **Nội dung:** 
    - `repository/postgres/`: Implement interface repository đã định nghĩa ở UseCase bằng GORM/Postgres.
    - `external/`: Gọi các API của service khác.

### 2.4. Tầng Delivery / Transport (`internal/delivery/`)
- **Lý thuyết:** Cổng giao tiếp với thế giới bên ngoài (Interface Adapters).
- **Nội dung:** 
    - `grpc/`: gRPC handlers, chuyển đổi proto sang domain model.
    - `http/`: Gin handlers (nếu có).

### 2.5. Composition Root (`cmd/main.go` và `internal/app/`)
- **Lý thuyết:** Nơi khởi tạo và kết nối (Wiring) mọi thứ.
- **Nội dung:** Khởi tạo DB, Valkey, Repo, UseCase, Handler và "nối dây" chúng bằng Manual DI.

---

## 3. Quy chuẩn "Interface belong to Consumer"

Theo **ADR 0001**, chúng ta áp dụng quy tắc: **Interface phải nằm ở nơi nó được sử dụng (Consumer), không phải nơi nó được thực thi (Producer).**

- **Sai:** Định nghĩa `Repository` interface trong `infrastructure`.
- **Đúng:** Định nghĩa `Repository` interface trong `usecase`. Tầng `infrastructure` chỉ đơn giản là thực thi (implement) nó.

---
**Xem thêm:**
- [Tiêu chuẩn Hệ thống & Shared Libraries](./system-wide-standards.md) để biết chi tiết về `pkg/` và các mẫu thiết kế chung (Outbox, Idempotency).
- [ADR 0001: Clean Architecture](../adrs/0001-use-clean-architecture.md).
