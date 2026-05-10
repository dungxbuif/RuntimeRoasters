# Technical Design: Sprint 11 — Real-world Commerce (Reference Aligned)

Mục tiêu: Chuyển đổi từ Mock Payment sang giao dịch thực tế qua Stripe, áp dụng lớp bảo vệ ACL và cơ chế Idempotency tuyệt đối.

---

## 1. Stripe Integration via ACL (Anti-Corruption Layer)
*Tham chiếu: UrbanX Payment Gateways & IPaymentGateway*

Chúng ta sẽ "hóa giải" sự phụ thuộc vào Stripe SDK bằng cách:
- **Interface**: `pkg/auth/payment` (hoặc internal payment) định nghĩa `ProcessPayment(ctx, amount, currency, method_id)`.
- **Implementation**: `StripeGateway` thực thi interface này, đóng gói toàn bộ logic gọi `stripe-go` SDK.
- **Benefit**: Có thể đổi sang PayPal/VNPay chỉ bằng cách cấu hình lại DI mà không sửa đổi `PaymentUseCase`.

---

## 2. Bảo vệ giao dịch (Stripe Webhook & Idempotency)
*Tham chiếu: UrbanX Webhook Security*

Giao dịch tiền tệ yêu cầu độ tin cậy "3 tầng":
### 2.1. Webhook Signature Verification
- Payment Service cung cấp endpoint `/v1/payments/webhook`.
- Sử dụng Stripe Webhook Secret để verify **HMAC signature** của từng request, chống giả mạo request từ bên ngoài.

### 2.2. Transactional Inbox (Tầng 2 Idempotency)
- Stripe có thể gửi Webhook lặp lại.
- **Logic**: Kiểm tra `Stripe_Event_ID` trong bảng `inbox_events` của Payment DB trước khi xử lý.

---

## 3. Saga Phase 2: Distributed Refund
*Tham chiếu: UrbanX Compensation Flow*

Khi một bước sau Payment (VD: Logistics hết xe) thất bại:
1. **Trigger**: Logistics Service bắn `ShipmentFailed`.
2. **Action**: Payment Service nghe sự kiện này và gọi `stripe.Refund`.
3. **Consistency**: Đảm bảo khách hàng được hoàn tiền tự động nếu hệ thống không thể giao hàng.

---
*TechLead Signed-off: 2026-05-10*
