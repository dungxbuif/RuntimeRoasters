'use client';

import { useAuth } from '../AuthProvider';
import { ReactNode, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { APP_ROUTES } from '@/constants/routes';

interface AuthGuardProps {
  children: ReactNode;
  fallback?: ReactNode;
  requireAuth?: boolean;
}

export function AuthGuard({ children, fallback = null, requireAuth = true }: AuthGuardProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated && requireAuth && !fallback) {
      router.push(APP_ROUTES.LOGIN);
    }
  }, [isLoading, isAuthenticated, requireAuth, fallback, router]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-surface" data-e2e="auth-loader">
        <div className="text-[10px] font-black uppercase tracking-[0.4em] text-primary animate-pulse italic">
          Verifying Identity Cluster...
        </div>
      </div>
    );
  }

  if (!isAuthenticated && requireAuth) {
    return fallback ? <>{fallback}</> : null;
  }

  return <>{children}</>;
}
