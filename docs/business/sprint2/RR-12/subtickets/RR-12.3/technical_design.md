# RR-12.3: Gin Authorization Middleware

## 1. Mục tiêu (Goal)
Triển khai Middleware cho Gin framework để kiểm tra quyền truy cập cho các RESTful API endpoint.

## 2. Ngữ cảnh (Context)
- **Package:** `src/pkg/base/casbin`
- **File cần tạo:** `src/pkg/base/casbin/middleware_gin.go`
- **Phụ thuộc:** `pkg/base/identity`, `casbin.Enforcer`, Gin Gonic.

## 3. Các bước triển khai (Step-by-Step Implementation)
1. **Định nghĩa Middleware:** Tạo hàm `NewGinMiddleware` nhận `casbin.Enforcer` làm tham số.
2. **Lấy Identity:** Trích xuất `Claims` từ context của Gin (thường được lưu bởi một Auth Middleware trước đó).
3. **Xác định Hành động (act):** Gọi hàm helper `MapHTTPMethodToAction(c.Request.Method)` đã tạo ở bước RR-12.1.
4. **Xác định Resource (obj):** 
   - Sử dụng `c.FullPath()` hoặc `c.Request.URL.Path` làm Resource ID.
   - *Lưu ý:* Nếu có thể, hãy ánh xạ URL path sang gRPC Method tương đương để sử dụng chung chính sách trong `policy.csv`.
5. **Thực thi phân quyền:** Gọi `enforcer.Enforce(claims.Role, resource, action)`.
6. **Xử lý phản hồi:**
   - Nếu `false`: Sử dụng `c.AbortWithStatusJSON(403, ...)` để trả về lỗi Forbidden.
   - Nếu `true`: Gọi `c.Next()`.

## 4. Quy chuẩn tuân thủ (Patterns to Follow)
- **Surgical Update:** Middleware phải nhẹ và không gây block request quá lâu.
- **Consistency:** Đảm bảo cách ánh xạ Action giữa Gin và gRPC là thống nhất.

## 5. Xác minh (Verification)
- **Manual Test:** Sử dụng `curl` hoặc Postman gọi tới một endpoint REST với các token có role khác nhau (ví dụ: `farmer` vs `guest`).
- **Log:** Kiểm tra log để đảm bảo middleware được kích hoạt và đưa ra quyết định đúng.
