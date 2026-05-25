# RR-12: Technical Design — Resilient Centralized Authorization

**Branch:** `feat/RR-12`
**Status:** `DRAFT`
**Author:** Tech Lead
**Date:** 2026-05-15

---

## 1. Context & Goal

Sau khi **RR-11** hoàn thiện việc xác thực danh tính (Authentication), hệ thống cần một cơ chế phân quyền (Authorization) mạnh mẽ. Thay vì mỗi service tự quản lý policy riêng lẻ (Decentralized), chúng ta chuyển sang mô hình **Centralized Authorization** nhưng đảm bảo tính **Resilient** (Khả năng chịu lỗi cao).

Mục tiêu cốt lõi:
- **Centralized Management:** Mọi chính sách phân quyền được quản lý tập trung tại `auth-service`.
- **Distributed Enforcement:** Các service (Readers) thực thi quyền tại chỗ với hiệu năng cực cao (Local Memory).
- **High Availability:** Hệ thống phân quyền vẫn hoạt động ngay cả khi `auth-service` hoặc Kafka gặp sự cố.

---

## 2. Architecture: Centralized Writer & Distributed Readers

Hệ thống được thiết kế theo mô hình **Single Writer - Multiple Readers**:

### 2.1 Centralized Writer (Auth Service)
- **Vai trò:** Chủ sở hữu (Owner) duy nhất của cơ sở dữ liệu chính sách (Postgres).
- **Trách nhiệm:** 
    - Cung cấp API để CRUD các chính sách (Policies).
    - Lưu trữ chính sách vào Postgres sử dụng `gorm-adapter/v3`.
    - Phát tín hiệu thay đổi (Change Signals) qua Kafka.
    - Cung cấp gRPC endpoint `GetPolicies` để các service khác lấy bản snapshot.

### 2.2 Distributed Readers (Business Services)
- **Vai trò:** Thực thi kiểm tra quyền (Enforcement).
- **Trách nhiệm:**
    - Duy trì một bản sao chính sách trong bộ nhớ (Local Memory Enforcer) bằng `casbin/v3`.
    - Tự động đồng bộ hóa với `auth-service` thông qua chiến lược 3 trụ cột.

---

## 3. Chiến lược Resilience 3 trụ cột (The 3-Pillar Strategy)

Để đảm bảo các Reader luôn có dữ liệu chính sách mới nhất và không bao giờ bị "chết đứng", chúng ta áp dụng 3 cơ chế đồng bộ song song:

### Trụ cột 1: gRPC Snapshot (Bootstrapping & Self-healing)
- **Khi khởi động:** Reader gọi gRPC `GetPolicies` tới `auth-service` để tải toàn bộ chính sách hiện có. Nếu `auth-service` không khả dụng, Reader sẽ thử lại với cơ chế **Exponential Backoff**.
- **Vai trò:** Đảm bảo Reader có dữ liệu nền tảng để bắt đầu phục vụ.

### Trụ cột 2: Kafka Live (Real-time Updates)
- **Khi có thay đổi:** `auth-service` publish một message vào Kafka topic `auth.policy.changed`.
- **Tại Reader:** Một `Watcher` (sử dụng `casbin-go-cloud-watcher`) lắng nghe Kafka. Khi nhận tín hiệu, Reader lập tức gọi `LoadPolicy()` để cập nhật bản sao trong memory.
- **Vai trò:** Đảm bảo độ trễ cập nhật chính sách ở mức mili giây (Near real-time).

### Trụ cột 3: Polling Fallback (Safety Net)
- **Định kỳ:** Reader tự động chạy một Background Ticker (ví dụ: mỗi 5-10 phút) để gọi lại gRPC `GetPolicies`.
- **Vai trò:** "Lưới an toàn" để bù đắp cho các message Kafka bị mất hoặc các lỗi đồng bộ tiềm tàng, đảm bảo tính nhất quán cuối cùng (Eventual Consistency).

---

## 4. Casbin Configuration

### 4.1 Model Definition (`model.conf`)
Sử dụng mô hình RBAC với Hierarchy và Matcher thông minh.

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
# 1. Admin bypass
# 2. Check role hierarchy & object match (support glob) & action match
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))
```

---

## 5. Package Structure (`pkg/base/casbin`)

Gói dùng chung này sẽ được các Reader import để triển khai Resilience Engine:

```text
src/pkg/base/casbin/
├── engine.go           # Resilience Engine (gRPC + Kafka + Polling)
├── watcher_kafka.go    # Tích hợp Kafka Watcher
├── enforcer_v3.go      # Wrapper cho Casbin v3 Enforcer
├── middleware_gin.go   # Gin AuthZ Middleware
└── interceptor_grpc.go # gRPC AuthZ Interceptor
```

---

## 6. Guidelines cho Developer

1. **Phân quyền tại Service biên:** Mọi service phải sử dụng `interceptor_grpc.go` để bảo vệ các hàm API.
2. **Quyền mặc định:** Các quyền hệ thống cốt lõi phải được định nghĩa trong `default_policies.csv` tại `auth-service`.
3. **Giám sát:** Cần theo dõi các metric về thời gian đồng bộ chính sách và trạng thái kết nối tới Kafka/gRPC của `auth-service`.

---

## 7. Risks & Mitigation

| Rủi ro | Giải pháp |
|---|---|
| **Kafka Down:** Reader không nhận được update tức thì. | Polling Fallback (Trụ cột 3) sẽ tự động đồng bộ sau một khoảng thời gian ngắn. |
| **Auth Service Down:** Reader mới khởi động không lấy được snapshot. | Sử dụng Retry với Backoff. Nếu vẫn lỗi, Reader có thể nạp một bộ chính sách "hardcoded" tối thiểu để duy trì các chức năng khẩn cấp. |
| **Data Inconsistency:** Dữ liệu memory khác với DB. | Cơ chế Polling đảm bảo hệ thống tự chữa lành (Self-healing). |
