# Frontend Coding Guidelines — Runtime Roasters

This document aggregates the frontend development rules and patterns based on project feedback and requirements. All future code changes must strictly adhere to these rules.

## 1. Service Architecture (OOP Class)
All service logic (Auth, Storage, API) must be implemented as **OOP Classes**.
- Do not use discrete exported functions for complex logic.
- Export a single instance of the class (Singleton pattern).

**Example:**
```typescript
class AuthService {
  async login(...) { ... }
}
export const authService = new AuthService();
```

## 2. Constant Management (No Magic Strings)
The use of magic strings in code is strictly prohibited. All static values must be defined in the `src/constants/` directory.

- **ENV**: Environment variables (`constants/env.ts`).
- **API_ENDPOINTS**: All API URLs (`constants/api.ts`).
- **STORAGE_KEYS**: LocalStorage/SessionStorage keys (`constants/storage.ts`).
- **APP_ROUTES**: Page routes within the application (`constants/routes.ts`).
- **AUTH_PARAMS**: Query parameters related to Auth (`constants/auth.ts`).

## 3. Libraries & State Management
- **API Client**: Always use `Axios`. Avoid using `fetch` directly except in extremely special cases.
- **Data Fetching**: Use `TanStack Query` (React Query) for all server interactions to manage cache and loading states.
- **Styling**: Vanilla CSS or TailwindCSS v4 (as per requirements).

## 4. Utilities & Helpers
- **Environment**: Use the `isBrowser` utility instead of `typeof window !== 'undefined'`.
- **URL Handling**: Use the `joinPaths` function to concatenate URL components, avoiding errors from extra or missing slashes (`/`).
- **Type Safety**: Avoid using `any`. Always define clear interfaces/types.

## 5. Aesthetics (Industrial Premium)
- The interface must embody an **Industrial Premium** style:
    - Use Glassmorphism (backdrop-blur).
    - Harmonious colors, high contrast but elegant (slate, emerald, primary).
    - Subtle animation effects with `framer-motion`.
    - Modern typography (Google Fonts).

## 6. Lints & Types
- Code must pass `npm run lint` and `npx tsc --noEmit` checks.
- Do not suppress lint errors with comments unless absolutely necessary.

## 7. Auth Patterns (Seamless Identity)
The project uses the standard OIDC model but optimizes the user experience (Zero-Consent flow for internal apps).

- **Auto-Accept Consent**: The `/consent` page is implemented as a **Server Component**.
    - It automatically checks the `client_id` against a `TRUSTED_CLIENTS` list.
    - It executes the `acceptOAuth2ConsentRequest` command directly from the server-side.
    - Users never see the authorization screen, providing a smooth SSO-like feel.
- **Hydra Admin Access**: Always perform Admin tasks (Accept Login/Consent) from the server-side via API routes or Server Components to secure tokens.

---
> [!TIP]
> Use the **Master Checklist** in `docs/master_checklist.md` to verify system-wide security and integration points after making changes.

> [!IMPORTANT]
> When reviewing code or performing new tasks, Antigravity must automatically check if the code violates the above rules (especially magic strings and OOP services).
