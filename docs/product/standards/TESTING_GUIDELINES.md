# 📖 Runtime Roasters: Testing & Verification Guidelines

Tài liệu này quy định tiêu chuẩn viết và xác thực kiểm thử (Testing Standards) cho toàn bộ hệ thống, đảm bảo mọi thay đổi đều có bằng chứng (Evidence) rõ ràng.

---

## 🏗️ 1. The Enterprise Test Pyramid
Dự án áp dụng mô hình kiểm thử 3 lớp:

1.  **Lớp 1: E2E Verification (Playwright)**
    *   **Mục tiêu:** Xác thực luồng nghiệp vụ từ góc nhìn người dùng (UI ➔ Gateway ➔ Backend).
    *   **Vị trí:** `src/apps/client-app/e2e/*.spec.ts`
    *   **Yêu cầu:** Phải bao phủ các trạng thái Loading, Success, Error và các rào chắn quyền (Role Scoping).

2.  **Lớp 2: Data Integrity Audit (Go Integration)**
    *   **Mục tiêu:** Kiểm tra trực tiếp Database và luồng Kafka sau khi nghiệp vụ hoàn tất.
    *   **Vị trí:** `src/pkg/testing/integration/*_test.go`
    *   **Yêu cầu:** Tuyệt đối không hardcode DB credentials. Sử dụng `require` cho các lỗi chặn đứng và `assert` cho các kiểm tra giá trị.

3.  **Lớp 3: Business Logic (Unit Test)**
    *   **Mục tiêu:** Test các hàm xử lý tính toán (ví dụ: Roasting Loss, ID Formatting).
    *   **Vị trí:** Cùng thư mục với code nghiệp vụ (`*_test.go`).

---

## ✍️ 2. Hướng dẫn viết Test Case mới

### Quy tắc đặt tên (Naming Convention)
Mọi Test Case phải được gán một ID duy nhất khớp với **Master Test Tracking Matrix**:
*   `TC-[Flow ID].[Sub ID]` (Ví dụ: `TC-1.2`)

### Cấu trúc một E2E Test chuẩn:
```typescript
test('TC-X.Y: [Mô tả ngắn gọn]', async ({ page }) => {
  // 1. Arrange: Chuẩn bị môi trường/data
  // 2. Act: Thực hiện hành động trên UI
  // 3. Assert: Kiểm tra kết quả mong đợi
  await expect(page.locator('text=Success')).toBeVisible();
});
```

---

## ✅ 3. Quy trình Xác thực (Verification Workflow)

Mọi Task sau khi hoàn thành **BẮT BUỘC** phải đi qua các bước:

1.  **Cập nhật Tracking:** Điền kết quả vào `docs/product/standards/TEST_TRACKING.md`.
2.  **Chạy Automation:** 
    *   FE: `cd src/apps/client-app && npm run test:e2e`
    *   BE: `cd src && go test -v ./pkg/testing/integration/...`
3.  **Tạo Validation Report:** Sử dụng template `docs/templates/validation-report-v2.md`.

---

## 🛠️ 4. Tooling & Commands
*   **Debug UI:** `npx playwright test --debug`
*   **Reset Data:** `task env:reset` (Luôn chạy trước khi test luồng khởi tạo).
*   **Kafka Monitor:** `http://localhost:8090` (Để verify sự kiện thực tế).

---
*Tiêu chuẩn V1.0 | Áp dụng cho mọi Agent và Nhà phát triển.*
