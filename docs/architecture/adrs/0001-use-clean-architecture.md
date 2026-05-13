# ADR 0001: Áp dụng Clean Architecture (Consumer-Owned Interfaces)

## Trạng thái
**Accepted**

## Bối cảnh (Context)
Dự án RuntimeRoasters bắt đầu từ Sprint 1 với mục tiêu xây dựng một nền tảng Microservices dễ bảo trì, dễ mở rộng, và đặc biệt là dễ viết Unit Test. Các framework truyền thống (như MVC) thường dẫn đến việc các package phụ thuộc chéo vào nhau (Circular dependencies), khiến code nghiệp vụ (Business logic) bị khóa chặt vào database.

## Quyết định (Decision)
Áp dụng **Clean Architecture Style** cho toàn bộ Microservices.
Điểm khác biệt cốt lõi:
- **Consumer-Owned Interfaces:** Tầng `usecase` (consumer) sẽ tự khai báo các Interface mà nó cần (ví dụ: `FarmRepository`). 
- Tầng `infrastructure` (ví dụ: `postgres`) sẽ implement các Interface này (nhờ cơ chế Duck Typing của Go) thay vì khai báo Interface tại tầng `domain` hay `infrastructure`.
- `domain/` chỉ chứa Pure Go Structs và không import bất kỳ package nào ngoài Standard Library.

## Hậu quả (Consequences)
- **Tích cực:** Tầng `usecase` hoàn toàn tách biệt khỏi cách dữ liệu được lưu trữ. Dễ dàng dùng thư viện `mockery` để generate mock cho Unit Test.
- **Tiêu cực:** Các Developer mới chuyển từ ngôn ngữ khác (Java, C#) có thể mất thời gian làm quen với việc Interface không nằm cùng file với Struct implement nó.

## Nguồn tham khảo
- **Sprint:** Sprint 1 (Bootstrap Phase).
- Xem chi tiết triển khai tại: [Clean Architecture Concepts](../README.md#-clean-architecture-framework--runtimeroasters).
