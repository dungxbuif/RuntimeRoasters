# RR-12.2: gRPC Authorization Interceptor

## 1. Mục tiêu (Goal)
Xây dựng một gRPC Interceptor để tự động kiểm tra quyền truy cập (Authorization) cho mọi request gRPC dựa trên thông tin định danh (Identity) trong context.

## 2. Ngữ cảnh (Context)
- **Package:** `src/pkg/base/casbin`
- **File cần tạo:** `src/pkg/base/casbin/interceptor_grpc.go`
- **Phụ thuộc:** `pkg/base/identity` (để lấy thông tin User), `casbin.Enforcer`.

## 3. Các bước triển khai (Step-by-Step Implementation)
1. **Định nghĩa Interceptor:** Tạo hàm `UnaryServerInterceptor` nhận `casbin.Enforcer` làm tham số đầu vào.
2. **Trích xuất thông tin định danh:** Sử dụng `identity.FromContext(ctx)` để lấy `Claims` (đã được nạp từ Interceptor Authentication ở RR-11).
3. **Xác định Resource (obj):** Sử dụng `info.FullMethod` (ví dụ: `/farm.v1.FarmService/CreateFarm`) làm Resource ID.
4. **Xác định Hành động (act):** 
   - Phân tích `info.FullMethod`. Nếu chứa tiền tố `Get` hoặc `List` -> hành động là `read`.
   - Các trường hợp khác mặc định là `write` hoặc `delete`.
5. **Kiểm tra quyền:** Gọi `enforcer.Enforce(claims.Role, fullMethod, action)`.
6. **Xử lý kết quả:**
   - Nếu trả về `true`: Cho phép tiếp tục (`handler(ctx, req)`).
   - Nếu trả về `false`: Trả về lỗi `codes.PermissionDenied` (gRPC error code).

## 4. Quy chuẩn tuân thủ (Patterns to Follow)
- **Resource ID:** Phải sử dụng `info.FullMethod` để đồng bộ với định nghĩa trong `policy.csv`.
- **Error Handling:** Trả về lỗi gRPC chuẩn để các client/gateway có thể xử lý đúng.

## 5. Xác minh (Verification)
- **Unit Test:** Sử dụng mock context chứa Identity và mock Enforcer để kiểm tra logic interceptor.
- **Manual Test:** Kiểm tra log của service khi gọi một hàm mà role hiện tại không có quyền.
