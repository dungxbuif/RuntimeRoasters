# RR-11.5: Verification & E2E Testing

## 🎯 Goal
Xác nhận toàn bộ luồng bảo mật hoạt động đúng như mong đợi trên môi trường thực tế.

## 📋 Tasks
- [ ] Test HTTP: Dùng `curl` gửi request kèm valid token ➔ Mong đợi 200 OK.
- [ ] Test HTTP: Dùng `curl` gửi request kèm expired token ➔ Mong đợi 401 Unauthorized.
- [ ] Test gRPC: Kiểm tra việc truyền identity qua Metadata.
- [ ] Kiểm tra SigNoz: Đảm bảo `user.id` xuất hiện trong Span attributes của các request đã authenticated.

## 🔍 Definition of Done
- Dashboard demo hiển thị đúng dữ liệu khi đã login.
- Logs hiển thị trace đúng kèm identity thông tin.
