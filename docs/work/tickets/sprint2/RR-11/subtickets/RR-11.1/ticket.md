# RR-11.1: Identity Core & Dependencies

## 🎯 Goal
Thiết lập nền tảng kiểu dữ liệu Identity chung cho toàn hệ thống và cài đặt các thư viện cần thiết.

## 📋 Tasks
- [ ] Chạy `go get github.com/golang-jwt/jwt/v5` tại thư mục `src`.
- [ ] Tạo package `pkg/base/identity`.
- [ ] Định nghĩa `type Claims struct` chứa các trường: `Subject`, `Role`, `OrgID`, `JTI`.
- [ ] Triển khai `InjectContext(ctx, claims)` và `FromContext(ctx)` sử dụng private context key để đảm bảo type-safety.

## 🔍 Definition of Done
- `go.mod` đã cập nhật thư viện jwt v5.
- Code `pkg/base/identity` pass unit test (nếu có) và không có magic strings.
