# RuntimeRoasters (OriginFlow) \- Farm-to-Cup Coffee Supply Chain Platform

**RuntimeRoasters** là một nền tảng quản lý chuỗi cung ứng và logistics mô phỏng vòng đời của hạt cà phê từ nông trại đến tách cà phê bán lẻ (Farm-to-Cup).

Dự án được xây dựng như một bản showcase kỹ thuật dành cho mục đích giáo dục và thực hành kiến trúc **Microservices**. Hệ thống được phát triển 100% bằng **Golang**, tuân thủ nghiêm ngặt các nguyên tắc Domain-Driven Design (DDD) và áp dụng các mẫu thiết kế phân tán hạng nặng để đảm bảo tính mở rộng, độ tin cậy và khả năng truy xuất nguồn gốc minh bạch.

---

## 📖 Câu chuyện Doanh nghiệp

Trong văn hóa Việt, cà phê không chỉ là thức uống mà là "sợi dây" kết nối xã hội. Tuy nhiên, sự đứt gãy thông tin giữa người nông dân và tách cà phê trên tay khách hàng vẫn còn hiện hữu. **RuntimeRoasters** ra đời để số hóa toàn bộ "mạch sống" này.

Dự án mô phỏng một hệ sinh thái nơi mỗi hạt cà phê đều có một danh tính số. Từ những "mạch" dữ liệu nhỏ nhất ở nông trường, đi qua các "mạch" xử lý tại nhà máy, đến các "mạch" logistics xuyên suốt đất nước để cuối cùng hội tụ tại những tách cà phê thơm ngon.

Các điểm chạm chính:

1. **Thượng nguồn (Farm):** Nông dân cập nhật diện tích canh tác và khai báo các mẻ cà phê vừa thu hoạch.  
2. **Trung nguồn (Processing & Warehouse):** Quản đốc nhà máy tiếp nhận hạt thô, tiến hành bóc vỏ, phơi khô, rang xay và đóng gói thành lô thành phẩm (Batch ID). Hàng được lưu trữ tại Kho tổng.  
3. **Vận tải (Logistics):** Khi có lệnh điều phối, tài xế nhận cuốc xe, vận chuyển hàng từ nhà máy đến các điểm bán lẻ, liên tục cập nhật tọa độ GPS theo thời gian thực.  
4. **Hạ nguồn (Retail):** Quản lý cửa hàng theo dõi tồn kho tại điểm bán, gửi yêu cầu nhập hàng và tiếp nhận hàng hóa từ tài xế.  
5. **Truy xuất & Quản trị (Traceability & Audit):** Người dùng cuối có thể quét mã QR để xem toàn bộ hành trình của ly cà phê. Quản trị viên có công cụ giám sát toàn cảnh và audit log chống gian lận.

---

## 🧩 Thành phần Hệ thống

| Service | Chức năng chính | Patterns Áp dụng |
| :---- | :---- | :---- |
| **Gateway** | Chốt chặn Authentication, Rate Limiting. | Gateway Pattern |
| **Farm** | Quản lý nông hộ, vườn cây và thu hoạch. | Outbox Pattern |
| **Process** | Chế biến mẻ rang, cấp Batch ID. | Event-Driven |
| **Logistics** | Điều phối xe, tracking GPS qua Valkey. | Real-time Geo |
| **Warehouse** | Giữ chỗ hàng (Reserve), xuất/nhập kho. | Saga (Participant) |
| **Retail** | Cửa hàng đặt hàng, quản lý tiêu thụ. | Saga (Orchestrator) |
| **Trace** | Tổng hợp hành trình hạt cà phê vào Elastic. | CQRS |
| **Audit** | Lưu toàn bộ lịch sử Kafka vào Apache Cassandra. | Event Sourcing (Lite) |

---

## 💻 Ngăn xếp Công nghệ

Hệ thống tận dụng tối đa tài nguyên máy chủ Proxmox:

* **Ngôn ngữ & Giao tiếp:** Go (Golang) 1.22+; Gin (External REST); gRPC (Internal RPC).  
* **Message Broker:** **Apache Kafka**.  
* **Cơ sở dữ liệu đa phương thức:**  
  * **PostgreSQL:** Source of Truth (Transaction, ACID).  
  * **Elasticsearch:** Search & Traceability (CQRS Read-model).  
  * **Apache Cassandra:** Event Store & Audit Log.  
* **Caching & Real-time:** **Valkey** (Thay thế Redis để quản lý Distributed Lock, Cache và GPS Tracking).  
* **Bảo mật & IAM:** Ory Kratos (Identity); Casbin (Decentralized Authorization).  
* **Quan sát:** Prometheus & Grafana.

---

## 🏗️ Các Mẫu Thiết kế Kỹ thuật Nâng cao

1. **Saga Pattern (Choreography):** Quản lý luồng đơn hàng cung ứng từ Cửa hàng ➔ Kho ➔ Logistics. Hệ thống tự động hoàn tiền/hủy đơn (Compensating Actions) khi có lỗi phát sinh ở bất kỳ công đoạn nào.
2. **Transactional Outbox & Inbox:**
   * **Outbox:** Đảm bảo tính nguyên tử (Atomicity) giữa việc lưu DB và gửi event.
   * **Inbox (Idempotency):** Sử dụng `Message_ID` để tránh xử lý trùng lặp sự kiện, đảm bảo tính nhất quán cuối cùng (Eventual Consistency).
3. **CQRS & Real-time Traceability:** 
   * Tách biệt luồng ghi (PostgreSQL) và luồng đọc (Elasticsearch).
   * **CQRS Sync Flow:** Khi có sự kiện nghiệp vụ (ví dụ: `ArticlePublished`, `BatchCreated`), **Trace Service** sẽ tiêu thụ event này từ Kafka và thực hiện logic **Upsert** (Update or Insert) vào Elasticsearch.
   * Điều này đảm bảo dữ liệu truy xuất hành trình (Traceability) luôn sẵn sàng với tốc độ tìm kiếm full-text cực nhanh mà không gây tải cho DB giao dịch.
4. **Standardized API Response (RFC 7807):**
   * Toàn bộ hệ thống áp dụng chuẩn **Problem Details for HTTP APIs** cho các phản hồi lỗi.
   * Điều này giúp client (Dashboard, Mobile App) có thể xử lý lỗi một cách nhất quán: `type`, `title`, `status`, `detail`, và `instance`.
5. **Decentralized Authorization:** Từng service ôm một bộ luật Casbin riêng thay vì kiểm tra quyền tập trung tại Gateway.
6. **Data Integrity (Hash Chaining):** Sử dụng chuỗi băm dữ liệu lưu tại Apache Cassandra để ngăn chặn sự can thiệp, sửa đổi dữ liệu từ bên trong.

---

## 🛡️ Bảo mật & Xác thực: Zero Trust & Gateway Offloading

* **Gateway Offloading (Auth Centralization):** API Gateway đóng vai trò là chốt chặn duy nhất xác thực JWT thông qua **Ory Kratos**. 
    * Gateway thực hiện *Token Introspection*, sau đó giải mã và inject các định danh người dùng (`X-User-ID`, `X-User-Role`, `X-User-Permissions`) vào HTTP Header trước khi chuyển tiếp request vào mạng nội bộ.
    * Các Microservices phía sau chỉ cần tin tưởng vào Header này (đã được bảo vệ bởi mTLS), giúp tinh giản logic xác thực trong từng service Go.
* **mTLS (Mutual TLS) cho gRPC:** Các Microservices gọi nhau bằng gRPC bắt buộc phải xác thực lẫn nhau bằng chứng chỉ số (Certificates).
* **Chữ ký Webhook (HMAC):** Mọi payload từ bên thứ 3 (Stripe, IoT GPS) phải đi kèm chữ ký HMAC. Webhook Service sẽ xác thực chữ ký này trước khi đẩy vào hệ thống Kafka.

---

## 🔍 Distributed Tracing (Dấu vết Phân tán)

Sử dụng **OpenTelemetry & Jaeger** để giám sát luồng đi của request. Ngay khi request chạm vào API Gateway, một `Trace-ID` duy nhất được sinh ra và lan truyền qua tất cả các service thông qua HTTP Headers và gRPC Context. Giao diện Jaeger sẽ hiển thị biểu đồ Gantt Chart, cho phép truy vết chính xác request đã đi qua những service nào, mất bao lâu, và phát sinh lỗi ở dòng code nào trong trường hợp giao dịch đa dịch vụ bị gián đoạn.  
