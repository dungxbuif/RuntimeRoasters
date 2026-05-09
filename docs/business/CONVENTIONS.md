# 📋 Quy chuẩn Quản lý Công việc (Jira-in-Markdown)

Tài liệu này quy định cấu trúc và quy trình quản lý các đầu việc (tickets) trong dự án RuntimeRoasters thông qua Git/Markdown, mô phỏng trải nghiệm Jira để đảm bảo sự tách bạch giữa Nghiệp vụ và Kỹ thuật.

---

## 🏗️ 1. Cấu trúc Thư mục (Folder Structure)

Mỗi đầu việc (Ticket) bắt buộc phải nằm trong một thư mục riêng biệt.

### A. Đối với Ticket đơn (Single Ticket)
```text
docs/business/sprintX/RR-x/
├── ticket.md             # Tài liệu BA (What/Why)
└── technical_design.md    # Tài liệu Kỹ thuật (How)
```

### B. Đối với Ticket lớn (Epic/Complex Ticket)
Nếu một Ticket quá lớn cần chia nhỏ, sử dụng thư mục `subtickets`:
```text
docs/business/sprintX/RR-x/
├── ticket.md             # BA Epic level
├── technical_design.md    # Tech Lead/Architect Strategy
└── subtickets/
    └── RR-x.y/           # Thư mục con cho từng sub-ticket
        ├── ticket.md     # BA Sub-ticket level
        └── technical_design.md # Dev Implementation Plan
```

---

## 🎭 2. Phân định Vai trò & Nội dung

### 📝 ticket.md (Vai trò: BA / Product Owner)
- **Ngôn ngữ:** Nghiệp vụ, phi kỹ thuật.
- **Nội dung:**
    - **User Story:** "Dưới vai trò là... tôi muốn... để..."
    - **Business Value:** Tại sao tính năng này quan trọng?
    - **Acceptance Criteria (AC):** Các kịch bản kiểm thử nghiệp vụ (Scenario-based).
    - **Priority & Type:** Độ ưu tiên và phân loại (Feature, Bug, Infra).

### 🛠️ technical_design.md (Vai trò: Developer / Tech Lead)
- **Ngôn ngữ:** Kỹ thuật chuyên sâu.
- **Nội dung:**
    - **Strategy:** Cách tiếp cận (Chọn library nào? Tại sao?).
    - **Implementation Steps:** Các bước code cụ thể (Copy boilerplate, setup DI, SQL Schema).
    - **Verification:** Cách thức kiểm tra kỹ thuật (Unit test, Curl command, DB check).
    - **Trade-offs:** Các đánh giá về hiệu năng hoặc giới hạn kỹ thuật.

---

## 🔄 3. Quy trình Vòng đời Ticket (Lifecycle)

1.  **Giai đoạn Khởi tạo (Refinement):** BA viết `ticket.md` để mô tả yêu cầu.
2.  **Giai đoạn Thiết kế (Technical Planning):** Developer đọc `ticket.md`, phân tích và viết `technical_design.md`. 
    - *Mục đích:* Chốt phạm vi (Scope) và giải pháp trước khi code.
3.  **Giai đoạn Thực hiện (Implementation):** Developer tiến hành code dựa trên `technical_design.md`.
4.  **Giai đoạn Hoàn tất (Closing):** 
    - (Optional) Cập nhật lại các thay đổi thực tế vào `technical_design.md` nếu có sai lệch so với kế hoạch ban đầu.
    - Chuyển trạng thái trong Kanban board của Sprint (`main.md`).

---

## 📏 4. Quy tắc Chung (General Rules)

1.  **Ngôn ngữ:** Sử dụng tiếng Việt chuyên nghiệp cho toàn bộ tài liệu.
2.  **Liên kết:** Sử dụng đường dẫn tương đối (Relative path) để đảm bảo các liên kết không bị hỏng khi di chuyển thư mục.
3.  **Kích thước Ticket:** 
    - Nếu ticket đủ nhỏ: Viết trực tiếp `technical_design.md` bên dưới folder task.
    - Nếu ticket to: Bắt buộc tạo thư mục `subtickets` để quản lý.
4.  **Tính nhất quán:** Không tạo file Markdown lẻ tẻ bên ngoài các thư mục task (trừ các file index như `main.md` hay `sprint-planning.md`).
