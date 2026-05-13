# 🤖 PROJECT AGENT CONTEXT: Runtime Roasters

Tài liệu này là "nguồn sự thật duy nhất" (Single Source of Truth) dành cho AI Agents khi làm việc trên dự án **Runtime Roasters**. Mọi hành động, thiết kế và review code PHẢI tuân thủ các chỉ dẫn dưới đây.

---

## ☕ 1. Tầm nhìn & Sứ mệnh (Project Mission)
**Runtime Roasters** là một nền tảng quản lý chuỗi cung ứng cà phê (Farm-to-Cup) mô phỏng production-grade.
- **Mục tiêu:** Showcase kiến trúc Microservices phức tạp trong hệ sinh thái Go.
- **Giá trị cốt lõi:** Minh bạch, Tin cậy, Quan sát được (Observability).

---

## 🏗️ 2. Trụ cột Kiến trúc (Architectural Pillars)

### 2.1 Microservices Monorepo
- **Cấu trúc:** `api/` (Proto), `pkg/` (Shared Libs), `src/apps/` (Services).
- **Giao tiếp:** gRPC (Nội bộ), REST/KrakenD (Bên ngoài), Kafka (Bất đồng bộ).

### 2.2 Clean Architecture (Mandatory)
Tuân thủ **ADR 0001** và triết lý **Go Clean Arch**:
- **Domain:** Pure Go, không import package ngoài Standard Library.
- **UseCase:** Định nghĩa logic nghiệp vụ. **QUAN TRỌNG:** Interface thuộc về Consumer (định nghĩa tại UseCase, thực thi tại Infra).
- **Infrastructure:** Chi tiết thực thi (DB, API, Broker).
- **Composition Root:** Manual DI tại `cmd/main.go` hoặc `internal/app/init.go`. **KHÔNG dùng DI Framework.**

### 2.3 Event-Driven & Consistency
- **Saga Pattern:** Điều phối đa bước (Choreography) với cơ chế Hoàn tác (Compensating Actions).
- **Transactional Outbox:** Đảm bảo tính nguyên tử giữa DB state và Kafka event.
- **Dual Idempotency:** Bảo vệ 2 lớp: Valkey (Sync/HTTP) + Inbox Pattern (Async/Kafka).
- **CQRS:** Tách biệt Postgres (Write) và Elasticsearch (Read/Traceability).

---

## 🛠️ 3. Tiêu chuẩn Kỹ thuật (Technical Standards)

- **Ngôn ngữ:** Go 1.25+ (Go Workspaces).
- **Error Handling:** Tuân thủ **RFC 7807/9457 (Problem Details)**.
- **Security:** Mô hình 3 lớp (API Gateway -> Casbin Service-level -> GormScoper Data-level).
- **Observability:** 
    - Luôn duy trì context propagation (Trace-ID).
    - Sử dụng `logger.FromContext(ctx)` để ghi log kèm Trace-ID.
- **Database:** PostgreSQL (GORM), Elasticsearch, Apache Cassandra, Valkey.

---

## 📋 4. Quy trình Phát triển (Workflow & Rules)

### 4.1 Quản lý Công việc (Git-as-Jira)
- Mỗi task nằm trong `docs/business/sprintX/RR-x/`.
- `ticket.md`: Yêu cầu nghiệp vụ (What/Why).
- `technical_design.md`: Kế hoạch thực thi kỹ thuật (How) - **PHẢI viết trước khi code.**

### 4.2 Documentation (Diátaxis Framework)
- `docs/architecture/README.md`: Kiến trúc tổng thể & Nguyên tắc cốt lõi (The WHY).
- `docs/engineering/README.md`: Cẩm nang kỹ thuật & Tra cứu (The WHAT/HOW).
- `docs/architecture/flows/`: Sơ đồ tương tác (The INTERACTIONS).
- `docs/architecture/adrs/`: Nhật ký quyết định kiến trúc.

### 4.3 Code Review (Tech Lead Role)
- Sử dụng skill tại `.antigravity/skills/tech-lead-reviewer/SKILL.md`.
- Ưu tiên tính Readability, Clean Arch, và Testability.

---

## 📂 5. Bản đồ Thư mục (Project Map)

- `/api/`: Chứa các file `.proto` định nghĩa contract gRPC.
- `/pkg/`: Thư viện dùng chung (Auth, Kafka, Logger, DB, Telemetry).
- `/src/apps/`: Các microservices độc lập (auth, farm, warehouse, etc.).
- `/deployments/`: Cấu hình Docker Compose, Infra seed scripts.
- `/docs/`: Toàn bộ tài liệu kiến trúc và nghiệp vụ.

---

## ⚠️ 6. Chỉ dẫn AI Đặc thù (AI-Specific Directives)

1. **Nghiêm cấm "Magic":** Luôn ưu tiên sự tường minh (Explicit over Implicit). Không dùng Reflection, Monkey Patching hay Hidden Logic.
2. **Test First:** Mọi bug fix phải đi kèm Test Case tái hiện lỗi. Mọi feature mới phải có Unit Test (dùng Mockery).
3. **Context is King:** Khi đọc code, hãy đọc file `GEMINI.md` hoặc `AGENT.md` trong thư mục tương ứng để hiểu bối cảnh local.
4. **Surgical Updates:** Khi sửa code, chỉ thay đổi những phần liên quan trực tiếp, không refactor lan man trừ khi được yêu cầu.

---
*Cập nhật lần cuối: 2026-05-13 bởi TechLead Agent*
