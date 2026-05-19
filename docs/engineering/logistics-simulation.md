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

The **Driver Simulator** (`src/scripts/simulate_drivers.go`) breathes life into the system.

### How it works:
1. **Load Paths:** It reads `routes.json` to get the set of available road paths.
2. **Assign Routes:** For each active driver, it assigns a random or specific route.
3. **Replay Path:** It iterates through the coordinate list at a simulated speed.
4. **Update Location:** Every few seconds, it calls the Logistics Service gRPC API `UpdateLocation` with the driver's current coordinates.
5. **Freshness Tracking:** The Logistics Service stores these coordinates in **Valkey GEO** and sets a TTL key to track driver "liveness".

---

## 4. UI Visualization Flow

The Frontend leverages this data to create a modern, animated tracking experience.

1. **Static Routes:** The UI fetches the path coordinates (via API or static asset) to draw lines on the map.
2. **Live Updates:** The UI connects to the **Monitor Service** (via SSE or WebSocket).
3. **Event Pipeline:** 
   - Simulator -> `UpdateLocation` (gRPC)
   - Logistics Service -> Valkey (GEO) & Kafka (`logistics.gps.updated`)
   - Monitor Service -> Consumer (Kafka) -> Client (SSE/WS)
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
3. **Start the simulator:**
   ```bash
   go run src/scripts/simulate_drivers.go
   ```
