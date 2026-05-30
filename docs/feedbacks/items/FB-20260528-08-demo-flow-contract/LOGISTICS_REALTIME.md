# Logistics Realtime Get/Publish

## Current Implementation State

- UI reads shipments through logistics backend APIs.
- UI reads locations through logistics backend APIs.
- Driver Client replays route points from static `/data/routes.json`.
- Browser posts GPS updates to the logistics API.
- Backend validates driver identity and shipment scope where claims are available.
- Backend stores accepted GPS in Valkey GEO data and writes last-seen TTL.
- Backend updates vehicle coordinates in Postgres.
- Backend publishes `logistics.gps.updated` to Kafka.
- Trace and audit services can consume emitted events.
- Current dashboard uses polling/browser state for the visible moving marker; it is not a full WebSocket realtime client yet.

## Data Path

```text
Driver browser route replay
  -> POST /v1/logistics/gps
  -> logistics-service UpdateDriverLocation
  -> validate driver/shipment scope
  -> Valkey GEO + last_seen TTL
  -> vehicle coordinate update
  -> Kafka logistics.gps.updated
  -> trace/audit/topology consumers
```

## Simulated Pieces

- GPS movement source is browser route replay, not a physical mobile device.
- Route coordinates are static demo data from `routes.json`.
- Dashboard moving marker can be browser-local while polling catches backend state.

## Omitted Pieces

- Real driver mobile app.
- Device GPS permissions and sensor telemetry.
- WebSocket/SSE map streaming as the primary logistics UI path.
- ETA, route optimization, geofencing, and map matching.
- Active-route-only rendering cleanup is still backlog where static routes appear without an active shipment.

## Recommended Demo Script

1. Log in as `WAREHOUSE_MGR` in one browser profile.
2. Ensure a shipment is assigned to a driver.
3. Log in as `DRIVER` in a separate profile.
4. Open `/dashboard/driver`.
5. Start route simulation.
6. Confirm pickup or delivery milestones in order.
7. Open `/dashboard/logistics` in a logistics/admin profile.
8. Verify shipment state, moving marker behavior, and Kafka `logistics.gps.updated`.
9. Open trace/topology views to confirm GPS and shipment events appear when consumers are running.

## Review Language

Use: "The browser simulates the device movement, while the backend validates and publishes accepted GPS updates."

Avoid: "The system uses real vehicle telematics."
