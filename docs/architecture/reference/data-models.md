# Database Design

## Overview

## Transactional Patterns (Reliability)

### Table: `outbox_events`
Dùng để lưu trữ sự kiện trước khi gửi lên Kafka.

| Cột | Kiểu | Mô tả |
| :--- | :--- | :--- |
| `id` | UUID | Unique ID của sự kiện. |
| `event_type` | STRING | Loại sự kiện (CloudEvents type). |
| `payload` | JSONB | Nội dung nghiệp vụ. |
| `status` | STRING | `PENDING`, `PROCESSED`, `FAILED`. |
| `retry_count` | INT | Số lần thử lại. |

### Table: `inbox_events`
Dùng để chống xử lý lặp sự kiện (Idempotency).

| Cột | Kiểu | Mô tả |
| :--- | :--- | :--- |
| `id` | UUID | Unique ID từ Kafka/CloudEvents (Primary Key). |
| `processed_at` | TIMESTAMPTZ | Thời điểm xử lý thành công. |

## Elasticsearch Index (Traceability)

### Index: `coffee_traceability`
Lưu trữ Read Model của toàn bộ vòng đời sản phẩm. Được denormalize từ nhiều sự kiện.

| Trường | Kiểu | Mô tả |
| :--- | :--- | :--- |
| `batch_id` | KEYWORD | ID duy nhất của mẻ hàng (Primary Key). |
| `product_name` | TEXT | Tên sản phẩm thương mại. |
| `coffee_type` | KEYWORD | Loại hạt (Arabica, Robusta...). |
| `origin.farm_name` | TEXT | Tên nông trại khởi nguồn. |
| `origin.location` | GEO_POINT | Tọa độ GPS vùng trồng. |
| `origin.region` | KEYWORD | Vùng địa lý (Cầu Đất, v.v.). |
| `processing.method` | KEYWORD | Phương pháp chế biến (Honey, Washed...). |
| `processing.roast_level` | KEYWORD | Mức độ rang. |
| `timeline` | NESTED | Mảng các cột mốc: `{status, timestamp, description}`. |
| `last_updated` | DATE | Thời điểm cập nhật cuối cùng. |

## Kafka Event Schemas (Async Messaging)

### Topic: `auth.user.events`
Sử dụng để đồng bộ thông tin định danh và quyền hạn trên toàn hệ thống.

#### Event: `USER_CREATED`
| Trường | Kiểu | Mô tả |
| :--- | :--- | :--- |
| `event_id` | UUID | ID duy nhất của sự kiện. |
| `event_type` | STRING | Luôn là `USER_CREATED`. |
| `payload.user_id` | UUID | Subject ID từ Ory Kratos. |
| `payload.email` | STRING | Email người dùng. |
| `payload.name` | STRING | Tên hiển thị. |
| `payload.role` | STRING | Role cố định (`FARM_ADMIN`, `FARM_MANAGER`). |
| `occurred_at` | TIMESTAMPTZ | Thời điểm phát sinh sự kiện. |
