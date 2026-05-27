import { useState, useEffect, createContext, useContext, ReactNode } from 'react';
// @ts-expect-error casbin.js does not ship complete TypeScript declarations for this import shape.
import { Enforcer, newEnforcer } from 'casbin.js';
import { useAuth } from './AuthProvider';
import { authService } from '@/services/auth.service';
import React from 'react';

interface CasbinContextType {
  enforcer: Enforcer | null;
  can: (action: string, resource: string) => boolean;
}

const CasbinContext = createContext<CasbinContextType>({
  enforcer: null,
  can: () => false,
});

export const CasbinProvider = ({ children }: { children: ReactNode }) => {
  const { user, isAuthenticated } = useAuth();
  const [enforcer, setEnforcer] = useState<Enforcer | null>(null);

  useEffect(() => {
    const initCasbin = async () => {
      if (isAuthenticated && user?.role) {
        try {
          const { policies } = await authService.getPolicies();
          
          // Basic Casbin model definition compatible with casbin.js
          const model = `
          [request_definition]
          r = sub, obj, act

          [policy_definition]
          p = sub, obj, act

          [role_definition]
          g = _, _

          [policy_effect]
          e = some(where (p.eft == allow))

          [matchers]
          m = g(r.sub, p.sub) && (r.obj == p.obj || p.obj == '*') && (r.act == p.act || p.act == '*')
          `;
          
          const e = await newEnforcer(model, policies.join('\\n'));
          setEnforcer(e);
        } catch (error) {
          console.error('[Casbin] Failed to initialize Casbin policies:', error);
        }
      } else {
        setEnforcer(null);
      }
    };

    initCasbin();
  }, [isAuthenticated, user?.role]);

  const can = (action: string, resource: string): boolean => {
    if (!enforcer || !user) return false;
    // We enforce with the user's role as the subject
    return enforcer.enforceSync(user.role, resource, action);
  };

  return (
    <CasbinContext.Provider value={{ enforcer, can }}>
      {children}
    </CasbinContext.Provider>
  );
};

export const useCasbin = () => useContext(CasbinContext);

export const PermissionGuard = ({
  action,
  resource,
  children,
  fallback = null,
}: {
  action: string;
  resource: string;
  children: ReactNode;
  fallback?: ReactNode;
}) => {
  const { can, enforcer } = useCasbin();
  
  if (!enforcer) {
    return <>{fallback}</>;
  }

  if (can(action, resource)) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
};
