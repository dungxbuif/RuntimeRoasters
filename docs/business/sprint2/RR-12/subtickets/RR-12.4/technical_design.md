# RR-12.4: Cấu hình Chính sách (Policy) & Triển khai

## 1. Mục tiêu (Goal)
Định nghĩa mô hình phân quyền (Model) và các chính sách (Policies) thực tế cho hệ thống, đồng thời đảm bảo chúng được nạp đúng vào môi trường chạy (Docker).

## 2. Ngữ cảnh (Context)
- **Thư mục:** `deployments/casbin/`
- **File cần tạo:** `model.conf`, `policy.csv`
- **File cần sửa:** `deployments/docker-compose.dev.yaml`

## 3. Các bước triển khai (Step-by-Step Implementation)
1. **Tạo `model.conf`:**
   - Sử dụng mô hình RBAC với phân cấp (`role_definition` dùng `g = _, _`).
   - Định nghĩa `matchers` hỗ trợ `g(r.sub, p.sub)` và kiểm tra `obj`, `act`.
2. **Tạo `policy.csv`:**
   - Định nghĩa phân cấp: `g, admin, farmer` (Admin kế thừa mọi quyền của Farmer).
   - Định nghĩa quyền cơ bản cho `farmer`:
     - `p, farmer, /farm.v1.FarmService/CreateFarm, write`
     - `p, farmer, /farm.v1.FarmService/ListFarms, read`
3. **Cập nhật Docker Compose:**
   - Mount thư mục `deployments/casbin/` vào một đường dẫn cố định trong container (ví dụ: `/app/config/casbin/`).
   - Đảm bảo các service như `demo-service` có biến môi trường trỏ tới các file này.

## 4. Quy chuẩn tuân thủ (Patterns to Follow)
- **RBAC Hierarchy:** Luôn sử dụng kế thừa để giảm thiểu việc lặp lại các dòng chính sách (DRY - Don't Repeat Yourself).
- **Security:** File `policy.csv` trong môi trường Dev có thể commit, nhưng trong Production cần được thay thế bằng DB Adapter (Postgres).

## 5. Xác minh (Verification)
- **Container Check:** Chạy `docker exec -it <container_id> ls /app/config/casbin/` để kiểm tra file đã tồn tại.
- **Casbin Editor:** Dùng [Casbin Online Editor](https://casbin.org/editor/) để dán nội dung file model và policy vào kiểm tra logic matcher trước khi deploy.
