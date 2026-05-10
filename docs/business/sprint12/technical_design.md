# Technical Design: Sprint 12 — Control Plane & Final Hardening (The Grand Finale)

Mục tiêu: Trực quan hóa toàn bộ "vũ trụ" Microservices và gia cố lớp giáp hạ tầng cuối cùng.

---

## 1. System Mesh Visualization (Monitor Service)
*Tham chiếu: Triết lý "Trong suốt" của UrbanX + Dashboard đồ họa*

Xây dựng bộ não giám sát luồng dữ liệu thời gian thực:
### 1.1. Monitor Service (SSE Orchestrator)
- **Source**: Lắng nghe toàn bộ Kafka Topics (`*`).
- **Stream**: Sử dụng **Server-Sent Events (SSE)** để đẩy thông tin sự kiện xuống trình duyệt. 
- **Payload**: `{ source_service, target_service, event_type, status }`.

### 1.2. Frontend: React Flow Topology
- Vẽ sơ đồ kiến trúc động.
- Khi nhận SSE event: Kích hoạt hiệu ứng **Moving Dots** trên các Edges (đường nối) giữa các services để người xem thấy được data đang "chảy" trong hệ thống.

---

## 2. Chaos Control Panel (Resiliency Showcase)
*Tham chiếu: Khả năng tự phục hồi của hệ thống*

Xây dựng bộ điều khiển thử nghiệm lỗi:
- **API**: `/v1/admin/chaos/inject-fault`.
- **Logic**: Cho phép Admin làm "sập" giả lập một service hoặc DB.
- **Quan sát**: Người xem sẽ thấy trên Dashboard luồng Saga tự động kích hoạt **Compensating Transactions** (Rollback) khi gặp lỗi.

---

## 3. Infrastructure Hardening (Deferred Tasks)
*Chốt chặn cuối cùng cho tiêu chuẩn Enterprise*

### 3.1. PgBouncer Deployment
- Triển khai cụm PgBouncer (Transaction mode) để tối ưu pool kết nối PostgreSQL cho toàn bộ 6+ Microservices.
- Cấu hình Auth Type MD5/Docker Secrets.

### 3.2. Zero Trust gRPC mTLS
- Cấu hình xác thực 2 chiều (Mutual TLS) cho mọi giao tiếp gRPC nội bộ.
- Mọi service phải có Certificate hợp lệ mới được phép gọi nhau, ngăn chặn tấn công "Man-in-the-middle" trong mạng nội bộ.

---
*TechLead Signed-off: 2026-05-10*
