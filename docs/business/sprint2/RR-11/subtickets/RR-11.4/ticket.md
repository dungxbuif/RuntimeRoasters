# RR-11.4: Refactor & Service Integration

## 🎯 Goal
Refactor `pkg/base/app` để hỗ trợ interceptor linh hoạt và tích hợp bảo mật vào `demo-service`.

## 📋 Tasks
- [ ] Refactor `pkg/base/app.go`:
    - Thêm `GRPCServerOptions []grpc.ServerOption` vào struct `Options`.
    - Cập nhật `NewApp` để áp dụng các options này khi khởi tạo `grpcServer`.
- [ ] Cập nhật `src/apps/demo-service/.env` với các biến: `INTERNAL_SECRET`, `JWKS_URL`, `EXPECTED_ISSUER`.
- [ ] Tích hợp middleware vào `apps/demo-service/internal/app/app.go`:
    - Mount `GinAuthMiddleware` cho các routes cần bảo vệ.
    - Inject `GRPCUnaryInterceptor` vào `base.NewApp`.

## 🔍 Definition of Done
- `demo-service` khởi động thành công và fetch được JWKS.
- gRPC server của demo-service áp dụng đúng interceptor.
