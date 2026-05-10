# 📋 Master Strategy Discussion: Sprints 4 — 7

**Trạng thái:** `🚧 IN DISCUSSION`
**Thành phần tham gia:** PO, TechLead, BA, Thư ký (Agent)
**Mục tiêu:** Thống nhất luồng nghiệp vụ tổng thể và cách áp dụng các kỹ thuật cao cấp cho giai đoạn còn lại của Runtime Roasters.

---

## 🏗️ 1. Bức tranh Tổng thể (The Big Picture)

Sau khi hoàn tất **Sprint 3 (Farm Management)**, chúng ta đã có "Gốc" của chuỗi cung ứng. Giai đoạn tiếp theo là xây dựng "Thân" và "Ngọn", biến các Microservices rời rạc thành một thực thể thống nhất.

### Lộ trình 4 giai đoạn tới:
1.  **Giao dịch Phân tán (Sprint 4)**: Đảm bảo Đặt hàng & Giữ hàng nhất quán (Saga).
2.  **Vận tải Thời gian thực (Sprint 5)**: Điều phối xe và theo dõi GPS (Redis Geo).
3.  **Minh bạch Dữ liệu (Sprint 6)**: Truy xuất nguồn gốc và Giám sát (Elasticsearch/Cassandra).
4.  **Trải nghiệm & Khả năng Phục hồi (Sprint 7)**: Dashboard điều khiển trung tâm và Zero Trust (Chaos/mTLS).

---

## 📝 2. Biên bản Thảo luận (Discussion Log)

*(Thư ký sẽ cập nhật nội dung tại đây dựa trên ý kiến của PO và team)*

### 📝 2. Biên bản Thảo luận (Discussion Log)

#### Quyết định quan trọng:
- **Mock Payment**: Tạm thời sử dụng `mock-payment-service` trong các giai đoạn đầu để tập trung hoàn thiện luồng Saga. Tích hợp Stripe thật sự sẽ dời xuống các Sprint cuối.
- **Epic-based Sprints**: Chia nhỏ lộ trình. Mỗi Sprint sẽ tập trung hoàn toàn vào một Epic kỹ thuật/nghiệp vụ duy nhất để đảm bảo chất lượng và tính tập trung.

#### Lộ trình Epic mới (Sprint 4 — 12):

1. **Sprint 4 (Epic: Security & Core Infra Hardening)**
   - Nội dung: Triển khai PgBouncer, hoàn tất E2E Auth Integration và Token Revocation (Logout).
   - Mục tiêu: Nền tảng hạ tầng và bảo mật đạt chuẩn "Production-ready".

2. **Sprint 5 (Epic: Distributed Order Orchestration)**
   - Nội dung: Retail Service, Transactional Outbox Pattern.
   - Mục tiêu: Ghi đơn hàng và phát tán sự kiện một cách nguyên tử.

3. **Sprint 6 (Epic: Inventory Consistency & Saga Lite)**
   - Nội dung: Warehouse Service, Saga Participant, Inbox Pattern (Idempotency T2), Mock Payment logic.
   - Mục tiêu: Hoàn thành luồng Saga 2 bước (Retail-Warehouse) đảm bảo tính lũy đẳng.

4. **Sprint 7 (Epic: Logistics & Real-time Delivery)**
   - Nội dung: Logistics Service, Basic Shipment management.
   - Mục tiêu: Tự động điều phối vận chuyển khi kho đã giữ hàng.

5. **Sprint 8 (Epic: Geographic Intelligence)**
   - Nội dung: Tích hợp Redis Geo, GPS Simulator, Real-time Tracking UI.
   - Mục tiêu: Theo dõi vị trí tài xế trên bản đồ theo thời gian thực.

6. **Sprint 9 (Epic: System-wide Traceability & Search)**
   - Nội dung: Trace Service, CQRS Read Model, Elasticsearch Integration.
   - Mục tiêu: Tìm kiếm và truy xuất toàn bộ lịch sử hạt cà phê trong < 100ms.

7. **Sprint 10 (Epic: Reliability & Observability)**
   - Nội dung: Full OpenTelemetry instrumentation, Prometheus/Grafana, Cassandra Audit Log.
   - Mục tiêu: Hệ thống hoàn toàn minh bạch (Transparent) và có Audit log không thể sửa đổi.

8. **Sprint 11 (Epic: Real-world Commerce Integration)**
   - Nội dung: Xây dựng Payment Service thật, tích hợp Stripe API, hoàn thiện luồng Saga 3 bước.
   - Mục tiêu: Xử lý thanh toán thực tế và cơ chế hoàn tiền (Refund) tự động.

9. **Sprint 12 (Epic: Control Plane & The Grand Finale)**
   - Nội dung: Monitor Service (SSE), React Flow System Mesh, Chaos Control Panel, mTLS.
   - Mục tiêu: Trực quan hóa toàn bộ "vũ trụ" Microservices và chứng minh khả năng tự phục hồi.

---

## ✅ 3. Quyết định & Action Items (Decisions)

*(Ghi lại các chốt chặn sau thảo luận)*

#### Chi tiết Chủ đề 1: Sprint 4 - Security & Core Infra Hardening

**1. PgBouncer: Chiến lược Connection Pooling**
- **Kỹ thuật (TechLead)**:
    - Sử dụng image `bitnami/pgbouncer`.
    - Mode: `transaction` (phù hợp với Go microservices dùng short-lived connections).
    - Cấu hình: Giới hạn `default_pool_size=20` tới Postgres, nhưng cho phép `max_client_conn=500` từ phía microservices.
    - Cập nhật: Toàn bộ `.env` của service sẽ đổi từ port `54321` (Postgres direct) sang `6432` (PgBouncer).
- **Business (BA)**: Đảm bảo tính ổn định của hệ thống khi mở rộng quy mô (scale) số lượng mẻ hàng và giao dịch đặt mua.

**2. E2E Auth Integration: Kiểm chứng "Two-Gate Model"**
- **Kỹ thuật (TechLead)**: Kiểm soát 3 chốt chặn:
    - `KrakenD (Gate 1)`: Verify scope (VD: `farm:read`).
    - `gRPC Interceptor`: Trích xuất identity từ metadata cho các cuộc gọi liên dịch vụ.
    - `OTel Correlation`: Gắn `user_id` vào mọi Trace span trên SigNoz.
- **Business (BA)**: Bảo vệ tuyệt đối dữ liệu nông trại, đảm bảo quyền riêng tư và bảo mật giữa các Farmer.

**3. Token Revocation: Luồng Logout An toàn**
- **Kỹ thuật (TechLead)**: Triển khai Distributed Blacklist.
    - Cơ chế: Lưu `jti` bị thu hồi vào Redis với TTL.
    - Logic: Middleware check Redis `EXISTS`. Chọn phương án **Fail-closed** (Chặn request nếu Redis sập) để tối đa bảo mật.
- **Business (BA)**: Cung cấp tính năng Logout thực sự cho người dùng, bảo vệ tài khoản trong trường hợp bị xâm nhập.
