# 🧪 Runtime Roasters: Master System Validation Manual

Tài liệu này là **Cẩm nang Xác thực Toàn diện** (Master Validation Manual) cho hệ thống Runtime Roasters. Tài liệu được thiết kế theo cấu trúc mô-đun (Modular), tuân thủ nghiêm ngặt [Testing & Verification Guidelines](./TESTING_GUIDELINES.md).

---

## 🛠️ Global Configuration & Resources
*Sử dụng các thông tin này làm hằng số (Constants) xuyên suốt mọi giai đoạn xác thực.*

| Resource | Value / Endpoint | Purpose |
| :--- | :--- | :--- |
| **Admin Email** | `admin@runtimeroasters.com` | Tài khoản quản trị tối cao |
| **Default Password**| `Hello@123` | Mật khẩu cho mọi tài khoản Seed/Demo |
| **Frontend URL** | `http://localhost:3000` | Giao diện người dùng |
| **API Gateway** | `http://localhost:8081` | KrakenD Entrypoint |
| **Kafka UI** | `http://localhost:8090` | Giám sát luồng sự kiện (Events) |
| **SigNoz (OTel)** | `http://localhost:3301` | Tracing & Metrics |
| **Kibana (Search)** | `http://localhost:5601` | Kiểm tra CQRS Read Model |

---

## 🏁 Flow 1.0: Fresh System Setup & ADMIN Bootstrap
*Đây là Flow nền tảng, đảm bảo "Luật chơi" và "Thế giới dữ liệu" được khởi tạo chính xác.*

### Bước 1: Khởi động & Dọn sạch môi trường
*   **Hành động:** Chạy lệnh `task env:reset`
*   **Xác thực hạ tầng:**
    *   Truy cập **System Explorer** (`/dashboard/explorer`): Phải thấy Swagger UI và KrakenD Gateway phản hồi 200.
    *   Truy cập **Kafka UI** (`:8090`): Phải thấy các topic `auth.*` và `logistics.*`.
*   **Trạng thái DB:**
    *   `identity_db`: Phải có 1 identity duy nhất (`admin@runtimeroasters.com`).
    *   `farm_db`, `retail_db`: Các bảng nghiệp vụ phải trống.

### Bước 2: Đăng nhập & Kích hoạt Seeding (ADMIN Only)
*   **Hành động:**
    1.  Mở `localhost:3000`, đăng nhập bằng `admin@runtimeroasters.com` / `Hello@123`.
    2.  **UI Requirement:** Sau khi đăng nhập, hệ thống điều hướng về **Manage Users**.
    3.  **Bootstrap Trigger:** Một **System Bootstrap Modal** phải tự động hiện lên nếu dữ liệu chưa được nạp.
    4.  *Ghi chú:* Nếu modal không hiện, có thể truy cập `http://localhost:3000/dashboard?bootstrap=true` để cưỡng bức hiển thị.
*   **Expectation:** Bấm **"Initialize Data"**. Modal báo `Success` và tự đóng sau 2 giây.

### Bước 3: Xác thực Dữ liệu Master (Dữ liệu nền tảng)
Kiểm tra tính hiện diện của dữ liệu tại các trang quản trị:

#### A. Manage Users (Identity List)
*   **UI Requirement:** Phải có bộ lọc **"Filter Role"** và ô **"Search"**.
*   **Xác thực:** Lọc theo Role `Logistics Driver` ➔ Phải thấy 3 tài xế: `driver@...`, `driver.beta@...`, `driver.hcm@...`.
*   **Status:** Tất cả user seed phải có trạng thái **"Active"**.

#### B. Manage Farms (Farm Registry)
*   **Xác thực:** Danh sách phải đủ **6 Farms**.
*   **Data Check:**
    *   `FARM-CAUDAT-001` (K'Ho Coffee Farm)
    *   `FARM-CAUDAT-002` (Cau Dat Arabica)
    *   `FARM-BMT-001` (Aeroco Coffee)
    *   ... (và 3 farm khác).

#### C. Intelligence Hub (Global KPI)
*   **Hành động:** Truy cập menu **Intelligence Hub** (biểu đồ analytics).
*   **UI Requirement:** Phải thấy 4 thẻ KPI chính:
    1.  **Pilot Farms:** Hiển thị số 6.
    2.  **Active Personnel:** Hiển thị tổng số User đã seed (khoảng 10-15).
    3.  **Retail Stores:** Hiển thị số 5.
    4.  **Fleet Status:** Hiển thị "3 Active".
*   **System Integrity:** Các dịch vụ `Identity`, `Kafka`, `Traceability` phải báo trạng thái **"Healthy"** (màu xanh).

---

## 🖥️ UI Visibility Rules (Dành cho ADMIN)
*Giữ giao diện tối giản để tập trung vào quản trị.*

| Màn hình | Item nên hiện | Item nên ẩn |
| :--- | :--- | :--- |
| **Sidebar** | Users, Farms, Account, Intelligence Hub, Trace, Explorer | Các dashboard vận hành lẻ (Harvest, Logistics Monitor) |
| **Intelligence Hub** | KPI Cards, Chain Integrity Status, Anomaly Watch | Các biểu đồ chi tiết của từng lô hàng |
| **Manage Users** | Search, Role Filter, Identity Table, Create Button | Plaintext Passwords |

---

## 🛰️ Monitoring Endpoints (Hậu kiểm kỹ thuật)
*Xác nhận luồng dữ liệu ngầm.*

1.  **Kafka Flow:** Topic `logistics.locations.updated` ➔ Phải có message cho 14 địa điểm.
2.  **Identity Admin:** `http://localhost:4434/admin/identities` ➔ Trả về JSON danh sách toàn bộ User.
3.  **Search Read Model:** `http://localhost:9200/trace_events/_search` ➔ Phải bắt đầu có dữ liệu index.

---
*Manual Version: 1.1 (Synced with UI) | Owner: Architect/ADMIN*
