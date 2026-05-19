# 📚 Domain Knowledge Base: Structure & Strategy

This document outlines the organization of the detailed domain knowledge base for Runtime Roasters. While `BUSINESS_SPECIFICATION.md` provides a high-level executive summary, the files in this directory dive deep into the specific business rules, calculations, and state machines for each operational domain.

## Directory Structure

To maintain a clean and scalable knowledge base, the domain logic is categorized into six primary modules:

1.  **`01-FARM_OPERATIONS.md`**: Focuses on upstream activities (planting, harvesting, yield verification).
2.  **`02-PROCESSING_INVENTORY.md`**: Details midstream activities (roasting, weight loss calculations, warehousing, FIFO).
3.  **`03-ORDER_FULFILLMENT.md`**: Covers downstream activities (customer orders, payment states, and reservation logic).
4.  **`04-LOGISTICS_TRACKING.md`**: Explains the rules for driver assignment, routing, and real-time movement.
5.  **`05-ROLE_UI_MATRIX.md`**: Defines which role sees which dashboard actions and how UI actions advance the business flow.
6.  **`06_PUBLIC_TRACE_DEMO.md`**: Defines the public QR traceability showcase backed by prepared real trace data.

## Canonical Production-Demo Flow

The current canonical showcase flow is:

1. Farm Manager declares a harvest.
2. Warehouse receives a pickup request, not an immediate intake.
3. Warehouse Manager dispatches a vehicle/driver.
4. Driver Client simulates movement on seeded routes and posts GPS/status updates to backend.
5. Driver confirms pickup, return, delivery, and completion milestones.
6. Warehouse creates intake only after pickup return/receipt.
7. Existing processing and inventory flow continues.
8. Paid retail order reserves stock and triggers warehouse-to-store delivery.
9. Socket/SSE updates keep dashboards and the public architecture showcase visually in sync.
10. Public QR trace page shows prepared real product journeys from trace-service/Elasticsearch.

The demo intentionally allows a `DRIVER` user to drive the movement simulation from the browser. Backend services remain the persisted source of truth: they validate assignment, store accepted updates, emit events, and broadcast realtime state.

## Role Principle

`ADMIN` is a setup, assignment, and overview role. It creates accounts and business nodes, assigns managers, and views aggregate health. Detailed operational actions belong to the assigned manager roles and `DRIVER`.

Public/no-auth UI is allowed only for the root Client App architecture/topology showcase with sanitized data. Private dashboards and driver simulation flows require authentication and role/entity scoping.
