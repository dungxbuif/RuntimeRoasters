# Sprint 4: The Resilient Farm (Transactional Outbox)

**Trạng thái:** ✅ Hoàn thành (Completed)
**Mục tiêu:** Đảm bảo mọi mẻ thu hoạch được ghi nhận 100% không mất dữ liệu ngay cả khi hệ thống phân tán gặp sự cố. Triển khai mẫu thiết kế Transactional Outbox Blueprint cho toàn dự án.

---

## 📋 Trạng thái Ticket (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-4.0](./RR-4.0/ticket.md) | [Tech] Refactor Farm Service: Migration & DI cleanup | ✅ Done | Tech Lead |
| [RR-4.1](./RR-4.1/ticket.md) | [BA] Khai báo mẻ thu hoạch (Harvesting Management) | ✅ Done | Farmer |
| [RR-4.2](./RR-4.2/ticket.md) | [Tech] Transactional Outbox: Event Reliability | ✅ Done | Tech Lead |
| [RR-22](./RR-22.md) | [Tech] Trace Service & CQRS Bootstrap | ✅ Done | Backend |

---

## 💡 Tầm nhìn Nghiệp vụ (Business Vision)
- **Data Integrity:** "Một hạt rơi, hệ thống biết". Tuyệt đối không để mất dữ liệu thu hoạch khi chuyển giao sang nhà máy.
- **Real-time Awareness:** Các bộ phận phía sau (Warehouse) nhận được thông báo ngay khi có hàng rời khỏi nông trại.
- **Professionalism:** Áp dụng chuẩn CloudEvents 1.0 để giao tiếp giữa các service chuyên nghiệp và dễ mở rộng.

## 📊 Kết quả đạt được (Sprint Result)
- Triển khai thành công bảng `outbox_events` và Relay Worker trong Farm Service.
- Hoàn thiện API thu hoạch tích hợp cơ chế GORM Transaction (Atomic Write).
- Sự kiện thu hoạch đã được bắn lên Kafka thành công và sẵn sàng để Warehouse tiêu thụ.
