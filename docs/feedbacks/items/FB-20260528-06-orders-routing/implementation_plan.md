# Implementation Plan

- **Expected Changes**: Create a simple page at `app/(dashboard)/dashboard/orders/page.tsx` that uses `redirect('/dashboard/retail')` or displays a unified order list.
- **Impacted Scope**: Dashboard routing.
- **Required Validation**: Navigate to `/dashboard/orders` and confirm it resolves without 404.
- **Expected Impact**: Fixes the broken link.
