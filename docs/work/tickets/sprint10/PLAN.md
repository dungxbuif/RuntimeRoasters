# RR-URG-04/05 Stripe Webhook Demo Plan

## Summary

Implement Stripe demo payment pass/fail through the official webhook endpoint,
available to all authenticated roles. Do not add a separate simulate endpoint.
The client fetches a seeded demo signing key, builds a Stripe-compatible payload,
signs it, and posts to `POST /v1/webhooks/stripe`.

Decision: write the Stripe webhook auth code in-repo instead of installing the
Stripe SDK. This flow only needs raw-body HMAC verification and Event payload
parsing, which Stripe documents clearly, and the repo already has HMAC helpers.
Avoiding a new SDK keeps this demo path small and deterministic.

Implementation status: RR-URG-04/05 webhook demo and Sprint 6/7 socket/topology
vertical slice are implemented in code. Remaining hardening is live WebSocket
verification through KrakenD in the full local stack and any multi-replica socket
fanout design beyond the single-replica Sprint 6/7 target.

Stripe docs checked:

- https://docs.stripe.com/webhooks
- https://docs.stripe.com/webhooks/signature

## Key Changes

### Database And Seed

- Add `payment_webhook_keys` table in the payment-service migration:
  - `id UUID PRIMARY KEY`
  - `provider VARCHAR(40) NOT NULL`
  - `key_name VARCHAR(120) NOT NULL`
  - `signing_secret VARCHAR(255) NOT NULL`
  - `active BOOLEAN NOT NULL DEFAULT true`
  - `demo BOOLEAN NOT NULL DEFAULT true`
  - `created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`
  - `updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`
  - unique `(provider, key_name)`
- Seed one active demo Stripe key in the migration:
  - provider: `STRIPE`
  - key name: `demo-stripe-local`
  - signing secret: `whsec_rr_demo_stripe_local`
- Align migration columns with the current `Payment` model if needed:
  - `provider`
  - `provider_ref`
  - `store_id`
  - `items`
  - `refund_ref`
  - `checkout_url`
  - `simulated`

### Payment Service

- Add repository/usecase access for the active demo Stripe webhook key.
- Add authenticated endpoint:
  - `GET /v1/payments/demo/stripe-webhook-key`
- Return:

```json
{
  "provider": "STRIPE",
  "key_name": "demo-stripe-local",
  "signing_secret": "whsec_rr_demo_stripe_local",
  "webhook_url": "/v1/webhooks/stripe",
  "algorithm": "HMAC-SHA256"
}
```

- Add canonical webhook route:
  - `POST /v1/webhooks/stripe`
- For this demo phase, require both app authentication and valid
  `Stripe-Signature`.
- Parse Stripe Event payloads:
  - `payment_intent.succeeded` updates payment to `SUCCEEDED` and publishes
    `payment.completed`.
  - `payment_intent.payment_failed` updates payment to `FAILED` and publishes
    `payment.failed`.
- Preserve idempotency by Stripe event `id`.

### Stripe Signature Code

- Support `Stripe-Signature: t=<unix>,v1=<hex>`.
- Verify HMAC-SHA256 over `timestamp + "." + raw_body`.
- Enforce a 5-minute timestamp tolerance.
- Use constant-time signature comparison.
- Keep provider-specific generic `X-Timestamp` / `X-Signature` only if needed
  for existing VNPay/simple tests.

### Stripe-Compatible Client Payload

Success payload:

```json
{
  "id": "evt_rr_demo_success_<uuid>",
  "object": "event",
  "type": "payment_intent.succeeded",
  "livemode": false,
  "created": 1779730000,
  "data": {
    "object": {
      "id": "pi_sim_<provider_ref>",
      "object": "payment_intent",
      "amount": 10000,
      "currency": "usd",
      "status": "succeeded",
      "metadata": {
        "order_id": "<order_id>",
        "payment_id": "<payment_id>",
        "store_id": "<store_id>"
      }
    }
  }
}
```

Failure payload:

```json
{
  "id": "evt_rr_demo_failed_<uuid>",
  "object": "event",
  "type": "payment_intent.payment_failed",
  "livemode": false,
  "created": 1779730000,
  "data": {
    "object": {
      "id": "pi_sim_<provider_ref>",
      "object": "payment_intent",
      "amount": 10000,
      "currency": "usd",
      "status": "requires_payment_method",
      "last_payment_error": {
        "type": "card_error",
        "code": "card_declined",
        "message": "Demo card declined"
      },
      "metadata": {
        "order_id": "<order_id>",
        "payment_id": "<payment_id>",
        "store_id": "<store_id>"
      }
    }
  }
}
```

### Frontend

- Show Pass/Fail webhook buttons for all logged-in roles on the Finance payment
  list.
- Button flow:
  - fetch `GET /v1/payments/demo/stripe-webhook-key`
  - build exact JSON Stripe Event
  - sign the raw JSON string
  - post to `POST /v1/webhooks/stripe`
  - refresh payments after a `200` response
- Use `data.object.id = payment.provider_ref`.
- Include `order_id`, `payment_id`, and `store_id` in Stripe metadata.

### Gateway And AuthZ

- Add KrakenD JWT routes for:
  - `GET /v1/payments/demo/stripe-webhook-key`
  - `POST /v1/webhooks/stripe`
- Add narrow Casbin policies for all active demo roles to:
  - read `/v1/payments/demo/stripe-webhook-key`
  - write `/v1/webhooks/stripe`
- Do not broaden payment mutation permissions elsewhere.

## RR-URG-04/05 Alignment

### RR-URG-05 Paid Order Flow

- Success webhook is the only trigger for reservation in this demo branch.
- Failure webhook must publish `payment.failed`; warehouse must not reserve
  stock.
- Fix the payment lookup route mismatch:
  - client/KrakenD use `/v1/payments/orders/{order_id}`
  - payment-service must serve the same path.

### RR-URG-04 Mandatory Return

- Retail order completion must wait for mandatory return-to-base.
- Update retail consumer to handle `logistics.driver.returned_to_base` for
  terminal `COMPLETED`, or make logistics emit `logistics.delivery.completed`
  only after return-to-base.
- Do not mark the retail order completed at
  `logistics.delivery.driver_confirmed`.

## Test Plan

### Payment Service

- Seeded Stripe key loads from DB.
- Valid Stripe success signature updates payment and publishes
  `payment.completed`.
- Valid Stripe failure signature updates payment and publishes `payment.failed`.
- Bad signature is rejected.
- Old timestamp is rejected.
- Duplicate Stripe event ID is idempotent.

### End-To-End Flow

- Paid order + Fail button:
  - payment becomes `FAILED`
  - order becomes rejected/failed
  - no stock reservation is created
- Paid order + Pass button:
  - payment becomes `SUCCEEDED`
  - warehouse reserves stock
  - logistics delivery starts
- Delivery confirmation alone does not complete order.
- Return-to-base completes order.

### Commands

```bash
cd src
GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/payment-service/... ./apps/warehouse-service/... ./apps/retail-service/... ./apps/logistics-service/...
```

```bash
cd src/apps/client-app
npm run lint
```

## Assumptions

- "Public API for client get the key" means available to every authenticated app
  role, not no-auth internet public.
- This phase supports Stripe only; VNPay stays unchanged.
- The seeded webhook key is demo-only and intentionally exposed to the browser.

---

# Sprint 6/7 Socket + Topology Revision Plan

## Summary

Add a new `socket-service` for realtime demo fanout, but keep it DB-free. The
service uses Gorilla WebSocket, Kafka, environment-provided internal API keys,
Valkey/Redis for ephemeral socket session metadata, and trace-derived topology
data. It does not own durable history or topology storage.

History/query belongs to `trace-service`, because trace-service already owns
PostgreSQL plus Elasticsearch read models. OTel is used for trace correlation
and enrichment, not as the primary business event source. Kafka/domain events
remain the source for demo timeline and replay.

Topology config is not admin-editable now or later. The public topology is
derived from trace-service data: CloudEvents, trace read models, topic mappings,
and observed flows. Socket-service serves live streams only; trace-service owns
queryable topology/history.

## Key Changes

### Socket Service

- Add `src/apps/socket-service`.
- Use Gorilla WebSocket.
- No DB, no migrations, no service-owned persistence.
- Use Valkey/Redis only for ephemeral sticky session, presence, and reconnect
  smoothing state.
- Consume Kafka topics from `events.TraceableTopics`.
- Consume `socket.broadcast.requested` for direct display events from services.
- Accept internal service push through `POST /internal/v1/socket/events`.
- Publish accepted internal push into Kafka topic `socket.broadcast.requested`.
- Keep only runtime memory:
  - connected clients
  - per-client bounded queues
  - API key hash map loaded from env
  - Valkey-backed session and presence metadata with TTL
  - topic/flow subscription state
  - optional short-lived recent-event buffer for reconnect smoothing
- Expose public sanitized WebSocket stream for root topology.
- Expose authenticated private WebSocket stream for dashboards.

### Internal API Key Design

- API keys are owned by socket-service config, not auth-service.
- Generate keys with a repo script or documented command before runtime.
- Store plaintext keys only in caller service env:
  - `SOCKET_API_KEY=rrsock_<service>_<secret>`
- Store hashes/scopes in socket-service env:
  - `SOCKET_INTERNAL_API_KEYS=farm-service:<hash>:events:publish,warehouse-service:<hash>:events:publish`
- Socket validates `X-RR-API-Key` by hashing and constant-time compare.
- Rotation is env/deployment based: update key envs and restart affected
  services.
- Fail closed:
  - missing key returns `401`
  - unknown key returns `401`
  - missing scope returns `403`
- API keys are for service-to-socket publish paths only, not end-user auth.

### Public And Private APIs

Public sanitized APIs:

- `GET /v1/realtime/public/topology/ws?flow_id=...`
- `GET /v1/traces/public/topology/config`
- `GET /v1/traces/public/topology/history?flow_id=&cursor=&limit=`

Private authenticated APIs:

- `GET /v1/realtime/stream?scope=dashboard`
- `GET /v1/traces/topology/history?flow_id=&entity_id=&cursor=&limit=`

Internal API-key protected API:

- `POST /internal/v1/socket/events`

`POST /internal/v1/socket/events` publishes to Kafka topic
`socket.broadcast.requested`; socket-service still does not persist.

### Trace-Derived Topology, History, And Query

- `trace-service` stores and queries topology/demo history.
- Extend trace projection with fields needed for topology playback:
  - `flow_id`
  - `node_id`
  - `edge_id`
  - `pattern`
  - `source_service`
  - `visibility`
  - `trace_id`
  - sanitized display payload
- Derive these fields from CloudEvent topic/source/entity IDs.
- Add shared topic mapping in code, not admin-editable config:
  - topic to `flow_id`
  - topic to source node
  - topic to target node
  - topic to pattern
  - topic to public/private visibility
- Public topology config endpoint is generated from trace-service observed
  flows plus code-owned canonical mapping.
- If no trace events exist yet, return canonical empty flows with nodes/edges
  from mapping and `last_seen_at = null`.
- PostgreSQL remains fallback/source read model.
- Elasticsearch remains fast query/search projection.
- Public trace history endpoint returns sanitized records only.

### Private Stream Scope

Private stream scope means filtering realtime data by business permission, not
adding a new OAuth concept.

- `ADMIN`: all private events.
- `STORE_MGR`: only matching `store_id` from token `store_ids`.
- `WAREHOUSE_MGR`: only matching `warehouse_id` from token `warehouse_ids`.
- `DRIVER`: only assigned `driver_id` or assigned `shipment_id`.
- Other roles: only matching their domain IDs, otherwise deny or sanitize.

Sprint 6/7 default:

- JWT and Casbin route permission are required for private stream.
- Then filter events by `store_id`, `warehouse_id`, and `driver_id` when those
  fields exist.
- Events without a private scope are visible only to `ADMIN` or converted to a
  public sanitized event.

### OTel Versus Manual Events

- Use Kafka/domain events for business/demo state:
  - order created
  - payment intent/completed/failed
  - warehouse reserved
  - logistics assigned/status/GPS/returned
  - notification/socket requested
- Use OTel only to attach `trace_id`, service spans, latency, and correlation.
- Do not depend on querying OTel/SigNoz as the primary public demo history
  source.
- If needed later, trace-service may enrich stored topology events with OTel
  data by `trace_id`.

### KrakenD Realtime Plan

- First attempt: route realtime through KrakenD.
- Spike WebSocket upgrade support in the current KrakenD image/config.
- If WebSocket proxy works:
  - expose `/v1/realtime/public/topology/ws`
  - expose `/v1/realtime/stream`
- If WebSocket proxy does not work:
  - client connects direct to `NEXT_PUBLIC_SOCKET_URL` in dev/demo.
- Document the chosen fallback in `docs/engineering/TROUBLESHOOTING.md` or the active ticket evidence.

## UI Topology Design

- Keep one `ArchitectureTopology` component.
- Root `/` and `/dashboard/topology-mesh` reuse the same component.
- UI loads topology config from trace-service.
- UI loads history from trace-service.
- UI subscribes realtime from socket-service.
- WebSocket failure falls back to polling trace history.

Canonical node IDs:

- `client.web`
- `gateway.krakend`
- `identity.kratos`
- `identity.hydra`
- `service.auth`
- `service.farm`
- `service.retail`
- `service.payment`
- `service.warehouse`
- `service.logistics`
- `service.trace`
- `service.audit`
- `service.socket`
- `infra.kafka`
- `infra.postgres`
- `infra.elasticsearch`
- `infra.cassandra`
- `infra.valkey`
- `infra.otel`

Canonical flow IDs:

- `flow.auth.oidc-login`
- `flow.farm.harvest-to-pickup`
- `flow.warehouse.intake-processing-stock`
- `flow.retail.paid-order-fulfillment`
- `flow.logistics.delivery-return-to-base`
- `flow.trace.public-qr`
- `flow.realtime.socket-push-pull`
- `flow.authz.policy-sync`

## Docs And Tests

- Update `docs/architecture/ARCHITECTURE.md` with socket-service, env-owned internal API
  keys, trace-derived topology/history, and OTel correlation design.
- Update `docs/requirements/SPEC.md` with public/private realtime behavior and
  push/pull demo UX.
- Update `docs/engineering/LOCAL_DEVELOPMENT.md` or `docs/engineering/TROUBLESHOOTING.md` with local run/debug steps.
- Update `README.md` with socket-service port map.
- Add ADR: socket-service is transport/config cache only; trace-service owns
  history.

Test coverage:

- Socket sanitizer and event normalizer.
- API key hash validation, scope validation, and fail-closed behavior.
- Trace topology history query and public sanitization.
- Trace topic-to-flow/node/edge mapping.
- Empty trace state still returns canonical topology skeleton.
- Private history and stream filtering by role/entity scope.
- UI config-driven topology render.
- UI live WebSocket update and history fallback.
- Playwright public root no-auth topology replay.
- KrakenD realtime route verification or documented fallback verification.

## Assumptions

- Socket-service must not use DB.
- No SSE in Sprint 6/7; realtime transport is WebSocket only.
- Socket-service may use Valkey/Redis for ephemeral session/presence/reconnect
  state only.
- Socket internal API keys are env/deployment managed.
- Topology is never runtime admin-editable.
- Trace-service is the correct owner for queryable topology/demo history.
- Kafka/domain events are required for business timeline; OTel is correlation
  metadata.

## Final Design Locks

### Realtime Transport

- Use Gorilla WebSocket only.
- Do not implement SSE.
- KrakenD realtime spike focuses only on WebSocket upgrade support.
- If KrakenD cannot support WebSocket reliably, the demo uses direct
  `NEXT_PUBLIC_SOCKET_URL` to socket-service and documents the limitation.

### Socket Session State

- Socket-service remains DB-free.
- Use Valkey/Redis for short-lived socket coordination:
  - `socket:session:<session_id>` stores user/role/channel/flow/connected
    instance with TTL.
  - `socket:presence:<role_or_scope>` stores active connection counters with
    TTL.
  - `socket:last-event:<flow_id>` stores latest compact display event or event
    pointer for reconnect smoothing with TTL.
- Valkey data is not durable business history.
- Trace-service remains the durable history/query owner.

### Scaling And Sticky Routing

- Sprint 6/7 default is one socket-service replica.
- If multiple socket replicas run, sticky routing is required so one browser
  session stays on one instance.
- Valkey stores session metadata, but it does not solve cross-replica Kafka
  fanout by itself.
- Multi-replica socket fanout requires a later design: per-instance Kafka
  consumer groups or Valkey pub/sub.
- Do not claim horizontal socket scaling is solved in Sprint 6/7.

### Internal API Key Generation

- Add a keygen command/script that outputs plaintext key plus hash.
- Plaintext keys live only in caller service env as `SOCKET_API_KEY`.
- Socket-service env stores service name, key hash, and scopes.
- Validate `X-RR-API-Key` by hashing and constant-time compare.
- Missing/invalid key returns `401`.
- Valid key without `events:publish` returns `403`.

### Trace Projection Schema

- Add trace projection fields:
  - `flow_id`
  - `node_id`
  - `edge_id`
  - `pattern`
  - `source_service`
  - `visibility`
  - `display_payload`
- PostgreSQL stores source/fallback projection.
- Elasticsearch stores fast query projection.

### Public Sanitizer

- Use allowlist, not blacklist.
- Public payload may include only:
  - `topic`
  - `status`
  - `flow_id`
  - `node_id`
  - `edge_id`
  - `source_service`
  - `trace_id`
  - `occurred_at`
  - demo-safe short IDs
- Drop raw payload by default.

### GPS And High-Frequency Events

- Socket-service coalesces `logistics.gps.updated`.
- Send latest location/status per shipment/driver at a bounded interval, for
  example 500ms.
- Trace-service history policy remains separate.
- Topology projection may store compact/latest display payload for high-volume
  GPS events.

### `socket.broadcast.requested` Contract

- Add explicit event payload struct in `pkg/events`.
- Fields:
  - `event_id`
  - `channel`
  - `flow_id`
  - `node_id`
  - `edge_id`
  - `visibility`
  - `status`
  - `source_service`
  - `payload`
  - `occurred_at`
---
 1. Gap về Tính toàn vẹn dữ liệu (Transactional Integrity)
  Hệ thống đặt mục tiêu Goal G3: Transactional Outbox + Inbox nhưng triển khai thực tế đang rất sơ sài:

   * Vấn đề Outbox: Hiện tại, các service (như retail-service, payment-service) gọi trực tiếp producer.Publish ngay trong hoặc sau transaction DB. 
       * Pain Point: Nếu hệ thống crash ngay sau khi commit DB nhưng trước khi gửi được tin nhắn tới Kafka, sự kiện sẽ bị mất vĩnh viễn. Một hệ thống chuẩn Prod cần có bảng outbox và một worker (như Debezium hoặc poller) để đảm bảo "At-least-once delivery".
   * Vấn đề Inbox (Idempotency): Chỉ có payment-service và trace-service thực hiện kiểm tra event_id trùng lặp một cách thủ công. 
       * Pain Point: Việc triển khai không đồng bộ. Các service quan trọng như warehouse và logistics nếu nhận trùng event từ Kafka (do network retry) sẽ xử lý lại từ đầu, dẫn đến sai lệch tồn kho hoặc trạng thái đơn hàng.

  2. Gap về Kiến trúc CQRS & Projection (Goal G4)
  trace-service đang đóng vai trò là Read-model cho toàn hệ thống nhưng thiết kế đang bị "thắt nút cổ chai":

   * Vấn đề Coupling: Hàm HandleEvent trong trace-service thực hiện việc rebuildDocument (ghi vào SQL và đẩy sang Elasticsearch) đồng bộ ngay trong luồng nhận event.
       * Pain Point: Nếu Elasticsearch phản hồi chậm hoặc bị down, toàn bộ luồng tiêu thụ sự kiện (Kafka Consumer) của trace-service sẽ bị block. Đúng ra luồng ghi (SQL) và luồng chiếu (ES Projection) phải tách biệt và bất đối xứng (Asynchronous Projection).

  3. Gap về Độ tin cậy của Kafka (Resiliency)
   * Vấn đề Xử lý lỗi: Trong pkg/kafka/consumer.go, nếu hàm handler trả về lỗi, Consumer chỉ log lỗi và... bỏ qua tin nhắn đó để đọc tin tiếp theo.
       * Pain Point: Không có cơ chế Retry (Exponential Backoff) và không có Dead Letter Queue (DLQ). Các lỗi logic tạm thời (ví dụ: lỗi DB lock, lỗi network service khác) sẽ khiến business flow bị đứt gãy mà không có cách nào hồi phục tự động.

  4. Gap về Quan sát & Telemetry (Goal G5)
  Mặc dù có SigNoz và OpenTelemetry trong infrastructure:

   * Vấn đề Tracing Propagation: Các file .env.example đang để OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317.
       * Pain Point: Khi chạy trong Docker, localhost trỏ vào chính container đó chứ không phải OTel Collector. Hầu hết các service sẽ không thể gửi được Trace data lên SigNoz nếu không sửa thành otel-collector:4317. Ngoài ra, việc truyền Trace-ID qua Kafka headers đã có code nhưng
         chưa được kiểm thử thực tế (E2E Tracing).

  5. Gap về Bảo mật & Authorization (Goal G7)
   * Vấn đề Internal API Key: Các service gọi lẫn nhau (như retail gọi socket) dùng chung một secret đơn giản.
       * Pain Point: Chưa có cơ chế xoay vòng key (rotation) hoặc phân quyền chi tiết cho từng service (Service-to-Service AuthZ).
   * Vấn đề Webhook Security: Webhook Stripe đã có HMAC verify nhưng key đang được fix cứng hoặc demo.
       * Pain Point: Thiếu cơ chế xử lý Idempotency triệt để cho Webhook ở tầng Gateway trước khi xuống Service.

  6. Gap về Frontend (Trải nghiệm người dùng)
  Như bạn đã nhận ra:
   * Vấn đề Mock UI: Dashboard của Retail và Logistics hiện tại chủ yếu là đồ họa tĩnh để "show" concept Saga/Map.
       * Pain Point: Không thể thực hiện một buổi Demo thực tế "End-to-End" (từ lúc nhấn mua hàng cho đến lúc thấy xe chạy trên bản đồ) mà không phải can thiệp bằng các công cụ bên ngoài (như Postman/CURL).

  ---

  Tổng kết: Danh sách ưu tiên cần xử lý để "Ready to Prod"

   1. Cấp thiết nhất: Implement logic Auto-Milestones cho Simulation và Action Panel cho Retail để khép kín luồng Demo.
   2. Kỹ thuật: Chuẩn hóa bảng outbox cho các service lõi để đảm bảo không mất event.
   3. Kỹ thuật: Bổ sung cơ chế Retry/DLQ cho Kafka Consumer để tăng độ bền vững (Resiliency).
   4. Cấu hình: Sửa lại các OTel Endpoint và verify luồng Trace đổ về SigNoz.

  Bạn muốn tôi tập trung xử lý "điểm đau" nào trước trong danh sách này? (Tôi khuyên nên bắt đầu từ mục 1 để có cái nhìn trực quan nhất).
