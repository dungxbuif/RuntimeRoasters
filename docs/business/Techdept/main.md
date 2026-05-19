Chưa thể nói là “đã thống nhất hoàn toàn” theo design ban đầu. Phần core đã được đưa vào code: JWT claims, store_ids, HMAC
  webhook, replay protection, record-level scoping. Nhưng còn một gap quan trọng ở Casbin matcher khiến route-level RBAC cho role
  như STORE_MGR có rủi ro không hoạt động đúng như policy mô tả.

  Auth Flow Hiện Tại

  Gate 1 là JWT validation:

  - HTTP services dùng security.NewHTTPGuards(...).
  - Middleware verify RS256 JWT qua JWKS, issuer http://localhost:4444/.
  - Claims được inject vào context để service layer đọc.

  Code chính:

  - src/pkg/base/security/http.go
  - src/pkg/base/auth/transport/http/middleware.go
  - src/pkg/base/auth/token/jwt.go

  Claims đã mở rộng:

  - sub
  - email
  - role
  - org_id
  - store_ids
  - jti

  JWT parser đọc được cả root claim và nested ext.

  Identity / Hydra / Kratos

  Kratos schema đã có optional store_ids:

  - deployments/kratos/identity.schema.json

  Auth-service khi tạo user đã lưu org_id, store_ids vào traits:

  - src/apps/auth-service/internal/usecase/user_usecase.go

  Consent accept route đã lấy Kratos traits và đưa vào Hydra token session:

  - src/apps/client-app/src/lib/ory/identity-admin.ts
  - src/apps/client-app/src/app/api/auth/consent/accept/route.ts

  Điểm này hiện đã khớp design: source of truth cho store access là store_ids trong identity/JWT.

  Record-Level Authorization

  Đã có helper chung:

  - src/pkg/base/identity/scope.go

  Quy tắc:

  - ADMIN: được xem tất cả ở service-level helper.
  - STORE_MGR: chỉ được thao tác store nằm trong StoreIDs.
  - STORE_MGR không có store_ids: fail closed.
  - Record không có store_id: chỉ admin được xem.

  Đã áp vào:

  - Retail: list stores, create order, get order.
  - Payment: list payments, get payment by order.
  - Logistics: list shipments, deliver shipment, update driver location.
  - Trace: trace events/documents theo store_id.
  - Audit: audit logs theo store_id.

  Webhook Auth

  Webhook không dùng Bearer JWT, đúng với design provider callback.

  Đã triển khai:

  - Stripe: Stripe-Signature: t=<unix>,v1=<hmac_sha256(secret, t + "." + raw_body)>
  - VNPay: X-VNPAY-Timestamp + X-VNPAY-Signature
  - Timestamp tolerance: 5 phút.
  - Missing secret/signature/timestamp hoặc sai signature đều bị reject.
  - Replay/idempotency qua WebhookEvent unique (provider, event_id).

  Code chính:

  - src/apps/payment-service/internal/provider/provider.go
  - src/apps/payment-service/internal/usecase/service.go
  - src/apps/payment-service/internal/domain/models.go

  KrakenD vẫn để webhook public ở route-auth layer và đã forward thêm VNPay timestamp header:

  - deployments/krakend/krakend.json

  Casbin Gap Quan Trọng

  Đây là phần chưa thống nhất hoàn toàn.

  Policy đã đổi sang STORE_MGR:

  - src/apps/auth-service/internal/infrastructure/casbin/default_policies.csv

  Nhưng matcher hiện là:

  m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))

  Trong HTTP middleware, r.sub đang là claims.Role, ví dụ STORE_MGR, không phải user id. Vì matcher chỉ dùng g(r.sub, p.sub) mà
  không có r.sub == p.sub, policy dạng:

  p, STORE_MGR, /v1/stores, read

  có thể không match nếu không có grouping g, STORE_MGR, STORE_MGR.

  Tác động:

  - STORE_MGR route-level RBAC có thể bị deny trước khi vào service.
  - ADMIN cũng không chắc có full /v1/* theo policy p, ADMIN, /v1/*, .*, vì matcher không so sánh trực tiếp ADMIN == ADMIN.
  - ADMIN chỉ chắc chắn inherit các role được khai báo g, ADMIN, STORE_MGR, g, ADMIN, FARM_MANAGER, ...

  Fix nên làm:

  - Sửa matcher thành dạng có direct role policy match:

  m = r.sub == p.sub || g(r.sub, p.sub) || g(r.sub, "admin")

  và vẫn giữ keyMatch / regexMatch quanh object/action.

  Test Status

  Đã pass:

  - JWT claim extraction, gồm email, store_ids, nested ext.
  - Identity scope helper.
  - Payment webhook valid/invalid/idempotent.
  - Targeted backend packages cho auth/payment/retail/logistics/trace/audit.

  Chưa pass full go test ./... vì lỗi sẵn có ngoài scope:

  - apps/warehouse-service copy có space trong import path.
  - scripts có nhiều main trong cùng package.

  Frontend lint cũng fail bởi lỗi sẵn có ngoài phần auth claim route.

  Đánh Giá Cuối

  Backend service-level authorization đã đi đúng design. Webhook protection cũng đã đúng hướng thực tế hơn JWT. Nhưng để nói Auth
  “đã thống nhất hoàn toàn”, cần xử lý Casbin matcher trước, vì hiện Gate 2 có thể chặn sai role policy. Sau khi sửa matcher và
  chạy integration qua KrakenD với JWT thật của STORE_MGR, cơ chế Auth mới đạt trạng thái thống nhất end-to-end.