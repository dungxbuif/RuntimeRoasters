# Database Design

## Overview

## PostgreSQL Schema (Transational)

## Apache Cassandra Schema (Audit Log)

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
| `payload.role` | STRING | Role cố định (`farm_admin`, `farm_manager`). |
| `occurred_at` | TIMESTAMPTZ | Thời điểm phát sinh sự kiện. |
