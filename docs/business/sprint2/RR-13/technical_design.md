# Dev Notes - [RR-13] End-to-End Secure Integration

## 🛠️ Verification Steps
1. Start full stack: `docker compose up`.
2. Login via Client App → obtain JWT.
3. Inspect JWT payload: confirm `sub`, `role`, `org_id`, `scope` claims present.
4. Test Gate 1: call endpoint with a JWT missing the required scope → expect `403` from KrakenD.
5. Test Gate 2: call endpoint with correct scope but wrong role → expect `403` from Service.
6. Test full happy path: correct scope + correct role → expect `200` with data.
7. Logout → replay JWT → expect `401` (blacklist check).
8. Check SigNoz: confirm `user.id` attribute in trace spans.
9. Simulate key rotation: restart Identity Server → confirm services recover without dropped requests.
