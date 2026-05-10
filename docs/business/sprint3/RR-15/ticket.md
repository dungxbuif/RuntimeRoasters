# [RR-15] Farm Service Bootstrapping

- **Summary:** Khởi tạo khung sườn (skeleton) cho Farm Service dựa trên kiến trúc chuẩn của dự án.
- **Priority:** `HIGH`
- **Type:** Feature

---

## 🔍 Acceptance Criteria

### Scenario 1: Khởi tạo Project Structure
- **Given:** Tôi có mã nguồn của `demo-service`.
- **When:** Tôi tạo thư mục `src/apps/farm-service` và copy cấu trúc từ `demo-service`.
- **Then:** Thư mục mới phải có đầy đủ `cmd/`, `internal/`, `config/` và `Makefile`.

### Scenario 2: Cấu hình Môi trường & Connectivity
- **Given:** Farm Service đã được khởi tạo.
- **When:** Tôi cấu hình `.env` và chạy service qua Docker Compose.
- **Then:** Service phải kết nối thành công tới Postgres (database `farm_db`) và Jaeger.

### Scenario 3: Gateway Integration
- **Given:** Farm Service đang chạy tại port nội bộ.
- **When:** Tôi cấu hình KrakenD để forward request tới Farm Service.
- **Then:** Request qua Gateway (`/api/v1/farms/...`) phải nhận được phản hồi từ Farm Service.

---

## 👨‍💻 Developer Implementation Guide

### 1. Folder Structure & Roles
- `src/apps/farm-service/cmd/main.go`: Entrypoint duy nhất. Gọi `InitializeApp`.
- `src/apps/farm-service/internal/app/`: Chứa logic nối dây DI (Wire).
- `src/apps/farm-service/config/`: Định nghĩa `Config` struct (nhúng `BaseConfig`).
- `api/runtime/farm/v1/farm.proto`: Hợp đồng API gRPC/REST.

### 2. Required API Contract (`farm.proto`)
- **Package:** `runtime.farm.v1`
- **Messages:** `Farm`, `CreateFarmRequest/Response`, `GetFarmRequest/Response`, `ListFarmsRequest/Response`, `UpdateFarmRequest/Response`, `DeleteFarmRequest/Response`.
- **Technique:** `Farm` message chỉ bao gồm các thông tin cơ bản: name, location, area, type, owner_id.

### 3. Wiring Technique (DI)
- Sử dụng **Google Wire**.
- Phải inject: `provider.KeyProvider` (cho JWT), `casbin.Engine` (cho AuthZ), và `database.DB` (cho GORM).
- **Ràng buộc:** `GRPCServerOptions` PHẢI bao gồm `authgrpc.GRPCUnaryInterceptor` và `casbingrpc.GRPCUnaryInterceptor`.

### 4. Integration Commands
```bash
# Sau khi tạo .proto
cd api && buf generate

# Sau khi tạo wire.go
cd src/apps/farm-service/internal/app && wire
```
