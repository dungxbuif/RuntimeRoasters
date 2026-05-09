# [RR-15] Technical Design: Farm Service Bootstrap & Domain

**Status:** `DRAFT`
**Author:** Tech Lead

---

## 1. Context & Goal
Khởi tạo service `farm-service` dựa trên template từ `demo-service`. Thiết lập Domain entities và cấu trúc thư mục Clean Architecture.

---

## 2. Technical Decisions

### 2.1 Service Bootstrap
- **Path:** `src/apps/farm-service`
- **Pattern:** Copy boilerplate từ `src/apps/demo-service` bao gồm:
    - `cmd/main.go`
    - `wire.go` / `wire_gen.go`
    - `internal/` structure (domain, usecase, repository, delivery).

### 2.2 Domain Entities
- **Package:** `internal/domain`
- **Structs:**
    - `Farm`: Định nghĩa các trường `ID`, `Name`, `Location`, `Area`, `Type`, `OwnerID`, `CreatedAt`, `UpdatedAt`.

### 2.3 Dependency Injection
- Sử dụng **Google Wire** để quản lý dependency.
- Các provider cần thiết: `Config`, `Logger`, `DB (Postgres)`, `FarmRepository`, `FarmUseCase`, `FarmDelivery`.

---

## 3. Implementation Plan
1. Tạo thư mục `src/apps/farm-service`.
2. Sao chép và refactor code từ `demo-service` (đổi tên package, config prefix).
3. Định nghĩa `Farm` struct trong `internal/domain/farm.go`.
4. Cập nhật `wire.go` để bao gồm các thành phần của Farm Service.
