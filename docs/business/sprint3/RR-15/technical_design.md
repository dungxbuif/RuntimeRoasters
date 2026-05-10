# Technical Design: Farm Service Bootstrapping (RR-15)

## 1. Overview
Farm Service là service nghiệp vụ đầu tiên của hệ thống sau khi đã hoàn thiện hạ tầng bảo mật. Service này sẽ quản lý thông tin các nông hộ và vùng trồng cà phê.

## 2. Component Design

### 2.1 Project Structure
Tiếp tục sử dụng Clean Architecture và copy boilerplate từ `demo-service` để đảm bảo tính nhất quán:
- `cmd/main.go`: Entrypoint.
- `internal/domain`: Entity và Repository interfaces.
- `internal/usecase`: Business logic.
- `internal/delivery`: gRPC Handlers.
- `internal/infrastructure`: Repo implementation, DB migrations.

### 2.2 Security Model
- **AuthN:** Sử dụng `pkg/base/auth` (JWKS Validation).
- **AuthZ:** Sử dụng `pkg/base/casbin` (Resilient Reader).
- **Policies:** Sẽ được khai báo trong `auth-service` và sync về Farm Service.

### 2.3 Data Model
- Table `farms`:
    - `id`: UUID (Primary Key)
    - `name`: String
    - `location`: String (Geo-coordinates or address)
    - `owner_id`: UUID (Reference to Identity Server `sub`)
    - `created_at/updated_at`: Timestamps

## 3. Integration
- **KrakenD:** Mở port :8083 (gRPC 50053) nội bộ cho Farm Service.
- **Database:** Sử dụng `farm_db` trong cụm Postgres chung.
- **OTel:** Export trace về SigNoz qua OTLP.

## 4. Verification Plan
- Unit tests cho Usecase và Repository.
- Integration test cho gRPC handler với mock auth.
