# 🏗️ Clean Architecture Framework — RuntimeRoasters

Tài liệu này là bản onboarding kỹ thuật cho bất kỳ ai (hoặc AI) bắt đầu làm việc trong dự án. Nó trả lời hai câu hỏi cốt lõi: **"Hiện tại codebase trông như thế nào?"** và **"Chúng ta đang xây dựng gì và theo trật tự nào?"**

---

## 1. Hiện Trạng Codebase (Bootstrap Phase)

### Cấu trúc hiện tại của Repo

Chúng ta đã hoàn tất giai đoạn Bootstrap (Sprint 1), thiết lập nền móng vững chắc cho toàn bộ dự án dưới dạng Monorepo tại thư mục `src/`.

```text
/RuntimeRoasters
├── go.work                  ← Go Workspaces root
├── src/                     ← Source code root
│   ├── pkg/                 ← Shared Framework đã hoàn thiện (The Foundation)
│   └── apps/farm-service/   ← Service đầu tiên đã được dựng khung The Clean Architecture
│
└── docs/                    ← Tài liệu thiết kế kiến trúc và quy trình
```

> **Nguyên lý cốt lõi của Kiến trúc hiện tại:** Tất cả các service tuân thủ Clean Architecture chuẩn (v4 style), nơi Interfaces thuộc về **consumer (người dùng nó)**, chứ không phải implementor. Do đó, các interface của Repository được định nghĩa tại lớp `usecase/`, giúp `domain/` hoàn toàn phi phụ thuộc. Đây là phong cách viết Go chuẩn mực: *"Accept interfaces, return structs."*

### Trạng thái `pkg/` hiện tại

| Package | File | Trạng thái | Ghi chú |
| :--- | :--- | :--- | :--- |
| `pkg/config` | `config.go` | ✅ Done | Viper + BaseConfig |
| `pkg/logger` | `logger.go` | ✅ Done | Zap + OpenTelemetry Context |
| `pkg/telemetry` | `provider.go`, `kafka.go` | 🚧 Sprint 1 | OTel Tracer + Meter + Kafka propagation. Xem [telemetry.md](./telemetry.md) |
| `pkg/base` | `app.go` | ✅ Done | Bootstrap lifecycle manager |
| `pkg/errs` | `problem.go` | ✅ Done | RFC 9457 Problem Details |
| `pkg/database` | `postgres.go` | ✅ Done | pgx + sqlx + WithTx() |
| `pkg/redis` | `client.go` | ✅ Done | Redis Sentinel-aware client |
| `pkg/middleware` | — | ❌ Sprint 2 | Auth header, Idempotency |
| `pkg/outbox` | — | ❌ Sprint 2 | Outbox worker |
| `pkg/inbox` | — | ❌ Sprint 2 | Inbox dedup |
| `pkg/lock` | — | ❌ Sprint 2 | Distributed lock |

---

## 2. Cấu Trúc Monorepo Mục Tiêu

```
/RuntimeRoasters
│
├── api/                              # gRPC Proto (shared contract)
│   └── farm/v1/
│       ├── farm.proto
│       └── farm.pb.go               # auto-generated — không edit tay
│
├── pkg/                              # Internal shared framework (The Foundation)
│   ├── config/
│   │   └── config.go                # BaseConfig + LoadConfig()
│   ├── logger/
│   │   └── logger.go                # Zap + FromContext(ctx) → inject TraceID
│   ├── base/
│   │   └── bootstrap.go            # App lifecycle (HTTP + gRPC + Health + Shutdown)
│   ├── errs/
│   │   └── problem.go              # RFC 9457 error format + sentinel errors
│   ├── database/
│   │   └── postgres.go             # pgx wrapper + SetMaxOpenConns + WithTx()
│   ├── redis/
│   │   └── client.go               # Redis Sentinel-aware client + redisotel instrumentation
│   ├── telemetry/                   # Sprint 1 — Xem docs/architecture/telemetry.md
│   │   ├── provider.go             # OTel TracerProvider + MeterProvider + W3C Propagator
│   │   └── kafka.go                # InjectKafkaHeaders / ExtractKafkaHeaders
│   ├── middleware/
│   │   ├── auth.go                 # Extract X-User-ID / X-Role từ header
│   │   └── idempotency.go          # Idempotency-Key via Redis
│   ├── outbox/                      # Sprint 2
│   │   └── worker.go
│   ├── inbox/                       # Sprint 2
│   │   └── inbox.go
│   └── lock/                        # Sprint 2
│       └── distributed_lock.go
│
├── apps/                             # Các Microservices
│   ├── farm-service/                 # Service đầu tiên (Sprint 1)
│   │   ├── cmd/
│   │   │   └── main.go              # Composition Root — ~30 lines
│   │   ├── config/
│   │   │   └── config.go            # Embed BaseConfig + farm-specific fields
│   │   └── internal/
│   │       ├── domain/              # Layer 1: Pure Go Entities ONLY. No interfaces.
│   │       │   ├── farm.go          # Farm, HarvestBatch structs + domain methods
│   │       │   └── errors.go        # Domain sentinel errors
│   │       ├── usecase/             # Layer 2: Business Logic + Repository Interfaces
│   │       │   ├── farm_usecase.go  # Declares FarmRepository interface (v4!)
│   │       │   └── batch_usecase.go # Declares BatchRepository interface (v4!)
│   │       └── infrastructure/      # Layer 3: Concrete implementations
│   │           ├── postgres/        # Implements usecase Repository interfaces
│   │           │   ├── farm_repo.go
│   │           │   └── batch_repo.go
│   │           ├── grpc/            # gRPC server handlers
│   │           │   └── handler.go
│   │           ├── http/            # Gin REST handlers
│   │           │   └── router.go
│   │           └── worker/          # Outbox Worker (Sprint 2)
│   │               └── outbox_worker.go
│   │
│   ├── trace-service/               # Sprint 2
│   ├── warehouse-service/           # Sprint 3
│   └── gateway/                     # KrakenD config
│
├── deployments/
│   ├── docker-compose.yaml          # Dev environment
│   └── docker-stack.yaml           # Docker Swarm (production-grade demo)
│
└── docs/
    └── architecture/
```

---

## 3. Thiết Kế Chi Tiết: `pkg/` Framework

### 3.1 `pkg/config` — Mở Rộng Pattern

`BaseConfig` nằm trong `pkg/`. Mỗi service tạo config riêng **embed** `BaseConfig`:

```go
// pkg/config/config.go (giữ nguyên, chỉ extend)
type BaseConfig struct {
    AppName  string `mapstructure:"APP_NAME"`
    AppEnv   string `mapstructure:"APP_ENV"`   // "development" | "production"
    AppPort  int    `mapstructure:"APP_PORT"`
    LogLevel string `mapstructure:"LOG_LEVEL"` // "debug" | "info" | "warn"

    // Telemetry — chuẩn OTel env var names; empty = disabled
    OTLPEndpoint      string  `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
    TracingSampleRate float64 `mapstructure:"OTEL_TRACES_SAMPLE_RATE"` // 1.0 dev, 0.1 prod
}

func LoadConfig(path string, name string, out interface{}) error { ... }
```

```go
// apps/farm-service/config/config.go
type Config struct {
    pkg_config.BaseConfig                        // embed — inherit AppName, AppEnv, ...
    DatabaseURL string `mapstructure:"DATABASE_URL"`
    RedisAddr   string `mapstructure:"REDIS_ADDR"`
    GRPCPort    int    `mapstructure:"GRPC_PORT"`
}

func Load() Config {
    var cfg Config
    pkg_config.LoadConfig(".", ".env", &cfg)
    return cfg
}
```

**Quy tắc:** `pkg/config.BaseConfig` chỉ chứa những field mà **mọi** service đều cần. Field đặc thù của service nằm trong `apps/<service>/config/`.

---

### 3.2 `pkg/logger` — Thêm TraceID Injection

```go
// pkg/logger/logger.go — thêm function này
func FromContext(ctx context.Context) *zap.Logger {
    if Log == nil {
        InitLogger("development", "debug")
    }
    traceID, _ := ctx.Value(TraceIDKey).(string)
    if traceID == "" {
        return Log
    }
    return Log.With(zap.String("trace_id", traceID))
}

// Dùng trong service:
// logger.FromContext(ctx).Info("farm created", zap.String("id", farmID))
// Output: {"timestamp":"2026-04-20T...","trace_id":"abc-123","msg":"farm created","id":"f-001"}
```

---

### 3.3 `pkg/errs` — RFC 9457 Problem Details

Mục đích: **Mọi service trả về cùng một format lỗi JSON**. Frontend không cần xử lý nhiều format khác nhau.

```go
// pkg/errs/problem.go

// Format chuẩn RFC 9457 — Content-Type: application/problem+json
type Problem struct {
    Type     string       `json:"type"`             // "about:blank" by default
    Title    string       `json:"title"`
    Status   int          `json:"status"`
    Detail   string       `json:"detail,omitempty"`
    Instance string       `json:"instance,omitempty"` // request path (e.g. /farms/123)
    TraceID  string       `json:"trace_id,omitempty"`
    Errors   []SubProblem `json:"errors,omitempty"`   // RFC 9457 §3.1 extension
}

// SubProblem — một entry trong errors array, dùng cho validation
type SubProblem struct {
    Detail  string `json:"detail"`
    Pointer string `json:"pointer,omitempty"` // JSON Pointer RFC 6901 (e.g. /body/email)
}

// Sentinel errors — domain layer dùng những lỗi này
var (
    ErrNotFound     = errors.New("not_found")
    ErrConflict     = errors.New("conflict")
    ErrUnauthorized = errors.New("unauthorized")
    ErrValidation   = errors.New("validation_failed")
    ErrForbidden    = errors.New("forbidden")
    ErrInternal     = errors.New("internal_error")
)

// HTTP Status mapping
var httpStatusMap = map[error]int{
    ErrNotFound:     http.StatusNotFound,             // 404
    ErrConflict:     http.StatusConflict,             // 409
    ErrUnauthorized: http.StatusUnauthorized,         // 401
    ErrValidation:   http.StatusBadRequest,           // 400
    ErrForbidden:    http.StatusForbidden,            // 403
    ErrInternal:     http.StatusInternalServerError,  // 500
}

// GinErrorHandler middleware — trả về application/problem+json
func GinErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) == 0 { return }

        err := c.Errors.Last().Err
        status := StatusCode(err)

        c.Render(status, problemRenderer{Problem{
            Type:     "about:blank",
            Title:    http.StatusText(status),
            Status:   status,
            Detail:   err.Error(),
            Instance: c.Request.URL.Path,
            TraceID:  c.GetHeader("X-Trace-ID"),
        }})
    }
}
```

---

### 3.4 `pkg/database` — Postgres Wrapper

```go
// pkg/database/postgres.go

type PostgresConfig struct {
    URL          string
    MaxOpenConns int // default: 20 (HA Design)
    MaxIdleConns int // default: 5
}

type DB struct {
    *sqlx.DB
}

func NewPostgres(cfg PostgresConfig) (*DB, error) {
    db, err := sqlx.Connect("pgx", cfg.URL)
    if err != nil { return nil, err }

    if cfg.MaxOpenConns == 0 { cfg.MaxOpenConns = 20 }
    if cfg.MaxIdleConns == 0 { cfg.MaxIdleConns = 5 }

    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(time.Hour)

    return &DB{db}, nil
}

// WithTx — dùng cho Outbox Pattern (ghi business data + outbox trong 1 TX)
func (db *DB) WithTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil { return err }

    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}
```

---

### 3.5 `pkg/redis` — Redis Client (Sentinel-Aware)

```go
// pkg/redis/client.go

type RedisConfig struct {
    Addr          string   // single node (dev): "localhost:6379"
    SentinelAddrs []string // sentinel nodes (prod): ["s1:26379","s2:26379","s3:26379"]
    MasterName    string   // sentinel master name: "mymaster"
    Password      string
}

func NewClient(cfg RedisConfig) *redis.Client {
    if len(cfg.SentinelAddrs) > 0 {
        // Production: Sentinel mode (HA — tự động failover)
        return redis.NewFailoverClient(&redis.FailoverOptions{
            MasterName:    cfg.MasterName,
            SentinelAddrs: cfg.SentinelAddrs,
            Password:      cfg.Password,
        })
    }
    // Development: single node
    return redis.NewClient(&redis.Options{
        Addr:     cfg.Addr,
        Password: cfg.Password,
    })
}
```

---

### 3.6 `pkg/base` — Bootstrap (The Crown Jewel)

Đây là điểm then chốt của framework. Mọi service gọi `base.Bootstrap()` một lần và nhận được Health Checks + Graceful Shutdown + Logging tự động. Service author chỉ viết business logic.

```go
// pkg/base/bootstrap.go

type Options struct {
    Name     string
    Config   config.BaseConfig
    HTTPPort int
    GRPCPort int
}

type App struct {
    httpServer *http.Server
    grpcServer *grpc.Server
    logger     *zap.Logger
    ginEngine  *gin.Engine
}

func Bootstrap(opts Options) *App {
    // 1. Init logger
    logger.InitLogger(opts.Config.AppEnv, opts.Config.LogLevel)
    log := logger.GetLogger()
    log.Info("bootstrapping service", zap.String("name", opts.Name))

    // 2. Create Gin engine with default middleware
    gin.SetMode(map[string]string{"production": gin.ReleaseMode}[opts.Config.AppEnv])
    engine := gin.New()
    engine.Use(gin.Recovery())
    engine.Use(errs.GinErrorHandler())

    // 3. Health check endpoints (HA Design — Section 2)
    engine.GET("/health/live", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    return &App{
        ginEngine:  engine,
        grpcServer: grpc.NewServer(),
        logger:     log,
    }
}

func (a *App) RegisterHTTP(routes func(*gin.Engine)) {
    routes(a.ginEngine)
}

func (a *App) RegisterGRPC(desc *grpc.ServiceDesc, impl interface{}) {
    a.grpcServer.RegisterService(desc, impl)
}

// RegisterReadiness — service đăng ký hàm check dependency của mình
func (a *App) RegisterReadiness(check func() error) {
    a.ginEngine.GET("/health/ready", func(c *gin.Context) {
        if err := check(); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "degraded", "reason": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
}

func (a *App) Run(httpPort, grpcPort int) {
    // Start HTTP server
    go func() { a.httpServer.ListenAndServe() }()

    // Start gRPC server
    go func() {
        lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
        a.grpcServer.Serve(lis)
    }()

    // Graceful shutdown (HA Design — Section 2)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    a.logger.Info("shutting down gracefully...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    a.httpServer.Shutdown(ctx)
    a.grpcServer.GracefulStop()
    a.logger.Info("shutdown complete")
}
```

---

## 4. Clean Architecture: Farm Service Standard

### Quy Tắc Phụ Thuộc (The Dependency Rule)

```
domain/         ← Pure entities + domain errors ONLY. Zero imports from this project.
    ↑
usecase/        ← Declares its OWN repository interfaces. Imports domain/ only.
    ↑
infrastructure/ ← Implements usecase interfaces. Imports usecase/ + domain/ + pkg/*

main.go         ← Wires everything (Composition Root). Only place that knows all concrete types.
```

> **Nguyên tắc cốt lõi:** `domain/` KHÔNG khai báo Repository interfaces.
> Lớp `usecase/` sẽ tự định nghĩa các interface mà nó *cần* sử dụng — điều này giúp code dễ test và đúng chuẩn idiom của Go.
> Lớp `infrastructure/postgres/` sẽ thỏa mãn các interface đó một cách ngầm định (duck typing).

---

### Layer 1: `domain/` — Pure Go Entities Only

> **Quy tắc:** Domain chỉ chứa **các entities và định nghĩa lỗi nghiệp vụ (domain errors)**. Không có bất kỳ phụ thuộc (import) nào khác tại đây.

```go
// apps/farm-service/internal/domain/farm.go
package domain

import "time"

type Farm struct {
    ID        string
    Name      string
    Location  string
    OwnerID   string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type HarvestBatch struct {
    ID        string
    FarmID    string
    WeightKg  float64
    Variety   string
    Status    BatchStatus
    Version   int       // Optimistic lock
    CreatedAt time.Time
    UpdatedAt time.Time
}

type BatchStatus string

const (
    StatusPending    BatchStatus = "PENDING"
    StatusProcessing BatchStatus = "PROCESSING"
    StatusDone       BatchStatus = "DONE"
)

// Domain method — business rule lives IN the entity, not in UseCase
func (b *HarvestBatch) CanTransitionTo(next BatchStatus) error {
    if b.Status == StatusDone {
        return ErrBatchAlreadyCompleted
    }
    return nil
}
```

```go
// apps/farm-service/internal/domain/errors.go
package domain

import "github.com/dungxbuif/RuntimeRoasters/pkg/errs"

// Domain-specific sentinel errors — mapped to pkg/errs
var (
    ErrFarmNotFound        = errs.ErrNotFound
    ErrBatchNotFound       = errs.ErrNotFound
    ErrBatchConflict       = errs.ErrConflict
    ErrOptimisticLock      = errs.ErrConflict
    ErrBatchAlreadyCompleted = errs.ErrConflict
)
```

> **Lưu ý:** Không có file `repository.go` trong `domain/`. Các định nghĩa giao tiếp dữ liệu được chuyển lên `usecase/`.

---

### Layer 2: `usecase/` — Business Logic + Khai báo Repository Interfaces

> **Quy tắc:** Tầng `usecase` khai báo chính các repository interfaces mà nó cần để thực thi logic.
> Tầng Infrastructure sẽ implement chúng. Tầng Domain không cần biết sự tồn tại của chúng.

```go
// apps/farm-service/internal/usecase/farm_usecase.go
package usecase

import (
    "context"
    "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
    "github.com/google/uuid"
)

// Giao tiếp dữ liệu được định nghĩa ở đây (nơi tiêu thụ interfaces)
//
//go:generate mockery --name FarmRepository
type FarmRepository interface {
    Create(ctx context.Context, farm *domain.Farm) error
    GetByID(ctx context.Context, id string) (*domain.Farm, error)
    ListByOwner(ctx context.Context, ownerID string) ([]*domain.Farm, error)
}

//go:generate mockery --name BatchRepository
type BatchRepository interface {
    Create(ctx context.Context, batch *domain.HarvestBatch) error
    GetByID(ctx context.Context, id string) (*domain.HarvestBatch, error)
    UpdateStatus(ctx context.Context, id string, status domain.BatchStatus, version int) error
}

type FarmService struct {
    farmRepo  FarmRepository
    batchRepo BatchRepository
}

func NewFarmService(f FarmRepository, b BatchRepository) *FarmService {
    return &FarmService{farmRepo: f, batchRepo: b}
}

func (s *FarmService) CreateFarm(ctx context.Context, name, location, ownerID string) (*domain.Farm, error) {
    farm := &domain.Farm{
        ID:      uuid.NewString(),
        Name:    name,
        Location: location,
        OwnerID: ownerID,
    }
    if err := s.farmRepo.Create(ctx, farm); err != nil {
        return nil, err
    }
    return farm, nil
}

func (s *FarmService) CreateBatch(ctx context.Context, farmID string, weightKg float64, variety string) (*domain.HarvestBatch, error) {
    // 1. Validate farm exists
    if _, err := s.farmRepo.GetByID(ctx, farmID); err != nil {
        return nil, domain.ErrFarmNotFound
    }
    // 2. Build entity
    batch := &domain.HarvestBatch{
        ID:       uuid.NewString(),
        FarmID:   farmID,
        WeightKg: weightKg,
        Variety:  variety,
        Status:   domain.StatusPending,
        Version:  1,
    }
    if err := s.batchRepo.Create(ctx, batch); err != nil {
        return nil, err
    }
    return batch, nil
    // Sprint 2: write Outbox event here (same TX via WithTx)
}
```

---

### Layer 3: `infrastructure/` — Implements `usecase` Interfaces

> **Quy tắc:** Khai báo `postgres.FarmRepository` tự động thỏa mãn `usecase.FarmRepository` thông qua duck typing trong Go.
> Package infrastructure import `usecase` (để dùng type trong interface) và `domain` (để truyền tải entities).

```go
// apps/farm-service/internal/infrastructure/postgres/farm_repo.go
package postgres

import (
    "context"
    "database/sql"
    "github.com/dungxbuif/RuntimeRoasters/pkg/database"
    "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
)

// FarmRepository implements usecase.FarmRepository (implicitly — duck typing)
type FarmRepository struct {
    db *database.DB
}

func NewFarmRepository(db *database.DB) *FarmRepository {
    return &FarmRepository{db: db}
}

func (r *FarmRepository) Create(ctx context.Context, farm *domain.Farm) error {
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO farms (id, name, location, owner_id, created_at)
         VALUES ($1, $2, $3, $4, NOW())`,
        farm.ID, farm.Name, farm.Location, farm.OwnerID,
    )
    return err
}

func (r *FarmRepository) GetByID(ctx context.Context, id string) (*domain.Farm, error) {
    var farm domain.Farm
    err := r.db.GetContext(ctx, &farm,
        `SELECT id, name, location, owner_id, created_at FROM farms WHERE id = $1`, id,
    )
    if err == sql.ErrNoRows {
        return nil, domain.ErrFarmNotFound
    }
    return &farm, err
}
```

```go
// apps/farm-service/internal/infrastructure/grpc/handler.go
package grpc

import (
    "context"
    farmv1 "github.com/dungxbuif/RuntimeRoasters/api/farm/v1"
    "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
)

type FarmHandler struct {
    farmv1.UnimplementedFarmServiceServer
    svc *usecase.FarmService  // ← FarmService, not FarmUseCase
}

func NewFarmHandler(svc *usecase.FarmService) *FarmHandler {
    return &FarmHandler{svc: svc}
}

func (h *FarmHandler) CreateFarm(ctx context.Context, req *farmv1.CreateFarmRequest) (*farmv1.CreateFarmResponse, error) {
    farm, err := h.svc.CreateFarm(ctx, req.Name, req.Location, req.OwnerId)
    if err != nil { return nil, err }
    return &farmv1.CreateFarmResponse{Id: farm.ID}, nil
}
```

```go
// apps/farm-service/internal/infrastructure/http/router.go
package http

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
    "github.com/dungxbuif/RuntimeRoasters/pkg/errs"
)

func NewRouter(svc *usecase.FarmService) func(*gin.Engine) {
    return func(e *gin.Engine) {
        v1 := e.Group("/api/v1")
        v1.POST("/farms",             createFarmHandler(svc))
        v1.GET("/farms/:id",          getFarmHandler(svc))
        v1.POST("/farms/:id/batches", createBatchHandler(svc))
    }
}

func createFarmHandler(svc *usecase.FarmService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Name     string `json:"name"     binding:"required"`
            Location string `json:"location" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.Error(errs.ErrValidation)
            return
        }
        // ownerID được trích xuất từ JWT context sau khi validate in-memory
        ownerID := c.GetString("user_id") 
        farm, err := svc.CreateFarm(c.Request.Context(), req.Name, req.Location, ownerID)
        if err != nil {
            c.Error(err)
            return
        }
        c.JSON(http.StatusCreated, farm)
    }
}
```

---

### Composition Root: `cmd/main.go` & Dependency Injection

Dự án sử dụng **Google Wire** để thực hiện Dependency Injection (DI) tại thời điểm biên dịch. Điều này giúp loại bỏ việc khởi tạo thủ công rườm rà trong `main.go` và đảm bảo tính an toàn.

#### 1. File định nghĩa: `internal/app/wire.go`

```go
// +build wireinject

package app

import (
    "github.com/google/wire"
    // ... imports
)

func InitApp(cfg config.Config, db *database.DB, rdb *redis.Client) (*base.App, error) {
    wire.Build(
        postgres.NewFarmRepository,
        postgres.NewBatchRepository,
        usecase.NewFarmService,
        infhttp.NewRouter,
        infgrpc.NewHandler,
        // ...
    )
    return nil, nil
}
```

#### 2. Entrypoint: `cmd/main.go`

```go
func main() {
    // 1. Config & Logger
    cfg := svcconfig.Load()
    app := base.NewApp(...)

    // 2. Infrastructure connections
    db, _ := database.NewPostgres(...)
    rdb   := redis.NewClient(...)

    // 3. Dependency Injection via Wire
    // Chạy command 'wire ./internal/app' để sinh code
    application, err := app.InitApp(cfg, db, rdb)
    if err != nil { panic(err) }

    // 4. Run
    application.Run(cfg.AppPort, cfg.GRPCPort)
}
```

---

## 5. Trật Tự Build (Sprint 2 onward)

### `pkg/` Foundation — ✅ DONE (Sprint 1)

```
✅ pkg/errs       RFC 9457 + sentinel errors
✅ pkg/config     Viper + BaseConfig
✅ pkg/logger     Zap + FromContext()
✅ pkg/database   pgx + sqlx + WithTx()
✅ pkg/redis      Redis Sentinel-aware client
✅ pkg/base       Bootstrap + Health + Graceful Shutdown
```

### Farm Service — Next Steps

```
Step 1:  domain/              Entities + domain errors ONLY (không interface)
Step 2:  usecase/             Declare interfaces + business logic
Step 3:  infra/postgres/      Implement usecase interfaces with SQL
Step 4:  infra/http/          Gin handler → FarmService
Step 5:  infra/grpc/          gRPC handler → FarmService (after proto)
Step 6:  cmd/main.go          Wire everything — Composition Root
```

### gRPC Contract (parallel or after Step 4)

```
api/farm/v1/farm.proto    Schema-first contract
protoc generate           Auto-generate Go stubs (do NOT edit by hand)
```

---

## 6. Mental Model Cốt Lõi

```
pkg/         = The shared framework. Written once, inherited by all services.
               → Think of it as the "Rails" of this project.

domain/      = What the business IS. Pure Go structs + domain methods.
               → Can be whiteboarded and coded directly. Zero framework imports.

usecase/     = What the business DOES. Declares interfaces it needs, implements logic.
               → This is where the "business story" is written as code.

infra/       = HOW data is stored and served. Postgres, gRPC, HTTP, Kafka.
               → Swap Postgres for MongoDB? Only this layer changes.

main.go      = The glue. Reads config, wires all layers, runs the app.
               → Composition Root. Zero business logic allowed here.
```

---

## 7. Quy Tắc Không Vi Phạm

1. **`domain/` không được import gì ngoài stdlib và `pkg/errs`.** Không có logic về Postgres, Gin tại đây.
2. **`domain/` không chứa Repository interfaces.** (Interface thuộc về Consumer Layer - `usecase/`).
3. **`usecase/` không được biết Postgres tồn tại.** Nó chỉ gọi interface mà nó tự khai báo.
4. **Không hardcode** — mọi cấu hình phải nạp qua `config/`.
5. **Mọi state phải ra ngoài container** — không có global `map`, không cache thông qua memory array, không `os.WriteFile`.
6. **`main.go` là nơi duy nhất** (Composition Root) biết concrete implementation nào đang được khởi tạo.
