# Sprint 5: Logistics & Real-time Tracking

**Goal:** Triển khai hệ thống Logistics, quản lý lộ trình vận chuyển và theo dõi vị trí (GPS) theo thời gian thực sử dụng Redis Geo và Kafka.

---

## 🎯 Sprint Goal

> **Hệ thống có thể điều phối vận chuyển, cập nhật tọa độ GPS liên tục và cho phép theo dõi thời gian thực.**

---

## 📋 Roadmap & Flows

Sprint này tập trung vào khả năng xử lý dữ liệu thời gian thực:

### 1. Logistics Service
- Quản lý chuyến xe và tài xế.
- Giao việc (Assign Shipment) khi nhận được sự kiện `StockReserved`.

### 2. Redis Geo Integration
- Lưu trữ và truy vấn vị trí hiện tại của xe chở hàng bằng lệnh không gian của Redis.

### 3. GPS Simulator
- Giả lập dữ liệu GPS liên tục gửi về hệ thống để test tính năng thời gian thực.

### 4. Zero Trust Security
- Áp dụng mTLS (Mutual TLS) cho các kết nối gRPC nội bộ giữa các services để tăng cường bảo mật.

---

## 🎫 Tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-23](./RR-23.md) | Logistics Service Core (Shipments & Drivers) | 🕒 To Do |
| [RR-24](./RR-24.md) | Redis Geo Tracking Integration | 🕒 To Do |
| [RR-25](./RR-25.md) | GPS Simulator & Event Publishing | 🕒 To Do |
