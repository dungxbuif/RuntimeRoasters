# Identity & SSO Implementation (Sprint Technical Doc)

## Overview
This document records the technical decisions, fixes, and edge cases encountered during the implementation of the Zero-Consent OIDC Single Sign-On (SSO) flow using Ory Kratos, Ory Hydra, and Next.js (Client App).

## Key Implementation Details

### 1. Kratos & Hydra Direct Integration (`oauth2_provider`)
**Problem:** Initially, passing `login_challenge` to Kratos caused a `500 Internal Server Error` (`refusing to parse login_challenge query parameter because oauth2_provider.url is invalid or unset`).
**Solution:** Configured Kratos to talk directly to Hydra via `kratos.yaml`:
```yaml
oauth2_provider:
  url: http://hydra:4445/
```
**Result:** Kratos natively understands the OAuth2 flow. When authentication succeeds, Kratos automatically calls Hydra's Admin API to accept the login request and returns a `redirect_browser_to` URL (Consent/Callback), eliminating the need for manual backend wrappers in Next.js.

### 2. Resolving the Infinite SSO Loop (Logout Fix)
**Problem:** Clicking "Logout" cleared `localStorage` and redirected to `/login`. However, Kratos still maintained a valid session cookie. When `/login` detected `session_already_available`, the app redirected to Home, which redirected to Hydra, which redirected to `/login?login_challenge`. Kratos automatically accepted this challenge and logged the user back in.
**Solution:** Replaced standard `localStorage.clear()` with the official Kratos Logout Flow:
```typescript
const handleLogout = async () => {
  storageService.clearAll();
  const { default: kratos } = await import('@/lib/ory/kratos');
  const { data } = await kratos.createBrowserLogoutFlow();
  window.location.href = data.logout_url; // Revokes cookie and redirects to /login
};
```

### 3. Token Exchange Endpoint Shift
**Problem:** Next.js Server needed to exchange the authorization code for an access token. Using `HYDRA_ADMIN_URL` (port 4445) returned a `404 Not Found`.
**Solution:** The `/oauth2/token` endpoint is only exposed on the **Public Port** (4444). The `/api/auth/callback/route.ts` was updated to use `HYDRA_PUBLIC_URL`.

### 4. Automated OAuth2 Client Bootstrapping
**Problem:** Hydra database was empty on first start, causing `invalid_client` errors during the OAuth2 challenge.
**Solution:** Created `client-app.json` and utilized `hydra import clients` to seed the database consistently on startup.

### 5. Next.js Routing & Isolation
Used Route Groups (`(auth)` and `(dashboard)`) to enforce layout isolation, preventing unnecessary wrapper rendering on auth screens.
