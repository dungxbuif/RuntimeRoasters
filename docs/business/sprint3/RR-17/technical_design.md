# [RR-17] Technical Design: Farm UseCase & API (gRPC/REST)

**Status:** `DRAFT`
**Author:** Tech Lead

---

## 1. Context & Goal
Triển khai logic nghiệp vụ tại tầng UseCase và phơi bày API qua gRPC Server. REST API sẽ được tự động hỗ trợ qua KrakenD Gateway mapping.

---

## 2. API Definition (Proto)
- **File:** `api/runtime/farm/v1/farm.proto`

```protobuf
service FarmService {
  rpc CreateFarm(CreateFarmRequest) returns (CreateFarmResponse);
  rpc GetFarm(GetFarmRequest) returns (GetFarmResponse);
  rpc ListFarms(ListFarmsRequest) returns (ListFarmsResponse);
}
```

---

## 3. Implementation Details

### 3.1 UseCase Layer
- **Package:** `internal/usecase`
- **Responsibilities:**
    - Input validation (sử dụng custom logic hoặc library).
    - Mapping từ DTO (gRPC requests) sang Domain Entity.
    - Gọi Repository để thực hiện lưu trữ/truy vấn.
    - Xử lý lỗi nghiệp vụ và chuyển đổi sang Domain Errors.

### 3.2 Delivery Layer (gRPC)
- **Package:** `internal/delivery/grpc`
- **Responsibilities:**
    - Triển khai interface được sinh ra từ Proto.
    - Tích hợp Middleware/Interceptor:
        - `AuthenticationInterceptor` (từ RR-11).
        - `AuthorizationInterceptor` (Casbin - từ RR-12).
    - Trích xuất `Claims` và truyền `OwnerID` xuống tầng UseCase.

---

## 4. Error Handling
- Sử dụng `pkg/errs` để trả về các mã lỗi RFC 9457 (Problem Details) thống nhất cho REST và mapping tương ứng sang gRPC Status Codes.
