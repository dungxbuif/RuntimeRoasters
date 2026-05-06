# Identity Flow Test Checklist (Sprint 2)

Tài liệu này hướng dẫn các bước kiểm chứng luồng xác thực OIDC (Ory Kratos & Hydra) sau khi đã tối ưu hóa kiến trúc Frontend và Server-side Consent.

---

## 🛠️ Bước 1: Chuẩn bị hạ tầng
- [ X ] **Docker Services**: Chạy lệnh khởi động hạ tầng:
  ```bash
  docker compose -f deployments/docker-compose.dev.yaml up -d
  ```
- [ X ] **Health Check**: Kiểm tra các container quan trọng phải ở trạng thái `Healthy`:
  - `rr-postgres`
  - `rr-kratos`
  - `rr-hydra`
  - `rr-identity` (Nginx Proxy)
- [ X ] **Clean State**: Xóa sạch Cookie và LocalStorage của domain `localhost:3000` trên trình duyệt.

---

## 🔐 Bước 2: Kiểm chứng luồng Đăng nhập (OIDC Flow)
- [ ] **Redirect Auth**: Truy cập `http://localhost:3000/`.
  - *Kết quả*: Phải tự động redirect sang trang Login `/login`.
- [ ] **Identity Input**: Nhập thông tin Admin:
  - **Email**: `admin@runtimeroasters.com`
  - **Password**: `AdminPassword123!`
- [ ] **Seamless Consent (Quan trọng)**: Quan sát quá trình sau khi bấm Login.
  - *Kết quả*: Bạn sẽ thấy URL nháy qua `/consent` và `/api/auth/callback` rất nhanh. 
  - *Yêu cầu*: **KHÔNG** được hiển thị giao diện nút "Allow/Deny". Người dùng phải vào thẳng Dashboard.
- [ ] **Session Persistence**: F5 lại trang Dashboard.
  - *Kết quả*: Phải giữ được trạng thái đăng nhập, không bị đá ra ngoài.

---

## 📊 Bước 3: Kiểm tra Token & Bảo mật
- [ ] **LocalStorage Check**: Mở DevTools (F12) > Application > Local Storage.
  - *Kết quả*: Phải thấy key `rr_access_token` chứa chuỗi JWT.
- [ ] **JWT Payload Verification**: Copy token dán vào [JWT.io](https://jwt.io).
  - [ ] `sub`: UUID của người dùng.
  - [ ] `role`: Phải là **`ADMIN`** (viết hoa).
  - [ ] `client_id`: Phải là `client-app`.
- [ ] **API Authorization**: Bấm nút **"Call Secured API"** trên Dashboard.
  - *Kết quả*: Kiểm tra Network tab, request gửi đi phải có Header `Authorization: Bearer <token>`.

---

## 🚪 Bước 4: Kiểm tra Logout
- [ ] **Logout Action**: Bấm nút Logout trên Dashboard.
  - *Kết quả*: LocalStorage `rr_access_token` bị xóa, redirect về `/login`.
- [ ] **Route Guard**: Thử truy cập lại `http://localhost:3000/` bằng tay.
  - *Kết quả*: Phải bị redirect về `/login`.

---

## 📝 Ghi chú cho Reviewer
- Toàn bộ logic Accept Consent đã được chuyển sang **Server Component** (`src/app/(auth)/consent/page.tsx`).
- Biến môi trường nhạy cảm được quản lý trong file `.env.local`.
- Danh sách Trusted Clients được cấu hình tại `src/constants/auth.ts`.
