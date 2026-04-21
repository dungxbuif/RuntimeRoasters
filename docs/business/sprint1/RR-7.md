# [RR-7] Client App — Control Plane

- **Summary:** Xây dựng Control Plane dashboard trong `apps/client-app/` cho admin giám sát toàn bộ hệ thống.
- **Priority:** `HIGH`
- **Depends on:** RR-5, RR-6

---

## User Story

> As an admin, I want a control plane dashboard showing live service health and an embedded observability UI, so that I can monitor the system from a single interface without needing to access internal ports directly.

---

## Acceptance Criteria

### Scenario 1: Admin auth gate
- **Given:** Client app đang chạy tại `localhost:3000`.
- **When:** Tôi truy cập `localhost:3000/control/services` mà chưa đăng nhập.
- **Then:** Tôi bị redirect về `/login` — không thể truy cập bất kỳ trang `/control/*` nào.

### Scenario 2: Service health dashboard
- **Given:** Tôi đã đăng nhập với quyền admin.
- **When:** Tôi truy cập `/control/services`.
- **Then:** Trang hiển thị health cards cho từng service (demo-service, krakend, ...) với trạng thái Running/Degraded và latency từ health check endpoint.

### Scenario 3: API Explorer load swagger
- **Given:** Demo Service đang chạy.
- **When:** Tôi truy cập `/control/api-explorer`.
- **Then:** Swagger UI hiển thị đúng nội dung từ `demo.swagger.json`. Tôi có thể Execute một request qua KrakenD và nhận response.

### Scenario 4: SigNoz accessible qua proxy
- **Given:** SigNoz đang chạy (port 3301 internal only).
- **When:** Tôi truy cập `/control/observability`.
- **Then:** SigNoz UI hiển thị đúng trong iframe. Tôi có thể xem traces, metrics, logs. Port 3301 vẫn không accessible trực tiếp từ browser.

### Scenario 5: SigNoz proxy hoạt động qua env var
- **Given:** `SIGNOZ_INTERNAL_URL=http://signoz:3301` được set trong docker-compose.
- **When:** Client app xử lý request `/signoz/*`.
- **Then:** Request được forward đúng đến SigNoz internal URL — không hardcode hostname.
