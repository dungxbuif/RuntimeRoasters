# Sprint 9: Financial Integrity

**Epic Goal:** Tích hợp thanh toán thực tế và tự động hóa quy trình hoàn tiền (Refund).

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-9.1](./RR-9.1/ticket.md) | [BA] Thanh toán đơn hàng qua Stripe | 🕒 To Do | Customer |
| [RR-9.2](./RR-9.2/ticket.md) | [Tech] Stripe Webhook & Automated Refund Saga | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Stripe API Integration:** Xử lý PaymentIntents và Webhooks.
- **HMAC Validation:** Xác thực chữ ký tin nhắn từ Payment Gateway.
- **Full Saga Rollback:** Hoàn tiền tự động nếu chuỗi cung ứng gặp sự cố sau khi đã trừ tiền.
