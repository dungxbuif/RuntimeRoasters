# RR-URG-08: Role Authorization, Demo Accounts, And Seeded Journey Data

## Priority

P2. Can run in parallel with backend work once event contracts are stable.

## Problem

The demo needs real accounts, assignments, vehicles, routes, inventory, paid orders, and prepared trace products. Role scoping must match the UI flow.

## Scope

- Verify role policies.
- Fix known Casbin direct-role matcher issue if still present.
- Seed accounts and assignments.
- Seed prepared data for QR trace.
- Seed routes/vehicles/drivers/products.

## Implementation Details

### 1. Roles

Required roles:

- `ADMIN`
- `FARM_MANAGER`
- `WAREHOUSE_MGR`
- `STORE_MGR`
- `DRIVER`
- optional `PROCESSOR`

Role behavior:

- `ADMIN`: setup, assignment, overview only by default.
- `FARM_MANAGER`: assigned farms.
- `WAREHOUSE_MGR`: assigned warehouses.
- `STORE_MGR`: assigned stores via `store_ids`.
- `DRIVER`: assigned shipments only.

Checklist:

- [ ] Policies defined.
- [ ] Role inheritance/direct matching works.
- [ ] `STORE_MGR` missing `store_ids` sees no data.
- [ ] `DRIVER` without assigned shipment sees no simulation controls.

### 2. Casbin Gap

Known risk:

- matcher may only use `g(r.sub,p.sub)` and not direct `r.sub == p.sub`.

Required matcher support:

- direct role match: `r.sub == p.sub`
- inherited role match: `g(r.sub, p.sub)`

Checklist:

- [ ] Matcher fixed if needed.
- [ ] Unit tests added.
- [ ] Role policy regression tests pass.

### 3. Demo Accounts

Verify or seed:

- `admin@runtimeroasters.com`
- `manager.caudat@runtimeroasters.com`
- Warehouse Manager account.
- `driver@runtimeroasters.com`
- `mgr.hn.hoankiem@runtimeroasters.com`
- optional `processor@runtimeroasters.com`

Checklist:

- [ ] Passwords documented.
- [ ] Roles documented.
- [ ] Assignments documented.
- [ ] JWT claims include role/email/org/store_ids as needed.

### 4. Seeded Operational Data

Seed:

- farms.
- warehouses.
- stores.
- vehicles.
- drivers.
- routes.
- inventory SKUs.
- paid order capable product.
- public trace products.

Checklist:

- [ ] At least one complete Farm -> Warehouse -> Retail journey can be run manually.
- [ ] At least two public QR trace products exist.
- [ ] Routes cover warehouse -> farm -> warehouse and warehouse -> store -> warehouse.
- [ ] Vehicle capacity enough for seeded cargo.

### 5. Prepared Public Trace Data

Options:

- replay seed events through services.
- load trace-service projection fixtures.
- run demo script once and snapshot public trace products.

Preferred:

- replay seed events so trace-service/Elasticsearch contains realistic data.

Checklist:

- [ ] Public trace code generated.
- [ ] Elasticsearch document exists.
- [ ] QR product list returns seeded products.

## Acceptance Criteria

1. All demo accounts can log in.
2. Each role sees only intended dashboards/actions.
3. Seeded data supports farm pickup, paid order delivery, and public QR trace.
4. Public QR trace data is real prepared trace data.
5. Auth tests prove fail-closed scoping.

## Test Checklist

- [ ] Auth unit tests.
- [ ] Casbin matcher tests.
- [ ] JWT claims tests.
- [ ] Manual login for each demo account.
- [ ] Manual scoping check for store manager A vs B.
- [ ] Seed reset script produces expected baseline.
