# Hướng dẫn Kỹ thuật: OpenTelemetry & Tracing Flow

Tài liệu này giải thích cách Trace ID được sinh ra, truyền tải và hiển thị trong hệ thống Runtime Roasters (Sprint 1).

## 1. Bản đồ luồng đi (Trace Flow)
Khi bạn nhấn "Trigger API Call" trên Client App, một Trace ID duy nhất được luân chuyển như sau:

1.  **Client App**: Gửi request HTTP tới `localhost:8081`.
2.  **KrakenD (Gateway)**: 
    *   **Khởi tạo**: Sinh ra `Trace ID` đầu tiên.
    *   **Propagate**: Chèn Trace ID vào header `traceparent`.
    *   **File**: `deployments/krakend/krakend.json` (`telemetry/opentelemetry`).
3.  **Demo Service (HTTP Layer)**:
    *   **Nhận diện**: Middleware `otelgin` đọc header và nhận diện Trace ID cũ.
    *   **Span 1**: Tạo span cho request HTTP.
    *   **File**: `src/pkg/base/app.go` (hàm `NewApp`).
4.  **Demo Service (Internal Bridge)**:
    *   **Chuyển đổi**: `grpc-gateway` gọi gRPC nội bộ. 
    *   **Span 2**: `otelgrpc.NewClientHandler` tạo span cho bước chuyển đổi này.
5.  **Demo Service (gRPC Layer)**:
    *   **Xử lý**: `otelgrpc.NewServerHandler` nhận request gRPC và tạo span cuối cùng cho logic nghiệp vụ.
    *   **File**: `src/pkg/base/app.go` (hàm `NewApp` & `RegisterGateway`).

## 2. SigNoz thu nhận và hiển thị
*   **Collector**: SigNoz OTel Collector lắng nghe tại cổng `4317` (gRPC OTLP) và `4318` (HTTP OTLP).
*   **Storage**: Các service "đẩy" (push) dữ liệu span về collector, collector ghi vào ClickHouse.
*   **Visualization**: SigNoz UI tại `localhost:3301` nhóm các spans có cùng `Trace ID` thành waterfall/timeline.

## 3. Cách mở rộng Tracing
Để thêm một span mới (ví dụ khi gọi Database), bạn chỉ cần:
1. Đảm bảo biến `ctx` (context) được truyền từ layer này sang layer khác.
2. Sử dụng thư viện `otel.Tracer` để bắt đầu một span mới từ `ctx`.

---
*Tài liệu được khởi tạo sau khi xác thực thành công 8 spans trong Sprint 1.*
