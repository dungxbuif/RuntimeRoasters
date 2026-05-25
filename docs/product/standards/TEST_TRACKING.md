# 📊 Master Test Tracking Matrix

Bảng này theo dõi tiến độ và kết quả xác thực của toàn bộ hệ thống Runtime Roasters.

---

## 🏁 Flow 1.0: Fresh System Setup & ADMIN Bootstrap
**Mục tiêu:** Verify foundational integrity and master data provisioning.

| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-1.1` | **Infra** | `task env:reset` - Infrastructure Ready (Kafka/DB) | 🟢 PASS | Logs |
| `TC-1.2` | **E2E** | Admin Login & Bootstrap Modal Detection | 🟡 PENDING | Automation |
| `TC-1.3` | **E2E** | Manual Seeding Trigger (Initialize Data) | 🟡 PENDING | Automation |
| `TC-1.4` | **Integration**| Identity Sync (15+ Users created in Kratos) | 🟡 PENDING | API Audit |
| `TC-1.5` | **Integration**| Farm Sync (6 Farms created in Farm DB) | 🟢 PASS | `bootstrap_test.go` |
| `TC-1.6` | **Integration**| Logistics Sync (14 Locations & Fleet in Kafka) | 🟢 PASS | `bootstrap_test.go` |
| `TC-1.7` | **UI/UX** | Admin Visibility (KPIs & Health Nominal) | 🟡 PENDING | Screenshots |

---

## 🚜 Flow 2.0: Upstream Supply Chain (Planned)
| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-2.1` | **E2E** | Harvest Declaration (ID Format & Outbox) | ⚪ IDLE | - |
| `TC-2.2` | **E2E** | Driver Dispatch & GPS Simulation | ⚪ IDLE | - |

---

## 🛠️ Automation Status
- **E2E (Playwright):** `src/apps/client-app/e2e/flow-1.0-bootstrap.spec.ts` (Implementing...)
- **Backend (Go/SQL):** `src/pkg/testing/integration/bootstrap_test.go` (Planned)

*Ghi chú: Cập nhật Status sau mỗi phiên chạy Automation.*
