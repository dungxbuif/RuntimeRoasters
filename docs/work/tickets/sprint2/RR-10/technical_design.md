# Dev Notes - [RR-10] Client-Side Auth & Login Flow

## 🛠️ Technical Implementation Details
- **Identity Flow:** Ory Kratos login → session token → Token Exchange to JWT (native Kratos OIDC Token Exchange).
- **Token Storage:** HttpOnly Cookie (prevents XSS token theft).
- **Gateway Role:** KrakenD acts as Traffic Coordinator — checks App **Scopes** (Gate 1) via Native Validator JSON config. No LUA scripts or Go Plugins.
- **State Management:** Zustand or simple Context to store User profile (decoded from JWT: `sub`, `role`, `org_id`).
- **Route Protection:** Next.js Middleware at edge; check cookie presence.
- **KrakenD Config:** `scopes` claim validation is declarative JSON — no custom auth plugin.

## 🏗️ Architectural Approach
The Frontend (`client-app`) acts as the Login Provider for Ory Hydra, utilizing Ory Kratos to handle the actual identity verification. This establishes a "Zero-Consent" OIDC SSO flow.

### Flow Breakdown:
1. **Unauthenticated Access:** User attempts to access a protected route (e.g., Dashboard).
2. **Redirect to Hydra:** App redirects to Hydra's `/oauth2/auth` endpoint.
3. **Hydra Challenge:** Hydra redirects back to the App's Login UI with a `login_challenge`.
4. **Kratos Authentication:** App uses `@ory/client` to render the Kratos login form and submit credentials.
5. **Accept Challenge:** Upon successful Kratos authentication, the App fetches the user identity and accepts the Hydra challenge.
6. **Token Exchange:** Hydra redirects to the App's callback, where the App exchanges the authorization code for an Access Token (JWT).

## 📂 Directory Structure
```text
src/app/
├── (auth)/
│   ├── login/page.tsx       # Kratos Login UI & Flow Management
│   └── consent/page.tsx     # Auto-consent handler
├── (dashboard)/
│   ├── layout.tsx           # Dashboard layout (Sidebar, Header)
│   └── page.tsx             # Main dashboard (Protected)
├── api/auth/
│   ├── callback/route.ts    # Token exchange endpoint
│   ├── login/accept/route.ts# Accept Hydra Login Server Action
│   └── consent/accept/route.ts
```

## 🎨 UI/UX Specifications
- **LoginCard:** Centered glassmorphic card without harsh borders.
- **KratosForm:** Dynamically parses `ui.nodes` from Kratos.
- **Loading State:** Themed spinner with "Initializing Terminal UI..." matching the primary brand color (`primary/20`).

## 🛡️ Security Measures
- Tokens are temporarily stored in LocalStorage for demo purposes (can be upgraded to HttpOnly cookies later).
- Cross-Site Request Forgery (CSRF) protection provided natively by Kratos.

## 📝 Implementation Review & Troubleshooting

### Issues Encountered & Solutions
- **Kratos 500 Error:** Resolved by adding `oauth2_provider` config to `kratos.yaml` pointing to Hydra's Admin API.
- **Logout Loop:** Fixed by using `kratos.createBrowserLogoutFlow()` instead of just clearing LocalStorage.
- **Token Exchange 404:** Switched from Hydra Admin URL (4445) to Public URL (4444) for token exchange.
- **Invalid Client:** Ensured client-app was seeded in Hydra using `hydra import clients`.

## ⚠️ Technical Debt
- **Token Storage:** Migrate JWTs from LocalStorage to HttpOnly secure cookies.
- **Token Refresh:** Implement `refresh_token` flow with a background refresh interceptor.
