# Technical Design: Token Revocation (RR-14)

## 1. Overview
Hệ thống sử dụng chiến lược **Distributed Blacklist** với Redis để vô hiệu hóa các JWT vẫn còn hạn sử dụng (Signature valid) nhưng đã bị người dùng logout hoặc admin thu hồi.

## 2. Component Design

### 2.1 Blacklist Interface (`pkg/base/auth/token`)
Định nghĩa một interface trừu tượng để các middleware có thể sử dụng mà không phụ thuộc trực tiếp vào Redis.

```go
type BlacklistChecker interface {
    IsRevoked(ctx context.Context, jti string) (bool, error)
}
```

### 2.2 Redis Implementation (`pkg/base/auth/blacklist`)
Triển khai interface sử dụng `go-redis`.

- **Key Format:** `blacklist:jti:{jti}`
- **TTL:** Tính toán dựa trên `exp` của token.
- **Logic:** Nếu key tồn tại trong Redis -> Token bị thu hồi.

### 2.3 Middleware Logic
Cả gRPC Interceptor và Gin Middleware sẽ thực hiện check sau khi verify signature thành công:
1. Verify Signature & Parse Token.
2. Lấy `jti` từ claims.
3. Kiểm tra `BlacklistChecker.IsRevoked(ctx, jti)`.
4. Nếu `true` -> Return `401 Unauthorized`.

## 3. Deployment
Redis đã có sẵn trong `docker-compose.dev.yaml`.

## 4. Security Considerations
- **Fail-Closed:** Nếu Redis sập, middleware sẽ trả về lỗi (deny access) để đảm bảo an toàn tối đa.
- **Performance:** Sử dụng `EXISTS` command trong Redis, độ trễ cực thấp (< 1ms).
