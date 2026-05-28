# Root Cause Analysis

**Symptom**: Accessing `/dashboard/orders` returns a 404.
**Actual Cause**: The Next.js router does not have a directory for `app/(dashboard)/dashboard/orders`. Orders are currently managed inside `/dashboard/retail/orders` or similar.
