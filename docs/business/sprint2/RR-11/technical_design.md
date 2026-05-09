# Dev Notes - [RR-11] Backend Security Core - JWT Validation

## 🛠️ Technical Implementation Details
- **Library:** `github.com/golang-jwt/jwt/v5`.
- **Algorithm:** RS256 (Asymmetric Offline Validation — no round-trip to Identity Server per request).
- **Packages created:**
  - `pkg/base/auth` — JWKS fetcher, cache, middleware (HTTP + gRPC interceptor).
  - `pkg/base/identity` — Type-safe context key (private type), `FromContext(ctx)`, `InjectContext(ctx, id)`.
- **JWKS Security:** Pass `X-Internal-Secret` header when fetching JWKS.
- **JWKS Caching Strategy:** stale-while-revalidate with background goroutine refresh.
- **OTel:** Inject `user.id` into active span attributes in the middleware.

## 🏗️ Architectural Approach

### 1. Codebase Analysis & Integration Points
- `pkg/base/app.go`: Mount point for middleware and gRPC interceptors.
- `pkg/errs/problem.go`: Reused RFC 9457 Problem Details and `ErrUnauthorized` sentinel.
- `deployments/identity/nginx.conf`: Source of JWKS at `http://identity:4434/.well-known/jwks.json`.

### 2. Solution Architecture

#### Package Layout
```
src/pkg/base/
├── auth/
│   ├── jwks.go         # JWKS fetcher + Stale-While-Revalidate cache
│   ├── middleware.go   # Gin HTTP middleware
│   └── interceptor.go  # gRPC UnaryServerInterceptor
└── identity/
    ├── identity.go     # Claims struct + type-safe context key
    └── context.go      # InjectContext/FromContext helpers
```

#### JWKS Cache - Stale-While-Revalidate Design
1. Return cached keys immediately using RLock.
2. If age > TTL and not already refreshing:
   - Launch background goroutine to fetch new keys.
   - Update cache under WLock.
3. Requests are never blocked during refresh.

### 3. Required Refactoring (`pkg/base/app.go`)
- Added `GRPCServerOptions` to `Options` struct to allow injectable interceptors.
- Updated `NewApp()` to accept and apply these options.

## 🚀 Environment Configuration
Added to `src/apps/demo-service/.env`:
- `INTERNAL_SECRET`: Shared secret for JWKS fetch.
- `JWKS_URL`: `http://identity:4434/.well-known/jwks.json`.
- `JWKS_CACHE_TTL`: `5m`.
- `EXPECTED_ISSUER`: `http://localhost:4444/`.

## 📝 Sub-Ticket Breakdown
1. **RR-11.1:** Identity Core & Dependencies.
2. **RR-11.2:** JWKS Provider with Cache.
3. **RR-11.3:** Auth Middleware (HTTP & gRPC).
4. **RR-11.4:** Refactor & Service Integration.
5. **RR-11.5:** Verification & E2E Testing.

## ⚠️ Risks & Considerations
- **Risk 1:** Missing `INTERNAL_SECRET` in service config causes failure.
- **Risk 2:** Hardcoded gRPC server creation prevents interceptor injection (Refactored).
- **Risk 3:** JWT Issuer mismatches between Hydra and service config.
- **Risk 4:** Kafka dependency for revocation is deferred to MVP+.
