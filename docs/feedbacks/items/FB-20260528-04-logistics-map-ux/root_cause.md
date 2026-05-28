# Root Cause Analysis

**Symptom**: Logistics map is in dark mode. Static routes are visible at all times.
**Actual Cause**: The `LogisticsMap.tsx` hardcodes the Mapbox/Leaflet URL to a dark theme (`dark_all`). The `Polyline` component renders all static routes from `routes.json` regardless of vehicle activity.
