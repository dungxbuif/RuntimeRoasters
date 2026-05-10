# Core Framework Specification (`pkg/`)

This document serves as the technical reference for the shared foundational packages used across all microservices in the RuntimeRoasters ecosystem.

## `pkg/config` — Mở Rộng Pattern

`BaseConfig` nằm trong `pkg/`. Mỗi service tạo config riêng **embed** `BaseConfig`:

```go
// pkg/config/config.go
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

## `pkg/logger` — Thêm TraceID Injection

```go
// pkg/logger/logger.go
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

## `pkg/errs` — RFC 9457 Problem Details

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

## `pkg/database` — Postgres Wrapper

```go
// pkg/database/postgres.go

type PostgresConfig struct {
    URL          string
    MaxOpenConns int // default: 20 (HA Design)
    MaxIdleConns int // default: 5
}

type DB struct {
    *gorm.DB
}

func NewPostgres(cfg PostgresConfig) (*DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
    if err != nil { return nil, err }

    sqlDB, _ := db.DB.DB()
    if cfg.MaxOpenConns == 0 { cfg.MaxOpenConns = 20 }
    if cfg.MaxIdleConns == 0 { cfg.MaxIdleConns = 5 }

    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return &DB{db}, nil
}

// WithTx — dùng cho cơ chế Transaction của GORM
func (db *DB) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
    return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return fn(tx)
    })
}
```

---

## `pkg/redis` — Redis Client (Sentinel-Aware)

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

## `pkg/base` — Bootstrap (The Crown Jewel)

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
