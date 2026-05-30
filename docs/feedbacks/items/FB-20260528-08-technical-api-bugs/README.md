# Technical API Bugs (405, Network Error, CORS)

## Status: REOPENED

## Description
Manual testing revealed several API issues related to the Orders flow:
1. `GET /v1/orders` returns `405 Method Not Allowed`.
2. `POST /v1/orders` returns `Network Error`.
3. CORS error on `http://localhost:8081/v1/orders`.

## Update (2026-05-29)
- Ticket reopened to track systematic connectivity and gateway configuration fixes.
- Initial fixes for routing and CORS headers applied.
