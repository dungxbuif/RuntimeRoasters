# FE Auth Library — Design & Integration Blueprint

## 1. Overview
The **RuntimeRoasters Auth Library** is a centralized, plug-and-play authentication package for the Frontend. It eliminates manual `if-else` checks and token management logic within business components. It leverages **Ory Kratos/Hydra** and provides a declarative API for protecting routes and UI fragments.

## 2. Core Components

### `AuthContext` & `AuthProvider`
- **Role**: Maintains the global authentication state (`isAuthenticated`, `user`, `token`, `isLoading`).
- **Logic**: Automatically initializes the session from `storageService` or Kratos on mount.
- **Location**: `src/lib/auth/AuthProvider.tsx`

### `useAuth()` Hook
- **Role**: The primary way to access auth state.
- **API**: `const { user, isAuthenticated, login, logout } = useAuth();`

### `<AuthGuard />` (Component)
- **Role**: Declarative wrapper to protect parts of a page.
- **Behavior**: If unauthenticated, it can either hide children, show a fallback, or redirect to `/login`.
- **Usage**:
  ```tsx
  <AuthGuard fallback={<PublicTeaser />}>
    <SensitiveDashboard />
  </AuthGuard>
  ```

### `<RoleGuard />` (Component)
- **Role**: RBAC (Role-Based Access Control) wrapper.
- **Usage**:
  ```tsx
  <RoleGuard roles={['ADMIN', 'FARM_ADMIN']}>
    <DeleteUserButton />
  </RoleGuard>
  ```

---

## 3. Library Structure (Implementation)

### File: `src/lib/auth/types.ts`
```typescript
export type UserRole = 'ADMIN' | 'FARM_ADMIN' | 'FARM_MANAGER' | 'FARMER' | 'GUEST';

export interface AuthUser {
  id: string;
  email: string;
  role: UserRole;
  name?: string;
}

export interface AuthState {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  token: string | null;
}
```

### File: `src/lib/auth/AuthProvider.tsx`
(Implementation follows standard React Context pattern with automated token validation).

---

## 4. Integration Plan
1.  **Wrap Root**: Add `AuthProvider` to `src/components/common/Providers.tsx`.
2.  **Refactor Dashboard**: Use `AuthGuard` in `(dashboard)/page.tsx` instead of manual `useEffect`.
3.  **Refactor Sidebar**: Use `RoleGuard` to show/hide management links dynamically.
