# Phase 07 Research: Real-time Logistics Simulation

## 1. Seed Data Requirements
The user requires specific locations in Vietnam to be seeded:
- **Farmers:** Cau Dat (Dalat), BMT, Pleiku.
- **Roasteries:** Song Than (Binh Duong), Hoa Lac (Hanoi), Hoa Khanh (Da Nang).
- **Retailers:** District 1 (HCM), Hoan Kiem (Hanoi), Hai Chau (Da Nang).

These locations should be stored in the `logistics-service` for route calculation and distance estimation.

## 2. Route Finding Tool
To get realistic road routes between locations, we will use **OSRM (Open Source Routing Machine)**.
- **API:** `http://router.project-osrm.org/route/v1/driving/{coordinates}`
- **Output:** A list of coordinates (polyline or array) representing the actual road path.
- **Format:** GeoJSON or simple JSON array of `[lat, lng]`.

### Tool Implementation: `src/scripts/generate_routes.go`
- Inputs: A list of start/end point pairs (e.g., FARM_CAU_DAT -> WH_SONG_THAN).
- Logic: Call OSRM API for each pair.
- Output: `src/apps/logistics-service/testdata/routes.json`.

## 3. Real-time Simulation Strategy
The simulator (`src/scripts/simulate_drivers.go`) will:
1. Load `routes.json`.
2. For each active driver, pick a route (e.g., a shipment assignment).
3. "Step" through the coordinates in the route at a simulated speed.
4. Call the Logistics Service `UpdateLocation` gRPC API with the current position.
5. Handle looping or returning to base.

## 4. UI Visualization
The UI will:
1. Fetch the static route data (or the service will provide it via API).
2. Use a map library (e.g., Leaflet or Mapbox GL JS) to draw the polyline.
3. Display the driver's current position (received via WebSocket/SSE from `Monitor Service`) on the route.

## 5. Implementation Changes
- **07-01:** Add `locations` table and seed SQL.
- **07-03:** Add `generate_routes.go` and update `simulate_drivers.go`.
- **Logistics Domain:** Add `Location` entity.
