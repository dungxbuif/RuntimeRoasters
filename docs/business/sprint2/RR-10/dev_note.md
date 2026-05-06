# RR-10: Developer Notes & Post-Implementation Review

## 1. Implementation Summary
The SSO flow integration between Next.js, Kratos, and Hydra was successfully implemented. However, several critical integration edge-cases were discovered and resolved during development.

## 2. Issues Encountered & Solutions

### A. Kratos 500 Internal Server Error on `login_challenge`
- **Issue:** Calling `kratos.createBrowserLoginFlow` with a `loginChallenge` resulted in Kratos throwing a 500 error: `refusing to parse login_challenge query parameter because oauth2_provider.url is invalid or unset`.
- **Solution:** Added the `oauth2_provider` configuration to `kratos.yaml` pointing to Hydra's Admin API (`http://hydra:4445/`). This enabled Kratos to parse the challenge and natively handle the Hydra handshake.

### B. Infinite SSO Loop on Logout
- **Issue:** The custom `handleLogout` function only cleared `localStorage`. When redirecting to `/login`, Kratos still held a valid session cookie, triggering a `session_already_available` 400 error. The error boundary blindly redirected to Home, which bounced back to Hydra, creating an infinite redirect loop.
- **Solution:** Replaced standard `localStorage.clear()` with the official `kratos.createBrowserLogoutFlow()` to forcefully revoke the session cookie before returning to the login screen.

### C. Token Exchange 404 Error
- **Issue:** Using `HYDRA_ADMIN_URL` (port 4445) for token exchange in `/api/auth/callback/route.ts` returned a `404 Not Found`.
- **Solution:** The token exchange must occur on the public OAuth2 interface. Switched to using `HYDRA_PUBLIC_URL` (port 4444).

### D. Missing OAuth2 Client (`invalid_client`)
- **Issue:** Hydra reported the `client-app` didn't exist during the initial redirect.
- **Solution:** Created `client-app.json` and executed `hydra import clients` via Docker to seed the client configuration.

## 3. Future Technical Debt
- **Token Storage:** Currently storing JWTs in LocalStorage for simplicity. In a production environment, this should be migrated to HttpOnly secure cookies using Next.js Server Actions.
- **Token Refresh:** Implementation of the `refresh_token` flow is pending and should be handled by a silent background refresh interceptor in Axios.
