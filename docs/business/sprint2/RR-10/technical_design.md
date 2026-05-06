# RR-10: Technical Design & Implementation Plan
**Feature:** Client-Side Auth & Login Flow

## 1. Architectural Approach
The Frontend (`client-app`) will act as the Login Provider for Ory Hydra, utilizing Ory Kratos to handle the actual identity verification. This establishes a "Zero-Consent" OIDC SSO flow.

### Flow Breakdown:
1. **Unauthenticated Access:** User attempts to access a protected route (e.g., Dashboard).
2. **Redirect to Hydra:** App redirects to Hydra's `/oauth2/auth` endpoint.
3. **Hydra Challenge:** Hydra redirects back to the App's Login UI with a `login_challenge`.
4. **Kratos Authentication:** App uses `@ory/client` to render the Kratos login form and submit credentials.
5. **Accept Challenge:** Upon successful Kratos authentication, the App fetches the user identity and accepts the Hydra challenge.
6. **Token Exchange:** Hydra redirects to the App's callback, where the App exchanges the authorization code for an Access Token (JWT).

## 2. Frontend Infrastructure
- **Framework:** Next.js 15 (App Router).
- **Libraries:** `@ory/client`, `axios`, `framer-motion`, `lucide-react`.
- **Styling:** Tailwind CSS 4 with a custom "Industrial Premium" aesthetic (Dark/Light hybrid, gradients, glassmorphism).

## 3. Directory Structure
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

## 4. UI/UX Specifications
- **LoginCard:** Centered glassmorphic card without harsh borders.
- **KratosForm:** Dynamically parses `ui.nodes` from Kratos.
- **Loading State:** Themed spinner with "Initializing Terminal UI..." matching the primary brand color (`primary/20`).

## 5. Security Measures
- Tokens are temporarily stored in LocalStorage for demo purposes (can be upgraded to HttpOnly cookies later).
- Cross-Site Request Forgery (CSRF) protection provided natively by Kratos.
