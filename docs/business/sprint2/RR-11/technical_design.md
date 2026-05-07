# RR-11: Technical Design — Backend Security Core (JWT Validation)

**Branch:** `feat/RR-11`
**Status:** `READY FOR REVIEW`
**Author:** Senior Dev
**Date:** 2026-05-06

---

## 1. Context & Goal

Sau khi **RR-10** hoàn thành, Client App đã có thể nhận JWT (Access Token) từ Hydra và đính kèm vào mỗi API request qua header `Authorization: Bearer <token>`.

Nhiệm vụ của **RR-11** là xây dựng tầng bảo mật trên **phía Backend (Go services)** để:
- Xác thực chữ ký JWT **offline** bằng Public Key (không gọi lại Hydra mỗi request).
- Inject identity của caller vào `context.Context` để UseCase layer truy xuất type-safe.
- Inject identity của caller vào `context.Context` để UseCase layer truy xuất type-safe.

---

## 2. Codebase Analysis (Pre-work)

### 2.1. Điểm tích hợp đã có

Sau khi đọc kỹ codebase, tôi xác định các điểm quan trọng:

| Component | File | Liên quan đến RR-11 |
|---|---|---|
| `base.App` | `pkg/base/app.go` | **Nơi mount middleware.** `engine.Use(...)` phải được gọi tại đây, trước `RegisterHTTP`. `grpcServer` cũng cần thêm `UnaryInterceptor`. |
| `errs.GinErrorHandler()` | `pkg/errs/problem.go` | RFC 9457 Problem Details đã sẵn có. Middleware auth chỉ cần `c.AbortWithStatusJSON`. |
| `errs.ErrUnauthorized` | `pkg/errs/problem.go` | Sentinel error đã có — **tái sử dụng, không tạo error string mới.** |
| `otelgin.Middleware` | `pkg/base/app.go` | OTel đã setup. Middleware auth cần gọi `span.SetAttributes(attribute.String("user.id", ...))` trên span hiện tại. |
| **JWKS Endpoint** | `deployments/identity/nginx.conf` | **CRITICAL:** JWKS tại `http://identity:4434/.well-known/jwks.json`. Nginx yêu cầu header `X-Internal-Secret` để đi qua Gatekeeper. |

### 2.2. Thứ gì CHƯA có — cần tạo mới

- `pkg/base/auth/` — Toàn bộ package chưa tồn tại.
- `pkg/base/identity/` — Toàn bộ package chưa tồn tại.
- `github.com/golang-jwt/jwt/v5` — **Chưa có trong `go.mod`**, phải `go get` trước.

### 2.3. Risk Flags ⚠️

**[RISK-1] Thiếu `INTERNAL_SECRET` trong demo-service config**
`src/apps/demo-service/.env` không có biến này. Service sẽ fail silent hoặc panic khi fetch JWKS. **Phải bổ sung trước khi tích hợp middleware.**

**[RISK-2] `grpcServer` hardcode trong `NewApp()`**
`pkg/base/app.go:75-77` tạo `grpc.NewServer()` không có option để inject `UnaryInterceptor` từ bên ngoài. Cần refactor để nhận `[]grpc.ServerOption`.

**[RISK-3] JWT Issuer bị hardcode trong Hydra config**
Hydra phát JWT với `iss: http://localhost:4444/`. Middleware phải validate `iss` claim. Địa chỉ này thay đổi khi deploy → cần biến `EXPECTED_ISSUER` trong config service. Không được hardcode.

**[RISK-4] Kafka bị comment-out**
Ticket đề cập Kafka cho token revocation events, nhưng Kafka đang disabled trong compose. Với scope MVP, **chỉ dùng Redis Blacklist là đúng hướng** và không cần bật Kafka.

---

## 3. Solution Architecture

### 3.1. Package Layout

```
src/pkg/base/
├── auth/
│   ├── jwks.go         # JWKS fetcher + Stale-While-Revalidate cache
│   ├── middleware.go   # Gin HTTP middleware (GinAuthMiddleware)
│   └── interceptor.go  # gRPC UnaryServerInterceptor
└── identity/
    ├── identity.go     # Claims struct + type-safe context key
    └── context.go      # InjectContext(ctx, id), FromContext(ctx)
```

### 3.2. Data Flow

```
HTTP Request (Bearer Token)
        │
        ▼
[GinAuthMiddleware]
  ├── Extract token from Authorization header
  ├── Parse & validate signature via JWKS cache (RS256)
  ├── Validate: exp, iss, aud
  ├── Set OTel span attribute: user.id
  └── identity.InjectContext(ctx, claims)
        │
        │ (on fail → 401 + Problem Details JSON, c.Abort())
        ▼
[Gin Handler / gRPC-Gateway]
        │
        ▼
[UseCase Layer]
  └── id, ok := identity.FromContext(ctx)
```

### 3.3. JWKS Cache — Stale-While-Revalidate Design

```go
type jwksCache struct {
    mu         sync.RWMutex
    keys       []publicKey
    fetchedAt  time.Time
    ttl        time.Duration  // default: 5m
    refreshing atomic.Bool    // prevents thundering herd
    fetchFn    func() ([]publicKey, error)
}

// Read logic:
// 1. Return cached keys immediately (RLock).
// 2. If age > ttl AND not already refreshing:
//    → Set refreshing=true, launch background goroutine.
//    → Goroutine fetches new keys, replaces cache under WLock.
// 3. NO request is blocked during refresh.
```

### 3.4. Identity Context Design

```go
// pkg/base/identity/identity.go

type Claims struct {
    Subject string // sub = Kratos identity ID
    Role    string // custom claim từ Hydra
    OrgID   string // custom claim từ Hydra
    JTI     string // JWT ID — dùng cho blacklist check
}

// KHÔNG dùng string key (tránh collision giữa các package)
type contextKey struct{}

func InjectContext(ctx context.Context, c Claims) context.Context {
    return context.WithValue(ctx, contextKey{}, c)
}

func FromContext(ctx context.Context) (Claims, bool) {
    c, ok := ctx.Value(contextKey{}).(Claims)
    return c, ok
}
```

## 3.5. Security Consideration
Mặc dù RR-11 tập trung vào Offline Validation, nhưng việc bảo mật JWKS endpoint qua Nginx (`identity` service) là bắt buộc để đảm bảo chỉ các internal services mới có thể discovery key set.

---

## 4. Refactoring Required — `pkg/base/app.go`

### 4.1. Thêm `GRPCServerOptions` vào `Options`

```go
// BEFORE
type Options struct {
    Name   string
    Config config.BaseConfig
}

// AFTER
type Options struct {
    Name              string
    Config            config.BaseConfig
    GRPCServerOptions []grpc.ServerOption // NEW — injectable interceptors
}

// Trong NewApp():
// BEFORE
grpcServer: grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
),

// AFTER
serverOpts := append(
    []grpc.ServerOption{grpc.StatsHandler(otelgrpc.NewServerHandler())},
    opts.GRPCServerOptions...,
)
grpcServer: grpc.NewServer(serverOpts...),
```

### 4.2. HTTP Middleware mounting

Không cần thêm method vào `base.App`. Middleware sẽ được mount tại service level trong `RegisterHTTP` callback:
```go
a.Base.RegisterHTTP(func(e *gin.Engine) {
    // Mount auth middleware trước khi đăng ký routes
    e.Use(auth.GinMiddleware(jwksCache, rdb, cfg.ExpectedIssuer))
    // routes...
})
```

---

## 5. Environment Variables cần bổ sung

Cho `src/apps/demo-service/.env`:
```dotenv
INTERNAL_SECRET=dev-secret-change-in-production
JWKS_URL=http://identity:4434/.well-known/jwks.json
JWKS_CACHE_TTL=5m
EXPECTED_ISSUER=http://localhost:4444/
```

---

## 6. New Dependencies

```bash
cd src
go get github.com/golang-jwt/jwt/v5
```

Không cần `jwx` — `golang-jwt/jwt/v5` với custom `Keyfunc` là đủ và nhẹ hơn.

---

## 7. Sub-Task Breakdown (Granular)

Để tập trung scope và dễ dàng review, ticket RR-11 được chia nhỏ thành các subtasks sau:

1. [**RR-11.1: Identity Core & Dependencies**](./subtasks/RR-11.1.md)
   - `go get jwt/v5`, triển khai `pkg/base/identity` (Claims & Context).
2. [**RR-11.2: JWKS Provider with Cache**](./subtasks/RR-11.2.md)
   - Triển khai `pkg/base/auth/jwks.go` với cơ chế Stale-While-Revalidate.
3. [**RR-11.3: Auth Middleware (HTTP & gRPC)**](./subtasks/RR-11.3.md)
   - Xây dựng Gin Middleware và gRPC Interceptor core logic.
4. [**RR-11.4: Refactor & Service Integration**](./subtasks/RR-11.4.md)
   - Refactor `base.App` và tích hợp vào `demo-service`.
5. [**RR-11.5: Verification & E2E Testing**](./subtasks/RR-11.5.md)
   - Manual test, kiểm tra OTel attributes và SigNoz.

---

## 8. Open Questions (Cần confirm trước khi code)

1. **`X-Internal-Secret` management:** Hiện tại hardcode trong compose. Prod strategy?
2. **gRPC-to-gRPC call (Service A → Service B):** Service A có tự động forward token không? Hay dùng service account riêng? Cần xác nhận pattern trước khi viết gRPC interceptor.

---

## 9. Initialization & Deployment Strategy (Resilience)

Vì hệ thống sử dụng cơ chế **Fail-Fast** tại thời điểm khởi động (Service sẽ Panic nếu không fetch được JWKS), chúng ta cần đảm bảo hạ tầng Identity sẵn sàng trước khi Microservices chạy.

### 9.1. Phía Code (Application Level)
- **Exponential Backoff Retry:** Hàm `NewJWKSCache` thực hiện thử lại (retry) 5 lần với thời gian chờ tăng dần (2s, 4s, 8s, 16s, 32s). Điều này giúp service tự phục hồi nếu hạ tầng Identity chỉ bị chậm trễ nhẹ.

### 9.2. Phía Infrastructure (DevOps Level)
- **Healthchecks:** Service `identity` (Nginx) được cấu hình `healthcheck` dựa trên endpoint `/health/ready` của Kratos.
- **Dependency Orchestration:** 
    - **Docker Compose:** Sử dụng `depends_on` với `condition: service_healthy`.
    - **Kubernetes:** Khuyến nghị sử dụng **Init Containers** để thăm dò endpoint JWKS của Identity service trước khi khởi chạy container chính.

