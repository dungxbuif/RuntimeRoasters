# 📊 Master Test Tracking Matrix

Bảng này theo dõi tiến độ và kết quả xác thực của toàn bộ hệ thống Runtime Roasters theo `VERIFICATION_PROTOCOL.md`.

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

## 🚜 Flow 2.0: Upstream Supply Chain (Farm ➔ Warehouse)
**Mục tiêu:** Validate harvest logging, dispatch logistics, and physical intake discipline.

| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-2.1` | **E2E** | Harvest Declaration (ID Format & Outbox) | ⚪ IDLE | - |
| `TC-2.2` | **E2E** | Warehouse Dispatch & Driver Assignment | ⚪ IDLE | - |
| `TC-2.3` | **E2E** | Driver GPS Simulation (Heartbeat & Offline Buffer) | ⚪ IDLE | - |
| `TC-2.4` | **E2E** | Mandatory Return & Warehouse Intake Creation | ⚪ IDLE | - |

---

## 🏭 Flow 3.0: Downstream SAGA (Warehouse ➔ Retail)
**Mục tiêu:** Validate value transformation, distributed SAGA coordination, and fulfillment.

| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-3.1` | **E2E** | Roasting Anomaly Note Enforcement (>5% loss) | ⚪ IDLE | - |
| `TC-3.2` | **E2E** | Retail Order Creation & SAGA Initiation | ⚪ IDLE | - |
| `TC-3.3` | **Integration**| Stripe Webhook Pass ➔ payment.completed & stock reserved | ⚪ IDLE | - |
| `TC-3.4` | **Integration**| Stripe Webhook Fail ➔ payment.failed & order rejected | ⚪ IDLE | - |
| `TC-3.5` | **Integration**| SAGA Rollback (Stock Reservation Fail ➔ Refund) | ⚪ IDLE | - |
| `TC-3.6` | **E2E** | Retail Delivery (GPS Simulation & Return to Base) | ⚪ IDLE | - |

---

## 🔬 Flow 4.0: Observability & Trace Integrity
**Mục tiêu:** Validate distributed tracing context across polyglot microservices.

| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-4.1` | **Infra** | W3C Trace ID contiguous propagation (SigNoz UI) | ⚪ IDLE | - |
| `TC-4.2` | **Infra** | Elasticsearch `trace_events` index contains correct trace metadata | ⚪ IDLE | - |

---

## ⚡ Flow 5.0: Real-time Awareness (WebSocket)
**Mục tiêu:** Validate low-latency updates and Role/Scope-based message filtering.

| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-5.1` | **API** | Public Topology Config Pull (No-Auth HTTP 200) | ⚪ IDLE | - |
| `TC-5.2` | **API** | Public Topology History Replay via REST | ⚪ IDLE | - |
| `TC-5.3` | **E2E** | WebSocket Push topology fanout on Kafka events | ⚪ IDLE | - |
| `TC-5.4` | **Security** | Private Stream Role Scoping (Cross-Manager isolation) | ⚪ IDLE | - |
| `TC-5.5` | **Security** | Internal Socket Publish API (AuthZ Fail Closed checks) | ⚪ IDLE | - |

---

## 🌐 Flow 6.0: Public Transparency Show
**Mục tiêu:** Validate end-consumer QR tracking and strict PII sanitization.

| Test ID | Type | Description | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| `TC-6.1` | **E2E** | Public Trace UI renders full Farm-to-Cup journey (<3s) | ⚪ IDLE | - |
| `TC-6.2` | **Security** | PII & Financial Data Sanitization in public payloads | ⚪ IDLE | - |

---

## 🛠️ Automation Status
- **E2E (Playwright):** `src/apps/client-app/e2e/flow-1.0-bootstrap.spec.ts` (Implementing...)
- **Backend (Go/SQL):** `src/pkg/testing/integration/bootstrap_test.go` (Planned)

*Ghi chú: Cập nhật Status sau mỗi phiên chạy Automation hoặc xác nhận Manual.*
