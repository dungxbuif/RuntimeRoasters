# Implementation Plan

- **Expected Changes**:
  1. Change TileLayer URL to a light theme variant.
  2. Filter the `routes` array to only render a `Polyline` if its `route.id` matches an active shipment's route.
- **Impacted Scope**: Logistics Map and Driver Map components.
- **Required Validation**: View `/dashboard/logistics` and confirm light theme and hidden static routes.
- **Expected Impact**: Cleaner UI, aligned with overall Light mode preference.
