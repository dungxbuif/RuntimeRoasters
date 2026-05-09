# RR-12.1: Thiết lập nền tảng Casbin (Foundation Setup)

## 1. Mục tiêu (Goal)
Thiết lập cơ sở hạ tầng cho việc phân quyền bằng Casbin, bao gồm việc khởi tạo Enforcer và các hàm tiện ích để ánh xạ (mapping) hành động.

## 2. Ngữ cảnh (Context)
- **Package:** `src/pkg/base/casbin`
- **Thư viện:** `github.com/casbin/casbin/v2`
- **File cần tạo:** `src/pkg/base/casbin/enforcer.go`

## 3. Các bước triển khai (Step-by-Step Implementation)
1. **Định nghĩa cấu trúc Config:** Tạo struct để chứa đường dẫn tới file `model.conf` và `policy.csv`.
2. **Khởi tạo Enforcer:**
   - Sử dụng `casbin.NewEnforcer(modelPath, policyPath)` để tạo một instance.
   - Bọc instance này trong một cấu trúc hoặc cung cấp hàm Factory để có thể sử dụng dưới dạng Singleton hoặc inject qua Wire.
3. **Triển khai hàm LoadPolicy:** Đảm bảo Enforcer có thể nạp lại chính sách từ file khi cần thiết.
4. **Hàm helper ánh xạ HTTP Method:**
   - Viết hàm `MapHTTPMethodToAction(method string) string`.
   - Quy tắc: `GET` -> `read`, `POST`/`PUT`/`PATCH` -> `write`, `DELETE` -> `delete`.
5. **Đăng ký hàm custom (tùy chọn):** Nếu cần dùng các hàm như `keyMatch` trong matcher của Casbin.

## 4. Quy chuẩn tuân thủ (Patterns to Follow)
- **Kiến trúc:** Tuân thủ Clean Architecture, giữ logic phân quyền độc lập với business logic.
- **Naming:** Sử dụng các danh từ chuẩn (read, write, delete) cho hành động (Actions).

## 5. Xác minh (Verification)
- **Unit Test:** Viết test kiểm tra hàm ánh xạ HTTP Method.
- **Integration Test:** Tạo một Enforcer tạm thời với file model/policy giả lập để kiểm tra hàm `Enforce`.
