# 🏗️ Clean Architecture Framework — RuntimeRoasters

Tài liệu này định nghĩa tiêu chuẩn codebase và các mô hình thiết kế áp dụng cho mọi microservice trong dự án.

---

## 1. Triết lý Kiến trúc (The Dependency Rule)

```
domain/         ← Pure entities + domain errors ONLY. Zero imports from this project.
    ↑
usecase/        ← Declares its OWN repository interfaces. Imports domain/ only.
    ↑
infrastructure/ ← Implements usecase interfaces. Imports usecase/ + domain/ + pkg/*

main.go         ← Wires everything (Composition Root). Only place that knows all concrete types.
```

---

## 2. Tiêu chuẩn Thư viện Dùng chung (`pkg/`)

| Package | Vai trò | Tầng bảo vệ / Cơ chế |
| :--- | :--- | :--- |
| `pkg/base/auth` | Authentication | **Three-Gate Model**, gRPC Metadata Propagator. |
| `pkg/base/casbin` | Authorization | **Distributed Enforcement**, gRPC Snapshot + Kafka Sync. |
| `pkg/database` | Persistence | **GORM Wrapper**, PgBouncer-ready (Port 6432). |
| `pkg/kafka` | Messaging | **CloudEvents Standard**, Transactional Outbox Support. |
| `pkg/errs` | Error Handling | **RFC 9457**, Global Problem Details. |

---

## 3. Các Mô hình Thiết kế Quan trọng

### 3.1. Dual Idempotency (Lũy đẳng kép)
- **Tầng 1 (API)**: `Idempotency-Key` lưu vào Redis. Cache toàn bộ response (Body + Status Code).
- **Tầng 2 (Consumer)**: `Transactional Inbox` lưu `message_id` vào Postgres trong cùng transaction nghiệp vụ.

### 3.2. Transactional Outbox (Relay Worker)
- Mọi service phát hành sự kiện phải ghi vào bảng `outbox_events` trước khi worker ngầm đẩy lên Kafka. 
- Đảm bảo tính nguyên tử (Atomicity) giữa thay đổi trạng thái DB và phát sự kiện.

### 3.3. Defensive Programming (Chốt chặn số âm)
*Tham chiếu: UrbanX*
- Tuyệt đối không thực hiện phép tính thay đổi số lượng nhạy cảm mà không bọc trong `Math.Max(0, ...)`.

---

## 4. Tiêu chuẩn Bảo mật 3 Tầng (Three-Gate Auth)

1.  **Gate 1 (Gateway)**: Verify chữ ký và `scope`.
2.  **Gate 2 (Service)**: Casbin kiểm tra `role` (RBAC) và Repository kiểm tra `owner_id` (ABAC).
3.  **Gate 3 (Internal)**: Tự động trích xuất và chuyển tiếp JWT qua gRPC Metadata.

---
*Cập nhật lần cuối: 2026-05-10 bởi TechLead*
