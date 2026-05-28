# Multi-Account Demo

## Problem

Using multiple personas in the same browser profile can produce misleading auth behavior because the profile shares:

- Kratos browser session cookies.
- Hydra/OIDC browser state.
- app access token in `localStorage`.
- React/query cache and role-scoped UI state.

## Recommended Browser Setup

| Context | Persona | Purpose |
| --- | --- | --- |
| Profile A | `ADMIN` | Bootstrap, aggregate dashboard, topology. |
| Profile B | `STORE_MGR` | Store order creation and incoming delivery visibility. |
| Profile C | `WAREHOUSE_MGR` | Warehouse stock, dispatch, and pickup work. |
| Profile D | `DRIVER` | Driver Client route replay and milestone confirmation. |
| Public/incognito | Anonymous | Root topology or public trace views. |

## Login/Logout Cautions

- Do not log in as two personas in the same browser profile during a demo.
- If using the same profile is unavoidable, explicitly log out, clear app token state, and reset the Kratos session before switching personas.
- Incognito windows in the same browser can still share state depending on browser behavior; prefer separate named profiles for live review.
- When testing root `/` callback behavior, verify the URL token is removed after the app stores it.

## Recommended Demo Sequence

1. Open all browser profiles before the demo.
2. Log each profile into exactly one persona.
3. Use `ADMIN` profile to initialize data if needed.
4. Use `STORE_MGR` profile to create a retail order.
5. Use `Finance` or admin-capable profile to pass/fail the demo webhook if required.
6. Use `WAREHOUSE_MGR` profile to inspect stock/dispatch state.
7. Use `DRIVER` profile to run route simulation and confirm milestones.
8. Use public/incognito profile for root topology or public trace surfaces.

## Expected Outcome

Each profile has isolated cookies and `localStorage`, so role-scoped pages and APIs behave consistently.
