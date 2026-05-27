'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode, useCallback } from 'react';
import { AuthContextType, AuthState, AuthUser, UserRole } from './types';
import { storageService } from '@/services/storage.service';
import { ENV } from '@/constants/env';
import { OIDC_CONFIG } from '@/constants/auth';
import { useRouter } from 'next/navigation';
import { APP_ROUTES } from '@/constants/routes';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
    token: null,
  });
  const router = useRouter();

  const initialize = useCallback(async () => {
    try {
      const token = storageService.getAccessToken();
      if (!token) {
        setState({ user: null, isAuthenticated: false, isLoading: false, token: null });
        return;
      }

      try {
        const { authService } = await import('@/services/auth.service');
        try {
          const { user: apiUser } = await authService.getMe();
          if (apiUser) {
            const user: AuthUser = {
              id: apiUser.id,
              email: apiUser.email,
              role: (apiUser.role as UserRole) || 'GUEST',
              name: apiUser.name,
            };
            setState({
              user,
              isAuthenticated: true,
              isLoading: false,
              token,
            });
            return;
          }
        } catch (meError) {
          console.warn('[Auth] getMe failed, trying session fallback', meError);
          // Fallback to Kratos session
          const session = await authService.getSession();
          if (session.identity) {
            const traits = session.identity.traits as { email: string; name?: string; role?: string };
            const user: AuthUser = {
              id: session.identity.id,
              email: traits.email,
              role: (traits.role as UserRole) || 'GUEST',
              name: traits.name,
            };
            
            setState({
              user,
              isAuthenticated: true,
              isLoading: false,
              token,
            });
            return;
          }
        }
      } catch (e) {
        console.warn('[Auth] Session validation failed', e);
        storageService.clearAll();
      }

      setState({ user: null, isAuthenticated: false, isLoading: false, token: null });
    } catch (error) {
      console.error('[Auth] Initialization failed', error);
      setState(prev => ({ ...prev, isLoading: false }));
    }
  }, []);

  useEffect(() => {
    let isMounted = true;
    const load = async () => {
      if (isMounted) {
        await initialize();
      }
    };
    load();
    return () => { isMounted = false; };
  }, [initialize]);

  const login = useCallback(() => {
    const authUrl = new URL(`${ENV.HYDRA_PUBLIC_URL}/oauth2/auth`);
    authUrl.searchParams.append('client_id', OIDC_CONFIG.CLIENT_ID);
    authUrl.searchParams.append('response_type', 'code');
    authUrl.searchParams.append('scope', 'openid offline_access');
    authUrl.searchParams.append('redirect_uri', `${ENV.APP_URL}/api/auth/callback`);
    authUrl.searchParams.append('state', crypto.randomUUID());
    window.location.href = authUrl.toString();
  }, []);

  const logout = useCallback(async () => {
    storageService.clearAll();
    setState({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      token: null,
    });
    
    try {
      const { default: kratos } = await import('@/lib/ory/kratos');
      const { data } = await kratos.createBrowserLogoutFlow();
      window.location.href = data.logout_url;
    } catch {
      router.push(APP_ROUTES.LOGIN);
    }
  }, [router]);

  const refreshSession = useCallback(async () => {
    await initialize();
  }, [initialize]);

  const value = React.useMemo(() => ({ 
    ...state, 
    login, 
    logout, 
    refreshSession 
  }), [state, login, logout, refreshSession]);

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
