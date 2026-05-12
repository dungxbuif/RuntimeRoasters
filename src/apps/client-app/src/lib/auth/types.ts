export type UserRole = 'ADMIN' | 'FARM_ADMIN' | 'FARM_MANAGER' | 'FARMER' | 'PROCESSOR' | 'DRIVER' | 'GUEST';

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
