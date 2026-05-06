# Frontend Coding Guidelines — Runtime Roasters

Tài liệu này tổng hợp các quy tắc và pattern phát triển Frontend dựa trên các feedback và yêu cầu của dự án. Tất cả các thay đổi code sau này phải tuân thủ nghiêm ngặt các quy tắc này.

## 1. Kiến trúc Service (OOP Class)
Tất cả các service logic (Auth, Storage, API) phải được triển khai theo dạng **OOP Class**.
- Không sử dụng các hàm export rời rạc cho logic phức tạp.
- Export một instance duy nhất của class (Singleton pattern).

**Ví dụ:**
```typescript
class AuthService {
  async login(...) { ... }
}
export const authService = new AuthService();
```

## 2. Quản lý hằng số (No Magic Strings)
Tuyệt đối không sử dụng magic strings trong code. Tất cả các giá trị tĩnh phải được định nghĩa trong thư mục `src/constants/`.

- **ENV**: Biến môi trường (`constants/env.ts`).
- **API_ENDPOINTS**: Tất cả URL API (`constants/api.ts`).
- **STORAGE_KEYS**: Các key lưu trữ LocalStorage/SessionStorage (`constants/storage.ts`).
- **APP_ROUTES**: Các đường dẫn trang trong ứng dụng (`constants/routes.ts`).
- **AUTH_PARAMS**: Các tham số query liên quan đến Auth (`constants/auth.ts`).

## 4. Thư viện & State Management
- **API Client**: Luôn sử dụng `Axios`. Tránh dùng `fetch` trực tiếp trừ trường hợp cực kỳ đặc biệt.
- **Data Fetching**: Sử dụng `TanStack Query` (React Query) cho mọi tương tác với server để quản lý cache và loading state.
- **Styling**: Vanilla CSS hoặc TailwindCSS v4 (theo yêu cầu).

## 4. Utilities & Helpers
- **Môi trường**: Sử dụng util `isBrowser` thay cho `typeof window !== 'undefined'`.
- **URL Handling**: Sử dụng hàm `joinPaths` để nối các thành phần của URL, tránh lỗi thừa/thiếu dấu gạch chéo (`/`).
- **Type Safety**: Tránh sử dụng `any`. Luôn định nghĩa interface/type rõ ràng.

## 5. Aesthetics (Industrial Premium)
- Giao diện phải mang đậm phong cách **Industrial Premium**:
    - Sử dụng Glassmorphism (backdrop-blur).
    - Màu sắc hài hòa, độ tương phản cao nhưng sang trọng (slate, emerald, primary).
    - Hiệu ứng animation tinh tế với `framer-motion`.
    - Typography hiện đại (Google Fonts).

## 6. Lints & Types
- Code phải vượt qua các bước kiểm tra `npm run lint` và `npx tsc --noEmit`.
- Không bỏ qua (suppress) các lỗi lints bằng comment trừ khi có lý do bất khả kháng.

## 7. Auth Patterns (Seamless Identity)
Dự án sử dụng mô hình OIDC chuẩn nhưng tối ưu hóa trải nghiệm người dùng (Zero-Consent flow cho app nội bộ).

- **Auto-Accept Consent**: Trang `/consent` được triển khai dưới dạng **Server Component**.
    - Nó tự động kiểm tra `client_id` dựa trên danh sách `TRUSTED_CLIENTS`.
    - Thực hiện lệnh `acceptOAuth2ConsentRequest` trực tiếp từ server-side.
    - Người dùng không bao giờ nhìn thấy màn hình cấp quyền, mang lại cảm giác mượt mà như SSO.
- **Hydra Admin Access**: Luôn thực hiện các tác vụ Admin (Accept Login/Consent) từ server-side thông qua API routes hoặc Server Components để bảo mật Token.

---
> [!TIP]
> Sử dụng [Checklist kiểm chứng Identity](file:///Users/dungxbuif/workspace/RuntimeRoasters/docs/identity_test_checklist.md) để test luồng đăng nhập sau khi thay đổi code.

> [!IMPORTANT]
> Khi review code hoặc thực hiện task mới, Antigravity phải tự động kiểm tra xem code có vi phạm các quy tắc trên không (đặc biệt là magic strings và OOP service).
