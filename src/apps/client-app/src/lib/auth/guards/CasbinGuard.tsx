'use client';

import { ReactNode, useEffect, useState } from 'react';
import { useAuth } from '../AuthProvider';
import { useCasbin } from '../casbin';

interface CasbinGuardProps {
  children: ReactNode;
  obj: string;
  act: string;
  fallback?: ReactNode;
}

/**
 * CasbinGuard conditionally renders children if the user's policies allow the requested action.
 * This is the preferred way to handle fine-grained UI authorization.
 */
export function CasbinGuard({ children, obj, act, fallback = null }: CasbinGuardProps) {
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const { enforcer, isLoading: casbinLoading } = useCasbin();
  const [allowed, setAllowed] = useState<boolean>(false);
  const [checking, setChecking] = useState<boolean>(true);

  useEffect(() => {
    async function check() {
      if (authLoading || casbinLoading || !enforcer || !user) {
        return;
      }

      try {
        const isAllowed = await enforcer.enforce(user.role, obj, act);
        setAllowed(isAllowed);
      } catch (err) {
        console.error('[CasbinGuard] Enforcement error:', err);
        setAllowed(false);
      } finally {
        setChecking(false);
      }
    }

    check();
  }, [authLoading, casbinLoading, enforcer, user, obj, act]);

  if (authLoading || casbinLoading || checking) return null;

  if (!isAuthenticated || !allowed) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
