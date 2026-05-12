# RuntimeRoasters Architecture Documentation

Welcome to the central hub for the RuntimeRoasters system architecture. This documentation is organized following Enterprise Big Tech standards (Diátaxis Framework + ADRs) to ensure maintainability, clear separation of concerns, and structured onboarding.

## 📚 Organization Rules (How to maintain this doc)

To prevent this directory from becoming a chaotic "God Document" dump, all contributors MUST adhere to the following rules:

1. **Diátaxis Framework (Separation of Concerns):**
   - **`concepts/` (Lý thuyết / The WHY):** Dành cho các bài viết giải thích tư tưởng, triết lý thiết kế (e.g., *Tại sao dùng Clean Architecture? Chiến lược Observability là gì?*). Không chứa code chi tiết ở đây.
   - **`reference/` (Tra cứu / The WHAT):** Dành cho tài liệu mang tính chất thông tin cứng (e.g., *Danh sách port, cấu hình `pkg/logger`, API Schema, Database Schema*).
   - **`flows/` (Luồng nghiệp vụ / The HOW):** Chứa các sơ đồ Sequence, Data Flow mô tả cách các thành phần tương tác trong một Use Case cụ thể (e.g., *Luồng đăng nhập Kratos, Luồng checkout*).

2. **ADRs (Architecture Decision Records):**
   - **Bắt buộc:** Mọi thay đổi lớn về kiến trúc, công nghệ, hoặc quy trình (e.g., *Đổi từ MongoDB sang Postgres, chọn GORM thay vì sqlx*) ĐỀU PHẢI được ghi nhận bằng 1 file markdown trong `adrs/`.
   - **Mục đích:** Tránh việc tranh cãi lại các quyết định cũ, và giúp Dev mới hiểu bối cảnh lịch sử của dự án.
   - **Format:** Tuân thủ chuẩn `[Trạng thái] - [Bối cảnh] - [Quyết định] - [Hậu quả]`. Gắn tag rõ ràng tới Sprint hoặc Ticket liên quan.

3. **No Duplication:**
   - Nếu `blueprint.md` đã vẽ sơ đồ tổng thể, đừng vẽ lại nó ở file khác. Hãy trỏ link tới nó.
   - Code snippet chỉ được dùng để minh họa, không copy toàn bộ source code vào tài liệu vì code sẽ out-of-date rất nhanh.

---

## 🗺️ Master Index

### 1. High-Level Overviews
- [01_system_blueprint.md](./01_system_blueprint.md) - Sơ đồ kiến trúc tổng thể, hạ tầng và Roadmap.

### 2. 🧠 Concepts (Lý thuyết cốt lõi)
- [clean-architecture.md](./concepts/clean-architecture.md) - Triết lý Clean Architecture, phân lớp và Directory Mapping.
- [system-wide-standards.md](./concepts/system-wide-standards.md) - Thư viện dùng chung, Outbox, Idempotency và Manual DI.
- [resilient-authz-sync.md](./concepts/resilient-authz-sync.md) - Kiến trúc đồng bộ quyền (Casbin + Kafka + gRPC).
- [observability-strategy.md](./concepts/observability-strategy.md) - Chiến lược giám sát (SigNoz, 4 pillars).

### 3. 📖 Reference (Tra cứu)
- [core-framework-pkg.md](./reference/core-framework-pkg.md) - Hướng dẫn sử dụng `pkg/` (Config, Logger, Errs, Database).
- [infrastructure-port-map.md](./reference/infrastructure-port-map.md) - Danh sách Port, môi trường.
- [data-models.md](./reference/data-models.md) - Database schemas và models.

### 4. 🔀 Flows (Luồng tương tác)
- [identity-authentication.md](./flows/identity-authentication.md) - Luồng đăng nhập và xác thực.
- [system-data-flow.md](./flows/system-data-flow.md) - Dòng chảy dữ liệu qua các services.

### 5. 📜 ADRs (Nhật ký Quyết định Kiến trúc)
- Xem thư mục: [`adrs/`](./adrs/)
