# Logistics Simulation Guide

This document explains how the real-time logistics simulation works in Runtime Roasters, from seeding locations to visualizing driver movement on actual road paths.

---

## 1. Seeding Strategy

The system uses a set of predefined locations in Vietnam to provide a realistic "Farm-to-Cup" experience. These locations are seeded into the `logistics-service` database.

### Core Locations (POIs)
- **Farmers:** Specific farms in Cầu Đất (Dalat), Buôn Ma Thuột, and Pleiku.
- **Roasteries:** Processing centers in major industrial zones (Sóng Thần, Hòa Lạc, Hòa Khánh).
- **Retailers:** Retail stores in central districts of HCM, Hanoi, and Da Nang.

**Seed File:** `deployments/logistics-seed.sql`

---

## 2. Realistic Route Generation (OSRM)

Unlike simple straight-line movement, our simulation uses actual road data. We use the **Open Source Routing Machine (OSRM)** to calculate paths between our seeded locations.

### The Route Generator Tool
Located at `src/scripts/generate_routes.go`, this tool:
1. Iterates through all logical pairs of locations (e.g., Farm -> Roastery).
2. Calls the OSRM Public API: `http://router.project-osrm.org/route/v1/driving/{lng,lat;lng,lat}`.
3. Extracts the coordinate list (polyline) representing the actual road path.
4. Exports the results to a static JSON file: `src/apps/logistics-service/testdata/routes.json`.

**Why static files?** To ensure the simulation remains stable, offline-capable for demos, and prevents hitting rate limits on public APIs during development.

---

## 3. Real-time Movement Simulation

The primary production-demo simulator is the authenticated **Driver Client**. The driver browser animates seeded route points and posts GPS/status updates to backend like a real driver device.

### How it works:
1. **Login:** User authenticates as `DRIVER`.
2. **Assigned Shipment:** Driver Client fetches only assigned shipments.
3. **Load Paths:** UI reads the seeded route from `routes.json` or an API backed by the same route data.
4. **Replay Path:** UI iterates through the coordinate list at a configured demo speed, usually 30-45 seconds per leg.
5. **Update Location:** Every few seconds, UI posts the current coordinate and shipment status to Logistics Service.
6. **Backend Validation:** Logistics Service validates that the driver is assigned to the shipment before accepting the update.
7. **Freshness Tracking:** Logistics Service stores current coordinates in **Valkey GEO** and sets a TTL key to track driver liveness.
8. **Persistence:** Shipment state and milestone confirmations remain persisted in Postgres.

`src/scripts/simulate_drivers.go` was researched/documented as a possible standalone simulator, but its implementation must be verified before use. If absent, keep Driver Client simulation as the main demo path and implement a backend/script fallback only for unattended demos.

---

## 4. UI Visualization Flow

The Frontend leverages this data to create a modern, animated tracking experience.

1. **Static Routes:** The UI fetches the path coordinates (via API or static asset) to draw lines on the map.
2. **Live Updates:** Dashboards connect to the socket/realtime service or trace-service SSE/WebSocket when private operational data is streamed.
3. **Event Pipeline:** 
   - Driver Client -> Logistics Service (`POST /v1/logistics/drivers/location` or equivalent)
   - Logistics Service -> Postgres shipment state + Valkey GEO/liveness
   - Logistics Service -> Kafka (`logistics.gps.updated`, `logistics.shipment.status_changed`)
   - Socket/Trace Service -> Consumer (Kafka) -> Client (SSE/WS)
4. **Rendering:** The UI updates the driver icon's position on the map based on the live coordinates, snapping them to the predefined road path for a smooth experience.

---

## 5. Running the Simulation

1. **Seed the database:**
   ```bash
   cat deployments/logistics-seed.sql | docker exec -i rr-postgres psql -U postgres -d logistics_db
   ```
2. **Generate routes (if POIs changed):**
   ```bash
   go run src/scripts/generate_routes.go
   ```
3. **Run the demo simulation:**
   - Log in with `driver@runtimeroasters.com`.
   - Open the assigned Driver/Logistics shipment screen.
   - Click Start to replay the seeded route.
   - Watch Logistics/Warehouse/Trace/Architecture screens receive realtime updates.

4. **Optional unattended simulator:**
   ```bash
   go run src/scripts/simulate_drivers.go
   ```
   This command is optional and only valid if the script exists in the current implementation.
