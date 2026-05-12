'use client';

import { useAuth } from '../AuthProvider';
import { UserRole } from '../types';
import { ReactNode } from 'react';

interface RoleGuardProps {
  children: ReactNode;
  roles: UserRole[];
  fallback?: ReactNode;
}

/**
 * RoleGuard conditionally renders children if the user has one of the allowed roles.
 */
export function RoleGuard({ children, roles, fallback = null }: RoleGuardProps) {
  const { user, isAuthenticated, isLoading } = useAuth();

  if (isLoading) return null;

  if (!isAuthenticated || !user || !roles.includes(user.role)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
