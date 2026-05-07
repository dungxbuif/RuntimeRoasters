# Identity & Auth Flow Architecture

This diagram illustrates the interactions between the Next.js Client App, Ory Kratos (Identity Provider), and Ory Hydra (OAuth2 Provider) during a standard SSO login cycle.

## Components & Endpoints
- **Client App (Next.js)**
  - Home/Dashboard: `/`
  - Login Page: `/login`
  - Callback Handler: `/api/auth/callback`
- **Ory Hydra (OIDC Provider)**
  - Auth Endpoint: `http://localhost:4444/oauth2/auth`
  - Token Endpoint: `http://localhost:4444/oauth2/token`
  - Admin API: `http://hydra:4445`
- **Ory Kratos (Identity Provider)**
  - Browser Login Flow: `http://localhost:4433/self-service/login/browser`
  - Browser Logout Flow: `http://localhost:4433/self-service/logout/browser`

## Standard Login Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant NextJS as Next.js (Client App)
    participant Hydra as Ory Hydra (OAuth2)
    participant Kratos as Ory Kratos (Identity)

    User->>NextJS: Access / (Dashboard)
    NextJS->>NextJS: Check LocalStorage (No Token)
    NextJS->>Hydra: Redirect to /oauth2/auth?client_id=...
    Hydra->>NextJS: Redirect to /login?login_challenge=abc
    NextJS->>Kratos: GET /self-service/login/browser (with login_challenge)
    Kratos-->>NextJS: Return Flow ID & Form Nodes
    NextJS->>User: Render Login Form
    User->>NextJS: Submit Credentials (Email/Password)
    NextJS->>Kratos: POST /self-service/login?flow=...
    Kratos->>Kratos: Validate Credentials & Issue Session Cookie
    Kratos->>Hydra: Admin API: Accept Login Request
    Hydra-->>Kratos: Return redirect_to (Consent/Callback)
    Kratos-->>NextJS: 422 browser_location_change_required (redirect_browser_to)
    NextJS->>Hydra: Redirect to redirect_browser_to
    Hydra->>NextJS: Redirect to /api/auth/callback?code=xyz
    NextJS->>Hydra: POST /oauth2/token (Exchange code for JWT)
    Hydra-->>NextJS: Return Access Token & ID Token
    NextJS->>NextJS: Save JWT to LocalStorage
    NextJS->>User: Redirect to / (Dashboard)
```
