# Admin UI and Policy Gaps

## Status: REOPENED

## Description
1. ADMIN being able to declare harvests (Business logic error).
2. ADMIN being able to see/access retail orders creation page.
3. Missing UI for Admin to create Warehouses and Retail stores.
4. Business gap: Admin should create entities and assign managers, while Warehouse Manager handles logistics assignment.

## Update (2026-05-29)
- Ticket reopened due to newly discovered business logic gaps and permission issues reported during manual testing.
- Fixed: Removed ADMIN permission to create harvests and orders in Casbin.
- Fixed: UI button guards for ADMIN on harvests and orders.
- Remaining: Missing Admin UI for Warehouse/Retail creation.
