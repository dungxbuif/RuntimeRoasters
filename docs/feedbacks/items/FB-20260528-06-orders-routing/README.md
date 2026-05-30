# Orders API & UI Issues

## Status: REOPENED

## Description
1. UI instability on the orders page.
2. `GET /v1/orders` returns 405.
3. `POST /v1/orders` returns Network Error.
4. CORS issues on `/v1/orders`.

## Update (2026-05-29)
- Ticket reopened to address specific API failures and UI instability.
- Fixed: Added GET route to KrakenD and Retail Service.
- Fixed: Added Idempotency-Key to KrakenD CORS headers.
- Fixed: Restrict page access to STORE_MGR only.
