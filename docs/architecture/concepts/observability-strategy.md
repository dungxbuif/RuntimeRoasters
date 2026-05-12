# Thiết kế Telemetry — Aspire-like Observability cho Go

Tài liệu này mô tả chi tiết kiến trúc Observability của RuntimeRoasters: thiết kế package, thư viện, cấu hình hạ tầng, và hướng dẫn implement cho cả **Backend (Go)** và **Frontend (React/Next.js)**.

Mục tiêu: mang lại trải nghiệm "zero-config" giống **.NET Aspire** — developer chỉ cần gọi `base.NewApp()` và tự động có đầy đủ Tracing, Metrics, Log Correlation.

---

## 1. Kiến trúc Tổng quan

```
┌─────────────────────────────────────────────────────────────────┐
│                        RUNTIME ROASTERS                          │
│                                                                   │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │  client-app  │    │ demo-service │    │other-service │       │
│  │  (Next.js)   │    │  (Go + Gin)  │    │  (Go + gRPC) │       │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘       │
│         │ OTLP/HTTP         │ OTLP/gRPC          │ OTLP/gRPC    │
│         └──────────────┬────┴────────────────────┘              │
│                        ▼                                          │
│              ┌──────────────────┐                                │
│              │     SigNoz       │  ← traces + metrics + logs     │
│              │  :4317 (gRPC)    │    (embedded collector)        │
│              │  :4318 (HTTP)    │                                │
│              │  :3301 (UI) ─────┼──── EXPOSED to host ─────────  │
│              └──────────────────┘                                │
│                        ▲                                          │
│              ┌──────────────────┐                                │
│              │   ClickHouse     │  ← SigNoz storage backend      │
│              └──────────────────┘                                │
│                                                                   │
│  ┌─────────────────────────────────────────────┐                 │
│  │            control-app (Next.js)             │                │
│  │   /control/services                          │                │
│  │   /control/api-explorer                      │                │
│  └─────────────────────────────────────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
```

### SigNoz thay thế 5 containers

| Bỏ | Thêm | Lý do |
| :--- | :--- | :--- |
| OTel Collector | SigNoz (embedded collector) | SigNoz tích hợp sẵn OTLP receiver |
| Jaeger | SigNoz (traces UI) | Waterfall chart built-in |
| Prometheus | SigNoz (metrics UI) | ClickHouse-backed metrics |
| Loki | SigNoz (logs UI) | Structured log ingestion |
| Grafana | SigNoz (unified UI) | Single observability platform |

### Port Map (Observability)

| Port | Vai trò | Expose ra host? |
| :--- | :--- | :--- |
| `4317` | OTLP gRPC — Go services push traces/metrics/logs | **CÓ** |
| `4318` | OTLP HTTP — browser/frontend push traces | **CÓ** |
| `3301` | SigNoz UI | **CÓ** — Truy cập trực tiếp |

---

## 2. SigNoz UI — Truy cập trực tiếp

SigNoz UI (port 3301) được expose trực tiếp ra host để admin truy cập. Toàn bộ telemetry data (traces, metrics, logs) được SigNoz thu thập qua OTLP endpoints (4317, 4318).

---

## 3. Mục tiêu Thiết kế (The Four Pillars)

1. **Auto-Instrumentation:** HTTP (Gin), gRPC, SQL (Postgres), Valkey — tự động tạo spans mà không cần code thủ công trong business logic.
2. **Context Propagation:** `Trace-ID` lan truyền liên tục qua HTTP Headers, gRPC Metadata, **và Kafka Headers** — không bị đứt gãy tại bất kỳ biên giới nào.
3. **Log Correlation:** `logger.FromContext(ctx)` tự động đính `trace_id` + `span_id` từ OpenTelemetry vào mọi dòng log.
4. **Structured & Automatic Logging:** Mọi request (HTTP/gRPC) và lỗi hệ thống đều được tự động ghi log dưới dạng JSON với schema thống nhất.
5. **Centralized Hub:** Services gửi OTLP trực tiếp đến SigNoz — không cần OTel Collector riêng.

---

## 4. Backend (Go)

### 4.1 Thư viện cần cài

```bash
# Core SDK
go get go.opentelemetry.io/otel@latest
go get go.opentelemetry.io/otel/sdk@latest
go get go.opentelemetry.io/otel/sdk/metric@latest

# OTLP exporters (gửi về SigNoz qua gRPC)
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@latest
go get go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc@latest

# W3C Trace Context propagation
go get go.opentelemetry.io/otel/propagation@latest

# Auto-instrumentation: Gin
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin@latest

# Auto-instrumentation: gRPC
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc@latest

# Auto-instrumentation: HTTP client
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@latest

# Auto-instrumentation: SQL/Postgres
go get github.com/XSAM/otelsql@latest

# Auto-instrumentation: Valkey
go get github.com/valkey/go-valkey/extra/valkeyotel/v9@latest
```

---

### 4.2 Package mới: `pkg/telemetry`

#### `pkg/telemetry/provider.go`

```go
package telemetry

import (
    "context"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type Config struct {
    ServiceName    string
    ServiceVersion string
    OTLPEndpoint   string  // "signoz:4317"; empty = disabled
    SamplingRatio  float64 // 1.0 = AlwaysSample (dev), 0.1 (prod)
}

type ShutdownFunc func(ctx context.Context) error

// Init bootstraps Tracer + Meter providers and sets the global W3C propagator.
// Returns a shutdown function — call it during graceful stop to flush pending data.
func Init(cfg Config) (ShutdownFunc, error) {
    if cfg.OTLPEndpoint == "" {
        return func(ctx context.Context) error { return nil }, nil
    }

    ctx := context.Background()

    conn, err := grpc.NewClient(cfg.OTLPEndpoint,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, err
    }

    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceVersion(cfg.ServiceVersion),
        ),
        resource.WithHost(),
        resource.WithProcess(),
    )
    if err != nil {
        return nil, err
    }

    // -- Tracer Provider --
    traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
    if err != nil {
        return nil, err
    }

    sampler := sdktrace.TraceIDRatioBased(cfg.SamplingRatio)
    if cfg.SamplingRatio >= 1.0 {
        sampler = sdktrace.AlwaysSample()
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(traceExporter,
            sdktrace.WithBatchTimeout(5*time.Second),
            sdktrace.WithMaxExportBatchSize(512),
        ),
        sdktrace.WithSampler(sampler),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)

    // -- Meter Provider --
    metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
    if err != nil {
        return nil, err
    }

    mp := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(metricExporter,
            metric.WithInterval(15*time.Second),
        )),
        metric.WithResource(res),
    )
    otel.SetMeterProvider(mp)

    // -- Propagator: W3C TraceContext + Baggage --
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return func(ctx context.Context) error {
        _ = tp.Shutdown(ctx)
        _ = mp.Shutdown(ctx)
        return conn.Close()
    }, nil
}
```

#### `pkg/telemetry/kafka.go`

```go
package telemetry

import (
    "context"
    "go.opentelemetry.io/otel"
)

type kafkaCarrier map[string][]byte

func (c kafkaCarrier) Get(key string) string {
    if v, ok := c[key]; ok { return string(v) }
    return ""
}
func (c kafkaCarrier) Set(key, val string) { c[key] = []byte(val) }
func (c kafkaCarrier) Keys() []string {
    keys := make([]string, 0, len(c))
    for k := range c { keys = append(keys, k) }
    return keys
}

// InjectKafkaHeaders extracts the active span context and returns it as Kafka headers.
func InjectKafkaHeaders(ctx context.Context) map[string][]byte {
    carrier := kafkaCarrier{}
    otel.GetTextMapPropagator().Inject(ctx, carrier)
    return carrier
}

// ExtractKafkaHeaders restores the trace context from Kafka message headers.
func ExtractKafkaHeaders(ctx context.Context, headers map[string][]byte) context.Context {
    return otel.GetTextMapPropagator().Extract(ctx, kafkaCarrier(headers))
}
```

---

### 4.3 Cập nhật `pkg/base/app.go`

```go
type Options struct {
    Name    string
    Version string
    Config  config.BaseConfig
}

func NewApp(opts Options) *App {
    logger.InitLogger(opts.Config.AppEnv, opts.Config.LogLevel)
    log := logger.GetLogger()

    // 1. Init OTel trước tất cả — gửi về SigNoz
    shutdownOTel, err := telemetry.Init(telemetry.Config{
        ServiceName:    opts.Name,
        ServiceVersion: opts.Version,
        OTLPEndpoint:   opts.Config.OTLPEndpoint,  // "signoz:4317" in docker
        SamplingRatio:  opts.Config.TracingSampleRate,
    })
    if err != nil {
        log.Fatal("telemetry init failed", zap.Error(err))
    }

    engine := gin.New()
    engine.Use(gin.Recovery())
    engine.Use(otelgin.Middleware(opts.Name)) // OTel middleware trước GinErrorHandler
    engine.Use(errs.GinErrorHandler())

    grpcServer := grpc.NewServer(
        grpc.StatsHandler(otelgrpc.NewServerHandler()),
    )

    return &App{
        ginEngine:    engine,
        grpcServer:   grpcServer,
        logger:       log,
        shutdownOTel: shutdownOTel,
    }
}

func (a *App) Run(httpPort, grpcPort int) {
    // ...
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Flush OTel data TRƯỚC khi đóng server
    _ = a.shutdownOTel(ctx)
    a.httpServer.Shutdown(ctx)
    a.grpcServer.GracefulStop()
}
```

---

### 4.4 Cập nhật `pkg/logger/logger.go`

```go
import "go.opentelemetry.io/otel/trace"

// FromContext extracts the TraceID and SpanID from OpenTelemetry context and attaches it to the logger.
// Fallback to manual TraceIDKey if OpenTelemetry is not present.
func FromContext(ctx context.Context) *zap.Logger {
	logger := GetLogger()
	
	// OpenTelemetry integration
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		logger = logger.With(zap.String("trace_id", spanContext.TraceID().String()))
	}
	if spanContext.HasSpanID() {
		logger = logger.With(zap.String("span_id", spanContext.SpanID().String()))
	}

	// Fallback to manual context key
	if !spanContext.HasTraceID() {
		if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
			logger = logger.With(zap.String("trace_id", traceID))
		}
	}

	return logger
}
```

#### 4.5 Tự động ghi log Requests (Middlewares)

Hệ thống sử dụng các Middlewares để tự động ghi log mọi request kèm theo thông tin định danh và latency.

##### HTTP (Gin Logger Middleware)
Mọi request HTTP đều được ghi log với các trường: `status`, `method`, `path`, `ip`, `latency`, `trace_id`.

```go
// src/pkg/logger/gin_middleware.go
func GinLoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        // ... xử lý request ...
        log := FromContext(c.Request.Context())
        log.Info(path, 
            zap.Int("status", status),
            zap.Duration("latency", time.Since(start)),
            zap.String("trace_id", traceID),
        )
    }
}
```

##### gRPC (Logger Interceptor)
Mọi cuộc gọi gRPC được ghi log kèm theo `grpc.method`, `grpc.code` và `latency`.

#### 4.6 Ghi log trong Error Handler (RFC 9457)

`GinErrorHandler` không chỉ trả về lỗi cho client mà còn thực hiện ghi log tập trung.

```go
// src/pkg/errs/problem.go
func GinErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            log := logger.FromContext(c.Request.Context())
            
            // Tự động phân loại level log dựa trên HTTP Status
            if status >= 500 {
                log.Error("internal server error", zap.Error(err))
            } else {
                log.Warn("client request error", zap.Error(err))
            }
        }
    }
}
```

#### 4.7 Chuẩn hóa Định dạng Log (Log Schema)

Để đảm bảo khả năng truy vấn hiệu quả, toàn bộ log phải tuân thủ schema sau:

| Field | Type | Mô tả |
| :--- | :--- | :--- |
| `ts` | ISO8601 | Timestamp (Production) |
| `level` | String | INFO, WARN, ERROR, DEBUG |
| `msg` | String | Thông báo chính |
| `trace_id` | String | OTel TraceID để liên kết Traces |
| `span_id` | String | OTel SpanID |
| `latency` | Duration | Thời gian xử lý request |
| `error` | String | Chi tiết lỗi (nếu có) |

---

### 4.5 Cập nhật `pkg/database/postgres.go`

```go
import (
    "github.com/XSAM/otelsql"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewPostgres(cfg PostgresConfig) (*DB, error) {
    driverName, err := otelsql.Register("pgx",
        otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
        otelsql.WithSQLCommenter(true),
        otelsql.WithSpanOptions(otelsql.SpanOptions{RecordError: true}),
    )
    // ...
    if err := otelsql.RecordStats(db.DB); err != nil {
        return nil, err
    }
}
```

---

### 4.6 Cập nhật `pkg/valkey/client.go`

```go
import "github.com/valkey/go-valkey/extra/valkeyotel/v9"

func NewClient(cfg Config) *valkey.Client {
    // ...
    valkeyotel.InstrumentTracing(client)
    valkeyotel.InstrumentMetrics(client)
    return client
}
```

---

### 4.7 Config trong `.env`

```env
OTEL_EXPORTER_OTLP_ENDPOINT=signoz:4317   # empty = telemetry disabled
OTEL_TRACES_SAMPLE_RATE=1.0               # 1.0 dev, 0.1 prod
```

---

## 5. Frontend (React / Next.js)

### 5.1 Thư viện cần cài

```bash
npm install \
  @opentelemetry/api \
  @opentelemetry/sdk-trace-web \
  @opentelemetry/sdk-trace-base \
  @opentelemetry/exporter-trace-otlp-http \
  @opentelemetry/instrumentation \
  @opentelemetry/instrumentation-fetch \
  @opentelemetry/instrumentation-document-load \
  @opentelemetry/context-zone \
  @opentelemetry/resources \
  @opentelemetry/semantic-conventions
```

---

### 5.2 `src/lib/telemetry.ts`

```typescript
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { Resource } from '@opentelemetry/resources';
import { SEMRESATTRS_SERVICE_NAME } from '@opentelemetry/semantic-conventions';

export function initTelemetry() {
    if (typeof window === 'undefined') return;

    const provider = new WebTracerProvider({
        resource: new Resource({
            [SEMRESATTRS_SERVICE_NAME]: 'rr-web-frontend',
        }),
    });

    // Gửi traces về SigNoz qua OTLP HTTP port 4318
    provider.addSpanProcessor(
        new BatchSpanProcessor(
            new OTLPTraceExporter({
                url: process.env.NEXT_PUBLIC_OTEL_ENDPOINT ?? 'http://localhost:4318/v1/traces',
            }),
        )
    );

    provider.register({ contextManager: new ZoneContextManager() });

    // Auto-instrument fetch và page loads
    new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: [/localhost:808[0-9]/],
    }).enable();
    new DocumentLoadInstrumentation().enable();
}
```

---

## 6. Infrastructure — Docker Compose

SigNoz chạy cùng ClickHouse. Go services trỏ OTLP endpoint đến `signoz:4317`. Browser trỏ đến `http://localhost:4318`.

```yaml
clickhouse:
  image: clickhouse/clickhouse-server:24.1.2-alpine
  environment:
    CLICKHOUSE_DB: signoz_metrics
    CLICKHOUSE_USER: admin
    CLICKHOUSE_PASSWORD: <password>

signoz:
  image: signoz/signoz:latest
  environment:
    - CLICKHOUSE_HOST=clickhouse
    - CLICKHOUSE_PORT=9000
  ports:
    - "4317:4317"   # OTLP gRPC — exposed (Go services local)
    - "4318:4318"   # OTLP HTTP — exposed (browser)
    - "3301:3301"   # SigNoz UI — EXPOSED for direct access
  depends_on:
    - clickhouse
```

---

## 7. End-to-End Trace Flow

```
Browser (React)
  │
  ├─ DocumentLoad span: "page /farms loaded"
  │
  ├─ fetch GET /v1/demo/ping (qua KrakenD)
  │    traceparent: auto-injected by FetchInstrumentation
  │
KrakenD (port 8081)
  │
  └─ forward → demo-service HTTP :8080
       │
demo-service (Go + Gin)
  │
  ├─ otelgin.Middleware: extract traceparent → span "GET /v1/demo/ping"
  │    trace_id: cùng trace với browser!
  │
  ├─ logger.FromContext(ctx).Info("ping received")
  │    → log JSON tự động có: "trace_id": "...", "span_id": "..."
  │
  ├─ db.QueryContext(ctx, ...) → otelsql span "db.query"
  ├─ rdb.Set(ctx, ...) → valkeyotel span "valkey.set"
  └─ produce Kafka → telemetry.InjectKafkaHeaders(ctx)

  → Toàn bộ waterfall trace visible trong SigNoz UI tại localhost:3301
```

---

## 8. Thứ tự Implement (Sprint 1)

| Bước | Việc cần làm | Ticket |
| :--- | :--- | :--- |
| 1 | SigNoz + ClickHouse vào docker-compose (expose 3301) | RR-5 |
| 2 | `go get` các thư viện OTel | RR-6 |
| 3 | Implement `pkg/telemetry/provider.go` + `kafka.go` | RR-6 |
| 4 | Cập nhật `pkg/base/app.go` — tích hợp OTel + otelgin + otelgrpc | RR-6 |
| 5 | Cập nhật `pkg/logger`, `pkg/database`, `pkg/valkey` | RR-6 |
| 6 | Setup Frontend OTel (`initTelemetry`) | RR-7 |
| 7 | Dashboard Health Monitoring (không proxy SigNoz) | RR-7 |
