# Root Cause Analysis

**Symptom**: Admin lacks UI to view/manage Warehouses/Stores. Sidebar is missing the Users menu link. Policy management is overly complex.
**Actual Cause**: The `Sidebar.tsx` omitted the global Users link outside of its own context. The `/dashboard/resources` page is incomplete. `SPEC.md` had a gap between designed creation vs actual implementation (Auto-seeding).
**Note**: `SPEC.md` and `AGENTS.md` have already been updated to reflect the Keep-It-Simple policy.
