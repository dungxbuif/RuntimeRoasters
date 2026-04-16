# 📋 Sprint 2: The Distributed Pulse — Kafka, Trace Service & Traceability GraphQL Engine

**Status:** `IN_PLANNING` | **Timeline:** `Scope-based` | **PO:** `User` | **Tech Lead:** `Gemini`

---

## 🏗️ Kanban Board

| 🕒 To Do | 🚧 In Progress | 🔍 Review | ✅ Done |
| :--- | :--- | :--- | :--- |
| [RR-9: Kafka & Outbox Setup](./RR-9.md) | | | |
| [RR-10: Trace Service — Kafka Consumer](./RR-10.md) | | | |
| [RR-11: Traceability GraphQL Engine](./RR-11.md) | | | |
| [RR-12: QR Lifecycle UI](./RR-12.md) | | | |

---

## 📝 Sprint Goal

Thiết lập xương sống sự kiện bất đồng bộ (Kafka + Outbox), xây dựng Trace Service tiêu thụ sự kiện vào Elasticsearch, và triển khai **Traceability GraphQL Engine** — lớp truy vấn đồ thị vòng đời sản phẩm phục vụ tính năng quét mã QR.

> **Ghi chú Kiến trúc:** Tham khảo [`docs/architecture/graphql-integration.md`](../../architecture/graphql-integration.md) để biết thiết kế đầy đủ của cả GraphQL BFF (Future Phase) và Traceability Engine (Sprint 2 này).

---

## 🎯 Sprint Objectives

1. **Kafka Backbone:** Cài đặt Transactional Outbox tại Farm Service. Event từ Farm được publish lên Kafka một cách đảm bảo.
2. **Trace Service Consumer:** Trace Service tiêu thụ event từ Kafka, chuẩn hóa và lưu vào Elasticsearch.
3. **Traceability GraphQL Engine:** Expose GraphQL endpoint `/trace/graphql` chuyên biệt cho truy vấn vòng đời sản phẩm theo cấu trúc đồ thị.
4. **QR UI:** Giao diện quét mã QR hiển thị hành trình từ Ly cà phê ➔ Chuyến xe ➔ Mẻ rang ➔ Lô thu hoạch.
