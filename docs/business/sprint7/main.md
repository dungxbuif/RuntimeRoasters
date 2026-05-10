# Sprint 7: Control Plane Visualization Dashboard

**Goal:** Hoàn thiện giao diện đồ họa trực quan (System Mesh) cho toàn bộ hệ thống, tích hợp khả năng giám sát luồng dữ liệu thời gian thực và quản trị "Chaos".

---

## 🎯 Sprint Goal

> **Hệ thống có một dashboard trung tâm cho phép người xem nhìn thấy data flow chuyển động giữa các microservices và thực hiện các thử nghiệm khả năng phục hồi (Resiliency).**

---

## 📋 Roadmap & Flows

Sprint cuối cùng tập trung vào trải nghiệm người xem (Showcase Experience):

### 1. Monitor Service (Real-time Broadcast)
- Consume tất cả events từ Kafka.
- Sử dụng Server-Sent Events (SSE) để bắn data realtime lên Frontend mà không cần polling.

### 2. System Mesh Visualization (React Flow)
- Xây dựng bản đồ hệ thống (Topology map) bằng React Flow.
- Hiển thị các "đường dẫn dữ liệu" (Animated Edges) khi có event phát sinh giữa các service.

### 3. Chaos Control Panel
- Cho phép Admin giả lập lỗi (ngắt một service, làm chậm network, set tồn kho ảo = 0).
- Quan sát cách hệ thống (Saga Rollback, Outbox Pattern) phản ứng ngay trên Dashboard.

### 4. Consumer QR Experience
- Hoàn thiện trang đích cho khách hàng cuối khi quét mã QR.
- Hiển thị "Hành trình hạt cà phê" với hiệu ứng animation đẹp mắt.

### 5. Zero Trust Security (Advanced)
- Áp dụng mTLS (Mutual TLS) cho các kết nối gRPC nội bộ giữa các services để tăng cường bảo mật tối đa.

---

## 🎫 Tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-31](./RR-31.md) | Monitor Service & SSE Broadcasting | 🕒 To Do |
| [RR-32](./RR-32.md) | Frontend: React Flow System Topology | 🕒 To Do |
| [RR-33](./RR-33.md) | Chaos Control Panel (Resiliency UI) | 🕒 To Do |
| [RR-34](./RR-34.md) | Consumer QR Traceability Flow (Final UI) | 🕒 To Do |
| [RR-26](../sprint5/RR-26.md) | Security: Internal gRPC mTLS | 🕒 To Do |
