---
phase: 07-real-time-logistics
reviewed: 2025-05-15T10:30:00Z
depth: standard
files_reviewed_list:
  - .planning/phases/07-real-time-logistics/07-01-PLAN.md
  - .planning/phases/07-real-time-logistics/07-02-PLAN.md
  - .planning/phases/07-real-time-logistics/07-03-PLAN.md
  - .planning/phases/07-real-time-logistics/07-04-PLAN.md
findings:
  critical: 3
  warning: 4
  info: 2
  total: 9
status: issues_found
---

# Phase 07: Code Review Report (Real-time Logistics)

**Reviewed:** 2025-05-15
**Depth:** standard
**Files Reviewed:** 4
**Status:** issues_found

## Summary

The implementation plans for the Logistics Service provide a solid foundation for geospatial tracking and shipment management. However, there are critical gaps in the **driver assignment lifecycle** and **data freshness** in Valkey that could lead to "ghost" assignments and stuck shipments. Additionally, some proposed patterns violate the project's **Clean Architecture** guidelines.

## Critical Issues

### CR-01: Deadlock in Shipment Assignment (No Retry Mechanism)

**File:** `.planning/phases/07-real-time-logistics/07-02-PLAN.md` & `07-04-PLAN.md`
**Issue:** The shipment assignment is triggered only once by the `warehouse.stock.reserved` Kafka event. If no driver is found within the radius at that exact moment, or if the assigned driver rejects the shipment (`ASSIGNED -> PENDING` transition), the shipment stays in `PENDING` state indefinitely. There is no background worker or retry loop to re-process `PENDING` shipments.
**Fix:** Implement a "Dispatching Worker" (e.g., a cron job or a background goroutine) that periodically queries the database for `PENDING` shipments and attempts to find available drivers.

### CR-02: Stale Driver Locations in Valkey (Ghost Drivers)

**File:** `.planning/phases/07-real-time-logistics/07-03-PLAN.md`
**Issue:** The plan uses `GEOADD` to store driver locations in Valkey. Valkey/Redis GEO sets do not support per-member TTL. If a driver goes offline or loses connectivity, their last known location remains in the `drivers:locations` set forever. `GEOSEARCH` will continue to return these stale drivers, leading to failed assignment attempts.
**Fix:** Maintain a separate key for each driver with a TTL (e.g., `driver:active:<id>` with 5-minute TTL) updated during `UpdateLocation`. In `FindNearestAvailableDriver`, cross-reference the results from `GEOSEARCH` with these TTL keys to ensure only "fresh" drivers are considered.

### CR-03: Race Condition in Driver Availability

**File:** `.planning/phases/07-real-time-logistics/07-02-PLAN.md`
**Issue:** The plan does not explicitly state that a driver's `IsAvailable` status in the database is updated to `false` when they are assigned a shipment. This allows the system to assign multiple shipments to the same driver simultaneously if they appear in the `GEOSEARCH` results for different orders.
**Fix:** The `AssignShipment` use case must atomically check availability and set `IsAvailable = false` for the selected driver within a database transaction.

## Warnings

### WR-01: Clean Architecture Violation (Interfaces in Domain)

**File:** `.planning/phases/07-real-time-logistics/07-01-PLAN.md:Task 2`
**Issue:** The plan proposes defining "repository interfaces" inside `internal/domain/shipment.go`. According to the `tech-lead-reviewer` skill rules, interfaces MUST be defined in the layer that consumes them (usually `usecase`), and `domain` must have zero external imports (no "Repo" in Domain).
**Fix:** Move repository interface definitions to the `usecase` layer or a dedicated port layer. Keep only pure entities in `internal/domain`.

### WR-02: Use of Deprecated Redis/Valkey Commands

**File:** `.planning/phases/07-real-time-logistics/07-02-PLAN.md:Task 2`
**Issue:** The plan mentions `GEORADIUS`. This command is deprecated in favor of `GEOSEARCH` in modern Valkey/Redis versions.
**Fix:** Use `GEOSEARCH` with `BYRADIUS` or `BYBOX` for more efficient and future-proof queries.

### WR-03: High-Frequency Kafka Event Overload

**File:** `.planning/phases/07-real-time-logistics/07-03-PLAN.md:Task 2`
**Issue:** Publishing a Kafka event `logistics.gps.updated` for every single GPS ping without throttling or a clearly defined consumer can lead to unnecessary resource consumption and potential bottlenecks.
**Fix:** Implement throttling (e.g., publish only every 30 seconds or if the driver has moved > 50 meters) unless a high-resolution real-time consumer is explicitly required.

### WR-04: Incomplete State Machine Failure Paths

**File:** `.planning/phases/07-real-time-logistics/07-04-PLAN.md:Task 1`
**Issue:** The state machine handles rejections but misses other critical failure modes such as `PICKED_UP -> FAILED` (accidents/breakdowns) or `CANCELLED` (order cancelled by user).
**Fix:** Add `FAILED` and `CANCELLED` states to the `ShipmentStatus` enum and implement corresponding transition logic.

## Info

### IN-01: Idempotency Implementation

**File:** `.planning/phases/07-real-time-logistics/07-01-PLAN.md`
**Issue:** Good use of `processed_kafka_messages` table and transactional outbox pattern ensures reliable and idempotent processing of warehouse events.

### IN-02: Driver Simulator

**File:** `.planning/phases/07-real-time-logistics/07-03-PLAN.md`
**Issue:** The inclusion of a dedicated simulation script (`src/scripts/simulate_drivers.go`) is a great practice for testing complex geospatial and real-time logic.

---

_Reviewed: 2025-05-15_
_Reviewer: gsd-code-reviewer (Tech Lead Role)_
_Depth: standard_
