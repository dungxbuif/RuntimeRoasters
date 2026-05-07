# RR-11.3: Auth Middleware (HTTP & gRPC)

## 🎯 Goal
Phát triển các "cổng bảo mật" (Middleware/Interceptor) cho cả hai giao thức HTTP và gRPC.

## 📋 Tasks
- [x] Triển khai `GinAuthMiddleware` trong `pkg/base/auth/middleware.go`:
    - Trích xuất Bearer token.
    - Parse & Verify RS256 signature bằng `jwksCache`.
    - Validate `exp`, `iss`, `aud`.
    - Inject `Claims` vào request context.
- [x] Triển khai `GRPCUnaryInterceptor` trong `pkg/base/auth/interceptor.go`:
    - Trích xuất token từ gRPC Metadata.
    - Logic verify tương tự HTTP.
- [x] Ghi `user.id` vào OTel span attributes cho cả hai loại middleware.

## 🔍 Definition of Done
- Middleware trả về đúng `401 Unauthorized` (RFC 9457) nếu token lỗi.
- Identity được truyền vào context thành công.
