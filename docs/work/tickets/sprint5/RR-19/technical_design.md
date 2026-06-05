# Technical Design: [RR-19] Chuỗi cung ứng hạt (Epic Strategy)

**Role:** Tech Lead
**Context:** Kiến trúc tổng thể cho luồng Batching và Traceability trong Warehouse Service.

---

## 1. Architectural Patterns

### 1.1. Event-Driven Bridge
Warehouse Service đóng vai trò là cầu nối giữa Farm (Production) và Retail (Consumption).
- **Inbound:** Lắng nghe `farm.harvest.created`.
- **Outbound:** Phát hành `warehouse.stock.updated`, `warehouse.stock.reserved`.

### 1.2. Automated Traceability ID
Hệ thống sẽ cung cấp một `BatchService` chuyên biệt để sinh mã theo logic tại `docs/business/BATCH_LOGIC.md`. Service này sử dụng Redis/Valkey `INCR` để đảm bảo tính duy nhất của `Sequence` trong ngày.

---

## 2. Weight Transformation Logic

Hệ thống sẽ thực hiện tính toán tự động tại tầng UseCase:
- `ExpectedYield = IntakeWeight * (1 - DefaultLossRate)`
- `ActualLoss = (1 - YieldWeight / IntakeWeight) * 100`

Mọi logic tính toán được bọc trong Domain service `WeightCalculator` để dễ dàng Unit Test mà không cần DB.

---

## 3. Sub-ticket Implementation Strategy

### RR-19.1 (Intake)
- **Tech:** Kafka Consumer cho `farm.harvest.events`.
- **Atomic Operation:** Ghi bản ghi `processing_batches` + Đánh dấu Inbox trong 1 Transaction.

### RR-19.2 (Processing)
- **Tech:** gRPC API `UpdateBatchProgress`.
- **Validation:** Trigger State Machine transition. Nếu status là `STOCKED`, tự động trigger logic cập nhật kho.

### RR-19.3 (Stocking)
- **Tech:** UPSERT logic cho `inventory_items`.
- **SKU Generation:** `SKU = COFFEE_TYPE + "-" + ORIGIN_CODE + "-ROASTED"`.
