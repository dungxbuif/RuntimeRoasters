# Demo Simulation Strategy

## Summary

- Simulate at the system edge; keep the backend core real.
- Services, DB state, Kafka events, auth, Casbin, stock reservation, dispatch, logistics state, trace, and audit must run through real APIs/events.
- Simulation represents real-world external actors only: payment providers, driver GPS devices, long-running time, and initial historical demo data.

## Canonical Simulation Boundaries

- **Payment simulation**: Finance demo controls send signed Stripe-compatible webhooks through `finance.service.ts`. Retail must not post unsigned mock payment payloads.
- **Driver/GPS simulation**: Driver Client replays `/data/routes.json` and posts each tick to the logistics API. The browser simulates device movement; logistics validates, stores, and publishes accepted GPS.
- **Warehouse/dispatch flow**: Dispatch queues are real runtime state. Dispatch requests are created after stock reservation, and Warehouse Manager assignment publishes real events for logistics.
- **Historical bootstrap**: Seeders may insert read-model history for immediate dashboard readiness. Live demos still use the runtime event chain.
- **Time compression**: Roasting, delivery travel, and return-to-base may be accelerated by timers or demo buttons, but they must not bypass domain state transitions or Kafka events.

## Flow Labels

- `REAL_RUNTIME`: OIDC, RBAC/Casbin, retail orders, payment rows, reservation, dispatch requests, shipment transitions, Kafka, trace, audit.
- `BACKEND_SIMULATED`: first-run historical dashboard rows, seeded users/resources/routes, provider settlement emulator.
- `BROWSER_SIMULATED`: driver GPS movement, route replay, timed milestone progression.
- `UI_MOCK_ONLY`: visual fallback only; never the source of truth for the main demo.

## Verification Expectations

- Go tests cover retail SAGA transitions, warehouse reservation/dispatch events, logistics assigned-shipment handling, and payment webhook behavior.
- Frontend checks must pass `npm run lint` and `npm run build`.
- Manual acceptance follows: store creates order -> finance passes payment -> warehouse dispatches -> driver route replay runs -> retail confirms receipt -> trace/audit/topology show resulting events.
