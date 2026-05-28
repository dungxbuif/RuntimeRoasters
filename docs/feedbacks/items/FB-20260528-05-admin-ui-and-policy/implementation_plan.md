# Implementation Plan

- **Expected Changes**:
  1. Add `Users` link to `Sidebar.tsx` globally for ADMIN.
  2. Implement data tables in `/dashboard/resources` for Warehouses and Stores, displaying their auto-seeded vehicles and assigned managers.
  3. Ensure user creation UI supports the "Delete and Re-create" flow for role assignment.
- **Impacted Scope**: Admin Resource Management UI.
- **Required Validation**: Login as ADMIN, click Users in sidebar, and view resources in the Resources tab.
- **Expected Impact**: Admin has full operational visibility over nodes and users.
