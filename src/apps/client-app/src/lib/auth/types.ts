export type UserRole = 'ADMIN' | 'FARM_MANAGER' | 'WAREHOUSE_MGR' | 'STORE_MGR' | 'PROCESSOR' | 'DRIVER' | 'GUEST';

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

export interface AuthContextType extends AuthState {
  login: () => void;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
}
