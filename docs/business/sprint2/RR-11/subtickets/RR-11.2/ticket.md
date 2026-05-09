# RR-11.2: JWKS Provider with Cache

## 🎯 Goal
Xây dựng cơ chế fetch Public Keys từ Identity Server một cách an toàn và hiệu năng cao.

## 📋 Tasks
- [x] Tạo file `pkg/base/auth/jwks.go`.
- [x] Triển khai `Fetcher` gửi kèm header `X-Internal-Secret`.
- [x] Triển khai `jwksCache` với cơ chế **Stale-While-Revalidate**:
    - Trả về key cũ ngay lập tức nếu đã hết hạn (TTL).
    - Khởi chạy goroutine ngầm để update key mới mà không làm block request hiện tại.
- [x] Cấu hình TTL thông qua biến môi trường (default 5m).

## 🔍 Definition of Done
- Cache hoạt động đúng: không block request khi refresh.
- Có log cảnh báo nếu không fetch được JWKS (ví dụ: sai secret).
---
# Chi tiết triển khai JWT Validation (RR-11)

Tài liệu này giải thích sâu về mặt kỹ thuật, các quyết định thiết kế và giải thuật được sử dụng trong hệ thống xác thực JWT của Runtime Roasters.

---

## 1. Cấu trúc dữ liệu (Struct Design)

### `JWKSCache` Struct
Cấu trúc này được thiết kế để tối ưu cho môi trường High-Concurrency (truy cập đồng thời cao).

```go
type JWKSCache struct {
    mu         sync.RWMutex              // Điều phối truy cập đồng thời
    keys       map[string]*rsa.PublicKey // Bộ nhớ đệm lưu trữ khóa công khai
    fetchedAt  time.Time                 // Thời điểm lấy dữ liệu cuối cùng
    refreshing bool                      // Cờ ngăn chặn Thundering Herd
    // ... config fields
}
```

*   **`sync.RWMutex`**: Sử dụng `RWMutex` thay vì `Mutex` thường vì số lượng request **Đọc** (xác thực token) nhiều hơn rất nhiều so với request **Ghi** (cập nhật khóa). `RLock` cho phép nhiều goroutine đọc cùng lúc mà không gây nghẽn.
*   **`map[string]*rsa.PublicKey`**: Lưu trữ trực tiếp đối tượng đã được parse giúp việc xác thực tốn cực ít tài nguyên (vài micro-giây).
*   **`refreshing bool`**: Đóng vai trò như một "lồng khóa" (locking mechanism) để đảm bảo tại một thời điểm chỉ có duy nhất một worker thực hiện gọi network lấy dữ liệu, tránh lãng phí tài nguyên và gây áp lực lên Identity Server.

---

## 2. Chiến lược Caching: Stale-While-Revalidate (SWR)

Đây là kỹ thuật cốt lõi giúp hệ thống đạt hiệu năng cao và tính sẵn sàng (Availability).

*   **Vấn đề:** Cache truyền thống thường gây ra hiện tượng "Latency Spike" khi cache hết hạn vì request phải chờ đợi fetch dữ liệu mới.
*   **Giải pháp (SWR):** 
    1. Khi request tới và cache đã hết hạn (TTL), nếu trong cache **vẫn còn dữ liệu cũ (stale)**, hệ thống sẽ trả về dữ liệu đó ngay lập tức.
    2. Một goroutine chạy ngầm được kích hoạt để thực hiện việc cập nhật (revalidate) cache.
*   **Kết quả:** 
    - Latency (P99) cực kỳ ổn định vì không bao giờ bị chặn bởi network call.
    - Khả năng chịu lỗi cao: Nếu Identity Server gặp sự cố tạm thời, hệ thống vẫn có thể hoạt động dựa trên dữ liệu stale.

---

## 3. Giải thuật giải mã RSA (`decodeRSA`)

Các khóa trong chuẩn JWKS bao gồm hai thành phần: `n` (Modulus) và `e` (Exponent) ở dạng Base64URL.

*   **Sử dụng `math/big`**: Thành phần `n` của RSA rất lớn (thường 2048-4096 bit), do đó bắt buộc phải dùng `big.Int` để tính toán.
*   **Chuyển đổi thành phần `e`**: Thuật toán quét qua mảng byte của `e` để chuyển đổi về kiểu `int` chuẩn của Go, đảm bảo tương thích hoàn toàn với thư viện `crypto/rsa`.

---

## 4. Nguyên tắc thiết kế (Design Principles)

### Functional Options Pattern (Interceptor)
Sử dụng `opts ...InterceptorOption` trong gRPC Interceptor.
*   **Lợi ích:** Tuân thủ nguyên tắc **Open/Closed** (SOLID). Cho phép mở rộng cấu hình (như thêm logging, whitelist, timeout) mà không làm thay đổi chữ ký hàm (breaking changes) đối với các service đang sử dụng.

### Dependency Inversion (KeyProvider Interface)
Tách biệt logic xác thực JWT khỏi logic lấy khóa.
*   **Unit Testing:** Dễ dàng tạo Mock để kiểm thử logic auth mà không cần mạng.
*   **Linh hoạt:** Có thể thay đổi nguồn cấp khóa (từ JWKS sang Local File, Vault, v.v.) mà không cần sửa đổi mã nguồn của Middleware hay JWT Parser.

---

## 5. Luồng xử lý (Call Flow)

1.  **Request tới (HTTP/gRPC)**.
2.  **Trích xuất Token** từ Header/Metadata.
3.  **Gọi `VerifyAndParseJWT`**:
    - Gọi `KeyProvider.GetPublicKey(kid)`.
    - `JWKSCache` kiểm tra bộ nhớ đệm.
    - Nếu hết hạn ➔ Trả key cũ + Chạy goroutine `refresh()`.
4.  **Xác thực chữ ký RS256** bằng Public Key nhận được.
5.  **Trích xuất Claims** và inject vào `Context`.
6.  **Ghi Log/Tracing** (user.id) vào OpenTelemetry.
7.  **Chuyển tiếp** request tới Business Logic.
