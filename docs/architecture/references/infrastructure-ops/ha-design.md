# ĐẶC TẢ THIẾT KẾ ĐỘ SẴN SÀNG CAO (HIGH AVAILABILITY - HA) RUNTIME ROASTERS

> **Triết lý: HA from Design — không phải HA from Patch**
>
> HA không phải là tính năng được thêm vào sau khi hệ thống đã chạy. Đây là một **tập hợp các ràng buộc thiết kế** được áp đặt ngay từ dòng code đầu tiên. Mọi service trong RuntimeRoasters được xây dựng theo nguyên tắc: **"Add a container = done."** — tức là scale ngang từ 1 lên N instances chỉ là thao tác hạ tầng, không yêu cầu thay đổi bất kỳ dòng logic nghiệp vụ nào.
>
> Tài liệu này là bộ **checklist kỹ thuật bắt buộc** cho mọi service, đảm bảo hệ thống không có Điểm lỗi đơn lẻ (Single Point of Failure - SPOF) ngay từ giai đoạn khởi tạo.

---

## 1. Nguyên Tắc Thiết Kế Lõi (Core HA Principles)

* **Stateless by Default:** Mọi Go Microservices (Farm, Logistics, Warehouse...) phải hoàn toàn phi trạng thái. Không lưu trữ session, file cục bộ hay cache trên RAM của container. Mọi trạng thái phải được đẩy xuống Data Layer (Postgres) hoặc Cache Layer (Redis).
* **Redundancy (Tính Dư Thừa):** Mọi thành phần (Service, Database, Broker) phải chạy tối thiểu từ 2 phiên bản (Replicas) trở lên trên các node/máy chủ vật lý khác nhau.
* **Fail-Fast & Graceful Degradation:** Khi một service nội bộ quá tải, các service phụ thuộc phải cắt kết nối sớm (Circuit Breaker) để tự bảo vệ, trả về dữ liệu cache hoặc báo lỗi nhanh thay vì chờ đợi và sập dây chuyền.

---

## 2. HA Tại Lớp Ứng Dụng (Go Microservices & Gateway)

Được áp dụng đồng loạt cho API Gateway, Webhook Service và các Go Services nội bộ:

* **Replication & Load Balancing:**
    * Triển khai tối thiểu 2-3 replicas cho mỗi service.
    * Sử dụng API Gateway (KrakenD/Traefik) làm Load Balancer (Round Robin hoặc Least Connections) để phân tải đều các request HTTP vào các replicas.
    * Đối với giao tiếp gRPC nội bộ, sử dụng Client-side Load Balancing tích hợp sẵn trong thư viện gRPC của Go kết hợp với Service Discovery (DNS nội bộ của Docker/Consul).
* **Health Checks (Liveness & Readiness Probes):**
    * *Liveness Probe:* Mỗi service phải expose endpoint `/health/live`. Nếu trả về lỗi, container orchestration (Docker Swarm/K8s) sẽ tự động kill và restart container đó.
    * *Readiness Probe:* Expose endpoint `/health/ready`. Service chỉ trả về HTTP 200 khi kết nối DB, Kafka và Redis của nó đã sẵn sàng. Nếu lỗi, Load Balancer sẽ ngừng gửi traffic đến replica này.
* **Graceful Shutdown (Tắt An Toàn):**
    * Khi nhận tín hiệu `SIGTERM` (scale down hoặc update), Go service phải ngừng nhận request mới, hoàn tất xử lý các request đang dở dang (ví dụ: chờ Kafka commit xong) trong một khoảng timeout nhất định (vd: 30s) trước khi tắt hẳn, tránh mất mát dữ liệu đang bay (in-flight data).
* **Triệt Tiêu Trạng Thái Cục Bộ (Statelessness Checklist):**
    * *Không lưu Session trong RAM:* Tuyệt đối không dùng biến `map` hay struct cục bộ của Go để lưu trạng thái user giữa các request. Xác thực hoàn toàn dựa vào **JWT** — API Gateway parse token và inject `X-User-ID`, `X-Role` vào Header. Service chỉ đọc Header, xử lý, và quên ngay lập tức. Request tiếp theo có thể rơi vào Instance khác và hoạt động hoàn toàn bình thường.
    * *Không cache trong biến global Go:* Mọi thao tác caching (query danh sách config hệ thống, rate-limit counter...) phải đi qua **Redis**. Nếu Instance A cập nhật config, Instance B nhận được ngay lập tức vì cùng đọc từ một nguồn. Biến `var cache = map[string]any{}` ở cấp package là anti-pattern trong môi trường multi-instance.
    * *Không lưu file lên ổ cứng container:* Nếu service có tính năng upload ảnh/tài liệu, không được `os.WriteFile()` vào filesystem của container. Container có thể bị kill và thay thế bất kỳ lúc nào — dữ liệu local sẽ mất theo. Phải **stream trực tiếp luồng byte lên Object Storage** (AWS S3 / MinIO) và chỉ lưu URL vào DB.
      ```go
      // ✅ HA-safe: stream thẳng lên MinIO, không chạm filesystem
      _, err = s.minioClient.PutObject(ctx, "rr-uploads", objectKey,
          fileReader, fileSize,
          minio.PutObjectOptions{ContentType: contentType},
      )
      ```

---

## 3. HA Tại Lớp Dữ Liệu (Data Layer)

* **PostgreSQL (Sát thủ giao dịch):**
    * Thiết lập cụm Primary-Replica (Active-Passive) sử dụng **Patroni** hoặc **Repmgr**.
    * Bật Streaming Replication đồng bộ (Synchronous Replication) để đảm bảo dữ liệu ghi vào Primary chắc chắn có mặt ở Replica trước khi báo thành công.
    * Sử dụng **PgBouncer** phía trước cụm DB để làm Connection Pooling, ngăn chặn tình trạng cạn kiệt kết nối (Connection Exhaustion) khi các Go services scale lên số lượng lớn.
    * **Giới hạn Connection Pool tại tầng Go:** Dù có PgBouncer, mỗi instance Go cũng phải tự giới hạn pool để không chiếm dụng quá nhiều slot khi scale ngang. Cấu hình bắt buộc trong `pkg/database`:
      ```go
      // pkg/database/postgres.go
      db.SetMaxOpenConns(20)            // Tối đa 20 kết nối mở đồng thời / instance
      db.SetMaxIdleConns(5)             // Giữ 5 kết nối idle để tái sử dụng
      db.SetConnMaxLifetime(time.Hour)  // Recycle kết nối sau 1 giờ, tránh stale conn
      ```
      Với 5 instances, tổng tối đa là 100 kết nối vào PgBouncer — nằm trong giới hạn an toàn của Postgres (`max_connections = 200`).
    * **Optimistic Locking (Khóa Lạc Quan) — Chống Ghi Đè Đồng Thời:** Khi 2 instances cùng đọc và muốn cập nhật cùng một row (ví dụ: cùng update trạng thái `HarvestBatch`), không dùng `SELECT FOR UPDATE` (pessimistic, gây lock contention). Thay vào đó, thêm cột `version` vào mỗi bảng nghiệp vụ quan trọng:
      ```sql
      -- Migration: thêm version vào bảng cần optimistic lock
      ALTER TABLE harvest_batches ADD COLUMN version INT NOT NULL DEFAULT 1;
      ```
      ```go
      // UseCase layer: luôn truyền version hiện tại vào câu UPDATE
      result, err := db.ExecContext(ctx, `
          UPDATE harvest_batches
          SET status = $1, version = version + 1
          WHERE id = $2 AND version = $3
      `, newStatus, batchID, currentVersion)

      rowsAffected, _ := result.RowsAffected()
      if rowsAffected == 0 {
          // Instance khác đã update trước — version đã đổi
          // Retry hoặc trả về lỗi Conflict (HTTP 409)
          return ErrOptimisticLockConflict
      }
      ```
      Instance nào đến chậm sẽ nhận `0 rows affected` và xử lý conflict — không bao giờ ghi đè dữ liệu của instance kia.
    * **PgBouncer — Connection Pool Tập Trung:** Mỗi Go instance tự giới hạn pool (`SetMaxOpenConns(20)`), nhưng khi scale ngang lên 10+ instances, tổng số kết nối trực tiếp vào Postgres có thể vượt `max_connections`. PgBouncer ngồi giữa Go services và Postgres, gom tất cả kết nối vào một pool nhỏ hơn và tái sử dụng chúng.

      ```
      Topology:
      ┌─────────────────────────────────────────┐
      │  Go Service Instances                   │
      │  [Farm x3] [Warehouse x3] [Retail x3]  │
      │   └─ mỗi instance: MaxOpenConns=20     │
      │   └─ tổng: 9 * 20 = 180 connections    │
      └───────────────┬─────────────────────────┘
                      │ tất cả đổ vào
                      ▼
              ┌──────────────┐
              │  PgBouncer   │  pool_size = 50 (transaction mode)
              │  :5432       │  Gom 180 req → 50 real connections
              └──────┬───────┘
                     │
                     ▼
              ┌──────────────┐
              │  PostgreSQL  │  max_connections = 100 (an toàn)
              │  Primary     │
              └──────────────┘
      ```

      Cấu hình `pgbouncer.ini` cốt lõi:
      ```ini
      [databases]
      farm_db     = host=postgres-primary port=5432 dbname=farm_db
      warehouse_db = host=postgres-primary port=5432 dbname=warehouse_db

      [pgbouncer]
      pool_mode         = transaction   ; Tốt nhất cho stateless Go services
      max_client_conn   = 500           ; Số kết nối từ Go services đến PgBouncer
      default_pool_size = 20            ; Số kết nối thật từ PgBouncer đến Postgres/db
      server_idle_timeout = 600
      ```

      `transaction` mode: PgBouncer chỉ giữ kết nối thật đến Postgres trong suốt 1 transaction. Giữa 2 request, kết nối được trả về pool — phù hợp hoàn hảo với stateless Go services dùng short-lived transactions.

* **Redis (Cache & Distributed Lock):**
    * Triển khai cấu hình **Redis Sentinel** (tối thiểu 3 nodes: 1 Master, 2 Slaves + 3 Sentinels).
    * Sentinel sẽ giám sát và tự động Promote (thăng cấp) Slave lên Master nếu Master hiện tại bị sập, thay đổi cấu hình định tuyến hoàn toàn trong suốt (transparent) với các Go services.
* **Apache Cassandra (Audit Log):** Triển khai cấu hình **Cluster** với nhiều nodes (Seed nodes) để đảm bảo không mất dữ liệu và có khả năng chịu lỗi cao (High Availability).
* **Elasticsearch (Traceability):** Triển khai cấu hình **Cluster** với cấu hình Sharding và Replicas (ví dụ: `number_of_replicas: 1`). Đảm bảo khi mất một Data Node, dữ liệu truy xuất không bị gián đoạn.

---

## 4. HA Tại Lớp Giao Tiếp Sự Kiện (Event Layer)

* **Apache Kafka / Redpanda Cluster:**
    * Triển khai Cluster tối thiểu 3 Brokers.
    * **Topic Replication Factor = 3:** Mọi topic sinh ra bắt buộc phải có hệ số nhân bản là 3.
    * **Min In-Sync Replicas (min.insync.replicas) = 2:** Yêu cầu ít nhất 2 brokers phải xác nhận đã nhận được Event thì mới tính là ghi thành công.
* **Cấu Hình Phía Producer (Go Services):**
    * Cài đặt `acks=all`: Producer chỉ xác nhận lưu Outbox thành công khi toàn bộ các In-Sync Replicas của Kafka đã lưu Event.
    * Bật cơ chế `enable.idempotence=true` trên Kafka Producer để chống mất hoặc lặp message do lỗi đường truyền.
* **Cấu Hình Phía Consumer:**
    * Sử dụng **Consumer Groups**: Nếu triển khai 3 replicas của `Warehouse Service`, cả 3 container này chung một `group_id`. Kafka sẽ tự động phân bổ các Partition cho 3 container này. Nếu 1 container sập, Kafka tự động Rebalance (phân bổ lại) tải cho 2 container còn lại, không gián đoạn việc tiêu thụ sự kiện.

---

## 5. HA Cho Tác Vụ Nền & Cronjob (Background Worker Safety)

Đây là điểm dễ bị bỏ sót nhất khi scale. Giả sử `Payment Service` có một background worker chạy mỗi phút để quét DB tìm giao dịch treo và gọi hoàn tiền. Với 5 instances, cả 5 workers đều thức dậy cùng lúc → cùng quét → cùng gọi Stripe Refund 5 lần → thảm họa.

* **Distributed Lock bằng Redis (Redlock Pattern):**
    * Trước khi thực thi task dùng chung, instance phải xin cấp "chìa khóa" từ Redis bằng lệnh nguyên tử `SET resource_name random_value NX PX ttl`.
    * `NX` (Not eXists) đảm bảo chỉ **một** instance nhận được `OK` trong cùng một thời điểm. 4 instances còn lại nhận `nil` và tự động skip.
    * TTL (`PX`) đảm bảo nếu instance giữ khóa bị chết đột ngột, khóa tự giải phóng sau timeout — tránh deadlock vĩnh viễn.

    ```go
    // pkg/lock/distributed_lock.go
    type DistributedLock struct {
        client    *redis.Client
        key       string
        ttl       time.Duration
        lockValue string // UUID riêng của lần acquire này
    }

    func NewLock(client *redis.Client, key string, ttl time.Duration) *DistributedLock {
        return &DistributedLock{client: client, key: key, ttl: ttl, lockValue: uuid.NewString()}
    }

    // Acquire: trả về true nếu giành được khóa, false nếu instance khác đang giữ
    func (l *DistributedLock) Acquire(ctx context.Context) (bool, error) {
        ok, err := l.client.SetNX(ctx, l.key, l.lockValue, l.ttl).Result()
        return ok, err
    }

    // Release: chỉ xóa nếu mình là người đang giữ (atomic via Lua script)
    func (l *DistributedLock) Release(ctx context.Context) error {
        script := `
            if redis.call("GET", KEYS[1]) == ARGV[1] then
                return redis.call("DEL", KEYS[1])
            end
            return 0`
        return l.client.Eval(ctx, script, []string{l.key}, l.lockValue).Err()
    }
    ```

* **Áp dụng vào Background Worker:**
    ```go
    // apps/payment-service/internal/worker/refund_worker.go
    func (w *RefundWorker) Run(ctx context.Context) {
        ticker := time.NewTicker(1 * time.Minute)
        for range ticker.C {
            lock := lock.NewLock(w.redis, "payment:refund-worker:lock", 55*time.Second)

            acquired, err := lock.Acquire(ctx)
            if err != nil || !acquired {
                // Có instance khác đang chạy — instance này nghỉ, không làm gì
                continue
            }

            // Chỉ 1 instance trong cluster đến được đây
            w.processStuckRefunds(ctx)
            lock.Release(ctx)
        }
    }
    ```
    TTL = 55s (nhỏ hơn interval 60s) đảm bảo khóa được giải phóng trước chu kỳ tiếp theo, kể cả khi `Release()` không được gọi do panic.

* **Kết quả:** Scale lên 100 instances, Refund Worker vẫn chỉ chạy đúng 1 lần mỗi phút — **Exactly-once execution** mà không cần bất kỳ điều phối bên ngoài nào.

---

## 6. Cơ Chế Bảo Vệ Và Tự Phục Hồi (Resiliency)

* **Circuit Breaker Pattern (Ngắt Mạch):**
    * Tích hợp thư viện (như `sony/gobreaker` hoặc `go-resiliency`) cho các lệnh gọi đồng bộ gRPC/HTTP giữa các service.
    * Nếu tỷ lệ lỗi gọi sang `Payment Service` vượt quá 50% trong 10 giây, ngắt mạch, chặn mọi request mới và lập tức trả về lỗi. Sau 30s, gửi các request thăm dò (Half-Open) để kiểm tra xem service đích đã phục hồi chưa.
* **Retry Mechanism & Exponential Backoff:**
    * Đối với các lỗi có thể phục hồi (Network Timeout, DB Deadlock), áp dụng Retry kèm theo độ trễ tăng dần (Exponential Backoff + Jitter) để tránh "dội bom" (Thundering Herd) khiến hệ thống đích càng thêm sập.
* **Dead Letter Queue (DLQ):**
    * Các event không thể xử lý sau N lần retry (lỗi schema, lỗi logic không hồi phục) phải được đẩy vào topic DLQ riêng biệt (ví dụ: `farm.harvest.created.dlq`) thay vì block consumer. Điều này đảm bảo 1 event lỗi không làm toàn bộ partition bị đình trệ.

---

## 7. Tiêu Chuẩn Hóa & Framework Nội Bộ (Standardization)

Để đảm bảo tính đồng nhất và khả năng phục hồi nhanh (MTTR), toàn bộ Microservices được xây dựng dựa trên một bộ Framework nội bộ đặt tại thư viện `pkg/`:

* **Service Bootstrap (`pkg/base`):** Mọi service mới đều khởi tạo thông qua hàm `base.Bootstrap()` để tự động kích hoạt:
    * **Health Checks:** Expose `/health/live` và `/health/ready` theo chuẩn ở Section 2.
    * **Structured Logging:** Tích hợp sẵn Zap logger với việc tự động inject `Trace-ID` từ context.
    * **Observability:** Khởi tạo OpenTelemetry provider (Tracing & Metrics) hướng về Jaeger/Prometheus.
    * **Graceful Shutdown:** Tự động lắng nghe tín hiệu OS để đóng kết nối DB, gRPC, HTTP một cách an toàn.
* **Chuẩn hóa Phản hồi Lỗi (RFC 9457):**
    * Sử dụng middleware tập trung để chuyển đổi các lỗi nghiệp vụ (Domain Errors) thành format **Problem Details**.
    * Đảm bảo hệ thống luôn trả về thông tin lỗi có cấu trúc, giúp việc debug và giám sát (Monitoring) trở nên chính xác hơn.
* **Chiến lược Testing:**
    * Áp dụng **Integration Testing** với **`testcontainers-go`**.
    * Mọi service có kết nối Database hoặc Message Broker bắt buộc phải có test suite chạy với infra thật trong container để kiểm thử các kịch bản HA (vd: rớt kết nối DB giữa chừng).

---

## 8. Triển Khai HA Với Docker Swarm (Current Phase Deployment)

Docker Swarm là lựa chọn cho giai đoạn hiện tại — đơn giản hơn Kubernetes nhưng đủ mạnh để thực hiện đầy đủ các nguyên tắc HA đã định nghĩa ở trên. Swarm biến `docker-compose.yaml` thành một kế hoạch triển khai phân tán thực sự.

### Tại sao Docker Swarm?

| Tiêu chí | Docker Compose | Docker Swarm | Kubernetes |
| :--- | :--- | :--- | :--- |
| Multi-node | Không | **Có** | Có |
| Auto-restart khi node chết | Không | **Có** | Có |
| Rolling update zero-downtime | Không | **Có** | Có |
| Health check routing | Không | **Có** | Có |
| Độ phức tạp vận hành | Thấp | **Thấp-Trung** | Cao |
| Phù hợp showcase | Dev only | **Production-grade demo** | Production thực |

### Kiến Trúc Swarm Cluster

```
┌─────────────────────────────────────────────────────┐
│                  Docker Swarm Cluster               │
│                                                     │
│  Manager Node (1)          Worker Nodes (2+)        │
│  ┌─────────────┐           ┌──────────┐             │
│  │ Swarm Mgr   │           │ Worker 1 │             │
│  │ (Raft)      │──────────►│          │             │
│  │             │           └──────────┘             │
│  │ Stack Deploy│           ┌──────────┐             │
│  │ Service Mgmt│──────────►│ Worker 2 │             │
│  └─────────────┘           └──────────┘             │
│                                                     │
│  Overlay Network: rr-network (encrypted)            │
│  Services communicate by name (DNS auto-resolve)    │
└─────────────────────────────────────────────────────┘
```

### Docker Swarm Stack File (`deployments/docker-stack.yaml`)

```yaml
version: "3.9"

networks:
  rr-network:
    driver: overlay
    encrypted: true          # mTLS tự động giữa các nodes

volumes:
  postgres-data:
    driver: local
  redis-data:
    driver: local

services:

  # ── API Gateway ─────────────────────────────────────
  gateway:
    image: rr/gateway:latest
    networks: [rr-network]
    ports:
      - "8080:8080"
    deploy:
      replicas: 2
      update_config:
        parallelism: 1       # Rolling update: update 1 replica tại một thời điểm
        delay: 10s
        failure_action: rollback
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      placement:
        max_replicas_per_node: 1  # Không đặt 2 replicas trên cùng 1 node
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/health/live"]
      interval: 10s
      timeout: 3s
      retries: 3
      start_period: 15s

  # ── Farm Service ─────────────────────────────────────
  farm-service:
    image: rr/farm-service:latest
    networks: [rr-network]
    environment:
      - DB_HOST=pgbouncer         # Trỏ vào PgBouncer, không phải Postgres trực tiếp
      - REDIS_ADDR=redis-master:6379
      - KAFKA_BROKERS=redpanda:9092
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
        failure_action: rollback
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      placement:
        max_replicas_per_node: 1
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8081/health/ready"]
      interval: 10s
      timeout: 3s
      retries: 3
      start_period: 20s

  # ── PgBouncer ────────────────────────────────────────
  pgbouncer:
    image: bitnami/pgbouncer:latest
    networks: [rr-network]
    environment:
      - POSTGRESQL_HOST=postgres
      - POSTGRESQL_PORT=5432
      - PGBOUNCER_POOL_MODE=transaction
      - PGBOUNCER_MAX_CLIENT_CONN=500
      - PGBOUNCER_DEFAULT_POOL_SIZE=20
    deploy:
      replicas: 1              # PgBouncer là stateless — có thể scale nếu cần
      restart_policy:
        condition: on-failure

  # ── PostgreSQL ───────────────────────────────────────
  postgres:
    image: postgres:16-alpine
    networks: [rr-network]
    volumes:
      - postgres-data:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/pg_password
    secrets: [pg_password]
    deploy:
      replicas: 1
      placement:
        constraints:
          - node.labels.role == db    # Ghim DB vào node chuyên dụng

  # ── Redis (Sentinel Mode) ────────────────────────────
  redis-master:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    networks: [rr-network]
    volumes:
      - redis-data:/data
    deploy:
      replicas: 1
      placement:
        constraints:
          - node.labels.role == cache

  redis-sentinel:
    image: redis:7-alpine
    command: >
      sh -c "redis-sentinel /etc/redis/sentinel.conf"
    networks: [rr-network]
    deploy:
      replicas: 3               # 3 Sentinels để đạt quorum
      placement:
        max_replicas_per_node: 1

secrets:
  pg_password:
    external: true              # Quản lý bằng `docker secret create`
```

### Các Lệnh Swarm Cốt Lõi

```bash
# Khởi tạo Swarm (chạy trên Manager node)
docker swarm init --advertise-addr <MANAGER_IP>

# Thêm Worker node vào cluster (chạy trên Worker)
docker swarm join --token <TOKEN> <MANAGER_IP>:2377

# Deploy/Update toàn bộ stack
docker stack deploy -c deployments/docker-stack.yaml rr

# Scale một service (không cần restart service khác)
docker service scale rr_farm-service=5

# Rolling update (zero-downtime)
docker service update --image rr/farm-service:v1.2.0 rr_farm-service

# Xem trạng thái các services
docker stack services rr

# Xem logs của service cụ thể
docker service logs -f rr_farm-service
```

### Mapping HA Design → Swarm Config

| HA Principle (Sections trên) | Swarm Implementation |
| :--- | :--- |
| Redundancy (Section 1) | `deploy.replicas: N` + `max_replicas_per_node: 1` |
| Health Checks (Section 2) | `healthcheck` block — Swarm tự route traffic qua node healthy |
| Graceful Shutdown (Section 2) | Swarm gửi `SIGTERM` → chờ `stop_grace_period` (default 10s) |
| Stateless (Section 2) | Không có `volumes` mount vào service container |
| DB Connection Pool (Section 3) | `pgbouncer` service đứng giữa — tất cả service trỏ vào đó |
| Redis Sentinel (Section 3) | `redis-sentinel` replicas=3 |
| Kafka Consumer Group (Section 4) | `group.id` cố định — Swarm scale replicas tự động |
| Zero Trust Network (mTLS) | `networks.encrypted: true` → overlay network mã hóa |
| Secrets Management | `docker secret` — không có password trong env plain text |