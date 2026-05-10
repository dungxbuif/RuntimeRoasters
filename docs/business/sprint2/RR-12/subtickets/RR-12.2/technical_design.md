# RR-12.2: Auth Service gRPC Snapshot API Design

## 1. Mục tiêu (Goal)
Thiết lập gRPC endpoint để phân phối bản snapshot của toàn bộ chính sách (Policies) cho các service khác trong hệ thống.

## 2. API Definition (Proto)
Vị trí: `api/runtime/auth/v1/auth.proto`

```proto
syntax = "proto3";
package runtime.auth.v1;

message PolicyRule {
  string p_type = 1; // "p" hoặc "g"
  string v0 = 2;     // sub / role
  string v1 = 3;     // obj / parent_role
  string v2 = 4;     // act
  string v3 = 5;
  string v4 = 6;
  string v5 = 7;
}

message GetPoliciesRequest {}

message GetPoliciesResponse {
  repeated PolicyRule rules = 1;
}

service AuthService {
  rpc GetPolicies(GetPoliciesRequest) returns (GetPoliciesResponse);
}
```

## 3. Implementation Logic

### 3.1 Data Access
- Auth Service sử dụng `enforcer.GetAdapter().(*gormadapter.Adapter)` để truy cập trực tiếp vào DB hoặc gọi `enforcer.GetPolicy()` và `enforcer.GetGroupingPolicy()`.
- **Ưu tiên:** Sử dụng `enforcer.GetModel().GetPolicy("p", "p")` và `enforcer.GetModel().GetPolicy("g", "g")` để lấy dữ liệu từ memory đã load.

### 3.2 Mapping Logic
Chuyển đổi từ định dạng Casbin sang gRPC Message:
```go
func (s *AuthServer) GetPolicies(ctx context.Context, req *pb.GetPoliciesRequest) (*pb.GetPoliciesResponse, error) {
    rules := [][]string{}
    // Lấy policy định nghĩa quyền (p)
    rules = append(rules, s.enforcer.GetPolicy()...)
    // Lấy policy định nghĩa group/role (g)
    rules = append(rules, s.enforcer.GetGroupingPolicy()...)
    
    resp := &pb.GetPoliciesResponse{}
    for _, r := range rules {
        resp.Rules = append(resp.Rules, &pb.PolicyRule{
            PType: r[0], // Lưu ý: Cần logic map PType phù hợp
            V0: r[0], V1: r[1], V2: r[2], ...
        })
    }
    return resp, nil
}
```
*Lưu ý: Casbin v3 trả về slice không bao gồm p_type ở phần tử đầu, cần cẩn thận khi mapping.*

## 4. Resilience Loop cho gRPC Server
- Đảm bảo gRPC server khởi động sau khi Database đã sẵn sàng.
- Sử dụng `pkg/logger` để ghi lại các request bootstrapping từ các service khác nhằm mục đích auditing.

## 5. Xác minh (Verification)
Sử dụng `grpcurl`:
```bash
grpcurl -plaintext localhost:50051 runtime.auth.v1.AuthService/GetPolicies
```
Kết quả mong đợi: Một danh sách JSON chứa toàn bộ các rules hiện có trong Postgres.
